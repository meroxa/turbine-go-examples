package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"

	turbine "github.com/meroxa/turbine-go"
	"github.com/meroxa/turbine-go/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	db, err := v.Resources("demobagel")
	if err != nil {
		return err
	}

	rr, err := db.Records("user_activity", nil) // rr is a collection of records, can't be inspected directly
	if err != nil {
		return err
	}

	res, _ := v.Process(rr, Anonymize{})
	// second return is dead-letter queue

	s3, err := v.Resources("s3")
	if err != nil {
		return err
	}
	err = s3.Write(res, "data-app-archive")
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	for i, r := range rr {
		e := fmt.Sprintf("%s", r.Payload.Get("email"))
		if e == "" {
			log.Printf("unable to find email value in %d record\n", i)
			break
		}
		hashedEmail := consistentHash(e)
		err := r.Payload.Set("email", hashedEmail)
		if err != nil {
			log.Println("error setting value: ", err)
			break
		}
		rr[i] = r
	}
	return rr, nil
}

func consistentHash(s string) string {
	h := md5.Sum([]byte(s))
	return hex.EncodeToString(h[:])
}
