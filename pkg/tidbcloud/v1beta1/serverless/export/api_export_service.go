/*
TiDB Cloud Serverless Export Open API

TiDB Cloud Serverless Export Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package export

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// ExportServiceAPIService ExportServiceAPI service
type ExportServiceAPIService service

type ApiExportServiceCancelExportRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	exportId   string
}

func (r ApiExportServiceCancelExportRequest) Execute() (*Export, *http.Response, error) {
	return r.ApiService.ExportServiceCancelExportExecute(r)
}

/*
ExportServiceCancelExport Cancel a specific export job.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export to be retrieved.
	@return ApiExportServiceCancelExportRequest
*/
func (a *ExportServiceAPIService) ExportServiceCancelExport(ctx context.Context, clusterId string, exportId string) ApiExportServiceCancelExportRequest {
	return ApiExportServiceCancelExportRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return Export
func (a *ExportServiceAPIService) ExportServiceCancelExportExecute(r ApiExportServiceCancelExportRequest) (*Export, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Export
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceCancelExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}:cancel"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

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

type ApiExportServiceCreateExportRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	body       *ExportServiceCreateExportBody
}

func (r ApiExportServiceCreateExportRequest) Body(body ExportServiceCreateExportBody) ApiExportServiceCreateExportRequest {
	r.body = &body
	return r
}

func (r ApiExportServiceCreateExportRequest) Execute() (*Export, *http.Response, error) {
	return r.ApiService.ExportServiceCreateExportExecute(r)
}

/*
ExportServiceCreateExport Create an export job

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@return ApiExportServiceCreateExportRequest
*/
func (a *ExportServiceAPIService) ExportServiceCreateExport(ctx context.Context, clusterId string) ApiExportServiceCreateExportRequest {
	return ApiExportServiceCreateExportRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return Export
func (a *ExportServiceAPIService) ExportServiceCreateExportExecute(r ApiExportServiceCreateExportRequest) (*Export, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Export
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceCreateExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
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
	localVarPostBody = r.body
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

type ApiExportServiceDeleteExportRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	exportId   string
}

func (r ApiExportServiceDeleteExportRequest) Execute() (*Export, *http.Response, error) {
	return r.ApiService.ExportServiceDeleteExportExecute(r)
}

/*
ExportServiceDeleteExport Delete an export job

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export to be retrieved.
	@return ApiExportServiceDeleteExportRequest
*/
func (a *ExportServiceAPIService) ExportServiceDeleteExport(ctx context.Context, clusterId string, exportId string) ApiExportServiceDeleteExportRequest {
	return ApiExportServiceDeleteExportRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return Export
func (a *ExportServiceAPIService) ExportServiceDeleteExportExecute(r ApiExportServiceDeleteExportRequest) (*Export, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodDelete
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Export
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceDeleteExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

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

type ApiExportServiceDownloadExportRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	exportId   string
	body       *map[string]interface{}
}

func (r ApiExportServiceDownloadExportRequest) Body(body map[string]interface{}) ApiExportServiceDownloadExportRequest {
	r.body = &body
	return r
}

func (r ApiExportServiceDownloadExportRequest) Execute() (*DownloadExportsResponse, *http.Response, error) {
	return r.ApiService.ExportServiceDownloadExportExecute(r)
}

/*
ExportServiceDownloadExport Generate download url

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export to be retrieved.
	@return ApiExportServiceDownloadExportRequest

Deprecated
*/
func (a *ExportServiceAPIService) ExportServiceDownloadExport(ctx context.Context, clusterId string, exportId string) ApiExportServiceDownloadExportRequest {
	return ApiExportServiceDownloadExportRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return DownloadExportsResponse
//
// Deprecated
func (a *ExportServiceAPIService) ExportServiceDownloadExportExecute(r ApiExportServiceDownloadExportRequest) (*DownloadExportsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DownloadExportsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceDownloadExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}:download"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
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
	localVarPostBody = r.body
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

type ApiExportServiceDownloadExportFilesRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	exportId   string
	body       *ExportServiceDownloadExportFilesBody
}

func (r ApiExportServiceDownloadExportFilesRequest) Body(body ExportServiceDownloadExportFilesBody) ApiExportServiceDownloadExportFilesRequest {
	r.body = &body
	return r
}

func (r ApiExportServiceDownloadExportFilesRequest) Execute() (*DownloadExportFilesResponse, *http.Response, error) {
	return r.ApiService.ExportServiceDownloadExportFilesExecute(r)
}

/*
ExportServiceDownloadExportFiles Generate export files download url

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export.
	@return ApiExportServiceDownloadExportFilesRequest
*/
func (a *ExportServiceAPIService) ExportServiceDownloadExportFiles(ctx context.Context, clusterId string, exportId string) ApiExportServiceDownloadExportFilesRequest {
	return ApiExportServiceDownloadExportFilesRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return DownloadExportFilesResponse
func (a *ExportServiceAPIService) ExportServiceDownloadExportFilesExecute(r ApiExportServiceDownloadExportFilesRequest) (*DownloadExportFilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodPost
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *DownloadExportFilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceDownloadExportFiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}/files:download"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}
	if r.body == nil {
		return localVarReturnValue, nil, reportError("body is required and must be specified")
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
	localVarPostBody = r.body
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

type ApiExportServiceGetExportRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	exportId   string
}

func (r ApiExportServiceGetExportRequest) Execute() (*Export, *http.Response, error) {
	return r.ApiService.ExportServiceGetExportExecute(r)
}

/*
ExportServiceGetExport Retrieves details of an export job.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export to be retrieved.
	@return ApiExportServiceGetExportRequest
*/
func (a *ExportServiceAPIService) ExportServiceGetExport(ctx context.Context, clusterId string, exportId string) ApiExportServiceGetExportRequest {
	return ApiExportServiceGetExportRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return Export
func (a *ExportServiceAPIService) ExportServiceGetExportExecute(r ApiExportServiceGetExportRequest) (*Export, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *Export
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceGetExport")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

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

type ApiExportServiceListExportFilesRequest struct {
	ctx         context.Context
	ApiService  *ExportServiceAPIService
	clusterId   string
	exportId    string
	pageSize    *int32
	pageToken   *string
	generateUrl *bool
}

// Optional. The maximum number to return.
func (r ApiExportServiceListExportFilesRequest) PageSize(pageSize int32) ApiExportServiceListExportFilesRequest {
	r.pageSize = &pageSize
	return r
}

// Optional. The page token from the previous response for pagination.
func (r ApiExportServiceListExportFilesRequest) PageToken(pageToken string) ApiExportServiceListExportFilesRequest {
	r.pageToken = &pageToken
	return r
}

// Optional. Whether return the file download urls, default is false
func (r ApiExportServiceListExportFilesRequest) GenerateUrl(generateUrl bool) ApiExportServiceListExportFilesRequest {
	r.generateUrl = &generateUrl
	return r
}

func (r ApiExportServiceListExportFilesRequest) Execute() (*ListExportFilesResponse, *http.Response, error) {
	return r.ApiService.ExportServiceListExportFilesExecute(r)
}

/*
ExportServiceListExportFiles List export files

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The ID of the cluster.
	@param exportId Required. The ID of the export.
	@return ApiExportServiceListExportFilesRequest
*/
func (a *ExportServiceAPIService) ExportServiceListExportFiles(ctx context.Context, clusterId string, exportId string) ApiExportServiceListExportFilesRequest {
	return ApiExportServiceListExportFilesRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
		exportId:   exportId,
	}
}

// Execute executes the request
//
//	@return ListExportFilesResponse
func (a *ExportServiceAPIService) ExportServiceListExportFilesExecute(r ApiExportServiceListExportFilesRequest) (*ListExportFilesResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ListExportFilesResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceListExportFiles")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports/{exportId}/files"
	localVarPath = strings.Replace(localVarPath, "{"+"clusterId"+"}", url.PathEscape(parameterValueToString(r.clusterId, "clusterId")), -1)
	localVarPath = strings.Replace(localVarPath, "{"+"exportId"+"}", url.PathEscape(parameterValueToString(r.exportId, "exportId")), -1)

	localVarHeaderParams := make(map[string]string)
	localVarQueryParams := url.Values{}
	localVarFormParams := url.Values{}

	if r.pageSize != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageSize", r.pageSize, "", "")
	}
	if r.pageToken != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "pageToken", r.pageToken, "", "")
	}
	if r.generateUrl != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "generateUrl", r.generateUrl, "", "")
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

type ApiExportServiceListExportsRequest struct {
	ctx        context.Context
	ApiService *ExportServiceAPIService
	clusterId  string
	pageSize   *int32
	pageToken  *string
	orderBy    *string
}

// Optional. The maximum number of clusters to return. Default is 10.
func (r ApiExportServiceListExportsRequest) PageSize(pageSize int32) ApiExportServiceListExportsRequest {
	r.pageSize = &pageSize
	return r
}

// Optional. The page token from the previous response for pagination.
func (r ApiExportServiceListExportsRequest) PageToken(pageToken string) ApiExportServiceListExportsRequest {
	r.pageToken = &pageToken
	return r
}

// Optional. List exports order by, separated by comma, default is ascending. Example: \&quot;foo, bar desc\&quot;. Supported field: create_time
func (r ApiExportServiceListExportsRequest) OrderBy(orderBy string) ApiExportServiceListExportsRequest {
	r.orderBy = &orderBy
	return r
}

func (r ApiExportServiceListExportsRequest) Execute() (*ListExportsResponse, *http.Response, error) {
	return r.ApiService.ExportServiceListExportsExecute(r)
}

/*
ExportServiceListExports Provides a list of export jobs.

	@param ctx context.Context - for authentication, logging, cancellation, deadlines, tracing, etc. Passed from http.Request or context.Background().
	@param clusterId Required. The cluster ID to list exports for.
	@return ApiExportServiceListExportsRequest
*/
func (a *ExportServiceAPIService) ExportServiceListExports(ctx context.Context, clusterId string) ApiExportServiceListExportsRequest {
	return ApiExportServiceListExportsRequest{
		ApiService: a,
		ctx:        ctx,
		clusterId:  clusterId,
	}
}

// Execute executes the request
//
//	@return ListExportsResponse
func (a *ExportServiceAPIService) ExportServiceListExportsExecute(r ApiExportServiceListExportsRequest) (*ListExportsResponse, *http.Response, error) {
	var (
		localVarHTTPMethod  = http.MethodGet
		localVarPostBody    interface{}
		formFiles           []formFile
		localVarReturnValue *ListExportsResponse
	)

	localBasePath, err := a.client.cfg.ServerURLWithContext(r.ctx, "ExportServiceAPIService.ExportServiceListExports")
	if err != nil {
		return localVarReturnValue, nil, &GenericOpenAPIError{error: err.Error()}
	}

	localVarPath := localBasePath + "/v1beta1/clusters/{clusterId}/exports"
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
	if r.orderBy != nil {
		parameterAddToHeaderOrQuery(localVarQueryParams, "orderBy", r.orderBy, "", "")
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
