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

package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"

	"github.com/google/go-github/v49/github"
	"golang.org/x/oauth2"
)

// releaseTagPattern matches the tags this project publishes (for example
// "v1.2.3" or "v1.2.3-rc.1+dev"). Anything outside it must never reach
// install.sh, because the value is interpolated into an executable script.
var releaseTagPattern = regexp.MustCompile(`^v[0-9][0-9A-Za-z.+-]*$`)

func main() {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("GITHUB_TOKEN")},
	)
	tc := oauth2.NewClient(ctx, ts)

	client := github.NewClient(tc)

	latestRelease, _, err := client.Repositories.GetLatestRelease(ctx, "tidbcloud", "tidbcloud-cli")
	if err != nil {
		log.Fatalf("failed to get latest release: %s", err.Error())
	}

	if latestRelease == nil || latestRelease.TagName == nil {
		log.Fatal("Can't get the latest release from Github api.")
	}

	if err := updateManifest(*latestRelease.TagName); err != nil {
		log.Fatal(err)
	}
}

func updateManifest(latestTag string) error {
	if !releaseTagPattern.MatchString(latestTag) {
		return fmt.Errorf("refusing unsafe release tag %q", latestTag)
	}

	// update the version in the installation script
	file, err := os.ReadFile("install.sh")
	if err != nil {
		return fmt.Errorf("failed to open install.sh: %w", err)
	}
	versionLine := regexp.MustCompile(`(?m)^version="[^"]*"$`)
	if len(versionLine.FindAllIndex(file, -1)) != 1 {
		return fmt.Errorf("install.sh must contain exactly one version assignment")
	}
	file = versionLine.ReplaceAllLiteral(file, []byte(fmt.Sprintf("version=\"%s\"", strings.TrimPrefix(latestTag, "v"))))
	if err := os.WriteFile("latest-version", []byte(latestTag), 0600); err != nil {
		return fmt.Errorf("failed to write latest-version: %w", err)
	}
	err = os.WriteFile("install.sh", file, 0600) //nolint:gosec // The path is the constant "install.sh".
	if err != nil {
		return fmt.Errorf("failed to write install.sh: %w", err)
	}
	return nil
}
