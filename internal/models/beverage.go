package models

// Beverage A composed entity in the system. It has a name, and a list of
// Ingredient required to make it along with their quantities.
type Beverage struct {
	Name string
	// Storing the ingredients as a map for faster lookups
	ingredients map[Ingredient]uint
}
