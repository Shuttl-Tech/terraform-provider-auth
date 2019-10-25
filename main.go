package main

import (
	"fmt"

	"github.com/Shuttl-Tech/terraform-provider-auth/auth"
	"github.com/hashicorp/terraform-plugin-sdk/plugin"
)

const Version = "0.0.5"

var (
	Name         string
	GitCommit    string
	HumanVersion = fmt.Sprintf("%s v%s (%s)", Name, Version, GitCommit)
)

func main() {
	plugin.Serve(&plugin.ServeOpts{
		ProviderFunc: auth.Provider,
	})
}
