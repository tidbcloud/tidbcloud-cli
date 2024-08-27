# ApiUpdateSqlUserReq

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BuiltinRole** | Pointer to **string** |  | [optional] 
**CustomRoles** | Pointer to **[]string** |  | [optional] 
**Password** | Pointer to **string** |  | [optional] 

## Methods

### NewApiUpdateSqlUserReq

`func NewApiUpdateSqlUserReq() *ApiUpdateSqlUserReq`

NewApiUpdateSqlUserReq instantiates a new ApiUpdateSqlUserReq object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiUpdateSqlUserReqWithDefaults

`func NewApiUpdateSqlUserReqWithDefaults() *ApiUpdateSqlUserReq`

NewApiUpdateSqlUserReqWithDefaults instantiates a new ApiUpdateSqlUserReq object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBuiltinRole

`func (o *ApiUpdateSqlUserReq) GetBuiltinRole() string`

GetBuiltinRole returns the BuiltinRole field if non-nil, zero value otherwise.

### GetBuiltinRoleOk

`func (o *ApiUpdateSqlUserReq) GetBuiltinRoleOk() (*string, bool)`

GetBuiltinRoleOk returns a tuple with the BuiltinRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuiltinRole

`func (o *ApiUpdateSqlUserReq) SetBuiltinRole(v string)`

SetBuiltinRole sets BuiltinRole field to given value.

### HasBuiltinRole

`func (o *ApiUpdateSqlUserReq) HasBuiltinRole() bool`

HasBuiltinRole returns a boolean if a field has been set.

### GetCustomRoles

`func (o *ApiUpdateSqlUserReq) GetCustomRoles() []string`

GetCustomRoles returns the CustomRoles field if non-nil, zero value otherwise.

### GetCustomRolesOk

`func (o *ApiUpdateSqlUserReq) GetCustomRolesOk() (*[]string, bool)`

GetCustomRolesOk returns a tuple with the CustomRoles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomRoles

`func (o *ApiUpdateSqlUserReq) SetCustomRoles(v []string)`

SetCustomRoles sets CustomRoles field to given value.

### HasCustomRoles

`func (o *ApiUpdateSqlUserReq) HasCustomRoles() bool`

HasCustomRoles returns a boolean if a field has been set.

### GetPassword

`func (o *ApiUpdateSqlUserReq) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ApiUpdateSqlUserReq) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ApiUpdateSqlUserReq) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *ApiUpdateSqlUserReq) HasPassword() bool`

HasPassword returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


