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

type KiosLocationHandler struct {
	kiosLocationService service.KiosLocationService
	userLogService      service.UserLogService
}

func NewKiosLocationHandler(kiosLocationService service.KiosLocationService, userLogService service.UserLogService) *KiosLocationHandler {
	return &KiosLocationHandler{kiosLocationService: kiosLocationService, userLogService: userLogService}
}

// ListKiosLocations godoc
// @Summary List Lokasi Kios
// @Description Menampilkan daftar lokasi kios, bisa difilter dengan pencarian kios
// @Tags Kios Locations
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan kios"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.KiosLocationResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/kios_locations [get]
func (h *KiosLocationHandler) ListKiosLocations(c *gin.Context) {
	filter := repository.KiosLocationFilter{Search: c.Query("search")}

	result, err := h.kiosLocationService.ListKiosLocations(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar lokasi kios", result)
}

// GetKiosLocation godoc
// @Summary Detail Lokasi Kios
// @Description Menampilkan detail satu lokasi kios
// @Tags Kios Locations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kios Location ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.KiosLocationResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/kios_locations/{id} [get]
func (h *KiosLocationHandler) GetKiosLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id lokasi kios tidak valid")
		return
	}

	result, err := h.kiosLocationService.GetKiosLocation(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "lokasi kios tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail lokasi kios", result)
}

// CreateKiosLocation godoc
// @Summary Buat Lokasi Kios
// @Description Membuat lokasi kios baru
// @Tags Kios Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateKiosLocationRequest true "Data lokasi kios baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.KiosLocationResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/kios_locations [post]
func (h *KiosLocationHandler) CreateKiosLocation(c *gin.Context) {
	var req dto.CreateKiosLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.kiosLocationService.CreateKiosLocation(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "kios_locations",
			fmt.Sprintf("membuat lokasi kios %s", result.Kios), req)
	}

	utils.Success(c, http.StatusCreated, "lokasi kios berhasil dibuat", result)
}

// UpdateKiosLocation godoc
// @Summary Ubah Lokasi Kios
// @Description Mengubah data lokasi kios
// @Tags Kios Locations
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kios Location ID (UUID)"
// @Param request body dto.UpdateKiosLocationRequest true "Data lokasi kios yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.KiosLocationResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/kios_locations/{id} [put]
func (h *KiosLocationHandler) UpdateKiosLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id lokasi kios tidak valid")
		return
	}

	var req dto.UpdateKiosLocationRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.kiosLocationService.UpdateKiosLocation(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "lokasi kios tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "kios_locations",
			fmt.Sprintf("mengubah lokasi kios %s", newResp.Kios), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "lokasi kios berhasil diperbarui", newResp)
}

// DeleteKiosLocation godoc
// @Summary Hapus Lokasi Kios
// @Description Menghapus lokasi kios
// @Tags Kios Locations
// @Produce json
// @Security BearerAuth
// @Param id path string true "Kios Location ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/kios_locations/{id} [delete]
func (h *KiosLocationHandler) DeleteKiosLocation(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id lokasi kios tidak valid")
		return
	}

	if err := h.kiosLocationService.DeleteKiosLocation(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "lokasi kios tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "kios_locations",
			fmt.Sprintf("menghapus lokasi kios id %s", id))
	}

	utils.Success(c, http.StatusOK, "lokasi kios berhasil dihapus", nil)
}
