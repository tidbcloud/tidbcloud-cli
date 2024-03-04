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

// TargetTargetType  - TARGET_UNSPECIFIED: The target of the export is unknown.
//   - LOCAL: Local target.
//   - S3: S3 target.
//
// swagger:model TargetTargetType
type TargetTargetType string

func NewTargetTargetType(value TargetTargetType) *TargetTargetType {
	return &value
}

// Pointer returns a pointer to a freshly-allocated TargetTargetType.
func (m TargetTargetType) Pointer() *TargetTargetType {
	return &m
}

const (

	// TargetTargetTypeLOCAL captures enum value "LOCAL"
	TargetTargetTypeLOCAL TargetTargetType = "LOCAL"

	// TargetTargetTypeS3 captures enum value "S3"
	TargetTargetTypeS3 TargetTargetType = "S3"
)

// for schema
var targetTargetTypeEnum []interface{}

func init() {
	var res []TargetTargetType
	if err := json.Unmarshal([]byte(`["LOCAL","S3"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		targetTargetTypeEnum = append(targetTargetTypeEnum, v)
	}
}

func (m TargetTargetType) validateTargetTargetTypeEnum(path, location string, value TargetTargetType) error {
	if err := validate.EnumCase(path, location, value, targetTargetTypeEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this target target type
func (m TargetTargetType) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateTargetTargetTypeEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this target target type based on context it is used
func (m TargetTargetType) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
