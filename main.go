package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/akerl/go-lambda/apigw/events"
	"github.com/akerl/go-lambda/mux"
	"github.com/akerl/go-lambda/s3"
)

type target struct {
	URL    string `json:"url"`
	Method string `json:"method,omitempty"`
}

type hook struct {
	Targets []target `json:"targets"`
}

type config struct {
	Hooks map[string]hook `json:"hooks"`
}

var c config

func handler(req events.Request) (events.Response, error) {
	hookID := req.PathParameters["hook"]
	if hookID == "" {
		return events.Fail("hook id not set")
	}
	h, ok := config.Hooks[hookID]
	if !ok {
		return events.Fail(fmt.Sprintf("hook not found in config: %s", hookID))
	}
	log.Printf("processing hook: %s", hookID)

	client := http.Client{}

	for _, t := range h.Targets {
		req, err := http.NewRequest(t.Method, t.URL, nil)
		if err != nil {
			return events.Fail(fmt.Sprintf("failed to parse request: %s", t.URL))
		}
		_, err = client.Do(req)
		if err != nil {
			return events.Fail(fmt.Sprintf("failed to hit url: %s (%s)", t.URL, err))
		}
	}
	log.Printf("successfully hit %d urls", len(h.Targets))
}

func loadConfig() {
	cf, err := s3.GetConfigFromEnv(&c)
	if err != nil {
		log.Print(err)
		return
	}
	cf.OnError = func(_ *s3.ConfigFile, err error) {
		log.Print(err)
	}
	cf.Autoreload(60)
}

func main() {
	loadConfig()
	d := mux.NewReceiver(mux.NoCheck, mux.NoAuth, handler, mux.NoError)
	mux.Start(d)
}
