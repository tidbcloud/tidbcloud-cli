/*
TiDB Cloud Dedicated Open API

TiDB Cloud Dedicated Open API.

API version: v1beta1
*/

// Code generated by OpenAPI Generator (https://openapi-generator.tech); DO NOT EDIT.

package dedicated

import (
	"encoding/json"
)

// ClusterStorageNodeSettingStorageType Spec https://pingcap.feishu.cn/wiki/R6dcwK0Q3i9XSgkgG1Scoc82nbf.   - Standard: A cost-effective storage type suitable for most workloads, offering balanced performance and affordability.   - Single gp3 disk for TiKV storage Ideal for general-purpose clusters where affordability and reliability are priorities.  - Standard_Optimized: An enhanced version of Standard storage that delivers approximately 10% better overall cluster performance, balancing cost and performance.   - Separate gp3 disk for TiKV raft-log storage   - gp3 disk for TiKV data storage Recommended as the default choice for most clusters, offering an improved balance of cost and performance.  - Performance_Optimized: High-performance storage designed for workloads requiring low latency and greater stability. Provides superior performance for performance-sensitive clusters.   - Separate io2 disk for TiKV raft-log storage   - gp3 disk for TiKV data storage Suitable for clusters where consistent I/O performance and reduced latency are critical.  - Performance: Premium storage option offering the highest levels of performance and stability. Uses a single high-performance disk to minimize jitter and maximize IOPS.   - Single io2 disk for storage Designed for I/O jitter-sensitive clusters with extreme performance demands. Available via allowlist only for select customers.  - Standard_Premium: Reserved.  - Performance_Premium: Reserved.
type ClusterStorageNodeSettingStorageType string

// List of ClusterStorageNodeSettingStorageType
const (
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_STANDARD              ClusterStorageNodeSettingStorageType = "Standard"
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_STANDARD_OPTIMIZED    ClusterStorageNodeSettingStorageType = "Standard_Optimized"
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_PERFORMANCE_OPTIMIZED ClusterStorageNodeSettingStorageType = "Performance_Optimized"
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_PERFORMANCE           ClusterStorageNodeSettingStorageType = "Performance"
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_STANDARD_PREMIUM      ClusterStorageNodeSettingStorageType = "Standard_Premium"
	CLUSTERSTORAGENODESETTINGSTORAGETYPE_PERFORMANCE_PREMIUM   ClusterStorageNodeSettingStorageType = "Performance_Premium"

	// Unknown value for handling new enum values gracefully
	ClusterStorageNodeSettingStorageType_UNKNOWN ClusterStorageNodeSettingStorageType = "UNKNOWN"
)

// All allowed values of ClusterStorageNodeSettingStorageType enum
var AllowedClusterStorageNodeSettingStorageTypeEnumValues = []ClusterStorageNodeSettingStorageType{
	"Standard",
	"Standard_Optimized",
	"Performance_Optimized",
	"Performance",
	"Standard_Premium",
	"Performance_Premium",
}

func (v *ClusterStorageNodeSettingStorageType) UnmarshalJSON(src []byte) error {
	var value string
	err := json.Unmarshal(src, &value)
	if err != nil {
		return err
	}
	enumTypeValue := ClusterStorageNodeSettingStorageType(value)
	for _, existing := range AllowedClusterStorageNodeSettingStorageTypeEnumValues {
		if existing == enumTypeValue {
			*v = enumTypeValue
			return nil
		}
	}

	// Instead of returning an error, assign UNKNOWN value
	*v = ClusterStorageNodeSettingStorageType_UNKNOWN
	return nil
}

// NewClusterStorageNodeSettingStorageTypeFromValue returns a pointer to a valid ClusterStorageNodeSettingStorageType
// for the value passed as argument, or UNKNOWN if the value is not in the enum list
func NewClusterStorageNodeSettingStorageTypeFromValue(v string) *ClusterStorageNodeSettingStorageType {
	ev := ClusterStorageNodeSettingStorageType(v)
	if ev.IsValid() {
		return &ev
	}
	unknown := ClusterStorageNodeSettingStorageType_UNKNOWN
	return &unknown
}

// IsValid return true if the value is valid for the enum, false otherwise
func (v ClusterStorageNodeSettingStorageType) IsValid() bool {
	for _, existing := range AllowedClusterStorageNodeSettingStorageTypeEnumValues {
		if existing == v {
			return true
		}
	}
	return false
}

// Ptr returns reference to ClusterStorageNodeSettingStorageType value
func (v ClusterStorageNodeSettingStorageType) Ptr() *ClusterStorageNodeSettingStorageType {
	return &v
}

type NullableClusterStorageNodeSettingStorageType struct {
	value *ClusterStorageNodeSettingStorageType
	isSet bool
}

func (v NullableClusterStorageNodeSettingStorageType) Get() *ClusterStorageNodeSettingStorageType {
	return v.value
}

func (v *NullableClusterStorageNodeSettingStorageType) Set(val *ClusterStorageNodeSettingStorageType) {
	v.value = val
	v.isSet = true
}

func (v NullableClusterStorageNodeSettingStorageType) IsSet() bool {
	return v.isSet
}

func (v *NullableClusterStorageNodeSettingStorageType) Unset() {
	v.value = nil
	v.isSet = false
}

func NewNullableClusterStorageNodeSettingStorageType(val *ClusterStorageNodeSettingStorageType) *NullableClusterStorageNodeSettingStorageType {
	return &NullableClusterStorageNodeSettingStorageType{value: val, isSet: true}
}

func (v NullableClusterStorageNodeSettingStorageType) MarshalJSON() ([]byte, error) {
	return json.Marshal(v.value)
}

func (v *NullableClusterStorageNodeSettingStorageType) UnmarshalJSON(src []byte) error {
	v.isSet = true
	return json.Unmarshal(src, &v.value)
}
