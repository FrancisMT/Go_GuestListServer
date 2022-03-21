package utilstest

import (
	"guestListChallenge/src/utils"
	"strconv"
	"strings"
	"testing"
)

// TestGetHoursAndMinutesString Tests the correctness of the GetHoursAndMinutesString return
func TestGetHoursAndMinutesString(t *testing.T) {
	timeString := utils.GetHoursAndMinutesString()
	timeStringContents := strings.Split(timeString, ":")

	if len(timeStringContents) != 2 {
		t.Error("Wrong hour format")
	} else {
		for _, timeString := range timeStringContents {
			if _, err := strconv.Atoi(timeString); err != nil {
				t.Error("Wrong hour format")
			}
		}
	}
}
