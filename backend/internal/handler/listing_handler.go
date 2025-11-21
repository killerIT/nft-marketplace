package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaomait/backend/internal/service"
)

// ListingHandler 挂单处理器
type ListingHandler struct {
	service *service.ListingService
}

// NewListingHandler 创建挂单处理器
func NewListingHandler(service *service.ListingService) *ListingHandler {
	return &ListingHandler{service: service}
}

// GetActiveListings 获取活跃挂单
// @Summary 获取活跃挂单列表
// @Tags Listing
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/listings [get]
func (h *ListingHandler) GetActiveListings(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	listings, total, err := h.service.GetActiveListings(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get active listings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": listings,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetListing 获取单个挂单
// @Summary 获取挂单详情
// @Tags Listing
// @Param id path int true "Listing ID"
// @Success 200 {object} service.ListingResponse
// @Router /api/v1/listings/{id} [get]
func (h *ListingHandler) GetListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid listing ID",
		})
		return
	}

	listing, err := h.service.GetListing(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Listing not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": listing,
	})
}

// CreateListing 创建挂单
// @Summary 创建挂单
// @Tags Listing
// @Accept json
// @Param listing body service.CreateListingRequest true "挂单信息"
// @Success 201 {object} service.ListingResponse
// @Router /api/v1/listings [post]
func (h *ListingHandler) CreateListing(c *gin.Context) {
	var req service.CreateListingRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	listing, err := h.service.CreateListing(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create listing",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    listing,
		"message": "Listing created successfully",
	})
}

// CancelListing 取消挂单
// @Summary 取消挂单
// @Tags Listing
// @Param id path int true "Listing ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/listings/{id} [delete]
func (h *ListingHandler) CancelListing(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid listing ID",
		})
		return
	}

	// TODO: 从 JWT 或请求中获取用户地址
	seller := c.GetHeader("X-User-Address")
	if seller == "" {
		c.JSON(http.StatusUnauthorized, gin.H{
			"error": "User address is required",
		})
		return
	}

	if err := h.service.CancelListing(c.Request.Context(), uint(id), seller); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to cancel listing",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Listing cancelled successfully",
	})
}

// GetUserListings 获取用户的挂单
// @Summary 获取用户的挂单
// @Tags Listing
// @Param address path string true "用户地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/listings/user/{address} [get]
func (h *ListingHandler) GetUserListings(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Address is required",
		})
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	listings, total, err := h.service.GetUserListings(c.Request.Context(), address, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user listings",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": listings,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// SearchListings 搜索挂单
// @Summary 搜索挂单
// @Tags Listing
// @Param contract query string false "合约地址"
// @Param min_price query string false "最低价格"
// @Param max_price query string false "最高价格"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/listings/search [get]
func (h *ListingHandler) SearchListings(c *gin.Context) {
	contract := c.Query("contract")
	minPrice := c.Query("min_price")
	maxPrice := c.Query("max_price")

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	// TODO: 实现搜索逻辑
	c.JSON(http.StatusOK, gin.H{
		"data": []interface{}{},
		"filters": gin.H{
			"contract":  contract,
			"min_price": minPrice,
			"max_price": maxPrice,
		},
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       0,
			"total_pages": 0,
		},
	})
}

// GetMarketStats 获取市场统计
// @Summary 获取市场统计信息
// @Tags Stats
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/stats [get]
func (h *ListingHandler) GetMarketStats(c *gin.Context) {
	stats, err := h.service.GetMarketStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get market stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}

// GetCollectionStats 获取系列统计
// @Summary 获取系列统计信息
// @Tags Stats
// @Param address path string true "合约地址"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/stats/collections/{address} [get]
func (h *ListingHandler) GetCollectionStats(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Contract address is required",
		})
		return
	}

	// TODO: 实现系列统计逻辑
	c.JSON(http.StatusOK, gin.H{
		"data": gin.H{
			"contract_address": address,
			"total_items":      0,
			"active_listings":  0,
			"floor_price":      "0",
			"total_volume":     "0",
			"owners":           0,
		},
	})
}
