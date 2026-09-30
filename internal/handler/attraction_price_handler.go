package handler

import (
	"errors"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type AttractionPriceHandler struct {
	attractionPriceService service.AttractionPriceService
	userLogService         service.UserLogService
}

func NewAttractionPriceHandler(attractionPriceService service.AttractionPriceService, userLogService service.UserLogService) *AttractionPriceHandler {
	return &AttractionPriceHandler{attractionPriceService: attractionPriceService, userLogService: userLogService}
}

// ListAttractionPrices godoc
// @Summary List Tarif Wahana
// @Description Menampilkan daftar tarif wahana/atraksi, bisa difilter dengan pencarian nama
// @Tags Attraction Prices
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan nama"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.AttractionPriceResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/attraction-prices [get]
func (h *AttractionPriceHandler) ListAttractionPrices(c *gin.Context) {
	filter := repository.AttractionPriceFilter{Search: c.Query("search")}

	result, err := h.attractionPriceService.ListAttractionPrices(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tarif wahana", result)
}

// GetAttractionPrice godoc
// @Summary Detail Tarif Wahana
// @Description Menampilkan detail satu tarif wahana/atraksi
// @Tags Attraction Prices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction Price ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.AttractionPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/attraction-prices/{id} [get]
func (h *AttractionPriceHandler) GetAttractionPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif wahana tidak valid")
		return
	}

	result, err := h.attractionPriceService.GetAttractionPrice(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "tarif wahana tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail tarif wahana", result)
}

// CreateAttractionPrice godoc
// @Summary Buat Tarif Wahana
// @Description Membuat tarif wahana/atraksi baru
// @Tags Attraction Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAttractionPriceRequest true "Data tarif wahana baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.AttractionPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/attraction-prices [post]
func (h *AttractionPriceHandler) CreateAttractionPrice(c *gin.Context) {
	var req dto.CreateAttractionPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.attractionPriceService.CreateAttractionPrice(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "attraction-prices",
			fmt.Sprintf("membuat tarif wahana %s", result.Name), req)
	}

	utils.Success(c, http.StatusCreated, "tarif wahana berhasil dibuat", result)
}

// UpdateAttractionPrice godoc
// @Summary Ubah Tarif Wahana
// @Description Mengubah data tarif wahana/atraksi
// @Tags Attraction Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction Price ID (UUID)"
// @Param request body dto.UpdateAttractionPriceRequest true "Data tarif wahana yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.AttractionPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/attraction-prices/{id} [put]
func (h *AttractionPriceHandler) UpdateAttractionPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif wahana tidak valid")
		return
	}

	var req dto.UpdateAttractionPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.attractionPriceService.UpdateAttractionPrice(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif wahana tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "attraction-prices",
			fmt.Sprintf("mengubah tarif wahana %s", newResp.Name), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "tarif wahana berhasil diperbarui", newResp)
}

// DeleteAttractionPrice godoc
// @Summary Hapus Tarif Wahana
// @Description Menghapus tarif wahana/atraksi
// @Tags Attraction Prices
// @Produce json
// @Security BearerAuth
// @Param id path string true "Attraction Price ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/attraction-prices/{id} [delete]
func (h *AttractionPriceHandler) DeleteAttractionPrice(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif wahana tidak valid")
		return
	}

	if err := h.attractionPriceService.DeleteAttractionPrice(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif wahana tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "attraction-prices",
			fmt.Sprintf("menghapus tarif wahana id %s", id))
	}

	utils.Success(c, http.StatusOK, "tarif wahana berhasil dihapus", nil)
}
