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

// checks if the ApiUpdateSqlUserReq type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiUpdateSqlUserReq{}

// ApiUpdateSqlUserReq struct for ApiUpdateSqlUserReq
type ApiUpdateSqlUserReq struct {
	BuiltinRole          *string  `json:"builtinRole,omitempty"`
	CustomRoles          []string `json:"customRoles,omitempty"`
	Password             *string  `json:"password,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _ApiUpdateSqlUserReq ApiUpdateSqlUserReq

// NewApiUpdateSqlUserReq instantiates a new ApiUpdateSqlUserReq object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiUpdateSqlUserReq() *ApiUpdateSqlUserReq {
	this := ApiUpdateSqlUserReq{}
	return &this
}

// NewApiUpdateSqlUserReqWithDefaults instantiates a new ApiUpdateSqlUserReq object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiUpdateSqlUserReqWithDefaults() *ApiUpdateSqlUserReq {
	this := ApiUpdateSqlUserReq{}
	return &this
}

// GetBuiltinRole returns the BuiltinRole field value if set, zero value otherwise.
func (o *ApiUpdateSqlUserReq) GetBuiltinRole() string {
	if o == nil || IsNil(o.BuiltinRole) {
		var ret string
		return ret
	}
	return *o.BuiltinRole
}

// GetBuiltinRoleOk returns a tuple with the BuiltinRole field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiUpdateSqlUserReq) GetBuiltinRoleOk() (*string, bool) {
	if o == nil || IsNil(o.BuiltinRole) {
		return nil, false
	}
	return o.BuiltinRole, true
}

// HasBuiltinRole returns a boolean if a field has been set.
func (o *ApiUpdateSqlUserReq) HasBuiltinRole() bool {
	if o != nil && !IsNil(o.BuiltinRole) {
		return true
	}

	return false
}

// SetBuiltinRole gets a reference to the given string and assigns it to the BuiltinRole field.
func (o *ApiUpdateSqlUserReq) SetBuiltinRole(v string) {
	o.BuiltinRole = &v
}

// GetCustomRoles returns the CustomRoles field value if set, zero value otherwise.
func (o *ApiUpdateSqlUserReq) GetCustomRoles() []string {
	if o == nil || IsNil(o.CustomRoles) {
		var ret []string
		return ret
	}
	return o.CustomRoles
}

// GetCustomRolesOk returns a tuple with the CustomRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiUpdateSqlUserReq) GetCustomRolesOk() ([]string, bool) {
	if o == nil || IsNil(o.CustomRoles) {
		return nil, false
	}
	return o.CustomRoles, true
}

// HasCustomRoles returns a boolean if a field has been set.
func (o *ApiUpdateSqlUserReq) HasCustomRoles() bool {
	if o != nil && !IsNil(o.CustomRoles) {
		return true
	}

	return false
}

// SetCustomRoles gets a reference to the given []string and assigns it to the CustomRoles field.
func (o *ApiUpdateSqlUserReq) SetCustomRoles(v []string) {
	o.CustomRoles = v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *ApiUpdateSqlUserReq) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiUpdateSqlUserReq) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *ApiUpdateSqlUserReq) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *ApiUpdateSqlUserReq) SetPassword(v string) {
	o.Password = &v
}

func (o ApiUpdateSqlUserReq) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiUpdateSqlUserReq) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.BuiltinRole) {
		toSerialize["builtinRole"] = o.BuiltinRole
	}
	if !IsNil(o.CustomRoles) {
		toSerialize["customRoles"] = o.CustomRoles
	}
	if !IsNil(o.Password) {
		toSerialize["password"] = o.Password
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *ApiUpdateSqlUserReq) UnmarshalJSON(data []byte) (err error) {
	varApiUpdateSqlUserReq := _ApiUpdateSqlUserReq{}

	err = json.Unmarshal(data, &varApiUpdateSqlUserReq)

	if err != nil {
		return err
	}

	*o = ApiUpdateSqlUserReq(varApiUpdateSqlUserReq)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "builtinRole")
		delete(additionalProperties, "customRoles")
		delete(additionalProperties, "password")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableApiUpdateSqlUserReq struct {
	value *ApiUpdateSqlUserReq
	isSet bool
}

func (v NullableApiUpdateSqlUserReq) Get() *ApiUpdateSqlUserReq {
	return v.value
}

func (v *NullableApiUpdateSqlUserReq) Set(val *ApiUpdateSqlUserReq) {
	v.value = val
	v.isSet = true
}

func (v NullableApiUpdateSqlUserReq) IsSet() bool {
	return v.isSet
}

func (v *NullableApiUpdateSqlUserReq) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiUpdateSqlUserReq(val *ApiUpdateSqlUserReq) *NullableApiUpdateSqlUserReq {
	return &NullableApiUpdateSqlUserReq{value: val, isSet: true}
}

func (v NullableApiUpdateSqlUserReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiUpdateSqlUserReq) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
