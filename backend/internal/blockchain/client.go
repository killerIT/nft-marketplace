package blockchain

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/big"
	"strings"
	"time"
	"unicode"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/ethclient"
)

// MarketItemCreatedEvent 市场项创建事件
type MarketItemCreatedEvent struct {
	ItemId      *big.Int
	NftContract common.Address
	TokenId     *big.Int
	Seller      common.Address
	Price       *big.Int
}

// MarketItemSoldEvent 市场项售出事件
type MarketItemSoldEvent struct {
	ItemId *big.Int
	Buyer  common.Address
	Price  *big.Int
}

// Client 区块链客户端
type Client struct {
	ethClient       *ethclient.Client
	marketplaceAddr common.Address
	contractABI     abi.ABI
}

// 合约 ABI (简化版本)
const marketplaceABI = `[
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "itemId", "type": "uint256"},
			{"indexed": true, "name": "nftContract", "type": "address"},
			{"indexed": true, "name": "tokenId", "type": "uint256"},
			{"indexed": false, "name": "seller", "type": "address"},
			{"indexed": false, "name": "price", "type": "uint256"}
		],
		"name": "MarketItemCreated",
		"type": "event"
	},
	{
		"anonymous": false,
		"inputs": [
			{"indexed": true, "name": "itemId", "type": "uint256"},
			{"indexed": true, "name": "buyer", "type": "address"},
			{"indexed": false, "name": "price", "type": "uint256"}
		],
		"name": "MarketItemSold",
		"type": "event"
	},
	{
		"inputs": [
			{"name": "itemId", "type": "uint256"}
		],
		"name": "getMarketItem",
		"outputs": [
			{
				"components": [
					{"name": "itemId", "type": "uint256"},
					{"name": "nftContract", "type": "address"},
					{"name": "tokenId", "type": "uint256"},
					{"name": "seller", "type": "address"},
					{"name": "owner", "type": "address"},
					{"name": "price", "type": "uint256"},
					{"name": "sold", "type": "bool"},
					{"name": "listedAt", "type": "uint256"}
				],
				"name": "",
				"type": "tuple"
			}
		],
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [],
		"name": "fetchActiveItems",
		"outputs": [
			{
				"components": [
					{"name": "itemId", "type": "uint256"},
					{"name": "nftContract", "type": "address"},
					{"name": "tokenId", "type": "uint256"},
					{"name": "seller", "type": "address"},
					{"name": "owner", "type": "address"},
					{"name": "price", "type": "uint256"},
					{"name": "sold", "type": "bool"},
					{"name": "listedAt", "type": "uint256"}
				],
				"name": "",
				"type": "tuple[]"
			}
		],
		"stateMutability": "view",
		"type": "function"
	}
]`

// NewClient 创建新的区块链客户端
func NewClient(rpcURL, marketplaceAddress string) (*Client, error) {
	client, err := ethclient.Dial(rpcURL)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Ethereum node: %w", err)
	}

	contractABI, err := abi.JSON(strings.NewReader(marketplaceABI))
	if err != nil {
		return nil, fmt.Errorf("failed to parse contract ABI: %w", err)
	}

	return &Client{
		ethClient:       client,
		marketplaceAddr: common.HexToAddress(marketplaceAddress),
		contractABI:     contractABI,
	}, nil
}

// GetBlockNumber 获取当前区块号
func (c *Client) GetBlockNumber(ctx context.Context) (uint64, error) {
	return c.ethClient.BlockNumber(ctx)
}

// GetMarketItem 获取市场项详情
func (c *Client) GetMarketItem(ctx context.Context, itemId *big.Int) (map[string]interface{}, error) {
	data, err := c.contractABI.Pack("getMarketItem", itemId)
	if err != nil {
		return nil, fmt.Errorf("failed to pack data: %w", err)
	}

	msg := ethereum.CallMsg{
		To:   &c.marketplaceAddr,
		Data: data,
	}

	result, err := c.ethClient.CallContract(ctx, msg, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to call contract: %w", err)
	}
	// 使用 UnpackIntoMap 方法
	resultMap := make(map[string]interface{})
	err = c.contractABI.UnpackIntoMap(resultMap, "getMarketItem", result)
	if err != nil {
		return nil, fmt.Errorf("failed to unpack result: %w", err)
	}

	log.Printf("Market item data: %+v", resultMap)

	// 解析特殊的 resultMap 结构
	var itemData interface{}
	for _, value := range resultMap {
		itemData = value
		break // 只取第一个值
	}

	log.Printf("Market itemData: %+v", itemData)

	// 定义结构体类型
	/*type MarketItemStruct struct {
		ItemId      *big.Int       `json:"itemId"`
		NftContract common.Address `json:"nftContract"`
		TokenId     *big.Int       `json:"tokenId"`
		Seller      common.Address `json:"seller"`
		Owner       common.Address `json:"owner"`
		Price       *big.Int       `json:"price"`
		Sold        bool           `json:"sold"`
		ListedAt    *big.Int       `json:"listedAt"`
	}*/
	return ConvertViaJSON(itemData)

}

// 方法3：JSON 方式（最通用）
func ConvertViaJSON(itemData interface{}) (map[string]interface{}, error) {
	jsonBytes, err := json.Marshal(itemData)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal: %w", err)
	}

	var rawMap map[string]interface{}
	err = json.Unmarshal(jsonBytes, &rawMap)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal: %w", err)
	}

	// 转换为蛇形命名
	result := make(map[string]interface{})
	for key, value := range rawMap {
		result[key] = value
	}

	return result, nil
}

// 辅助函数：将驼峰命名转为蛇形命名
func toSnakeCase(str string) string {
	var result []rune
	for i, r := range str {
		if i > 0 && 'A' <= r && r <= 'Z' {
			result = append(result, '_')
		}
		result = append(result, unicode.ToLower(r))
	}
	return string(result)
}

// ListenMarketItemCreated 监听 MarketItemCreated 事件（带重连机制）
func (c *Client) ListenMarketItemCreated(ctx context.Context) <-chan *MarketItemCreatedEvent {
	eventChan := make(chan *MarketItemCreatedEvent)

	go func() {
		defer close(eventChan)

		query := ethereum.FilterQuery{
			Addresses: []common.Address{c.marketplaceAddr},
			Topics:    [][]common.Hash{{c.contractABI.Events["MarketItemCreated"].ID}},
		}

		for {
			// 检查 context 是否已取消
			select {
			case <-ctx.Done():
				log.Println("MarketItemCreated listener stopped")
				return
			default:
			}

			logs := make(chan types.Log)
			sub, err := c.ethClient.SubscribeFilterLogs(ctx, query, logs)
			if err != nil {
				log.Printf("Failed to subscribe to MarketItemCreated logs, retrying in 5s: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Println("MarketItemCreated listener connected")

			// 处理事件循环
		eventLoop:
			for {
				select {
				case <-ctx.Done():
					sub.Unsubscribe()
					log.Println("MarketItemCreated listener stopped")
					return
				case err := <-sub.Err():
					log.Printf("MarketItemCreated subscription error: %v, reconnecting...", err)
					sub.Unsubscribe()
					time.Sleep(5 * time.Second)
					break eventLoop // 退出内层循环，重新订阅
				case vLog := <-logs:
					event := &MarketItemCreatedEvent{}
					err := c.contractABI.UnpackIntoInterface(event, "MarketItemCreated", vLog.Data)
					if err != nil {
						log.Printf("Failed to unpack MarketItemCreated event: %v", err)
						continue
					}

					// 解析 indexed 参数
					event.ItemId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
					event.NftContract = common.BytesToAddress(vLog.Topics[2].Bytes())
					event.TokenId = new(big.Int).SetBytes(vLog.Topics[3].Bytes())

					eventChan <- event
				}
			}
		}
	}()

	return eventChan
}

// ListenMarketItemSold 监听 MarketItemSold 事件（带重连机制）
func (c *Client) ListenMarketItemSold(ctx context.Context) <-chan *MarketItemSoldEvent {
	eventChan := make(chan *MarketItemSoldEvent)

	go func() {
		defer close(eventChan)

		query := ethereum.FilterQuery{
			Addresses: []common.Address{c.marketplaceAddr},
			Topics:    [][]common.Hash{{c.contractABI.Events["MarketItemSold"].ID}},
		}

		for {
			// 检查 context 是否已取消
			select {
			case <-ctx.Done():
				log.Println("MarketItemSold listener stopped")
				return
			default:
			}

			logs := make(chan types.Log)
			sub, err := c.ethClient.SubscribeFilterLogs(ctx, query, logs)
			if err != nil {
				log.Printf("Failed to subscribe to MarketItemSold logs, retrying in 5s: %v", err)
				time.Sleep(5 * time.Second)
				continue
			}

			log.Println("MarketItemSold listener connected")

			// 处理事件循环
		eventLoop:
			for {
				select {
				case <-ctx.Done():
					sub.Unsubscribe()
					log.Println("MarketItemSold listener stopped")
					return
				case err := <-sub.Err():
					log.Printf("MarketItemSold subscription error: %v, reconnecting...", err)
					sub.Unsubscribe()
					time.Sleep(5 * time.Second)
					break eventLoop // 退出内层循环，重新订阅
				case vLog := <-logs:
					event := &MarketItemSoldEvent{}
					err := c.contractABI.UnpackIntoInterface(event, "MarketItemSold", vLog.Data)
					if err != nil {
						log.Printf("Failed to unpack MarketItemSold event: %v", err)
						continue
					}

					// 解析 indexed 参数
					event.ItemId = new(big.Int).SetBytes(vLog.Topics[1].Bytes())
					event.Buyer = common.BytesToAddress(vLog.Topics[2].Bytes())

					eventChan <- event
				}
			}
		}
	}()

	return eventChan
}

// GetTransactionReceipt 获取交易回执
func (c *Client) GetTransactionReceipt(ctx context.Context, txHash common.Hash) (*types.Receipt, error) {
	return c.ethClient.TransactionReceipt(ctx, txHash)
}

// Close 关闭客户端
func (c *Client) Close() {
	c.ethClient.Close()
}
