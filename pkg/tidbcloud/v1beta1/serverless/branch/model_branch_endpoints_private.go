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

// checks if the BranchEndpointsPrivate type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &BranchEndpointsPrivate{}

// BranchEndpointsPrivate Message for Private Endpoint for this branch.
type BranchEndpointsPrivate struct {
	// Output Only. Host Name of Public Endpoint.
	Host *string `json:"host,omitempty"`
	// Output Only. Port of Public Endpoint.
	Port *int32 `json:"port,omitempty"`
	Aws *BranchEndpointsPrivateAWS `json:"aws,omitempty"`
	Gcp *BranchEndpointsPrivateGCP `json:"gcp,omitempty"`
}

// NewBranchEndpointsPrivate instantiates a new BranchEndpointsPrivate object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewBranchEndpointsPrivate() *BranchEndpointsPrivate {
	this := BranchEndpointsPrivate{}
	return &this
}

// NewBranchEndpointsPrivateWithDefaults instantiates a new BranchEndpointsPrivate object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewBranchEndpointsPrivateWithDefaults() *BranchEndpointsPrivate {
	this := BranchEndpointsPrivate{}
	return &this
}

// GetHost returns the Host field value if set, zero value otherwise.
func (o *BranchEndpointsPrivate) GetHost() string {
	if o == nil || IsNil(o.Host) {
		var ret string
		return ret
	}
	return *o.Host
}

// GetHostOk returns a tuple with the Host field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BranchEndpointsPrivate) GetHostOk() (*string, bool) {
	if o == nil || IsNil(o.Host) {
		return nil, false
	}
	return o.Host, true
}

// HasHost returns a boolean if a field has been set.
func (o *BranchEndpointsPrivate) HasHost() bool {
	if o != nil && !IsNil(o.Host) {
		return true
	}

	return false
}

// SetHost gets a reference to the given string and assigns it to the Host field.
func (o *BranchEndpointsPrivate) SetHost(v string) {
	o.Host = &v
}

// GetPort returns the Port field value if set, zero value otherwise.
func (o *BranchEndpointsPrivate) GetPort() int32 {
	if o == nil || IsNil(o.Port) {
		var ret int32
		return ret
	}
	return *o.Port
}

// GetPortOk returns a tuple with the Port field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BranchEndpointsPrivate) GetPortOk() (*int32, bool) {
	if o == nil || IsNil(o.Port) {
		return nil, false
	}
	return o.Port, true
}

// HasPort returns a boolean if a field has been set.
func (o *BranchEndpointsPrivate) HasPort() bool {
	if o != nil && !IsNil(o.Port) {
		return true
	}

	return false
}

// SetPort gets a reference to the given int32 and assigns it to the Port field.
func (o *BranchEndpointsPrivate) SetPort(v int32) {
	o.Port = &v
}

// GetAws returns the Aws field value if set, zero value otherwise.
func (o *BranchEndpointsPrivate) GetAws() BranchEndpointsPrivateAWS {
	if o == nil || IsNil(o.Aws) {
		var ret BranchEndpointsPrivateAWS
		return ret
	}
	return *o.Aws
}

// GetAwsOk returns a tuple with the Aws field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BranchEndpointsPrivate) GetAwsOk() (*BranchEndpointsPrivateAWS, bool) {
	if o == nil || IsNil(o.Aws) {
		return nil, false
	}
	return o.Aws, true
}

// HasAws returns a boolean if a field has been set.
func (o *BranchEndpointsPrivate) HasAws() bool {
	if o != nil && !IsNil(o.Aws) {
		return true
	}

	return false
}

// SetAws gets a reference to the given BranchEndpointsPrivateAWS and assigns it to the Aws field.
func (o *BranchEndpointsPrivate) SetAws(v BranchEndpointsPrivateAWS) {
	o.Aws = &v
}

// GetGcp returns the Gcp field value if set, zero value otherwise.
func (o *BranchEndpointsPrivate) GetGcp() BranchEndpointsPrivateGCP {
	if o == nil || IsNil(o.Gcp) {
		var ret BranchEndpointsPrivateGCP
		return ret
	}
	return *o.Gcp
}

// GetGcpOk returns a tuple with the Gcp field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *BranchEndpointsPrivate) GetGcpOk() (*BranchEndpointsPrivateGCP, bool) {
	if o == nil || IsNil(o.Gcp) {
		return nil, false
	}
	return o.Gcp, true
}

// HasGcp returns a boolean if a field has been set.
func (o *BranchEndpointsPrivate) HasGcp() bool {
	if o != nil && !IsNil(o.Gcp) {
		return true
	}

	return false
}

// SetGcp gets a reference to the given BranchEndpointsPrivateGCP and assigns it to the Gcp field.
func (o *BranchEndpointsPrivate) SetGcp(v BranchEndpointsPrivateGCP) {
	o.Gcp = &v
}

func (o BranchEndpointsPrivate) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o BranchEndpointsPrivate) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Host) {
		toSerialize["host"] = o.Host
	}
	if !IsNil(o.Port) {
		toSerialize["port"] = o.Port
	}
	if !IsNil(o.Aws) {
		toSerialize["aws"] = o.Aws
	}
	if !IsNil(o.Gcp) {
		toSerialize["gcp"] = o.Gcp
	}
	return toSerialize, nil
}

type NullableBranchEndpointsPrivate struct {
	value *BranchEndpointsPrivate
	isSet bool
}

func (v NullableBranchEndpointsPrivate) Get() *BranchEndpointsPrivate {
	return v.value
}

func (v *NullableBranchEndpointsPrivate) Set(val *BranchEndpointsPrivate) {
	v.value = val
	v.isSet = true
}

func (v NullableBranchEndpointsPrivate) IsSet() bool {
	return v.isSet
}

func (v *NullableBranchEndpointsPrivate) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableBranchEndpointsPrivate(val *BranchEndpointsPrivate) *NullableBranchEndpointsPrivate {
	return &NullableBranchEndpointsPrivate{value: val, isSet: true}
}

func (v NullableBranchEndpointsPrivate) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableBranchEndpointsPrivate) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


