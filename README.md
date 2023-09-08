# machogen

This project allows you to dynamically compile a mach-o binary that performs a sequence of user-defined commands.
The purpose of this project is to allow safe replication of malicious behavior for validating detection logic without
needing to execute malware directly.

## Usage

### Building

- `go build machogen.go`

### Execution

This can be executed in two ways:

1. Specify a comma-separated list of commands
- `./machogen -commands "echo hello, echo world"`

2. Read from a JSON that contains an array of commands (see `commands.json` for an example)
- `./machogen -json commands.json`


## Compiled Binary

The generated binary can be executed via `./generated_binary`.
