package main

import (
	"fmt"
	"log"
	"time"
	"up-planilhas-go/core"
)

func main() {
	start := time.Now()
	db, err := core.NewDbConnection()
	if err != nil {
		log.Fatalf("Error on DB connection")
	}

	worker := core.NewWorker(db)
	err = worker.Start()

	elapsed := time.Since(start)
	fmt.Printf("Tempo de processamento com %d workers: %s", worker.Workers, elapsed)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
