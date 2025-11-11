package orchestrator

import (
	"fmt"
	"net/http"
)

func Orchestrate(r *http.Request) error {
	fmt.Println("Hello, World!")
	return nil

}
