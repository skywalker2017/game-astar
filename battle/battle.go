package sprite

import (
	"context"
	"encoding/json"
	"fmt"
)

var battleList []*Battle

type Battle struct {
	userId    uint64
	ctx       context.Context
	cancel    context.CancelFunc
	matrix    [][]int
	nodes     [][]*Node
	Defenders []*Defender
	Attackers []*Attacker
	log       []string
	Result    *BattleResult
}

type BattleResult struct {
	DestroyedPTC int
	ResourceMap  map[string]Resource
}

func GetVersion() string {
	return "0.0.1"
}

func InitBattle(size, buildingSize int) int {
	ctx, cancel := context.WithCancel(context.Background())
	mp, nodes, _ := buildMatrix(size, buildingSize)
	index := len(battleList)
	battleList = append(battleList, &Battle{
		matrix: mp,
		nodes:  nodes,
		ctx:    ctx,
		cancel: cancel,
	})
	return index
}

func CreateBattle(size, buildingSize int, userId uint64) *Battle {
	ctx, cancel := context.WithCancel(context.Background())
	mp, nodes, _ := buildMatrix(size, buildingSize)
	return &Battle{
		userId: userId,
		matrix: mp,
		nodes:  nodes,
		ctx:    ctx,
		cancel: cancel,
	}
}

func buildMatrix(size int, buildingSize int) ([][]int, [][]*Node, [][]*ElfNode) {

	matrix := make([][]int, size)
	nodes := make([][]*Node, size)
	elfMp := make([][]*ElfNode, size)
	for i := range matrix {
		matrix[i] = make([]int, size)
		nodes[i] = make([]*Node, size)
		elfMp[i] = make([]*ElfNode, size)
		for j := range matrix[i] {
			nodes[i][j] = &Node{
				P: NewPoint(i, j),
				h: make([]int, buildingSize),
			}
			elfMp[i][j] = &ElfNode{
				Node: nodes[i][j],
			}
		}
	}
	/*for b, obj := range buildings {
		//init Node.cc
		if obj.DefenderType == wall {
			nodes[obj.XMin][obj.YMin].cc = obj.Health
		}
		//init h[]
		for i := 0; i < size; i++ {
			for j := 0; j < size; j++ {
				if i < obj.XMin && j < obj.YMin {
					nodes[i][j].h[b] = Distance(i, j, obj.XMin, obj.YMin)
				}
				if i < obj.XMin && j >= obj.YMin && j <= obj.YMax {
					nodes[i][j].h[b] = straightCost * (obj.XMin - i)
				}
				if i < obj.XMin && j > obj.YMax {
					nodes[i][j].h[b] = Distance(i, j, obj.XMin, obj.YMax)
				}

				if i >= obj.XMin && i <= obj.XMax && j < obj.YMin {
					nodes[i][j].h[b] = straightCost * (obj.YMin - j)
				}
				if i >= obj.XMin && i <= obj.XMax && j >= obj.YMin && j <= obj.YMax {
					matrix[i][j] = int(obj.DefenderType)
					nodes[i][j].dIndex = obj.Index
				}
				if i >= obj.XMin && i <= obj.XMax && j > obj.YMax {
					nodes[i][j].h[b] = straightCost * (j - obj.YMax)
				}

				if i > obj.XMax && j < obj.YMin {
					nodes[i][j].h[b] = Distance(i, j, obj.XMax, obj.YMin)
				}
				if i > obj.XMax && j >= obj.YMin && j <= obj.YMax {
					nodes[i][j].h[b] = straightCost * (i - obj.XMax)
				}
				if i > obj.XMax && j > obj.YMax {
					nodes[i][j].h[b] = Distance(i, j, obj.XMax, obj.YMax)
				}
			}
		}
	}*/
	return matrix, nodes, elfMp
}

func GetBattle(index int) *Battle {
	if len(battleList) <= index {
		return nil
	}
	return battleList[index]
}

func (b *Battle) CalResult() {

}

func (b *Battle) GetStatus() (string, error) {
	buf, err := json.Marshal(b)
	if err != nil {
		return "", err
	}
	return string(buf), nil
}

func (b *Battle) AddDefenderOri(obj *Defender) int {
	obj.battle = b
	size := len(b.nodes)
	lock.Lock()
	index := len(b.Defenders)
	obj.Index = index
	b.Defenders = append(b.Defenders, obj)
	lock.Unlock()
	if obj.DefenderType == wall {
		b.nodes[obj.XMin][obj.YMin].cc = obj.Health
	}
	//init h[]
	for i := 0; i < size; i++ {
		for j := 0; j < size; j++ {
			if i < obj.XMin && j < obj.YMin {
				b.nodes[i][j].h[index] = Distance(i, j, obj.XMin, obj.YMin)
			}
			if i < obj.XMin && j >= obj.YMin && j <= obj.YMax {
				b.nodes[i][j].h[index] = straightCost * (obj.XMin - i)
			}
			if i < obj.XMin && j > obj.YMax {
				b.nodes[i][j].h[index] = Distance(i, j, obj.XMin, obj.YMax)
			}

			if i >= obj.XMin && i <= obj.XMax && j < obj.YMin {
				b.nodes[i][j].h[index] = straightCost * (obj.YMin - j)
			}
			if i >= obj.XMin && i <= obj.XMax && j >= obj.YMin && j <= obj.YMax {
				b.matrix[i][j] = int(obj.DefenderType)
				b.nodes[i][j].dIndex = index
			}
			if i >= obj.XMin && i <= obj.XMax && j > obj.YMax {
				b.nodes[i][j].h[index] = straightCost * (j - obj.YMax)
			}

			if i > obj.XMax && j < obj.YMin {
				b.nodes[i][j].h[index] = Distance(i, j, obj.XMax, obj.YMin)
			}
			if i > obj.XMax && j >= obj.YMin && j <= obj.YMax {
				b.nodes[i][j].h[index] = straightCost * (i - obj.XMax)
			}
			if i > obj.XMax && j > obj.YMax {
				b.nodes[i][j].h[index] = Distance(i, j, obj.XMax, obj.YMax)
			}
		}
	}
	return index
}

// AddDefender add defender to the battle, it should be the only method that concern about explicit resource type
func (b *Battle) AddDefender(xMin, yMin, dSize, damage, atkRange, atkFrequency, health int, defenderType int, water int, money int) int {
	resourceMap := make(map[string]Resource)
	//todo abstract to convertUtil
	resourceMap[ResourceMoney] = Resource{
		Preserve:    money,
		PreserveOri: money,
	}
	resourceMap[ResourceWater] = Resource{
		Preserve:    water,
		PreserveOri: water,
	}
	obj := NewDefender(xMin, yMin, dSize, damage, atkRange, atkFrequency, health, defenderType, resourceMap)
	return b.AddDefenderOri(obj)
}

func (b *Battle) AddAttackerOri(attacker *Attacker) int {
	lock.Lock()
	attacker.Index = len(b.Attackers)
	b.Attackers = append(b.Attackers, attacker)
	lock.Unlock()
	attacker.battle = b
	attacker.path = []*Node{attacker.battle.nodes[attacker.Pos.P.X][attacker.Pos.P.Y]}
	return attacker.Index
}

func (b *Battle) AddAttacker(x, y, atkRange, damage, deathDamage, atkTypeInt, atkPriorityInt, atkFrequency, health, speed, suicideInt int) int {
	suicide := !(suicideInt == 0)
	attacker := NewAttacker(x, y, atkRange, damage, deathDamage, atkTypeInt, atkPriorityInt, atkFrequency, health, speed, suicide)
	return b.AddAttackerOri(attacker)
}

func (b *Battle) GetDefender(index int) *Defender {
	if index < 0 || len(b.Defenders) <= index {
		return nil
	}
	return b.Defenders[index]
}

func (b *Battle) GetAttacker(index int) *Attacker {
	if index < 0 || len(b.Attackers) <= index {
		return nil
	}
	return b.Attackers[index]
}

func (b *Battle) GetNode(x, y int) *Node {
	if x >= len(b.nodes) || y >= len(b.nodes) {
		return nil
	}
	return b.nodes[x][y]
}

func (b *Battle) GetMap(x, y int) int {
	if x >= len(b.nodes) || y >= len(b.nodes) {
		return -1
	}
	return b.matrix[x][y]
}

func (b *Battle) logStore() {
	log, err := b.GetStatus()
	if err != nil {
		fmt.Printf("log failed:%v", err)
	}
	b.log = append(b.log, log)
}

func (b *Battle) Play() bool {

	playing := false
	for _, d := range b.Defenders {
		if d.CanPlay() {
			playing = true
			go d.play()
		}
	}
	if !playing {
		return playing
	}
	for _, a := range b.Attackers {
		go a.play()
	}
	b.logStore()
	return playing
}

func (b *Battle) Start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	b.ctx = ctx
	b.cancel = cancel
	for {
		if !b.Play() {
			return
		}
	}
	// todo fix if need dynamic play
	/*ticker := time.NewTicker(timePeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
			case <-b.ctx.Done():
				// Stop the goroutine when the cancellation signal is received
				return
			}
		}
	}()*/
}

func (b *Battle) stop() {
	if b.cancel != nil {
		b.cancel()
	}
}
