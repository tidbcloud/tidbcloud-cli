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
	"fmt"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/output"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	api "github.com/tidbcloud/tidbcloud-cli/pkg/tidbcloud/v1beta2/nextgen"

	"github.com/spf13/cobra"
)

func cmekCmd(h *internal.Helper) *cobra.Command {
	root := &cobra.Command{Use: "cmek", Short: "Configure customer-managed encryption for Premium instances"}
	principal := &cobra.Command{
		Use:   "principal",
		Short: "Get the IAM principal for your KMS key policy",
		Long:  "Get the IAM principal required to grant TiDB Cloud access to your AWS or Alibaba Cloud KMS key. Configure the key policy before creating an instance with --encryption cmek.",
		Args:  cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			format, err := validatedOutputFormat(cmd)
			if err != nil {
				return err
			}
			regionID, err := cmd.Flags().GetString(flag.Region)
			if err != nil {
				return err
			}
			provider, err := cmekProvider(regionID)
			if err != nil {
				return err
			}
			client, err := h.NextGenClient()
			if err != nil {
				return err
			}
			result, err := client.GetCmekAccessIAMPrincipal(cmd.Context(), regionID)
			if err != nil {
				return err
			}
			if format == output.JsonFormat || !h.IOStreams.CanPrompt {
				return output.PrintJson(h.IOStreams.Out, result)
			}
			aws := result.GetAws()
			aliyun := result.GetAliyun()
			accountID, externalID := aws.GetAccountId(), aws.GetExternalId()
			if provider == "alicloud" {
				accountID, externalID = aliyun.GetAccountId(), ""
			}
			return output.PrintHumanTable(h.IOStreams.Out,
				[]output.Column{"RegionID", "CloudProvider", "AccountID", "ExternalID"},
				[]output.Row{{regionID, provider, accountID, externalID}})
		},
	}
	principal.Flags().StringP(flag.Region, flag.RegionShort, "", "The region ID, for example aws-us-west-2.")
	principal.Flags().StringP(flag.Output, flag.OutputShort, output.HumanFormat, flag.OutputHelp)
	_ = principal.MarkFlagRequired(flag.Region)
	root.AddCommand(principal)
	return root
}

func cmekProvider(regionID string) (string, error) {
	provider, region, ok := strings.Cut(strings.TrimPrefix(regionID, "regions/"), "-")
	if !ok || region == "" || (provider != "aws" && provider != "alicloud") {
		return "", fmt.Errorf("CMEK requires an AWS or Alibaba Cloud region, for example aws-us-west-2")
	}
	return provider, nil
}

func verifiedCMEK(ctx context.Context, client cloud.NextGenClient, regionID, keyARN string) (*api.V1beta2CustomerManagedEncryptionKey, error) {
	principal, err := client.GetCmekAccessIAMPrincipal(ctx, regionID)
	if err != nil {
		return nil, err
	}
	key := api.NewV1beta2CustomerManagedEncryptionKey(regionID, api.V1BETA1SERVICEPLAN_PREMIUM)
	provider, _ := cmekProvider(regionID) // Validated before obtaining the client.
	switch provider {
	case "aws":
		aws := principal.GetAws()
		if aws.GetAccountId() == "" {
			return nil, fmt.Errorf("the API returned no AWS CMEK principal")
		}
		key.AwsKms = api.NewCustomerManagedEncryptionKeyAwsKms(keyARN)
		key.AwsPrincipal = &aws
	case "alicloud":
		aliyun := principal.GetAliyun()
		if aliyun.GetAccountId() == "" {
			return nil, fmt.Errorf("the API returned no Alibaba Cloud CMEK principal")
		}
		key.AliyunKms = api.NewCustomerManagedEncryptionKeyAliyunKms(keyARN)
		key.AliyunPrincipal = &aliyun
	}
	result, err := client.VerifyCmekAccessIAMPrincipal(ctx, key)
	if err != nil {
		return nil, err
	}
	if !result.GetValid() {
		return nil, fmt.Errorf("CMEK access verification failed: %s", result.GetFailReason())
	}
	return key, nil
}
