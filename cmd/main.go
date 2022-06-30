package main

import "coffee-machine-assign/internal/serve"

// main The entrypoint of the program. Interacts with the Coffee Machine CLI.
func main() {
	cli := serve.NewCoffeeMachineCLI(5, 1000)
	cli.Start()
	cli.Run()
	cli.Stop()
}
