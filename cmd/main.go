package main

import (
	"cf-observer/internal/config"
	"flag"
	"fmt"
	"log"
	"os"
)

func main() {
	if len(os.Args) < 2 {
		printRootUsage()
		os.Exit(1)
	}

	if os.Args[1] == "-h" || os.Args[1] == "--help" || os.Args[1] == "help" {
		printRootUsage()
		return
	}

	cmdArg := os.Args[1]
	switch cmdArg {
	// NOTE: Run the daemon in the background when the project is ready for deployment
	case "init":
		initCmd := flag.NewFlagSet("init", flag.ExitOnError)
		force := initCmd.Bool("force", false, "Overwrite existing config file")
		_ = initCmd.Parse(os.Args[2:])

		path, err := config.InitConfigDir(*force)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("Initialized configuration at %s\n", path)
	case "start":
		startCmd := flag.NewFlagSet("start", flag.ExitOnError)
		configFile := startCmd.String("config", "", "path to config file")
		_ = startCmd.Parse(os.Args[2:])

		conf, err := config.LoadConfigFile(*configFile)
		if err != nil {
			log.Fatal(err)
		}

		if err := conf.Validate(); err != nil {
			log.Fatal(err)
		}

		config.AppConfig = conf
		fmt.Println("observer started")
	case "stop":
	default:
		printRootUsage()
		os.Exit(1)
	}

}

func printRootUsage() {
	fmt.Println("Usage: cf-observer <command> [options]")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  init      Initialize the config directory and default configuration")
	fmt.Println("  start     Start the observer daemon")
	fmt.Println("  stop      Stop the observer daemon")
	fmt.Println()
	fmt.Println("Run 'cf-observer <command> -h' for command-specific help")
}
