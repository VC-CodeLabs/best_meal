
public class Food {
    String name = new String();
    int cost;
    int satisfaction;
    String category;
    
    public void setName(String inputName){
        name = inputName;        
    }
    
    public String getName(){
        return name;
    }
    
    public void setCost(int inputCost){
        cost = inputCost;        
    }
    
    public int getCost(){
        return cost;
    }
    
    public void setSatisfaction(int inputSatisfaction){
        satisfaction = inputSatisfaction;        
    }
    
    public int getSatisfaction(){
        return satisfaction;
    }
    
    public void setCategory(String inputCategory){
        category = inputCategory;        
    }
    
    public String getCategory(){
        return category;
    }

}
