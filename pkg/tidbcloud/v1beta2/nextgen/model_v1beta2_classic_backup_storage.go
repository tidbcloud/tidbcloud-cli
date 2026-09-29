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

// checks if the V1beta2ClassicBackupStorage type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta2ClassicBackupStorage{}

// V1beta2ClassicBackupStorage The storage configuration for a classic backup.
type V1beta2ClassicBackupStorage struct {
	// The cloud provider where the backup object storage is located.  **Note**: Currently, only `\"aws\"` and `\"alicloud\"` are supported.
	CloudProvider RegionCloudProvider `json:"cloudProvider"`
	// The region where the backup object storage is located, in the format of `{cloud_provider}-{region_code}`. For example, `alicloud-cn-shanghai`.  This field is required when `cloud_provider` is `\"alicloud\"`.
	RegionId *string `json:"regionId,omitempty"`
	// The URI of the backup bucket.  - For Amazon S3, use the format `s3://{bucket_name}/{path}`. For example, `s3://my-backup-bucket/tidb-backups/`.  - For Alibaba Cloud OSS, use the format `oss://{bucket_name}/{path}`. For example, `oss://my-backup-bucket/tidb-backups/`.
	BucketUri string `json:"bucketUri"`
	// The authentication credentials for accessing the backup storage.
	Auth                 V1beta2BackupStorageAuth `json:"auth"`
	AdditionalProperties map[string]interface{}
}

type _V1beta2ClassicBackupStorage V1beta2ClassicBackupStorage

// NewV1beta2ClassicBackupStorage instantiates a new V1beta2ClassicBackupStorage object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta2ClassicBackupStorage(cloudProvider RegionCloudProvider, bucketUri string, auth V1beta2BackupStorageAuth) *V1beta2ClassicBackupStorage {
	this := V1beta2ClassicBackupStorage{}
	this.CloudProvider = cloudProvider
	this.BucketUri = bucketUri
	this.Auth = auth
	return &this
}

// NewV1beta2ClassicBackupStorageWithDefaults instantiates a new V1beta2ClassicBackupStorage object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta2ClassicBackupStorageWithDefaults() *V1beta2ClassicBackupStorage {
	this := V1beta2ClassicBackupStorage{}
	return &this
}

// GetCloudProvider returns the CloudProvider field value
func (o *V1beta2ClassicBackupStorage) GetCloudProvider() RegionCloudProvider {
	if o == nil {
		var ret RegionCloudProvider
		return ret
	}

	return o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value
// and a boolean to check if the value has been set.
func (o *V1beta2ClassicBackupStorage) GetCloudProviderOk() (*RegionCloudProvider, bool) {
	if o == nil {
		return nil, false
	}
	return &o.CloudProvider, true
}

// SetCloudProvider sets field value
func (o *V1beta2ClassicBackupStorage) SetCloudProvider(v RegionCloudProvider) {
	o.CloudProvider = v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *V1beta2ClassicBackupStorage) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2ClassicBackupStorage) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *V1beta2ClassicBackupStorage) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *V1beta2ClassicBackupStorage) SetRegionId(v string) {
	o.RegionId = &v
}

// GetBucketUri returns the BucketUri field value
func (o *V1beta2ClassicBackupStorage) GetBucketUri() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.BucketUri
}

// GetBucketUriOk returns a tuple with the BucketUri field value
// and a boolean to check if the value has been set.
func (o *V1beta2ClassicBackupStorage) GetBucketUriOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.BucketUri, true
}

// SetBucketUri sets field value
func (o *V1beta2ClassicBackupStorage) SetBucketUri(v string) {
	o.BucketUri = v
}

// GetAuth returns the Auth field value
func (o *V1beta2ClassicBackupStorage) GetAuth() V1beta2BackupStorageAuth {
	if o == nil {
		var ret V1beta2BackupStorageAuth
		return ret
	}

	return o.Auth
}

// GetAuthOk returns a tuple with the Auth field value
// and a boolean to check if the value has been set.
func (o *V1beta2ClassicBackupStorage) GetAuthOk() (*V1beta2BackupStorageAuth, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Auth, true
}

// SetAuth sets field value
func (o *V1beta2ClassicBackupStorage) SetAuth(v V1beta2BackupStorageAuth) {
	o.Auth = v
}

func (o V1beta2ClassicBackupStorage) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta2ClassicBackupStorage) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["cloudProvider"] = o.CloudProvider
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	toSerialize["bucketUri"] = o.BucketUri
	toSerialize["auth"] = o.Auth

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta2ClassicBackupStorage) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"cloudProvider",
		"bucketUri",
		"auth",
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

	varV1beta2ClassicBackupStorage := _V1beta2ClassicBackupStorage{}

	err = json.Unmarshal(data, &varV1beta2ClassicBackupStorage)

	if err != nil {
		return err
	}

	*o = V1beta2ClassicBackupStorage(varV1beta2ClassicBackupStorage)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "bucketUri")
		delete(additionalProperties, "auth")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta2ClassicBackupStorage struct {
	value *V1beta2ClassicBackupStorage
	isSet bool
}

func (v NullableV1beta2ClassicBackupStorage) Get() *V1beta2ClassicBackupStorage {
	return v.value
}

func (v *NullableV1beta2ClassicBackupStorage) Set(val *V1beta2ClassicBackupStorage) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta2ClassicBackupStorage) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta2ClassicBackupStorage) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta2ClassicBackupStorage(val *V1beta2ClassicBackupStorage) *NullableV1beta2ClassicBackupStorage {
	return &NullableV1beta2ClassicBackupStorage{value: val, isSet: true}
}

func (v NullableV1beta2ClassicBackupStorage) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta2ClassicBackupStorage) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
