package main

import (
	"fmt"

	"github.com/cfpb/rhobot/healthcheck"
)

func main() {
	fmt.Printf(healthcheck.Hello())
}
