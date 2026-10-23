// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package nextgen

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"testing/iotest"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/telemetry"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/stretchr/testify/require"
)

func TestPasswordCommandResetsPremiumRootPassword(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_PREMIUM)
	id := "tidb-1"
	instance.TidbId = &id
	var captured *api.TidbServiceResetRootPasswordBody
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(_ context.Context, actualID string, body *api.TidbServiceResetRootPasswordBody) error {
			require.Equal(t, id, actualID)
			captured = body
			return nil
		},
	}
	command := passwordCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetIn(strings.NewReader("new-password"))
	command.SetArgs([]string{"--cluster-id", id, "--password-stdin"})

	require.NoError(t, command.ExecuteContext(context.Background()))
	require.NotNil(t, captured)
	require.Equal(t, "new-password", captured.GetRootPassword())
}

func TestPasswordCommandSupportsInteractivePasswordWithExplicitInstance(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_ESSENTIAL_V2)
	id := "tidb-1"
	instance.TidbId = &id
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(_ context.Context, actualID string, body *api.TidbServiceResetRootPasswordBody) error {
			require.Equal(t, id, actualID)
			require.Equal(t, "new-password", body.GetRootPassword())
			return nil
		},
	}
	h := helperWithNextGenClient(client)
	h.IOStreams.CanPrompt = true
	command := passwordCmdWithDependencies(
		h,
		essentialV2Plan,
		func(context.Context, cloud.NextGenClient, planSpec, int32) (*api.Nextgenv1beta2Tidb, error) {
			t.Fatal("instance selector must not be called when --cluster-id is set")
			return nil, nil
		},
		func() (string, error) { return "new-password", nil },
	)
	command.SetArgs([]string{"--cluster-id", id})

	require.NoError(t, command.ExecuteContext(context.Background()))
	require.Equal(t, "true", command.Annotations[telemetry.InteractiveMode])
}

func TestPasswordCommandPreservesLineBreaks(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", plan.servicePlan)
			password := "1234\r\n5678"
			calls := 0
			client := &fakeNextGenClient{
				get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
				password: func(_ context.Context, id string, body *api.TidbServiceResetRootPasswordBody) error {
					calls++
					require.Equal(t, "tidb-1", id)
					require.Equal(t, password, body.GetRootPassword())
					return nil
				},
			}
			command := passwordCmd(helperWithNextGenClient(client), plan)
			command.SetIn(strings.NewReader(password))
			command.SetArgs([]string{"-c", "tidb-1", "--password-stdin"})
			require.NoError(t, command.ExecuteContext(context.Background()))
			require.Equal(t, 1, calls)
		})
	}
}

func TestPasswordCommandRejectsInvalidPasswordBeforeRequest(t *testing.T) {
	clientCalls := 0
	h := helperWithNextGenClient(&fakeNextGenClient{})
	h.NextGenClient = func() (cloud.NextGenClient, error) {
		clientCalls++
		return &fakeNextGenClient{}, nil
	}
	command := passwordCmd(h, premiumPlan)
	command.SetIn(strings.NewReader("short"))
	command.SetArgs([]string{"--cluster-id", "tidb-1", "--password-stdin"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, "root password must be between 8 and 64 bytes")
	require.Zero(t, clientCalls)
}

func TestPasswordCommandRequiresInstanceForNonInteractiveUse(t *testing.T) {
	tests := []struct {
		name string
		args []string
	}{
		{name: "omitted", args: []string{"--password-stdin"}},
		{name: "empty", args: []string{"--cluster-id=", "--password-stdin"}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			clientCalls := 0
			h := helperWithNextGenClient(&fakeNextGenClient{})
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				clientCalls++
				return &fakeNextGenClient{}, nil
			}
			command := passwordCmd(h, premiumPlan)
			command.SetArgs(tt.args)

			err := command.ExecuteContext(context.Background())
			require.ErrorContains(t, err, `--cluster-id is required when --password-stdin is specified`)
			require.Zero(t, clientCalls)
		})
	}
}

func TestPasswordCommandRejectsCrossPlanInstance(t *testing.T) {
	instance := api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", api.V1BETA1SERVICEPLAN_ESSENTIAL_V2)
	id := "tidb-1"
	instance.TidbId = &id
	resetCalls := 0
	client := &fakeNextGenClient{
		get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) { return instance, nil },
		password: func(context.Context, string, *api.TidbServiceResetRootPasswordBody) error {
			resetCalls++
			return nil
		},
	}
	command := passwordCmd(helperWithNextGenClient(client), premiumPlan)
	command.SetIn(strings.NewReader("new-password"))
	command.SetArgs([]string{"--cluster-id", id, "--password-stdin"})

	err := command.ExecuteContext(context.Background())
	require.ErrorContains(t, err, `only manages "Premium" instances`)
	require.Zero(t, resetCalls)
}

func TestValidateRootPasswordAcceptsEncryptedPayload(t *testing.T) {
	require.NoError(t, validateRootPassword("rsa_oaep_sha256:ciphertext"))
	require.Error(t, validateRootPassword("rsa_oaep_sha256:"))
}

func TestValidateRootPasswordMatchesBackendByteLimit(t *testing.T) {
	for _, tc := range []struct {
		name     string
		password string
		valid    bool
	}{
		{name: "below minimum", password: "1234567"},
		{name: "minimum", password: "12345678", valid: true},
		{name: "maximum", password: strings.Repeat("a", 64), valid: true},
		{name: "above maximum", password: strings.Repeat("a", 65)},
		{name: "multibyte below eight runes", password: "测试密码", valid: true},
		{name: "multibyte maximum", password: strings.Repeat("密", 21) + "a", valid: true},
		{name: "multibyte above maximum", password: strings.Repeat("密", 30)},
		{name: "line feed", password: "1234\n5678", valid: true},
		{name: "carriage return", password: "1234\r5678", valid: true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			err := validateRootPassword(tc.password)
			if tc.valid {
				require.NoError(t, err)
			} else {
				require.Error(t, err)
			}
		})
	}
}

func TestPasswordStdinInput(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name, input, want, wantError string
		}{
			{name: "no newline", input: "new-password", want: "new-password"},
			{name: "LF", input: "new-password\n", want: "new-password"},
			{name: "CRLF", input: "new-password\r\n", want: "new-password"},
			{name: "spaces", input: " padded-secret \n", want: " padded-secret "},
			{name: "one newline only", input: "new-password\n\n", want: "new-password\n"},
			{name: "bare CR preserved", input: "new-password\r", want: "new-password\r"},
			{name: "encrypted", input: "rsa_oaep_sha256:ciphertext\n", want: "rsa_oaep_sha256:ciphertext"},
			{name: "empty", wantError: "between 8 and 64 bytes"},
			{name: "newline only", input: "\n", wantError: "between 8 and 64 bytes"},
			{name: "too short", input: "short", wantError: "between 8 and 64 bytes"},
			{name: "too long", input: strings.Repeat("x", 65), wantError: "between 8 and 64 bytes"},
			{name: "oversized input", input: strings.Repeat("x", (1<<20)+1), wantError: "exceeds 1 MiB"},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				calls := 0
				client := &fakeNextGenClient{
					get: func(context.Context, string) (*api.Nextgenv1beta2Tidb, error) {
						require.Empty(t, tc.wantError, "invalid input must fail before an API call")
						return api.NewNextgenv1beta2Tidb("test-instance", "aws-us-west-2", plan.servicePlan), nil
					},
					password: func(_ context.Context, id string, body *api.TidbServiceResetRootPasswordBody) error {
						calls++
						require.Equal(t, "tidb-1", id)
						require.Equal(t, tc.want, body.GetRootPassword())
						return nil
					},
				}
				h := helperWithNextGenClient(client)
				h.IOStreams.CanPrompt = true // stdin mode must not prompt even in a terminal.
				command := passwordCmdWithDependencies(h, plan, nil, func() (string, error) {
					t.Fatal("stdin mode must not prompt")
					return "", nil
				})
				command.SetIn(strings.NewReader(tc.input))
				command.SetArgs([]string{"-c", "tidb-1", "--password-stdin"})
				err := command.Execute()
				if tc.wantError != "" {
					require.ErrorContains(t, err, tc.wantError)
					require.Zero(t, calls)
				} else {
					require.NoError(t, err)
					require.Equal(t, 1, calls)
					require.NotContains(t, command.Annotations, telemetry.InteractiveMode)
					require.NotContains(t, fmt.Sprint(h.IOStreams.Out), tc.want)
				}
			})
		}
	}
}

func TestPasswordStdinRejectsReadErrorAndOldFlag(t *testing.T) {
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		for _, tc := range []struct {
			name      string
			args      []string
			reader    io.Reader
			wantError string
		}{
			{"read error", []string{"-c", "tidb-1", "--password-stdin"}, iotest.ErrReader(errors.New("input unavailable")), "read password from standard input"},
			{"old flag", []string{"-c", "tidb-1", "--password", "secret-value"}, strings.NewReader(""), "unknown flag: --password"},
			{"no terminal", []string{"-c", "tidb-1"}, strings.NewReader(""), "use --password-stdin with --cluster-id"},
			{"stdin false", []string{"-c", "tidb-1", "--password-stdin=false"}, strings.NewReader(""), "use --password-stdin with --cluster-id"},
		} {
			t.Run(plan.commandName+"/"+tc.name, func(t *testing.T) {
				command := passwordCmd(helperWithNextGenClient(nil), plan)
				command.SetIn(tc.reader)
				command.SetArgs(tc.args)
				err := command.Execute()
				require.ErrorContains(t, err, tc.wantError)
				require.NotContains(t, err.Error(), "secret-value")
			})
		}
	}
}

// Expose a real terminal descriptor, but fail immediately if the code attempts
// to read it. This keeps the regression test from hanging with the old behavior.
type terminalPasswordInput struct {
	fd uintptr
}

func (r terminalPasswordInput) Fd() uintptr { return r.fd }
func (r terminalPasswordInput) Read([]byte) (int, error) {
	return 0, errors.New("unexpected read from terminal")
}

func TestPasswordStdinRejectsTerminalBeforeRead(t *testing.T) {
	terminal, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		t.Skipf("controlling terminal unavailable: %v", err)
	}
	defer terminal.Close()
	for _, plan := range []planSpec{premiumPlan, essentialV2Plan} {
		t.Run(plan.commandName, func(t *testing.T) {
			h := helperWithNextGenClient(nil)
			h.IOStreams.CanPrompt = true
			h.NextGenClient = func() (cloud.NextGenClient, error) {
				t.Fatal("terminal input must be rejected before creating a client")
				return nil, nil
			}
			command := passwordCmd(h, plan)
			command.SetIn(terminalPasswordInput{fd: terminal.Fd()})
			command.SetArgs([]string{"-c", "tidb-1", "--password-stdin"})
			err := command.Execute()
			require.ErrorContains(t, err, "--password-stdin requires a pipe or file")
			require.ErrorContains(t, err, "hidden input")
			require.NotContains(t, err.Error(), "unexpected read")
		})
	}
}

func TestReadRootPasswordAcceptsFileAndPipe(t *testing.T) {
	path := filepath.Join(t.TempDir(), "password")
	require.NoError(t, os.WriteFile(path, []byte("synthetic-password\r\n"), 0600))
	file, err := os.Open(path)
	require.NoError(t, err)
	defer file.Close()
	password, err := readRootPassword(file)
	require.NoError(t, err)
	require.Equal(t, "synthetic-password", password)

	reader, writer, err := os.Pipe()
	require.NoError(t, err)
	defer reader.Close()
	defer writer.Close()
	_, err = io.WriteString(writer, "synthetic-password\n")
	require.NoError(t, err)
	require.NoError(t, writer.Close())
	password, err = readRootPassword(reader)
	require.NoError(t, err)
	require.Equal(t, "synthetic-password", password)
}
