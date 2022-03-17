package data

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
		Name: "Deep Water",
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"medium_water": {
		ID: 2,
		Name: "Medium Water",	
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"shallow_water": {
		ID: 3,
		Name: "Shallow Water",	
		Solid: true,
		Opaque: false,
		Speed: 0,
	},
	"land": {
		ID: 5,
		Name: "Land (grass)",	
		Solid: false,
		Opaque: false,
		Speed: 1,
	},
	"foothills": {
		ID: 6,
		Name: "Foothills",	
		Solid: false,
		Opaque: false,
		Speed: 0.66,
	},
	"low_mountain": {
		ID: 7,
		Name: "Low Mountain",	
		Solid: false,
		Opaque: false,
		Speed: 0.33,
	},
	"high_mountain": {
		ID: 8,
		Name: "High Mountain",	
		Solid: true,
		Opaque: true,
	},	
	"light_forest": {
		ID: 9,
		Name: "Light Forest",	
		Solid: false,
		Opaque: false,
		Speed: 0.7,
	},	
	"heavy_forest": {
		ID: 10,
		Name: "Heavy Forest",	
		Solid: false,
		Opaque: true,
		Speed: 0.4,
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
		Name: "Cactus (desert)",	
		Solid: false,
		Opaque: false,
		Speed: 0.6,
	},
	"desert_sand": {
		ID: 13,
		Name: "Desert",	
		Solid: false,
		Opaque: false,
		Speed: 0.75,
	},	
	"marsh": {
		ID: 14,
		Name: "Marsh/Wetland",	
		Solid: false,
		Opaque: false,
		Speed: 0.4,
	},
}