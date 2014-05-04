// Copyright (c) 2014 The cider-collector-circleci AUTHORS
//
// Use of this source code is governed by The MIT License
// that can be found in the LICENSE file.

package main

import receiver "github.com/salsita-cider/cider-webhook-receiver"

func main() {
	receiver.ListenAndServe(&CircleCiWebhookHandler{
		receiver.Logger,
		func(eventType string, eventObject interface{}) error {
			return receiver.PubSub.Publish(eventType, eventObject)
		},
	})
}
