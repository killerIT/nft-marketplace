package config

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

// Config 应用配置结构
type Config struct {
	// 服务器配置
	ServerPort  string
	Environment string // development, staging, production

	// 数据库配置
	DBHost     string
	DBPort     string
	DBUser     string
	DBPassword string
	DBName     string
	DBSSLMode  string
	DBTimezone string

	// 数据库连接池配置
	DBMaxOpenConns    int
	DBMaxIdleConns    int
	DBConnMaxLifetime time.Duration
	DBConnMaxIdleTime time.Duration

	// Redis 配置
	RedisHost     string
	RedisPort     string
	RedisPassword string
	RedisDB       int

	// 区块链配置
	EthereumRPC        string
	MarketplaceAddress string
	NFTContractAddress string
	ChainID            int64

	// 区块链同步配置
	StartBlock          uint64
	BlockConfirmations  uint64
	SyncBatchSize       uint64
	EventProcessWorkers int

	// API 配置
	RateLimitPerMinute int
	MaxPageSize        int
	DefaultPageSize    int

	// JWT 配置
	JWTSecret     string
	JWTExpiration time.Duration

	// CORS 配置
	AllowedOrigins []string
	AllowedMethods []string
	AllowedHeaders []string

	// 文件存储配置
	StorageProvider string // local, s3, ipfs
	S3Bucket        string
	S3Region        string
	S3AccessKey     string
	S3SecretKey     string
	IPFSGateway     string

	// 日志配置
	LogLevel  string // debug, info, warn, error
	LogFormat string // json, text

	// 监控配置
	EnableMetrics bool
	MetricsPort   string
	EnablePprof   bool
	PprofPort     string

	// 第三方服务
	EtherscanAPIKey     string
	InfuraProjectID     string
	AlchemyAPIKey       string
	CoinMarketCapAPIKey string

	// 邮件配置
	SMTPHost     string
	SMTPPort     int
	SMTPUser     string
	SMTPPassword string
	SMTPFrom     string

	// 缓存配置
	CacheTTL          time.Duration
	EnableRedisCache  bool
	EnableMemoryCache bool

	// 安全配置
	EnableRateLimit    bool
	TrustedProxies     []string
	MaxRequestBodySize int64
}

// Load 从环境变量加载配置
func Load() *Config {
	return &Config{
		// 服务器配置
		ServerPort:  getEnv("SERVER_PORT", "8080"),
		Environment: getEnv("ENVIRONMENT", "development"),

		// 数据库配置
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", ""),
		DBName:     getEnv("DB_NAME", "nft_marketplace"),
		DBSSLMode:  getEnv("DB_SSL_MODE", "disable"),
		DBTimezone: getEnv("DB_TIMEZONE", "UTC"),

		// 数据库连接池配置
		DBMaxOpenConns:    getEnvAsInt("DB_MAX_OPEN_CONNS", 25),
		DBMaxIdleConns:    getEnvAsInt("DB_MAX_IDLE_CONNS", 5),
		DBConnMaxLifetime: getEnvAsDuration("DB_CONN_MAX_LIFETIME", 5*time.Minute),
		DBConnMaxIdleTime: getEnvAsDuration("DB_CONN_MAX_IDLE_TIME", 10*time.Minute),

		// Redis 配置
		RedisHost:     getEnv("REDIS_HOST", "localhost"),
		RedisPort:     getEnv("REDIS_PORT", "6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       getEnvAsInt("REDIS_DB", 0),

		// 区块链配置
		EthereumRPC:        getEnv("ETHEREUM_RPC", ""),
		MarketplaceAddress: getEnv("MARKETPLACE_ADDRESS", ""),
		NFTContractAddress: getEnv("NFT_CONTRACT_ADDRESS", ""),
		ChainID:            getEnvAsInt64("CHAIN_ID", 11155111),

		// 区块链同步配置
		StartBlock:          getEnvAsUint64("START_BLOCK", 0),
		BlockConfirmations:  getEnvAsUint64("BLOCK_CONFIRMATIONS", 12),
		SyncBatchSize:       getEnvAsUint64("SYNC_BATCH_SIZE", 1000),
		EventProcessWorkers: getEnvAsInt("EVENT_PROCESS_WORKERS", 5),

		// API 配置
		RateLimitPerMinute: getEnvAsInt("RATE_LIMIT_PER_MINUTE", 100),
		MaxPageSize:        getEnvAsInt("MAX_PAGE_SIZE", 100),
		DefaultPageSize:    getEnvAsInt("DEFAULT_PAGE_SIZE", 20),

		// JWT 配置
		JWTSecret:     getEnv("JWT_SECRET", "your-secret-key-change-in-production"),
		JWTExpiration: getEnvAsDuration("JWT_EXPIRATION", 24*time.Hour),

		// CORS 配置
		AllowedOrigins: getEnvAsSlice("ALLOWED_ORIGINS", []string{"*"}),
		AllowedMethods: getEnvAsSlice("ALLOWED_METHODS", []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"}),
		AllowedHeaders: getEnvAsSlice("ALLOWED_HEADERS", []string{"Origin", "Content-Type", "Authorization"}),

		// 文件存储配置
		StorageProvider: getEnv("STORAGE_PROVIDER", "local"),
		S3Bucket:        getEnv("S3_BUCKET", ""),
		S3Region:        getEnv("S3_REGION", "us-east-1"),
		S3AccessKey:     getEnv("S3_ACCESS_KEY", ""),
		S3SecretKey:     getEnv("S3_SECRET_KEY", ""),
		IPFSGateway:     getEnv("IPFS_GATEWAY", "https://ipfs.io"),

		// 日志配置
		LogLevel:  getEnv("LOG_LEVEL", "info"),
		LogFormat: getEnv("LOG_FORMAT", "json"),

		// 监控配置
		EnableMetrics: getEnvAsBool("ENABLE_METRICS", true),
		MetricsPort:   getEnv("METRICS_PORT", "9090"),
		EnablePprof:   getEnvAsBool("ENABLE_PPROF", false),
		PprofPort:     getEnv("PPROF_PORT", "6060"),

		// 第三方服务
		EtherscanAPIKey:     getEnv("ETHERSCAN_API_KEY", ""),
		InfuraProjectID:     getEnv("INFURA_PROJECT_ID", ""),
		AlchemyAPIKey:       getEnv("ALCHEMY_API_KEY", ""),
		CoinMarketCapAPIKey: getEnv("COINMARKETCAP_API_KEY", ""),

		// 邮件配置
		SMTPHost:     getEnv("SMTP_HOST", ""),
		SMTPPort:     getEnvAsInt("SMTP_PORT", 587),
		SMTPUser:     getEnv("SMTP_USER", ""),
		SMTPPassword: getEnv("SMTP_PASSWORD", ""),
		SMTPFrom:     getEnv("SMTP_FROM", "noreply@nftmarketplace.com"),

		// 缓存配置
		CacheTTL:          getEnvAsDuration("CACHE_TTL", 5*time.Minute),
		EnableRedisCache:  getEnvAsBool("ENABLE_REDIS_CACHE", true),
		EnableMemoryCache: getEnvAsBool("ENABLE_MEMORY_CACHE", true),

		// 安全配置
		EnableRateLimit:    getEnvAsBool("ENABLE_RATE_LIMIT", true),
		TrustedProxies:     getEnvAsSlice("TRUSTED_PROXIES", []string{}),
		MaxRequestBodySize: getEnvAsInt64("MAX_REQUEST_BODY_SIZE", 10*1024*1024), // 10MB
	}
}

// GetDSN 返回数据库 DSN 连接字符串
func (c *Config) GetDSN() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=%s TimeZone=%s",
		c.DBHost,
		c.DBPort,
		c.DBUser,
		c.DBPassword,
		c.DBName,
		c.DBSSLMode,
		c.DBTimezone,
	)
}

// GetRedisAddr 返回 Redis 地址
func (c *Config) GetRedisAddr() string {
	return fmt.Sprintf("%s:%s", c.RedisHost, c.RedisPort)
}

// IsProduction 判断是否为生产环境
func (c *Config) IsProduction() bool {
	return c.Environment == "production"
}

// IsDevelopment 判断是否为开发环境
func (c *Config) IsDevelopment() bool {
	return c.Environment == "development"
}

// IsStaging 判断是否为测试环境
func (c *Config) IsStaging() bool {
	return c.Environment == "staging"
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.DBHost == "" {
		return fmt.Errorf("DB_HOST is required")
	}

	if c.DBName == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	if c.EthereumRPC == "" {
		return fmt.Errorf("ETHEREUM_RPC is required")
	}

	if c.MarketplaceAddress == "" {
		return fmt.Errorf("MARKETPLACE_ADDRESS is required")
	}

	if c.IsProduction() && c.JWTSecret == "your-secret-key-change-in-production" {
		return fmt.Errorf("JWT_SECRET must be changed in production")
	}

	return nil
}

// Print 打印配置信息（隐藏敏感信息）
func (c *Config) Print() {
	fmt.Println("=== Application Configuration ===")
	fmt.Printf("Environment: %s\n", c.Environment)
	fmt.Printf("Server Port: %s\n", c.ServerPort)
	fmt.Printf("Database: %s@%s:%s/%s\n", c.DBUser, c.DBHost, c.DBPort, c.DBName)
	fmt.Printf("Redis: %s:%s\n", c.RedisHost, c.RedisPort)
	fmt.Printf("Ethereum RPC: %s\n", c.EthereumRPC)
	fmt.Printf("Marketplace Address: %s\n", c.MarketplaceAddress)
	fmt.Printf("Chain ID: %d\n", c.ChainID)
	fmt.Printf("Log Level: %s\n", c.LogLevel)
	fmt.Printf("Metrics Enabled: %v\n", c.EnableMetrics)
	fmt.Println("=================================")
}

// ===== 辅助函数 =====

// getEnv 获取环境变量，如果不存在则返回默认值
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt 获取整数类型的环境变量
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := os.Getenv(key)
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsInt64 获取 int64 类型的环境变量
func getEnvAsInt64(key string, defaultValue int64) int64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseInt(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsUint64 获取 uint64 类型的环境变量
func getEnvAsUint64(key string, defaultValue uint64) uint64 {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseUint(valueStr, 10, 64); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool 获取布尔类型的环境变量
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := os.Getenv(key)
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsDuration 获取时间间隔类型的环境变量
func getEnvAsDuration(key string, defaultValue time.Duration) time.Duration {
	valueStr := os.Getenv(key)
	if value, err := time.ParseDuration(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsSlice 获取字符串切片类型的环境变量（逗号分隔）
func getEnvAsSlice(key string, defaultValue []string) []string {
	valueStr := os.Getenv(key)
	if valueStr == "" {
		return defaultValue
	}

	var result []string
	for _, v := range splitAndTrim(valueStr, ",") {
		if v != "" {
			result = append(result, v)
		}
	}

	if len(result) == 0 {
		return defaultValue
	}

	return result
}

// splitAndTrim 分割字符串并去除空格
func splitAndTrim(s, sep string) []string {
	var result []string
	for _, item := range splitString(s, sep) {
		if trimmed := trimSpace(item); trimmed != "" {
			result = append(result, trimmed)
		}
	}
	return result
}

// splitString 简单的字符串分割
func splitString(s, sep string) []string {
	if s == "" {
		return []string{}
	}

	var result []string
	start := 0

	for i := 0; i < len(s); i++ {
		if s[i:i+len(sep)] == sep {
			result = append(result, s[start:i])
			start = i + len(sep)
			i += len(sep) - 1
		}
	}

	result = append(result, s[start:])
	return result
}

// trimSpace 去除字符串首尾空格
func trimSpace(s string) string {
	start := 0
	end := len(s)

	for start < end && (s[start] == ' ' || s[start] == '\t' || s[start] == '\n' || s[start] == '\r') {
		start++
	}

	for end > start && (s[end-1] == ' ' || s[end-1] == '\t' || s[end-1] == '\n' || s[end-1] == '\r') {
		end--
	}

	return s[start:end]
}
