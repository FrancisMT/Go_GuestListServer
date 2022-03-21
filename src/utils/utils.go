package utils

import (
	"strconv"
	"time"
)

// GetHoursAndMinutesString Generates a string with the format hours::minutes
func GetHoursAndMinutesString() (result string) {
	result = strconv.Itoa(time.Now().Hour()) + ":" + strconv.Itoa(time.Now().Minute())
	return
}
