package openapi

import (
	"fmt"
	"net/http"
	"strings"
	"text/scanner"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/go-openapi/jsonpointer"
	"magalu.cloud/core"
	mgcHttpPkg "magalu.cloud/core/http"
)

func getRemainder(s *scanner.Scanner) string {
	var result string
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		result += s.TokenText()
	}
	return result
}

type linkRtExpResolver interface {
	resolve() (value any, found bool, err error)
}

type linkRtExpression struct {
	str                string
	httpResult         mgcHttpPkg.HttpResult
	findParameterValue func(location, name string) (any, bool)
}

var _ linkRtExpResolver = (*linkRtExpression)(nil)

type linkRtExpSource struct {
	s                  *scanner.Scanner
	findParameterValue func(location, name string) (any, bool)
	header             http.Header
	body               core.Value
}

var _ linkRtExpResolver = (*linkRtExpSource)(nil)

type linkRtExpHeader struct {
	s      *scanner.Scanner
	header http.Header
}

var _ linkRtExpResolver = (*linkRtExpHeader)(nil)

type linkRtExpQuery struct {
	s                  *scanner.Scanner
	filter             string
	findParameterValue func(location, name string) (any, bool)
}

var _ linkRtExpResolver = (*linkRtExpQuery)(nil)

type linkRtExpPath = linkRtExpQuery

type linkRtExpBody struct {
	s    *scanner.Scanner
	data core.Value
}

var _ linkRtExpResolver = (*linkRtExpBody)(nil)

func (o *linkRtExpHeader) resolveChild() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case ".":
			break Loop
		default:
			result := o.header.Get(text)
			if result == "" {
				return nil, false, nil
			}
			return result, true, nil
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpHeader) resolve() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case ".":
			return o.resolveChild()
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpBody) resolveChild() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case "/":
			jpStr := text + getRemainder(o.s)
			jpHandler, err := jsonpointer.New(jpStr)
			if err != nil {
				return nil, false, fmt.Errorf("malformed json pointer on link runtime expression: %s", jpStr)
			}

			resolved, _, err := jpHandler.Get(o.data)
			if err != nil {
				return nil, false, fmt.Errorf("unable to resolve JSON Pointer '%s' on data %#v", jpStr, o.data)
			}
			return resolved, true, nil
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpBody) resolve() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case "#":
			return o.resolveChild()
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpQuery) resolveChild() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case ".":
			break Loop
		default:
			resolved, ok := o.findParameterValue(o.filter, text)
			if !ok {
				return nil, false, nil
			}

			return resolved, true, nil
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpQuery) resolve() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case ".":
			return o.resolveChild()
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpSource) resolveChild() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case "header":
			return (&linkRtExpHeader{o.s, o.header}).resolve()
		case "query":
			return (&linkRtExpQuery{o.s, openapi3.ParameterInQuery, o.findParameterValue}).resolve()
		case "path":
			return (&linkRtExpPath{o.s, openapi3.ParameterInPath, o.findParameterValue}).resolve()
		case "body":
			return (&linkRtExpBody{o.s, o.body}).resolve()
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpSource) resolve() (any, bool, error) {
Loop:
	for tok := o.s.Scan(); tok != scanner.EOF; tok = o.s.Scan() {
		text := o.s.TokenText()
		switch text {
		case ".":
			return o.resolveChild()
		default:
			break Loop
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpression) resolveChild(s *scanner.Scanner) (any, bool, error) {
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		text := s.TokenText()
		switch text {
		case "url":
			return o.httpResult.Request().URL.String(), true, nil
		case "method":
			return o.httpResult.Request().Method, true, nil
		case "statusCode":
			return o.httpResult.Response().StatusCode, true, nil
		case "request":
			return (&linkRtExpSource{s, o.findParameterValue, o.httpResult.Request().Header, o.httpResult.RequestBody()}).resolve()
		case "response":
			return (&linkRtExpSource{s, o.findParameterValue, o.httpResult.Response().Header, o.httpResult.ResponseBody()}).resolve()
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}

func (o *linkRtExpression) resolve() (any, bool, error) {
	s := &scanner.Scanner{}
	s.Init(strings.NewReader(o.str))

Loop:
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		text := s.TokenText()
		switch text {
		case "$":
			return o.resolveChild(s)
		default:
			// Treat entire string as raw value

			// TODO: Use regex to detect pattern
			if strings.Contains(o.str, "{") && strings.Contains(o.str, "}") {
				break Loop
			}

			return o.str, true, nil
		}
	}
	return nil, false, fmt.Errorf("malformed link runtime expression")
}
