package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	u "./url"
	"github.com/VojtechVitek/go-trello"
	"github.com/benrowe/trello-downloader/services"
	"github.com/spf13/viper"
)

// import "flag"

// Webhook d

type trelloLabel struct {
	name    string
	service services.Service
}

type internalWebhook struct {
	baseURL u.URL
}

// config interface
type config interface {
	GetString(string) string
	GetStringMap(key string) map[string]interface{}
	UnmarshalKey(key string, rawVal interface{}) error
}

func loadConfiguration() config {
	viper.SetConfigName("config")

	viper.AddConfigPath(".")
	viper.AddConfigPath("config")

	viper.ReadInConfig()

	return viper.GetViper()
}

//
func extractLabelServices(b trello.Board, c config, s map[string]services.Service) map[string]trelloLabel {
	newmap := make(map[string]trelloLabel)

	labels := c.GetStringMap("trello.labels")

	for something := range labels {
		serviceName := c.GetString("trello.labels." + something + ".service")
		newmap[something] = trelloLabel{something, s[serviceName]}
	}

	return newmap
}

func main() {
	// boot
	config := loadConfiguration()
	validateConfig(config)
	services := getDownloadServices(config)
	client := getTrello(config.GetString("trello.api.appKey"), config.GetString("trello.api.token"))
	board := loadBoard(client, config.GetString("trello.boardID"))

	validateBoard(board, config)

	labels := extractLabelServices(board, config, services)

	var webhook Webhook

	go handleWebhookEvents()
	fmt.Print(labels)
	go registerTrelloWebhook(&webhook, board, client)
	// go registerServicesWebhook(labels)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	end := make(chan bool)
	// handle signals
	go func() {
		sig := <-signalChannel
		switch sig {
		case os.Interrupt:
			fmt.Println("interrupt")
			deleteWebhook(webhook, client)
			end <- true
		case syscall.SIGTERM:
			fmt.Println("sigterm")
			deleteWebhook(webhook, client)
			end <- true
		}
	}()
	// force the application to end when this channel receives an update
	<-end
}

// retrieve a list of all the available download services we might need to support
func getDownloadServices(c config) map[string]services.Service {

	s := map[string]services.Service{}

	for key := range c.GetStringMap("downloadServices") {
		d, err := services.Make(key, c.GetString("downloadServices."+key+".label"), c.GetString("downloadServices."+key+".baseUrl"))
		if err == nil {
			s[key] = d
		}

	}
	return s
}

// Get the trello client
func getTrello(appKey string, token string) trello.Client {
	client, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		panic(err)
	}
	return *client
}

// validate that we have a valid config file
func validateConfig(config config) {
	// validate trello api details
	if len(config.GetString("trello.api.appKey")) != 32 {
		panic("config: trello.api.appKey invalid")
	}

	if len(config.GetString("trello.api.token")) != 64 {
		panic("config: trello.api.token invalid")
	}
	// make sure we have at least one label correctly registered

}

// retrieve an instance of the trello board
func loadBoard(client trello.Client, boardID string) trello.Board {
	board, err := client.Board(boardID)
	if err != nil {
		panic(err)
	}
	return *board
}

// validate the state of the board against the provided configuration
func validateBoard(board trello.Board, config config) {
	// ensure listed labels in config exist

}

func deleteWebhook(w Webhook, t trello.Client) {
	err := w.Delete()
	if err != nil {
		panic(err)
	}
}

func registerTrelloWebhook(w *Webhook, board trello.Board, t trello.Client) {
	// try to register the webhook
	payload := url.Values{}
	payload.Add("idModel", board.Id)
	payload.Add("description", "something")
	payload.Add("callbackURL", "http://1.129.106.234:2600/?a")
	payload.Add("active", "1")
	body, err := t.Post("/webhooks", payload)
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, w); err != nil {
		panic(err)
	}
	fmt.Println("registered trello webhook")
}

// register an event hook with each of the unique download services
func registerServicesWebhook(labels map[string]trelloLabel) {
	// var registered []string
	// for labelName, label := range labels {
	// 	if !stringInSlice(labelName, registered) {
	// 		// _, ok := interface{}(label).(services.SupportsWebhookEvents)
	// 		// if ok {
	// 		// 	label.service.RegisterWebhook()
	// 		// } else {
	// 		// 	// have to result to polling the service
	// 		// }
	// 	}
	// }
}

func handleWebhookEvents() {
	fmt.Println("Starting web server")
	http.HandleFunc("/", trelloWebhookHandler)
	http.ListenAndServe(":2600", nil)
	// register a endpoint to handle trello update events + deligate that off to a specific handler

	// handle events from each service + associated label + action
	// forward those requests to the relevant service to be handled + a trello handler to update the state as necessary

}

func trelloWebhookHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(r.URL)
	body, _ := ioutil.ReadAll(r.Body)
	fmt.Print(string(body))
}

func printBoard(board trello.Board) {
	fmt.Printf("* %v %v (%v)\n", board.Name, board.ShortUrl, board.Id)

	// @trello Board Lists
	lists, err := board.Lists()
	if err != nil {
		log.Fatal(err)
	}

	for _, list := range lists {
		fmt.Println("   - ", list.Name)

		// @trello Board List Cards
		cards, _ := list.Cards()
		for _, card := range cards {
			fmt.Println("      + ", card.Name)
		}
	}
}

func findBoard(boards []trello.Board, boardName string) (trello.Board, error) {
	for i := range boards {
		if strings.ToLower(boards[i].Name) == boardName {
			return boards[i], nil
		}
	}
	err := fmt.Errorf("unknown board")
	return trello.Board{}, err
}

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

func stringInSlice(a string, list []string) bool {
	for _, b := range list {
		if a == b {
			return true
		}
	}
	return false
}
