# Telemetry

`ticloud`'s telemetry collects anonymous usage data to help improve the product. `ticloud` enables telemetry by default.

## What telemetry data is collected

`ticloud` telemetry tracks non-Personally-Identifiable Information (PII), which includes:

| Filed name                  | Description                                                 | Example value                                                                                                                                                                                                                |
|-----------------------------|-------------------------------------------------------------|------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| `timestamp`                 | The time when the command is executed                       | `2023-02-03T09:33:49.384304+08:00`                                                                                                                                                                                           |
| `source`                    | The name of the command                                     | `ticloud`                                                                                                                                                                                                                    |
| `properties`                | The map which contains properties                           | -                                                                                                                                                                                                                            |
| `properties["alias"]`       | The alias of the command                                    | `ls`                                                                                                                                                                                                                         |
| `properties["arch"]`        | The architecture of the machine                             | `arm64`                                                                                                                                                                                                                      |
| `properties["auth_method"]` | The authentication method used to access the TiDB Cloud API | `api_key`                                                                                                                                                                                                                    |
| `properties["command"]`     | The full command path of the command                        | `ticloud-project-list`                                                                                                                                                                                                       |
| `properties["duration"]`    | The execution duration of the command                       | `871`                                                                                                                                                                                                                        |
| `properties["error"]`       | The error message of the command                            | `[POST /api/internal/projects/{project_id}/clusters/{cluster_id}/upload_url][403] GenerateUploadURL default  \\u0026{Code:49900003 Details:[] Message:you do not have permission to access this project, project_id: 43511}` |
| `properties["flags"]`       | The names of flags used when executing the command          | `[debug, output]`                                                                                                                                                                                                            |
| `properties["git_commit"]`  | The git commit hash of the CLI                              | `dfdsfdsdfs`                                                                                                                                                                                                                 |
| `properties["installer"]`   | The installer of the CLI                                    | `TiUP`                                                                                                                                                                                                                       |
| `properties["interactive"]` | Whether the command is executed in interactive mode         | `false`                                                                                                                                                                                                                      |
| `properties["os"]`          | The operating system of the machine                         | `darwin`                                                                                                                                                                                                                     |
| `properties["project_id"]`  | The TiDB Cloud project ID which command used                | `43511`                                                                                                                                                                                                                      |
| `properties["result"]`      | The result of the command, `SUCCESS` or `ERROR`             | `SUCCESS`                                                                                                                                                                                                                    |
| `properties["terminal"]`    | The terminal type                                           | `tty`                                                                                                                                                                                                                        |
| `properties["version"]`     | The version of the CLI                                      | `v0.0.1`                                                                                                                                                                                                                     |
              "
To view the full content of telemetry, use flag `--debug` when running `ticloud` command.

```shell
[2023/02/03 10:02:36.932 +08:00] [DEBUG] [sender.go:46] ["sending telemetry events"] [body="[{\"timestamp\":\"2023-02-03T09:33:49.384304+08:00\",\"source\":\"ticloud\",\"properties\":{\"alias\":\"ls\",\"arch\":\"arm64\",\"auth_method\":\"api_key\",\"command\":\"ticloud-project-list\",\"duration\":871,\"flags\":[\"debug\"],\"git_commit\":\"dfdsfdsdfs\",\"interactive\":false,\"os\":\"darwin\",\"result\":\"SUCCESS\",\"terminal\":\"tty\",\"version\":\"v0.0.1\"}},{\"timestamp\":\"2023-02-03T10:02:36.93222+08:00\",\"source\":\"ticloud\",\"properties\":{\"alias\":\"ls\",\"arch\":\"arm64\",\"auth_method\":\"api_key\",\"command\":\"ticloud-project-list\",\"duration\":1077,\"flags\":[\"debug\",\"output\"],\"git_commit\":\"dfdsfdsdfs\",\"interactive\":false,\"os\":\"darwin\",\"result\":\"SUCCESS\",\"terminal\":\"tty\",\"version\":\"v0.0.1\"}}]"]
```

## Disable Telemetry for the CLI
To disable telemetry for the CLI, use the `ticloud config set` to set profile configuration:

```shell
ticloud config set telemetry-enabled false
```

> **Note:**
>
> After using this command, the configuration `telemetry-enabled` is only configured for the current profile. If you want to disable telemetry for all profiles, you need to set the configuration for each profile.
