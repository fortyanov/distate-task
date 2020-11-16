package webserver

import (
	"github.com/valyala/fasthttp"
)

func (s *WebServer) Recovery(next func(ctx *fasthttp.RequestCtx)) func(ctx *fasthttp.RequestCtx) {
	fn := func(ctx *fasthttp.RequestCtx) {
		defer func() {
			if rvr := recover(); rvr != nil {
				ctx.Error("recover", 500)
			}
		}()
		next(ctx)
	}
	return fn
}
