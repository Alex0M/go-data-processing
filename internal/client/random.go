package client

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"time"
)

type jsonContent struct {
	States  []string `json:"states"`
	Devices []string `json:"devices"`
	Streams []string `json:"streams"`
}

type pRand struct {
	pr *rand.Rand
}

func (r *pRand) Intn(n int) int {
	return r.pr.Intn(n)
}

var data = jsonContent{}
var privateRand *pRand

func init() {
	privateRand = &pRand{rand.New(rand.NewSource(time.Now().UnixNano()))}
}

func ReadConfigRandomData(f string) error {
	fData, err := os.ReadFile(f)
	if err != nil {
		return fmt.Errorf("error reading random config data from file: %s", err)
	}

	err = json.Unmarshal(fData, &data)
	if err != nil {
		return fmt.Errorf("error reading json random config data: %s", err)
	}

	return nil
}

func randomFrom(source []string) string {
	if len(source) == 0 {
		panic("source should no be nil in randomFrom function. Please check config file and make sure data is correct")
	}
	return source[privateRand.Intn(len(source))]
}

func RandomState() string {
	return randomFrom(data.States)
}

func RandomDevice() string {
	return randomFrom(data.Devices)
}

func RandomStream() string {
	return randomFrom(data.Streams)
}
