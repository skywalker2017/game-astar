package sprite

import (
	"encoding/json"
	"errors"
	"fmt"
)

func GetDefenderType(defenderType string) DefenderType {
	if defenderType == "DEFENSE" {
		return groundDefender
	}
	if defenderType == "RESOURCE" {
		return resourceDefender
	}
	return otherDefender
}

func GetDefenderLocation(locationStr string) (Point, error) {
	var locationArray [][]int
	err := json.Unmarshal([]byte(locationStr), &locationArray)
	if err != nil {
		return Point{}, err
	}
	if len(locationArray) == 0 || len(locationArray[0]) == 0 {
		return Point{}, errors.New(fmt.Sprintf("illegal locationStr:%s", locationStr))
	}
	return Point{
		X: locationArray[0][0],
		Y: locationArray[0][1],
	}, nil
}
