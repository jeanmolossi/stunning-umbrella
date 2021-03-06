package core

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"time"
)

var (
	ErrorType   = "Error"
	SuccessType = "Success"
)

type Worker struct {
	Db         *gorm.DB `gorm:"-"`
	UsrIDs     []string
	Workers    int
	StartIndex int
}

func NewWorker(db *gorm.DB) *Worker {
	return &Worker{
		Db:         db,
		Workers:    10,
		StartIndex: 1,
	}
}

func (w *Worker) ExtractIdsFromFile(filename string) error {
	file, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}

	usrIds := strings.Split(string(file), "\n")[w.StartIndex:]

	for _, usrId := range usrIds {
		w.UsrIDs = append(w.UsrIDs, usrId)
	}

	return nil
}

func (w *Worker) WorkerRunner(done chan bool) {
	filename := os.Getenv("FILENAME")

	err := w.ExtractIdsFromFile(filename)
	if err != nil {
		log.Fatalf("Erro ao extrair o arquivo")
	}

	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	in := make(chan string, runtime.NumCPU())
	returnChannel := make(chan string)
	wg := sync.WaitGroup{}

	concurrencyMax := w.Workers

	for process := 1; process <= concurrencyMax; process++ {
		wg.Add(1)
		go w.processUpdate(in, returnChannel, client, &wg)
	}

	go func() {
		for _, usrId := range w.UsrIDs {
			in <- usrId
		}
		close(in)
	}()

	go func() {
		for r := range returnChannel {
			if r == "update done" {
				done <- true
			}
		}
	}()

	<-returnChannel

	go func() {
		for d := range done {
			if d == true {
				<-done
			}
		}
	}()

	wg.Wait()
}

func (w *Worker) processUpdate(in chan string, returnChannel chan string, client *http.Client, wg *sync.WaitGroup) {
	defer wg.Done()
	for usrId := range in {
		logger := NewLogger(w.Db)
		r := NewRequester(usrId, client)
		err := r.DoUpdate()

		logger.RefID = usrId
		logger.Message = fmt.Sprintf("User [ %s ] has been updated", usrId)
		logger.Type = SuccessType

		if r.Error != "" || err != nil {
			log.Println(r.Error)
			logger.Message = r.Error
			logger.Type = ErrorType
		}

		err = logger.AddLog()
		if err != nil {
			log.Println(err)
			logger.Message = r.Error
			logger.Type = ErrorType
		}

		returnChannel <- usrId
	}

	returnChannel <- "update done"
}

func (w *Worker) Start() error {
	done := make(chan bool)

	workers, err := strconv.Atoi(os.Getenv("WORKERS"))
	if err != nil {
		return err
	}

	if workers > 0 {
		w.Workers = workers
	}
	w.WorkerRunner(done)

	return nil
}
