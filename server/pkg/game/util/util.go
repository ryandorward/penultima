package util

import (
	"math"
	"encoding/json"
	"fmt"
)

const (
	xpLevelFactor = 500
)

func Dist(x1, y1, x2, y2 int) float64 {
	return math.Sqrt(math.Pow(float64(x1-x2), 2) + math.Pow(float64(y1-y2), 2))
}

func Modifier(stat int) int {
	return (stat - 10) / 2
}

func WorthXP(level int) int {
	return level * 100
}

func XPForLevel(level int) int {
	return int(xpLevelFactor*math.Pow(float64(level), 2) - float64(xpLevelFactor*level))
}

func WrapMod(w, mod int) int{
	return (w%mod + mod)%mod;
}

func IntAbs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

func IntMin(a,b,c int) int {
	if a < b {	
		if a < c {
			return a
		} else {
			return c
		}
	} else {
		if b < c {
			return b
		} else {
			return c
		}
	}	
}

func WrapDiff(a,b, mod int) int{
	diff1 := IntAbs(a - b)
	diff2 := IntAbs(a - (b+mod))
	diff3 := IntAbs(a - (b-mod))
	return IntMin(diff1,diff2,diff3)
}

func PrettyPrint(i interface{}) {
	s, _ := json.MarshalIndent(i, "", "\t")
	fmt.Printf("%+v\n",string(s))
}

func Bool2Int (b bool) int {
	if b {
		return 1
	} else {
		return 0
	}
}