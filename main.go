package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform/command"
	"github.com/hashicorp/terraform/config/module"
	"github.com/hashicorp/terraform/svchost"
	"github.com/hashicorp/terraform/svchost/auth"
	"github.com/hashicorp/terraform/svchost/disco"

	"github.com/pgavlin/firewalker/il"
)

type noCredentials struct{}

func (noCredentials) ForHost(host svchost.Hostname) (auth.HostCredentials, error) {
	return nil, nil
}

func main() {
	credentials := noCredentials{}
	services := disco.NewDisco()
	services.SetCredentialsSource(credentials)
	moduleStorage := module.NewStorage(filepath.Join(command.DefaultDataDir, "modules"), services, credentials)

	mod, err := module.NewTreeModule("", ".")
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not load module: %v\n", err)
		os.Exit(-1)
	}

	if err = mod.Load(moduleStorage); err != nil {
		fmt.Fprintf(os.Stderr, "could not load module: %v\n", err)
		os.Exit(-1)
	}

	g, err := il.BuildGraph(mod.Config())
	if err != nil {
		fmt.Fprintf(os.Stderr, "could not build graph: %v\n", err)
		os.Exit(-1)
	}

	if err = generate(g, &nodeGenerator{projectName: "auto"}); err != nil {
		fmt.Fprintf(os.Stderr, "generation failed: %v\n", err)
		os.Exit(-1)
	}
}
