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

package util

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"os"
	"os/user"
	"time"

	"github.com/go-sql-driver/mysql"
	isatty "github.com/mattn/go-isatty"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

func ExecuteSqlDialog(ctx context.Context, clusterType, userName, host, port string, pass *string, out io.Writer) error {
	h, err := newSQLHandler()
	if err != nil {
		return err
	}

	var dsn string
	if pass == nil {
		dsn, err = generateDsnWithoutPassword(clusterType, userName, host, port, h)
		if err != nil {
			return err
		}
	} else {
		dsn, err = generateDsnWithPassword(clusterType, userName, host, port, *pass)
		if err != nil {
			return err
		}
	}

	return runSQLDialog(ctx, h, dsn, dsn, 0)
}

// ExecuteSqlDialogWithTLS opens a SQL shell using a caller-provided TLS
// configuration. It is used by APIs that publish their own CA certificate.
// A positive connectTimeout bounds initial connection setup after password
// collection, not the interactive SQL session.
func ExecuteSqlDialogWithTLS(ctx context.Context, userName, host, port string, pass *string, tlsConfig *tls.Config, connectTimeout time.Duration) error {
	h, err := newSQLHandler()
	if err != nil {
		return err
	}

	const tlsConfigName = "tidb-nextgen"
	if err := mysql.RegisterTLSConfig(tlsConfigName, tlsConfig); err != nil {
		return err
	}

	var dsn string
	if pass == nil {
		dsn, err = h.Password(buildTiDBShellDSN(userName, host, port, nil, tlsConfigName))
		if err != nil && errors.Is(err, rline.ErrInterrupt) {
			return InterruptError
		}
		if err != nil {
			return err
		}
	} else {
		dsn = buildTiDBShellDSN(userName, host, port, pass, tlsConfigName)
	}

	return runSQLDialog(ctx, h, dsn, net.JoinHostPort(host, port), connectTimeout)
}

func newSQLHandler() (*handler.Handler, error) {
	u, err := user.Current()
	if err != nil {
		return nil, fmt.Errorf("can't get current user: %s", err.Error())
	}
	// https://github.com/xo/usql/commit/074448a65adcebe1391879a15ffe8c16493bf9fa
	interactive := isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stdin.Fd())
	cygwin := isatty.IsCygwinTerminal(os.Stdout.Fd()) && isatty.IsCygwinTerminal(os.Stdin.Fd())
	l, err := rline.New(interactive, cygwin, false, "", env.HistoryFile(u))
	if err != nil {
		return nil, fmt.Errorf("can't open history file: %s", err.Error())
	}
	wd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	return handler.New(l, u, wd, true), nil
}

type sqlDialog interface {
	Open(context.Context, ...string) error
	Run() error
}

func runSQLDialog(ctx context.Context, h sqlDialog, dsn, endpoint string, connectTimeout time.Duration) error {
	// Password collection is complete before this point. Only bound initial
	// connection setup, never the interactive SQL session or later queries.
	connectCtx := ctx
	cancel := func() {}
	if connectTimeout > 0 {
		connectCtx, cancel = context.WithTimeout(ctx, connectTimeout)
	}
	err := h.Open(connectCtx, dsn)
	// usql can swallow a timeout while fetching the server version after Ping.
	// Inspect the deadline before canceling it, and preserve the legacy path.
	if err == nil && connectTimeout > 0 {
		err = connectCtx.Err()
	}
	cancel()
	if err != nil {
		return fmt.Errorf("can't open connection to %s: %w", endpoint, err)
	}
	if err := h.Run(); err != io.EOF {
		return err
	}
	return nil
}

func buildTiDBShellDSN(userName, host, port string, pass *string, tlsConfigName string) string {
	userInfo := url.User(userName)
	if pass != nil {
		userInfo = url.UserPassword(userName, *pass)
	}
	query := url.Values{"tls": []string{tlsConfigName}}
	return (&url.URL{
		Scheme:   "tidb",
		User:     userInfo,
		Host:     net.JoinHostPort(host, port),
		RawQuery: query.Encode(),
	}).String()
}

func generateDsnWithPassword(clusterType string, userName string, host string, port string, pass string) (string, error) {
	var dsn string
	if clusterType == SERVERLESS {
		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if err != nil {
			return "", err
		}
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s?tls=tidb", userName, pass, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s?tls=skip-verify", userName, pass, host, port)
	} else {
		return "", fmt.Errorf("unsupported cluster type: %s", clusterType)
	}
	return dsn, nil
}

func generateDsnWithoutPassword(clusterType string, userName string, host string, port string, h *handler.Handler) (string, error) {
	var dsn string
	if clusterType == SERVERLESS {
		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if err != nil {
			return "", err
		}
		dsn = fmt.Sprintf("tidb://%s@%s:%s?tls=tidb", userName, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s@%s:%s?tls=skip-verify", userName, host, port)
	} else {
		return "", fmt.Errorf("unsupported cluster type: %s", clusterType)
	}

	// Prompt for password
	dsn, err := h.Password(dsn)
	if err != nil && errors.Is(err, rline.ErrInterrupt) {
		return "", InterruptError
	}
	return dsn, err
}
