package main

import (
	"fmt"
	"log"
	"time"
	"up-planilhas-go/core/database"
	"up-planilhas-go/core/worker"
)

func main() {
	start := time.Now()
	db, err := database.NewDbConnection()
	if err != nil {
		log.Fatalf("Error on DB connection")
	}

	w := worker.NewWorker(db)
	err = w.Start()

	elapsed := time.Since(start)
	fmt.Printf("Tempo de processamento com %d workers: %s\n", w.Workers, elapsed)

	if err != nil {
		log.Fatalf(err.Error())
	}
}
