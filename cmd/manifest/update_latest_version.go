// Copyright 2025 PingCAP, Inc.
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

	latestTag := *latestRelease.TagName
	err = os.WriteFile("latest-version", []byte(latestTag), 0600)
	if err != nil {
		log.Fatalf("Failed to write latest-version file: %s", err.Error())
	}

	// update the version in the installation script
	file, err := os.ReadFile("install.sh")
	if err != nil {
		log.Fatalf("failed to open install.sh: %s", err.Error())
	}

	latestTag = strings.TrimPrefix(latestTag, "v")
	file = regexp.MustCompile("version=\".*\"").ReplaceAll(file, []byte(fmt.Sprintf("version=\"%s\"", latestTag)))
	err = os.WriteFile("install.sh", file, 0600)
	if err != nil {
		log.Fatalf("Failed to write install.sh file: %s", err.Error())
	}
}
