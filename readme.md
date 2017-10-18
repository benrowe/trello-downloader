# Trello Downloader

Uses trello to manage state in the cloud for my download system

List webhooks: https://api.trello.com/1/members/me/tokens?webhooks=true&key=[APPLICATION_KEY]&token=[USER_TOKEN]

## Process flow

### Startup

- Gather required variables
  - trello 
    - api keys
    - mapping of list names to specific states
    - board to manage
  - mananged labels
    - label name
    - label key (if not matching the label name)
    - associated service
  - services
    - name
    - urls
      - search
      - add
      - webhook(?): depends on how that service sends back events
- register webhook with trello for human events
- create http mux to handle trello webhook requests
  - for each webhook request
    - validate webhook signature
    - filter out those requests that are:
      - not made by humans?
      - not in the relavant lists
      
- check the state of the current board:
  - verify the state of the items
    - within the lists specified
    - where the cards in those lists meet the minimum requirements (label, title)
- periodicly verify the state of the board against the state of the application.


### Shutdown
- de-register webhooks
- finish executing current activies/tidy-up