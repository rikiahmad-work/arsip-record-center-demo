package router

import (
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/kemenpppa/siarsip-api/internal/handler"
	"github.com/kemenpppa/siarsip-api/internal/middleware"
	"github.com/kemenpppa/siarsip-api/internal/model"
)

// RouterOptions mengelompokkan handler & konfigurasi tambahan untuk SetupRouter.
type RouterOptions struct {
	MasterHandler      *handler.MasterHandler
	OperasionalHandler *handler.OperasionalHandler
	LaporanHandler     *handler.LaporanHandler
	RoleHandler        *handler.RoleHandler
	MfaHandler         *handler.MfaHandler
	DocsHandler        *handler.DocsHandler
	EnableSwagger      bool
	AuditMiddleware    gin.HandlerFunc
	LoginRateLimit     int
	MFARateLimit       int
	APIRateLimit       int
	RateWindowMinutes  int
}

// SetupRouter mengatur seluruh route API sesuai matriks otorisasi 4 role (PRD Backend §4 & §7.5).
func SetupRouter(
	authHandler *handler.AuthHandler,
	userHandler *handler.UserHandler,
	arsipHandler *handler.ArsipHandler,
	authMiddleware *middleware.AuthMiddleware,
	roleMiddleware *middleware.RoleMiddleware,
	allowedOrigins string,
	opts ...RouterOptions,
) *gin.Engine {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(middleware.ErrorHandler())

	var opt RouterOptions
	if len(opts) > 0 {
		opt = opts[0]
	}
	if opt.LoginRateLimit <= 0 {
		opt.LoginRateLimit = 5
	}
	if opt.MFARateLimit <= 0 {
		opt.MFARateLimit = 5
	}
	if opt.APIRateLimit <= 0 {
		opt.APIRateLimit = 60
	}
	if opt.RateWindowMinutes <= 0 {
		opt.RateWindowMinutes = 1
	}

	// Parse comma-separated CORS origins dari environment
	origins := make([]string, 0)
	for _, o := range strings.Split(allowedOrigins, ",") {
		if trimmed := strings.TrimSpace(o); trimmed != "" {
			origins = append(origins, trimmed)
		}
	}
	if len(origins) == 0 {
		origins = []string{"http://localhost:3000"}
	}

	corsConfig := cors.Config{
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		AllowCredentials: true,
	}
	if allowedOrigins == "*" {
		corsConfig.AllowOriginFunc = func(origin string) bool {
			return true
		}
	} else {
		corsConfig.AllowOrigins = origins
	}
	r.Use(cors.New(corsConfig))

	if opt.AuditMiddleware != nil {
		r.Use(opt.AuditMiddleware)
	}

	window := time.Duration(opt.RateWindowMinutes) * time.Minute
	loginLimiter := middleware.NewRateLimiter(opt.LoginRateLimit, window)
	mfaLimiter := middleware.NewRateLimiter(opt.MFARateLimit, window)
	apiLimiter := middleware.NewRateLimiter(opt.APIRateLimit, window)

	// OpenAPI / Swagger UI (hanya didaftarkan jika EnableSwagger aktif)
	if opt.EnableSwagger {
		docsHandler := opt.DocsHandler
		if docsHandler == nil {
			docsHandler = handler.NewDocsHandler()
		}
		r.GET("/swagger", docsHandler.ServeSwaggerUI)
		r.GET("/api/docs", docsHandler.ServeSwaggerUI)
		r.GET("/api/docs/openapi.json", docsHandler.ServeOpenAPISpec)
	}

	api := r.Group("/api")
	{
		// === Public Routes ===
		api.POST("/auth/login", loginLimiter.Middleware(), authHandler.Login)
		if opt.MfaHandler != nil {
			api.POST("/auth/mfa/verify", mfaLimiter.Middleware(), opt.MfaHandler.Verify)
		}

		// === Protected Routes (perlu JWT token + general API rate limit) ===
		protected := api.Group("")
		protected.Use(apiLimiter.Middleware())
		protected.Use(authMiddleware.RequireAuth())
		{
			// 4.1 Auth
			protected.POST("/auth/logout", authHandler.Logout)
			protected.GET("/auth/profile", authHandler.Me)
			protected.GET("/auth/me", authHandler.Me)
			protected.PUT("/auth/change-password", authHandler.ChangePassword)
			if opt.MfaHandler != nil {
				protected.GET("/auth/mfa/status", opt.MfaHandler.Status)
				protected.POST("/auth/mfa/recovery-codes", mfaLimiter.Middleware(), opt.MfaHandler.ResetRecoveryCodes)
				protected.POST("/auth/mfa/reset", mfaLimiter.Middleware(), opt.MfaHandler.Reset)
			}

			// 4.2 User Management — hanya Administrator
			users := protected.Group("/users")
			users.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
			{
				users.GET("", userHandler.Index)
				users.GET("/stats", userHandler.Stats)
				users.GET("/:id", userHandler.Show)
				users.POST("", userHandler.Store)
				users.PUT("/:id", userHandler.Update)
				users.DELETE("/:id", userHandler.Destroy)
				users.PUT("/:id/restore", userHandler.Restore)
				if opt.MfaHandler != nil {
					users.POST("/:id/mfa/reset", mfaLimiter.Middleware(), opt.MfaHandler.ResetUser)
				}
			}

			// 4.2b Dynamic RBAC Roles & Permissions (/api/roles)
			if opt.RoleHandler != nil {
				roles := protected.Group("/roles")
				{
					roles.GET("", opt.RoleHandler.Index)

					rolesAdmin := roles.Group("")
					rolesAdmin.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
					{
						rolesAdmin.POST("", opt.RoleHandler.Store)
						rolesAdmin.PUT("/:id/permissions", opt.RoleHandler.UpdatePermissions)
						rolesAdmin.DELETE("/:id", opt.RoleHandler.Destroy)
					}
				}
			}

			// 4.3 Arsip — Read untuk role berizin archives.view; Write untuk Administrator & Operator; Delete/Restore khusus Administrator
			arsip := protected.Group("/arsip")
			arsip.Use(roleMiddleware.RequirePermission("archives.view"))
			{
				arsip.GET("", arsipHandler.Index)
				arsip.GET("/unor", arsipHandler.DaftarUNOR)
				arsip.GET("/stats", arsipHandler.Stats)
				arsip.GET("/export", arsipHandler.Export)
				arsip.GET("/:id", arsipHandler.Show)

				if opt.OperasionalHandler != nil {
					arsip.GET("/:id/dokumen", opt.OperasionalHandler.ListDokumen)
				}

				arsipWrite := arsip.Group("")
				arsipWrite.Use(roleMiddleware.RequireRole(model.RoleAdministrator, model.RoleOperator))
				{
					arsipWrite.POST("", arsipHandler.Store)
					arsipWrite.PUT("/:id", arsipHandler.Update)
					arsipWrite.POST("/import", arsipHandler.Import)
					if opt.OperasionalHandler != nil {
						arsipWrite.POST("/:id/dokumen", opt.OperasionalHandler.UploadDokumen)
					}
				}

				arsipAdmin := arsip.Group("")
				arsipAdmin.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
				{
					arsipAdmin.DELETE("/:id", arsipHandler.Destroy)
					arsipAdmin.PUT("/:id/restore", arsipHandler.Restore)
				}
			}

			// 4.4 Master Data — Read untuk role berizin settings.view; Write & Konfigurasi khusus Administrator
			if opt.MasterHandler != nil {
				master := protected.Group("/master")
				master.Use(roleMiddleware.RequirePermission("settings.view"))
				{
					master.GET("/kode-klasifikasi", opt.MasterHandler.ListKodeKlasifikasi)
					master.GET("/lokasi-simpan", opt.MasterHandler.ListLokasiSimpan)
					master.GET("/unit-kerja", opt.MasterHandler.ListUnitKerja)

					masterAdmin := master.Group("")
					masterAdmin.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
					{
						masterAdmin.POST("/kode-klasifikasi", opt.MasterHandler.StoreKodeKlasifikasi)
						masterAdmin.PUT("/kode-klasifikasi/:id", opt.MasterHandler.UpdateKodeKlasifikasi)
						masterAdmin.DELETE("/kode-klasifikasi/:id", opt.MasterHandler.DestroyKodeKlasifikasi)

						masterAdmin.POST("/lokasi-simpan", opt.MasterHandler.StoreLokasiSimpan)
						masterAdmin.PUT("/lokasi-simpan/:id", opt.MasterHandler.UpdateLokasiSimpan)
						masterAdmin.DELETE("/lokasi-simpan/:id", opt.MasterHandler.DestroyLokasiSimpan)

						masterAdmin.POST("/unit-kerja", opt.MasterHandler.StoreUnitKerja)
						masterAdmin.PUT("/unit-kerja/:id", opt.MasterHandler.UpdateUnitKerja)
						masterAdmin.DELETE("/unit-kerja/:id", opt.MasterHandler.DestroyUnitKerja)

						masterAdmin.GET("/konfigurasi", opt.MasterHandler.GetKonfigurasi)
						masterAdmin.PUT("/konfigurasi", opt.MasterHandler.UpdateKonfigurasi)
					}
				}
			}

			// 4.5, 4.6, 4.8 Peminjaman, Penyusutan, Dokumen, & Audit Log
			if opt.OperasionalHandler != nil {
				peminjaman := protected.Group("/peminjaman")
				peminjaman.Use(roleMiddleware.RequirePermission("loans.view"))
				{
					peminjaman.GET("", opt.OperasionalHandler.ListPeminjaman)
					peminjaman.GET("/:id", opt.OperasionalHandler.ShowPeminjaman)

					peminjamanWrite := peminjaman.Group("")
					peminjamanWrite.Use(roleMiddleware.RequireRole(model.RoleAdministrator, model.RoleOperator))
					{
						peminjamanWrite.POST("", opt.OperasionalHandler.StorePeminjaman)
						peminjamanWrite.PUT("/:id", opt.OperasionalHandler.UpdatePeminjaman)
						peminjamanWrite.PUT("/:id/kembali", opt.OperasionalHandler.KembalikanPeminjaman)
					}

					peminjamanAdmin := peminjaman.Group("")
					peminjamanAdmin.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
					{
						peminjamanAdmin.DELETE("/:id", opt.OperasionalHandler.DestroyPeminjaman)
					}
				}

				penyusutan := protected.Group("/penyusutan")
				penyusutan.Use(roleMiddleware.RequirePermission("disposal.view"))
				{
					penyusutan.GET("", opt.OperasionalHandler.ListPenyusutan)
					penyusutan.GET("/:id", opt.OperasionalHandler.ShowPenyusutan)

					penyusutanWrite := penyusutan.Group("")
					penyusutanWrite.Use(roleMiddleware.RequireRole(model.RoleAdministrator, model.RoleOperator))
					{
						penyusutanWrite.POST("", opt.OperasionalHandler.StorePenyusutan)
					}

					penyusutanApprove := penyusutan.Group("")
					penyusutanApprove.Use(roleMiddleware.RequireRole(model.RoleApprover, model.RoleAdministrator))
					{
						penyusutanApprove.PUT("/:id/approve", opt.OperasionalHandler.ApprovePenyusutan)
						penyusutanApprove.PUT("/:id/reject", opt.OperasionalHandler.RejectPenyusutan)
					}

					penyusutanAdmin := penyusutan.Group("")
					penyusutanAdmin.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
					{
						penyusutanAdmin.PUT("/:id/selesai", opt.OperasionalHandler.SelesaiPenyusutan)
					}
				}

				dokumen := protected.Group("/dokumen")
				{
					dokumen.GET("/:id/download", opt.OperasionalHandler.DownloadDokumen)

					dokumenWrite := dokumen.Group("")
					dokumenWrite.Use(roleMiddleware.RequireRole(model.RoleAdministrator, model.RoleOperator))
					{
						dokumenWrite.DELETE("/:id", opt.OperasionalHandler.DestroyDokumen)
					}
				}

				auditLogs := protected.Group("/audit-logs")
				auditLogs.Use(roleMiddleware.RequireRole(model.RoleAdministrator))
				{
					auditLogs.GET("", opt.OperasionalHandler.ListAuditLogs)
				}
			}

			// 4.7 & 4.8 Dashboard, Laporan DAI, & Cetak Label PDF (Semua Role sesuai Matriks RBAC)
			if opt.LaporanHandler != nil {
				dashboard := protected.Group("/dashboard")
				dashboard.Use(roleMiddleware.RequirePermission("dashboard.view"))
				{
					dashboard.GET("/stats", opt.LaporanHandler.DashboardStats)
					dashboard.GET("/chart", opt.LaporanHandler.DashboardChart)
				}

				laporan := protected.Group("/laporan")
				laporan.Use(roleMiddleware.RequirePermission("reports.view"))
				{
					laporan.GET("/rekap", opt.LaporanHandler.LaporanRekap)
					laporan.GET("/dai", opt.LaporanHandler.LaporanDAI)
					laporan.GET("/export", opt.LaporanHandler.LaporanExport)
				}

				label := protected.Group("/label")
				label.Use(roleMiddleware.RequirePermission("reports.export"))
				{
					label.GET("/box", opt.LaporanHandler.LabelBox)
					label.GET("/folder", opt.LaporanHandler.LabelFolder)
					label.GET("/inaktif", opt.LaporanHandler.LabelInaktif)
				}
			}
		}
	}

	return r
}
