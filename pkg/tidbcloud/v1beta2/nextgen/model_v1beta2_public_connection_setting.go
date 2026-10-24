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

// checks if the V1beta2PublicConnectionSetting type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta2PublicConnectionSetting{}

// V1beta2PublicConnectionSetting struct for V1beta2PublicConnectionSetting
type V1beta2PublicConnectionSetting struct {
	Name              *string                                      `json:"name,omitempty"`
	TidbId            *string                                      `json:"tidbId,omitempty"`
	Enabled           NullableBool                                 `json:"enabled,omitempty"`
	IpAccessList      []V1beta2PublicConnectionSettingIpAccessList `json:"ipAccessList,omitempty"`
	ClearIpAccessList NullableBool                                 `json:"clearIpAccessList,omitempty"`
	// The public egress IP addresses used by the TiDB Cloud instance.  This field is returned only when the public endpoint is enabled. You can use these IP addresses to configure allowlists on upstream services.
	EgressIps            []string `json:"egressIps,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta2PublicConnectionSetting V1beta2PublicConnectionSetting

// NewV1beta2PublicConnectionSetting instantiates a new V1beta2PublicConnectionSetting object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta2PublicConnectionSetting() *V1beta2PublicConnectionSetting {
	this := V1beta2PublicConnectionSetting{}
	return &this
}

// NewV1beta2PublicConnectionSettingWithDefaults instantiates a new V1beta2PublicConnectionSetting object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta2PublicConnectionSettingWithDefaults() *V1beta2PublicConnectionSetting {
	this := V1beta2PublicConnectionSetting{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta2PublicConnectionSetting) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2PublicConnectionSetting) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta2PublicConnectionSetting) SetName(v string) {
	o.Name = &v
}

// GetTidbId returns the TidbId field value if set, zero value otherwise.
func (o *V1beta2PublicConnectionSetting) GetTidbId() string {
	if o == nil || IsNil(o.TidbId) {
		var ret string
		return ret
	}
	return *o.TidbId
}

// GetTidbIdOk returns a tuple with the TidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2PublicConnectionSetting) GetTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbId) {
		return nil, false
	}
	return o.TidbId, true
}

// HasTidbId returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasTidbId() bool {
	if o != nil && !IsNil(o.TidbId) {
		return true
	}

	return false
}

// SetTidbId gets a reference to the given string and assigns it to the TidbId field.
func (o *V1beta2PublicConnectionSetting) SetTidbId(v string) {
	o.TidbId = &v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta2PublicConnectionSetting) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled.Get()) {
		var ret bool
		return ret
	}
	return *o.Enabled.Get()
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta2PublicConnectionSetting) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.Enabled.Get(), o.Enabled.IsSet()
}

// HasEnabled returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasEnabled() bool {
	if o != nil && o.Enabled.IsSet() {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given NullableBool and assigns it to the Enabled field.
func (o *V1beta2PublicConnectionSetting) SetEnabled(v bool) {
	o.Enabled.Set(&v)
}

// SetEnabledNil sets the value for Enabled to be an explicit nil
func (o *V1beta2PublicConnectionSetting) SetEnabledNil() {
	o.Enabled.Set(nil)
}

// UnsetEnabled ensures that no value is present for Enabled, not even an explicit nil
func (o *V1beta2PublicConnectionSetting) UnsetEnabled() {
	o.Enabled.Unset()
}

// GetIpAccessList returns the IpAccessList field value if set, zero value otherwise.
func (o *V1beta2PublicConnectionSetting) GetIpAccessList() []V1beta2PublicConnectionSettingIpAccessList {
	if o == nil || IsNil(o.IpAccessList) {
		var ret []V1beta2PublicConnectionSettingIpAccessList
		return ret
	}
	return o.IpAccessList
}

// GetIpAccessListOk returns a tuple with the IpAccessList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2PublicConnectionSetting) GetIpAccessListOk() ([]V1beta2PublicConnectionSettingIpAccessList, bool) {
	if o == nil || IsNil(o.IpAccessList) {
		return nil, false
	}
	return o.IpAccessList, true
}

// HasIpAccessList returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasIpAccessList() bool {
	if o != nil && !IsNil(o.IpAccessList) {
		return true
	}

	return false
}

// SetIpAccessList gets a reference to the given []V1beta2PublicConnectionSettingIpAccessList and assigns it to the IpAccessList field.
func (o *V1beta2PublicConnectionSetting) SetIpAccessList(v []V1beta2PublicConnectionSettingIpAccessList) {
	o.IpAccessList = v
}

// GetClearIpAccessList returns the ClearIpAccessList field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta2PublicConnectionSetting) GetClearIpAccessList() bool {
	if o == nil || IsNil(o.ClearIpAccessList.Get()) {
		var ret bool
		return ret
	}
	return *o.ClearIpAccessList.Get()
}

// GetClearIpAccessListOk returns a tuple with the ClearIpAccessList field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta2PublicConnectionSetting) GetClearIpAccessListOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.ClearIpAccessList.Get(), o.ClearIpAccessList.IsSet()
}

// HasClearIpAccessList returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasClearIpAccessList() bool {
	if o != nil && o.ClearIpAccessList.IsSet() {
		return true
	}

	return false
}

// SetClearIpAccessList gets a reference to the given NullableBool and assigns it to the ClearIpAccessList field.
func (o *V1beta2PublicConnectionSetting) SetClearIpAccessList(v bool) {
	o.ClearIpAccessList.Set(&v)
}

// SetClearIpAccessListNil sets the value for ClearIpAccessList to be an explicit nil
func (o *V1beta2PublicConnectionSetting) SetClearIpAccessListNil() {
	o.ClearIpAccessList.Set(nil)
}

// UnsetClearIpAccessList ensures that no value is present for ClearIpAccessList, not even an explicit nil
func (o *V1beta2PublicConnectionSetting) UnsetClearIpAccessList() {
	o.ClearIpAccessList.Unset()
}

// GetEgressIps returns the EgressIps field value if set, zero value otherwise.
func (o *V1beta2PublicConnectionSetting) GetEgressIps() []string {
	if o == nil || IsNil(o.EgressIps) {
		var ret []string
		return ret
	}
	return o.EgressIps
}

// GetEgressIpsOk returns a tuple with the EgressIps field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2PublicConnectionSetting) GetEgressIpsOk() ([]string, bool) {
	if o == nil || IsNil(o.EgressIps) {
		return nil, false
	}
	return o.EgressIps, true
}

// HasEgressIps returns a boolean if a field has been set.
func (o *V1beta2PublicConnectionSetting) HasEgressIps() bool {
	if o != nil && !IsNil(o.EgressIps) {
		return true
	}

	return false
}

// SetEgressIps gets a reference to the given []string and assigns it to the EgressIps field.
func (o *V1beta2PublicConnectionSetting) SetEgressIps(v []string) {
	o.EgressIps = v
}

func (o V1beta2PublicConnectionSetting) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta2PublicConnectionSetting) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbId) {
		toSerialize["tidbId"] = o.TidbId
	}
	if o.Enabled.IsSet() {
		toSerialize["enabled"] = o.Enabled.Get()
	}
	if !IsNil(o.IpAccessList) {
		toSerialize["ipAccessList"] = o.IpAccessList
	}
	if o.ClearIpAccessList.IsSet() {
		toSerialize["clearIpAccessList"] = o.ClearIpAccessList.Get()
	}
	if !IsNil(o.EgressIps) {
		toSerialize["egressIps"] = o.EgressIps
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta2PublicConnectionSetting) UnmarshalJSON(data []byte) (err error) {
	varV1beta2PublicConnectionSetting := _V1beta2PublicConnectionSetting{}

	err = json.Unmarshal(data, &varV1beta2PublicConnectionSetting)

	if err != nil {
		return err
	}

	*o = V1beta2PublicConnectionSetting(varV1beta2PublicConnectionSetting)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbId")
		delete(additionalProperties, "enabled")
		delete(additionalProperties, "ipAccessList")
		delete(additionalProperties, "clearIpAccessList")
		delete(additionalProperties, "egressIps")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta2PublicConnectionSetting struct {
	value *V1beta2PublicConnectionSetting
	isSet bool
}

func (v NullableV1beta2PublicConnectionSetting) Get() *V1beta2PublicConnectionSetting {
	return v.value
}

func (v *NullableV1beta2PublicConnectionSetting) Set(val *V1beta2PublicConnectionSetting) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta2PublicConnectionSetting) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta2PublicConnectionSetting) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta2PublicConnectionSetting(val *V1beta2PublicConnectionSetting) *NullableV1beta2PublicConnectionSetting {
	return &NullableV1beta2PublicConnectionSetting{value: val, isSet: true}
}

func (v NullableV1beta2PublicConnectionSetting) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta2PublicConnectionSetting) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
