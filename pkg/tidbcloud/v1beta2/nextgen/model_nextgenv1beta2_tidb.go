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

// checks if the Nextgenv1beta2Tidb type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Nextgenv1beta2Tidb{}

// Nextgenv1beta2Tidb struct for Nextgenv1beta2Tidb
type Nextgenv1beta2Tidb struct {
	// The unique identifier for the TiDB Cloud Premium instance, which is generated by the API and follows the format `tidbs/{tidbId}`.
	Name *string `json:"name,omitempty"`
	// The ID of the TiDB Cloud Premium instance.
	TidbId *string `json:"tidbId,omitempty"`
	// The user-defined name of the TiDB Cloud Premium instance.
	DisplayName string `json:"displayName" validate:"regexp=^[A-Za-z0-9][-A-Za-z0-9]{2,62}[A-Za-z0-9]$"`
	// The unique identifier of the region where the TiDB Cloud Premium instance is deployed, in the format of `{cloud_provider}-{region_code}`. For example, `aws-us-west-2`.
	RegionId string `json:"regionId"`
	// The cloud provider where the TiDB Cloud Premium instance is deployed.  - `\"aws\"`: Amazon Web Services  - `\"gcp\"`: Google Cloud  - `\"azure\"`: Microsoft Azure  - `\"alicloud\"`: Alibaba Cloud
	CloudProvider *RegionCloudProvider `json:"cloudProvider,omitempty"`
	// The display name of the region where the TiDB Cloud Premium instance is deployed. For example, `Oregon (us-west-2)`.
	RegionDisplayName *string `json:"regionDisplayName,omitempty"`
	// The current state of the TiDB Cloud Premium instance.
	State *V1beta1ClusterState `json:"state,omitempty"`
	// The root password of the TiDB Cloud Premium instance.  This field supports two input formats: - Plaintext password (legacy behavior) - RSA-OAEP-SHA256 encrypted payload prefixed with `rsa_oaep_sha256:`  If the marker prefix is present, the server parses and decrypts the encrypted payload. Any parse/decrypt failure returns an explicit parameter error without plaintext fallback. If the marker prefix is absent, the value is treated as plaintext for backward compatibility.  For plaintext input, the password must be between 8 and 64 characters long and can contain letters, numbers, and special characters.
	RootPassword *string `json:"rootPassword,omitempty" validate:"regexp=^(rsa_oaep_sha256:.+|.{8,64})$"`
	// The minimum number of Request Capacity Units (RCUs) for the TiDB Cloud Premium instance.  This field is read-only and is calculated by the service. It is omitted in Elastic Mode, where Baseline RCU is the provisioning and billing floor.
	MinRcu NullableString `json:"minRcu,omitempty"`
	// The customer-defined maximum number of Request Capacity Units (RCUs). Requirements depend on service_plan: - Premium and Essential_V2: optional. A positive value selects Max RCU Compatibility Mode. Omitting both   max_rcu and baseline_rcu selects Elastic Mode with the plan default Baseline RCU: 5,000 for Premium and   2,000 for Essential_V2. - BYOC: required and must be greater than zero. - Premium_Reserved: must be omitted. Zero and JSON null are treated as omission and are not clear operations.
	MaxRcu NullableString `json:"maxRcu,omitempty"`
	// The plan of the service.  - `Premium`: [TiDB Cloud Premium](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#tidb-cloud-premium)  - `Premium_Reserved`: TiDB Cloud Premium Reserved
	ServicePlan V1beta1ServicePlan `json:"servicePlan"`
	// The capacity template ID for Premium Reserved. This field is only accepted when `service_plan` is `Premium_Reserved`; it is not returned in responses.
	ReservedCapacityTemplateId *string `json:"reservedCapacityTemplateId,omitempty"`
	// The current Reserved capacity information. This field is only returned for Premium Reserved instances.
	ReservedCapacity *V1beta2ReservedCapacity `json:"reservedCapacity,omitempty"`
	// The resource pool selected for this TiDB Cloud instance. Currently, this field is only used for TiDB Cloud BYOC instances.
	ResourcePoolId NullableString `json:"resourcePoolId,omitempty"`
	// The current display name of the selected resource pool. Currently, this field is only returned for TiDB Cloud BYOC instances.
	ResourcePoolDisplayName NullableString `json:"resourcePoolDisplayName,omitempty"`
	// The high availability configuration for the TiDB Cloud instance. - `REGIONAL`: high availability across multiple availability zones within a region. For more information, see [Regional high availability architecture](https://docs.pingcap.com/tidbcloud/serverless-high-availability/#regional-high-availability-architecture). - `ZONAL`: high availability within a single availability zone. For more information, see [Zonal high availability architecture](https://docs.pingcap.com/tidbcloud/serverless-high-availability/#zonal-high-availability-architecture). TiDB Cloud Premium and Premium Reserved support only `REGIONAL`. If omitted for TiDB Cloud Essential, the selected resource pool determines the value.
	HighAvailabilityType *V1beta2TidbHighAvailabilityType `json:"highAvailabilityType,omitempty"`
	// The annotations for the TiDB Cloud Premium instance. The following lists some predefined annotations: - `tidb.cloud/has-set-password`: indicates whether the TiDB Cloud Premium instance has a root password set. - `tidb.cloud/available-features`: lists available features of the TiDB Cloud Premium instance.
	Annotations *map[string]string `json:"annotations,omitempty"`
	// The labels for the TiDB Cloud Premium instance. - `tidb.cloud/organization`: the ID of the organization where the TiDB Cloud Premium instance belongs. - `tidb.cloud/project`: the ID of the project where the TiDB Cloud Premium instance belongs.
	Labels *map[string]string `json:"labels,omitempty"`
	// The email address or public API key of the user who created the TiDB Cloud Premium instance.
	Creator *string `json:"creator,omitempty"`
	// The timestamp when the TiDB Cloud Premium instance was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// The timestamp when the TiDB Cloud Premium instance was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The scheduled time when the Resource Pool will automatically resume this instance's pool. This field is set only while the Resource Pool is paused.
	AutomaticResumeAt *time.Time `json:"automaticResumeAt,omitempty"`
	// The connection endpoints for accessing the TiDB Cloud Premium instance.
	Endpoints []TidbEndpoint `json:"endpoints,omitempty"`
	// The dual-layer data encryption configuration for the TiDB Cloud Premium instance.
	DualLayerDataEncryption *TidbDualLayerDataEncryption `json:"dualLayerDataEncryption,omitempty"`
	// The version of the TiDB instance.
	TidbVersion *string  `json:"tidbVersion,omitempty"`
	Tags        []string `json:"tags,omitempty"`
	// The timestamp when the capacity of the TiDB Cloud Premium instance was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CapacityUpdateTime *time.Time `json:"capacityUpdateTime,omitempty"`
	// The optional guaranteed provisioning commitment in Elastic Mode. Premium accepts 5,000, 12,500, or 50,000 RCU. Essential_V2 uses a fixed 2,000 RCU baseline. When neither max_rcu nor baseline_rcu is positive during creation, the service uses the plan default. Zero and JSON null are treated as omission. This field cannot be configured to a positive value together with a positive max_rcu.
	BaselineRcu NullableString `json:"baselineRcu,omitempty"`
	// The effective customer-facing capacity mode for Premium and Essential_V2. Other service plans, or capacity records that do not identify exactly one supported mode, return CAPACITY_MODE_UNSPECIFIED.
	CapacityMode         *TidbCapacityMode `json:"capacityMode,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Nextgenv1beta2Tidb Nextgenv1beta2Tidb

// NewNextgenv1beta2Tidb instantiates a new Nextgenv1beta2Tidb object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewNextgenv1beta2Tidb(displayName string, regionId string, servicePlan V1beta1ServicePlan) *Nextgenv1beta2Tidb {
	this := Nextgenv1beta2Tidb{}
	this.DisplayName = displayName
	this.RegionId = regionId
	this.ServicePlan = servicePlan
	return &this
}

// NewNextgenv1beta2TidbWithDefaults instantiates a new Nextgenv1beta2Tidb object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewNextgenv1beta2TidbWithDefaults() *Nextgenv1beta2Tidb {
	this := Nextgenv1beta2Tidb{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Nextgenv1beta2Tidb) SetName(v string) {
	o.Name = &v
}

// GetTidbId returns the TidbId field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetTidbId() string {
	if o == nil || IsNil(o.TidbId) {
		var ret string
		return ret
	}
	return *o.TidbId
}

// GetTidbIdOk returns a tuple with the TidbId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetTidbIdOk() (*string, bool) {
	if o == nil || IsNil(o.TidbId) {
		return nil, false
	}
	return o.TidbId, true
}

// HasTidbId returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasTidbId() bool {
	if o != nil && !IsNil(o.TidbId) {
		return true
	}

	return false
}

// SetTidbId gets a reference to the given string and assigns it to the TidbId field.
func (o *Nextgenv1beta2Tidb) SetTidbId(v string) {
	o.TidbId = &v
}

// GetDisplayName returns the DisplayName field value
func (o *Nextgenv1beta2Tidb) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *Nextgenv1beta2Tidb) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetRegionId returns the RegionId field value
func (o *Nextgenv1beta2Tidb) GetRegionId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetRegionIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.RegionId, true
}

// SetRegionId sets field value
func (o *Nextgenv1beta2Tidb) SetRegionId(v string) {
	o.RegionId = v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetCloudProvider() RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetCloudProviderOk() (*RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Nextgenv1beta2Tidb) SetCloudProvider(v RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetRegionDisplayName returns the RegionDisplayName field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetRegionDisplayName() string {
	if o == nil || IsNil(o.RegionDisplayName) {
		var ret string
		return ret
	}
	return *o.RegionDisplayName
}

// GetRegionDisplayNameOk returns a tuple with the RegionDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetRegionDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.RegionDisplayName) {
		return nil, false
	}
	return o.RegionDisplayName, true
}

// HasRegionDisplayName returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasRegionDisplayName() bool {
	if o != nil && !IsNil(o.RegionDisplayName) {
		return true
	}

	return false
}

// SetRegionDisplayName gets a reference to the given string and assigns it to the RegionDisplayName field.
func (o *Nextgenv1beta2Tidb) SetRegionDisplayName(v string) {
	o.RegionDisplayName = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetState() V1beta1ClusterState {
	if o == nil || IsNil(o.State) {
		var ret V1beta1ClusterState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetStateOk() (*V1beta1ClusterState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given V1beta1ClusterState and assigns it to the State field.
func (o *Nextgenv1beta2Tidb) SetState(v V1beta1ClusterState) {
	o.State = &v
}

// GetRootPassword returns the RootPassword field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetRootPassword() string {
	if o == nil || IsNil(o.RootPassword) {
		var ret string
		return ret
	}
	return *o.RootPassword
}

// GetRootPasswordOk returns a tuple with the RootPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetRootPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.RootPassword) {
		return nil, false
	}
	return o.RootPassword, true
}

// HasRootPassword returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasRootPassword() bool {
	if o != nil && !IsNil(o.RootPassword) {
		return true
	}

	return false
}

// SetRootPassword gets a reference to the given string and assigns it to the RootPassword field.
func (o *Nextgenv1beta2Tidb) SetRootPassword(v string) {
	o.RootPassword = &v
}

// GetMinRcu returns the MinRcu field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Nextgenv1beta2Tidb) GetMinRcu() string {
	if o == nil || IsNil(o.MinRcu.Get()) {
		var ret string
		return ret
	}
	return *o.MinRcu.Get()
}

// GetMinRcuOk returns a tuple with the MinRcu field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Nextgenv1beta2Tidb) GetMinRcuOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.MinRcu.Get(), o.MinRcu.IsSet()
}

// HasMinRcu returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasMinRcu() bool {
	if o != nil && o.MinRcu.IsSet() {
		return true
	}

	return false
}

// SetMinRcu gets a reference to the given NullableString and assigns it to the MinRcu field.
func (o *Nextgenv1beta2Tidb) SetMinRcu(v string) {
	o.MinRcu.Set(&v)
}

// SetMinRcuNil sets the value for MinRcu to be an explicit nil
func (o *Nextgenv1beta2Tidb) SetMinRcuNil() {
	o.MinRcu.Set(nil)
}

// UnsetMinRcu ensures that no value is present for MinRcu, not even an explicit nil
func (o *Nextgenv1beta2Tidb) UnsetMinRcu() {
	o.MinRcu.Unset()
}

// GetMaxRcu returns the MaxRcu field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Nextgenv1beta2Tidb) GetMaxRcu() string {
	if o == nil || IsNil(o.MaxRcu.Get()) {
		var ret string
		return ret
	}
	return *o.MaxRcu.Get()
}

// GetMaxRcuOk returns a tuple with the MaxRcu field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Nextgenv1beta2Tidb) GetMaxRcuOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.MaxRcu.Get(), o.MaxRcu.IsSet()
}

// HasMaxRcu returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasMaxRcu() bool {
	if o != nil && o.MaxRcu.IsSet() {
		return true
	}

	return false
}

// SetMaxRcu gets a reference to the given NullableString and assigns it to the MaxRcu field.
func (o *Nextgenv1beta2Tidb) SetMaxRcu(v string) {
	o.MaxRcu.Set(&v)
}

// SetMaxRcuNil sets the value for MaxRcu to be an explicit nil
func (o *Nextgenv1beta2Tidb) SetMaxRcuNil() {
	o.MaxRcu.Set(nil)
}

// UnsetMaxRcu ensures that no value is present for MaxRcu, not even an explicit nil
func (o *Nextgenv1beta2Tidb) UnsetMaxRcu() {
	o.MaxRcu.Unset()
}

// GetServicePlan returns the ServicePlan field value
func (o *Nextgenv1beta2Tidb) GetServicePlan() V1beta1ServicePlan {
	if o == nil {
		var ret V1beta1ServicePlan
		return ret
	}

	return o.ServicePlan
}

// GetServicePlanOk returns a tuple with the ServicePlan field value
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetServicePlanOk() (*V1beta1ServicePlan, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ServicePlan, true
}

// SetServicePlan sets field value
func (o *Nextgenv1beta2Tidb) SetServicePlan(v V1beta1ServicePlan) {
	o.ServicePlan = v
}

// GetReservedCapacityTemplateId returns the ReservedCapacityTemplateId field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetReservedCapacityTemplateId() string {
	if o == nil || IsNil(o.ReservedCapacityTemplateId) {
		var ret string
		return ret
	}
	return *o.ReservedCapacityTemplateId
}

// GetReservedCapacityTemplateIdOk returns a tuple with the ReservedCapacityTemplateId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetReservedCapacityTemplateIdOk() (*string, bool) {
	if o == nil || IsNil(o.ReservedCapacityTemplateId) {
		return nil, false
	}
	return o.ReservedCapacityTemplateId, true
}

// HasReservedCapacityTemplateId returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasReservedCapacityTemplateId() bool {
	if o != nil && !IsNil(o.ReservedCapacityTemplateId) {
		return true
	}

	return false
}

// SetReservedCapacityTemplateId gets a reference to the given string and assigns it to the ReservedCapacityTemplateId field.
func (o *Nextgenv1beta2Tidb) SetReservedCapacityTemplateId(v string) {
	o.ReservedCapacityTemplateId = &v
}

// GetReservedCapacity returns the ReservedCapacity field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetReservedCapacity() V1beta2ReservedCapacity {
	if o == nil || IsNil(o.ReservedCapacity) {
		var ret V1beta2ReservedCapacity
		return ret
	}
	return *o.ReservedCapacity
}

// GetReservedCapacityOk returns a tuple with the ReservedCapacity field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetReservedCapacityOk() (*V1beta2ReservedCapacity, bool) {
	if o == nil || IsNil(o.ReservedCapacity) {
		return nil, false
	}
	return o.ReservedCapacity, true
}

// HasReservedCapacity returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasReservedCapacity() bool {
	if o != nil && !IsNil(o.ReservedCapacity) {
		return true
	}

	return false
}

// SetReservedCapacity gets a reference to the given V1beta2ReservedCapacity and assigns it to the ReservedCapacity field.
func (o *Nextgenv1beta2Tidb) SetReservedCapacity(v V1beta2ReservedCapacity) {
	o.ReservedCapacity = &v
}

// GetResourcePoolId returns the ResourcePoolId field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Nextgenv1beta2Tidb) GetResourcePoolId() string {
	if o == nil || IsNil(o.ResourcePoolId.Get()) {
		var ret string
		return ret
	}
	return *o.ResourcePoolId.Get()
}

// GetResourcePoolIdOk returns a tuple with the ResourcePoolId field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Nextgenv1beta2Tidb) GetResourcePoolIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResourcePoolId.Get(), o.ResourcePoolId.IsSet()
}

// HasResourcePoolId returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasResourcePoolId() bool {
	if o != nil && o.ResourcePoolId.IsSet() {
		return true
	}

	return false
}

// SetResourcePoolId gets a reference to the given NullableString and assigns it to the ResourcePoolId field.
func (o *Nextgenv1beta2Tidb) SetResourcePoolId(v string) {
	o.ResourcePoolId.Set(&v)
}

// SetResourcePoolIdNil sets the value for ResourcePoolId to be an explicit nil
func (o *Nextgenv1beta2Tidb) SetResourcePoolIdNil() {
	o.ResourcePoolId.Set(nil)
}

// UnsetResourcePoolId ensures that no value is present for ResourcePoolId, not even an explicit nil
func (o *Nextgenv1beta2Tidb) UnsetResourcePoolId() {
	o.ResourcePoolId.Unset()
}

// GetResourcePoolDisplayName returns the ResourcePoolDisplayName field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Nextgenv1beta2Tidb) GetResourcePoolDisplayName() string {
	if o == nil || IsNil(o.ResourcePoolDisplayName.Get()) {
		var ret string
		return ret
	}
	return *o.ResourcePoolDisplayName.Get()
}

// GetResourcePoolDisplayNameOk returns a tuple with the ResourcePoolDisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Nextgenv1beta2Tidb) GetResourcePoolDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.ResourcePoolDisplayName.Get(), o.ResourcePoolDisplayName.IsSet()
}

// HasResourcePoolDisplayName returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasResourcePoolDisplayName() bool {
	if o != nil && o.ResourcePoolDisplayName.IsSet() {
		return true
	}

	return false
}

// SetResourcePoolDisplayName gets a reference to the given NullableString and assigns it to the ResourcePoolDisplayName field.
func (o *Nextgenv1beta2Tidb) SetResourcePoolDisplayName(v string) {
	o.ResourcePoolDisplayName.Set(&v)
}

// SetResourcePoolDisplayNameNil sets the value for ResourcePoolDisplayName to be an explicit nil
func (o *Nextgenv1beta2Tidb) SetResourcePoolDisplayNameNil() {
	o.ResourcePoolDisplayName.Set(nil)
}

// UnsetResourcePoolDisplayName ensures that no value is present for ResourcePoolDisplayName, not even an explicit nil
func (o *Nextgenv1beta2Tidb) UnsetResourcePoolDisplayName() {
	o.ResourcePoolDisplayName.Unset()
}

// GetHighAvailabilityType returns the HighAvailabilityType field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetHighAvailabilityType() V1beta2TidbHighAvailabilityType {
	if o == nil || IsNil(o.HighAvailabilityType) {
		var ret V1beta2TidbHighAvailabilityType
		return ret
	}
	return *o.HighAvailabilityType
}

// GetHighAvailabilityTypeOk returns a tuple with the HighAvailabilityType field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetHighAvailabilityTypeOk() (*V1beta2TidbHighAvailabilityType, bool) {
	if o == nil || IsNil(o.HighAvailabilityType) {
		return nil, false
	}
	return o.HighAvailabilityType, true
}

// HasHighAvailabilityType returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasHighAvailabilityType() bool {
	if o != nil && !IsNil(o.HighAvailabilityType) {
		return true
	}

	return false
}

// SetHighAvailabilityType gets a reference to the given V1beta2TidbHighAvailabilityType and assigns it to the HighAvailabilityType field.
func (o *Nextgenv1beta2Tidb) SetHighAvailabilityType(v V1beta2TidbHighAvailabilityType) {
	o.HighAvailabilityType = &v
}

// GetAnnotations returns the Annotations field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetAnnotations() map[string]string {
	if o == nil || IsNil(o.Annotations) {
		var ret map[string]string
		return ret
	}
	return *o.Annotations
}

// GetAnnotationsOk returns a tuple with the Annotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetAnnotationsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Annotations) {
		return nil, false
	}
	return o.Annotations, true
}

// HasAnnotations returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasAnnotations() bool {
	if o != nil && !IsNil(o.Annotations) {
		return true
	}

	return false
}

// SetAnnotations gets a reference to the given map[string]string and assigns it to the Annotations field.
func (o *Nextgenv1beta2Tidb) SetAnnotations(v map[string]string) {
	o.Annotations = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *Nextgenv1beta2Tidb) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetCreator returns the Creator field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetCreator() string {
	if o == nil || IsNil(o.Creator) {
		var ret string
		return ret
	}
	return *o.Creator
}

// GetCreatorOk returns a tuple with the Creator field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetCreatorOk() (*string, bool) {
	if o == nil || IsNil(o.Creator) {
		return nil, false
	}
	return o.Creator, true
}

// HasCreator returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasCreator() bool {
	if o != nil && !IsNil(o.Creator) {
		return true
	}

	return false
}

// SetCreator gets a reference to the given string and assigns it to the Creator field.
func (o *Nextgenv1beta2Tidb) SetCreator(v string) {
	o.Creator = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *Nextgenv1beta2Tidb) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetUpdateTime returns the UpdateTime field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetUpdateTime() time.Time {
	if o == nil || IsNil(o.UpdateTime) {
		var ret time.Time
		return ret
	}
	return *o.UpdateTime
}

// GetUpdateTimeOk returns a tuple with the UpdateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetUpdateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdateTime) {
		return nil, false
	}
	return o.UpdateTime, true
}

// HasUpdateTime returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasUpdateTime() bool {
	if o != nil && !IsNil(o.UpdateTime) {
		return true
	}

	return false
}

// SetUpdateTime gets a reference to the given time.Time and assigns it to the UpdateTime field.
func (o *Nextgenv1beta2Tidb) SetUpdateTime(v time.Time) {
	o.UpdateTime = &v
}

// GetAutomaticResumeAt returns the AutomaticResumeAt field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetAutomaticResumeAt() time.Time {
	if o == nil || IsNil(o.AutomaticResumeAt) {
		var ret time.Time
		return ret
	}
	return *o.AutomaticResumeAt
}

// GetAutomaticResumeAtOk returns a tuple with the AutomaticResumeAt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetAutomaticResumeAtOk() (*time.Time, bool) {
	if o == nil || IsNil(o.AutomaticResumeAt) {
		return nil, false
	}
	return o.AutomaticResumeAt, true
}

// HasAutomaticResumeAt returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasAutomaticResumeAt() bool {
	if o != nil && !IsNil(o.AutomaticResumeAt) {
		return true
	}

	return false
}

// SetAutomaticResumeAt gets a reference to the given time.Time and assigns it to the AutomaticResumeAt field.
func (o *Nextgenv1beta2Tidb) SetAutomaticResumeAt(v time.Time) {
	o.AutomaticResumeAt = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetEndpoints() []TidbEndpoint {
	if o == nil || IsNil(o.Endpoints) {
		var ret []TidbEndpoint
		return ret
	}
	return o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetEndpointsOk() ([]TidbEndpoint, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given []TidbEndpoint and assigns it to the Endpoints field.
func (o *Nextgenv1beta2Tidb) SetEndpoints(v []TidbEndpoint) {
	o.Endpoints = v
}

// GetDualLayerDataEncryption returns the DualLayerDataEncryption field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetDualLayerDataEncryption() TidbDualLayerDataEncryption {
	if o == nil || IsNil(o.DualLayerDataEncryption) {
		var ret TidbDualLayerDataEncryption
		return ret
	}
	return *o.DualLayerDataEncryption
}

// GetDualLayerDataEncryptionOk returns a tuple with the DualLayerDataEncryption field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetDualLayerDataEncryptionOk() (*TidbDualLayerDataEncryption, bool) {
	if o == nil || IsNil(o.DualLayerDataEncryption) {
		return nil, false
	}
	return o.DualLayerDataEncryption, true
}

// HasDualLayerDataEncryption returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasDualLayerDataEncryption() bool {
	if o != nil && !IsNil(o.DualLayerDataEncryption) {
		return true
	}

	return false
}

// SetDualLayerDataEncryption gets a reference to the given TidbDualLayerDataEncryption and assigns it to the DualLayerDataEncryption field.
func (o *Nextgenv1beta2Tidb) SetDualLayerDataEncryption(v TidbDualLayerDataEncryption) {
	o.DualLayerDataEncryption = &v
}

// GetTidbVersion returns the TidbVersion field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetTidbVersion() string {
	if o == nil || IsNil(o.TidbVersion) {
		var ret string
		return ret
	}
	return *o.TidbVersion
}

// GetTidbVersionOk returns a tuple with the TidbVersion field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetTidbVersionOk() (*string, bool) {
	if o == nil || IsNil(o.TidbVersion) {
		return nil, false
	}
	return o.TidbVersion, true
}

// HasTidbVersion returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasTidbVersion() bool {
	if o != nil && !IsNil(o.TidbVersion) {
		return true
	}

	return false
}

// SetTidbVersion gets a reference to the given string and assigns it to the TidbVersion field.
func (o *Nextgenv1beta2Tidb) SetTidbVersion(v string) {
	o.TidbVersion = &v
}

// GetTags returns the Tags field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetTags() []string {
	if o == nil || IsNil(o.Tags) {
		var ret []string
		return ret
	}
	return o.Tags
}

// GetTagsOk returns a tuple with the Tags field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetTagsOk() ([]string, bool) {
	if o == nil || IsNil(o.Tags) {
		return nil, false
	}
	return o.Tags, true
}

// HasTags returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasTags() bool {
	if o != nil && !IsNil(o.Tags) {
		return true
	}

	return false
}

// SetTags gets a reference to the given []string and assigns it to the Tags field.
func (o *Nextgenv1beta2Tidb) SetTags(v []string) {
	o.Tags = v
}

// GetCapacityUpdateTime returns the CapacityUpdateTime field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetCapacityUpdateTime() time.Time {
	if o == nil || IsNil(o.CapacityUpdateTime) {
		var ret time.Time
		return ret
	}
	return *o.CapacityUpdateTime
}

// GetCapacityUpdateTimeOk returns a tuple with the CapacityUpdateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetCapacityUpdateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CapacityUpdateTime) {
		return nil, false
	}
	return o.CapacityUpdateTime, true
}

// HasCapacityUpdateTime returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasCapacityUpdateTime() bool {
	if o != nil && !IsNil(o.CapacityUpdateTime) {
		return true
	}

	return false
}

// SetCapacityUpdateTime gets a reference to the given time.Time and assigns it to the CapacityUpdateTime field.
func (o *Nextgenv1beta2Tidb) SetCapacityUpdateTime(v time.Time) {
	o.CapacityUpdateTime = &v
}

// GetBaselineRcu returns the BaselineRcu field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Nextgenv1beta2Tidb) GetBaselineRcu() string {
	if o == nil || IsNil(o.BaselineRcu.Get()) {
		var ret string
		return ret
	}
	return *o.BaselineRcu.Get()
}

// GetBaselineRcuOk returns a tuple with the BaselineRcu field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Nextgenv1beta2Tidb) GetBaselineRcuOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.BaselineRcu.Get(), o.BaselineRcu.IsSet()
}

// HasBaselineRcu returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasBaselineRcu() bool {
	if o != nil && o.BaselineRcu.IsSet() {
		return true
	}

	return false
}

// SetBaselineRcu gets a reference to the given NullableString and assigns it to the BaselineRcu field.
func (o *Nextgenv1beta2Tidb) SetBaselineRcu(v string) {
	o.BaselineRcu.Set(&v)
}

// SetBaselineRcuNil sets the value for BaselineRcu to be an explicit nil
func (o *Nextgenv1beta2Tidb) SetBaselineRcuNil() {
	o.BaselineRcu.Set(nil)
}

// UnsetBaselineRcu ensures that no value is present for BaselineRcu, not even an explicit nil
func (o *Nextgenv1beta2Tidb) UnsetBaselineRcu() {
	o.BaselineRcu.Unset()
}

// GetCapacityMode returns the CapacityMode field value if set, zero value otherwise.
func (o *Nextgenv1beta2Tidb) GetCapacityMode() TidbCapacityMode {
	if o == nil || IsNil(o.CapacityMode) {
		var ret TidbCapacityMode
		return ret
	}
	return *o.CapacityMode
}

// GetCapacityModeOk returns a tuple with the CapacityMode field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Nextgenv1beta2Tidb) GetCapacityModeOk() (*TidbCapacityMode, bool) {
	if o == nil || IsNil(o.CapacityMode) {
		return nil, false
	}
	return o.CapacityMode, true
}

// HasCapacityMode returns a boolean if a field has been set.
func (o *Nextgenv1beta2Tidb) HasCapacityMode() bool {
	if o != nil && !IsNil(o.CapacityMode) {
		return true
	}

	return false
}

// SetCapacityMode gets a reference to the given TidbCapacityMode and assigns it to the CapacityMode field.
func (o *Nextgenv1beta2Tidb) SetCapacityMode(v TidbCapacityMode) {
	o.CapacityMode = &v
}

func (o Nextgenv1beta2Tidb) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Nextgenv1beta2Tidb) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.TidbId) {
		toSerialize["tidbId"] = o.TidbId
	}
	toSerialize["displayName"] = o.DisplayName
	toSerialize["regionId"] = o.RegionId
	if !IsNil(o.CloudProvider) {
		toSerialize["cloudProvider"] = o.CloudProvider
	}
	if !IsNil(o.RegionDisplayName) {
		toSerialize["regionDisplayName"] = o.RegionDisplayName
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.RootPassword) {
		toSerialize["rootPassword"] = o.RootPassword
	}
	if o.MinRcu.IsSet() {
		toSerialize["minRcu"] = o.MinRcu.Get()
	}
	if o.MaxRcu.IsSet() {
		toSerialize["maxRcu"] = o.MaxRcu.Get()
	}
	toSerialize["servicePlan"] = o.ServicePlan
	if !IsNil(o.ReservedCapacityTemplateId) {
		toSerialize["reservedCapacityTemplateId"] = o.ReservedCapacityTemplateId
	}
	if !IsNil(o.ReservedCapacity) {
		toSerialize["reservedCapacity"] = o.ReservedCapacity
	}
	if o.ResourcePoolId.IsSet() {
		toSerialize["resourcePoolId"] = o.ResourcePoolId.Get()
	}
	if o.ResourcePoolDisplayName.IsSet() {
		toSerialize["resourcePoolDisplayName"] = o.ResourcePoolDisplayName.Get()
	}
	if !IsNil(o.HighAvailabilityType) {
		toSerialize["highAvailabilityType"] = o.HighAvailabilityType
	}
	if !IsNil(o.Annotations) {
		toSerialize["annotations"] = o.Annotations
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.Creator) {
		toSerialize["creator"] = o.Creator
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.UpdateTime) {
		toSerialize["updateTime"] = o.UpdateTime
	}
	if !IsNil(o.AutomaticResumeAt) {
		toSerialize["automaticResumeAt"] = o.AutomaticResumeAt
	}
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.DualLayerDataEncryption) {
		toSerialize["dualLayerDataEncryption"] = o.DualLayerDataEncryption
	}
	if !IsNil(o.TidbVersion) {
		toSerialize["tidbVersion"] = o.TidbVersion
	}
	if !IsNil(o.Tags) {
		toSerialize["tags"] = o.Tags
	}
	if !IsNil(o.CapacityUpdateTime) {
		toSerialize["capacityUpdateTime"] = o.CapacityUpdateTime
	}
	if o.BaselineRcu.IsSet() {
		toSerialize["baselineRcu"] = o.BaselineRcu.Get()
	}
	if !IsNil(o.CapacityMode) {
		toSerialize["capacityMode"] = o.CapacityMode
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Nextgenv1beta2Tidb) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"displayName",
		"regionId",
		"servicePlan",
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

	varNextgenv1beta2Tidb := _Nextgenv1beta2Tidb{}

	err = json.Unmarshal(data, &varNextgenv1beta2Tidb)

	if err != nil {
		return err
	}

	*o = Nextgenv1beta2Tidb(varNextgenv1beta2Tidb)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "tidbId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "regionDisplayName")
		delete(additionalProperties, "state")
		delete(additionalProperties, "rootPassword")
		delete(additionalProperties, "minRcu")
		delete(additionalProperties, "maxRcu")
		delete(additionalProperties, "servicePlan")
		delete(additionalProperties, "reservedCapacityTemplateId")
		delete(additionalProperties, "reservedCapacity")
		delete(additionalProperties, "resourcePoolId")
		delete(additionalProperties, "resourcePoolDisplayName")
		delete(additionalProperties, "highAvailabilityType")
		delete(additionalProperties, "annotations")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "creator")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "updateTime")
		delete(additionalProperties, "automaticResumeAt")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "dualLayerDataEncryption")
		delete(additionalProperties, "tidbVersion")
		delete(additionalProperties, "tags")
		delete(additionalProperties, "capacityUpdateTime")
		delete(additionalProperties, "baselineRcu")
		delete(additionalProperties, "capacityMode")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableNextgenv1beta2Tidb struct {
	value *Nextgenv1beta2Tidb
	isSet bool
}

func (v NullableNextgenv1beta2Tidb) Get() *Nextgenv1beta2Tidb {
	return v.value
}

func (v *NullableNextgenv1beta2Tidb) Set(val *Nextgenv1beta2Tidb) {
	v.value = val
	v.isSet = true
}

func (v NullableNextgenv1beta2Tidb) IsSet() bool {
	return v.isSet
}

func (v *NullableNextgenv1beta2Tidb) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableNextgenv1beta2Tidb(val *Nextgenv1beta2Tidb) *NullableNextgenv1beta2Tidb {
	return &NullableNextgenv1beta2Tidb{value: val, isSet: true}
}

func (v NullableNextgenv1beta2Tidb) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableNextgenv1beta2Tidb) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
