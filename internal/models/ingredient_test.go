package models

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestIngredient_String(t *testing.T) {
	ingredient := &Ingredient{
		Name: "coffee",
	}
	require.EqualValues(t, ingredient.Name, ingredient.String())
}
