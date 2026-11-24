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
	"net"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"unicode"

	"github.com/tidbcloud/tidbcloud-cli/internal/redact"
)

var (
	errDownloadURL           = errors.New("download URL must use HTTPS and include a host")
	errDownloadRedirect      = errors.New("download redirect must remain on the same HTTPS origin")
	errDownloadRedirectLimit = errors.New("stopped after 10 redirects")
)

// GetResponse returns the response of a given AWS pre-signed URL.
func GetResponse(url string, debug bool) (*http.Response, error) {
	httpClient := &http.Client{
		Transport:     downloadTransport{inner: redact.NewDebugTransport(http.DefaultTransport, debug)},
		CheckRedirect: checkDownloadRedirect,
	}
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, &downloadRequestError{cause: err}
	}
	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		// Storage error bodies can echo the signed request or its credentials.
		return nil, fmt.Errorf("download failed with HTTP status %d", resp.StatusCode)
	}
	if resp.ContentLength <= 0 {
		resp.Body.Close()
		return nil, fmt.Errorf("file is empty")
	}
	return resp, nil
}

// Request errors may include the full URL, including a redirect target.
// Render known failure reasons and keep the original error for errors.Is/As.
type downloadRequestError struct{ cause error }

func (e *downloadRequestError) Error() string {
	var dns *net.DNSError
	var cert *tls.CertificateVerificationError
	var timeout net.Error
	switch {
	case errors.Is(e.cause, errDownloadURL):
		return errDownloadURL.Error()
	case errors.Is(e.cause, errDownloadRedirect):
		return errDownloadRedirect.Error()
	case errors.Is(e.cause, errDownloadRedirectLimit):
		return errDownloadRedirectLimit.Error()
	case errors.Is(e.cause, context.Canceled):
		return "download request canceled"
	case errors.As(e.cause, &timeout) && timeout.Timeout():
		return "download request timed out"
	case errors.As(e.cause, &dns):
		return "download request failed: could not resolve the host"
	case errors.As(e.cause, &cert):
		return "download request failed: TLS certificate verification failed"
	default:
		return "download request failed: check the URL, network, and TLS settings"
	}
}

func (e *downloadRequestError) Unwrap() error { return e.cause }

// Validate every request before sending pre-signed URL credentials.
type downloadTransport struct {
	inner http.RoundTripper
}

func (t downloadTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	if req.URL.Scheme != "https" || req.URL.Hostname() == "" {
		return nil, errDownloadURL
	}
	return t.inner.RoundTrip(req)
}

// Storage URLs can differ from the API origin, but redirects must remain on
// the original download URL's HTTPS origin.
func checkDownloadRedirect(req *http.Request, via []*http.Request) error {
	if len(via) >= 10 {
		return errDownloadRedirectLimit
	}
	origin := via[0].URL
	originPort, targetPort := origin.Port(), req.URL.Port()
	if originPort == "" {
		originPort = "443"
	}
	if targetPort == "" {
		targetPort = "443"
	}
	if req.URL.Scheme != "https" || origin.Scheme != "https" ||
		!strings.EqualFold(req.URL.Hostname(), origin.Hostname()) || targetPort != originPort {
		return errDownloadRedirect
	}
	return nil
}

// CreateFile creates a file if it does not exist
func CreateFile(path, fileName string) (*os.File, error) {
	filePath, err := safeDownloadPath(fileName)
	if err != nil {
		return nil, err
	}

	if path == "" {
		path = "."
	}
	root, err := os.OpenRoot(path)
	if err != nil {
		return nil, err
	}
	defer root.Close()

	file, err := root.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_EXCL, 0600)
	if err != nil {
		if os.IsExist(err) {
			return nil, fmt.Errorf("file already exists")
		}
		return nil, fmt.Errorf("download destination %q cannot be created: %w", fileName, err)
	}
	return file, nil
}

func safeDownloadPath(fileName string) (string, error) {
	if fileName == "" {
		return "", downloadDestinationError(fileName, "cannot be empty")
	}
	if fileName == "." || fileName == ".." {
		return "", downloadDestinationError(fileName, "must refer to a file")
	}
	if filepath.IsAbs(fileName) || isWindowsAbs(fileName) {
		return "", downloadDestinationError(fileName, "must be relative to the output path")
	}
	if containsControlCharacter(fileName) {
		return "", downloadDestinationError(fileName, "contains unsupported characters")
	}

	clean := filepath.Clean(fileName)
	if clean == "." || clean == ".." {
		return "", downloadDestinationError(fileName, "must refer to a file")
	}
	if strings.HasPrefix(clean, ".."+string(os.PathSeparator)) || filepath.IsAbs(clean) {
		return "", downloadDestinationError(fileName, "is outside the output path")
	}

	return clean, nil
}

func downloadDestinationError(fileName, reason string) error {
	return fmt.Errorf("download destination %q %s", fileName, reason)
}

func containsControlCharacter(s string) bool {
	return strings.ContainsRune(s, 0) || strings.IndexFunc(s, unicode.IsControl) >= 0
}

func isWindowsAbs(path string) bool {
	if runtime.GOOS == "windows" {
		return false
	}
	if len(path) >= 3 && path[1] == ':' && (path[2] == '\\' || path[2] == '/') {
		c := path[0]
		return ('a' <= c && c <= 'z') || ('A' <= c && c <= 'Z')
	}
	return strings.HasPrefix(path, `\\`)
}

// CreateFolder creates a folder if it does not exist
func CreateFolder(path string) error {
	if path == "" {
		return nil
	}
	if _, err := os.Stat(path); err != nil {
		if os.IsNotExist(err) {
			err = os.MkdirAll(path, 0700)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
