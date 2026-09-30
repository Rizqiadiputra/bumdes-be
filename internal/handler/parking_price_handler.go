package handler

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type ParkingPriceHandler struct {
	parkingPriceService service.ParkingPriceService
	userLogService      service.UserLogService
}

func NewParkingPriceHandler(parkingPriceService service.ParkingPriceService, userLogService service.UserLogService) *ParkingPriceHandler {
	return &ParkingPriceHandler{parkingPriceService: parkingPriceService, userLogService: userLogService}
}

// ListParkingPrices godoc
// @Summary List Tarif Parkir
// @Description Menampilkan daftar tarif parkir, bisa difilter dengan pencarian nama
// @Tags Parking Prices
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan nama"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.ParkingPriceResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/parking-prices [get]
func (h *ParkingPriceHandler) ListParkingPrices(c *gin.Context) {
	filter := repository.ParkingPriceFilter{Search: c.Query("search")}

	result, err := h.parkingPriceService.ListParkingPrices(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tarif parkir", result)
}

// GetParkingPrice godoc
// @Summary Detail Tarif Parkir
// @Description Menampilkan detail satu tarif parkir
// @Tags Parking Prices
// @Produce json
// @Security BearerAuth
// @Param id path int true "Parking Price ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.ParkingPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/parking-prices/{id} [get]
func (h *ParkingPriceHandler) GetParkingPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif parkir tidak valid")
		return
	}

	result, err := h.parkingPriceService.GetParkingPrice(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "tarif parkir tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail tarif parkir", result)
}

// CreateParkingPrice godoc
// @Summary Buat Tarif Parkir
// @Description Membuat tarif parkir baru
// @Tags Parking Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateParkingPriceRequest true "Data tarif parkir baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.ParkingPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/parking-prices [post]
func (h *ParkingPriceHandler) CreateParkingPrice(c *gin.Context) {
	var req dto.CreateParkingPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.parkingPriceService.CreateParkingPrice(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "parking-prices",
			fmt.Sprintf("membuat tarif parkir %s", result.Name), req)
	}

	utils.Success(c, http.StatusCreated, "tarif parkir berhasil dibuat", result)
}

// UpdateParkingPrice godoc
// @Summary Ubah Tarif Parkir
// @Description Mengubah data tarif parkir
// @Tags Parking Prices
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Parking Price ID"
// @Param request body dto.UpdateParkingPriceRequest true "Data tarif parkir yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.ParkingPriceResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/parking-prices/{id} [put]
func (h *ParkingPriceHandler) UpdateParkingPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif parkir tidak valid")
		return
	}

	var req dto.UpdateParkingPriceRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.parkingPriceService.UpdateParkingPrice(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif parkir tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "parking-prices",
			fmt.Sprintf("mengubah tarif parkir %s", newResp.Name), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "tarif parkir berhasil diperbarui", newResp)
}

// DeleteParkingPrice godoc
// @Summary Hapus Tarif Parkir
// @Description Menghapus tarif parkir
// @Tags Parking Prices
// @Produce json
// @Security BearerAuth
// @Param id path int true "Parking Price ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/parking-prices/{id} [delete]
func (h *ParkingPriceHandler) DeleteParkingPrice(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tarif parkir tidak valid")
		return
	}

	if err := h.parkingPriceService.DeleteParkingPrice(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tarif parkir tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "parking-prices",
			fmt.Sprintf("menghapus tarif parkir id %d", id))
	}

	utils.Success(c, http.StatusOK, "tarif parkir berhasil dihapus", nil)
}
