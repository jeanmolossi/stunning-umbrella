package main

import (
	"log"
	"up-planilhas-go/core"
)

func main() {
	done := make(chan bool)

	db, err := core.NewDbConnection()
	if err != nil {
		log.Fatalf("Error on DB connection")
	}

	worker := core.NewWorker(db)
	worker.Workers = 175
	worker.WorkerRunner(done)
}
