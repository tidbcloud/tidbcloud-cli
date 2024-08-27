# ApiOpenApiListMspCustomerRsp

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**MspCustomers** | Pointer to [**[]ApiOpenApiMspCustomer**](ApiOpenApiMspCustomer.md) | The list of matching MSP Customers. | [optional] 
**NextPageToken** | Pointer to **string** | &#x60;next_page_token&#x60; can be sent in a subsequent call to fetch more results | [optional] 

## Methods

### NewApiOpenApiListMspCustomerRsp

`func NewApiOpenApiListMspCustomerRsp() *ApiOpenApiListMspCustomerRsp`

NewApiOpenApiListMspCustomerRsp instantiates a new ApiOpenApiListMspCustomerRsp object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewApiOpenApiListMspCustomerRspWithDefaults

`func NewApiOpenApiListMspCustomerRspWithDefaults() *ApiOpenApiListMspCustomerRsp`

NewApiOpenApiListMspCustomerRspWithDefaults instantiates a new ApiOpenApiListMspCustomerRsp object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetMspCustomers

`func (o *ApiOpenApiListMspCustomerRsp) GetMspCustomers() []ApiOpenApiMspCustomer`

GetMspCustomers returns the MspCustomers field if non-nil, zero value otherwise.

### GetMspCustomersOk

`func (o *ApiOpenApiListMspCustomerRsp) GetMspCustomersOk() (*[]ApiOpenApiMspCustomer, bool)`

GetMspCustomersOk returns a tuple with the MspCustomers field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetMspCustomers

`func (o *ApiOpenApiListMspCustomerRsp) SetMspCustomers(v []ApiOpenApiMspCustomer)`

SetMspCustomers sets MspCustomers field to given value.

### HasMspCustomers

`func (o *ApiOpenApiListMspCustomerRsp) HasMspCustomers() bool`

HasMspCustomers returns a boolean if a field has been set.

### GetNextPageToken

`func (o *ApiOpenApiListMspCustomerRsp) GetNextPageToken() string`

GetNextPageToken returns the NextPageToken field if non-nil, zero value otherwise.

### GetNextPageTokenOk

`func (o *ApiOpenApiListMspCustomerRsp) GetNextPageTokenOk() (*string, bool)`

GetNextPageTokenOk returns a tuple with the NextPageToken field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNextPageToken

`func (o *ApiOpenApiListMspCustomerRsp) SetNextPageToken(v string)`

SetNextPageToken sets NextPageToken field to given value.

### HasNextPageToken

`func (o *ApiOpenApiListMspCustomerRsp) HasNextPageToken() bool`

HasNextPageToken returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


