package main

import (
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/kemenpppa/siarsip-api/internal/config"
	"github.com/kemenpppa/siarsip-api/internal/database"
	"github.com/kemenpppa/siarsip-api/internal/database/seeders"
	"github.com/kemenpppa/siarsip-api/internal/handler"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"github.com/kemenpppa/siarsip-api/internal/repository"
	"github.com/kemenpppa/siarsip-api/internal/router"
	"github.com/kemenpppa/siarsip-api/internal/service"
	"github.com/kemenpppa/siarsip-api/pkg/helper"
)

func main() {
	// Parse command-line flags
	seedFlag := flag.Bool("seed", false, "Jalankan seeder untuk membuat akun default")
	migrateFlag := flag.Bool("migrate", false, "Jalankan migrasi database saja")
	flag.Parse()

	// Load konfigurasi
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("Gagal memuat konfigurasi: %v", err)
	}

	// Koneksi database
	db, err := database.NewMySQL(cfg)
	if err != nil {
		log.Fatalf("Gagal koneksi database: %v", err)
	}

	// Auto migrate
	if err := database.AutoMigrate(db); err != nil {
		log.Fatalf("Gagal migrasi database: %v", err)
	}
	fmt.Println("✅ Migrasi database selesai")

	// Mode: hanya migrasi
	if *migrateFlag {
		fmt.Println("Migrasi selesai. Keluar...")
		os.Exit(0)
	}

	// Mode: seed akun default + master data + data arsip
	if *seedFlag {
		fmt.Println("🌱 Menjalankan seeder...")
		if err := seeders.SeedKodeKlasifikasi(db); err != nil {
			log.Fatalf("Gagal menjalankan seeder kode klasifikasi: %v", err)
		}
		if err := seeders.SeedMasterPendukung(db); err != nil {
			log.Fatalf("Gagal menjalankan seeder master pendukung: %v", err)
		}
		if err := seeders.SeedDefaultUsers(db); err != nil {
			log.Fatalf("Gagal menjalankan seeder user: %v", err)
		}
		if err := seeders.SeedArsipDefault(db); err != nil {
			log.Fatalf("Gagal menjalankan seeder arsip: %v", err)
		}
		fmt.Println("✅ Seeder selesai. Keluar...")
		os.Exit(0)
	}

	// Auto seed jika diaktifkan via environment variable (berguna saat deploy ke cloud)
	if os.Getenv("AUTO_SEED") == "true" {
		fmt.Println("🌱 Menjalankan auto-seeder...")
		_ = seeders.SeedKodeKlasifikasi(db)
		_ = seeders.SeedMasterPendukung(db)
		_ = seeders.SeedDefaultUsers(db)
		_ = seeders.SeedArsipDefault(db)
		fmt.Println("✅ Auto-seeder selesai")
	}

	// Mode: normal (jalankan server)
	userRepo := repository.NewUserRepository(db)
	masterRepo := repository.NewMasterRepository(db)
	arsipRepo := repository.NewArsipRepository(db)
	opRepo := repository.NewOperasionalRepository(db)
	mfaRepo := repository.NewMfaRepository(db)
	hashHelper := helper.NewHashHelper()

	authMiddleware := middleware.NewAuthMiddleware(cfg, db)
	roleMiddleware := middleware.NewRoleMiddleware(db)
	auditMiddleware := middleware.AuditMiddleware(opRepo)

	jwtSvc := service.NewJWTService(cfg, db)
	mfaSvc := service.NewMfaService(cfg, mfaRepo, userRepo, authMiddleware)
	authSvc := service.NewAuthService(userRepo, jwtSvc, hashHelper, mfaSvc, cfg)
	userSvc := service.NewUserService(userRepo, hashHelper, mfaRepo)
	masterSvc := service.NewMasterService(masterRepo)
	arsipSvc := service.NewArsipService(arsipRepo, masterRepo)
	opSvc := service.NewOperasionalService(opRepo, arsipRepo, cfg.Storage)
	laporanSvc := service.NewLaporanService(db)

	authHandler := handler.NewAuthHandler(authSvc)
	userHandler := handler.NewUserHandler(userSvc)
	masterHandler := handler.NewMasterHandler(masterSvc)
	arsipHandler := handler.NewArsipHandler(arsipSvc)
	opHandler := handler.NewOperasionalHandler(opSvc)
	laporanHandler := handler.NewLaporanHandler(laporanSvc)
	roleHandler := handler.NewRoleHandler(db)
	mfaHandler := handler.NewMfaHandler(mfaSvc, userRepo)
	docsHandler := handler.NewDocsHandler()
	_ = handler.EnsureDefaultRoles(db)

	r := router.SetupRouter(
		authHandler,
		userHandler,
		arsipHandler,
		authMiddleware,
		roleMiddleware,
		cfg.CORS.AllowedOrigins,
		router.RouterOptions{
			MasterHandler:      masterHandler,
			OperasionalHandler: opHandler,
			LaporanHandler:     laporanHandler,
			RoleHandler:        roleHandler,
			MfaHandler:         mfaHandler,
			DocsHandler:        docsHandler,
			EnableSwagger:      cfg.EnableSwagger,
			AuditMiddleware:    auditMiddleware,
			LoginRateLimit:     cfg.RateLimit.Login,
			MFARateLimit:       cfg.MFA.RateLimitVerify,
			APIRateLimit:       cfg.RateLimit.API,
			RateWindowMinutes:  cfg.RateLimit.WindowMinutes,
		},
	)

	port := cfg.AppPort
	if port == "" {
		port = "8000"
	}

	fmt.Printf("🚀 SiARSIP API Server berjalan di port %s\n", port)
	fmt.Printf("📋 CORS allowed origins: %s\n", cfg.CORS.AllowedOrigins)

	if err := r.Run(":" + port); err != nil {
		log.Fatalf("Gagal menjalankan server: %v", err)
		os.Exit(1)
	}
}
