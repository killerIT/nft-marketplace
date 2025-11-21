package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaomait/backend/internal/service"
)

// NFTHandler NFT 处理器
type NFTHandler struct {
	service *service.NFTService
}

// NewNFTHandler 创建 NFT 处理器
func NewNFTHandler(service *service.NFTService) *NFTHandler {
	return &NFTHandler{service: service}
}

// GetNFTs 获取 NFT 列表
// @Summary 获取 NFT 列表
// @Tags NFT
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts [get]
func (h *NFTHandler) GetNFTs(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	nfts, total, err := h.service.GetNFTs(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get NFTs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nfts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetNFT 获取单个 NFT
// @Summary 获取 NFT 详情
// @Tags NFT
// @Param id path int true "NFT ID"
// @Success 200 {object} service.NFTResponse
// @Router /api/v1/nfts/{id} [get]
func (h *NFTHandler) GetNFT(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid NFT ID",
		})
		return
	}

	nft, err := h.service.GetNFT(c.Request.Context(), uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "NFT not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nft,
	})
}

// CreateNFT 创建 NFT
// @Summary 创建 NFT
// @Tags NFT
// @Accept json
// @Param nft body service.CreateNFTRequest true "NFT 信息"
// @Success 201 {object} service.NFTResponse
// @Router /api/v1/nfts [post]
func (h *NFTHandler) CreateNFT(c *gin.Context) {
	var req service.CreateNFTRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error":   "Invalid request body",
			"details": err.Error(),
		})
		return
	}

	nft, err := h.service.CreateNFT(c.Request.Context(), &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to create NFT",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"data":    nft,
		"message": "NFT created successfully",
	})
}

// GetUserNFTs 获取用户的 NFT
// @Summary 获取用户的 NFT
// @Tags NFT
// @Param address path string true "用户地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/user/{address} [get]
func (h *NFTHandler) GetUserNFTs(c *gin.Context) {
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

	nfts, total, err := h.service.GetUserNFTs(c.Request.Context(), address, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user NFTs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nfts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetNFTsByContract 获取合约的 NFT
// @Summary 获取合约的所有 NFT
// @Tags NFT
// @Param address path string true "合约地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/contract/{address} [get]
func (h *NFTHandler) GetNFTsByContract(c *gin.Context) {
	address := c.Param("address")
	if address == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Contract address is required",
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

	nfts, total, err := h.service.GetNFTsByContract(c.Request.Context(), address, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get NFTs by contract",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nfts,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// SearchNFTs 搜索 NFT
// @Summary 搜索 NFT
// @Tags NFT
// @Param q query string true "搜索关键词"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/search [get]
func (h *NFTHandler) SearchNFTs(c *gin.Context) {
	query := c.Query("q")
	if query == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Search query is required",
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

	nfts, total, err := h.service.SearchNFTs(c.Request.Context(), query, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to search NFTs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data":  nfts,
		"query": query,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetTrendingNFTs 获取热门 NFT
// @Summary 获取热门 NFT
// @Tags NFT
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/trending [get]
func (h *NFTHandler) GetTrendingNFTs(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	nfts, err := h.service.GetTrendingNFTs(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get trending NFTs",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": nfts,
	})
}

// LikeNFT 点赞 NFT
// @Summary 点赞 NFT
// @Tags NFT
// @Param id path int true "NFT ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/{id}/like [post]
func (h *NFTHandler) LikeNFT(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid NFT ID",
		})
		return
	}

	if err := h.service.LikeNFT(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to like NFT",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "NFT liked successfully",
	})
}

// UnlikeNFT 取消点赞 NFT
// @Summary 取消点赞 NFT
// @Tags NFT
// @Param id path int true "NFT ID"
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/nfts/{id}/unlike [post]
func (h *NFTHandler) UnlikeNFT(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid NFT ID",
		})
		return
	}

	if err := h.service.UnlikeNFT(c.Request.Context(), uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to unlike NFT",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "NFT unliked successfully",
	})
}
