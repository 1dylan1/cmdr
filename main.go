package main

import (
	"bufio"
	"commander/rcon"
	"flag"
	"fmt"
	"os"
	"strings"
)

func main() {

	serverAddr := flag.String("addr", "", "RCON server address/port (eg: 192.168.1.10:27015)")
	password := flag.String("pwd", "", "RCON password")
	command := flag.String("cmd", "", "RCON command to execute on program startup (eg: status)")
	interactive := flag.Bool("it", false, "Enable interactive shell mode")
	flag.Parse()

	if *password == "" {
		fmt.Printf("Error: RCON password cannot be empty: Use --pwd <password> | type %s --help for usage.\n", os.Args[0])
		os.Exit(1)
	}
	if *serverAddr == "" {
		fmt.Printf("Error: RCON address cannot be empty: Use --addr <server address> | type %s --help for usage\n", os.Args[0])
		os.Exit(1)
	}
	
	rconClient, err := rcon.NewRconClient(*serverAddr, *password)
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
				fmt.Println(output)
			}
		}

		fmt.Println("Entering interactive RCON shell. Type 'exit' or 'quit' to leave.")
		scanner := bufio.NewScanner(os.Stdin)
		for {
			fmt.Print("cmd> ")
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
