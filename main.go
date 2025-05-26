package main

import (
	"bufio"
	"commander/rcon"
	"commander/config"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	serverAddr := flag.String("addr", "", "RCON server address/port (eg: 192.168.1.10:27015)")
	password := flag.String("pwd", "", "RCON password")
	command := flag.String("cmd", "", "RCON command to execute on program startup (eg: status)")
	environment := flag.String("env", "", "Predefined environment from config file to use")
	interactive := flag.Bool("it", false, "Enable interactive shell mode")
	// configFile := flag.String("config", "", "Config file to use (defaults to your .config/cmdr/config file)")	
	flag.Parse()

	config, err := config.LoadConfig()
	if err != nil {
		fmt.Printf("Error getting config file: %v", err)
		os.Exit(1)
	}
	
	var address string
	var pwd string

	if defaultServer, found := config.Servers["default"]; found {
		address = defaultServer.Address
		pwd = defaultServer.Password
	}

	if *environment != "" {
		if server, found := config.Servers[*environment]; found {
			address = server.Address
			pwd = server.Password
		}
	}

	if *password != "" && *serverAddr != "" {
		address = *serverAddr
		pwd = *password
	}

	rconClient, err := rcon.NewRconClient(address, pwd)
	if err != nil {
		fmt.Printf("error opening rcon connection: %v\n", err)
		return
	}
	defer rconClient.Close()

	var output string
	
	if *interactive {
		if *command != "" {
			fmt.Printf("Executing command: %s\n", *command)
			output, err = rconClient.ExecuteCommand(*command)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", *command, err)
			} else {
				fmt.Print(output)
			}
		}

		fmt.Println("Entering interactive RCON shell. Type 'exit' or 'quit' to leave.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("cmdr> ")
			if !scanner.Scan() {
				break
			}
			input := scanner.Text()

			if input == "exit" || input == "quit" {
				fmt.Println("Exiting RCON shell.")
				break
			}
			if strings.TrimSpace(input) == "" {
				continue
			}

			output, err = rconClient.ExecuteCommand(input)
			if err != nil {
				fmt.Printf("Error executing command '%s': %v\n", input, err)
			} else {
				fmt.Println(output)
			}
		}

		if err = scanner.Err(); err != nil {
			fmt.Fprintf(os.Stderr, "Error reading input: %v\n", err)
		}
	} else {
		if *command == "" {
			fmt.Printf("Error: no command specified Use --cmd <rcon_command> or --it for interactive mode. Type %s --help for usage.\n", os.Args[0])
			os.Exit(1)
		}
		fmt.Printf("Executing single command: %s\n", *command)
		output, err = rconClient.ExecuteCommand(*command)
		if err != nil {
			fmt.Printf("Error executing command '%s': %v\n", *command, err)
			os.Exit(1)
		} else {
			fmt.Println(output)
		}
	}
}
