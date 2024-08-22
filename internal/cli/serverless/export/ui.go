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

//func initialS3AccessKeyInputModel() ui.TextInputModel {
//	m := ui.TextInputModel{
//		Inputs: make([]textinput.Model, len(S3AccessKeyInputFields)),
//	}
//	for k, v := range S3AccessKeyInputFields {
//		t := textinput.New()
//		switch k {
//		case flag.S3URI:
//			t.Placeholder = "Input your S3 uri in s3://<bucket>/<path> format"
//			t.Focus()
//			t.PromptStyle = config.FocusedStyle
//			t.TextStyle = config.FocusedStyle
//		case flag.S3AccessKeyID:
//			t.Placeholder = "Input your S3 access key id"
//		case flag.S3SecretAccessKey:
//			t.Placeholder = "Input your S3 secret access key"
//		}
//		m.Inputs[v] = t
//	}
//	return m
//}
//
//func initialS3RoleArnInputModel() ui.TextInputModel {
//	m := ui.TextInputModel{
//		Inputs: make([]textinput.Model, len(S3RoleArnInputFields)),
//	}
//	for k, v := range S3RoleArnInputFields {
//		t := textinput.New()
//		switch k {
//		case flag.S3URI:
//			t.Placeholder = "Input your S3 uri in s3://<bucket>/<path> format"
//			t.Focus()
//			t.PromptStyle = config.FocusedStyle
//			t.TextStyle = config.FocusedStyle
//		case flag.S3RoleArn:
//			t.Placeholder = "Input your S3 role arn"
//		}
//		m.Inputs[v] = t
//	}
//	return m
//}
//
//func initialGSCServiceAccountKeyInputModel() ui.TextInputModel {
//	m := ui.TextInputModel{
//		Inputs: make([]textinput.Model, len(GSCServiceAccountKeyInputFields)),
//	}
//	for k, v := range GSCServiceAccountKeyInputFields {
//		t := textinput.New()
//		switch k {
//		case flag.GCSURI:
//			t.Placeholder = "Input your S3 uri in s3://<bucket>/<path> format"
//			t.Focus()
//			t.PromptStyle = config.FocusedStyle
//			t.TextStyle = config.FocusedStyle
//		case flag.GCSServiceAccountKey:
//			t.Placeholder = "Input your S3 role arn"
//		}
//		m.Inputs[v] = t
//	}
//	return m
//}
//
//func initialAzBlobSasTokenInputModel() ui.TextInputModel {
//	m := ui.TextInputModel{
//		Inputs: make([]textinput.Model, len(AuthTypeAzBlobSasToken)),
//	}
//	for k, v := range AuthTypeAzBlobSasToken {
//		t := textinput.New()
//		switch k {
//		case flag.S3URI:
//			t.Placeholder = "Input your S3 uri in s3://<bucket>/<path> format"
//			t.Focus()
//			t.PromptStyle = config.FocusedStyle
//			t.TextStyle = config.FocusedStyle
//		case flag.S3RoleArn:
//			t.Placeholder = "Input your S3 role arn"
//		}
//		m.Inputs[v] = t
//	}
//	return m
//}
//
//func GetExternalStorageInput(authType AuthType) (tea.Model, error) {
//
//	p := tea.NewProgram(initialS3AccessKeyInputModel())
//	inputModel, err := p.Run()
//	if err != nil {
//		return nil, errors.Trace(err)
//	}
//	if inputModel.(ui.TextInputModel).Interrupted {
//		return nil, util.InterruptError
//	}
//	return inputModel, nil
//}
