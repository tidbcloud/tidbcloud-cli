package cli

import (
	"os"

	"tidbcloud-cli/internal"
	"tidbcloud-cli/internal/service/aws/s3"

	"github.com/spf13/cobra"
)

func UploadCmd(h *internal.Helper) *cobra.Command {
	var projectCmd = &cobra.Command{
		Use: "upload",
		RunE: func(cmd *cobra.Command, args []string) error {
			client, err := h.Client()
			if err != nil {
				return err
			}
			file, err := os.Open("/Users/leon/Downloads/sysbench-1.0.20.tar.gz")
			if err != nil {
				return err
			}

			key := "yx_test"
			stat, err := file.Stat()
			if err != nil {
				return err
			}
			length := stat.Size()
			clusterID := "10092819840638843800"
			err = s3.NewUploader(client).Upload(cmd.Context(),
				&s3.PutObjectInput{
					Key:           &key,
					ContentLength: &length,
					ClusterID:     &clusterID,
					Body:          file,
				})
			if err != nil {
				return err
			}

			return nil
		},
	}

	return projectCmd
}
