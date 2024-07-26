# PrivateAWS

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**ServiceName** | Pointer to **string** | Output Only. Service Name for Private Link Service. | [optional] [readonly] 
**AvailabilityZone** | Pointer to **[]string** | Output Only. Availability Zone for Private Link Service. | [optional] [readonly] 

## Methods

### NewPrivateAWS

`func NewPrivateAWS() *PrivateAWS`

NewPrivateAWS instantiates a new PrivateAWS object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewPrivateAWSWithDefaults

`func NewPrivateAWSWithDefaults() *PrivateAWS`

NewPrivateAWSWithDefaults instantiates a new PrivateAWS object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetServiceName

`func (o *PrivateAWS) GetServiceName() string`

GetServiceName returns the ServiceName field if non-nil, zero value otherwise.

### GetServiceNameOk

`func (o *PrivateAWS) GetServiceNameOk() (*string, bool)`

GetServiceNameOk returns a tuple with the ServiceName field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceName

`func (o *PrivateAWS) SetServiceName(v string)`

SetServiceName sets ServiceName field to given value.

### HasServiceName

`func (o *PrivateAWS) HasServiceName() bool`

HasServiceName returns a boolean if a field has been set.

### GetAvailabilityZone

`func (o *PrivateAWS) GetAvailabilityZone() []string`

GetAvailabilityZone returns the AvailabilityZone field if non-nil, zero value otherwise.

### GetAvailabilityZoneOk

`func (o *PrivateAWS) GetAvailabilityZoneOk() (*[]string, bool)`

GetAvailabilityZoneOk returns a tuple with the AvailabilityZone field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAvailabilityZone

`func (o *PrivateAWS) SetAvailabilityZone(v []string)`

SetAvailabilityZone sets AvailabilityZone field to given value.

### HasAvailabilityZone

`func (o *PrivateAWS) HasAvailabilityZone() bool`

HasAvailabilityZone returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


