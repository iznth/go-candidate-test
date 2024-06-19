package main

func IsFibNumber(valToCheck, previousValue, currentValue int) bool {
	if valToCheck < 0 {
		return false
	}
	if valToCheck == currentValue {
		return true
	}
	if currentValue > valToCheck {
		return false
	}
	return IsFibNumber(valToCheck, currentValue, previousValue+currentValue)
}
