/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// checks if the V1beta1UpdateAuditLogConfigRequestAuditLogConfig type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1UpdateAuditLogConfigRequestAuditLogConfig{}

// V1beta1UpdateAuditLogConfigRequestAuditLogConfig struct for V1beta1UpdateAuditLogConfigRequestAuditLogConfig
type V1beta1UpdateAuditLogConfigRequestAuditLogConfig struct {
	ClusterId            string       `json:"clusterId"`
	Enabled              NullableBool `json:"enabled,omitempty"`
	BucketUri            *string      `json:"bucketUri,omitempty"`
	BucketRegionId       *string      `json:"bucketRegionId,omitempty"`
	AwsRoleArn           *string      `json:"awsRoleArn,omitempty"`
	AzureSasToken        *string      `json:"azureSasToken,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1UpdateAuditLogConfigRequestAuditLogConfig V1beta1UpdateAuditLogConfigRequestAuditLogConfig

// NewV1beta1UpdateAuditLogConfigRequestAuditLogConfig instantiates a new V1beta1UpdateAuditLogConfigRequestAuditLogConfig object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1UpdateAuditLogConfigRequestAuditLogConfig(clusterId string) *V1beta1UpdateAuditLogConfigRequestAuditLogConfig {
	this := V1beta1UpdateAuditLogConfigRequestAuditLogConfig{}
	this.ClusterId = clusterId
	return &this
}

// NewV1beta1UpdateAuditLogConfigRequestAuditLogConfigWithDefaults instantiates a new V1beta1UpdateAuditLogConfigRequestAuditLogConfig object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1UpdateAuditLogConfigRequestAuditLogConfigWithDefaults() *V1beta1UpdateAuditLogConfigRequestAuditLogConfig {
	this := V1beta1UpdateAuditLogConfigRequestAuditLogConfig{}
	return &this
}

// GetClusterId returns the ClusterId field value
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetClusterId(v string) {
	o.ClusterId = v
}

// GetEnabled returns the Enabled field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetEnabled() bool {
	if o == nil || IsNil(o.Enabled.Get()) {
		var ret bool
		return ret
	}
	return *o.Enabled.Get()
}

// GetEnabledOk returns a tuple with the Enabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetEnabledOk() (*bool, bool) {
	if o == nil {
		return nil, false
	}
	return o.Enabled.Get(), o.Enabled.IsSet()
}

// HasEnabled returns a boolean if a field has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) HasEnabled() bool {
	if o != nil && o.Enabled.IsSet() {
		return true
	}

	return false
}

// SetEnabled gets a reference to the given NullableBool and assigns it to the Enabled field.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetEnabled(v bool) {
	o.Enabled.Set(&v)
}

// SetEnabledNil sets the value for Enabled to be an explicit nil
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetEnabledNil() {
	o.Enabled.Set(nil)
}

// UnsetEnabled ensures that no value is present for Enabled, not even an explicit nil
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) UnsetEnabled() {
	o.Enabled.Unset()
}

// GetBucketUri returns the BucketUri field value if set, zero value otherwise.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetBucketUri() string {
	if o == nil || IsNil(o.BucketUri) {
		var ret string
		return ret
	}
	return *o.BucketUri
}

// GetBucketUriOk returns a tuple with the BucketUri field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetBucketUriOk() (*string, bool) {
	if o == nil || IsNil(o.BucketUri) {
		return nil, false
	}
	return o.BucketUri, true
}

// HasBucketUri returns a boolean if a field has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) HasBucketUri() bool {
	if o != nil && !IsNil(o.BucketUri) {
		return true
	}

	return false
}

// SetBucketUri gets a reference to the given string and assigns it to the BucketUri field.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetBucketUri(v string) {
	o.BucketUri = &v
}

// GetBucketRegionId returns the BucketRegionId field value if set, zero value otherwise.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetBucketRegionId() string {
	if o == nil || IsNil(o.BucketRegionId) {
		var ret string
		return ret
	}
	return *o.BucketRegionId
}

// GetBucketRegionIdOk returns a tuple with the BucketRegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetBucketRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.BucketRegionId) {
		return nil, false
	}
	return o.BucketRegionId, true
}

// HasBucketRegionId returns a boolean if a field has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) HasBucketRegionId() bool {
	if o != nil && !IsNil(o.BucketRegionId) {
		return true
	}

	return false
}

// SetBucketRegionId gets a reference to the given string and assigns it to the BucketRegionId field.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetBucketRegionId(v string) {
	o.BucketRegionId = &v
}

// GetAwsRoleArn returns the AwsRoleArn field value if set, zero value otherwise.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetAwsRoleArn() string {
	if o == nil || IsNil(o.AwsRoleArn) {
		var ret string
		return ret
	}
	return *o.AwsRoleArn
}

// GetAwsRoleArnOk returns a tuple with the AwsRoleArn field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetAwsRoleArnOk() (*string, bool) {
	if o == nil || IsNil(o.AwsRoleArn) {
		return nil, false
	}
	return o.AwsRoleArn, true
}

// HasAwsRoleArn returns a boolean if a field has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) HasAwsRoleArn() bool {
	if o != nil && !IsNil(o.AwsRoleArn) {
		return true
	}

	return false
}

// SetAwsRoleArn gets a reference to the given string and assigns it to the AwsRoleArn field.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetAwsRoleArn(v string) {
	o.AwsRoleArn = &v
}

// GetAzureSasToken returns the AzureSasToken field value if set, zero value otherwise.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetAzureSasToken() string {
	if o == nil || IsNil(o.AzureSasToken) {
		var ret string
		return ret
	}
	return *o.AzureSasToken
}

// GetAzureSasTokenOk returns a tuple with the AzureSasToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) GetAzureSasTokenOk() (*string, bool) {
	if o == nil || IsNil(o.AzureSasToken) {
		return nil, false
	}
	return o.AzureSasToken, true
}

// HasAzureSasToken returns a boolean if a field has been set.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) HasAzureSasToken() bool {
	if o != nil && !IsNil(o.AzureSasToken) {
		return true
	}

	return false
}

// SetAzureSasToken gets a reference to the given string and assigns it to the AzureSasToken field.
func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) SetAzureSasToken(v string) {
	o.AzureSasToken = &v
}

func (o V1beta1UpdateAuditLogConfigRequestAuditLogConfig) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1UpdateAuditLogConfigRequestAuditLogConfig) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["clusterId"] = o.ClusterId
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

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"clusterId",
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

	varV1beta1UpdateAuditLogConfigRequestAuditLogConfig := _V1beta1UpdateAuditLogConfigRequestAuditLogConfig{}

	err = json.Unmarshal(data, &varV1beta1UpdateAuditLogConfigRequestAuditLogConfig)

	if err != nil {
		return err
	}

	*o = V1beta1UpdateAuditLogConfigRequestAuditLogConfig(varV1beta1UpdateAuditLogConfigRequestAuditLogConfig)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "enabled")
		delete(additionalProperties, "bucketUri")
		delete(additionalProperties, "bucketRegionId")
		delete(additionalProperties, "awsRoleArn")
		delete(additionalProperties, "azureSasToken")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig struct {
	value *V1beta1UpdateAuditLogConfigRequestAuditLogConfig
	isSet bool
}

func (v NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) Get() *V1beta1UpdateAuditLogConfigRequestAuditLogConfig {
	return v.value
}

func (v *NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) Set(val *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig(val *V1beta1UpdateAuditLogConfigRequestAuditLogConfig) *NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig {
	return &NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig{value: val, isSet: true}
}

func (v NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1UpdateAuditLogConfigRequestAuditLogConfig) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
