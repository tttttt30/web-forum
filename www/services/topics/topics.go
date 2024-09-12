package topics

import (
	"context"
	"web-forum/internal"
	"web-forum/system"
	"web-forum/www/services/account"
	"web-forum/www/services/paginator"
)

var ctx = context.Background()

func Get(forumId int, page int) (*internal.Paginator, error) {
	const errorFunction = "topics.Get"

	tx, rows, topics, err := paginator.Query("topics", "forum_id", forumId, page)
	defer tx.Commit(ctx)

	if err != nil {
		return nil, system.ErrLog(errorFunction, err.Error())
	}

	var tempUsers []int
	var tempTopics []internal.Topic

	for rows.Next() {
		topic := internal.Topic{}

		scanErr := rows.Scan(&topic.Id, &topic.ForumId, &topic.Name, &topic.Creator, &topic.CreateTime, &topic.UpdateTime, &topic.MessageCount)

		if scanErr != nil {
			system.ErrLog(errorFunction, scanErr.Error())
			continue
		}

		tempUsers = append(tempUsers, topic.Creator)
		tempTopics = append(tempTopics, topic)
	}

	usersInfo := account.GetFromSlice(tempUsers, tx)

	for i := 0; i < len(tempTopics); i++ {
		topic := tempTopics[i]

		updateTime := ""

		if topic.UpdateTime.Valid {
			updateTime = topic.UpdateTime.Time.Format("2006-01-02 15:04:05")
		}

		creatorAccount, ok := usersInfo[topic.Creator]

		if !ok {
			system.ErrLog("topics.Get", "Не найден креатор топика в БД?")
			continue
		}

		aboutTopic := map[string]interface{}{
			"topic_id":    topic.Id,
			"forum_id":    forumId,
			"topic_name":  topic.Name,
			"username":    creatorAccount.Username,
			"create_time": topic.CreateTime.Format("2006-01-02 15:04:05"),
			"update_time": updateTime,

			// Поскольку 1 сообщение - это сообщение самого топика.
			"message_count": topic.MessageCount,
		}

		if creatorAccount.Avatar.Valid {
			aboutTopic["avatar"] = creatorAccount.Avatar.String
		}

		topics.Objects = append(topics.Objects, aboutTopic)
	}

	return &topics, nil
}
