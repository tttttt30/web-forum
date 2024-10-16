package topics_messages

import (
	"context"
	"fmt"
	"web-forum/internal/app/models"
	"web-forum/internal/app/services/account"
	"web-forum/internal/app/services/paginator"
	"web-forum/pkg/stuff"
)

var ctx = context.Background()

func Get(topic *models.Topic, page int) *models.Paginator {
	const errorFunction = "topics_messages.Get"

	var queryCount = fmt.Sprintf("select message_count from topics where id = %d", (*topic).Id)
	var preQuery = models.PaginatorPreQuery{
		TableName:     "messages",
		OutputColumns: "id, topic_id, account_id, message, create_time, update_time",
		WhereColumn:   "topic_id",
		WhereValue:    (*topic).Id,
		Page:          page,
		QueryCount: models.PaginatorQueryCount{
			Query: queryCount,
		},
	}
	var paginatorList = models.Paginator{}

	tx, rows, paginatorList, err := paginator.Query(preQuery)
	if tx != nil {
		defer tx.Commit(ctx)
	}

	if err != nil {
		stuff.ErrLog(errorFunction, err)
		return &paginatorList
	}

	var tempUsers []int
	var tempMessages []models.Message

	for rows.Next() {
		var msg models.Message
		var scanErr = rows.Scan(&msg.Id, &msg.TopicId, &msg.CreatorId, &msg.Message, &msg.CreateTime, &msg.UpdateTime)

		if scanErr != nil {
			stuff.ErrLog(errorFunction, scanErr)
			continue
		}

		tempUsers = append(tempUsers, msg.CreatorId)
		tempMessages = append(tempMessages, msg)
	}

	var usersInfo = account.GetFromSlice(tempUsers, tx)
	var parentId = (*topic).ParentId

	for i := 0; i < len(tempMessages); i++ {
		var msg = tempMessages[i]
		var acc, ok = usersInfo[msg.CreatorId]

		if !ok {
			stuff.ErrLog(errorFunction, fmt.Errorf("не найден креатор сообщения в бд? > %s(ID): %s(MSG)", msg.CreatorId, msg.TopicId))
			continue
		}

		messageInfo := map[string]interface{}{
			"MessageId":         msg.Id,
			"CreatorId":         msg.CreatorId,
			"Username":          acc.Username,
			"Message":           msg.Message,
			"CreateTime":        msg.CreateTime.Format("2006-01-02 15:04:05"),
			"IsParentedMessage": parentId == msg.Id,
		}

		if msg.UpdateTime.Valid {
			messageInfo["UpdateTime"] = msg.UpdateTime.Time.Format("2006-01-02 15:04:05")
		}

		if acc.Avatar.Valid {
			messageInfo["Avatar"] = acc.Avatar.String
		}

		if acc.SignText.Valid {
			messageInfo["SignText"] = acc.SignText.String
		}

		paginatorList.Objects = append(paginatorList.Objects, messageInfo)
	}

	return &paginatorList
}
