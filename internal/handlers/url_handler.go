package handlers

import (
	"net/http"

	"github.com/4Noyis/url-shortener/internal/dto"
	"github.com/4Noyis/url-shortener/internal/service"
	"github.com/gin-gonic/gin"
)

type URLHandler struct {
	urlService *service.URLService
}

func NewURLHandler(urlService *service.URLService) *URLHandler {
	return &URLHandler{
		urlService: urlService,
	}
}

func (h *URLHandler) ShortenURL(c *gin.Context) {
	var req dto.ShortenURLRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: err.Error(),
		})
		return
	}

	url, err := h.urlService.ShortenURLWithTTL(req.LongURL, req.TTLSeconds)
	if err != nil {
		if err.Error() == "URL already exists" {
			c.JSON(http.StatusConflict, dto.ErrorResponse{
				Error:   "url_exists",
				Message: "This URL has already been shortened",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to shorten URL",
		})
		return
	}

	response := dto.ShortenURLResponse{
		ShortURL:  url.ShortURL,
		LongURL:   url.LongURL,
		CreatedAt: url.CreatedAt,
		ExpiresAt: url.ExpiresAt,
	}

	c.JSON(http.StatusCreated, response)
}

func (h *URLHandler) RedirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	if shortURL == "" {
		c.JSON(http.StatusBadRequest, dto.ErrorResponse{
			Error:   "invalid_request",
			Message: "Short URL parameter is required",
		})
		return
	}

	longURL, err := h.urlService.RedirectURL(shortURL)
	if err != nil {
		if err.Error() == "failed to get URL: short URL not found" || 
		   err.Error() == "failed to get URL: short URL has expired" {
			c.JSON(http.StatusNotFound, dto.ErrorResponse{
				Error:   "url_not_found",
				Message: "Short URL not found or has expired",
			})
			return
		}

		c.JSON(http.StatusInternalServerError, dto.ErrorResponse{
			Error:   "internal_error",
			Message: "Failed to redirect URL",
		})
		return
	}

	c.Redirect(http.StatusMovedPermanently, longURL)
}