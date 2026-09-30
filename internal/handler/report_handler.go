package handler

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
	"github.com/liyansasongko/bumdes-be/internal/utils"
)

const reportDateLayout = "2006-01-02"

var validRevenueUnits = map[string]bool{
	service.RevenueUnitParkir: true,
	service.RevenueUnitTiket:  true,
	service.RevenueUnitWahana: true,
}

type ReportHandler struct {
	reportService service.ReportService
}

func NewReportHandler(reportService service.ReportService) *ReportHandler {
	return &ReportHandler{reportService: reportService}
}

// GetRevenueReport godoc
// @Summary Laporan Pendapatan
// @Description Laporan pendapatan gabungan dari parkir, tiket, dan wahana pada rentang tanggal tertentu. Bisa difilter berdasarkan unit usaha dan jenis pendapatan.
// @Tags Reports
// @Produce json
// @Security BearerAuth
// @Param start_date query string false "Tanggal awal, format YYYY-MM-DD (default: 1 bulan sebelum end_date)"
// @Param end_date query string false "Tanggal akhir, format YYYY-MM-DD (default: hari ini)"
// @Param unit query string false "Filter unit usaha (Parkir, Tiket, Wahana)"
// @Param jenis query string false "Filter jenis pendapatan sesuai unit yang dipilih"
// @Success 200 {object} utils.SuccessResponse{data=dto.RevenueReportResponse}
// @Failure 400 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Failure 500 {object} utils.ErrorResponse
// @Router /reports/revenue [get]
func (h *ReportHandler) GetRevenueReport(c *gin.Context) {
	endDate := time.Now()
	if v := c.Query("end_date"); v != "" {
		parsed, err := time.Parse(reportDateLayout, v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "format end_date tidak valid, gunakan YYYY-MM-DD")
			return
		}
		endDate = parsed
	}

	startDate := endDate.AddDate(0, -1, 0)
	if v := c.Query("start_date"); v != "" {
		parsed, err := time.Parse(reportDateLayout, v)
		if err != nil {
			utils.Error(c, http.StatusBadRequest, "format start_date tidak valid, gunakan YYYY-MM-DD")
			return
		}
		startDate = parsed
	}

	if startDate.After(endDate) {
		utils.Error(c, http.StatusBadRequest, "start_date tidak boleh setelah end_date")
		return
	}

	unit := c.Query("unit")
	if unit != "" && !validRevenueUnits[unit] {
		utils.Error(c, http.StatusBadRequest, "unit tidak valid, gunakan Parkir, Tiket, atau Wahana")
		return
	}

	filter := repository.RevenueReportFilter{
		StartDate: startDate,
		EndDate:   endDate,
		Unit:      unit,
		Jenis:     c.Query("jenis"),
	}

	result, err := h.reportService.GetRevenueReport(filter)
	if err != nil {
		utils.Error(c, http.StatusInternalServerError, "terjadi kesalahan pada server")
		return
	}

	utils.Success(c, http.StatusOK, "berhasil mengambil laporan pendapatan", result)
}
