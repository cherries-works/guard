package utils

func GetIndexOf[T comparable](arr []T, item T) int {
	index := 0

	for i, _item := range arr {
		if _item == item {
			index = i
		}
	}

	return index
}

func Includes[T comparable](arr []T, item T) bool {
	for _, _item := range arr {
		if _item == item {
			return true
		}
	}

	return false
}
