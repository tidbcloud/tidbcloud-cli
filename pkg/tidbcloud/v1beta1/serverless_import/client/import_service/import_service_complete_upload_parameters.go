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

// NewImportServiceCompleteUploadParams creates a new ImportServiceCompleteUploadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportServiceCompleteUploadParams() *ImportServiceCompleteUploadParams {
	return &ImportServiceCompleteUploadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportServiceCompleteUploadParamsWithTimeout creates a new ImportServiceCompleteUploadParams object
// with the ability to set a timeout on a request.
func NewImportServiceCompleteUploadParamsWithTimeout(timeout time.Duration) *ImportServiceCompleteUploadParams {
	return &ImportServiceCompleteUploadParams{
		timeout: timeout,
	}
}

// NewImportServiceCompleteUploadParamsWithContext creates a new ImportServiceCompleteUploadParams object
// with the ability to set a context for a request.
func NewImportServiceCompleteUploadParamsWithContext(ctx context.Context) *ImportServiceCompleteUploadParams {
	return &ImportServiceCompleteUploadParams{
		Context: ctx,
	}
}

// NewImportServiceCompleteUploadParamsWithHTTPClient creates a new ImportServiceCompleteUploadParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportServiceCompleteUploadParamsWithHTTPClient(client *http.Client) *ImportServiceCompleteUploadParams {
	return &ImportServiceCompleteUploadParams{
		HTTPClient: client,
	}
}

/*
ImportServiceCompleteUploadParams contains all the parameters to send to the API endpoint

	for the import service complete upload operation.

	Typically these are written to a http.Request.
*/
type ImportServiceCompleteUploadParams struct {

	/* ClusterID.

	   The ID of the cluster to import into
	*/
	ClusterID string

	/* Parts.

	   The parts have been uploaded
	*/
	Parts []*models.CompletePart

	/* UploadID.

	   The ID of the upload
	*/
	UploadID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import service complete upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceCompleteUploadParams) WithDefaults() *ImportServiceCompleteUploadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import service complete upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceCompleteUploadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithTimeout(timeout time.Duration) *ImportServiceCompleteUploadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithContext(ctx context.Context) *ImportServiceCompleteUploadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithHTTPClient(client *http.Client) *ImportServiceCompleteUploadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithClusterID(clusterID string) *ImportServiceCompleteUploadParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithParts adds the parts to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithParts(parts []*models.CompletePart) *ImportServiceCompleteUploadParams {
	o.SetParts(parts)
	return o
}

// SetParts adds the parts to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetParts(parts []*models.CompletePart) {
	o.Parts = parts
}

// WithUploadID adds the uploadID to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) WithUploadID(uploadID string) *ImportServiceCompleteUploadParams {
	o.SetUploadID(uploadID)
	return o
}

// SetUploadID adds the uploadId to the import service complete upload params
func (o *ImportServiceCompleteUploadParams) SetUploadID(uploadID string) {
	o.UploadID = uploadID
}

// WriteToRequest writes these params to a swagger request
func (o *ImportServiceCompleteUploadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// path param clusterId
	if err := r.SetPathParam("clusterId", o.ClusterID); err != nil {
		return err
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
