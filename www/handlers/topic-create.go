package handlers

import (
	"fmt"
	"net/http"
	"web-forum/www/services/category"
	"web-forum/www/templates"
)

func TopicCreate(stdRequest *http.Request) {
	forums, err := category.GetAll()

	if err != nil {
		panic(err)
	}

	var categorys []interface{}
	currentCategory := stdRequest.FormValue("category")

	for _, output := range *forums {
		forumId := output.Id

		categorys = append(categorys, map[string]interface{}{
			"Id":         output.Id,
			"Name":       output.Name,
			"IsSelected": fmt.Sprint(forumId) == currentCategory,
		})
	}

	templates.ContentAdd(stdRequest, templates.CreateNewTopic, map[string]interface{}{
		"Categorys": categorys,
	})
}
