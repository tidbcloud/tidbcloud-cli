/*
TiDB Cloud Premium API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium) instances. For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/api/tidb-cloud-api-overview/).  The public connection setting endpoints also support TiDB Cloud Essential V2 instances.  # Overview  The TiDB Cloud Premium API (v1beta2) provides [REST](https://en.wikipedia.org/wiki/REST) endpoints to manage TiDB Cloud Premium instances and related resources.  You can use this API to manage the following resources:  - **TiDB Cloud Premium instance**: manage the lifecycle and configuration of TiDB Cloud Premium instances, including passwords, CA certificates, and cloud provider information. - **Public connection setting**: enable or disable public connections and manage IP access lists for Premium and Essential V2 instances. - **Customer-managed encryption key**: retrieve the IAM principal and verify KMS access before creating a Premium instance on AWS or Alibaba Cloud. - **Backup**: manage backups for TiDB Cloud Premium instances, including backup-based restore. - **Region**: retrieve available regions for deploying TiDB Cloud Premium instances.  # Get Started  This guide helps you make your first API call to the TiDB Cloud Premium API. You will learn how to authenticate a request, build a request, and interpret the response.  1. Create a [TiDB Cloud account](https://tidbcloud.com/signup) if you do not already have one. 2. In the [TiDB Cloud console](https://tidbcloud.com/), go to **Organization** > **API Keys** and create an API key. For more information, see [API key management](#section/Authentication/API-key-management). 3. Make your first API call.   To get all TiDB Cloud Premium instances in your organization, run the following command in your terminal. Replace `YOUR_PUBLIC_KEY` and `YOUR_PRIVATE_KEY` with your own key values.   ```bash  curl --digest \\    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\    --request GET \\    --url 'https://cloud.tidbapi.com/v1beta2/tidbs' \\    --header 'Accept: application/json'  ```  4. The API returns a JSON list of your TiDB Cloud Premium instances. If none exist, the response contains an empty list.  # Authentication  The TiDB Cloud API supports [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) with API keys and OAuth Bearer tokens. OAuth clients send the token in the `Authorization: Bearer <token>` header. It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. 5. Configure the role and scope for the API key. For more information about the permissions of a role, see [User roles](https://docs.pingcap.com/tidbcloud/manage-user-access/#user-roles). 6. Click **Generate API Key**. Copy and save the public key and the private key. 7. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 8. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Update Role**. 4. You can update the description and role of the API key. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # API Changelog  This changelog lists all changes to the TiDB Cloud Premium API (v1beta2).  <!-- In reverse chronological order -->  ## 20260924  - Published CMEK IAM principal retrieval and verification APIs for Premium instances on AWS and Alibaba Cloud.  ## 20260923  - Added APIs to get and update public connection settings for Premium and Essential V2 instances.  ## 20260428  - Initial release of the TiDB Cloud Premium API (v1beta2), including the following resources and endpoints:  * TiDB Cloud Premium instance   * [List TiDB Cloud Premium instances](#tag/TiDB-Instance/operation/TidbService_ListTidbs)   * [Create a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_CreateTidb)   * [Get a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetTidb)   * [Delete a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_DeleteTidb)   * [Update a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_UpdateTidb)   * [Reset the root password of a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_ResetRootPassword)   * [Get cloud provider information for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCloudProviderInfo)   * [Get the CA certificate download URL for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCaCertificateDownloadUrl)  * Backup   * [List backups for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_ListTidbBackups)   * [Delete a backup for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_DeleteTidbBackup)   * [Restore a TiDB Cloud Premium instance from a backup](#tag/Backup/operation/TidbService_RestoreTidb)   * [Get restore status for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_GetRestoreStatus)  * Region   * [List regions](#tag/Region/operation/RegionService_ListRegions)

API version: v1beta2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nextgen

import (
	"encoding/json"
	"fmt"
)

// checks if the Nextgenv1beta2PrivateEndpointConnection type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Nextgenv1beta2PrivateEndpointConnection{}

// Nextgenv1beta2PrivateEndpointConnection struct for Nextgenv1beta2PrivateEndpointConnection
type Nextgenv1beta2PrivateEndpointConnection struct {
	Name                        *string `json:"name,omitempty"`
	Id                          *string `json:"id,omitempty"`
	TidbId                      *string `json:"tidbId,omitempty"`
	PrivateEndpointConnectionId *string `json:"privateEndpointConnectionId,omitempty"`
	// The region where the private endpoint connection is created.
	RegionId          string  `json:"regionId"`
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// The cloud provider for the region.
	CloudProvider *RegionCloudProvider `json:"cloudProvider,omitempty"`
	// The endpoint ID of the private link connection. For AWS, it's VPC endpoint ID. For GCP, it's private service connect endpoint ID. For Azure, it's private endpoint resource ID. For AliCloud, it's endpoint ID.
	EndpointId    string                                                `json:"endpointId"`
	EndpointState *Nextgenv1beta2PrivateEndpointConnectionEndpointState `json:"endpointState,omitempty"`
	Message       *string                                               `json:"message,omitempty"`
	// The private link service name.
	PrivateLinkService *string `json:"privateLinkService,omitempty"`
	// The private link service state.
	PrivateLinkServiceState *Nextgenv1beta2PrivateLinkServiceState `json:"privateLinkServiceState,omitempty"`
	Host                    *string                                `json:"host,omitempty"`
	Port                    *string                                `json:"port,omitempty"`
	// Only for AliCloud. Use private endpoint's Domain name to generate new TiDB DNS name.
	PrivateDomainName *string `json:"privateDomainName,omitempty"`
	// Only for Azure. The private IP address of the private endpoint in the user's vNet.
	PrivateIpAddress     *string `json:"privateIpAddress,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Nextgenv1beta2PrivateEndpointConnection Nextgenv1beta2PrivateEndpointConnection

// NewNextgenv1beta2PrivateEndpointConnection instantiates a new Nextgenv1beta2PrivateEndpointConnection object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNextgenv1beta2PrivateEndpointConnection(regionId string, endpointId string) *Nextgenv1beta2PrivateEndpointConnection {
	this := Nextgenv1beta2PrivateEndpointConnection{}
	this.RegionId = regionId
	this.EndpointId = endpointId
	return &this
}

// NewNextgenv1beta2PrivateEndpointConnectionWithDefaults instantiates a new Nextgenv1beta2PrivateEndpointConnection object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNextgenv1beta2PrivateEndpointConnectionWithDefaults() *Nextgenv1beta2PrivateEndpointConnection {
	this := Nextgenv1beta2PrivateEndpointConnection{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetName(v string) {
	o.Name = &v
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetId(v string) {
	o.Id = &v
}

// GetTidbId returns the TidbId field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetTidbId() string {
	if o == nil || IsNil(o.TidbId) {
		var ret string
		return ret
	}
	return *o.TidbId
}

// GetTidbIdOk returns a tuple with the TidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbId) {
		return nil, false
	}
	return o.TidbId, true
}

// HasTidbId returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasTidbId() bool {
	if o != nil && !IsNil(o.TidbId) {
		return true
	}

	return false
}

// SetTidbId gets a reference to the given string and assigns it to the TidbId field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetTidbId(v string) {
	o.TidbId = &v
}

// GetPrivateEndpointConnectionId returns the PrivateEndpointConnectionId field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateEndpointConnectionId() string {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		var ret string
		return ret
	}
	return *o.PrivateEndpointConnectionId
}

// GetPrivateEndpointConnectionIdOk returns a tuple with the PrivateEndpointConnectionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateEndpointConnectionIdOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateEndpointConnectionId) {
		return nil, false
	}
	return o.PrivateEndpointConnectionId, true
}

// HasPrivateEndpointConnectionId returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPrivateEndpointConnectionId() bool {
	if o != nil && !IsNil(o.PrivateEndpointConnectionId) {
		return true
	}

	return false
}

// SetPrivateEndpointConnectionId gets a reference to the given string and assigns it to the PrivateEndpointConnectionId field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPrivateEndpointConnectionId(v string) {
	o.PrivateEndpointConnectionId = &v
}

// GetRegionId returns the RegionId field value
func (o *Nextgenv1beta2PrivateEndpointConnection) GetRegionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetRegionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionId, true
}

// SetRegionId sets field value
func (o *Nextgenv1beta2PrivateEndpointConnection) SetRegionId(v string) {
	o.RegionId = v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetCloudProvider() RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetCloudProviderOk() (*RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetCloudProvider(v RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetEndpointId returns the EndpointId field value
func (o *Nextgenv1beta2PrivateEndpointConnection) GetEndpointId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndpointId
}

// GetEndpointIdOk returns a tuple with the EndpointId field value
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetEndpointIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndpointId, true
}

// SetEndpointId sets field value
func (o *Nextgenv1beta2PrivateEndpointConnection) SetEndpointId(v string) {
	o.EndpointId = v
}

// GetEndpointState returns the EndpointState field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetEndpointState() Nextgenv1beta2PrivateEndpointConnectionEndpointState {
	if o == nil || IsNil(o.EndpointState) {
		var ret Nextgenv1beta2PrivateEndpointConnectionEndpointState
		return ret
	}
	return *o.EndpointState
}

// GetEndpointStateOk returns a tuple with the EndpointState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetEndpointStateOk() (*Nextgenv1beta2PrivateEndpointConnectionEndpointState, bool) {
	if o == nil || IsNil(o.EndpointState) {
		return nil, false
	}
	return o.EndpointState, true
}

// HasEndpointState returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasEndpointState() bool {
	if o != nil && !IsNil(o.EndpointState) {
		return true
	}

	return false
}

// SetEndpointState gets a reference to the given Nextgenv1beta2PrivateEndpointConnectionEndpointState and assigns it to the EndpointState field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetEndpointState(v Nextgenv1beta2PrivateEndpointConnectionEndpointState) {
	o.EndpointState = &v
}

// GetMessage returns the Message field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetMessage() string {
	if o == nil || IsNil(o.Message) {
		var ret string
		return ret
	}
	return *o.Message
}

// GetMessageOk returns a tuple with the Message field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetMessageOk() (*string, bool) {
	if o == nil || IsNil(o.Message) {
		return nil, false
	}
	return o.Message, true
}

// HasMessage returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasMessage() bool {
	if o != nil && !IsNil(o.Message) {
		return true
	}

	return false
}

// SetMessage gets a reference to the given string and assigns it to the Message field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetMessage(v string) {
	o.Message = &v
}

// GetPrivateLinkService returns the PrivateLinkService field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateLinkService() string {
	if o == nil || IsNil(o.PrivateLinkService) {
		var ret string
		return ret
	}
	return *o.PrivateLinkService
}

// GetPrivateLinkServiceOk returns a tuple with the PrivateLinkService field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateLinkServiceOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateLinkService) {
		return nil, false
	}
	return o.PrivateLinkService, true
}

// HasPrivateLinkService returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPrivateLinkService() bool {
	if o != nil && !IsNil(o.PrivateLinkService) {
		return true
	}

	return false
}

// SetPrivateLinkService gets a reference to the given string and assigns it to the PrivateLinkService field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPrivateLinkService(v string) {
	o.PrivateLinkService = &v
}

// GetPrivateLinkServiceState returns the PrivateLinkServiceState field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateLinkServiceState() Nextgenv1beta2PrivateLinkServiceState {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		var ret Nextgenv1beta2PrivateLinkServiceState
		return ret
	}
	return *o.PrivateLinkServiceState
}

// GetPrivateLinkServiceStateOk returns a tuple with the PrivateLinkServiceState field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateLinkServiceStateOk() (*Nextgenv1beta2PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.PrivateLinkServiceState) {
		return nil, false
	}
	return o.PrivateLinkServiceState, true
}

// HasPrivateLinkServiceState returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPrivateLinkServiceState() bool {
	if o != nil && !IsNil(o.PrivateLinkServiceState) {
		return true
	}

	return false
}

// SetPrivateLinkServiceState gets a reference to the given Nextgenv1beta2PrivateLinkServiceState and assigns it to the PrivateLinkServiceState field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPrivateLinkServiceState(v Nextgenv1beta2PrivateLinkServiceState) {
	o.PrivateLinkServiceState = &v
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPort() string {
	if o == nil || IsNil(o.Port) {
		var ret string
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPortOk() (*string, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given string and assigns it to the Port field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPort(v string) {
	o.Port = &v
}

// GetPrivateDomainName returns the PrivateDomainName field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateDomainName() string {
	if o == nil || IsNil(o.PrivateDomainName) {
		var ret string
		return ret
	}
	return *o.PrivateDomainName
}

// GetPrivateDomainNameOk returns a tuple with the PrivateDomainName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateDomainNameOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateDomainName) {
		return nil, false
	}
	return o.PrivateDomainName, true
}

// HasPrivateDomainName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPrivateDomainName() bool {
	if o != nil && !IsNil(o.PrivateDomainName) {
		return true
	}

	return false
}

// SetPrivateDomainName gets a reference to the given string and assigns it to the PrivateDomainName field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPrivateDomainName(v string) {
	o.PrivateDomainName = &v
}

// GetPrivateIpAddress returns the PrivateIpAddress field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateIpAddress() string {
	if o == nil || IsNil(o.PrivateIpAddress) {
		var ret string
		return ret
	}
	return *o.PrivateIpAddress
}

// GetPrivateIpAddressOk returns a tuple with the PrivateIpAddress field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) GetPrivateIpAddressOk() (*string, bool) {
	if o == nil || IsNil(o.PrivateIpAddress) {
		return nil, false
	}
	return o.PrivateIpAddress, true
}

// HasPrivateIpAddress returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateEndpointConnection) HasPrivateIpAddress() bool {
	if o != nil && !IsNil(o.PrivateIpAddress) {
		return true
	}

	return false
}

// SetPrivateIpAddress gets a reference to the given string and assigns it to the PrivateIpAddress field.
func (o *Nextgenv1beta2PrivateEndpointConnection) SetPrivateIpAddress(v string) {
	o.PrivateIpAddress = &v
}

func (o Nextgenv1beta2PrivateEndpointConnection) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Nextgenv1beta2PrivateEndpointConnection) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.TidbId) {
		toSerialize["tidbId"] = o.TidbId
	}
	if !IsNil(o.PrivateEndpointConnectionId) {
		toSerialize["privateEndpointConnectionId"] = o.PrivateEndpointConnectionId
	}
	toSerialize["regionId"] = o.RegionId
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}
	toSerialize["endpointId"] = o.EndpointId
	if !IsNil(o.EndpointState) {
		toSerialize["endpointState"] = o.EndpointState
	}
	if !IsNil(o.Message) {
		toSerialize["message"] = o.Message
	}
	if !IsNil(o.PrivateLinkService) {
		toSerialize["privateLinkService"] = o.PrivateLinkService
	}
	if !IsNil(o.PrivateLinkServiceState) {
		toSerialize["privateLinkServiceState"] = o.PrivateLinkServiceState
	}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}
	if !IsNil(o.PrivateDomainName) {
		toSerialize["privateDomainName"] = o.PrivateDomainName
	}
	if !IsNil(o.PrivateIpAddress) {
		toSerialize["privateIpAddress"] = o.PrivateIpAddress
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Nextgenv1beta2PrivateEndpointConnection) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"regionId",
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

	varNextgenv1beta2PrivateEndpointConnection := _Nextgenv1beta2PrivateEndpointConnection{}

	err = json.Unmarshal(data, &varNextgenv1beta2PrivateEndpointConnection)

	if err != nil {
		return err
	}

	*o = Nextgenv1beta2PrivateEndpointConnection(varNextgenv1beta2PrivateEndpointConnection)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "id")
		delete(additionalProperties, "tidbId")
		delete(additionalProperties, "privateEndpointConnectionId")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "endpointId")
		delete(additionalProperties, "endpointState")
		delete(additionalProperties, "message")
		delete(additionalProperties, "privateLinkService")
		delete(additionalProperties, "privateLinkServiceState")
		delete(additionalProperties, "host")
		delete(additionalProperties, "port")
		delete(additionalProperties, "privateDomainName")
		delete(additionalProperties, "privateIpAddress")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNextgenv1beta2PrivateEndpointConnection struct {
	value *Nextgenv1beta2PrivateEndpointConnection
	isSet bool
}

func (v NullableNextgenv1beta2PrivateEndpointConnection) Get() *Nextgenv1beta2PrivateEndpointConnection {
	return v.value
}

func (v *NullableNextgenv1beta2PrivateEndpointConnection) Set(val *Nextgenv1beta2PrivateEndpointConnection) {
	v.value = val
	v.isSet = true
}

func (v NullableNextgenv1beta2PrivateEndpointConnection) IsSet() bool {
	return v.isSet
}

func (v *NullableNextgenv1beta2PrivateEndpointConnection) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNextgenv1beta2PrivateEndpointConnection(val *Nextgenv1beta2PrivateEndpointConnection) *NullableNextgenv1beta2PrivateEndpointConnection {
	return &NullableNextgenv1beta2PrivateEndpointConnection{value: val, isSet: true}
}

func (v NullableNextgenv1beta2PrivateEndpointConnection) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNextgenv1beta2PrivateEndpointConnection) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
