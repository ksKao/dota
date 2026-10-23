package lib

import (
	"fmt"
	"math"
)

func FormatMinutesSeconds(timestamp float32) string {
	sign := ""
	if timestamp < 0 {
		sign = "-"
	}

	t := math.Abs(float64(timestamp))
	minutes := int(t / 60)
	seconds := int(math.Mod(t, 60))

	return fmt.Sprintf("%s%d:%02d", sign, minutes, seconds)
}
