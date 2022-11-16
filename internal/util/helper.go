package util

import "tidbcloud-cli/internal/iostream"

const (
	DefaultPageSize = 100
)

type Helper struct {
	Client        func() CloudClient
	QueryPageSize int64
	IOStreams     *iostream.IOStreams
}
