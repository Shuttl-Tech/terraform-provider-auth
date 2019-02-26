package main

import (
	"fmt"
	"github.com/Shuttl-Tech/terraform-provider-auth/auth"
	"github.com/hashicorp/terraform/plugin"
)

const Version = "0.0.1"

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
