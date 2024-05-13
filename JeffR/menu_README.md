use -f={menuJsonFileSpec} parameter to override the default menu.json e.g.
`go run JeffR_BestMeal_Solution.go -f=menuRandMax.json`

| Test .json | Notes | Result |
| ---   | ---   | ---    |
| menu | a small example derived from Alek's README (subset of full) | Fried Calamari, Beer, Lasagna, Cheesecake, tC=21, tS=18 |
| menuBadCategory | category=Cheese | error="Unknown food category..." |
| menuEmptyFile | no content, zero length | error="Bad menu json: unexpected end of input" |
| menuEmptyJson | {} | error="No food in menu??" | 
| menuEmptyJsonArray | [] | error="Bad menu ... cannot unmarshal" |
| menuFull  | the full example from Alek's README | Bruschetta, Coffee, Steak, Tiramisu, tC=24, tS=22** |
| menuLooseCategory | case/whitespace within category is ignored | **==menuFull |
| menuMalformed | [{]} (bad json) | error="Bad menu ... invalid character" |
| menuMismatchedJson | valid json but doesn't match menu foods+budget structure | error="No food in menu??" |
| menuNoApps | no food w/ category=="Appetizer" | error="No apps in menu" |
| menuNoBudget | missing budget field in json | error="You need a budget" (same when budget < 0) |
| menuNoDessert | no food w/ category=="Dessert" | error="No dessert in menu" |
| menuNoDrinks | no food w/ category=="Drink" | error="No drink in menu" |
| menuNoFoods | foods==[] | error="No food in menu??" |
| menuNoMains | no food w/ category=="Main Course" | error="No mains in menu" |
| menuNullJson | file content is literally == "null" sans quotes | error="No food in menu??" |
| menuRandMax | randomized set of data with 50 items per category | **==menuFull, tC=76, tS=200 |
| menuReversedJson | json fields are in reverse order | **==menuFull |
| menuTooPoor1 | no meals fit budget; matches Alek's readme example | error="...you need another 35 buck(s)..." |
| menuTooPoor2 | no meals fit budget | error="...you need another 1 buck(s)..." |
| 
