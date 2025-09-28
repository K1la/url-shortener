package handler

import (
	"fmt"
	"github.com/K1la/url-shortener/internal/api/response"
	"github.com/K1la/url-shortener/internal/model"
	"github.com/wb-go/wbf/ginext"
	"github.com/wb-go/wbf/zlog"
	"time"
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
		URL:       req.URL,
		ShortURL:  req.UserShortURL,
		CreatedAt: time.Now(),
	}

	res, err := h.service.CreateShortURL(c.Request.Context(), url)
	if err != nil {
		zlog.Logger.Error().Err(err).Msg("create short url error")
		response.BadRequest(c.Writer, fmt.Errorf("create short url error: %s", err.Error()))
		return
	}

	zlog.Logger.Debug().Msgf("create short url success: %+v", res)

}
