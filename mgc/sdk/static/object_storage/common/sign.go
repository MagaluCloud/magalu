package common

import (
	"bytes"
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"slices"

	"magalu.cloud/core"
)

type HeaderMap = map[string]any

type SignatureParameters struct {
	Algorithm     string
	AccessKey     string
	Credential    string
	Scope         string
	Date          string
	ShortDate     string
	PayloadHash   string
	SignedHeaders []string // must be lower-case and sorted, see https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
}

func NewSignatureParameters(
	accessKey string,
	signingTime time.Time,
	payloadHash string,
	signedHeaders []string,
) SignatureParameters {
	shortDate := signingTime.Format(shortTimeFormat)
	scope := strings.Join([]string{
		shortDate,
		signingRegion,
		signingService,
		requestSuffix,
	}, "/")

	return SignatureParameters{
		Algorithm:     signingAlgorithm,
		AccessKey:     accessKey,
		Credential:    fmt.Sprintf("%s/%s", accessKey, scope),
		Scope:         scope,
		Date:          signingTime.Format(longTimeFormat),
		ShortDate:     shortDate,
		PayloadHash:   payloadHash,
		SignedHeaders: signedHeaders,
	}
}

type SignatureContext struct {
	Parameters SignatureParameters

	// Request
	HTTPMethod       string
	CanonicalURI     string
	CanonicalQuery   string
	CanonicalHeaders string

	SignedHeaders string

	// these are set by Sign
	Signature string
}

func NewSignatureContext(
	params SignatureParameters,
	req *http.Request,
) *SignatureContext {
	canonicalHeaders := buildCanonicalHeaders(req, params.SignedHeaders)

	return &SignatureContext{
		Parameters: params,

		// Request
		HTTPMethod:       req.Method,
		CanonicalURI:     req.URL.EscapedPath(),
		CanonicalQuery:   buildCanonicalQuery(req.URL.Query()),
		CanonicalHeaders: canonicalHeaders,

		SignedHeaders: strings.Join(params.SignedHeaders, ";"),
	}
}

// sorted as per https://docs.aws.amazon.com/AmazonS3/latest/API/sig-v4-header-based-auth.html
func buildCanonicalQuery(query url.Values) string {
	type p struct {
		key, value string
	}
	pairs := make([]p, 0, len(query))
	for key, values := range query {
		for _, value := range values {
			pairs = append(pairs, p{url.QueryEscape(key), url.QueryEscape(value)})
		}
	}

	slices.SortFunc(pairs, func(a, b p) int {
		return strings.Compare(a.key, b.key)
	})

	var result string
	for _, pair := range pairs {
		if result != "" {
			result += "&"
		}
		result += pair.key + "=" + pair.value
	}

	return result
}

func setContentHeader(req *http.Request, unsigned bool) (payloadHash string, err error) {
	if unsigned {
		req.Header.Set(contentSHAKey, unsignedPayloadHeader)
		return unsignedPayloadHeader, nil
	}
	// TODO: if we are able to receive a Seeker interface we must treat differently
	payloadHash, err = getPayloadHash(req)
	if err != nil {
		return "", err
	}
	req.Header.Set(contentSHAKey, payloadHash)
	return payloadHash, nil
}

/*
Computes the hash of the payload from the current request. We need to clone the
request in order to safely read the body stream. If body is empty (i.e., GET requests),
a default hashed string (emptyStringSHA256) is used based on Sig V4 specs.
*/
func getPayloadHash(req *http.Request) (string, error) {
	// Need to clone in order to safely consume body reader
	if req.Body == nil {
		return emptyStringSHA256, nil
	}
	bodyReader, err := req.GetBody()
	if err != nil {
		return "", err
	}

	defer bodyReader.Close()
	return core.SHA256Hex(bodyReader)
}

// setMD5Checksum computes the MD5 of the request payload and sets it to the
// Content-MD5 header. Returning the MD5 base64 encoded string or error.
//
// If the MD5 is already set as the Content-MD5 header, that value will be
// returned, and nothing else will be done.
//
// If the payload is empty, no MD5 will be computed. No error will be returned.
// Empty payloads do not have an MD5 value.
//
// Replaces the smithy-go middleware for httpChecksum trait.
func setMD5Checksum(req *http.Request) error {
	if req.Body == nil {
		return nil
	}

	if v := req.Header.Get(contentMD5Header); len(v) != 0 {
		return nil
	}
	if req.GetBody == nil {
		return fmt.Errorf("programming error: object storage operation must define a GetBody function in the request to set the MD5 Checksum")
	}

	body, err := req.GetBody()
	if err != nil {
		return err
	}

	defer body.Close()

	h := md5.New()

	_, err = io.Copy(h, body)

	if err != nil {
		return err
	}
	checksum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set(contentMD5Header, checksum)

	return nil
}

func getSignedHeaders(req *http.Request, ignoredHeaders map[string]struct{}) []string {
	signedHeaders := make([]string, 0, len(req.Header))

	for k := range req.Header {
		if _, ok := ignoredHeaders[k]; ok {
			continue
		}
		signedHeaders = append(signedHeaders, strings.ToLower(k))
	}

	slices.Sort(signedHeaders)
	return signedHeaders
}

func buildCanonicalHeaders(req *http.Request, signedHeaders []string) (canonicalHeaders string) {
	for _, k := range signedHeaders {
		v := req.Header.Values(k)

		line := fmt.Sprintf("%s:%s", strings.ToLower(k), strings.Join(v, ","))
		canonicalHeaders = fmt.Sprintf("%s%s\n", canonicalHeaders, line)
	}
	return
}

/*
Canonical string is composed by request elements that identifies that the signed string
relates to a specific request. All elements must be encoded in standard format, which
means:

- URI: escaped component from the provided URL.
- Headers: Trim out leading, trailing, and dedup inner spaces from signed header values.
- Query: must be sorted before hashing for consistency
*/
func buildCanonicalString(ctx *SignatureContext) string {
	return strings.Join([]string{
		ctx.HTTPMethod,
		ctx.CanonicalURI,
		ctx.CanonicalQuery,
		ctx.CanonicalHeaders,
		ctx.SignedHeaders,
		ctx.Parameters.PayloadHash,
	}, "\n")
}

/*
Builds the string to be signed with the derived key. The signed key will be embedded
into the "Authorization" header. The string to sign is a combination of the hashed
canonical request string plus extra information about the request, such as:

1. Signing Algorithm
2. Signing Time
3. Credentials Scope (region, service name, etc.)
*/
func buildStringToSign(ctx *SignatureContext) (string, error) {
	canonicalStr := buildCanonicalString(ctx)
	canonicalSHA, err := core.SHA256Hex(bytes.NewReader([]byte(canonicalStr)))
	if err != nil {
		return "", fmt.Errorf("failed to compute SHA from canonical str: %w", err)
	}
	return strings.Join([]string{
		ctx.Parameters.Algorithm,
		ctx.Parameters.Date,
		ctx.Parameters.Scope,
		canonicalSHA,
	}, "\n"), nil
}

/*
Derive a signing key by performing a succession of keyed hash operations
(HMAC operations) on the request date, Region, and service, with the secret access key
as the key for the initial hashing operation:

	Secret -> Date -> Region -deriveKey> Service -> Request Suffix
*/
func deriveKey(secretKey, shortTime string) []byte {
	hmacDate := core.HMACSHA256String([]byte(secretPrefix+secretKey), shortTime)
	hmacRegion := core.HMACSHA256String(hmacDate, signingRegion)
	hmacService := core.HMACSHA256String(hmacRegion, signingService)
	return core.HMACSHA256String(hmacService, requestSuffix)
}

/*
The Authorization header value carries more information than just the plain
signed strToSign. This function builds the header contents as specified with:

"Credential=": Your access key and the scope information, which includes the date,
Region, and service that were used to calculate the signature:

	<access-key>/<date>/<region>/<service>/<pre-defined-suffix>

Where:
- <date> value is specified using YYYYMMDD format.
- <aws-service> value is s3 when sending request to Amazon S3.

"Credential=": Your access key and the scope information, which includes the date,
Region, and service that were used to calculate the signature:

	<access-key>/<date>/<region>/<service>/<pre-defined-suffix>

"SignedHeaders=": A semicolon-separated list of request headers that you used to
compute Signature. The list includes header names only, and the header names must be in
lowercase:

	host;range;x-amz-date

"Signature=": The strToSign signed by the derived key, a 256-bit signature expressed as
64 lowercase hexadecimal characters:

	fe5f80f77d5fa3beca038a248ff027d0445342fe2855ddc963176630326f1024
*/

func buildAuthorizationHeader(ctx *SignatureContext) string {
	return fmt.Sprintf(
		"%s Credential=%s, SignedHeaders=%s, Signature=%s",
		signingAlgorithm,
		ctx.Parameters.Credential,
		ctx.SignedHeaders,
		ctx.Signature,
	)
}

func sign(ctx *SignatureContext, secretKey string) (err error) {
	strToSign, err := buildStringToSign(ctx)
	if err != nil {
		return
	}
	signKey := deriveKey(secretKey, ctx.Parameters.ShortDate)
	ctx.Signature = hex.EncodeToString(core.HMACSHA256String(signKey, strToSign))
	return
}

func signHeaders(req *http.Request, accessKey, secretKey string, unsignedPayload bool, ignoredHeaders map[string]struct{}) (err error) {
	payloadHash, err := setContentHeader(req, unsignedPayload)
	if err != nil {
		return
	}

	if req.Header.Get("Host") == "" {
		req.Header.Set("Host", req.Host)
	}

	if req.Body != nil && req.Header.Get("Content-Type") == "" {
		req.Header.Set("Content-Type", "application/octet-stream")
	}

	signingTime := time.Now().UTC()

	req.Header.Set(headerDateKey, signingTime.Format(longTimeFormat))

	if _, ok := ignoredHeaders[contentMD5Header]; !ok {
		if err := setMD5Checksum(req); err != nil {
			return fmt.Errorf("unable to compute checksum of the body content: %w", err)
		}
	}

	signedHeaders := getSignedHeaders(req, ignoredHeaders)
	params := NewSignatureParameters(accessKey, signingTime, payloadHash, signedHeaders)

	ctx := NewSignatureContext(params, req)
	if err = sign(ctx, secretKey); err != nil {
		return
	}

	authorization := buildAuthorizationHeader(ctx)
	req.Header.Set(authorizationHeaderKey, authorization)
	return nil
}

func SignedUrl(req *http.Request, accessKey, secretKey string, expirationTime time.Duration) (url *url.URL, err error) {
	params := NewSignatureParameters(accessKey, time.Now().UTC(), unsignedPayloadHeader, defaultSignedHeaders)

	if req.Header.Get("Host") == "" {
		req.Header.Set("Host", req.Host)
	}

	url = req.URL
	q := url.Query()
	q.Set("X-Amz-Expires", fmt.Sprintf("%d", int(expirationTime.Seconds())))
	q.Set("X-Amz-Algorithm", params.Algorithm)
	q.Set("X-Amz-Credential", params.Credential)
	q.Set("X-Amz-Date", params.Date)
	q.Set("X-Amz-SignedHeaders", strings.Join(params.SignedHeaders, ";"))
	url.RawQuery = q.Encode()

	ctx := NewSignatureContext(params, req)
	if err = sign(ctx, secretKey); err != nil {
		return
	}

	q.Set("X-Amz-Signature", ctx.Signature)

	url.RawQuery = q.Encode()
	return
}
