# `cmdr` - Source protocol RCON CLI Tool

`cmdr` is a simple to use command-line interface tool written in Golang that allows you to execute commands directly on your server, either as a single command or through the interactive shell.

---

Now available on the AUR: `yay install cmdr` (https://aur.archlinux.org/packages/cmdr)
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
cmdr [OPTIONS] [COMMANDS]

```

### Flags 

`--addr <server_address:port>` (optional): The IP address and port of your RCON server (eg: 192.168.1.10:27015) if not using the default.

`--pwd <password>` (optional): Your rcon password for the server if not using the default.

`--cmd <command>` (optional): a single RCON command to execute on startup.

`--it` (optional): enable interactive shell mode.

`--cfg` (optional): configuration file to use. Defaults to the one specified in the make install if not provided.

`--env` (optional): which environment from configuration file to use, will use 'default' from the config file if not specified.


### Examples

Executing a single command:

`cmdr --addr 127.0.0.1:27015 --pwd your_rcon_password --cmd "status"`

Entering interactive shell mode with a start command included:

`cmdr --addr 127.0.0.1:27015 --pwd your_rcon_password --cmd "status" --it`

Entering an interactive shell from the minecraft environment, specified in the default config:

`cmdr --env minecraft -it`

Executing a command from a custom configuration, with a different minecraft environment: 

`cmdr --cfg /home/jdoe/servers.yaml --env minecraft --cmd "say hello"`

A config file must match this structure, with each key being the environment's name. An example file is also provided in the repo.

Example yaml file:
```yaml
default:
    address: "192.168.10.1:25565"
    password: "password"
minecraft:
    address: "example-rcon.com:448"
    password: "admin"
```


### Interactive mode

Once in interactive mode, you'll see a `cmd>` prompt. Type your RCON commands and press enter.
To exit the interactive shell, type `exit` or `quit`.
