/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package branch

import (
	"encoding/json"
)

// checks if the EndpointsPublic type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EndpointsPublic{}

// EndpointsPublic Message for Public Endpoint for this branch.
type EndpointsPublic struct {
	Host *string `json:"host,omitempty"`
	Port *int32 `json:"port,omitempty"`
	Disabled *bool `json:"disabled,omitempty"`
}

// NewEndpointsPublic instantiates a new EndpointsPublic object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEndpointsPublic() *EndpointsPublic {
	this := EndpointsPublic{}
	return &this
}

// NewEndpointsPublicWithDefaults instantiates a new EndpointsPublic object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEndpointsPublicWithDefaults() *EndpointsPublic {
	this := EndpointsPublic{}
	return &this
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *EndpointsPublic) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointsPublic) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *EndpointsPublic) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *EndpointsPublic) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *EndpointsPublic) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointsPublic) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *EndpointsPublic) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *EndpointsPublic) SetPort(v int32) {
	o.Port = &v
}

// GetDisabled returns the Disabled field value if set, zero value otherwise.
func (o *EndpointsPublic) GetDisabled() bool {
	if o == nil || IsNil(o.Disabled) {
		var ret bool
		return ret
	}
	return *o.Disabled
}

// GetDisabledOk returns a tuple with the Disabled field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *EndpointsPublic) GetDisabledOk() (*bool, bool) {
	if o == nil || IsNil(o.Disabled) {
		return nil, false
	}
	return o.Disabled, true
}

// HasDisabled returns a boolean if a field has been set.
func (o *EndpointsPublic) HasDisabled() bool {
	if o != nil && !IsNil(o.Disabled) {
		return true
	}

	return false
}

// SetDisabled gets a reference to the given bool and assigns it to the Disabled field.
func (o *EndpointsPublic) SetDisabled(v bool) {
	o.Disabled = &v
}

func (o EndpointsPublic) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EndpointsPublic) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}
	if !IsNil(o.Disabled) {
		toSerialize["disabled"] = o.Disabled
	}
	return toSerialize, nil
}

type NullableEndpointsPublic struct {
	value *EndpointsPublic
	isSet bool
}

func (v NullableEndpointsPublic) Get() *EndpointsPublic {
	return v.value
}

func (v *NullableEndpointsPublic) Set(val *EndpointsPublic) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointsPublic) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointsPublic) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointsPublic(val *EndpointsPublic) *NullableEndpointsPublic {
	return &NullableEndpointsPublic{value: val, isSet: true}
}

func (v NullableEndpointsPublic) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointsPublic) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


