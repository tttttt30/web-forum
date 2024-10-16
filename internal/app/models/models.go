package models

import (
	"database/sql"
	"os"
	"time"
)

const LoginMinLength = 4
const LoginMaxLength = 32

const PasswordMinLength = 8
const PasswordMaxLength = 64

const EmailMinLength = 4
const EmailMaxLength = 64

const UsernameMinLength = 4
const UsernameMaxLength = 24

const AvatarsFilePath = "www/staticfiles/imgs/avatars/"
const AvatarsSize = 256.0
const MaxPaginatorMessages = 10
const HowMuchPagesWillBeVisibleInPaginator = 9 // Только нечётные числа!!!

var HmacSecret = []byte(os.Getenv("HMAC_SECRET"))

type Category struct {
	Id          int    `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	TopicsCount int
}

type Topic struct {
	Id           int
	ForumId      int
	Name         string
	Creator      int
	CreateTime   time.Time
	UpdateTime   sql.NullTime
	MessageCount int
	ParentId     int
}

type MessageCreate struct {
	TopicId int    `json:"topic_id"`
	Message string `json:"message"`
}

type TopicCreate struct {
	Name       string `json:"name"`
	Message    string `json:"message"`
	CategoryId int    `json:"category_id"`
}

type MessageDelete struct {
	Id int `json:"id"`
}

type DeleteObject struct {
	Id int `json:"id"`
}

type Message struct {
	Id         int
	TopicId    int
	CreatorId  int
	Message    string
	CreateTime time.Time
	UpdateTime sql.NullTime
}

type ProfileMessage struct {
	TopicId    int
	TopicName  string
	Message    string
	CreateTime string
}

type PaginatorQueryCount struct {
	PreparedValue int
	Query         string
}

type PaginatorPreQuery struct {
	TableName       string
	OutputColumns   string
	WhereColumn     string
	WhereOperator   string
	OrderReverse    bool
	WhereValue      any
	Page            int
	ColumnsToOutput string

	QueryCount PaginatorQueryCount
}

type PaginatorArrows struct {
	Activated bool
	WhichPage int
}

type Paginator struct {
	PagesArray  []int
	Objects     []interface{} // Здесь наши обрезанные объекты
	CurrentPage int           // Текущая страница
	AllPages    int           // Все страницы
	Err         error

	Left  PaginatorArrows
	Right PaginatorArrows
}

type CountStruct struct {
	Users    int
	Topics   int
	Messages int
}
