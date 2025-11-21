package service

import (
	"context"
	"fmt"
	"time"

	"github.com/xiaomait/backend/internal/blockchain"
	"github.com/xiaomait/backend/internal/repository"
)

// TransactionService 交易服务
type TransactionService struct {
	repo     *repository.TransactionRepository
	bcClient *blockchain.Client
}

// NewTransactionService 创建交易服务
func NewTransactionService(repo *repository.TransactionRepository, bcClient *blockchain.Client) *TransactionService {
	return &TransactionService{
		repo:     repo,
		bcClient: bcClient,
	}
}

// TransactionResponse 交易响应
type TransactionResponse struct {
	ID             uint      `json:"id"`
	TxHash         string    `json:"tx_hash"`
	BlockNumber    uint64    `json:"block_number"`
	BlockTimestamp time.Time `json:"block_timestamp"`
	TxType         string    `json:"tx_type"`
	ListingID      *uint     `json:"listing_id,omitempty"`
	NFTContract    string    `json:"nft_contract"`
	TokenID        string    `json:"token_id"`
	FromAddress    string    `json:"from_address"`
	ToAddress      string    `json:"to_address"`
	Value          string    `json:"value"`
	GasPrice       string    `json:"gas_price"`
	GasUsed        uint64    `json:"gas_used"`
	PlatformFee    string    `json:"platform_fee"`
	Status         string    `json:"status"`
	CreatedAt      time.Time `json:"created_at"`
}

// GetTransaction 获取交易
func (s *TransactionService) GetTransaction(ctx context.Context, txHash string) (*TransactionResponse, error) {
	tx, err := s.repo.GetByHash(txHash)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return s.toResponse(tx), nil
}

// GetTransactionByID 根据 ID 获取交易
func (s *TransactionService) GetTransactionByID(ctx context.Context, id uint) (*TransactionResponse, error) {
	tx, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("failed to get transaction: %w", err)
	}

	return s.toResponse(tx), nil
}

// GetTransactions 获取交易列表
func (s *TransactionService) GetTransactions(ctx context.Context, page, pageSize int) ([]*TransactionResponse, int64, error) {
	txs, total, err := s.repo.GetAll(page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get transactions: %w", err)
	}

	responses := make([]*TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = s.toResponse(&tx)
	}

	return responses, total, nil
}

// GetUserTransactions 获取用户的交易
func (s *TransactionService) GetUserTransactions(ctx context.Context, address string, page, pageSize int) ([]*TransactionResponse, int64, error) {
	txs, total, err := s.repo.GetByAddress(address, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get user transactions: %w", err)
	}

	responses := make([]*TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = s.toResponse(&tx)
	}

	return responses, total, nil
}

// GetNFTTransactions 获取 NFT 的交易历史
func (s *TransactionService) GetNFTTransactions(ctx context.Context, nftContract, tokenID string, page, pageSize int) ([]*TransactionResponse, int64, error) {
	txs, total, err := s.repo.GetByNFT(nftContract, tokenID, page, pageSize)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to get NFT transactions: %w", err)
	}

	responses := make([]*TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = s.toResponse(&tx)
	}

	return responses, total, nil
}

// GetRecentTransactions 获取最近的交易
func (s *TransactionService) GetRecentTransactions(ctx context.Context, limit int) ([]*TransactionResponse, error) {
	txs, err := s.repo.GetRecent(limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get recent transactions: %w", err)
	}

	responses := make([]*TransactionResponse, len(txs))
	for i, tx := range txs {
		responses[i] = s.toResponse(&tx)
	}

	return responses, nil
}

// RecordSale 记录销售事件
func (s *TransactionService) RecordSale(event *blockchain.MarketItemSoldEvent) error {
	// 检查是否已存在
	// existing, _ := s.repo.GetByHash(event.TxHash)
	// if existing != nil {
	// 	return nil // 已存在，跳过
	// }

	tx := &repository.Transaction{
		TxHash:         "", // 需要从事件中获取
		BlockNumber:    0,  // 需要从事件中获取
		BlockTimestamp: time.Now(),
		TxType:         "sale",
		FromAddress:    event.Buyer.Hex(),
		ToAddress:      event.Buyer.Hex(),
		Value:          event.Price.String(),
		ValueNumeric:   event.Price.String(),
		Status:         "confirmed",
	}

	return s.repo.Create(tx)
}

// GetTotalVolume 获取总交易额
func (s *TransactionService) GetTotalVolume(ctx context.Context) (string, error) {
	volume, err := s.repo.GetTotalVolume()
	if err != nil {
		return "0", fmt.Errorf("failed to get total volume: %w", err)
	}
	return volume, nil
}

// GetVolumeByContract 获取合约的交易额
func (s *TransactionService) GetVolumeByContract(ctx context.Context, nftContract string) (string, error) {
	volume, err := s.repo.GetVolumeByContract(nftContract)
	if err != nil {
		return "0", fmt.Errorf("failed to get volume by contract: %w", err)
	}
	return volume, nil
}

// GetTransactionStats 获取交易统计
func (s *TransactionService) GetTransactionStats(ctx context.Context) (map[string]interface{}, error) {
	stats := make(map[string]interface{})

	// 总交易数
	listCount, _ := s.repo.CountByType("list")
	saleCount, _ := s.repo.CountByType("sale")
	cancelCount, _ := s.repo.CountByType("cancel")

	stats["total_listings"] = listCount
	stats["total_sales"] = saleCount
	stats["total_cancellations"] = cancelCount
	stats["total_transactions"] = listCount + saleCount + cancelCount

	// 总交易额
	totalVolume, _ := s.repo.GetTotalVolume()
	stats["total_volume"] = totalVolume

	return stats, nil
}

// toResponse 转换为响应对象
func (s *TransactionService) toResponse(tx *repository.Transaction) *TransactionResponse {
	return &TransactionResponse{
		ID:             tx.ID,
		TxHash:         tx.TxHash,
		BlockNumber:    tx.BlockNumber,
		BlockTimestamp: tx.BlockTimestamp,
		TxType:         tx.TxType,
		ListingID:      tx.ListingID,
		NFTContract:    tx.NFTContract,
		TokenID:        tx.TokenID,
		FromAddress:    tx.FromAddress,
		ToAddress:      tx.ToAddress,
		Value:          tx.Value,
		GasPrice:       tx.GasPrice,
		GasUsed:        tx.GasUsed,
		PlatformFee:    tx.PlatformFee,
		Status:         tx.Status,
		CreatedAt:      tx.CreatedAt,
	}
}
