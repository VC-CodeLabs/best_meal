open System.IO
open System.Text.Json



type Food = {
    name : string
    cost : int
    satisfaction : int 
    category : string
}
type  MenuOrder = {
    foods: Food list
    budget: int

}

let readMenuFile (filePath: string) : MenuOrder = 
    try 
        let menuFile = File.ReadAllText(filePath)
        let jsonOptions = JsonSerializerOptions()
        let menuOrder = JsonSerializer.Deserialize<MenuOrder>(menuFile,jsonOptions)
        match menuOrder with
        | data -> data
    with 
    | ex-> failwithf "can't read JSON file : %s " ex.Message

let filterMenuItems (menuOrder: MenuOrder) = 
    let filterOnCategory = 
        menuOrder.foods
        |> List.groupBy(fun food -> food.category)
        |> List.map (fun (category, foods) ->
            category, List.sortByDescending (fun food -> food.satisfaction) foods)

    let rec buildMeal currentFood remainCategories remainBudget =
        match remainCategories with
        | [] -> Some currentFood 
        | (category, foods)::rest ->
            // Try each food item from the current category
            let possMeals =
                foods
                |> List.filter (fun food -> food.cost <= remainBudget) 
                |> List.map (fun food ->
                    let newFood = food.name :: currentFood
                    let newBudget = remainBudget - food.cost
                    buildMeal newFood rest newBudget)

            //find best meal in budg
            match possMeals |> List.filter Option.isSome with
            | [] -> None 
            | meals -> meals |> List.maxBy Option.get 

    let emptyList = []
    let categories = filterOnCategory
    let budget = menuOrder.budget

    // Construct the  meal tuple
    match buildMeal emptyList categories budget with
    | Some Meal ->
        let totalCost = List.sumBy (fun foodName ->
            menuOrder.foods |> List.find (fun food -> food.name = foodName) |> fun food -> food.cost) Meal
        let totalSatisfaction = List.sumBy (fun foodName ->
            menuOrder.foods |> List.find (fun food -> food.name = foodName) |> fun food -> food.satisfaction) Meal
        (Meal, totalCost, totalSatisfaction)
    | None -> failwithf "Not enough money found to buy a whole meal. Broke nerd"

        
let filePath = "C:/Users/VividCloud/Projects/johntruemancodelabs/menu.json"
let menu = readMenuFile filePath

let (Meal, totalCost, totalSatisfaction) = filterMenuItems menu
printfn "{"
printfn "  \"meal\": %A," Meal
printfn "  \"cost\": %d," totalCost
printfn "  \"satisfaction\": %d" totalSatisfaction
printfn "}"
    