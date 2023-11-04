package utils

func If[T any](trueCondition bool, v1, v2 T) T {
	if trueCondition {
		return v1
	}
	return v2
}
