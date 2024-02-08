package main

import (
	"fmt"
	"math"
	"math/rand"
	"regexp"
	"strconv"

	"github.com/TwiN/go-color"
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
	biome    int
	trees    bool
	items    map[string]Item
	exits    Exits
	dark     bool
	monsters int
	valid    bool
}

type Exits struct {
	north bool
	south bool
	east  bool
	west  bool
	down  bool
	up    bool
}

func (e Exits) getExit(s string) bool {
	switch s {
	case "north":
		return e.north
	case "south":
		return e.south
	case "west":
		return e.west
	case "east":
		return e.east
	case "up":
		return e.up
	case "down":
		return e.down
	default:
		return false
	}
}

func (e *Exits) setExit(s string, v bool) {
	switch s {
	case "north":
		e.north = v
	case "south":
		e.south = v
	case "west":
		e.west = v
	case "east":
		e.east = v
	case "up":
		e.up = v
	case "down":
		e.down = v
	}
}

func (r Room) getExits() []string {
	exits := []string{}

	if r.exits.north {
		exits = append(exits, "north")
	}

	if r.exits.south {
		exits = append(exits, "south")
	}

	if r.exits.west {
		exits = append(exits, "west")
	}

	if r.exits.east {
		exits = append(exits, "east")
	}

	if r.exits.up {
		exits = append(exits, "up")
	}

	if r.exits.down {
		exits = append(exits, "down")
	}

	return exits
}

func getTimeOfDay() float64 {
	return math.Mod(float64(turn/3), float64(len(dayCycle))) + 1.0
}

func isSunny() bool {
	return getTimeOfDay() < 10
}

type RoomCoord struct {
	x int
	y int
	z int
}

func getRoom(x int, y int, z int, dontCreate bool) RoomCoord {
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
				coords := getRoom(x, y, z, true)
				adj := roomMap[coords.x][coords.y][coords.z]

				if adj.valid {
					room.exits.setExit(sDir, adj.exits.getExit(sOpp))
				} else {
					if rand.Intn(3) == 0 {
						room.exits.setExit(sDir, true)
					}
				}
			}

			if y == -1 {
				coords := getRoom(x, y+1, z, false)
				above := roomMap[coords.x][coords.y][coords.z]

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

		roomMap[x][y][z] = room
	}

	return RoomCoord{x: x, y: y, z: z}
}

func itemizeNum(t []int) string {
	ts := []string{}

	for _, v := range t {
		ts = append(ts, strconv.Itoa(v))
	}

	return itemizeStr(ts)
}

func itemizeStr(t []string) string {
	if len(t) == 0 {
		return "nothing"
	}

	text := ""

	for i := 0; i < len(t); i++ {
		text += t[i]

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
			"look at the ([A-z ]+)",
			"look at ([A-z ]+)",
			"look",
			"inspect ([A-z ]+)",
			"inspect the ([A-z ]+)",
			"inspect",
		},
		"inventory": {
			"check self",
			"check inventory",
			"inventory",
			"i",
		},
		"go": {
			"go ([A-z]+)",
			"travel ([A-z]+)",
			"walk ([A-z]+)",
			"run ([A-z]+)",
			"go",
		},
		"dig": {
			"dig ([A-z]+) using ([A-z ]+)",
			"dig ([A-z]+) with ([A-z ]+)",
			"dig ([A-z]+)",
			"dig",
		},
		"take": {
			"pick up the ([A-z ]+)",
			"pick up ([A-z ]+)",
			"pickup ([A-z ]+)",
			"take the ([A-z ]+)",
			"take ([A-z ]+)",
			"take",
		},
		"drop": {
			"put down the ([A-z ]+)",
			"put down ([A-z ]+)",
			"drop the ([A-z ]+)",
			"drop ([A-z ]+)",
			"drop",
		},
		"place": {
			"place the ([A-z ]+)",
			"place ([A-z ]+)",
			"place",
		},
		"cbreak": {
			"punch the ([A-z ]+)",
			"punch ([A-z ]+)",
			"punch",
			"break the ([A-z ]+) with the ([A-z ]+)",
			"break ([A-z ]+) with ([A-z ]+)",
			"break the ([A-z ]+)",
			"break ([A-z ]+)",
			"break",
		},
		"mine": {
			"mine the ([A-z ]+) with the ([A-z ]+)",
			"mine ([A-z ]+) with ([A-z ]+)",
			"mine ([A-z ]+)",
			"mine",
		},
		"attack": {
			"attack the ([A-z ]+) with the ([A-z ]+)",
			"attack ([A-z ]+) with ([A-z ]+)",
			"attack ([A-z ]+)",
			"attack",
			"kill the ([A-z ]+) with the ([A-z ]+)",
			"kill ([A-z ]+) with ([A-z ]+)",
			"kill ([A-z ]+)",
			"kill",
			"hit the ([A-z ]+) with the ([A-z ]+)",
			"hit ([A-z ]+) with ([A-z ]+)",
			"hit ([A-z ]+)",
			"hit",
		},
		"craft": {
			"craft a ([A-z ]+)",
			"craft some ([A-z ]+)",
			"craft ([A-z ]+)",
			"craft",
			"make a ([A-z ]+)",
			"make some ([A-z ]+)",
			"make ([A-z ]+)",
			"make",
		},
		"build": {
			"build ([A-z ]+) out of ([A-z ]+)",
			"build ([A-z ]+) from ([A-z ]+)",
			"build ([A-z ]+)",
			"build",
		},
		"eat": {
			"eat a ([A-z ]+)",
			"eat the ([A-z ]+)",
			"eat ([A-z ]+)",
			"eat",
		},
		"help": {
			"help me",
			"help",
		},
		"exit": {
			"exit",
			"quit",
			"goodbye",
			"good bye",
			"bye",
			"farewell",
		},
	}

	commands = map[string]func([]string){
		"noinput": func(_ []string) {
			responses := []string{
				"Speak up.",
				"Enunciate.",
				"Project your voice.",
				"Don't be shy.",
				"Use your words.",
			}

			fmt.Println(randomChoice(responses))
		},
		"badinput": func(_ []string) {
			responses := []string{
				"I don't understand.",
				"I don't understand you.",
				"You can't do that.",
				"Nope.",
				"Huh?",
				"Say again?",
				"That's crazy talk.",
				"Speak clearly.",
				"I'll think about it.",
				"Let me get back to you on that one.",
				"That doens't make any sense.",
				"What?",
			}

			fmt.Println(randomChoice(responses))
		},
		"wait": func(_ []string) {
			fmt.Println("Time passes...")
		},
		"look": func(vals []string) {
			var target string

			if len(vals) == 0 {
				target = ""
			} else {
				target = vals[0]
			}

			coords := getRoom(x, y, z, false)
			room := roomMap[coords.x][coords.y][coords.z]

			if room.dark {
				fmt.Println("It is pitch dark.")
				return
			}

			if target == "" {
				if y == 0 {
					fmt.Printf("You are standing %s. ", biomes[room.biome])
					fmt.Println(dayCycle[int(getTimeOfDay())])
				} else {
					fmt.Print("You are underground. ")
					exits := room.getExits()

					if len(exits) != 0 {
						fmt.Printf("You can travel %s.", itemizeStr(exits))
					} else {
						fmt.Println()
					}
				}

				if len(room.items) > 0 {
					items := []string{}

					for i := range room.items {
						items = append(items, i)
					}

					fmt.Printf("There is %s here.\n", itemizeStr(items))
				}

				if room.trees {
					fmt.Println("There are trees here.")
				}
			} else {
				if room.trees && (target == "tree" || target == "trees") {
					fmt.Println("The trees look easy to break.")
				} else if target == "self" || target == "myself" {
					fmt.Println("Very handsome.")
				} else {
					item, ok := room.items[target]

					if !ok {
						item, ok = inventory[target]
					}

					if ok {
						if item.desc == "" {
							fmt.Printf("You see nothing special about %s.\n", target)
						} else {
							fmt.Println(item.desc)
						}
					} else {
						fmt.Printf("You don't see any %s here.\n", target)
					}
				}
			}
		},
		"go": func(vals []string) {
			var dir string

			if len(vals) == 0 {
				dir = ""
			} else {
				dir = vals[0]
			}

			coords := getRoom(x, y, z, false)
			room := roomMap[coords.x][coords.y][coords.z]

			if dir == "" {
				fmt.Println("Go where?")
				return
			}

			if nGoWest != -1 {
				if dir == "west" {
					nGoWest += 1

					if nGoWest > len(goWest) {
						nGoWest = 0
					}

					fmt.Println(goWest[nGoWest])
				} else {
					if nGoWest > 0 || turn > 6 {
						nGoWest = -1
					}
				}
			}

			if !room.exits.getExit(dir) {
				fmt.Println("You can't go that way.")
				return
			}

			if dir == "north" {
				z += 1
			} else if dir == "south" {
				z -= 1
			} else if dir == "east" {
				x -= 1
			} else if dir == "west" {
				x += 1
			} else if dir == "up" {
				y += 1
			} else if dir == "down" {
				y -= 1
			} else {
				fmt.Println("I don't understand that direction.")
				return
			}

			timeInRoom = 0
			lookComm([]string{})
		},
		"dig": func(vals []string) {
			var dir string
			var tool string

			if len(vals) == 0 {
				dir = ""
				tool = ""
			} else if len(vals) == 1 {
				dir = vals[0]
				tool = ""
			} else {
				dir = vals[0]
				tool = vals[1]
			}

			coords := getRoom(x, y, z, false)
			room := roomMap[coords.x][coords.y][coords.z]

			if dir == "" {
				fmt.Println("Dig where?")
				return
			}

			var iTool Item
			var fTool bool

			if tool != "" {
				iTool, fTool = inventory[tool]

				if !fTool {
					fmt.Printf("You're not carrying a %s.\n", tool)
					return
				}
			}

			actuallyDigging := !room.exits.getExit(dir)

			if actuallyDigging {
				if !fTool || iTool.toolType != Pick {
					fmt.Println("You need to use a pickaxe to dig through stone.")
					return
				}
			}

			setCoordExit := func(x int, y int, z int, exit string, val bool) {
				tmp := getRoom(x, y, z, false)
				tmp2 := roomMap[tmp.x][tmp.y][tmp.z]
				tmp2.exits.setExit(exit, val)
				roomMap[tmp.x][tmp.y][tmp.z] = tmp2
			}

			if dir == "north" {
				room.exits.north = true
				z += 1
				setCoordExit(x, y, z, "south", true)
			} else if dir == "south" {
				room.exits.south = true
				z -= 1
				setCoordExit(x, y, z, "north", true)
			} else if dir == "east" {
				room.exits.east = true
				x -= 1
				setCoordExit(x, y, z, "west", true)
			} else if dir == "west" {
				room.exits.west = true
				x += 1
				setCoordExit(x, y, z, "east", true)
			} else if dir == "up" {
				if y == 0 {
					fmt.Println("You can't dig that way.")
					return
				}

				room.exits.up = true

				if y == -1 {
					room.items["an exit to the surface"] = items["an exit to the surface"]
				}

				y += 1
				coords1 := getRoom(x, y, z, false)
				room1 := roomMap[coords1.x][coords1.y][coords1.z]
				room1.exits.down = true

				if y == 0 {
					room1.items["a cave entrance"] = items["a cave entrance"]
				}

				roomMap[coords1.x][coords1.y][coords1.z] = room
			} else if dir == "down" {
				if y <= -3 {
					fmt.Println("You hit bedrock.")
					return
				}

				room.exits.down = true

				if y == 0 {
					room.items["a cave entrance"] = items["a cave entrance"]
				}

				y -= 1

				coords1 := getRoom(x, y, z, false)
				room1 := roomMap[coords1.x][coords1.y][coords1.z]
				room1.exits.up = true

				if y == -1 {
					room.items["an exit to the surface"] = items["an exit to the surface"]
				}

				roomMap[coords1.x][coords1.y][coords1.z] = room1
			} else {
				fmt.Println("I don't understand that direction.")
				return
			}

			if actuallyDigging {
				if (dir == "down" && y == -1) || (dir == "up" && y == 0) {
					inventory["some dirt"] = items["some dirt"]
					inventory["some stone"] = items["some stone"]
					fmt.Printf("You dig %s using %s and collect some dirt and stone.\n", dir, tool)
				} else {
					inventory["some stone"] = items["some stone"]
					fmt.Printf("You dig %s using %s and collect some stone.\n", dir, tool)
				}
			}

			timeInRoom = 0
			lookComm([]string{})
			roomMap[coords.x][coords.y][coords.z] = room
		},
		"inventory": func(_ []string) {
			vals := []string{}

			for i := range inventory {
				vals = append(vals, i)
			}

			fmt.Printf("You are carrying %s.\n", itemizeStr(vals))
		},
		"drop": func(vals []string) {
			var item string

			if len(vals) == 0 {
				item = ""
			} else {
				item = vals[0]
			}

			dropComm(item)
		},
		"place": func(vals []string) {
			var item string

			if len(vals) == 0 {
				item = ""
			} else {
				item = vals[0]
			}

			if item == "" {
				fmt.Println("Place what?")
				return
			}

			if item == "torch" || item == "a torch" {
				coords := getRoom(x, y, z, false)
				room := roomMap[coords.x][coords.y][coords.z]
				_, sTorches := inventory["some torches"]
				_, torch := inventory["a torch"]

				if sTorches || torch {
					delete(inventory, "a torch")
					room.items["a torch"] = items["a torch"]

					if room.dark {
						fmt.Println("The cave lights up under the torchflame.")
						room.dark = false
					} else if y == 0 && !isSunny() {
						fmt.Println("The night gets a little brighter.")
					} else {
						fmt.Println("Placed.")
					}
				} else {
					fmt.Println("You don't have torches.")
				}

				roomMap[coords.x][coords.y][coords.z] = room
				return
			}

			dropComm(item)
		},
		"take": func(vals []string) {
			var item string

			if len(vals) == 0 {
				item = ""
			} else {
				item = vals[0]
			}

			if item == "" {
				fmt.Println("Take what?")
				return
			}

			coords := getRoom(x, y, z, false)
			room := roomMap[coords.x][coords.y][coords.z]
			iItem, fItem := room.items[item]

			if fItem {
				if iItem.heavy {
					fmt.Printf("You can't carry %s.\n", item)
				} else if iItem.ore {
					fmt.Println("You need to mine this ore.")
				} else {
					if !iItem.infinite {
						delete(room.items, item)
					}

					inventory[item] = iItem

					_, sTorches := inventory["some torches"]
					_, torch := inventory["torch"]

					if sTorches && torch {
						delete(inventory, "a torch")
					}

					if item == "a torch" && y < 0 {
						room.dark = true
						fmt.Println("The cave plunges into darkness.")
					} else {
						fmt.Println("Taken.")
					}
				}
			} else {
				fmt.Printf("You don't see a %s here.\n", item)
			}

			roomMap[coords.x][coords.y][coords.z] = room
		},
		"mine": func(vals []string) {
			var item string
			var tool string

			if len(vals) == 0 {
				item = ""
				tool = ""
			} else if len(vals) == 1 {
				item = vals[0]
				tool = ""
			} else {
				item = vals[0]
				tool = vals[1]
			}

			if item == "" {
				fmt.Println("Mine what?")
				return
			}

			if tool == "" {
				fmt.Printf("Mine %s with what?\n", item)
				return
			}

			cbreakComm(item, tool)
		},
		"attack": func(vals []string) {
			var item string
			var tool string

			if len(vals) == 0 {
				item = ""
				tool = ""
			} else if len(vals) == 1 {
				item = vals[0]
				tool = ""
			} else {
				item = vals[0]
				tool = vals[1]
			}

			if item == "" {
				fmt.Println("Attack what?")
				return
			}

			cbreakComm(item, tool)
		},
		"cbreak": func(vals []string) {
			var item string
			var tool string

			if len(vals) == 0 {
				item = ""
				tool = ""
			} else if len(vals) == 1 {
				item = vals[0]
				tool = ""
			} else {
				item = vals[0]
				tool = vals[1]
			}

			cbreakComm(item, tool)
		},
	}
)

func doCommand(text string) {
	if text == "" {
		commands["noinput"]([]string{})
		return
	}

	for command, t := range matches {
		for _, match := range t {
			re := regexp.MustCompile("^" + match + "$")
			captures := re.FindStringSubmatch(text)

			if len(captures) != 0 {
				fnCommand := commands[command]

				if len(captures) == 1 && captures[0] == match {
					fnCommand([]string{})
				} else {
					fnCommand(captures[1:])
				}

				return
			}
		}
	}

	commands["badinput"]([]string{})
}

func lookComm(vals []string) {

}

func dropComm(item string) {
	if item == "" {
		fmt.Println("Drop what?")
		return
	}

	coords := getRoom(x, y, z, false)
	room := roomMap[coords.x][coords.y][coords.z]
	iItem, fItem := inventory[item]

	if fItem {
		if iItem.droppable {
			room.items[item] = iItem
			delete(inventory, item)
			fmt.Println("Dropped.")
		} else {
			fmt.Println("You can't drop that.")
		}
	} else {
		fmt.Printf("You don't have a %s.\n", item)
	}

	roomMap[coords.x][coords.y][coords.z] = room
}

func cbreakComm(item string, tool string) {
	if item == "" {
		fmt.Println("Break what?")
		return
	}

	var iTool Item
	var fTool bool

	if tool != "" {
		iTool, fTool = inventory[tool]

		if !fTool {
			fmt.Printf("You're not carrying a %s.\n", tool)
			return
		}
	}

	coords := getRoom(x, y, z, false)
	room := roomMap[coords.x][coords.y][coords.z]

	if item == "tree" || item == "trees" || item == "a tree" {
		fmt.Println("The tree breaks into blocks of wood, which you pick up.")
		inventory["some wood"] = items["some wood"]
		return
	} else if item == "self" || item == "myself" {
		fmt.Println(color.Ize(color.Red, "You have died."))
		running = false
	}

	iItem, fItem := room.items[item]

	if fItem {
		if iItem.ore {
			if !fTool {
				fmt.Println("You need a tool to break this ore.")
				return
			}

			if iTool.tool {
				if iTool.toolLevel < iItem.toolLevel {
					fmt.Printf("%s is not strong enough to break this ore.\n", tool)
				} else if iTool.toolType != iItem.toolType {
					fmt.Println("You need a different kind of tool to break this ore.")
				} else {
					fmt.Printf("The ore breaks, dropping %s, which you pick up.", item)
					inventory[item] = items[item]
					if !iItem.infinite {
						delete(room.items, item)
					}
				}
			} else {
				fmt.Printf("You can't break %s with %s.\n", item, tool)
			}
		} else if iItem.creature {
			//https://github.com/dan200/ComputerCraft/blob/master/src/main/resources/assets/computercraft/lua/rom/programs/fun/adventure.lua#L1007
		}
	}
}

func randomChoice[T any](t []T) T {
	return t[rand.Intn(len(t))]
}

func main() {
	fmt.Println(color.Ize(color.Red, "This is red"))
}
