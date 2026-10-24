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
	"time"
)

// checks if the V1beta2RestoreTidbRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta2RestoreTidbRequest{}

// V1beta2RestoreTidbRequest struct for V1beta2RestoreTidbRequest
type V1beta2RestoreTidbRequest struct {
	// The configuration of the TiDB Cloud Premium instance to create.  **Note**: Currently, only TiDB Cloud Premium instances with `servicePlan` set to `Premium` are supported.
	Tidb Nextgenv1beta2Tidb `json:"tidb"`
	// The ID of the source TiDB Cloud Premium instance containing the backup data.
	SourceTidbId *string            `json:"sourceTidbId,omitempty"`
	RestoreMode  V1beta2RestoreMode `json:"restoreMode"`
	// The ID of the backup to restore from.  This field is required when `restoreMode` is `\"SNAPSHOT\"` or `\"DEDICATED_SNAPSHOT\"`.
	BackupId *string `json:"backupId,omitempty"`
	// The point in time to restore data to, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.  This field is required when `restoreMode` is `\"PITR\"`.
	PointInTime *time.Time `json:"pointInTime,omitempty"`
	// The classic backup storage configuration.  This field is required when `restoreMode` is `\"CLASSIC_SNAPSHOT\"`.
	ClassicBackupStorage *V1beta2ClassicBackupStorage `json:"classicBackupStorage,omitempty"`
	// If set to `true`, the request is validated but not executed. Defaults to `false`.
	ValidateOnly         *bool `json:"validateOnly,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta2RestoreTidbRequest V1beta2RestoreTidbRequest

// NewV1beta2RestoreTidbRequest instantiates a new V1beta2RestoreTidbRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta2RestoreTidbRequest(tidb Nextgenv1beta2Tidb, restoreMode V1beta2RestoreMode) *V1beta2RestoreTidbRequest {
	this := V1beta2RestoreTidbRequest{}
	this.Tidb = tidb
	this.RestoreMode = restoreMode
	return &this
}

// NewV1beta2RestoreTidbRequestWithDefaults instantiates a new V1beta2RestoreTidbRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta2RestoreTidbRequestWithDefaults() *V1beta2RestoreTidbRequest {
	this := V1beta2RestoreTidbRequest{}
	return &this
}

// GetTidb returns the Tidb field value
func (o *V1beta2RestoreTidbRequest) GetTidb() Nextgenv1beta2Tidb {
	if o == nil {
		var ret Nextgenv1beta2Tidb
		return ret
	}

	return o.Tidb
}

// GetTidbOk returns a tuple with the Tidb field value
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetTidbOk() (*Nextgenv1beta2Tidb, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Tidb, true
}

// SetTidb sets field value
func (o *V1beta2RestoreTidbRequest) SetTidb(v Nextgenv1beta2Tidb) {
	o.Tidb = v
}

// GetSourceTidbId returns the SourceTidbId field value if set, zero value otherwise.
func (o *V1beta2RestoreTidbRequest) GetSourceTidbId() string {
	if o == nil || IsNil(o.SourceTidbId) {
		var ret string
		return ret
	}
	return *o.SourceTidbId
}

// GetSourceTidbIdOk returns a tuple with the SourceTidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetSourceTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.SourceTidbId) {
		return nil, false
	}
	return o.SourceTidbId, true
}

// HasSourceTidbId returns a boolean if a field has been set.
func (o *V1beta2RestoreTidbRequest) HasSourceTidbId() bool {
	if o != nil && !IsNil(o.SourceTidbId) {
		return true
	}

	return false
}

// SetSourceTidbId gets a reference to the given string and assigns it to the SourceTidbId field.
func (o *V1beta2RestoreTidbRequest) SetSourceTidbId(v string) {
	o.SourceTidbId = &v
}

// GetRestoreMode returns the RestoreMode field value
func (o *V1beta2RestoreTidbRequest) GetRestoreMode() V1beta2RestoreMode {
	if o == nil {
		var ret V1beta2RestoreMode
		return ret
	}

	return o.RestoreMode
}

// GetRestoreModeOk returns a tuple with the RestoreMode field value
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetRestoreModeOk() (*V1beta2RestoreMode, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RestoreMode, true
}

// SetRestoreMode sets field value
func (o *V1beta2RestoreTidbRequest) SetRestoreMode(v V1beta2RestoreMode) {
	o.RestoreMode = v
}

// GetBackupId returns the BackupId field value if set, zero value otherwise.
func (o *V1beta2RestoreTidbRequest) GetBackupId() string {
	if o == nil || IsNil(o.BackupId) {
		var ret string
		return ret
	}
	return *o.BackupId
}

// GetBackupIdOk returns a tuple with the BackupId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetBackupIdOk() (*string, bool) {
	if o == nil || IsNil(o.BackupId) {
		return nil, false
	}
	return o.BackupId, true
}

// HasBackupId returns a boolean if a field has been set.
func (o *V1beta2RestoreTidbRequest) HasBackupId() bool {
	if o != nil && !IsNil(o.BackupId) {
		return true
	}

	return false
}

// SetBackupId gets a reference to the given string and assigns it to the BackupId field.
func (o *V1beta2RestoreTidbRequest) SetBackupId(v string) {
	o.BackupId = &v
}

// GetPointInTime returns the PointInTime field value if set, zero value otherwise.
func (o *V1beta2RestoreTidbRequest) GetPointInTime() time.Time {
	if o == nil || IsNil(o.PointInTime) {
		var ret time.Time
		return ret
	}
	return *o.PointInTime
}

// GetPointInTimeOk returns a tuple with the PointInTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetPointInTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.PointInTime) {
		return nil, false
	}
	return o.PointInTime, true
}

// HasPointInTime returns a boolean if a field has been set.
func (o *V1beta2RestoreTidbRequest) HasPointInTime() bool {
	if o != nil && !IsNil(o.PointInTime) {
		return true
	}

	return false
}

// SetPointInTime gets a reference to the given time.Time and assigns it to the PointInTime field.
func (o *V1beta2RestoreTidbRequest) SetPointInTime(v time.Time) {
	o.PointInTime = &v
}

// GetClassicBackupStorage returns the ClassicBackupStorage field value if set, zero value otherwise.
func (o *V1beta2RestoreTidbRequest) GetClassicBackupStorage() V1beta2ClassicBackupStorage {
	if o == nil || IsNil(o.ClassicBackupStorage) {
		var ret V1beta2ClassicBackupStorage
		return ret
	}
	return *o.ClassicBackupStorage
}

// GetClassicBackupStorageOk returns a tuple with the ClassicBackupStorage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetClassicBackupStorageOk() (*V1beta2ClassicBackupStorage, bool) {
	if o == nil || IsNil(o.ClassicBackupStorage) {
		return nil, false
	}
	return o.ClassicBackupStorage, true
}

// HasClassicBackupStorage returns a boolean if a field has been set.
func (o *V1beta2RestoreTidbRequest) HasClassicBackupStorage() bool {
	if o != nil && !IsNil(o.ClassicBackupStorage) {
		return true
	}

	return false
}

// SetClassicBackupStorage gets a reference to the given V1beta2ClassicBackupStorage and assigns it to the ClassicBackupStorage field.
func (o *V1beta2RestoreTidbRequest) SetClassicBackupStorage(v V1beta2ClassicBackupStorage) {
	o.ClassicBackupStorage = &v
}

// GetValidateOnly returns the ValidateOnly field value if set, zero value otherwise.
func (o *V1beta2RestoreTidbRequest) GetValidateOnly() bool {
	if o == nil || IsNil(o.ValidateOnly) {
		var ret bool
		return ret
	}
	return *o.ValidateOnly
}

// GetValidateOnlyOk returns a tuple with the ValidateOnly field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2RestoreTidbRequest) GetValidateOnlyOk() (*bool, bool) {
	if o == nil || IsNil(o.ValidateOnly) {
		return nil, false
	}
	return o.ValidateOnly, true
}

// HasValidateOnly returns a boolean if a field has been set.
func (o *V1beta2RestoreTidbRequest) HasValidateOnly() bool {
	if o != nil && !IsNil(o.ValidateOnly) {
		return true
	}

	return false
}

// SetValidateOnly gets a reference to the given bool and assigns it to the ValidateOnly field.
func (o *V1beta2RestoreTidbRequest) SetValidateOnly(v bool) {
	o.ValidateOnly = &v
}

func (o V1beta2RestoreTidbRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta2RestoreTidbRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["tidb"] = o.Tidb
	if !IsNil(o.SourceTidbId) {
		toSerialize["sourceTidbId"] = o.SourceTidbId
	}
	toSerialize["restoreMode"] = o.RestoreMode
	if !IsNil(o.BackupId) {
		toSerialize["backupId"] = o.BackupId
	}
	if !IsNil(o.PointInTime) {
		toSerialize["pointInTime"] = o.PointInTime
	}
	if !IsNil(o.ClassicBackupStorage) {
		toSerialize["classicBackupStorage"] = o.ClassicBackupStorage
	}
	if !IsNil(o.ValidateOnly) {
		toSerialize["validateOnly"] = o.ValidateOnly
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta2RestoreTidbRequest) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"tidb",
		"restoreMode",
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

	varV1beta2RestoreTidbRequest := _V1beta2RestoreTidbRequest{}

	err = json.Unmarshal(data, &varV1beta2RestoreTidbRequest)

	if err != nil {
		return err
	}

	*o = V1beta2RestoreTidbRequest(varV1beta2RestoreTidbRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "tidb")
		delete(additionalProperties, "sourceTidbId")
		delete(additionalProperties, "restoreMode")
		delete(additionalProperties, "backupId")
		delete(additionalProperties, "pointInTime")
		delete(additionalProperties, "classicBackupStorage")
		delete(additionalProperties, "validateOnly")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta2RestoreTidbRequest struct {
	value *V1beta2RestoreTidbRequest
	isSet bool
}

func (v NullableV1beta2RestoreTidbRequest) Get() *V1beta2RestoreTidbRequest {
	return v.value
}

func (v *NullableV1beta2RestoreTidbRequest) Set(val *V1beta2RestoreTidbRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta2RestoreTidbRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta2RestoreTidbRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta2RestoreTidbRequest(val *V1beta2RestoreTidbRequest) *NullableV1beta2RestoreTidbRequest {
	return &NullableV1beta2RestoreTidbRequest{value: val, isSet: true}
}

func (v NullableV1beta2RestoreTidbRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta2RestoreTidbRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
