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

type UserHandler struct {
	userService    service.UserService
	userLogService service.UserLogService
}

func NewUserHandler(userService service.UserService, userLogService service.UserLogService) *UserHandler {
	return &UserHandler{userService: userService, userLogService: userLogService}
}

// Me godoc
// @Summary Profil Saya
// @Description Menampilkan data user yang sedang login beserta role dan permission-nya
// @Tags Me
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=dto.UserResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /me [get]
func (h *UserHandler) Me(c *gin.Context) {
	userID, ok := c.Get(middleware.ContextUserIDKey)
	if !ok {
		utils.Error(c, http.StatusUnauthorized, "token tidak valid")
		return
	}

	result, err := h.userService.GetMe(userID.(uuid.UUID))
	if err != nil {
		utils.Error(c, http.StatusNotFound, "user tidak ditemukan")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil profil", result)
}

// ListAccounts godoc
// @Summary List Akun
// @Description Menampilkan daftar seluruh akun/user beserta role masing-masing. Bisa difilter dengan pencarian nama/email dan nama role.
// @Tags Account Management
// @Produce json
// @Security BearerAuth
// @Param search query string false "Cari berdasarkan nama atau email"
// @Param name_role query string false "Filter berdasarkan nama role, isi 'all' untuk semua role"
// @Success 200 {object} utils.SuccessResponse{data=[]dto.UserResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /accounts [get]
func (h *UserHandler) ListAccounts(c *gin.Context) {
	filter := repository.UserFilter{
		Search:   c.Query("search"),
		NameRole: c.Query("name_role"),
	}

	result, err := h.userService.ListAccounts(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar akun", result)
}

// CreateAccount godoc
// @Summary Buat Akun
// @Description Membuat akun/user baru
// @Tags Account Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param request body dto.CreateAccountRequest true "Data akun baru"
// @Success 201 {object} utils.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /accounts [post]
func (h *UserHandler) CreateAccount(c *gin.Context) {
	var req dto.CreateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	result, err := h.userService.CreateAccount(req)
	if err != nil {
		switch {
		case errors.Is(err, service.ErrEmailAlreadyExists):
			utils.Error(c, http.StatusConflict, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogCreate(meta, actorID.(uuid.UUID), "accounts",
			fmt.Sprintf("membuat akun %s", result.Email), redactPassword(req))
	}

	utils.Success(c, http.StatusCreated, "akun berhasil dibuat", result)
}

// UpdateAccount godoc
// @Summary Ubah Akun
// @Description Mengubah data akun/user
// @Tags Account Management
// @Accept json
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Param request body dto.UpdateAccountRequest true "Data akun yang diubah"
// @Success 200 {object} utils.SuccessResponse{data=dto.UserResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 409 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /accounts/{id} [put]
func (h *UserHandler) UpdateAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id akun tidak valid")
		return
	}

	var req dto.UpdateAccountRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, http.StatusBadRequest, "data yang dikirim tidak valid")
		return
	}

	oldResp, newResp, err := h.userService.UpdateAccount(id, req)
	if err != nil {
		switch {
		case errors.Is(err, gorm.ErrRecordNotFound):
			utils.Error(c, http.StatusNotFound, "akun tidak ditemukan")
		case errors.Is(err, service.ErrEmailAlreadyExists):
			utils.Error(c, http.StatusConflict, err.Error())
		default:
			utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		}
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogUpdate(meta, actorID.(uuid.UUID), "accounts",
			fmt.Sprintf("mengubah akun %s", newResp.Email), oldResp, newResp)
	}

	utils.Success(c, http.StatusOK, "akun berhasil diperbarui", newResp)
}

// DeleteAccount godoc
// @Summary Hapus Akun
// @Description Menghapus akun/user
// @Tags Account Management
// @Produce json
// @Security BearerAuth
// @Param id path string true "User ID (UUID)"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /accounts/{id} [delete]
func (h *UserHandler) DeleteAccount(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		utils.Error(c, http.StatusBadRequest, "id akun tidak valid")
		return
	}

	if err := h.userService.DeleteAccount(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			utils.Error(c, http.StatusNotFound, "akun tidak ditemukan")
			return
		}
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	actorID, ok := c.Get(middleware.ContextUserIDKey)
	if ok {
		meta := utils.GetRequestMeta(c)
		_ = h.userLogService.LogDelete(meta, actorID.(uuid.UUID), "accounts",
			fmt.Sprintf("menghapus akun id %s", id))
	}

	utils.Success(c, http.StatusOK, "akun berhasil dihapus", nil)
}

func redactPassword(req dto.CreateAccountRequest) dto.CreateAccountRequest {
	if req.Password != "" {
		req.Password = "***"
	}
	return req
}
