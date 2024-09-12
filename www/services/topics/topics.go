package topics

import (
	"context"
	"fmt"
	"log"
	"math"
	"web-forum/internal"
	"web-forum/system"
	"web-forum/system/db"
	"web-forum/www/services/account"
)

var ctx = context.Background()

func Get(forumId int, page int) (*internal.Paginator, error) {
	const errorFunction = "topics.Get"
	var topics internal.Paginator

	tx, err := db.Postgres.Begin(ctx)
	defer tx.Commit(ctx)

	if err != nil {
		log.Fatal(fmt.Errorf("%s: %w", errorFunction, err))
	}

	var topicsCount float64
	queryRow := tx.QueryRow(ctx, "SELECT COUNT(*) FROM topics WHERE forum_id=$1;", forumId)
	countTopicsErr := queryRow.Scan(&topicsCount)
	pagesCount := math.Ceil(topicsCount / internal.MaxPaginatorTopics)

	if countTopicsErr != nil {
		log.Fatal(fmt.Errorf("%s: %w", errorFunction, countTopicsErr))
	}

	fmtQuery := fmt.Sprintf("SELECT * FROM topics WHERE forum_id = $1 ORDER BY id DESC LIMIT %d OFFSET %d;", internal.MaxPaginatorTopics, (page-1)*internal.MaxPaginatorTopics)
	rows, err := tx.Query(ctx, fmtQuery, forumId)
	defer rows.Close()

	if err != nil {
		log.Fatal("[functions:35]", err)
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

	topics.CurrentPage = page
	topics.AllPages = int(pagesCount)

	return &topics, nil
}
