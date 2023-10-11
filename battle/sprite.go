package sprite

import (
	"context"
	"time"
)

const timePeriod = 50 * time.Millisecond //Millisecond
const rangeAttackUpmost = 300
const subNodeWidth = 100
const baseSpeed = 100

type DefenderType int
type LivingStatus int
type AtkStatus int
type AtkType int
type MoveType int

// 0-other, 1-wall, 2-training, 3-resource, 4-ground defender, 5-air defender
const (
	otherDefender DefenderType = iota
	wall
	trainingDefender
	resourceDefender
	groundDefender
	airDefender
)

const (
	normalAtk AtkType = iota
	rangeAtk
)

// 0-stand 1-ground, 2-fly, 3-bounce, 4-dig
const (
	stand MoveType = iota
	moveGround
	moveFly
	moveBounce
	moveDig
)

// -1-destroyed 1-存活
const (
	destroyed LivingStatus = iota - 1
	undeployed
	deployed
)

// 0-free, 1-searching, 2-moving, 3-attacking, 4-wallAttacking
const (
	free AtkStatus = iota
	searching
	moving
	attacking
	wallAttacking //no need to wait for the Target to be destroyed before recalculating the path
)

const (
	ResourceMoney = "MONEY"
	ResourceWater = "WATER"
)

type sprite interface {
	GetIndex() int
	GetBattle() *Battle
	GetTarget() int
	GetLivingStatus() int
	GetAtkStatus() int
	Destroyed()
	Hurt(damage int)
	Attack()
	Move()
	Search()
	start(ctx context.Context)
	stop()
}

type Resource struct {
	Preserve    int
	PreserveOri int
}
