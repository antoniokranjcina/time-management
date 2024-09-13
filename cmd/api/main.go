package main

import (
	"fmt"
	"time-management/internal/server"
)

func main() {
	newServer := server.NewServer()
	err := newServer.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
