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
	"crypto/tls"
	"errors"
	"net"
	"net/url"
	"os/user"
	"strings"
	"testing"

	"github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/require"
	"github.com/xo/dburl"
	"github.com/xo/usql/handler"
	"github.com/xo/usql/rline"
)

func TestShellDSNCredentialsAndTLS(t *testing.T) {
	for _, clusterType := range []string{SERVERLESS, DEDICATED} {
		for _, host := range []string{"db.example", "127.0.0.1"} {
			for _, password := range []string{"", "p@ss:word?/#%&=", "p%41ss", "p@evil.example:4000?allowAllFiles=true#"} {
				t.Run(clusterType+"/"+host+"/"+password, func(t *testing.T) {
					username := "tenant.user@?/#%"
					dsn, err := generateDsnWithPassword(clusterType, username, host, "4000", password)
					require.NoError(t, err)
					assertShellDSN(t, dsn, clusterType, username, password, host)

					// Exercise usql's real password prompt and DSN reserialization.
					line := &rline.Rline{Pw: func(string) (string, error) { return password, nil }}
					h := handler.New(line, &user.User{Username: "local-user"}, t.TempDir(), true)
					dsn, err = generateDsnWithoutPassword(clusterType, username, host, "4000", h)
					require.NoError(t, err)
					assertShellDSN(t, dsn, clusterType, username, password, host)
				})
			}
		}
	}
}

func assertShellDSN(t *testing.T, dsn, clusterType, username, password, host string) {
	t.Helper()
	u, err := dburl.Parse(dsn)
	require.NoError(t, err)
	require.Equal(t, host, u.Hostname())
	require.Equal(t, "4000", u.Port())
	require.Len(t, u.Query(), 1)
	cfg, err := mysql.ParseDSN(u.DSN)
	require.NoError(t, err)
	require.Equal(t, username, cfg.User)
	require.Equal(t, password, cfg.Passwd)
	require.Equal(t, net.JoinHostPort(host, "4000"), cfg.Addr)
	require.NotNil(t, cfg.TLS)
	require.False(t, cfg.AllowAllFiles)
	require.False(t, cfg.MultiStatements)
	require.False(t, cfg.InterpolateParams)
	require.False(t, cfg.AllowCleartextPasswords)
	require.False(t, cfg.AllowFallbackToPlaintext)
	if clusterType == SERVERLESS {
		require.Equal(t, "tidb", cfg.TLSConfig)
		require.Equal(t, host, cfg.TLS.ServerName)
		require.Equal(t, uint16(tls.VersionTLS12), cfg.TLS.MinVersion)
		require.False(t, cfg.TLS.InsecureSkipVerify)
	} else {
		require.Equal(t, "skip-verify", cfg.TLSConfig)
	}
}

func TestShellDSNRejectsInvalidAddress(t *testing.T) {
	for _, clusterType := range []string{SERVERLESS, DEDICATED} {
		for _, address := range []struct{ host, port string }{
			{"attacker.example?allowAllFiles=true&x=", "4000"},
			{"attacker.example:4000?allowAllFiles=true#", "4000"},
			{"attacker.example&multiStatements=true", "4000"},
			{"user@attacker.example", "4000"},
			{"db.example/path", "4000"},
			{"db.example:4000", "4000"},
			{"https://db.example", "4000"},
			{"db.example%3fallowAllFiles=true", "4000"},
			{"db.example\n", "4000"},
			{"", "4000"}, {"db..example", "4000"},
			{"-db.example", "4000"}, {"db-.example", "4000"},
			{strings.Repeat("a", 64) + ".example", "4000"},
			{"db.example", "4000?allowAllFiles=true#"},
			{"db.example", ""}, {"db.example", "0"},
			{"db.example", "-1"}, {"db.example", "65536"},
		} {
			t.Run(clusterType+"/"+address.host+"/"+address.port, func(t *testing.T) {
				dsn, err := generateDsnWithPassword(clusterType, "user", address.host, address.port, "test-password")
				require.Error(t, err)
				require.Empty(t, dsn)
				// A nil handler ensures invalid addresses are rejected before prompting.
				dsn, err = generateDsnWithoutPassword(clusterType, "user", address.host, address.port, nil)
				require.Error(t, err)
				require.Empty(t, dsn)
			})
		}
	}
}

func TestShellDSNIPv6AndPortBounds(t *testing.T) {
	for _, port := range []string{"1", "65535"} {
		dsn, err := generateDsnWithPassword(SERVERLESS, "user", "2001:db8::1", port, "password")
		require.NoError(t, err)
		u, err := url.Parse(dsn)
		require.NoError(t, err)
		require.Equal(t, "2001:db8::1", u.Hostname())
		require.Equal(t, port, u.Port())
	}
}

func TestShellDSNPasswordInterrupt(t *testing.T) {
	line := &rline.Rline{Pw: func(string) (string, error) { return "", rline.ErrInterrupt }}
	h := handler.New(line, &user.User{Username: "local-user"}, t.TempDir(), true)
	_, err := generateDsnWithoutPassword(SERVERLESS, "user", "db.example", "4000", h)
	require.ErrorIs(t, err, InterruptError)
}

func TestConnectionErrorOmitsDSN(t *testing.T) {
	dsn, err := generateDsnWithPassword(SERVERLESS, "user", "example.com", "4000", "sql-secret")
	require.NoError(t, err)
	for _, cause := range []error{errors.New("connection refused"), errors.New("invalid connection " + dsn)} {
		message := connectionError(dsn, cause).Error()
		require.NotContains(t, message, "sql-secret")
		require.NotContains(t, message, dsn)
		require.Contains(t, message, "can't open connection")
	}
	require.Contains(t, connectionError(dsn, errors.New("connection refused")).Error(), "connection refused")
}
