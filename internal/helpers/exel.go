package helpers

import "github.com/360EntSecGroup-Skylar/excelize"

type Exel struct {
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
