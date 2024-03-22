package generator

import (
	mgcSdkPkg "magalu.cloud/sdk"
)

func generateCommon(p string, sdk *mgcSdkPkg.Sdk, ctx *GeneratorContext) (err error) {
	if err = generateGoMod(p, sdk, ctx); err != nil {
		return
	}

	if err = generateHelpers(p, sdk, ctx); err != nil {
		return
	}

	if err = generateClient(p, sdk, ctx); err != nil {
		return
	}

	return nil
}