package spec

import (
	"crypto/md5"
	"encoding/hex"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	markdownReport "github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/spec/markdown-report"
	"github.com/MagaluCloud/magalu/mgc/spec_manipulator/cmd/tui"
	"github.com/google/uuid"
	"github.com/pterm/pterm"

	// oasChanges "github.com/pb33f/openapi-changes"
	"github.com/spf13/cobra"
)

func diffCheckerCmd() *cobra.Command {
	var dir string
	var menu string

	cmd := &cobra.Command{
		Use:   "diff [dir] [menu]",
		Short: "Download available spec",
		Run: func(cmd *cobra.Command, args []string) {

			_ = verificarEAtualizarDiretorio(dir)

			var currentConfig []specList
			var err error

			if menu != "" {
				currentConfig, err = loadList(menu)
			} else {
				currentConfig, err = getConfigToRun()
			}
			if err != nil {
				return
			}
			spinner := tui.NewSpinner()
			spinner.Start("Downloading ...")
			for _, v := range currentConfig {
				spinner.UpdateText("Downloading " + v.File)

				dirTmp := filepath.Join(dir, "tmp")
				os.MkdirAll(dirTmp, 0755)

				tmpFile := filepath.Join(dirTmp, v.File)

				if !strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = getAndSaveFile(v.Url, tmpFile, v.Menu)
					if err != nil {
						return
					}
				}

				if strings.Contains(v.Url, "gitlab.luizalabs.com") {
					err = downloadGitlab(v.Url, tmpFile)
					if err != nil {
						return
					}
				}

				justRunValidate(dirTmp, v)

				//
				runMarkdownReport(filepath.Join(dir, v.File), tmpFile)
			}
			spinner.Success("Specs downloaded successfully")
		},
	}
	cmd.Flags().StringVarP(&dir, "dir", "d", "", "Directory to save the converted specs")
	cmd.Flags().StringVarP(&menu, "menu", "m", "", "Menu to download the specs")
	return cmd
}

func runMarkdownReport(left, right string) error {

	updateChan := make(chan *markdownReport.ProgressUpdate)
	errorChan := make(chan markdownReport.ProgressError)
	doneChan := make(chan bool)
	failed := false
	baseFlag := ""

	noColorFlag := false
	cdnFlag := false
	remoteFlag := false
	reportFile := "report.md"
	extRefs := false

	if noColorFlag {
		pterm.DisableStyling()
		pterm.DisableColor()
	}

	listenForUpdates := func(updateChan chan *markdownReport.ProgressUpdate, errorChan chan markdownReport.ProgressError) {
		var spinner *pterm.SpinnerPrinter
		if !noColorFlag {
			spinner, _ = pterm.DefaultSpinner.Start("starting work.")

			spinner.InfoPrinter = &pterm.PrefixPrinter{
				MessageStyle: &pterm.Style{pterm.FgLightCyan},
				Prefix: pterm.Prefix{
					Style: &pterm.Style{pterm.FgBlack, pterm.BgLightMagenta},
					Text:  " SPEC ",
				},
			}
			spinner.SuccessPrinter = &pterm.PrefixPrinter{
				MessageStyle: &pterm.Style{pterm.FgLightCyan},
				Prefix: pterm.Prefix{
					Style: &pterm.Style{pterm.FgBlack, pterm.BgLightCyan},
					Text:  " DONE ",
				},
			}
		}

		var warnings []string

		for {
			select {
			case update, ok := <-updateChan:
				if ok {
					if !noColorFlag {
						if !update.Completed {
							spinner.UpdateText(update.Message)
						} else {
							spinner.Info(update.Message)
						}
					}
					if update.Warning {
						warnings = append(warnings, update.Message)
					}
				} else {
					if !failed {
						if !noColorFlag {
							spinner.Success("completed processing")
							spinner.Stop()
							pterm.Println()
							pterm.Println()
						}
					} else {
						if !noColorFlag {
							spinner.Fail("failed to complete. sorry!")
							pterm.Println()
							pterm.Println()
						}
					}
					if len(warnings) > 0 {
						pterm.Warning.Print("warnings reported during processing")
						pterm.Println()
						dupes := make(map[string]bool)
						for w := range warnings {
							sum := md5.Sum([]byte(warnings[w]))
							md5 := hex.EncodeToString(sum[:])
							if !dupes[md5] {
								dupes[md5] = true
							} else {
								continue
							}
							pterm.Println(fmt.Sprintf("⚠️  %s", pterm.FgYellow.Sprint(warnings[w])))
						}
					}
					doneChan <- true
					return
				}
			case err := <-errorChan:
				if err.Fatal {
					if !noColorFlag {
						spinner.Fail(fmt.Sprintf("Stopped: %s", err.Message))
						spinner.Stop()
					} else {
						pterm.Error.Println(err)
					}
					// doneChan <- true
					//return
				} else {
					warnings = append(warnings, err.Message)
				}
			}
		}
	}

	var urlErr error

	go listenForUpdates(updateChan, errorChan)

	// check if the first arg is a URL, if so download it, if not - assume it's a file.
	left, urlErr = markdownReport.CheckURL(left, errorChan)
	if urlErr != nil {
		pterm.Error.Println(urlErr.Error())
		return urlErr
	}

	// check if the second arg is a URL, if so download it, if not - assume it's a file.
	right, urlErr = markdownReport.CheckURL(right, errorChan)
	if urlErr != nil {
		pterm.Error.Println(urlErr.Error())
		return urlErr
	}

	report, errs := RunLeftRightMarkDownReport(left, right, cdnFlag, updateChan, errorChan, baseFlag, remoteFlag, extRefs)
	<-doneChan
	if len(errs) > 0 {
		for e := range errs {
			pterm.Error.Println(errs[e].Error())
		}
		return errors.New("unable to process specifications")
	}

	return markdownReport.WriteReportFile(reportFile, report)
}

func RunLeftRightMarkDownReport(left, right string, useCDN bool,
	progressChan chan *markdownReport.ProgressUpdate, errorChan chan markdownReport.ProgressError, base string, remote, extRefs bool) ([]byte, []error) {

	var leftBytes, rightBytes []byte
	var errs []error
	var err error

	leftBytes, err = os.ReadFile(left)
	if err != nil {
		markdownReport.SendFatalError("extraction",
			fmt.Sprintf("cannot read original spec: %s", err.Error()), errorChan)
		close(progressChan)
		return nil, []error{err}
	}
	rightBytes, err = os.ReadFile(right)
	if err != nil {
		markdownReport.SendFatalError("extraction",
			fmt.Sprintf("cannot read modified spec: %s", err.Error()), errorChan)
		close(progressChan)
		return nil, []error{err}
	}

	commits := []*markdownReport.Commit{
		{
			Hash:       uuid.New().String()[:6],
			Message:    fmt.Sprintf("Original: %s. Modified: %s", left, right),
			CommitDate: time.Now(),
			Data:       rightBytes,
			FilePath:   right,
		}, {
			Hash:       uuid.New().String()[:6],
			Message:    fmt.Sprintf("Original file: %s", left),
			CommitDate: time.Now(),
			Data:       leftBytes,
			FilePath:   left,
		},
	}

	commits, errs = markdownReport.BuildCommitChangelog(commits, progressChan, errorChan, base, remote, extRefs)
	if len(errs) > 0 {
		close(progressChan)
		return nil, errs
	}
	generator := markdownReport.NewMarkdownReport(false, time.Now(), commits)

	close(progressChan)
	return generator.GenerateReport(), nil
}
