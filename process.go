package main

import (
	"context"
	"sync"

	"github.com/rotisserie/eris"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

func process(ctx context.Context, instances []Instance) ([]Instance, error) { //nolint:funlen
	log := log.With().Str("Function", "process").Logger()
	defer log.Debug().Msg("Return")

	instanceChan := make(chan Instance, 128)
	resultChan := make(chan Instance, 128)
	processed := make([]Instance, 0)

	g, ctx := errgroup.WithContext(ctx)

	g.Go(func() error {
		defer close(instanceChan)
		defer log.Debug().Msg("Close instanceChan")

		for _, instance := range instances {
			instanceChan <- instance
		}

		return nil
	})

	g.Go(func() error {
		defer close(resultChan)
		defer log.Debug().Msg("Close resultChan")

		if err := processRoutine(ctx, instanceChan, resultChan); err != nil {
			return eris.Wrap(err, "processRoutine")
		}

		return nil
	})

	p := newProgressbar()

	bar := p.Bar("[INFO] Run HealthChecks: \t\t")
	bar.SetTotal(int64(len(instances)), false)

	var wg sync.WaitGroup

	wg.Add(1)

	go func() {
		defer wg.Done()

		for {
			select {
			case <-ctx.Done():
				log.Debug().Msg("Done")

				bar.SetTotal(-1, true)

				return
			case instance, ok := <-resultChan:
				if !ok {
					bar.SetTotal(-1, true)

					return
				}

				processed = append(processed, instance)

				bar.IncrBy(1)
			}
		}
	}()

	wg.Wait()

	err := g.Wait()

	if err != nil {
		return nil, eris.Wrapf(err, "errGroup error")
	}

	p.Wait()

	return processed, nil
}

func processRoutine(ctx context.Context, instanceChan chan Instance, resultChan chan Instance) error {
	log := log.With().Str("Function", "processRoutine").Logger()
	defer log.Debug().Msg("Return")

	g, ctx := errgroup.WithContext(ctx)

	for i := 0; i < NumThreadsProcessRoutine; i++ {
		threadNum := i

		g.Go(func() error {
			if err := runHealthCheck(ctx, instanceChan, resultChan, threadNum); err != nil {
				return eris.Wrap(err, "compute.getSmartNicAndIlomHealth")
			}

			return nil
		})
	}

	if err := g.Wait(); err != nil {
		return eris.Wrap(err, "g.Wait")
	}

	return nil
}
