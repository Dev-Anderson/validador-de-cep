package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

type cepValid struct {
	Cep   string
	Valid string
}

func validCep(cep string) (bool, error) {
	var valid struct {
		Error bool `json:"erro"`
	}

	res, err := http.Get(fmt.Sprintf("https://viacep.com.br/ws/%v/json", cep))
	if err != nil {
		return false, fmt.Errorf("valid cep: %w", err)
	}

	defer res.Body.Close()

	b, err := io.ReadAll(res.Body)
	if err != nil {
		return false, err
	}

	if err := json.Unmarshal(b, &valid); err != nil {
		return false, fmt.Errorf("valid cep: %w", err)
	}

	return !valid.Error, nil
}

func validCeps(ceps []string) ([]cepValid, error) {
	for _, cep := range ceps {
		go func(cep string) {
			valid, err := validCep(cep)
			fmt.Println(valid, err, time.Since(time.Now()))
		}(cep)
	}
	fmt.Println("Loop rodado com sucesso!")

	return nil, nil
}

func main() {
	c, err := os.ReadFile("./cep.txt")
	if err != nil {
		log.Fatal(err)
	}

	ceps := strings.Split(string(c), "\n")
	_, err = validCeps(ceps)
	if err != nil {
		log.Fatal(err)
	}
	time.Sleep(time.Second * 10)

}
