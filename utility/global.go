package utility

import (
	"fmt"
	"log"
	"time"

	"github.com/brianvoe/gofakeit"
	"github.com/fatih/color"
)

// Print task on terminal and log
// 1=TASK;2=DONE;3=INFO;4=WARNING;5=ERROR;
func InfoPrint(status int, msg string) {
	var info, logInfo string
	t := time.Now().Format(Dmyhms)
	switch status {
	case 1:
		logInfo = "START"
		info = color.HiCyanString(logInfo)
	case 2:
		logInfo = "DONE"
		info = color.HiGreenString(logInfo)
	case 3:
		logInfo = "INFO"
		info = color.HiBlueString(logInfo)
	case 4:
		logInfo = "WARNING"
		info = color.HiYellowString(logInfo)
	case 5:
		logInfo = "ERROR"
		info = color.HiRedString(logInfo)
	}
	log.Printf("%s : %s\n", logInfo, msg)
	fmt.Printf("%s [%s] %s\n", t, info, msg)
}

func DateRandom(minYear string, maxYear string) time.Time {
	min, _ := time.Parse(Dmy, minYear)
	max, _ := time.Parse(Dmy, maxYear)
	return gofakeit.DateRange(min, max)
}

func StatusString(i int) string {
	switch i {
	case 1:
		return "New"
	case 2:
		return "Great"
	case 3:
		return "Good"
	case 4:
		return "Bad"
	default:
		return "Undefined"
	}
}

func AvailabilityString(i int) string {
	switch i {
	case 1:
		return "Available"
	case 2:
		return "Read Only"
	case 3:
		return "Good"
	default:
		return "Undefined"
	}
}

func Compare(s []string, str string) bool {
	for _, v := range s {
		if v == str {
			return true
		}
	}
	return false
}

func RemoveIndex(s []string, index int) []string {
	ret := make([]string, 0)
	ret = append(ret, s[:index]...)
	return append(ret, s[index+1:]...)
}
