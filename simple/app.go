package main

import (
	"crypto/md5"
	"encoding/hex"
	"log"

	"github.com/meroxa/turbine"
	"github.com/meroxa/turbine/runner"
)

func main() {
	runner.Start(App{})
}

var _ turbine.App = (*App)(nil)

type App struct{}

func (a App) Run(v turbine.Turbine) error {
	db, err := v.Resources("demopg")
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
	err = s3.Write(res, "data-app-archive", nil)
	if err != nil {
		return err
	}

	return nil
}

type Anonymize struct{}

func (f Anonymize) Process(rr []turbine.Record) ([]turbine.Record, []turbine.RecordWithError) {
	for i, r := range rr {
		hashedEmail := consistentHash(r.Payload.Get("email").(string))
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
