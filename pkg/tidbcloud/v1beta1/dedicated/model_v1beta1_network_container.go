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

// checks if the V1beta1NetworkContainer type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1NetworkContainer{}

// V1beta1NetworkContainer struct for V1beta1NetworkContainer
type V1beta1NetworkContainer struct {
	Name               *string `json:"name,omitempty"`
	NetworkContainerId *string `json:"networkContainerId,omitempty"`
	// The labels of the cluster. If there is no \"project_id\" in labels, resource should be in the default project of the creator's organization.
	Labels   *map[string]string `json:"labels,omitempty"`
	RegionId string             `json:"regionId"`
	// If not set, the default cidr of the region will be used.
	CidrNotation      *string                       `json:"cidrNotation,omitempty"`
	CloudProvider     *V1beta1RegionCloudProvider   `json:"cloudProvider,omitempty"`
	State             *V1beta1NetworkContainerState `json:"state,omitempty"`
	RegionDisplayName *string                       `json:"regionDisplayName,omitempty"`
	// For AWS, it is the vpc id. For GCP, it is the network name. For Azure, it is the vnet name.
	VpcId                *string `json:"vpcId,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1NetworkContainer V1beta1NetworkContainer

// NewV1beta1NetworkContainer instantiates a new V1beta1NetworkContainer object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1NetworkContainer(regionId string) *V1beta1NetworkContainer {
	this := V1beta1NetworkContainer{}
	this.RegionId = regionId
	return &this
}

// NewV1beta1NetworkContainerWithDefaults instantiates a new V1beta1NetworkContainer object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1NetworkContainerWithDefaults() *V1beta1NetworkContainer {
	this := V1beta1NetworkContainer{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1NetworkContainer) SetName(v string) {
	o.Name = &v
}

// GetNetworkContainerId returns the NetworkContainerId field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetNetworkContainerId() string {
	if o == nil || IsNil(o.NetworkContainerId) {
		var ret string
		return ret
	}
	return *o.NetworkContainerId
}

// GetNetworkContainerIdOk returns a tuple with the NetworkContainerId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetNetworkContainerIdOk() (*string, bool) {
	if o == nil || IsNil(o.NetworkContainerId) {
		return nil, false
	}
	return o.NetworkContainerId, true
}

// HasNetworkContainerId returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasNetworkContainerId() bool {
	if o != nil && !IsNil(o.NetworkContainerId) {
		return true
	}

	return false
}

// SetNetworkContainerId gets a reference to the given string and assigns it to the NetworkContainerId field.
func (o *V1beta1NetworkContainer) SetNetworkContainerId(v string) {
	o.NetworkContainerId = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *V1beta1NetworkContainer) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetRegionId returns the RegionId field value
func (o *V1beta1NetworkContainer) GetRegionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetRegionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionId, true
}

// SetRegionId sets field value
func (o *V1beta1NetworkContainer) SetRegionId(v string) {
	o.RegionId = v
}

// GetCidrNotation returns the CidrNotation field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetCidrNotation() string {
	if o == nil || IsNil(o.CidrNotation) {
		var ret string
		return ret
	}
	return *o.CidrNotation
}

// GetCidrNotationOk returns a tuple with the CidrNotation field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetCidrNotationOk() (*string, bool) {
	if o == nil || IsNil(o.CidrNotation) {
		return nil, false
	}
	return o.CidrNotation, true
}

// HasCidrNotation returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasCidrNotation() bool {
	if o != nil && !IsNil(o.CidrNotation) {
		return true
	}

	return false
}

// SetCidrNotation gets a reference to the given string and assigns it to the CidrNotation field.
func (o *V1beta1NetworkContainer) SetCidrNotation(v string) {
	o.CidrNotation = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *V1beta1NetworkContainer) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetState() V1beta1NetworkContainerState {
	if o == nil || IsNil(o.State) {
		var ret V1beta1NetworkContainerState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetStateOk() (*V1beta1NetworkContainerState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta1NetworkContainerState and assigns it to the State field.
func (o *V1beta1NetworkContainer) SetState(v V1beta1NetworkContainerState) {
	o.State = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *V1beta1NetworkContainer) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetVpcId returns the VpcId field value if set, zero value otherwise.
func (o *V1beta1NetworkContainer) GetVpcId() string {
	if o == nil || IsNil(o.VpcId) {
		var ret string
		return ret
	}
	return *o.VpcId
}

// GetVpcIdOk returns a tuple with the VpcId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1NetworkContainer) GetVpcIdOk() (*string, bool) {
	if o == nil || IsNil(o.VpcId) {
		return nil, false
	}
	return o.VpcId, true
}

// HasVpcId returns a boolean if a field has been set.
func (o *V1beta1NetworkContainer) HasVpcId() bool {
	if o != nil && !IsNil(o.VpcId) {
		return true
	}

	return false
}

// SetVpcId gets a reference to the given string and assigns it to the VpcId field.
func (o *V1beta1NetworkContainer) SetVpcId(v string) {
	o.VpcId = &v
}

func (o V1beta1NetworkContainer) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1NetworkContainer) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.NetworkContainerId) {
		toSerialize["networkContainerId"] = o.NetworkContainerId
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["regionId"] = o.RegionId
	if !IsNil(o.CidrNotation) {
		toSerialize["cidrNotation"] = o.CidrNotation
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.VpcId) {
		toSerialize["vpcId"] = o.VpcId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1NetworkContainer) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"regionId",
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

	varV1beta1NetworkContainer := _V1beta1NetworkContainer{}

	err = json.Unmarshal(data, &varV1beta1NetworkContainer)

	if err != nil {
		return err
	}

	*o = V1beta1NetworkContainer(varV1beta1NetworkContainer)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "networkContainerId")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "cidrNotation")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "state")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "vpcId")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1NetworkContainer struct {
	value *V1beta1NetworkContainer
	isSet bool
}

func (v NullableV1beta1NetworkContainer) Get() *V1beta1NetworkContainer {
	return v.value
}

func (v *NullableV1beta1NetworkContainer) Set(val *V1beta1NetworkContainer) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1NetworkContainer) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1NetworkContainer) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1NetworkContainer(val *V1beta1NetworkContainer) *NullableV1beta1NetworkContainer {
	return &NullableV1beta1NetworkContainer{value: val, isSet: true}
}

func (v NullableV1beta1NetworkContainer) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1NetworkContainer) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
