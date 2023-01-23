package controllers

import (
	"fmt"
	"worker/internal/helpers"

	"github.com/sirupsen/logrus"
)

type ImportController struct {
	Db    helpers.Db
	Array helpers.Arrays
	Exel  helpers.Exel
}

func (c ImportController) CheckQueues() (bool, error) {
	queues, err := c.Db.GetQueues(4, "import")

	if err != nil {
		logrus.Fatalf(err.Error())
		return false, nil
	}

	if len(queues) == 0 {
		return false, nil
	}
	for _, queue := range queues {
		err := c.ImportRows(queue.ID)

		if err != nil {
			c.Db.UpdateQueue(queue.ID, 3, err.Error())
			return false, err
		}
	}
	return true, nil
}

func (c ImportController) ImportRows(queueId int) error {
	importable, err := c.Db.GetImportable(queueId)

	if err != nil {
		return err
	}

	for _, query := range importable {
		err := c.Db.Exec(query.Query)

		if err != nil {
			return err
		}
	}

	c.Db.UpdateQueue(queueId, 5, nil)

	err = c.Db.Exec(fmt.Sprintf("DELETE FROM import_queues WHERE queueId = %v", queueId))
	if err != nil {
		return err
	}

	return nil
}
