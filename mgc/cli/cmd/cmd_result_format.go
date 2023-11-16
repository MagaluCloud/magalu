package cmd

import (
	"fmt"
	"io"
	"os"

	"github.com/spf13/cobra"
	"magalu.cloud/core"
	mgcSdk "magalu.cloud/sdk"
)

func formatResult(sdk *mgcSdk.Sdk, cmd *cobra.Command, result core.Result) error {
	output := getOutputFor(sdk, cmd, result)

	if resultWithReader, ok := core.ResultAs[core.ResultWithReader](result); ok {
		return handleResultWithReader(resultWithReader.Reader(), output)
	}

	if resultWithValue, ok := core.ResultAs[core.ResultWithValue](result); ok {
		return handleResultWithValue(resultWithValue, output)
	}

	return fmt.Errorf("unsupported result: %T %+v", result, result)
}

func handleResultWithReader(reader io.Reader, outFile string) (err error) {
	if closer, ok := reader.(io.Closer); ok {
		defer closer.Close()
	}

	var writer io.WriteCloser
	if outFile == "" || outFile == "-" {
		writer = os.Stdout
	} else {
		writer, err = os.OpenFile(outFile, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
		if err != nil {
			return err
		}
	}

	n, err := io.Copy(writer, reader)
	defer writer.Close()
	if err != nil {
		return fmt.Errorf("Wrote %d bytes. Error: %w\n", n, err)
	}
	return nil
}

func handleResultWithValue(result core.ResultWithValue, output string) (err error) {
	value := result.Value()
	if value == nil {
		return nil
	}

	err = result.ValidateSchema()
	if err != nil {
		logger().Warnw("result validation failed", "error", err.Error())
	}

	name, options := parseOutputFormatter(output)
	if name == "" {
		if formatter, ok := core.ResultAs[core.ResultWithDefaultFormatter](result); ok {
			fmt.Println(formatter.DefaultFormatter())
			return nil
		}
	}

	formatter, err := getOutputFormatter(name, options)
	if err != nil {
		return err
	}
	return formatter.Format(value, options)
}

func handleSimpleResultValue(value core.Value, output string) error {
	name, options := parseOutputFormatter(output)
	formatter, err := getOutputFormatter(name, options)

	if err != nil {
		return err
	}

	return formatter.Format(value, options)
}
