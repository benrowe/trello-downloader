# Trello Downloader

Uses trello to manage state in the cloud for my download system

## Process flow

### Startup

- Gather required variables
  - trello api keys
  - username/board to manage
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
- register webhook for human events

### Shutdown
- de-register webhooks