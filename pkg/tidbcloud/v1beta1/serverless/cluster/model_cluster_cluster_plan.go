/*
TiDB Cloud Starter and Essential API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Starter](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-serverless) and [TiDB Cloud Essential](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#essential) clusters. For [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters, use the [TiDB Cloud Dedicated API](https://docs.pingcap.com/tidbcloud/api/v1beta1/dedicated/). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://serverless.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List all clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://serverless.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Starter and Essential API, including the following resources and endpoints:   - Cluster:  - [List TiDB Cloud Starter and Essential clusters](#tag/Cluster/operation/ClusterService_ListClusters)  - [Create a new TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_CreateCluster)  - [Get details of a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_GetCluster)  - [Delete a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)  - [Update a TiDB Cloud Starter or Essential cluster](#tag/Cluster/operation/ClusterService_PartialUpdateCluster)  - [List available regions for an organization](#tag/Cluster/operation/ClusterService_ListRegions)  - Branch:  - [List branches](#tag/Branch/operation/BranchService_ListBranches)  - [Create a branch](#tag/Branch/operation/BranchService_CreateBranch)  - [Get details of a branch](#tag/Branch/operation/BranchService_GetBranch)  - [Delete a branch](#tag/Branch/operation/BranchService_DeleteBranch)  - [Reset a branch](#tag/Branch/operation/BranchService_ResetBranch)  - Export:  - [List export tasks for a cluster](#tag/Export/operation/ExportService_ListExports)  - [Create an export task](#tag/Export/operation/ExportService_CreateExport)  - [Get details of an export task](#tag/Export/operation/ExportService_GetExport)  - [Delete an export task](#tag/Export/operation/ExportService_DeleteExport)  - [Cancel an export task](#tag/Export/operation/ExportService_CancelExport)  - Import:  - [List import tasks for a cluster](#tag/Import/operation/ImportService_ListImports)  - [Create an import task](#tag/Import/operation/ImportService_CreateImport)  - [Get an import task](#tag/Import/operation/ImportService_GetImport)  - [Cancel an import task](#tag/Import/operation/ImportService_CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
)

// ClusterClusterPlan The plan of the cluster.  - `STARTER`: [TiDB Cloud Starter](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-serverless) - `ESSENTIAL`: [TiDB Cloud Essential](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#essential).
type ClusterClusterPlan string

// List of ClusterClusterPlan
const (
	CLUSTERCLUSTERPLAN_STARTER   ClusterClusterPlan = "STARTER"
	CLUSTERCLUSTERPLAN_ESSENTIAL ClusterClusterPlan = "ESSENTIAL"
)

// All allowed values of ClusterClusterPlan enum
var AllowedClusterClusterPlanEnumValues = []ClusterClusterPlan{
	"STARTER",
	"ESSENTIAL",
}

func (v *ClusterClusterPlan) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ClusterClusterPlan(value)
	for _, existing := range AllowedClusterClusterPlanEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = ClusterClusterPlan(value)
	return nil
}

// NewClusterClusterPlanFromValue returns a pointer to a valid ClusterClusterPlan for the value passed as argument
func NewClusterClusterPlanFromValue(v string) *ClusterClusterPlan {
	ev := ClusterClusterPlan(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ClusterClusterPlan) IsValid() bool {
	for _, existing := range AllowedClusterClusterPlanEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ClusterClusterPlan value
func (v ClusterClusterPlan) Ptr() *ClusterClusterPlan {
	return &v
}

type NullableClusterClusterPlan struct {
	value *ClusterClusterPlan
	isSet bool
}

func (v NullableClusterClusterPlan) Get() *ClusterClusterPlan {
	return v.value
}

func (v *NullableClusterClusterPlan) Set(val *ClusterClusterPlan) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterClusterPlan) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterClusterPlan) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterClusterPlan(val *ClusterClusterPlan) *NullableClusterClusterPlan {
	return &NullableClusterClusterPlan{value: val, isSet: true}
}

func (v NullableClusterClusterPlan) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterClusterPlan) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
