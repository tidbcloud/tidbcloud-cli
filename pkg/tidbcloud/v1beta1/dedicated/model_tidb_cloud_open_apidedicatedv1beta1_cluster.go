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
	"time"
)

// checks if the TidbCloudOpenApidedicatedv1beta1Cluster type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApidedicatedv1beta1Cluster{}

// TidbCloudOpenApidedicatedv1beta1Cluster Cluster represents a dedicated TiDB cluster.
type TidbCloudOpenApidedicatedv1beta1Cluster struct {
	// The unique identifier for the TiDB cluster, which is generated by the API and follows the format `clusters/{clusterId}`.
	Name *string `json:"name,omitempty"`
	// The ID of the cluster.
	ClusterId *string `json:"clusterId,omitempty"`
	// The user-defined name of the cluster.
	DisplayName string `json:"displayName" validate:"regexp=^[A-Za-z0-9][-A-Za-z0-9]{2,62}[A-Za-z0-9]$"`
	// The region where the cluster is deployed, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-west-2`.
	RegionId string `json:"regionId"`
	// Key-value pairs used to label the cluster. If the `tidb.cloud/project` label is not specified, the cluster is associated with the default project in the creator's organization.   **Note**: Currently, only the `tidb.cloud/project` label key can be specified when creating a new cluster.
	Labels *map[string]string `json:"labels,omitempty"`
	// The configuration for [TiDB nodes](https://docs.pingcap.com/tidbcloud/tidb-cloud-glossary/#tidb-node) in a TiDB Cloud Dedicated cluster.  To view available node specs for a specific region and cloud provider, use the [List node specs](#tag/RegionService/operation/RegionService_ListNodeSpecs) endpoint.
	TidbNodeSetting V1beta1ClusterTidbNodeSetting `json:"tidbNodeSetting"`
	// The configuration for [TiKV nodes](https://docs.pingcap.com/tidbcloud/tidb-cloud-glossary/#tikv-node) in a TiDB Cloud Dedicated cluster.  To view available node specs for a specific region and cloud provider, use the [List node specs](#tag/RegionService/operation/RegionService_ListNodeSpecs) endpoint.
	TikvNodeSetting V1beta1ClusterStorageNodeSetting `json:"tikvNodeSetting"`
	// The configuration for [TiFlash nodes](https://docs.pingcap.com/tidbcloud/tidb-cloud-glossary/#tiflash-node) in a TiDB Cloud Dedicated cluster. If not set, TiFlash is disabled.  To view available node specs for a specific region and cloud provider, use the [List node specs](#tag/RegionService/operation/RegionService_ListNodeSpecs) endpoint.
	TiflashNodeSetting *V1beta1ClusterStorageNodeSetting `json:"tiflashNodeSetting,omitempty"`
	// The port for cluster connections. All network endpoints in the cluster use this port.
	Port int32 `json:"port"`
	// The root password of the cluster. It must be between 8 and 64 characters long and can contain letters, numbers, and special characters.
	RootPassword *string `json:"rootPassword,omitempty" validate:"regexp=^.{8,64}$"`
	// The current state of the cluster.
	State *Commonv1beta1ClusterState `json:"state,omitempty"`
	// The TiDB version of the cluster.
	Version *string `json:"version,omitempty"`
	// The email address or public API key of the user who create the cluster.
	CreatedBy *string `json:"createdBy,omitempty"`
	// The timestamp when the cluster was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// The timestamp when the cluster was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The pause plan configuration of the cluster.
	PausePlan *Dedicatedv1beta1ClusterPausePlan `json:"pausePlan,omitempty"`
	// The display name of the region where the cluster is located. For example, `N. Virginia (us-east-1)`.
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// The cloud provider where the cluster is located.  - `\"aws\"`: Amazon Web Services  - `\"gcp\"`: Google Cloud  - `\"azure\"`: Microsoft Azure  - `\"alicloud\"`: Alibaba Cloud
	CloudProvider *V1beta1RegionCloudProvider `json:"cloudProvider,omitempty"`
	// The annotations for the cluster. The following lists some predefined annotations: - `tidb.cloud/has-set-password`: indicates whether the cluster has a root password set. - `tidb.cloud/available-features`: lists available features of the cluster. - `tidb.cloud/insufficient-vm-resource`: indicates insufficient virtual machine resources during during cluster creation or modification.
	Annotations          *map[string]string `json:"annotations,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApidedicatedv1beta1Cluster TidbCloudOpenApidedicatedv1beta1Cluster

// NewTidbCloudOpenApidedicatedv1beta1Cluster instantiates a new TidbCloudOpenApidedicatedv1beta1Cluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApidedicatedv1beta1Cluster(displayName string, regionId string, tidbNodeSetting V1beta1ClusterTidbNodeSetting, tikvNodeSetting V1beta1ClusterStorageNodeSetting, port int32) *TidbCloudOpenApidedicatedv1beta1Cluster {
	this := TidbCloudOpenApidedicatedv1beta1Cluster{}
	this.DisplayName = displayName
	this.RegionId = regionId
	this.TidbNodeSetting = tidbNodeSetting
	this.TikvNodeSetting = tikvNodeSetting
	this.Port = port
	return &this
}

// NewTidbCloudOpenApidedicatedv1beta1ClusterWithDefaults instantiates a new TidbCloudOpenApidedicatedv1beta1Cluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApidedicatedv1beta1ClusterWithDefaults() *TidbCloudOpenApidedicatedv1beta1Cluster {
	this := TidbCloudOpenApidedicatedv1beta1Cluster{}
	var port int32 = 4000
	this.Port = port
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDisplayName returns the DisplayName field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetRegionId returns the RegionId field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRegionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRegionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionId, true
}

// SetRegionId sets field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetRegionId(v string) {
	o.RegionId = v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetTidbNodeSetting returns the TidbNodeSetting field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTidbNodeSetting() V1beta1ClusterTidbNodeSetting {
	if o == nil {
		var ret V1beta1ClusterTidbNodeSetting
		return ret
	}

	return o.TidbNodeSetting
}

// GetTidbNodeSettingOk returns a tuple with the TidbNodeSetting field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTidbNodeSettingOk() (*V1beta1ClusterTidbNodeSetting, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TidbNodeSetting, true
}

// SetTidbNodeSetting sets field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetTidbNodeSetting(v V1beta1ClusterTidbNodeSetting) {
	o.TidbNodeSetting = v
}

// GetTikvNodeSetting returns the TikvNodeSetting field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTikvNodeSetting() V1beta1ClusterStorageNodeSetting {
	if o == nil {
		var ret V1beta1ClusterStorageNodeSetting
		return ret
	}

	return o.TikvNodeSetting
}

// GetTikvNodeSettingOk returns a tuple with the TikvNodeSetting field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTikvNodeSettingOk() (*V1beta1ClusterStorageNodeSetting, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TikvNodeSetting, true
}

// SetTikvNodeSetting sets field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetTikvNodeSetting(v V1beta1ClusterStorageNodeSetting) {
	o.TikvNodeSetting = v
}

// GetTiflashNodeSetting returns the TiflashNodeSetting field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTiflashNodeSetting() V1beta1ClusterStorageNodeSetting {
	if o == nil || IsNil(o.TiflashNodeSetting) {
		var ret V1beta1ClusterStorageNodeSetting
		return ret
	}
	return *o.TiflashNodeSetting
}

// GetTiflashNodeSettingOk returns a tuple with the TiflashNodeSetting field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetTiflashNodeSettingOk() (*V1beta1ClusterStorageNodeSetting, bool) {
	if o == nil || IsNil(o.TiflashNodeSetting) {
		return nil, false
	}
	return o.TiflashNodeSetting, true
}

// HasTiflashNodeSetting returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasTiflashNodeSetting() bool {
	if o != nil && !IsNil(o.TiflashNodeSetting) {
		return true
	}

	return false
}

// SetTiflashNodeSetting gets a reference to the given V1beta1ClusterStorageNodeSetting and assigns it to the TiflashNodeSetting field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetTiflashNodeSetting(v V1beta1ClusterStorageNodeSetting) {
	o.TiflashNodeSetting = &v
}

// GetPort returns the Port field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetPort() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.Port
}

// GetPortOk returns a tuple with the Port field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetPortOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Port, true
}

// SetPort sets field value
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetPort(v int32) {
	o.Port = v
}

// GetRootPassword returns the RootPassword field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRootPassword() string {
	if o == nil || IsNil(o.RootPassword) {
		var ret string
		return ret
	}
	return *o.RootPassword
}

// GetRootPasswordOk returns a tuple with the RootPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRootPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.RootPassword) {
		return nil, false
	}
	return o.RootPassword, true
}

// HasRootPassword returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasRootPassword() bool {
	if o != nil && !IsNil(o.RootPassword) {
		return true
	}

	return false
}

// SetRootPassword gets a reference to the given string and assigns it to the RootPassword field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetRootPassword(v string) {
	o.RootPassword = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetState() Commonv1beta1ClusterState {
	if o == nil || IsNil(o.State) {
		var ret Commonv1beta1ClusterState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetStateOk() (*Commonv1beta1ClusterState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Commonv1beta1ClusterState and assigns it to the State field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetState(v Commonv1beta1ClusterState) {
	o.State = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetVersion(v string) {
	o.Version = &v
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetUpdateTime returns the UpdateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetUpdateTime() time.Time {
	if o == nil || IsNil(o.UpdateTime) {
		var ret time.Time
		return ret
	}
	return *o.UpdateTime
}

// GetUpdateTimeOk returns a tuple with the UpdateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetUpdateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdateTime) {
		return nil, false
	}
	return o.UpdateTime, true
}

// HasUpdateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasUpdateTime() bool {
	if o != nil && !IsNil(o.UpdateTime) {
		return true
	}

	return false
}

// SetUpdateTime gets a reference to the given time.Time and assigns it to the UpdateTime field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetUpdateTime(v time.Time) {
	o.UpdateTime = &v
}

// GetPausePlan returns the PausePlan field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetPausePlan() Dedicatedv1beta1ClusterPausePlan {
	if o == nil || IsNil(o.PausePlan) {
		var ret Dedicatedv1beta1ClusterPausePlan
		return ret
	}
	return *o.PausePlan
}

// GetPausePlanOk returns a tuple with the PausePlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetPausePlanOk() (*Dedicatedv1beta1ClusterPausePlan, bool) {
	if o == nil || IsNil(o.PausePlan) {
		return nil, false
	}
	return o.PausePlan, true
}

// HasPausePlan returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasPausePlan() bool {
	if o != nil && !IsNil(o.PausePlan) {
		return true
	}

	return false
}

// SetPausePlan gets a reference to the given Dedicatedv1beta1ClusterPausePlan and assigns it to the PausePlan field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetPausePlan(v Dedicatedv1beta1ClusterPausePlan) {
	o.PausePlan = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetAnnotations returns the Annotations field value if set, zero value otherwise.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetAnnotations() map[string]string {
	if o == nil || IsNil(o.Annotations) {
		var ret map[string]string
		return ret
	}
	return *o.Annotations
}

// GetAnnotationsOk returns a tuple with the Annotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) GetAnnotationsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Annotations) {
		return nil, false
	}
	return o.Annotations, true
}

// HasAnnotations returns a boolean if a field has been set.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) HasAnnotations() bool {
	if o != nil && !IsNil(o.Annotations) {
		return true
	}

	return false
}

// SetAnnotations gets a reference to the given map[string]string and assigns it to the Annotations field.
func (o *TidbCloudOpenApidedicatedv1beta1Cluster) SetAnnotations(v map[string]string) {
	o.Annotations = &v
}

func (o TidbCloudOpenApidedicatedv1beta1Cluster) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApidedicatedv1beta1Cluster) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	toSerialize["displayName"] = o.DisplayName
	toSerialize["regionId"] = o.RegionId
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["tidbNodeSetting"] = o.TidbNodeSetting
	toSerialize["tikvNodeSetting"] = o.TikvNodeSetting
	if !IsNil(o.TiflashNodeSetting) {
		toSerialize["tiflashNodeSetting"] = o.TiflashNodeSetting
	}
	toSerialize["port"] = o.Port
	if !IsNil(o.RootPassword) {
		toSerialize["rootPassword"] = o.RootPassword
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["createdBy"] = o.CreatedBy
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.UpdateTime) {
		toSerialize["updateTime"] = o.UpdateTime
	}
	if !IsNil(o.PausePlan) {
		toSerialize["pausePlan"] = o.PausePlan
	}
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}
	if !IsNil(o.Annotations) {
		toSerialize["annotations"] = o.Annotations
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApidedicatedv1beta1Cluster) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"displayName",
		"regionId",
		"tidbNodeSetting",
		"tikvNodeSetting",
		"port",
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

	varTidbCloudOpenApidedicatedv1beta1Cluster := _TidbCloudOpenApidedicatedv1beta1Cluster{}

	err = json.Unmarshal(data, &varTidbCloudOpenApidedicatedv1beta1Cluster)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApidedicatedv1beta1Cluster(varTidbCloudOpenApidedicatedv1beta1Cluster)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "tidbNodeSetting")
		delete(additionalProperties, "tikvNodeSetting")
		delete(additionalProperties, "tiflashNodeSetting")
		delete(additionalProperties, "port")
		delete(additionalProperties, "rootPassword")
		delete(additionalProperties, "state")
		delete(additionalProperties, "version")
		delete(additionalProperties, "createdBy")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "updateTime")
		delete(additionalProperties, "pausePlan")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "annotations")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApidedicatedv1beta1Cluster struct {
	value *TidbCloudOpenApidedicatedv1beta1Cluster
	isSet bool
}

func (v NullableTidbCloudOpenApidedicatedv1beta1Cluster) Get() *TidbCloudOpenApidedicatedv1beta1Cluster {
	return v.value
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1Cluster) Set(val *TidbCloudOpenApidedicatedv1beta1Cluster) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApidedicatedv1beta1Cluster) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1Cluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApidedicatedv1beta1Cluster(val *TidbCloudOpenApidedicatedv1beta1Cluster) *NullableTidbCloudOpenApidedicatedv1beta1Cluster {
	return &NullableTidbCloudOpenApidedicatedv1beta1Cluster{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApidedicatedv1beta1Cluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApidedicatedv1beta1Cluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
