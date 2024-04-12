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

package internal

import (
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/service/aws/s3"
	"tidbcloud-cli/internal/service/cloud"
)

const (
	DefaultPageSize = 100
)

type Helper struct {
	Client        func() (cloud.TiDBCloudClient, error)
	Uploader      func(client cloud.TiDBCloudClient) s3.Uploader
	QueryPageSize int64
	IOStreams     *iostream.IOStreams
}
