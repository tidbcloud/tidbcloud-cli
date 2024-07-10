# BranchUsage

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**RequestUnit** | Pointer to **string** | Output Only. The latest value of Request Unit Metric for this cluster. | [optional] [readonly] 
**RowStorage** | Pointer to **float64** | Output Only. The latest value of Row Storage Metric for this cluster. | [optional] [readonly] 
**ColumnarStorage** | Pointer to **float64** | Output Only. The latest value of Columnar Storage Metric for this cluster. | [optional] [readonly] 

## Methods

### NewBranchUsage

`func NewBranchUsage() *BranchUsage`

NewBranchUsage instantiates a new BranchUsage object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewBranchUsageWithDefaults

`func NewBranchUsageWithDefaults() *BranchUsage`

NewBranchUsageWithDefaults instantiates a new BranchUsage object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetRequestUnit

`func (o *BranchUsage) GetRequestUnit() string`

GetRequestUnit returns the RequestUnit field if non-nil, zero value otherwise.

### GetRequestUnitOk

`func (o *BranchUsage) GetRequestUnitOk() (*string, bool)`

GetRequestUnitOk returns a tuple with the RequestUnit field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRequestUnit

`func (o *BranchUsage) SetRequestUnit(v string)`

SetRequestUnit sets RequestUnit field to given value.

### HasRequestUnit

`func (o *BranchUsage) HasRequestUnit() bool`

HasRequestUnit returns a boolean if a field has been set.

### GetRowStorage

`func (o *BranchUsage) GetRowStorage() float64`

GetRowStorage returns the RowStorage field if non-nil, zero value otherwise.

### GetRowStorageOk

`func (o *BranchUsage) GetRowStorageOk() (*float64, bool)`

GetRowStorageOk returns a tuple with the RowStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRowStorage

`func (o *BranchUsage) SetRowStorage(v float64)`

SetRowStorage sets RowStorage field to given value.

### HasRowStorage

`func (o *BranchUsage) HasRowStorage() bool`

HasRowStorage returns a boolean if a field has been set.

### GetColumnarStorage

`func (o *BranchUsage) GetColumnarStorage() float64`

GetColumnarStorage returns the ColumnarStorage field if non-nil, zero value otherwise.

### GetColumnarStorageOk

`func (o *BranchUsage) GetColumnarStorageOk() (*float64, bool)`

GetColumnarStorageOk returns a tuple with the ColumnarStorage field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetColumnarStorage

`func (o *BranchUsage) SetColumnarStorage(v float64)`

SetColumnarStorage sets ColumnarStorage field to given value.

### HasColumnarStorage

`func (o *BranchUsage) HasColumnarStorage() bool`

HasColumnarStorage returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


