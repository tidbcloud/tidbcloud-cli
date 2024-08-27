# ApiGetDbuserRsp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Dbuser** | Pointer to **string** | The username connect to the cluster | [optional] 
**Jwt** | Pointer to **string** | JWT to connect to the cluster | [optional] 

## Methods

### NewApiGetDbuserRsp

`func NewApiGetDbuserRsp() *ApiGetDbuserRsp`

NewApiGetDbuserRsp instantiates a new ApiGetDbuserRsp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiGetDbuserRspWithDefaults

`func NewApiGetDbuserRspWithDefaults() *ApiGetDbuserRsp`

NewApiGetDbuserRspWithDefaults instantiates a new ApiGetDbuserRsp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetDbuser

`func (o *ApiGetDbuserRsp) GetDbuser() string`

GetDbuser returns the Dbuser field if non-nil, zero value otherwise.

### GetDbuserOk

`func (o *ApiGetDbuserRsp) GetDbuserOk() (*string, bool)`

GetDbuserOk returns a tuple with the Dbuser field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDbuser

`func (o *ApiGetDbuserRsp) SetDbuser(v string)`

SetDbuser sets Dbuser field to given value.

### HasDbuser

`func (o *ApiGetDbuserRsp) HasDbuser() bool`

HasDbuser returns a boolean if a field has been set.

### GetJwt

`func (o *ApiGetDbuserRsp) GetJwt() string`

GetJwt returns the Jwt field if non-nil, zero value otherwise.

### GetJwtOk

`func (o *ApiGetDbuserRsp) GetJwtOk() (*string, bool)`

GetJwtOk returns a tuple with the Jwt field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetJwt

`func (o *ApiGetDbuserRsp) SetJwt(v string)`

SetJwt sets Jwt field to given value.

### HasJwt

`func (o *ApiGetDbuserRsp) HasJwt() bool`

HasJwt returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


