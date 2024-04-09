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

// V1beta1ExportOptions v1beta1 export options
//
// swagger:model v1beta1ExportOptions
type V1beta1ExportOptions struct {

	// Optional. The compression of the export.
	Compression ExportOptionsCompressionType `json:"compression,omitempty"`

	// Optional. The specify database of the export.
	Database string `json:"database,omitempty"`

	// Optional. The exported file type.
	FileType ExportOptionsFileType `json:"fileType,omitempty"`

	// Optional. The specify table of the export.
	Table string `json:"table,omitempty"`
}

// Validate validates this v1beta1 export options
func (m *V1beta1ExportOptions) Validate(formats strfmt.Registry) error {
	var res []error

	if err := m.validateCompression(formats); err != nil {
		res = append(res, err)
	}

	if err := m.validateFileType(formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ExportOptions) validateCompression(formats strfmt.Registry) error {
	if swag.IsZero(m.Compression) { // not required
		return nil
	}

	if err := m.Compression.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("compression")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("compression")
		}
		return err
	}

	return nil
}

func (m *V1beta1ExportOptions) validateFileType(formats strfmt.Registry) error {
	if swag.IsZero(m.FileType) { // not required
		return nil
	}

	if err := m.FileType.Validate(formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("fileType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("fileType")
		}
		return err
	}

	return nil
}

// ContextValidate validate this v1beta1 export options based on the context it is used
func (m *V1beta1ExportOptions) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	var res []error

	if err := m.contextValidateCompression(ctx, formats); err != nil {
		res = append(res, err)
	}

	if err := m.contextValidateFileType(ctx, formats); err != nil {
		res = append(res, err)
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

func (m *V1beta1ExportOptions) contextValidateCompression(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.Compression) { // not required
		return nil
	}

	if err := m.Compression.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("compression")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("compression")
		}
		return err
	}

	return nil
}

func (m *V1beta1ExportOptions) contextValidateFileType(ctx context.Context, formats strfmt.Registry) error {

	if swag.IsZero(m.FileType) { // not required
		return nil
	}

	if err := m.FileType.ContextValidate(ctx, formats); err != nil {
		if ve, ok := err.(*errors.Validation); ok {
			return ve.ValidateName("fileType")
		} else if ce, ok := err.(*errors.CompositeError); ok {
			return ce.ValidateName("fileType")
		}
		return err
	}

	return nil
}

// MarshalBinary interface implementation
func (m *V1beta1ExportOptions) MarshalBinary() ([]byte, error) {
	if m == nil {
		return nil, nil
	}
	return swag.WriteJSON(m)
}

// UnmarshalBinary interface implementation
func (m *V1beta1ExportOptions) UnmarshalBinary(b []byte) error {
	var res V1beta1ExportOptions
	if err := swag.ReadJSON(b, &res); err != nil {
		return err
	}
	*m = res
	return nil
}
