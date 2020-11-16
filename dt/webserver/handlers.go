package webserver

import (
	"distate-task/dt/utils"

	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/valyala/fasthttp"
	"go.uber.org/zap"
)

type SetDocument struct {
	Name   string        `json:"name"`
	Date   utils.TimeISO `json:"date"`
	Number int           `json:"number"`
	Sum    string        `json:"sum"`
}

type SetDocRequest struct {
	Documents []SetDocument `json:"documents"`
}

type ResDocument struct {
	Id     int64     `json:"id"`
	Name   string    `json:"name"`
	Date   time.Time `json:"date"`
	Number int       `json:"number"`
	Sum    string    `json:"sum"`
}

func (s *WebServer) setDocHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		req := &SetDocRequest{Documents: make([]SetDocument, 0, 10)}
		pBody := ctx.PostBody()

		err := json.Unmarshal(pBody, &req)
		if err != nil {
			s.Log.Info("Incorrect request body in setDocHandler", zap.String("request", string(pBody)))
			jsonError(ctx, "Incorrect body", 400)
			return
		}

		succesDocs := make(map[int]SetDocument)
		var wg sync.WaitGroup
		wg.Add(len(req.Documents))
		for _, doc := range req.Documents {
			go func(doc SetDocument) {
				defer wg.Done()

				//docDate := &doc.Date
				//if doc.Date == utils.TimeISO(*new(time.Time)) {
				//	docDate = nil
				//}
				succesDocs[doc.Number] = doc

				tx, err := s.conn.Begin(ctx)
				if err != nil {
					s.Log.Error("transaction begin error", zap.Error(err))
					return
				}
				defer tx.Rollback(context.Background())

				_, err = tx.Exec(ctx,
					`INSERT INTO "distate"."document" (name, date, number, sum)
					VALUES ($1, $2, $3, $4)`,
					//doc.Name, docDate, doc.Number, doc.Sum)
					doc.Name, doc.Date, doc.Number, doc.Sum)
				if err != nil {
					s.Log.Error("exec error", zap.Error(err))
					delete(succesDocs, doc.Number)
					return
				}

				err = tx.Commit(context.Background())
				if err != nil {
					s.Log.Error("transaction commit error", zap.Error(err))
					delete(succesDocs, doc.Number)
					return
				}
			}(doc)
		}
		wg.Wait()

		jsonDocumentsBody(ctx, succesDocs)
	}
}

func (s *WebServer) getDocHandler() func(ctx *fasthttp.RequestCtx) {
	return func(ctx *fasthttp.RequestCtx) {
		ctx.Response.Header.SetContentType("application/json")
		ctx.Response.SetStatusCode(200)

		ctx.Response.SetBodyStreamWriter(func(w *bufio.Writer) {
			rows, err := s.conn.Query(ctx,
				`SELECT id,
						name,
						date,
						number,
						sum
			FROM "distate"."document"`)
			if err != nil {
				s.Log.Error("getDocHandler query error", zap.Error(err))
				jsonError(ctx, "Server not acceptable", 406)
				return
			}
			defer rows.Close()

			for rows.Next() {
				doc := new(ResDocument)
				err := rows.Scan(&doc.Id, &doc.Name, &doc.Date, &doc.Number, &doc.Sum)
				if err != nil {
					s.Log.Error("DB scan error", zap.Error(err))
					//jsonError(ctx, "Server not acceptable", 406)
					return
				}

				chunk, _ := json.Marshal(doc)
				fmt.Fprint(w, string(chunk))
				if err := w.Flush(); err != nil {
					return
				}
				//time.Sleep(time.Second)
			}
		})
	}
}
