package models

import (
	"encoding/gob"
	"fmt"
	"os"
	"strings"
	"time"
)

// CoffeeMachine Represents a coffee machine. It has N outlets for serving
// ListBeverages in parallel. It has a list of ListBeverages (immutable) that can be
// served from it. It has a list of Ingredient (mutable) for making those ListBeverages.
type CoffeeMachine struct {
	// NumOutlets The number of outlets that the machine has for serving ListBeverages.
	NumOutlets uint
	Beverages  map[string]*Beverage
	Ingredient map[Ingredient]uint
	// MaxCapacity Indicates the maximum quantity of each ingredient the machine can
	// hold. Can be set while instantiating the machine.
	MaxCapacity uint
}

// NewCoffeeMachine Instantiate a new Coffee machine instance with n outlets and
// maxIngredientCap as the maximum capacity of the machine to store an
// ingredient. Returns a pointer to a machine.
func NewCoffeeMachine(n uint, maxIngredientCap uint) *CoffeeMachine {
	cm := &CoffeeMachine{
		NumOutlets:  n,
		Beverages:   make(map[string]*Beverage, 0),
		Ingredient:  make(map[Ingredient]uint, 0),
		MaxCapacity: maxIngredientCap,
	}
	return cm
}

// SeedData Seed the machine with some dummy data.
func (cm *CoffeeMachine) SeedData() {
	cm.Beverages["ginger tea"] = &Beverage{
		Name: "ginger tea",
		ingredients: map[Ingredient]uint{
			{Name: "hot water"}:        50,
			{Name: "hot milk"}:         10,
			{Name: "tea leaves syrup"}: 10,
			{Name: "ginger syrup"}:     5,
			{Name: "sugar syrup"}:      10,
		},
	}

	cm.Beverages["elaichi tea"] = &Beverage{
		Name: "elaichi tea",
		ingredients: map[Ingredient]uint{
			{Name: "hot water"}:        50,
			{Name: "hot milk"}:         10,
			{Name: "tea leaves syrup"}: 10,
			{Name: "elaichi syrup"}:    5,
			{Name: "sugar syrup"}:      10,
		},
	}

	cm.Beverages["coffee"] = &Beverage{
		Name: "coffee",
		ingredients: map[Ingredient]uint{
			{Name: "hot water"}:    50,
			{Name: "hot milk"}:     10,
			{Name: "coffee syrup"}: 10,
			{Name: "sugar syrup"}:  10,
		},
	}

	cm.Beverages["hot milk"] = &Beverage{
		Name: "hot milk",
		ingredients: map[Ingredient]uint{
			{Name: "milk"}: 50,
		},
	}

	cm.Beverages["hot water"] = &Beverage{
		Name: "hot water",
		ingredients: map[Ingredient]uint{
			{Name: "water"}: 50,
		},
	}

	for _, b := range cm.Beverages {
		for i, _ := range b.ingredients {
			cm.Ingredient[i] = cm.MaxCapacity
		}
	}
}

// Start Seed the machine, and other setup (if required)
func (cm *CoffeeMachine) Start() {
	cm.SeedData()
}

// Stop Save the last machine state to a Golang binary (GOB) file, and then exit
// gracefully.
func (cm *CoffeeMachine) Stop() {
	file, err := os.Create("./data/machine.gob")
	if err != nil {
		panic(err)
	}
	defer func(file *os.File) {
		err = file.Close()
		if err != nil {
			fmt.Println("Error closing the machine state file.")
			return
		}
	}(file)

	enc := gob.NewEncoder(file)
	err = enc.Encode(*cm)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Printf("\nUnable to save machine state\n")
		return
	}
}

// ListOptions Returns an array of options that the machine can perform.
func (cm *CoffeeMachine) ListOptions() []string {
	return []string{
		"1. List ListBeverages",
		"2. List Ingredient",
		"3. Pour a beverage",
		"4. Add an ingredient",
		"5. Remove an ingredient",
		"6. Refill an ingredient",
		"7. Refill all Ingredient",
		"8. Stop coffee machine",
	}
}

// PourBeverage When an outlet is free, pour a beverage asynchronously. The
// notifications are sent to the channel c. It is assumed that the machine takes
// 5 second to pour one drink.
func (cm *CoffeeMachine) PourBeverage(c chan<- string, b *Beverage, delay bool) {
	// Take the Ingredient from the machine's ingredient container.
	for i, q := range b.ingredients {
		cm.Ingredient[i] -= q
	}

	// Simulate pouring a drink.
	c <- fmt.Sprintf("Serving %s", b.Name)
	if delay {
		time.Sleep(time.Duration(5 * time.Second))
	}
	c <- fmt.Sprintf("Served %s", b.Name)
}

// CheckIngredients Check the availability of Ingredient for a given beverage.
// If enough Ingredient are not available, returns false.
func (cm *CoffeeMachine) CheckIngredients(b *Beverage) error {
	for i, q := range b.ingredients {
		if cm.Ingredient[i] < q {
			return fmt.Errorf("Not enough %s to make %s\n", i.Name, b.Name)
		}
	}
	return nil
}

// ListBeverages Returns a collection of all the ListBeverages in the machine.
func (cm *CoffeeMachine) ListBeverages() map[string]*Beverage {
	return cm.Beverages
}

// Ingredients Returns a collection of all the Ingredient in the system and
// their present quantity.
func (cm *CoffeeMachine) Ingredients() map[Ingredient]uint {
	return cm.Ingredient
}

// AddIngredient Add an ingredient to the machine. Constructs the ingredient
// struct, and sets the quantity to be the max specified of the machine.
// Returns true if successfully added, false if already exists.
func (cm *CoffeeMachine) AddIngredient(s string) bool {
	// Normalize and create the name
	i := Ingredient{strings.ToLower(s)}
	// Return false if exists
	if _, ok := cm.Ingredient[i]; ok {
		return false
	}
	cm.Ingredient[i] = cm.MaxCapacity
	return true
}

// RemoveIngredient Remove a given ingredient from the machine entirely. Returns
// true if removed, false if doesn't exist in the machine.
func (cm *CoffeeMachine) RemoveIngredient(s string) bool {
	i := Ingredient{strings.ToLower(s)}
	if _, ok := cm.Ingredient[i]; !ok {
		return false
	}
	delete(cm.Ingredient, i)
	return true
}

// RefillIngredient Refills the specified ingredient to it's maximum capacity.
// Returns true if refilled, false if it does not exist in the machine.
func (cm *CoffeeMachine) RefillIngredient(s string) bool {
	i := Ingredient{strings.ToLower(s)}
	if _, ok := cm.Ingredient[i]; !ok {
		return false
	}
	cm.Ingredient[i] = cm.MaxCapacity
	return true
}

// RefillAllIngredients Refills all the Ingredient currently in the machine to
// their maximum capacity.
func (cm *CoffeeMachine) RefillAllIngredients() {
	for i := range cm.Ingredient {
		cm.Ingredient[i] = cm.MaxCapacity
	}
}
