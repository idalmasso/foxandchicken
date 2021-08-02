package game

import "github.com/idalmasso/foxandchicken/server/game/common"

func (g *GameRoom) addGameObject(gameObject *GameObject) {
	cell := g.getCellNum(gameObject.Position.X, gameObject.Position.Y)
	g.cellsToGameObjectmap[cell][gameObject] = struct{}{}
}
func (g *GameRoom) gameObjectsInPoint(point common.Vector2) []*GameObject {
	objects := make([]*GameObject, 0)
	cell := g.getCellNum(point.X, point.Y)
	objects = g.pointIntersectObjectsInCell(point, objects, cell)
	cells := g.getCellNeightbours(cell)
	for _, cell = range cells {
		objects = g.pointIntersectObjectsInCell(point, objects, cell)
	}
	return objects
}

func (g *GameRoom) pointIntersectObjectsInCell(point common.Vector2, objects []*GameObject, cell int) []*GameObject {
	for p := range g.cellsToGameObjectmap[cell] {
		v := common.Vector2{X: point.X - p.Position.X, Y: point.Y - p.Position.Y}
		if v.Magnitude() < p.size*p.size {
			objects = append(objects, p)
		}
	}
	return objects
}

func (g *GameRoom) objectMove(gameObject *GameObject, from, to common.Vector2) {
	fromCell := g.getCellNum(from.X, from.Y)
	toCell := g.getCellNum(to.X, to.Y)
	if fromCell == toCell {
		return
	}
	delete(g.cellsToGameObjectmap[fromCell], gameObject)
	g.cellsToGameObjectmap[toCell][gameObject] = struct{}{}
}

//Get the cellnum in 1-based grid
func (g *GameRoom) getCellNum(x, y float64) int {
	if x < 0 || x > g.sizeX || y < 0 || y > g.sizeY {
		return -1
	}

	intX, intY := int(x), int(y)
	if intX == int(g.sizeX) {
		intX--
	}
	if intY == int(g.sizeY) {
		intY--
	}
	return intX + intY*int(g.sizeX)
}

//get the 8 behaviours cells
func (g *GameRoom) getCellNeightbours(cellNum int) []int {
	if cellNum < 0 || cellNum >= int(g.sizeX*g.sizeY) {
		return make([]int, 0)
	}
	size := 0
	bottom, top, left, right := true, true, true, true

	//First bottom line
	if cellNum < int(g.sizeX) {
		bottom = false
	}
	//Last top line
	if cellNum > (int(g.sizeX)-1)*int(g.sizeY) {
		top = false
	}
	//first column
	if cellNum%int(g.sizeX) == 0 {
		left = false
	}
	//last column
	if (cellNum+1)%int(g.sizeX) == 0 {
		right = false
	}
	size = 8
	if !top || !bottom {
		size -= 3
		if !left || !right {
			size -= 2
		}
	} else {
		if !left || !right {
			size -= 3
		}
	}
	neightbours := make([]int, size)
	index := 0
	if top {
		if left {
			neightbours[index] = cellNum + int(g.sizeX) - 1
			index++
		}
		neightbours[index] = cellNum + int(g.sizeX)
		index++
		if right {
			neightbours[index] = cellNum + int(g.sizeX) + 1
			index++
		}
	}
	if left {
		neightbours[index] = cellNum - 1
		index++
	}
	if right {
		neightbours[index] = cellNum + 1
		index++
	}
	if bottom {
		if left {
			neightbours[index] = cellNum - int(g.sizeX) - 1
			index++
		}
		neightbours[index] = cellNum - int(g.sizeX)
		index++
		if right {
			neightbours[index] = cellNum - int(g.sizeX) + 1
			index++
		}
	}
	return neightbours
}
