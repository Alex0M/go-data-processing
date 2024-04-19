package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"out/internal/models"
	"time"

	"github.com/google/uuid"
)

type Client struct {
	ConfigData   []*models.Stream
	Frequency    int
	ClientNumber int
	RandomMode   bool
	ReportingAPI string
	Token        string
}

func NewClient(num, freq int, data []*models.Stream, mode bool, uri, token string) *Client {
	return &Client{
		ConfigData:   data,
		Frequency:    freq,
		ClientNumber: num,
		RandomMode:   mode,
		ReportingAPI: uri,
		Token:        token,
	}
}

func (c *Client) Run() {
	for i := 0; i < c.ClientNumber; i++ {
		if c.RandomMode {
			go func() {
				s := c.setRandomStreamValue()
				c.sendHeardBeat(s)
			}()
		} else {
			go func(n int) {
				s := c.ConfigData[n]
				c.sendHeardBeat(s)
			}(i)
		}
	}
}

func ReadClientDataFromFile(f string) ([]*models.Stream, error) {
	data := []*models.Stream{}
	fData, err := os.ReadFile(f)
	if err != nil {
		return nil, fmt.Errorf("error reading client config data from file: %s", err)
	}

	err = json.Unmarshal(fData, &data)
	if err != nil {
		return nil, fmt.Errorf("error reading json client config data: %s", err)
	}

	return data, nil
}

func (c *Client) setRandomStreamValue() *models.Stream {
	if c.RandomMode {
		return &models.Stream{
			ClientID: fmt.Sprintf("user-%s", uuid.New()),
			Geo:      RandomState(),
			Content:  RandomStream(),
			Device:   RandomDevice(),
		}
	}
	return &models.Stream{}
}

func (c *Client) sendHeardBeat(s *models.Stream) {
	for {
		s.Timestamp = time.Now().Unix()
		jsonData, err := json.Marshal(s)
		if err != nil {
			panic(err)
		}
		c.sentPostRequest(jsonData)
		time.Sleep(time.Duration(c.Frequency) * time.Second)
	}
}

func (c *Client) sentPostRequest(data []byte) {
	httpClient := &http.Client{}
	req, err := http.NewRequest(http.MethodPost, c.ReportingAPI, bytes.NewBuffer(data))
	if err != nil {
		log.Printf("error creating request to tracking API: %s", err)
		return
	}
	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))

	resp, err := httpClient.Do(req)
	if err != nil {
		log.Printf("error sending data to tracking API: %s", err)
		return
	}
	defer resp.Body.Close()

	log.Printf("send data: %s", string(data))
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Printf("error reading response from tracking API: %s", err)
		return
	}

	log.Printf("tracking API response: %s", string(body))
}
