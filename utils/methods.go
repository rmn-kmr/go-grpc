package utils

import (
	"strings"
	"time"
)

func GetTimeFromEpoch(epoch int64) time.Time {
	return time.Unix(epoch, 0)
}

// GetFirstMiddleLastName Return first, middle and last name
func GetFirstMiddleLastName(name string) (string, string, string) {
	nameArray := strings.Split(name, " ")
	if len(nameArray) == 1 {
		return nameArray[0], " ", " "
	} else if len(nameArray) == 2 {
		return nameArray[0], " ", nameArray[1]
	} else if len(nameArray) > 2 {
		middleName := ""
		for i, namePart := range nameArray {
			if i == 0 || i == len(nameArray)-1 || namePart == "" || namePart == " " {
				continue
			}
			if middleName != "" {
				middleName = middleName + " " + namePart
			} else {
				middleName = namePart
			}
		}
		return nameArray[0], middleName, nameArray[len(nameArray)-1]
	}

	return "", "", ""
}

// IfThenElse evaluates a condition, if true returns the first parameter otherwise the second
func IfThenElse(condition bool, a interface{}, b interface{}) interface{} {
	if condition {
		return a
	}
	return b
}

func StringInSlice(a string, list []string) bool {
	for _, b := range list {
		if b == a {
			return true
		}
	}
	return false
}
