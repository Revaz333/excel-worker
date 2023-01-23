package controllers

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
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

var seededRand *rand.Rand = rand.New(
	rand.NewSource(time.Now().UnixNano()))

var haveSlug = map[string][]map[string]interface{}{
	"service_categories":        {{"from": "name", "to": "slug"}},
	"service_subcategories":     {{"from": "name", "to": "slug"}},
	"service_subcategory_types": {{"from": "name", "to": "slug"}},
	"companies":                 {{"from": "name", "to": "slug"}},
	"articles":                  {{"from": "name", "to": "slug"}},
}

var ids []string

var useMethod = map[string][]map[string]interface{}{
	"users": {{"field": "password", "method": Md5}},
}

var generateToken = map[string][]map[string]interface{}{
	"users": {{"field": "token", "rune": "email"}},
}

func (c QueueController) CheckQueue() (bool, error) {
	queues, err := c.Db.GetQueues(0, "import")

	if err != nil {
		logrus.Fatalf(err.Error())
		return false, nil
	}

	if len(queues) == 0 {
		return false, nil
	}

	for _, queue := range queues {
		return c.createImportQueue(&queue), nil
	}
	return true, nil
}

func (c QueueController) createImportQueue(data *models.WorkerQueues) bool {
	c.Db.UpdateQueue(data.ID, 1, nil)
	path := fmt.Sprintf("%v", data.FilePath)
	rows, err := c.Exel.ReadXls(path)

	if err != nil {
		c.Db.UpdateQueue(data.ID, 3, err.Error())
		return false
	}

	c.Db.UpdateQueue(data.ID, 2, nil)
	types := c.Db.GetTypes(data.TableName)

	for _, row := range rows {
		err := c.buildQueue(data.ID, data.TableName, row.(map[string]interface{}), types)

		if err != nil {
			c.Db.UpdateQueue(data.ID, 3, err.Error())
			return false
		}
	}

	c.Db.UpdateQueue(data.ID, 4, nil)
	return true
}

func (c QueueController) buildQueue(queueId int, table string, data map[string]interface{}, types []string) error {
	chunkedArrays := c.Array.ArrayChunk(data["rows"].([][]string), 1000)
	idIndex := c.Array.ArraySearch(data["columns"].([]string), "id")
	for _, val := range chunkedArrays {
		err := c.convertAndSave(
			queueId, table, idIndex.(int), data["columns"].([]string), val.([][]string), types,
		)

		if err != nil {
			return err
		}
	}

	if len(ids) > 0 {
		err := c.BuildDelete(table, ids, queueId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c QueueController) convertAndSave(
	queueId int, table string, idIndex int, columns []string, rows [][]string, types []string,
) error {
	var sqlValues []string
	var updates [][]string

	for _, val := range rows {
		if len(val[idIndex]) > 0 {
			ids = append(ids, val[idIndex])
		}

		if len(val[idIndex]) > 0 && c.Db.CheckRow(table, val[idIndex]) {
			updates = append(updates, val)
			continue
		}

		sqlVal := c.BuildRow(table, val, columns, types)
		sqlValues = append(sqlValues, sqlVal)
	}

	if len(sqlValues) > 0 {
		sql := fmt.Sprintf("INSERT INTO %v (%v) VALUES %v",
			table, strings.Join(columns, ","), strings.Join(sqlValues, ","),
		)

		if err := c.Db.Insert("import_queues", sql, queueId); err != nil {
			return err
		}
	}

	if len(updates) > 0 {
		err := c.BuildUpdate(table, idIndex, columns, updates, types, queueId)
		if err != nil {
			return err
		}
	}

	return nil
}

func (c QueueController) BuildRow(table string, row []string, columns []string, types []string) string {
	var sqlValue []string

	if _, ok := generateToken[table]; ok {
		for _, method := range generateToken[table] {
			from := c.Array.ArraySearch(columns, method["rune"])
			to := c.Array.ArraySearch(columns, method["field"])
			res := GenerateToken(
				50, row[from.(int)],
			)

			row[to.(int)] = fmt.Sprintf("%s", res)
		}
	}

	if _, ok := useMethod[table]; ok {
		for _, method := range useMethod[table] {
			index := c.Array.ArraySearch(columns, method["field"])
			res := Call(
				method["method"], []interface{}{row[index.(int)]},
			)

			row[index.(int)] = fmt.Sprintf("%s", res)
		}
	}

	if _, ok := haveSlug[table]; ok {
		for _, el := range haveSlug[table] {
			fromIndex := c.Array.ArraySearch(columns, el["from"])
			toIndex := c.Array.ArraySearch(columns, el["to"])

			row[toIndex.(int)] = Replace(translit.Ru(row[fromIndex.(int)]), "-")
		}
	}

	for i, field := range row {
		if strings.Index(types[i], "date") != -1 && field == "" {
			sqlValue = append(sqlValue, "CURRENT_TIMESTAMP")
		} else if field == "" {
			sqlValue = append(sqlValue, "NULL")
		} else {
			sqlValue = append(sqlValue, "'"+strings.Replace(field, "'", "", -1)+"'")
		}
	}

	return "(" + strings.Join(sqlValue, ",") + ")"
}

func (c QueueController) BuildUpdate(
	table string, idIndex int, columns []string, rows [][]string, types []string, queueId int,
) error {
	for _, val := range rows {
		var sqlVal []string

		if _, ok := generateToken[table]; ok {
			for _, method := range generateToken[table] {
				from := c.Array.ArraySearch(columns, method["rune"])
				to := c.Array.ArraySearch(columns, method["field"])
				res := GenerateToken(
					50, val[from.(int)],
				)

				val[to.(int)] = fmt.Sprintf("%s", res)
			}
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

				val[toIndex.(int)] = Replace(translit.Ru(val[fromIndex.(int)]), "-")
			}
		}

		for i, field := range val {
			if strings.Index(types[i], "date") != -1 && field == "" {
				sqlVal = append(sqlVal, fmt.Sprintf("%s = %v", columns[i], "CURRENT_TIMESTAMP"))
			} else if field == "" {
				sqlVal = append(sqlVal, fmt.Sprintf("%s = %v", columns[i], "NULL"))
			} else {
				sqlVal = append(sqlVal, fmt.Sprintf("%s = %v", columns[i], "'"+strings.Replace(field, "'", "", -1)+"'"))
			}
		}
		sql := fmt.Sprintf("UPDATE %s SET %v WHERE id = %v", table, strings.Join(sqlVal, ","), val[idIndex])
		if err := c.Db.Insert("import_queues", sql, queueId); err != nil {
			return err
		}
	}

	return nil
}

func (c QueueController) BuildDelete(table string, ids []string, queueId int) error {
	sql := fmt.Sprintf("DELETE FROM %s  WHERE id NOT IN (%v)", table, strings.Join(ids, ","))
	if err := c.Db.Exec(sql); err != nil {
		return err
	}

	return nil
}
