package controllers

import (
	"fmt"
	"worker/internal/helpers"
	"worker/models"

	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/sirupsen/logrus"
)

const savePath = "/var/www/element/public/upload/exel/"

type ExportController struct {
	Db    helpers.Db
	Array helpers.Arrays
	Exel  helpers.Exel
}

func (c ExportController) CheckQueues() (bool, error) {
	queues, err := c.Db.GetQueues(0, "export")

	if err != nil {
		logrus.Fatalf(err.Error())
		return false, nil
	}

	if len(queues) == 0 {
		return false, nil
	}
	for _, queue := range queues {
		err := c.ExportTable(queue)

		if err != nil {
			c.Db.UpdateQueue(queue.ID, 3, err.Error())
			return false, err
		}
	}

	return true, nil
}

func (c ExportController) ExportTable(queue models.WorkerQueues) error {
	c.Db.UpdateQueue(queue.ID, 1, nil)
	cols := queue.Columns

	if cols == "" {
		cols = "*"
	}

	data, columns, err := c.Db.GetTableData(queue.TableName, cols)

	if err != nil {
		return err
	}

	f := excelize.NewFile()

	for i, col := range columns {
		letter := c.Exel.GetLetter(i)
		f.SetCellValue("Sheet1", letter+"1", col)
	}

	i := 2
	exportedRows := 1
	total := len(data)
	for _, row := range data {
		fmt.Println(exportedRows, "/", total)
		values := row.(map[string]interface{})

		for j := 0; j < len(values); j++ {

			letter := c.Exel.GetLetter(j)
			val := values[columns[j]]
			if val == "%!s(<nil>)" {
				val = " "
			}

			f.SetCellValue("Sheet1", letter+fmt.Sprintf("%v", i), val)
		}
		exportedRows++
		i++
	}

	filePath := savePath + queue.TableName + Md5(GenerateToken(10, queue.TableName)) + ".xlsx"
	if err := f.SaveAs(filePath); err != nil {
		return err
	}

	query := fmt.Sprintf("UPDATE worker_queues SET file_path = '%s' WHERE id = %v", filePath, queue.ID)
	err = c.Db.Exec(query)

	if err != nil {
		return err
	}

	c.Db.UpdateQueue(queue.ID, 5, nil)

	return nil
}
