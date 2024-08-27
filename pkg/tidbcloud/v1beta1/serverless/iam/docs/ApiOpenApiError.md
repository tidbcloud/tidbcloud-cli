# ApiOpenApiError

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **string** |  | [optional] 
**Error** | Pointer to **map[string]interface{}** |  | [optional] 
**MsgPrefix** | Pointer to **string** |  | [optional] 
**Status** | Pointer to **int32** |  | [optional] 

## Methods

### NewApiOpenApiError

`func NewApiOpenApiError() *ApiOpenApiError`

NewApiOpenApiError instantiates a new ApiOpenApiError object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiOpenApiErrorWithDefaults

`func NewApiOpenApiErrorWithDefaults() *ApiOpenApiError`

NewApiOpenApiErrorWithDefaults instantiates a new ApiOpenApiError object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *ApiOpenApiError) GetCode() string`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ApiOpenApiError) GetCodeOk() (*string, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ApiOpenApiError) SetCode(v string)`

SetCode sets Code field to given value.

### HasCode

`func (o *ApiOpenApiError) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetError

`func (o *ApiOpenApiError) GetError() map[string]interface{}`

GetError returns the Error field if non-nil, zero value otherwise.

### GetErrorOk

`func (o *ApiOpenApiError) GetErrorOk() (*map[string]interface{}, bool)`

GetErrorOk returns a tuple with the Error field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetError

`func (o *ApiOpenApiError) SetError(v map[string]interface{})`

SetError sets Error field to given value.

### HasError

`func (o *ApiOpenApiError) HasError() bool`

HasError returns a boolean if a field has been set.

### GetMsgPrefix

`func (o *ApiOpenApiError) GetMsgPrefix() string`

GetMsgPrefix returns the MsgPrefix field if non-nil, zero value otherwise.

### GetMsgPrefixOk

`func (o *ApiOpenApiError) GetMsgPrefixOk() (*string, bool)`

GetMsgPrefixOk returns a tuple with the MsgPrefix field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMsgPrefix

`func (o *ApiOpenApiError) SetMsgPrefix(v string)`

SetMsgPrefix sets MsgPrefix field to given value.

### HasMsgPrefix

`func (o *ApiOpenApiError) HasMsgPrefix() bool`

HasMsgPrefix returns a boolean if a field has been set.

### GetStatus

`func (o *ApiOpenApiError) GetStatus() int32`

GetStatus returns the Status field if non-nil, zero value otherwise.

### GetStatusOk

`func (o *ApiOpenApiError) GetStatusOk() (*int32, bool)`

GetStatusOk returns a tuple with the Status field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetStatus

`func (o *ApiOpenApiError) SetStatus(v int32)`

SetStatus sets Status field to given value.

### HasStatus

`func (o *ApiOpenApiError) HasStatus() bool`

HasStatus returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


