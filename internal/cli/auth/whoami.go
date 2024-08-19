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
	"encoding/json"
	"fmt"
	"sync"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/config/store"
	"tidbcloud-cli/internal/flag"
	ver "tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
	"github.com/spf13/cobra"
	_ "github.com/xo/usql/drivers/mysql"
	"github.com/zalando/go-keyring"
)

type WhoamiOpts struct {
	client *resty.Client
}

var wg sync.WaitGroup

func WhoamiCmd(h *internal.Helper) *cobra.Command {
	opts := WhoamiOpts{
		client: resty.New(),
	}
	var whoamiCmd = &cobra.Command{
		Use:   "whoami",
		Short: "Display information about the current user",
		Args:  cobra.NoArgs,
		Example: fmt.Sprintf(`  To display information about the current user:
  $ %[1]s auth whoami`, config.CliName),
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
				if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, store.ErrNotSupported) {
					fmt.Fprintln(h.IOStreams.Out, color.YellowString("You are not logged in."))
					return nil
				}
				return err
			}

			results := make(chan Result, 2)
			wg.Add(2)

			go func() {
				defer wg.Done()
				results <- getUserInfo(ctx, opts.client, token)
			}()

			go func() {
				defer wg.Done()
				results <- getOrgInfo(ctx, opts.client, token)
			}()

			go func() {
				wg.Wait()
				close(results)
			}()

			var userInfo *UserInfo
			var orgInfo *OrgInfo

			for result := range results {
				if result.Error != nil {
					return result.Error
				}

				switch result.API {
				case userInfoPath:
					userInfo = result.Data.(*UserInfo)
				case orgPath:
					orgInfo = result.Data.(*OrgInfo)
				}
			}

			fmt.Fprintln(h.IOStreams.Out, "Email:", userInfo.Email)
			fmt.Fprintln(h.IOStreams.Out, "User Name:", userInfo.Username)
			fmt.Fprintln(h.IOStreams.Out, "Org Name:", orgInfo.Orgname)

			return nil
		},
	}

	return whoamiCmd
}

type UserInfo struct {
	Email    string `json:"email"`
	Username string `json:"username"`
}

type OrgInfo struct {
	Orgname string `json:"orgname"`
}

type Result struct {
	API   string
	Data  interface{}
	Error error
}

func getUserInfo(ctx context.Context, client *resty.Client, token string) Result {
	API := userInfoPath
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
		Get(fmt.Sprintf("%s%s", config.GetOAuthEndpoint(), userInfoPath))

	if err != nil {
		return Result{API, nil, err}
	}

	if !resp.IsSuccess() {
		return Result{API, nil, err}
	}

	var userInfo UserInfo
	err = json.Unmarshal(resp.Body(), &userInfo)
	if err != nil {
		return Result{API, nil, err}
	}

	return Result{API, &userInfo, nil}
}

func getOrgInfo(ctx context.Context, client *resty.Client, token string) Result {
	API := orgPath
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
		Get(fmt.Sprintf("%s%s", config.GetIAMEndpoint(), orgPath))

	if err != nil {
		return Result{API, nil, err}
	}

	if !resp.IsSuccess() {
		return Result{API, nil, err}
	}

	var orgInfo OrgInfo
	err = json.Unmarshal(resp.Body(), &orgInfo)
	if err != nil {
		return Result{API, nil, err}
	}

	return Result{API, &orgInfo, nil}
}
