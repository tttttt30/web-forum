package templates

import (
	"bytes"
	"html/template"
	"reflect"
)

func ParseFiles(page string) (*template.Template, error) {
	return template.ParseFiles(page, "frontend/template/not-authorized.html")
}

func ContentAdd(infoToSend *map[string]interface{}, tmpl *template.Template, content any) {
	if reflect.ValueOf(content).Kind() == reflect.Map {
		for k, v := range *infoToSend {
			content.(map[string]interface{})[k] = v
		}
	}

	newBytesBuffer := new(bytes.Buffer)
	tmpl.Execute(newBytesBuffer, content)
	(*infoToSend)["Content"] = template.HTML(newBytesBuffer.String())
}

var IndexTemplate = template.Must(ParseFiles("frontend/template/index.html"))
var LoginTemplate = template.Must(ParseFiles("frontend/template/login.html"))
var RegisterTemplate = template.Must(ParseFiles("frontend/template/register.html"))
var ProfileSettingsTemplate = template.Must(ParseFiles("frontend/template/profile-settings.html"))
var ForumTemplate = template.Must(ParseFiles("frontend/template/forum.html"))
var TopicsTemplate = template.Must(ParseFiles("frontend/template/topics.html"))
var TopicTemplate = template.Must(ParseFiles("frontend/template/topic.html"))
var CreateNewTopicTemplate = template.Must(ParseFiles("frontend/template/create-new-topic.html"))
var FAQTemplate = template.Must(ParseFiles("frontend/template/faq.html"))
var UsersTemplate = template.Must(ParseFiles("frontend/template/users.html"))