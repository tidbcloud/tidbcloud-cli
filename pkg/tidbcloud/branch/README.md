We need to change the googlerpcStatus in branch_openapi.swagger.json

```
    "googlerpcStatus": {
      "type": "object",
      "properties": {
        "code": {
          "type": "integer",
          "format": "int32"
        },
        "message": {
          "type": "string"
        },
        "details": {
          "type": "array",
          "items": {
            "type": "object",
            "$ref": "#/definitions/protobufAny"
          }
        }
      }
    }
```

to

```
    "googlerpcStatus": {
      "type": "object",
      "properties": {
        "error": {
          "type": "object",
          "properties": {
            "code": {
              "type": "integer",
              "format": "int32"
            },
            "message": {
              "type": "string"
            },
            "details": {
              "type": "array",
              "items": {
                "type": "object",
                "$ref": "#/definitions/protobufAny"
              }
            }
          }
        }
      }
    }
```

Reason: open API's error output is different from `googlerpcStatus` (default error output)
