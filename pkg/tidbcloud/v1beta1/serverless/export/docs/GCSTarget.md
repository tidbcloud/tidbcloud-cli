# GCSTarget

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**Uri** | **string** | The GCS URI of the export target. | 
**AuthType** | [**ExportGcsAuthTypeEnum**](ExportGcsAuthTypeEnum.md) | The auth method of the export target. | 
**ServiceAccountKey** | Pointer to **string** |  | [optional] 

## Methods

### NewGCSTarget

`func NewGCSTarget(uri string, authType ExportGcsAuthTypeEnum, ) *GCSTarget`

NewGCSTarget instantiates a new GCSTarget object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewGCSTargetWithDefaults

`func NewGCSTargetWithDefaults() *GCSTarget`

NewGCSTargetWithDefaults instantiates a new GCSTarget object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetUri

`func (o *GCSTarget) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *GCSTarget) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *GCSTarget) SetUri(v string)`

SetUri sets Uri field to given value.


### GetAuthType

`func (o *GCSTarget) GetAuthType() ExportGcsAuthTypeEnum`

GetAuthType returns the AuthType field if non-nil, zero value otherwise.

### GetAuthTypeOk

`func (o *GCSTarget) GetAuthTypeOk() (*ExportGcsAuthTypeEnum, bool)`

GetAuthTypeOk returns a tuple with the AuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthType

`func (o *GCSTarget) SetAuthType(v ExportGcsAuthTypeEnum)`

SetAuthType sets AuthType field to given value.


### GetServiceAccountKey

`func (o *GCSTarget) GetServiceAccountKey() string`

GetServiceAccountKey returns the ServiceAccountKey field if non-nil, zero value otherwise.

### GetServiceAccountKeyOk

`func (o *GCSTarget) GetServiceAccountKeyOk() (*string, bool)`

GetServiceAccountKeyOk returns a tuple with the ServiceAccountKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetServiceAccountKey

`func (o *GCSTarget) SetServiceAccountKey(v string)`

SetServiceAccountKey sets ServiceAccountKey field to given value.

### HasServiceAccountKey

`func (o *GCSTarget) HasServiceAccountKey() bool`

HasServiceAccountKey returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


