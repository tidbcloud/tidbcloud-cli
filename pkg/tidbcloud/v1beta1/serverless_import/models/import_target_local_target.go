// Code generated by go-swagger; DO NOT EDIT.

package models

// This file was generated by the swagger tool.
// Editing this file might prove futile when you re-run the swagger generate command

import (
	"context"

	"github.com/go-openapi/errors"
	"github.com/go-openapi/strfmt"
	"github.com/go-openapi/swag"
)

// ImportTargetLocalTarget import target local target
//
// swagger:model ImportTargetLocalTarget
type ImportTargetLocalTarget struct {

	// Optional. The file name to import.
	FileName string `json:"fileName,omitempty"`

	// Optional. The table to import to.
	TargetTable *V1beta1Table `json:"targetTable,omitempty"`

	// Optional. The upload id of import source file.
	UploadID string `json:"uploadId,omitempty"`
}

// Validate validates this import target local target
func (m *ImportTargetLocalTarget) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateTargetTable(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImportTargetLocalTarget) validateTargetTable(formats strfmt.Registry) error {
	if swag.IsZero(m.TargetTable) { // not required
		return nil
	}

	if m.TargetTable != nil {
		if err := m.TargetTable.Validate(formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("targetTable")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("targetTable")
			}
			return err
		}
	}

	return nil
}

// ContextValidate validate this import target local target based on the context it is used
func (m *ImportTargetLocalTarget) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateTargetTable(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *ImportTargetLocalTarget) contextValidateTargetTable(ctx context.Context, formats strfmt.Registry) error {

	if m.TargetTable != nil {

		if swag.IsZero(m.TargetTable) { // not required
			return nil
		}

		if err := m.TargetTable.ContextValidate(ctx, formats); err != nil {
			if ve, ok := err.(*errors.Validation); ok {
				return ve.ValidateName("targetTable")
			} else if ce, ok := err.(*errors.CompositeError); ok {
				return ce.ValidateName("targetTable")
			}
			return err
		}
	}

	return nil
}

// MarshalBinary interface implementation
func (m *ImportTargetLocalTarget) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *ImportTargetLocalTarget) UnmarshalBinary(b []byte) error {
	var res ImportTargetLocalTarget
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
