/*
TiDB Cloud Premium API

*TiDB Cloud API is in beta.*  This API manages [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium) instances. For more information about TiDB Cloud API, see [TiDB Cloud API Overview](https://docs.pingcap.com/api/tidb-cloud-api-overview/).  The public connection setting endpoints also support TiDB Cloud Essential V2 instances.  # Overview  The TiDB Cloud Premium API (v1beta2) provides [REST](https://en.wikipedia.org/wiki/REST) endpoints to manage TiDB Cloud Premium instances and related resources.  You can use this API to manage the following resources:  - **TiDB Cloud Premium instance**: manage the lifecycle and configuration of TiDB Cloud Premium instances, including passwords, CA certificates, and cloud provider information. - **Public connection setting**: enable or disable public connections and manage IP access lists for Premium and Essential V2 instances. - **Customer-managed encryption key**: retrieve the IAM principal and verify KMS access before creating a Premium instance on AWS or Alibaba Cloud. - **Backup**: manage backups for TiDB Cloud Premium instances, including backup-based restore. - **Region**: retrieve available regions for deploying TiDB Cloud Premium instances.  # Get Started  This guide helps you make your first API call to the TiDB Cloud Premium API. You will learn how to authenticate a request, build a request, and interpret the response.  1. Create a [TiDB Cloud account](https://tidbcloud.com/signup) if you do not already have one. 2. In the [TiDB Cloud console](https://tidbcloud.com/), go to **Organization** > **API Keys** and create an API key. For more information, see [API key management](#section/Authentication/API-key-management). 3. Make your first API call.   To get all TiDB Cloud Premium instances in your organization, run the following command in your terminal. Replace `YOUR_PUBLIC_KEY` and `YOUR_PRIVATE_KEY` with your own key values.   ```bash  curl --digest \\    --user 'YOUR_PUBLIC_KEY:YOUR_PRIVATE_KEY' \\    --request GET \\    --url 'https://cloud.tidbapi.com/v1beta2/tidbs' \\    --header 'Accept: application/json'  ```  4. The API returns a JSON list of your TiDB Cloud Premium instances. If none exist, the response contains an empty list.  # Authentication  The TiDB Cloud API supports [HTTP Digest Authentication](https://en.wikipedia.org/wiki/Digest_access_authentication) with API keys and OAuth Bearer tokens. OAuth clients send the token in the `Authorization: Bearer <token>` header. It protects your private key from being sent over the network. For more details about HTTP Digest Authentication, refer to the [IETF RFC](https://datatracker.ietf.org/doc/html/rfc7616).  ## API key overview  - The API key contains a public key and a private key, which act as the username and password required in the HTTP Digest Authentication. The private key only displays upon the key creation. - The API key belongs to your organization and acts as the `Organization Owner` role. You can check [permissions of owner](https://docs.pingcap.com/tidbcloud/manage-user-access#configure-member-roles). - You must provide the correct API key in every request. Otherwise, TiDB Cloud responds with a `401` error.  ## API key management  ### Create an API key  Only the **owner** of an organization can create an API key.  To create an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **Create API Key**. 4. Enter a description for your API key. 5. Configure the role and scope for the API key. For more information about the permissions of a role, see [User roles](https://docs.pingcap.com/tidbcloud/manage-user-access/#user-roles). 6. Click **Generate API Key**. Copy and save the public key and the private key. 7. Make sure that you have copied and saved the private key in a secure location. The private key only displays upon the creation. After leaving this page, you will not be able to get the full private key again. 8. Click **Done**.  ### View details of an API key  To view details of an API key, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. You can view the details of the API keys on the page.  ### Edit an API key  Only the **owner** of an organization can modify an API key.  To edit an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to change, and then click **Update Role**. 4. You can update the description and role of the API key. 5. Click **Update**.  ### Delete an API key  Only the **owner** of an organization can delete an API key.  To delete an API key in an organization, perform the following steps:  1. In the [TiDB Cloud console](https://tidbcloud.com), switch to your target organization using the combo box in the upper-left corner. 2. In the left navigation pane, click **Organization Settings** > **API Keys**. 3. On the **API Keys** page, click **...** in the API key row that you want to delete, and then click **Delete**. 4. Click **I understand, delete it.**  # API Changelog  This changelog lists all changes to the TiDB Cloud Premium API (v1beta2).  <!-- In reverse chronological order -->  ## 20260924  - Published CMEK IAM principal retrieval and verification APIs for Premium instances on AWS and Alibaba Cloud.  ## 20260923  - Added APIs to get and update public connection settings for Premium and Essential V2 instances.  ## 20260428  - Initial release of the TiDB Cloud Premium API (v1beta2), including the following resources and endpoints:  * TiDB Cloud Premium instance   * [List TiDB Cloud Premium instances](#tag/TiDB-Instance/operation/TidbService_ListTidbs)   * [Create a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_CreateTidb)   * [Get a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetTidb)   * [Delete a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_DeleteTidb)   * [Update a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_UpdateTidb)   * [Reset the root password of a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_ResetRootPassword)   * [Get cloud provider information for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCloudProviderInfo)   * [Get the CA certificate download URL for a TiDB Cloud Premium instance](#tag/TiDB-Instance/operation/TidbService_GetCaCertificateDownloadUrl)  * Backup   * [List backups for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_ListTidbBackups)   * [Delete a backup for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_DeleteTidbBackup)   * [Restore a TiDB Cloud Premium instance from a backup](#tag/Backup/operation/TidbService_RestoreTidb)   * [Get restore status for a TiDB Cloud Premium instance](#tag/Backup/operation/TidbService_GetRestoreStatus)  * Region   * [List regions](#tag/Region/operation/RegionService_ListRegions)

API version: v1beta2
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package nextgen

import (
	"encoding/json"
	"time"
)

// checks if the V1beta2Backup type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta2Backup{}

// V1beta2Backup struct for V1beta2Backup
type V1beta2Backup struct {
	// The unique identifier of the backup.
	Id *string `json:"id,omitempty"`
	// The ID of the TiDB Cloud Premium instance to which the backup belongs.
	TidbId *string `json:"tidbId,omitempty"`
	// The ID of the organization to which the backup belongs.
	OrganizationId *string `json:"organizationId,omitempty"`
	// The service plan of the backup. The possible value is `Premium`.
	ServicePlan *V1beta1ServicePlan `json:"servicePlan,omitempty"`
	// The name of the backup.
	DisplayName *string `json:"displayName,omitempty"`
	// The description of the backup.
	Description *string `json:"description,omitempty"`
	// The timestamp when the backup was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// The size of the backup in bytes.
	SizeBytes *string `json:"sizeBytes,omitempty"`
	// The current state of the backup.
	State *V1beta2BackupState `json:"state,omitempty"`
	// The type of the backup.
	Type *V1beta2BackupType `json:"type,omitempty"`
	// The trigger type of the backup.
	TriggerType *BackupTriggerType `json:"triggerType,omitempty"`
	// The timestamp when the backup expires, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	ExpirationTime *time.Time `json:"expirationTime,omitempty"`
	// The backup timestamp.
	BackupTs *string `json:"backupTs,omitempty"`
	// The ID of the region where the backup is stored.
	RegionId *string `json:"regionId,omitempty"`
	// Whether the source cluster used a customer-managed key at backup time. Combine with encryption_key_id: \"0\" means not encrypted; non-zero key id with false means SMEK; non-zero key id with true means CMEK.
	CustomerManaged *bool `json:"customerManaged,omitempty"`
	// Source cluster encryption key ID at backup time. \"0\" means no encryption key was used.
	EncryptionKeyId      *string `json:"encryptionKeyId,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta2Backup V1beta2Backup

// NewV1beta2Backup instantiates a new V1beta2Backup object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta2Backup() *V1beta2Backup {
	this := V1beta2Backup{}
	return &this
}

// NewV1beta2BackupWithDefaults instantiates a new V1beta2Backup object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta2BackupWithDefaults() *V1beta2Backup {
	this := V1beta2Backup{}
	return &this
}

// GetId returns the Id field value if set, zero value otherwise.
func (o *V1beta2Backup) GetId() string {
	if o == nil || IsNil(o.Id) {
		var ret string
		return ret
	}
	return *o.Id
}

// GetIdOk returns a tuple with the Id field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetIdOk() (*string, bool) {
	if o == nil || IsNil(o.Id) {
		return nil, false
	}
	return o.Id, true
}

// HasId returns a boolean if a field has been set.
func (o *V1beta2Backup) HasId() bool {
	if o != nil && !IsNil(o.Id) {
		return true
	}

	return false
}

// SetId gets a reference to the given string and assigns it to the Id field.
func (o *V1beta2Backup) SetId(v string) {
	o.Id = &v
}

// GetTidbId returns the TidbId field value if set, zero value otherwise.
func (o *V1beta2Backup) GetTidbId() string {
	if o == nil || IsNil(o.TidbId) {
		var ret string
		return ret
	}
	return *o.TidbId
}

// GetTidbIdOk returns a tuple with the TidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbId) {
		return nil, false
	}
	return o.TidbId, true
}

// HasTidbId returns a boolean if a field has been set.
func (o *V1beta2Backup) HasTidbId() bool {
	if o != nil && !IsNil(o.TidbId) {
		return true
	}

	return false
}

// SetTidbId gets a reference to the given string and assigns it to the TidbId field.
func (o *V1beta2Backup) SetTidbId(v string) {
	o.TidbId = &v
}

// GetOrganizationId returns the OrganizationId field value if set, zero value otherwise.
func (o *V1beta2Backup) GetOrganizationId() string {
	if o == nil || IsNil(o.OrganizationId) {
		var ret string
		return ret
	}
	return *o.OrganizationId
}

// GetOrganizationIdOk returns a tuple with the OrganizationId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetOrganizationIdOk() (*string, bool) {
	if o == nil || IsNil(o.OrganizationId) {
		return nil, false
	}
	return o.OrganizationId, true
}

// HasOrganizationId returns a boolean if a field has been set.
func (o *V1beta2Backup) HasOrganizationId() bool {
	if o != nil && !IsNil(o.OrganizationId) {
		return true
	}

	return false
}

// SetOrganizationId gets a reference to the given string and assigns it to the OrganizationId field.
func (o *V1beta2Backup) SetOrganizationId(v string) {
	o.OrganizationId = &v
}

// GetServicePlan returns the ServicePlan field value if set, zero value otherwise.
func (o *V1beta2Backup) GetServicePlan() V1beta1ServicePlan {
	if o == nil || IsNil(o.ServicePlan) {
		var ret V1beta1ServicePlan
		return ret
	}
	return *o.ServicePlan
}

// GetServicePlanOk returns a tuple with the ServicePlan field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetServicePlanOk() (*V1beta1ServicePlan, bool) {
	if o == nil || IsNil(o.ServicePlan) {
		return nil, false
	}
	return o.ServicePlan, true
}

// HasServicePlan returns a boolean if a field has been set.
func (o *V1beta2Backup) HasServicePlan() bool {
	if o != nil && !IsNil(o.ServicePlan) {
		return true
	}

	return false
}

// SetServicePlan gets a reference to the given V1beta1ServicePlan and assigns it to the ServicePlan field.
func (o *V1beta2Backup) SetServicePlan(v V1beta1ServicePlan) {
	o.ServicePlan = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *V1beta2Backup) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *V1beta2Backup) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *V1beta2Backup) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetDescription returns the Description field value if set, zero value otherwise.
func (o *V1beta2Backup) GetDescription() string {
	if o == nil || IsNil(o.Description) {
		var ret string
		return ret
	}
	return *o.Description
}

// GetDescriptionOk returns a tuple with the Description field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetDescriptionOk() (*string, bool) {
	if o == nil || IsNil(o.Description) {
		return nil, false
	}
	return o.Description, true
}

// HasDescription returns a boolean if a field has been set.
func (o *V1beta2Backup) HasDescription() bool {
	if o != nil && !IsNil(o.Description) {
		return true
	}

	return false
}

// SetDescription gets a reference to the given string and assigns it to the Description field.
func (o *V1beta2Backup) SetDescription(v string) {
	o.Description = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *V1beta2Backup) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *V1beta2Backup) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *V1beta2Backup) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetSizeBytes returns the SizeBytes field value if set, zero value otherwise.
func (o *V1beta2Backup) GetSizeBytes() string {
	if o == nil || IsNil(o.SizeBytes) {
		var ret string
		return ret
	}
	return *o.SizeBytes
}

// GetSizeBytesOk returns a tuple with the SizeBytes field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetSizeBytesOk() (*string, bool) {
	if o == nil || IsNil(o.SizeBytes) {
		return nil, false
	}
	return o.SizeBytes, true
}

// HasSizeBytes returns a boolean if a field has been set.
func (o *V1beta2Backup) HasSizeBytes() bool {
	if o != nil && !IsNil(o.SizeBytes) {
		return true
	}

	return false
}

// SetSizeBytes gets a reference to the given string and assigns it to the SizeBytes field.
func (o *V1beta2Backup) SetSizeBytes(v string) {
	o.SizeBytes = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *V1beta2Backup) GetState() V1beta2BackupState {
	if o == nil || IsNil(o.State) {
		var ret V1beta2BackupState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetStateOk() (*V1beta2BackupState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *V1beta2Backup) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta2BackupState and assigns it to the State field.
func (o *V1beta2Backup) SetState(v V1beta2BackupState) {
	o.State = &v
}

// GetType returns the Type field value if set, zero value otherwise.
func (o *V1beta2Backup) GetType() V1beta2BackupType {
	if o == nil || IsNil(o.Type) {
		var ret V1beta2BackupType
		return ret
	}
	return *o.Type
}

// GetTypeOk returns a tuple with the Type field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetTypeOk() (*V1beta2BackupType, bool) {
	if o == nil || IsNil(o.Type) {
		return nil, false
	}
	return o.Type, true
}

// HasType returns a boolean if a field has been set.
func (o *V1beta2Backup) HasType() bool {
	if o != nil && !IsNil(o.Type) {
		return true
	}

	return false
}

// SetType gets a reference to the given V1beta2BackupType and assigns it to the Type field.
func (o *V1beta2Backup) SetType(v V1beta2BackupType) {
	o.Type = &v
}

// GetTriggerType returns the TriggerType field value if set, zero value otherwise.
func (o *V1beta2Backup) GetTriggerType() BackupTriggerType {
	if o == nil || IsNil(o.TriggerType) {
		var ret BackupTriggerType
		return ret
	}
	return *o.TriggerType
}

// GetTriggerTypeOk returns a tuple with the TriggerType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetTriggerTypeOk() (*BackupTriggerType, bool) {
	if o == nil || IsNil(o.TriggerType) {
		return nil, false
	}
	return o.TriggerType, true
}

// HasTriggerType returns a boolean if a field has been set.
func (o *V1beta2Backup) HasTriggerType() bool {
	if o != nil && !IsNil(o.TriggerType) {
		return true
	}

	return false
}

// SetTriggerType gets a reference to the given BackupTriggerType and assigns it to the TriggerType field.
func (o *V1beta2Backup) SetTriggerType(v BackupTriggerType) {
	o.TriggerType = &v
}

// GetExpirationTime returns the ExpirationTime field value if set, zero value otherwise.
func (o *V1beta2Backup) GetExpirationTime() time.Time {
	if o == nil || IsNil(o.ExpirationTime) {
		var ret time.Time
		return ret
	}
	return *o.ExpirationTime
}

// GetExpirationTimeOk returns a tuple with the ExpirationTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetExpirationTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.ExpirationTime) {
		return nil, false
	}
	return o.ExpirationTime, true
}

// HasExpirationTime returns a boolean if a field has been set.
func (o *V1beta2Backup) HasExpirationTime() bool {
	if o != nil && !IsNil(o.ExpirationTime) {
		return true
	}

	return false
}

// SetExpirationTime gets a reference to the given time.Time and assigns it to the ExpirationTime field.
func (o *V1beta2Backup) SetExpirationTime(v time.Time) {
	o.ExpirationTime = &v
}

// GetBackupTs returns the BackupTs field value if set, zero value otherwise.
func (o *V1beta2Backup) GetBackupTs() string {
	if o == nil || IsNil(o.BackupTs) {
		var ret string
		return ret
	}
	return *o.BackupTs
}

// GetBackupTsOk returns a tuple with the BackupTs field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetBackupTsOk() (*string, bool) {
	if o == nil || IsNil(o.BackupTs) {
		return nil, false
	}
	return o.BackupTs, true
}

// HasBackupTs returns a boolean if a field has been set.
func (o *V1beta2Backup) HasBackupTs() bool {
	if o != nil && !IsNil(o.BackupTs) {
		return true
	}

	return false
}

// SetBackupTs gets a reference to the given string and assigns it to the BackupTs field.
func (o *V1beta2Backup) SetBackupTs(v string) {
	o.BackupTs = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *V1beta2Backup) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *V1beta2Backup) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *V1beta2Backup) SetRegionId(v string) {
	o.RegionId = &v
}

// GetCustomerManaged returns the CustomerManaged field value if set, zero value otherwise.
func (o *V1beta2Backup) GetCustomerManaged() bool {
	if o == nil || IsNil(o.CustomerManaged) {
		var ret bool
		return ret
	}
	return *o.CustomerManaged
}

// GetCustomerManagedOk returns a tuple with the CustomerManaged field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetCustomerManagedOk() (*bool, bool) {
	if o == nil || IsNil(o.CustomerManaged) {
		return nil, false
	}
	return o.CustomerManaged, true
}

// HasCustomerManaged returns a boolean if a field has been set.
func (o *V1beta2Backup) HasCustomerManaged() bool {
	if o != nil && !IsNil(o.CustomerManaged) {
		return true
	}

	return false
}

// SetCustomerManaged gets a reference to the given bool and assigns it to the CustomerManaged field.
func (o *V1beta2Backup) SetCustomerManaged(v bool) {
	o.CustomerManaged = &v
}

// GetEncryptionKeyId returns the EncryptionKeyId field value if set, zero value otherwise.
func (o *V1beta2Backup) GetEncryptionKeyId() string {
	if o == nil || IsNil(o.EncryptionKeyId) {
		var ret string
		return ret
	}
	return *o.EncryptionKeyId
}

// GetEncryptionKeyIdOk returns a tuple with the EncryptionKeyId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta2Backup) GetEncryptionKeyIdOk() (*string, bool) {
	if o == nil || IsNil(o.EncryptionKeyId) {
		return nil, false
	}
	return o.EncryptionKeyId, true
}

// HasEncryptionKeyId returns a boolean if a field has been set.
func (o *V1beta2Backup) HasEncryptionKeyId() bool {
	if o != nil && !IsNil(o.EncryptionKeyId) {
		return true
	}

	return false
}

// SetEncryptionKeyId gets a reference to the given string and assigns it to the EncryptionKeyId field.
func (o *V1beta2Backup) SetEncryptionKeyId(v string) {
	o.EncryptionKeyId = &v
}

func (o V1beta2Backup) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta2Backup) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Id) {
		toSerialize["id"] = o.Id
	}
	if !IsNil(o.TidbId) {
		toSerialize["tidbId"] = o.TidbId
	}
	if !IsNil(o.OrganizationId) {
		toSerialize["organizationId"] = o.OrganizationId
	}
	if !IsNil(o.ServicePlan) {
		toSerialize["servicePlan"] = o.ServicePlan
	}
	if !IsNil(o.DisplayName) {
		toSerialize["displayName"] = o.DisplayName
	}
	if !IsNil(o.Description) {
		toSerialize["description"] = o.Description
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.SizeBytes) {
		toSerialize["sizeBytes"] = o.SizeBytes
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Type) {
		toSerialize["type"] = o.Type
	}
	if !IsNil(o.TriggerType) {
		toSerialize["triggerType"] = o.TriggerType
	}
	if !IsNil(o.ExpirationTime) {
		toSerialize["expirationTime"] = o.ExpirationTime
	}
	if !IsNil(o.BackupTs) {
		toSerialize["backupTs"] = o.BackupTs
	}
	if !IsNil(o.RegionId) {
		toSerialize["regionId"] = o.RegionId
	}
	if !IsNil(o.CustomerManaged) {
		toSerialize["customerManaged"] = o.CustomerManaged
	}
	if !IsNil(o.EncryptionKeyId) {
		toSerialize["encryptionKeyId"] = o.EncryptionKeyId
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta2Backup) UnmarshalJSON(data []byte) (err error) {
	varV1beta2Backup := _V1beta2Backup{}

	err = json.Unmarshal(data, &varV1beta2Backup)

	if err != nil {
		return err
	}

	*o = V1beta2Backup(varV1beta2Backup)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "id")
		delete(additionalProperties, "tidbId")
		delete(additionalProperties, "organizationId")
		delete(additionalProperties, "servicePlan")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "description")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "sizeBytes")
		delete(additionalProperties, "state")
		delete(additionalProperties, "type")
		delete(additionalProperties, "triggerType")
		delete(additionalProperties, "expirationTime")
		delete(additionalProperties, "backupTs")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "customerManaged")
		delete(additionalProperties, "encryptionKeyId")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta2Backup struct {
	value *V1beta2Backup
	isSet bool
}

func (v NullableV1beta2Backup) Get() *V1beta2Backup {
	return v.value
}

func (v *NullableV1beta2Backup) Set(val *V1beta2Backup) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta2Backup) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta2Backup) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta2Backup(val *V1beta2Backup) *NullableV1beta2Backup {
	return &NullableV1beta2Backup{value: val, isSet: true}
}

func (v NullableV1beta2Backup) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta2Backup) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
