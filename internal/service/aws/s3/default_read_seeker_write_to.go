//go:build !windows
// +build !windows

package s3

import "github.com/aws/aws-sdk-go-v2/feature/s3/manager"

func defaultUploadBufferProvider() manager.ReadSeekerWriteToProvider {
	return nil
}
