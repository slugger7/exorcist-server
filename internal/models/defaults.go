package models

import "strconv"

func DefualtBool(strVal string, def bool) bool {
	val, err := strconv.ParseBool(strVal)
	if err != nil {
		return def
	}
	return val
}

func DefaultInt(strVal string, def int) int {
	if strVal != "" {
		val, err := strconv.Atoi(strVal)
		if err == nil {
			return val
		}
	}

	return def
}
