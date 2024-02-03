package main

import (
	"backend/internal/server"
	"fmt"
)

func main() {

	server := server.NewServer()
	fmt.Printf("Server is running on port %s\n", server.Addr)
	err := server.ListenAndServe()
	if err != nil {
		panic(fmt.Sprintf("cannot start server: %s", err))
	}
}
