package middleware

import (
	"context"
	"net/http"
	"web-forum/www/handlers"
	"web-forum/www/templates"
)

func Page(uri string, title string, newFunc func(*http.Request)) {
	http.HandleFunc(uri, func(writer http.ResponseWriter, reader *http.Request) {
		if uri == "/" && reader.URL.Path != "/" {
			http.NotFound(writer, reader)
			return
		}

		infoToSend, accountData := handlers.Base(reader)
		infoToSend["Title"] = title

		ctx := context.WithValue(reader.Context(), "InfoToSend", infoToSend)
		ctx = context.WithValue(ctx, "AccountData", accountData)

		reader = reader.WithContext(ctx)

		newFunc(reader)

		templates.Index.Execute(writer, infoToSend)
	})

	// templates.ContentAdd(infoToSend, templates.FAQ, nil)
}