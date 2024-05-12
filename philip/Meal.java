
public class Meal {
    Food appetizer = new Food();
    Food drink = new Food();
    Food mainCourse = new Food();
    Food dessert = new Food();
    int totalSatisfaction = -1;
    int totalCost = 0;
    
    public Meal() {
        
    }
    
    public Meal (Food app, Food dr, Food mc, Food des) {
        appetizer = app;
        drink = dr;
        mainCourse = mc;
        dessert = des;
        totalSatisfaction = appetizer.getSatisfaction() + drink.getSatisfaction() + mainCourse.getSatisfaction() + dessert.getSatisfaction();
        totalCost = appetizer.getCost() + drink.getCost() + mainCourse.getCost() + dessert.getCost();
    }
    
    public Meal (Food app, Food dr, Food mc, Food des, int sat) {
        appetizer = app;
        drink = dr;
        mainCourse = mc;
        dessert = des;
        totalSatisfaction = sat; 
        totalCost = appetizer.getCost() + drink.getCost() + mainCourse.getCost() + dessert.getCost();
    }
    
    public Meal (Food app, Food dr, Food mc, Food des, int sat, int cost) {
        appetizer = app;
        drink = dr;
        mainCourse = mc;
        dessert = des;
        totalSatisfaction = sat; 
        totalCost = cost;
    }
    
    public Food getAppetizer() {
        return appetizer;
    }
    
    public Food getDrink() {
        return drink;
    }
    
    public Food getMainCourse() {
        return mainCourse;
    }
    
    public Food getDessert() {
        return dessert;
    }
    
    public int getSatisfaction() {
        return totalSatisfaction;
    }
    
    public int getCost() {
        return totalCost;
    }
    
    
    public void setAppetizer(Food app) {
        appetizer = app;
    }
    
    public void setDrink(Food dr) {
        drink = dr;
    }
    
    public void setMainCourse(Food mc) {
        mainCourse = mc;
    }
    
    public void setdessert(Food des) {
        dessert = des;
    }
    
    public void setSatisfaction(int sat) {
        totalSatisfaction = sat;
    }
    
    public void setCost(int cost) {
        totalCost = cost;
    }
            

}
