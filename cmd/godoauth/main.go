package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/n1tr0g/godoauth"
)

var (
	name    string = "godoauth - Go Docker Token Authenticator"
	version string = "v0.0.1"
	commit  string
)

var (
	confFile    string
	showVersion bool
)

func init() {
	flag.StringVar(&confFile, "config", "config.yaml", "Go Docker Token Auth Config file")
	flag.BoolVar(&showVersion, "version", false, "show the version and exit")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage of Go Docker Token Auth (version %v):\n", version)
		flag.PrintDefaults()
	}
}

func main() {
	flag.Parse()

	if showVersion {
		fmt.Fprintln(os.Stderr, os.Args[0], version)
		return
	}

	var config godoauth.Config
	if err := config.LoadFromFile(confFile); err != nil {
		fmt.Fprintln(os.Stderr, "error parsing config file: ", err)
		os.Exit(1)
	}

	if err := config.LoadCerts(); err != nil {
		fmt.Fprintln(os.Stderr, "error while loading/veryfing certs: ", err)
	}

	fmt.Printf("Starting %s version: %s\n", name, version)

	server, err := godoauth.NewServer(&config)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error while creating new server: ", err)
		os.Exit(1)
	}
	server.Start()

	// waiting for a termination signal to clean up
	interruptChan := make(chan os.Signal)
	signal.Notify(interruptChan, os.Interrupt, syscall.SIGTERM, syscall.SIGHUP)
	<-interruptChan
}
