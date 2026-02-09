# Flag inference

1. Infer flags if the user does not provide the command with flags.

Skip this if the user has already provided flags.

Infer flags from the Swagger spec and SDK client parameters. For example, see the branch create command in `internal/cli/serverless/branch/create.go`.

The SDK client parameter `Branch` is as follows:

```
type Branch struct {
	// The unique identifier for the branch.
	Name *string `json:"name,omitempty"`
	// The system-generated ID of the branch.
	BranchId *string `json:"branchId,omitempty"`
	// The user-defined name of the branch.
	DisplayName string `json:"displayName"`
	// The ID of the cluster to which the branch belongs.
	ClusterId *string `json:"clusterId,omitempty"`
	// The ID of the branch parent.
	ParentId *string `json:"parentId,omitempty"`
	// The email address of the user who create the branch.
	CreatedBy *string `json:"createdBy,omitempty"`
	// The state of the branch.
	State *BranchState `json:"state,omitempty"`
	// The connection endpoints for accessing the branch.
	Endpoints *BranchEndpoints `json:"endpoints,omitempty"`
	// The unique prefix automatically generated for SQL usernames on this cluster. TiDB Cloud uses this prefix to distinguish between clusters. For more information, see [User name prefix](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#user-name-prefix).
	UserPrefix NullableString `json:"userPrefix,omitempty"`
	// The timestamp when the branch was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	CreateTime *time.Time `json:"createTime,omitempty"`
	// The timestamp when the branch was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.
	UpdateTime *time.Time `json:"updateTime,omitempty"`
	// The annotations for the branch.
	Annotations *map[string]string `json:"annotations,omitempty"`
	// The display name of the parent branch from which the branch was created.
	ParentDisplayName *string `json:"parentDisplayName,omitempty"`
	// The point in time on the parent branch from which the branch is created. The timestamp is truncated to seconds without rounding.
	ParentTimestamp NullableTime `json:"parentTimestamp,omitempty"`
	// The root password of the branch. It must be between 8 and 64 characters long and can contain letters, numbers, and special characters.
	RootPassword         *string `json:"rootPassword,omitempty" validate:"regexp=^.{8,64}$"`
	AdditionalProperties map[string]interface{}
}
```

The part of the swagger is as follows:

```
        "parameters": [
          {
            "name": "clusterId",
            "description": "The ID of the cluster to which the branch belongs.",
            "in": "path",
            "required": true,
            "type": "string"
          },
          {
            "name": "branch",
            "description": "The branch being created.",
            "in": "body",
            "required": true,
            "schema": {
              "$ref": "#/definitions/Branch"
            }
          }
        ],

"Branch": {
      "type": "object",
      "properties": {
        "name": {
          "type": "string",
          "description": "The unique identifier for the branch.",
          "readOnly": true
        },
        "branchId": {
          "type": "string",
          "description": "The system-generated ID of the branch.",
          "readOnly": true
        },
        "displayName": {
          "type": "string",
          "description": "The user-defined name of the branch."
        },
        "clusterId": {
          "type": "string",
          "description": "The ID of the cluster to which the branch belongs.",
          "readOnly": true
        },
        "parentId": {
          "type": "string",
          "description": "The ID of the branch parent."
        },
        "createdBy": {
          "type": "string",
          "description": "The email address of the user who create the branch.",
          "readOnly": true
        },
        "state": {
          "description": "The state of the branch.",
          "readOnly": true,
          "allOf": [
            {
              "$ref": "#/definitions/Branch.State"
            }
          ]
        },
        "endpoints": {
          "description": "The connection endpoints for accessing the branch.",
          "readOnly": true,
          "allOf": [
            {
              "$ref": "#/definitions/Branch.Endpoints"
            }
          ]
        },
        "userPrefix": {
          "type": "string",
          "x-nullable": true,
          "description": "The unique prefix automatically generated for SQL usernames on this cluster. TiDB Cloud uses this prefix to distinguish between clusters. For more information, see [User name prefix](https://docs.pingcap.com/tidbcloud/select-cluster-tier/#user-name-prefix).",
          "readOnly": true
        },
        "createTime": {
          "type": "string",
          "format": "date-time",
          "description": "The timestamp when the branch was created, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.",
          "readOnly": true
        },
        "updateTime": {
          "type": "string",
          "format": "date-time",
          "description": "The timestamp when the branch was last updated, in the [ISO 8601](https://en.wikipedia.org/wiki/ISO_8601) format.",
          "readOnly": true
        },
        "annotations": {
          "type": "object",
          "additionalProperties": {
            "type": "string"
          },
          "description": "The annotations for the branch."
        },
        "parentDisplayName": {
          "type": "string",
          "description": "The display name of the parent branch from which the branch was created.",
          "readOnly": true
        },
        "parentTimestamp": {
          "type": "string",
          "format": "date-time",
          "x-nullable": true,
          "description": "The point in time on the parent branch from which the branch is created. The timestamp is truncated to seconds without rounding."
        },
        "rootPassword": {
          "type": "string",
          "example": "my-shining-password",
          "description": "The root password of the branch. It must be between 8 and 64 characters long and can contain letters, numbers, and special characters.",
          "maxLength": 64,
          "minLength": 8,
          "pattern": "^.{8,64}$"
        }
      },
      "description": "Message for branch.",
      "required": [
        "displayName"
      ]
    }
```

Generally, flags should exclude readOnly fields in the Swagger spec and include other fields. Exceptions may occur, so this rule is not mandatory.

```
	createCmd.Flags().StringP(flag.DisplayName, flag.DisplayNameShort, "", "The displayName of the branch to be created.")
	createCmd.Flags().StringP(flag.ClusterID, flag.ClusterIDShort, "", "The ID of the cluster, in which the branch will be created.")
	createCmd.Flags().StringP(flag.ParentID, "", "", "The ID of the branch parent, default is cluster id.")
	createCmd.Flags().StringP(flag.ParentTimestamp, "", "", "The timestamp of the parent branch, default is current time. (RFC3339 format, e.g., 2024-01-01T00:00:00Z)")
```

2. Infer flag types from the SDK client parameter types.

For example: if the parameter is a string type, use `Flags().String`.

Flag inference may be inaccurate, so:
- Always ask the user when unsure about flags or flag types.
- Always ask the user to confirm the final result.
