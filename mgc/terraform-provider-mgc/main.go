package main

import (
	"context"
	"flag"
	"log"

	"github.com/MagaluCloud/magalu/mgc/terraform-provider-mgc/mgc"
	"github.com/hashicorp/terraform-plugin-framework/providerserver"
)

var (
	// these will be set by the goreleaser configuration
	// to appropriate values for the compiled binary.
	version string = "dev"
	commit  string = "none"
	date    string = "unknown"
)

func main() {
	var debug bool

	flag.BoolVar(&debug, "debug", false, "set to true to run the provider with support for debuggers like delve")
	flag.Parse()

	opts := providerserver.ServeOpts{
		Address: "registry.terraform.io/magalucloud/mgc",
		Debug:   debug,
	}

	err := providerserver.Serve(context.Background(), mgc.New(version, commit, date), opts)

	if err != nil {
		log.Fatal(err.Error())
	}
}
