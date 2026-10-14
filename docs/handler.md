# HANDLER Package

After decoding the byte stream in the `RESP` package, the `HANDLER` package parses the stream to execute certain commands.

## Architecture

* **`command.go`**: The entry point of the package
  * `ParseCommand()` takes in a struct `Value` (defined in the `RESP` package) and parses the incoming command into a struct `Command` with name and arguments. 
    * The command must be an array of bulk strings, where the first element is the command name and the remaining elements are the arguments.

* **`dispatch.go`**: 
  * `ExecuteCommand()` receives the parsed command and arguments from `ParseCommand()`, and executes said command. 
  * Uses a hashmap with string keys and function values for direct lookup.

* **`ping.go`**: 
  * Contains the individual functions to handle specific commands. 
  * Includes: 
    * `handlePing()` 
    * `handleEcho()` 
    * `handleSet()` 
    * `handleGet()`

## Available commands

- `PING`: Tests server availability.
- `ECHO`: Tests server transmission with a specific message.
- `SET`: Add value to key-value store.
- `GET`: Retrieves value from key-value store.
