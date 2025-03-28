/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
)

// checks if the DatabaseAuditLogServiceUpdateAuditLogConfigRequest type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &DatabaseAuditLogServiceUpdateAuditLogConfigRequest{}

// DatabaseAuditLogServiceUpdateAuditLogConfigRequest struct for DatabaseAuditLogServiceUpdateAuditLogConfigRequest
type DatabaseAuditLogServiceUpdateAuditLogConfigRequest struct {
	Enabled              NullableBool                   `json:"enabled,omitempty"`
	BucketUri            *string                        `json:"bucketUri,omitempty"`
	BucketRegionId       *string                        `json:"bucketRegionId,omitempty"`
	AwsRoleArn           *string                        `json:"awsRoleArn,omitempty"`
	AzureSasToken        *string                        `json:"azureSasToken,omitempty"`
	BucketManager        *Dedicatedv1beta1BucketManager `json:"bucketManager,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _DatabaseAuditLogServiceUpdateAuditLogConfigRequest DatabaseAuditLogServiceUpdateAuditLogConfigRequest

// NewDatabaseAuditLogServiceUpdateAuditLogConfigRequest instantiates a new DatabaseAuditLogServiceUpdateAuditLogConfigRequest object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewDatabaseAuditLogServiceUpdateAuditLogConfigRequest() *DatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	this := DatabaseAuditLogServiceUpdateAuditLogConfigRequest{}
	return &this
}

// NewDatabaseAuditLogServiceUpdateAuditLogConfigRequestWithDefaults instantiates a new DatabaseAuditLogServiceUpdateAuditLogConfigRequest object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewDatabaseAuditLogServiceUpdateAuditLogConfigRequestWithDefaults() *DatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	this := DatabaseAuditLogServiceUpdateAuditLogConfigRequest{}
	return &this
}

// GetEnabled returns the Enabled field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled.Get()) {
		var ret bool
		return ret
	}
	return *o.Enabled.Get()
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.Enabled.Get(), o.Enabled.IsSet()
}

// HasEnabled returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasEnabled() bool {
	if o != nil && o.Enabled.IsSet() {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given NullableBool and assigns it to the Enabled field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetEnabled(v bool) {
	o.Enabled.Set(&v)
}

// SetEnabledNil sets the value for Enabled to be an explicit nil
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetEnabledNil() {
	o.Enabled.Set(nil)
}

// UnsetEnabled ensures that no value is present for Enabled, not even an explicit nil
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) UnsetEnabled() {
	o.Enabled.Unset()
}

// GetBucketUri returns the BucketUri field value if set, zero value otherwise.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketUri() string {
	if o == nil || IsNil(o.BucketUri) {
		var ret string
		return ret
	}
	return *o.BucketUri
}

// GetBucketUriOk returns a tuple with the BucketUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketUriOk() (*string, bool) {
	if o == nil || IsNil(o.BucketUri) {
		return nil, false
	}
	return o.BucketUri, true
}

// HasBucketUri returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasBucketUri() bool {
	if o != nil && !IsNil(o.BucketUri) {
		return true
	}

	return false
}

// SetBucketUri gets a reference to the given string and assigns it to the BucketUri field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetBucketUri(v string) {
	o.BucketUri = &v
}

// GetBucketRegionId returns the BucketRegionId field value if set, zero value otherwise.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketRegionId() string {
	if o == nil || IsNil(o.BucketRegionId) {
		var ret string
		return ret
	}
	return *o.BucketRegionId
}

// GetBucketRegionIdOk returns a tuple with the BucketRegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.BucketRegionId) {
		return nil, false
	}
	return o.BucketRegionId, true
}

// HasBucketRegionId returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasBucketRegionId() bool {
	if o != nil && !IsNil(o.BucketRegionId) {
		return true
	}

	return false
}

// SetBucketRegionId gets a reference to the given string and assigns it to the BucketRegionId field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetBucketRegionId(v string) {
	o.BucketRegionId = &v
}

// GetAwsRoleArn returns the AwsRoleArn field value if set, zero value otherwise.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetAwsRoleArn() string {
	if o == nil || IsNil(o.AwsRoleArn) {
		var ret string
		return ret
	}
	return *o.AwsRoleArn
}

// GetAwsRoleArnOk returns a tuple with the AwsRoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetAwsRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.AwsRoleArn) {
		return nil, false
	}
	return o.AwsRoleArn, true
}

// HasAwsRoleArn returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasAwsRoleArn() bool {
	if o != nil && !IsNil(o.AwsRoleArn) {
		return true
	}

	return false
}

// SetAwsRoleArn gets a reference to the given string and assigns it to the AwsRoleArn field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetAwsRoleArn(v string) {
	o.AwsRoleArn = &v
}

// GetAzureSasToken returns the AzureSasToken field value if set, zero value otherwise.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetAzureSasToken() string {
	if o == nil || IsNil(o.AzureSasToken) {
		var ret string
		return ret
	}
	return *o.AzureSasToken
}

// GetAzureSasTokenOk returns a tuple with the AzureSasToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetAzureSasTokenOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSasToken) {
		return nil, false
	}
	return o.AzureSasToken, true
}

// HasAzureSasToken returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasAzureSasToken() bool {
	if o != nil && !IsNil(o.AzureSasToken) {
		return true
	}

	return false
}

// SetAzureSasToken gets a reference to the given string and assigns it to the AzureSasToken field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetAzureSasToken(v string) {
	o.AzureSasToken = &v
}

// GetBucketManager returns the BucketManager field value if set, zero value otherwise.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketManager() Dedicatedv1beta1BucketManager {
	if o == nil || IsNil(o.BucketManager) {
		var ret Dedicatedv1beta1BucketManager
		return ret
	}
	return *o.BucketManager
}

// GetBucketManagerOk returns a tuple with the BucketManager field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) GetBucketManagerOk() (*Dedicatedv1beta1BucketManager, bool) {
	if o == nil || IsNil(o.BucketManager) {
		return nil, false
	}
	return o.BucketManager, true
}

// HasBucketManager returns a boolean if a field has been set.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) HasBucketManager() bool {
	if o != nil && !IsNil(o.BucketManager) {
		return true
	}

	return false
}

// SetBucketManager gets a reference to the given Dedicatedv1beta1BucketManager and assigns it to the BucketManager field.
func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) SetBucketManager(v Dedicatedv1beta1BucketManager) {
	o.BucketManager = &v
}

func (o DatabaseAuditLogServiceUpdateAuditLogConfigRequest) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o DatabaseAuditLogServiceUpdateAuditLogConfigRequest) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if o.Enabled.IsSet() {
		toSerialize["enabled"] = o.Enabled.Get()
	}
	if !IsNil(o.BucketUri) {
		toSerialize["bucketUri"] = o.BucketUri
	}
	if !IsNil(o.BucketRegionId) {
		toSerialize["bucketRegionId"] = o.BucketRegionId
	}
	if !IsNil(o.AwsRoleArn) {
		toSerialize["awsRoleArn"] = o.AwsRoleArn
	}
	if !IsNil(o.AzureSasToken) {
		toSerialize["azureSasToken"] = o.AzureSasToken
	}
	if !IsNil(o.BucketManager) {
		toSerialize["bucketManager"] = o.BucketManager
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) UnmarshalJSON(data []byte) (err error) {
	varDatabaseAuditLogServiceUpdateAuditLogConfigRequest := _DatabaseAuditLogServiceUpdateAuditLogConfigRequest{}

	err = json.Unmarshal(data, &varDatabaseAuditLogServiceUpdateAuditLogConfigRequest)

	if err != nil {
		return err
	}

	*o = DatabaseAuditLogServiceUpdateAuditLogConfigRequest(varDatabaseAuditLogServiceUpdateAuditLogConfigRequest)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "enabled")
		delete(additionalProperties, "bucketUri")
		delete(additionalProperties, "bucketRegionId")
		delete(additionalProperties, "awsRoleArn")
		delete(additionalProperties, "azureSasToken")
		delete(additionalProperties, "bucketManager")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest struct {
	value *DatabaseAuditLogServiceUpdateAuditLogConfigRequest
	isSet bool
}

func (v NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) Get() *DatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	return v.value
}

func (v *NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) Set(val *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) {
	v.value = val
	v.isSet = true
}

func (v NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) IsSet() bool {
	return v.isSet
}

func (v *NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest(val *DatabaseAuditLogServiceUpdateAuditLogConfigRequest) *NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest {
	return &NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest{value: val, isSet: true}
}

func (v NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableDatabaseAuditLogServiceUpdateAuditLogConfigRequest) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
