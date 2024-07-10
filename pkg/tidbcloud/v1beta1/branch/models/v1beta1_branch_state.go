// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"
	"encoding/json"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/validate"
)

// V1beta1BranchState Output Only. Branch State.
//
//   - STATE_UNSPECIFIED: The state of the branch is unknown.
//   - CREATING: The branch is being created.
//   - ACTIVE: The branch is active and running.
//   - DELETED: The branch is being deleted.
//   - MAINTENANCE: The branch is under maintenance.
//
// swagger:model v1beta1BranchState
type V1beta1BranchState string

func NewV1beta1BranchState(value V1beta1BranchState) *V1beta1BranchState {
	return &value
}

// Pointer returns a pointer to a freshly-allocated V1beta1BranchState.
func (m V1beta1BranchState) Pointer() *V1beta1BranchState {
	return &m
}

const (

	// V1beta1BranchStateCREATING captures enum value "CREATING"
	V1beta1BranchStateCREATING V1beta1BranchState = "CREATING"

	// V1beta1BranchStateACTIVE captures enum value "ACTIVE"
	V1beta1BranchStateACTIVE V1beta1BranchState = "ACTIVE"

	// V1beta1BranchStateDELETED captures enum value "DELETED"
	V1beta1BranchStateDELETED V1beta1BranchState = "DELETED"

	// V1beta1BranchStateMAINTENANCE captures enum value "MAINTENANCE"
	V1beta1BranchStateMAINTENANCE V1beta1BranchState = "MAINTENANCE"
)

// for schema
var v1beta1BranchStateEnum []interface{}

func init() {
	var res []V1beta1BranchState
	if err := json.Unmarshal([]byte(`["CREATING","ACTIVE","DELETED","MAINTENANCE"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		v1beta1BranchStateEnum = append(v1beta1BranchStateEnum, v)
	}
}

func (m V1beta1BranchState) validateV1beta1BranchStateEnum(path, location string, value V1beta1BranchState) error {
	if err := validate.EnumCase(path, location, value, v1beta1BranchStateEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this v1beta1 branch state
func (m V1beta1BranchState) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateV1beta1BranchStateEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this v1beta1 branch state based on context it is used
func (m V1beta1BranchState) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}