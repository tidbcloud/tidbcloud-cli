# ApiSqlUser

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthMethod** | Pointer to **string** |  | [optional] 
**BuiltinRole** | Pointer to **string** |  | [optional] 
**CustomRoles** | Pointer to **[]string** |  | [optional] 
**UserName** | Pointer to **string** |  | [optional] 

## Methods

### NewApiSqlUser

`func NewApiSqlUser() *ApiSqlUser`

NewApiSqlUser instantiates a new ApiSqlUser object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiSqlUserWithDefaults

`func NewApiSqlUserWithDefaults() *ApiSqlUser`

NewApiSqlUserWithDefaults instantiates a new ApiSqlUser object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthMethod

`func (o *ApiSqlUser) GetAuthMethod() string`

GetAuthMethod returns the AuthMethod field if non-nil, zero value otherwise.

### GetAuthMethodOk

`func (o *ApiSqlUser) GetAuthMethodOk() (*string, bool)`

GetAuthMethodOk returns a tuple with the AuthMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthMethod

`func (o *ApiSqlUser) SetAuthMethod(v string)`

SetAuthMethod sets AuthMethod field to given value.

### HasAuthMethod

`func (o *ApiSqlUser) HasAuthMethod() bool`

HasAuthMethod returns a boolean if a field has been set.

### GetBuiltinRole

`func (o *ApiSqlUser) GetBuiltinRole() string`

GetBuiltinRole returns the BuiltinRole field if non-nil, zero value otherwise.

### GetBuiltinRoleOk

`func (o *ApiSqlUser) GetBuiltinRoleOk() (*string, bool)`

GetBuiltinRoleOk returns a tuple with the BuiltinRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuiltinRole

`func (o *ApiSqlUser) SetBuiltinRole(v string)`

SetBuiltinRole sets BuiltinRole field to given value.

### HasBuiltinRole

`func (o *ApiSqlUser) HasBuiltinRole() bool`

HasBuiltinRole returns a boolean if a field has been set.

### GetCustomRoles

`func (o *ApiSqlUser) GetCustomRoles() []string`

GetCustomRoles returns the CustomRoles field if non-nil, zero value otherwise.

### GetCustomRolesOk

`func (o *ApiSqlUser) GetCustomRolesOk() (*[]string, bool)`

GetCustomRolesOk returns a tuple with the CustomRoles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomRoles

`func (o *ApiSqlUser) SetCustomRoles(v []string)`

SetCustomRoles sets CustomRoles field to given value.

### HasCustomRoles

`func (o *ApiSqlUser) HasCustomRoles() bool`

HasCustomRoles returns a boolean if a field has been set.

### GetUserName

`func (o *ApiSqlUser) GetUserName() string`

GetUserName returns the UserName field if non-nil, zero value otherwise.

### GetUserNameOk

`func (o *ApiSqlUser) GetUserNameOk() (*string, bool)`

GetUserNameOk returns a tuple with the UserName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserName

`func (o *ApiSqlUser) SetUserName(v string)`

SetUserName sets UserName field to given value.

### HasUserName

`func (o *ApiSqlUser) HasUserName() bool`

HasUserName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


