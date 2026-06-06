package router

import (
	"net/http"
	"path/filepath"

	"github.com/gin-gonic/gin"
	"github.com/instaagrammeta/crm-real/backend-go/internal/auth"
	"github.com/instaagrammeta/crm-real/backend-go/internal/config"
	"github.com/instaagrammeta/crm-real/backend-go/internal/handlers"
	"github.com/instaagrammeta/crm-real/backend-go/internal/middleware"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services"
	"github.com/instaagrammeta/crm-real/backend-go/internal/services/telegram"
	"github.com/instaagrammeta/crm-real/backend-go/internal/ws"
	"gorm.io/gorm"
)

// New builds the *gin.Engine and registers all HTTP + WS routes.
func New(
	cfg *config.Config,
	db *gorm.DB,
	hub *ws.Hub,
	notifier *services.Notifier,
	tg *telegram.Bot,
) *gin.Engine {

	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.New()
	r.Use(middleware.Recovery())
	r.Use(middleware.Logger())
	r.Use(middleware.CORS(cfg.CORSOrigins))
	r.MaxMultipartMemory = int64(cfg.MaxUploadMB) << 20

	jwtMgr := auth.NewManager(cfg.JWTSecret, cfg.JWTExpireHours)

	// ====================================================================
	// Public routes
	// ====================================================================
	r.GET("/healthz", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"ok": true})
	})

	// uploaded files served from disk
	uploadAbs, _ := filepath.Abs(cfg.UploadDir)
	r.Static("/uploads", uploadAbs)

	// WebSocket endpoint (auth via ?token=)
	r.GET("/ws", hub.HandleWS(jwtMgr))

	// ====================================================================
	// Handlers
	// ====================================================================
	authH := &handlers.AuthHandler{DB: db, JWTMgr: jwtMgr}
	usersH := &handlers.UsersHandler{DB: db, Cfg: cfg, TG: tg}
	tasksH := &handlers.TasksHandler{DB: db, Notifier: notifier}
	lidsH := &handlers.LidsHandler{DB: db, Notifier: notifier}
	kanbanH := &handlers.KanbanHandler{DB: db, Notifier: notifier}
	reqH := &handlers.RequestsHandler{DB: db, Notifier: notifier}
	housesH := &handlers.HousesHandler{DB: db}
	realtyH := &handlers.RealtyHandler{DB: db}
	objektH := &handlers.ObjektHandler{DB: db}
	simH := &handlers.SimHandler{DB: db}
	postsH := &handlers.PostsHandler{DB: db}
	foldersH := &handlers.FoldersHandler{DB: db, Cfg: cfg}
	ipotekaH := &handlers.IpotekaHandler{DB: db}
	rasrochkaH := &handlers.RasrochkaHandler{DB: db}
	chatH := &handlers.ChatHandler{DB: db, Cfg: cfg, Notifier: notifier}
	dashH := &handlers.DashboardHandler{DB: db}
	notifH := &handlers.NotificationsHandler{DB: db}

	api := r.Group("/api")

	// ===== Auth (public) =====
	api.POST("/login", authH.Login)
	api.POST("/logout", authH.Logout)

	// ===== Public read for landing pages (banks, installment objects) =====
	api.GET("/banks", ipotekaH.ListBanks)
	api.GET("/banks/by-slug/:slug", ipotekaH.GetBankBySlug)
	api.GET("/banks/:id", ipotekaH.GetBank)
	api.GET("/banks/:id/conditions", ipotekaH.ListConditions)
	api.GET("/installment-objects", rasrochkaH.ListObjects)
	api.GET("/installment-objects/by-slug/:slug", rasrochkaH.GetObjectBySlug)
	api.GET("/installment-objects/:id", rasrochkaH.GetObject)
	api.GET("/installment-objects/:id/conditions", rasrochkaH.ListConditions)

	// ===== Authenticated =====
	authed := api.Group("")
	authed.Use(middleware.Auth(jwtMgr))

	// auth-self
	authed.GET("/check-auth", authH.CheckAuth)
	authed.GET("/me", authH.Me)
	authed.POST("/me/telegram/link", usersH.GenerateTelegramLink)
	authed.GET("/me/telegram", usersH.ListTelegramSubscribers)
	authed.DELETE("/me/telegram/:id", usersH.UnlinkTelegram)

	// users
	authed.GET("/users", usersH.List)
	authed.GET("/users/list", usersH.Listing)

	// tasks
	authed.GET("/tasks", tasksH.List)
	authed.POST("/tasks", tasksH.Create)
	authed.PUT("/tasks/:id", tasksH.Update)
	authed.DELETE("/tasks/:id", tasksH.Delete)

	// lids
	authed.GET("/lids", lidsH.List)
	authed.POST("/lids", lidsH.Create)
	authed.PUT("/lids/:id", lidsH.Update)
	authed.DELETE("/lids/:id", lidsH.Delete)

	// kanban (lids board)
	authed.GET("/kanban/boards", kanbanH.ListBoards)
	authed.POST("/kanban/boards", kanbanH.CreateBoard)
	authed.PUT("/kanban/boards/:id", kanbanH.UpdateBoard)
	authed.DELETE("/kanban/boards/:id", kanbanH.DeleteBoard)
	authed.PUT("/kanban/boards/:id/archive", kanbanH.ArchiveBoard)
	authed.GET("/kanban/boards/:id/columns", kanbanH.ListColumns)

	authed.POST("/kanban/columns", kanbanH.CreateColumn)
	authed.PUT("/kanban/columns/:id", kanbanH.UpdateColumn)
	authed.DELETE("/kanban/columns/:id", kanbanH.DeleteColumn)
	authed.PUT("/kanban/columns/:id/order", kanbanH.UpdateColumnOrder)

	authed.GET("/kanban/leads", kanbanH.ListLeads)
	authed.POST("/kanban/leads", kanbanH.CreateLead)
	authed.PUT("/kanban/leads/:id", kanbanH.UpdateLead)
	authed.DELETE("/kanban/leads/:id", kanbanH.DeleteLead)
	authed.POST("/kanban/leads/move", kanbanH.MoveLead)
	authed.PUT("/kanban/leads/:id/order", kanbanH.UpdateLeadOrder)

	// requests (legacy)
	authed.GET("/requests", reqH.List)
	authed.POST("/requests", reqH.Create)
	authed.PUT("/requests/:id", reqH.Update)
	authed.DELETE("/requests/:id", reqH.Delete)

	// requests (new kanban-style)
	authed.GET("/requests-boards", reqH.ListBoards)
	authed.POST("/requests-boards", reqH.CreateBoard)
	authed.PUT("/requests-boards/:id", reqH.UpdateBoard)
	authed.DELETE("/requests-boards/:id", reqH.DeleteBoard)
	authed.GET("/requests-boards/:id/columns", reqH.ListColumnsForBoard)

	authed.POST("/requests-columns", reqH.CreateColumn)
	authed.PUT("/requests-columns/:id", reqH.UpdateColumn)
	authed.DELETE("/requests-columns/:id", reqH.DeleteColumn)

	authed.GET("/requests-board/:id/requests", reqH.ListItemsForBoard)
	authed.POST("/requests-new", reqH.CreateItem)
	authed.PUT("/requests-update/:id", reqH.UpdateItem)
	authed.DELETE("/requests-delete/:id", reqH.DeleteItem)
	authed.POST("/requests-move", reqH.MoveItem)

	// houses
	authed.GET("/houses", housesH.List)
	authed.POST("/houses", housesH.Create)
	authed.PUT("/houses/:id", housesH.Update)
	authed.DELETE("/houses/:id", housesH.Delete)

	// realty (объекты с координатами)
	authed.GET("/realty/objects", realtyH.ListObjects)
	authed.POST("/realty/objects", realtyH.CreateObject)
	authed.PUT("/realty/objects/:id", realtyH.UpdateObject)
	authed.DELETE("/realty/objects/:id", realtyH.DeleteObject)
	authed.GET("/realty/objects/:id/blocks", realtyH.ListBlocks)
	authed.POST("/realty/blocks", realtyH.CreateBlock)
	authed.PUT("/realty/blocks/:id", realtyH.UpdateBlock)
	authed.DELETE("/realty/blocks/:id", realtyH.DeleteBlock)
	authed.GET("/realty/pricing", realtyH.ListPricing)
	authed.POST("/realty/pricing", realtyH.CreatePricing)
	authed.PUT("/realty/pricing/:id", realtyH.UpdatePricing)
	authed.DELETE("/realty/pricing/:id", realtyH.DeletePricing)
	authed.GET("/realty/layouts", realtyH.ListLayouts)
	authed.POST("/realty/layouts", realtyH.CreateLayout)
	authed.PUT("/realty/layouts/:id", realtyH.UpdateLayout)
	authed.DELETE("/realty/layouts/:id", realtyH.DeleteLayout)

	// objekt — шахматка
	authed.GET("/objekt/projects", objektH.ListProjects)
	authed.POST("/objekt/projects", objektH.CreateProject)
	authed.PUT("/objekt/projects/:id", objektH.UpdateProject)
	authed.DELETE("/objekt/projects/:id", objektH.DeleteProject)
	authed.GET("/objekt/projects/:id/blocks", objektH.ListBlocks)
	authed.GET("/objekt/projects/:id/apartments", objektH.ListApartments)
	authed.POST("/objekt/blocks", objektH.CreateBlock)
	authed.PUT("/objekt/blocks/:id", objektH.UpdateBlock)
	authed.DELETE("/objekt/blocks/:id", objektH.DeleteBlock)
	authed.PUT("/objekt/apartments/:id", objektH.UpdateApartment)
	authed.PATCH("/objekt/apartments/:id/status", objektH.UpdateApartmentStatus)

	// sim cards
	authed.GET("/company-phones", simH.ListPhones)
	authed.POST("/company-phones", simH.CreatePhone)
	authed.PUT("/company-phones/:id", simH.UpdatePhone)
	authed.DELETE("/company-phones/:id", simH.DeletePhone)
	authed.GET("/sim-cards", simH.ListSims)
	authed.POST("/sim-cards", simH.CreateSim)
	authed.PUT("/sim-cards/:id", simH.UpdateSim)
	authed.DELETE("/sim-cards/:id", simH.DeleteSim)
	authed.GET("/sim-tariffs", simH.ListTariffs)
	authed.POST("/sim-tariffs", simH.CreateTariff)
	authed.PUT("/sim-tariffs/:id", simH.UpdateTariff)
	authed.DELETE("/sim-tariffs/:id", simH.DeleteTariff)
	authed.GET("/tariff-payments", simH.ListPayments)
	authed.DELETE("/tariff-payments/:id", simH.DeletePayment)
	authed.GET("/sim-phones/stats", simH.Stats)

	// posts
	authed.GET("/posts", postsH.List)
	authed.POST("/posts", postsH.Create)
	authed.PUT("/posts/:id", postsH.Update)
	authed.DELETE("/posts/:id", postsH.Delete)

	// folders / files
	authed.GET("/folders", foldersH.List)
	authed.POST("/folders", foldersH.Create)
	authed.PUT("/folders/:id", foldersH.Rename)
	authed.DELETE("/folders/:id", foldersH.Delete)
	authed.GET("/folders/:id/files", foldersH.ListFiles)
	authed.POST("/folders/:id/files", foldersH.UploadFile)
	authed.DELETE("/files/:id", foldersH.DeleteFile)
	authed.GET("/download/:id", foldersH.Download)

	// chat
	authed.GET("/messages", chatH.List)
	authed.POST("/messages", chatH.Send)
	authed.PUT("/messages/:id", chatH.Update)
	authed.DELETE("/messages/:id", chatH.Delete)

	// dashboard
	authed.GET("/dashboard/stats", dashH.Stats)

	// notifications
	authed.GET("/notifications", notifH.List)
	authed.POST("/notifications/read-all", notifH.MarkAllRead)
	authed.POST("/notifications/:id/read", notifH.MarkRead)
	authed.DELETE("/notifications/:id", notifH.Delete)
	authed.GET("/notifications/unread-count", notifH.UnreadCount)

	// ===== Admin only =====
	admin := authed.Group("")
	admin.Use(middleware.AdminOnly())

	admin.POST("/users", usersH.Create)
	admin.PUT("/users/:id", usersH.Update)
	admin.DELETE("/users/:id", usersH.Delete)

	// ipoteka — admin write
	admin.POST("/banks", ipotekaH.CreateBank)
	admin.PUT("/banks/:id", ipotekaH.UpdateBank)
	admin.DELETE("/banks/:id", ipotekaH.DeleteBank)
	admin.POST("/banks/:id/conditions", ipotekaH.CreateCondition)
	admin.PUT("/conditions/:id", ipotekaH.UpdateCondition)
	admin.DELETE("/conditions/:id", ipotekaH.DeleteCondition)

	// rasrochka — admin write
	admin.POST("/installment-objects", rasrochkaH.CreateObject)
	admin.PUT("/installment-objects/:id", rasrochkaH.UpdateObject)
	admin.DELETE("/installment-objects/:id", rasrochkaH.DeleteObject)
	admin.POST("/installment-objects/:id/conditions", rasrochkaH.CreateCondition)
	admin.PUT("/installment-conditions/:id", rasrochkaH.UpdateCondition)
	admin.DELETE("/installment-conditions/:id", rasrochkaH.DeleteCondition)

	return r
}
