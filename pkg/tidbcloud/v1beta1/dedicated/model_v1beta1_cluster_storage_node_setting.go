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

// checks if the V1beta1ClusterStorageNodeSetting type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1ClusterStorageNodeSetting{}

// V1beta1ClusterStorageNodeSetting struct for V1beta1ClusterStorageNodeSetting
type V1beta1ClusterStorageNodeSetting struct {
	// The number of nodes in the cluster.
	NodeCount int32 `json:"nodeCount"`
	// The node spec key of the nodes in the cluster. For example, `8C32G`.
	NodeSpecKey string `json:"nodeSpecKey"`
	// The version tag of the node spec resource. The performance and price of the component may vary based on the version tag.
	NodeSpecVersion *string `json:"nodeSpecVersion,omitempty"`
	// The storage size of each node in gibibytes (GiB). To view supported storage size ranges, use the [List node specs](#tag/RegionService/operation/RegionService_ListNodeSpecs) or [Get a node spec](#tag/RegionService/operation/RegionService_GetNodeSpec) endpoint.
	StorageSizeGi int32 `json:"storageSizeGi"`
	// The storage type of the node. For more information, see [TiKV node storage types](https://docs.pingcap.com/tidbcloud/size-your-cluster/#tikv-node-storage-types) and [TiFlash node storage types](https://docs.pingcap.com/tidbcloud/size-your-cluster/#tiflash-node-storage-types).
	StorageType *ClusterStorageNodeSettingStorageType `json:"storageType,omitempty"`
	// The IOPS (input/output operations per second) of the Raft store for the node. If not specified, the default IOPS of Raft store is used. It is recommended to leave it blank to use the default settings for optimal performance.
	RaftStoreIops NullableInt32 `json:"raftStoreIops,omitempty"`
	// The display name of the node spec. For example, `8 vCPU, 32 GiB`.
	NodeSpecDisplayName *string `json:"nodeSpecDisplayName,omitempty"`
	// The progress of the node changing.
	NodeChangingProgress *ClusterNodeChangingProgress `json:"nodeChangingProgress,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1ClusterStorageNodeSetting V1beta1ClusterStorageNodeSetting

// NewV1beta1ClusterStorageNodeSetting instantiates a new V1beta1ClusterStorageNodeSetting object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1ClusterStorageNodeSetting(nodeCount int32, nodeSpecKey string, storageSizeGi int32) *V1beta1ClusterStorageNodeSetting {
	this := V1beta1ClusterStorageNodeSetting{}
	this.NodeCount = nodeCount
	this.NodeSpecKey = nodeSpecKey
	this.StorageSizeGi = storageSizeGi
	return &this
}

// NewV1beta1ClusterStorageNodeSettingWithDefaults instantiates a new V1beta1ClusterStorageNodeSetting object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1ClusterStorageNodeSettingWithDefaults() *V1beta1ClusterStorageNodeSetting {
	this := V1beta1ClusterStorageNodeSetting{}
	return &this
}

// GetNodeCount returns the NodeCount field value
func (o *V1beta1ClusterStorageNodeSetting) GetNodeCount() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.NodeCount
}

// GetNodeCountOk returns a tuple with the NodeCount field value
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeCountOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeCount, true
}

// SetNodeCount sets field value
func (o *V1beta1ClusterStorageNodeSetting) SetNodeCount(v int32) {
	o.NodeCount = v
}

// GetNodeSpecKey returns the NodeSpecKey field value
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecKey() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.NodeSpecKey
}

// GetNodeSpecKeyOk returns a tuple with the NodeSpecKey field value
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecKeyOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.NodeSpecKey, true
}

// SetNodeSpecKey sets field value
func (o *V1beta1ClusterStorageNodeSetting) SetNodeSpecKey(v string) {
	o.NodeSpecKey = v
}

// GetNodeSpecVersion returns the NodeSpecVersion field value if set, zero value otherwise.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecVersion() string {
	if o == nil || IsNil(o.NodeSpecVersion) {
		var ret string
		return ret
	}
	return *o.NodeSpecVersion
}

// GetNodeSpecVersionOk returns a tuple with the NodeSpecVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecVersionOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecVersion) {
		return nil, false
	}
	return o.NodeSpecVersion, true
}

// HasNodeSpecVersion returns a boolean if a field has been set.
func (o *V1beta1ClusterStorageNodeSetting) HasNodeSpecVersion() bool {
	if o != nil && !IsNil(o.NodeSpecVersion) {
		return true
	}

	return false
}

// SetNodeSpecVersion gets a reference to the given string and assigns it to the NodeSpecVersion field.
func (o *V1beta1ClusterStorageNodeSetting) SetNodeSpecVersion(v string) {
	o.NodeSpecVersion = &v
}

// GetStorageSizeGi returns the StorageSizeGi field value
func (o *V1beta1ClusterStorageNodeSetting) GetStorageSizeGi() int32 {
	if o == nil {
		var ret int32
		return ret
	}

	return o.StorageSizeGi
}

// GetStorageSizeGiOk returns a tuple with the StorageSizeGi field value
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetStorageSizeGiOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StorageSizeGi, true
}

// SetStorageSizeGi sets field value
func (o *V1beta1ClusterStorageNodeSetting) SetStorageSizeGi(v int32) {
	o.StorageSizeGi = v
}

// GetStorageType returns the StorageType field value if set, zero value otherwise.
func (o *V1beta1ClusterStorageNodeSetting) GetStorageType() ClusterStorageNodeSettingStorageType {
	if o == nil || IsNil(o.StorageType) {
		var ret ClusterStorageNodeSettingStorageType
		return ret
	}
	return *o.StorageType
}

// GetStorageTypeOk returns a tuple with the StorageType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetStorageTypeOk() (*ClusterStorageNodeSettingStorageType, bool) {
	if o == nil || IsNil(o.StorageType) {
		return nil, false
	}
	return o.StorageType, true
}

// HasStorageType returns a boolean if a field has been set.
func (o *V1beta1ClusterStorageNodeSetting) HasStorageType() bool {
	if o != nil && !IsNil(o.StorageType) {
		return true
	}

	return false
}

// SetStorageType gets a reference to the given ClusterStorageNodeSettingStorageType and assigns it to the StorageType field.
func (o *V1beta1ClusterStorageNodeSetting) SetStorageType(v ClusterStorageNodeSettingStorageType) {
	o.StorageType = &v
}

// GetRaftStoreIops returns the RaftStoreIops field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1ClusterStorageNodeSetting) GetRaftStoreIops() int32 {
	if o == nil || IsNil(o.RaftStoreIops.Get()) {
		var ret int32
		return ret
	}
	return *o.RaftStoreIops.Get()
}

// GetRaftStoreIopsOk returns a tuple with the RaftStoreIops field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1ClusterStorageNodeSetting) GetRaftStoreIopsOk() (*int32, bool) {
	if o == nil {
		return nil, false
	}
	return o.RaftStoreIops.Get(), o.RaftStoreIops.IsSet()
}

// HasRaftStoreIops returns a boolean if a field has been set.
func (o *V1beta1ClusterStorageNodeSetting) HasRaftStoreIops() bool {
	if o != nil && o.RaftStoreIops.IsSet() {
		return true
	}

	return false
}

// SetRaftStoreIops gets a reference to the given NullableInt32 and assigns it to the RaftStoreIops field.
func (o *V1beta1ClusterStorageNodeSetting) SetRaftStoreIops(v int32) {
	o.RaftStoreIops.Set(&v)
}

// SetRaftStoreIopsNil sets the value for RaftStoreIops to be an explicit nil
func (o *V1beta1ClusterStorageNodeSetting) SetRaftStoreIopsNil() {
	o.RaftStoreIops.Set(nil)
}

// UnsetRaftStoreIops ensures that no value is present for RaftStoreIops, not even an explicit nil
func (o *V1beta1ClusterStorageNodeSetting) UnsetRaftStoreIops() {
	o.RaftStoreIops.Unset()
}

// GetNodeSpecDisplayName returns the NodeSpecDisplayName field value if set, zero value otherwise.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecDisplayName() string {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		var ret string
		return ret
	}
	return *o.NodeSpecDisplayName
}

// GetNodeSpecDisplayNameOk returns a tuple with the NodeSpecDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeSpecDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.NodeSpecDisplayName) {
		return nil, false
	}
	return o.NodeSpecDisplayName, true
}

// HasNodeSpecDisplayName returns a boolean if a field has been set.
func (o *V1beta1ClusterStorageNodeSetting) HasNodeSpecDisplayName() bool {
	if o != nil && !IsNil(o.NodeSpecDisplayName) {
		return true
	}

	return false
}

// SetNodeSpecDisplayName gets a reference to the given string and assigns it to the NodeSpecDisplayName field.
func (o *V1beta1ClusterStorageNodeSetting) SetNodeSpecDisplayName(v string) {
	o.NodeSpecDisplayName = &v
}

// GetNodeChangingProgress returns the NodeChangingProgress field value if set, zero value otherwise.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeChangingProgress() ClusterNodeChangingProgress {
	if o == nil || IsNil(o.NodeChangingProgress) {
		var ret ClusterNodeChangingProgress
		return ret
	}
	return *o.NodeChangingProgress
}

// GetNodeChangingProgressOk returns a tuple with the NodeChangingProgress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1ClusterStorageNodeSetting) GetNodeChangingProgressOk() (*ClusterNodeChangingProgress, bool) {
	if o == nil || IsNil(o.NodeChangingProgress) {
		return nil, false
	}
	return o.NodeChangingProgress, true
}

// HasNodeChangingProgress returns a boolean if a field has been set.
func (o *V1beta1ClusterStorageNodeSetting) HasNodeChangingProgress() bool {
	if o != nil && !IsNil(o.NodeChangingProgress) {
		return true
	}

	return false
}

// SetNodeChangingProgress gets a reference to the given ClusterNodeChangingProgress and assigns it to the NodeChangingProgress field.
func (o *V1beta1ClusterStorageNodeSetting) SetNodeChangingProgress(v ClusterNodeChangingProgress) {
	o.NodeChangingProgress = &v
}

func (o V1beta1ClusterStorageNodeSetting) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1ClusterStorageNodeSetting) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["nodeCount"] = o.NodeCount
	toSerialize["nodeSpecKey"] = o.NodeSpecKey
	if !IsNil(o.NodeSpecVersion) {
		toSerialize["nodeSpecVersion"] = o.NodeSpecVersion
	}
	toSerialize["storageSizeGi"] = o.StorageSizeGi
	if !IsNil(o.StorageType) {
		toSerialize["storageType"] = o.StorageType
	}
	if o.RaftStoreIops.IsSet() {
		toSerialize["raftStoreIops"] = o.RaftStoreIops.Get()
	}
	if !IsNil(o.NodeSpecDisplayName) {
		toSerialize["nodeSpecDisplayName"] = o.NodeSpecDisplayName
	}
	if !IsNil(o.NodeChangingProgress) {
		toSerialize["nodeChangingProgress"] = o.NodeChangingProgress
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1ClusterStorageNodeSetting) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"nodeCount",
		"nodeSpecKey",
		"storageSizeGi",
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

	varV1beta1ClusterStorageNodeSetting := _V1beta1ClusterStorageNodeSetting{}

	err = json.Unmarshal(data, &varV1beta1ClusterStorageNodeSetting)

	if err != nil {
		return err
	}

	*o = V1beta1ClusterStorageNodeSetting(varV1beta1ClusterStorageNodeSetting)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "nodeCount")
		delete(additionalProperties, "nodeSpecKey")
		delete(additionalProperties, "nodeSpecVersion")
		delete(additionalProperties, "storageSizeGi")
		delete(additionalProperties, "storageType")
		delete(additionalProperties, "raftStoreIops")
		delete(additionalProperties, "nodeSpecDisplayName")
		delete(additionalProperties, "nodeChangingProgress")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1ClusterStorageNodeSetting struct {
	value *V1beta1ClusterStorageNodeSetting
	isSet bool
}

func (v NullableV1beta1ClusterStorageNodeSetting) Get() *V1beta1ClusterStorageNodeSetting {
	return v.value
}

func (v *NullableV1beta1ClusterStorageNodeSetting) Set(val *V1beta1ClusterStorageNodeSetting) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1ClusterStorageNodeSetting) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1ClusterStorageNodeSetting) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1ClusterStorageNodeSetting(val *V1beta1ClusterStorageNodeSetting) *NullableV1beta1ClusterStorageNodeSetting {
	return &NullableV1beta1ClusterStorageNodeSetting{value: val, isSet: true}
}

func (v NullableV1beta1ClusterStorageNodeSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1ClusterStorageNodeSetting) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
