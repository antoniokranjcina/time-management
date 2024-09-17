package main

import (
	"fmt"
	"time-management/internal"
)

func main() {
	server := internal.NewServer()
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
