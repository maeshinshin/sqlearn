package main

import (
	"flag"
	"os"
	"strings"

	"github.com/maeshinshin/sqlearn/backend/internal/handler"
	"github.com/maeshinshin/sqlearn/backend/internal/logger"
	"github.com/maeshinshin/sqlearn/backend/internal/server"
	"github.com/maeshinshin/sqlearn/backend/problems"
)

func main() {
	var debug bool
	var port string
	flag.BoolVar(&debug, "d", false, "Enable debug mode (shorthand)")
	flag.BoolVar(&debug, "debug", false, "Enable debug mode")
	flag.StringVar(&port, "p", "8080", "Port to run the server on (shorthand)")
	flag.StringVar(&port, "port", "8080", "Port to run the server on")

	flag.Parse()

	debugEnv := strings.ToLower(os.Getenv("DEBUG"))
	debug = debug || debugEnv == "true" || debugEnv == "1"

	logger.InitLogger(debug)

	problemHandler, err := handler.NewProblemHandler(problems.FS)
	if err != nil {
		panic(err)
	}

	svr := server.NewServer(
		server.WithDebug(debug),
		server.WithPort(port),
		server.WithHandler(problemHandler),
	)

	svr.Start()
}
