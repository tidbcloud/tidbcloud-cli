/*
TiDB Cloud Serverless Database Audit Logging Open API

TiDB Cloud Serverless Database Audit Logging Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package auditlog

import (
	"encoding/json"
	"fmt"
)

// checks if the AuditLogServiceDownloadAuditLogsBody type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &AuditLogServiceDownloadAuditLogsBody{}

// AuditLogServiceDownloadAuditLogsBody struct for AuditLogServiceDownloadAuditLogsBody
type AuditLogServiceDownloadAuditLogsBody struct {
	// Required. The name of the audit logs to download. Up to 100 audit logs can be downloaded at the same time.
	AuditLogNames        []string `json:"auditLogNames"`
	AdditionalProperties map[string]interface{}
}

type _AuditLogServiceDownloadAuditLogsBody AuditLogServiceDownloadAuditLogsBody

// NewAuditLogServiceDownloadAuditLogsBody instantiates a new AuditLogServiceDownloadAuditLogsBody object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewAuditLogServiceDownloadAuditLogsBody(auditLogNames []string) *AuditLogServiceDownloadAuditLogsBody {
	this := AuditLogServiceDownloadAuditLogsBody{}
	this.AuditLogNames = auditLogNames
	return &this
}

// NewAuditLogServiceDownloadAuditLogsBodyWithDefaults instantiates a new AuditLogServiceDownloadAuditLogsBody object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewAuditLogServiceDownloadAuditLogsBodyWithDefaults() *AuditLogServiceDownloadAuditLogsBody {
	this := AuditLogServiceDownloadAuditLogsBody{}
	return &this
}

// GetAuditLogNames returns the AuditLogNames field value
func (o *AuditLogServiceDownloadAuditLogsBody) GetAuditLogNames() []string {
	if o == nil {
		var ret []string
		return ret
	}

	return o.AuditLogNames
}

// GetAuditLogNamesOk returns a tuple with the AuditLogNames field value
// and a boolean to check if the value has been set.
func (o *AuditLogServiceDownloadAuditLogsBody) GetAuditLogNamesOk() ([]string, bool) {
	if o == nil {
		return nil, false
	}
	return o.AuditLogNames, true
}

// SetAuditLogNames sets field value
func (o *AuditLogServiceDownloadAuditLogsBody) SetAuditLogNames(v []string) {
	o.AuditLogNames = v
}

func (o AuditLogServiceDownloadAuditLogsBody) MarshalJSON() ([]byte, error) {
	toSerialize, err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o AuditLogServiceDownloadAuditLogsBody) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	toSerialize["auditLogNames"] = o.AuditLogNames

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *AuditLogServiceDownloadAuditLogsBody) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"auditLogNames",
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

	varAuditLogServiceDownloadAuditLogsBody := _AuditLogServiceDownloadAuditLogsBody{}

	err = json.Unmarshal(data, &varAuditLogServiceDownloadAuditLogsBody)

	if err != nil {
		return err
	}

	*o = AuditLogServiceDownloadAuditLogsBody(varAuditLogServiceDownloadAuditLogsBody)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "auditLogNames")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableAuditLogServiceDownloadAuditLogsBody struct {
	value *AuditLogServiceDownloadAuditLogsBody
	isSet bool
}

func (v NullableAuditLogServiceDownloadAuditLogsBody) Get() *AuditLogServiceDownloadAuditLogsBody {
	return v.value
}

func (v *NullableAuditLogServiceDownloadAuditLogsBody) Set(val *AuditLogServiceDownloadAuditLogsBody) {
	v.value = val
	v.isSet = true
}

func (v NullableAuditLogServiceDownloadAuditLogsBody) IsSet() bool {
	return v.isSet
}

func (v *NullableAuditLogServiceDownloadAuditLogsBody) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableAuditLogServiceDownloadAuditLogsBody(val *AuditLogServiceDownloadAuditLogsBody) *NullableAuditLogServiceDownloadAuditLogsBody {
	return &NullableAuditLogServiceDownloadAuditLogsBody{value: val, isSet: true}
}

func (v NullableAuditLogServiceDownloadAuditLogsBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableAuditLogServiceDownloadAuditLogsBody) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
