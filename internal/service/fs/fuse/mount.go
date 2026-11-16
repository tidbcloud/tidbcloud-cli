//go:build !windows

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

// Package fuse provides FUSE filesystem implementation for TiDB Cloud FS.
package fuse

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/tidbcloud/tidbcloud-cli/internal/service/fs"

	gofuse "github.com/hanwen/go-fuse/v2/fuse"
)

// MountOptions configures the FUSE mount.
type MountOptions struct {
	Client     *fs.Client // FS client (with auth and headers configured)
	MountPoint string     // Local mount point
	ReadOnly   bool       // Mount as read-only (MVP default)
	Debug      bool       // Enable FUSE debug logging
	AllowOther bool       // Allow other users to access mount
}

// Validate checks the mount options.
func (o *MountOptions) Validate() error {
	if o.Client == nil {
		return fmt.Errorf("FS client is required")
	}
	if o.MountPoint == "" {
		return fmt.Errorf("mount point is required")
	}
	return nil
}

// Mount creates and serves a FUSE mount. It blocks until the filesystem
// is unmounted or a signal (SIGINT, SIGTERM) is received.
func Mount(opts *MountOptions) error {
	if err := opts.Validate(); err != nil {
		return err
	}

	// Create mount point if not exists
	if err := os.MkdirAll(opts.MountPoint, 0o755); err != nil {
		return fmt.Errorf("create mount point: %w", err)
	}

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if _, err := opts.Client.List("/"); err != nil {
		return fmt.Errorf("cannot reach FS server: %w", err)
	}
	_ = ctx

	// Create FUSE filesystem
	fsys := NewReadOnlyFS(opts.Client, opts)

	// Configure mount options
	fuseOpts := &gofuse.MountOptions{
		FsName:       "ticloudfs",
		Name:         "ticloudfs",
		MaxReadAhead: 128 * 1024,
		Debug:        opts.Debug,
		AllowOther:   opts.AllowOther,
	}
	if opts.ReadOnly {
		fuseOpts.Options = append(fuseOpts.Options, "ro")
	}

	// Create FUSE server
	server, err := gofuse.NewServer(fsys, opts.MountPoint, fuseOpts)
	if err != nil {
		return fmt.Errorf("fuse mount: %w", err)
	}

	// Handle signals
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		fmt.Fprintf(os.Stderr, "\nUnmounting %s...\n", opts.MountPoint)
		if err := server.Unmount(); err != nil {
			fmt.Fprintf(os.Stderr, "Unmount error: %v\n", err)
		}
	}()

	fmt.Fprintf(os.Stderr, "Mounted on %s (server: %s)\n", opts.MountPoint, opts.Client.GetBaseURL())
	server.Serve()
	return nil
}
