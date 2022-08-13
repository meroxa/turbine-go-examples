package main

import (
	"log"

	// Dependencies of Turbine
	"github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
	"github.com/meroxa/turbine-go/transforms"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	source, err := v.Resources("mongo")
	if err != nil {
		return err
	}

	rr, err := source.Records("events", nil)
	if err != nil {
		return err
	}

	res := v.Process(rr, Flatten{})

	dest, err := v.Resources("destination_name")
	if err != nil {
		return err
	}

	err = dest.Write(res, "collection_archive")
	if err != nil {
		return err
	}

	return nil
}

type Flatten struct{}

func (f Flatten) Process(stream []turbine.Record) []turbine.Record {
	for i, r := range stream {
		err := transforms.Flatten(&r.Payload)
		if err != nil {
			log.Printf("error: %s", err.Error())
		}

		stream[i] = r
	}
	return stream
}
