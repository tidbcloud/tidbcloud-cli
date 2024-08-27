# ApiCreateSqlUserReq

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**AuthMethod** | Pointer to **string** | available values [mysql_native_password] . | [optional] 
**AutoPrefix** | Pointer to **bool** | if autoPrefix is true ,username and  builtinRole will automatically add the serverless token prefix. | [optional] 
**BuiltinRole** | Pointer to **string** | The builtinRole of the sql user,available values [role_admin,role_readonly,role_readwrite] . if cluster is serverless and autoPrefix is false, the builtinRole[role_readonly,role_readwrite] must be start with serverless token. | [optional] 
**CustomRoles** | Pointer to **[]string** | if cluster is serverless ,customRoles roles do not need to be prefixed. | [optional] 
**Password** | Pointer to **string** |  | [optional] 
**UserName** | Pointer to **string** | The username of the sql user, if cluster is serverless and autoPrefix is false, the userName must be start with serverless token. | [optional] 

## Methods

### NewApiCreateSqlUserReq

`func NewApiCreateSqlUserReq() *ApiCreateSqlUserReq`

NewApiCreateSqlUserReq instantiates a new ApiCreateSqlUserReq object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiCreateSqlUserReqWithDefaults

`func NewApiCreateSqlUserReqWithDefaults() *ApiCreateSqlUserReq`

NewApiCreateSqlUserReqWithDefaults instantiates a new ApiCreateSqlUserReq object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetAuthMethod

`func (o *ApiCreateSqlUserReq) GetAuthMethod() string`

GetAuthMethod returns the AuthMethod field if non-nil, zero value otherwise.

### GetAuthMethodOk

`func (o *ApiCreateSqlUserReq) GetAuthMethodOk() (*string, bool)`

GetAuthMethodOk returns a tuple with the AuthMethod field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthMethod

`func (o *ApiCreateSqlUserReq) SetAuthMethod(v string)`

SetAuthMethod sets AuthMethod field to given value.

### HasAuthMethod

`func (o *ApiCreateSqlUserReq) HasAuthMethod() bool`

HasAuthMethod returns a boolean if a field has been set.

### GetAutoPrefix

`func (o *ApiCreateSqlUserReq) GetAutoPrefix() bool`

GetAutoPrefix returns the AutoPrefix field if non-nil, zero value otherwise.

### GetAutoPrefixOk

`func (o *ApiCreateSqlUserReq) GetAutoPrefixOk() (*bool, bool)`

GetAutoPrefixOk returns a tuple with the AutoPrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAutoPrefix

`func (o *ApiCreateSqlUserReq) SetAutoPrefix(v bool)`

SetAutoPrefix sets AutoPrefix field to given value.

### HasAutoPrefix

`func (o *ApiCreateSqlUserReq) HasAutoPrefix() bool`

HasAutoPrefix returns a boolean if a field has been set.

### GetBuiltinRole

`func (o *ApiCreateSqlUserReq) GetBuiltinRole() string`

GetBuiltinRole returns the BuiltinRole field if non-nil, zero value otherwise.

### GetBuiltinRoleOk

`func (o *ApiCreateSqlUserReq) GetBuiltinRoleOk() (*string, bool)`

GetBuiltinRoleOk returns a tuple with the BuiltinRole field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBuiltinRole

`func (o *ApiCreateSqlUserReq) SetBuiltinRole(v string)`

SetBuiltinRole sets BuiltinRole field to given value.

### HasBuiltinRole

`func (o *ApiCreateSqlUserReq) HasBuiltinRole() bool`

HasBuiltinRole returns a boolean if a field has been set.

### GetCustomRoles

`func (o *ApiCreateSqlUserReq) GetCustomRoles() []string`

GetCustomRoles returns the CustomRoles field if non-nil, zero value otherwise.

### GetCustomRolesOk

`func (o *ApiCreateSqlUserReq) GetCustomRolesOk() (*[]string, bool)`

GetCustomRolesOk returns a tuple with the CustomRoles field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCustomRoles

`func (o *ApiCreateSqlUserReq) SetCustomRoles(v []string)`

SetCustomRoles sets CustomRoles field to given value.

### HasCustomRoles

`func (o *ApiCreateSqlUserReq) HasCustomRoles() bool`

HasCustomRoles returns a boolean if a field has been set.

### GetPassword

`func (o *ApiCreateSqlUserReq) GetPassword() string`

GetPassword returns the Password field if non-nil, zero value otherwise.

### GetPasswordOk

`func (o *ApiCreateSqlUserReq) GetPasswordOk() (*string, bool)`

GetPasswordOk returns a tuple with the Password field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPassword

`func (o *ApiCreateSqlUserReq) SetPassword(v string)`

SetPassword sets Password field to given value.

### HasPassword

`func (o *ApiCreateSqlUserReq) HasPassword() bool`

HasPassword returns a boolean if a field has been set.

### GetUserName

`func (o *ApiCreateSqlUserReq) GetUserName() string`

GetUserName returns the UserName field if non-nil, zero value otherwise.

### GetUserNameOk

`func (o *ApiCreateSqlUserReq) GetUserNameOk() (*string, bool)`

GetUserNameOk returns a tuple with the UserName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUserName

`func (o *ApiCreateSqlUserReq) SetUserName(v string)`

SetUserName sets UserName field to given value.

### HasUserName

`func (o *ApiCreateSqlUserReq) HasUserName() bool`

HasUserName returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


