package internal

import (
	"tidbcloud-cli/internal/config"
	"tidbcloud-cli/internal/iostream"
	"tidbcloud-cli/internal/util"
)

const (
	DefaultPageSize = 100
)

type Helper struct {
	Client        func() util.CloudClient
	QueryPageSize int64
	IOStreams     *iostream.IOStreams
	Config        *config.Config
}
