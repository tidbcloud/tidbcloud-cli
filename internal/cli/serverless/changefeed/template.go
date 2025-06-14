package changefeed

import (
	"fmt"

	"github.com/fatih/color"
	"github.com/spf13/cobra"

	"github.com/tidbcloud/tidbcloud-cli/internal"
)

const (
	KafkaInfoTemplate = `{
  "broker": {
    "address": "localhost:9092"
  },
  "authentication": {
    "mechanism": "PLAIN",
    "username": "user",
    "password": "pass"
  },
  "dataFormat": {
    "type": "canal-json"
  },
  "topicPartitionConfig": {
    "topic": "my_topic",
    "partitionNum": 3
  }
}`

	CDCFilterTemplate = `{
  "filterRule": ["test.t1", "test.t2"],
  "mode": "all",
  "eventFilterRule": [
    {
      "type": "insert",
      "value": "true"
    }
  ]
}`
)

func TemplateCmd(h *internal.Helper) *cobra.Command {
	return &cobra.Command{
		Use:   "template",
		Short: "Show changefeed KafkaInfo and CDCFilter JSON templates",
		Example: `  Show changefeed templates:
  $ tidbcloud serverless changefeed template`,
		Args: cobra.NoArgs,
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("KafkaInfo JSON template:"))
			fmt.Fprintln(h.IOStreams.Out, KafkaInfoTemplate)
			fmt.Fprintln(h.IOStreams.Out, color.GreenString("CDCFilter JSON template:"))
			fmt.Fprintln(h.IOStreams.Out, CDCFilterTemplate)
			return nil
		},
	}
}
