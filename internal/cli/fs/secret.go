// Copyright 2026 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package fs

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"os/exec"
	"sort"
	"strings"
	"text/tabwriter"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/fs"

	"github.com/spf13/cobra"
)

const (
	defaultAuditLimit   = 100
	maxClientAuditLimit = 1000
)

func secretCmd(h *internal.Helper) *cobra.Command {
	var secretCmd = &cobra.Command{
		Use:   "secret",
		Short: "Secret management operations",
		Long:  `Manage secrets in TiDB Cloud FS vault.`,
	}

	secretCmd.AddCommand(secretSetCmd(h))
	secretCmd.AddCommand(secretGetCmd(h))
	secretCmd.AddCommand(secretExecCmd(h))
	secretCmd.AddCommand(secretLsCmd(h))
	secretCmd.AddCommand(secretRmCmd(h))
	secretCmd.AddCommand(secretGrantCmd(h))
	secretCmd.AddCommand(secretRevokeCmd(h))
	secretCmd.AddCommand(secretAuditCmd(h))

	return secretCmd
}

func secretSetCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "set <name> <field=value|field=@file|field=-> ...",
		Short: "Create or update a secret",
		Example: `  ticloud fs secret set myapp DATABASE_URL=mysql://...
  ticloud fs secret set myapp key=@secret.txt password=-`,
		Args: cobra.MinimumNArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			if err := validateSecretName(name); err != nil {
				return err
			}
			fields, err := parseSecretFields(args[1:])
			if err != nil {
				return err
			}

			ctx := context.Background()
			if _, err := client.CreateVaultSecret(ctx, name, fields); err != nil {
				// If conflict, update instead
				if !strings.Contains(err.Error(), "conflict") {
					return suggestInitIfTenantNotFound(fmt.Errorf("set secret: %w", err))
				}
				_, err = client.UpdateVaultSecret(ctx, name, fields)
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("update secret: %w", err))
				}
			}
			return nil
		},
	}
}

func secretGetCmd(h *internal.Helper) *cobra.Command {
	var asJSON, asEnv bool
	var cmd = &cobra.Command{
		Use:   "get <name[/field]>",
		Short: "Read a secret or one field",
		Example: `  ticloud fs secret get myapp
  ticloud fs secret get myapp/password --json
  ticloud fs secret get myapp/api_key --env`,
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			ref := args[0]
			name, field, err := parseSecretRef(ref)
			if err != nil {
				return err
			}
			if asJSON && asEnv {
				return fmt.Errorf("--json and --env are mutually exclusive")
			}

			ctx := context.Background()
			if field != "" {
				value, err := client.ReadVaultSecretField(ctx, name, field)
				if err != nil {
					return suggestInitIfTenantNotFound(fmt.Errorf("get secret field: %w", err))
				}
				switch {
				case asEnv:
					envKey, err := normalizeSecretEnvKey(field)
					if err != nil {
						return err
					}
					fmt.Fprintf(h.IOStreams.Out, "%s=%s\n", envKey, value)
				case asJSON:
					return writeJSON(h, map[string]string{field: value})
				default:
					fmt.Fprintln(h.IOStreams.Out, value)
				}
				return nil
			}

			fields, err := client.ReadVaultSecret(ctx, name)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("get secret: %w", err))
			}
			if asEnv {
				return printEnv(h, fields)
			}
			return writeJSON(h, fields)
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&asEnv, "env", false, "Output as environment variables")
	return cmd
}

func secretExecCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:     "exec <name> -- <command...>",
		Short:   "Run a command with secret fields injected as env vars",
		Example: `  ticloud fs secret exec myapp -- ./run.sh`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			if err := validateSecretName(name); err != nil {
				return err
			}

			sep := -1
			for i := 1; i < len(args); i++ {
				if args[i] == "--" {
					sep = i
					break
				}
			}
			if sep < 0 || sep == len(args)-1 {
				return fmt.Errorf("usage: ticloud fs secret exec <name> -- <command...>")
			}
			cmdArgs := args[sep+1:]

			fields, err := client.ReadVaultSecret(context.Background(), name)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("get secret: %w", err))
			}

			c := exec.Command(cmdArgs[0], cmdArgs[1:]...)
			c.Stdin = os.Stdin
			c.Stdout = h.IOStreams.Out
			c.Stderr = h.IOStreams.Err
			envMap, err := buildSecretEnvMap(fields)
			if err != nil {
				return err
			}
			c.Env = mergeEnv(os.Environ(), envMap)
			return c.Run()
		},
	}
}

func secretLsCmd(h *internal.Helper) *cobra.Command {
	var asJSON bool
	var cmd = &cobra.Command{
		Use:     "ls",
		Short:   "List secrets",
		Example: `  ticloud fs secret ls --json`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			secrets, err := client.ListVaultSecrets(context.Background())
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("list secrets: %w", err))
			}

			names := make([]string, 0, len(secrets))
			for _, sec := range secrets {
				names = append(names, sec.Name)
			}
			sort.Strings(names)

			if asJSON {
				return writeJSON(h, map[string]any{"secrets": names})
			}
			for _, name := range names {
				fmt.Fprintln(h.IOStreams.Out, name)
			}
			return nil
		},
	}
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	return cmd
}

func secretRmCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:     "rm <name>",
		Short:   "Delete a secret",
		Example: `  ticloud fs secret rm myapp`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			name := args[0]
			if err := validateSecretName(name); err != nil {
				return err
			}
			if err := client.DeleteVaultSecret(context.Background(), name); err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("delete secret: %w", err))
			}
			return nil
		},
	}
}

func secretGrantCmd(h *internal.Helper) *cobra.Command {
	var (
		agentID   string
		taskID    string
		ttlRaw    string
		asJSON    bool
		tokenOnly bool
	)
	var cmd = &cobra.Command{
		Use:     "grant --agent <id> --ttl <duration> [--task <id>] <scope...>",
		Short:   "Issue a scoped capability token",
		Example: `  ticloud fs secret grant --agent myagent --ttl 1h myapp/password`,
		Args:    cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			scope := args
			for _, entry := range scope {
				if _, _, err := parseSecretRef(entry); err != nil {
					return fmt.Errorf("invalid scope %q: %w", entry, err)
				}
			}
			if asJSON && tokenOnly {
				return fmt.Errorf("--json and --token-only are mutually exclusive")
			}
			if agentID == "" {
				return fmt.Errorf("--agent is required")
			}
			if ttlRaw == "" {
				return fmt.Errorf("--ttl is required")
			}
			ttl, err := time.ParseDuration(ttlRaw)
			if err != nil {
				return fmt.Errorf("invalid --ttl %q: %w", ttlRaw, err)
			}
			if ttl <= 0 {
				return fmt.Errorf("--ttl must be positive")
			}

			resp, err := client.IssueVaultToken(context.Background(), agentID, taskID, scope, ttl)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("grant token: %w", err))
			}
			switch {
			case tokenOnly:
				fmt.Fprintln(h.IOStreams.Out, resp.Token)
			case asJSON:
				return writeJSON(h, resp)
			default:
				fmt.Fprintf(h.IOStreams.Out, "token=%s\n", resp.Token)
				fmt.Fprintf(h.IOStreams.Out, "token_id=%s\n", resp.TokenID)
				fmt.Fprintf(h.IOStreams.Out, "expires_at=%s\n", resp.ExpiresAt.Format(time.RFC3339))
			}
			return nil
		},
	}
	cmd.Flags().StringVar(&agentID, "agent", "", "Agent ID")
	cmd.Flags().StringVar(&taskID, "task", "", "Task ID")
	cmd.Flags().StringVar(&ttlRaw, "ttl", "", "Token TTL (e.g., 1h, 30m)")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	cmd.Flags().BoolVar(&tokenOnly, "token-only", false, "Output only the token string")
	return cmd
}

func secretRevokeCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:     "revoke <token-id>",
		Short:   "Revoke a capability token",
		Example: `  ticloud fs secret revoke tok_abc123`,
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}
			if err := client.RevokeVaultToken(context.Background(), args[0]); err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("revoke token: %w", err))
			}
			return nil
		},
	}
}

func secretAuditCmd(h *internal.Helper) *cobra.Command {
	var (
		secretName string
		agentID    string
		sinceRaw   string
		limit      = defaultAuditLimit
		asJSON     bool
	)
	var cmd = &cobra.Command{
		Use:   "audit [--secret <name>] [--agent <id>] [--since <duration>] [--limit <n>]",
		Short: "Query vault audit events",
		Example: `  ticloud fs secret audit --limit 50
  ticloud fs secret audit --secret myapp --since 24h`,
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := newClient(cmd)
			if err != nil {
				return err
			}

			queryLimit := limit
			if agentID != "" || sinceRaw != "" {
				queryLimit = maxClientAuditLimit
			}
			if queryLimit > maxClientAuditLimit {
				queryLimit = maxClientAuditLimit
			}

			events, err := client.QueryVaultAudit(context.Background(), secretName, queryLimit)
			if err != nil {
				return suggestInitIfTenantNotFound(fmt.Errorf("audit query: %w", err))
			}

			if sinceRaw != "" {
				d, err := time.ParseDuration(sinceRaw)
				if err != nil {
					return fmt.Errorf("invalid --since %q: %w", sinceRaw, err)
				}
				if d <= 0 {
					return fmt.Errorf("--since must be positive")
				}
				events = filterAuditEvents(events, agentID, time.Now().Add(-d))
			} else if agentID != "" {
				events = filterAuditEvents(events, agentID, time.Time{})
			}
			if len(events) > limit {
				events = events[:limit]
			}

			if asJSON {
				return writeJSON(h, map[string]any{"events": events})
			}
			printAudit(h, events)
			return nil
		},
	}
	cmd.Flags().StringVar(&secretName, "secret", "", "Filter by secret name")
	cmd.Flags().StringVar(&agentID, "agent", "", "Filter by agent ID")
	cmd.Flags().StringVar(&sinceRaw, "since", "", "Filter events newer than duration (e.g., 24h)")
	cmd.Flags().IntVar(&limit, "limit", defaultAuditLimit, "Maximum number of events")
	cmd.Flags().BoolVar(&asJSON, "json", false, "Output as JSON")
	return cmd
}

func validateSecretName(name string) error {
	if name == "" {
		return fmt.Errorf("secret name is required")
	}
	if strings.Contains(name, "/") {
		return fmt.Errorf("secret name %q must be flat; use <name/field> only for reads and scopes", name)
	}
	if strings.Contains(name, "*") {
		return fmt.Errorf("wildcard scope entries are not supported: %q", name)
	}
	return nil
}

func parseSecretRef(raw string) (string, string, error) {
	if raw == "" {
		return "", "", fmt.Errorf("secret reference is required")
	}
	parts := strings.SplitN(raw, "/", 2)
	name := parts[0]
	if err := validateSecretName(name); err != nil {
		return "", "", err
	}
	if len(parts) == 1 {
		return name, "", nil
	}
	field := parts[1]
	if field == "" {
		return "", "", fmt.Errorf("field name is required in %q", raw)
	}
	if strings.Contains(field, "*") {
		return "", "", fmt.Errorf("wildcard scope entries are not supported: %q", raw)
	}
	return name, field, nil
}

func parseSecretFields(args []string) (map[string]string, error) {
	fields := make(map[string]string, len(args))
	var stdinValue []byte
	var stdinRead bool
	for _, arg := range args {
		key, valueSpec, ok := strings.Cut(arg, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("field assignment must be field=value, field=@file, or field=-: %q", arg)
		}
		var value string
		switch {
		case valueSpec == "-":
			if !stdinRead {
				data, err := io.ReadAll(os.Stdin)
				if err != nil {
					return nil, fmt.Errorf("read stdin: %w", err)
				}
				stdinValue = data
				stdinRead = true
			}
			value = string(stdinValue)
		case strings.HasPrefix(valueSpec, "@"):
			data, err := os.ReadFile(valueSpec[1:])
			if err != nil {
				return nil, fmt.Errorf("read %s: %w", valueSpec[1:], err)
			}
			value = string(data)
		default:
			value = valueSpec
		}
		fields[key] = value
	}
	return fields, nil
}

func writeJSON(h *internal.Helper, v any) error {
	enc := json.NewEncoder(h.IOStreams.Out)
	enc.SetIndent("", "  ")
	return enc.Encode(v)
}

func printEnv(h *internal.Helper, fields map[string]string) error {
	envMap, err := buildSecretEnvMap(fields)
	if err != nil {
		return err
	}
	keys := make([]string, 0, len(envMap))
	for k := range envMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		fmt.Fprintf(h.IOStreams.Out, "%s=%s\n", key, envMap[key])
	}
	return nil
}

func printAudit(h *internal.Helper, events []fs.VaultAuditEvent) {
	w := tabwriter.NewWriter(h.IOStreams.Out, 0, 4, 2, ' ', 0)
	_, _ = fmt.Fprintln(w, "TIME\tAGENT\tACTION\tSECRET\tFIELD")
	for _, ev := range events {
		_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t%s\t%s\n",
			ev.Timestamp.Format(time.RFC3339),
			ev.AgentID,
			ev.EventType,
			ev.SecretName,
			ev.FieldName,
		)
	}
	_ = w.Flush()
}

func buildSecretEnvMap(fields map[string]string) (map[string]string, error) {
	env := make(map[string]string, len(fields))
	owners := make(map[string]string, len(fields))
	for field, value := range fields {
		envKey, err := normalizeSecretEnvKey(field)
		if err != nil {
			return nil, err
		}
		if prevField, exists := owners[envKey]; exists {
			return nil, fmt.Errorf("secret fields %q and %q both normalize to env var %q", prevField, field, envKey)
		}
		owners[envKey] = field
		env[envKey] = value
	}
	return env, nil
}

func filterAuditEvents(events []fs.VaultAuditEvent, agentID string, since time.Time) []fs.VaultAuditEvent {
	filtered := make([]fs.VaultAuditEvent, 0, len(events))
	for _, ev := range events {
		if agentID != "" && ev.AgentID != agentID {
			continue
		}
		if !since.IsZero() && ev.Timestamp.Before(since) {
			continue
		}
		filtered = append(filtered, ev)
	}
	return filtered
}

func mergeEnv(base []string, overrides map[string]string) []string {
	merged := make(map[string]string, len(base)+len(overrides))
	for _, entry := range base {
		key, value, ok := strings.Cut(entry, "=")
		if ok {
			merged[key] = value
		}
	}
	for k, v := range overrides {
		merged[k] = v
	}
	keys := make([]string, 0, len(merged))
	for k := range merged {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	env := make([]string, 0, len(keys))
	for _, k := range keys {
		env = append(env, k+"="+merged[k])
	}
	return env
}

func normalizeSecretEnvKey(field string) (string, error) {
	if field == "" {
		return "", fmt.Errorf("secret field name is required")
	}
	var b strings.Builder
	b.Grow(len(field) + 1)
	for i := 0; i < len(field); i++ {
		ch := field[i]
		switch {
		case ch >= 'a' && ch <= 'z':
			b.WriteByte(ch - ('a' - 'A'))
		case ch >= 'A' && ch <= 'Z', ch >= '0' && ch <= '9', ch == '_':
			b.WriteByte(ch)
		default:
			b.WriteByte('_')
		}
	}
	key := b.String()
	if key == "" {
		return "", fmt.Errorf("secret field name is required")
	}
	if key[0] >= '0' && key[0] <= '9' {
		key = "_" + key
	}
	return key, nil
}
