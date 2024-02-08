package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
)

var (
	biomes = []string{
		"in a forest",
		"in a pine forest",
		"knee deep in a swamp",
		"in a mountain range",
		"in a desert",
		"in a grassy plain",
		"in a frozen tundra",
	}
)

func hasTrees(biome int) bool {
	return biome < 3
}

func hasStone(biome int) bool {
	return biome == 3
}

func hasRivers(biome int) bool {
	return biome != 2 && biome != 4
}

type ToolType int

const (
	NoneToolType ToolType = iota
	Pick
	Sword
	Shovel
)

type Item struct {
	droppable bool
	desc      string
	heavy     bool
	creature  bool
	drops     []string
	aliases   []string
	hitDrops  []string
	monster   bool
	nocturnal bool
	material  bool
	tool      bool
	toolLevel int
	toolType  ToolType
	ore       bool
	infinite  bool
	food      bool
}

var (
	items = map[string]Item{
		"no tea": {
			droppable: false,
			desc:      "Pull youreslf together man.",
		},
		"a pig": {
			heavy:    true,
			creature: true,
			drops:    []string{"some pork"},
			aliases:  []string{"pig"},
			desc:     "The pig has a square nose.",
		},
		"a cow": {
			heavy:    true,
			creature: true,
			aliases:  []string{"cow"},
			desc:     "The cow stares at you blankly.",
		},
		"a sheep": {
			heavy:    true,
			creature: true,
			hitDrops: []string{"some wool"},
			aliases:  []string{"sheep"},
			desc:     "The sheep is fluffy.",
		},
		"a chicken": {
			heavy:    true,
			creature: true,
			drops:    []string{"some chicken"},
			aliases:  []string{"chicken"},
			desc:     "The chicken looks delicious.",
		},
		"a creeper": {
			heavy:    true,
			creature: true,
			monster:  true,
			aliases:  []string{"creeper"},
			desc:     "The creeper needs a hug.",
		},
		"a skeleton": {
			heavy:     true,
			creature:  true,
			monster:   true,
			aliases:   []string{"skeleton"},
			nocturnal: true,
			desc:      "The head bone's connected to the neck bone, the neck bone's connected to the chest bone, the chest bone's connected to the arm bone, the arm bone's connected to the bow, and the bow is pointed at you.",
		},
		"a zombie": {
			heavy:     true,
			creature:  true,
			monster:   true,
			aliases:   []string{"zombie"},
			nocturnal: true,
			desc:      "All he wants to do is eat your brains.",
		},
		"a spider": {
			heavy:    true,
			creature: true,
			monster:  true,
			aliases:  []string{"spider"},
			desc:     "Dozens of eyes stare back at you.",
		},
		"a cave entrance": {
			heavy:   true,
			aliases: []string{"cave entrance", "cave", "entrance"},
			desc:    "The entrance to the cave is dark, but it looks like you can climb down.",
		},
		"an exit to the surface": {
			heavy:   true,
			aliases: []string{"exit to the surface", "exit", "opening"},
			desc:    "You can just see the sky through the opening.",
		},
		"a river": {
			heavy:   true,
			aliases: []string{"river"},
			desc:    "The river flows majestically towards the horizon. It doesn't do anything else.",
		},
		"some wood": {
			aliases:  []string{"wood"},
			material: true,
			desc:     "You could easily craft this wood into planks.",
		},
		"some planks": {
			aliases: []string{"planks", "wooden planks", "wood planks"},
			desc:    "You could easily craft these planks into sticks.",
		},
		"some sticks": {
			aliases: []string{"sticks", "wooden sticks", "wood sticks"},
			desc:    "A perfect handle for torches or a pickaxe.",
		},
		"a crafting table": {
			aliases: []string{"crafting table", "craft table", "work bench", "workbench", "crafting bench", "table"},
			desc:    "It's a crafting table. I shouldn't tell you this, but these don't actually do anything in this game, you can craft tools whenever you like.",
		},
		"a furnace": {
			aliases: []string{"furnace"},
			desc:    "It's a furnace. Between you and me, these don't actually do anything in this game.",
		},
		"a wooden pickaxe": {
			aliases:   []string{"pickaxe", "pick", "wooden pick", "wooden pickaxe", "wood pick", "wood pickaxe"},
			tool:      true,
			toolLevel: 2,
			toolType:  Pick,
			desc:      "The pickaxe looks good for breaking stone and coal.",
		},
		"a stone pickaxe": {
			aliases:   []string{"pickaxe", "pick", "stone pick", "stone pickaxe"},
			tool:      true,
			toolLevel: 2,
			toolType:  Pick,
			desc:      "The pickaxe looks good for breaking iron.",
		},
		"an iron pickaxe": {
			aliases:   []string{"pickaxe", "pick", "iron pick", "iron pickaxe"},
			tool:      true,
			toolLevel: 3,
			toolType:  Pick,
			desc:      "The pickaxe looks strong enough to break diamond.",
		},
		"a diamond pickaxe": {
			aliases:   []string{"pickaxe", "pick", "diamond pick", "diamond pickaxe"},
			tool:      true,
			toolLevel: 4,
			toolType:  Pick,
			desc:      "Best. Pickaxe. Ever.",
		},
		"a wooden sword": {
			aliases:   []string{"sword", "wooden sword", "wood sword"},
			tool:      true,
			toolLevel: 1,
			toolType:  Sword,
			desc:      "Flimsy, but better than nothing.",
		},
		"a stone sword": {
			aliases:   []string{"sword", "stone sword"},
			tool:      true,
			toolLevel: 2,
			toolType:  Sword,
			desc:      "A pretty good sword.",
		},
		"an iron sword": {
			aliases:   []string{"sword", "iron sword"},
			tool:      true,
			toolLevel: 3,
			toolType:  Sword,
			desc:      "This sword can slay any enemy.",
		},
		"a diamond sword": {
			aliases:   []string{"sword", "diamond sword"},
			tool:      true,
			toolLevel: 4,
			toolType:  Sword,
			desc:      "Best. Sword. Ever.",
		},
		"a wooden shovel": {
			aliases:   []string{"shovel", "wooden shovel", "wood shovel"},
			tool:      true,
			toolLevel: 1,
			toolType:  Shovel,
			desc:      "Good for digging holes.",
		},
		"a stone shovel": {
			aliases:   []string{"shovel", "stone shovel"},
			tool:      true,
			toolLevel: 2,
			toolType:  Shovel,
			desc:      "Good for digging holes.",
		},
		"an iron shovel": {
			aliases:   []string{"shovel", "iron shovel"},
			tool:      true,
			toolLevel: 3,
			toolType:  Shovel,
			desc:      "Good for digging holes.",
		},
		"a diamond shovel": {
			aliases:   []string{"shovel", "diamond shovel"},
			tool:      true,
			toolLevel: 4,
			toolType:  Shovel,
			desc:      "Good for digging holes.",
		},
		"some coal": {
			aliases:   []string{"coal"},
			ore:       true,
			toolLevel: 1,
			toolType:  Pick,
			desc:      "That coal looks useful for building torches, if only you had a pickaxe to mine it.",
		},
		"some dirt": {
			aliases:  []string{"dirt"},
			material: true,
			desc:     "Why not build a mud hut?",
		},
		"some stone": {
			aliases:   []string{"stone", "cobblestone"},
			material:  true,
			ore:       true,
			infinite:  true,
			toolLevel: 1,
			toolType:  Pick,
			desc:      "Stone is useful for building things, and making stone pickaxes.",
		},
		"some iron": {
			aliases:   []string{"iron"},
			material:  true,
			ore:       true,
			toolLevel: 2,
			toolType:  Pick,
			desc:      "That iron looks might strong, you'll need a stone pickaxe to mine it.",
		},
		"some diamond": {
			aliases:   []string{"diamond", "diamonds"},
			material:  true,
			ore:       true,
			toolLevel: 3,
			toolType:  Pick,
			desc:      "Sparkly, rare, and impossible to mine without an iron pickaxe.",
		},
		"some torches": {
			aliases: []string{"torches", "torch"},
			desc:    "These won't run out of a while.",
		},
		"a torch": {
			aliases: []string{"torch"},
			desc:    "Fire, fire, burn so bright, won't you light my cave tonight?",
		},
		"some wool": {
			aliases:  []string{"wool"},
			material: true,
			desc:     "Soft and good for building.",
		},
		"some pork": {
			aliases: []string{"pork", "porkchops"},
			food:    true,
			desc:    "Delicious and nutricious.",
		},
		"some chicken": {
			aliases: []string{"chicken"},
			food:    true,
			desc:    "Finger licking good.",
		},
	}
	animals = []string{
		"a pig", "a cow", "a sheep", "a chicken",
	}
	monsters = []string{
		"a creeper", "a skeleton", "a zombie", "a spider",
	}
	recipes = map[string][]string{
		"some planks":      []string{"some wood"},
		"some sticks":      []string{"some planks"},
		"a crafting table": []string{"some planks"},
		"a furnace":        []string{"some stone"},
		"some torches":     []string{"some sticks", "some coal"},

		"a wooden pickaxe":  []string{"some planks", "some sticks"},
		"a stone pickaxe":   []string{"some stone", "some sticks"},
		"an iron pickaxe":   []string{"some iron", "some sticks"},
		"a diamond pickaxe": []string{"some diamond", "some sticks"},

		"a wooden sword":  []string{"some planks", "some sticks"},
		"a stone sword":   []string{"some stone", "some sticks"},
		"an iron sword":   []string{"some iron", "some sticks"},
		"a diamond sword": []string{"some diamond", "some sticks"},

		"a wooden shovel":  []string{"some planks", "some sticks"},
		"a stone shovel":   []string{"some stone", "some sticks"},
		"an iron shovel":   []string{"some iron", "some sticks"},
		"a diamond shovel": []string{"some diamond", "some sticks"},
	}

	goWest = []string{
		"(life is peaceful there)",
		"(lots of open air)",
		"(to begin life anew)",
		"(this is what we'll do)",
		"(sun in winter time)",
		"(we will do just fine)",
		"(where the skies are blue)",
		"(this and more we'll do)",
	}
	nGoWest   = 0
	running   = true
	x         = 0
	y         = 0
	z         = 0
	inventory = map[string]Item{
		"no tea": items["no tea"],
	}
	turn       = 0
	timeInRoom = 0
	injured    = false
	dayCycle   = []string{
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"It is daytime.",
		"The sun is setting.",
		"It is night.",
		"It is night.",
		"It is night.",
		"It is night.",
		"It is night.",
		"The sun is rising.",
	}
	roomMap = map[int]map[int]map[int]Room{}
)

type Room struct {
	biome int
	trees bool
	items map[string]Item
	exits struct {
		north bool
		south bool
		east  bool
		west  bool
		down  bool
		up    bool
	}
	dark     bool
	monsters int
	valid    bool
}

func getTimeOfDay() float64 {
	return math.Mod(float64(turn/3), float64(len(dayCycle))) + 1.0
}

func isSunny() bool {
	return getTimeOfDay() < 10
}

func getRoom(x int, y int, z int, dontCreate bool) Room {
	xVal, ok := roomMap[x]

	if !ok {
		xVal = make(map[int]map[int]Room)
	}

	yVal, ok := xVal[y]

	if !ok {
		yVal = make(map[int]Room)
	}

	_, ok = yVal[z]

	if !ok && !dontCreate {
		room := Room{
			items: make(map[string]Item),
			exits: struct {
				north bool
				south bool
				east  bool
				west  bool
				down  bool
				up    bool
			}{},
			monsters: 0,
		}
		roomMap[x][y][z] = Room{}

		if y == 0 {
			room.biome = rand.Intn(len(biomes))
			room.trees = hasTrees(room.biome)

			if rand.Intn(3) == 0 {
				for i := 0; i < rand.Intn(1); i++ {
					animal := animals[rand.Intn(len(animals))]
					room.items[animal] = items[animal]
				}
			}

			if rand.Intn(5) == 0 || hasStone(room.biome) {
				room.items["some stone"] = items["some stone"]
			}

			if rand.Intn(8) == 0 {
				room.items["some coal"] = items["some coal"]
			}

			if rand.Intn(8) == 0 && hasRivers(room.biome) {
				room.items["a river"] = items["a river"]
			}

			room.exits.north = true
			room.exits.south = true
			room.exits.east = true
			room.exits.west = true

			if rand.Intn(8) == 0 {
				room.exits.down = true
				room.items["a cave entrance"] = items["a cave entrance"]
			}

			room.valid = true
		} else {
			tryExit := func(sDir string, sOpp string, x int, y int, z int) {
				adj := getRoom(x, y, z, true)

				if adj.valid {
					switch sOpp {
					case "north":
						if adj.exits.north {
							switch sDir {
							case "north":
								room.exits.north = true
							case "south":
								room.exits.south = true
							case "west":
								room.exits.west = true
							case "east":
								room.exits.east = true
							case "down":
								room.exits.down = true
							case "up":
								room.exits.up = true
							}
						}
					case "south":
						if adj.exits.south {
							switch sDir {
							case "north":
								room.exits.north = true
							case "south":
								room.exits.south = true
							case "west":
								room.exits.west = true
							case "east":
								room.exits.east = true
							case "down":
								room.exits.down = true
							case "up":
								room.exits.up = true
							}
						}
					case "west":
						if adj.exits.west {
							switch sDir {
							case "north":
								room.exits.north = true
							case "south":
								room.exits.south = true
							case "west":
								room.exits.west = true
							case "east":
								room.exits.east = true
							case "down":
								room.exits.down = true
							case "up":
								room.exits.up = true
							}
						}
					case "east":
						if adj.exits.east {
							switch sDir {
							case "north":
								room.exits.north = true
							case "south":
								room.exits.south = true
							case "west":
								room.exits.west = true
							case "east":
								room.exits.east = true
							case "down":
								room.exits.down = true
							case "up":
								room.exits.up = true
							}
						}
					case "down":
						if adj.exits.down {
							switch sDir {
							case "north":
								room.exits.north = true
							case "south":
								room.exits.south = true
							case "west":
								room.exits.west = true
							case "east":
								room.exits.east = true
							case "down":
								room.exits.down = true
							case "up":
								room.exits.up = true
							}
						}
					}
				} else {
					if rand.Intn(3) == 0 {
						switch sDir {
						case "north":
							room.exits.north = true
						case "south":
							room.exits.south = true
						case "west":
							room.exits.west = true
						case "east":
							room.exits.east = true
						case "down":
							room.exits.down = true
						case "up":
							room.exits.up = true
						}
					}
				}
			}

			if y == -1 {
				above := getRoom(x, y+1, z, false)

				if above.exits.down {
					room.exits.up = true
					room.items["an exit to the surface"] = items["an exit to the surface"]
				}
			} else {
				tryExit("up", "down", x, y+1, z)
			}

			if y > -3 {
				tryExit("down", "up", x, y-1, z)
			}

			tryExit("east", "west", x-1, y, z)
			tryExit("west", "east", x+1, y, z)
			tryExit("north", "south", x, y, z+1)
			tryExit("south", "north", x, y, z-1)

			room.items["some stone"] = items["some stone"]

			if rand.Intn(3) == 0 {
				room.items["some coal"] = items["some coal"]
			}

			if rand.Intn(8) == 0 {
				room.items["some iron"] = items["some iron"]
			}

			if y == -3 && rand.Intn(15) == 0 {
				room.items["some diamond"] = items["some diamond"]
			}

			room.dark = true
		}
	}

	return roomMap[x][y][z]
}

func itemize(t []int) string {
	if len(t) == 0 {
		return "nothing"
	}

	text := ""

	for i := 0; i < len(t); i++ {
		text += strconv.Itoa(t[i])

		if i < len(t)-1 {
			if i < len(t)-2 {
				text += " and "
			} else {
				text += ", "
			}
		}
	}

	return text
}

var (
	matches = map[string][]string{
		"wait": {"wait"},
		"look": {
			"look at the ([%a ]+)",
			"look at ([%a ]+)",
			"look",
			"inspect ([%a ]+)",
			"inspect the ([%a ]+)",
			"inspect",
		},
		"inventory": {
			"check self",
			"check inventory",
			"inventory",
			"i",
		},
		"go": {
			"go (%a+)",
			"travel (%a+)",
			"walk (%a+)",
			"run (%a+)",
			"go",
		},
		"dig": {
			"dig (%a+) using ([%a ]+)",
			"dig (%a+) with ([%a ]+)",
			"dig (%a+)",
			"dig",
		},
		"take": {
			"pick up the ([%a ]+)",
			"pick up ([%a ]+)",
			"pickup ([%a ]+)",
			"take the ([%a ]+)",
			"take ([%a ]+)",
			"take",
		},
		"drop": {
			"put down the ([%a ]+)",
			"put down ([%a ]+)",
			"drop the ([%a ]+)",
			"drop ([%a ]+)",
			"drop",
		},
		"place": {
			"place the ([%a ]+)",
			"place ([%a ]+)",
			"place",
		},
		"cbreak": {
			"punch the ([%a ]+)",
			"punch ([%a ]+)",
			"punch",
			"break the ([%a ]+) with the ([%a ]+)",
			"break ([%a ]+) with ([%a ]+)",
			"break the ([%a ]+)",
			"break ([%a ]+)",
			"break",
		},
		"mine": {
			"mine the ([%a ]+) with the ([%a ]+)",
		},
	}
)

func main() {
	fmt.Println("Hello world!")
}
