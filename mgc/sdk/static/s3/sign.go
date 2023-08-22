package s3

import (
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
func setContentHeader(req *http.Request) (payloadHash string, err error) {
	// TODO: handle signed payloads
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
	nr := req.Clone(req.Context())
	buf, err := io.ReadAll(nr.Body)
	if err != nil {
		return "", fmt.Errorf("could not read request body: %w", err)
	}

	return core.SHA256Hex(buf), nil
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

	if req.Header.Get("Content-Length") == "" {
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
	fmt.Println("headers: ", signedHeaders)
	fmt.Println("canHeaders: ", canonicalHeaders)
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

func sign(req *http.Request, accessKey, secretKey string, ignoredHeaders HeaderMap) error {
	signingTime := time.Now().UTC()
	payloadHash, err := setContentHeader(req)
	if err != nil {
		return err
	}

	// Set date header based on the custom key provided
	req.Header.Set(headerDateKey, signingTime.Format(longTimeFormat))

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

	// TODO: log and let the filters decide if shown or not
	LogSigningInfo(canonicalStr, strToSign)

	return nil
}
