package repository

import (
	"time"

	"gorm.io/gorm"
)

// Listing 挂单模型
type Listing struct {
	ID          uint      `gorm:"primaryKey" json:"id"`
	ItemID      uint64    `gorm:"uniqueIndex;not null" json:"item_id"`
	NFTContract string    `gorm:"index;not null" json:"nft_contract"`
	TokenID     string    `gorm:"index;not null" json:"token_id"`
	Seller     string    `gorm:"index;not null" json:"seller"`
	Price       string    `gorm:"not null" json:"price"`
	Status      string    `gorm:"index;not null;default:'active'" json:"status"` // active, sold, cancelled
	TxHash      string    `gorm:"index" json:"tx_hash"`
	ListedAt    time.Time `gorm:"not null" json:"listed_at"`
	SoldAt      *time.Time `json:"sold_at,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// ListingRepository 挂单仓储
type ListingRepository struct {
	db *gorm.DB
}

// NewListingRepository 创建挂单仓储
func NewListingRepository(db *gorm.DB) *ListingRepository {
	return &ListingRepository{db: db}
}

// Create 创建挂单
func (r *ListingRepository) Create(listing *Listing) error {
	return r.db.Create(listing).Error
}

// CreateIfNotExists 创建挂单（如果不存在）- 防止并发重复插入
func (r *ListingRepository) CreateIfNotExists(listing *Listing) error {
	// 使用 FirstOrCreate 来处理并发情况
	// 如果 item_id 已存在，则不插入；否则插入新记录
	result := r.db.Where("item_id = ?", listing.ItemID).FirstOrCreate(listing)
	return result.Error
}

// GetByID 根据 ID 获取挂单
func (r *ListingRepository) GetByID(id uint) (*Listing, error) {
	var listing Listing
	err := r.db.First(&listing, id).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

// GetByItemID 根据 ItemID 获取挂单
func (r *ListingRepository) GetByItemID(itemID uint64) (*Listing, error) {
	var listing Listing
	err := r.db.Where("item_id = ?", itemID).First(&listing).Error
	if err != nil {
		return nil, err
	}
	return &listing, nil
}

// GetActiveListings 获取活跃挂单（分页）
func (r *ListingRepository) GetActiveListings(page, pageSize int) ([]Listing, int64, error) {
	var listings []Listing
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Listing{}).Where("status = ?", "active").Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("status = ?", "active").
		Order("listed_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&listings).Error

	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

// GetBySeller 根据卖家获取挂单
func (r *ListingRepository) GetBySeller(seller string) ([]Listing, error) {
	var listings []Listing
	err := r.db.Where("seller = ?", seller).Order("listed_at DESC").Find(&listings).Error
	return listings, err
}

// GetBySellerPaginated 根据卖家获取挂单（分页）
func (r *ListingRepository) GetBySellerPaginated(seller string, page, pageSize int) ([]Listing, int64, error) {
	var listings []Listing
	var total int64

	offset := (page - 1) * pageSize

	// 计算总数
	if err := r.db.Model(&Listing{}).Where("seller = ?", seller).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := r.db.Where("seller = ?", seller).
		Order("listed_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&listings).Error

	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}

// UpdateStatus 更新状态
func (r *ListingRepository) UpdateStatus(id uint, status string) error {
	updates := map[string]interface{}{
		"status": status,
	}

	if status == "sold" {
		now := time.Now()
		updates["sold_at"] = &now
	}

	return r.db.Model(&Listing{}).Where("id = ?", id).Updates(updates).Error
}

// CountActiveListings 统计活跃挂单数量
func (r *ListingRepository) CountActiveListings() (int64, error) {
	var count int64
	err := r.db.Model(&Listing{}).Where("status = ?", "active").Count(&count).Error
	return count, err
}

// CountTotalListings 统计总挂单数量
func (r *ListingRepository) CountTotalListings() (int64, error) {
	var count int64
	err := r.db.Model(&Listing{}).Count(&count).Error
	return count, err
}

// GetTotalVolume 获取总交易额
func (r *ListingRepository) GetTotalVolume() (string, error) {
	var result struct {
		Total string
	}

	err := r.db.Model(&Listing{}).
		Select("COALESCE(SUM(CAST(price AS NUMERIC)), 0) as total").
		Where("status = ?", "sold").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Total, nil
}

// GetAveragePrice 获取平均价格
func (r *ListingRepository) GetAveragePrice() (string, error) {
	var result struct {
		Avg string
	}

	err := r.db.Model(&Listing{}).
		Select("COALESCE(AVG(CAST(price AS NUMERIC)), 0) as avg").
		Where("status = ?", "active").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Avg, nil
}

// GetMinPrice 获取最低价格（地板价）
func (r *ListingRepository) GetMinPrice() (string, error) {
	var result struct {
		Min string
	}

	err := r.db.Model(&Listing{}).
		Select("COALESCE(MIN(CAST(price AS NUMERIC)), 0) as min").
		Where("status = ?", "active").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Min, nil
}

// GetMaxPrice 获取最高价格
func (r *ListingRepository) GetMaxPrice() (string, error) {
	var result struct {
		Max string
	}

	err := r.db.Model(&Listing{}).
		Select("COALESCE(MAX(CAST(price AS NUMERIC)), 0) as max").
		Where("status = ?", "active").
		Scan(&result).Error

	if err != nil {
		return "0", err
	}

	return result.Max, nil
}

// GetRecentListings 获取最近挂单
func (r *ListingRepository) GetRecentListings(limit int) ([]Listing, error) {
	var listings []Listing
	err := r.db.Where("status = ?", "active").
		Order("listed_at DESC").
		Limit(limit).
		Find(&listings).Error
	return listings, err
}

// SearchListings 搜索挂单
func (r *ListingRepository) SearchListings(nftContract string, minPrice, maxPrice string, page, pageSize int) ([]Listing, int64, error) {
	var listings []Listing
	var total int64

	offset := (page - 1) * pageSize

	query := r.db.Model(&Listing{}).Where("status = ?", "active")

	if nftContract != "" {
		query = query.Where("nft_contract = ?", nftContract)
	}

	if minPrice != "" {
		query = query.Where("CAST(price AS NUMERIC) >= ?", minPrice)
	}

	if maxPrice != "" {
		query = query.Where("CAST(price AS NUMERIC) <= ?", maxPrice)
	}

	// 计算总数
	if err := query.Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// 获取数据
	err := query.Order("listed_at DESC").
		Offset(offset).
		Limit(pageSize).
		Find(&listings).Error

	if err != nil {
		return nil, 0, err
	}

	return listings, total, nil
}