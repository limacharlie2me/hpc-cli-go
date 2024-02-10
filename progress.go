package main

import (
	"github.com/vbauerster/mpb/v8"
	"github.com/vbauerster/mpb/v8/decor"

	app "bitbucket.oci.oraclecorp.com/one-punch/compute-go-app"
)

type progressbar struct {
	multi *mpb.Progress
}

func newProgressbar() *progressbar {
	options := []mpb.ContainerOption{mpb.WithWidth(32)}

	if app.Debug {
		options = append(options, mpb.WithOutput(nil))
	}

	return &progressbar{
		mpb.New(options...)}
}

func (p *progressbar) Bar(name string) *mpb.Bar {
	return p.multi.AddBar(0,
		mpb.PrependDecorators(
			decor.Name(name+"\t"),
			decor.Counters(0, " %d / %d "),
		),
		mpb.AppendDecorators(
			decor.OnComplete(decor.Name("", decor.WC{}), "Complete"),
		),
	)
}

func (p *progressbar) Wait() {
	p.multi.Wait()
}
