package sprite

import (
	"context"
	"sync"
	"sync/atomic"
	"time"
)

// sync
var lock sync.Mutex

type Attacker struct {
	Index        int
	ctx          context.Context
	cancel       context.CancelFunc
	Pos          *ElfNode
	SubPos       *SubNode
	PathPos      int //Index in path
	atkRange     int //range that could attack
	atkFrequency int //attack once when play atkFrequency times
	suicide      bool
	damage       int
	deathDamage  int          //damage when dead
	atkPriority  DefenderType //0-normal 1-wall, 2-resourceDefender, 3-groundDefender, 4-airDefender
	atkType      AtkType      //0-normalAtk, 1-rangeAtk
	speed        int          //0-stand
	moveType     MoveType     //0-stand 1-ground, 2-fly, 3-bounce, 4-dig
	LivingStatus LivingStatus //0-undeploy, 1-deployed, -1-destroyed
	AtkStatus    AtkStatus    //0-free, 1-searching, 2-moving, 3-attacking, 4-wall destroy
	heading      *Node
	path         []*Node
	Target       int //sprite id in Defenders
	Health       int
	battle       *Battle
	playCount    uint32 // record the times of the play function be called
}

func NewAttacker(x, y, atkRange, damage, deathDamage, atkTypeInt, atkPriorityInt, atkFrequency, health, speed int, suicide bool) *Attacker {
	attacker := &Attacker{
		Pos: newElfNode(x/subNodeWidth, y/subNodeWidth),
		SubPos: &SubNode{
			P: NewPoint(x%subNodeWidth, y%subNodeWidth),
		},
		PathPos:      0,
		atkRange:     atkRange,
		suicide:      suicide,
		damage:       damage,
		deathDamage:  deathDamage,
		atkPriority:  DefenderType(atkPriorityInt),
		atkFrequency: atkFrequency,
		atkType:      AtkType(atkTypeInt), //0-normalAtk, 1-rangeAtk
		speed:        speed,
		moveType:     moveGround,
		LivingStatus: deployed,
		AtkStatus:    free,
		heading:      nil,
		Target:       -1,
		Health:       health,
	}
	return attacker
}

func (a *Attacker) GetBattle() *Battle {
	return a.battle
}

func (a *Attacker) GetNextPath() *Node {
	if len(a.path) <= a.PathPos {
		return a.path[len(a.path)-1]
	}
	return a.path[a.PathPos+1]
}

func (a *Attacker) play() {
	atomic.AddUint32(&a.playCount, 1)
	if a.LivingStatus != deployed {
		return
	}
	a.Move()
	a.Attack()
	a.Search()
}

func (a *Attacker) start(parentCtx context.Context) {
	if a.battle == nil {
		return
	}
	a.path = []*Node{a.battle.nodes[a.Pos.P.X][a.Pos.P.Y]}
	ticker := time.NewTicker(timePeriod)
	ctx, cancel := context.WithCancel(parentCtx)
	a.ctx = ctx
	a.cancel = cancel
	go func() {
		var ite int
		for {
			select {
			case <-ticker.C:
				ite++
				a.Move()
				if ite%10 == 0 {
					a.Attack()
					a.Search()
				}
			case <-a.ctx.Done():
				// Stop the goroutine when the cancellation signal is received
				return
			}
		}
	}()
}

func (a *Attacker) getSubPoint() (x, y int) {
	x = a.Pos.Node.P.X*straightCost + a.SubPos.P.X
	y = a.Pos.Node.P.Y*straightCost + a.SubPos.P.Y
	return x, y
}

func (a *Attacker) stop() {
	if a.cancel != nil {
		a.cancel()
	}
}

func (a *Attacker) Destroyed() {
	a.LivingStatus = destroyed
	if a.deathDamage != 0 {
		x, y := a.getSubPoint()
		a.rangeAttack(x, y, a.deathDamage)
	}
	a.stop()
}

func (a *Attacker) Hurt(damage int) {
	if a.LivingStatus == destroyed {
		return
	}
	a.Health -= damage
	if a.Health <= 0 {
		a.Destroyed()
		return
	}
}

func (a *Attacker) rangeAttack(x, y, damage int) {
	if a.battle == nil {
		return
	}
	for _, defender := range a.battle.Defenders {
		xMin, xMax, yMin, yMax := defender.getSubPos()
		dis := calDistanceToRectangle(x, y, xMin, xMax, yMin, yMax)
		if dis < rangeAttackUpmost {
			defender.Hurt((damage*rangeAttackUpmost - damage*dis) / rangeAttackUpmost)
		}
	}
}

func (a *Attacker) Attack() {
	if a.atkFrequency == 0 || int(a.playCount)%a.atkFrequency != 0 {
		return
	}
	if a.battle == nil {
		return
	}
	if a.AtkStatus != attacking && a.AtkStatus != wallAttacking {
		return
	}
	target := a.battle.Defenders[a.Target]
	if a.atkType == normalAtk {
		target.Hurt(a.damage)
		/*if Target.Health <= 0 {
			a.AtkStatus = free
		}*/
	}
	if a.atkType == rangeAtk {
		x, y := a.getSubPoint()
		xMin, xMax, yMin, yMax := target.getSubPos()
		px, py := findClosetPointInRectangle(x, y, xMin, xMax, yMin, yMax)
		a.rangeAttack(px, py, a.damage)
	}
	//if suicide
	if a.suicide {
		a.Destroyed()
	}
}

func (a *Attacker) getNextPos() *Node {
	if a.PathPos+1 >= len(a.path) {
		return nil
	}
	return a.path[a.PathPos+1]
}

func (a *Attacker) subMove(direction Direction) bool {
	dx, dy := a.SubPos.getStep(int(direction))
	dx, dy = dx*a.speed, dy*a.speed
	if dx < baseSpeed && dx > 0 {
		dx = baseSpeed
	}
	if dy < baseSpeed && dy > 0 {
		dy = baseSpeed
	}
	dx, dy = dx/baseSpeed, dy/baseSpeed
	a.SubPos.P.X += dx * Moves[direction][0]
	a.SubPos.P.Y += dy * Moves[direction][1]
	if a.SubPos.P.X < 0 || a.SubPos.P.X > 99 || a.SubPos.P.Y < 0 || a.SubPos.P.Y > 99 {
		nextNode := a.getNextPos()
		if nextNode != nil {
			a.PathPos += 1
			a.Pos.Node = a.path[a.PathPos]
			if a.SubPos.P.X < 0 {
				a.SubPos.P.X += subNodeWidth
			}
			if a.SubPos.P.X > subNodeWidth-1 {
				a.SubPos.P.X -= subNodeWidth
			}
			if a.SubPos.P.Y < 0 {
				a.SubPos.P.Y += subNodeWidth
			}
			if a.SubPos.P.Y > subNodeWidth-1 {
				a.SubPos.P.Y -= subNodeWidth
			}
			return false
		}
		if nextNode == nil {
			if a.SubPos.P.X < 0 {
				a.SubPos.P.X = 0
			}
			if a.SubPos.P.X > subNodeWidth-1 {
				a.SubPos.P.X = subNodeWidth - 1
			}
			if a.SubPos.P.Y < 0 {
				a.SubPos.P.Y = 0
			}
			if a.SubPos.P.Y > subNodeWidth-1 {
				a.SubPos.P.Y = subNodeWidth - 1
			}
			return (a.SubPos.P.X == 0 || a.SubPos.P.X == subNodeWidth-1) && (a.SubPos.P.Y == 0 || a.SubPos.P.Y == subNodeWidth-1)
		}
	}
	return false
}

func (a *Attacker) Move() {
	if a.battle == nil {
		return
	}
	if a.AtkStatus != moving || a.Target < 0 {
		return
	}

	direction := CalDirectionBetweenAttackerAndDefender(a, a.battle.Defenders[a.Target])

	dis := CalDistanceBetweenAttackerAndDefender(a, a.battle.Defenders[a.Target])
	if dis <= a.atkRange {
		//check distance
		a.AtkStatus = a.attackStatus()
	} else {
		//normal move
		nextNode := a.getNextPos()
		//attack move
		if nextNode != nil {
			direction = GetDirectionInt(a.Pos.Node, nextNode)
		}
		isEnd := a.subMove(Direction(direction))
		if isEnd {
			a.AtkStatus = a.attackStatus()
		}
	}
}

func (a *Attacker) Search() {
	if a.AtkStatus != free || a.battle == nil {
		return
	}
	btMp := a.battle
	elfMp := newElfMp(btMp.nodes)
	var path []*Node
	var target int
	a.AtkStatus = searching
	if a.atkPriority == wall {
		path, target = AstarBomber(a, btMp.matrix, elfMp, btMp.Defenders)
	} else {
		path, target = AstarBuilding(a, btMp.matrix, elfMp, btMp.Defenders)
	}
	if path != nil {
		a.path = path
		a.Target = target
	}
	a.PathPos = 0
	a.AtkStatus = moving
}

func (a *Attacker) attackStatus() AtkStatus {
	target := a.battle.Defenders[a.Target]
	distance := CalDistanceBetweenAttackerAndDefender(a, target)
	if distance <= a.atkRange {
		if target.DefenderType != wall {
			return attacking
		}
		return wallAttacking
	}
	return free
}

func (a *Attacker) getNode() *Node {
	if a.path == nil || len(a.path) <= a.PathPos {
		return nil
	}
	return a.path[a.PathPos]
}
