// Code generated by go-swagger; DO NOT EDIT.

package import_service

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"net/http"
	"time"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/runtime"
	cr "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
)

// NewImportServiceAbortMultipartUploadParams creates a new ImportServiceAbortMultipartUploadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportServiceAbortMultipartUploadParams() *ImportServiceAbortMultipartUploadParams {
	return &ImportServiceAbortMultipartUploadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportServiceAbortMultipartUploadParamsWithTimeout creates a new ImportServiceAbortMultipartUploadParams object
// with the ability to set a timeout on a request.
func NewImportServiceAbortMultipartUploadParamsWithTimeout(timeout time.Duration) *ImportServiceAbortMultipartUploadParams {
	return &ImportServiceAbortMultipartUploadParams{
		timeout: timeout,
	}
}

// NewImportServiceAbortMultipartUploadParamsWithContext creates a new ImportServiceAbortMultipartUploadParams object
// with the ability to set a context for a request.
func NewImportServiceAbortMultipartUploadParamsWithContext(ctx context.Context) *ImportServiceAbortMultipartUploadParams {
	return &ImportServiceAbortMultipartUploadParams{
		Context: ctx,
	}
}

// NewImportServiceAbortMultipartUploadParamsWithHTTPClient creates a new ImportServiceAbortMultipartUploadParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportServiceAbortMultipartUploadParamsWithHTTPClient(client *http.Client) *ImportServiceAbortMultipartUploadParams {
	return &ImportServiceAbortMultipartUploadParams{
		HTTPClient: client,
	}
}

/*
ImportServiceAbortMultipartUploadParams contains all the parameters to send to the API endpoint

	for the import service abort multipart upload operation.

	Typically these are written to a http.Request.
*/
type ImportServiceAbortMultipartUploadParams struct {

	/* ClusterID.

	   The ID of the cluster to import into
	*/
	ClusterID string

	/* FileName.

	   The name of the file to import
	*/
	FileName string

	/* UploadID.

	   The ID of the upload
	*/
	UploadID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import service abort multipart upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceAbortMultipartUploadParams) WithDefaults() *ImportServiceAbortMultipartUploadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import service abort multipart upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceAbortMultipartUploadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithTimeout(timeout time.Duration) *ImportServiceAbortMultipartUploadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithContext(ctx context.Context) *ImportServiceAbortMultipartUploadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithHTTPClient(client *http.Client) *ImportServiceAbortMultipartUploadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithClusterID(clusterID string) *ImportServiceAbortMultipartUploadParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithFileName adds the fileName to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithFileName(fileName string) *ImportServiceAbortMultipartUploadParams {
	o.SetFileName(fileName)
	return o
}

// SetFileName adds the fileName to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetFileName(fileName string) {
	o.FileName = fileName
}

// WithUploadID adds the uploadID to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) WithUploadID(uploadID string) *ImportServiceAbortMultipartUploadParams {
	o.SetUploadID(uploadID)
	return o
}

// SetUploadID adds the uploadId to the import service abort multipart upload params
func (o *ImportServiceAbortMultipartUploadParams) SetUploadID(uploadID string) {
	o.UploadID = uploadID
}

// WriteToRequest writes these params to a swagger request
func (o *ImportServiceAbortMultipartUploadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param clusterId
	if err := r.SetPathParam("clusterId", o.ClusterID); err != nil {
		return err
	}

	// query param fileName
	qrFileName := o.FileName
	qFileName := qrFileName
	if qFileName != "" {

		if err := r.SetQueryParam("fileName", qFileName); err != nil {
			return err
		}
	}

	// query param uploadId
	qrUploadID := o.UploadID
	qUploadID := qrUploadID
	if qUploadID != "" {

		if err := r.SetQueryParam("uploadId", qUploadID); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
