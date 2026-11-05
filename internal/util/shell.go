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
	"strconv"
	"strings"

	"github.com/tidbcloud/tidbcloud-cli/internal/redact"

	"github.com/go-sql-driver/mysql"
	isatty "github.com/mattn/go-isatty"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

func ExecuteSqlDialog(ctx context.Context, clusterType, userName, host, port string, pass *string, out io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("can't get current user: %s", err.Error())
	}
	// https://github.com/xo/usql/commit/074448a65adcebe1391879a15ffe8c16493bf9fa
	interactive := isatty.IsTerminal(os.Stdout.Fd()) && isatty.IsTerminal(os.Stdin.Fd())
	cygwin := isatty.IsCygwinTerminal(os.Stdout.Fd()) && isatty.IsCygwinTerminal(os.Stdin.Fd())
	l, err := rline.New(interactive, cygwin, false, "", env.HistoryFile(u))
	if err != nil {
		return fmt.Errorf("can't open history file: %s", err.Error())
	}
	wd, err := os.Getwd()
	if err != nil {
		return err
	}
	h := handler.New(l, u, wd, true)

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

	if err = h.Open(ctx, dsn); err != nil {
		return connectionError(dsn, err)
	}

	if err = h.Run(); err != io.EOF {
		return err
	}
	return nil
}

func generateDsnWithPassword(clusterType string, userName string, host string, port string, pass string) (string, error) {
	return generateDsn(clusterType, userName, host, port, &pass)
}

func generateDsn(clusterType, userName, host, port string, pass *string) (string, error) {
	if !validSQLHost(host) {
		return "", errors.New("SQL host must be a hostname or IP address without a port")
	}
	portNumber, err := strconv.ParseUint(port, 10, 16)
	if err != nil || portNumber == 0 {
		return "", errors.New("SQL port must be an integer between 1 and 65535")
	}
	var tlsMode string
	if clusterType == SERVERLESS {
		err := mysql.RegisterTLSConfig("tidb", &tls.Config{
			MinVersion: tls.VersionTLS12,
			ServerName: host,
		})
		if err != nil {
			return "", err
		}
		tlsMode = "tidb"
	} else if clusterType == DEDICATED {
		tlsMode = "skip-verify"
	} else {
		return "", fmt.Errorf("unsupported cluster type: %s", clusterType)
	}
	u := url.URL{
		Scheme:   "tidb",
		User:     url.User(userName),
		Host:     net.JoinHostPort(host, strconv.FormatUint(portNumber, 10)),
		RawQuery: url.Values{"tls": {tlsMode}}.Encode(),
	}
	if pass != nil {
		u.User = url.UserPassword(userName, *pass)
	}
	return u.String(), nil
}

func validSQLHost(host string) bool {
	if net.ParseIP(host) != nil {
		return true
	}
	host = strings.TrimSuffix(host, ".")
	if host == "" || len(host) > 253 {
		return false
	}
	for _, label := range strings.Split(host, ".") {
		if len(label) == 0 || len(label) > 63 || label[0] == '-' || label[len(label)-1] == '-' {
			return false
		}
		for _, c := range label {
			if !(c >= 'a' && c <= 'z' || c >= 'A' && c <= 'Z' || c >= '0' && c <= '9' || c == '-') {
				return false
			}
		}
	}
	return true
}

func generateDsnWithoutPassword(clusterType string, userName string, host string, port string, h *handler.Handler) (string, error) {
	dsn, err := generateDsn(clusterType, userName, host, port, nil)
	if err != nil {
		return "", err
	}

	// Prompt for password
	dsn, err = h.Password(dsn)
	if err != nil && errors.Is(err, rline.ErrInterrupt) {
		return "", InterruptError
	}
	return dsn, err
}

func connectionError(dsn string, err error) error {
	// The driver may also include the original DSN in its error.
	detail := strings.ReplaceAll(err.Error(), dsn, "<DSN omitted>")
	return fmt.Errorf("can't open connection: %s", redact.Text(detail))
}
