package middleware

import (
	"context"
	"net/http"

	"github.com/chaihaobo/gocommon/restkit"
	"github.com/gin-gonic/gin"
	"github.com/gobwas/glob"
	"github.com/samber/lo"

	"github.com/chaihaobo/be-template/application"
	"github.com/chaihaobo/be-template/constant"
	"github.com/chaihaobo/be-template/resource"
)

const (
	headerKeyAuthorization = "Authorization"
)

type (
	httpEndpoint struct {
		Method string
		Path   string
	}
)

func (h *httpEndpoint) Match(request *http.Request) bool {
	return request.Method == h.Method && glob.MustCompile(h.Path).Match(request.URL.Path)
}

var (
	authWhitelist = []*httpEndpoint{
		{Method: http.MethodPost, Path: "/user/login"},
		{Method: http.MethodGet, Path: "/health"},
		{Method: http.MethodGet, Path: "/*/open/*"},
	}
)

func AuthMiddleware(res resource.Resource, app application.Application) gin.HandlerFunc {
	return func(gctx *gin.Context) {
		if isWhiteListRequest(gctx.Request) {
			gctx.Next()
			return
		}
		token := gctx.GetHeader(headerKeyAuthorization)
		if token == "" {
			restkit.HTTPWriteErr(gctx.Writer, constant.ErrUnauthorized)
			gctx.Abort()
			return
		}

		userClaims, err := app.User().TokenManger().Verify(token)
		if err != nil {
			res.Logger().Error(gctx, "verify jwt token failed", err)
			restkit.HTTPWriteErr(gctx.Writer, constant.ErrUnauthorized)
			gctx.Abort()
		}
		ctx := context.WithValue(gctx.Request.Context(), constant.ContextKeyUserID, userClaims.ID)
		gctx.Request = gctx.Request.WithContext(ctx)
		gctx.Next()
	}
}

func isWhiteListRequest(request *http.Request) bool {
	if _, ok := lo.Find(authWhitelist, func(endpoint *httpEndpoint) bool {
		return endpoint.Match(request)
	}); ok {
		return true
	}
	return false
}
