package service

import (
	"context"
	"fmt"
	"github.com/ethereum/go-ethereum/common"
	"log"
	"math/big"
	"time"

	"github.com/xiaomait/backend/internal/blockchain"
	"github.com/xiaomait/backend/internal/repository"
)

// ListingService 挂单服务
type ListingService struct {
	repo     *repository.ListingRepository
	bcClient *blockchain.Client
}

// NewListingService 创建挂单服务
func NewListingService(repo *repository.ListingRepository, bcClient *blockchain.Client) *ListingService {
	return &ListingService{
		repo:     repo,
		bcClient: bcClient,
	}
}

// CreateListingRequest 创建挂单请求
type CreateListingRequest struct {
	ItemID      uint64 `json:"item_id" binding:"required"`
	NFTContract string `json:"nft_contract" binding:"required"`
	TokenID     string `json:"token_id" binding:"required"`
	Seller      string `json:"seller" binding:"required"`
	Price       string `json:"price" binding:"required"`
	TxHash      string `json:"tx_hash" binding:"required"`
}

// ListingResponse 挂单响应
type ListingResponse struct {
	ID          uint      `json:"id"`
	ItemID      uint64    `json:"item_id"`
	NFTContract string    `json:"nft_contract"`
	TokenID     string    `json:"token_id"`
	Seller      string    `json:"seller"`
	Price       string    `json:"price"`
	Status      string    `json:"status"`
	ListedAt    time.Time `json:"listed_at"`
	CreatedAt   time.Time `json:"created_at"`
}

// CreateListing 创建挂单
func (s *ListingService) CreateListing(ctx context.Context, req *CreateListingRequest) (*ListingResponse, error) {
	// 验证链上数据
	itemID := big.NewInt(int64(req.ItemID))
	itemData, err := s.bcClient.GetMarketItem(ctx, itemID)
	if err != nil {
		return nil, fmt.Errorf("failed to verify on-chain data: %w", err)
	}
	log.Printf("Market itemData: %+v", itemData)

	chainNFTContract := itemData["nftContract"].(string)
	reqNFTContract := req.NFTContract
	// 检查数据一致性
	if common.HexToAddress(chainNFTContract) != common.HexToAddress(reqNFTContract) {
		return nil, fmt.Errorf("nft contract mismatch")
	}

	listing := &repository.Listing{
		ItemID:      req.ItemID,
		NFTContract: req.NFTContract,
		TokenID:     req.TokenID,
		Seller:      req.Seller,
		Price:       req.Price,
		Status:      "active",
		TxHash:      req.TxHash,
		ListedAt:    time.Now(),
	}

	if err := s.repo.Create(listing); err != nil {
		return nil, fmt.Errorf("failed to create listing: %w", err)
	}

	return s.toResponse(listing), nil
}

// GetListing 获取挂单
func (s *ListingService) GetListing(ctx context.Context, id uint) (*ListingResponse, error) {
	listing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get listing: %w", err)
	}

	return s.toResponse(listing), nil
}

// GetActiveListings 获取活跃挂单
func (s *ListingService) GetActiveListings(ctx context.Context, page, pageSize int) ([]*ListingResponse, int64, error) {
	listings, total, err := s.repo.GetActiveListings(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get active listings: %w", err)
	}

	responses := make([]*ListingResponse, len(listings))
	for i, listing := range listings {
		responses[i] = s.toResponse(&listing)
	}

	return responses, total, nil
}

// GetUserListings 获取用户挂单
func (s *ListingService) GetUserListings(ctx context.Context, address string, page, pageSize int) ([]*ListingResponse, int64, error) {
	listings, total, err := s.repo.GetBySellerPaginated(address, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user listings: %w", err)
	}

	responses := make([]*ListingResponse, len(listings))
	for i, listing := range listings {
		responses[i] = s.toResponse(&listing)
	}

	return responses, total, nil
}

// CancelListing 取消挂单
func (s *ListingService) CancelListing(ctx context.Context, id uint, seller string) error {
	listing, err := s.repo.GetByID(id)
	if err != nil {
		return fmt.Errorf("failed to get listing: %w", err)
	}

	if listing.Seller != seller {
		return fmt.Errorf("unauthorized: not the seller")
	}

	if listing.Status != "active" {
		return fmt.Errorf("listing is not active")
	}

	if err := s.repo.UpdateStatus(id, "cancelled"); err != nil {
		return fmt.Errorf("failed to cancel listing: %w", err)
	}

	return nil
}

// UpdateFromEvent 从区块链事件更新挂单
func (s *ListingService) UpdateFromEvent(event *blockchain.MarketItemCreatedEvent) error {
	listing := &repository.Listing{
		ItemID:      event.ItemId.Uint64(),
		NFTContract: event.NftContract.Hex(),
		TokenID:     event.TokenId.String(),
		Seller:      event.Seller.Hex(),
		Price:       event.Price.String(),
		Status:      "active",
		ListedAt:    time.Now(),
	}

	// 使用 CreateIfNotExists 防止并发重复插入
	return s.repo.CreateIfNotExists(listing)
}

// GetMarketStats 获取市场统计
func (s *ListingService) GetMarketStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 活跃挂单数量
	activeCount, err := s.repo.CountActiveListings()
	if err != nil {
		return nil, fmt.Errorf("failed to count active listings: %w", err)
	}
	stats["active_listings"] = activeCount

	// 总挂单数量
	totalCount, err := s.repo.CountTotalListings()
	if err != nil {
		return nil, fmt.Errorf("failed to count total listings: %w", err)
	}
	stats["total_listings"] = totalCount

	// 已售出数量
	stats["sold_listings"] = totalCount - activeCount

	// 总交易额
	totalVolume, err := s.repo.GetTotalVolume()
	if err != nil {
		return nil, fmt.Errorf("failed to get total volume: %w", err)
	}
	stats["total_volume"] = totalVolume

	// 平均价格
	avgPrice, err := s.repo.GetAveragePrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get average price: %w", err)
	}
	stats["average_price"] = avgPrice

	// 最低价格
	minPrice, err := s.repo.GetMinPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get min price: %w", err)
	}
	stats["floor_price"] = minPrice

	// 最高价格
	maxPrice, err := s.repo.GetMaxPrice()
	if err != nil {
		return nil, fmt.Errorf("failed to get max price: %w", err)
	}
	stats["ceiling_price"] = maxPrice

	return stats, nil
}

// toResponse 转换为响应对象
func (s *ListingService) toResponse(listing *repository.Listing) *ListingResponse {
	return &ListingResponse{
		ID:          listing.ID,
		ItemID:      listing.ItemID,
		NFTContract: listing.NFTContract,
		TokenID:     listing.TokenID,
		Seller:      listing.Seller,
		Price:       listing.Price,
		Status:      listing.Status,
		ListedAt:    listing.ListedAt,
		CreatedAt:   listing.CreatedAt,
	}
}
