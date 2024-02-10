package main

import (
	"context"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"

	sashay "bitbucket.oci.oraclecorp.com/one-punch/sashay2"
)

func runHealthCheck(ctx context.Context, instanceChan chan Instance, resultChan chan Instance, threadNum int) error { //nolint:cyclop
	log := log.With().Str("Function", "runHealthCheck").Int("Thread", threadNum).Logger()
	defer log.Debug().Msg("Return")

	bareMetalStatsAPI, err := sashay.BareMetalStatsAPI(ctx, cli.Region)
	if err != nil {
		return eris.Wrapf(err, "BareMetalStatsAPI: %s", cli.Region)
	}

	for {
		select {
		case <-ctx.Done():
			log.Debug().Msg("Done")

			return nil
		case instance, ok := <-instanceChan:
			if !ok {
				return nil
			}

			log.Debug().Int("Instance CHAN", len(instanceChan)).Int("Result CHAN", len(resultChan)).Msg("")

			if cli.HealthCheck {
				if instance.response != 200 {
					break
				}

				smartNic, _, err := bareMetalStatsAPI.GetHostSmartNicData(ctx).Identity(instance.HostID).Execute()
				if err != nil {
					return eris.Wrap(err, "GetHostSmartNicData")
				}

				ilom, _, err := bareMetalStatsAPI.GetHostData(ctx).Identity(instance.HostID).Execute()
				if err != nil {
					return eris.Wrap(err, "GetHostData")
				}

				instance.HealthCheck.Smartnic = smartNic
				instance.HealthCheck.Ilom = ilom
			}

			resultChan <- instance
		}
	}
}
