package openapi

import (
	"net/http"

	apiClient "github.com/c4pt0r/go-tidbcloud-sdk-v1/client"
	httpTransport "github.com/go-openapi/runtime/client"
	"github.com/go-openapi/strfmt"
	"github.com/icholy/digest"
)

func NewApiClient(publicKey string, privateKey string) *apiClient.GoTidbcloud {
	httpclient := &http.Client{
		Transport: &digest.Transport{
			Username: publicKey,
			Password: privateKey,
		},
	}
	return apiClient.New(httpTransport.NewWithClient("api.tidbcloud.com", "/", []string{"https"}, httpclient), strfmt.Default)
}
