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

// checks if the Commonv1beta1Region type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &Commonv1beta1Region{}

// Commonv1beta1Region A representation of a region for deploying TiDB clusters.
type Commonv1beta1Region struct {
	Name *string `json:"name,omitempty" validate:"regexp=^regions\\/(aws|gcp|azure)-(.+)$"`
	// Format: {cloud_provider}-{region_code} Region code: us-west-2, asia-east1.
	RegionId *string `json:"regionId,omitempty"`
	// The cloud provider for the region.
	CloudProvider *V1beta1RegionCloudProvider `json:"cloudProvider,omitempty"`
	// User-friendly display name of the region.
	DisplayName *string `json:"displayName,omitempty"`
	// Optional provider name for the region. Only used for serverless cluster. Deprecated.
	Provider NullableString `json:"provider,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _Commonv1beta1Region Commonv1beta1Region

// NewCommonv1beta1Region instantiates a new Commonv1beta1Region object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewCommonv1beta1Region() *Commonv1beta1Region {
	this := Commonv1beta1Region{}
	return &this
}

// NewCommonv1beta1RegionWithDefaults instantiates a new Commonv1beta1Region object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewCommonv1beta1RegionWithDefaults() *Commonv1beta1Region {
	this := Commonv1beta1Region{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *Commonv1beta1Region) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Commonv1beta1Region) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *Commonv1beta1Region) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *Commonv1beta1Region) SetName(v string) {
	o.Name = &v
}

// GetRegionId returns the RegionId field value if set, zero value otherwise.
func (o *Commonv1beta1Region) GetRegionId() string {
	if o == nil || IsNil(o.RegionId) {
		var ret string
		return ret
	}
	return *o.RegionId
}

// GetRegionIdOk returns a tuple with the RegionId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Commonv1beta1Region) GetRegionIdOk() (*string, bool) {
	if o == nil || IsNil(o.RegionId) {
		return nil, false
	}
	return o.RegionId, true
}

// HasRegionId returns a boolean if a field has been set.
func (o *Commonv1beta1Region) HasRegionId() bool {
	if o != nil && !IsNil(o.RegionId) {
		return true
	}

	return false
}

// SetRegionId gets a reference to the given string and assigns it to the RegionId field.
func (o *Commonv1beta1Region) SetRegionId(v string) {
	o.RegionId = &v
}

// GetCloudProvider returns the CloudProvider field value if set, zero value otherwise.
func (o *Commonv1beta1Region) GetCloudProvider() V1beta1RegionCloudProvider {
	if o == nil || IsNil(o.CloudProvider) {
		var ret V1beta1RegionCloudProvider
		return ret
	}
	return *o.CloudProvider
}

// GetCloudProviderOk returns a tuple with the CloudProvider field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Commonv1beta1Region) GetCloudProviderOk() (*V1beta1RegionCloudProvider, bool) {
	if o == nil || IsNil(o.CloudProvider) {
		return nil, false
	}
	return o.CloudProvider, true
}

// HasCloudProvider returns a boolean if a field has been set.
func (o *Commonv1beta1Region) HasCloudProvider() bool {
	if o != nil && !IsNil(o.CloudProvider) {
		return true
	}

	return false
}

// SetCloudProvider gets a reference to the given V1beta1RegionCloudProvider and assigns it to the CloudProvider field.
func (o *Commonv1beta1Region) SetCloudProvider(v V1beta1RegionCloudProvider) {
	o.CloudProvider = &v
}

// GetDisplayName returns the DisplayName field value if set, zero value otherwise.
func (o *Commonv1beta1Region) GetDisplayName() string {
	if o == nil || IsNil(o.DisplayName) {
		var ret string
		return ret
	}
	return *o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *Commonv1beta1Region) GetDisplayNameOk() (*string, bool) {
	if o == nil || IsNil(o.DisplayName) {
		return nil, false
	}
	return o.DisplayName, true
}

// HasDisplayName returns a boolean if a field has been set.
func (o *Commonv1beta1Region) HasDisplayName() bool {
	if o != nil && !IsNil(o.DisplayName) {
		return true
	}

	return false
}

// SetDisplayName gets a reference to the given string and assigns it to the DisplayName field.
func (o *Commonv1beta1Region) SetDisplayName(v string) {
	o.DisplayName = &v
}

// GetProvider returns the Provider field value if set, zero value otherwise (both if not set or set to explicit null).
func (o *Commonv1beta1Region) GetProvider() string {
	if o == nil || IsNil(o.Provider.Get()) {
		var ret string
		return ret
	}
	return *o.Provider.Get()
}

// GetProviderOk returns a tuple with the Provider field value if set, nil otherwise
// and a boolean to check if the value has been set.
// NOTE: If the value is an explicit nil, `nil, true` will be returned
func (o *Commonv1beta1Region) GetProviderOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return o.Provider.Get(), o.Provider.IsSet()
}

// HasProvider returns a boolean if a field has been set.
func (o *Commonv1beta1Region) HasProvider() bool {
	if o != nil && o.Provider.IsSet() {
		return true
	}

	return false
}

// SetProvider gets a reference to the given NullableString and assigns it to the Provider field.
func (o *Commonv1beta1Region) SetProvider(v string) {
	o.Provider.Set(&v)
}
// SetProviderNil sets the value for Provider to be an explicit nil
func (o *Commonv1beta1Region) SetProviderNil() {
	o.Provider.Set(nil)
}

// UnsetProvider ensures that no value is present for Provider, not even an explicit nil
func (o *Commonv1beta1Region) UnsetProvider() {
	o.Provider.Unset()
}

func (o Commonv1beta1Region) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o Commonv1beta1Region) ToMap() (map[string]interface{}, error) {
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

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *Commonv1beta1Region) UnmarshalJSON(data []byte) (err error) {
	varCommonv1beta1Region := _Commonv1beta1Region{}

	err = json.Unmarshal(data, &varCommonv1beta1Region)

	if err != nil {
		return err
	}

	*o = Commonv1beta1Region(varCommonv1beta1Region)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "regionId")
		delete(additionalProperties, "cloudProvider")
		delete(additionalProperties, "displayName")
		delete(additionalProperties, "provider")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableCommonv1beta1Region struct {
	value *Commonv1beta1Region
	isSet bool
}

func (v NullableCommonv1beta1Region) Get() *Commonv1beta1Region {
	return v.value
}

func (v *NullableCommonv1beta1Region) Set(val *Commonv1beta1Region) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonv1beta1Region) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonv1beta1Region) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonv1beta1Region(val *Commonv1beta1Region) *NullableCommonv1beta1Region {
	return &NullableCommonv1beta1Region{value: val, isSet: true}
}

func (v NullableCommonv1beta1Region) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonv1beta1Region) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


