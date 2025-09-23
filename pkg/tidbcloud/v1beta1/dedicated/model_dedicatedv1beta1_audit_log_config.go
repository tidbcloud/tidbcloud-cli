/*
TiDB Cloud Dedicated API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://dedicated.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:  * Cluster    * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)    * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)    * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)    * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)    * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)    * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)    * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)    * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)    * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)    * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)   * Region    * [List regions](#tag/Region/operation/RegionService_ListRegions)    * [Get a region](#tag/Region/operation/RegionService_GetRegion)    * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)    * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)    * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)   * Private Endpoint Connection    * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)    * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)    * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)    * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)    * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)   * Import    * [List import tasks](#tag/Import/operation/ListImports)    * [Create an import task](#tag/Import/operation/CreateImport)    * [Get an import task](#tag/Import/operation/GetImport)    * [Cancel an import task](#tag/Import/operation/CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// checks if the Dedicatedv1beta1AuditLogConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1AuditLogConfig{}

// Dedicatedv1beta1AuditLogConfig struct for Dedicatedv1beta1AuditLogConfig
type Dedicatedv1beta1AuditLogConfig struct {
	Name      *string `json:"name,omitempty"`
	ClusterId string  `json:"clusterId"`
	// Default is disabled.
	Enabled *bool `json:"enabled,omitempty"`
	// Required if bucket_manager is CUSTOMER.
	BucketUri      *string `json:"bucketUri,omitempty"`
	BucketRegionId *string `json:"bucketRegionId,omitempty"`
	AwsRoleArn     *string `json:"awsRoleArn,omitempty"`
	AzureSasToken  *string `json:"azureSasToken,omitempty"`
	// The bucket_manager field is used to indicate who manages the bucket. If this field is not set, the bucket is managed by the customer by default.
	BucketManager *AuditLogConfigBucketManager `json:"bucketManager,omitempty"`
	// Default is redacted.
	Unredacted *bool `json:"unredacted,omitempty"`
	// If unset, default to TEXT.
	Format *V1beta1AuditLogConfigFormat `json:"format,omitempty"`
	// If unset, default to max_size_mib=10, max_age_second=86400 (1 day).
	RotationPolicy *V1beta1AuditLogConfigRotationPolicy `json:"rotationPolicy,omitempty"`
	// Available only when write check failed.
	BucketWriteCheck *V1beta1AuditLogConfigBucketWriteCheck `json:"bucketWriteCheck,omitempty"`
	// Default to false (and not shown), which means the legacy audit log implementation are not used.
	Legacy               NullableBool `json:"legacy,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1AuditLogConfig Dedicatedv1beta1AuditLogConfig

// NewDedicatedv1beta1AuditLogConfig instantiates a new Dedicatedv1beta1AuditLogConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1AuditLogConfig(clusterId string) *Dedicatedv1beta1AuditLogConfig {
	this := Dedicatedv1beta1AuditLogConfig{}
	this.ClusterId = clusterId
	return &this
}

// NewDedicatedv1beta1AuditLogConfigWithDefaults instantiates a new Dedicatedv1beta1AuditLogConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1AuditLogConfigWithDefaults() *Dedicatedv1beta1AuditLogConfig {
	this := Dedicatedv1beta1AuditLogConfig{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1AuditLogConfig) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value
func (o *Dedicatedv1beta1AuditLogConfig) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *Dedicatedv1beta1AuditLogConfig) SetClusterId(v string) {
	o.ClusterId = v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled) {
		var ret bool
		return ret
	}
	return *o.Enabled
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetEnabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Enabled) {
		return nil, false
	}
	return o.Enabled, true
}

// HasEnabled returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasEnabled() bool {
	if o != nil && !IsNil(o.Enabled) {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given bool and assigns it to the Enabled field.
func (o *Dedicatedv1beta1AuditLogConfig) SetEnabled(v bool) {
	o.Enabled = &v
}

// GetBucketUri returns the BucketUri field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketUri() string {
	if o == nil || IsNil(o.BucketUri) {
		var ret string
		return ret
	}
	return *o.BucketUri
}

// GetBucketUriOk returns a tuple with the BucketUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketUriOk() (*string, bool) {
	if o == nil || IsNil(o.BucketUri) {
		return nil, false
	}
	return o.BucketUri, true
}

// HasBucketUri returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasBucketUri() bool {
	if o != nil && !IsNil(o.BucketUri) {
		return true
	}

	return false
}

// SetBucketUri gets a reference to the given string and assigns it to the BucketUri field.
func (o *Dedicatedv1beta1AuditLogConfig) SetBucketUri(v string) {
	o.BucketUri = &v
}

// GetBucketRegionId returns the BucketRegionId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketRegionId() string {
	if o == nil || IsNil(o.BucketRegionId) {
		var ret string
		return ret
	}
	return *o.BucketRegionId
}

// GetBucketRegionIdOk returns a tuple with the BucketRegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.BucketRegionId) {
		return nil, false
	}
	return o.BucketRegionId, true
}

// HasBucketRegionId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasBucketRegionId() bool {
	if o != nil && !IsNil(o.BucketRegionId) {
		return true
	}

	return false
}

// SetBucketRegionId gets a reference to the given string and assigns it to the BucketRegionId field.
func (o *Dedicatedv1beta1AuditLogConfig) SetBucketRegionId(v string) {
	o.BucketRegionId = &v
}

// GetAwsRoleArn returns the AwsRoleArn field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetAwsRoleArn() string {
	if o == nil || IsNil(o.AwsRoleArn) {
		var ret string
		return ret
	}
	return *o.AwsRoleArn
}

// GetAwsRoleArnOk returns a tuple with the AwsRoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetAwsRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.AwsRoleArn) {
		return nil, false
	}
	return o.AwsRoleArn, true
}

// HasAwsRoleArn returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasAwsRoleArn() bool {
	if o != nil && !IsNil(o.AwsRoleArn) {
		return true
	}

	return false
}

// SetAwsRoleArn gets a reference to the given string and assigns it to the AwsRoleArn field.
func (o *Dedicatedv1beta1AuditLogConfig) SetAwsRoleArn(v string) {
	o.AwsRoleArn = &v
}

// GetAzureSasToken returns the AzureSasToken field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetAzureSasToken() string {
	if o == nil || IsNil(o.AzureSasToken) {
		var ret string
		return ret
	}
	return *o.AzureSasToken
}

// GetAzureSasTokenOk returns a tuple with the AzureSasToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetAzureSasTokenOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSasToken) {
		return nil, false
	}
	return o.AzureSasToken, true
}

// HasAzureSasToken returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasAzureSasToken() bool {
	if o != nil && !IsNil(o.AzureSasToken) {
		return true
	}

	return false
}

// SetAzureSasToken gets a reference to the given string and assigns it to the AzureSasToken field.
func (o *Dedicatedv1beta1AuditLogConfig) SetAzureSasToken(v string) {
	o.AzureSasToken = &v
}

// GetBucketManager returns the BucketManager field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketManager() AuditLogConfigBucketManager {
	if o == nil || IsNil(o.BucketManager) {
		var ret AuditLogConfigBucketManager
		return ret
	}
	return *o.BucketManager
}

// GetBucketManagerOk returns a tuple with the BucketManager field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketManagerOk() (*AuditLogConfigBucketManager, bool) {
	if o == nil || IsNil(o.BucketManager) {
		return nil, false
	}
	return o.BucketManager, true
}

// HasBucketManager returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasBucketManager() bool {
	if o != nil && !IsNil(o.BucketManager) {
		return true
	}

	return false
}

// SetBucketManager gets a reference to the given AuditLogConfigBucketManager and assigns it to the BucketManager field.
func (o *Dedicatedv1beta1AuditLogConfig) SetBucketManager(v AuditLogConfigBucketManager) {
	o.BucketManager = &v
}

// GetUnredacted returns the Unredacted field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetUnredacted() bool {
	if o == nil || IsNil(o.Unredacted) {
		var ret bool
		return ret
	}
	return *o.Unredacted
}

// GetUnredactedOk returns a tuple with the Unredacted field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetUnredactedOk() (*bool, bool) {
	if o == nil || IsNil(o.Unredacted) {
		return nil, false
	}
	return o.Unredacted, true
}

// HasUnredacted returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasUnredacted() bool {
	if o != nil && !IsNil(o.Unredacted) {
		return true
	}

	return false
}

// SetUnredacted gets a reference to the given bool and assigns it to the Unredacted field.
func (o *Dedicatedv1beta1AuditLogConfig) SetUnredacted(v bool) {
	o.Unredacted = &v
}

// GetFormat returns the Format field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetFormat() V1beta1AuditLogConfigFormat {
	if o == nil || IsNil(o.Format) {
		var ret V1beta1AuditLogConfigFormat
		return ret
	}
	return *o.Format
}

// GetFormatOk returns a tuple with the Format field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetFormatOk() (*V1beta1AuditLogConfigFormat, bool) {
	if o == nil || IsNil(o.Format) {
		return nil, false
	}
	return o.Format, true
}

// HasFormat returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasFormat() bool {
	if o != nil && !IsNil(o.Format) {
		return true
	}

	return false
}

// SetFormat gets a reference to the given V1beta1AuditLogConfigFormat and assigns it to the Format field.
func (o *Dedicatedv1beta1AuditLogConfig) SetFormat(v V1beta1AuditLogConfigFormat) {
	o.Format = &v
}

// GetRotationPolicy returns the RotationPolicy field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetRotationPolicy() V1beta1AuditLogConfigRotationPolicy {
	if o == nil || IsNil(o.RotationPolicy) {
		var ret V1beta1AuditLogConfigRotationPolicy
		return ret
	}
	return *o.RotationPolicy
}

// GetRotationPolicyOk returns a tuple with the RotationPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetRotationPolicyOk() (*V1beta1AuditLogConfigRotationPolicy, bool) {
	if o == nil || IsNil(o.RotationPolicy) {
		return nil, false
	}
	return o.RotationPolicy, true
}

// HasRotationPolicy returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasRotationPolicy() bool {
	if o != nil && !IsNil(o.RotationPolicy) {
		return true
	}

	return false
}

// SetRotationPolicy gets a reference to the given V1beta1AuditLogConfigRotationPolicy and assigns it to the RotationPolicy field.
func (o *Dedicatedv1beta1AuditLogConfig) SetRotationPolicy(v V1beta1AuditLogConfigRotationPolicy) {
	o.RotationPolicy = &v
}

// GetBucketWriteCheck returns the BucketWriteCheck field value if set, zero value otherwise.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketWriteCheck() V1beta1AuditLogConfigBucketWriteCheck {
	if o == nil || IsNil(o.BucketWriteCheck) {
		var ret V1beta1AuditLogConfigBucketWriteCheck
		return ret
	}
	return *o.BucketWriteCheck
}

// GetBucketWriteCheckOk returns a tuple with the BucketWriteCheck field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1AuditLogConfig) GetBucketWriteCheckOk() (*V1beta1AuditLogConfigBucketWriteCheck, bool) {
	if o == nil || IsNil(o.BucketWriteCheck) {
		return nil, false
	}
	return o.BucketWriteCheck, true
}

// HasBucketWriteCheck returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasBucketWriteCheck() bool {
	if o != nil && !IsNil(o.BucketWriteCheck) {
		return true
	}

	return false
}

// SetBucketWriteCheck gets a reference to the given V1beta1AuditLogConfigBucketWriteCheck and assigns it to the BucketWriteCheck field.
func (o *Dedicatedv1beta1AuditLogConfig) SetBucketWriteCheck(v V1beta1AuditLogConfigBucketWriteCheck) {
	o.BucketWriteCheck = &v
}

// GetLegacy returns the Legacy field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Dedicatedv1beta1AuditLogConfig) GetLegacy() bool {
	if o == nil || IsNil(o.Legacy.Get()) {
		var ret bool
		return ret
	}
	return *o.Legacy.Get()
}

// GetLegacyOk returns a tuple with the Legacy field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Dedicatedv1beta1AuditLogConfig) GetLegacyOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.Legacy.Get(), o.Legacy.IsSet()
}

// HasLegacy returns a boolean if a field has been set.
func (o *Dedicatedv1beta1AuditLogConfig) HasLegacy() bool {
	if o != nil && o.Legacy.IsSet() {
		return true
	}

	return false
}

// SetLegacy gets a reference to the given NullableBool and assigns it to the Legacy field.
func (o *Dedicatedv1beta1AuditLogConfig) SetLegacy(v bool) {
	o.Legacy.Set(&v)
}

// SetLegacyNil sets the value for Legacy to be an explicit nil
func (o *Dedicatedv1beta1AuditLogConfig) SetLegacyNil() {
	o.Legacy.Set(nil)
}

// UnsetLegacy ensures that no value is present for Legacy, not even an explicit nil
func (o *Dedicatedv1beta1AuditLogConfig) UnsetLegacy() {
	o.Legacy.Unset()
}

func (o Dedicatedv1beta1AuditLogConfig) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1AuditLogConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	toSerialize["clusterId"] = o.ClusterId
	if !IsNil(o.Enabled) {
		toSerialize["enabled"] = o.Enabled
	}
	if !IsNil(o.BucketUri) {
		toSerialize["bucketUri"] = o.BucketUri
	}
	if !IsNil(o.BucketRegionId) {
		toSerialize["bucketRegionId"] = o.BucketRegionId
	}
	if !IsNil(o.AwsRoleArn) {
		toSerialize["awsRoleArn"] = o.AwsRoleArn
	}
	if !IsNil(o.AzureSasToken) {
		toSerialize["azureSasToken"] = o.AzureSasToken
	}
	if !IsNil(o.BucketManager) {
		toSerialize["bucketManager"] = o.BucketManager
	}
	if !IsNil(o.Unredacted) {
		toSerialize["unredacted"] = o.Unredacted
	}
	if !IsNil(o.Format) {
		toSerialize["format"] = o.Format
	}
	if !IsNil(o.RotationPolicy) {
		toSerialize["rotationPolicy"] = o.RotationPolicy
	}
	if !IsNil(o.BucketWriteCheck) {
		toSerialize["bucketWriteCheck"] = o.BucketWriteCheck
	}
	if o.Legacy.IsSet() {
		toSerialize["legacy"] = o.Legacy.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1AuditLogConfig) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"clusterId",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err
	}

	for _, requiredProperty := range requiredProperties {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varDedicatedv1beta1AuditLogConfig := _Dedicatedv1beta1AuditLogConfig{}

	err = json.Unmarshal(data, &varDedicatedv1beta1AuditLogConfig)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1AuditLogConfig(varDedicatedv1beta1AuditLogConfig)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "enabled")
		delete(additionalProperties, "bucketUri")
		delete(additionalProperties, "bucketRegionId")
		delete(additionalProperties, "awsRoleArn")
		delete(additionalProperties, "azureSasToken")
		delete(additionalProperties, "bucketManager")
		delete(additionalProperties, "unredacted")
		delete(additionalProperties, "format")
		delete(additionalProperties, "rotationPolicy")
		delete(additionalProperties, "bucketWriteCheck")
		delete(additionalProperties, "legacy")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1AuditLogConfig struct {
	value *Dedicatedv1beta1AuditLogConfig
	isSet bool
}

func (v NullableDedicatedv1beta1AuditLogConfig) Get() *Dedicatedv1beta1AuditLogConfig {
	return v.value
}

func (v *NullableDedicatedv1beta1AuditLogConfig) Set(val *Dedicatedv1beta1AuditLogConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1AuditLogConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1AuditLogConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1AuditLogConfig(val *Dedicatedv1beta1AuditLogConfig) *NullableDedicatedv1beta1AuditLogConfig {
	return &NullableDedicatedv1beta1AuditLogConfig{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1AuditLogConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1AuditLogConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
