package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func newMockMachine() *CoffeeMachine {
	machine := NewCoffeeMachine(5, 1000)
	machine.SeedData()
	return machine
}

var cm *CoffeeMachine = newMockMachine()

func TestCoffeeMachine_ListOptions(t *testing.T) {
	opts := cm.ListOptions()
	require.Equal(t, 8, len(opts))
	for _, opt := range opts {
		require.NotEmpty(t, opt)
	}
}

func TestCoffeeMachine_ListBeverages(t *testing.T) {
	opts := cm.ListBeverages()
	require.Equal(t, 5, len(opts))
	for _, opt := range opts {
		require.NotEmpty(t, opt)
	}
}

func TestCoffeeMachine_Ingredients(t *testing.T) {
	opts := cm.Ingredients()
	for _, opt := range opts {
		require.NotEmpty(t, opt)
	}
}

func TestCoffeeMachine_CheckIngredientsEnough(t *testing.T) {
	b := cm.Beverages["coffee"]
	err := cm.CheckIngredients(b)
	require.NoError(t, err)
}

func TestCoffeeMachine_CheckIngredientsNotEnough(t *testing.T) {
	b := cm.Beverages["coffee"]
	b.ingredients[Ingredient{Name: "sugar syrup"}] += 1000
	err := cm.CheckIngredients(b)
	require.Error(t, err)
}

func TestCoffeeMachine_PourBeverage(t *testing.T) {
	c := make(chan string)
	b := cm.Beverages["coffee"]
	oldIngredients := cm.Ingredients()
	go cm.PourBeverage(c, b, false)
	newIngredients := cm.Ingredients()
	for i := range oldIngredients {
		require.LessOrEqual(t, oldIngredients[i], newIngredients[i])
	}
}

func TestCoffeeMachine_AddIngredient(t *testing.T) {
	_, found := cm.Ingredients()[Ingredient{"olive oil"}]
	require.False(t, found)

	added := cm.AddIngredient("olive oil")
	val, found := cm.Ingredients()[Ingredient{"olive oil"}]

	require.True(t, added)
	require.True(t, found)
	require.Equal(t, uint(1000), val)
}

func TestCoffeeMachine_RemoveIngredient(t *testing.T) {
	_, found := cm.Ingredients()[Ingredient{Name: "coffee syrup"}]
	require.True(t, found)

	removed := cm.RemoveIngredient("coffee syrup")
	_, found = cm.Ingredients()[Ingredient{Name: "coffee syrup"}]

	require.True(t, removed)
	require.False(t, found)
}

func TestCoffeeMachine_RefillIngredient(t *testing.T) {
	cm.Ingredient[Ingredient{Name: "coffee syrup"}] = 50
	require.Equal(t, uint(50), cm.Ingredient[Ingredient{"coffee syrup"}])
	cm.RefillIngredient("coffee syrup")

	require.Equal(t, uint(1000), cm.Ingredient[Ingredient{"coffee syrup"}])
}
