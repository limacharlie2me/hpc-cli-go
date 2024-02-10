package main

import (
	"github.com/rotisserie/eris"

	app "bitbucket.oci.oraclecorp.com/one-punch/compute-go-app"
)

func main() {
	app.Run(
		"hpc-cli",
		version,
		app.WithKongFlagCLI(&cli),
		app.WithFunction(run),
	)
}

func run() error {
	ctx, cancel, g := app.ErrGroupContext()

	g.Go(func() error {
		err := hpc(ctx)

		cancel()

		return eris.Wrap(err, "hpc")
	})

	if err := g.Wait(); err != nil {
		return eris.Wrap(err, "Wait")
	}

	return nil
}
