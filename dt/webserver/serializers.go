package webserver

import (
	"encoding/json"

	"github.com/valyala/fasthttp"
)

type errorBody struct {
	Error string `json:"error"`
}

type documentsBody struct {
	Documents []SetDocument `json:"documents"`
}

func jsonError(ctx *fasthttp.RequestCtx, message string, statusCode int) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(statusCode)

	body, _ := json.Marshal(&errorBody{Error: message})
	ctx.Response.SetBody(body)
}

func jsonDocumentsBody(ctx *fasthttp.RequestCtx, successDocs map[int]SetDocument) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(201)

	docs := make([]SetDocument, 0, len(successDocs))
	for _, doc := range successDocs {
		docs = append(docs, doc)
	}

	body, _ := json.Marshal(&documentsBody{Documents: docs})
	ctx.Response.SetBody(body)
}
