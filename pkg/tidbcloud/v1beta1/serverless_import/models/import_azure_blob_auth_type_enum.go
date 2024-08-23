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

// ImportAzureBlobAuthTypeEnum  - SAS_TOKEN: The access method is sas token.
//
// swagger:model ImportAzureBlobAuthType.Enum
type ImportAzureBlobAuthTypeEnum string

func NewImportAzureBlobAuthTypeEnum(value ImportAzureBlobAuthTypeEnum) *ImportAzureBlobAuthTypeEnum {
	return &value
}

// Pointer returns a pointer to a freshly-allocated ImportAzureBlobAuthTypeEnum.
func (m ImportAzureBlobAuthTypeEnum) Pointer() *ImportAzureBlobAuthTypeEnum {
	return &m
}

const (

	// ImportAzureBlobAuthTypeEnumSASTOKEN captures enum value "SAS_TOKEN"
	ImportAzureBlobAuthTypeEnumSASTOKEN ImportAzureBlobAuthTypeEnum = "SAS_TOKEN"
)

// for schema
var importAzureBlobAuthTypeEnumEnum []interface{}

func init() {
	var res []ImportAzureBlobAuthTypeEnum
	if err := json.Unmarshal([]byte(`["SAS_TOKEN"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		importAzureBlobAuthTypeEnumEnum = append(importAzureBlobAuthTypeEnumEnum, v)
	}
}

func (m ImportAzureBlobAuthTypeEnum) validateImportAzureBlobAuthTypeEnumEnum(path, location string, value ImportAzureBlobAuthTypeEnum) error {
	if err := validate.EnumCase(path, location, value, importAzureBlobAuthTypeEnumEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this import azure blob auth type enum
func (m ImportAzureBlobAuthTypeEnum) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateImportAzureBlobAuthTypeEnumEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this import azure blob auth type enum based on context it is used
func (m ImportAzureBlobAuthTypeEnum) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
