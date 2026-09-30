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

type TenantTypeHandler struct {
	tenantTypeService service.TenantTypeService
	userLogService    service.UserLogService
}

func NewTenantTypeHandler(tenantTypeService service.TenantTypeService, userLogService service.UserLogService) *TenantTypeHandler {
	return &TenantTypeHandler{tenantTypeService: tenantTypeService, userLogService: userLogService}
}

// ListTenantTypes godoc
// @Summary List Tenant Type
// @Description Menampilkan daftar tipe tenant, bisa difilter dengan pencarian type
// @Tags Tenant Types
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan type"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.TenantTypeResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/tenant-types [get]
func (h *TenantTypeHandler) ListTenantTypes(c *gin.Context) {
	filter := repository.TenantTypeFilter{Search: c.Query("search")}

	result, err := h.tenantTypeService.ListTenantTypes(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar tipe tenant", result)
}

// GetTenantType godoc
// @Summary Detail Tenant Type
// @Description Menampilkan detail satu tipe tenant
// @Tags Tenant Types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant Type ID (UUID)"
// @Success 200 {object} utils.SuccessResponse{data=dto.TenantTypeResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/tenant-types/{id} [get]
func (h *TenantTypeHandler) GetTenantType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tipe tenant tidak valid")
		return
	}

	result, err := h.tenantTypeService.GetTenantType(id)
	if err != nil {
		utils.Error(c, http.StatusNotFound, "tipe tenant tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail tipe tenant", result)
}

// CreateTenantType godoc
// @Summary Buat Tenant Type
// @Description Membuat tipe tenant baru
// @Tags Tenant Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateTenantTypeRequest true "Data tipe tenant baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.TenantTypeResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/tenant-types [post]
func (h *TenantTypeHandler) CreateTenantType(c *gin.Context) {
	var req dto.CreateTenantTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.tenantTypeService.CreateTenantType(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "tenant-types",
			fmt.Sprintf("membuat tipe tenant %s", result.Type), req)
	}

	utils.Success(c, http.StatusCreated, "tipe tenant berhasil dibuat", result)
}

// UpdateTenantType godoc
// @Summary Ubah Tenant Type
// @Description Mengubah data tipe tenant
// @Tags Tenant Types
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant Type ID (UUID)"
// @Param request body dto.UpdateTenantTypeRequest true "Data tipe tenant yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.TenantTypeResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/tenant-types/{id} [put]
func (h *TenantTypeHandler) UpdateTenantType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tipe tenant tidak valid")
		return
	}

	var req dto.UpdateTenantTypeRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.tenantTypeService.UpdateTenantType(id, req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tipe tenant tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "tenant-types",
			fmt.Sprintf("mengubah tipe tenant %s", newResp.Type), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "tipe tenant berhasil diperbarui", newResp)
}

// DeleteTenantType godoc
// @Summary Hapus Tenant Type
// @Description Menghapus tipe tenant
// @Tags Tenant Types
// @Produce json
// @Security BearerAuth
// @Param id path string true "Tenant Type ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/tenant-types/{id} [delete]
func (h *TenantTypeHandler) DeleteTenantType(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id tipe tenant tidak valid")
		return
	}

	if err := h.tenantTypeService.DeleteTenantType(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "tipe tenant tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "tenant-types",
			fmt.Sprintf("menghapus tipe tenant id %s", id))
	}

	utils.Success(c, http.StatusOK, "tipe tenant berhasil dihapus", nil)
}
