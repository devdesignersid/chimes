# Chimes

Chimes is a command-line interface (CLI) application written in Go. It sends desktop notifications to remind users at the due time of each reminder.

## Features

- Set reminders with due times.
- Receive desktop notifications when reminders are due.

## Installation

To install Chimes, you need to have Go installed on your machine. Once you have Go installed, you can clone this repository and build the project.

## Building the Project

You can build the project using the provided Makefile. The Makefile includes the following commands:

- `make build`: Builds the project for Darwin, Linux, and Windows.
- `make run`: Runs the project.
- `make lint`: Lints the project using golangci-lint.
- `make clean`: Cleans up the project.
- `make dep`: Downloads the project dependencies.
- `make install-linter`: Installs golangci-lint.

## Usage

After building the project, you can run the application with the `make run` command. This will start the application and you can begin setting reminders.

## Commands

- `add`: Adds a reminder to the system. It takes in various parameters such as 'memo', 'priority', 'date', 'days', 'months', 'years', 'hours', 'minutes', 'seconds', 'repeat', and 'repeatInterval'.
- `list`: Lists all available reminders. It has two optional flags: 'was-due' to filter reminders that were due in the past, and 'will-be-due' to filter reminders that will be due in the future.
- `update`: Updates a reminder in the system. It takes in various parameters such as 'memo', 'priority', 'date', 'days', 'months', 'years', 'hours', 'minutes', 'seconds', 'repeat', and 'repeatInterval'.
- `stop`: Stops the reminder system.
- `delete`: Deletes a reminder from the system. It takes in the 'id' parameter which is the ID of the reminder to be deleted.

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## License

This project is licensed under the MIT License - see the LICENSE.md file for details.
