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

// checks if the Dedicatedv1beta1TidbNodeGroup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1TidbNodeGroup{}

// Dedicatedv1beta1TidbNodeGroup struct for Dedicatedv1beta1TidbNodeGroup
type Dedicatedv1beta1TidbNodeGroup struct {
	// The resource name of the TiDB group, in the format of `tidbNodeGroups/{tidb_node_group_id}`.
	Name *string `json:"name,omitempty"`
	// The unique ID of the TiDB group.
	TidbNodeGroupId *string `json:"tidbNodeGroupId,omitempty"`
	// The ID of the cluster that the TiDB node group belongs to.  - This field is **optional** when creating a cluster with the default TiDB group.  - This field is **required** when creating a non-default TiDB group.
	ClusterId *string `json:"clusterId,omitempty"`
	// The display name of the TiDB group.
	DisplayName *string `json:"displayName,omitempty"`
	// The number of TiDB nodes in the TiDB group. It must be greater than or equal to `1`.
	NodeCount int32 `json:"nodeCount"`
	// The endpoints of the TiDB group.
	Endpoints []Dedicatedv1beta1TidbNodeGroupEndpoint `json:"endpoints,omitempty"`
	// The node spec key of the TiDB group. For example, `8C32G`.
	NodeSpecKey *string `json:"nodeSpecKey,omitempty"`
	// The version tag of the node spec resource. The performance and price of the component may vary based on the version tag.
	NodeSpecVersion *string `json:"nodeSpecVersion,omitempty"`
	// The display name of the node spec of the TiDB group. For example, `8 vCPU, 32 GiB`.
	NodeSpecDisplayName *string `json:"nodeSpecDisplayName,omitempty"`
	// Indicates whether this is the default TiDB node group.
	IsDefaultGroup *bool `json:"isDefaultGroup,omitempty"`
	// The current state of the TiDB group.
	State *Dedicatedv1beta1TidbNodeGroupState `json:"state,omitempty"`
	// The progress of node configuration changes.
	NodeChangingProgress *ClusterNodeChangingProgress `json:"nodeChangingProgress,omitempty"`
	// Configures TiProxy settings for this TiDB group. If not specified, the default TiProxy settings is used.
	TiproxySetting *Dedicatedv1beta1TidbNodeGroupTiProxySetting `json:"tiproxySetting,omitempty"`
	// The progress of TiProxy node configuration changes.
	TiproxyNodeChangingProgress *ClusterNodeChangingProgress `json:"tiproxyNodeChangingProgress,omitempty"`
	AdditionalProperties        map[string]interface{}
}

type _Dedicatedv1beta1TidbNodeGroup Dedicatedv1beta1TidbNodeGroup

// NewDedicatedv1beta1TidbNodeGroup instantiates a new Dedicatedv1beta1TidbNodeGroup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1TidbNodeGroup(nodeCount int32) *Dedicatedv1beta1TidbNodeGroup {
	this := Dedicatedv1beta1TidbNodeGroup{}
	this.NodeCount = nodeCount
	return &this
}

// NewDedicatedv1beta1TidbNodeGroupWithDefaults instantiates a new Dedicatedv1beta1TidbNodeGroup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1TidbNodeGroupWithDefaults() *Dedicatedv1beta1TidbNodeGroup {
	this := Dedicatedv1beta1TidbNodeGroup{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetName(v string) {
	o.Name = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		return nil, false
	}
	return o.TidbNodeGroupId, true
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasTidbNodeGroupId() bool {
	if o != nil && !IsNil(o.TidbNodeGroupId) {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given string and assigns it to the TidbNodeGroupId field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetNodeCount returns the NodeCount field value
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeCount, true
}

// SetNodeCount sets field value
func (o *Dedicatedv1beta1TidbNodeGroup) SetNodeCount(v int32) {
	o.NodeCount = v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetEndpoints() []Dedicatedv1beta1TidbNodeGroupEndpoint {
	if o == nil || IsNil(o.Endpoints) {
		var ret []Dedicatedv1beta1TidbNodeGroupEndpoint
		return ret
	}
	return o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetEndpointsOk() ([]Dedicatedv1beta1TidbNodeGroupEndpoint, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []Dedicatedv1beta1TidbNodeGroupEndpoint and assigns it to the Endpoints field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetEndpoints(v []Dedicatedv1beta1TidbNodeGroupEndpoint) {
	o.Endpoints = v
}

// GetNodeSpecKey returns the NodeSpecKey field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecKey() string {
	if o == nil || IsNil(o.NodeSpecKey) {
		var ret string
		return ret
	}
	return *o.NodeSpecKey
}

// GetNodeSpecKeyOk returns a tuple with the NodeSpecKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecKeyOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecKey) {
		return nil, false
	}
	return o.NodeSpecKey, true
}

// HasNodeSpecKey returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasNodeSpecKey() bool {
	if o != nil && !IsNil(o.NodeSpecKey) {
		return true
	}

	return false
}

// SetNodeSpecKey gets a reference to the given string and assigns it to the NodeSpecKey field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetNodeSpecKey(v string) {
	o.NodeSpecKey = &v
}

// GetNodeSpecVersion returns the NodeSpecVersion field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecVersion() string {
	if o == nil || IsNil(o.NodeSpecVersion) {
		var ret string
		return ret
	}
	return *o.NodeSpecVersion
}

// GetNodeSpecVersionOk returns a tuple with the NodeSpecVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecVersionOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecVersion) {
		return nil, false
	}
	return o.NodeSpecVersion, true
}

// HasNodeSpecVersion returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasNodeSpecVersion() bool {
	if o != nil && !IsNil(o.NodeSpecVersion) {
		return true
	}

	return false
}

// SetNodeSpecVersion gets a reference to the given string and assigns it to the NodeSpecVersion field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetNodeSpecVersion(v string) {
	o.NodeSpecVersion = &v
}

// GetNodeSpecDisplayName returns the NodeSpecDisplayName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecDisplayName() string {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		var ret string
		return ret
	}
	return *o.NodeSpecDisplayName
}

// GetNodeSpecDisplayNameOk returns a tuple with the NodeSpecDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeSpecDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		return nil, false
	}
	return o.NodeSpecDisplayName, true
}

// HasNodeSpecDisplayName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasNodeSpecDisplayName() bool {
	if o != nil && !IsNil(o.NodeSpecDisplayName) {
		return true
	}

	return false
}

// SetNodeSpecDisplayName gets a reference to the given string and assigns it to the NodeSpecDisplayName field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetNodeSpecDisplayName(v string) {
	o.NodeSpecDisplayName = &v
}

// GetIsDefaultGroup returns the IsDefaultGroup field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetIsDefaultGroup() bool {
	if o == nil || IsNil(o.IsDefaultGroup) {
		var ret bool
		return ret
	}
	return *o.IsDefaultGroup
}

// GetIsDefaultGroupOk returns a tuple with the IsDefaultGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetIsDefaultGroupOk() (*bool, bool) {
	if o == nil || IsNil(o.IsDefaultGroup) {
		return nil, false
	}
	return o.IsDefaultGroup, true
}

// HasIsDefaultGroup returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasIsDefaultGroup() bool {
	if o != nil && !IsNil(o.IsDefaultGroup) {
		return true
	}

	return false
}

// SetIsDefaultGroup gets a reference to the given bool and assigns it to the IsDefaultGroup field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetIsDefaultGroup(v bool) {
	o.IsDefaultGroup = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetState() Dedicatedv1beta1TidbNodeGroupState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1TidbNodeGroupState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetStateOk() (*Dedicatedv1beta1TidbNodeGroupState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1TidbNodeGroupState and assigns it to the State field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetState(v Dedicatedv1beta1TidbNodeGroupState) {
	o.State = &v
}

// GetNodeChangingProgress returns the NodeChangingProgress field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeChangingProgress() ClusterNodeChangingProgress {
	if o == nil || IsNil(o.NodeChangingProgress) {
		var ret ClusterNodeChangingProgress
		return ret
	}
	return *o.NodeChangingProgress
}

// GetNodeChangingProgressOk returns a tuple with the NodeChangingProgress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetNodeChangingProgressOk() (*ClusterNodeChangingProgress, bool) {
	if o == nil || IsNil(o.NodeChangingProgress) {
		return nil, false
	}
	return o.NodeChangingProgress, true
}

// HasNodeChangingProgress returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasNodeChangingProgress() bool {
	if o != nil && !IsNil(o.NodeChangingProgress) {
		return true
	}

	return false
}

// SetNodeChangingProgress gets a reference to the given ClusterNodeChangingProgress and assigns it to the NodeChangingProgress field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetNodeChangingProgress(v ClusterNodeChangingProgress) {
	o.NodeChangingProgress = &v
}

// GetTiproxySetting returns the TiproxySetting field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTiproxySetting() Dedicatedv1beta1TidbNodeGroupTiProxySetting {
	if o == nil || IsNil(o.TiproxySetting) {
		var ret Dedicatedv1beta1TidbNodeGroupTiProxySetting
		return ret
	}
	return *o.TiproxySetting
}

// GetTiproxySettingOk returns a tuple with the TiproxySetting field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTiproxySettingOk() (*Dedicatedv1beta1TidbNodeGroupTiProxySetting, bool) {
	if o == nil || IsNil(o.TiproxySetting) {
		return nil, false
	}
	return o.TiproxySetting, true
}

// HasTiproxySetting returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasTiproxySetting() bool {
	if o != nil && !IsNil(o.TiproxySetting) {
		return true
	}

	return false
}

// SetTiproxySetting gets a reference to the given Dedicatedv1beta1TidbNodeGroupTiProxySetting and assigns it to the TiproxySetting field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetTiproxySetting(v Dedicatedv1beta1TidbNodeGroupTiProxySetting) {
	o.TiproxySetting = &v
}

// GetTiproxyNodeChangingProgress returns the TiproxyNodeChangingProgress field value if set, zero value otherwise.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTiproxyNodeChangingProgress() ClusterNodeChangingProgress {
	if o == nil || IsNil(o.TiproxyNodeChangingProgress) {
		var ret ClusterNodeChangingProgress
		return ret
	}
	return *o.TiproxyNodeChangingProgress
}

// GetTiproxyNodeChangingProgressOk returns a tuple with the TiproxyNodeChangingProgress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) GetTiproxyNodeChangingProgressOk() (*ClusterNodeChangingProgress, bool) {
	if o == nil || IsNil(o.TiproxyNodeChangingProgress) {
		return nil, false
	}
	return o.TiproxyNodeChangingProgress, true
}

// HasTiproxyNodeChangingProgress returns a boolean if a field has been set.
func (o *Dedicatedv1beta1TidbNodeGroup) HasTiproxyNodeChangingProgress() bool {
	if o != nil && !IsNil(o.TiproxyNodeChangingProgress) {
		return true
	}

	return false
}

// SetTiproxyNodeChangingProgress gets a reference to the given ClusterNodeChangingProgress and assigns it to the TiproxyNodeChangingProgress field.
func (o *Dedicatedv1beta1TidbNodeGroup) SetTiproxyNodeChangingProgress(v ClusterNodeChangingProgress) {
	o.TiproxyNodeChangingProgress = &v
}

func (o Dedicatedv1beta1TidbNodeGroup) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1TidbNodeGroup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbNodeGroupId) {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	toSerialize["nodeCount"] = o.NodeCount
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.NodeSpecKey) {
		toSerialize["nodeSpecKey"] = o.NodeSpecKey
	}
	if !IsNil(o.NodeSpecVersion) {
		toSerialize["nodeSpecVersion"] = o.NodeSpecVersion
	}
	if !IsNil(o.NodeSpecDisplayName) {
		toSerialize["nodeSpecDisplayName"] = o.NodeSpecDisplayName
	}
	if !IsNil(o.IsDefaultGroup) {
		toSerialize["isDefaultGroup"] = o.IsDefaultGroup
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.NodeChangingProgress) {
		toSerialize["nodeChangingProgress"] = o.NodeChangingProgress
	}
	if !IsNil(o.TiproxySetting) {
		toSerialize["tiproxySetting"] = o.TiproxySetting
	}
	if !IsNil(o.TiproxyNodeChangingProgress) {
		toSerialize["tiproxyNodeChangingProgress"] = o.TiproxyNodeChangingProgress
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1TidbNodeGroup) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"nodeCount",
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

	varDedicatedv1beta1TidbNodeGroup := _Dedicatedv1beta1TidbNodeGroup{}

	err = json.Unmarshal(data, &varDedicatedv1beta1TidbNodeGroup)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1TidbNodeGroup(varDedicatedv1beta1TidbNodeGroup)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "nodeCount")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "nodeSpecKey")
		delete(additionalProperties, "nodeSpecVersion")
		delete(additionalProperties, "nodeSpecDisplayName")
		delete(additionalProperties, "isDefaultGroup")
		delete(additionalProperties, "state")
		delete(additionalProperties, "nodeChangingProgress")
		delete(additionalProperties, "tiproxySetting")
		delete(additionalProperties, "tiproxyNodeChangingProgress")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1TidbNodeGroup struct {
	value *Dedicatedv1beta1TidbNodeGroup
	isSet bool
}

func (v NullableDedicatedv1beta1TidbNodeGroup) Get() *Dedicatedv1beta1TidbNodeGroup {
	return v.value
}

func (v *NullableDedicatedv1beta1TidbNodeGroup) Set(val *Dedicatedv1beta1TidbNodeGroup) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1TidbNodeGroup) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1TidbNodeGroup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1TidbNodeGroup(val *Dedicatedv1beta1TidbNodeGroup) *NullableDedicatedv1beta1TidbNodeGroup {
	return &NullableDedicatedv1beta1TidbNodeGroup{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1TidbNodeGroup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1TidbNodeGroup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
