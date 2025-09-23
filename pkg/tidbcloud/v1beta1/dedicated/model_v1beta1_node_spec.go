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

// checks if the V1beta1NodeSpec type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1NodeSpec{}

// V1beta1NodeSpec All fields are output only.
type V1beta1NodeSpec struct {
	// The name of the node spec resource, in the format of `regions/{region_id}/componentTypes/{component_type}/nodeSpecs/{node_spec_key}`. For example, `regions/aws-us-west-2/componentTypes/TIKV/nodeSpecs/8C32G`.
	Name *string `json:"name,omitempty"`
	// The region ID of the node spec resource, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-west-2`.
	RegionId *string `json:"regionId,omitempty"`
	// The component type of the node spec.
	ComponentType *Dedicatedv1beta1ComponentType `json:"componentType,omitempty"`
	// The key of the node spec. For example, `8C32G`.
	NodeSpecKey *string `json:"nodeSpecKey,omitempty"`
	// The display name of the node spec. For example, `8 vCPU, 32 GiB`.
	DisplayName *string `json:"displayName,omitempty"`
	// The number of virtual CPUs (vCPUs) allocated to the node spec. For example, `8`.
	VCpu *int32 `json:"vCpu,omitempty"`
	// The amount of memory in gibibytes (GiB) allocated to the node spec. For example, `32`.
	MemorySizeGi *int32 `json:"memorySizeGi,omitempty"`
	// The default storage size of the node spec resource in GiB.
	DefaultStorageSizeGi *int32 `json:"defaultStorageSizeGi,omitempty"`
	// The maximum storage size of the node spec resource in GiB.
	MaxStorageSizeGi *int32 `json:"maxStorageSizeGi,omitempty"`
	// The minimum storage size of the node spec resource in GiB.
	MinStorageSizeGi *int32 `json:"minStorageSizeGi,omitempty"`
	// The default number of nodes for the node spec resource.
	DefaultNodeCount *int32 `json:"defaultNodeCount,omitempty"`
	// The storage types supported by the node spec resource.
	StorageTypes []ClusterStorageNodeSettingStorageType `json:"storageTypes,omitempty"`
	// The maximum IOPS for Raft log storage of the node spec resource. Currently, this parameter is only useful when overriding IOPS for Raft log storage.
	MaxRaftStoreIops NullableInt32 `json:"maxRaftStoreIops,omitempty"`
	// The minimum IOPS for Raft log storage of the node spec resource. Currently, this parameter is only useful when overriding IOPS for Raft log storage.
	MinRaftStoreIops NullableInt32 `json:"minRaftStoreIops,omitempty"`
	// The default IOPS for Raft log storage of the node spec resource. Currently, this parameter is only useful when overriding IOPS for Raft log storage.
	DefaultRaftStoreIops NullableInt32 `json:"defaultRaftStoreIops,omitempty"`
	// The version tag of the node spec resource. The performance and price of the component may vary based on the version tag.
	Version *string `json:"version,omitempty"`
	// Indicates whether this is the default node spec.
	Default              *bool `json:"default,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1NodeSpec V1beta1NodeSpec

// NewV1beta1NodeSpec instantiates a new V1beta1NodeSpec object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1NodeSpec() *V1beta1NodeSpec {
	this := V1beta1NodeSpec{}
	return &this
}

// NewV1beta1NodeSpecWithDefaults instantiates a new V1beta1NodeSpec object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1NodeSpecWithDefaults() *V1beta1NodeSpec {
	this := V1beta1NodeSpec{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1NodeSpec) SetName(v string) {
	o.Name = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *V1beta1NodeSpec) SetRegionId(v string) {
	o.RegionId = &v
}

// GetComponentType returns the ComponentType field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetComponentType() Dedicatedv1beta1ComponentType {
	if o == nil || IsNil(o.ComponentType) {
		var ret Dedicatedv1beta1ComponentType
		return ret
	}
	return *o.ComponentType
}

// GetComponentTypeOk returns a tuple with the ComponentType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetComponentTypeOk() (*Dedicatedv1beta1ComponentType, bool) {
	if o == nil || IsNil(o.ComponentType) {
		return nil, false
	}
	return o.ComponentType, true
}

// HasComponentType returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasComponentType() bool {
	if o != nil && !IsNil(o.ComponentType) {
		return true
	}

	return false
}

// SetComponentType gets a reference to the given Dedicatedv1beta1ComponentType and assigns it to the ComponentType field.
func (o *V1beta1NodeSpec) SetComponentType(v Dedicatedv1beta1ComponentType) {
	o.ComponentType = &v
}

// GetNodeSpecKey returns the NodeSpecKey field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetNodeSpecKey() string {
	if o == nil || IsNil(o.NodeSpecKey) {
		var ret string
		return ret
	}
	return *o.NodeSpecKey
}

// GetNodeSpecKeyOk returns a tuple with the NodeSpecKey field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetNodeSpecKeyOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecKey) {
		return nil, false
	}
	return o.NodeSpecKey, true
}

// HasNodeSpecKey returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasNodeSpecKey() bool {
	if o != nil && !IsNil(o.NodeSpecKey) {
		return true
	}

	return false
}

// SetNodeSpecKey gets a reference to the given string and assigns it to the NodeSpecKey field.
func (o *V1beta1NodeSpec) SetNodeSpecKey(v string) {
	o.NodeSpecKey = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *V1beta1NodeSpec) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetVCpu returns the VCpu field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetVCpu() int32 {
	if o == nil || IsNil(o.VCpu) {
		var ret int32
		return ret
	}
	return *o.VCpu
}

// GetVCpuOk returns a tuple with the VCpu field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetVCpuOk() (*int32, bool) {
	if o == nil || IsNil(o.VCpu) {
		return nil, false
	}
	return o.VCpu, true
}

// HasVCpu returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasVCpu() bool {
	if o != nil && !IsNil(o.VCpu) {
		return true
	}

	return false
}

// SetVCpu gets a reference to the given int32 and assigns it to the VCpu field.
func (o *V1beta1NodeSpec) SetVCpu(v int32) {
	o.VCpu = &v
}

// GetMemorySizeGi returns the MemorySizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetMemorySizeGi() int32 {
	if o == nil || IsNil(o.MemorySizeGi) {
		var ret int32
		return ret
	}
	return *o.MemorySizeGi
}

// GetMemorySizeGiOk returns a tuple with the MemorySizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetMemorySizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.MemorySizeGi) {
		return nil, false
	}
	return o.MemorySizeGi, true
}

// HasMemorySizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasMemorySizeGi() bool {
	if o != nil && !IsNil(o.MemorySizeGi) {
		return true
	}

	return false
}

// SetMemorySizeGi gets a reference to the given int32 and assigns it to the MemorySizeGi field.
func (o *V1beta1NodeSpec) SetMemorySizeGi(v int32) {
	o.MemorySizeGi = &v
}

// GetDefaultStorageSizeGi returns the DefaultStorageSizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetDefaultStorageSizeGi() int32 {
	if o == nil || IsNil(o.DefaultStorageSizeGi) {
		var ret int32
		return ret
	}
	return *o.DefaultStorageSizeGi
}

// GetDefaultStorageSizeGiOk returns a tuple with the DefaultStorageSizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetDefaultStorageSizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.DefaultStorageSizeGi) {
		return nil, false
	}
	return o.DefaultStorageSizeGi, true
}

// HasDefaultStorageSizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasDefaultStorageSizeGi() bool {
	if o != nil && !IsNil(o.DefaultStorageSizeGi) {
		return true
	}

	return false
}

// SetDefaultStorageSizeGi gets a reference to the given int32 and assigns it to the DefaultStorageSizeGi field.
func (o *V1beta1NodeSpec) SetDefaultStorageSizeGi(v int32) {
	o.DefaultStorageSizeGi = &v
}

// GetMaxStorageSizeGi returns the MaxStorageSizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetMaxStorageSizeGi() int32 {
	if o == nil || IsNil(o.MaxStorageSizeGi) {
		var ret int32
		return ret
	}
	return *o.MaxStorageSizeGi
}

// GetMaxStorageSizeGiOk returns a tuple with the MaxStorageSizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetMaxStorageSizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.MaxStorageSizeGi) {
		return nil, false
	}
	return o.MaxStorageSizeGi, true
}

// HasMaxStorageSizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasMaxStorageSizeGi() bool {
	if o != nil && !IsNil(o.MaxStorageSizeGi) {
		return true
	}

	return false
}

// SetMaxStorageSizeGi gets a reference to the given int32 and assigns it to the MaxStorageSizeGi field.
func (o *V1beta1NodeSpec) SetMaxStorageSizeGi(v int32) {
	o.MaxStorageSizeGi = &v
}

// GetMinStorageSizeGi returns the MinStorageSizeGi field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetMinStorageSizeGi() int32 {
	if o == nil || IsNil(o.MinStorageSizeGi) {
		var ret int32
		return ret
	}
	return *o.MinStorageSizeGi
}

// GetMinStorageSizeGiOk returns a tuple with the MinStorageSizeGi field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetMinStorageSizeGiOk() (*int32, bool) {
	if o == nil || IsNil(o.MinStorageSizeGi) {
		return nil, false
	}
	return o.MinStorageSizeGi, true
}

// HasMinStorageSizeGi returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasMinStorageSizeGi() bool {
	if o != nil && !IsNil(o.MinStorageSizeGi) {
		return true
	}

	return false
}

// SetMinStorageSizeGi gets a reference to the given int32 and assigns it to the MinStorageSizeGi field.
func (o *V1beta1NodeSpec) SetMinStorageSizeGi(v int32) {
	o.MinStorageSizeGi = &v
}

// GetDefaultNodeCount returns the DefaultNodeCount field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetDefaultNodeCount() int32 {
	if o == nil || IsNil(o.DefaultNodeCount) {
		var ret int32
		return ret
	}
	return *o.DefaultNodeCount
}

// GetDefaultNodeCountOk returns a tuple with the DefaultNodeCount field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetDefaultNodeCountOk() (*int32, bool) {
	if o == nil || IsNil(o.DefaultNodeCount) {
		return nil, false
	}
	return o.DefaultNodeCount, true
}

// HasDefaultNodeCount returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasDefaultNodeCount() bool {
	if o != nil && !IsNil(o.DefaultNodeCount) {
		return true
	}

	return false
}

// SetDefaultNodeCount gets a reference to the given int32 and assigns it to the DefaultNodeCount field.
func (o *V1beta1NodeSpec) SetDefaultNodeCount(v int32) {
	o.DefaultNodeCount = &v
}

// GetStorageTypes returns the StorageTypes field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetStorageTypes() []ClusterStorageNodeSettingStorageType {
	if o == nil || IsNil(o.StorageTypes) {
		var ret []ClusterStorageNodeSettingStorageType
		return ret
	}
	return o.StorageTypes
}

// GetStorageTypesOk returns a tuple with the StorageTypes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetStorageTypesOk() ([]ClusterStorageNodeSettingStorageType, bool) {
	if o == nil || IsNil(o.StorageTypes) {
		return nil, false
	}
	return o.StorageTypes, true
}

// HasStorageTypes returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasStorageTypes() bool {
	if o != nil && !IsNil(o.StorageTypes) {
		return true
	}

	return false
}

// SetStorageTypes gets a reference to the given []ClusterStorageNodeSettingStorageType and assigns it to the StorageTypes field.
func (o *V1beta1NodeSpec) SetStorageTypes(v []ClusterStorageNodeSettingStorageType) {
	o.StorageTypes = v
}

// GetMaxRaftStoreIops returns the MaxRaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeSpec) GetMaxRaftStoreIops() int32 {
	if o == nil || IsNil(o.MaxRaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.MaxRaftStoreIops.Get()
}

// GetMaxRaftStoreIopsOk returns a tuple with the MaxRaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeSpec) GetMaxRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.MaxRaftStoreIops.Get(), o.MaxRaftStoreIops.IsSet()
}

// HasMaxRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasMaxRaftStoreIops() bool {
	if o != nil && o.MaxRaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetMaxRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the MaxRaftStoreIops field.
func (o *V1beta1NodeSpec) SetMaxRaftStoreIops(v int32) {
	o.MaxRaftStoreIops.Set(&v)
}

// SetMaxRaftStoreIopsNil sets the value for MaxRaftStoreIops to be an explicit nil
func (o *V1beta1NodeSpec) SetMaxRaftStoreIopsNil() {
	o.MaxRaftStoreIops.Set(nil)
}

// UnsetMaxRaftStoreIops ensures that no value is present for MaxRaftStoreIops, not even an explicit nil
func (o *V1beta1NodeSpec) UnsetMaxRaftStoreIops() {
	o.MaxRaftStoreIops.Unset()
}

// GetMinRaftStoreIops returns the MinRaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeSpec) GetMinRaftStoreIops() int32 {
	if o == nil || IsNil(o.MinRaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.MinRaftStoreIops.Get()
}

// GetMinRaftStoreIopsOk returns a tuple with the MinRaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeSpec) GetMinRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.MinRaftStoreIops.Get(), o.MinRaftStoreIops.IsSet()
}

// HasMinRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasMinRaftStoreIops() bool {
	if o != nil && o.MinRaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetMinRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the MinRaftStoreIops field.
func (o *V1beta1NodeSpec) SetMinRaftStoreIops(v int32) {
	o.MinRaftStoreIops.Set(&v)
}

// SetMinRaftStoreIopsNil sets the value for MinRaftStoreIops to be an explicit nil
func (o *V1beta1NodeSpec) SetMinRaftStoreIopsNil() {
	o.MinRaftStoreIops.Set(nil)
}

// UnsetMinRaftStoreIops ensures that no value is present for MinRaftStoreIops, not even an explicit nil
func (o *V1beta1NodeSpec) UnsetMinRaftStoreIops() {
	o.MinRaftStoreIops.Unset()
}

// GetDefaultRaftStoreIops returns the DefaultRaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1NodeSpec) GetDefaultRaftStoreIops() int32 {
	if o == nil || IsNil(o.DefaultRaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.DefaultRaftStoreIops.Get()
}

// GetDefaultRaftStoreIopsOk returns a tuple with the DefaultRaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1NodeSpec) GetDefaultRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.DefaultRaftStoreIops.Get(), o.DefaultRaftStoreIops.IsSet()
}

// HasDefaultRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasDefaultRaftStoreIops() bool {
	if o != nil && o.DefaultRaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetDefaultRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the DefaultRaftStoreIops field.
func (o *V1beta1NodeSpec) SetDefaultRaftStoreIops(v int32) {
	o.DefaultRaftStoreIops.Set(&v)
}

// SetDefaultRaftStoreIopsNil sets the value for DefaultRaftStoreIops to be an explicit nil
func (o *V1beta1NodeSpec) SetDefaultRaftStoreIopsNil() {
	o.DefaultRaftStoreIops.Set(nil)
}

// UnsetDefaultRaftStoreIops ensures that no value is present for DefaultRaftStoreIops, not even an explicit nil
func (o *V1beta1NodeSpec) UnsetDefaultRaftStoreIops() {
	o.DefaultRaftStoreIops.Unset()
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *V1beta1NodeSpec) SetVersion(v string) {
	o.Version = &v
}

// GetDefault returns the Default field value if set, zero value otherwise.
func (o *V1beta1NodeSpec) GetDefault() bool {
	if o == nil || IsNil(o.Default) {
		var ret bool
		return ret
	}
	return *o.Default
}

// GetDefaultOk returns a tuple with the Default field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NodeSpec) GetDefaultOk() (*bool, bool) {
	if o == nil || IsNil(o.Default) {
		return nil, false
	}
	return o.Default, true
}

// HasDefault returns a boolean if a field has been set.
func (o *V1beta1NodeSpec) HasDefault() bool {
	if o != nil && !IsNil(o.Default) {
		return true
	}

	return false
}

// SetDefault gets a reference to the given bool and assigns it to the Default field.
func (o *V1beta1NodeSpec) SetDefault(v bool) {
	o.Default = &v
}

func (o V1beta1NodeSpec) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1NodeSpec) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	if !IsNil(o.ComponentType) {
		toSerialize["componentType"] = o.ComponentType
	}
	if !IsNil(o.NodeSpecKey) {
		toSerialize["nodeSpecKey"] = o.NodeSpecKey
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	if !IsNil(o.VCpu) {
		toSerialize["vCpu"] = o.VCpu
	}
	if !IsNil(o.MemorySizeGi) {
		toSerialize["memorySizeGi"] = o.MemorySizeGi
	}
	if !IsNil(o.DefaultStorageSizeGi) {
		toSerialize["defaultStorageSizeGi"] = o.DefaultStorageSizeGi
	}
	if !IsNil(o.MaxStorageSizeGi) {
		toSerialize["maxStorageSizeGi"] = o.MaxStorageSizeGi
	}
	if !IsNil(o.MinStorageSizeGi) {
		toSerialize["minStorageSizeGi"] = o.MinStorageSizeGi
	}
	if !IsNil(o.DefaultNodeCount) {
		toSerialize["defaultNodeCount"] = o.DefaultNodeCount
	}
	if !IsNil(o.StorageTypes) {
		toSerialize["storageTypes"] = o.StorageTypes
	}
	if o.MaxRaftStoreIops.IsSet() {
		toSerialize["maxRaftStoreIops"] = o.MaxRaftStoreIops.Get()
	}
	if o.MinRaftStoreIops.IsSet() {
		toSerialize["minRaftStoreIops"] = o.MinRaftStoreIops.Get()
	}
	if o.DefaultRaftStoreIops.IsSet() {
		toSerialize["defaultRaftStoreIops"] = o.DefaultRaftStoreIops.Get()
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.Default) {
		toSerialize["default"] = o.Default
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1NodeSpec) UnmarshalJSON(data []byte) (err error) {
	varV1beta1NodeSpec := _V1beta1NodeSpec{}

	err = json.Unmarshal(data, &varV1beta1NodeSpec)

	if err != nil {
		return err
	}

	*o = V1beta1NodeSpec(varV1beta1NodeSpec)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "componentType")
		delete(additionalProperties, "nodeSpecKey")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "vCpu")
		delete(additionalProperties, "memorySizeGi")
		delete(additionalProperties, "defaultStorageSizeGi")
		delete(additionalProperties, "maxStorageSizeGi")
		delete(additionalProperties, "minStorageSizeGi")
		delete(additionalProperties, "defaultNodeCount")
		delete(additionalProperties, "storageTypes")
		delete(additionalProperties, "maxRaftStoreIops")
		delete(additionalProperties, "minRaftStoreIops")
		delete(additionalProperties, "defaultRaftStoreIops")
		delete(additionalProperties, "version")
		delete(additionalProperties, "default")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1NodeSpec struct {
	value *V1beta1NodeSpec
	isSet bool
}

func (v NullableV1beta1NodeSpec) Get() *V1beta1NodeSpec {
	return v.value
}

func (v *NullableV1beta1NodeSpec) Set(val *V1beta1NodeSpec) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1NodeSpec) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1NodeSpec) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1NodeSpec(val *V1beta1NodeSpec) *NullableV1beta1NodeSpec {
	return &NullableV1beta1NodeSpec{value: val, isSet: true}
}

func (v NullableV1beta1NodeSpec) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1NodeSpec) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
