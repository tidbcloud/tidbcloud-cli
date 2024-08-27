# ApiBasicResp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Code** | Pointer to **int32** |  | [optional] 
**Message** | Pointer to **string** |  | [optional] 

## Methods

### NewApiBasicResp

`func NewApiBasicResp() *ApiBasicResp`

NewApiBasicResp instantiates a new ApiBasicResp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiBasicRespWithDefaults

`func NewApiBasicRespWithDefaults() *ApiBasicResp`

NewApiBasicRespWithDefaults instantiates a new ApiBasicResp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetCode

`func (o *ApiBasicResp) GetCode() int32`

GetCode returns the Code field if non-nil, zero value otherwise.

### GetCodeOk

`func (o *ApiBasicResp) GetCodeOk() (*int32, bool)`

GetCodeOk returns a tuple with the Code field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetCode

`func (o *ApiBasicResp) SetCode(v int32)`

SetCode sets Code field to given value.

### HasCode

`func (o *ApiBasicResp) HasCode() bool`

HasCode returns a boolean if a field has been set.

### GetMessage

`func (o *ApiBasicResp) GetMessage() string`

GetMessage returns the Message field if non-nil, zero value otherwise.

### GetMessageOk

`func (o *ApiBasicResp) GetMessageOk() (*string, bool)`

GetMessageOk returns a tuple with the Message field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMessage

`func (o *ApiBasicResp) SetMessage(v string)`

SetMessage sets Message field to given value.

### HasMessage

`func (o *ApiBasicResp) HasMessage() bool`

HasMessage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


