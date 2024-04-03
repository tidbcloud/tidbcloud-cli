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

	"tidbcloud-cli/pkg/tidbcloud/v1beta1/serverless_import/models"
)

// NewImportServiceCompleteMultipartUploadParams creates a new ImportServiceCompleteMultipartUploadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportServiceCompleteMultipartUploadParams() *ImportServiceCompleteMultipartUploadParams {
	return &ImportServiceCompleteMultipartUploadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportServiceCompleteMultipartUploadParamsWithTimeout creates a new ImportServiceCompleteMultipartUploadParams object
// with the ability to set a timeout on a request.
func NewImportServiceCompleteMultipartUploadParamsWithTimeout(timeout time.Duration) *ImportServiceCompleteMultipartUploadParams {
	return &ImportServiceCompleteMultipartUploadParams{
		timeout: timeout,
	}
}

// NewImportServiceCompleteMultipartUploadParamsWithContext creates a new ImportServiceCompleteMultipartUploadParams object
// with the ability to set a context for a request.
func NewImportServiceCompleteMultipartUploadParamsWithContext(ctx context.Context) *ImportServiceCompleteMultipartUploadParams {
	return &ImportServiceCompleteMultipartUploadParams{
		Context: ctx,
	}
}

// NewImportServiceCompleteMultipartUploadParamsWithHTTPClient creates a new ImportServiceCompleteMultipartUploadParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportServiceCompleteMultipartUploadParamsWithHTTPClient(client *http.Client) *ImportServiceCompleteMultipartUploadParams {
	return &ImportServiceCompleteMultipartUploadParams{
		HTTPClient: client,
	}
}

/*
ImportServiceCompleteMultipartUploadParams contains all the parameters to send to the API endpoint

	for the import service complete multipart upload operation.

	Typically these are written to a http.Request.
*/
type ImportServiceCompleteMultipartUploadParams struct {

	/* ClusterID.

	   The ID of the cluster to import into
	*/
	ClusterID string

	/* FileName.

	   The name of the file to import
	*/
	FileName string

	/* Parts.

	   The parts have been uploaded
	*/
	Parts []*models.V1beta1CompletePart

	/* UploadID.

	   The ID of the upload
	*/
	UploadID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import service complete multipart upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceCompleteMultipartUploadParams) WithDefaults() *ImportServiceCompleteMultipartUploadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import service complete multipart upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceCompleteMultipartUploadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithTimeout(timeout time.Duration) *ImportServiceCompleteMultipartUploadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithContext(ctx context.Context) *ImportServiceCompleteMultipartUploadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithHTTPClient(client *http.Client) *ImportServiceCompleteMultipartUploadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithClusterID(clusterID string) *ImportServiceCompleteMultipartUploadParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithFileName adds the fileName to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithFileName(fileName string) *ImportServiceCompleteMultipartUploadParams {
	o.SetFileName(fileName)
	return o
}

// SetFileName adds the fileName to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetFileName(fileName string) {
	o.FileName = fileName
}

// WithParts adds the parts to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithParts(parts []*models.V1beta1CompletePart) *ImportServiceCompleteMultipartUploadParams {
	o.SetParts(parts)
	return o
}

// SetParts adds the parts to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetParts(parts []*models.V1beta1CompletePart) {
	o.Parts = parts
}

// WithUploadID adds the uploadID to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) WithUploadID(uploadID string) *ImportServiceCompleteMultipartUploadParams {
	o.SetUploadID(uploadID)
	return o
}

// SetUploadID adds the uploadId to the import service complete multipart upload params
func (o *ImportServiceCompleteMultipartUploadParams) SetUploadID(uploadID string) {
	o.UploadID = uploadID
}

// WriteToRequest writes these params to a swagger request
func (o *ImportServiceCompleteMultipartUploadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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
	if o.Parts != nil {
		if err := r.SetBodyParam(o.Parts); err != nil {
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