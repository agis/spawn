package main

import (
	"context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/agis/spawn"
)

func TestMain(m *testing.M) {
	cmd := spawn.New(main)

	// start server
	ctx, cancel := context.WithCancel(context.Background())
	err := cmd.Start(ctx)
	if err != nil {
		log.Fatal(err)
	}

	// make sure server is up before running the tests
	for i := 0; i < 3; i++ {
		conn, err := net.Dial("tcp", ":8080")
		if err != nil {
			time.Sleep(500 * time.Millisecond)
			continue
		}
		conn.Close()
	}

	result := m.Run()

	// shutdown server
	cancel()
	err = cmd.Wait()
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

	resBody, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		t.Fatal(err)
	}

	if string(resBody) != "Hello!" {
		t.Fatalf("expected response to be 'Hello!', got '%s'", resBody)
	}
}
