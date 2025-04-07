package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"reflect"
	"testing"
)

type dummyTransport struct{}

func (o dummyTransport) RoundTrip(*http.Request) (*http.Response, error) {
	return &http.Response{}, nil
}

func TestClientCreation(t *testing.T) {
	client := NewClient(dummyTransport{})
	if client == nil {
		t.Fail()
	}
}

func TestContext(t *testing.T) {
	ctx := context.Background()
	if ClientFromContext(ctx) != nil {
		t.Error("corehttp.ClientFromContext() should not return a client from an empty context")
	}
	client := NewClient(dummyTransport{})
	ctx = NewClientContext(ctx, client)
	if ClientFromContext(ctx) == nil {
		t.Error("corehttp.ClientFromContext() failed to retrieve client from valid context")
	}
}

type dummyResponseBodyStruct struct {
	Data string `json:"data"`
}

func TestDecodeJSON(t *testing.T) {
	expectedData := "some string"
	dummyResponse := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(fmt.Sprintf("{\"data\": \"%s\"}", expectedData))),
	}
	decoded := new(dummyResponseBodyStruct)
	err := DecodeJSON(dummyResponse, &decoded)
	if err != nil {
		t.Errorf("DecodeJSON function failed: %s", err)
	}
	if decoded.Data != "some string" {
		t.Errorf("DecodeJSON function failed. 'dummyResponseBodyStruct.Data' expected %s but got %s", expectedData, decoded.Data)
	}
}

type testTypesDecodeJson struct {
	Name     string  `json:"name"`
	Latitude string  `json:"latitude"`
	CPUCount float64 `json:"cpu_count"`
	RAM      int     `json:"ram"`
	Tops     []struct {
		Read  int `json:"read"`
		Write int `json:"write"`
	} `json:"tops"`
}

func TestDecodeJSONComplex(t *testing.T) {
	expectedData := `{
					  "name": "play",
					  "latitude": "2",
					  "cpu_count": 1.32,
					  "ram": 16,
					  "tops": [
					    {
					      "read": 1000,
					      "write": 1000
					    }
					  ]
					}`
	dummyResponse := &http.Response{
		Body: io.NopCloser(bytes.NewBufferString(expectedData)),
	}
	decoded := new(testTypesDecodeJson)
	err := DecodeJSON(dummyResponse, &decoded)
	if err != nil {
		t.Errorf("DecodeJSON function failed: %s", err)
	}

	latitudeReturn := string("2")
	if decoded.Latitude != latitudeReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %s but got %s", decoded.Latitude, latitudeReturn)
	}

	nameReturn := string("play")
	if decoded.Name != nameReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %s but got %s", decoded.Name, nameReturn)
	}

	cpuReturn := float64(1.32)
	if decoded.CPUCount != cpuReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %v but got %v", decoded.CPUCount, cpuReturn)
	}

	ramReturn := int(16)
	if decoded.RAM != ramReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %v but got %v", decoded.RAM, ramReturn)
	}

	readReturn := int(1000)
	if decoded.Tops[0].Read != readReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %v but got %v", decoded.Tops[0].Read, readReturn)
	}

	writeReturn := int(1000)
	if decoded.Tops[0].Write != writeReturn {
		t.Errorf("DecodeJSON function failed. 'testTypesDecodeJson' expected %v but got %v", decoded.Tops[0].Write, writeReturn)
	}

}

func TestNewHttpErrorFromResponse(t *testing.T) {
	dummyResponse := &http.Response{
		Body:       io.NopCloser(bytes.NewBufferString("some value")),
		StatusCode: 123,
		Status:     "not ok",
		Header:     http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"1234"}},
		Request: &http.Request{Header: http.Header{"X-Request-Id": []string{"1234"}},
			Response: &http.Response{
				Header: http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"1234"}},
			}},
	}
	dummyRequest := &http.Request{
		Header: http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"1234"}},
	}
	httpErr := NewHttpErrorFromResponse(dummyResponse, dummyRequest)

	expectedHttpErrr := &HttpError{
		Code:    123,
		Status:  "not ok",
		Headers: http.Header{"Content-Type": []string{"application/json"}, "X-Request-Id": []string{"1234"}},
		Payload: bytes.NewBufferString("some value").Bytes(),
		Message: "not ok",
		Slug:    "unknown",
	}

	expected := &IdentifiableHttpError{
		HttpError: expectedHttpErrr,
		RequestID: "1234",
	}
	if !reflect.DeepEqual(httpErr, expected) {
		t.Errorf("NewHttpErrorFromResponse returned %+v, but expected %+v", *httpErr, *expected)
	}

	dummyResponse.Body = io.NopCloser(bytes.NewBufferString("{\"slug\": \"the slug\",\"message\": \"the message\"}"))
	expected.Message = "the message"
	expected.Slug = "the slug"
	expected.Payload = bytes.NewBufferString("{\"slug\": \"the slug\",\"message\": \"the message\"}").Bytes()

	httpErr = NewHttpErrorFromResponse(dummyResponse, dummyRequest)
	if !reflect.DeepEqual(httpErr, expected) {
		t.Errorf("NewHttpErrorFromResponse failed to decode response's 'data' and 'message' fields properly\nInput: %+v\nOutput: %+v\nExpected: %+v", *dummyResponse, *httpErr, *expected)
	}
}

func TestUnwrapResponse(t *testing.T) {
	t.Run("non-2xx status code", func(t *testing.T) {
		for i := 100; i < 600; i++ {
			if i >= 200 && i < 300 {
				continue
			}

			resp := &http.Response{StatusCode: i, Body: io.NopCloser(bytes.NewBufferString(""))}
			req := &http.Request{}
			_, err := UnwrapResponse[any](resp, req)
			httpErr, ok := err.(*IdentifiableHttpError)
			if !ok {
				t.Fatalf("expected IdentifiableHttpError when status code is %v, but was unable to convert %#v to *HttpError", i, err)
				return
			}

			expectedErr := NewHttpErrorFromResponse(resp, req)
			if !reflect.DeepEqual(httpErr, expectedErr) {
				t.Fatalf("expected err == %#v when status code is %v, got %#v instead", expectedErr, i, err)
			}
		}
	})

	t.Run("empty body status code", func(t *testing.T) {
		resp := &http.Response{StatusCode: 204}
		req := &http.Request{}
		var expectedStr string
		resultStr, err := UnwrapResponse[string](resp, req)
		if err != nil || resultStr != expectedStr {
			t.Fatalf("expected err == nil and zero value return, got instead err == '%v' and result '%v'", err, resultStr)
		}

		var expectedAny any
		resultAny, err := UnwrapResponse[any](resp, req)
		if err != nil || resultAny != expectedAny {
			t.Fatalf("expected err == nil and zero value return, got instead err == '%v' and result '%v'", err, resultAny)
		}

		var expectedInt int
		resultInt, err := UnwrapResponse[int](resp, req)
		if err != nil || resultInt != expectedInt {
			t.Fatalf("expected err == nil and zero value return, got instead err == '%v' and result '%v'", err, resultInt)
		}

		var expectedBool bool
		resultBool, err := UnwrapResponse[bool](resp, req)
		if err != nil || resultBool != expectedBool {
			t.Fatalf("expected err == nil and zero value return, got instead err == '%v' and result '%v'", err, resultBool)
		}
	})

	t.Run("multipart response", func(t *testing.T) {
		header := http.Header{}
		header.Set("Content-Type", `multipart/form-data; boundary="XXX"`)
		bodyText := `--XXX
Content-Disposition: form-data; name="foo"

dummy text
--XXX
Content-Disposition: form-data; name="bar"

more dummy text
`
		resp := &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       io.NopCloser(bytes.NewBufferString(bodyText)),
		}
		req := &http.Request{}

		part, err := UnwrapResponse[*multipart.Part](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping multipart response to *multipart.Part: %v", err)
		}

		bytesRead, err := io.ReadAll(part)
		if err != nil {
			t.Fatalf("error when reading multipart part: %v", err)
		}

		expectedStrRead := "dummy text"
		if strRead := string(bytesRead[:]); strRead != expectedStrRead {
			t.Fatalf("multipart part expected '%v' but got %v instead", expectedStrRead, err)
		}

		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[any](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping multipart response to any: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[int](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or '*multipart.Part', got nil instead for int")
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[string](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or '*multipart.Part', got nil instead for string")
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[bool](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or '*multipart.Part', got nil instead for bool")
		}
		type dummyStruct struct{}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[dummyStruct](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or '*multipart.Part', got nil instead for dummyStruct")
		}
	})

	t.Run("json response", func(t *testing.T) {
		header := http.Header{}
		header.Set("Content-Type", "application/json")
		bodyText := `{"str": "strValue"}`
		resp := &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       io.NopCloser(bytes.NewBufferString(bodyText)),
		}

		req := &http.Request{}

		type dummyRespStruct struct {
			Str string `json:"str"`
		}

		result, err := UnwrapResponse[dummyRespStruct](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping json response to dummy struct: %v", err)
		}

		if result.Str != "strValue" {
			t.Fatalf("expected result struct to have 'strValue' in 'str' field, got '%s' instead", result.Str)
		}

		type invalidDummyRespStruct struct {
			Field string
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[invalidDummyRespStruct](resp, req)
		if err == nil {
			t.Fatalf("unwrapping response with text '%s' to invalid struct should fail, error was %v instead", bodyText, err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		anyResult, err := UnwrapResponse[any](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping json response to any: %v", err)
		}
		if _, ok := anyResult.(map[string]any); !ok {
			t.Fatalf("decoding to any with body text '%s' should result in a map[string]any, got %T instead", bodyText, anyResult)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[int](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got %v instead for int", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[string](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got %v instead for string", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[bool](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got %v instead for bool", err)
		}
	})

	t.Run("xml response", func(t *testing.T) {
		header := http.Header{}
		header.Set("Content-Type", "application/xml")
		bodyText := `<dummyRespStruct><str>strValue</str></dummyRespStruct>`
		resp := &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       io.NopCloser(bytes.NewBufferString(bodyText)),
		}
		req := &http.Request{}
		type dummyRespStruct struct {
			Str string `xml:"str"`
		}

		result, err := UnwrapResponse[dummyRespStruct](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping xml response to dummy struct: %v", err)
		}

		if result.Str != "strValue" {
			t.Fatalf("expected result struct to have 'strValue' in 'str' field, got '%s' instead", result.Str)
		}

		// type invalidDummyRespStruct struct {
		// 	Field string
		// }
		// resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		// _, err = UnwrapResponse[invalidDummyRespStruct](resp)
		// if err == nil {
		// 	t.Fatalf("unwrapping response with text '%s' to invalid struct should fail, error was %v instead", bodyText, err)
		// }
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[any](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping xml response to any: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[string](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping xml response to string: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[int](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any', a decodable struct or a slice got nil instead for int")
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[bool](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any', a decodable struct or a slice, got nil instead for bool")
		}
	})

	t.Run("default body", func(t *testing.T) {
		header := http.Header{}
		header.Set("Content-Type", "text/html")
		bodyText := `<root><str>strValue</str></root>`
		resp := &http.Response{
			StatusCode: 200,
			Header:     header,
			Body:       io.NopCloser(bytes.NewBufferString(bodyText)),
		}
		req := &http.Request{}

		result, err := UnwrapResponse[io.ReadCloser](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping body as ReadCloser: %v", err)
		}

		bytesRead, err := io.ReadAll(result)
		if err != nil {
			t.Fatalf("error when reading result body ReadCloser: %v", err)
		}

		strRead := string(bytesRead[:])
		if strRead != bodyText {
			t.Fatalf("result body ReadCloser doesn't match body content. Expected '%s', but got '%s'", bodyText, strRead)
		}

		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[any](resp, req)
		if err != nil {
			t.Fatalf("error when unwrapping default response to any: %v", err)
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[int](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got nil instead for int")
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[string](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got nil instead for string")
		}
		resp.Body = io.NopCloser(bytes.NewBufferString(bodyText))
		_, err = UnwrapResponse[bool](resp, req)
		if err == nil {
			t.Fatalf("should return error when T is not 'any' or a decodable struct, got nil instead for bool")
		}
	})
}

func TestConvertComplexJSONNumbers(t *testing.T) {
	input := map[string]interface{}{
		"types": map[string]interface{}{
			"name":      "play",
			"latitude":  "2",
			"cpu_count": 1.32,
			"tops": map[string]interface{}{
				"read":  json.Number("1000"),
				"write": json.Number("1000"),
			},
		},
	}

	expected := map[string]interface{}{
		"types": map[string]interface{}{
			"name":      "play",
			"latitude":  "2",
			"cpu_count": 1.32,
			"tops": map[string]interface{}{
				"read":  int64(1000),
				"write": int64(1000),
			},
		},
	}

	inputValue := reflect.ValueOf(input)
	err := convertJSONNumbers(inputValue)
	if err != nil {
		t.Fatalf("convertJSONNumbers() error = %v", err)
	}

	if !reflect.DeepEqual(input, expected) {
		t.Errorf("convertJSONNumbers() = %v, want %v", input, expected)
	}

	if tops, ok := input["types"].(map[string]interface{})["tops"].(map[string]interface{}); ok {
		if read, ok := tops["read"].(int64); !ok || read != 1000 {
			t.Errorf("Expected 'read' to be int64(1000), got %v", tops["read"])
		}
		if write, ok := tops["write"].(int64); !ok || write != 1000 {
			t.Errorf("Expected 'write' to be int64(1000), got %v", tops["write"])
		}
	} else {
		t.Error("Expected structure not found in the result")
	}
}
