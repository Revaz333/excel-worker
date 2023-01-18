package helpers

import (
	"database/sql"
	"fmt"
	"worker/models"

	"github.com/sirupsen/logrus"
)

type Db struct {
	Db *sql.DB
}

func (d Db) GetQueues() ([]models.WorkerQueues, error) {
	sql := "SELECT * FROM worker_queues WHERE status = ?"
	queues := []models.WorkerQueues{}
	rows, err := d.Db.Query(sql, 0) //Table("worker_queues").Where("status =	?", 0).Select("*").Scan(&queues)

	if err != nil {
		logrus.Fatalf(err.Error())
		return []models.WorkerQueues{}, nil
	}

	for rows.Next() {
		var r models.WorkerQueues
		err := rows.Scan(&r.ID, &r.Driver, &r.FilePath, &r.TableName, &r.Error, &r.Status)
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

func (d Db) GetTypes(table string) []models.Types {
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

	for _, val := range types {
		fmt.Println(val.Type)
	}
	return types
}
