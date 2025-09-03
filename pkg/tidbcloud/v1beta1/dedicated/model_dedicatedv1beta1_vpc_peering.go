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

// checks if the Dedicatedv1beta1VpcPeering type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Dedicatedv1beta1VpcPeering{}

// Dedicatedv1beta1VpcPeering struct for Dedicatedv1beta1VpcPeering
type Dedicatedv1beta1VpcPeering struct {
	Name         *string `json:"name,omitempty"`
	VpcPeeringId *string `json:"vpcPeeringId,omitempty"`
	// The labels of the vpc peering. It always contains the `project_id` label.
	Labels            *map[string]string `json:"labels,omitempty"`
	TidbCloudRegionId string             `json:"tidbCloudRegionId"`
	// Format: {cloud_provider}-{region_code} For AWS, it's required. For GCP, it's optional. Since GCP does not require region_id when creating VPC peering.
	CustomerRegionId *string `json:"customerRegionId,omitempty"`
	// In AWS, it is the account ID. In GCP, it is the project name.
	CustomerAccountId string `json:"customerAccountId"`
	// In AWS, it is the VPC ID. In GCP, it is the network name.
	CustomerVpcId          string                      `json:"customerVpcId"`
	CustomerVpcCidr        string                      `json:"customerVpcCidr"`
	TidbCloudCloudProvider *V1beta1RegionCloudProvider `json:"tidbCloudCloudProvider,omitempty"`
	// In AWS, it is the account ID. In GCP, it is the project name.
	TidbCloudAccountId *string `json:"tidbCloudAccountId,omitempty"`
	// In AWS, it is the VPC ID. In GCP, it is the network name.
	TidbCloudVpcId   *string                          `json:"tidbCloudVpcId,omitempty"`
	TidbCloudVpcCidr *string                          `json:"tidbCloudVpcCidr,omitempty"`
	State            *Dedicatedv1beta1VpcPeeringState `json:"state,omitempty"`
	// Only for AWS vpc peerings.
	AwsVpcPeeringConnectionId NullableString `json:"awsVpcPeeringConnectionId,omitempty"`
	AdditionalProperties      map[string]interface{}
}

type _Dedicatedv1beta1VpcPeering Dedicatedv1beta1VpcPeering

// NewDedicatedv1beta1VpcPeering instantiates a new Dedicatedv1beta1VpcPeering object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDedicatedv1beta1VpcPeering(tidbCloudRegionId string, customerAccountId string, customerVpcId string, customerVpcCidr string) *Dedicatedv1beta1VpcPeering {
	this := Dedicatedv1beta1VpcPeering{}
	this.TidbCloudRegionId = tidbCloudRegionId
	this.CustomerAccountId = customerAccountId
	this.CustomerVpcId = customerVpcId
	this.CustomerVpcCidr = customerVpcCidr
	return &this
}

// NewDedicatedv1beta1VpcPeeringWithDefaults instantiates a new Dedicatedv1beta1VpcPeering object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDedicatedv1beta1VpcPeeringWithDefaults() *Dedicatedv1beta1VpcPeering {
	this := Dedicatedv1beta1VpcPeering{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Dedicatedv1beta1VpcPeering) SetName(v string) {
	o.Name = &v
}

// GetVpcPeeringId returns the VpcPeeringId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetVpcPeeringId() string {
	if o == nil || IsNil(o.VpcPeeringId) {
		var ret string
		return ret
	}
	return *o.VpcPeeringId
}

// GetVpcPeeringIdOk returns a tuple with the VpcPeeringId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetVpcPeeringIdOk() (*string, bool) {
	if o == nil || IsNil(o.VpcPeeringId) {
		return nil, false
	}
	return o.VpcPeeringId, true
}

// HasVpcPeeringId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasVpcPeeringId() bool {
	if o != nil && !IsNil(o.VpcPeeringId) {
		return true
	}

	return false
}

// SetVpcPeeringId gets a reference to the given string and assigns it to the VpcPeeringId field.
func (o *Dedicatedv1beta1VpcPeering) SetVpcPeeringId(v string) {
	o.VpcPeeringId = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *Dedicatedv1beta1VpcPeering) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetTidbCloudRegionId returns the TidbCloudRegionId field value
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudRegionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.TidbCloudRegionId
}

// GetTidbCloudRegionIdOk returns a tuple with the TidbCloudRegionId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudRegionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.TidbCloudRegionId, true
}

// SetTidbCloudRegionId sets field value
func (o *Dedicatedv1beta1VpcPeering) SetTidbCloudRegionId(v string) {
	o.TidbCloudRegionId = v
}

// GetCustomerRegionId returns the CustomerRegionId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetCustomerRegionId() string {
	if o == nil || IsNil(o.CustomerRegionId) {
		var ret string
		return ret
	}
	return *o.CustomerRegionId
}

// GetCustomerRegionIdOk returns a tuple with the CustomerRegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetCustomerRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.CustomerRegionId) {
		return nil, false
	}
	return o.CustomerRegionId, true
}

// HasCustomerRegionId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasCustomerRegionId() bool {
	if o != nil && !IsNil(o.CustomerRegionId) {
		return true
	}

	return false
}

// SetCustomerRegionId gets a reference to the given string and assigns it to the CustomerRegionId field.
func (o *Dedicatedv1beta1VpcPeering) SetCustomerRegionId(v string) {
	o.CustomerRegionId = &v
}

// GetCustomerAccountId returns the CustomerAccountId field value
func (o *Dedicatedv1beta1VpcPeering) GetCustomerAccountId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerAccountId
}

// GetCustomerAccountIdOk returns a tuple with the CustomerAccountId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetCustomerAccountIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerAccountId, true
}

// SetCustomerAccountId sets field value
func (o *Dedicatedv1beta1VpcPeering) SetCustomerAccountId(v string) {
	o.CustomerAccountId = v
}

// GetCustomerVpcId returns the CustomerVpcId field value
func (o *Dedicatedv1beta1VpcPeering) GetCustomerVpcId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerVpcId
}

// GetCustomerVpcIdOk returns a tuple with the CustomerVpcId field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetCustomerVpcIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerVpcId, true
}

// SetCustomerVpcId sets field value
func (o *Dedicatedv1beta1VpcPeering) SetCustomerVpcId(v string) {
	o.CustomerVpcId = v
}

// GetCustomerVpcCidr returns the CustomerVpcCidr field value
func (o *Dedicatedv1beta1VpcPeering) GetCustomerVpcCidr() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.CustomerVpcCidr
}

// GetCustomerVpcCidrOk returns a tuple with the CustomerVpcCidr field value
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetCustomerVpcCidrOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CustomerVpcCidr, true
}

// SetCustomerVpcCidr sets field value
func (o *Dedicatedv1beta1VpcPeering) SetCustomerVpcCidr(v string) {
	o.CustomerVpcCidr = v
}

// GetTidbCloudCloudProvider returns the TidbCloudCloudProvider field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.TidbCloudCloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.TidbCloudCloudProvider
}

// GetTidbCloudCloudProviderOk returns a tuple with the TidbCloudCloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.TidbCloudCloudProvider) {
		return nil, false
	}
	return o.TidbCloudCloudProvider, true
}

// HasTidbCloudCloudProvider returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasTidbCloudCloudProvider() bool {
	if o != nil && !IsNil(o.TidbCloudCloudProvider) {
		return true
	}

	return false
}

// SetTidbCloudCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the TidbCloudCloudProvider field.
func (o *Dedicatedv1beta1VpcPeering) SetTidbCloudCloudProvider(v V1beta1RegionCloudProvider) {
	o.TidbCloudCloudProvider = &v
}

// GetTidbCloudAccountId returns the TidbCloudAccountId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudAccountId() string {
	if o == nil || IsNil(o.TidbCloudAccountId) {
		var ret string
		return ret
	}
	return *o.TidbCloudAccountId
}

// GetTidbCloudAccountIdOk returns a tuple with the TidbCloudAccountId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudAccountIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbCloudAccountId) {
		return nil, false
	}
	return o.TidbCloudAccountId, true
}

// HasTidbCloudAccountId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasTidbCloudAccountId() bool {
	if o != nil && !IsNil(o.TidbCloudAccountId) {
		return true
	}

	return false
}

// SetTidbCloudAccountId gets a reference to the given string and assigns it to the TidbCloudAccountId field.
func (o *Dedicatedv1beta1VpcPeering) SetTidbCloudAccountId(v string) {
	o.TidbCloudAccountId = &v
}

// GetTidbCloudVpcId returns the TidbCloudVpcId field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudVpcId() string {
	if o == nil || IsNil(o.TidbCloudVpcId) {
		var ret string
		return ret
	}
	return *o.TidbCloudVpcId
}

// GetTidbCloudVpcIdOk returns a tuple with the TidbCloudVpcId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudVpcIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbCloudVpcId) {
		return nil, false
	}
	return o.TidbCloudVpcId, true
}

// HasTidbCloudVpcId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasTidbCloudVpcId() bool {
	if o != nil && !IsNil(o.TidbCloudVpcId) {
		return true
	}

	return false
}

// SetTidbCloudVpcId gets a reference to the given string and assigns it to the TidbCloudVpcId field.
func (o *Dedicatedv1beta1VpcPeering) SetTidbCloudVpcId(v string) {
	o.TidbCloudVpcId = &v
}

// GetTidbCloudVpcCidr returns the TidbCloudVpcCidr field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudVpcCidr() string {
	if o == nil || IsNil(o.TidbCloudVpcCidr) {
		var ret string
		return ret
	}
	return *o.TidbCloudVpcCidr
}

// GetTidbCloudVpcCidrOk returns a tuple with the TidbCloudVpcCidr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetTidbCloudVpcCidrOk() (*string, bool) {
	if o == nil || IsNil(o.TidbCloudVpcCidr) {
		return nil, false
	}
	return o.TidbCloudVpcCidr, true
}

// HasTidbCloudVpcCidr returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasTidbCloudVpcCidr() bool {
	if o != nil && !IsNil(o.TidbCloudVpcCidr) {
		return true
	}

	return false
}

// SetTidbCloudVpcCidr gets a reference to the given string and assigns it to the TidbCloudVpcCidr field.
func (o *Dedicatedv1beta1VpcPeering) SetTidbCloudVpcCidr(v string) {
	o.TidbCloudVpcCidr = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Dedicatedv1beta1VpcPeering) GetState() Dedicatedv1beta1VpcPeeringState {
	if o == nil || IsNil(o.State) {
		var ret Dedicatedv1beta1VpcPeeringState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Dedicatedv1beta1VpcPeering) GetStateOk() (*Dedicatedv1beta1VpcPeeringState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Dedicatedv1beta1VpcPeeringState and assigns it to the State field.
func (o *Dedicatedv1beta1VpcPeering) SetState(v Dedicatedv1beta1VpcPeeringState) {
	o.State = &v
}

// GetAwsVpcPeeringConnectionId returns the AwsVpcPeeringConnectionId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Dedicatedv1beta1VpcPeering) GetAwsVpcPeeringConnectionId() string {
	if o == nil || IsNil(o.AwsVpcPeeringConnectionId.Get()) {
		var ret string
		return ret
	}
	return *o.AwsVpcPeeringConnectionId.Get()
}

// GetAwsVpcPeeringConnectionIdOk returns a tuple with the AwsVpcPeeringConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Dedicatedv1beta1VpcPeering) GetAwsVpcPeeringConnectionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AwsVpcPeeringConnectionId.Get(), o.AwsVpcPeeringConnectionId.IsSet()
}

// HasAwsVpcPeeringConnectionId returns a boolean if a field has been set.
func (o *Dedicatedv1beta1VpcPeering) HasAwsVpcPeeringConnectionId() bool {
	if o != nil && o.AwsVpcPeeringConnectionId.IsSet() {
		return true
	}

	return false
}

// SetAwsVpcPeeringConnectionId gets a reference to the given NullableString and assigns it to the AwsVpcPeeringConnectionId field.
func (o *Dedicatedv1beta1VpcPeering) SetAwsVpcPeeringConnectionId(v string) {
	o.AwsVpcPeeringConnectionId.Set(&v)
}

// SetAwsVpcPeeringConnectionIdNil sets the value for AwsVpcPeeringConnectionId to be an explicit nil
func (o *Dedicatedv1beta1VpcPeering) SetAwsVpcPeeringConnectionIdNil() {
	o.AwsVpcPeeringConnectionId.Set(nil)
}

// UnsetAwsVpcPeeringConnectionId ensures that no value is present for AwsVpcPeeringConnectionId, not even an explicit nil
func (o *Dedicatedv1beta1VpcPeering) UnsetAwsVpcPeeringConnectionId() {
	o.AwsVpcPeeringConnectionId.Unset()
}

func (o Dedicatedv1beta1VpcPeering) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Dedicatedv1beta1VpcPeering) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.VpcPeeringId) {
		toSerialize["vpcPeeringId"] = o.VpcPeeringId
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	toSerialize["tidbCloudRegionId"] = o.TidbCloudRegionId
	if !IsNil(o.CustomerRegionId) {
		toSerialize["customerRegionId"] = o.CustomerRegionId
	}
	toSerialize["customerAccountId"] = o.CustomerAccountId
	toSerialize["customerVpcId"] = o.CustomerVpcId
	toSerialize["customerVpcCidr"] = o.CustomerVpcCidr
	if !IsNil(o.TidbCloudCloudProvider) {
		toSerialize["tidbCloudCloudProvider"] = o.TidbCloudCloudProvider
	}
	if !IsNil(o.TidbCloudAccountId) {
		toSerialize["tidbCloudAccountId"] = o.TidbCloudAccountId
	}
	if !IsNil(o.TidbCloudVpcId) {
		toSerialize["tidbCloudVpcId"] = o.TidbCloudVpcId
	}
	if !IsNil(o.TidbCloudVpcCidr) {
		toSerialize["tidbCloudVpcCidr"] = o.TidbCloudVpcCidr
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if o.AwsVpcPeeringConnectionId.IsSet() {
		toSerialize["awsVpcPeeringConnectionId"] = o.AwsVpcPeeringConnectionId.Get()
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Dedicatedv1beta1VpcPeering) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"tidbCloudRegionId",
		"customerAccountId",
		"customerVpcId",
		"customerVpcCidr",
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

	varDedicatedv1beta1VpcPeering := _Dedicatedv1beta1VpcPeering{}

	err = json.Unmarshal(data, &varDedicatedv1beta1VpcPeering)

	if err != nil {
		return err
	}

	*o = Dedicatedv1beta1VpcPeering(varDedicatedv1beta1VpcPeering)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "vpcPeeringId")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "tidbCloudRegionId")
		delete(additionalProperties, "customerRegionId")
		delete(additionalProperties, "customerAccountId")
		delete(additionalProperties, "customerVpcId")
		delete(additionalProperties, "customerVpcCidr")
		delete(additionalProperties, "tidbCloudCloudProvider")
		delete(additionalProperties, "tidbCloudAccountId")
		delete(additionalProperties, "tidbCloudVpcId")
		delete(additionalProperties, "tidbCloudVpcCidr")
		delete(additionalProperties, "state")
		delete(additionalProperties, "awsVpcPeeringConnectionId")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDedicatedv1beta1VpcPeering struct {
	value *Dedicatedv1beta1VpcPeering
	isSet bool
}

func (v NullableDedicatedv1beta1VpcPeering) Get() *Dedicatedv1beta1VpcPeering {
	return v.value
}

func (v *NullableDedicatedv1beta1VpcPeering) Set(val *Dedicatedv1beta1VpcPeering) {
	v.value = val
	v.isSet = true
}

func (v NullableDedicatedv1beta1VpcPeering) IsSet() bool {
	return v.isSet
}

func (v *NullableDedicatedv1beta1VpcPeering) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDedicatedv1beta1VpcPeering(val *Dedicatedv1beta1VpcPeering) *NullableDedicatedv1beta1VpcPeering {
	return &NullableDedicatedv1beta1VpcPeering{value: val, isSet: true}
}

func (v NullableDedicatedv1beta1VpcPeering) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDedicatedv1beta1VpcPeering) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
