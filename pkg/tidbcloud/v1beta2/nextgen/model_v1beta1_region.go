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

// checks if the V1beta1Region type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1Region{}

// V1beta1Region A representation of a region for deploying TiDB clusters.
type V1beta1Region struct {
	// The unique name of the region, in the format of `regions/{region_id}`. For example, `regions/aws-us-west-2`.
	Name *string `json:"name,omitempty" validate:"regexp=^regions\\/(aws|gcp|azure)-(.+)$"`
	// The unique identifier for the region, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-west-2`.
	RegionId *string `json:"regionId,omitempty" validate:"regexp=^(aws|gcp|azure|alicloud)-[a-z0-9-]+$"`
	// The cloud provider that offers the region.  - `\"aws\"`: Amazon Web Services  - `\"gcp\"`: Google Cloud  - `\"azure\"`: Microsoft Azure  - `\"alicloud\"`: Alibaba Cloud
	CloudProvider *RegionCloudProvider `json:"cloudProvider,omitempty"`
	// A human-readable name for the region. For example, `Oregon (us-west-2)`.
	DisplayName *string `json:"displayName,omitempty"`
	// **Deprecated.** Use `cloudProvider` instead. The name of the cloud provider. For example, `aws`, `gcp`, `azure`, or `alicloud`.
	Provider NullableString `json:"provider,omitempty"`
	// The service plans available in this region. Premium, Essential_V2, and Premium_Reserved can appear together in the same logical region. Starter, Essential, and BYOC are standalone service plans and cannot share a logical region with any other service plan.
	ServicePlans          []V1beta1ServicePlan     `json:"servicePlans,omitempty"`
	SupportedServicePlans []V1beta1ServicePlanInfo `json:"supportedServicePlans,omitempty"`
	AdditionalProperties  map[string]interface{}
}

type _V1beta1Region V1beta1Region

// NewV1beta1Region instantiates a new V1beta1Region object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1Region() *V1beta1Region {
	this := V1beta1Region{}
	return &this
}

// NewV1beta1RegionWithDefaults instantiates a new V1beta1Region object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1RegionWithDefaults() *V1beta1Region {
	this := V1beta1Region{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1Region) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1Region) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1Region) SetName(v string) {
	o.Name = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *V1beta1Region) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *V1beta1Region) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *V1beta1Region) SetRegionId(v string) {
	o.RegionId = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *V1beta1Region) GetCloudProvider() RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetCloudProviderOk() (*RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *V1beta1Region) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given RegionCloudProvider and assigns it to the CloudProvider field.
func (o *V1beta1Region) SetCloudProvider(v RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *V1beta1Region) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *V1beta1Region) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *V1beta1Region) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetProvider returns the Provider field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1Region) GetProvider() string {
	if o == nil || IsNil(o.Provider.Get()) {
		var ret string
		return ret
	}
	return *o.Provider.Get()
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1Region) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Provider.Get(), o.Provider.IsSet()
}

// HasProvider returns a boolean if a field has been set.
func (o *V1beta1Region) HasProvider() bool {
	if o != nil && o.Provider.IsSet() {
		return true
	}

	return false
}

// SetProvider gets a reference to the given NullableString and assigns it to the Provider field.
func (o *V1beta1Region) SetProvider(v string) {
	o.Provider.Set(&v)
}

// SetProviderNil sets the value for Provider to be an explicit nil
func (o *V1beta1Region) SetProviderNil() {
	o.Provider.Set(nil)
}

// UnsetProvider ensures that no value is present for Provider, not even an explicit nil
func (o *V1beta1Region) UnsetProvider() {
	o.Provider.Unset()
}

// GetServicePlans returns the ServicePlans field value if set, zero value otherwise.
func (o *V1beta1Region) GetServicePlans() []V1beta1ServicePlan {
	if o == nil || IsNil(o.ServicePlans) {
		var ret []V1beta1ServicePlan
		return ret
	}
	return o.ServicePlans
}

// GetServicePlansOk returns a tuple with the ServicePlans field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetServicePlansOk() ([]V1beta1ServicePlan, bool) {
	if o == nil || IsNil(o.ServicePlans) {
		return nil, false
	}
	return o.ServicePlans, true
}

// HasServicePlans returns a boolean if a field has been set.
func (o *V1beta1Region) HasServicePlans() bool {
	if o != nil && !IsNil(o.ServicePlans) {
		return true
	}

	return false
}

// SetServicePlans gets a reference to the given []V1beta1ServicePlan and assigns it to the ServicePlans field.
func (o *V1beta1Region) SetServicePlans(v []V1beta1ServicePlan) {
	o.ServicePlans = v
}

// GetSupportedServicePlans returns the SupportedServicePlans field value if set, zero value otherwise.
func (o *V1beta1Region) GetSupportedServicePlans() []V1beta1ServicePlanInfo {
	if o == nil || IsNil(o.SupportedServicePlans) {
		var ret []V1beta1ServicePlanInfo
		return ret
	}
	return o.SupportedServicePlans
}

// GetSupportedServicePlansOk returns a tuple with the SupportedServicePlans field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1Region) GetSupportedServicePlansOk() ([]V1beta1ServicePlanInfo, bool) {
	if o == nil || IsNil(o.SupportedServicePlans) {
		return nil, false
	}
	return o.SupportedServicePlans, true
}

// HasSupportedServicePlans returns a boolean if a field has been set.
func (o *V1beta1Region) HasSupportedServicePlans() bool {
	if o != nil && !IsNil(o.SupportedServicePlans) {
		return true
	}

	return false
}

// SetSupportedServicePlans gets a reference to the given []V1beta1ServicePlanInfo and assigns it to the SupportedServicePlans field.
func (o *V1beta1Region) SetSupportedServicePlans(v []V1beta1ServicePlanInfo) {
	o.SupportedServicePlans = v
}

func (o V1beta1Region) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1Region) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	if o.Provider.IsSet() {
		toSerialize["provider"] = o.Provider.Get()
	}
	if !IsNil(o.ServicePlans) {
		toSerialize["servicePlans"] = o.ServicePlans
	}
	if !IsNil(o.SupportedServicePlans) {
		toSerialize["supportedServicePlans"] = o.SupportedServicePlans
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1Region) UnmarshalJSON(data []byte) (err error) {
	varV1beta1Region := _V1beta1Region{}

	err = json.Unmarshal(data, &varV1beta1Region)

	if err != nil {
		return err
	}

	*o = V1beta1Region(varV1beta1Region)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "provider")
		delete(additionalProperties, "servicePlans")
		delete(additionalProperties, "supportedServicePlans")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1Region struct {
	value *V1beta1Region
	isSet bool
}

func (v NullableV1beta1Region) Get() *V1beta1Region {
	return v.value
}

func (v *NullableV1beta1Region) Set(val *V1beta1Region) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1Region) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1Region) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1Region(val *V1beta1Region) *NullableV1beta1Region {
	return &NullableV1beta1Region{value: val, isSet: true}
}

func (v NullableV1beta1Region) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1Region) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
