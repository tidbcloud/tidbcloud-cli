# BranchEndpoints

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Public** | Pointer to [**EndpointsPublic**](EndpointsPublic.md) | Optional. Public Endpoint for this branch. | [optional] 
**Private** | Pointer to [**EndpointsPrivate**](EndpointsPrivate.md) | Output only. Private Endpoint for this branch. | [optional] 

## Methods

### NewBranchEndpoints

`func NewBranchEndpoints() *BranchEndpoints`

NewBranchEndpoints instantiates a new BranchEndpoints object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBranchEndpointsWithDefaults

`func NewBranchEndpointsWithDefaults() *BranchEndpoints`

NewBranchEndpointsWithDefaults instantiates a new BranchEndpoints object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPublic

`func (o *BranchEndpoints) GetPublic() EndpointsPublic`

GetPublic returns the Public field if non-nil, zero value otherwise.

### GetPublicOk

`func (o *BranchEndpoints) GetPublicOk() (*EndpointsPublic, bool)`

GetPublicOk returns a tuple with the Public field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPublic

`func (o *BranchEndpoints) SetPublic(v EndpointsPublic)`

SetPublic sets Public field to given value.

### HasPublic

`func (o *BranchEndpoints) HasPublic() bool`

HasPublic returns a boolean if a field has been set.

### GetPrivate

`func (o *BranchEndpoints) GetPrivate() EndpointsPrivate`

GetPrivate returns the Private field if non-nil, zero value otherwise.

### GetPrivateOk

`func (o *BranchEndpoints) GetPrivateOk() (*EndpointsPrivate, bool)`

GetPrivateOk returns a tuple with the Private field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPrivate

`func (o *BranchEndpoints) SetPrivate(v EndpointsPrivate)`

SetPrivate sets Private field to given value.

### HasPrivate

`func (o *BranchEndpoints) HasPrivate() bool`

HasPrivate returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


