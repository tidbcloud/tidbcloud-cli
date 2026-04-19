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
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/config"
	"github.com/tidbcloud/tidbcloud-cli/internal/config/store"
	"github.com/tidbcloud/tidbcloud-cli/internal/flag"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/cloud"
	"github.com/tidbcloud/tidbcloud-cli/internal/service/fs"

	"github.com/juju/errors"
	"github.com/spf13/cobra"
	"github.com/zalando/go-keyring"
)

// FSCmd returns the fs command.
func FSCmd(h *internal.Helper) *cobra.Command {
	var fsCmd = &cobra.Command{
		Use:   "fs",
		Short: "Filesystem operations for TiDB Cloud FS",
		Long: `Filesystem operations for TiDB Cloud FS.

This command provides filesystem-like operations for managing files
and directories in TiDB Cloud FS.

Path format: :/path/to/file
  - :/ denotes the root of the remote filesystem

Authentication:
  The fs command reuses your TiDB Cloud CLI authentication (OAuth or API key).
  Please login with 'ticloud auth login' or configure API keys with:
    ticloud config set public-key <public-key>
    ticloud config set private-key <private-key>

Database Association:
  FS operations can be associated with a specific database:
  - Serverless cluster: ticloud config set fs.cluster-id <cluster-id>
    (sent as X-TIDBCLOUD-CLUSTER-ID header)
  - TiDB Zero instance: ticloud config set fs.zero-instance-id <instance-id>
    (sent as X-TIDBCLOUD-ZERO-INSTANCE-ID header)
  If both are set, cluster-id takes precedence.

Configuration:
  - Server URL: ticloud config set fs-endpoint <url> (default: https://fs.tidbapi.com/)
  - Env override: TICLOUD_FS_ENDPOINT, TICLOUD_FS_CLUSTER_ID, TICLOUD_FS_ZERO_INSTANCE_ID

Examples:
  ticloud fs ls :/                    # List root directory
  ticloud fs ls -l :/myfolder        # List with details
  ticloud fs cat :/file.txt          # Display file contents
  ticloud fs cp local.txt :/remote/  # Upload file
  ticloud fs cp :/remote/file.txt .  # Download file
  ticloud fs mv :/old :/new          # Rename/move
  ticloud fs rm :/file.txt           # Remove file
  ticloud fs mkdir :/mydir           # Create directory
  ticloud fs stat :/file.txt         # Show metadata
  ticloud fs sql -q "SELECT 1"       # Execute SQL
  ticloud fs secret ls               # List secrets
  ticloud fs grep "pattern" :/       # Search content
  ticloud fs find :/ -name "*.txt"   # Find files
  ticloud fs shell                   # Interactive shell
`,
	}

	fsCmd.AddCommand(initCmd(h))
	fsCmd.AddCommand(lsCmd(h))
	fsCmd.AddCommand(catCmd(h))
	fsCmd.AddCommand(cpCmd(h))
	fsCmd.AddCommand(mkdirCmd(h))
	fsCmd.AddCommand(statCmd(h))
	fsCmd.AddCommand(mvCmd(h))
	fsCmd.AddCommand(rmCmd(h))
	fsCmd.AddCommand(sqlCmd(h))
	fsCmd.AddCommand(secretCmd(h))
	fsCmd.AddCommand(shellCmd(h))
	fsCmd.AddCommand(grepCmd(h))
	fsCmd.AddCommand(findCmd(h))
	fsCmd.AddCommand(mountCmd(h))
	fsCmd.AddCommand(umountCmd(h))

	fsCmd.PersistentFlags().String(flag.FSEndpoint, "", "FS server URL (overrides config)")
	fsCmd.PersistentFlags().String(flag.FSClusterID, "", "Associated TiDB Cloud cluster ID")
	fsCmd.PersistentFlags().String(flag.FSZeroInstanceID, "", "Associated TiDB Zero instance ID")

	return fsCmd
}

// newClient creates a new FS client using the active TiDB Cloud CLI authentication.
// It reuses OAuth tokens or API key pairs from the current profile.
func newClient(cmd *cobra.Command) (*fs.Client, error) {
	// Resolve server URL: flag > env > config > default
	server, _ := cmd.Flags().GetString(flag.FSEndpoint)
	if server == "" {
		server = os.Getenv("TICLOUD_FS_ENDPOINT")
	}
	if server == "" {
		server = config.GetFSEndpoint()
	}

	// Resolve database association: flag > env > config
	clusterID, _ := cmd.Flags().GetString(flag.FSClusterID)
	if clusterID == "" {
		clusterID = os.Getenv("TICLOUD_FS_CLUSTER_ID")
	}
	if clusterID == "" {
		clusterID = config.GetFSClusterID()
	}

	zeroInstanceID, _ := cmd.Flags().GetString(flag.FSZeroInstanceID)
	if zeroInstanceID == "" {
		zeroInstanceID = os.Getenv("TICLOUD_FS_ZERO_INSTANCE_ID")
	}
	if zeroInstanceID == "" {
		zeroInstanceID = config.GetFSZeroInstanceID()
	}

	// Build HTTP transport using the same auth logic as the root CLI.
	var transport http.RoundTripper
	publicKey, privateKey := config.GetPublicKey(), config.GetPrivateKey()
	if publicKey != "" && privateKey != "" {
		transport = cloud.NewDigestTransport(publicKey, privateKey)
	} else {
		if err := config.ValidateToken(); err != nil {
			return nil, errors.New("no valid auth found: please login with 'ticloud auth login' or configure API keys with 'ticloud config set public-key <key>' and 'ticloud config set private-key <key>'")
		}
		token, err := config.GetAccessToken()
		if err != nil {
			if errors.Is(err, keyring.ErrNotFound) || errors.Is(err, store.ErrNotSupported) {
				return nil, errors.New("no valid auth found: please login with 'ticloud auth login' or configure API keys")
			}
			return nil, errors.Trace(err)
		}
		transport = cloud.NewBearTokenTransport(token)
	}

	httpClient := &http.Client{
		Transport: transport,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) > 0 && req.URL.Host != via[0].URL.Host {
				req.Header.Del("Authorization")
			}
			if len(via) >= 10 {
				return fmt.Errorf("too many redirects")
			}
			return nil
		},
	}

	return fs.NewClient(server, httpClient, clusterID, zeroInstanceID), nil
}

// suggestInitIfTenantNotFound wraps errors that indicate the tenant has not been
// provisioned yet with a helpful hint to run 'ticloud fs init'.
func suggestInitIfTenantNotFound(err error) error {
	if err == nil {
		return nil
	}
	msg := err.Error()
	if strings.Contains(strings.ToLower(msg), "tenant not found") {
		return fmt.Errorf("%w\n\nFS tenant not found for this database. Please run 'ticloud fs init' first.", err)
	}
	return err
}

// RemotePath represents a parsed remote path.
type RemotePath struct {
	Context string
	Path    string
}

// ParseRemotePath parses a remote path string.
// Format: [context:]/path or :/path (uses default context)
func ParseRemotePath(s string) RemotePath {
	if !strings.HasPrefix(s, ":") {
		return RemotePath{Path: s}
	}
	s = s[1:] // remove leading :
	if !strings.HasPrefix(s, "/") {
		// Check for context prefix
		if idx := strings.Index(s, ":/"); idx != -1 {
			return RemotePath{
				Context: s[:idx],
				Path:    s[idx+1:],
			}
		}
	}
	return RemotePath{Path: s}
}

func (rp RemotePath) String() string {
	if rp.Context != "" {
		return ":" + rp.Context + ":" + rp.Path
	}
	return ":" + rp.Path
}

// IsRemote returns true if the path is a remote path (starts with :).
func IsRemote(s string) bool {
	return strings.HasPrefix(s, ":")
}

