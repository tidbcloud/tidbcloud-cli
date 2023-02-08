// Copyright 2023 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package telemetry

import (
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
	"github.com/pingcap/log"
	"go.uber.org/zap"
)

const url = "https://telemetry.pingcap.com/api/v1/ticloud/report"

type EventsSender interface {
	SendEvent(event interface{}) error
}

type Sender struct {
	client *resty.Client
}

func NewSender() *Sender {
	client := resty.New()
	client.SetTimeout(1 * time.Second)

	return &Sender{
		client: client,
	}
}

func (s *Sender) SendEvent(event interface{}) error {
	log.Debug("sending telemetry events", zap.Any("event", event))
	response, err := s.client.
		R().SetBody(event).
		Post(url)

	if err != nil {
		return fmt.Errorf("failed to send telemetry events: %w", err)
	}

	if !response.IsSuccess() {
		return errors.Errorf("failed to send telemetry events: %s", string(response.Body()))
	}

	return nil
}
