package sprite

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"testing"
	"time"
)

/*func TestAstar_solidity_matrix(t *testing.T) {

	// Find the path using A* algorithm
	var path []*Node
	buildings, _ := generateSprites()
	battleMap, elfMp := buildSolidityMatrix(buildings)
	Start := time.Now()

	//for i := 0; i < 100; i++ {
	/*matrix := GenerateMap(len(matrix), len(matrix[0]))
	elfMp := GenerateElfMap(matrix)
	startNode := battleMap.nodes[0][4]
	//endNode := matrix[39][39]
	startElfNode := &ElfNode{
		Node: startNode,
	}
	path = AstarTargetBuilding(startElfNode, 2, battleMap.matrix, elfMp)
	//}
	end := time.Now()
	fmt.Printf("%v\n", end.Sub(Start))
	Print(battleMap.matrix, path)
}*/

func TestAstar_priority_matrix(t *testing.T) {

	// Find the path using A* algorithm
	var path []*Node
	buildings, _ := generateSprites()
	battleMap := buildSolidityMatrix(buildings)
	start := time.Now()

	elf := &Attacker{
		Pos:          newElfNode(4, 4),
		atkRange:     0,
		atkPriority:  3,
		speed:        0,
		moveType:     0,
		LivingStatus: 0,
		AtkStatus:    0,
		heading:      nil,
		path:         nil,
	}
	path, _ = AstarBuilding(elf, battleMap.matrix, newElfMp(battleMap.nodes), battleMap.Defenders)
	//}
	end := time.Now()
	fmt.Printf("%v\n", end.Sub(start))
	Print(battleMap.matrix, path)
}

func TestAstar_battle(t *testing.T) {
	buildings, attackers := generateSprites()
	battle := buildSolidityMatrix(buildings)
	for _, atk := range attackers {
		battle.AddAttackerOri(atk)
		//atk.Start(battle.ctx)
	}
	/*for _, building := range buildings {
		building.Start(battle.ctx)
	}*/
	start := time.Now()
	battle.Start(context.Background())
	end := time.Now()
	displayLog(battle)
	fmt.Printf("cost:%d", end.Sub(start).Milliseconds())
	//viewer()
}

func TestDisplayBattle(t *testing.T) {
	file, _ := os.Open("/Users/wangyuan/Downloads/projects/yongcheng/shaman/game/rpc/battle/logres/shaman_log_file_2023-05-19 15:02:21.727710.txt")
	battle := &Battle{}
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		if !strings.HasPrefix(line, "{\"Defenders\"") {
			continue
		}
		json.Unmarshal([]byte(line), battle)
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

func displayLog(battle *Battle) {
	for _, log := range battle.log {
		logBattle := &Battle{}
		err := json.Unmarshal([]byte(log), logBattle)
		if err != nil {
			fmt.Printf("json unmarshal failed:%v", err)
			return
		}
		//replace living status
		for i, defender := range logBattle.Defenders {
			defender.XMin = battle.Defenders[i].XMin
			defender.YMin = battle.Defenders[i].YMin
			defender.XMax = battle.Defenders[i].XMax
			defender.YMax = battle.Defenders[i].YMax
			defender.DefenderType = battle.Defenders[i].DefenderType
		}
		PrintBattle(len(GetBattle(0).matrix), logBattle.Attackers, logBattle.Defenders)
	}
}

func printLog(battle *Battle) {
	for _, s := range battle.log {
		fmt.Printf("%s\n", s)
	}
}

func viewer() {
	size := len(GetBattle(0).matrix)
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		for {
			select {
			case <-ticker.C:
				PrintBattle(size, GetBattle(0).Attackers, GetBattle(0).Defenders)
			}
		}
	}()
}

func PrintBattle(size int, attackers []*Attacker, defenders []*Defender) {
	matrix := make([][]string, size)
	for i := range matrix {
		matrix[i] = make([]string, size)
		for j := range matrix[i] {
			matrix[i][j] = "  "
		}
	}

	for _, defender := range defenders {
		str := fmt.Sprintf("%d ", defender.DefenderType)
		if defender.LivingStatus == destroyed {
			str = "* "
		}
		for i := defender.XMin; i <= defender.XMax; i++ {
			for j := defender.YMin; j <= defender.YMax; j++ {
				matrix[i][j] = str
			}
		}
	}
	for _, attacker := range attackers {
		atkPos := attacker.Pos
		str := "0 "
		if attacker.LivingStatus == destroyed {
			str = "* "
		}
		if atkPos != nil {
			matrix[atkPos.P.X][atkPos.P.Y] = str
		}
	}

	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			fmt.Print(matrix[i][j])
		}
		fmt.Print("\n")
	}
	fmt.Printf("pos:%d, %d subPos:%d, %d",
		attackers[0].Pos.P.X, attackers[0].Pos.P.Y, attackers[0].SubPos.P.X, attackers[0].SubPos.P.Y)

}

func Print(matrix [][]int, path []*Node) {
	for i := 0; i < len(matrix); i++ {
		for j := 0; j < len(matrix[0]); j++ {
			if checkIn(i, j, path) {
				fmt.Printf("%d ", 0)
			} else {
				if matrix[i][j] != 0 {
					fmt.Printf("%d ", matrix[i][j])
				} else {
					fmt.Printf("  ")
				}

			}
		}
		fmt.Print("\n")
	}
	// Print the path

	for _, Node := range path {
		fmt.Printf("(%d,%d) ", Node.P.X, Node.P.Y)
	}
}

func TestAngle(t *testing.T) {
	subNode := SubNode{
		P: NewPoint(98, 1),
	}
	direction := downleft
	dx, dy := subNode.getStep(int(direction))
	fmt.Printf("%d, %d", dx, dy)
}

func TestDistance(t *testing.T) {
	battleId := InitBattle(20, 3)
	battle := GetBattle(battleId)
	var defenders []*Defender
	defenders = append(defenders, NewDefender(6, 5, 4, 0, 5000, 10, 400, 4, nil))
	defenders = append(defenders, NewDefender(11, 5, 4, 0, 5000, 10, 400, 3, nil))
	defenders = append(defenders, NewDefender(9, 12, 4, 0, 5000, 10, 400, 3, nil))
	e := ReflectCall("battle", "AddDefender", 0, 0, "6, 5,4,0,5000,400,4")
	e = ReflectCall("battle", "AddDefender", 0, 0, "11,5,4,0,5000,400,3")
	e = ReflectCall("battle", "AddDefender", 0, 0, "9,12,4,0,5000,400,3")

	//k := battle.GetMap(14, 7)
	k := ReflectCall("battle", "GetMap", 0, 0, "14,7")
	fmt.Printf("%v %v, %v", e, k, battle)
}

func buildSolidityMatrix(defenders []*Defender) *Battle {
	//mp, nodes, c := buildMatrix(Defenders, 20)
	battleIndex := InitBattle(20, len(defenders))
	battle := GetBattle(battleIndex)
	for _, defender := range defenders {
		battle.AddDefenderOri(defender)
	}
	return battle
}

func buildRandomMatrix() [][]int {
	matrix := [][]int{
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 1, 1, 1, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0},
		{0, 1, 1, 1, 1, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 1, 1, 1, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 1, 1, 1, 1, 0, 0, 1, 1, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 1, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
		{0, 0, 1, 1, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0},
	}
	return matrix
}

func generateSprites() ([]*Defender, []*Attacker) {
	type ap struct {
		num, x, y, atkRange, damage, deathDamage, priority, atkType int
		suicide                                                     bool
	}
	type wp struct {
		xMin, xMax, yMin, yMax int
	}
	type p struct {
		functional, size, posX, posY, damage, atkRange, health int
	}
	/*wArr := []wp{
		{6, 12, 9, 15},
		{8, 10, 11, 13},
		{17, 17, 5, 5},
	}*/

	//case which attacker could go inside the remains surrounded by the wall and attack another Defenders
	/*attackArr := []ap{
		{5, 8, 14, 10, 3},
	}
	defenderArr := []P{
		{4, 4, 6, 5, 0, 5, 40},
		{3, 4, 11, 5, 0, 5, 40},
		{3, 5, 7, 10, 0, 5, 40},
	}*/
	//case which attacker balance the wall cost and choose its path
	/*attackArr := []ap{
		{0, 19, 14, 10, 3},
	}
	defenderArr := []P{
		{4, 4, 6, 5, 0, 5, 40},
		{3, 4, 11, 5, 0, 5, 40},
		{3, 5, 7, 10, 0, 5, 40},
	}*/
	attackArr := []ap{
		{10, 1983, 2, 2, 10, 0, 0, 0, false},
		//{1, 17, 6, 60, 30, 0, 0, 0, false},
		//{1, 19, 19, 60, 30, 0, 0, 0, false},
		//{500, 19, 0, 60, 1, 0, 3, 0, false},
	}
	/*defenderArr := []P{
		{4, 4, 6, 5, 30, 500, 400},
		{3, 4, 11, 5, 0, 500, 400},
		{3, 1, 9, 12, 0, 500, 400},
	}*/
	defenderArr := []p{
		{4, 4, 6, 5, 0, 5000, 400},
		//{3, 4, 11, 5, 0, 5000, 400},
		{3, 4, 10, 12, 0, 5000, 400},
	}

	var defenders []*Defender
	var attackers []*Attacker
	for _, a := range attackArr {
		for i := 0; i < a.num; i++ {
			attackers = append(attackers, NewAttacker(a.x, a.y, a.atkRange, a.damage, a.deathDamage, a.atkType, a.priority, 10, 100, 100, a.suicide))
		}
	}
	//var Attackers []*Attacker
	for _, a := range defenderArr {
		defenders = append(defenders, NewDefender(a.posX, a.posY, a.size, a.damage, a.atkRange, 10, a.health, a.functional, nil))
	}
	/*for _, item := range wArr {
		for j := item.YMin; j <= item.YMax; j++ {
			defenders = append(defenders, NewDefender(item.XMin, j, 1, 0, 0, 400, int(wall)))
			defenders = append(defenders, NewDefender(item.XMax, j, 1, 0, 0, 400, int(wall)))
		}
		for i := item.XMin + 1; i < item.XMax; i++ {
			defenders = append(defenders, NewDefender(i, item.YMin, 1, 0, 0, 400, int(wall)))
			defenders = append(defenders, NewDefender(i, item.YMax, 1, 0, 0, 400, int(wall)))
		}
	}*/
	/*for _, b := range att {
		Attackers = append(Attackers, NewAttacker(b.posX, b.posY, DefenderType(b.priority)))
	}*/
	return defenders, attackers
}

func checkIn(x, y int, path []*Node) bool {
	for _, item := range path {
		if x == item.P.X && y == item.P.Y {
			return true
		}
	}
	return false
}

func TestSubNode_GetStep(t *testing.T) {
	s := &SubNode{
		P: NewPoint(1, 9),
	}
	direction := upleft
	dx, dy := s.getStep(int(direction))
	fmt.Printf("%d, %d", dx*Moves[direction][0], dy*Moves[direction][1])
}
