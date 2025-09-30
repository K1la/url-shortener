package handler

import (
	"context"
	"fmt"
	"net"

	"github.com/K1la/url-shortener/internal/api/response"
	"github.com/K1la/url-shortener/internal/model"
	"github.com/gin-gonic/gin"
	"github.com/mssola/user_agent"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
)

func (h *Handler) CreateURLShort(c *ginext.Context) {
	var req CreateRequest

	if err := c.BindJSON(&req); err != nil {
		zlog.Logger.Error().Err(err).Msg("failed to decode request body")
		response.BadRequest(c.Writer, fmt.Errorf("decode req body error: %s", err.Error()))
		return
	}

	zlog.Logger.Debug().Msgf("json decode req: %+v", req)

	if err := h.valid.Struct(&req); err != nil {
		zlog.Logger.Warn().Err(err).Msg("failed to validate request body")
		response.BadRequest(c.Writer, fmt.Errorf("validation error: %s", err.Error()))
		return
	}
	zlog.Logger.Debug().Msgf("req after h.valid: %+v", req)

	url := model.URL{
		URL:      req.URL,
		ShortURL: req.UserShortURL,
	}

	res, err := h.service.CreateShortURL(c.Request.Context(), url)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("create short url error")
		response.BadRequest(c.Writer, fmt.Errorf("create short url error: %s", err.Error()))
		return
	}

	zlog.Logger.Debug().Msgf("create short url success: %+v", res)

	response.Created(c.Writer, res)
}

func (h *Handler) createAnalytics(c *gin.Context, rUrl *model.RedirectClicks) *model.RedirectClicks {
	uaString := c.Request.UserAgent()
	ip := c.ClientIP()

	ua := user_agent.New(uaString)
	browser, _ := ua.Browser()

	device := "desktop"
	if ua.Mobile() {
		device = "mobile"
	} else if ua.Bot() {
		device = "bot"
	}
	ip, _, err := net.SplitHostPort(c.Request.RemoteAddr)
	if err != nil {
		ip = c.Request.RemoteAddr
	}

	return &model.RedirectClicks{
		ShortURL:  rUrl.ShortURL,
		UserAgent: uaString,
		Device:    device,
		OS:        ua.OS(),
		Browser:   browser,
		IP:        ip,
	}
}

func (h *Handler) saveAnalytics(rUrl *model.RedirectClicks) {
	ctx := context.Background()

	id, err := h.service.SaveAnalytics(ctx, rUrl)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("shortUrl", rUrl.ShortURL).Msg("save analytics failed")
		return
	}

	if err := h.service.InvalidateAnalyticsCache(ctx, rUrl.ShortURL); err != nil {
		zlog.Logger.Error().Err(err).Str("shortUrl", rUrl.ShortURL).Msg("failed to invalidate analytics cache")
	}

	zlog.Logger.Info().Str("id", id).Interface("rUrl", rUrl).Msg("save analytics success")
}
