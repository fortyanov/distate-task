package webserver

import (
	"context"
	"distate-task/dt/utils"
	"encoding/json"
	"fmt"
	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
	"golang.org/x/sync/errgroup"
	"time"
)

type SetDocRequest struct {
	Documents []Document `json:"documents"`
}

type Document struct {
	Name string `json:"name"`
	Date   utils.TimeISO `json:"date"`
	Number int           `json:"number"`
	Sum    string        `json:"sum"`
}

func (s *WebServer) setDocHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		req := &SetDocRequest{Documents: make([]Document, 0, 10)}
		pBody := ctx.PostBody()

		err := json.Unmarshal(pBody, &req)
		if err != nil {
			s.Log.Info("Incorrect request body in setDocHandler", zap.String("request", string(pBody)))
			jsonError(ctx, "Incorrect body", 400)
			return
		}

		succesDocs := make(map[int]Document)
		wg, wgctx := errgroup.WithContext(ctx)
		for _, doc := range req.Documents {
			doc := doc

			docDate := &doc.Date
			if doc.Date == utils.TimeISO(*new(time.Time)) {
				docDate = nil
			}

			succesDocs[doc.Number] = doc
			wg.Go(func() error {
				tx, err := s.conn.Begin(ctx)
				if err != nil {
					s.Log.Error("transaction begin error", zap.Error(err))
					return err
				}
				defer tx.Rollback(context.Background())

				tx.Exec(wgctx,
					`INSERT INTO "distate"."document" (name, date, number, sum)
					VALUES ($1, $2, $3, $4)`,
					doc.Name, docDate, doc.Number, doc.Sum)
				err = tx.Commit(context.Background())
				if err != nil {
					s.Log.Error("transaction commit error", zap.Error(err))
					delete(succesDocs, doc.Number)
					return err
				}
				return nil
			})
		}
		err = wg.Wait()
		if err != nil {
			s.Log.Error("documents group set error", zap.Error(err))
		}

		jsonDocumentsBody(ctx, succesDocs)
	}
}

func (s *WebServer) getDocHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		fmt.Fprintf(ctx, "Hello, world!\n\n")

		fmt.Fprintf(ctx, "Request method is %q\n", ctx.Method())
		fmt.Fprintf(ctx, "RequestURI is %q\n", ctx.RequestURI())
		fmt.Fprintf(ctx, "Requested path is %q\n", ctx.Path())
		fmt.Fprintf(ctx, "Host is %q\n", ctx.Host())
		fmt.Fprintf(ctx, "Query string is %q\n", ctx.QueryArgs())
		fmt.Fprintf(ctx, "User-Agent is %q\n", ctx.UserAgent())
		fmt.Fprintf(ctx, "Connection has been established at %s\n", ctx.ConnTime())
		fmt.Fprintf(ctx, "Request has been started at %s\n", ctx.Time())
		fmt.Fprintf(ctx, "Serial request number for the current connection is %d\n", ctx.ConnRequestNum())
		fmt.Fprintf(ctx, "Your ip is %q\n\n", ctx.RemoteIP())

		fmt.Fprintf(ctx, "Raw request is:\n---CUT---\n%s\n---CUT---", &ctx.Request)

		ctx.SetContentType("text/plain; charset=utf8")

		// Set arbitrary headers
		ctx.Response.Header.Set("X-My-Header", "my-header-value")

		// Set cookies
		var c fasthttp.Cookie
		c.SetKey("cookie-name")
		c.SetValue("cookie-value")
		ctx.Response.Header.SetCookie(&c)
	}
}

type errorBody struct {
	Error string `json:"error"`
}

type documentsBody struct {
	Documents []Document `json:"documents"`
}

func jsonError(ctx *fasthttp.RequestCtx, message string, statusCode int) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(statusCode)

	body, _ := json.Marshal(&errorBody{Error: message})
	ctx.Response.SetBody(body)
}

func jsonDocumentsBody(ctx *fasthttp.RequestCtx, successDocs map[int]Document) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(201)

	docs := make([]Document, 0, len(successDocs))
	for _, doc := range successDocs {
		docs = append(docs, doc)
	}

	body, _ := json.Marshal(&documentsBody{Documents: docs})
	ctx.Response.SetBody(body)
}
