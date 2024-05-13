requires GOTMPDIR setting  
bash: `export GOTMPDIR=~/Projects/tmp`  
cmd: `set GOTMPDIR=%USERPROFILE%\Projects\tmp`  
powershell: `$env:GOTMPDIR="$env:USERPROFILE\Projects\tmp"`

# Table of Contents
1. [Command-line](#command-line)
    * [Powershell](#powershell)
2. [Implementation Notes](#implementation-notes)
    * [Missing Categories](#missing-categories)
3. [Randomized Max Test Case](#randomized-max-test-case)

# Command-line
*As spec'd, with no command-line parameters, the solution expects menu.json to exist in the current directory, and only outputs bestMeal json or error json to the console.*

**@Alek: for validating the solution, it is assumed you won't use any command-line parameters, or at most `-f={testCaseMenuJsonFilespec}`**

to use a different file for input (such as the other .json test files provided), 
use the `-f={altMenuJsonFilespec}` parameter.

to enable debug output, use `-v`.  even with this enabled, you can redirect the output to a file to capture only the json result.

with the following command-line, an alternative menu file is used for input, verbose mode is enabled with output written to stderr, result.json will contain only the bestmeal or error json output.  

`go run JeffR_BestMeal_Solution.go -f=menuFull.json -v >result.json`  

to redirect the debug logging to a separate file, add `2>{logFilespec}`:  

`go run JeffR_BestMeal_Solution.go -f=menuFull.json -v >result.json 2>out.log`

(works in bash/cmd)

[TOC](#table-of-contents)

## Powershell

powershell requires the -f to be prefixed with path and quoted; redirecting log output in powershell produces weird error in the log that can be ignored.

`go run JeffR_BestMeal_Solution.go -f=".\menuFull.json" -v >result.json 2>out.log`

[TOC](#table-of-contents)

# Implementation notes
I had originally intended to march thru the inputs to assemble all possible meals with total cost and satisfaction for each, then sort the list, then find what's within budget... then I realized I only needed to track the most satisfying meal

checking for equal-satisfaction-but-lower-cost meals is included in the algorithm- IMO, spending less money is more satisfying

categories are stripped of whitespace and converted to lower case and stripped of pluralization, so 

- `"Main Courses"`
- `"mAiNcOuRsE"`
- `"   M A I N   C O U R S E   "`  

are all valid; use -c=false or change the CLEANSE_CATEGORIES var near the top of the solution .go file to disable this

There is no distinction between valid json that is otherwise empty (literally, `{}`) vs a foods array that is empty- both produce the same "No food..." error output

disabled-by-default debug output via logger helps verify the results (see `-v` command-line parameter)

I ignored the complaints in VSCode editor re errors starting with a capital letter. I'm sure they have their reasons.

[TOC](#table-of-contents)

## Missing categories
when detecting missing categories, only the first missing category is reported in order of precedence:  
a. Appetizer  
b. Drink  
c. Main Course  
d. Dessert  

[TOC](#table-of-contents)

# Randomized Max Test Case
code was written to generate a test case with max items per category with randomized-yet-unique values 1-50 for cost and satisfaction within each category and randomized order of items; 

the single item with max satisfaction in each category has the same name as the (presumably) best meal 
answer for the full menu example, all other names are obviously fake but unique;

the budget of $200 for this menu means we'll always find the max satisfaction of 200.

each time you run the generator code, the output will be different, but the best meal will always 
consist of the same items, satisfaction=200, and randomized cost.

the separate `*Generator*.go` routine was renamed to `*_go` to avoid collisions / complaints in the VSCode editor, tho after renaming it back to `*.go`, it and the solution technically have no problems co-existing.

see menuRandMax.json

[TOC](#table-of-contents)
