package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"os/signal"
	"strings"
	"syscall"

	"github.com/VojtechVitek/go-trello"
	"github.com/benrowe/trello-downloader/services"
	u "github.com/benrowe/trello-downloader/url"
	"github.com/spf13/viper"
)

// import "flag"

// Webhook d
type webhook struct {
	client      *trello.Client
	ID          string `json:"id"`
	Description string `json:"description"`
	IDModel     string `json:"idModel"`
	CallbackURL string `json:"callbackUrl"`
	Active      bool   `json:"active"`
}

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
	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	end := make(chan bool)
	config := loadConfiguration()
	validateConfig(config)
	services := getDownloadServices(config)
	client := getTrello(config.GetString("trello.api.appKey"), config.GetString("trello.api.token"))
	board := loadBoard(client, config.GetString("trello.boardID"))

	validateBoard(board, config)

	extractLabelServices(board, config, services)

	var webhook webhook

	go registerTrelloWebhook(&webhook, board, client)
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
	// go registerServicesWebhook()
	// go handleTrelloWebhooks()

	// args := len(os.Args)
	// argument := os.Args[args-1]
	// fmt.Println(argument)

	// appKey := os.Getenv("TRELLO_APP_KEY")
	// token := os.Getenv("TRELLO_TOKEN")
	// trello, err := trello.NewAuthClient(appKey, &token)
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// body, err := trello.Get("/webhooks")
	// if err != nil {
	// 	panic(err)
	// }
	// var webhooks []Webhook
	// err = json.Unmarshal(body, &webhooks)
	// for i := range webhooks {
	// 	webhooks[i].client = trello
	// }

	// fmt.Println(webhooks)

	// user, err := trello.Member("benrowe")

	// if err != nil {
	// 	log.Fatal(err)
	// }
	// fmt.Println(user.FullName, user.Bio)

	// boards, err := user.Boards()
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// board, err := findBoard(boards, argument)

	// if err != nil {
	// 	panic(err)
	// }
	// fmt.Println(board)
	// registerTrelloWebHook(board.Id, trello)

	// go printBoard(board)
	// go printBoard(board)
	// fmt.Scanln()

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

func deleteWebhook(w webhook, t trello.Client) {
	_, err := t.Delete("/webhooks/" + w.ID)
	if err != nil {
		panic(err)
	}
}

func registerTrelloWebhook(w *webhook, board trello.Board, t trello.Client) {
	payload := url.Values{}
	payload.Add("idModel", board.Id)
	payload.Add("description", "something")
	payload.Add("callbackURL", "http://benrowe.info/?a")
	payload.Add("active", "1")
	body, err := t.Post("/webhooks", payload)
	fmt.Println(string(body))
	if err != nil {
		panic(err)
	}

	if err = json.Unmarshal(body, w); err != nil {
		panic(err)
	}
}

func registerServicesWebhook() {

}

func handleTrelloWebhooks() {

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
