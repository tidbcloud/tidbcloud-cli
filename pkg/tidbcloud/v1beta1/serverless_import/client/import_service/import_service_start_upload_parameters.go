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
	"github.com/go-openapi/swag"
)

// NewImportServiceStartUploadParams creates a new ImportServiceStartUploadParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewImportServiceStartUploadParams() *ImportServiceStartUploadParams {
	return &ImportServiceStartUploadParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewImportServiceStartUploadParamsWithTimeout creates a new ImportServiceStartUploadParams object
// with the ability to set a timeout on a request.
func NewImportServiceStartUploadParamsWithTimeout(timeout time.Duration) *ImportServiceStartUploadParams {
	return &ImportServiceStartUploadParams{
		timeout: timeout,
	}
}

// NewImportServiceStartUploadParamsWithContext creates a new ImportServiceStartUploadParams object
// with the ability to set a context for a request.
func NewImportServiceStartUploadParamsWithContext(ctx context.Context) *ImportServiceStartUploadParams {
	return &ImportServiceStartUploadParams{
		Context: ctx,
	}
}

// NewImportServiceStartUploadParamsWithHTTPClient creates a new ImportServiceStartUploadParams object
// with the ability to set a custom HTTPClient for a request.
func NewImportServiceStartUploadParamsWithHTTPClient(client *http.Client) *ImportServiceStartUploadParams {
	return &ImportServiceStartUploadParams{
		HTTPClient: client,
	}
}

/*
ImportServiceStartUploadParams contains all the parameters to send to the API endpoint

	for the import service start upload operation.

	Typically these are written to a http.Request.
*/
type ImportServiceStartUploadParams struct {

	/* ClusterID.

	   The ID of the cluster to import into
	*/
	ClusterID string

	/* FileName.

	   The name of the file to import
	*/
	FileName string

	/* PartNumber.

	   The number of parts to split the file into

	   Format: int32
	*/
	PartNumber int32

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the import service start upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceStartUploadParams) WithDefaults() *ImportServiceStartUploadParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the import service start upload params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *ImportServiceStartUploadParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the import service start upload params
func (o *ImportServiceStartUploadParams) WithTimeout(timeout time.Duration) *ImportServiceStartUploadParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the import service start upload params
func (o *ImportServiceStartUploadParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the import service start upload params
func (o *ImportServiceStartUploadParams) WithContext(ctx context.Context) *ImportServiceStartUploadParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the import service start upload params
func (o *ImportServiceStartUploadParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the import service start upload params
func (o *ImportServiceStartUploadParams) WithHTTPClient(client *http.Client) *ImportServiceStartUploadParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the import service start upload params
func (o *ImportServiceStartUploadParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the import service start upload params
func (o *ImportServiceStartUploadParams) WithClusterID(clusterID string) *ImportServiceStartUploadParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the import service start upload params
func (o *ImportServiceStartUploadParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithFileName adds the fileName to the import service start upload params
func (o *ImportServiceStartUploadParams) WithFileName(fileName string) *ImportServiceStartUploadParams {
	o.SetFileName(fileName)
	return o
}

// SetFileName adds the fileName to the import service start upload params
func (o *ImportServiceStartUploadParams) SetFileName(fileName string) {
	o.FileName = fileName
}

// WithPartNumber adds the partNumber to the import service start upload params
func (o *ImportServiceStartUploadParams) WithPartNumber(partNumber int32) *ImportServiceStartUploadParams {
	o.SetPartNumber(partNumber)
	return o
}

// SetPartNumber adds the partNumber to the import service start upload params
func (o *ImportServiceStartUploadParams) SetPartNumber(partNumber int32) {
	o.PartNumber = partNumber
}

// WriteToRequest writes these params to a swagger request
func (o *ImportServiceStartUploadParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

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

	// query param partNumber
	qrPartNumber := o.PartNumber
	qPartNumber := swag.FormatInt32(qrPartNumber)
	if qPartNumber != "" {

		if err := r.SetQueryParam("partNumber", qPartNumber); err != nil {
			return err
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
