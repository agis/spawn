spawn
===============

Spawn makes it super-easy to spin up your Go server inline, inside its test
suite, for end-to-end testing.

It spawns your program from its [`TestMain`](https://golang.org/pkg/testing/#hdr-Main),
so that you can interact with it from your tests.


Usage
--------------
An example usage for [this simple HTTP server](examples/main.go) can be found below.
The complete runnable example is at [examples](examples/).

```go
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
```

Rationale
--------------
Writing an end-to-end test for a server typically involves:

1) compiling the server code
2) spinning up the binary
3) communicating with it from the tests
4) shutting the server down
5) verify everything went OK (server was closed cleanly etc.)

This package makes this process easy to do from within the tests of the server.
