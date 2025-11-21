package handler

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/xiaomait/backend/internal/service"
)

// TransactionHandler 交易处理器
type TransactionHandler struct {
	service *service.TransactionService
}

// NewTransactionHandler 创建交易处理器
func NewTransactionHandler(service *service.TransactionService) *TransactionHandler {
	return &TransactionHandler{service: service}
}

// GetTransactions 获取交易列表
// @Summary 获取交易列表
// @Tags Transaction
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions [get]
func (h *TransactionHandler) GetTransactions(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))

	if page < 1 {
		page = 1
	}
	if pageSize < 1 || pageSize > 100 {
		pageSize = 20
	}

	transactions, total, err := h.service.GetTransactions(c.Request.Context(), page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetTransaction 获取单个交易
// @Summary 获取交易详情
// @Tags Transaction
// @Param hash path string true "交易哈希"
// @Success 200 {object} service.TransactionResponse
// @Router /api/v1/transactions/{hash} [get]
func (h *TransactionHandler) GetTransaction(c *gin.Context) {
	txHash := c.Param("hash")
	if txHash == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Transaction hash is required",
		})
		return
	}

	transaction, err := h.service.GetTransaction(c.Request.Context(), txHash)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{
			"error":   "Transaction not found",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transaction,
	})
}

// GetUserTransactions 获取用户的交易
// @Summary 获取用户的交易历史
// @Tags Transaction
// @Param address path string true "用户地址"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/user/{address} [get]
func (h *TransactionHandler) GetUserTransactions(c *gin.Context) {
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

	transactions, total, err := h.service.GetUserTransactions(c.Request.Context(), address, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get user transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetNFTTransactions 获取 NFT 的交易历史
// @Summary 获取 NFT 的交易历史
// @Tags Transaction
// @Param contract path string true "合约地址"
// @Param tokenId path string true "Token ID"
// @Param page query int false "页码" default(1)
// @Param page_size query int false "每页数量" default(20)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/nft/{contract}/{tokenId} [get]
func (h *TransactionHandler) GetNFTTransactions(c *gin.Context) {
	contract := c.Param("contract")
	tokenID := c.Param("tokenId")

	if contract == "" || tokenID == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Contract address and token ID are required",
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

	transactions, total, err := h.service.GetNFTTransactions(c.Request.Context(), contract, tokenID, page, pageSize)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get NFT transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
		"nft": gin.H{
			"contract": contract,
			"token_id": tokenID,
		},
		"pagination": gin.H{
			"page":        page,
			"page_size":   pageSize,
			"total":       total,
			"total_pages": (total + int64(pageSize) - 1) / int64(pageSize),
		},
	})
}

// GetRecentTransactions 获取最近的交易
// @Summary 获取最近的交易
// @Tags Transaction
// @Param limit query int false "数量限制" default(10)
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/recent [get]
func (h *TransactionHandler) GetRecentTransactions(c *gin.Context) {
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))
	if limit < 1 || limit > 100 {
		limit = 10
	}

	transactions, err := h.service.GetRecentTransactions(c.Request.Context(), limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get recent transactions",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": transactions,
	})
}

// GetTransactionStats 获取交易统计
// @Summary 获取交易统计信息
// @Tags Transaction
// @Success 200 {object} map[string]interface{}
// @Router /api/v1/transactions/stats [get]
func (h *TransactionHandler) GetTransactionStats(c *gin.Context) {
	stats, err := h.service.GetTransactionStats(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error":   "Failed to get transaction stats",
			"details": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": stats,
	})
}
