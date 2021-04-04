package game

type gridMap struct {
	grid [][]int8	
}
func (g gridMap) InBounds(x, y int) bool {
	if (x < 0 || y < 0 || x > ViewWidth-1 || y > ViewHeight-1) {
		return false
	}
	return true;
}
func (g gridMap) IsOpaque(x, y int) bool {	
	if g.grid[x][y] == 8 { // high mountain
		return true
	}
	if g.grid[x][y] == 10 { // heavy forest
		return true
	}
	return false
}