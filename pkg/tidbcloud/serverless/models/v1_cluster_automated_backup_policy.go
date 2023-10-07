// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// V1ClusterAutomatedBackupPolicy Message for automated backup configuration for a cluster.
//
// swagger:model v1ClusterAutomatedBackupPolicy
type V1ClusterAutomatedBackupPolicy struct {

	// Optional. Number of days to retain automated backups.
	RetentionDays int32 `json:"retentionDays,omitempty"`

	// Optional. Cron expression for when automated backups should start.
	StartTime string `json:"startTime,omitempty"`
}

// Validate validates this v1 cluster automated backup policy
func (m *V1ClusterAutomatedBackupPolicy) Validate(formats strfmt.Registry) error {
	return nil
}

// ContextValidate validates this v1 cluster automated backup policy based on context it is used
func (m *V1ClusterAutomatedBackupPolicy) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}

// MarshalBinary interface implementation
func (m *V1ClusterAutomatedBackupPolicy) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1ClusterAutomatedBackupPolicy) UnmarshalBinary(b []byte) error {
	var res V1ClusterAutomatedBackupPolicy
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
