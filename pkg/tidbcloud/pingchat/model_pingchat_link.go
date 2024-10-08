/*
PingChat Swagger API

No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)

API version: 1.0
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package pingchat

import (
	"encoding/json"
)

// checks if the PingchatLink type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &PingchatLink{}

// PingchatLink struct for PingchatLink
type PingchatLink struct {
	Link  *string `json:"link,omitempty"`
	Title *string `json:"title,omitempty"`
}

// NewPingchatLink instantiates a new PingchatLink object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewPingchatLink() *PingchatLink {
	this := PingchatLink{}
	return &this
}

// NewPingchatLinkWithDefaults instantiates a new PingchatLink object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewPingchatLinkWithDefaults() *PingchatLink {
	this := PingchatLink{}
	return &this
}

// GetLink returns the Link field value if set, zero value otherwise.
func (o *PingchatLink) GetLink() string {
	if o == nil || IsNil(o.Link) {
		var ret string
		return ret
	}
	return *o.Link
}

// GetLinkOk returns a tuple with the Link field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PingchatLink) GetLinkOk() (*string, bool) {
	if o == nil || IsNil(o.Link) {
		return nil, false
	}
	return o.Link, true
}

// HasLink returns a boolean if a field has been set.
func (o *PingchatLink) HasLink() bool {
	if o != nil && !IsNil(o.Link) {
		return true
	}

	return false
}

// SetLink gets a reference to the given string and assigns it to the Link field.
func (o *PingchatLink) SetLink(v string) {
	o.Link = &v
}

// GetTitle returns the Title field value if set, zero value otherwise.
func (o *PingchatLink) GetTitle() string {
	if o == nil || IsNil(o.Title) {
		var ret string
		return ret
	}
	return *o.Title
}

// GetTitleOk returns a tuple with the Title field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *PingchatLink) GetTitleOk() (*string, bool) {
	if o == nil || IsNil(o.Title) {
		return nil, false
	}
	return o.Title, true
}

// HasTitle returns a boolean if a field has been set.
func (o *PingchatLink) HasTitle() bool {
	if o != nil && !IsNil(o.Title) {
		return true
	}

	return false
}

// SetTitle gets a reference to the given string and assigns it to the Title field.
func (o *PingchatLink) SetTitle(v string) {
	o.Title = &v
}

func (o PingchatLink) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o PingchatLink) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Link) {
		toSerialize["link"] = o.Link
	}
	if !IsNil(o.Title) {
		toSerialize["title"] = o.Title
	}
	return toSerialize, nil
}

type NullablePingchatLink struct {
	value *PingchatLink
	isSet bool
}

func (v NullablePingchatLink) Get() *PingchatLink {
	return v.value
}

func (v *NullablePingchatLink) Set(val *PingchatLink) {
	v.value = val
	v.isSet = true
}

func (v NullablePingchatLink) IsSet() bool {
	return v.isSet
}

func (v *NullablePingchatLink) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullablePingchatLink(val *PingchatLink) *NullablePingchatLink {
	return &NullablePingchatLink{value: val, isSet: true}
}

func (v NullablePingchatLink) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullablePingchatLink) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
