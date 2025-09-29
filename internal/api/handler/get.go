package handler

import (
	"errors"
	"fmt"
	"github.com/K1la/url-shortener/internal/api/response"
	"github.com/K1la/url-shortener/internal/model"
	"github.com/K1la/url-shortener/internal/repository"
	"github.com/gin-gonic/gin"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"net/http"
)

func (h *Handler) GetShortURL(c *ginext.Context) {
	shortUrl := c.Param("short_url")
	var RedirectInfo model.RedirectClicks
	RedirectInfo.ShortURL = shortUrl
	RedirectInfo.UserAgent = c.Request.Header.Get("User-Agent")
	zlog.Logger.Info().Str("header.Get: %v", RedirectInfo.UserAgent).Msg("useragent Header get")
	zlog.Logger.Info().Str("userAgent(): %v", c.Request.UserAgent()).Msg("useragent Header get")

	url, err := h.service.GetShortURL(c.Request.Context(), RedirectInfo)
	if err != nil {
		if errors.Is(err, repository.ErrShortURLNotFound) {
			zlog.Logger.Error().Err(err).Str("short_url", shortUrl).Msg("short url not found")
			response.BadRequest(c.Writer, fmt.Errorf("get short url from db: %w", err))
			return
		}

		zlog.Logger.Error().Err(err).Str("short_url", shortUrl).Msg("get short url")
		response.Internal(c.Writer, fmt.Errorf("get short url: %w", err))
		return
	}

	rUrl := h.createAnalytics(c, &RedirectInfo)
	zlog.Logger.Info().Interface("new redirect url info", rUrl).Msg("got new redirect url")

	go h.saveAnalytics(c, rUrl)

	http.Redirect(c.Writer, c.Request, url.URL, http.StatusFound)
}

func (h *Handler) saveAnalytics(c *gin.Context, rUrl *model.RedirectClicks) {
	id, err := h.service.SaveAnalytics(c.Request.Context(), rUrl)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("shortUrl", rUrl.ShortURL).Msg("save analytics failed")
		return
	}

	zlog.Logger.Info().Str("id", id).Interface("rUrl", rUrl).Msg("save analytics success")
}
