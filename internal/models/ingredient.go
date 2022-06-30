package models

// Ingredient Represents an ingredient present in the machine. Ingredients make up ListBeverages in the machine.
type Ingredient struct {
	Name string
}

// String Overriding the default to string method to just return the name of the ingredient
func (i Ingredient) String() string {
	return i.Name
}
