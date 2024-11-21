package helpers

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"worker/models"

	"github.com/sirupsen/logrus"
)

type Db struct {
	Db *sql.DB
}

func (d Db) GetQueues(status int, driver string) ([]models.WorkerQueues, error) {
	sql := "SELECT * FROM worker_queues WHERE status = ? AND driver = ?"
	queues := []models.WorkerQueues{}
	rows, err := d.Db.Query(sql, status, driver)

	if err != nil {
		logrus.Fatalf(err.Error())
		return []models.WorkerQueues{}, nil
	}

	for rows.Next() {
		var r models.WorkerQueues
		err := rows.Scan(&r.ID, &r.Driver, &r.FilePath, &r.TableName, &r.Error, &r.Status, &r.Columns)
		if err != nil {
			logrus.Fatalf(err.Error())
			return []models.WorkerQueues{}, nil
		}
		queues = append(queues, r)
	}

	return queues, nil
}

func (d Db) UpdateQueue(id int, status int, errors interface{}) {
	tx, err := d.Db.Begin()

	if err != nil {
		logrus.Fatal(err)
	}

	defer tx.Rollback()

	stmt, err := tx.Prepare("UPDATE worker_queues SET status = ?, error = ? WHERE id = ?")

	if err != nil {
		logrus.Fatal(err)
	}

	defer stmt.Close()

	_, err = stmt.Exec(status, errors, id)
	if err != nil {
		logrus.Fatal(err)
	}

	err = tx.Commit()
	if err != nil {
		logrus.Fatal(err)
	}
}

func (d Db) GetTypes(table string) []string {
	cols, err := d.Db.Query(fmt.Sprintf("SHOW COLUMNS FROM %s", table))

	if err != nil {
		logrus.Fatalf(err.Error())
	}
	types := []models.Types{}
	for cols.Next() {
		var t models.Types
		err := cols.Scan(&t.Field, &t.Type, &t.Null, &t.Key, &t.Default, &t.Extra)
		if err != nil {
			logrus.Fatalf(err.Error())
		}
		types = append(types, t)

	}
	var typesArr []string
	for _, val := range types {
		typesArr = append(typesArr, val.Type)
	}

	return typesArr
}

func (d Db) CheckRow(table string, id string) bool {
	var row models.Row
	d.Db.QueryRow(fmt.Sprintf("SELECT id FROM %v WHERE id = ?", table), id).Scan(&row.Id)

	return (row.Id != "")
}

func (d Db) Insert(table string, row string, queueId int) error {
	val, err := json.Marshal(row)

	if err != nil {
		return err
	}

	query := fmt.Sprintf("INSERT INTO %s (query, queueId) VALUES (%s, %v)", table, string(val), queueId)

	_, err = d.Db.Exec(query)

	if err != nil {
		logrus.Printf("Error %s when preparing SQL statement", err)
		return err
	}

	return nil
}

func (d Db) GetImportable(queueId int) ([]models.ImportQueue, error) {
	sql := "SELECT * FROM import_queues WHERE queueId = ?"
	queues := []models.ImportQueue{}
	rows, err := d.Db.Query(sql, queueId)

	if err != nil {
		logrus.Fatalf(err.Error())
		return []models.ImportQueue{}, nil
	}

	for rows.Next() {
		var r models.ImportQueue
		err := rows.Scan(&r.Id, &r.Query, &r.QueueId)
		if err != nil {
			logrus.Fatalf(err.Error())
			return []models.ImportQueue{}, nil
		}
		queues = append(queues, r)
	}

	return queues, nil
}

func (d Db) Exec(query string) error {
	_, err := d.Db.Exec(query)

	if err != nil {
		logrus.Printf("Error %s when preparing SQL statement", err)
		return err
	}
	return nil
}

func Explode(delimiter, text string) []string {
	if len(delimiter) > len(text) {
		return strings.Split(delimiter, text)
	} else {
		return strings.Split(text, delimiter)
	}
}

func (d Db) GetTableData(table string, cols string) ([]interface{}, []string, error) {
	rows, err := d.Db.Query(fmt.Sprintf("SELECT %s FROM %s", cols, table))
	columns, _ := rows.Columns()
	if err != nil {
		return []interface{}{}, []string{}, err
	}

	var allMaps []interface{}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		pointers := make([]interface{}, len(columns))
		for i, _ := range values {
			pointers[i] = &values[i]
		}

		err := rows.Scan(pointers...)
		if err != nil {
			logrus.Fatalf(err.Error())
		}

		resultMap := make(map[string]interface{})
		for i, val := range values {
			resultMap[columns[i]] = fmt.Sprintf("%s", val)

		}
		allMaps = append(allMaps, resultMap)
	}

	return allMaps, columns, nil
}
