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

// checks if the ApiListSqlUsersRsp type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &ApiListSqlUsersRsp{}

// ApiListSqlUsersRsp struct for ApiListSqlUsersRsp
type ApiListSqlUsersRsp struct {
	// `next_page_token` can be sent in a subsequent call to fetch more results
	NextPageToken *string `json:"nextPageToken,omitempty"`
	// SqlUsers []*SqlUser `json:\"sqlUsers\"`
	SqlUsers []ApiSqlUser `json:"sqlUsers,omitempty"`
}

// NewApiListSqlUsersRsp instantiates a new ApiListSqlUsersRsp object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewApiListSqlUsersRsp() *ApiListSqlUsersRsp {
	this := ApiListSqlUsersRsp{}
	return &this
}

// NewApiListSqlUsersRspWithDefaults instantiates a new ApiListSqlUsersRsp object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewApiListSqlUsersRspWithDefaults() *ApiListSqlUsersRsp {
	this := ApiListSqlUsersRsp{}
	return &this
}

// GetNextPageToken returns the NextPageToken field value if set, zero value otherwise.
func (o *ApiListSqlUsersRsp) GetNextPageToken() string {
	if o == nil || IsNil(o.NextPageToken) {
		var ret string
		return ret
	}
	return *o.NextPageToken
}

// GetNextPageTokenOk returns a tuple with the NextPageToken field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiListSqlUsersRsp) GetNextPageTokenOk() (*string, bool) {
	if o == nil || IsNil(o.NextPageToken) {
		return nil, false
	}
	return o.NextPageToken, true
}

// HasNextPageToken returns a boolean if a field has been set.
func (o *ApiListSqlUsersRsp) HasNextPageToken() bool {
	if o != nil && !IsNil(o.NextPageToken) {
		return true
	}

	return false
}

// SetNextPageToken gets a reference to the given string and assigns it to the NextPageToken field.
func (o *ApiListSqlUsersRsp) SetNextPageToken(v string) {
	o.NextPageToken = &v
}

// GetSqlUsers returns the SqlUsers field value if set, zero value otherwise.
func (o *ApiListSqlUsersRsp) GetSqlUsers() []ApiSqlUser {
	if o == nil || IsNil(o.SqlUsers) {
		var ret []ApiSqlUser
		return ret
	}
	return o.SqlUsers
}

// GetSqlUsersOk returns a tuple with the SqlUsers field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *ApiListSqlUsersRsp) GetSqlUsersOk() ([]ApiSqlUser, bool) {
	if o == nil || IsNil(o.SqlUsers) {
		return nil, false
	}
	return o.SqlUsers, true
}

// HasSqlUsers returns a boolean if a field has been set.
func (o *ApiListSqlUsersRsp) HasSqlUsers() bool {
	if o != nil && !IsNil(o.SqlUsers) {
		return true
	}

	return false
}

// SetSqlUsers gets a reference to the given []ApiSqlUser and assigns it to the SqlUsers field.
func (o *ApiListSqlUsersRsp) SetSqlUsers(v []ApiSqlUser) {
	o.SqlUsers = v
}

func (o ApiListSqlUsersRsp) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o ApiListSqlUsersRsp) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.NextPageToken) {
		toSerialize["nextPageToken"] = o.NextPageToken
	}
	if !IsNil(o.SqlUsers) {
		toSerialize["sqlUsers"] = o.SqlUsers
	}
	return toSerialize, nil
}

type NullableApiListSqlUsersRsp struct {
	value *ApiListSqlUsersRsp
	isSet bool
}

func (v NullableApiListSqlUsersRsp) Get() *ApiListSqlUsersRsp {
	return v.value
}

func (v *NullableApiListSqlUsersRsp) Set(val *ApiListSqlUsersRsp) {
	v.value = val
	v.isSet = true
}

func (v NullableApiListSqlUsersRsp) IsSet() bool {
	return v.isSet
}

func (v *NullableApiListSqlUsersRsp) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableApiListSqlUsersRsp(val *ApiListSqlUsersRsp) *NullableApiListSqlUsersRsp {
	return &NullableApiListSqlUsersRsp{value: val, isSet: true}
}

func (v NullableApiListSqlUsersRsp) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableApiListSqlUsersRsp) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
