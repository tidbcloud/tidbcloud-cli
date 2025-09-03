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

// checks if the Dedicatedv1beta1PrivateLinkService type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1PrivateLinkService{}

// Dedicatedv1beta1PrivateLinkService All fields are output only.
type Dedicatedv1beta1PrivateLinkService struct {
	// The name of the private link service.
	Name *string `json:"name,omitempty"`
	// The ID of the TiDB group to which the private link service belongs.
	TidbNodeGroupId *string `json:"tidbNodeGroupId,omitempty"`
	// The service name of the private link service, which varies by cloud provider:  - AWS: the service name of the private link service.  - Google Cloud: the resource name of the service attachment.  - Azure: the resource ID of the private link service.
	ServiceName *string `json:"serviceName,omitempty"`
	// The DNS name of the private link service, which varies by cloud provider:  - AWS: the fully qualified domain name (FQDN) shared across all private endpoints, regardless of VPC location.  - Google Cloud: the zone name (suffix of the FQDN) shared across all private endpoints in a single VPC network. The format of FQDN is `<endpoint_name>.<service_dns_name>`.  - Azure: the zone name shared across public internet. The format of FQDN is `<endpoint_name>-<random_hash>.<service_dns_name>`.
	ServiceDnsName *string `json:"serviceDnsName,omitempty"`
	// (AWS only) The availability zones where the private link service is available. For more information, see [`DescribeVpcEndpointServices`](https://docs.aws.amazon.com/AWSEC2/latest/APIReference/API_DescribeVpcEndpointServices.html).
	AvailableZones []string `json:"availableZones,omitempty"`
	// The state of the private link service.  - `\"CREATING\"`: the private link service is being created.  - `\"ACTIVE\"`: the private link service is ready to use.  - `\"DELETING\"`: the private link service is being deleted.
	State *Dedicatedv1beta1PrivateLinkServiceState `json:"state,omitempty"`
	// The ID of the region where the private link service is located, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-east-1`.
	RegionId *string `json:"regionId,omitempty"`
	// The display name of the region where the private link service is located. For example, `N. Virginia (us-east-1)`.
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// The cloud provider where the private link service is located.  - `\"aws\"`: Amazon Web Services  - `\"gcp\"`: Google Cloud  - `\"azure\"`: Microsoft Azure  - `\"alicloud\"`: Alibaba Cloud
	CloudProvider        *V1beta1RegionCloudProvider `json:"cloudProvider,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Dedicatedv1beta1PrivateLinkService Dedicatedv1beta1PrivateLinkService

// NewDedicatedv1beta1PrivateLinkService instantiates a new Dedicatedv1beta1PrivateLinkService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1PrivateLinkService() *Dedicatedv1beta1PrivateLinkService {
	this := Dedicatedv1beta1PrivateLinkService{}
	return &this
}

// NewDedicatedv1beta1PrivateLinkServiceWithDefaults instantiates a new Dedicatedv1beta1PrivateLinkService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1PrivateLinkServiceWithDefaults() *Dedicatedv1beta1PrivateLinkService {
	this := Dedicatedv1beta1PrivateLinkService{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1PrivateLinkService) SetName(v string) {
	o.Name = &v
}

// GetTidbNodeGroupId returns the TidbNodeGroupId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetTidbNodeGroupId() string {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupId
}

// GetTidbNodeGroupIdOk returns a tuple with the TidbNodeGroupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetTidbNodeGroupIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupId) {
		return nil, false
	}
	return o.TidbNodeGroupId, true
}

// HasTidbNodeGroupId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasTidbNodeGroupId() bool {
	if o != nil && !IsNil(o.TidbNodeGroupId) {
		return true
	}

	return false
}

// SetTidbNodeGroupId gets a reference to the given string and assigns it to the TidbNodeGroupId field.
func (o *Dedicatedv1beta1PrivateLinkService) SetTidbNodeGroupId(v string) {
	o.TidbNodeGroupId = &v
}

// GetServiceName returns the ServiceName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceName() string {
	if o == nil || IsNil(o.ServiceName) {
		var ret string
		return ret
	}
	return *o.ServiceName
}

// GetServiceNameOk returns a tuple with the ServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceName) {
		return nil, false
	}
	return o.ServiceName, true
}

// HasServiceName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasServiceName() bool {
	if o != nil && !IsNil(o.ServiceName) {
		return true
	}

	return false
}

// SetServiceName gets a reference to the given string and assigns it to the ServiceName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetServiceName(v string) {
	o.ServiceName = &v
}

// GetServiceDnsName returns the ServiceDnsName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceDnsName() string {
	if o == nil || IsNil(o.ServiceDnsName) {
		var ret string
		return ret
	}
	return *o.ServiceDnsName
}

// GetServiceDnsNameOk returns a tuple with the ServiceDnsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetServiceDnsNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceDnsName) {
		return nil, false
	}
	return o.ServiceDnsName, true
}

// HasServiceDnsName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasServiceDnsName() bool {
	if o != nil && !IsNil(o.ServiceDnsName) {
		return true
	}

	return false
}

// SetServiceDnsName gets a reference to the given string and assigns it to the ServiceDnsName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetServiceDnsName(v string) {
	o.ServiceDnsName = &v
}

// GetAvailableZones returns the AvailableZones field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetAvailableZones() []string {
	if o == nil || IsNil(o.AvailableZones) {
		var ret []string
		return ret
	}
	return o.AvailableZones
}

// GetAvailableZonesOk returns a tuple with the AvailableZones field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetAvailableZonesOk() ([]string, bool) {
	if o == nil || IsNil(o.AvailableZones) {
		return nil, false
	}
	return o.AvailableZones, true
}

// HasAvailableZones returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasAvailableZones() bool {
	if o != nil && !IsNil(o.AvailableZones) {
		return true
	}

	return false
}

// SetAvailableZones gets a reference to the given []string and assigns it to the AvailableZones field.
func (o *Dedicatedv1beta1PrivateLinkService) SetAvailableZones(v []string) {
	o.AvailableZones = v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetState() Dedicatedv1beta1PrivateLinkServiceState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1PrivateLinkServiceState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetStateOk() (*Dedicatedv1beta1PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1PrivateLinkServiceState and assigns it to the State field.
func (o *Dedicatedv1beta1PrivateLinkService) SetState(v Dedicatedv1beta1PrivateLinkServiceState) {
	o.State = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *Dedicatedv1beta1PrivateLinkService) SetRegionId(v string) {
	o.RegionId = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *Dedicatedv1beta1PrivateLinkService) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Dedicatedv1beta1PrivateLinkService) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1PrivateLinkService) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Dedicatedv1beta1PrivateLinkService) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Dedicatedv1beta1PrivateLinkService) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

func (o Dedicatedv1beta1PrivateLinkService) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1PrivateLinkService) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbNodeGroupId) {
		toSerialize["tidbNodeGroupId"] = o.TidbNodeGroupId
	}
	if !IsNil(o.ServiceName) {
		toSerialize["serviceName"] = o.ServiceName
	}
	if !IsNil(o.ServiceDnsName) {
		toSerialize["serviceDnsName"] = o.ServiceDnsName
	}
	if !IsNil(o.AvailableZones) {
		toSerialize["availableZones"] = o.AvailableZones
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1PrivateLinkService) UnmarshalJSON(data []byte) (err error) {
	varDedicatedv1beta1PrivateLinkService := _Dedicatedv1beta1PrivateLinkService{}

	err = json.Unmarshal(data, &varDedicatedv1beta1PrivateLinkService)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1PrivateLinkService(varDedicatedv1beta1PrivateLinkService)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbNodeGroupId")
		delete(additionalProperties, "serviceName")
		delete(additionalProperties, "serviceDnsName")
		delete(additionalProperties, "availableZones")
		delete(additionalProperties, "state")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1PrivateLinkService struct {
	value *Dedicatedv1beta1PrivateLinkService
	isSet bool
}

func (v NullableDedicatedv1beta1PrivateLinkService) Get() *Dedicatedv1beta1PrivateLinkService {
	return v.value
}

func (v *NullableDedicatedv1beta1PrivateLinkService) Set(val *Dedicatedv1beta1PrivateLinkService) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1PrivateLinkService) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1PrivateLinkService) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1PrivateLinkService(val *Dedicatedv1beta1PrivateLinkService) *NullableDedicatedv1beta1PrivateLinkService {
	return &NullableDedicatedv1beta1PrivateLinkService{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1PrivateLinkService) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1PrivateLinkService) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
