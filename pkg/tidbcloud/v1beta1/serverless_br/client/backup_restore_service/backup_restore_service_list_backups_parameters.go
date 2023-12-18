// Code generated by go-swagger; DO NOT EDIT.

package backup_restore_service

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

// NewBackupRestoreServiceListBackupsParams creates a new BackupRestoreServiceListBackupsParams object,
// with the default timeout for this client.
//
// Default values are not hydrated, since defaults are normally applied by the API server side.
//
// To enforce default values in parameter, use SetDefaults or WithDefaults.
func NewBackupRestoreServiceListBackupsParams() *BackupRestoreServiceListBackupsParams {
	return &BackupRestoreServiceListBackupsParams{
		timeout: cr.DefaultTimeout,
	}
}

// NewBackupRestoreServiceListBackupsParamsWithTimeout creates a new BackupRestoreServiceListBackupsParams object
// with the ability to set a timeout on a request.
func NewBackupRestoreServiceListBackupsParamsWithTimeout(timeout time.Duration) *BackupRestoreServiceListBackupsParams {
	return &BackupRestoreServiceListBackupsParams{
		timeout: timeout,
	}
}

// NewBackupRestoreServiceListBackupsParamsWithContext creates a new BackupRestoreServiceListBackupsParams object
// with the ability to set a context for a request.
func NewBackupRestoreServiceListBackupsParamsWithContext(ctx context.Context) *BackupRestoreServiceListBackupsParams {
	return &BackupRestoreServiceListBackupsParams{
		Context: ctx,
	}
}

// NewBackupRestoreServiceListBackupsParamsWithHTTPClient creates a new BackupRestoreServiceListBackupsParams object
// with the ability to set a custom HTTPClient for a request.
func NewBackupRestoreServiceListBackupsParamsWithHTTPClient(client *http.Client) *BackupRestoreServiceListBackupsParams {
	return &BackupRestoreServiceListBackupsParams{
		HTTPClient: client,
	}
}

/*
BackupRestoreServiceListBackupsParams contains all the parameters to send to the API endpoint

	for the backup restore service list backups operation.

	Typically these are written to a http.Request.
*/
type BackupRestoreServiceListBackupsParams struct {

	/* ClusterID.

	   Required. The cluster ID to list backups for.
	*/
	ClusterID string

	/* PageSize.

	   Optional. The maximum number of clusters to return.

	   Format: int32
	*/
	PageSize *int32

	/* PageToken.

	   Optional. The page token from the previous response for pagination.
	*/
	PageToken *string

	timeout    time.Duration
	Context    context.Context
	HTTPClient *http.Client
}

// WithDefaults hydrates default values in the backup restore service list backups params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BackupRestoreServiceListBackupsParams) WithDefaults() *BackupRestoreServiceListBackupsParams {
	o.SetDefaults()
	return o
}

// SetDefaults hydrates default values in the backup restore service list backups params (not the query body).
//
// All values with no default are reset to their zero value.
func (o *BackupRestoreServiceListBackupsParams) SetDefaults() {
	// no default values defined for this parameter
}

// WithTimeout adds the timeout to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithTimeout(timeout time.Duration) *BackupRestoreServiceListBackupsParams {
	o.SetTimeout(timeout)
	return o
}

// SetTimeout adds the timeout to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetTimeout(timeout time.Duration) {
	o.timeout = timeout
}

// WithContext adds the context to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithContext(ctx context.Context) *BackupRestoreServiceListBackupsParams {
	o.SetContext(ctx)
	return o
}

// SetContext adds the context to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetContext(ctx context.Context) {
	o.Context = ctx
}

// WithHTTPClient adds the HTTPClient to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithHTTPClient(client *http.Client) *BackupRestoreServiceListBackupsParams {
	o.SetHTTPClient(client)
	return o
}

// SetHTTPClient adds the HTTPClient to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetHTTPClient(client *http.Client) {
	o.HTTPClient = client
}

// WithClusterID adds the clusterID to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithClusterID(clusterID string) *BackupRestoreServiceListBackupsParams {
	o.SetClusterID(clusterID)
	return o
}

// SetClusterID adds the clusterId to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetClusterID(clusterID string) {
	o.ClusterID = clusterID
}

// WithPageSize adds the pageSize to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithPageSize(pageSize *int32) *BackupRestoreServiceListBackupsParams {
	o.SetPageSize(pageSize)
	return o
}

// SetPageSize adds the pageSize to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetPageSize(pageSize *int32) {
	o.PageSize = pageSize
}

// WithPageToken adds the pageToken to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) WithPageToken(pageToken *string) *BackupRestoreServiceListBackupsParams {
	o.SetPageToken(pageToken)
	return o
}

// SetPageToken adds the pageToken to the backup restore service list backups params
func (o *BackupRestoreServiceListBackupsParams) SetPageToken(pageToken *string) {
	o.PageToken = pageToken
}

// WriteToRequest writes these params to a swagger request
func (o *BackupRestoreServiceListBackupsParams) WriteToRequest(r runtime.ClientRequest, reg strfmt.Registry) error {

	if err := r.SetTimeout(o.timeout); err != nil {
		return err
	}
	var res []error

	// query param clusterId
	qrClusterID := o.ClusterID
	qClusterID := qrClusterID
	if qClusterID != "" {

		if err := r.SetQueryParam("clusterId", qClusterID); err != nil {
			return err
		}
	}

	if o.PageSize != nil {

		// query param pageSize
		var qrPageSize int32

		if o.PageSize != nil {
			qrPageSize = *o.PageSize
		}
		qPageSize := swag.FormatInt32(qrPageSize)
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
