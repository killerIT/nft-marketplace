package repository

import (
	"time"

	"gorm.io/gorm"
)

// Transaction 交易模型
type Transaction struct {
	ID               uint      `gorm:"primaryKey" json:"id"`
	TxHash           string    `gorm:"uniqueIndex;not null" json:"tx_hash"`
	BlockNumber      uint64    `gorm:"index;not null" json:"block_number"`
	BlockTimestamp   time.Time `gorm:"index;not null" json:"block_timestamp"`
	TxType           string    `gorm:"index;not null" json:"tx_type"` // list, sale, cancel, transfer, mint
	ListingID        *uint     `gorm:"index" json:"listing_id"`
	NFTContract      string    `gorm:"index;not null" json:"nft_contract"`
	TokenID          string    `gorm:"index;not null" json:"token_id"`
	FromAddress      string    `gorm:"index;not null" json:"from_address"`
	ToAddress        string    `gorm:"index" json:"to_address"`
	Value            string    `json:"value"`
	ValueNumeric     string    `gorm:"type:numeric(78,0)" json:"value_numeric"`
	GasPrice         string    `json:"gas_price"`
	GasUsed          uint64    `json:"gas_used"`
	PlatformFee      string    `json:"platform_fee"`
	Status           string    `gorm:"default:'confirmed'" json:"status"` // pending, confirmed, failed
	LogIndex         int       `json:"log_index"`
	TransactionIndex int       `json:"transaction_index"`
	CreatedAt        time.Time `json:"created_at"`
	UpdatedAt        time.Time `json:"updated_at"`
}

// TableName 指定表名
func (Transaction) TableName() string {
	return "transactions"
}

// TransactionRepository 交易仓储
type TransactionRepository struct {
	db *gorm.DB
}

// NewTransactionRepository 创建交易仓储
func NewTransactionRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

// Create 创建交易记录
func (r *TransactionRepository) Create(tx *Transaction) error {
	return r.db.Create(tx).Error
}

// GetByHash 根据交易哈希获取交易
func (r *TransactionRepository) GetByHash(txHash string) (*Transaction, error) {
	var tx Transaction
	err := r.db.Where("tx_hash = ?", txHash).First(&tx).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetByID 根据 ID 获取交易
func (r *TransactionRepository) GetByID(id uint) (*Transaction, error) {
	var tx Transaction
	err := r.db.First(&tx, id).Error
	if err != nil {
		return nil, err
	}
	return &tx, nil
}

// GetByAddress 根据地址获取交易（发送或接收）
func (r *TransactionRepository) GetByAddress(address string, page, pageSize int) ([]Transaction, int64, error) {
	var txs []Transaction
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Transaction{}).
		Where("from_address = ? OR to_address = ?", address, address).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("from_address = ? OR to_address = ?", address, address).
		Order("block_timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&txs).Error

	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetByNFT 根据 NFT 获取交易历史
func (r *TransactionRepository) GetByNFT(nftContract, tokenID string, page, pageSize int) ([]Transaction, int64, error) {
	var txs []Transaction
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Transaction{}).
		Where("nft_contract = ? AND token_id = ?", nftContract, tokenID).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("nft_contract = ? AND token_id = ?", nftContract, tokenID).
		Order("block_timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&txs).Error

	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetByType 根据类型获取交易
func (r *TransactionRepository) GetByType(txType string, page, pageSize int) ([]Transaction, int64, error) {
	var txs []Transaction
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Transaction{}).Where("tx_type = ?", txType).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("tx_type = ?", txType).
		Order("block_timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&txs).Error

	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetRecent 获取最近的交易
func (r *TransactionRepository) GetRecent(limit int) ([]Transaction, error) {
	var txs []Transaction
	err := r.db.Order("block_timestamp DESC").Limit(limit).Find(&txs).Error
	return txs, err
}

// GetAll 获取所有交易（分页）
func (r *TransactionRepository) GetAll(page, pageSize int) ([]Transaction, int64, error) {
	var txs []Transaction
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Transaction{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Order("block_timestamp DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&txs).Error

	if err != nil {
		return nil, 0, err
	}

	return txs, total, nil
}

// GetTotalVolume 获取总交易额
func (r *TransactionRepository) GetTotalVolume() (string, error) {
	var result struct {
		Total string
	}

	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(CAST(value_numeric AS NUMERIC)), 0) as total").
		Where("tx_type = ? AND status = ?", "sale", "confirmed").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Total, nil
}

// GetVolumeByContract 获取合约的交易额
func (r *TransactionRepository) GetVolumeByContract(nftContract string) (string, error) {
	var result struct {
		Total string
	}

	err := r.db.Model(&Transaction{}).
		Select("COALESCE(SUM(CAST(value_numeric AS NUMERIC)), 0) as total").
		Where("nft_contract = ? AND tx_type = ? AND status = ?", nftContract, "sale", "confirmed").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Total, nil
}

// CountByType 统计指定类型的交易数量
func (r *TransactionRepository) CountByType(txType string) (int64, error) {
	var count int64
	err := r.db.Model(&Transaction{}).Where("tx_type = ? AND status = ?", txType, "confirmed").Count(&count).Error
	return count, err
}

// GetDailyVolume 获取每日交易额（最近 N 天）
func (r *TransactionRepository) GetDailyVolume(days int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}

	query := `
		SELECT 
			DATE(block_timestamp) as date,
			COUNT(*) as tx_count,
			COALESCE(SUM(CAST(value_numeric AS NUMERIC)), 0) as volume
		FROM transactions
		WHERE tx_type = 'sale' 
		AND status = 'confirmed'
		AND block_timestamp >= NOW() - INTERVAL '? days'
		GROUP BY DATE(block_timestamp)
		ORDER BY date DESC
	`

	err := r.db.Raw(query, days).Scan(&results).Error
	return results, err
}

// Update 更新交易
func (r *TransactionRepository) Update(tx *Transaction) error {
	return r.db.Save(tx).Error
}

// UpdateStatus 更新交易状态
func (r *TransactionRepository) UpdateStatus(id uint, status string) error {
	return r.db.Model(&Transaction{}).Where("id = ?", id).Update("status", status).Error
}
