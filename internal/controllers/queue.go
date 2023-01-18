package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"worker/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/mdigger/translit"
	"gorm.io/gorm"
)

type QueueController struct {
	Db *gorm.DB
}

var haveSlug = map[string][]map[string]interface{}{
	"service_categories":        {{"from": "name", "to": "slug"}},
	"service_subcategories":     {{"from": "name", "to": "slug"}},
	"service_subcategory_types": {{"from": "name", "to": "slug"}},
	"companies":                 {{"from": "name", "to": "slug"}},
	"articles":                  {{"from": "name", "to": "slug"}},
}

var useMethod = map[string][]map[string]interface{}{
	"users": {{"field": "password", "method": Md5}},
}

func (c QueueController) CheckQueue() (bool, error) {
	queues := []models.WorkerQueues{}
	c.Db.Table("worker_queues").Where("status =	?", 0).Select("*").Scan(&queues)

	if len(queues) == 0 {
		return false, nil
	}

	for _, queue := range queues {
		if queue.Driver == "import" {
			return c.createImportQueue(&queue), nil
		} else {
			fmt.Println("export")
		}
	}
	return true, nil
}

func (c QueueController) createImportQueue(data *models.WorkerQueues) bool {
	c.Db.Model(&models.WorkerQueues{}).Where("id = ?", data.ID).Updates(
		map[string]interface{}{"status": 1},
	)

	rows, err := c.readXls(data.FilePath)

	if err != nil {
		c.Db.Model(&models.WorkerQueues{}).Where("id = ?", data.ID).Updates(
			map[string]interface{}{"error": err.Error(), "status": 3},
		)
		return false
	}

	c.Db.Model(&models.WorkerQueues{}).Where("id = ?", data.ID).Updates(
		map[string]interface{}{"status": 2},
	)

	for _, row := range rows {
		// fmt.Println(row)
		_, err := c.buildQueue(data.ID, data.TableName, row.(map[string]interface{}))

		if err != nil {
			panic(err)
		}
	}

	// // fmt.Println(res)

	return false
}

func (c QueueController) buildQueue(queueId int, table string, data map[string]interface{}) (bool, error) {
	chunkedArrays := arrayChunk(data["rows"].([][]string), 1000)
	idIndex := arraySearch(data["columns"].([]string), "id")
	fmt.Println(idIndex)
	for _, val := range chunkedArrays {
		c.convertAndSave(
			table, idIndex.(int), data["columns"].([]string), val.([][]string),
		)
	}

	return false, nil
}

func (c QueueController) convertAndSave(table string, idIndex int, columns []string, rows [][]string) {
	var ids []string

	// if _, ok := useMethod[table]; ok {
	// 	index = arraySearch(columns, useMethod[table]["field"])
	// }
	// var slice []interface{}
	for _, val := range rows {
		if val[idIndex] != "" {
			ids = append(ids, val[idIndex])
		}

		if _, ok := useMethod[table]; ok {
			for _, method := range useMethod[table] {
				index := arraySearch(columns, method["field"])
				res := call(
					method["method"], []interface{}{val[index.(int)]},
				)

				val[index.(int)] = fmt.Sprintf("%s", res)
			}
		}

		if _, ok := haveSlug[table]; ok {
			for _, el := range haveSlug[table] {
				fromIndex := arraySearch(columns, el["from"])
				toIndex := arraySearch(columns, el["to"])

				val[toIndex.(int)] = translit.Ru(val[fromIndex.(int)])
			}
		}

		for _, field := range val {
			fmt.Println(field + "dsdc")
		}
	}

	// fmt.Println(slice)
	// arr, _ := json.Marshal(slice)
	// importQueue := models.ImportQueue{
	// 	Query: string(arr),
	// }
	// res := c.Db.Create(&importQueue)

	// if res.Error != nil {
	// 	panic(res.Error)
	// }
}

func (c QueueController) readXls(path string) ([]interface{}, error) {
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

func Md5(data string) string {
	response := md5.Sum([]byte(data))
	return fmt.Sprintf("%s", hex.EncodeToString(response[:]))
}

func arrayChunk(slice [][]string, size int) []interface{} {
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

func arraySearch(slice []string, search interface{}) interface{} {
	for i, v := range slice {
		if v == search {
			return i
		}
	}

	return false
}

func ArrayColumn(input map[string]map[string]interface{}, columnKey string) []interface{} {
	columns := make([]interface{}, 0, len(input))
	for _, val := range input {
		if v, ok := val[columnKey]; ok {
			columns = append(columns, v)
		}
	}

	return columns
}

func call(fn interface{}, args []interface{}) interface{} {
	method := reflect.ValueOf(fn)
	var inputs []reflect.Value

	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	return method.Call(inputs)[0]
}
