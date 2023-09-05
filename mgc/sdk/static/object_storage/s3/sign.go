package s3

import (
	"crypto/md5"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"strings"
	"time"

	"magalu.cloud/core"
)

type HeaderMap = map[string]any

// TODO: refactor into a round tripper
func setContentHeader(req *http.Request, unsigned bool) (payloadHash string, err error) {
	if unsigned {
		req.Header.Set(contentSHAKey, unsignedPayloadHeader)
		return unsignedPayloadHeader, nil
	}
	payloadHash, err = getPayloadHash(req)
	if err != nil {
		return "", err
	}
	req.Header.Set(contentSHAKey, payloadHash)
	return payloadHash, nil
}

// CopySeekableBody copies the seekable body to an io.Writer
// Copied from https://github.com/aws/aws-sdk-go/blob/main/aws/types.go#L244
func CopySeekableBody(dst io.Writer, src io.ReadSeeker) (int64, error) {
	curPos, err := src.Seek(0, io.SeekCurrent)
	if err != nil {
		return 0, err
	}

	// copy errors may be assumed to be from the body.
	n, err := io.Copy(dst, src)
	if err != nil {
		return n, err
	}

	// seek back to the first position after reading to reset
	// the body for transmission.
	_, err = src.Seek(curPos, io.SeekStart)
	if err != nil {
		return n, err
	}

	return n, nil
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
	nr := req.Clone(req.Context())
	buf, err := io.ReadAll(nr.Body)
	if err != nil {
		return "", fmt.Errorf("could not read request body: %w", err)
	}

	return core.SHA256Hex(buf), nil
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

	h := md5.New()
	switch req.Body.(type) {
	case io.ReadSeeker:
		length, err := CopySeekableBody(h, req.Body.(io.ReadSeeker))
		if err != nil {
			return err
		}
		req.ContentLength = length
	default:
		// Creates a deep copy of body that is safe to use
		bodyReader, err := req.GetBody()
		if err != nil {
			return err
		}

		defer bodyReader.Close()
		buf, err := io.ReadAll(bodyReader)
		if err != nil {
			return err
		}
		h.Write(buf)
	}
	checksum := base64.StdEncoding.EncodeToString(h.Sum(nil))
	req.Header.Set(contentMD5Header, checksum)
	return nil
}

// buildCredentialScope builds the Signature Version 4 (SigV4) signing scope
func buildCredentialScope(shortTime string) string {
	return strings.Join([]string{
		shortTime,
		signingRegion,
		signingService,
		requestSuffix,
	}, "/")
}

func buildCredentialStr(parts ...string) string {
	return strings.Join(parts, "/")
}

func buildCanonicalHeaders(req *http.Request, ignoredHeaders HeaderMap) (string, string) {
	signedHeaders := []string{}
	canonicalHeaders := ""

	if req.Header.Get("Host") == "" {
		req.Header.Set("Host", req.Host)
	}

	if req.Body != nil && req.Header.Get("Content-Length") == "" {
		req.Header.Set("Content-Length", strconv.FormatInt(req.ContentLength, 10))
	}

	sortedHeaderKeys := make([]string, 0, len(req.Header))
	for k := range req.Header {
		sortedHeaderKeys = append(sortedHeaderKeys, k)
	}
	sort.Strings(sortedHeaderKeys)
	for _, k := range sortedHeaderKeys {
		v := req.Header.Values(k)
		if _, ok := ignoredHeaders[k]; ok {
			continue // ignored header
		}

		line := fmt.Sprintf("%s:%s", strings.ToLower(k), strings.Join(v, ","))
		signedHeaders = append(signedHeaders, strings.ToLower(k))
		canonicalHeaders = fmt.Sprintf("%s%s\n", canonicalHeaders, line)
	}
	return strings.Join(signedHeaders, ";"), canonicalHeaders
}

/*
Canonical string is composed by request elements that identifies that the signed string
relates to a specific request. All elements must be encoded in standard format, which
means:

- URI: escaped component from the provided URL.
- Headers: Trim out leading, trailing, and dedup inner spaces from signed header values.
- Query: must be sorted before hashing for consistency
*/
func buildCanonicalString(method, uri, query, signedHeaders, canonicalHeaders, payloadHash string) string {
	return strings.Join([]string{
		method,
		uri,
		query,
		canonicalHeaders,
		signedHeaders,
		payloadHash,
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
func buildStringToSign(credScope, canonicalStr, longTime string) string {
	return strings.Join([]string{
		signingAlgorithm,
		longTime,
		credScope,
		core.SHA256Hex([]byte(canonicalStr)),
	}, "\n")
}

/*
Derive a signing key by performing a succession of keyed hash operations
(HMAC operations) on the request date, Region, and service, with the secret access key
as the key for the initial hashing operation:

	Secret -> Date -> Region -deriveKey> Service -> Request Suffix
*/
func deriveKey(prefix, secretKey, shortTime string) []byte {
	hmacDate := core.HMACSHA256String([]byte(prefix+secretKey), shortTime)
	hmacRegion := core.HMACSHA256String(hmacDate, signingRegion)
	hmacService := core.HMACSHA256String(hmacRegion, signingService)
	return core.HMACSHA256String(hmacService, "aws4_request")
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
func buildAuthorizationHeader(credentialStr, signedHeadersStr, signingSignature string) string {
	const credential = "Credential="
	const signedHeaders = "SignedHeaders="
	const signature = "Signature="
	return fmt.Sprintf("%s %s%s, %s%s, %s%s", signingAlgorithm, credential, credentialStr, signedHeaders, signedHeadersStr, signature, signingSignature)
}

func sign(req *http.Request, accessKey, secretKey string, unsignedPayload bool, ignoredHeaders HeaderMap) error {
	signingTime := time.Now().UTC()
	payloadHash, err := setContentHeader(req, unsignedPayload)
	if err != nil {
		return err
	}

	// Set date header based on the custom key provided
	req.Header.Set(headerDateKey, signingTime.Format(longTimeFormat))

	if err := setMD5Checksum(req); err != nil {
		return fmt.Errorf("Unable to compute checksum of the body content: %w", err)
	}

	// Sort Each Query Key's Values
	query := req.URL.Query()
	for key := range query {
		sort.Strings(query[key])
	}

	credScope := buildCredentialScope(signingTime.Format(shortTimeFormat))
	credStr := buildCredentialStr(accessKey, credScope)

	signedHeadersStr, canonicalHeaderStr := buildCanonicalHeaders(req, ignoredHeaders)
	canonicalStr := buildCanonicalString(
		req.Method,
		req.URL.EscapedPath(),
		url.QueryEscape(req.URL.RawQuery),
		signedHeadersStr,
		canonicalHeaderStr,
		payloadHash,
	)
	strToSign := buildStringToSign(credScope, canonicalStr, signingTime.Format(longTimeFormat))
	signKey := deriveKey(secretPrefix, secretKey, signingTime.Format(shortTimeFormat))
	signature := hex.EncodeToString(core.HMACSHA256String(signKey, strToSign))

	signedAuthorization := buildAuthorizationHeader(credStr, signedHeadersStr, signature)
	req.Header.Set(authorizationHeaderKey, signedAuthorization)

	return nil
}
