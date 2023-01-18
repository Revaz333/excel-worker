package helpers

type Arrays struct {
}

func (a Arrays) ArrayChunk(slice [][]string, size int) []interface{} {
	var divided []interface{}

	for i := 0; i < len(slice); i += size {
		end := i + size

		if end > len(slice) {
			end = len(slice)
		}

		divided = append(divided, slice[i:end])
	}

	return divided
}

func (a Arrays) ArraySearch(slice []string, search interface{}) interface{} {
	for i, v := range slice {
		if v == search {
			return i
		}
	}

	return false
}

func (a Arrays) ArrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}

	return columns
}
