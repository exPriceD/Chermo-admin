package main

import (
	"fmt"
	"github.com/exPriceD/Chermo-admin/internal/server"
)

func main() {

	srv := server.NewServer()

	err := srv.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
