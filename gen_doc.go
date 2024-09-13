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

package main

import (
	"log"

	"github.com/tidbcloud/tidbcloud-cli/internal"
	"github.com/tidbcloud/tidbcloud-cli/internal/cli"

	"github.com/spf13/cobra/doc"
)

func main() {
	h := &internal.Helper{}
	ticloud := cli.RootCmd(h)
	ticloud.DisableAutoGenTag = true
	ticloud.InitDefaultCompletionCmd()
	err := doc.GenMarkdownTree(ticloud, "./docs/generate_doc")
	if err != nil {
		log.Fatal(err)
	}
}
