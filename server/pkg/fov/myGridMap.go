package fov
import (
	"app/pkg/game/model"
)

type MyGridMap struct {
	Grid [][]*model.Tile 
}
func (g MyGridMap) InBounds(x, y int) bool {
	if (x < 0 || y < 0 || x > 16 || y > 16) { // @todo: make view size dynamic. viewWidth-1, viewHeight-1
		return false
	}
	return true;
}
func (g MyGridMap) IsOpaque(x, y int) bool {	
	// return get_tile_opacity(g.grid[x][y]) < 0.5 
	return g.Grid[x][y].Opaque	
} 