/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
	"fmt"
)

// checks if the V1beta1AuditLogFilterRule type satisfies the MappedNullable interface at compile time
var _ MappedNullable = &V1beta1AuditLogFilterRule{}

// V1beta1AuditLogFilterRule struct for V1beta1AuditLogFilterRule
type V1beta1AuditLogFilterRule struct {
	Name *string `json:"name,omitempty"`
	AuditLogFilterRuleId *string `json:"auditLogFilterRuleId,omitempty"`
	ClusterId string `json:"clusterId"`
	UserExpr *string `json:"userExpr,omitempty"`
	DbExpr *string `json:"dbExpr,omitempty"`
	TableExpr *string `json:"tableExpr,omitempty"`
	AccessTypeList []string `json:"accessTypeList,omitempty"`
	AdditionalProperties map[string]interface{}
}

type _V1beta1AuditLogFilterRule V1beta1AuditLogFilterRule

// NewV1beta1AuditLogFilterRule instantiates a new V1beta1AuditLogFilterRule object
// This constructor will assign default values to properties that have it defined,
// and makes sure properties required by API are set, but the set of arguments
// will change when the set of required properties is changed
func NewV1beta1AuditLogFilterRule(clusterId string) *V1beta1AuditLogFilterRule {
	this := V1beta1AuditLogFilterRule{}
	this.ClusterId = clusterId
	return &this
}

// NewV1beta1AuditLogFilterRuleWithDefaults instantiates a new V1beta1AuditLogFilterRule object
// This constructor will only assign default values to properties that have it defined,
// but it doesn't guarantee that properties required by API are set
func NewV1beta1AuditLogFilterRuleWithDefaults() *V1beta1AuditLogFilterRule {
	this := V1beta1AuditLogFilterRule{}
	return &this
}

// GetName returns the Name field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetName() string {
	if o == nil || IsNil(o.Name) {
		var ret string
		return ret
	}
	return *o.Name
}

// GetNameOk returns a tuple with the Name field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetNameOk() (*string, bool) {
	if o == nil || IsNil(o.Name) {
		return nil, false
	}
	return o.Name, true
}

// HasName returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasName() bool {
	if o != nil && !IsNil(o.Name) {
		return true
	}

	return false
}

// SetName gets a reference to the given string and assigns it to the Name field.
func (o *V1beta1AuditLogFilterRule) SetName(v string) {
	o.Name = &v
}

// GetAuditLogFilterRuleId returns the AuditLogFilterRuleId field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetAuditLogFilterRuleId() string {
	if o == nil || IsNil(o.AuditLogFilterRuleId) {
		var ret string
		return ret
	}
	return *o.AuditLogFilterRuleId
}

// GetAuditLogFilterRuleIdOk returns a tuple with the AuditLogFilterRuleId field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetAuditLogFilterRuleIdOk() (*string, bool) {
	if o == nil || IsNil(o.AuditLogFilterRuleId) {
		return nil, false
	}
	return o.AuditLogFilterRuleId, true
}

// HasAuditLogFilterRuleId returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasAuditLogFilterRuleId() bool {
	if o != nil && !IsNil(o.AuditLogFilterRuleId) {
		return true
	}

	return false
}

// SetAuditLogFilterRuleId gets a reference to the given string and assigns it to the AuditLogFilterRuleId field.
func (o *V1beta1AuditLogFilterRule) SetAuditLogFilterRuleId(v string) {
	o.AuditLogFilterRuleId = &v
}

// GetClusterId returns the ClusterId field value
func (o *V1beta1AuditLogFilterRule) GetClusterId() string {
	if o == nil {
		var ret string
		return ret
	}

	return o.ClusterId
}

// GetClusterIdOk returns a tuple with the ClusterId field value
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetClusterIdOk() (*string, bool) {
	if o == nil {
		return nil, false
	}
	return &o.ClusterId, true
}

// SetClusterId sets field value
func (o *V1beta1AuditLogFilterRule) SetClusterId(v string) {
	o.ClusterId = v
}

// GetUserExpr returns the UserExpr field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetUserExpr() string {
	if o == nil || IsNil(o.UserExpr) {
		var ret string
		return ret
	}
	return *o.UserExpr
}

// GetUserExprOk returns a tuple with the UserExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetUserExprOk() (*string, bool) {
	if o == nil || IsNil(o.UserExpr) {
		return nil, false
	}
	return o.UserExpr, true
}

// HasUserExpr returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasUserExpr() bool {
	if o != nil && !IsNil(o.UserExpr) {
		return true
	}

	return false
}

// SetUserExpr gets a reference to the given string and assigns it to the UserExpr field.
func (o *V1beta1AuditLogFilterRule) SetUserExpr(v string) {
	o.UserExpr = &v
}

// GetDbExpr returns the DbExpr field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetDbExpr() string {
	if o == nil || IsNil(o.DbExpr) {
		var ret string
		return ret
	}
	return *o.DbExpr
}

// GetDbExprOk returns a tuple with the DbExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetDbExprOk() (*string, bool) {
	if o == nil || IsNil(o.DbExpr) {
		return nil, false
	}
	return o.DbExpr, true
}

// HasDbExpr returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasDbExpr() bool {
	if o != nil && !IsNil(o.DbExpr) {
		return true
	}

	return false
}

// SetDbExpr gets a reference to the given string and assigns it to the DbExpr field.
func (o *V1beta1AuditLogFilterRule) SetDbExpr(v string) {
	o.DbExpr = &v
}

// GetTableExpr returns the TableExpr field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetTableExpr() string {
	if o == nil || IsNil(o.TableExpr) {
		var ret string
		return ret
	}
	return *o.TableExpr
}

// GetTableExprOk returns a tuple with the TableExpr field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetTableExprOk() (*string, bool) {
	if o == nil || IsNil(o.TableExpr) {
		return nil, false
	}
	return o.TableExpr, true
}

// HasTableExpr returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasTableExpr() bool {
	if o != nil && !IsNil(o.TableExpr) {
		return true
	}

	return false
}

// SetTableExpr gets a reference to the given string and assigns it to the TableExpr field.
func (o *V1beta1AuditLogFilterRule) SetTableExpr(v string) {
	o.TableExpr = &v
}

// GetAccessTypeList returns the AccessTypeList field value if set, zero value otherwise.
func (o *V1beta1AuditLogFilterRule) GetAccessTypeList() []string {
	if o == nil || IsNil(o.AccessTypeList) {
		var ret []string
		return ret
	}
	return o.AccessTypeList
}

// GetAccessTypeListOk returns a tuple with the AccessTypeList field value if set, nil otherwise
// and a boolean to check if the value has been set.
func (o *V1beta1AuditLogFilterRule) GetAccessTypeListOk() ([]string, bool) {
	if o == nil || IsNil(o.AccessTypeList) {
		return nil, false
	}
	return o.AccessTypeList, true
}

// HasAccessTypeList returns a boolean if a field has been set.
func (o *V1beta1AuditLogFilterRule) HasAccessTypeList() bool {
	if o != nil && !IsNil(o.AccessTypeList) {
		return true
	}

	return false
}

// SetAccessTypeList gets a reference to the given []string and assigns it to the AccessTypeList field.
func (o *V1beta1AuditLogFilterRule) SetAccessTypeList(v []string) {
	o.AccessTypeList = v
}

func (o V1beta1AuditLogFilterRule) MarshalJSON() ([]byte, error) {
	toSerialize,err := o.ToMap()
	if err != nil {
		return []byte{}, err
	}
	return json.Marshal(toSerialize)
}

func (o V1beta1AuditLogFilterRule) ToMap() (map[string]interface{}, error) {
	toSerialize := map[string]interface{}{}
	if !IsNil(o.Name) {
		toSerialize["name"] = o.Name
	}
	if !IsNil(o.AuditLogFilterRuleId) {
		toSerialize["auditLogFilterRuleId"] = o.AuditLogFilterRuleId
	}
	toSerialize["clusterId"] = o.ClusterId
	if !IsNil(o.UserExpr) {
		toSerialize["userExpr"] = o.UserExpr
	}
	if !IsNil(o.DbExpr) {
		toSerialize["dbExpr"] = o.DbExpr
	}
	if !IsNil(o.TableExpr) {
		toSerialize["tableExpr"] = o.TableExpr
	}
	if !IsNil(o.AccessTypeList) {
		toSerialize["accessTypeList"] = o.AccessTypeList
	}

	for key, value := range o.AdditionalProperties {
		toSerialize[key] = value
	}

	return toSerialize, nil
}

func (o *V1beta1AuditLogFilterRule) UnmarshalJSON(data []byte) (err error) {
	// This validates that all required properties are included in the JSON object
	// by unmarshalling the object into a generic map with string keys and checking
	// that every required field exists as a key in the generic map.
	requiredProperties := []string{
		"clusterId",
	}

	allProperties := make(map[string]interface{})

	err = json.Unmarshal(data, &allProperties)

	if err != nil {
		return err;
	}

	for _, requiredProperty := range(requiredProperties) {
		if _, exists := allProperties[requiredProperty]; !exists {
			return fmt.Errorf("no value given for required property %v", requiredProperty)
		}
	}

	varV1beta1AuditLogFilterRule := _V1beta1AuditLogFilterRule{}

	err = json.Unmarshal(data, &varV1beta1AuditLogFilterRule)

	if err != nil {
		return err
	}

	*o = V1beta1AuditLogFilterRule(varV1beta1AuditLogFilterRule)

	additionalProperties := make(map[string]interface{})

	if err = json.Unmarshal(data, &additionalProperties); err == nil {
		delete(additionalProperties, "name")
		delete(additionalProperties, "auditLogFilterRuleId")
		delete(additionalProperties, "clusterId")
		delete(additionalProperties, "userExpr")
		delete(additionalProperties, "dbExpr")
		delete(additionalProperties, "tableExpr")
		delete(additionalProperties, "accessTypeList")
		o.AdditionalProperties = additionalProperties
	}

	return err
}

type NullableV1beta1AuditLogFilterRule struct {
	value *V1beta1AuditLogFilterRule
	isSet bool
}

func (v NullableV1beta1AuditLogFilterRule) Get() *V1beta1AuditLogFilterRule {
	return v.value
}

func (v *NullableV1beta1AuditLogFilterRule) Set(val *V1beta1AuditLogFilterRule) {
	v.value = val
	v.isSet = true
}

func (v NullableV1beta1AuditLogFilterRule) IsSet() bool {
	return v.isSet
}

func (v *NullableV1beta1AuditLogFilterRule) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableV1beta1AuditLogFilterRule(val *V1beta1AuditLogFilterRule) *NullableV1beta1AuditLogFilterRule {
	return &NullableV1beta1AuditLogFilterRule{value: val, isSet: true}
}

func (v NullableV1beta1AuditLogFilterRule) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableV1beta1AuditLogFilterRule) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}


