package main

import (
	"errors"
	"fmt"
	"math"
	"strings"
)

// note the go.mod is used to support having the solution in a folder

func main() {
	// variable for testing support
	inputFile := "menu.json"
	findAndEmitBestMeal(inputFile)
}

func findAndEmitBestMeal(inputFile string) {
	bestMeals, err := findBestMeal(inputFile)

	if err != nil {
		emitBestMealError(err)
	}

	bestMeal := bestMeals[0]

	emitBestMeal(bestMeal)
}

func emitBestMealError(err error) {

	// write the error in json format
	bestMealError := BestMealError{err.Error()}
	fmt.Printf("{\n\t\"error\": \"%s\"\n}\n", bestMealError.error)
}

func emitBestMeal(bestMeal BestMeal) {

	// write the best meal in json format
	fmt.Printf("{\n"+
		"\t\"selectedFoods\": %s,\n"+
		"\t\"totalCost\": %d,\n"+
		"\t\"totalSatisfaction\": %d\n"+
		"}\n",
		fmt.Sprintf("[ \"%s\", \"%s\", \"%s\", \"%s\" ]",
			bestMeal.selectedFoods[APP_INDEX],
			bestMeal.selectedFoods[DRINK_INDEX],
			bestMeal.selectedFoods[MAIN_INDEX],
			bestMeal.selectedFoods[DESSERT_INDEX],
		),
		bestMeal.totalCost,
		bestMeal.totalSatisfaction,
	)

}

////////////////////////////////////////////////////////////////////
// define the output

type BestMeal struct {
	selectedFoods     [4]string
	totalCost         int
	totalSatisfaction int
}

/*
// a few error conditions
type BestMealErr int

const ( // BEST_MEAL_ERR
	// e.g. missing category
	BAD_INPUT BestMealErr = 1
	// nothing available within your budget
	TOO_POOR
)
*/

type BestMealError struct {
	error string
}

////////////////////////////////////////////////////////////////////
// define the input

type MenuItem struct {
	name         string
	cost         int
	satisfaction int
	category     string
}

type Menu struct {
	foods  []MenuItem
	budget int
}

func findBestMeal(inputFile string) ([]BestMeal, error) {

	menu, err := loadMenuAndBudget(inputFile)

	if err != nil {
		return nil, err
	}

	meals, err := findMostSatisfyingMeal(menu.foods, menu.budget)

	if err != nil {
		return nil, err
	}

	mostSatisfyingMeal := meals[0]

	bestMeal := BestMeal{
		[4]string{
			menu.foods[mostSatisfyingMeal.appetizer].name,
			menu.foods[mostSatisfyingMeal.drink].name,
			menu.foods[mostSatisfyingMeal.mainCourse].name,
			menu.foods[mostSatisfyingMeal.dessert].name,
		},
		mostSatisfyingMeal.totalCost, mostSatisfyingMeal.totalSatisfaction}

	return []BestMeal{bestMeal}, nil

}

func loadMenuAndBudget(inputFile string) (Menu, error) {

	menu := Menu{make([]MenuItem, 0), 0}

	menu.foods = append(menu.foods, MenuItem{"Fried Calamari", 6, 5, "Appetizer"})
	menu.foods = append(menu.foods, MenuItem{"Bruschetta", 4, 4, "Appetizer"})

	menu.foods = append(menu.foods, MenuItem{"Soda", 1, 1, "Drink"})
	menu.foods = append(menu.foods, MenuItem{"Beer", 3, 2, "Drink"})

	menu.foods = append(menu.foods, MenuItem{"Lasagna", 8, 7, "Main Course"})
	menu.foods = append(menu.foods, MenuItem{"Burger", 6, 5, "Main Course"})

	menu.foods = append(menu.foods, MenuItem{"Cheesecake", 4, 4, "Dessert"})
	menu.foods = append(menu.foods, MenuItem{"Ice Cream", 2, 2, "Dessert"})

	menu.budget = 25

	return menu, nil
}

type Meal struct {
	appetizer         int
	drink             int
	mainCourse        int
	dessert           int
	totalCost         int
	totalSatisfaction int
}

func findMostSatisfyingMeal(foods []MenuItem, budget int) ([]Meal, error) {

	if /* foods == nil || */ len(foods) == 0 {
		return nil, errors.New("Nothing in the menu??")
	}

	// "The number of items per category on the menu can range from 0 to 50." -- main README from Alek
	// with a max of 50 items in each of 4 categories, the max possible meals would be 50^4 or 6.25M

	// assemble the list of items by category
	apps := make([]int, 0)
	drinks := make([]int, 0)
	mains := make([]int, 0)
	desserts := make([]int, 0)

	for i, food := range foods {

		// fmt.Println(food)
		category := strings.ToLower(strings.ReplaceAll(food.category, " ", ""))
		switch category {
		case APPETIZER:
			{
				apps = append(apps, i)
			}
		case DRINK:
			{
				drinks = append(drinks, i)
			}
		case MAIN_COURSE:
			{
				mains = append(mains, i)
			}
		case DESSERT:
			{
				desserts = append(desserts, i)
			}
		default:
			// TODO Unknown category
			return nil, errors.New("Unknown category: " + food.category)
		}
	}

	// verify we have something in each category
	if len(apps) == 0 {
		return nil, errors.New("No apps in menu")
	}

	if len(drinks) == 0 {
		return nil, errors.New("No drinks in menu")
	}

	if len(mains) == 0 {
		return nil, errors.New("No mains in menu")
	}

	if len(desserts) == 0 {
		return nil, errors.New("No desserts in menu")
	}

	// now that we have validated input broken down by category,
	// assemble the list of meals **within budget**

	foundMostSatisfyingMeal := false
	maxSatisfaction := 0
	lowestCost := math.MaxInt
	var mostSatisfyingMeal Meal

	for _, app := range apps {
		for _, drink := range drinks {
			for _, main := range mains {
				for _, dessert := range desserts {

					// fmt.Printf("Meal: %d %d %d %d\n", app, drink, main, dessert)
					totalCost := 0 +
						foods[app].cost +
						foods[drink].cost +
						foods[main].cost +
						foods[dessert].cost

					if totalCost <= budget {

						// this meal is within budget
						// compute our total satisfaction for this meal

						totalSatisfaction := 0 +
							foods[app].satisfaction +
							foods[drink].satisfaction +
							foods[main].satisfaction +
							foods[dessert].satisfaction

						// highest satisfaction and lowest cost so far?
						if totalSatisfaction > maxSatisfaction ||
							// it seems like lower cost also figures into satisfaction
							(totalSatisfaction >= maxSatisfaction && totalCost < lowestCost) {
							// save off the meal
							mostSatisfyingMeal = Meal{app, drink, main, dessert, totalCost, totalSatisfaction}
							maxSatisfaction = totalSatisfaction
							lowestCost = totalCost
							foundMostSatisfyingMeal = true
						}
					}
				}
			}
		}
	}

	if foundMostSatisfyingMeal {
		return []Meal{mostSatisfyingMeal}, nil
	} else {
		return nil, errors.New("No meals for your budget")
	}
}

const ( // MenuItemCategory
	APPETIZER   string = "appetizer"
	DRINK       string = "drink"
	MAIN_COURSE string = "maincourse"
	DESSERT     string = "dessert"
)

const ( // MenuItemCategoryIndex
	APP_INDEX = iota
	DRINK_INDEX
	MAIN_INDEX
	DESSERT_INDEX
)
