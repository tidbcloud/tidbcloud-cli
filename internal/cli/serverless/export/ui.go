package export

import (
	"github.com/charmbracelet/bubbles/textinput"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/juju/errors"
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/flag"
	"tidbcloud-cli/internal/ui"
	"tidbcloud-cli/internal/util"
)

var InputDescription = map[string]string{
	flag.S3URI:                "Input your S3 uri in s3://<bucket>/<path> format",
	flag.S3AccessKeyID:        "Input your S3 access key id",
	flag.S3SecretAccessKey:    "Input your S3 secret access key",
	flag.S3RoleArn:            "Input your S3 role arn",
	flag.AzureBlobURI:         "Input your Azure Blob uri in azure://<account>.blob.core.windows.net/<container>/<path> format",
	flag.AzureBlobSASToken:    "Input your Azure Blob SAS token",
	flag.GCSURI:               "Input your GCS uri in gcs://<bucket>/<path> format",
	flag.GCSServiceAccountKey: "Input your base64 encoded GCS service account key",
}

func InitialInputModel(inputs []string) (ui.TextInputModel, error) {
	m := ui.TextInputModel{
		Inputs: make([]textinput.Model, len(inputs)),
	}
	for i, input := range inputs {
		t := textinput.New()
		if i == 0 {
			t.Focus()
			t.PromptStyle = config.FocusedStyle
			t.TextStyle = config.FocusedStyle
		}
		t.Placeholder = InputDescription[input]
		m.Inputs[i] = t
	}
	p := tea.NewProgram(m)
	inputModel, err := p.Run()
	finalModel := inputModel.(ui.TextInputModel)
	if err != nil {
		return finalModel, errors.Trace(err)
	}
	if finalModel.Interrupted {
		return finalModel, util.InterruptError
	}
	return finalModel, nil
}

func GetSelectedParquetCompression() (string, error) {
	compressions := make([]interface{}, 0, 4)
	compressions = append(compressions, "SNAPPY", "GZIP", "NONE")
	model, err := ui.InitialSelectModel(compressions, "Choose the parquet compression algorithm:")
	if err != nil {
		return "", errors.Trace(err)
	}

	p := tea.NewProgram(model)
	fileTypeModel, err := p.Run()
	if err != nil {
		return "", errors.Trace(err)
	}
	if m, _ := fileTypeModel.(ui.SelectModel); m.Interrupted {
		return "", util.InterruptError
	}
	compression := fileTypeModel.(ui.SelectModel).GetSelectedItem()
	if compression == nil {
		return "", errors.New("no compression algorithm selected")
	}
	return compression.(string), nil
}
