package sprite

import (
	"encoding/json"
	"math"
	"reflect"
	"strconv"
	"strings"
)

func calDistanceToRectangle(ax, ay, dxMin, dxMax, dyMin, dyMax int) int {
	if ax < dxMin && ay < dyMin {
		return SubDistance(ax, ay, dxMin, dyMin)
	}
	if ax < dxMin && ay >= dyMin && ay <= dyMax {
		return dxMin - ax
	}
	if ax < dxMin && ay > dyMax {
		return SubDistance(ax, ay, dxMin, dyMax)
	}

	if ax >= dxMin && ax <= dxMax && ay < dyMin {
		return dyMin - ay
	}
	if ax >= dxMin && ax <= dxMax && ay >= dyMin && ay <= dyMax {
		return 0
	}
	if ax >= dxMin && ax <= dxMax && ay > dyMax {
		return ay - dyMax
	}

	if ax > dxMax && ay < dyMin {
		return SubDistance(ax, ay, dxMax, dyMin)
	}
	if ax > dxMax && ay >= dyMin && ay <= dyMax {
		return ax - dxMax
	}
	if ax > dxMax && ay > dyMax {
		return SubDistance(ax, ay, dxMax, dyMax)
	}
	return 0
}

func CalDirectionBetweenAttackerAndDefender(attacker *Attacker, defender *Defender) int {
	ax, ay := attacker.getSubPoint()
	dxMin, dxMax, dyMin, dyMax := defender.getSubPos()
	if ax < dxMin && ay < dyMin {
		return int(downright)
	}
	if ax < dxMin && ay >= dyMin && ay <= dyMax {
		return int(down)
	}
	if ax < dxMin && ay > dyMax {
		return int(downleft)
	}

	if ax >= dxMin && ax <= dxMax && ay < dyMin {
		return int(right)
	}
	if ax >= dxMin && ax <= dxMax && ay >= dyMin && ay <= dyMax {
		return int(in)
	}
	if ax >= dxMin && ax <= dxMax && ay > dyMax {
		return int(left)
	}

	if ax > dxMax && ay < dyMin {
		return int(upright)
	}
	if ax > dxMax && ay >= dyMin && ay <= dyMax {
		return int(up)
	}
	if ax > dxMax && ay > dyMax {
		return int(upleft)
	}
	return int(in)
}

func CalDistanceBetweenAttackerAndDefender(attacker *Attacker, defender *Defender) int {
	ax, ay := attacker.getSubPoint()
	dxMin, dxMax, dyMin, dyMax := defender.getSubPos()
	if ax < dxMin && ay < dyMin {
		return SubDistance(ax, ay, dxMin, dyMin)
	}
	if ax < dxMin && ay >= dyMin && ay <= dyMax {
		return dxMin - ax
	}
	if ax < dxMin && ay > dyMax {
		return SubDistance(ax, ay, dxMin, dyMax)
	}

	if ax >= dxMin && ax <= dxMax && ay < dyMin {
		return dyMin - ay
	}
	if ax >= dxMin && ax <= dxMax && ay >= dyMin && ay <= dyMax {
		return 0
	}
	if ax >= dxMin && ax <= dxMax && ay > dyMax {
		return ay - dyMax
	}

	if ax > dxMax && ay < dyMin {
		return SubDistance(ax, ay, dxMax, dyMin)
	}
	if ax > dxMax && ay >= dyMin && ay <= dyMax {
		return ax - dxMax
	}
	if ax > dxMax && ay > dyMax {
		return SubDistance(ax, ay, dxMax, dyMax)
	}
	return 0
}

func findClosetPointInRectangle(x, y, dxMin, dxMax, dyMin, dyMax int) (int, int) {
	if x < dxMin && y < dyMin {
		return dxMin, dyMin
	}
	if x < dxMin && y >= dyMin && y <= dyMax {
		return dxMin, y
	}
	if x < dxMin && y > dyMax {
		return dxMin, dyMax
	}

	if x >= dxMin && x <= dxMax && y < dyMin {
		return x, dyMin
	}
	if x >= dxMin && x <= dxMax && y >= dyMin && y <= dyMax {
		return x, y
	}
	if x >= dxMin && x <= dxMax && y > dyMax {
		return x, dyMax
	}

	if x > dxMax && y < dyMin {
		return dxMax, dyMin
	}
	if x > dxMax && y >= dyMin && y <= dyMax {
		return dxMax, y
	}
	if x > dxMax && y > dyMax {
		return dxMax, dyMax
	}
	return x, y
}

func calWiderDistanceToDefender(x, y int, defender *Defender) int {
	dxMin, dxMax, dyMin, dyMax := defender.getSubPos()
	if x < defender.XMin && y < defender.YMin {
		return SubDistance(straightCost*x+99, straightCost*y+99, dxMin, dyMin)
	}
	if x < defender.XMin && y >= defender.YMin && y <= defender.YMax {
		return dxMin - (straightCost*x + 99)
		//return straightCost * (defender.XMin - (X + 1))
	}
	if x < defender.XMin && y > defender.YMax {
		return SubDistance(straightCost*x+99, straightCost*y, dxMin, dyMax)
	}

	if x >= defender.XMin && x <= defender.XMax && y < defender.YMin {
		return dyMin - (straightCost*y + 99)
	}
	if x >= defender.XMin && x <= defender.XMax && y >= defender.YMin && y <= defender.YMax {
		return 0
	}
	if x >= defender.XMin && x <= defender.XMax && y > defender.YMax {
		return straightCost*y - dyMax
	}

	if x > defender.XMax && y < defender.YMin {
		return SubDistance(straightCost*x, straightCost*y+99, dxMax, dyMin)
	}
	if x > defender.XMax && y >= defender.YMin && y <= defender.YMax {
		return straightCost*x - dxMax
	}
	if x > defender.XMax && y > defender.YMax {
		return SubDistance(straightCost*x, straightCost*y, dxMax, dyMax)
	}
	return 0
}

// GetDirectionInt on which direction that A should move to
func GetDirectionInt(a *Node, b *Node) int {
	if a.P.X < b.P.X && a.P.Y < b.P.Y {
		return int(downright)
	}
	if a.P.X < b.P.X && a.P.Y == b.P.Y {
		return int(down)
	}
	if a.P.X < b.P.X && a.P.Y > b.P.Y {
		return int(downleft)
	}

	if a.P.X == b.P.X && a.P.Y < b.P.Y {
		return int(right)
	}
	if a.P.X == b.P.X && a.P.Y == b.P.Y {
		return int(in)
	}
	if a.P.X == b.P.X && a.P.Y > b.P.Y {
		return int(left)
	}

	if a.P.X > b.P.X && a.P.Y < b.P.Y {
		return int(upright)
	}
	if a.P.X > b.P.X && a.P.Y == b.P.Y {
		return int(up)
	}
	if a.P.X > b.P.X && a.P.Y > b.P.Y {
		return int(upleft)
	}
	return int(in)
}

func Distance(ax, ay, bx, by int) int {
	return int(math.Sqrt(math.Pow(float64(straightCost)*math.Abs(float64(ax-bx)), 2) + math.Pow(float64(straightCost)*math.Abs(float64(ay-by)), 2)))
}

func SubDistance(ax, ay, bx, by int) int {
	return int(math.Sqrt(math.Pow(math.Abs(float64(ax-bx)), 2) + math.Pow(math.Abs(float64(ay-by)), 2)))
}

func IsInMap(x, y, size int) bool {
	return x >= 0 && x < size && y >= 0 && y < size
}

func HasPriority(buildings []*Defender, elf *Attacker) bool {
	for _, building := range buildings {
		if building.isPriority(elf) && building.LivingStatus == deployed {
			return true
		}
	}
	return false
}

func ReflectCallHello(class, methodStr string, battleIndex int, index int, paramStr string) string {
	return "helloworld"
}

func ReflectCall(class, methodStr string, battleIndex int, index int, paramStr string) string {
	var target any
	if class == "battle" {
		target = GetBattle(battleIndex)
	} else if class == "attacker" {
		battle := GetBattle(battleIndex)
		target = battle.GetAttacker(index)
	} else if class == "defender" {
		battle := GetBattle(battleIndex)
		target = battle.GetDefender(index)
	} else {
		return "error:not support"
	}

	params := strings.Split(paramStr, ",")

	// Get the method you want to call using reflect
	method := reflect.ValueOf(target).MethodByName(methodStr)

	// Get the types of the arguments
	argTypes := make([]reflect.Type, method.Type().NumIn())
	for i := 0; i < method.Type().NumIn(); i++ {
		argTypes[i] = method.Type().In(i)
	}

	// Create new zero values of the argument types
	argValues := make([]reflect.Value, len(argTypes))
	for i, at := range argTypes {
		//refV := reflect.New(argTypes[i].Elem())
		if at.Kind() == reflect.Int {
			if len(params) <= i {
				return "error:wrong param length"
			}
			num, err := strconv.Atoi(params[i])
			if err != nil {
				return "error:" + err.Error()
			}
			argValues[i] = reflect.ValueOf(num)
		} else if at.Kind() == reflect.String {
			argValues[i] = reflect.ValueOf(params[i])
		} else {
			return "error:only support int and string"
		}
	}

	// Call the method with the arguments
	result := method.Call(argValues)
	if len(result) > 0 {
		buf, _ := json.Marshal(result[0].Interface())
		return string(buf)
	}
	return ""
}
