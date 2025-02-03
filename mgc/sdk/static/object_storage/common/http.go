package common

import (
	"bytes"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strings"

	mgcHttpPkg "github.com/MagaluCloud/magalu/mgc/core/http"
	"github.com/MagaluCloud/magalu/mgc/core/utils"
	"github.com/MagaluCloud/magalu/mgc/core/xml"
)

type XMLError struct {
	Message string `xml:"Message"`
	Code    string `xml:"Code"`
}

func UnwrapResponse[T any](resp *http.Response, req *http.Request) (result T, err error) {
	if err = ExtractErr(resp, req); err != nil {
		return
	}

	if resp.StatusCode == 204 {
		return
	}

	contentLength := resp.Header.Get("Content-Length")
	if contentLength == "0" {
		return
	}

	contentType, params, _ := mime.ParseMediaType(resp.Header.Get("Content-Type"))

	switch {
	default:
		err = utils.AssignToT(&result, resp.Body)
		return
	case strings.HasPrefix(contentType, "multipart/"):
		body, bodyErr := mgcHttpPkg.BodyReaderSafe(resp)
		if bodyErr != nil {
			err = fmt.Errorf("error when reading response body: %w", bodyErr)
			return
		}
		// TODO: do we have multi-part downloads? or just single?
		// If multi, then we need to return a multipart reader...
		// return multipart.NewReader(resp.Body, params["boundary"]), nil
		r := multipart.NewReader(body, params["boundary"])
		nextPart, npErr := r.NextPart()
		err = npErr
		if err != nil {
			return
		}
		err = utils.AssignToT(&result, nextPart)
		return
	case contentType == "application/json":
		err = mgcHttpPkg.DecodeJSON(resp, &result)
	// TODO: Don't assume that empty or text/plain or text/html Content-Type is xml. We currently do this because the server
	// has some endpoints that don't return wrong Content-Types (or none at all), but when those are fixed we should
	// remove this check for `""` and `text/plain`
	case contentType == "application/xml", contentType == "text/plain", contentType == "text/html", contentType == "":
		err = mgcHttpPkg.DecodeXML(resp, &result)
	}

	return
}

func ExtractErr(resp *http.Response, req *http.Request) error {
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return newIdentifiableHttpError(resp, req)
	}
	return nil
}

func newIdentifiableHttpError(resp *http.Response, req *http.Request) *mgcHttpPkg.IdentifiableHttpError {
	slug := "unknown"
	message := resp.Status

	defer resp.Body.Close()

	payload, _ := io.ReadAll(resp.Body)
	contentType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		logger().Debugw("ignored invalid response", "Content-Type", resp.Header.Get("Content-Type"), "error", err.Error())
	}

	if contentType == "application/xml" {
		data := XMLError{}
		decoder := xml.NewDecoder(bytes.NewBuffer(payload))

		if err := decoder.Decode(&data); err == nil {
			// fmt.Printf("%#v\n", data)
			if data.Message != "" {
				message = data.Message
			}
			if data.Code != "" {
				slug = data.Code
			}
		}
	}

	httpError := &mgcHttpPkg.HttpError{
		Code:    resp.StatusCode,
		Status:  resp.Status,
		Headers: resp.Header,
		Payload: payload,
		Message: message,
		Slug:    slug,
	}

	return mgcHttpPkg.NewIdentifiableHttpError(httpError, req, resp)

}
