package main

import (
	"context"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/agis/spawn"
)

func TestMain(m *testing.M) {
	// start the server
	server := spawn.New(main)
	ctx, cancel := context.WithCancel(context.Background())
	server.Start(ctx)

	// wait a bit for it to become ready
	time.Sleep(500 * time.Millisecond)

	// execute the test suite
	result := m.Run()

	// cleanly shutdown server
	cancel()
	err := server.Wait()
	if err != nil {
		log.Fatal(err)
	}

	os.Exit(result)
}

func TestFoo(t *testing.T) {
	res, err := http.Get("http://localhost:8080/foo")
	if err != nil {
		t.Fatal(err)
	}
	defer res.Body.Close()

	resBody, err := ioutil.ReadAll(res.Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(resBody) != "Hello!" {
		t.Fatalf("expected response to be 'Hello!', got '%s'", resBody)
	}
}
