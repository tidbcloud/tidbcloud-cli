/*
Acccount OPENAPI

This is account open api.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package iam

import (
	"encoding/json"
)

// checks if the ApiGetDbuserRsp type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiGetDbuserRsp{}

// ApiGetDbuserRsp struct for ApiGetDbuserRsp
type ApiGetDbuserRsp struct {
	// The username connect to the cluster
	Dbuser *string `json:"dbuser,omitempty"`
	// JWT to connect to the cluster
	Jwt                  *string `json:"jwt,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ApiGetDbuserRsp ApiGetDbuserRsp

// NewApiGetDbuserRsp instantiates a new ApiGetDbuserRsp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiGetDbuserRsp() *ApiGetDbuserRsp {
	this := ApiGetDbuserRsp{}
	return &this
}

// NewApiGetDbuserRspWithDefaults instantiates a new ApiGetDbuserRsp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiGetDbuserRspWithDefaults() *ApiGetDbuserRsp {
	this := ApiGetDbuserRsp{}
	return &this
}

// GetDbuser returns the Dbuser field value if set, zero value otherwise.
func (o *ApiGetDbuserRsp) GetDbuser() string {
	if o == nil || IsNil(o.Dbuser) {
		var ret string
		return ret
	}
	return *o.Dbuser
}

// GetDbuserOk returns a tuple with the Dbuser field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiGetDbuserRsp) GetDbuserOk() (*string, bool) {
	if o == nil || IsNil(o.Dbuser) {
		return nil, false
	}
	return o.Dbuser, true
}

// HasDbuser returns a boolean if a field has been set.
func (o *ApiGetDbuserRsp) HasDbuser() bool {
	if o != nil && !IsNil(o.Dbuser) {
		return true
	}

	return false
}

// SetDbuser gets a reference to the given string and assigns it to the Dbuser field.
func (o *ApiGetDbuserRsp) SetDbuser(v string) {
	o.Dbuser = &v
}

// GetJwt returns the Jwt field value if set, zero value otherwise.
func (o *ApiGetDbuserRsp) GetJwt() string {
	if o == nil || IsNil(o.Jwt) {
		var ret string
		return ret
	}
	return *o.Jwt
}

// GetJwtOk returns a tuple with the Jwt field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiGetDbuserRsp) GetJwtOk() (*string, bool) {
	if o == nil || IsNil(o.Jwt) {
		return nil, false
	}
	return o.Jwt, true
}

// HasJwt returns a boolean if a field has been set.
func (o *ApiGetDbuserRsp) HasJwt() bool {
	if o != nil && !IsNil(o.Jwt) {
		return true
	}

	return false
}

// SetJwt gets a reference to the given string and assigns it to the Jwt field.
func (o *ApiGetDbuserRsp) SetJwt(v string) {
	o.Jwt = &v
}

func (o ApiGetDbuserRsp) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiGetDbuserRsp) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Dbuser) {
		toSerialize["dbuser"] = o.Dbuser
	}
	if !IsNil(o.Jwt) {
		toSerialize["jwt"] = o.Jwt
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ApiGetDbuserRsp) UnmarshalJSON(data []byte) (err error) {
	varApiGetDbuserRsp := _ApiGetDbuserRsp{}

	err = json.Unmarshal(data, &varApiGetDbuserRsp)

	if err != nil {
		return err
	}

	*o = ApiGetDbuserRsp(varApiGetDbuserRsp)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "dbuser")
		delete(additionalProperties, "jwt")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableApiGetDbuserRsp struct {
	value *ApiGetDbuserRsp
	isSet bool
}

func (v NullableApiGetDbuserRsp) Get() *ApiGetDbuserRsp {
	return v.value
}

func (v *NullableApiGetDbuserRsp) Set(val *ApiGetDbuserRsp) {
	v.value = val
	v.isSet = true
}

func (v NullableApiGetDbuserRsp) IsSet() bool {
	return v.isSet
}

func (v *NullableApiGetDbuserRsp) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiGetDbuserRsp(val *ApiGetDbuserRsp) *NullableApiGetDbuserRsp {
	return &NullableApiGetDbuserRsp{value: val, isSet: true}
}

func (v NullableApiGetDbuserRsp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiGetDbuserRsp) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
