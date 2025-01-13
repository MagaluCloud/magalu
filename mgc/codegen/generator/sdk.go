package generator

import (
	"path/filepath"

	mgcSdkPkg "github.com/MagaluCloud/magalu/mgc/sdk"
)

func GenerateSdk(outputDir string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	p, err := filepath.Abs(outputDir)
	if err != nil {
		return err
	}
	err = createDir(ctx, p)
	if err != nil {
		return err
	}

	err = generateCommon(p, sdk, ctx)
	if err != nil {
		return err
	}

	return generateGroups(p, sdk, ctx)
}
