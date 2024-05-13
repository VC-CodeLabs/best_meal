// the primary algorithm is in the findMostSatisfayingMeal function-
// see the godoc on and within that function for details
package main

//

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

// the name of the input menu json file;
//
//	define a variable here at the top for easy testing support-
//	NOTE this can also be overridden with -f={menuJsonFileSpec}
var inputFile string = "menu.json"

func main() {

	//
	// process command-line, if any;
	// if no parameters are used,
	// the default "menu.json" is expected to be in the current directory,
	// the only output will an error or best meal, both in json format.
	//

	// allow the menu.json input filespec to be customized without changing code
	filePtr := flag.String("f", inputFile, "the input menu json filespec")
	// verbose mode supports troubleshooting
	verbosePtr := flag.Bool("v", VERBOSE, "enable verbose output for troubleshooting")
	// the above *declares our support for command line params,
	// we must actually [flag.Parse]() to process all parameters
	flag.Parse()
	if verbosePtr != nil {
		// we had the -v parameter
		VERBOSE = *verbosePtr
		if VERBOSE {
			// -v or -v=true
			log.Printf("VERBOSE mode enabled.\n")
		}
	}
	if filePtr != nil {
		// we had the -f={menuJsonFileSpec} parameter
		inputFile = *filePtr
		if VERBOSE {
			log.Printf("Reading menu json from `%s`\n", inputFile)
		}
	}

	// find the best meal and report on it, else report an error
	findAndEmitBestMeal(inputFile)
}

// the global-ish flag for verbose mode-
// if enabled, logs inner workings for algo
var VERBOSE = false

// find the best meal and report on it, else report an error
func findAndEmitBestMeal(inputFile string) {
	bestMeals, err := findBestMeal(inputFile)

	if err != nil {
		emitBestMealError(err)
	} else {

		bestMeal := bestMeals[0]

		emitBestMeal(bestMeal)
	}
}

// write the error that occurred while trying to find best meal to the console
func emitBestMealError(err error) {

	// construct the best meal error object
	bestMealError := BestMealError{err.Error()}

	//// original code to "manually" output in json:
	//// fmt.Printf("{\n\t\"error\": \"%s\"\n}\n", bestMealError.Error)

	// pretty print the output in json format, indenting with tabs
	json, _ := json.MarshalIndent(&bestMealError, "", "\t")
	fmt.Println(string(json))
}

// write the best meal we found to the console
func emitBestMeal(bestMeal BestMeal) {

	//// original code "manually" output best meal in json format
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

	// pretty print the output in json format, indenting with tabs
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

// a MenuItem represents an entry in the menu.foods array
type MenuItem struct {
	// the name of this menu item e.g. "Steak"
	Name string
	// the cost of this menu item in dollars
	Cost int
	// the satisfaction score of this menu item
	Satisfaction int
	// the category of this menu item e.g. "Main Course"
	Category string
}

// define the structure of the menu input:
//
//	an array of foods + our budget
type Menu struct {
	// the set of foods, including their cost, satisfaction and category
	Foods []MenuItem
	// the budget we have to spend on a four-part meal from this menu in dollars
	Budget int
}

func findBestMeal(inputFile string) ([]BestMeal, error) {

	menu, err := loadMenuAndBudget(inputFile)

	if err != nil {
		return nil, err
	}

	// do some basic error checking on the input
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

	// (eventually) our menu: foods + budget
	menu := Menu{make([]MenuItem, 0), 0}

	//// originally had the menu hard-wired into the code for a q&d test
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

	// given the limit on input (200 items),
	// should be no issue reading the entire text file into memory
	menuBytes, err := os.ReadFile(inputFile)

	if err != nil {
		// bad file or access issue
		return menu, err
	}

	// menuString := string(menuBytes)

	// convert the json to in-memory object representation
	err = json.Unmarshal(menuBytes, &menu)

	if err != nil {
		// some issue with the json
		err = errors.New("Bad menu json: " + err.Error())
	}

	// menu loaded from json, return it w/ no error
	return menu, err
}

// the Meal struct is used to track a particular instance of a meal-
// the combination of a specific appetizer, drink, main course and dessert;
// for the food items, we're actually tracking the index in the original menu item
// rather than a separate copy.
//
// we also track the totalCost and totalSatisfaction,
// the sum of the corresponding fields from the four food items in this meal
type Meal struct {
	// menu.foods index for appetizer
	appetizer int
	// menu.foods index for drink
	drink int
	// menu.foods index for main course
	mainCourse int
	// menu.foods index for dessert
	dessert int
	// sum of menu.foods[meal.appetizer|drink|mainCourse|dessert].cost
	totalCost int
	// sum of menu.foods[meal.appetizer|drink|mainCourse|dessert].satisfaction
	totalSatisfaction int
}

// *the* solution algorithm- given a list of food items and budget in dollars,
// find the most satisfying aka best meal
//
//   - break down the foods by category
//   - validate each category has at least one food
//   - for each combination of foods, one from each category, aka a meal:
//   - compute the total cost
//   - if total cost fits our budget:
//   - compute the total satisfaction
//   - if total satisfaction is greater than prior most-satisfying max,
//   - *OR* if total satisfaction equals and total cost is lower than prior most-satisfying meal,
//   - save the meal as most-satisfying so far
//
// as noted in the solution-specific readme,
// I factored cost into the equation-
// given two meals with the same total satisfaction,
// the meal with the lower total cost will be favored.
func findMostSatisfyingMeal(foods []MenuItem, budget int) ([]Meal, error) {

	if /* foods == nil || */ len(foods) == 0 {
		return nil, errors.New("Nothing in the menu??")
	}

	// "The number of items per category on the menu can range from 0 to 50." -- main README from Alek
	// with a max of 50 items in each of 4 categories, the max possible meals would be 50^4 or 6.25M

	// assemble the list of food items by category;
	// the values we're tracking here are indexes in the set of foods
	apps := make([]int, 0)
	drinks := make([]int, 0)
	mains := make([]int, 0)
	desserts := make([]int, 0)

	for i, food := range foods {

		if VERBOSE {
			log.Printf("Input foods[%d]: %+v\n", i, food)
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
			return nil, fmt.Errorf("Unknown foods[%d] category: %+v", i, food)
		}
	}

	// verify we have something in each category- note only the first gap found is reported
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

	if VERBOSE {
		log.Printf("Checking %d apps x %d drinks x %d mains x %d desserts = %d meals within $%d budget\n",
			len(apps), len(drinks), len(mains), len(desserts),
			len(apps)*len(drinks)*len(mains)*len(desserts),
			budget)
	}

	// now that we have validated input broken down by category,
	// find the most satisfying meal within our budget

	// have we found at least one candidate
	foundMostSatisfyingMeal := false
	// track the maximum total satisfaction we've found so far
	maxSatisfaction := 0
	// track the
	lowestCost := math.MaxInt
	// track the cheapest four-part meal regardless of budget;
	// if no meal is available within budget,
	// we'll let them know how far short they are
	minimumCost := math.MaxInt

	// track the actual most satisfying meal within budget found so far
	// as we work our way thru the menu food items
	var mostSatisfyingMeal Meal

	//
	mealCounter := 0

	// for each combination of foods, one from each category

	for _, app := range apps {
		for _, drink := range drinks {
			for _, main := range mains {
				for _, dessert := range desserts {

					if VERBOSE {
						log.Printf("Meal #%d foods[] indexes: app=[%d] drink=[%d] main=[%d] dessert=[%d]\n",
							mealCounter, app, drink, main, dessert)
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

						// highest satisfaction
						if totalSatisfaction > maxSatisfaction ||
							// equal satisfaction but lower cost (which is also satisfying)
							(totalSatisfaction >= maxSatisfaction && totalCost < lowestCost) {

							// save off the candidate meal and stats
							mostSatisfyingMeal = Meal{app, drink, main, dessert, totalCost, totalSatisfaction}
							maxSatisfaction = totalSatisfaction
							lowestCost = totalCost
							foundMostSatisfyingMeal = true

							// since we found at least one candidate, we'll no longer be tracking minimum cost
							minimumCost = -1

							if VERBOSE {
								log.Printf("** Most Satisfying + Lowest Cost (so far): %s totalCost=%d totalSatisfaction=%d\n", foodNames(foods, app, drink, main, dessert), totalCost, totalSatisfaction)
							}
						} else {
							if VERBOSE {
								log.Printf("Less Satisfying: %s totalCost=%d totalSatisfaction=%d\n", foodNames(foods, app, drink, main, dessert), totalCost, totalSatisfaction)
							}
						}

					} else {
						if VERBOSE {
							log.Printf("Over budget: %s totalCost=%d\n", foodNames(foods, app, drink, main, dessert), totalCost)
						}
					}

					if !foundMostSatisfyingMeal {
						if totalCost < minimumCost {
							minimumCost = totalCost
							if VERBOSE {
								log.Printf("Cheapest meal cost: %d\n", minimumCost)
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

// convert the indexes for foods in a four-part meal to the food names
func foodNames(foods []MenuItem, app int, drink int, main int, dessert int) [4]string {
	return [4]string{
		foods[app].Name,
		foods[drink].Name,
		foods[main].Name,
		foods[dessert].Name,
	}
}

// the food type names used to categorize the foods in the menu
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
