package sprite

import "math"

var straightCost = subNodeWidth
var obliqueCost = 141
var straightStep = 10
var obliqueStep = 7

type Point struct {
	X int
	Y int
}

func NewPoint(x, y int) *Point {
	return &Point{
		X: x,
		Y: y,
	}
}

type Node struct {
	/*X int // X-coordinate of the Node
	Y int // Y-coordinate of the Node*/
	P *Point
	//g      int     // cost to move from Start Node to this Node
	h      []int // heuristic cost to move from this Node to specific building
	cc     int   // const cost to move across this Node (for wall)
	sc     int   // speed cost when across this Node (for speed cut field)
	dIndex int   // Defender Index in this node
	//f      int     // total cost of this Node: f = g + h
	//parent *Node   // parent Node of this Node
}

func (n *Node) GetH(index int) int {
	if index >= len(n.h) {
		return -1
	}
	return n.h[index]
}

func (n *Node) minH() int {
	minH := n.h[0]
	//minI := 0
	for _, h := range n.h {
		if h < minH {
			minH = h
		}
	}
	return minH
}

func (n *Node) minElfH(elf *Attacker, buildings []*Defender) int {
	minH := math.MaxInt
	inFavor := false
	for i, h := range n.h {
		if buildings[i].LivingStatus == destroyed {
			continue
		}
		favor := buildings[i].isPriority(elf) && buildings[i].LivingStatus == deployed
		if !inFavor {
			if h < minH {
				minH = h
			}
			inFavor = favor
		}
		if inFavor {
			if h < minH && favor {
				minH = h
			}
		}
	}
	return minH
}

var Moves = [][]int{
	{-1, 0},  // up
	{0, -1},  // left
	{1, 0},   // down
	{0, 1},   // right
	{-1, -1}, // up-left
	{-1, 1},  // up-right
	{1, -1},  // down-left
	{1, 1},   // down-right
}

type Direction int

const (
	in Direction = iota - 1
	up
	left
	down
	right
	upleft
	upright
	downleft
	downright
)

// SubNode sub position of a Node, used for distance calculation and movements
type SubNode struct {
	P *Point
}

func (s *SubNode) getD(dx int) int {
	if dx == 0 {
		return 10
	}
	if dx == 1 || dx == 2 || dx == 3 || dx == 4 {
		return 9
	}
	if dx == 5 || dx == 6 {
		return 8
	}
	if dx == 7 {
		return 7
	}
	if dx == 8 {
		return 6
	}
	if dx == 9 {
		return 4
	}
	return 0
}

func (s *SubNode) getStep(directionInt int) (int, int) {
	direction := Direction(directionInt)
	tx, ty := s.getTarget(directionInt)

	var dx, dy int

	if direction >= upleft {
		if SubDistance(s.P.X, s.P.Y, tx, ty) < straightStep {
			return (tx-s.P.X)*Moves[direction][0] + 1, (ty-s.P.Y)*Moves[direction][1] + 1
		} else {
			dx = (tx - s.P.X) * Moves[direction][0]
			dy = (ty - s.P.Y) * Moves[direction][1]
			if dx >= obliqueStep && dy >= obliqueStep {
				return obliqueStep, obliqueStep
			}
			if dx > straightStep {
				dx = straightStep
			}
			return dx, s.getD(dx)
		}
	} else {
		if direction == up || direction == down {
			return straightStep, 0
		} else {
			return 0, straightStep
		}
	}
}

func (s *SubNode) getTarget(directionInt int) (int, int) {
	direction := Direction(directionInt)
	if direction == in {
		return s.P.X, s.P.Y
	}
	if direction == up {
		return -straightStep, s.P.Y
	}
	if direction == left {
		return s.P.X, -straightStep
	}
	if direction == down {
		return straightCost + straightStep - 1, s.P.Y
	}
	if direction == right {
		return s.P.X, straightCost + straightStep - 1
	}
	if direction == upleft {
		return 0, 0
	}
	if direction == upright {
		return 0, straightCost - 1
	}
	if direction == downleft {
		return straightCost - 1, 0
	}
	if direction == downright {
		return straightCost - 1, straightCost - 1
	}
	return s.P.X, s.P.Y
}

type ElfNode struct {
	*Node
	g      int // cost to move from Start Node to this Node
	f      int // total cost of this Node: f = g + h
	parent *ElfNode
}

func newElfNode(x, y int) *ElfNode {
	return &ElfNode{
		Node: &Node{
			P: NewPoint(x, y),
		},
	}
}

func newElfMp(nodes [][]*Node) [][]*ElfNode {
	elfMp := make([][]*ElfNode, len(nodes))
	for raw := range nodes {
		elfMp[raw] = make([]*ElfNode, len(nodes))
		for col := range nodes[raw] {
			elfMp[raw][col] = &ElfNode{
				Node: nodes[raw][col],
			}
		}
	}
	return elfMp
}

// Define a function to implement A* algorithm

func isSearchEndBomber(currentNode *Node, attacker *Attacker, defender *Defender) bool {

	return defender.LivingStatus != destroyed && calWiderDistanceToDefender(currentNode.P.X, currentNode.P.Y, defender) < attacker.atkRange /*defender.isIn(currentNode.X, currentNode.Y)*/
}

func isSearchEnd(currentNode *Node, attacker *Attacker, defender *Defender, needPriority bool) bool {
	if needPriority && !defender.isPriority(attacker) {
		return false
	}
	//end neighbour should be a building
	if defender.DefenderType <= wall && attacker.atkPriority != wall {
		return false
	}

	return defender.LivingStatus != destroyed && calWiderDistanceToDefender(currentNode.P.X, currentNode.P.Y, defender) < attacker.atkRange /*defender.isIn(currentNode.X, currentNode.Y)*/
}

// AstarBomber calculate paths to all Defenders, and chose the closet one that should over path a wall, or find the closet wall to destroy
func AstarBomber(elf *Attacker, matrix [][]int, mp [][]*ElfNode, buildings []*Defender) ([]*Node, int) {
	//hasPriority := HasPriority(buildings, elf)
	// Initialize the open and closed lists
	openList := make([]*ElfNode, 0)
	closedList := make([]*ElfNode, 0)
	skipBuildings := make(map[int]bool, 0)

	// Add the Start Node to the open list
	openList = append(openList, elf.Pos)

	var preTarget int
	var preTargetPath []*Node

	// Loop until the open list is empty
	for len(openList) > 0 {
		// Get the Node with the lowest f value from the open list
		currentNode := openList[0]
		currentIndex := 0
		for i, Node := range openList {
			if Node.f < currentNode.f {
				currentNode = Node
				currentIndex = i
			}
		}

		// Remove the current Node from the open list and add it to the closed list
		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedList = append(closedList, currentNode)

		// If the current Node is the end Node, we have found the path
		for b, _ := range currentNode.Node.h {
			target := b
			if isSearchEndBomber(currentNode.Node, elf, buildings[b]) {
				targetIsWall := buildings[b].DefenderType == wall
				//if buildings[b].LivingStatus != destroyed && d <= obliqueCost && (!hasPriority || buildings[b].isPriority(elf)) {
				path := make([]*Node, 0)
				current := currentNode
				wallInPath := false
				for current != nil {
					path = append(path, current.Node)
					if current.cc != 0 {
						wallInPath = true
						//change Target
						target = current.dIndex
						path = make([]*Node, 0)
					}
					current = current.parent
					//cut the path if a node's const cost not 0 (like wall exist)
				}
				if !wallInPath {
					skipBuildings[b] = true
				}
				// Reverse the path and return it
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				if wallInPath && !targetIsWall && !skipBuildings[b] {
					return path, target
				}
				if preTargetPath == nil || (targetIsWall && buildings[preTarget].DefenderType != wall) {
					preTarget = target
					preTargetPath = path
				}
			}
		}

		// Generate the neighbors of the current Node
		for i, move := range Moves {
			// Get the coordinates of the neighbor
			neighborX := currentNode.P.X + move[0]
			neighborY := currentNode.P.Y + move[1]
			cost := straightCost
			if Direction(i) > upleft {
				cost = obliqueCost
			}
			if !IsInMap(neighborX, neighborY, len(matrix)) {
				continue
			}
			neighborNode := mp[neighborX][neighborY]

			// Check if the neighbor is within the bounds of the// Check if the neighbor is within the bounds of the matrix and is not an obstacle
			if !IsInMap(neighborX, neighborY, len(matrix)) || matrix[neighborX][neighborY] > int(wall) {
				continue
			}
			cost += mp[neighborX][neighborY].cc

			// Create a new Node for the neighbor
			if neighborNode.parent == nil && neighborNode != elf.Pos {
				neighborNode.parent = currentNode
				// Calculate the g, h, and f values for the neighbor Node
				neighborNode.g = currentNode.g + cost
				//neighborNode.h[0] = int(math.Sqrt(math.Pow(float64(straightCost)*math.Abs(float64(neighborX-endNode.X)), 2) + math.Pow(float64(straightCost)*math.Abs(float64(neighborY-endNode.Y)), 2)))
				neighborNode.f = neighborNode.g + neighborNode.minElfH(elf, buildings)
			}

			// Check if the neighbor is already in the closed list
			inClosedList := false
			for _, Node := range closedList {
				if Node.P.X == neighborNode.P.X && Node.P.Y == neighborNode.P.Y {
					inClosedList = true
					break
				}
			}

			// If the neighbor is already in the closed list, skip it
			if inClosedList {
				continue
			}

			// Check if the neighbor is already in the open list
			inOpenList := false
			for _, Node := range openList {
				if Node.P.X == neighborNode.P.X && Node.P.Y == neighborNode.P.Y {
					inOpenList = true
					break
				}
			}

			// If the neighbor is not in the open list, add it
			if !inOpenList {
				openList = append(openList, neighborNode)
			} else {
				// If the neighbor is already in the open list, check if this path to the neighbor is better
				for i, node := range openList {
					if node.P.X == neighborNode.P.X && node.P.Y == neighborNode.P.Y {
						cost := straightCost
						if node.P.X != currentNode.P.X && node.P.Y != currentNode.P.Y {
							cost = obliqueCost
						}
						// Add const cost if it is a wall
						cost += node.cc
						if currentNode.g+cost < node.g {
							openList[i].g = currentNode.g + cost
							openList[i].f = openList[i].g + openList[i].minElfH(elf, buildings)
							openList[i].parent = currentNode
						}
						break
					}
				}
			}
		}
	}
	// If we reach this P, there is no path from the Start Node to the end Node
	return preTargetPath, preTarget
}

func AstarBuilding(elf *Attacker, matrix [][]int, mp [][]*ElfNode, buildings []*Defender) ([]*Node, int) {
	hasPriority := HasPriority(buildings, elf)
	// Initialize the open and closed lists
	openList := make([]*ElfNode, 0)
	closedList := make([]*ElfNode, 0)

	// Add the Start Node to the open list
	openList = append(openList, elf.Pos)

	// Loop until the open list is empty
	for len(openList) > 0 {
		// Get the Node with the lowest f value from the open list
		currentNode := openList[0]
		currentIndex := 0
		for i, Node := range openList {
			if Node.f < currentNode.f {
				currentNode = Node
				currentIndex = i
			}
		}

		// Remove the current Node from the open list and add it to the closed list
		openList = append(openList[:currentIndex], openList[currentIndex+1:]...)
		closedList = append(closedList, currentNode)

		// If the current Node is the end Node, we have found the path
		for b, _ := range currentNode.Node.h {
			target := b
			if isSearchEnd(currentNode.Node, elf, buildings[b], hasPriority) {
				//if buildings[b].LivingStatus != destroyed && d <= obliqueCost && (!hasPriority || buildings[b].isPriority(elf)) {
				path := make([]*Node, 0)
				current := currentNode
				for current != nil {
					path = append(path, current.Node)
					if current.cc != 0 {
						//change Target
						target = current.dIndex
						path = make([]*Node, 0)
					}
					current = current.parent
					//cut the path if a node's const cost not 0 (like wall exist)
				}
				// Reverse the path and return it
				for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
					path[i], path[j] = path[j], path[i]
				}
				return path, target
			}
		}

		// Generate the neighbors of the current Node
		for i, move := range Moves {
			// Get the coordinates of the neighbor
			neighborX := currentNode.P.X + move[0]
			neighborY := currentNode.P.Y + move[1]
			cost := straightCost
			if Direction(i) > upleft {
				cost = obliqueCost
			}
			if !IsInMap(neighborX, neighborY, len(matrix)) {
				continue
			}
			neighborNode := mp[neighborX][neighborY]

			// Check if the neighbor is within the bounds of the// Check if the neighbor is within the bounds of the matrix and is not an obstacle
			if !IsInMap(neighborX, neighborY, len(matrix)) || matrix[neighborX][neighborY] > int(wall) {
				continue
			}
			cost += mp[neighborX][neighborY].cc

			// Create a new Node for the neighbor
			if neighborNode.parent == nil && neighborNode != elf.Pos {
				neighborNode.parent = currentNode
				// Calculate the g, h, and f values for the neighbor Node
				neighborNode.g = currentNode.g + cost
				//neighborNode.h[0] = int(math.Sqrt(math.Pow(float64(straightCost)*math.Abs(float64(neighborX-endNode.X)), 2) + math.Pow(float64(straightCost)*math.Abs(float64(neighborY-endNode.Y)), 2)))
				neighborNode.f = neighborNode.g + neighborNode.minElfH(elf, buildings)
			}

			// Check if the neighbor is already in the closed list
			inClosedList := false
			for _, Node := range closedList {
				if Node.P.X == neighborNode.P.X && Node.P.Y == neighborNode.P.Y {
					inClosedList = true
					break
				}
			}

			// If the neighbor is already in the closed list, skip it
			if inClosedList {
				continue
			}

			// Check if the neighbor is already in the open list
			inOpenList := false
			for _, Node := range openList {
				if Node.P.X == neighborNode.P.X && Node.P.Y == neighborNode.P.Y {
					inOpenList = true
					break
				}
			}

			// If the neighbor is not in the open list, add it
			if !inOpenList {
				openList = append(openList, neighborNode)
			} else {
				// If the neighbor is already in the open list, check if this path to the neighbor is better
				for i, node := range openList {
					if node.P.X == neighborNode.P.X && node.P.Y == neighborNode.P.Y {
						cost := straightCost
						if node.P.X != currentNode.P.X && node.P.Y != currentNode.P.Y {
							cost = obliqueCost
						}
						// Add const cost if it is a wall
						cost += node.cc
						if currentNode.g+cost < node.g {
							openList[i].g = currentNode.g + cost
							openList[i].f = openList[i].g + openList[i].minElfH(elf, buildings)
							openList[i].parent = currentNode
						}
						break
					}
				}
			}
		}
	}
	// If we reach this P, there is no path from the Start Node to the end Node
	return nil, -1
}
