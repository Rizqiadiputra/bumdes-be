package handler

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"gorm.io/gorm"

	"github.com/liyansasongko/bumdes-be/internal/dto"
	"github.com/liyansasongko/bumdes-be/internal/middleware"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

const parkingDateLayout = "2006-01-02"

type ParkingHandler struct {
	parkingService service.ParkingService
	userLogService service.UserLogService
}

func NewParkingHandler(parkingService service.ParkingService, userLogService service.UserLogService) *ParkingHandler {
	return &ParkingHandler{parkingService: parkingService, userLogService: userLogService}
}

// ListParkings godoc
// @Summary List Parkir
// @Description Menampilkan daftar transaksi parkir. Default menampilkan transaksi hari ini, bisa difilter dengan date dan lokasi.
// @Tags Parking
// @Produce json
// @Security BearerAuth
// @Param date query string false "Filter tanggal, format YYYY-MM-DD (default: hari ini)"
// @Param lokasi query string false "Filter berdasarkan location"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.ParkingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /parkings [get]
func (h *ParkingHandler) ListParkings(c *gin.Context) {
	filterDate := time.Now()
	if v := c.Query("date"); v != "" {
		parsed, err := time.Parse(parkingDateLayout, v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "format date tidak valid, gunakan YYYY-MM-DD")
			return
		}
		filterDate = parsed
	}

	filter := repository.ParkingFilter{
		Date:     filterDate,
		Location: c.Query("lokasi"),
	}

	result, err := h.parkingService.ListParkings(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar parkir", result)
}

// GetParking godoc
// @Summary Detail Parkir
// @Description Menampilkan detail satu transaksi parkir
// @Tags Parking
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parking ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.ParkingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /parkings/{id} [get]
func (h *ParkingHandler) GetParking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id parkir tidak valid")
		return
	}

	result, err := h.parkingService.GetParking(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "data parkir tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail parkir", result)
}

// CreateParking godoc
// @Summary Buat Transaksi Parkir
// @Description Membuat transaksi parkir baru. Field type dan price otomatis diambil dari tarif parkir (parking_price_id) yang dipilih, dikali jumlah kendaraan (count). Field user_id otomatis diisi dari user yang sedang login.
// @Tags Parking
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateParkingRequest true "Data transaksi parkir baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.ParkingResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /parkings [post]
func (h *ParkingHandler) CreateParking(c *gin.Context) {
	var req dto.CreateParkingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "token tidak valid")
		return
	}

	result, err := h.parkingService.CreateParking(req, actorID.(uuid.UUID))
	if err != nil {
		switch {
		case errors.Is(err, service.ErrParkingPriceNotFound), errors.Is(err, service.ErrPaymentMethodNotFound):
			utils.Error(c, http.StatusNotFound, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	meta := utils.GetRequestMeta(c)
	_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "parkings",
		fmt.Sprintf("membuat transaksi parkir %s di %s", result.Type, result.Location), req)

	utils.Success(c, http.StatusCreated, "transaksi parkir berhasil dibuat", result)
}

// DeleteParking godoc
// @Summary Hapus Transaksi Parkir
// @Description Menghapus transaksi parkir
// @Tags Parking
// @Produce json
// @Security BearerAuth
// @Param id path string true "Parking ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /parkings/{id} [delete]
func (h *ParkingHandler) DeleteParking(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id parkir tidak valid")
		return
	}

	if err := h.parkingService.DeleteParking(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "data parkir tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "parkings",
			fmt.Sprintf("menghapus transaksi parkir id %s", id))
	}

	utils.Success(c, http.StatusOK, "transaksi parkir berhasil dihapus", nil)
}
