package s3

type Config struct {
	AccessKeyID string `json:"accessKeyId" jsonschema:"description=Access Key ID for S3 Credentials"`
	SecretKey   string `json:"secretKey" jsonschema:"description=Secret Key for S3 Credentials"`
	Token       string `json:"token,omitempty" jsonschema:"description=Token for S3 Credentials"`
	Region      string `json:"region,omitempty" jsonschema:"description=Region to reach the service,default=br-ne-1,enum=br-ne-1,enum=br-ne-2,enum=br-se-1"`
	MaxWorkers  int    `json:"max_workers,omitempty" jsonschema:"description=Max workers used to recursively execute jobs from the queue,default=10"`
	QueueSize   int64  `json:"queue_size,omitempty" jsonschema:"description=Size of the queue of jobs used to recursively delete files,default=1000"`
}
