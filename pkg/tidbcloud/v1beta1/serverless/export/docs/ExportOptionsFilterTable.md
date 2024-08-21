# ExportOptionsFilterTable

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Patterns** | Pointer to **[]string** | Optional. The table-filter expressions. | [optional] 
**Where** | Pointer to **string** | Optional. Export only selected records. | [optional] 

## Methods

### NewExportOptionsFilterTable

`func NewExportOptionsFilterTable() *ExportOptionsFilterTable`

NewExportOptionsFilterTable instantiates a new ExportOptionsFilterTable object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportOptionsFilterTableWithDefaults

`func NewExportOptionsFilterTableWithDefaults() *ExportOptionsFilterTable`

NewExportOptionsFilterTableWithDefaults instantiates a new ExportOptionsFilterTable object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetPatterns

`func (o *ExportOptionsFilterTable) GetPatterns() []string`

GetPatterns returns the Patterns field if non-nil, zero value otherwise.

### GetPatternsOk

`func (o *ExportOptionsFilterTable) GetPatternsOk() (*[]string, bool)`

GetPatternsOk returns a tuple with the Patterns field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetPatterns

`func (o *ExportOptionsFilterTable) SetPatterns(v []string)`

SetPatterns sets Patterns field to given value.

### HasPatterns

`func (o *ExportOptionsFilterTable) HasPatterns() bool`

HasPatterns returns a boolean if a field has been set.

### GetWhere

`func (o *ExportOptionsFilterTable) GetWhere() string`

GetWhere returns the Where field if non-nil, zero value otherwise.

### GetWhereOk

`func (o *ExportOptionsFilterTable) GetWhereOk() (*string, bool)`

GetWhereOk returns a tuple with the Where field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetWhere

`func (o *ExportOptionsFilterTable) SetWhere(v string)`

SetWhere sets Where field to given value.

### HasWhere

`func (o *ExportOptionsFilterTable) HasWhere() bool`

HasWhere returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


