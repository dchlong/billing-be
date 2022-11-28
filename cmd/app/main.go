package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	_ "github.com/golang-migrate/migrate/v4/source/file"

	"github.com/dchlong/billing-be/internal"
)

func main() {
	migrate()
	server, cleanup, err := internal.InitializeServer()
	if err != nil {
		log.Printf("could not initialize server, error: %+v", err)
		panic(err)
	}

	go func() {
		err = server.Start()
		if err != nil {
			log.Printf("could not start server, error: %+v", err)
			panic(err)
		}
	}()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGINT, syscall.SIGTERM)
	s := <-c
	log.Printf("Receive os.Signal: %s", s.String())
	cleanup()
}

func migrate() {
	log.Printf("Migrating....")
	migrationTool, cleanup, err := internal.InitializeMigrationTool()
	if err != nil {
		log.Printf("could not initialize migration tool, error: %+v", err)
		panic(err)
	}

	err = migrationTool.Migrate()
	if err != nil {
		log.Printf("migrate failed, error: %+v", err)
		panic(err)
	}

	log.Printf("Migration done...")
	cleanup()
}
