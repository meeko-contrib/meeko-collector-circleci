// Copyright (c) 2014 The cider-collector-circleci AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/cider/go-cider/cider/services/logging"
)

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type CircleCiEvent struct {
	Payload map[string]interface{}
}

type CircleCiWebhookHandler struct {
	logger  *logging.Service
	forward func(eventType string, eventObject interface{}) error
}

func (handler *CircleCiWebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the event object.
	bodyReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer bodyReader.Close()

	eventBody, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, "Request Payload Too Large", http.StatusRequestEntityTooLarge)
		handler.logger.Warnf("POST from %v: Request payload too large", r.URL)
		return
	}

	// Unmarshal the event object.
	var event CircleCiEvent
	if err := json.Unmarshal(eventBody, &event); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		handler.logger.Warnf("POST from %v: Invalid json: %v", r.URL, err)
		return
	}

	if event.Payload == nil {
		http.Error(w, "Payload Missing", http.StatusBadRequest)
		handler.logger.Warnf("POST from %v: Payload key missing", r.URL)
	}

	// Publish the event.
	if err := handler.forward("circleci.build", event.Payload); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		handler.logger.Critical(err)
	}

	handler.logger.Infof("POST from %v: Forwarding circleci.build", r.URL)
	w.WriteHeader(http.StatusAccepted)
}
