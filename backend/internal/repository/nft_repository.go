package repository

import (
	"time"

	"gorm.io/gorm"
)

// NFT NFT 模型
type NFT struct {
	ID              uint      `gorm:"primaryKey" json:"id"`
	ContractAddress string    `gorm:"index;not null" json:"contract_address"`
	TokenID         string    `gorm:"index;not null" json:"token_id"`
	Owner           string    `gorm:"index;not null" json:"owner"`
	Creator         string    `gorm:"index" json:"creator"`
	Name            string    `json:"name"`
	Description     string    `json:"description"`
	ImageURL        string    `json:"image_url"`
	MetadataURI     string    `json:"metadata_uri"`
	Metadata        string    `gorm:"type:jsonb" json:"metadata"`           // JSON 字符串
	Status          string    `gorm:"index;default:'active'" json:"status"` // active, burned, transferred
	ViewCount       int64     `gorm:"default:0" json:"view_count"`
	LikeCount       int64     `gorm:"default:0" json:"like_count"`
	MintedAt        time.Time `json:"minted_at"`
	CreatedAt       time.Time `json:"created_at"`
	UpdatedAt       time.Time `json:"updated_at"`
}

// TableName 指定表名
func (NFT) TableName() string {
	return "nfts"
}

// NFTRepository NFT 仓储
type NFTRepository struct {
	db *gorm.DB
}

// NewNFTRepository 创建 NFT 仓储
func NewNFTRepository(db *gorm.DB) *NFTRepository {
	return &NFTRepository{db: db}
}

// Create 创建 NFT
func (r *NFTRepository) Create(nft *NFT) error {
	return r.db.Create(nft).Error
}

// GetByID 根据 ID 获取 NFT
func (r *NFTRepository) GetByID(id uint) (*NFT, error) {
	var nft NFT
	err := r.db.First(&nft, id).Error
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

// GetByContractAndToken 根据合约地址和 Token ID 获取 NFT
func (r *NFTRepository) GetByContractAndToken(contractAddress, tokenID string) (*NFT, error) {
	var nft NFT
	err := r.db.Where("contract_address = ? AND token_id = ?", contractAddress, tokenID).First(&nft).Error
	if err != nil {
		return nil, err
	}
	return &nft, nil
}

// GetByOwner 根据所有者获取 NFT 列表
func (r *NFTRepository) GetByOwner(owner string, page, pageSize int) ([]NFT, int64, error) {
	var nfts []NFT
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&NFT{}).Where("owner = ? AND status = ?", owner, "active").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("owner = ? AND status = ?", owner, "active").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&nfts).Error

	if err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// GetByContract 根据合约地址获取 NFT 列表
func (r *NFTRepository) GetByContract(contractAddress string, page, pageSize int) ([]NFT, int64, error) {
	var nfts []NFT
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&NFT{}).Where("contract_address = ? AND status = ?", contractAddress, "active").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("contract_address = ? AND status = ?", contractAddress, "active").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&nfts).Error

	if err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// GetAll 获取所有 NFT（分页）
func (r *NFTRepository) GetAll(page, pageSize int) ([]NFT, int64, error) {
	var nfts []NFT
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&NFT{}).Where("status = ?", "active").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("status = ?", "active").
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&nfts).Error

	if err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// Update 更新 NFT
func (r *NFTRepository) Update(nft *NFT) error {
	return r.db.Save(nft).Error
}

// UpdateOwner 更新所有者
func (r *NFTRepository) UpdateOwner(id uint, newOwner string) error {
	return r.db.Model(&NFT{}).Where("id = ?", id).Update("owner", newOwner).Error
}

// IncrementViewCount 增加浏览次数
func (r *NFTRepository) IncrementViewCount(id uint) error {
	return r.db.Model(&NFT{}).Where("id = ?", id).UpdateColumn("view_count", gorm.Expr("view_count + ?", 1)).Error
}

// IncrementLikeCount 增加点赞次数
func (r *NFTRepository) IncrementLikeCount(id uint) error {
	return r.db.Model(&NFT{}).Where("id = ?", id).UpdateColumn("like_count", gorm.Expr("like_count + ?", 1)).Error
}

// DecrementLikeCount 减少点赞次数
func (r *NFTRepository) DecrementLikeCount(id uint) error {
	return r.db.Model(&NFT{}).Where("id = ?", id).Where("like_count > ?", 0).UpdateColumn("like_count", gorm.Expr("like_count - ?", 1)).Error
}

// Delete 删除 NFT（软删除，更新状态）
func (r *NFTRepository) Delete(id uint) error {
	return r.db.Model(&NFT{}).Where("id = ?", id).Update("status", "burned").Error
}

// Search 搜索 NFT
func (r *NFTRepository) Search(query string, page, pageSize int) ([]NFT, int64, error) {
	var nfts []NFT
	var total int64

	offset := (page - 1) * pageSize

	searchQuery := "%" + query + "%"

	// 计算总数
	if err := r.db.Model(&NFT{}).
		Where("status = ? AND (name ILIKE ? OR description ILIKE ?)", "active", searchQuery, searchQuery).
		Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("status = ? AND (name ILIKE ? OR description ILIKE ?)", "active", searchQuery, searchQuery).
		Order("created_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&nfts).Error

	if err != nil {
		return nil, 0, err
	}

	return nfts, total, nil
}

// GetTrending 获取热门 NFT（按浏览量和点赞数）
func (r *NFTRepository) GetTrending(limit int) ([]NFT, error) {
	var nfts []NFT
	err := r.db.Where("status = ?", "active").
		Order("(view_count + like_count * 2) DESC").
		Limit(limit).
		Find(&nfts).Error
	return nfts, err
}

// CountByOwner 统计用户拥有的 NFT 数量
func (r *NFTRepository) CountByOwner(owner string) (int64, error) {
	var count int64
	err := r.db.Model(&NFT{}).Where("owner = ? AND status = ?", owner, "active").Count(&count).Error
	return count, err
}

// CountByContract 统计合约的 NFT 数量
func (r *NFTRepository) CountByContract(contractAddress string) (int64, error) {
	var count int64
	err := r.db.Model(&NFT{}).Where("contract_address = ? AND status = ?", contractAddress, "active").Count(&count).Error
	return count, err
}
