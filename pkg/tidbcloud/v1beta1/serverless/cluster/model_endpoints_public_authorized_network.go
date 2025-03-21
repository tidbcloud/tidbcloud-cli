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
)

// checks if the EndpointsPublicAuthorizedNetwork type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &EndpointsPublicAuthorizedNetwork{}

// EndpointsPublicAuthorizedNetwork struct for EndpointsPublicAuthorizedNetwork
type EndpointsPublicAuthorizedNetwork struct {
	StartIpAddress       string `json:"startIpAddress"`
	EndIpAddress         string `json:"endIpAddress"`
	DisplayName          string `json:"displayName"`
	AdditionalProperties map[string]interface{}
}

type _EndpointsPublicAuthorizedNetwork EndpointsPublicAuthorizedNetwork

// NewEndpointsPublicAuthorizedNetwork instantiates a new EndpointsPublicAuthorizedNetwork object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewEndpointsPublicAuthorizedNetwork(startIpAddress string, endIpAddress string, displayName string) *EndpointsPublicAuthorizedNetwork {
	this := EndpointsPublicAuthorizedNetwork{}
	this.StartIpAddress = startIpAddress
	this.EndIpAddress = endIpAddress
	this.DisplayName = displayName
	return &this
}

// NewEndpointsPublicAuthorizedNetworkWithDefaults instantiates a new EndpointsPublicAuthorizedNetwork object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewEndpointsPublicAuthorizedNetworkWithDefaults() *EndpointsPublicAuthorizedNetwork {
	this := EndpointsPublicAuthorizedNetwork{}
	return &this
}

// GetStartIpAddress returns the StartIpAddress field value
func (o *EndpointsPublicAuthorizedNetwork) GetStartIpAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.StartIpAddress
}

// GetStartIpAddressOk returns a tuple with the StartIpAddress field value
// and a boolean to check if the value has been set.
func (o *EndpointsPublicAuthorizedNetwork) GetStartIpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.StartIpAddress, true
}

// SetStartIpAddress sets field value
func (o *EndpointsPublicAuthorizedNetwork) SetStartIpAddress(v string) {
	o.StartIpAddress = v
}

// GetEndIpAddress returns the EndIpAddress field value
func (o *EndpointsPublicAuthorizedNetwork) GetEndIpAddress() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.EndIpAddress
}

// GetEndIpAddressOk returns a tuple with the EndIpAddress field value
// and a boolean to check if the value has been set.
func (o *EndpointsPublicAuthorizedNetwork) GetEndIpAddressOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.EndIpAddress, true
}

// SetEndIpAddress sets field value
func (o *EndpointsPublicAuthorizedNetwork) SetEndIpAddress(v string) {
	o.EndIpAddress = v
}

// GetDisplayName returns the DisplayName field value
func (o *EndpointsPublicAuthorizedNetwork) GetDisplayName() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.DisplayName
}

// GetDisplayNameOk returns a tuple with the DisplayName field value
// and a boolean to check if the value has been set.
func (o *EndpointsPublicAuthorizedNetwork) GetDisplayNameOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.DisplayName, true
}

// SetDisplayName sets field value
func (o *EndpointsPublicAuthorizedNetwork) SetDisplayName(v string) {
	o.DisplayName = v
}

func (o EndpointsPublicAuthorizedNetwork) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o EndpointsPublicAuthorizedNetwork) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["startIpAddress"] = o.StartIpAddress
	toSerialize["endIpAddress"] = o.EndIpAddress
	toSerialize["displayName"] = o.DisplayName

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *EndpointsPublicAuthorizedNetwork) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"startIpAddress",
		"endIpAddress",
		"displayName",
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

	varEndpointsPublicAuthorizedNetwork := _EndpointsPublicAuthorizedNetwork{}

	err = json.Unmarshal(data, &varEndpointsPublicAuthorizedNetwork)

	if err != nil {
		return err
	}

	*o = EndpointsPublicAuthorizedNetwork(varEndpointsPublicAuthorizedNetwork)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "startIpAddress")
		delete(additionalProperties, "endIpAddress")
		delete(additionalProperties, "displayName")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableEndpointsPublicAuthorizedNetwork struct {
	value *EndpointsPublicAuthorizedNetwork
	isSet bool
}

func (v NullableEndpointsPublicAuthorizedNetwork) Get() *EndpointsPublicAuthorizedNetwork {
	return v.value
}

func (v *NullableEndpointsPublicAuthorizedNetwork) Set(val *EndpointsPublicAuthorizedNetwork) {
	v.value = val
	v.isSet = true
}

func (v NullableEndpointsPublicAuthorizedNetwork) IsSet() bool {
	return v.isSet
}

func (v *NullableEndpointsPublicAuthorizedNetwork) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableEndpointsPublicAuthorizedNetwork(val *EndpointsPublicAuthorizedNetwork) *NullableEndpointsPublicAuthorizedNetwork {
	return &NullableEndpointsPublicAuthorizedNetwork{value: val, isSet: true}
}

func (v NullableEndpointsPublicAuthorizedNetwork) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableEndpointsPublicAuthorizedNetwork) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
