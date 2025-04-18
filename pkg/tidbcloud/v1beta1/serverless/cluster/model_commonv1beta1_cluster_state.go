/*
TiDB Cloud Serverless Open API

TiDB Cloud Serverless Open API

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package cluster

import (
	"encoding/json"
)

// Commonv1beta1ClusterState Enum of possible states of a cluster.   - CREATING: Cluster is being created.  - DELETING: Cluster is being deleted.  - ACTIVE: Cluster is active for use.  - RESTORING: Cluster data is being restored.  - MAINTENANCE: Cluster is under maintenance.  - DELETED: Cluster has been deleted.  - INACTIVE: Cluster is not active, but not being deleted.  - UPGRADING: Cluster is being updated. Only for Dedicated Cluster.  - IMPORTING: Cluster is being imported. Only for Dedicated Cluster.  - MODIFYING: Cluster is being modified. Only for Dedicated Cluster.  - PAUSING: Cluster is being paused. Only for Dedicated Cluster.  - PAUSED: Cluster is paused. Only for Dedicated Cluster.  - RESUMING: Cluster is resuming. Only for Dedicated Cluster.
type Commonv1beta1ClusterState string

// List of commonv1beta1ClusterState
const (
	COMMONV1BETA1CLUSTERSTATE_CREATING    Commonv1beta1ClusterState = "CREATING"
	COMMONV1BETA1CLUSTERSTATE_DELETING    Commonv1beta1ClusterState = "DELETING"
	COMMONV1BETA1CLUSTERSTATE_ACTIVE      Commonv1beta1ClusterState = "ACTIVE"
	COMMONV1BETA1CLUSTERSTATE_RESTORING   Commonv1beta1ClusterState = "RESTORING"
	COMMONV1BETA1CLUSTERSTATE_MAINTENANCE Commonv1beta1ClusterState = "MAINTENANCE"
	COMMONV1BETA1CLUSTERSTATE_DELETED     Commonv1beta1ClusterState = "DELETED"
	COMMONV1BETA1CLUSTERSTATE_INACTIVE    Commonv1beta1ClusterState = "INACTIVE"
	COMMONV1BETA1CLUSTERSTATE_UPGRADING   Commonv1beta1ClusterState = "UPGRADING"
	COMMONV1BETA1CLUSTERSTATE_IMPORTING   Commonv1beta1ClusterState = "IMPORTING"
	COMMONV1BETA1CLUSTERSTATE_MODIFYING   Commonv1beta1ClusterState = "MODIFYING"
	COMMONV1BETA1CLUSTERSTATE_PAUSING     Commonv1beta1ClusterState = "PAUSING"
	COMMONV1BETA1CLUSTERSTATE_PAUSED      Commonv1beta1ClusterState = "PAUSED"
	COMMONV1BETA1CLUSTERSTATE_RESUMING    Commonv1beta1ClusterState = "RESUMING"
)

// All allowed values of Commonv1beta1ClusterState enum
var AllowedCommonv1beta1ClusterStateEnumValues = []Commonv1beta1ClusterState{
	"CREATING",
	"DELETING",
	"ACTIVE",
	"RESTORING",
	"MAINTENANCE",
	"DELETED",
	"INACTIVE",
	"UPGRADING",
	"IMPORTING",
	"MODIFYING",
	"PAUSING",
	"PAUSED",
	"RESUMING",
}

func (v *Commonv1beta1ClusterState) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := Commonv1beta1ClusterState(value)
	for _, existing := range AllowedCommonv1beta1ClusterStateEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	*v = Commonv1beta1ClusterState(value)
	return nil
}

// NewCommonv1beta1ClusterStateFromValue returns a pointer to a valid Commonv1beta1ClusterState for the value passed as argument
func NewCommonv1beta1ClusterStateFromValue(v string) *Commonv1beta1ClusterState {
	ev := Commonv1beta1ClusterState(v)
	return &ev
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v Commonv1beta1ClusterState) IsValid() bool {
	for _, existing := range AllowedCommonv1beta1ClusterStateEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to commonv1beta1ClusterState value
func (v Commonv1beta1ClusterState) Ptr() *Commonv1beta1ClusterState {
	return &v
}

type NullableCommonv1beta1ClusterState struct {
	value *Commonv1beta1ClusterState
	isSet bool
}

func (v NullableCommonv1beta1ClusterState) Get() *Commonv1beta1ClusterState {
	return v.value
}

func (v *NullableCommonv1beta1ClusterState) Set(val *Commonv1beta1ClusterState) {
	v.value = val
	v.isSet = true
}

func (v NullableCommonv1beta1ClusterState) IsSet() bool {
	return v.isSet
}

func (v *NullableCommonv1beta1ClusterState) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableCommonv1beta1ClusterState(val *Commonv1beta1ClusterState) *NullableCommonv1beta1ClusterState {
	return &NullableCommonv1beta1ClusterState{value: val, isSet: true}
}

func (v NullableCommonv1beta1ClusterState) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableCommonv1beta1ClusterState) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
