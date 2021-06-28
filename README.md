# ginkin

Create command-line interface for Gin server APIs by mix Gin and Kingpin.

## Example

For example in `example/main.go`, one can use following command line interface to call API function
directly or start server by default:

- `go run main.go` will start http api server handle requests to "http://:3000/version".
- `go run main.go version` will print out version info directly.
- `go run main.go user/list` will print out user list.
- `go run main.go user/:user john` will print out user list.
- `go run main.go user/add "['jane']"` will add users.
- `go run main.go user/:user#del john` will del user.
