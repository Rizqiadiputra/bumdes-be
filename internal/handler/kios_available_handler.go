package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

type KiosAvailableHandler struct {
	kiosAvailableService service.KiosAvailableService
}

func NewKiosAvailableHandler(kiosAvailableService service.KiosAvailableService) *KiosAvailableHandler {
	return &KiosAvailableHandler{kiosAvailableService: kiosAvailableService}
}

// ListKiosAvailables godoc
// @Summary List Ketersediaan Kios
// @Description Menampilkan daftar kios beserta status ketersediaan. Status Terisi memakai price dari data tenant yang menyewa, status Tersedia memakai price dari kios_locations.
// @Tags Kios Available
// @Produce json
// @Security BearerAuth
// @Success 200 {object} utils.SuccessResponse{data=[]dto.KiosAvailableResponse}
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /kios_availables [get]
func (h *KiosAvailableHandler) ListKiosAvailables(c *gin.Context) {
	result, err := h.kiosAvailableService.ListKiosAvailables()
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil daftar ketersediaan kios", result)
}
