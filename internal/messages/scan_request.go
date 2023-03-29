package messages

type ScanRequest struct {
	BucketName string
	ObjectKey  string
}

func NewScanRequest(bucketName, objectKey string) *ScanRequest {
	return &ScanRequest{BucketName: bucketName, ObjectKey: objectKey}
}
