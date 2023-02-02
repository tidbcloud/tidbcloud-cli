package telemetry

import (
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/juju/errors"
)

const url = "https://telemetry.pingca.com/api/v1/ticloud/report"

type EventsSender interface {
	SendEvents(body interface{}) error
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

func (s *Sender) SendEvents(body interface{}) error {
	response, err := s.client.
		R().SetBody(body).
		Post(url)

	if err != nil {
		return errors.Annotate(err, "failed to send telemetry events")
	}

	if !response.IsSuccess() {
		return errors.Errorf("failed to send telemetry events: %s", string(response.Body()))
	}

	return nil
}
