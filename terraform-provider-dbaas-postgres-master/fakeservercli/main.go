package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"terraform-provider-dbaas/fakeserver"
	"time"
)

func main() {
	rand.Seed(time.Now().UnixNano())
	defaultPort := 49152 + rand.Intn(65535-49152) // Random port in dynamic range

	databases := make(map[string]fakeserver.DatabaseProperties)

	port := flag.Int("port", defaultPort, "The port the fakeserver will listen on")
	debug := flag.Bool("debug", false, "Enable debug output")

	flag.Parse()

	svr := fakeserver.NewFakeServer(*port, databases, true, *debug)

	fmt.Printf("Starting server on port %d...\n", *port)
	fmt.Println("Database Endpoint: /api/database/{name}")

	internalServer := svr.GetServer()
	if err := internalServer.ListenAndServe(); err != nil {
		fmt.Printf("Error with the server: %s", err)
		os.Exit(1)
	}
}
