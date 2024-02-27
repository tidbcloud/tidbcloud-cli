package s3

import "github.com/aws/aws-sdk-go-v2/feature/s3/manager"

// BufferedReadSeekerWriteToPool is only used by default when executing in Windows environments.
// It exists to work around constraints with Go standard library's os file implementations for Windows.
func defaultUploadBufferProvider() manager.ReadSeekerWriteToProvider {
	return manager.NewBufferedReadSeekerWriteToPool(1024 * 1024)
}
