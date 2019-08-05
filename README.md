spawn
===============

A utility package that makes writing end-to-end tests for Go servers easier.

It spawns your program from its [`TestMain`](https://golang.org/pkg/testing/#hdr-Main),
so that you can interact with it from your tests.

For usage instructions refer to the [examples](examples/) and the
[documentation](https://godoc.org/github.com/agis/spawn).

Rationale
--------------
Writing an end-to-end test for a server typically involves:

1) compiling the server code
2) spinning up the binary
3) communicating with it from the tests
4) shut the server down
5) verify everything went OK (server was closed cleanly etc.)

This package makes this process easy to do from within the tests of the server.
