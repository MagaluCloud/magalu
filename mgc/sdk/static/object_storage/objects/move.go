package objects

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"magalu.cloud/core"
	mgcSchemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
	"magalu.cloud/sdk/static/object_storage/common"
)

type moveParams struct {
	Source      mgcSchemaPkg.URI `json:"src" jsonschema:"description=A file to move to the Destination,example=./hello.txt" mgc:"positional"`
	Destination mgcSchemaPkg.URI `json:"dst" jsonschema:"description=Destination to put the file into,example=s3://my-bucket/test.txt" mgc:"positional"`
}

var getMove = utils.NewLazyLoader[core.Executor](func() core.Executor {
	executor := core.NewStaticExecute(
		core.DescriptorSpec{
			Name:    "move",
			Summary: "Moves one object from source to destination",
			Description: `Moves one object from a source to a destination.
It can be either local or remote but not both local (Local -> Remote, Remote -> Local, Remote -> Remote)`,
		},
		move,
	)

	return core.NewExecuteResultOutputOptions(executor, func(exec core.Executor, result core.Result) string {
		return "template=Moved from {{.src}} to {{.dst}}\n"
	})
})

func move(ctx context.Context, params moveParams, cfg common.Config) (moveParams, error) {
	srcIsRemote := isRemote(params.Source)
	dstIsRemote := isRemote(params.Destination)

	if !srcIsRemote && dstIsRemote {
		return moveLocalRemote(ctx, params, cfg)
	}
	if srcIsRemote && !dstIsRemote {
		return moveRemoteLocal(ctx, params, cfg)
	}
	if srcIsRemote && dstIsRemote {
		return moveRemote(ctx, params, cfg)
	}
	if !srcIsRemote && !dstIsRemote {
		return params, core.UsageError{Err: fmt.Errorf("operation not supported, this command cannot be used to move a local source to a local destination")}
	}

	return params, nil
}

func moveRemote(ctx context.Context, params moveParams, cfg common.Config) (moveParams, error) {
	_, err := copy(ctx, common.CopyObjectParams{Source: params.Source, Destination: params.Destination}, cfg)
	if err != nil {
		return params, err
	}

	_, err = deleteObject(ctx, common.DeleteObjectParams{Destination: params.Source}, cfg)
	if err != nil {
		return params, err
	}

	return params, err
}

func moveLocalRemote(ctx context.Context, params moveParams, cfg common.Config) (moveParams, error) {
	srcAbs, err := filepath.Abs(params.Source.String())
	if err != nil {
		return params, err
	}

	_, err = upload(ctx, uploadParams{Source: mgcSchemaPkg.FilePath(params.Source), Destination: params.Destination}, cfg)
	if err != nil {
		return params, err
	}

	err = os.Remove(srcAbs)
	if err != nil {
		return params, err
	}

	return params, nil
}

func moveRemoteLocal(ctx context.Context, params moveParams, cfg common.Config) (moveParams, error) {
	_, err := download(ctx, common.DownloadObjectParams{Source: params.Source, Destination: mgcSchemaPkg.FilePath(params.Destination)}, cfg)
	if err != nil {
		return params, err
	}

	_, err = deleteObject(ctx, common.DeleteObjectParams{Destination: params.Source}, cfg)
	if err != nil {
		return params, err
	}

	return params, nil
}
