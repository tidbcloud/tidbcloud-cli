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

package util

import (
	"context"
	"errors"
	"io"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/require"
)

func TestBuildTiDBShellDSNEscapesCredentials(t *testing.T) {
	password := "pa#ss?/word%@"
	dsn := buildTiDBShellDSN("user@tenant", "2001:db8::1", "4000", &password, "tidb-nextgen")

	parsed, err := url.Parse(dsn)
	require.NoError(t, err)
	require.Equal(t, "user@tenant", parsed.User.Username())
	actualPassword, ok := parsed.User.Password()
	require.True(t, ok)
	require.Equal(t, password, actualPassword)
	require.Equal(t, "[2001:db8::1]:4000", parsed.Host)
	require.Equal(t, "tidb-nextgen", parsed.Query().Get("tls"))
}

type fakeSQLDialog struct {
	open func(context.Context, ...string) error
	run  func() error
}

func (h fakeSQLDialog) Open(ctx context.Context, params ...string) error {
	return h.open(ctx, params...)
}
func (h fakeSQLDialog) Run() error { return h.run() }

func TestSQLDialogConnectionTimeout(t *testing.T) {
	h := fakeSQLDialog{
		open: func(ctx context.Context, params ...string) error {
			require.Equal(t, []string{"dsn"}, params)
			_, hasDeadline := ctx.Deadline()
			require.True(t, hasDeadline)
			<-ctx.Done()
			return ctx.Err()
		},
		run: func() error { t.Fatal("must not start SQL shell after connection failure"); return nil },
	}
	err := runSQLDialog(context.Background(), h, "dsn", "private.example.com:4000", 10*time.Millisecond)
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.ErrorContains(t, err, "private.example.com:4000")
}

func TestSQLDialogTimeoutDoesNotLimitSession(t *testing.T) {
	parent := context.Background()
	var connectCtx context.Context
	h := fakeSQLDialog{
		open: func(ctx context.Context, _ ...string) error {
			connectCtx = ctx
			return nil
		},
		run: func() error {
			require.ErrorIs(t, connectCtx.Err(), context.Canceled)
			require.NoError(t, parent.Err())
			return io.EOF
		},
	}
	require.NoError(t, runSQLDialog(parent, h, "dsn", "private.example.com:4000", time.Second))
}

func TestSQLDialogRejectsTimeoutSwallowedByOpen(t *testing.T) {
	h := fakeSQLDialog{
		open: func(ctx context.Context, _ ...string) error {
			<-ctx.Done()
			// usql's version query reports an unknown version but returns nil.
			return nil
		},
		run: func() error { t.Fatal("must not enter SQL shell after connection timeout"); return nil },
	}
	err := runSQLDialog(context.Background(), h, "dsn", "private.example.com:4000", 10*time.Millisecond)
	require.ErrorIs(t, err, context.DeadlineExceeded)
}

func TestSQLDialogWithoutTimeoutPreservesContextAndErrors(t *testing.T) {
	parent, cancel := context.WithCancel(context.Background())
	defer cancel()
	queryErr := errors.New("query failed")
	h := fakeSQLDialog{
		open: func(ctx context.Context, _ ...string) error {
			require.Same(t, parent, ctx)
			_, hasDeadline := ctx.Deadline()
			require.False(t, hasDeadline)
			return nil
		},
		run: func() error {
			require.NoError(t, parent.Err())
			return queryErr
		},
	}
	require.ErrorIs(t, runSQLDialog(parent, h, "dsn", "public.example.com:4000", 0), queryErr)
}
