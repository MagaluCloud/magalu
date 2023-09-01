package s3

const (
	authorizationHeaderKey = "Authorization"

	contentMD5Header = "Content-Md5"

	// ContentSHAKey is the SHA256 of request body
	contentSHAKey = "X-Amz-Content-Sha256"

	// EmptyStringSHA256 is the hex encoded sha256 value of an empty string
	emptyStringSHA256 = "e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b855"

	// headerDateKey is the UTC timestamp for the request in the format YYYYMMDD'T'HHMMSS'Z'
	headerDateKey = "X-Amz-Date"

	// TimeFormat is the time format to be used in the X-Amz-Date header or query parameter
	longTimeFormat = "20060102T150405Z"

	requestSuffix = "aws4_request"

	secretPrefix = "AWS4"

	signingAlgorithm = "AWS4-HMAC-SHA256"

	// Default region to sign payload
	signingRegion = "us-east-1"

	// Default service name to sign payload
	signingService = "s3"

	// ShortTimeFormat is the shorten time format used in the credential scope
	shortTimeFormat = "20060102"

	templateUrl = "https://{{region}}.api.magalu.cloud/magaluobjects"

	unsignedPayloadHeader = "UNSIGNED-PAYLOAD"
)
