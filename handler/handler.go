// Copyright (c) 2014 The meeko-collector-circleci AUTHORS
//
// Use of this source code is governed by the MIT License
// that can be found in the LICENSE file.

package handler

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
)

const EventType = "circleci.build"

const (
	statusUnprocessableEntity = 422
	maxBodySize               = int64(10 << 20)
)

type WebhookBody struct {
	Payload map[string]interface{}
}

type WebhookHandler struct {
	Logger  Logger
	Forward func(eventType string, eventObject interface{}) error
}

func (handler *WebhookHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Read the event object.
	bodyReader := http.MaxBytesReader(w, r.Body, maxBodySize)
	defer bodyReader.Close()

	body, err := ioutil.ReadAll(bodyReader)
	if err != nil {
		http.Error(w, "Request Payload Too Large", http.StatusRequestEntityTooLarge)
		handler.Logger.Warnf("POST from %v: Request payload too large", r.RemoteAddr)
		return
	}

	// Unmarshal the event object.
	var payload WebhookBody
	if err := json.Unmarshal(body, &payload); err != nil {
		http.Error(w, "Invalid Json", http.StatusBadRequest)
		handler.Logger.Warnf("POST from %v: Invalid json: %v", r.RemoteAddr, err)
		return
	}

	if payload.Payload == nil {
		http.Error(w, "Payload Field Missing", http.StatusBadRequest)
		handler.Logger.Warnf("POST from %v: Payload field missing", r.RemoteAddr)
		return
	}

	// Publish the event.
	if err := handler.Forward(EventType, payload.Payload); err != nil {
		http.Error(w, "Event Not Published", http.StatusInternalServerError)
		handler.Logger.Critical(err)
		return
	}

	handler.Logger.Infof("POST from %v: Forwarding %v", r.RemoteAddr, EventType)
	w.WriteHeader(http.StatusAccepted)
}
