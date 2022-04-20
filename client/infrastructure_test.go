package client

import (
	"fmt"
	"github.com/bobo-boom/dcrnsugar/config"
	"net/http"
	"testing"
)

func TestClient_SendRequest(t *testing.T) {

	config := &config.DefaultConfig

	client, err := New(config)
	if err != nil {
		fmt.Println(err)
	}
	address := "Dsj9YNJQBhwan9KPqf2Ki4i8943g6yToTzj"
	client.GenerateGetBalanceUrl(address)
	request, err := http.NewRequest("GET", address, nil)
	if err != nil {
		fmt.Println(err)
	}
	details := &RequestDetails{
		HttpRequest: request,
	}
	for i := 0; i < 10000; i++ {
		respData, err := client.SendRequest(details)
		if err != nil {
			fmt.Printf("index  %d, err: %v", i, err)
		}

		fmt.Printf("index: %v data : %v", i, string(respData))
	}

}
func TestClient_String(t *testing.T) {
	config := &config.DefaultConfig

	client, err := New(config)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("url : ", client.String())
}

func TestClient_GenerateGetBalanceUrl(t *testing.T) {
	config := &config.DefaultConfig

	client, err := New(config)
	if err != nil {
		fmt.Println(err)
	}
	address := "Dsj9YNJQBhwan9KPqf2Ki4i8943g6yToTzj"
	fmt.Println("url : ", client.GenerateGetBalanceUrl(address))

}
