package serve

import (
	"bufio"
	"coffee-machine-assign/internal/models"
	"fmt"
	"os"
	"strconv"
	"strings"
	"sync"
)

// CoffeeMachineCLI A CLI is one of the methods of interacting with a coffee
// machine. Others may include HTTP, gRPC, GUI etc.
type CoffeeMachineCLI struct {
	cm *models.CoffeeMachine
	// Semaphore to prevent more servings of drinks than outlets.
	sem chan int
}

// NewCoffeeMachineCLI Returns a new instance of a coffee machine CLI
// application.
func NewCoffeeMachineCLI(outlets uint, maxCap uint) CoffeeMachineCLI {
	return CoffeeMachineCLI{
		cm: models.NewCoffeeMachine(outlets, maxCap),
		// sem is a buffered channel
		sem: make(chan int, outlets),
	}
}

// Start Set up the coffee machine, and other setup is done (if required)
func (cli *CoffeeMachineCLI) Start() {
	cli.cm.Start()
	fmt.Printf("Welcome to Chai Point!\n\n")
}

// Stop Tear down the CLI app and coffee machine, saving the coffee machine's
// state to a file.
func (cli *CoffeeMachineCLI) Stop() {
	fmt.Println("\nStopping the machine...")
	fmt.Println("Saving current state...")
	cli.cm.Stop()
}

// Run The main method in the CLI. It loops continuously and servers the user
// requests.
func (cli *CoffeeMachineCLI) Run() {
	// Make a channel for notifications from the coffee machine (de-clutters Stdin)
	c := make(chan string)
	// Create a waiting group to make sure all requests are served before the machine
	// is shut down.
	wg := sync.WaitGroup{}

	// Flag variable
	run := true

	for run {
		cli.ListOptions()

		// List all the notifications
		fmt.Println("\nNotifications : ")
		listNotifications(c)

		cli.DisplayLowIngredients()

		opt := cli.SelectOption()

		err := cli.ServiceOption(&wg, c, opt)

		if err != nil {
			fmt.Printf("%s\n\n", err.Error())
		}

		// Stop running the app if option 8 is selected.
		if opt == 8 {
			run = false
			continue
		}

		// Waits for ENTER before showing the prompt again.
		fmt.Println("\nPress ENTER to continue...")
		bufio.NewReader(os.Stdin).ReadBytes('\n')
	}

	// Wait for all the waiting group goroutines to finish before closing the
	// channel.
	go func() {
		wg.Wait()
		close(c)
	}()

	// Dump out any remaining notifications.
	for s := range c {
		fmt.Print(s)
	}
}

// ListOptions List all the options the coffee machine supports
func (cli *CoffeeMachineCLI) ListOptions() {
	opts := cli.cm.ListOptions()
	for _, opt := range opts {
		fmt.Println(opt)
	}
}

// ListBeverages Produce a list of all beverages the machine can serve.
func (cli *CoffeeMachineCLI) ListBeverages() {
	fmt.Print("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n")
	m := cli.cm.ListBeverages()
	for k := range m {
		fmt.Println(strings.ToTitle(k))
	}
	fmt.Print("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n\n")
}

// ListIngredients Produce a list of all ingredients currently in the system,
// along with their quantity.
func (cli *CoffeeMachineCLI) ListIngredients() {
	fmt.Print("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n")
	m := cli.cm.Ingredients()
	for k, q := range m {
		fmt.Printf("Name: %s\tQuantity: %d\n", strings.ToTitle(k.Name), q)
	}
	fmt.Print("-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-\n\n")
}

// SelectOption Helper function to parse user input for selecting one of the
// listed options.
func (cli *CoffeeMachineCLI) SelectOption() uint {
	optLen := len(cli.cm.ListOptions())
	strInput := getStringInput("\nSelect an action : ")
	// Parse a number from string input.
	input, err := strconv.ParseInt(strInput, 10, 64)
	if err != nil {
		fmt.Println("\nInvalid option...")
	}

	// If the selected option is out of range, continue asking till correct input is provided.
	for input <= 0 || int(input) > optLen {
		fmt.Printf("Please select a valid option (1 - %d): ", optLen)
		strInput = getStringInput("")
		input, err = strconv.ParseInt(strInput, 10, 64)
		if err != nil {
			fmt.Println("\nInvalid option...")
		}
	}
	fmt.Println()
	return uint(input)
}

// ServiceOption A helper function to switch on the option selected by user and
// then call the appropriate function.
func (cli *CoffeeMachineCLI) ServiceOption(wg *sync.WaitGroup, c chan string, opt uint) error {
	switch opt {
	// List all beverages
	case 1:
		cli.ListBeverages()

	// List all ingredients
	case 2:
		cli.ListIngredients()

	// Serve a beverage to the user
	case 3:
		err := cli.ServeBeverage(wg, c)
		if err != nil {
			return err
		}

	// Add an ingredient to the machine
	case 4:
		ing := getStringInput("Enter ingredient name : ")
		ok := cli.cm.AddIngredient(ing)
		if !ok {
			return fmt.Errorf("ingredient already added")
		}

	// Remove an ingredient from the machine.
	case 5:
		ing := getStringInput("Enter ingredient name : ")
		ok := cli.cm.RemoveIngredient(ing)
		if !ok {
			return fmt.Errorf("ingredient not in machine")
		}

	// Refill the given ingredient.
	case 6:
		ing := getStringInput("Enter ingredient name : ")
		ok := cli.cm.RefillIngredient(ing)
		if !ok {
			return fmt.Errorf("ingredient not in machine. please add")
		}

	// Refill all the ingredients in the machine.
	case 7:
		cli.cm.RefillAllIngredients()
	}
	return nil
}

// ServeBeverage Method responsible for asking user for a beverage name,
// validating it, and serving it.
func (cli *CoffeeMachineCLI) ServeBeverage(wg *sync.WaitGroup, c chan string) error {
	// List all the beverages to pick from.
	beverages := cli.cm.ListBeverages()
	cli.ListBeverages()

	// Get the beverage name, and check if it is available or not.
	b := getStringInput("Select a beverage : ")
	b = strings.ToLower(b)
	_, ok := beverages[b]
	for !ok {
		b = getStringInput("Please select a listed beverage : ")
		_, ok = beverages[b]
	}

	// Check whether we have enough ingredients to serve the beverage, if not return
	err := cli.cm.CheckIngredients(beverages[b])
	if err != nil {
		fmt.Println(err.Error())
		return nil
	}

	// Add to the waiting group.
	wg.Add(1)
	// Wait to acquire the semaphore
	cli.sem <- 1
	// Fire off a goroutine to serve the drink. and then decrease the waiting group.
	go func() {
		cli.cm.PourBeverage(c, beverages[b], true)
		wg.Done()
	}()
	// Release the semaphore.
	<-cli.sem

	return nil
}

// DisplayLowIngredients Displays low ingredients in the machine. Ingredient is
// counted as low if it has 10% or less quantity remaining in the machine.
func (cli *CoffeeMachineCLI) DisplayLowIngredients() {
	igs := cli.cm.Ingredients()
	low := make(map[models.Ingredient]uint, 0)
	for i, q := range igs {
		if q <= 100 {
			low[i] = q
		}
	}
	// Return prematurely if no ingredients are low.
	if len(low) == 0 {
		return
	}
	fmt.Print("Ingredients running low: ")
	for l := range low {
		fmt.Printf("\t%s", strings.ToTitle(l.Name))
	}
}

// Helper function to list the notifications in the channel.
func listNotifications(c chan string) {
	select {
	case s, ok := <-c:
		if ok {
			fmt.Println(s)
			// Recurse till the channel is empty.
			listNotifications(c)
		}
	default:
		fmt.Printf("No new notifications\n")
	}
}

// Helper function to parse user input from the os.Stdin, and displaying the
// given message.
func getStringInput(msg string) string {
	fmt.Print(msg)
	// Create a new scanner
	scanner := bufio.NewScanner(os.Stdin)
	// Scan stdin
	scanner.Scan()
	// Convert to string
	input := scanner.Text()
	// Return as lower cased.
	return strings.ToLower(input)
}
