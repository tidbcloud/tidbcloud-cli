/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
	"fmt"
	"time"
)

// checks if the TidbCloudOpenApiserverlessv1beta1Cluster type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &TidbCloudOpenApiserverlessv1beta1Cluster{}

// TidbCloudOpenApiserverlessv1beta1Cluster Message for a serverless TiDB cluster resource.
type TidbCloudOpenApiserverlessv1beta1Cluster struct {
	// Output_only. The unique name of the cluster.
	Name *string `json:"name,omitempty"`
	// Output_only. The unique ID of the cluster.
	ClusterId *string `json:"clusterId,omitempty"`
	// Required. User friendly display name of the cluster.
	DisplayName string `json:"displayName"`
	// Required. Region where the cluster will be created.
	Region Commonv1beta1Region `json:"region"`
	// Optional. The spending limit for the cluster.
	SpendingLimit *ClusterSpendingLimit `json:"spendingLimit,omitempty"`
	// Optional. Automated backup policy to set on the cluster.
	AutomatedBackupPolicy *V1beta1ClusterAutomatedBackupPolicy `json:"automatedBackupPolicy,omitempty"`
	// Optional. The endpoints for connecting to the cluster.
	Endpoints    *V1beta1ClusterEndpoints `json:"endpoints,omitempty"`
	RootPassword *string                  `json:"rootPassword,omitempty" validate:"regexp=^.{8,64}$"`
	// Optional. Encryption settings for the cluster.
	EncryptionConfig *V1beta1ClusterEncryptionConfig `json:"encryptionConfig,omitempty"`
	// Output_only. The TiDB version of the cluster.
	Version *string `json:"version,omitempty"`
	// Output_only. The email of the creator of the cluster.
	CreatedBy *string `json:"createdBy,omitempty"`
	// Output_only. The unique prefix in SQL user name.
	UserPrefix *string `json:"userPrefix,omitempty"`
	// Output_only. The current state of the cluster.
	State *Commonv1beta1ClusterState `json:"state,omitempty"`
	// Output_only. Usage details of the cluster.
	Usage *V1beta1ClusterUsage `json:"usage,omitempty"`
	// Optional. The labels for the cluster. tidb.cloud/organization. The label for the cluster organization id. tidb.cloud/project. The label for the cluster project id.
	Labels *map[string]string `json:"labels,omitempty"`
	// OUTPUT_ONLY. The annotations for the cluster. tidb.cloud/has-set-password. The annotation for whether the cluster has set password. tidb.cloud/available-features. The annotation for the available features of the cluster.
	Annotations *map[string]string `json:"annotations,omitempty"`
	// Output_only. Timestamp when the cluster was created.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// Output_only. Timestamp when the cluster was last updated.
	UpdateTime           *time.Time `json:"updateTime,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _TidbCloudOpenApiserverlessv1beta1Cluster TidbCloudOpenApiserverlessv1beta1Cluster

// NewTidbCloudOpenApiserverlessv1beta1Cluster instantiates a new TidbCloudOpenApiserverlessv1beta1Cluster object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewTidbCloudOpenApiserverlessv1beta1Cluster(displayName string, region Commonv1beta1Region) *TidbCloudOpenApiserverlessv1beta1Cluster {
	this := TidbCloudOpenApiserverlessv1beta1Cluster{}
	this.DisplayName = displayName
	this.Region = region
	return &this
}

// NewTidbCloudOpenApiserverlessv1beta1ClusterWithDefaults instantiates a new TidbCloudOpenApiserverlessv1beta1Cluster object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewTidbCloudOpenApiserverlessv1beta1ClusterWithDefaults() *TidbCloudOpenApiserverlessv1beta1Cluster {
	this := TidbCloudOpenApiserverlessv1beta1Cluster{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetName(v string) {
	o.Name = &v
}

// GetClusterId returns the ClusterId field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterId() string {
	if o == nil || IsNil(o.ClusterId) {
		var ret string
		return ret
	}
	return *o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetClusterIdOk() (*string, bool) {
	if o == nil || IsNil(o.ClusterId) {
		return nil, false
	}
	return o.ClusterId, true
}

// HasClusterId returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasClusterId() bool {
	if o != nil && !IsNil(o.ClusterId) {
		return true
	}

	return false
}

// SetClusterId gets a reference to the given string and assigns it to the ClusterId field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetClusterId(v string) {
	o.ClusterId = &v
}

// GetDisplayName returns the DisplayName field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetDisplayName(v string) {
	o.DisplayName = v
}

// GetRegion returns the Region field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRegion() Commonv1beta1Region {
	if o == nil {
		var ret Commonv1beta1Region
		return ret
	}

	return o.Region
}

// GetRegionOk returns a tuple with the Region field value
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRegionOk() (*Commonv1beta1Region, bool) {
	if o == nil {
		return nil, false
	}
	return &o.Region, true
}

// SetRegion sets field value
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetRegion(v Commonv1beta1Region) {
	o.Region = v
}

// GetSpendingLimit returns the SpendingLimit field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetSpendingLimit() ClusterSpendingLimit {
	if o == nil || IsNil(o.SpendingLimit) {
		var ret ClusterSpendingLimit
		return ret
	}
	return *o.SpendingLimit
}

// GetSpendingLimitOk returns a tuple with the SpendingLimit field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetSpendingLimitOk() (*ClusterSpendingLimit, bool) {
	if o == nil || IsNil(o.SpendingLimit) {
		return nil, false
	}
	return o.SpendingLimit, true
}

// HasSpendingLimit returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasSpendingLimit() bool {
	if o != nil && !IsNil(o.SpendingLimit) {
		return true
	}

	return false
}

// SetSpendingLimit gets a reference to the given ClusterSpendingLimit and assigns it to the SpendingLimit field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetSpendingLimit(v ClusterSpendingLimit) {
	o.SpendingLimit = &v
}

// GetAutomatedBackupPolicy returns the AutomatedBackupPolicy field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutomatedBackupPolicy() V1beta1ClusterAutomatedBackupPolicy {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		var ret V1beta1ClusterAutomatedBackupPolicy
		return ret
	}
	return *o.AutomatedBackupPolicy
}

// GetAutomatedBackupPolicyOk returns a tuple with the AutomatedBackupPolicy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAutomatedBackupPolicyOk() (*V1beta1ClusterAutomatedBackupPolicy, bool) {
	if o == nil || IsNil(o.AutomatedBackupPolicy) {
		return nil, false
	}
	return o.AutomatedBackupPolicy, true
}

// HasAutomatedBackupPolicy returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAutomatedBackupPolicy() bool {
	if o != nil && !IsNil(o.AutomatedBackupPolicy) {
		return true
	}

	return false
}

// SetAutomatedBackupPolicy gets a reference to the given V1beta1ClusterAutomatedBackupPolicy and assigns it to the AutomatedBackupPolicy field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAutomatedBackupPolicy(v V1beta1ClusterAutomatedBackupPolicy) {
	o.AutomatedBackupPolicy = &v
}

// GetEndpoints returns the Endpoints field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEndpoints() V1beta1ClusterEndpoints {
	if o == nil || IsNil(o.Endpoints) {
		var ret V1beta1ClusterEndpoints
		return ret
	}
	return *o.Endpoints
}

// GetEndpointsOk returns a tuple with the Endpoints field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEndpointsOk() (*V1beta1ClusterEndpoints, bool) {
	if o == nil || IsNil(o.Endpoints) {
		return nil, false
	}
	return o.Endpoints, true
}

// HasEndpoints returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasEndpoints() bool {
	if o != nil && !IsNil(o.Endpoints) {
		return true
	}

	return false
}

// SetEndpoints gets a reference to the given V1beta1ClusterEndpoints and assigns it to the Endpoints field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetEndpoints(v V1beta1ClusterEndpoints) {
	o.Endpoints = &v
}

// GetRootPassword returns the RootPassword field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRootPassword() string {
	if o == nil || IsNil(o.RootPassword) {
		var ret string
		return ret
	}
	return *o.RootPassword
}

// GetRootPasswordOk returns a tuple with the RootPassword field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetRootPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.RootPassword) {
		return nil, false
	}
	return o.RootPassword, true
}

// HasRootPassword returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasRootPassword() bool {
	if o != nil && !IsNil(o.RootPassword) {
		return true
	}

	return false
}

// SetRootPassword gets a reference to the given string and assigns it to the RootPassword field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetRootPassword(v string) {
	o.RootPassword = &v
}

// GetEncryptionConfig returns the EncryptionConfig field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEncryptionConfig() V1beta1ClusterEncryptionConfig {
	if o == nil || IsNil(o.EncryptionConfig) {
		var ret V1beta1ClusterEncryptionConfig
		return ret
	}
	return *o.EncryptionConfig
}

// GetEncryptionConfigOk returns a tuple with the EncryptionConfig field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetEncryptionConfigOk() (*V1beta1ClusterEncryptionConfig, bool) {
	if o == nil || IsNil(o.EncryptionConfig) {
		return nil, false
	}
	return o.EncryptionConfig, true
}

// HasEncryptionConfig returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasEncryptionConfig() bool {
	if o != nil && !IsNil(o.EncryptionConfig) {
		return true
	}

	return false
}

// SetEncryptionConfig gets a reference to the given V1beta1ClusterEncryptionConfig and assigns it to the EncryptionConfig field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetEncryptionConfig(v V1beta1ClusterEncryptionConfig) {
	o.EncryptionConfig = &v
}

// GetVersion returns the Version field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetVersion() string {
	if o == nil || IsNil(o.Version) {
		var ret string
		return ret
	}
	return *o.Version
}

// GetVersionOk returns a tuple with the Version field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetVersionOk() (*string, bool) {
	if o == nil || IsNil(o.Version) {
		return nil, false
	}
	return o.Version, true
}

// HasVersion returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasVersion() bool {
	if o != nil && !IsNil(o.Version) {
		return true
	}

	return false
}

// SetVersion gets a reference to the given string and assigns it to the Version field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetVersion(v string) {
	o.Version = &v
}

// GetCreatedBy returns the CreatedBy field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreatedBy() string {
	if o == nil || IsNil(o.CreatedBy) {
		var ret string
		return ret
	}
	return *o.CreatedBy
}

// GetCreatedByOk returns a tuple with the CreatedBy field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreatedByOk() (*string, bool) {
	if o == nil || IsNil(o.CreatedBy) {
		return nil, false
	}
	return o.CreatedBy, true
}

// HasCreatedBy returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasCreatedBy() bool {
	if o != nil && !IsNil(o.CreatedBy) {
		return true
	}

	return false
}

// SetCreatedBy gets a reference to the given string and assigns it to the CreatedBy field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetCreatedBy(v string) {
	o.CreatedBy = &v
}

// GetUserPrefix returns the UserPrefix field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUserPrefix() string {
	if o == nil || IsNil(o.UserPrefix) {
		var ret string
		return ret
	}
	return *o.UserPrefix
}

// GetUserPrefixOk returns a tuple with the UserPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUserPrefixOk() (*string, bool) {
	if o == nil || IsNil(o.UserPrefix) {
		return nil, false
	}
	return o.UserPrefix, true
}

// HasUserPrefix returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasUserPrefix() bool {
	if o != nil && !IsNil(o.UserPrefix) {
		return true
	}

	return false
}

// SetUserPrefix gets a reference to the given string and assigns it to the UserPrefix field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetUserPrefix(v string) {
	o.UserPrefix = &v
}

// GetState returns the State field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetState() Commonv1beta1ClusterState {
	if o == nil || IsNil(o.State) {
		var ret Commonv1beta1ClusterState
		return ret
	}
	return *o.State
}

// GetStateOk returns a tuple with the State field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetStateOk() (*Commonv1beta1ClusterState, bool) {
	if o == nil || IsNil(o.State) {
		return nil, false
	}
	return o.State, true
}

// HasState returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasState() bool {
	if o != nil && !IsNil(o.State) {
		return true
	}

	return false
}

// SetState gets a reference to the given Commonv1beta1ClusterState and assigns it to the State field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetState(v Commonv1beta1ClusterState) {
	o.State = &v
}

// GetUsage returns the Usage field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUsage() V1beta1ClusterUsage {
	if o == nil || IsNil(o.Usage) {
		var ret V1beta1ClusterUsage
		return ret
	}
	return *o.Usage
}

// GetUsageOk returns a tuple with the Usage field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUsageOk() (*V1beta1ClusterUsage, bool) {
	if o == nil || IsNil(o.Usage) {
		return nil, false
	}
	return o.Usage, true
}

// HasUsage returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasUsage() bool {
	if o != nil && !IsNil(o.Usage) {
		return true
	}

	return false
}

// SetUsage gets a reference to the given V1beta1ClusterUsage and assigns it to the Usage field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetUsage(v V1beta1ClusterUsage) {
	o.Usage = &v
}

// GetLabels returns the Labels field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetLabels() map[string]string {
	if o == nil || IsNil(o.Labels) {
		var ret map[string]string
		return ret
	}
	return *o.Labels
}

// GetLabelsOk returns a tuple with the Labels field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetLabelsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Labels) {
		return nil, false
	}
	return o.Labels, true
}

// HasLabels returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasLabels() bool {
	if o != nil && !IsNil(o.Labels) {
		return true
	}

	return false
}

// SetLabels gets a reference to the given map[string]string and assigns it to the Labels field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetLabels(v map[string]string) {
	o.Labels = &v
}

// GetAnnotations returns the Annotations field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAnnotations() map[string]string {
	if o == nil || IsNil(o.Annotations) {
		var ret map[string]string
		return ret
	}
	return *o.Annotations
}

// GetAnnotationsOk returns a tuple with the Annotations field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetAnnotationsOk() (*map[string]string, bool) {
	if o == nil || IsNil(o.Annotations) {
		return nil, false
	}
	return o.Annotations, true
}

// HasAnnotations returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasAnnotations() bool {
	if o != nil && !IsNil(o.Annotations) {
		return true
	}

	return false
}

// SetAnnotations gets a reference to the given map[string]string and assigns it to the Annotations field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetAnnotations(v map[string]string) {
	o.Annotations = &v
}

// GetCreateTime returns the CreateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreateTime() time.Time {
	if o == nil || IsNil(o.CreateTime) {
		var ret time.Time
		return ret
	}
	return *o.CreateTime
}

// GetCreateTimeOk returns a tuple with the CreateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetCreateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.CreateTime) {
		return nil, false
	}
	return o.CreateTime, true
}

// HasCreateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasCreateTime() bool {
	if o != nil && !IsNil(o.CreateTime) {
		return true
	}

	return false
}

// SetCreateTime gets a reference to the given time.Time and assigns it to the CreateTime field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetCreateTime(v time.Time) {
	o.CreateTime = &v
}

// GetUpdateTime returns the UpdateTime field value if set, zero value otherwise.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUpdateTime() time.Time {
	if o == nil || IsNil(o.UpdateTime) {
		var ret time.Time
		return ret
	}
	return *o.UpdateTime
}

// GetUpdateTimeOk returns a tuple with the UpdateTime field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) GetUpdateTimeOk() (*time.Time, bool) {
	if o == nil || IsNil(o.UpdateTime) {
		return nil, false
	}
	return o.UpdateTime, true
}

// HasUpdateTime returns a boolean if a field has been set.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) HasUpdateTime() bool {
	if o != nil && !IsNil(o.UpdateTime) {
		return true
	}

	return false
}

// SetUpdateTime gets a reference to the given time.Time and assigns it to the UpdateTime field.
func (o *TidbCloudOpenApiserverlessv1beta1Cluster) SetUpdateTime(v time.Time) {
	o.UpdateTime = &v
}

func (o TidbCloudOpenApiserverlessv1beta1Cluster) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o TidbCloudOpenApiserverlessv1beta1Cluster) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.ClusterId) {
		toSerialize["clusterId"] = o.ClusterId
	}
	toSerialize["displayName"] = o.DisplayName
	toSerialize["region"] = o.Region
	if !IsNil(o.SpendingLimit) {
		toSerialize["spendingLimit"] = o.SpendingLimit
	}
	if !IsNil(o.AutomatedBackupPolicy) {
		toSerialize["automatedBackupPolicy"] = o.AutomatedBackupPolicy
	}
	if !IsNil(o.Endpoints) {
		toSerialize["endpoints"] = o.Endpoints
	}
	if !IsNil(o.RootPassword) {
		toSerialize["rootPassword"] = o.RootPassword
	}
	if !IsNil(o.EncryptionConfig) {
		toSerialize["encryptionConfig"] = o.EncryptionConfig
	}
	if !IsNil(o.Version) {
		toSerialize["version"] = o.Version
	}
	if !IsNil(o.CreatedBy) {
		toSerialize["createdBy"] = o.CreatedBy
	}
	if !IsNil(o.UserPrefix) {
		toSerialize["userPrefix"] = o.UserPrefix
	}
	if !IsNil(o.State) {
		toSerialize["state"] = o.State
	}
	if !IsNil(o.Usage) {
		toSerialize["usage"] = o.Usage
	}
	if !IsNil(o.Labels) {
		toSerialize["labels"] = o.Labels
	}
	if !IsNil(o.Annotations) {
		toSerialize["annotations"] = o.Annotations
	}
	if !IsNil(o.CreateTime) {
		toSerialize["createTime"] = o.CreateTime
	}
	if !IsNil(o.UpdateTime) {
		toSerialize["updateTime"] = o.UpdateTime
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *TidbCloudOpenApiserverlessv1beta1Cluster) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"displayName",
		"region",
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

	varTidbCloudOpenApiserverlessv1beta1Cluster := _TidbCloudOpenApiserverlessv1beta1Cluster{}

	err = json.Unmarshal(data, &varTidbCloudOpenApiserverlessv1beta1Cluster)

	if err != nil {
		return err
	}

	*o = TidbCloudOpenApiserverlessv1beta1Cluster(varTidbCloudOpenApiserverlessv1beta1Cluster)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "region")
		delete(additionalProperties, "spendingLimit")
		delete(additionalProperties, "automatedBackupPolicy")
		delete(additionalProperties, "endpoints")
		delete(additionalProperties, "rootPassword")
		delete(additionalProperties, "encryptionConfig")
		delete(additionalProperties, "version")
		delete(additionalProperties, "createdBy")
		delete(additionalProperties, "userPrefix")
		delete(additionalProperties, "state")
		delete(additionalProperties, "usage")
		delete(additionalProperties, "labels")
		delete(additionalProperties, "annotations")
		delete(additionalProperties, "createTime")
		delete(additionalProperties, "updateTime")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableTidbCloudOpenApiserverlessv1beta1Cluster struct {
	value *TidbCloudOpenApiserverlessv1beta1Cluster
	isSet bool
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) Get() *TidbCloudOpenApiserverlessv1beta1Cluster {
	return v.value
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) Set(val *TidbCloudOpenApiserverlessv1beta1Cluster) {
	v.value = val
	v.isSet = true
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) IsSet() bool {
	return v.isSet
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableTidbCloudOpenApiserverlessv1beta1Cluster(val *TidbCloudOpenApiserverlessv1beta1Cluster) *NullableTidbCloudOpenApiserverlessv1beta1Cluster {
	return &NullableTidbCloudOpenApiserverlessv1beta1Cluster{value: val, isSet: true}
}

func (v NullableTidbCloudOpenApiserverlessv1beta1Cluster) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableTidbCloudOpenApiserverlessv1beta1Cluster) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
