package markdown_report

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/pterm/pterm"
)

func CheckURL(urlString string, errorChan chan ProgressError) (string, error) {
	u, urlErr := url.Parse(urlString)
	if urlErr == nil && strings.HasPrefix(u.Scheme, "http") {
		// download the file
		resp, httpErr := http.Get(urlString)
		if httpErr != nil {
			errorChan <- ProgressError{
				Job:     "download",
				Message: fmt.Sprintf("error downloading file '%s': %s", urlString, httpErr.Error()),
			}
			return urlString, httpErr
		}
		bits, _ := io.ReadAll(resp.Body)

		if len(bits) <= 0 {
			errorChan <- ProgressError{
				Job:     "download",
				Message: fmt.Sprintf("downloaded file '%s' is empty", urlString),
			}
			return urlString, fmt.Errorf("downloaded file '%s' is empty", urlString)
		}
		tmpFile, _ := os.CreateTemp("", "left.yaml")
		_, wErr := tmpFile.Write(bits)
		if wErr != nil {
			errorChan <- ProgressError{
				Job:     "download",
				Message: fmt.Sprintf("downloaded file '%s' cannot be written: %s", urlString, wErr.Error()),
			}
			return urlString, fmt.Errorf("downloaded file '%s' is empty", urlString)
		}
		return tmpFile.Name(), nil
	}
	return urlString, nil
}

func WriteReportFile(reportFile string, report []byte) error {
	err := os.WriteFile(reportFile, report, 0744)
	if err != nil {
		pterm.Error.Println(err.Error())
		return err
	}
	pterm.Success.Printf("report written to file '%s' (%dkb)", reportFile, len(report)/1024)
	pterm.Println()
	pterm.Println()
	return nil
}
