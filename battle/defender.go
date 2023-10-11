package sprite

import (
	"context"
	"math"
	"sync/atomic"
	"time"
)

type Defender struct {
	Index                  int
	ctx                    context.Context
	cancel                 context.CancelFunc
	XMin, XMax, YMin, YMax int
	DefenderType           DefenderType // 1-wall, 2-normal, 3-resourceDefender, 4-groundDefender, 5-airDefender
	atkRange               int
	atkFrequency           int
	damage                 int          //vary with level
	passLoss               int          //-1-nopass, wall(vary with level)
	LivingStatus           LivingStatus //0-undeployed, 1-deployed, -1-destroyed
	AtkStatus              AtkStatus    //0-free, 1-searching, 2-moving, 3-attacking, 4-wall destroy
	Health                 int          // present health value
	healthOri              int          // health value before attacked
	Target                 int          // sprite id in attackerMap
	battle                 *Battle
	ResourceMap            map[string]Resource //resource contain
	playCount              uint32              // record the times of the play function be called
}

func NewDefender(xMin, yMin, size, damage, atkRange, atkFrequency, health int, defenderType int, resourceMap map[string]Resource) *Defender {
	defender := &Defender{
		XMin:         xMin,
		XMax:         xMin + size - 1,
		YMin:         yMin,
		YMax:         yMin + size - 1,
		DefenderType: DefenderType(defenderType),
		atkRange:     atkRange,
		atkFrequency: atkFrequency,
		damage:       damage,
		LivingStatus: deployed,
		AtkStatus:    free,
		Target:       -1,
		Health:       health,
		healthOri:    health,
		ResourceMap:  resourceMap,
	}
	return defender
}

func (d *Defender) GetBattle() *Battle {
	return d.battle
}

func (d *Defender) CanPlay() bool {
	return d.DefenderType != wall && d.LivingStatus == deployed
}

func (d *Defender) play() {
	atomic.AddUint32(&d.playCount, 1)
	if d.LivingStatus != deployed {
		return
	}
	if d.DefenderType != groundDefender && d.DefenderType != airDefender {
		return
	}
	d.Attack()
	d.Search()
}

func (d *Defender) start(parentCtx context.Context) {
	ctx, cancel := context.WithCancel(parentCtx)
	d.ctx = ctx
	d.cancel = cancel
	if d.DefenderType != groundDefender && d.DefenderType != airDefender {
		return
	}

	ticker := time.NewTicker(timePeriod)
	go func() {
		for {
			select {
			case <-ticker.C:
				d.Attack()
				d.Search()
			case <-d.ctx.Done():
				// Stop the goroutine when the cancellation signal is received
				return
			}
		}
	}()
}

func (d *Defender) stop() {
	if d.cancel != nil {
		d.cancel()
	}
}

func (d *Defender) getSubPos() (xMin, xMax, yMin, yMax int) {
	return d.XMin * straightCost, (d.XMax+1)*straightCost - 1, d.YMin * straightCost, (d.YMax+1)*straightCost - 1
}

func (d *Defender) GetSubPointMin() *Point {
	xMin, _, yMin, _ := d.getSubPos()
	return &Point{
		X: xMin,
		Y: yMin,
	}
}

func (d *Defender) GetSubPointMax() *Point {
	_, xMax, _, yMax := d.getSubPos()
	return &Point{
		X: xMax,
		Y: yMax,
	}
}

func (d *Defender) Hurt(damage int) {
	d.Health -= damage
	if d.DefenderType == wall && d.battle != nil {
		d.battle.nodes[d.XMin][d.YMin].cc -= damage
	}
	var ptc float64
	ptc = float64(damage) / float64(d.healthOri)
	for _, resource := range d.ResourceMap {
		resource.Preserve -= int(math.Round(ptc * float64(resource.PreserveOri)))
	}
	if d.Health <= 0 {
		for _, resource := range d.ResourceMap {
			resource.Preserve = 0
		}
		d.Destroyed()
		return
	}
}

func (d *Defender) GetDefenderType() int {
	return int(d.DefenderType)
}

func (d *Defender) GetAtkRange() int {
	return d.atkRange
}

func (d *Defender) Attack() {
	if d.atkFrequency == 0 || int(d.playCount)%d.atkFrequency != 0 {
		return
	}
	if d.AtkStatus != attacking || d.battle == nil {
		return
	}
	target := d.battle.Attackers[d.Target]
	atkPos := target.getNode()
	if atkPos == nil {
		return
	}
	dis := calWiderDistanceToDefender(atkPos.P.X, atkPos.P.Y, d)
	if dis <= d.atkRange {
		target.Hurt(d.damage)
	}
	if target.Health <= 0 || dis > d.atkRange {
		d.AtkStatus = free
	}
}

func (d *Defender) Move() {
	return
}

func (d *Defender) Destroyed() {
	if d.battle == nil {
		return
	}
	d.LivingStatus = destroyed
	bm := d.battle
	for i := d.XMin; i <= d.XMax; i++ {
		for j := d.YMin; j <= d.YMax; j++ {
			// node destroy
			bm.matrix[i][j] = 0
			bm.nodes[i][j].cc = 0
		}
	}
	d.stop()

	// end the game
	canEnd := true
	for _, defender := range bm.Defenders {
		if defender.LivingStatus == deployed && defender.DefenderType != wall {
			canEnd = false
		}
	}
	if canEnd {
		bm.stop()
		return
	}
	// set Attackers' status to free
	for _, attacker := range bm.Attackers {
		if attacker.Target == d.Index {
			attacker.AtkStatus = free
		}
	}
}

func (d *Defender) isPriority(elf *Attacker) bool {
	return d.DefenderType == elf.atkPriority
}

func (d *Defender) Search() {

	if d.AtkStatus != free || d.battle == nil {
		return
	}
	var preTarget *Attacker
	preDis := d.atkRange
	for _, attacker := range d.battle.Attackers {
		if attacker.LivingStatus == destroyed {
			continue
		}
		atkPos := attacker.path[attacker.PathPos]
		dis := calWiderDistanceToDefender(atkPos.P.X, atkPos.P.Y, d)
		if dis > d.atkRange {
			continue
		}
		//if dis <= d.atkRange {
		if attacker.AtkStatus == attacking {
			preTarget = attacker
			//preDis = dis
			break
		}
		if preDis > dis {
			preTarget = attacker
			preDis = dis
		}
		//}
	}
	if preTarget != nil {
		d.AtkStatus = attacking
		d.Target = preTarget.Index
	}
}
