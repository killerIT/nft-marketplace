package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"

	"github.com/xiaomait/backend/internal/blockchain"
	"github.com/xiaomait/backend/internal/config"
	"github.com/xiaomait/backend/internal/handler"
	"github.com/xiaomait/backend/internal/repository"
	"github.com/xiaomait/backend/internal/service"
)

func main() {
	// 加载 .env 文件
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// 加载配置
	cfg := config.Load()

	// 验证配置
	if err := cfg.Validate(); err != nil {
		log.Fatalf("Invalid configuration: %v", err)
	}

	// 打印配置信息
	cfg.Print()

	// 初始化数据库
	db, err := initDB(cfg)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	log.Println("✓ Database connected successfully")

	// 初始化区块链客户端
	blockchainClient, err := blockchain.NewClient(cfg.EthereumRPC, cfg.MarketplaceAddress)
	if err != nil {
		log.Fatalf("Failed to initialize blockchain client: %v", err)
	}
	log.Println("✓ Blockchain client initialized")

	// 初始化仓储层
	nftRepo := repository.NewNFTRepository(db)
	listingRepo := repository.NewListingRepository(db)
	txRepo := repository.NewTransactionRepository(db)

	// 初始化服务层
	nftService := service.NewNFTService(nftRepo, blockchainClient)
	listingService := service.NewListingService(listingRepo, blockchainClient)
	txService := service.NewTransactionService(txRepo, blockchainClient)

	// 初始化处理器
	nftHandler := handler.NewNFTHandler(nftService)
	listingHandler := handler.NewListingHandler(listingService)
	txHandler := handler.NewTransactionHandler(txService)

	// 启动区块链事件监听器
	if cfg.IsDevelopment() || cfg.IsStaging() {
		go startEventListener(blockchainClient, listingService, txService)
		log.Println("✓ Event listeners started")
	}

	// 初始化 Gin 路由
	router := setupRouter(cfg, nftHandler, listingHandler, txHandler)

	// 创建 HTTP 服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%s", cfg.ServerPort),
		Handler:        router,
		ReadTimeout:    15 * time.Second,
		WriteTimeout:   15 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	// 启动服务器
	go func() {
		log.Printf("🚀 Server starting on http://localhost:%s", cfg.ServerPort)
		log.Printf("📊 Health check: http://localhost:%s/health", cfg.ServerPort)
		log.Printf("📚 API docs: http://localhost:%s/api/v1", cfg.ServerPort)

		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 启动 Metrics 服务器（如果启用）
	if cfg.EnableMetrics {
		go startMetricsServer(cfg.MetricsPort)
	}

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("🛑 Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	// 关闭数据库连接
	sqlDB, err := db.DB()
	if err == nil {
		sqlDB.Close()
	}

	// 关闭区块链客户端
	blockchainClient.Close()

	log.Println("✓ Server exited gracefully")
}

// initDB 初始化数据库连接
func initDB(cfg *config.Config) (*gorm.DB, error) {
	// 构建 DSN
	dsn := cfg.GetDSN()

	// 配置 GORM 日志
	var gormLogger logger.Interface
	switch cfg.LogLevel {
	case "debug":
		gormLogger = logger.Default.LogMode(logger.Info)
	case "info":
		gormLogger = logger.Default.LogMode(logger.Warn)
	default:
		gormLogger = logger.Default.LogMode(logger.Error)
	}

	// 连接数据库
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().UTC()
		},
		PrepareStmt:              true, // 预编译 SQL
		DisableNestedTransaction: false,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}

	// 获取底层 SQL DB
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 配置连接池
	sqlDB.SetMaxOpenConns(cfg.DBMaxOpenConns)
	sqlDB.SetMaxIdleConns(cfg.DBMaxIdleConns)
	sqlDB.SetConnMaxLifetime(cfg.DBConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(cfg.DBConnMaxIdleTime)

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := sqlDB.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	// 自动迁移（开发环境）
	/*if cfg.IsDevelopment() {
		if err := autoMigrate(db); err != nil {
			return nil, fmt.Errorf("failed to auto migrate: %w", err)
		}
		log.Println("✓ Database auto-migration completed")
	}*/

	// 打印连接池状态
	printDBStats(sqlDB)

	return db, nil
}

// autoMigrate 自动迁移数据库表
func autoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&repository.NFT{},
		&repository.Listing{},
		&repository.Transaction{},
		// 添加其他模型...
	)
}

// printDBStats 打印数据库连接池状态
func printDBStats(db *sql.DB) {
	stats := db.Stats()
	log.Printf("Database Pool Stats:")
	log.Printf("  - MaxOpenConns: %d", stats.MaxOpenConnections)
	log.Printf("  - OpenConns: %d", stats.OpenConnections)
	log.Printf("  - InUse: %d", stats.InUse)
	log.Printf("  - Idle: %d", stats.Idle)
}

// setupRouter 设置路由
func setupRouter(
	cfg *config.Config,
	nftHandler *handler.NFTHandler,
	listingHandler *handler.ListingHandler,
	txHandler *handler.TransactionHandler,
) *gin.Engine {
	// 设置 Gin 模式
	if cfg.IsProduction() {
		gin.SetMode(gin.ReleaseMode)
	}

	router := gin.New()

	// 中间件
	router.Use(gin.Logger())
	router.Use(gin.Recovery())

	// CORS 配置
	router.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.AllowedOrigins,
		AllowMethods:     cfg.AllowedMethods,
		AllowHeaders:     cfg.AllowedHeaders,
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	// 限制请求体大小
	router.MaxMultipartMemory = cfg.MaxRequestBodySize

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":      "ok",
			"environment": cfg.Environment,
			"timestamp":   time.Now().UTC(),
		})
	})

	// 系统信息
	router.GET("/info", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"version":     "1.0.0",
			"environment": cfg.Environment,
			"chain_id":    cfg.ChainID,
			"marketplace": cfg.MarketplaceAddress,
		})
	})

	// API 路由
	v1 := router.Group("/api/v1")
	{
		// NFT 路由
		nfts := v1.Group("/nfts")
		{
			nfts.GET("", nftHandler.GetNFTs)
			nfts.GET("/:id", nftHandler.GetNFT)
			nfts.POST("", nftHandler.CreateNFT)
			nfts.GET("/user/:address", nftHandler.GetUserNFTs)
			nfts.GET("/contract/:address", nftHandler.GetNFTsByContract)
		}

		// 挂单路由
		listings := v1.Group("/listings")
		{
			listings.GET("", listingHandler.GetActiveListings)
			listings.GET("/:id", listingHandler.GetListing)
			listings.POST("", listingHandler.CreateListing)
			listings.DELETE("/:id", listingHandler.CancelListing)
			listings.GET("/user/:address", listingHandler.GetUserListings)
			listings.GET("/search", listingHandler.SearchListings)
		}

		// 交易路由
		transactions := v1.Group("/transactions")
		{
			transactions.GET("", txHandler.GetTransactions)
			transactions.GET("/:hash", txHandler.GetTransaction)
			transactions.GET("/user/:address", txHandler.GetUserTransactions)
			transactions.GET("/nft/:contract/:tokenId", txHandler.GetNFTTransactions)
		}

		// 市场统计
		stats := v1.Group("/stats")
		{
			stats.GET("", listingHandler.GetMarketStats)
			stats.GET("/collections/:address", listingHandler.GetCollectionStats)
		}
	}

	return router
}

// startEventListener 启动事件监听器
func startEventListener(
	client *blockchain.Client,
	listingService *service.ListingService,
	txService *service.TransactionService,
) {
	log.Println("Starting blockchain event listener...")
	// 监听 MarketItemCreated 事件
	go func() {
		events := client.ListenMarketItemCreated()
		log.Println("MarketItemCreated listener started")
		for event := range events {
			log.Printf("📝 MarketItemCreated: ItemID=%d, Price=%s",
				event.ItemId, event.Price.String())

			if err := listingService.UpdateFromEvent(event); err != nil {
				log.Printf("Error updating listing from event: %v", err)
			}
		}
	}()

	// 监听 MarketItemSold 事件
	go func() {
		events := client.ListenMarketItemSold()
		for event := range events {
			log.Printf("💰 MarketItemSold: ItemID=%d, Buyer=%s",
				event.ItemId, event.Buyer.Hex())

			if err := txService.RecordSale(event); err != nil {
				log.Printf("Error recording sale: %v", err)
			}
		}
	}()

	log.Println("✓ Event listeners are running")
}

// startMetricsServer 启动 Metrics 服务器
func startMetricsServer(port string) {
	mux := http.NewServeMux()

	mux.HandleFunc("/metrics", func(w http.ResponseWriter, r *http.Request) {
		// 这里可以集成 Prometheus metrics
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "# Metrics endpoint\n")
	})

	addr := fmt.Sprintf(":%s", port)
	log.Printf("📊 Metrics server starting on http://localhost:%s/metrics", port)

	if err := http.ListenAndServe(addr, mux); err != nil {
		log.Printf("Metrics server error: %v", err)
	}
}
