// Code generated by go-swagger; DO NOT EDIT.

package account

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

// NewGetV1beta1ProjectsParams creates a new GetV1beta1ProjectsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewGetV1beta1ProjectsParams() *GetV1beta1ProjectsParams {
	return &GetV1beta1ProjectsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewGetV1beta1ProjectsParamsWithTimeout creates a new GetV1beta1ProjectsParams object
// with the ability to set a timeout on a request.
func NewGetV1beta1ProjectsParamsWithTimeout(timeout time.Duration) *GetV1beta1ProjectsParams {
	return &GetV1beta1ProjectsParams{
		timeout: timeout,
	}
}

// NewGetV1beta1ProjectsParamsWithContext creates a new GetV1beta1ProjectsParams object
// with the ability to set a context for a request.
func NewGetV1beta1ProjectsParamsWithContext(ctx context.Context) *GetV1beta1ProjectsParams {
	return &GetV1beta1ProjectsParams{
		Context: ctx,
	}
}

// NewGetV1beta1ProjectsParamsWithHTTPClient creates a new GetV1beta1ProjectsParams object
// with the ability to set a custom HTTPClient for a request.
func NewGetV1beta1ProjectsParamsWithHTTPClient(client *http.Client) *GetV1beta1ProjectsParams {
	return &GetV1beta1ProjectsParams{
		HTTPClient: client,
	}
}

/*
GetV1beta1ProjectsParams contains all the parameters to send to the API endpoint

	for the get v1beta1 projects operation.

	Typically these are written to a http.Request.
*/
type GetV1beta1ProjectsParams struct {

	/* PageSize.

	   The page size of the next page. If `pageSize` is set to 0, it returns 100 records in one page.
	*/
	PageSize *int64

	/* PageToken.

	   The page token of the next page.
	*/
	PageToken *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the get v1beta1 projects params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetV1beta1ProjectsParams) WithDefaults() *GetV1beta1ProjectsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the get v1beta1 projects params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *GetV1beta1ProjectsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) WithTimeout(timeout time.Duration) *GetV1beta1ProjectsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) WithContext(ctx context.Context) *GetV1beta1ProjectsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) WithHTTPClient(client *http.Client) *GetV1beta1ProjectsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithPageSize adds the pageSize to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) WithPageSize(pageSize *int64) *GetV1beta1ProjectsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) SetPageSize(pageSize *int64) {
	o.PageSize = pageSize
}

// WithPageToken adds the pageToken to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) WithPageToken(pageToken *string) *GetV1beta1ProjectsParams {
	o.SetPageToken(pageToken)
	return o
}

// SetPageToken adds the pageToken to the get v1beta1 projects params
func (o *GetV1beta1ProjectsParams) SetPageToken(pageToken *string) {
	o.PageToken = pageToken
}

// WriteToRequest writes these params to a swagger request
func (o *GetV1beta1ProjectsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	if o.PageSize != nil {

		// query param pageSize
		var qrPageSize int64

		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt64(qrPageSize)
		if qPageSize != "" {

			if err := r.SetQueryParam("pageSize", qPageSize); err != nil {
				return err
			}
		}
	}

	if o.PageToken != nil {

		// query param pageToken
		var qrPageToken string

		if o.PageToken != nil {
			qrPageToken = *o.PageToken
		}
		qPageToken := qrPageToken
		if qPageToken != "" {

			if err := r.SetQueryParam("pageToken", qPageToken); err != nil {
				return err
			}
		}
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}
