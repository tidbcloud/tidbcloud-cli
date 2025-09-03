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

// checks if the TidbCloudOpenApiserverlessv1beta1ListClustersResponse type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}

// TidbCloudOpenApiserverlessv1beta1ListClustersResponse Responses message to the request for listing of TiDB Cloud Starter and Essential clusters.
type TidbCloudOpenApiserverlessv1beta1ListClustersResponse struct {
	// A list of clusters that match the query.
	Clusters []TidbCloudOpenApiserverlessv1beta1Cluster `json:"clusters,omitempty"`
	// The token to retrieve the next page of results. Use this value as the `pageToken` parameter in the next request. This field is empty when there are no more pages.
	NextPageToken *string `json:"nextPageToken,omitempty"`
	// The total number of clusters that match the query.
	TotalSize            *int64 `json:"totalSize,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApiserverlessv1beta1ListClustersResponse TidbCloudOpenApiserverlessv1beta1ListClustersResponse

// NewTidbCloudOpenApiserverlessv1beta1ListClustersResponse instantiates a new TidbCloudOpenApiserverlessv1beta1ListClustersResponse object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApiserverlessv1beta1ListClustersResponse() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	this := TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}
	return &this
}

// NewTidbCloudOpenApiserverlessv1beta1ListClustersResponseWithDefaults instantiates a new TidbCloudOpenApiserverlessv1beta1ListClustersResponse object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApiserverlessv1beta1ListClustersResponseWithDefaults() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	this := TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}
	return &this
}

// GetClusters returns the Clusters field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetClusters() []TidbCloudOpenApiserverlessv1beta1Cluster {
	if o == nil || IsNil(o.Clusters) {
		var ret []TidbCloudOpenApiserverlessv1beta1Cluster
		return ret
	}
	return o.Clusters
}

// GetClustersOk returns a tuple with the Clusters field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetClustersOk() ([]TidbCloudOpenApiserverlessv1beta1Cluster, bool) {
	if o == nil || IsNil(o.Clusters) {
		return nil, false
	}
	return o.Clusters, true
}

// HasClusters returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasClusters() bool {
	if o != nil && !IsNil(o.Clusters) {
		return true
	}

	return false
}

// SetClusters gets a reference to the given []TidbCloudOpenApiserverlessv1beta1Cluster and assigns it to the Clusters field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetClusters(v []TidbCloudOpenApiserverlessv1beta1Cluster) {
	o.Clusters = v
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

// GetTotalSize returns the TotalSize field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetTotalSize() int64 {
	if o == nil || IsNil(o.TotalSize) {
		var ret int64
		return ret
	}
	return *o.TotalSize
}

// GetTotalSizeOk returns a tuple with the TotalSize field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) GetTotalSizeOk() (*int64, bool) {
	if o == nil || IsNil(o.TotalSize) {
		return nil, false
	}
	return o.TotalSize, true
}

// HasTotalSize returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) HasTotalSize() bool {
	if o != nil && !IsNil(o.TotalSize) {
		return true
	}

	return false
}

// SetTotalSize gets a reference to the given int64 and assigns it to the TotalSize field.
func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) SetTotalSize(v int64) {
	o.TotalSize = &v
}

func (o TidbCloudOpenApiserverlessv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApiserverlessv1beta1ListClustersResponse) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Clusters) {
		toSerialize["clusters"] = o.Clusters
	}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}
	if !IsNil(o.TotalSize) {
		toSerialize["totalSize"] = o.TotalSize
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) UnmarshalJSON(data []byte) (err error) {
	varTidbCloudOpenApiserverlessv1beta1ListClustersResponse := _TidbCloudOpenApiserverlessv1beta1ListClustersResponse{}

	err = json.Unmarshal(data, &varTidbCloudOpenApiserverlessv1beta1ListClustersResponse)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApiserverlessv1beta1ListClustersResponse(varTidbCloudOpenApiserverlessv1beta1ListClustersResponse)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusters")
		delete(additionalProperties, "nextPageToken")
		delete(additionalProperties, "totalSize")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse struct {
	value *TidbCloudOpenApiserverlessv1beta1ListClustersResponse
	isSet bool
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Get() *TidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	return v.value
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Set(val *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse(val *TidbCloudOpenApiserverlessv1beta1ListClustersResponse) *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse {
	return &NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1ListClustersResponse) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
