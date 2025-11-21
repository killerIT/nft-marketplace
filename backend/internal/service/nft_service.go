package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/xiaomait/backend/internal/blockchain"
	"github.com/xiaomait/backend/internal/repository"
)

// NFTService NFT 服务
type NFTService struct {
	repo     *repository.NFTRepository
	bcClient *blockchain.Client
}

// NewNFTService 创建 NFT 服务
func NewNFTService(repo *repository.NFTRepository, bcClient *blockchain.Client) *NFTService {
	return &NFTService{
		repo:     repo,
		bcClient: bcClient,
	}
}

// CreateNFTRequest 创建 NFT 请求
type CreateNFTRequest struct {
	ContractAddress string                 `json:"contract_address" binding:"required"`
	TokenID         string                 `json:"token_id" binding:"required"`
	Owner           string                 `json:"owner" binding:"required"`
	Creator         string                 `json:"creator"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ImageURL        string                 `json:"image_url"`
	MetadataURI     string                 `json:"metadata_uri"`
	Metadata        map[string]interface{} `json:"metadata"`
}

// NFTResponse NFT 响应
type NFTResponse struct {
	ID              uint                   `json:"id"`
	ContractAddress string                 `json:"contract_address"`
	TokenID         string                 `json:"token_id"`
	Owner           string                 `json:"owner"`
	Creator         string                 `json:"creator"`
	Name            string                 `json:"name"`
	Description     string                 `json:"description"`
	ImageURL        string                 `json:"image_url"`
	MetadataURI     string                 `json:"metadata_uri"`
	Metadata        map[string]interface{} `json:"metadata"`
	Status          string                 `json:"status"`
	ViewCount       int64                  `json:"view_count"`
	LikeCount       int64                  `json:"like_count"`
	MintedAt        time.Time              `json:"minted_at"`
	CreatedAt       time.Time              `json:"created_at"`
	UpdatedAt       time.Time              `json:"updated_at"`
}

// CreateNFT 创建 NFT
func (s *NFTService) CreateNFT(ctx context.Context, req *CreateNFTRequest) (*NFTResponse, error) {
	// 检查是否已存在
	existing, _ := s.repo.GetByContractAndToken(req.ContractAddress, req.TokenID)
	if existing != nil {
		return nil, fmt.Errorf("NFT already exists")
	}

	// 序列化 metadata
	metadataJSON, err := json.Marshal(req.Metadata)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal metadata: %w", err)
	}

	nft := &repository.NFT{
		ContractAddress: req.ContractAddress,
		TokenID:         req.TokenID,
		Owner:           req.Owner,
		Creator:         req.Creator,
		Name:            req.Name,
		Description:     req.Description,
		ImageURL:        req.ImageURL,
		MetadataURI:     req.MetadataURI,
		Metadata:        string(metadataJSON),
		Status:          "active",
		MintedAt:        time.Now(),
	}

	if err := s.repo.Create(nft); err != nil {
		return nil, fmt.Errorf("failed to create NFT: %w", err)
	}

	return s.toResponse(nft), nil
}

// GetNFT 获取 NFT
func (s *NFTService) GetNFT(ctx context.Context, id uint) (*NFTResponse, error) {
	nft, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT: %w", err)
	}

	// 增加浏览次数
	go s.repo.IncrementViewCount(id)

	return s.toResponse(nft), nil
}

// GetNFTByContractAndToken 根据合约和 Token ID 获取 NFT
func (s *NFTService) GetNFTByContractAndToken(ctx context.Context, contractAddress, tokenID string) (*NFTResponse, error) {
	nft, err := s.repo.GetByContractAndToken(contractAddress, tokenID)
	if err != nil {
		return nil, fmt.Errorf("failed to get NFT: %w", err)
	}

	// 增加浏览次数
	go s.repo.IncrementViewCount(nft.ID)

	return s.toResponse(nft), nil
}

// GetNFTs 获取 NFT 列表
func (s *NFTService) GetNFTs(ctx context.Context, page, pageSize int) ([]*NFTResponse, int64, error) {
	nfts, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get NFTs: %w", err)
	}

	responses := make([]*NFTResponse, len(nfts))
	for i, nft := range nfts {
		responses[i] = s.toResponse(&nft)
	}

	return responses, total, nil
}

// GetUserNFTs 获取用户的 NFT
func (s *NFTService) GetUserNFTs(ctx context.Context, owner string, page, pageSize int) ([]*NFTResponse, int64, error) {
	nfts, total, err := s.repo.GetByOwner(owner, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user NFTs: %w", err)
	}

	responses := make([]*NFTResponse, len(nfts))
	for i, nft := range nfts {
		responses[i] = s.toResponse(&nft)
	}

	return responses, total, nil
}

// GetNFTsByContract 获取合约的 NFT
func (s *NFTService) GetNFTsByContract(ctx context.Context, contractAddress string, page, pageSize int) ([]*NFTResponse, int64, error) {
	nfts, total, err := s.repo.GetByContract(contractAddress, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get NFTs by contract: %w", err)
	}

	responses := make([]*NFTResponse, len(nfts))
	for i, nft := range nfts {
		responses[i] = s.toResponse(&nft)
	}

	return responses, total, nil
}

// SearchNFTs 搜索 NFT
func (s *NFTService) SearchNFTs(ctx context.Context, query string, page, pageSize int) ([]*NFTResponse, int64, error) {
	nfts, total, err := s.repo.Search(query, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to search NFTs: %w", err)
	}

	responses := make([]*NFTResponse, len(nfts))
	for i, nft := range nfts {
		responses[i] = s.toResponse(&nft)
	}

	return responses, total, nil
}

// GetTrendingNFTs 获取热门 NFT
func (s *NFTService) GetTrendingNFTs(ctx context.Context, limit int) ([]*NFTResponse, error) {
	nfts, err := s.repo.GetTrending(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get trending NFTs: %w", err)
	}

	responses := make([]*NFTResponse, len(nfts))
	for i, nft := range nfts {
		responses[i] = s.toResponse(&nft)
	}

	return responses, nil
}

// UpdateNFTOwner 更新 NFT 所有者
func (s *NFTService) UpdateNFTOwner(ctx context.Context, id uint, newOwner string) error {
	if err := s.repo.UpdateOwner(id, newOwner); err != nil {
		return fmt.Errorf("failed to update NFT owner: %w", err)
	}
	return nil
}

// LikeNFT 点赞 NFT
func (s *NFTService) LikeNFT(ctx context.Context, id uint) error {
	if err := s.repo.IncrementLikeCount(id); err != nil {
		return fmt.Errorf("failed to like NFT: %w", err)
	}
	return nil
}

// UnlikeNFT 取消点赞 NFT
func (s *NFTService) UnlikeNFT(ctx context.Context, id uint) error {
	if err := s.repo.DecrementLikeCount(id); err != nil {
		return fmt.Errorf("failed to unlike NFT: %w", err)
	}
	return nil
}

// toResponse 转换为响应对象
func (s *NFTService) toResponse(nft *repository.NFT) *NFTResponse {
	var metadata map[string]interface{}
	if nft.Metadata != "" {
		json.Unmarshal([]byte(nft.Metadata), &metadata)
	}

	return &NFTResponse{
		ID:              nft.ID,
		ContractAddress: nft.ContractAddress,
		TokenID:         nft.TokenID,
		Owner:           nft.Owner,
		Creator:         nft.Creator,
		Name:            nft.Name,
		Description:     nft.Description,
		ImageURL:        nft.ImageURL,
		MetadataURI:     nft.MetadataURI,
		Metadata:        metadata,
		Status:          nft.Status,
		ViewCount:       nft.ViewCount,
		LikeCount:       nft.LikeCount,
		MintedAt:        nft.MintedAt,
		CreatedAt:       nft.CreatedAt,
		UpdatedAt:       nft.UpdatedAt,
	}
}
