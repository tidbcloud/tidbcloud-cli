// Copyright 2024 PingCAP, Inc.
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
	"fmt"
	"io"
	"os"
	"os/user"

	"github.com/go-sql-driver/mysql"
	"github.com/xo/usql/env"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

func ExecuteSqlDialog(ctx context.Context, clusterType, userName, host, port string, pass *string, out io.Writer) error {
	u, err := user.Current()
	if err != nil {
		return fmt.Errorf("can't get current user: %s", err.Error())
	}
	l, err := rline.New(false, "", env.HistoryFile(u))
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
		return fmt.Errorf("can't open connection to %s: %s", dsn, err.Error())
	}

	if err = h.Run(); err != io.EOF {
		return err
	}
	return nil
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
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s/test?tls=tidb", userName, pass, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s:%s@%s:%s/test?tls=skip-verify", userName, pass, host, port)
	} else {
		return "", fmt.Errorf("unsupproted cluster type: %s", clusterType)
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
		dsn = fmt.Sprintf("tidb://%s@%s:%s/test?tls=tidb", userName, host, port)
	} else if clusterType == DEDICATED {
		dsn = fmt.Sprintf("tidb://%s@%s:%s/test?tls=skip-verify", userName, host, port)
	} else {
		return "", fmt.Errorf("unsupproted cluster type: %s", clusterType)
	}

	// Prompt for password
	dsn, err := h.Password(dsn)
	if err != nil {
		return "", err
	}
	return dsn, nil
}
