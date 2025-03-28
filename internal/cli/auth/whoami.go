// Copyright 2025 PingCAP, Inc.
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
	"errors"
	"fmt"
	"sync"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	ver "github.com/tidbcloud/tidbcloud-cli/internal/version"

	"github.com/fatih/color"
	"github.com/go-resty/resty/v2"
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

			wg.Wait()
			close(results)

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

			if config.GetPublicKey() != "" && config.GetPrivateKey() != "" {
				color.HiYellow("\nDetect an API key already set in %s profile! Note it will take precedence over auth token.", config.ActiveProfileName())
				color.HiYellow(fmt.Sprintf("Use `%s config create --profile-name <profile-name>` to create a new profile and login again.", config.CliName))
			}
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
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
		Get(fmt.Sprintf("%s%s", config.GetOAuthEndpoint(), userInfoPath))

	if err != nil {
		return Result{userInfoPath, nil, err}
	}

	if !resp.IsSuccess() {
		return Result{userInfoPath, nil, fmt.Errorf("%s get %d response", userInfoPath, resp.StatusCode())}
	}

	var userInfo UserInfo
	err = json.Unmarshal(resp.Body(), &userInfo)
	if err != nil {
		return Result{userInfoPath, nil, err}
	}

	return Result{userInfoPath, &userInfo, nil}
}

func getOrgInfo(ctx context.Context, client *resty.Client, token string) Result {
	iamEndpoint := config.GetIAMEndpoint()
	if iamEndpoint == "" {
		iamEndpoint = cloud.DefaultIAMEndpoint
	}
	resp, err := client.R().
		SetContext(ctx).
		SetHeader("Authorization", "Bearer "+token).
		SetHeader("user-agent", fmt.Sprintf("%s/%s", config.CliName, ver.Version)).
		Get(fmt.Sprintf("%s%s", iamEndpoint, orgPath))

	if err != nil {
		return Result{orgPath, nil, err}
	}

	if !resp.IsSuccess() {
		return Result{orgPath, nil, fmt.Errorf("%s get %d response", orgPath, resp.StatusCode())}
	}

	var orgInfo OrgInfo
	err = json.Unmarshal(resp.Body(), &orgInfo)
	if err != nil {
		return Result{orgPath, nil, err}
	}

	return Result{orgPath, &orgInfo, nil}
}
