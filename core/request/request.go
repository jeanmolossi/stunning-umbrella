package request

import (
	"errors"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strings"
)

type Requester struct {
	UsrId   string
	Error   string
	Client  *http.Client
	baseURL string
}

func init() {
	err := godotenv.Load()
	if err != nil {
		log.Fatalf("Erro ao carregar .env file")
	}
}

func NewRequester(usrId string, client *http.Client) *Requester {
	return &Requester{
		UsrId:   usrId,
		Client:  client,
		baseURL: os.Getenv("BASE_URL"),
	}
}

func (r *Requester) getUrl() string {
	return fmt.Sprintf("%s/api/recruiter/%s/access-level", r.baseURL, r.UsrId)
}

func (r *Requester) ping() (string, string, string) {
	return fmt.Sprintf("%s/ping", r.baseURL), "GET", ""
}

func (r *Requester) mountRequest() (*http.Request, error) {
	url := r.getUrl()
	method := "PUT"
	bodyJson := `{ "nivelAcesso": "admin" }`
	// url, method, bodyJson = r.ping()
	payload := strings.NewReader(bodyJson)

	request, err := http.NewRequest(method, url, payload)
	if err != nil {
		r.Error = fmt.Sprintf("Erro ao criar request: [%s]", r.UsrId)
		return nil, err
	}

	Cookie := fmt.Sprintf("")

	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Accept", "application/json")
	if Cookie != "" {
		request.Header.Add("Cookie", Cookie)
	}

	return request, nil
}

func (r *Requester) DoUpdate() error {
	request, err := r.mountRequest()
	if err != nil && r.Error != "" {
		r.Error = err.Error()
		log.Fatalf(err.Error())
		return err
	}

	response, err := r.Client.Do(request)
	if err != nil {
		r.Error = fmt.Sprintf("Erro ao executar a request USR_ID [ %s ] -  %v", r.UsrId, err.Error())
		return err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		errMsg := fmt.Sprintf("Erro com a request do USR_ID [ %s ] - %s", r.UsrId, response.Status)
		r.Error = errMsg
		return errors.New(errMsg)
	}

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		r.Error = fmt.Sprintf("Erro ao ler o corpo da resposta. USR_ID [%s]:: %v", r.UsrId, err.Error())
		return err
	}

	if len(string(body)) > 1000 {
		errorMsg := fmt.Sprintf("A resposta retornou 200 mas o corpo da resposta parece incorreto: [ %s ]", r.UsrId)
		r.Error = errorMsg
		return errors.New(errorMsg)
	}

	if r.Error == "" {
		log.Printf("Sucesso, USR_ID [ %s ] processado: %v", r.UsrId, string(body))
	}

	return nil
}
