// Code generated by go-swagger; DO NOT EDIT.

package export_service

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

// NewExportServiceDownloadExportParams creates a new ExportServiceDownloadExportParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewExportServiceDownloadExportParams() *ExportServiceDownloadExportParams {
	return &ExportServiceDownloadExportParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewExportServiceDownloadExportParamsWithTimeout creates a new ExportServiceDownloadExportParams object
// with the ability to set a timeout on a request.
func NewExportServiceDownloadExportParamsWithTimeout(timeout time.Duration) *ExportServiceDownloadExportParams {
	return &ExportServiceDownloadExportParams{
		timeout: timeout,
	}
}

// NewExportServiceDownloadExportParamsWithContext creates a new ExportServiceDownloadExportParams object
// with the ability to set a context for a request.
func NewExportServiceDownloadExportParamsWithContext(ctx context.Context) *ExportServiceDownloadExportParams {
	return &ExportServiceDownloadExportParams{
		Context: ctx,
	}
}

// NewExportServiceDownloadExportParamsWithHTTPClient creates a new ExportServiceDownloadExportParams object
// with the ability to set a custom HTTPClient for a request.
func NewExportServiceDownloadExportParamsWithHTTPClient(client *http.Client) *ExportServiceDownloadExportParams {
	return &ExportServiceDownloadExportParams{
		HTTPClient: client,
	}
}

/*
ExportServiceDownloadExportParams contains all the parameters to send to the API endpoint

	for the export service download export operation.

	Typically these are written to a http.Request.
*/
type ExportServiceDownloadExportParams struct {

	// Body.
	Body interface{}

	/* ClusterID.

	   Required. The ID of the cluster.
	*/
	ClusterID string

	/* ExportID.

	   Required. The ID of the export to be retrieved.
	*/
	ExportID string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the export service download export params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ExportServiceDownloadExportParams) WithDefaults() *ExportServiceDownloadExportParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the export service download export params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ExportServiceDownloadExportParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the export service download export params
func (o *ExportServiceDownloadExportParams) WithTimeout(timeout time.Duration) *ExportServiceDownloadExportParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the export service download export params
func (o *ExportServiceDownloadExportParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the export service download export params
func (o *ExportServiceDownloadExportParams) WithContext(ctx context.Context) *ExportServiceDownloadExportParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the export service download export params
func (o *ExportServiceDownloadExportParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the export service download export params
func (o *ExportServiceDownloadExportParams) WithHTTPClient(client *http.Client) *ExportServiceDownloadExportParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the export service download export params
func (o *ExportServiceDownloadExportParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the export service download export params
func (o *ExportServiceDownloadExportParams) WithBody(body interface{}) *ExportServiceDownloadExportParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the export service download export params
func (o *ExportServiceDownloadExportParams) SetBody(body interface{}) {
	o.Body = body
}

// WithClusterID adds the clusterID to the export service download export params
func (o *ExportServiceDownloadExportParams) WithClusterID(clusterID string) *ExportServiceDownloadExportParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the export service download export params
func (o *ExportServiceDownloadExportParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithExportID adds the exportID to the export service download export params
func (o *ExportServiceDownloadExportParams) WithExportID(exportID string) *ExportServiceDownloadExportParams {
	o.SetExportID(exportID)
	return o
}

// SetExportID adds the exportId to the export service download export params
func (o *ExportServiceDownloadExportParams) SetExportID(exportID string) {
	o.ExportID = exportID
}

// WriteToRequest writes these params to a swagger request
func (o *ExportServiceDownloadExportParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error
	if o.Body != nil {
		if err := r.SetBodyParam(o.Body); err != nil {
			return err
		}
	}

	// path param clusterId
	if err := r.SetPathParam("clusterId", o.ClusterID); err != nil {
		return err
	}

	// path param exportId
	if err := r.SetPathParam("exportId", o.ExportID); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}