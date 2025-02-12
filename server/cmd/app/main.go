package main

// Boot up backend services to listen on some local port
// The services will handle all of the downstream logic to on http requests etc.

import (
	"github.com/obrown4/credit-stack/server/internal/app"
)

func main() {
	app.Boot()
}
