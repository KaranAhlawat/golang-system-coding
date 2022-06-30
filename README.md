# Coffee Machine - System Coding

This Golang project implements a coffee machine system which serves a beverage from the listed beverage to the user.
The project currently has a CLI interface for interacting with an instance of the coffee machine.

## Running the project
After changing directories to the one where the project is, type in the following commands

`go build -o build/coffee ./cmd/main.go`

and then run

`./build/coffee`

This will start a menu-driven CLI program for the user to interact with. 

On exit, the current state of the coffee machine is saved to a Go binary file, located in the `data` directory.

## Functionality
The coffee machine CLI has the following functionality the user can interactively trigger

1. List beverages
2. List ingredients and their quantities
3. Pour a beverage
4. Add an ingredient
5. Remove an ingredient
6. Refill an ingredient
7. Refill all ingredients
8. Stop the machine (saves state to binary file)

In addition, the CLI has input validation, and only proceeds if the input is correct, and the resource exists (if required).

The machine can serve upto N drinks at the time, where N is the number of outlets specified at the time of creating the machine instance in the `main.go` file. This is done through the user of goroutines and channels in go to achieve concurrent behaviour and semaphores-like behavior.

## Testing

To run the test suite for the project, type in the following command

`go test -v ./internal/models`

This will run all the test files for the project, giving a verbose output as to which tests were run, how many passed and how many failed, and how long each test took.
