/*
TiDB Cloud Dedicated API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Dedicated](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-dedicated) clusters. For TiDB Cloud Starter or TiDB Cloud Essential clusters, use the [TiDB Cloud Starter and Essential API](). For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/tidbcloud/api-overview/).  # Overview  The TiDB Cloud API is a [REST interface](https://en.wikipedia.org/wiki/Representational_state_transfer) that provides you with programmatic access to manage clusters and related resources within TiDB Cloud.  The API has the following features:  - **JSON entities.** All entities are expressed in JSON. - **HTTPS-only.** You can only access the API via HTTPS, ensuring all the data sent over the network is encrypted with TLS. - **Key-based access and digest authentication.** Before you access TiDB Cloud API, you must generate an API key. All requests are authenticated through [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication), ensuring the API key is never sent over the network.  # Get Started  This guide helps you make your first API call to TiDB Cloud API. You'll learn how to authenticate a request, build a request, and interpret the response.  ## Prerequisites  To complete this guide, you need to perform the following tasks:  - Create a [TiDB Cloud account](https://tidbcloud.com/free-trial) - Install [curl](https://curl.se/)  ## Step 1. Create an API key  To create an API key, log in to your TiDB Cloud console. Navigate to the [**API Keys**](https://tidbcloud.com/org-settings/api-keys) page of your organization, and create an API key.  An API key contains a public key and a private key. Copy and save them in a secure location. You will need to use the API key later in this guide.  For more details about creating API key, refer to [API Key Management](#section/Authentication/API-Key-Management).  ## Step 2. Make your first API call  ### Build an API call  TiDB Cloud API call consists of the following components:  - **A host**. The host for TiDB Cloud API is <https://dedicated.tidbapi.com>. - **An API Key**. The public key and the private key are required for authentication. - **A request**. When submitting data to a resource via `POST`, `PATCH`, or `PUT`, you must submit your payload in JSON.  In this guide, you call the [List clusters](#tag/Cluster/operation/ClusterService_ListClusters) endpoint. For the detailed description of the endpoint, see the [API reference](#tag/Cluster/operation/ClusterService_ListClusters).  ### Call an API endpoint  To get all clusters in your organization, run the following command in your terminal. Remember to change `YOUR_PUBLIC_KEY` to your public key and `YOUR_PRIVATE_KEY` to your private key.  ```shell curl --digest \\  --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\  --request GET \\  --url 'https://dedicated.tidbapi.com/v1beta1/clusters' ```  ## Step 3. Check the response  After making the API call, if the status code in response is `200` and you see details about all clusters in your organization, your request is successful.  # Authentication  The TiDB Cloud API uses [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication). It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, the TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. The role of the API key is always `Organization Owner` currently. 5. Click **Next**. Copy and save the public key and the private key. 6. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 7. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Edit**. 4. You can update the API key description. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # Rate Limiting  The TiDB Cloud API allows up to 100 requests per minute per API key. If you exceed the rate limit, the API returns a `429` error. For more quota, you can [submit a request](https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519) to contact our support team.  Each API request returns the following headers about the limit.  - `X-Ratelimit-Limit-Minute`: The number of requests allowed per minute. It is 100 currently. - `X-Ratelimit-Remaining-Minute`: The number of remaining requests in the current minute. When it reaches `0`, the API returns a `429` error and indicates that you exceed the rate limit. - `X-Ratelimit-Reset`: The time in seconds at which the current rate limit resets.  If you exceed the rate limit, an error response returns like this.  ``` > HTTP/2 429 > date: Fri, 22 Jul 2022 05:28:37 GMT > content-type: application/json > content-length: 66 > x-ratelimit-reset: 23 > x-ratelimit-remaining-minute: 0 > x-ratelimit-limit-minute: 100 > x-kong-response-latency: 2 > server: kong/2.8.1  > {\"details\":[],\"code\":49900007,\"message\":\"The request exceeded the limit of 100 times per apikey per minute. For more quota, please contact us: https://support.pingcap.com/hc/en-us/requests/new?ticket_form_id=7800003722519\"} ```  # API Changelog  This changelog lists all changes to the TiDB Cloud API.  <!-- In reverse chronological order -->  ## 20250812  - Initial release of the TiDB Cloud Dedicated API, including the following resources and endpoints:  * Cluster    * [List clusters](#tag/Cluster/operation/ClusterService_ListClusters)    * [Create a cluster](#tag/Cluster/operation/ClusterService_CreateCluster)    * [Get a cluster](#tag/Cluster/operation/ClusterService_GetCluster)    * [Delete a cluster](#tag/Cluster/operation/ClusterService_DeleteCluster)    * [Update a cluster](#tag/Cluster/operation/ClusterService_UpdateCluster)    * [Pause a cluster](#tag/Cluster/operation/ClusterService_PauseCluster)    * [Resume a cluster](#tag/Cluster/operation/ClusterService_ResumeCluster)    * [Reset the root password of a cluster](#tag/Cluster/operation/ClusterService_ResetRootPassword)    * [List node quotas for your organization](#tag/Cluster/operation/ClusterService_ShowNodeQuota)    * [Get log redaction policy](#tag/Cluster/operation/ClusterService_GetLogRedactionPolicy)   * Region    * [List regions](#tag/Region/operation/RegionService_ListRegions)    * [Get a region](#tag/Region/operation/RegionService_GetRegion)    * [List cloud providers](#tag/Region/operation/RegionService_ShowCloudProviders)    * [List node specs](#tag/Region/operation/RegionService_ListNodeSpecs)    * [Get a node spec](#tag/Region/operation/RegionService_GetNodeSpec)   * Private Endpoint Connection    * [Get private link service for a TiDB node group](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateLinkService)    * [Create a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_CreatePrivateEndpointConnection)    * [List private endpoint connections](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_ListPrivateEndpointConnections)    * [Get a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_GetPrivateEndpointConnection)    * [Delete a private endpoint connection](#tag/Private-Endpoint-Connection/operation/PrivateEndpointConnectionService_DeletePrivateEndpointConnection)   * Import    * [List import tasks](#tag/Import/operation/ListImports)    * [Create an import task](#tag/Import/operation/CreateImport)    * [Get an import task](#tag/Import/operation/GetImport)    * [Cancel an import task](#tag/Import/operation/CancelImport)

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
)

// checks if the V1beta1NodeInstance type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1NodeInstance{}

// V1beta1NodeInstance All fields are output only.
type V1beta1NodeInstance struct {
	// The name of the node instance resource which is formatted as: `clusters/{cluster_id}/nodeInstances/{instance_id}`.
	Name *string `json:"name,omitempty"`
	// The ID of the cluster to which the node instance belongs.
	ClusterId *string `json:"clusterId,omitempty"`
	// The ID of the node instance. It is formatted as: - `tidb-{index}` for TiDB instances in the default TiDB group - `{tidb_group_name}-tidb-{index}` for TiDB instances in non-default TiDB groups - `{component_type}-{index}` for other instances.
	InstanceId *string `json:"instanceId,omitempty"`
	// The component type of the node instance. It can be one of the following values: - `TIDB`: TiDB node instance. - `TIKV`: TiKV node instances. - `TIFLASH`: TiFlash node instances. - `PD`: PD node instances. - `PROXY`: Proxy node instances.
	ComponentType *Dedicatedv1beta1ComponentType `json:"componentType,omitempty"`
	// The state of the node instance. - `CREATING`: indicates the node instance is being created. - `AVAILABLE`: indicates the node instance is available for use. - `DELETING`: indicates the node instance is being deleted. - `UNAVAILABLE`: indicates the node instance is not available.
	State *V1beta1NodeInstanceState `json:"state,omitempty"`
	// The number of vCPUs of the node instance. e.g. `8`.
	VCpu *int32 `json:"vCpu,omitempty"`
	// The memory size of the node instance in GiB. e.g. `32`.
	MemorySizeGi *int32 `json:"memorySizeGi,omitempty"`
	// The availability zone of the node instance, e.g. `us-west-2a`.
	AvailabilityZone *string `json:"availabilityZone,omitempty"`
	// The storage size of the node instance in GiB. e.g. `200`.
	StorageSizeGi *int32 `json:"storageSizeGi,omitempty"`
	// The ID of the TiDB node group that the node instance belongs to.
	TidbNodeGroupId NullableString `json:"tidbNodeGroupId,omitempty"`
	// The display name of the TiDB node group that the node instance belongs to.
	TidbNodeGroupDisplayName NullableString `json:"tidbNodeGroupDisplayName,omitempty"`
	// Indicates whether the TiDB node group that the node instance belongs to is the default one.
	IsDefaultTidbNodeGroup NullableBool `json:"isDefaultTidbNodeGroup,omitempty"`
	// The IOPS of the raft store of the node instance. Only available for instances which have storage. If not set, the default IOPS of raft store will be used.
	RaftStoreIops NullableInt32 `json:"raftStoreIops,omitempty"`
	// The storage type of the node instance. Only available for instances which have storage.
	StorageType *ClusterStorageNodeSettingStorageType `json:"storageType,omitempty"`
	// The display name of the node spec of the node instance. e.g. `8 vCPU, 32 GiB`.
	NodeSpecDisplayName *string `json:"nodeSpecDisplayName,omitempty"`
	// The version tag of the node spec resource. The performance and price of the component may vary based on the version tag.
	NodeSpecVersion *string `json:"nodeSpecVersion,omitempty"`
	// The node spec key of the node instance. e.g. `8C32G`.
	NodeSpecKey          *string `json:"nodeSpecKey,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1NodeInstance V1beta1NodeInstance

// NewV1beta1NodeInstance instantiates a new V1beta1NodeInstance object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1NodeInstance() *V1beta1NodeInstance {
	this := V1beta1NodeInstance{}
	return &this
}

// NewV1beta1NodeInstanceWithDefaults instantiates a new V1beta1NodeInstance object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1NodeInstanceWithDefaults() *V1beta1NodeInstance {
	this := V1beta1NodeInstance{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1NodeInstance) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *V1beta1NodeInstance) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetInstanceId returns the InstanceId field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetInstanceId() string {
	if o == nil || IsNil(o.InstanceId) {
		var ret string
		return ret
	}
	return *o.InstanceId
}

// GetInstanceIdOk returns a tuple with the InstanceId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetInstanceIdOk() (*string, bool) {
	if o == nil || IsNil(o.InstanceId) {
		return nil, false
	}
	return o.InstanceId, true
}

// HasInstanceId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasInstanceId() bool {
	if o != nil && !IsNil(o.InstanceId) {
		return true
	}

	return false
}

// SetInstanceId gets a reference to the given string and assigns it to the InstanceId field.
func (o *V1beta1NodeInstance) SetInstanceId(v string) {
	o.InstanceId = &v
}

// GetComponentType returns the ComponentType field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetComponentType() Dedicatedv1beta1ComponentType {
	if o == nil || IsNil(o.ComponentType) {
		var ret Dedicatedv1beta1ComponentType
		return ret
	}
	return *o.ComponentType
}

// GetComponentTypeOk returns a tuple with the ComponentType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetComponentTypeOk() (*Dedicatedv1beta1ComponentType, bool) {
	if o == nil || IsNil(o.ComponentType) {
		return nil, false
	}
	return o.ComponentType, true
}

// HasComponentType returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasComponentType() bool {
	if o != nil && !IsNil(o.ComponentType) {
		return true
	}

	return false
}

// SetComponentType gets a reference to the given Dedicatedv1beta1ComponentType and assigns it to the ComponentType field.
func (o *V1beta1NodeInstance) SetComponentType(v Dedicatedv1beta1ComponentType) {
	o.ComponentType = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetState() V1beta1NodeInstanceState {
	if o == nil || IsNil(o.State) {
		var ret V1beta1NodeInstanceState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStateOk() (*V1beta1NodeInstanceState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta1NodeInstanceState and assigns it to the State field.
func (o *V1beta1NodeInstance) SetState(v V1beta1NodeInstanceState) {
	o.State = &v
}

// GetVCpu returns the VCpu field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetVCpu() int32 {
	if o == nil || IsNil(o.VCpu) {
		var ret int32
		return ret
	}
	return *o.VCpu
}

// GetVCpuOk returns a tuple with the VCpu field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetVCpuOk() (*int32, bool) {
	if o == nil || IsNil(o.VCpu) {
		return nil, false
	}
	return o.VCpu, true
}

// HasVCpu returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasVCpu() bool {
	if o != nil && !IsNil(o.VCpu) {
		return true
	}

	return false
}

// SetVCpu gets a reference to the given int32 and assigns it to the VCpu field.
func (o *V1beta1NodeInstance) SetVCpu(v int32) {
	o.VCpu = &v
}

// GetMemorySizeGi returns the MemorySizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetMemorySizeGi() int32 {
	if o == nil || IsNil(o.MemorySizeGi) {
		var ret int32
		return ret
	}
	return *o.MemorySizeGi
}

// GetMemorySizeGiOk returns a tuple with the MemorySizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetMemorySizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.MemorySizeGi) {
		return nil, false
	}
	return o.MemorySizeGi, true
}

// HasMemorySizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasMemorySizeGi() bool {
	if o != nil && !IsNil(o.MemorySizeGi) {
		return true
	}

	return false
}

// SetMemorySizeGi gets a reference to the given int32 and assigns it to the MemorySizeGi field.
func (o *V1beta1NodeInstance) SetMemorySizeGi(v int32) {
	o.MemorySizeGi = &v
}

// GetAvailabilityZone returns the AvailabilityZone field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetAvailabilityZone() string {
	if o == nil || IsNil(o.AvailabilityZone) {
		var ret string
		return ret
	}
	return *o.AvailabilityZone
}

// GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetAvailabilityZoneOk() (*string, bool) {
	if o == nil || IsNil(o.AvailabilityZone) {
		return nil, false
	}
	return o.AvailabilityZone, true
}

// HasAvailabilityZone returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasAvailabilityZone() bool {
	if o != nil && !IsNil(o.AvailabilityZone) {
		return true
	}

	return false
}

// SetAvailabilityZone gets a reference to the given string and assigns it to the AvailabilityZone field.
func (o *V1beta1NodeInstance) SetAvailabilityZone(v string) {
	o.AvailabilityZone = &v
}

// GetStorageSizeGi returns the StorageSizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetStorageSizeGi() int32 {
	if o == nil || IsNil(o.StorageSizeGi) {
		var ret int32
		return ret
	}
	return *o.StorageSizeGi
}

// GetStorageSizeGiOk returns a tuple with the StorageSizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStorageSizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.StorageSizeGi) {
		return nil, false
	}
	return o.StorageSizeGi, true
}

// HasStorageSizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasStorageSizeGi() bool {
	if o != nil && !IsNil(o.StorageSizeGi) {
		return true
	}

	return false
}

// SetStorageSizeGi gets a reference to the given int32 and assigns it to the StorageSizeGi field.
func (o *V1beta1NodeInstance) SetStorageSizeGi(v int32) {
	o.StorageSizeGi = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId.Get()) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId.Get()
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.TidbNodeGroupId.Get(), o.TidbNodeGroupId.IsSet()
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasTidbNodeGroupId() bool {
	if o != nil && o.TidbNodeGroupId.IsSet() {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given NullableString and assigns it to the TidbNodeGroupId field.
func (o *V1beta1NodeInstance) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId.Set(&v)
}

// SetTidbNodeGroupIdNil sets the value for TidbNodeGroupId to be an explicit nil
func (o *V1beta1NodeInstance) SetTidbNodeGroupIdNil() {
	o.TidbNodeGroupId.Set(nil)
}

// UnsetTidbNodeGroupId ensures that no value is present for TidbNodeGroupId, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetTidbNodeGroupId() {
	o.TidbNodeGroupId.Unset()
}

// GetTidbNodeGroupDisplayName returns the TidbNodeGroupDisplayName field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetTidbNodeGroupDisplayName() string {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName.Get()) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupDisplayName.Get()
}

// GetTidbNodeGroupDisplayNameOk returns a tuple with the TidbNodeGroupDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetTidbNodeGroupDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.TidbNodeGroupDisplayName.Get(), o.TidbNodeGroupDisplayName.IsSet()
}

// HasTidbNodeGroupDisplayName returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasTidbNodeGroupDisplayName() bool {
	if o != nil && o.TidbNodeGroupDisplayName.IsSet() {
		return true
	}

	return false
}

// SetTidbNodeGroupDisplayName gets a reference to the given NullableString and assigns it to the TidbNodeGroupDisplayName field.
func (o *V1beta1NodeInstance) SetTidbNodeGroupDisplayName(v string) {
	o.TidbNodeGroupDisplayName.Set(&v)
}

// SetTidbNodeGroupDisplayNameNil sets the value for TidbNodeGroupDisplayName to be an explicit nil
func (o *V1beta1NodeInstance) SetTidbNodeGroupDisplayNameNil() {
	o.TidbNodeGroupDisplayName.Set(nil)
}

// UnsetTidbNodeGroupDisplayName ensures that no value is present for TidbNodeGroupDisplayName, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetTidbNodeGroupDisplayName() {
	o.TidbNodeGroupDisplayName.Unset()
}

// GetIsDefaultTidbNodeGroup returns the IsDefaultTidbNodeGroup field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetIsDefaultTidbNodeGroup() bool {
	if o == nil || IsNil(o.IsDefaultTidbNodeGroup.Get()) {
		var ret bool
		return ret
	}
	return *o.IsDefaultTidbNodeGroup.Get()
}

// GetIsDefaultTidbNodeGroupOk returns a tuple with the IsDefaultTidbNodeGroup field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetIsDefaultTidbNodeGroupOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.IsDefaultTidbNodeGroup.Get(), o.IsDefaultTidbNodeGroup.IsSet()
}

// HasIsDefaultTidbNodeGroup returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasIsDefaultTidbNodeGroup() bool {
	if o != nil && o.IsDefaultTidbNodeGroup.IsSet() {
		return true
	}

	return false
}

// SetIsDefaultTidbNodeGroup gets a reference to the given NullableBool and assigns it to the IsDefaultTidbNodeGroup field.
func (o *V1beta1NodeInstance) SetIsDefaultTidbNodeGroup(v bool) {
	o.IsDefaultTidbNodeGroup.Set(&v)
}

// SetIsDefaultTidbNodeGroupNil sets the value for IsDefaultTidbNodeGroup to be an explicit nil
func (o *V1beta1NodeInstance) SetIsDefaultTidbNodeGroupNil() {
	o.IsDefaultTidbNodeGroup.Set(nil)
}

// UnsetIsDefaultTidbNodeGroup ensures that no value is present for IsDefaultTidbNodeGroup, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetIsDefaultTidbNodeGroup() {
	o.IsDefaultTidbNodeGroup.Unset()
}

// GetRaftStoreIops returns the RaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeInstance) GetRaftStoreIops() int32 {
	if o == nil || IsNil(o.RaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.RaftStoreIops.Get()
}

// GetRaftStoreIopsOk returns a tuple with the RaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeInstance) GetRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.RaftStoreIops.Get(), o.RaftStoreIops.IsSet()
}

// HasRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasRaftStoreIops() bool {
	if o != nil && o.RaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the RaftStoreIops field.
func (o *V1beta1NodeInstance) SetRaftStoreIops(v int32) {
	o.RaftStoreIops.Set(&v)
}

// SetRaftStoreIopsNil sets the value for RaftStoreIops to be an explicit nil
func (o *V1beta1NodeInstance) SetRaftStoreIopsNil() {
	o.RaftStoreIops.Set(nil)
}

// UnsetRaftStoreIops ensures that no value is present for RaftStoreIops, not even an explicit nil
func (o *V1beta1NodeInstance) UnsetRaftStoreIops() {
	o.RaftStoreIops.Unset()
}

// GetStorageType returns the StorageType field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetStorageType() ClusterStorageNodeSettingStorageType {
	if o == nil || IsNil(o.StorageType) {
		var ret ClusterStorageNodeSettingStorageType
		return ret
	}
	return *o.StorageType
}

// GetStorageTypeOk returns a tuple with the StorageType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetStorageTypeOk() (*ClusterStorageNodeSettingStorageType, bool) {
	if o == nil || IsNil(o.StorageType) {
		return nil, false
	}
	return o.StorageType, true
}

// HasStorageType returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasStorageType() bool {
	if o != nil && !IsNil(o.StorageType) {
		return true
	}

	return false
}

// SetStorageType gets a reference to the given ClusterStorageNodeSettingStorageType and assigns it to the StorageType field.
func (o *V1beta1NodeInstance) SetStorageType(v ClusterStorageNodeSettingStorageType) {
	o.StorageType = &v
}

// GetNodeSpecDisplayName returns the NodeSpecDisplayName field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetNodeSpecDisplayName() string {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		var ret string
		return ret
	}
	return *o.NodeSpecDisplayName
}

// GetNodeSpecDisplayNameOk returns a tuple with the NodeSpecDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetNodeSpecDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		return nil, false
	}
	return o.NodeSpecDisplayName, true
}

// HasNodeSpecDisplayName returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasNodeSpecDisplayName() bool {
	if o != nil && !IsNil(o.NodeSpecDisplayName) {
		return true
	}

	return false
}

// SetNodeSpecDisplayName gets a reference to the given string and assigns it to the NodeSpecDisplayName field.
func (o *V1beta1NodeInstance) SetNodeSpecDisplayName(v string) {
	o.NodeSpecDisplayName = &v
}

// GetNodeSpecVersion returns the NodeSpecVersion field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetNodeSpecVersion() string {
	if o == nil || IsNil(o.NodeSpecVersion) {
		var ret string
		return ret
	}
	return *o.NodeSpecVersion
}

// GetNodeSpecVersionOk returns a tuple with the NodeSpecVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetNodeSpecVersionOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecVersion) {
		return nil, false
	}
	return o.NodeSpecVersion, true
}

// HasNodeSpecVersion returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasNodeSpecVersion() bool {
	if o != nil && !IsNil(o.NodeSpecVersion) {
		return true
	}

	return false
}

// SetNodeSpecVersion gets a reference to the given string and assigns it to the NodeSpecVersion field.
func (o *V1beta1NodeInstance) SetNodeSpecVersion(v string) {
	o.NodeSpecVersion = &v
}

// GetNodeSpecKey returns the NodeSpecKey field value if set, zero value otherwise.
func (o *V1beta1NodeInstance) GetNodeSpecKey() string {
	if o == nil || IsNil(o.NodeSpecKey) {
		var ret string
		return ret
	}
	return *o.NodeSpecKey
}

// GetNodeSpecKeyOk returns a tuple with the NodeSpecKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeInstance) GetNodeSpecKeyOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecKey) {
		return nil, false
	}
	return o.NodeSpecKey, true
}

// HasNodeSpecKey returns a boolean if a field has been set.
func (o *V1beta1NodeInstance) HasNodeSpecKey() bool {
	if o != nil && !IsNil(o.NodeSpecKey) {
		return true
	}

	return false
}

// SetNodeSpecKey gets a reference to the given string and assigns it to the NodeSpecKey field.
func (o *V1beta1NodeInstance) SetNodeSpecKey(v string) {
	o.NodeSpecKey = &v
}

func (o V1beta1NodeInstance) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1NodeInstance) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	if !IsNil(o.InstanceId) {
		toSerialize["instanceId"] = o.InstanceId
	}
	if !IsNil(o.ComponentType) {
		toSerialize["componentType"] = o.ComponentType
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.VCpu) {
		toSerialize["vCpu"] = o.VCpu
	}
	if !IsNil(o.MemorySizeGi) {
		toSerialize["memorySizeGi"] = o.MemorySizeGi
	}
	if !IsNil(o.AvailabilityZone) {
		toSerialize["availabilityZone"] = o.AvailabilityZone
	}
	if !IsNil(o.StorageSizeGi) {
		toSerialize["storageSizeGi"] = o.StorageSizeGi
	}
	if o.TidbNodeGroupId.IsSet() {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId.Get()
	}
	if o.TidbNodeGroupDisplayName.IsSet() {
		toSerialize["tidbNodeGroupDisplayName"] = o.TidbNodeGroupDisplayName.Get()
	}
	if o.IsDefaultTidbNodeGroup.IsSet() {
		toSerialize["isDefaultTidbNodeGroup"] = o.IsDefaultTidbNodeGroup.Get()
	}
	if o.RaftStoreIops.IsSet() {
		toSerialize["raftStoreIops"] = o.RaftStoreIops.Get()
	}
	if !IsNil(o.StorageType) {
		toSerialize["storageType"] = o.StorageType
	}
	if !IsNil(o.NodeSpecDisplayName) {
		toSerialize["nodeSpecDisplayName"] = o.NodeSpecDisplayName
	}
	if !IsNil(o.NodeSpecVersion) {
		toSerialize["nodeSpecVersion"] = o.NodeSpecVersion
	}
	if !IsNil(o.NodeSpecKey) {
		toSerialize["nodeSpecKey"] = o.NodeSpecKey
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1NodeInstance) UnmarshalJSON(data []byte) (err error) {
	varV1beta1NodeInstance := _V1beta1NodeInstance{}

	err = json.Unmarshal(data, &varV1beta1NodeInstance)

	if err != nil {
		return err
	}

	*o = V1beta1NodeInstance(varV1beta1NodeInstance)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "instanceId")
		delete(additionalProperties, "componentType")
		delete(additionalProperties, "state")
		delete(additionalProperties, "vCpu")
		delete(additionalProperties, "memorySizeGi")
		delete(additionalProperties, "availabilityZone")
		delete(additionalProperties, "storageSizeGi")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "tidbNodeGroupDisplayName")
		delete(additionalProperties, "isDefaultTidbNodeGroup")
		delete(additionalProperties, "raftStoreIops")
		delete(additionalProperties, "storageType")
		delete(additionalProperties, "nodeSpecDisplayName")
		delete(additionalProperties, "nodeSpecVersion")
		delete(additionalProperties, "nodeSpecKey")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1NodeInstance struct {
	value *V1beta1NodeInstance
	isSet bool
}

func (v NullableV1beta1NodeInstance) Get() *V1beta1NodeInstance {
	return v.value
}

func (v *NullableV1beta1NodeInstance) Set(val *V1beta1NodeInstance) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1NodeInstance) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1NodeInstance) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1NodeInstance(val *V1beta1NodeInstance) *NullableV1beta1NodeInstance {
	return &NullableV1beta1NodeInstance{value: val, isSet: true}
}

func (v NullableV1beta1NodeInstance) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1NodeInstance) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
