# EndpointsPrivate

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Host** | Pointer to **string** | Output Only. Host Name of Public Endpoint. | [optional] [readonly] 
**Port** | Pointer to **int32** | Output Only. Port of Public Endpoint. | [optional] [readonly] 
**Aws** | Pointer to [**PrivateAWS**](PrivateAWS.md) |  | [optional] 
**Gcp** | Pointer to [**PrivateGCP**](PrivateGCP.md) |  | [optional] 

## Methods

### NewEndpointsPrivate

`func NewEndpointsPrivate() *EndpointsPrivate`

NewEndpointsPrivate instantiates a new EndpointsPrivate object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewEndpointsPrivateWithDefaults

`func NewEndpointsPrivateWithDefaults() *EndpointsPrivate`

NewEndpointsPrivateWithDefaults instantiates a new EndpointsPrivate object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetHost

`func (o *EndpointsPrivate) GetHost() string`

GetHost returns the Host field if non-nil, zero value otherwise.

### GetHostOk

`func (o *EndpointsPrivate) GetHostOk() (*string, bool)`

GetHostOk returns a tuple with the Host field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetHost

`func (o *EndpointsPrivate) SetHost(v string)`

SetHost sets Host field to given value.

### HasHost

`func (o *EndpointsPrivate) HasHost() bool`

HasHost returns a boolean if a field has been set.

### GetPort

`func (o *EndpointsPrivate) GetPort() int32`

GetPort returns the Port field if non-nil, zero value otherwise.

### GetPortOk

`func (o *EndpointsPrivate) GetPortOk() (*int32, bool)`

GetPortOk returns a tuple with the Port field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPort

`func (o *EndpointsPrivate) SetPort(v int32)`

SetPort sets Port field to given value.

### HasPort

`func (o *EndpointsPrivate) HasPort() bool`

HasPort returns a boolean if a field has been set.

### GetAws

`func (o *EndpointsPrivate) GetAws() PrivateAWS`

GetAws returns the Aws field if non-nil, zero value otherwise.

### GetAwsOk

`func (o *EndpointsPrivate) GetAwsOk() (*PrivateAWS, bool)`

GetAwsOk returns a tuple with the Aws field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAws

`func (o *EndpointsPrivate) SetAws(v PrivateAWS)`

SetAws sets Aws field to given value.

### HasAws

`func (o *EndpointsPrivate) HasAws() bool`

HasAws returns a boolean if a field has been set.

### GetGcp

`func (o *EndpointsPrivate) GetGcp() PrivateGCP`

GetGcp returns the Gcp field if non-nil, zero value otherwise.

### GetGcpOk

`func (o *EndpointsPrivate) GetGcpOk() (*PrivateGCP, bool)`

GetGcpOk returns a tuple with the Gcp field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetGcp

`func (o *EndpointsPrivate) SetGcp(v PrivateGCP)`

SetGcp sets Gcp field to given value.

### HasGcp

`func (o *EndpointsPrivate) HasGcp() bool`

HasGcp returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


