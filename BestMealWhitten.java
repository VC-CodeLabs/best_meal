import java.io.File;
import java.io.FileNotFoundException;
import java.util.Map;
import java.util.Scanner;
import java.util.ArrayList;
import java.util.List;
import java.util.EnumMap;

/*  For my approach, I calculate all possible valid meals (ie, one item from each 
    category), discarding all the meals where the total budget exceeds the allowed 
    budget. I keep track of which meal has the highest satisfaction while iterating,
    and display its information at the end.

    Though this is a "brute force" approach, given the constraints of the task (max
    50 food items) this will never be worse than 292,032 iterations (assuming even 
    distribution amongst categories, 13x13x12x12).

    I did my best to reduce code smells, including cognitive complexity, which is why
    there are so many helper functions.
*/

enum Category {
    APPETIZER,
    DRINK,
    MAINCOURSE,
    DESSERT
}

class Food {
    String name;
    int cost;
    int satisfaction;
    Category category;

    public Food(String name, int cost, int satisfaction, Category category) {
        this.name = name;
        this.cost = cost;
        this.satisfaction = satisfaction;
        this.category = category;
    }
}

public class BestMealWhitten {
    protected static List<Food> menu = new ArrayList<>();
    protected static int allowedBudget;
    protected static List<EnumMap<Category, Food>> allValidPossibleMeals = new ArrayList<>();

    protected static EnumMap<Category, Food> bestMeal = new EnumMap<>(Category.class);
    protected static int bestMealSatisfaction = 0;

    public static void main(String[] args) {
        // Read in the json file, create all the menu items, and set the budget
        populateMenu("src/menu.json");

        // Brute force all possible meals, combining one food from each category
        // Keep track of the most satisfying meal that's within budget
        findAllPossibleValidMeals();

        if (bestMeal == null) {
            System.out.println("Unable to find satisfactory meal :(");
        } else {
            System.out.println("Your meal is served, sir/madam. Bon appétit!");
            for (Category category : Category.values()) {
                Food food = bestMeal.get(category);
                System.out.printf("     %-12s: %s ($%d / %s)%n", category.name(), food.name, food.cost, food.satisfaction);
            }
            System.out.printf("Total budget used: $%d (allowed budget was $%d)%n", calculateTotalCost(bestMeal), allowedBudget);
            System.out.printf("Total satisfaction from meal: %d%n", bestMealSatisfaction);
        }
    }

    private static void populateMenu(String filename) {
        File myObj = new File(filename);
        try {
            Scanner myReader = new Scanner(myObj);
            while (myReader.hasNextLine()) {
                String data = myReader.nextLine();
                if (data.contains("name")) {
                    menu.add(createNewFood(data));
                } else if (data.contains("budget")) {
                    allowedBudget = Integer.parseInt(data.split(": ")[1]);
                }
            }
            myReader.close();
        } catch (FileNotFoundException e) {
            System.out.println("File Not Found. Make sure the file is named \"menu.json\" and is in the src directory.");
            e.printStackTrace();
        }
    }

    private static void findAllPossibleValidMeals() {
        // Split the menu into its component categories
        List<Food> appetizers = filterByCategory(Category.APPETIZER);
        List<Food> mainCourses = filterByCategory(Category.MAINCOURSE);
        List<Food> drinks = filterByCategory(Category.DRINK);
        List<Food> desserts = filterByCategory(Category.DESSERT);
        
        for (Food app : appetizers) {
            for (Food main : mainCourses) {
                for (Food drink : drinks) {
                    for (Food dessert : desserts) {
                        EnumMap<Category, Food> meal = createMeal(app, main, drink, dessert);
                        
                        // Check to make sure the meal is valid (ie under budget and containing
                        // one food of each category), and if it's better than our current 
                        // "bestMeal", swap it.
                        evaluateMeal(meal);
                    }
                }
            }
        }
    }
    
    private static void evaluateMeal(EnumMap<Category, Food> meal) {
        if (isMealValid(meal)) {
            int mealSatisfaction = calculateTotalSatisfaction(meal);
            // if the new meal is better, it becomes King of the Hill (or plate)
            if (mealSatisfaction > bestMealSatisfaction) {
                bestMealSatisfaction = mealSatisfaction;
                bestMeal = meal;
            }
        }
    }

    private static EnumMap<Category, Food> createMeal(Food app, Food main, Food drink, Food dessert) {
        EnumMap<Category, Food> meal = new EnumMap<>(Category.class);
        meal.put(Category.APPETIZER, app);
        meal.put(Category.MAINCOURSE, main);
        meal.put(Category.DRINK, drink);
        meal.put(Category.DESSERT, dessert);
        return meal;
    }

    private static List<Food> filterByCategory(Category category) {
        return menu.stream()
            .filter(food -> food.category == category)
            .toList();
    }

    private static boolean isMealValid(EnumMap<Category, Food> meal) {
        return allCategoriesPresent(meal) && calculateTotalCost(meal) <= allowedBudget;
    }

    private static Food createNewFood(String data) {
        String name = null;
        int cost = 0;
        int satisfaction = 0;
        Category category = null;

        // Parse the string that holds the Food item
        String[] parts = data.substring(1, data.length() - 1).split(", ");

        for (String part : parts) {
            String[] keyValue = part.split(": ");
            String key = keyValue[0].replaceAll("[\"{} ]", "");
            String value = keyValue[1].replaceAll("[\"{}]", "");
            
            switch (key) {
                case "name":
                    name = value;
                    break;
                case "cost":
                    cost = Integer.parseInt(value);
                    break;
                case "satisfaction":
                    satisfaction = Integer.parseInt(value);
                    break;
                case "category":
                    category = Category.valueOf(value.replace(" ", "").toUpperCase());
                    break;
                default:
                    break;
            }
        }
        
        return new Food(name, cost, satisfaction, category);
    }

    private static boolean allCategoriesPresent(EnumMap<Category, Food> meal) {
        for (Category category : Category.values()) {
            if (!meal.containsKey(category)) {
                return false;
            }
        }
        return true;
    }

    private static int calculateTotalCost(EnumMap<Category, Food> meal) {
        int totalCost = 0;
        for (Map.Entry<Category, Food> entry : meal.entrySet()) {
            Food food = entry.getValue();
            totalCost += food.cost;
        }
        return totalCost;
    }

    private static int calculateTotalSatisfaction(EnumMap<Category, Food> meal) {
        int totalSatisfaction = 0;
        for (Map.Entry<Category, Food> entry : meal.entrySet()) {
            Food food = entry.getValue();
            totalSatisfaction += food.satisfaction;
        }
        return totalSatisfaction;
    }
}
