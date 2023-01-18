package controllers

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"reflect"
	"worker/internal/helpers"
	"worker/models"

	"github.com/mdigger/translit"
	"github.com/sirupsen/logrus"
)

type QueueController struct {
	Db    helpers.Db
	Array helpers.Arrays
	Exel  helpers.Exel
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
	queues, err := c.Db.GetQueues()

	if err != nil {
		logrus.Fatalf(err.Error())
		return false, nil
	}

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
	c.Db.UpdateQueue(data.ID, 1, nil)

	rows, err := c.Exel.ReadXls(data.FilePath)

	if err != nil {
		c.Db.UpdateQueue(data.ID, 3, err.Error())
		return false
	}

	c.Db.UpdateQueue(data.ID, 2, nil)
	// cols := c.Db.GetTypes(data.TableName)
	for _, row := range rows {
		// fmt.Println(row)
		_, err := c.buildQueue(data.ID, data.TableName, row.(map[string]interface{}))

		if err != nil {
			logrus.Fatalf(err.Error())
		}
	}

	// // fmt.Println(res)

	return false
}

func (c QueueController) buildQueue(queueId int, table string, data map[string]interface{}) (bool, error) {
	chunkedArrays := c.Array.ArrayChunk(data["rows"].([][]string), 1000)
	idIndex := c.Array.ArraySearch(data["columns"].([]string), "id")
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
				index := c.Array.ArraySearch(columns, method["field"])
				res := Call(
					method["method"], []interface{}{val[index.(int)]},
				)

				val[index.(int)] = fmt.Sprintf("%s", res)
			}
		}

		if _, ok := haveSlug[table]; ok {
			for _, el := range haveSlug[table] {
				fromIndex := c.Array.ArraySearch(columns, el["from"])
				toIndex := c.Array.ArraySearch(columns, el["to"])

				val[toIndex.(int)] = translit.Ru(val[fromIndex.(int)])
			}
		}

		// for _, field := range val {
		// 	fmt.Println(field + "dsdc")
		// }
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

func Md5(data string) string {
	response := md5.Sum([]byte(data))
	return fmt.Sprintf("%s", hex.EncodeToString(response[:]))
}

func Call(fn interface{}, args []interface{}) interface{} {
	method := reflect.ValueOf(fn)
	var inputs []reflect.Value

	for _, v := range args {
		inputs = append(inputs, reflect.ValueOf(v))
	}

	return method.Call(inputs)[0]
}
