# TiDB Cloud SDK knowledge

TiDB Cloud CLI keeps the TiDB Cloud Open API SDK inside the project.

Key points about the SDK:

1. Swagger files are the single source of truth for the SDK; they live under `/pkg/tidbcloud`.
2. The openapi-generator library generates the SDK from Swagger. You do not need to know openapi-generator in detail; use the commands in the `Makefile` to invoke it.
3. `internal/service/cloud/api_client.go` combines the generated SDKs into a single client.
4. Mocks for `api_client.go` must also be generated.

It is the user's responsibility to add or change Swagger files under `pkg/tidbcloud/`.

It is the AI agent's responsibility to generate the SDK. Follow the workflow below:

1. If new swagger is added, add the generate client script in Makefile, refer the format under generate-v1beta1-serverless-client
2. Run the appropriate command based on where Swagger files were added or changed. For example, if `pkg/tidbcloud/v1beta1/serverless` changed, run `generate-v1beta1-serverless-client`. If you cannot determine which Swagger changed, run `make generate-v1beta1-client` to generate all.
3. Update `internal/service/cloud/api_client.go` according to the generated SDK. If new methods appear in the SDK, add new interface methods and implementations in `api_client.go` following the existing code style. If new API clients appear, wrap them in the `ClientDelegate` struct.