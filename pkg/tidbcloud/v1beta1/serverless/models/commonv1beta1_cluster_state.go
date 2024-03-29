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

// Commonv1beta1ClusterState Enum for current state of a cluster.
//
//   - CREATING: Cluster is being created.
//   - DELETING: Cluster is being deleted.
//   - ACTIVE: Cluster is active for use.
//   - RESTORING: Cluster data is being restored.
//   - MAINTENANCE: Cluster is under maintenance.
//   - DELETED: Cluster has been deleted.
//   - INACTIVE: Cluster is not active, but not being deleted.
//   - UPDATING: Cluster is being updated.
//
// Only for Dedicated Cluster.
//   - IMPORTING: Cluster is being imported.
//
// Only for Dedicated Cluster.
//   - MODIFYING: Cluster is being modified.
//
// Only for Dedicated Cluster.
//   - PAUSING: Cluster is being paused.
//
// Only for Dedicated Cluster.
//   - PAUSED: Cluster is paused.
//
// Only for Dedicated Cluster.
//   - RESUMING: Cluster is resuming.
//
// Only for Dedicated Cluster.
//   - SCALING: Cluster is scaling.
//
// Only for Dedicated Cluster.
// Only for Mgmt Internal API.
//
// swagger:model commonv1beta1ClusterState
type Commonv1beta1ClusterState string

func NewCommonv1beta1ClusterState(value Commonv1beta1ClusterState) *Commonv1beta1ClusterState {
	return &value
}

// Pointer returns a pointer to a freshly-allocated Commonv1beta1ClusterState.
func (m Commonv1beta1ClusterState) Pointer() *Commonv1beta1ClusterState {
	return &m
}

const (

	// Commonv1beta1ClusterStateCLUSTERSTATEUNSPECIFIED captures enum value "CLUSTER_STATE_UNSPECIFIED"
	Commonv1beta1ClusterStateCLUSTERSTATEUNSPECIFIED Commonv1beta1ClusterState = "CLUSTER_STATE_UNSPECIFIED"

	// Commonv1beta1ClusterStateCREATING captures enum value "CREATING"
	Commonv1beta1ClusterStateCREATING Commonv1beta1ClusterState = "CREATING"

	// Commonv1beta1ClusterStateDELETING captures enum value "DELETING"
	Commonv1beta1ClusterStateDELETING Commonv1beta1ClusterState = "DELETING"

	// Commonv1beta1ClusterStateACTIVE captures enum value "ACTIVE"
	Commonv1beta1ClusterStateACTIVE Commonv1beta1ClusterState = "ACTIVE"

	// Commonv1beta1ClusterStateRESTORING captures enum value "RESTORING"
	Commonv1beta1ClusterStateRESTORING Commonv1beta1ClusterState = "RESTORING"

	// Commonv1beta1ClusterStateMAINTENANCE captures enum value "MAINTENANCE"
	Commonv1beta1ClusterStateMAINTENANCE Commonv1beta1ClusterState = "MAINTENANCE"

	// Commonv1beta1ClusterStateDELETED captures enum value "DELETED"
	Commonv1beta1ClusterStateDELETED Commonv1beta1ClusterState = "DELETED"

	// Commonv1beta1ClusterStateINACTIVE captures enum value "INACTIVE"
	Commonv1beta1ClusterStateINACTIVE Commonv1beta1ClusterState = "INACTIVE"

	// Commonv1beta1ClusterStateUPDATING captures enum value "UPDATING"
	Commonv1beta1ClusterStateUPDATING Commonv1beta1ClusterState = "UPDATING"

	// Commonv1beta1ClusterStateIMPORTING captures enum value "IMPORTING"
	Commonv1beta1ClusterStateIMPORTING Commonv1beta1ClusterState = "IMPORTING"

	// Commonv1beta1ClusterStateMODIFYING captures enum value "MODIFYING"
	Commonv1beta1ClusterStateMODIFYING Commonv1beta1ClusterState = "MODIFYING"

	// Commonv1beta1ClusterStatePAUSING captures enum value "PAUSING"
	Commonv1beta1ClusterStatePAUSING Commonv1beta1ClusterState = "PAUSING"

	// Commonv1beta1ClusterStatePAUSED captures enum value "PAUSED"
	Commonv1beta1ClusterStatePAUSED Commonv1beta1ClusterState = "PAUSED"

	// Commonv1beta1ClusterStateRESUMING captures enum value "RESUMING"
	Commonv1beta1ClusterStateRESUMING Commonv1beta1ClusterState = "RESUMING"

	// Commonv1beta1ClusterStateSCALING captures enum value "SCALING"
	Commonv1beta1ClusterStateSCALING Commonv1beta1ClusterState = "SCALING"
)

// for schema
var commonv1beta1ClusterStateEnum []interface{}

func init() {
	var res []Commonv1beta1ClusterState
	if err := json.Unmarshal([]byte(`["CLUSTER_STATE_UNSPECIFIED","CREATING","DELETING","ACTIVE","RESTORING","MAINTENANCE","DELETED","INACTIVE","UPDATING","IMPORTING","MODIFYING","PAUSING","PAUSED","RESUMING","SCALING"]`), &res); err != nil {
		panic(err)
	}
	for _, v := range res {
		commonv1beta1ClusterStateEnum = append(commonv1beta1ClusterStateEnum, v)
	}
}

func (m Commonv1beta1ClusterState) validateCommonv1beta1ClusterStateEnum(path, location string, value Commonv1beta1ClusterState) error {
	if err := validate.EnumCase(path, location, value, commonv1beta1ClusterStateEnum, true); err != nil {
		return err
	}
	return nil
}

// Validate validates this commonv1beta1 cluster state
func (m Commonv1beta1ClusterState) Validate(formats strfmt.Registry) error {
	var res []error

	// value enum
	if err := m.validateCommonv1beta1ClusterStateEnum("", "body", m); err != nil {
		return err
	}

	if len(res) > 0 {
		return errors.CompositeValidationError(res...)
	}
	return nil
}

// ContextValidate validates this commonv1beta1 cluster state based on context it is used
func (m Commonv1beta1ClusterState) ContextValidate(ctx context.Context, formats strfmt.Registry) error {
	return nil
}
