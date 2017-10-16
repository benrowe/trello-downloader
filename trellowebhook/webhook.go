package trellowebhook

import trello "github.com/VojtechVitek/go-trello"

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
