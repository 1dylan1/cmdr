# `cmdr` - Source protocol RCON CLI Tool

`cmdr` is a simple to use command-line interface tool written in Golang that allows you to execute commands directly on your server, either as a single command or through the interactive shell.

---

## Features

**Single command execution** - run a specific RCON command directory from the terminal and get the output
**Interactive shell** - enter an interactive session to send/view multiple RCON commands without reconnecting


## Installation

To install `cmdr` to your `/usr/local/bin`, use the Makefile provided, and run `make install`. Once installed, you can run `cmdr` from any directory.

Alternatively, you can build the file with `make build` and, run `./cmdr`.

---

## Usage

`cmdr` offers two primary modes of operation: single command execution and an interactive shell. You can access the shell by adding the `--it` flag to your command.

### Basic Syntax

```bash
cmdr --addr <server_ip:port> --pwd <rcon_password [OPTIONS]

```

### Flags 

`--addr <server_address:port>` (required): The IP address and port of your RCON server (eg: 192.168.1.10:27015).
`--pwd <password>` (required): Your rcon password for the server.
`--cmd <command>` (optional): a single RCON command to execute on startup.
`--it` (optional): enable interactive shell mode.


### Examples

Executing a single command:
`cmdr --addr 127.0.0.1:27015 --pwd your_rcon_password --cmd "status"`

Entering interactive shell mode with a start command included:
`cmdr --addr 127.0.0.1:27015 --pwd your_rcon_password --cmd "status" --it`

### Interactive mode

Once in interactive mode, you'll see a `cmd>` prompt. Type your RCON commands and press enter.
To exit the interactive shell, type `exit` or `quit`.
