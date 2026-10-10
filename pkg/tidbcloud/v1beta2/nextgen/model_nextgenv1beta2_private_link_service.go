/*
TiDB Cloud Premium API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium) instances. For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/api/tidb-cloud-api-overview/).  The public connection setting endpoints also support TiDB Cloud Essential V2 instances.  # Overview  The TiDB Cloud Premium API (v1beta2) provides [REST](https://en.wikipedia.org/wiki/REST) endpoints to manage TiDB Cloud Premium instances and related resources.  You can use this API to manage the following resources:  - **TiDB Cloud Premium instance**: manage the lifecycle and configuration of TiDB Cloud Premium instances, including passwords, CA certificates, and cloud provider information. - **Public connection setting**: enable or disable public connections and manage IP access lists for Premium and Essential V2 instances. - **Customer-managed encryption key**: retrieve the IAM principal and verify KMS access before creating a Premium instance on AWS or Alibaba Cloud. - **Backup**: manage backups for TiDB Cloud Premium instances, including backup-based restore. - **Region**: retrieve available regions for deploying TiDB Cloud Premium instances.  # Get Started  This guide helps you make your first API call to the TiDB Cloud Premium API. You will learn how to authenticate a request, build a request, and interpret the response.  1. Create a [TiDB Cloud account](https://tidbcloud.com/signup) if you do not already have one. 2. In the [TiDB Cloud console](https://tidbcloud.com/), go to **Organization** > **API Keys** and create an API key. For more information, see [API key management](#section/Authentication/API-key-management). 3. Make your first API call.   To get all TiDB Cloud Premium instances in your organization, run the following command in your terminal. Replace `YOUR_PUBLIC_KEY` and `YOUR_PRIVATE_KEY` with your own key values.   ```bash  curl --digest \\    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\    --request GET \\    --url 'https://cloud.tidbapi.com/v1beta2/tidbs' \\    --header 'Accept: application/json'  ```  4. The API returns a JSON list of your TiDB Cloud Premium instances. If none exist, the response contains an empty list.  # Authentication  The TiDB Cloud API supports [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) with API keys and OAuth Bearer tokens. OAuth clients send the token in the `Authorization: Bearer <token>` header. It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. 5. Configure the role and scope for the API key. For more information about the permissions of a role, see [User roles](https://docs.pingcap.com/tidbcloud/manage-user-access/#user-roles). 6. Click **Generate API Key**. Copy and save the public key and the private key. 7. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 8. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Update Role**. 4. You can update the description and role of the API key. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # API Changelog  This changelog lists all changes to the TiDB Cloud Premium API (v1beta2).  <!-- In reverse chronological order -->  ## 20260924  - Published CMEK IAM principal retrieval and verification APIs for Premium instances on AWS and Alibaba Cloud.  ## 20260923  - Added APIs to get and update public connection settings for Premium and Essential V2 instances.  ## 20260428  - Initial release of the TiDB Cloud Premium API (v1beta2), including the following resources and endpoints:  * TiDB Cloud Premium instance   * [List TiDB Cloud Premium instances](#tag/TiDB-Instance/operation/TidbService_ListTidbs)   * [Create a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_CreateTidb)   * [Get a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetTidb)   * [Delete a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_DeleteTidb)   * [Update a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_UpdateTidb)   * [Reset the root password of a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_ResetRootPassword)   * [Get cloud provider information for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCloudProviderInfo)   * [Get the CA certificate download URL for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCaCertificateDownloadUrl)  * Backup   * [List backups for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_ListTidbBackups)   * [Delete a backup for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_DeleteTidbBackup)   * [Restore a TiDB Cloud Premium instance from a backup](#tag/Backup/operation/TidbService_RestoreTidb)   * [Get restore status for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_GetRestoreStatus)  * Region   * [List regions](#tag/Region/operation/RegionService_ListRegions)

API version: v1beta2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nextgen

import (
	"encoding/json"
)

// checks if the Nextgenv1beta2PrivateLinkService type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Nextgenv1beta2PrivateLinkService{}

// Nextgenv1beta2PrivateLinkService struct for Nextgenv1beta2PrivateLinkService
type Nextgenv1beta2PrivateLinkService struct {
	Name              *string              `json:"name,omitempty"`
	TidbId            *string              `json:"tidbId,omitempty"`
	RegionId          *string              `json:"regionId,omitempty"`
	RegionDisplayName *string              `json:"regionDisplayName,omitempty"`
	CloudProvider     *RegionCloudProvider `json:"cloudProvider,omitempty"`
	// For AWS, it's the service name of the Private Link Service. For GCP, it's the resource name of the service attachment. For Azure, it's service resource ID of the Private Link Service. For AliCloud, it's the service name of the Private Link Service.
	ServiceName *string `json:"serviceName,omitempty"`
	// For AWS, it's the fully qualified domain name (FQDN) shared for all private endpoints, despite which VPC the endpoint located in. For GCP, it's the zone name (suffix of FQDN) shared for all private endpoints located in a single VPC network. The format of FQDN is `<endpoint_name>.<service_dns_name>`. For Azure, it's the zone name shared across public internet. The format of FQDN is `<endpoint_name>-<random_hash>.<service_dns_name>`. For AliCloud, it's the endpoint domain name shared for all private endpoints. The format of FQDN is `<endpoint_name>.<service_dns_name>`.
	ServiceDnsName *string `json:"serviceDnsName,omitempty"`
	// Only available for AWS. Same as the `AvailabilityZones` field in response body of `github.com/aws/aws-sdk-go-v2/service/ec2.DescribeVpcEndpointServices` method.
	AvailableZones []string `json:"availableZones,omitempty"`
	// The port of the private link service.
	ServicePort          *int32                                 `json:"servicePort,omitempty"`
	State                *Nextgenv1beta2PrivateLinkServiceState `json:"state,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Nextgenv1beta2PrivateLinkService Nextgenv1beta2PrivateLinkService

// NewNextgenv1beta2PrivateLinkService instantiates a new Nextgenv1beta2PrivateLinkService object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNextgenv1beta2PrivateLinkService() *Nextgenv1beta2PrivateLinkService {
	this := Nextgenv1beta2PrivateLinkService{}
	return &this
}

// NewNextgenv1beta2PrivateLinkServiceWithDefaults instantiates a new Nextgenv1beta2PrivateLinkService object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNextgenv1beta2PrivateLinkServiceWithDefaults() *Nextgenv1beta2PrivateLinkService {
	this := Nextgenv1beta2PrivateLinkService{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Nextgenv1beta2PrivateLinkService) SetName(v string) {
	o.Name = &v
}

// GetTidbId returns the TidbId field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetTidbId() string {
	if o == nil || IsNil(o.TidbId) {
		var ret string
		return ret
	}
	return *o.TidbId
}

// GetTidbIdOk returns a tuple with the TidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbId) {
		return nil, false
	}
	return o.TidbId, true
}

// HasTidbId returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasTidbId() bool {
	if o != nil && !IsNil(o.TidbId) {
		return true
	}

	return false
}

// SetTidbId gets a reference to the given string and assigns it to the TidbId field.
func (o *Nextgenv1beta2PrivateLinkService) SetTidbId(v string) {
	o.TidbId = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *Nextgenv1beta2PrivateLinkService) SetRegionId(v string) {
	o.RegionId = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *Nextgenv1beta2PrivateLinkService) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetCloudProvider() RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetCloudProviderOk() (*RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Nextgenv1beta2PrivateLinkService) SetCloudProvider(v RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetServiceName returns the ServiceName field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetServiceName() string {
	if o == nil || IsNil(o.ServiceName) {
		var ret string
		return ret
	}
	return *o.ServiceName
}

// GetServiceNameOk returns a tuple with the ServiceName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetServiceNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceName) {
		return nil, false
	}
	return o.ServiceName, true
}

// HasServiceName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasServiceName() bool {
	if o != nil && !IsNil(o.ServiceName) {
		return true
	}

	return false
}

// SetServiceName gets a reference to the given string and assigns it to the ServiceName field.
func (o *Nextgenv1beta2PrivateLinkService) SetServiceName(v string) {
	o.ServiceName = &v
}

// GetServiceDnsName returns the ServiceDnsName field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetServiceDnsName() string {
	if o == nil || IsNil(o.ServiceDnsName) {
		var ret string
		return ret
	}
	return *o.ServiceDnsName
}

// GetServiceDnsNameOk returns a tuple with the ServiceDnsName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetServiceDnsNameOk() (*string, bool) {
	if o == nil || IsNil(o.ServiceDnsName) {
		return nil, false
	}
	return o.ServiceDnsName, true
}

// HasServiceDnsName returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasServiceDnsName() bool {
	if o != nil && !IsNil(o.ServiceDnsName) {
		return true
	}

	return false
}

// SetServiceDnsName gets a reference to the given string and assigns it to the ServiceDnsName field.
func (o *Nextgenv1beta2PrivateLinkService) SetServiceDnsName(v string) {
	o.ServiceDnsName = &v
}

// GetAvailableZones returns the AvailableZones field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetAvailableZones() []string {
	if o == nil || IsNil(o.AvailableZones) {
		var ret []string
		return ret
	}
	return o.AvailableZones
}

// GetAvailableZonesOk returns a tuple with the AvailableZones field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetAvailableZonesOk() ([]string, bool) {
	if o == nil || IsNil(o.AvailableZones) {
		return nil, false
	}
	return o.AvailableZones, true
}

// HasAvailableZones returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasAvailableZones() bool {
	if o != nil && !IsNil(o.AvailableZones) {
		return true
	}

	return false
}

// SetAvailableZones gets a reference to the given []string and assigns it to the AvailableZones field.
func (o *Nextgenv1beta2PrivateLinkService) SetAvailableZones(v []string) {
	o.AvailableZones = v
}

// GetServicePort returns the ServicePort field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetServicePort() int32 {
	if o == nil || IsNil(o.ServicePort) {
		var ret int32
		return ret
	}
	return *o.ServicePort
}

// GetServicePortOk returns a tuple with the ServicePort field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetServicePortOk() (*int32, bool) {
	if o == nil || IsNil(o.ServicePort) {
		return nil, false
	}
	return o.ServicePort, true
}

// HasServicePort returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasServicePort() bool {
	if o != nil && !IsNil(o.ServicePort) {
		return true
	}

	return false
}

// SetServicePort gets a reference to the given int32 and assigns it to the ServicePort field.
func (o *Nextgenv1beta2PrivateLinkService) SetServicePort(v int32) {
	o.ServicePort = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Nextgenv1beta2PrivateLinkService) GetState() Nextgenv1beta2PrivateLinkServiceState {
	if o == nil || IsNil(o.State) {
		var ret Nextgenv1beta2PrivateLinkServiceState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2PrivateLinkService) GetStateOk() (*Nextgenv1beta2PrivateLinkServiceState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Nextgenv1beta2PrivateLinkService) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Nextgenv1beta2PrivateLinkServiceState and assigns it to the State field.
func (o *Nextgenv1beta2PrivateLinkService) SetState(v Nextgenv1beta2PrivateLinkServiceState) {
	o.State = &v
}

func (o Nextgenv1beta2PrivateLinkService) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Nextgenv1beta2PrivateLinkService) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbId) {
		toSerialize["tidbId"] = o.TidbId
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
	if !IsNil(o.ServiceName) {
		toSerialize["serviceName"] = o.ServiceName
	}
	if !IsNil(o.ServiceDnsName) {
		toSerialize["serviceDnsName"] = o.ServiceDnsName
	}
	if !IsNil(o.AvailableZones) {
		toSerialize["availableZones"] = o.AvailableZones
	}
	if !IsNil(o.ServicePort) {
		toSerialize["servicePort"] = o.ServicePort
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Nextgenv1beta2PrivateLinkService) UnmarshalJSON(data []byte) (err error) {
	varNextgenv1beta2PrivateLinkService := _Nextgenv1beta2PrivateLinkService{}

	err = json.Unmarshal(data, &varNextgenv1beta2PrivateLinkService)

	if err != nil {
		return err
	}

	*o = Nextgenv1beta2PrivateLinkService(varNextgenv1beta2PrivateLinkService)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbId")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "serviceName")
		delete(additionalProperties, "serviceDnsName")
		delete(additionalProperties, "availableZones")
		delete(additionalProperties, "servicePort")
		delete(additionalProperties, "state")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNextgenv1beta2PrivateLinkService struct {
	value *Nextgenv1beta2PrivateLinkService
	isSet bool
}

func (v NullableNextgenv1beta2PrivateLinkService) Get() *Nextgenv1beta2PrivateLinkService {
	return v.value
}

func (v *NullableNextgenv1beta2PrivateLinkService) Set(val *Nextgenv1beta2PrivateLinkService) {
	v.value = val
	v.isSet = true
}

func (v NullableNextgenv1beta2PrivateLinkService) IsSet() bool {
	return v.isSet
}

func (v *NullableNextgenv1beta2PrivateLinkService) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNextgenv1beta2PrivateLinkService(val *Nextgenv1beta2PrivateLinkService) *NullableNextgenv1beta2PrivateLinkService {
	return &NullableNextgenv1beta2PrivateLinkService{value: val, isSet: true}
}

func (v NullableNextgenv1beta2PrivateLinkService) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNextgenv1beta2PrivateLinkService) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
