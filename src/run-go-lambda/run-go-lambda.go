package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/rpc"
	"os"
	"time"

	"github.com/cenkalti/backoff"

	"github.com/aws/aws-lambda-go/lambda/messages"
)

func main() {
	args := &messages.InvokeRequest{
		Payload:            readPayload(),
		RequestId:          "1",
		XAmznTraceId:       "1",
		Deadline:           messages.InvokeRequest_Timestamp{Seconds: 300, Nanos: 0},
		InvokedFunctionArn: "arn:aws:lambda:an-antarctica-1:123456789100:function:test",
	}

	client := connect()

	var response *messages.InvokeResponse
	err := client.Call("Function.Invoke", args, &response)
	if err != nil {
		log.Println("Invocation:", err)
		log.Fatal("Response:", response)
	}
}

func readPayload() []byte {
	payload, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}
	return payload
}

func connect() *rpc.Client {
	port := os.Getenv("_LAMBDA_SERVER_PORT")
	serverAddress := fmt.Sprintf("localhost:%s", port)
	log.Println("Test harness connecting to: " + serverAddress)

	var client *rpc.Client
	connect := func() error {
		var err error
		client, err = rpc.Dial("tcp", serverAddress)
		if err != nil {
			return err
		}
		return nil
	}
	err := backoff.Retry(connect, constantBackoff())
	if err != nil {
		log.Fatal(err)
	}
	return client
}

func constantBackoff() *backoff.ExponentialBackOff {
	algorithm := backoff.NewExponentialBackOff()
	algorithm.MaxElapsedTime = 8 * time.Second
	algorithm.Multiplier = 1
	return algorithm
}