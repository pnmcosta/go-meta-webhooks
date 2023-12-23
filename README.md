# Go Meta Webhooks

This project provides a type safe payload verification, validation, parsing and handling for [Meta's Webhooks](https://developers.facebook.com/docs/graph-api/webhooks/) Objects, Fields, and Values event notifications with the following features:

- Subscription [verification requests](https://developers.facebook.com/docs/graph-api/webhooks/getting-started#verification-requests)
- Event SHA256 [signature validation](https://developers.facebook.com/docs/graph-api/webhooks/getting-started#event-notifications)
- JSON Schema [payload validation](./schema.json)
- JSON payload struct unmarshalling
- Standard HTTP request support and [libraries](./samples/)
- Concurrent batch processing
- Composable handler interfaces
- Option configurations

## Install

```console
go get github.com/pnmcosta/go-meta-webhooks
```

## Import

```go
import gometawebhooks "github.com/pnmcosta/go-meta-webhooks"
```

## Use Cases

Wire-in this package with any HTTP server package, an example using [Echo](https://echo.labstack.com/) is provided in the [samples](./samples/) directory submodule of this repository. 

The example is an implementation of the [InstagramHandler](./handler_instagram.go) which covers supported Instagram field changes and messaging.

### Scoped Handlers

You can granually implement each handler for scoped support instead. For example, to only handle [InstagramMessageHandler](./messaging_instagram.go) event only instead:

```go
package main

import gometawebhooks "github.com/pnmcosta/go-meta-webhooks"

var _ gometawebhooks.InstagramMessageHandler = (*handler)(nil)

type handler struct{}

func (h handler) InstagramMessage(ctx context.Context, sender, recipient string, sent time.Time, message Message){
    // TODO: implement message handling
}

func main(){
    handler := handler{}
    hooks, err := gometawebhooks.New(
        gometawebhooks.Options.InstagramMessageHandler(handler),
    )

    // TODO: implement HTTP routes see samples directory.
}
```

### Unsupported Objects

Currently only Instagram Webhook objects and fields are explicitly supported, I plan on growing the package over time, however you can implement [EntryHandler](./events.go), [ChangesHandler](./changes.go) or [MessagingHandler](./messaging.go) to handle events for unsupported objects.
