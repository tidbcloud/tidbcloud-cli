// Copyright 2024 PingCAP, Inc.
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

package auth

import (
	"fmt"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	ver "tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
	"github.com/zalando/go-keyring"
)

type LogoutOpts struct {
	client *resty.Client
}

func LogoutCmd(h *internal.Helper) *cobra.Command {
	opts := LogoutOpts{
		client: resty.New(),
	}
	var logoutCmd = &cobra.Command{
		Use:   "logout",
		Short: "Log out of TiDB Cloud",
		Example: fmt.Sprintf(`  To log out of TiDB Cloud:
  $ %[1]s auth logout`, config.CliName),
		PreRunE: func(cmd *cobra.Command, args []string) error {
			debug, err := cmd.Flags().GetBool(flag.Debug)
			if err != nil {
				return err
			}
			opts.client.SetDebug(debug)
			return nil
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			token, err := config.GetAccessToken()
			if err != nil {
				if errors.Is(err, keyring.ErrNotFound) {
					fmt.Fprintln(h.IOStreams.Out, color.YellowString("You are not logged in."))
					return nil
				}
				return err
			}

			body := RevokeRequest{
				Token:         token,
				TokenTypeHint: tokenTypeHint,
				ClientID:      config.GetOAuthClientID(),
				ClientSecret:  config.GetOAuthClientSecret(),
			}

			resp, err := opts.client.R().
				SetContext(ctx).
				SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
				SetHeader("Content-type", "application/json").
				SetBody(body).
				Post(fmt.Sprintf("%s%s", config.GetOAuthEndpoint(), revokePath))
			if err != nil {
				return err
			}

			if !resp.IsSuccess() {
				return errors.Errorf("Failed to revoke access token, code: %s, body: %s", resp.Status(), string(resp.Body()))
			}

			err = config.DeleteAccessToken()
			if err != nil {
				return err
			}

			fmt.Fprintln(h.IOStreams.Out, color.GreenString("You have successfully logged out."))

			return nil
		},
	}

	return logoutCmd
}

type RevokeRequest struct {
	Token         string `json:"token"`
	TokenTypeHint string `json:"token_type_hint"`
	ClientID      string `json:"client_id"`
	ClientSecret  string `json:"client_secret"`
}
