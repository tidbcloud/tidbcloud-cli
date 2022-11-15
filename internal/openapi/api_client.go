// Copyright 2022 PingCAP, Inc.
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
