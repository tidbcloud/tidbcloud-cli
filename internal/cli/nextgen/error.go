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
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/redact"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
)

type rcuError struct {
	error
	hint string
}

func (e *rcuError) Error() string { return e.error.Error() + "; " + e.hint }
func (e *rcuError) Unwrap() error { return e.error }

func explainRCUError(err error, maxRCU int64, plan planSpec) error {
	var apiError *cloud.NextGenAPIError
	if !errors.As(err, &apiError) || apiError.Status != http.StatusBadRequest {
		return err
	}
	message, _ := apiError.Response["message"].(string)
	message = strings.TrimSuffix(strings.TrimSpace(message), ".")
	if !strings.EqualFold(message, "Min RCU must be less than or equal to max RCU") {
		return err
	}
	return &rcuError{error: err, hint: fmt.Sprintf(
		"--max-rcu=%d is below the service-calculated minimum. The default MAX_RCU minimum for %s is %d RCU; organization-specific limits may differ. The server did not return the actual minimum in its message. Confirm the applicable limit before retrying. MAX_RCU creation requires --max-rcu with no usable default; Elastic mode uses --baseline-rcu instead. Updates only change capacity when explicitly supplied",
		maxRCU, plan.displayName, plan.minimumMaxRCU)}
}

// JSONError is limited to plan commands that support --output. Re-reading flags
// with ParseAll handles -ojson and validation errors without mutating flag values.
// Shell and legacy commands retain their existing error rendering.
func JSONError(command *cobra.Command, args []string, err error) (string, bool) {
	if err == nil || !UsesJSONOutput(command, args) {
		return "", false
	}
	return marshalJSONError(err)
}

// UsesJSONOutput also lets credential failures avoid a human preamble on stdout.
func UsesJSONOutput(command *cobra.Command, args []string) bool {
	if command == nil || (command.Flags().Lookup(flag.Output) == nil && !isPlanCommand(command.Name())) {
		return false
	}
	plan := command
	for plan.Parent() != nil && !isPlanCommand(plan.Name()) {
		plan = plan.Parent()
	}
	if !isPlanCommand(plan.Name()) {
		return false
	}
	flags := pflag.NewFlagSet("error-output", pflag.ContinueOnError)
	flags.SetOutput(io.Discard)
	flags.ParseErrorsWhitelist.UnknownFlags = true
	flags.AddFlagSet(command.Flags())
	flags.AddFlagSet(command.InheritedFlags())
	// Action normalization can fail before a subcommand is selected.
	if isPlanCommand(command.Name()) {
		for _, child := range command.Commands() {
			if child.Hidden {
				flags.AddFlagSet(child.Flags())
			}
		}
	}
	format := ""
	// Only the new private endpoint commands opt into default JSON errors.
	// Keep the explicit-output behavior of existing commands unchanged.
	if parent := command.Parent(); parent != nil && parent.Name() == "private-endpoint" {
		if outputFlag := command.Flags().Lookup(flag.Output); outputFlag != nil {
			format = outputFlag.DefValue
		}
	}
	if parseErr := flags.ParseAll(args, func(f *pflag.Flag, value string) error {
		if f.Name == flag.Output {
			format = value
		}
		return nil
	}); parseErr != nil {
		return false
	}
	return format == output.JsonFormat
}

func marshalJSONError(err error) (string, bool) {
	message := err.Error()
	var transportError *url.Error
	if errors.As(err, &transportError) {
		safeURL := "<URL omitted>"
		if parsed, parseErr := url.Parse(transportError.URL); parseErr == nil {
			safeURL = redact.URL(parsed)
		}
		message = strings.ReplaceAll(message, transportError.URL, safeURL)
	}
	result := map[string]interface{}{"message": message}
	var apiError *cloud.NextGenAPIError
	if errors.As(err, &apiError) {
		result = map[string]interface{}{
			"message": http.StatusText(apiError.Status), "status": apiError.Status,
			"requestId": apiError.RequestID, "method": apiError.Method, "path": apiError.Path,
		}
		if apiError.Response != nil {
			result["response"] = apiError.Response
			if message, ok := apiError.Response["message"].(string); ok {
				result["message"] = message
			}
		}
		// HTTP status and response messages must not hide a local read failure or
		// rejected redirect. The reason is allowlisted by the Cloud error parser.
		if reason := apiError.LocalFailureReason(); reason != "" {
			result["message"] = reason
		}
	}
	var rcu *rcuError
	if errors.As(err, &rcu) {
		result["hint"] = rcu.hint
	}
	data, marshalErr := json.MarshalIndent(map[string]interface{}{"error": result}, "", "  ")
	if marshalErr != nil {
		return "", false
	}
	return string(data), true
}
