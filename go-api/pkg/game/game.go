package game

import (
	// "fmt"
	"errors" 
)

type TerrainMap interface {
	getTerrain(Coord) int8
}

type LocatableEntity interface {
	GetLocation() Coord
	GetID() string
	GetAvatar() int8
}

type Coord struct {
	X int
	Y int
}

func WrapMod (x, mod int) int{
	return (x%mod + mod)%mod;
}

func GetNewPosition(move int, position Coord) (Coord, error) {	
	switch move {
		case 38: // up  
			position.Y = WrapMod((position.Y - 1), WorldHeight); 		
		case 40: // down
			position.Y = WrapMod((position.Y + 1), WorldHeight); 		
		case 37: // left
			position.X = WrapMod((position.X - 1), WorldWidth);		
		case 39: // right
			position.X = WrapMod((position.X + 1), WorldWidth); 
		case 13: // enter == ping == no move, just return current pos				
			return position, nil
		default: 		
			return Coord{X: -1, Y: -1}, errors.New("Requested a non-move")
	}
	return position, nil;
}	

func IsLocationValid(location Coord, terrainMap TerrainMap, others []LocatableEntity ) bool {
	
	// Check if the tile is impassible	
	terrain := terrainMap.getTerrain(location)
	if (terrain) <= 3 { // water
		return false
	}
	if (terrain) == 8 { // high mountain
		return false
	}

	// Check if a player is in the way:
	for  _, other := range others { 	
		otherLocation := other.GetLocation()
		if (otherLocation.X == location.X) && (otherLocation.Y == location.Y) {								
			return false
		}								
	}		
	return true;
}