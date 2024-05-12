// TomsBestMeal.cpp
//

/*
Simple problem statement:

You're looking to help plan a perfect dining experience within a given budget.
You will be provided with a menu and a budget. The menu will include various dishes across four categories:
Appetizers, Drinks, Main Courses, and Desserts. Each dish has an associated cost and satisfaction level.
You must select at least one item from each category to maximize total satisfaction without exceeding the budget.

Methodology:
Use a 2D search to find the highest satisfaction for a meal without going over the budget.  Two solutions:
    1. Solution 1:  Do an exhaustive search
    2. Solution 2:  Find the optimal solution by sorting the input in a way that the highest satisfaction
       items are evaluated first.  Then setting a search depth parameter to stop the search before
       exhausting it.  Note:  it turns out that the max case is still very fast and no need for Solution 2.

- Read the menu.json file into an array of objects
- Organize the array:
    - Put into columns by category
    - If there is not at least one item for each category then end with error.  Output result in jason format.
    - Sort the items in each category by highest to lowest satisfaction (Solution 2)
- Evaluate the array:
    - Go through every combination of array, from top to bottom and left to right
    - Track highest satisfaction within budget
    - Stop when max depth is exceeded (Solution 2) (not implemented)
- Output result in jason format.  Could be a null result.

*/

#include <cstdio>
#include <cstdlib>
#include <time.h>
#include <atlstr.h>
#include <stdio.h>
#include <stdlib.h>

enum Category { APP = 0, DRINK, MAIN, DESSERT };

struct Item
{
    bool validItem;
    char name[100];
    Category cat;
    int sat;
    int cost;
};

int GetMenu(Item menu[4][50], bool sorted, int *p_a, int *p_dr, int *p_m, int *p_de);
void OutputResult(bool success, Item selected[], int satisfaction);

int main()
{
    Item menu[4][50];  // set to max size since it is not too big
    int budget, satisfaction = 0, temp;
    int a = 0, dr = 0, m = 0, de = 0;  // counters for the 4 categories
    Item selected[4];
    int i, j, k, l;

    // Read the menu.json file into an array of objects and get item counts
    // FALSE = no sorting of items
    budget = GetMenu(menu, FALSE, &a, &dr, &m, &de);

    // Look for incomplete menu
    if (a == 0 || dr == 0 || m == 0 || de == 0)
    {
        // Output error in json format
        OutputResult(FALSE, NULL, satisfaction);
        return(0);
    }

    // Find the maximum satisfaction within the budget
    // Go through every combination of array, from top to bottom and left to right
    // Track highest satisfaction within budget
    for (i = 0; i < a; i++)
    {
        for (j = 0; j < dr; j++)
        {
            for (k = 0; k < m; k++)
            {
                for (l = 0; l < de; l++)
                {
                    // Check to see if this menu has greater satisfaction than the saves one and within budget
                    if (((temp = menu[APP][i].sat + menu[DRINK][j].sat + menu[MAIN][k].sat + menu[DESSERT][l].sat) > satisfaction) &&
                        ((menu[APP][i].cost + menu[DRINK][j].cost + menu[MAIN][k].cost + menu[DESSERT][l].cost) <= budget))
                    {
                        // Found a menu that is better than the best so far, save it
                        memcpy(&selected[APP], &menu[APP][i], sizeof(Item));
                        memcpy(&selected[DRINK], &menu[DRINK][j], sizeof(Item));
                        memcpy(&selected[MAIN], &menu[MAIN][k], sizeof(Item));
                        memcpy(&selected[DESSERT], &menu[DESSERT][l], sizeof(Item));
                        satisfaction = temp;
                    }
                }

            }

        }

    }
    OutputResult(TRUE, selected, satisfaction);
    return(0);
}

#define MAX_TOKEN 6

// Simple json reader with no error checking (sorting on or off)
// Requires strict json formatting
// Assume line breaks only after full menu items
// Returns budget
int GetMenu(Item menu[4][50], bool sorted, int* p_a, int* p_dr, int* p_m, int* p_de)
{
    FILE* pFile;
    Item temp;
    char tokens[MAX_TOKEN][30] = {"name", "category", "cost", "satisfaction", "budget", "foods"};
    int token, i, j;
    char* pStart, *pEnd;
    char line[1000];
    char tempstr[30];
    bool endLine, endItem;
    int budget = 0;

    memset(&temp, 0, sizeof(Item));
    fopen_s(&pFile, "menu.json", "r");
    endItem = FALSE;
    while (fgets(line, 1000, pFile) != NULL)
    {
        //  Process the line of text
        endLine = FALSE;
        pStart = &line[0];
        while (!endLine)
        {
            // Look for the end of an item
            if (*pStart == '}' && temp.validItem)
            {
                endItem = TRUE;
                break;
            }

            // Make sure the line has content
            if (strstr(pStart, "\"") != NULL)
            {
                // Find beginning and end of token
                pStart = strstr(pStart, "\"") + 1;
                pEnd = strstr(pStart, "\"");
                pEnd[0] = '\0';

                // Get token
                for (i = 0; i < MAX_TOKEN; i++)
                    if (strcmp(pStart, &tokens[i][0]) == 0)
                        token = i;
                // Get value depending on token
                pStart = pEnd + 1;
                if (token == 0)
                {
                    // Name
                    pStart = strstr(pStart, "\"") + 1;
                    pEnd = strstr(pStart, "\"");
                    pEnd[0] = '\0';
                    strcpy_s(temp.name, pStart);
                    pStart = pEnd + 1;
                    temp.validItem = TRUE;  // once the name is found then it is valid
                }
                else if (token == 1)
                {
                    // Category
                    pStart = strstr(pStart, "\"") + 1;
                    pEnd = strstr(pStart, "\"");
                    pEnd[0] = '\0';
                    if (strcmp("Appetizer", pStart) == 0) temp.cat = APP;
                    else if (strcmp("Drink", pStart) == 0)  temp.cat = DRINK;
                    else if (strcmp("Main Course", pStart) == 0)  temp.cat = MAIN;
                    else if (strcmp("Dessert", pStart) == 0)  temp.cat = DESSERT;
                    pStart = pEnd + 1;
                }
                else if (token == 2 || token == 3 || token == 4)
                {
                    // Cost / Sat / Budget
                    pStart = strstr(pStart, ":") + 1;
                    pEnd = pStart;
                    for (i=0;;i++)  // loop forever
                    {
                        if (pStart[i] == '\n' || pStart[i] == ',' || pStart[i] == '}' || pStart[i] == '\0')
                        {
                            for (j = 0; j < i; j++)
                                tempstr[j] = pStart[j];
                            tempstr[i] = '\0';
                            pEnd = pStart + i;
                            break;
                        }
                    }
                    if (token == 2)
                        sscanf_s(tempstr, "%d", &temp.cost);
                    else if (token == 3)
                        sscanf_s(tempstr, "%d", &temp.sat);
                    else // 4
                        sscanf_s(pStart, "%d", &budget);
                }
                else if (token == 5)
                {
                    // Foods
                    endLine = TRUE;  // Nothing to do with foods token.  Just go to next line.
                }
            }
            else
            {
                // Check to see if there is an end to the item on this line
                if (strstr(pStart, "}") != NULL && temp.validItem)
                    endItem = TRUE;
                endLine = TRUE;  // no content on line.  Go to next line.
            }
        }

        // Save the item?
        if (endItem)
        {
            // Solution 2: NO SORTING FOR NOW
            switch (temp.cat)
            {
            case APP:       if (*p_a < 50) { memcpy(&menu[APP][*p_a], &temp, sizeof(Item)); (*p_a)++; } break;
            case DRINK:     if (*p_dr < 50) { memcpy(&menu[DRINK][*p_dr], &temp, sizeof(Item)); (*p_dr)++; } break;
            case MAIN:      if (*p_m < 50) { memcpy(&menu[MAIN][*p_m], &temp, sizeof(Item)); (*p_m)++; } break;
            case DESSERT:   if (*p_de < 50) { memcpy(&menu[DESSERT][*p_de], &temp, sizeof(Item)); (*p_de)++; } break;
            }
            memset(&temp, 0, sizeof(Item));
            endItem = FALSE;
        }
    }
    fclose(pFile);
    return(budget);
}

// Output the result to the display in json format
void OutputResult(bool success, Item selected[], int satisfaction)
{
    int i;
    int cost = 0;

    printf("\n");
    if (!success)
    {
        printf("{\n  \"error\": \"I would never eat here as the menu is very limited.  Alek, how about Chinese instead?\"\n}\n");
    }
    else if (satisfaction == 0)
    {
        printf("{\n  \"error\": \"I can't get no....sat-is-fac-tion!  Rolling Stones.  Alek, I need more money to eat here.\"\n}\n");
    }
    else
    {
        // Should be a valid result
        // Check to make sure all menu items are valid
        for (i = 0; i < 4; i++)
        {
            if (!selected[i].validItem)
            {
                printf("{]n  \"error\": \"Invalid result.  I obviously lost my coding touch long ago.\"\n}\n");
                return;
            }
            cost += selected[i].cost;
        }

        // All items are valid
        printf("{\n  \"selectedFoods\": [\"%s\", \"%s\", \"%s\", \"%s\"],\n", selected[0].name, selected[1].name, selected[2].name, selected[3].name);
        printf("  \"totalCost\": %d,\n", cost);
        printf("  \"totalSatisfaction\": %d\n}\n", satisfaction);
    }
    return;
}
