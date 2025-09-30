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
	shortUrl := c.Param("shorten")

	if shortUrl == "" {
		zlog.Logger.Warn().Msg("missing short url")
		response.BadRequest(c.Writer, fmt.Errorf("missing short_url"))
		return
	}

	var RedirectInfo model.RedirectClicks
	RedirectInfo.ShortURL = shortUrl
	zlog.Logger.Info().Str("shorturl", RedirectInfo.ShortURL).Msg("shorturl from request")

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

	go h.saveAnalytics(rUrl)

	http.Redirect(c.Writer, c.Request, url.URL, http.StatusFound)
}

func (h *Handler) GetAnalytics(c *gin.Context) {
	shortUrl := c.Param("shorten")
	zlog.Logger.Info().Str("short_url", shortUrl).Msg("short url from summary analytics")
	if shortUrl == "" {
		zlog.Logger.Warn().Msg("missing short url")
		response.BadRequest(c.Writer, fmt.Errorf("missing short_url"))
		return
	}

	summary, err := h.service.GetAnalyticsSummary(c.Request.Context(), shortUrl)
	if err != nil {
		zlog.Logger.Error().Err(err).Str("short url", shortUrl).Msg("failed to get link analytics")
		response.Internal(c.Writer, fmt.Errorf("failed to get link analytics summary: %w", err))
		return
	}

	response.OK(c.Writer, summary)
}
