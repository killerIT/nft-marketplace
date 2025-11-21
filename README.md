# NFT Marketplace - 生产级项目

一个完整的 NFT 市场平台，采用 **Hardhat + GoLand** 技术栈开发。


 项目亮点
智能合约层

✅ 生产级 Solidity 合约（OpenZeppelin 标准）
✅ 完整的安全机制（ReentrancyGuard、Ownable）
✅ 平台费用系统和统计功能
✅ 全面的单元测试（12+ 测试用例）
✅ Gas 优化和覆盖率报告

Go 后端架构

✅ 分层架构：Handler → Service → Repository
✅ 区块链集成：实时事件监听器
✅ 数据库管理：GORM + PostgreSQL
✅ RESTful API：完整的 CRUD 端点
✅ 错误处理：优雅的错误管理
✅ 优雅关闭：信号处理

DevOps & 监控

✅ Docker Compose：一键启动全部服务
✅ Prometheus + Grafana：实时监控
✅ PgAdmin：数据库管理工具
✅ 健康检查：所有服务的健康监测

🚀 核心功能

NFT 上架与交易

卖家上架 NFT
买家购买 NFT
平台收取 2.5% 手续费
取消挂单功能


市场统计

活跃挂单数量
总交易额
地板价/天花板价
平均价格


实时事件监听

MarketItemCreated 事件
MarketItemSold 事件
自动同步链上数据到数据库


查询功能

分页查询
用户挂单查询
交易历史查询
搜索和筛选



📚 学习要点
从合约开始

研究 NFTMarketplace.sol 的状态管理
理解 Hardhat 的测试框架
学习事件日志的设计模式

Go 后端集成

blockchain/client.go：如何用 Go 调用智能合约
service/ 层：业务逻辑的封装
repository/ 层：数据持久化模式

生产部署

Docker Compose 的服务编排
环境变量管理
监控系统搭建

🎓 使用建议
bash# 1. 启动所有服务
docker-compose up -d

# 2. 查看日志
docker-compose logs -f backend

# 3. 运行合约测试
cd contracts && npx hardhat test

# 4. 测试 API
curl http://localhost:8080/health
curl http://localhost:8080/api/v1/stats


## 🏗️ 技术栈

### 智能合约
- **Solidity 0.8.20**: 智能合约语言
- **Hardhat**: 开发、测试、部署框架
- **OpenZeppelin**: 安全的合约库
- **Ethers.js**: 以太坊交互库

### 后端
- **Go 1.21+**: 后端语言
- **Gin**: Web 框架
- **GORM**: ORM 框架
- **go-ethereum**: 以太坊 Go 客户端
- **PostgreSQL**: 关系数据库
- **Redis**: 缓存

### DevOps
- **Docker & Docker Compose**: 容器化部署
- **Prometheus**: 监控指标
- **Grafana**: 可视化监控

## 📁 项目结构

```
nft-marketplace/
├── contracts/                 # Hardhat 智能合约
│   ├── contracts/
│   │   ├── NFTMarketplace.sol
│   │   └── NFT.sol
│   ├── test/
│   │   └── NFTMarketplace.test.js
│   ├── scripts/
│   │   └── deploy.js
│   └── hardhat.config.js
│
├── backend/                   # Go 后端
│   ├── cmd/api/
│   │   └── main.go
│   ├── internal/
│   │   ├── config/          # 配置管理
│   │   ├── handler/         # HTTP 处理器
│   │   ├── service/         # 业务逻辑
│   │   ├── repository/      # 数据访问
│   │   └── blockchain/      # 区块链交互
│   └── go.mod
│
├── monitoring/               # 监控配置
│   ├── prometheus.yml
│   └── grafana/
│
└── docker-compose.yml       # Docker 编排
```

## 🚀 快速开始

### 1. 环境要求

- **Node.js** 18+
- **Go** 1.21+
- **Docker & Docker Compose**
- **PostgreSQL** 15+ (可用 Docker)

### 2. 安装依赖

#### 智能合约
```bash
cd contracts
npm install
```

#### Go 后端
```bash
cd backend
go mod download
```

### 3. 配置环境变量

创建 `.env` 文件：

```env
# 服务器配置
SERVER_PORT=8080
ENVIRONMENT=development

# 数据库配置
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=nft_marketplace

# Redis 配置
REDIS_HOST=localhost
REDIS_PORT=6379

# 以太坊配置
ETHEREUM_RPC=http://localhost:8545
MARKETPLACE_ADDRESS=0x5FbDB2315678afecb367f032d93F642f64180aa3

# API Keys (生产环境)
ETHERSCAN_API_KEY=your_etherscan_key
COINMARKETCAP_API_KEY=your_cmc_key
```

### 4. 使用 Docker Compose 启动

```bash
# 启动所有服务
docker-compose up -d

# 查看日志
docker-compose logs -f

# 停止服务
docker-compose down
```

服务端口：
- **后端 API**: http://localhost:8080
- **Hardhat 节点**: http://localhost:8545
- **PostgreSQL**: localhost:5432
- **Redis**: localhost:6379
- **Prometheus**: http://localhost:9090
- **Grafana**: http://localhost:3000
- **PgAdmin**: http://localhost:5050

## 🔧 开发流程

### 1. 编译合约

```bash
cd contracts
npx hardhat compile
```

### 2. 运行合约测试

```bash
npx hardhat test

# 带 gas 报告
REPORT_GAS=true npx hardhat test

# 生成覆盖率报告
npx hardhat coverage
```

### 3. 部署合约

```bash
# 部署到本地网络
npx hardhat run scripts/deploy.js --network localhost

# 部署到测试网
npx hardhat run scripts/deploy.js --network sepolia

# 验证合约
npx hardhat verify --network sepolia DEPLOYED_CONTRACT_ADDRESS
```

### 4. 运行 Go 后端

```bash
cd backend

# 开发模式（热重载）
go run cmd/api/main.go

# 编译
go build -o bin/api cmd/api/main.go

# 运行测试
go test ./...

# 运行带覆盖率的测试
go test -cover ./...
```

## 📡 API 文档

### NFT 端点

#### 获取所有 NFT
```http
GET /api/v1/nfts?page=1&page_size=20
```

#### 获取单个 NFT
```http
GET /api/v1/nfts/:id
```

#### 创建 NFT
```http
POST /api/v1/nfts
Content-Type: application/json

{
  "contract_address": "0x...",
  "token_id": "1",
  "owner": "0x...",
  "metadata_uri": "ipfs://..."
}
```

### 挂单端点

#### 获取活跃挂单
```http
GET /api/v1/listings?page=1&page_size=20
```

#### 创建挂单
```http
POST /api/v1/listings
Content-Type: application/json

{
  "item_id": 1,
  "nft_contract": "0x...",
  "token_id": "1",
  "seller": "0x...",
  "price": "1000000000000000000",
  "tx_hash": "0x..."
}
```

#### 取消挂单
```http
DELETE /api/v1/listings/:id
```

### 市场统计
```http
GET /api/v1/stats
```

响应：
```json
{
  "active_listings": 150,
  "total_listings": 500,
  "sold_listings": 350,
  "total_volume": "12500000000000000000000",
  "average_price": "500000000000000000",
  "floor_price": "100000000000000000",
  "ceiling_price": "5000000000000000000"
}
```

## 🧪 测试

### 合约测试
```bash
cd contracts
npx hardhat test
```

输出示例：
```
  NFTMarketplace
    Deployment
      ✓ Should set the correct owner
      ✓ Should set the correct platform fee
    Listing
      ✓ Should create a market item
      ✓ Should fail if price is zero
    Purchasing
      ✓ Should complete a sale
      ✓ Should fail if incorrect price sent

  12 passing (2s)
```

### Go 后端测试
```bash
cd backend
go test -v ./internal/...
```

## 📊 监控

### Prometheus 指标
访问 http://localhost:9090 查看：
- HTTP 请求延迟
- 数据库查询性能
- 区块链事件处理
- 错误率

### Grafana 仪表板
访问 http://localhost:3000 (admin/admin)

预配置仪表板：
- API 性能监控
- 数据库连接池
- 区块链事件统计
- 市场交易统计

## 🔐 安全最佳实践

### 智能合约
- ✅ 使用 OpenZeppelin 审计过的合约
- ✅ 实现 ReentrancyGuard
- ✅ 所有外部调用都有错误处理
- ✅ 使用 SafeMath (Solidity 0.8+)
- ✅ 完整的事件日志

### 后端
- ✅ 参数验证
- ✅ SQL 注入防护 (GORM)
- ✅ CORS 配置
- ✅ Rate Limiting
- ✅ 错误处理和日志

## 🚢 生产部署

### 1. 准备生产环境变量

```env
ENVIRONMENT=production
ETHEREUM_RPC=https://mainnet.infura.io/v3/YOUR_KEY
DB_PASSWORD=strong_password
```

### 2. 编译优化

```bash
# 合约
npx hardhat compile --optimizer

# Go 后端
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -ldflags '-w -s' -o bin/api cmd/api/main.go
```

### 3. 数据库迁移

```bash
# 使用 GORM AutoMigrate 或 migrate 工具
go run cmd/migrate/main.go
```

### 4. 部署到云平台

#### AWS ECS
```bash
# 构建镜像
docker build -t nft-marketplace-backend ./backend

# 推送到 ECR
aws ecr get-login-password | docker login --username AWS --password-stdin
docker push your-registry/nft-marketplace-backend
```

#### Kubernetes
```bash
kubectl apply -f k8s/deployment.yaml
kubectl apply -f k8s/service.yaml
```

## 📈 性能优化

### 数据库
- 索引优化（seller, nft_contract, status）
- 连接池配置
- 查询缓存（Redis）

### 区块链
- 批量事件处理
- 缓存合约调用结果
- 使用 WebSocket 订阅事件

### API
- 响应压缩
- CDN 静态资源
- 分页限制

## 🤝 贡献指南

1. Fork 项目
2. 创建特性分支 (`git checkout -b feature/AmazingFeature`)
3. 提交更改 (`git commit -m 'Add some AmazingFeature'`)
4. 推送到分支 (`git push origin feature/AmazingFeature`)
5. 开启 Pull Request

## 📄 许可证

MIT License

## 🔗 相关资源

- [Hardhat 文档](https://hardhat.org/docs)
- [Go Ethereum 文档](https://geth.ethereum.org/docs)
- [Gin 文档](https://gin-gonic.com/docs/)
- [GORM 文档](https://gorm.io/docs/)
- [OpenZeppelin 合约](https://docs.openzeppelin.com/contracts/)

## 💬 支持

有问题？欢迎提交 Issue 或加入我们的社区！

---

**Happy Coding! 🎉**