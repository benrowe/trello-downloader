package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/url"
	"os"
	"strings"

	"github.com/VojtechVitek/go-trello"
	"github.com/joho/godotenv"
)

// import "flag"

// Webhook d
type Webhook struct {
	client      *trello.Client
	ID          string `json:"id"`
	Description string `json:"description"`
	IDModel     string `json:"idModel"`
	CallbackURL string `json:"callbackUrl"`
	Active      bool   `json:"active"`
}

type label struct {
	name    string
	key     string
	service *service
}

type service struct {
	name string
	urls struct {
		search        string
		add           string
		requestUpdate string
	}
}

func loadEnv() {
	err := godotenv.Load()

	if err != nil {
		panic(err)
	}
}

func main() {

	loadEnv()

	args := len(os.Args)
	argument := os.Args[args-1]
	fmt.Println(argument)

	appKey := os.Getenv("TRELLO_APP_KEY")
	token := os.Getenv("TRELLO_TOKEN")
	trello, err := trello.NewAuthClient(appKey, &token)
	if err != nil {
		log.Fatal(err)
	}

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

	user, err := trello.Member("benrowe")

	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(user.FullName, user.Bio)

	boards, err := user.Boards()
	if err != nil {
		log.Fatal(err)
	}

	board, err := findBoard(boards, argument)

	if err != nil {
		panic(err)
	}
	fmt.Println(board)
	registerTrelloWebHook(board.Id, trello)

	// go printBoard(board)
	// go printBoard(board)
	// fmt.Scanln()

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

func registerTrelloWebHook(boardId string, client *trello.Client) {
	payload := url.Values{}
	payload.Add("IDModel", boardId)
	payload.Add("Description", "something")
	payload.Add("CallbackURL", "")
	payload.Add("Active", "1")
	// webhook := &Webhook{IDModel: "", Description: "", CallbackURL: "", Active: true}
	// payload, _ := json.Marshal(payload)
	// url.
	fmt.Println(boardId)
	body, err := client.Post("/webhooks", payload)
	if err != nil {
		panic(err)
	}

	newList := &Webhook{}

	if err = json.Unmarshal(body, newList); err != nil {
		panic(err)
	}

	fmt.Println(newList.ID)

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
