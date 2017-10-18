package main

import (
	"net/url"

	trello "github.com/VojtechVitek/go-trello"
)

type Webhook struct {
	client      *trello.Client
	ID          string `json:"id"`
	Description string `json:"description"`
	IDModel     string `json:"idModel"`
	CallbackURL string `json:"callbackUrl"`
	Active      bool   `json:"active"`
}

// Delete delete webhook
func (w *Webhook) Delete() error {
	_, err := w.client.Delete("/webhooks/" + w.ID)
	return err
}

// Save persist the webhook to trello
func (w *Webhook) Save() error {
	var active int
	if w.Active {
		active = 1
	}
	payload := url.Values{}
	payload.Add("idModel", w.IDModel)
	payload.Add("description", w.Description)
	payload.Add("callbackUrl", w.CallbackURL)
	payload.Add("active", string(active))

	var error error

	if w.ID != "" {
		_, err := w.client.Put("/webhooks/"+w.ID, payload)
		error = err
	} else {
		// create
		_, err := w.client.Post("/webhooks", payload)
		error = err
	}
	return error
}
