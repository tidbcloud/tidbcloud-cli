# ExportOptionsCSVFormat

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Separator** | Pointer to **string** | Separator of each value in CSV files. It is recommended to use &#39;|+|&#39; or other uncommon character combinations. Default is &#39;,&#39;. | [optional] 
**Delimiter** | Pointer to **NullableString** | Delimiter of string type variables in CSV files. Default is &#39;\&quot;&#39;. | [optional] 
**NullValue** | Pointer to **NullableString** | Representation of null values in CSV files. Default is \&quot;\\N\&quot;. | [optional] 
**SkipHeader** | Pointer to **bool** | Export CSV files of the tables without header. Default is false. | [optional] 

## Methods

### NewExportOptionsCSVFormat

`func NewExportOptionsCSVFormat() *ExportOptionsCSVFormat`

NewExportOptionsCSVFormat instantiates a new ExportOptionsCSVFormat object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewExportOptionsCSVFormatWithDefaults

`func NewExportOptionsCSVFormatWithDefaults() *ExportOptionsCSVFormat`

NewExportOptionsCSVFormatWithDefaults instantiates a new ExportOptionsCSVFormat object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetSeparator

`func (o *ExportOptionsCSVFormat) GetSeparator() string`

GetSeparator returns the Separator field if non-nil, zero value otherwise.

### GetSeparatorOk

`func (o *ExportOptionsCSVFormat) GetSeparatorOk() (*string, bool)`

GetSeparatorOk returns a tuple with the Separator field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSeparator

`func (o *ExportOptionsCSVFormat) SetSeparator(v string)`

SetSeparator sets Separator field to given value.

### HasSeparator

`func (o *ExportOptionsCSVFormat) HasSeparator() bool`

HasSeparator returns a boolean if a field has been set.

### GetDelimiter

`func (o *ExportOptionsCSVFormat) GetDelimiter() string`

GetDelimiter returns the Delimiter field if non-nil, zero value otherwise.

### GetDelimiterOk

`func (o *ExportOptionsCSVFormat) GetDelimiterOk() (*string, bool)`

GetDelimiterOk returns a tuple with the Delimiter field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetDelimiter

`func (o *ExportOptionsCSVFormat) SetDelimiter(v string)`

SetDelimiter sets Delimiter field to given value.

### HasDelimiter

`func (o *ExportOptionsCSVFormat) HasDelimiter() bool`

HasDelimiter returns a boolean if a field has been set.

### SetDelimiterNil

`func (o *ExportOptionsCSVFormat) SetDelimiterNil(b bool)`

 SetDelimiterNil sets the value for Delimiter to be an explicit nil

### UnsetDelimiter
`func (o *ExportOptionsCSVFormat) UnsetDelimiter()`

UnsetDelimiter ensures that no value is present for Delimiter, not even an explicit nil
### GetNullValue

`func (o *ExportOptionsCSVFormat) GetNullValue() string`

GetNullValue returns the NullValue field if non-nil, zero value otherwise.

### GetNullValueOk

`func (o *ExportOptionsCSVFormat) GetNullValueOk() (*string, bool)`

GetNullValueOk returns a tuple with the NullValue field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetNullValue

`func (o *ExportOptionsCSVFormat) SetNullValue(v string)`

SetNullValue sets NullValue field to given value.

### HasNullValue

`func (o *ExportOptionsCSVFormat) HasNullValue() bool`

HasNullValue returns a boolean if a field has been set.

### SetNullValueNil

`func (o *ExportOptionsCSVFormat) SetNullValueNil(b bool)`

 SetNullValueNil sets the value for NullValue to be an explicit nil

### UnsetNullValue
`func (o *ExportOptionsCSVFormat) UnsetNullValue()`

UnsetNullValue ensures that no value is present for NullValue, not even an explicit nil
### GetSkipHeader

`func (o *ExportOptionsCSVFormat) GetSkipHeader() bool`

GetSkipHeader returns the SkipHeader field if non-nil, zero value otherwise.

### GetSkipHeaderOk

`func (o *ExportOptionsCSVFormat) GetSkipHeaderOk() (*bool, bool)`

GetSkipHeaderOk returns a tuple with the SkipHeader field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetSkipHeader

`func (o *ExportOptionsCSVFormat) SetSkipHeader(v bool)`

SetSkipHeader sets SkipHeader field to given value.

### HasSkipHeader

`func (o *ExportOptionsCSVFormat) HasSkipHeader() bool`

HasSkipHeader returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


