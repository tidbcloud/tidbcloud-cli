# S3Target

## Properties

Name | Type | Description | Notes
------------ | ------------- | ------------- | -------------
**BucketUri** | Pointer to **string** | Optional. The bucket URI of the s3. DEPRECATED, use uri instead. | [optional] 
**Uri** | Pointer to **string** | Optional. The URI of the s3 folder. | [optional] 
**AuthType** | [**ExportS3AuthTypeEnum**](ExportS3AuthTypeEnum.md) | Required. The auth method of the export s3. | 
**AccessKey** | Pointer to [**S3TargetAccessKey**](S3TargetAccessKey.md) | Optional. The access key of the s3. | [optional] 
**RoleArn** | Pointer to **string** | Optional. The role arn of the s3. | [optional] 

## Methods

### NewS3Target

`func NewS3Target(authType ExportS3AuthTypeEnum, ) *S3Target`

NewS3Target instantiates a new S3Target object
This constructor will assign default values to properties that have it defined,
and makes sure properties required by API are set, but the set of arguments
will change when the set of required properties is changed

### NewS3TargetWithDefaults

`func NewS3TargetWithDefaults() *S3Target`

NewS3TargetWithDefaults instantiates a new S3Target object
This constructor will only assign default values to properties that have it defined,
but it doesn't guarantee that properties required by API are set

### GetBucketUri

`func (o *S3Target) GetBucketUri() string`

GetBucketUri returns the BucketUri field if non-nil, zero value otherwise.

### GetBucketUriOk

`func (o *S3Target) GetBucketUriOk() (*string, bool)`

GetBucketUriOk returns a tuple with the BucketUri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetBucketUri

`func (o *S3Target) SetBucketUri(v string)`

SetBucketUri sets BucketUri field to given value.

### HasBucketUri

`func (o *S3Target) HasBucketUri() bool`

HasBucketUri returns a boolean if a field has been set.

### GetUri

`func (o *S3Target) GetUri() string`

GetUri returns the Uri field if non-nil, zero value otherwise.

### GetUriOk

`func (o *S3Target) GetUriOk() (*string, bool)`

GetUriOk returns a tuple with the Uri field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetUri

`func (o *S3Target) SetUri(v string)`

SetUri sets Uri field to given value.

### HasUri

`func (o *S3Target) HasUri() bool`

HasUri returns a boolean if a field has been set.

### GetAuthType

`func (o *S3Target) GetAuthType() ExportS3AuthTypeEnum`

GetAuthType returns the AuthType field if non-nil, zero value otherwise.

### GetAuthTypeOk

`func (o *S3Target) GetAuthTypeOk() (*ExportS3AuthTypeEnum, bool)`

GetAuthTypeOk returns a tuple with the AuthType field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAuthType

`func (o *S3Target) SetAuthType(v ExportS3AuthTypeEnum)`

SetAuthType sets AuthType field to given value.


### GetAccessKey

`func (o *S3Target) GetAccessKey() S3TargetAccessKey`

GetAccessKey returns the AccessKey field if non-nil, zero value otherwise.

### GetAccessKeyOk

`func (o *S3Target) GetAccessKeyOk() (*S3TargetAccessKey, bool)`

GetAccessKeyOk returns a tuple with the AccessKey field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetAccessKey

`func (o *S3Target) SetAccessKey(v S3TargetAccessKey)`

SetAccessKey sets AccessKey field to given value.

### HasAccessKey

`func (o *S3Target) HasAccessKey() bool`

HasAccessKey returns a boolean if a field has been set.

### GetRoleArn

`func (o *S3Target) GetRoleArn() string`

GetRoleArn returns the RoleArn field if non-nil, zero value otherwise.

### GetRoleArnOk

`func (o *S3Target) GetRoleArnOk() (*string, bool)`

GetRoleArnOk returns a tuple with the RoleArn field if it's non-nil, zero value otherwise
and a boolean to check if the value has been set.

### SetRoleArn

`func (o *S3Target) SetRoleArn(v string)`

SetRoleArn sets RoleArn field to given value.

### HasRoleArn

`func (o *S3Target) HasRoleArn() bool`

HasRoleArn returns a boolean if a field has been set.


[[Back to Model list]](../README.md#documentation-for-models) [[Back to API list]](../README.md#documentation-for-api-endpoints) [[Back to README]](../README.md)


