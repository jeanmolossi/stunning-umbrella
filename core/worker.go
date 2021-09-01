package core

import (
	"io/ioutil"
	"log"
	"net/http"
	"runtime"
	"strings"
	"time"
)

func ExtractIdsFromFile(filename string) ([]string, error) {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return []string{}, err
	}

	usrIds := strings.Split(string(file), "\n")[1:]

	return usrIds, nil
}

func WorkerRunner() {
	usrIds, err := ExtractIdsFromFile("./mock.csv")
	if err != nil {
		log.Fatalf("Erro ao extrair o arquivo")
	}

	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	in := make(chan string, runtime.NumCPU())
	returnChannel := make(chan string)
	done := make(chan bool)

	concurrencyMax := 5

	for process := 1; process <= concurrencyMax; process++ {
		go processUpdate(in, returnChannel, client, process)
	}

	go func() {
		for _, usrId := range usrIds {
			in <- usrId
		}
		close(in)
		log.Println("Producer finished")
	}()

	for r := range returnChannel {
		if r == "update done" {
			done <- true
			break
		}
	}
	log.Println("Consumer finished")

}

func processUpdate(in chan string, returnChannel chan string, client *http.Client, workerId int) {
	log.Printf("Worker %d", workerId)
	for usrId := range in {
		log.Printf("Worker ID %d : usr : %s", workerId, usrId)
		r := NewRequester(usrId, client)
		err := r.DoUpdate()

		if r.Error != "" || err != nil {
			log.Println(r.Error)
		}

		returnChannel <- usrId
	}
	log.Printf("Worker %d finalizado", workerId)
	returnChannel <- "update done"
}
