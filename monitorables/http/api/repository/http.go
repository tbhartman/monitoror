package repository

import (
	"bufio"
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/monitoror/monitoror/monitorables/http/api"
	"github.com/monitoror/monitoror/monitorables/http/api/models"
	"github.com/monitoror/monitoror/monitorables/http/config"
)

type (
	httpRepository struct {
		httpClient *http.Client
		header     http.Header
	}
)

func NewHTTPRepository(config *config.HTTP) api.Repository {
	var certificates []tls.Certificate

	if config.Certificate != "" && config.Key != "" {
		cert, error := tls.LoadX509KeyPair(config.Certificate, config.Key)

		if error == nil {
			certificates = append(certificates, cert)
		}
	}

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: !config.SSLVerify, Certificates: certificates},
	}
	client := &http.Client{Transport: tr, Timeout: time.Duration(config.Timeout) * time.Millisecond}

	var header = make(http.Header)
	var headerScanner = bufio.NewScanner(strings.NewReader(config.Header))
	for headerScanner.Scan() {
		key, value, ok := strings.Cut(headerScanner.Text(), ":")
		if ok {
			header.Add(key, value)
		}
	}

	return &httpRepository{client, header}
}

func (r *httpRepository) Get(url string) (response *models.Response, err error) {
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}
	req.Header = r.header.Clone()
	resp, err := r.httpClient.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	response = &models.Response{
		StatusCode: resp.StatusCode,
		Body:       bytes,
	}

	return
}
