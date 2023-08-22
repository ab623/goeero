# Eero (Unofficial Client)

This repository covers 2 main functions
1. A library to interact with eero devices
2. A CLI which makes use of the library


## Quick Start - Using the CLI

The CLI is capable of the following actions:

* `auth` - Authenticate the user and store
* `account` - Get a JSON output of the customer account details
* `devices` - Get a JSON output of the devices from all the networks
* `networks` - Get a JSON output of the networks that setup on the Eero
* `session` - Show the session token if it is available

The client accepts an optional flag `--session-file` or `-s` which accepts a path for where to read/write the eero session token. If a path is not provided, it will default to the application directory in a file called `eero_session.txt`.

### Authentication

Running `./eeroclient auth` will launch the authentication flow and save the authentication token at the path specified. 

Tha authentication flow is a 2 stage process:
1. CLI will ask for a mobile number, where an OTP will be sent. (Note: Ensure the number contains the +xx prefix if that is what is registed with eero)
2. CLI will ask for the OTP, received via email or SMS, to complete the authentication flow


## Quick Start - Using the Library

Run the following command

```
go mod add github.com/ab623/goeero/eero
```
Then utilise the library in your code using

```go
import "github.com/ab623/goeero/eero"
```

## Building

To build the CLI locally run the following command

```bash
git clone https://github.com/ab623/goeero.git
go build -o eeroclient cli/eeroclient/main.go
```


## Recipes

### Combining outputs with JQ

Output from commands can be passed into `jq` to manipulate the output. 

For example to get the list of device display names and associated IP's in host file format the following command can be used:

`./eeroclient devices | jq -r '.[] | select(.ip | length > 0) | [.ip, .display_name] | @tsv' | sort -k2`

This will produce a file such as 

```
192.168.1.100 Johns iPhone
192.168.1.120 Davids Laptop
```
