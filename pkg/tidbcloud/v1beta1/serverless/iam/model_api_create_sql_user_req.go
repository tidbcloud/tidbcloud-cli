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

// checks if the ApiCreateSqlUserReq type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiCreateSqlUserReq{}

// ApiCreateSqlUserReq struct for ApiCreateSqlUserReq
type ApiCreateSqlUserReq struct {
	// available values [mysql_native_password] .
	AuthMethod *string `json:"authMethod,omitempty"`
	// if autoPrefix is true ,username and  builtinRole will automatically add the serverless token prefix.
	AutoPrefix *bool `json:"autoPrefix,omitempty"`
	// The builtinRole of the sql user,available values [role_admin,role_readonly,role_readwrite] . if cluster is serverless and autoPrefix is false, the builtinRole[role_readonly,role_readwrite] must be start with serverless token.
	BuiltinRole *string `json:"builtinRole,omitempty"`
	// if cluster is serverless ,customRoles roles do not need to be prefixed.
	CustomRoles []string `json:"customRoles,omitempty"`
	Password    *string  `json:"password,omitempty"`
	// The username of the sql user, if cluster is serverless and autoPrefix is false, the userName must be start with serverless token.
	UserName *string `json:"userName,omitempty"`
}

// NewApiCreateSqlUserReq instantiates a new ApiCreateSqlUserReq object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiCreateSqlUserReq() *ApiCreateSqlUserReq {
	this := ApiCreateSqlUserReq{}
	return &this
}

// NewApiCreateSqlUserReqWithDefaults instantiates a new ApiCreateSqlUserReq object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiCreateSqlUserReqWithDefaults() *ApiCreateSqlUserReq {
	this := ApiCreateSqlUserReq{}
	return &this
}

// GetAuthMethod returns the AuthMethod field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetAuthMethod() string {
	if o == nil || IsNil(o.AuthMethod) {
		var ret string
		return ret
	}
	return *o.AuthMethod
}

// GetAuthMethodOk returns a tuple with the AuthMethod field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetAuthMethodOk() (*string, bool) {
	if o == nil || IsNil(o.AuthMethod) {
		return nil, false
	}
	return o.AuthMethod, true
}

// HasAuthMethod returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasAuthMethod() bool {
	if o != nil && !IsNil(o.AuthMethod) {
		return true
	}

	return false
}

// SetAuthMethod gets a reference to the given string and assigns it to the AuthMethod field.
func (o *ApiCreateSqlUserReq) SetAuthMethod(v string) {
	o.AuthMethod = &v
}

// GetAutoPrefix returns the AutoPrefix field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetAutoPrefix() bool {
	if o == nil || IsNil(o.AutoPrefix) {
		var ret bool
		return ret
	}
	return *o.AutoPrefix
}

// GetAutoPrefixOk returns a tuple with the AutoPrefix field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetAutoPrefixOk() (*bool, bool) {
	if o == nil || IsNil(o.AutoPrefix) {
		return nil, false
	}
	return o.AutoPrefix, true
}

// HasAutoPrefix returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasAutoPrefix() bool {
	if o != nil && !IsNil(o.AutoPrefix) {
		return true
	}

	return false
}

// SetAutoPrefix gets a reference to the given bool and assigns it to the AutoPrefix field.
func (o *ApiCreateSqlUserReq) SetAutoPrefix(v bool) {
	o.AutoPrefix = &v
}

// GetBuiltinRole returns the BuiltinRole field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetBuiltinRole() string {
	if o == nil || IsNil(o.BuiltinRole) {
		var ret string
		return ret
	}
	return *o.BuiltinRole
}

// GetBuiltinRoleOk returns a tuple with the BuiltinRole field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetBuiltinRoleOk() (*string, bool) {
	if o == nil || IsNil(o.BuiltinRole) {
		return nil, false
	}
	return o.BuiltinRole, true
}

// HasBuiltinRole returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasBuiltinRole() bool {
	if o != nil && !IsNil(o.BuiltinRole) {
		return true
	}

	return false
}

// SetBuiltinRole gets a reference to the given string and assigns it to the BuiltinRole field.
func (o *ApiCreateSqlUserReq) SetBuiltinRole(v string) {
	o.BuiltinRole = &v
}

// GetCustomRoles returns the CustomRoles field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetCustomRoles() []string {
	if o == nil || IsNil(o.CustomRoles) {
		var ret []string
		return ret
	}
	return o.CustomRoles
}

// GetCustomRolesOk returns a tuple with the CustomRoles field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetCustomRolesOk() ([]string, bool) {
	if o == nil || IsNil(o.CustomRoles) {
		return nil, false
	}
	return o.CustomRoles, true
}

// HasCustomRoles returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasCustomRoles() bool {
	if o != nil && !IsNil(o.CustomRoles) {
		return true
	}

	return false
}

// SetCustomRoles gets a reference to the given []string and assigns it to the CustomRoles field.
func (o *ApiCreateSqlUserReq) SetCustomRoles(v []string) {
	o.CustomRoles = v
}

// GetPassword returns the Password field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetPassword() string {
	if o == nil || IsNil(o.Password) {
		var ret string
		return ret
	}
	return *o.Password
}

// GetPasswordOk returns a tuple with the Password field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetPasswordOk() (*string, bool) {
	if o == nil || IsNil(o.Password) {
		return nil, false
	}
	return o.Password, true
}

// HasPassword returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasPassword() bool {
	if o != nil && !IsNil(o.Password) {
		return true
	}

	return false
}

// SetPassword gets a reference to the given string and assigns it to the Password field.
func (o *ApiCreateSqlUserReq) SetPassword(v string) {
	o.Password = &v
}

// GetUserName returns the UserName field value if set, zero value otherwise.
func (o *ApiCreateSqlUserReq) GetUserName() string {
	if o == nil || IsNil(o.UserName) {
		var ret string
		return ret
	}
	return *o.UserName
}

// GetUserNameOk returns a tuple with the UserName field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiCreateSqlUserReq) GetUserNameOk() (*string, bool) {
	if o == nil || IsNil(o.UserName) {
		return nil, false
	}
	return o.UserName, true
}

// HasUserName returns a boolean if a field has been set.
func (o *ApiCreateSqlUserReq) HasUserName() bool {
	if o != nil && !IsNil(o.UserName) {
		return true
	}

	return false
}

// SetUserName gets a reference to the given string and assigns it to the UserName field.
func (o *ApiCreateSqlUserReq) SetUserName(v string) {
	o.UserName = &v
}

func (o ApiCreateSqlUserReq) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiCreateSqlUserReq) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.AuthMethod) {
		toSerialize["authMethod"] = o.AuthMethod
	}
	if !IsNil(o.AutoPrefix) {
		toSerialize["autoPrefix"] = o.AutoPrefix
	}
	if !IsNil(o.BuiltinRole) {
		toSerialize["builtinRole"] = o.BuiltinRole
	}
	if !IsNil(o.CustomRoles) {
		toSerialize["customRoles"] = o.CustomRoles
	}
	if !IsNil(o.Password) {
		toSerialize["password"] = o.Password
	}
	if !IsNil(o.UserName) {
		toSerialize["userName"] = o.UserName
	}
	return toSerialize, nil
}

type NullableApiCreateSqlUserReq struct {
	value *ApiCreateSqlUserReq
	isSet bool
}

func (v NullableApiCreateSqlUserReq) Get() *ApiCreateSqlUserReq {
	return v.value
}

func (v *NullableApiCreateSqlUserReq) Set(val *ApiCreateSqlUserReq) {
	v.value = val
	v.isSet = true
}

func (v NullableApiCreateSqlUserReq) IsSet() bool {
	return v.isSet
}

func (v *NullableApiCreateSqlUserReq) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiCreateSqlUserReq(val *ApiCreateSqlUserReq) *NullableApiCreateSqlUserReq {
	return &NullableApiCreateSqlUserReq{value: val, isSet: true}
}

func (v NullableApiCreateSqlUserReq) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiCreateSqlUserReq) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
