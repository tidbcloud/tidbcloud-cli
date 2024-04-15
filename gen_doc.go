package main

import (
	"log"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/cli"

	"github.com/spf13/cobra/doc"
)

func main() {
	h := &internal.Helper{}
	kubectl := cli.RootCmd(h)
	err := doc.GenMarkdownTree(kubectl, "./docs/generate_doc")
	if err != nil {
		log.Fatal(err)
	}
}
