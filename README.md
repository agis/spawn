spawn
===============

A utility package that makes writing end-to-end tests for Go servers easier.

It essentially spawns your program from its [`TestMain`](https://golang.org/pkg/testing/#hdr-Main), so that you can interact with it from your tests.

For usage instructions refer to the [examples](examples/).

Rationale
--------------
Writing an end-to-end test for a server in Go typically involves:

1) compiling the server
2) spinning up the binary
3) communicating with it somehow from the tests
4) shut the server down

This package makes the above steps easy to do right from within the Go tests of the server.
