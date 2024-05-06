package main

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"log"
	"math"
	"os"
	"strings"
)

// note the go.mod is used to support having the solution in a folder

// variable for testing support
var inputFile string = "menu.json"

func main() {

	filePtr := flag.String("f", inputFile, "the input menu json filespec")
	verbosePtr := flag.Bool("v", VERBOSE, "enable verbose output for troubleshooting")
	flag.Parse()
	if verbosePtr != nil {
		VERBOSE = *verbosePtr
		if VERBOSE {
			log.Printf("VERBOSE mode enabled.\n")
		}
	}
	if filePtr != nil {
		inputFile = *filePtr
		if VERBOSE {
			log.Printf("Reading menu json from `%s`\n", inputFile)
		}
	}

	findAndEmitBestMeal(inputFile)
}

var VERBOSE = false

func findAndEmitBestMeal(inputFile string) {
	bestMeals, err := findBestMeal(inputFile)

	if err != nil {
		emitBestMealError(err)
	} else {

		bestMeal := bestMeals[0]

		emitBestMeal(bestMeal)
	}
}

func emitBestMealError(err error) {

	// write the error in json format
	bestMealError := BestMealError{err.Error()}
	// fmt.Printf("{\n\t\"error\": \"%s\"\n}\n", bestMealError.Error)
	json, _ := json.MarshalIndent(&bestMealError, "", "\t")
	fmt.Println(string(json))
}

func emitBestMeal(bestMeal BestMeal) {

	// write the best meal in json format
	/*
		fmt.Printf("{\n"+
			"\t\"selectedFoods\": %s,\n"+
			"\t\"totalCost\": %d,\n"+
			"\t\"totalSatisfaction\": %d\n"+
			"}\n",
			fmt.Sprintf("[ \"%s\", \"%s\", \"%s\", \"%s\" ]",
				bestMeal.SelectedFoods[APP_INDEX],
				bestMeal.SelectedFoods[DRINK_INDEX],
				bestMeal.SelectedFoods[MAIN_INDEX],
				bestMeal.SelectedFoods[DESSERT_INDEX],
			),
			bestMeal.TotalCost,
			bestMeal.TotalSatisfaction,
		)
	*/
	json, _ := json.MarshalIndent(&bestMeal, "", "\t")
	fmt.Println(string(json))

}

////////////////////////////////////////////////////////////////////
// define the output

type BestMeal struct {
	SelectedFoods     [4]string `json:"selectedFoods"`
	TotalCost         int       `json:"totalCost"`
	TotalSatisfaction int       `json:"totalSatisfaction"`
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
	Error string `json:"error"`
}

////////////////////////////////////////////////////////////////////
// define the input

type MenuItem struct {
	Name         string
	Cost         int
	Satisfaction int
	Category     string
}

type Menu struct {
	Foods  []MenuItem
	Budget int
}

func findBestMeal(inputFile string) ([]BestMeal, error) {

	menu, err := loadMenuAndBudget(inputFile)

	if err != nil {
		return nil, err
	}

	if len(menu.Foods) == 0 {
		return nil, errors.New("No food in menu??")
	}

	if menu.Budget <= 0 {
		return nil, errors.New("You need a budget")
	}

	meals, err := findMostSatisfyingMeal(menu.Foods, menu.Budget)

	if err != nil {
		return nil, err
	}

	mostSatisfyingMeal := meals[0]

	bestMeal := BestMeal{
		[4]string{
			menu.Foods[mostSatisfyingMeal.appetizer].Name,
			menu.Foods[mostSatisfyingMeal.drink].Name,
			menu.Foods[mostSatisfyingMeal.mainCourse].Name,
			menu.Foods[mostSatisfyingMeal.dessert].Name,
		},
		mostSatisfyingMeal.totalCost, mostSatisfyingMeal.totalSatisfaction}

	return []BestMeal{bestMeal}, nil

}

func loadMenuAndBudget(inputFile string) (Menu, error) {

	menu := Menu{make([]MenuItem, 0), 0}

	/*
		menu.foods = append(menu.foods, MenuItem{"Fried Calamari", 6, 5, "Appetizer"})
		menu.foods = append(menu.foods, MenuItem{"Bruschetta", 4, 4, "Appetizer"})

		menu.foods = append(menu.foods, MenuItem{"Soda", 1, 1, "Drink"})
		menu.foods = append(menu.foods, MenuItem{"Beer", 3, 2, "Drink"})

		menu.foods = append(menu.foods, MenuItem{"Lasagna", 8, 7, "Main Course"})
		menu.foods = append(menu.foods, MenuItem{"Burger", 6, 5, "Main Course"})

		menu.foods = append(menu.foods, MenuItem{"Cheesecake", 4, 4, "Dessert"})
		menu.foods = append(menu.foods, MenuItem{"Ice Cream", 2, 2, "Dessert"})

		menu.budget = 25
	*/

	menuBytes, err := os.ReadFile(inputFile)

	if err != nil {
		return menu, err
	}

	// menuString := string(menuBytes)

	err = json.Unmarshal(menuBytes, &menu)

	if err != nil {
		err = errors.New("Bad menu json: " + err.Error())
	}

	return menu, err
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

		if VERBOSE {
			log.Printf("Input Food %+v\n", food)
		}

		category := strings.ToLower(strings.ReplaceAll(food.Category, " ", ""))
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
			return nil, errors.New("Unknown food category: " + food.Category)
		}
	}

	// verify we have something in each category- note only the first found is reported
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
	minimumCost := math.MaxInt
	var mostSatisfyingMeal Meal

	mealCounter := 0

	for _, app := range apps {
		for _, drink := range drinks {
			for _, main := range mains {
				for _, dessert := range desserts {

					if VERBOSE {
						log.Printf("Meal #%d Food indexes: %d %d %d %d\n", mealCounter, app, drink, main, dessert)
					}

					totalCost := 0 +
						foods[app].Cost +
						foods[drink].Cost +
						foods[main].Cost +
						foods[dessert].Cost

					if totalCost <= budget {

						// this meal is within budget
						// compute our total satisfaction for this meal

						totalSatisfaction := 0 +
							foods[app].Satisfaction +
							foods[drink].Satisfaction +
							foods[main].Satisfaction +
							foods[dessert].Satisfaction

						// highest satisfaction and lowest cost so far?
						if totalSatisfaction > maxSatisfaction ||
							// it seems like lower cost also figures into satisfaction
							(totalSatisfaction >= maxSatisfaction && totalCost < lowestCost) {
							// save off the meal
							mostSatisfyingMeal = Meal{app, drink, main, dessert, totalCost, totalSatisfaction}
							maxSatisfaction = totalSatisfaction
							lowestCost = totalCost
							foundMostSatisfyingMeal = true
							minimumCost = -1

							if VERBOSE {
								log.Printf("** Most Satisfying + Lowest Cost (so far): %s cost=%d satisfaction=%d\n", foodNames(foods, app, drink, main, dessert), totalCost, totalSatisfaction)
							}
						} else {
							if VERBOSE {
								log.Printf("Less Satisfying: %s cost=%d satisfaction=%d\n", foodNames(foods, app, drink, main, dessert), totalCost, totalSatisfaction)
							}
						}

					} else {
						if VERBOSE {
							log.Printf("Over budget: %s cost=%d\n", foodNames(foods, app, drink, main, dessert), totalCost)
						}
					}

					if !foundMostSatisfyingMeal {
						if totalCost < minimumCost {
							minimumCost = totalCost
							if VERBOSE {
								log.Printf("Minimal cost: %d\n", minimumCost)
							}
						}
					}

					mealCounter++
				}
			}
		}
	}

	if foundMostSatisfyingMeal {
		if VERBOSE {
			log.Printf("**** Most Satisfying + Lowest Cost Meal: %s totalCost=%d totalSatisfaction=%d\n",
				foodNames(foods,
					mostSatisfyingMeal.appetizer,
					mostSatisfyingMeal.drink,
					mostSatisfyingMeal.mainCourse,
					mostSatisfyingMeal.dessert),
				mostSatisfyingMeal.totalCost,
				mostSatisfyingMeal.totalSatisfaction)
		}
		return []Meal{mostSatisfyingMeal}, nil
	} else {
		if VERBOSE {
			log.Printf("*** Budget=%d vs Cheapest Meal=%d\n", budget, minimumCost)
		}
		return nil,
			fmt.Errorf(""+
				"Checked %d meal(s), none fit your budget- "+
				"you need another %d buck(s) to dine here :/",
				mealCounter, minimumCost-budget)
	}
}

func foodNames(foods []MenuItem, app int, drink int, main int, dessert int) [4]string {
	return [4]string{
		foods[app].Name,
		foods[drink].Name,
		foods[main].Name,
		foods[dessert].Name,
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
