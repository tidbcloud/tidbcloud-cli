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
	"context"
	"fmt"
	"runtime"
	"time"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/util"
	ver "tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
)

type LoginOpts struct {
	client          *resty.Client
	insecureStorage bool
}

func (l LoginOpts) NonInteractiveFlags() []string {
	return []string{}
}

func LoginCmd(h *internal.Helper) *cobra.Command {
	opts := LoginOpts{
		client: resty.New(),
	}

	var loginCmd = &cobra.Command{
		Use:   "login",
		Short: "Authenticate with the TiDB Cloud",
		Example: fmt.Sprintf(`  To start the login for your account:
  $ %[1]s auth login`, config.CliName),
		RunE: func(cmd *cobra.Command, args []string) error {
			ctx := cmd.Context()
			body := AuthRequest{
				ClientID: clientID,
			}
			result := AuthResponse{}
			response, err := opts.client.
				R().
				SetContext(ctx).
				SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
				SetHeader("Content-type", "application/json").
				SetResult(&result).
				SetBody(body).
				Post(authUrl)
			if err != nil {
				return err
			}
			if !response.IsSuccess() {
				return errors.Errorf("Failed to authorize device, code: %s, body: %s", response.Status(), string(response.Body()))
			}

			if h.IOStreams.CanPrompt {
				fmt.Fprintln(h.IOStreams.Out, "Attempting to automatically open the authorization page in your default browser.")
				fmt.Fprintln(h.IOStreams.Out, "If the browser does not open or you wish to use a different device to authorize this CLI, open the following URL:")
				fmt.Fprintln(h.IOStreams.Out, "\n", result.VerificationURIComplete)
				fmt.Fprintln(h.IOStreams.Out, "\n", "Confirmation Code: ", color.GreenString(result.UserCode))
				openCmd := util.OpenBrowser(runtime.GOOS, result.VerificationURIComplete)
				err = openCmd.Run()
				if err != nil {
					fmt.Fprintln(h.IOStreams.Err, "Failed to open a browser: ", err.Error())
				}
			} else {
				fmt.Fprintln(h.IOStreams.Out, "Please open the following URL in your browser:")
				fmt.Fprintln(h.IOStreams.Out, "\n", result.VerificationURIComplete)
				fmt.Fprintln(h.IOStreams.Out, "\n", "Confirmation Code: ", color.GreenString(result.UserCode))
			}

			now := time.Now()
			token, err := opts.requestForToken(ctx, result)
			if err != nil {
				return err
			}

			color.HiGreen("\nSuccessfully logged in.")

			err = config.SaveAccessToken(now.Add(time.Duration(token.ExpireIn)*time.Second), token.TokenType, token.AccessToken, opts.insecureStorage)
			if err != nil {
				return err
			}

			return nil
		},
	}

	loginCmd.Flags().BoolVar(&opts.insecureStorage, "insecure-storage", false, "Save authentication credentials in plain text instead of credential store")

	return loginCmd
}

func (l LoginOpts) requestForToken(ctx context.Context, result AuthResponse) (*TokenResponse, error) {
	start := time.Now()
	interval := time.Duration(result.PollingInterval) * time.Second
	body := TokenRequest{
		DeviceCode: result.DeviceCode,
		GrantType:  grantType,
		ClientID:   clientID,
	}
	var res TokenResponse
	var tokenError TokenError
	for {
		// This loop begins right after we open the user's browser to send an
		// authentication code. We don't request a token immediately because the
		// has to complete that authentication flow before we can provide a
		// token anyway.
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(interval):
			// Ready to check again.
		}

		resp, err := l.client.R().SetContext(ctx).
			SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
			SetHeader("Content-type", "application/json").
			SetBody(body).
			SetResult(&res).
			SetError(&tokenError).
			Post(accessTokenUrl)
		if err != nil {
			return nil, err
		}
		if !resp.IsSuccess() {
			if tokenError.Error != errPending && tokenError.Error != errSlowDown {
				return nil, errors.Errorf("Failed to get access token, code: %s, body: %s", resp.Status(), string(resp.Body()))
			}
			if tokenError.Error == errSlowDown {
				interval += 1 * time.Second
			}
		}

		if res.AccessToken != "" {
			// Successful authentication.
			return &res, nil
		}

		if time.Now().After(start.Add(time.Duration(result.ExpiresIn) * time.Second)) {
			return nil, errors.New("authentication timed out")
		}
	}
}

type AuthRequest struct {
	ClientID string `json:"client_id"`
}

type AuthResponse struct {
	DeviceCode              string `json:"device_code"`
	UserCode                string `json:"user_code"`
	VerificationURI         string `json:"verification_uri"`
	VerificationURIComplete string `json:"verification_uri_complete"`
	ExpiresIn               int    `json:"expires_in"`
	PollingInterval         int    `json:"interval"`
}

type TokenRequest struct {
	DeviceCode string `json:"device_code"`
	GrantType  string `json:"grant_type"`
	ClientID   string `json:"client_id"`
}

type TokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpireIn     int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type TokenError struct {
	Error string `json:"error"`
}
