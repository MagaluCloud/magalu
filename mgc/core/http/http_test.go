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

func TestConvertSliceJSONNumbers(t *testing.T) {
	t.Run("array of json.Number", func(t *testing.T) {
		// Caso 1: Array direto de json.Number
		input := map[string]interface{}{
			"example": map[string]interface{}{
				"allowed_values": []interface{}{
					json.Number("0"),
					json.Number("1"),
				},
			},
		}

		expected := map[string]interface{}{
			"example": map[string]interface{}{
				"allowed_values": []interface{}{
					int64(0),
					int64(1),
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

		if values, ok := input["example"].(map[string]interface{})["allowed_values"].([]interface{}); ok {
			if val0, ok := values[0].(int64); !ok || val0 != 0 {
				t.Errorf("Expected 'allowed_values[0]' to be int64(0), got %v of type %T", values[0], values[0])
			}
			if val1, ok := values[1].(int64); !ok || val1 != 1 {
				t.Errorf("Expected 'allowed_values[1]' to be int64(1), got %v of type %T", values[1], values[1])
			}
		} else {
			t.Error("Expected array structure not found in the result")
		}
	})

	t.Run("array with mixed types", func(t *testing.T) {
		// Caso 2: Array com tipos mistos
		input := map[string]interface{}{
			"mixed_array": []interface{}{
				json.Number("42"),
				json.Number("3.14"),
				"string",
				true,
				map[string]interface{}{
					"nested": json.Number("123"),
				},
			},
		}

		expected := map[string]interface{}{
			"mixed_array": []interface{}{
				int64(42),
				float64(3.14),
				"string",
				true,
				map[string]interface{}{
					"nested": int64(123),
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

		if arr, ok := input["mixed_array"].([]interface{}); ok {
			if val0, ok := arr[0].(int64); !ok || val0 != 42 {
				t.Errorf("Expected arr[0] to be int64(42), got %v of type %T", arr[0], arr[0])
			}
			if val1, ok := arr[1].(float64); !ok || val1 != 3.14 {
				t.Errorf("Expected arr[1] to be float64(3.14), got %v of type %T", arr[1], arr[1])
			}
			if val2, ok := arr[2].(string); !ok || val2 != "string" {
				t.Errorf("Expected arr[2] to be string('string'), got %v of type %T", arr[2], arr[2])
			}
			if val3, ok := arr[3].(bool); !ok || val3 != true {
				t.Errorf("Expected arr[3] to be bool(true), got %v of type %T", arr[3], arr[3])
			}
			if nestedMap, ok := arr[4].(map[string]interface{}); ok {
				if nested, ok := nestedMap["nested"].(int64); !ok || nested != 123 {
					t.Errorf("Expected arr[4]['nested'] to be int64(123), got %v of type %T", nestedMap["nested"], nestedMap["nested"])
				}
			} else {
				t.Errorf("Expected arr[4] to be a map, got %T", arr[4])
			}
		} else {
			t.Error("Expected array structure not found in the result")
		}
	})

	t.Run("array inside array", func(t *testing.T) {
		// Caso 3: Array dentro de array
		input := map[string]interface{}{
			"nested_arrays": []interface{}{
				[]interface{}{
					json.Number("1"),
					json.Number("2"),
				},
				[]interface{}{
					json.Number("3"),
					json.Number("4"),
				},
			},
		}

		expected := map[string]interface{}{
			"nested_arrays": []interface{}{
				[]interface{}{
					int64(1),
					int64(2),
				},
				[]interface{}{
					int64(3),
					int64(4),
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

		if outerArr, ok := input["nested_arrays"].([]interface{}); ok {
			if innerArr1, ok := outerArr[0].([]interface{}); ok {
				if val0, ok := innerArr1[0].(int64); !ok || val0 != 1 {
					t.Errorf("Expected innerArr1[0] to be int64(1), got %v of type %T", innerArr1[0], innerArr1[0])
				}
				if val1, ok := innerArr1[1].(int64); !ok || val1 != 2 {
					t.Errorf("Expected innerArr1[1] to be int64(2), got %v of type %T", innerArr1[1], innerArr1[1])
				}
			} else {
				t.Errorf("Expected outerArr[0] to be an array, got %T", outerArr[0])
			}

			if innerArr2, ok := outerArr[1].([]interface{}); ok {
				if val0, ok := innerArr2[0].(int64); !ok || val0 != 3 {
					t.Errorf("Expected innerArr2[0] to be int64(3), got %v of type %T", innerArr2[0], innerArr2[0])
				}
				if val1, ok := innerArr2[1].(int64); !ok || val1 != 4 {
					t.Errorf("Expected innerArr2[1] to be int64(4), got %v of type %T", innerArr2[1], innerArr2[1])
				}
			} else {
				t.Errorf("Expected outerArr[1] to be an array, got %T", outerArr[1])
			}
		} else {
			t.Error("Expected nested array structure not found in the result")
		}
	})
}
