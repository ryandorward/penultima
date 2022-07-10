package tiles

import (
	"app/pkg/game/model"
)

var Tiles = map[string]model.Tile {
	"nothing": {
		ID: 0,
		Name: "Nothing",
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"deep_water": {
		ID: 1,
		Name: "Very deep water",
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"medium_water": {
		ID: 2,
		Name: "Deep water",	
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"shallow_water": {
		ID: 3,
		Name: "Shallow water",	
		Solid: true,
		Opaque: false,
		Speed: 0.0005,
	},
	"land": {
		ID: 5,
		Name: "Grass",	
		Solid: false,
		Opaque: false,
		Speed: 0.88,
	},
	"foothills": {
		ID: 6,
		Name: "Foothills",	
		Solid: false,
		Opaque: false,
		Speed: 0.5,
	},
	"low_mountain": {
		ID: 7,
		Name: "Low mountains",	
		Solid: false,
		Opaque: false,
		Speed: 0.22,
	},
	"high_mountain": {
		ID: 8,
		Name: "High mountains",	
		Solid: true,
		Opaque: true,
		Speed: 0,
	},		
	"light_forest": {
		ID: 9,
		Name: "Light Forest",	
		Solid: false,
		Opaque: false,
		Speed: 0.6,
	},	
	"heavy_forest": {
		ID: 10,
		Name: "Thick forest",	
		Solid: false,
		Opaque: true,
		Speed: 0.3,
	},
	"beach": {
		ID: 11,
		Name: "Beach",	
		Solid: false,
		Opaque: false,
		Speed: 0.75,
	},
	"desert_cactus": {
		ID: 12,
		Name: "Desert cactus",	
		Solid: false,
		Opaque: false,
		Speed: 0.6,
	},
	"desert_sand": {
		ID: 13,
		Name: "Desert",	
		Solid: false,
		Opaque: false,
		Speed: 0.6,
	},	
	"marsh": {
		ID: 14,
		Name: "Wetlands",	
		Solid: false,
		Opaque: false,
		Speed: 0.3,
	},
	"brush": {
		ID: 15,
		Name: "Bushes",
		Solid: false,
		Opaque: false,
		Speed: 0.8,
	},
	"brush2": {
		ID: 16,
		Name: "Bushes",
		Solid: false,
		Opaque: false,
		Speed: 0.82,
	},
	"brush3": {
		ID: 17,
		Name: "Bushes",
		Solid: false,
		Opaque: false,
		Speed: 0.78,
	},
	"brush4": {
		ID: 18,
		Name: "Bushes",
		Solid: false,
		Opaque: false,
		Speed: 0.85,
	},
	"land_shallow_waterS": {
		ID: 19,
		Name: "Shoreline",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},

	"hut": {
		ID: 200,
		Name: "Hut",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"village": {
		ID: 201,
		Name: "Village",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},

	"town": {
		ID: 202,
		Name: "Towne",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"castle": {
		ID: 203,
		Name: "Castle",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"tower": {
		ID: 204,
		Name: "Tower",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"temple": {
		ID: 205,
		Name: "Temple",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"shrine": {
		ID: 206,
		Name: "Shrine",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"broken-shrine": {
		ID: 207,
		Name: "Broken Shrine",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"lighthouse": {
		ID: 208,
		Name: "Lighthouse",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},	
	"cave": {
		ID: 250,
		Name: "Cave",	
		Solid: false,
		Opaque: false,
		Speed: 0.22,
	},
	"blocked-cave": {
		ID: 251,
		Name: "Blocked cave",
		Solid: false,
		Opaque: false,
		Speed: 0.22,
	},
	"mine": {
		ID: 252,
		Name: "Mine",	
		Solid: false,
		Opaque: false,
		Speed: 0.22,
	},
	"dungeon": {
		ID: 253,
		Name: "Dungeon",	
		Solid: false,
		Opaque: false,
		Speed: 0.22,
	},

	"chicken": {
		ID: 254,
		Name: "Chicken",	
		Solid: false,
		Opaque: false,
		Speed: 0.8,
	},

	// from 1000+, indexing tiles from their position in the spritesheet
	// 1000 is the first tile in the spritesheet
	// spritesheet is 32 tiles wide, 16 tall - 512 tiles total
	// so reserving 1000-1512 for this indexing

	"pathNS": {
		ID: 1032,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathEW": {
		ID: 1033,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathNE": {
		ID: 1034,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathSE": {
		ID: 1035,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathSW": {
		ID: 1036,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathNW":{
		ID: 1037,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"pathIntersection": {
		ID: 1038,
		Name: "Path",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},

	"fruitTree": {
		ID: 1046,
		Name: "Fruit tree",
		Solid: false,
		Opaque: false,
		Speed: 0.6,
	},

	"woodFloorHoriz": {
		ID: 1064,
		Name: "Wood floor",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	
	"tileFloor": {
		ID: 1068,
		Name: "Floor",
		Solid: false,
		Opaque: false,
		Speed: 1,
	},

	"stones": {
		ID: 1076,
		Name: "Stones",
		Solid: true,
		Opaque: false,
		Speed: 0,
	},

	"rockWall": {
		ID: 1077, 
		Name: "Rock wall",
		Solid: true,
		Opaque: true,
		Speed: 0,
	},
	"stoneWallSecretDoor": {
		ID: 1078, 
		Name: "Stone wall",
		Solid: false,
		Opaque: true,
		Speed: 1,
	},
	"stoneWall": { 
		ID: 1079,
		Name: "Stone wall",
		Solid: true,
		Opaque: true,
		Speed: 0,
	},
	"moneyBag": {
		ID: 1258,
		Name: "Silver",
		Solid: false,
		Opaque: false,
		Speed: 0, 
	},
	"gem": {
		ID: 1264,
		Name: "Gem",
		Solid: false,
		Opaque: false,
		Speed: 0, 
	},

	"cookedChicken": {
		ID: 1271,
		Name: "Food",
		Solid: false,
		Opaque: false,
		Speed: 0, 
	},
	"deadBody": {
		ID: 1286,
		Name: "Dead Body",
		Solid: false,
		Opaque: false,
		Speed: 0.5,
	},

	"guard": {
		ID: 1369,
		Name: "Guard", 
		Solid: true,
		Opaque: true,
		Speed: 0,
	},

	"minstrel": {
		ID: 1349,
		Name: "Minstrel",
		Solid: true,
		Opaque: true,
		Speed: 0,
	},

	"orc": {
		ID: 1449,
		Name: "Orc",
		Solid: true,
		Opaque: true,
		Speed: 0, 
	},

	"snake": {
		ID: 1457,
		Name: "Snake",
		Solid: true,
		Opaque: true,
		Speed: 0, 
	},

	

};