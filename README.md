spawn
===============

Spawn makes it easy to spin up your Go server right from within its own test suite, for end-to-end testing.

Usage
--------------
An example usage for [this simple HTTP server](examples/main.go) can be found below.
The complete runnable example is at [examples](examples/).

```go
func TestMain(m *testing.M) {
	// start the server on localhost:8080 (we assume it accepts a `--port` argument)
	server := spawn.New(main, "--port", "8080")
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

func TestServerFoo(t *testing.T) {
	res, _ := http.Get("http://localhost:8080/foo")
	defer res.Body.Close()

	resBody, _ := ioutil.ReadAll(res.Body)

	if string(resBody) != "Hello!" {
		t.Fatalf("expected response to be 'Hello!', got '%s'", resBody)
	}
}

// more tests using the server
```

Rationale
--------------
Writing an end-to-end test for a server typically involves:

1) compiling the server code
2) spinning up the binary
3) communicating with it from the tests
4) shutting the server down
5) verify everything went OK (server was closed cleanly etc.)

This package aims to simplify this process.
