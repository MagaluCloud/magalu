package core

import (
	"bytes"
	"io"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/magiconair/properties/assert"
)

func TestCheckVersionWithOutdatedVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getHttp := func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 302,
			Header:     map[string][]string{"Location": {"https://github.com/MagaluCloud/mgccli/releases/2.0.0"}},
		}, nil
	}

	config := func(key string, value interface{}) error {
		return nil
	}

	vc := NewVersionChecker(getHttp, config, config)
	vc.CheckVersion("1.0.0")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	out := <-outC

	assert.Equal(
		t, out, "⚠️ You are using an outdated version of mgc cli. "+
			"Please update to the latest version: 2.0.0 \n\n\n",
	)
}

func TestCheckVersionWithLatestVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getHttp := func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 302,
			Header:     map[string][]string{"Location": {"https://github.com/MagaluCloud/mgccli/releases/2.0.0"}},
		}, nil
	}

	config := func(key string, value interface{}) error {
		return nil
	}

	vc := NewVersionChecker(getHttp, config, config)
	vc.CheckVersion("2.0.0")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	out := <-outC

	assert.Equal(t, out, "")
}

func TestCheckVersionWithInvalidCurrentVersion(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getHttp := func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 302,
			Header:     map[string][]string{"Location": {"https://github.com/MagaluCloud/mgccli/releases/2.0.0"}},
		}, nil
	}

	config := func(key string, value interface{}) error {
		return nil
	}
	vc := NewVersionChecker(getHttp, config, config)
	vc.CheckVersion("invalid")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	out := <-outC

	assert.Equal(t, out, "Invalid current version: Invalid Semantic Version\n")
}

func TestCheckVersionWithUpdateIntervalNotExceeded(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getHttp := func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 302,
			Header:     map[string][]string{"Location": {"https://github.com/MagaluCloud/mgccli/releases/2.0.0"}},
		}, nil
	}

	getConfig := func(key string, value interface{}) error {
		if key == versionLastCheck {
			v := value.(*string)
			*v = time.Now().Format(time.RFC3339)
		}
		return nil
	}

	vc := NewVersionChecker(getHttp, getConfig, nil)
	vc.CheckVersion("2.0.0")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	out := <-outC

	assert.Equal(t, out, "")
}

func TestCheckVersionWithUpdateIntervalExceeded(t *testing.T) {
	old := os.Stdout
	r, w, _ := os.Pipe()
	os.Stdout = w

	getHttp := func(url string) (*http.Response, error) {
		return &http.Response{
			StatusCode: 302,
			Header:     map[string][]string{"Location": {"https://github.com/MagaluCloud/mgccli/releases/2.0.0"}},
		}, nil
	}

	getConfig := func(key string, value interface{}) error {
		if key == versionLastCheck {
			v := value.(*string)
			*v = time.Now().Add(-3 * time.Hour).Format(time.RFC3339)
		}
		return nil
	}
	setConfig := func(key string, value interface{}) error {
		return nil
	}

	vc := NewVersionChecker(getHttp, getConfig, setConfig)
	vc.CheckVersion("1.0.0")

	outC := make(chan string)
	go func() {
		var buf bytes.Buffer
		_, _ = io.Copy(&buf, r)
		outC <- buf.String()
	}()

	_ = w.Close()
	os.Stdout = old
	out := <-outC

	assert.Equal(
		t, out, "⚠️ You are using an outdated version of mgc cli. "+
			"Please update to the latest version: 2.0.0 \n\n\n",
	)
}
