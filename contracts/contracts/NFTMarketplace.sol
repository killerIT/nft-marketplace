// SPDX-License-Identifier: MIT
pragma solidity ^0.8.20;

import "@openzeppelin/contracts/token/ERC721/ERC721.sol";
import "@openzeppelin/contracts/access/Ownable.sol";
// import "@openzeppelin/contracts/utils/Counters.sol";
import "@openzeppelin/contracts/security/ReentrancyGuard.sol"; // OpenZeppelin v4.x路径

/**
 * @title NFTMarketplace
 * @dev 生产级 NFT 市场合约
 */
contract NFTMarketplace is
    Ownable,
    ReentrancyGuard // 调整继承顺序
{
    // using Counters for Counters.Counter;

    // Counters.Counter private _itemIds;
    // Counters.Counter private _itemsSold;
    uint256 private _itemIds;
    uint256 private _itemsSold;

    uint256 public platformFee = 250; // 2.5% (基点)
    uint256 public constant DENOMINATOR = 10000;

    struct MarketItem {
        uint256 itemId;
        address nftContract;
        uint256 tokenId;
        address payable seller;
        address payable owner;
        uint256 price;
        bool sold;
        uint256 listedAt;
    }

    mapping(uint256 => MarketItem) private idToMarketItem;

    // 用户的活跃挂单
    mapping(address => uint256[]) private userActiveListings;

    event MarketItemCreated(
        uint256 indexed itemId,
        address indexed nftContract,
        uint256 indexed tokenId,
        address seller,
        uint256 price
    );

    event MarketItemSold(
        uint256 indexed itemId,
        address indexed buyer,
        uint256 price
    );

    event MarketItemCanceled(uint256 indexed itemId);

    event PlatformFeeUpdated(uint256 newFee);

    // 修复构造函数 - OpenZeppelin v4 兼容
    constructor() Ownable() {} // 移除参数，使用默认部署者为owner

    /**
     * @dev 上架 NFT
     */
    function createMarketItem(
        address nftContract,
        uint256 tokenId,
        uint256 price
    ) external nonReentrant returns (uint256) {
        require(price > 0, "Price must be greater than 0");
        require(
            IERC721(nftContract).ownerOf(tokenId) == msg.sender,
            "Not the owner"
        );
        require(
            IERC721(nftContract).isApprovedForAll(msg.sender, address(this)) ||
                IERC721(nftContract).getApproved(tokenId) == address(this),
            "Market not approved"
        );

        ++_itemIds;
        uint256 itemId = _itemIds;

        idToMarketItem[itemId] = MarketItem({
            itemId: itemId,
            nftContract: nftContract,
            tokenId: tokenId,
            seller: payable(msg.sender),
            owner: payable(address(0)),
            price: price,
            sold: false,
            listedAt: block.timestamp
        });

        userActiveListings[msg.sender].push(itemId);

        emit MarketItemCreated(itemId, nftContract, tokenId, msg.sender, price);

        return itemId;
    }

    /**
     * @dev 购买 NFT
     */
    function createMarketSale(uint256 itemId) external payable nonReentrant {
        MarketItem storage item = idToMarketItem[itemId];
        uint256 price = item.price;
        uint256 tokenId = item.tokenId;

        require(!item.sold, "Item already sold");
        require(msg.value == price, "Incorrect price");
        require(item.seller != msg.sender, "Cannot buy own item");

        // 计算平台费用
        uint256 fee = (price * platformFee) / DENOMINATOR;
        uint256 sellerProceeds = price - fee;

        // 转账给卖家 transfer() 方法有 2300 gas 限制，可能导致转账失败，且无法自定义错误处理。
        // item.seller.transfer(sellerProceeds);
        // 使用 call 代替 transfer
        (bool success, ) = item.seller.call{value: sellerProceeds}("");
        require(success, "Transfer to seller failed");
        // 转移 NFT
        IERC721(item.nftContract).safeTransferFrom(
            item.seller,
            msg.sender,
            tokenId
        );

        // 更新状态
        item.owner = payable(msg.sender);
        item.sold = true;
        ++_itemsSold;

        _removeFromUserListings(item.seller, itemId);

        emit MarketItemSold(itemId, msg.sender, price);
    }

    /**
     * @dev 取消挂单
     */
    function cancelMarketItem(uint256 itemId) external nonReentrant {
        MarketItem storage item = idToMarketItem[itemId];

        require(item.seller == msg.sender, "Not the seller");
        require(!item.sold, "Item already sold");

        item.sold = true; // 标记为已处理
        item.owner = item.seller;

        _removeFromUserListings(msg.sender, itemId);

        emit MarketItemCanceled(itemId);
    }

    /**
     * @dev 获取市场项详情
     */
    function getMarketItem(
        uint256 itemId
    ) external view returns (MarketItem memory) {
        return idToMarketItem[itemId];
    }

    /**
     * @dev 获取所有未售出的市场项
     */
    function fetchActiveItems() external view returns (MarketItem[] memory) {
        uint256 itemCount = _itemIds;
        uint256 unsoldCount = itemCount - _itemsSold;
        uint256 currentIndex = 0;

        MarketItem[] memory items = new MarketItem[](unsoldCount);

        for (uint256 i = 1; i <= itemCount; i++) {
            if (
                !idToMarketItem[i].sold && idToMarketItem[i].owner == address(0)
            ) {
                items[currentIndex] = idToMarketItem[i];
                currentIndex++;
            }
        }

        return items;
    }

    /**
     * @dev 获取用户的挂单
     */
    function fetchUserListings(
        address user
    ) external view returns (MarketItem[] memory) {
        uint256[] memory itemIds = userActiveListings[user];
        uint256 activeCount = 0;

        // 计算活跃挂单数量
        for (uint256 i = 0; i < itemIds.length; i++) {
            if (!idToMarketItem[itemIds[i]].sold) {
                activeCount++;
            }
        }

        MarketItem[] memory items = new MarketItem[](activeCount);
        uint256 currentIndex = 0;

        for (uint256 i = 0; i < itemIds.length; i++) {
            if (!idToMarketItem[itemIds[i]].sold) {
                items[currentIndex] = idToMarketItem[itemIds[i]];
                currentIndex++;
            }
        }

        return items;
    }

    /**
     * @dev 更新平台费用 (仅 owner)
     */
    function updatePlatformFee(uint256 newFee) external onlyOwner {
        require(newFee <= 1000, "Fee too high"); // 最高 10%
        platformFee = newFee;
        emit PlatformFeeUpdated(newFee);
    }

    /**
     * @dev 提取平台收益 (仅 owner)
     */
    function withdrawFees() external onlyOwner {
        uint256 balance = address(this).balance;
        require(balance > 0, "No funds to withdraw");
        payable(owner()).transfer(balance);
    }

    /**
     * @dev 从用户挂单列表中移除
     */
    function _removeFromUserListings(address user, uint256 itemId) private {
        uint256[] storage listings = userActiveListings[user];
        for (uint256 i = 0; i < listings.length; i++) {
            if (listings[i] == itemId) {
                listings[i] = listings[listings.length - 1];
                listings.pop();
                break;
            }
        }
    }

    /**
     * @dev 获取市场统计
     */
    function getMarketStats()
        external
        view
        returns (uint256 totalItems, uint256 soldItems, uint256 activeItems)
    {
        totalItems = _itemIds;
        soldItems = _itemsSold;
        activeItems = totalItems - soldItems;
    }
}
