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

type RevenueCategoryHandler struct {
	revenueCategoryService service.RevenueCategoryService
	userLogService         service.UserLogService
}

func NewRevenueCategoryHandler(revenueCategoryService service.RevenueCategoryService, userLogService service.UserLogService) *RevenueCategoryHandler {
	return &RevenueCategoryHandler{revenueCategoryService: revenueCategoryService, userLogService: userLogService}
}

// ListRevenueCategories godoc
// @Summary List Revenue Categories
// @Description Menampilkan daftar kategori pendapatan, bisa difilter dengan pencarian category_name atau source
// @Tags Revenue Categories
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan category_name atau source"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.RevenueCategoryResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/revenue-categories [get]
func (h *RevenueCategoryHandler) ListRevenueCategories(c *gin.Context) {
	filter := repository.RevenueCategoryFilter{Search: c.Query("search")}

	result, err := h.revenueCategoryService.ListRevenueCategories(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar kategori pendapatan", result)
}

// GetRevenueCategory godoc
// @Summary Detail Revenue Category
// @Description Menampilkan detail satu kategori pendapatan
// @Tags Revenue Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Revenue Category ID"
// @Success 200 {object} utils.SuccessResponse{data=dto.RevenueCategoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /master-data/revenue-categories/{id} [get]
func (h *RevenueCategoryHandler) GetRevenueCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id kategori pendapatan tidak valid")
		return
	}

	result, err := h.revenueCategoryService.GetRevenueCategory(uint(id))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "kategori pendapatan tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil detail kategori pendapatan", result)
}

// CreateRevenueCategory godoc
// @Summary Buat Revenue Category
// @Description Membuat kategori pendapatan baru
// @Tags Revenue Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateRevenueCategoryRequest true "Data kategori pendapatan baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.RevenueCategoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/revenue-categories [post]
func (h *RevenueCategoryHandler) CreateRevenueCategory(c *gin.Context) {
	var req dto.CreateRevenueCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.revenueCategoryService.CreateRevenueCategory(req)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "revenue-categories",
			fmt.Sprintf("membuat kategori pendapatan %s", result.CategoryName), req)
	}

	utils.Success(c, http.StatusCreated, "kategori pendapatan berhasil dibuat", result)
}

// UpdateRevenueCategory godoc
// @Summary Ubah Revenue Category
// @Description Mengubah data kategori pendapatan
// @Tags Revenue Categories
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path int true "Revenue Category ID"
// @Param request body dto.UpdateRevenueCategoryRequest true "Data kategori pendapatan yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.RevenueCategoryResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/revenue-categories/{id} [put]
func (h *RevenueCategoryHandler) UpdateRevenueCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id kategori pendapatan tidak valid")
		return
	}

	var req dto.UpdateRevenueCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.revenueCategoryService.UpdateRevenueCategory(uint(id), req)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "kategori pendapatan tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "revenue-categories",
			fmt.Sprintf("mengubah kategori pendapatan %s", newResp.CategoryName), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "kategori pendapatan berhasil diperbarui", newResp)
}

// DeleteRevenueCategory godoc
// @Summary Hapus Revenue Category
// @Description Menghapus kategori pendapatan
// @Tags Revenue Categories
// @Produce json
// @Security BearerAuth
// @Param id path int true "Revenue Category ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /master-data/revenue-categories/{id} [delete]
func (h *RevenueCategoryHandler) DeleteRevenueCategory(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id kategori pendapatan tidak valid")
		return
	}

	if err := h.revenueCategoryService.DeleteRevenueCategory(uint(id)); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "kategori pendapatan tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "revenue-categories",
			fmt.Sprintf("menghapus kategori pendapatan id %d", id))
	}

	utils.Success(c, http.StatusOK, "kategori pendapatan berhasil dihapus", nil)
}
