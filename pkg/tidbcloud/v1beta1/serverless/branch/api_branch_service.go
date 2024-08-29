/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package branch

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)


// BranchServiceAPIService BranchServiceAPI service
type BranchServiceAPIService service

type ApiBranchServiceCreateBranchRequest struct {
	ctx context.Context
	ApiService *BranchServiceAPIService
	clusterId string
	branch *Branch
}

// Required. The resource being created
func (r ApiBranchServiceCreateBranchRequest) Branch(branch Branch) ApiBranchServiceCreateBranchRequest {
	r.branch = &branch
	return r
}

func (r ApiBranchServiceCreateBranchRequest) Execute() (*Branch, *http.Response, error) {
	return r.ApiService.BranchServiceCreateBranchExecute(r)
}

/*
BranchServiceCreateBranch Creates a branch.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clusterId Required. The cluster ID of the branch
 @return ApiBranchServiceCreateBranchRequest
*/
func (a *BranchServiceAPIService) BranchServiceCreateBranch(ctx context.Context, clusterId string) ApiBranchServiceCreateBranchRequest {
	return ApiBranchServiceCreateBranchRequest{
		ApiService: a,
		ctx: ctx,
		clusterId: clusterId,
	}
}

// Execute executes the request
//  @return Branch
func (a *BranchServiceAPIService) BranchServiceCreateBranchExecute(r ApiBranchServiceCreateBranchRequest) (*Branch, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodPost
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Branch
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BranchServiceAPIService.BranchServiceCreateBranch")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/branches"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.branch == nil {
		return localVarReturnValue, nil, reportError("branch is required and must be specified")
	}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{"application/json"}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	// body params
	localVarPostBody = r.branch
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
			var v Status
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiBranchServiceDeleteBranchRequest struct {
	ctx context.Context
	ApiService *BranchServiceAPIService
	clusterId string
	branchId string
}

func (r ApiBranchServiceDeleteBranchRequest) Execute() (*Branch, *http.Response, error) {
	return r.ApiService.BranchServiceDeleteBranchExecute(r)
}

/*
BranchServiceDeleteBranch Deletes a branch.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clusterId Required. The cluster ID of the branch
 @param branchId Required. The branch ID
 @return ApiBranchServiceDeleteBranchRequest
*/
func (a *BranchServiceAPIService) BranchServiceDeleteBranch(ctx context.Context, clusterId string, branchId string) ApiBranchServiceDeleteBranchRequest {
	return ApiBranchServiceDeleteBranchRequest{
		ApiService: a,
		ctx: ctx,
		clusterId: clusterId,
		branchId: branchId,
	}
}

// Execute executes the request
//  @return Branch
func (a *BranchServiceAPIService) BranchServiceDeleteBranchExecute(r ApiBranchServiceDeleteBranchRequest) (*Branch, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodDelete
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Branch
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BranchServiceAPIService.BranchServiceDeleteBranch")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/branches/{branchId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"branchId"+"}", url.PathEscape(parameterValueToString(r.branchId, "branchId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
			var v Status
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiBranchServiceGetBranchRequest struct {
	ctx context.Context
	ApiService *BranchServiceAPIService
	clusterId string
	branchId string
	view *string
}

// Optional. The view of the branch to return. Defaults to FULL   - BASIC: Basic response contains basic information for a branch.  - FULL: FULL response contains all detailed information for a branch.
func (r ApiBranchServiceGetBranchRequest) View(view string) ApiBranchServiceGetBranchRequest {
	r.view = &view
	return r
}

func (r ApiBranchServiceGetBranchRequest) Execute() (*Branch, *http.Response, error) {
	return r.ApiService.BranchServiceGetBranchExecute(r)
}

/*
BranchServiceGetBranch Gets information about a branch.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clusterId Required. The cluster ID of the branch
 @param branchId Required. The branch ID
 @return ApiBranchServiceGetBranchRequest
*/
func (a *BranchServiceAPIService) BranchServiceGetBranch(ctx context.Context, clusterId string, branchId string) ApiBranchServiceGetBranchRequest {
	return ApiBranchServiceGetBranchRequest{
		ApiService: a,
		ctx: ctx,
		clusterId: clusterId,
		branchId: branchId,
	}
}

// Execute executes the request
//  @return Branch
func (a *BranchServiceAPIService) BranchServiceGetBranchExecute(r ApiBranchServiceGetBranchRequest) (*Branch, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *Branch
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BranchServiceAPIService.BranchServiceGetBranch")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/branches/{branchId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"branchId"+"}", url.PathEscape(parameterValueToString(r.branchId, "branchId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.view != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "view", r.view, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
			var v Status
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}

type ApiBranchServiceListBranchesRequest struct {
	ctx context.Context
	ApiService *BranchServiceAPIService
	clusterId string
	pageSize *int32
	pageToken *string
}

// Optional. Requested page size. Server may return fewer items than requested. If unspecified, server will pick an appropriate default.
func (r ApiBranchServiceListBranchesRequest) PageSize(pageSize int32) ApiBranchServiceListBranchesRequest {
	r.pageSize = &pageSize
	return r
}

// Optional. A token identifying a page of results the server should return.
func (r ApiBranchServiceListBranchesRequest) PageToken(pageToken string) ApiBranchServiceListBranchesRequest {
	r.pageToken = &pageToken
	return r
}

func (r ApiBranchServiceListBranchesRequest) Execute() (*ListBranchesResponse, *http.Response, error) {
	return r.ApiService.BranchServiceListBranchesExecute(r)
}

/*
BranchServiceListBranches Lists information about branches.

 @param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
 @param clusterId Required. The ID of the project to which the clusters belong.
 @return ApiBranchServiceListBranchesRequest
*/
func (a *BranchServiceAPIService) BranchServiceListBranches(ctx context.Context, clusterId string) ApiBranchServiceListBranchesRequest {
	return ApiBranchServiceListBranchesRequest{
		ApiService: a,
		ctx: ctx,
		clusterId: clusterId,
	}
}

// Execute executes the request
//  @return ListBranchesResponse
func (a *BranchServiceAPIService) BranchServiceListBranchesExecute(r ApiBranchServiceListBranchesRequest) (*ListBranchesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod   = http.MethodGet
		localVarPostBody     interface{}
		formFiles            []formFile
		localVarReturnValue  *ListBranchesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "BranchServiceAPIService.BranchServiceListBranches")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/branches"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "", "")
	}
	if r.pageToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageToken", r.pageToken, "", "")
	}
	// to determine the Content-Type header
	localVarHTTPContentTypes := []string{}

	// set Content-Type header
	localVarHTTPContentType := selectHeaderContentType(localVarHTTPContentTypes)
	if localVarHTTPContentType != "" {
		localVarHeaderParams["Content-Type"] = localVarHTTPContentType
	}

	// to determine the Accept header
	localVarHTTPHeaderAccepts := []string{"application/json"}

	// set Accept header
	localVarHTTPHeaderAccept := selectHeaderAccept(localVarHTTPHeaderAccepts)
	if localVarHTTPHeaderAccept != "" {
		localVarHeaderParams["Accept"] = localVarHTTPHeaderAccept
	}
	req, err := a.client.prepareRequest(r.ctx, localVarPath, localVarHTTPMethod, localVarPostBody, localVarHeaderParams, localVarQueryParams, localVarFormParams, formFiles)
	if err != nil {
		return localVarReturnValue, nil, err
	}

	localVarHTTPResponse, err := a.client.callAPI(req)
	if err != nil || localVarHTTPResponse == nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	localVarBody, err := io.ReadAll(localVarHTTPResponse.Body)
	localVarHTTPResponse.Body.Close()
	localVarHTTPResponse.Body = io.NopCloser(bytes.NewBuffer(localVarBody))
	if err != nil {
		return localVarReturnValue, localVarHTTPResponse, err
	}

	if localVarHTTPResponse.StatusCode >= 300 {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: localVarHTTPResponse.Status,
		}
			var v Status
			err = a.client.decode(&v, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
			if err != nil {
				newErr.error = err.Error()
				return localVarReturnValue, localVarHTTPResponse, newErr
			}
					newErr.error = formatErrorMessage(localVarHTTPResponse.Status, &v)
					newErr.model = v
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	err = a.client.decode(&localVarReturnValue, localVarBody, localVarHTTPResponse.Header.Get("Content-Type"))
	if err != nil {
		newErr := &GenericOpenAPIError{
			body:  localVarBody,
			error: err.Error(),
		}
		return localVarReturnValue, localVarHTTPResponse, newErr
	}

	return localVarReturnValue, localVarHTTPResponse, nil
}
