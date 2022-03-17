package data
import (
	"app/pkg/game/model"
)

type GridMap struct {
	Grid [][]*model.Tile
}
func (g GridMap) InBounds(x, y int) bool {
	if (x < 0 || y < 0 || x > 14 || y > 14) { // @todo: make view size dynamic. viewWidth-1, viewHeight-1
		return false
	}
	return true;
}
func (g GridMap) IsOpaque(x, y int) bool {	

	// return get_tile_opacity(g.grid[x][y]) < 0.5 
	return g.Grid[x][y].Opaque
	
}