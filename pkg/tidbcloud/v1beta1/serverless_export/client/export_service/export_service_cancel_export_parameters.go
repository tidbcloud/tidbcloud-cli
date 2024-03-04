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

// NewExportServiceCancelExportParams creates a new ExportServiceCancelExportParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewExportServiceCancelExportParams() *ExportServiceCancelExportParams {
	return &ExportServiceCancelExportParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewExportServiceCancelExportParamsWithTimeout creates a new ExportServiceCancelExportParams object
// with the ability to set a timeout on a request.
func NewExportServiceCancelExportParamsWithTimeout(timeout time.Duration) *ExportServiceCancelExportParams {
	return &ExportServiceCancelExportParams{
		timeout: timeout,
	}
}

// NewExportServiceCancelExportParamsWithContext creates a new ExportServiceCancelExportParams object
// with the ability to set a context for a request.
func NewExportServiceCancelExportParamsWithContext(ctx context.Context) *ExportServiceCancelExportParams {
	return &ExportServiceCancelExportParams{
		Context: ctx,
	}
}

// NewExportServiceCancelExportParamsWithHTTPClient creates a new ExportServiceCancelExportParams object
// with the ability to set a custom HTTPClient for a request.
func NewExportServiceCancelExportParamsWithHTTPClient(client *http.Client) *ExportServiceCancelExportParams {
	return &ExportServiceCancelExportParams{
		HTTPClient: client,
	}
}

/*
ExportServiceCancelExportParams contains all the parameters to send to the API endpoint

	for the export service cancel export operation.

	Typically these are written to a http.Request.
*/
type ExportServiceCancelExportParams struct {

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

// WithDefaults hydrates default values in the export service cancel export params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ExportServiceCancelExportParams) WithDefaults() *ExportServiceCancelExportParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the export service cancel export params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ExportServiceCancelExportParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithTimeout(timeout time.Duration) *ExportServiceCancelExportParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithContext(ctx context.Context) *ExportServiceCancelExportParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithHTTPClient(client *http.Client) *ExportServiceCancelExportParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithBody adds the body to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithBody(body interface{}) *ExportServiceCancelExportParams {
	o.SetBody(body)
	return o
}

// SetBody adds the body to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetBody(body interface{}) {
	o.Body = body
}

// WithClusterID adds the clusterID to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithClusterID(clusterID string) *ExportServiceCancelExportParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithExportID adds the exportID to the export service cancel export params
func (o *ExportServiceCancelExportParams) WithExportID(exportID string) *ExportServiceCancelExportParams {
	o.SetExportID(exportID)
	return o
}

// SetExportID adds the exportId to the export service cancel export params
func (o *ExportServiceCancelExportParams) SetExportID(exportID string) {
	o.ExportID = exportID
}

// WriteToRequest writes these params to a swagger request
func (o *ExportServiceCancelExportParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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