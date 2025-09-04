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

// checks if the PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}

// PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct for PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
type PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct {
	// The name of the private endpoint connection, in the format of `tidbNodeGroups/{tidb_node_group_id}/privateEndpointConnections/{private_endpoint_connection_id}`.
	Name *string `json:"name,omitempty"`
	// The unique ID of the private endpoint connection.
	PrivateEndpointConnectionId *string `json:"privateEndpointConnectionId,omitempty"`
	// The ID of the cluster to which the private endpoint connection belongs.
	ClusterId *string `json:"clusterId,omitempty"`
	// The display name of the cluster to which the private endpoint connection belongs.
	ClusterDisplayName *string `json:"clusterDisplayName,omitempty"`
	// The labels of the private link connection, including the mandatory `tidb.cloud/project` label identifying the project it belongs to.
	Labels *map[string]string `json:"labels,omitempty"`
	// The endpoint ID of the private link connection.  - AWS: the VPC endpoint ID for [AWS PrivateLink ](https://aws.amazon.com/privatelink/).  - Google Cloud: the endpoint ID for [Private Service Connect](https://cloud.google.com/vpc/docs/private-service-connect).  - Azure: the resource ID for [Azure Private Link](https://learn.microsoft.com/en-us/azure/private-link/private-link-overview).
	EndpointId string `json:"endpointId"`
	// (Azure only) The private IP address of the private endpoint in your virtual network. TiDB Cloud automatically creates a public DNS record that resolves to this IP address, enabling you to connect using the DNS name.
	PrivateIpAddress NullableString `json:"privateIpAddress,omitempty"`
	// The state of the private endpoint connection.  - `\"PENDING\"`: TiDB Cloud is asynchronously accepting the endpoint connection.  - `\"ACTIVE\"`: the private endpoint connection is ready to use.  - `\"DELETING\"`: the private endpoint connection is being deleted.  - `\"FAILED\"`: the private endpoint connection has failed. - `\"DISCOVERED\"`: the endpoint is created in your VPC but not registered with TiDB Cloud.
	EndpointState *Dedicatedv1beta1PrivateEndpointConnectionEndpointState `json:"endpointState,omitempty"`
	// The detailed message when the `endpointState` field is `\"FAILED\"`.
	Message *string `json:"message,omitempty"`
	// The ID of the region where the private endpoint connection is located, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-east-1`.
	RegionId *string `json:"regionId,omitempty"`
	// The display name of the region where the private endpoint connection is located. For example, `N. Virginia (us-east-1)`.
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// The cloud provider where the private endpoint connection is located.  - `\"aws\"`: Amazon Web Services  - `\"gcp\"`: Google Cloud  - `\"azure\"`: Microsoft Azure  - `\"alicloud\"`: Alibaba Cloud
	CloudProvider *V1beta1RegionCloudProvider `json:"cloudProvider,omitempty"`
	// The name of the private link service that the private endpoint connection is connected to.
	PrivateLinkServiceName *string `json:"privateLinkServiceName,omitempty"`
	// The state of the private link service that the private endpoint connection is connected to.
	PrivateLinkServiceState *Dedicatedv1beta1PrivateLinkServiceState `json:"privateLinkServiceState,omitempty"`
	// The display name of the TiDB node group that the private endpoint connection is connected to.
	TidbNodeGroupDisplayName *string `json:"tidbNodeGroupDisplayName,omitempty"`
	// (Google Cloud only) The project name used to identify the Google Cloud project that the private service connection belongs to.
	AccountId NullableString `json:"accountId,omitempty"`
	// The hostname for accessing the TiDB cluster through the private endpoint connection.
	Host *string `json:"host,omitempty"`
	// The port used to connect to the TiDB cluster through the private endpoint connection.
	Port                 *int32 `json:"port,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest

// NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest instantiates a new PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(endpointId string) *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	this := PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}
	this.EndpointId = endpointId
	return &this
}

// NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequestWithDefaults instantiates a new PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequestWithDefaults() *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	this := PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetName(v string) {
	o.Name = &v
}

// GetPrivateEndpointConnectionId returns the PrivateEndpointConnectionId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateEndpointConnectionId() string {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointConnectionId
}

// GetPrivateEndpointConnectionIdOk returns a tuple with the PrivateEndpointConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateEndpointConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		return nil, false
	}
	return o.PrivateEndpointConnectionId, true
}

// HasPrivateEndpointConnectionId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateEndpointConnectionId() bool {
	if o != nil && !IsNil(o.PrivateEndpointConnectionId) {
		return true
	}

	return false
}

// SetPrivateEndpointConnectionId gets a reference to the given string and assigns it to the PrivateEndpointConnectionId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateEndpointConnectionId(v string) {
	o.PrivateEndpointConnectionId = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetClusterDisplayName returns the ClusterDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterDisplayName() string {
	if o == nil || IsNil(o.ClusterDisplayName) {
		var ret string
		return ret
	}
	return *o.ClusterDisplayName
}

// GetClusterDisplayNameOk returns a tuple with the ClusterDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetClusterDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterDisplayName) {
		return nil, false
	}
	return o.ClusterDisplayName, true
}

// HasClusterDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasClusterDisplayName() bool {
	if o != nil && !IsNil(o.ClusterDisplayName) {
		return true
	}

	return false
}

// SetClusterDisplayName gets a reference to the given string and assigns it to the ClusterDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetClusterDisplayName(v string) {
	o.ClusterDisplayName = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetEndpointId returns the EndpointId field value
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndpointId, true
}

// SetEndpointId sets field value
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetEndpointId(v string) {
	o.EndpointId = v
}

// GetPrivateIpAddress returns the PrivateIpAddress field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateIpAddress() string {
	if o == nil || IsNil(o.PrivateIpAddress.Get()) {
		var ret string
		return ret
	}
	return *o.PrivateIpAddress.Get()
}

// GetPrivateIpAddressOk returns a tuple with the PrivateIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateIpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.PrivateIpAddress.Get(), o.PrivateIpAddress.IsSet()
}

// HasPrivateIpAddress returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateIpAddress() bool {
	if o != nil && o.PrivateIpAddress.IsSet() {
		return true
	}

	return false
}

// SetPrivateIpAddress gets a reference to the given NullableString and assigns it to the PrivateIpAddress field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateIpAddress(v string) {
	o.PrivateIpAddress.Set(&v)
}

// SetPrivateIpAddressNil sets the value for PrivateIpAddress to be an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateIpAddressNil() {
	o.PrivateIpAddress.Set(nil)
}

// UnsetPrivateIpAddress ensures that no value is present for PrivateIpAddress, not even an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnsetPrivateIpAddress() {
	o.PrivateIpAddress.Unset()
}

// GetEndpointState returns the EndpointState field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointState() Dedicatedv1beta1PrivateEndpointConnectionEndpointState {
	if o == nil || IsNil(o.EndpointState) {
		var ret Dedicatedv1beta1PrivateEndpointConnectionEndpointState
		return ret
	}
	return *o.EndpointState
}

// GetEndpointStateOk returns a tuple with the EndpointState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetEndpointStateOk() (*Dedicatedv1beta1PrivateEndpointConnectionEndpointState, bool) {
	if o == nil || IsNil(o.EndpointState) {
		return nil, false
	}
	return o.EndpointState, true
}

// HasEndpointState returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasEndpointState() bool {
	if o != nil && !IsNil(o.EndpointState) {
		return true
	}

	return false
}

// SetEndpointState gets a reference to the given Dedicatedv1beta1PrivateEndpointConnectionEndpointState and assigns it to the EndpointState field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetEndpointState(v Dedicatedv1beta1PrivateEndpointConnectionEndpointState) {
	o.EndpointState = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetMessage(v string) {
	o.Message = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetRegionId(v string) {
	o.RegionId = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetPrivateLinkServiceName returns the PrivateLinkServiceName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceName() string {
	if o == nil || IsNil(o.PrivateLinkServiceName) {
		var ret string
		return ret
	}
	return *o.PrivateLinkServiceName
}

// GetPrivateLinkServiceNameOk returns a tuple with the PrivateLinkServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceName) {
		return nil, false
	}
	return o.PrivateLinkServiceName, true
}

// HasPrivateLinkServiceName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateLinkServiceName() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceName) {
		return true
	}

	return false
}

// SetPrivateLinkServiceName gets a reference to the given string and assigns it to the PrivateLinkServiceName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateLinkServiceName(v string) {
	o.PrivateLinkServiceName = &v
}

// GetPrivateLinkServiceState returns the PrivateLinkServiceState field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceState() Dedicatedv1beta1PrivateLinkServiceState {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		var ret Dedicatedv1beta1PrivateLinkServiceState
		return ret
	}
	return *o.PrivateLinkServiceState
}

// GetPrivateLinkServiceStateOk returns a tuple with the PrivateLinkServiceState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPrivateLinkServiceStateOk() (*Dedicatedv1beta1PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		return nil, false
	}
	return o.PrivateLinkServiceState, true
}

// HasPrivateLinkServiceState returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPrivateLinkServiceState() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceState) {
		return true
	}

	return false
}

// SetPrivateLinkServiceState gets a reference to the given Dedicatedv1beta1PrivateLinkServiceState and assigns it to the PrivateLinkServiceState field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPrivateLinkServiceState(v Dedicatedv1beta1PrivateLinkServiceState) {
	o.PrivateLinkServiceState = &v
}

// GetTidbNodeGroupDisplayName returns the TidbNodeGroupDisplayName field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetTidbNodeGroupDisplayName() string {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName) {
		var ret string
		return ret
	}
	return *o.TidbNodeGroupDisplayName
}

// GetTidbNodeGroupDisplayNameOk returns a tuple with the TidbNodeGroupDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetTidbNodeGroupDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.TidbNodeGroupDisplayName) {
		return nil, false
	}
	return o.TidbNodeGroupDisplayName, true
}

// HasTidbNodeGroupDisplayName returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasTidbNodeGroupDisplayName() bool {
	if o != nil && !IsNil(o.TidbNodeGroupDisplayName) {
		return true
	}

	return false
}

// SetTidbNodeGroupDisplayName gets a reference to the given string and assigns it to the TidbNodeGroupDisplayName field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetTidbNodeGroupDisplayName(v string) {
	o.TidbNodeGroupDisplayName = &v
}

// GetAccountId returns the AccountId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetAccountId() string {
	if o == nil || IsNil(o.AccountId.Get()) {
		var ret string
		return ret
	}
	return *o.AccountId.Get()
}

// GetAccountIdOk returns a tuple with the AccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AccountId.Get(), o.AccountId.IsSet()
}

// HasAccountId returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasAccountId() bool {
	if o != nil && o.AccountId.IsSet() {
		return true
	}

	return false
}

// SetAccountId gets a reference to the given NullableString and assigns it to the AccountId field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetAccountId(v string) {
	o.AccountId.Set(&v)
}

// SetAccountIdNil sets the value for AccountId to be an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetAccountIdNil() {
	o.AccountId.Set(nil)
}

// UnsetAccountId ensures that no value is present for AccountId, not even an explicit nil
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnsetAccountId() {
	o.AccountId.Unset()
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) SetPort(v int32) {
	o.Port = &v
}

func (o PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.PrivateEndpointConnectionId) {
		toSerialize["privateEndpointConnectionId"] = o.PrivateEndpointConnectionId
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	if !IsNil(o.ClusterDisplayName) {
		toSerialize["clusterDisplayName"] = o.ClusterDisplayName
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["endpointId"] = o.EndpointId
	if o.PrivateIpAddress.IsSet() {
		toSerialize["privateIpAddress"] = o.PrivateIpAddress.Get()
	}
	if !IsNil(o.EndpointState) {
		toSerialize["endpointState"] = o.EndpointState
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
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
	if !IsNil(o.PrivateLinkServiceName) {
		toSerialize["privateLinkServiceName"] = o.PrivateLinkServiceName
	}
	if !IsNil(o.PrivateLinkServiceState) {
		toSerialize["privateLinkServiceState"] = o.PrivateLinkServiceState
	}
	if !IsNil(o.TidbNodeGroupDisplayName) {
		toSerialize["tidbNodeGroupDisplayName"] = o.TidbNodeGroupDisplayName
	}
	if o.AccountId.IsSet() {
		toSerialize["accountId"] = o.AccountId.Get()
	}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"endpointId",
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

	varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest := _PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{}

	err = json.Unmarshal(data, &varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest)

	if err != nil {
		return err
	}

	*o = PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(varPrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "privateEndpointConnectionId")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "clusterDisplayName")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "endpointId")
		delete(additionalProperties, "privateIpAddress")
		delete(additionalProperties, "endpointState")
		delete(additionalProperties, "message")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "privateLinkServiceName")
		delete(additionalProperties, "privateLinkServiceState")
		delete(additionalProperties, "tidbNodeGroupDisplayName")
		delete(additionalProperties, "accountId")
		delete(additionalProperties, "host")
		delete(additionalProperties, "port")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest struct {
	value *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest
	isSet bool
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Get() *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	return v.value
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Set(val *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) {
	v.value = val
	v.isSet = true
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) IsSet() bool {
	return v.isSet
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest(val *PrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest {
	return &NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest{value: val, isSet: true}
}

func (v NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePrivateEndpointConnectionServiceCreatePrivateEndpointConnectionRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
