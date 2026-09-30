// Command billing-cron menghasilkan billing_histories untuk tenant yang
// tanggal jatuh temponya sudah masuk H-BILLING_DUE_REMINDER_DAYS hari.
// Dijalankan sekali per eksekusi, dijadwalkan lewat cron OS / scheduler
// (mis. `0 1 * * *` untuk jalan tiap hari jam 01:00), bukan sebagai daemon.
package main

import (
	"log"
	"time"

	"github.com/liyansasongko/bumdes-be/internal/config"
	"github.com/liyansasongko/bumdes-be/internal/database"
	"github.com/liyansasongko/bumdes-be/internal/repository"
	"github.com/liyansasongko/bumdes-be/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := database.Connect(cfg)
	if err != nil {
		log.Fatalf("gagal terhubung ke database: %v", err)
	}

	tenantRepo := repository.NewTenantRepository(db)
	billingRepo := repository.NewBillingHistoryRepository(db)
	billingService := service.NewBillingService(billingRepo, tenantRepo, cfg.BillingDueReminderDays)

	created, err := billingService.GenerateDueBillings(time.Now())
	if err != nil {
		log.Fatalf("gagal membuat billing: %v", err)
	}

	log.Printf("billing cron selesai: %d tagihan baru dibuat (H-%d)", len(created), cfg.BillingDueReminderDays)
}
