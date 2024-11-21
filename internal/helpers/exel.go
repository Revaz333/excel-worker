package helpers

import (
	"github.com/360EntSecGroup-Skylar/excelize"
)

type Exel struct {
}

var exelAlphabet = []string{
	"", "A", "B", "C", "D", "E", "F", "G",
	"H", "I", "J", "K", "L", "M", "N",
	"O", "P", "Q", "R", "S", "T", "U",
	"V", "W", "X", "Y", "Z",
}

func (e Exel) ReadXls(path string) ([]interface{}, error) {
	data, err := excelize.OpenFile(path)

	if err != nil {
		return []interface{}{}, err
	}

	var rows []interface{}
	for _, sheet := range data.GetSheetMap() {
		row := data.GetRows(sheet)
		rows = append(rows, map[string]interface{}{"columns": row[0], "rows": row[1:]})
	}
	return rows, nil
}

func (e Exel) GetLetter(iteration int) string {
	lCount := len(exelAlphabet) - 2
	if iteration <= lCount {
		return exelAlphabet[iteration+1]
	}

	i := iteration / lCount

	remains := iteration - (lCount * i)

	return exelAlphabet[i] + exelAlphabet[remains]
}
