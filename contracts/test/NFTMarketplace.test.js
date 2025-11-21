const { expect } = require("chai");
const { ethers } = require("hardhat");

describe("NFTMarketplace", function () {
  let marketplace;
  let nft;
  let owner;
  let seller;
  let buyer;
  const PLATFORM_FEE = 250; // 2.5%

  beforeEach(async function () {
    [owner, seller, buyer] = await ethers.getSigners();

    // 部署 NFT 合约
    const NFT = await ethers.getContractFactory("NFT");
    nft = await NFT.deploy("TestNFT", "TNFT");
    await nft.waitForDeployment();

    // 部署市场合约
    const Marketplace = await ethers.getContractFactory("NFTMarketplace");
    marketplace = await Marketplace.deploy(owner.address);
    await marketplace.waitForDeployment();

    // Mint NFT 给 seller
    await nft.connect(seller).mint("https://example.com/token/1");
  });

  describe("Deployment", function () {
    it("Should set the correct owner", async function () {
      expect(await marketplace.owner()).to.equal(owner.address);
    });

    it("Should set the correct platform fee", async function () {
      expect(await marketplace.platformFee()).to.equal(PLATFORM_FEE);
    });
  });

  describe("Listing", function () {
    const tokenId = 1;
    const price = ethers.parseEther("1");

    it("Should create a market item", async function () {
      // Approve marketplace
      await nft.connect(seller).approve(marketplace.target, tokenId);

      // Create listing
      await expect(
        marketplace.connect(seller).createMarketItem(nft.target, tokenId, price)
      )
        .to.emit(marketplace, "MarketItemCreated")
        .withArgs(1, nft.target, tokenId, seller.address, price);

      const item = await marketplace.getMarketItem(1);
      expect(item.seller).to.equal(seller.address);
      expect(item.price).to.equal(price);
      expect(item.sold).to.equal(false);
    });

    it("Should fail if price is zero", async function () {
      await nft.connect(seller).approve(marketplace.target, tokenId);

      await expect(
        marketplace.connect(seller).createMarketItem(nft.target, tokenId, 0)
      ).to.be.revertedWith("Price must be greater than 0");
    });

    it("Should fail if not the owner", async function () {
      await expect(
        marketplace.connect(buyer).createMarketItem(nft.target, tokenId, price)
      ).to.be.revertedWith("Not the owner");
    });

    it("Should fail if marketplace not approved", async function () {
      await expect(
        marketplace.connect(seller).createMarketItem(nft.target, tokenId, price)
      ).to.be.revertedWith("Market not approved");
    });
  });

  describe("Purchasing", function () {
    const tokenId = 1;
    const price = ethers.parseEther("1");
    let itemId;

    beforeEach(async function () {
      await nft.connect(seller).approve(marketplace.target, tokenId);
      const tx = await marketplace.connect(seller).createMarketItem(
        nft.target,
        tokenId,
        price
      );
      const receipt = await tx.wait();
      itemId = 1;
    });

    it("Should complete a sale", async function () {
      const sellerBalanceBefore = await ethers.provider.getBalance(seller.address);

      await expect(
        marketplace.connect(buyer).createMarketSale(itemId, { value: price })
      )
        .to.emit(marketplace, "MarketItemSold")
        .withArgs(itemId, buyer.address, price);

      // Check NFT ownership
      expect(await nft.ownerOf(tokenId)).to.equal(buyer.address);

      // Check item status
      const item = await marketplace.getMarketItem(itemId);
      expect(item.sold).to.equal(true);
      expect(item.owner).to.equal(buyer.address);

      // Check seller received payment (minus platform fee)
      const sellerBalanceAfter = await ethers.provider.getBalance(seller.address);
      const expectedProceeds = price - (price * BigInt(PLATFORM_FEE)) / BigInt(10000);
      expect(sellerBalanceAfter - sellerBalanceBefore).to.equal(expectedProceeds);
    });

    it("Should fail if incorrect price sent", async function () {
      await expect(
        marketplace.connect(buyer).createMarketSale(itemId, { 
          value: ethers.parseEther("0.5") 
        })
      ).to.be.revertedWith("Incorrect price");
    });

    it("Should fail if item already sold", async function () {
      await marketplace.connect(buyer).createMarketSale(itemId, { value: price });

      await expect(
        marketplace.connect(buyer).createMarketSale(itemId, { value: price })
      ).to.be.revertedWith("Item already sold");
    });

    it("Should fail if seller tries to buy own item", async function () {
      await expect(
        marketplace.connect(seller).createMarketSale(itemId, { value: price })
      ).to.be.revertedWith("Cannot buy own item");
    });
  });

  describe("Cancellation", function () {
    const tokenId = 1;
    const price = ethers.parseEther("1");
    let itemId;

    beforeEach(async function () {
      await nft.connect(seller).approve(marketplace.target, tokenId);
      await marketplace.connect(seller).createMarketItem(nft.target, tokenId, price);
      itemId = 1;
    });

    it("Should cancel a listing", async function () {
      await expect(marketplace.connect(seller).cancelMarketItem(itemId))
        .to.emit(marketplace, "MarketItemCanceled")
        .withArgs(itemId);

      const item = await marketplace.getMarketItem(itemId);
      expect(item.sold).to.equal(true);
    });

    it("Should fail if not the seller", async function () {
      await expect(
        marketplace.connect(buyer).cancelMarketItem(itemId)
      ).to.be.revertedWith("Not the seller");
    });
  });

  describe("Fetching Items", function () {
    beforeEach(async function () {
      // Create multiple listings
      for (let i = 1; i <= 3; i++) {
        await nft.connect(seller).mint(`https://example.com/token/${i}`);
        await nft.connect(seller).approve(marketplace.target, i);
        await marketplace.connect(seller).createMarketItem(
          nft.target,
          i,
          ethers.parseEther(`${i}`)
        );
      }
    });

    it("Should fetch active items", async function () {
      const items = await marketplace.fetchActiveItems();
      expect(items.length).to.equal(3);
    });

    it("Should fetch user listings", async function () {
      const items = await marketplace.fetchUserListings(seller.address);
      expect(items.length).to.equal(3);
    });

    it("Should update active items after sale", async function () {
      await marketplace.connect(buyer).createMarketSale(1, {
        value: ethers.parseEther("1"),
      });

      const items = await marketplace.fetchActiveItems();
      expect(items.length).to.equal(2);
    });
  });

  describe("Platform Fee Management", function () {
    it("Should update platform fee", async function () {
      const newFee = 500; // 5%
      await expect(marketplace.connect(owner).updatePlatformFee(newFee))
        .to.emit(marketplace, "PlatformFeeUpdated")
        .withArgs(newFee);

      expect(await marketplace.platformFee()).to.equal(newFee);
    });

    it("Should fail if fee is too high", async function () {
      await expect(
        marketplace.connect(owner).updatePlatformFee(1001)
      ).to.be.revertedWith("Fee too high");
    });

    it("Should fail if not owner", async function () {
      await expect(
        marketplace.connect(seller).updatePlatformFee(500)
      ).to.be.reverted;
    });
  });

  describe("Withdrawals", function () {
    it("Should withdraw accumulated fees", async function () {
      // Create and complete a sale
      await nft.connect(seller).approve(marketplace.target, 1);
      await marketplace.connect(seller).createMarketItem(
        nft.target,
        1,
        ethers.parseEther("1")
      );
      await marketplace.connect(buyer).createMarketSale(1, {
        value: ethers.parseEther("1"),
      });

      const ownerBalanceBefore = await ethers.provider.getBalance(owner.address);
      const marketplaceBalance = await ethers.provider.getBalance(marketplace.target);

      const tx = await marketplace.connect(owner).withdrawFees();
      const receipt = await tx.wait();
      const gasUsed = receipt.gasUsed * receipt.gasPrice;

      const ownerBalanceAfter = await ethers.provider.getBalance(owner.address);

      expect(ownerBalanceAfter).to.equal(
        ownerBalanceBefore + marketplaceBalance - gasUsed
      );
    });

    it("Should fail if not owner", async function () {
      await expect(marketplace.connect(seller).withdrawFees()).to.be.reverted;
    });
  });

  describe("Market Stats", function () {
    it("Should return correct market statistics", async function () {
      // Create 3 listings
      for (let i = 1; i <= 3; i++) {
        await nft.connect(seller).mint(`https://example.com/token/${i}`);
        await nft.connect(seller).approve(marketplace.target, i);
        await marketplace.connect(seller).createMarketItem(
          nft.target,
          i,
          ethers.parseEther(`${i}`)
        );
      }

      // Sell 1 item
      await marketplace.connect(buyer).createMarketSale(1, {
        value: ethers.parseEther("1"),
      });

      const stats = await marketplace.getMarketStats();
      expect(stats.totalItems).to.equal(3);
      expect(stats.soldItems).to.equal(1);
      expect(stats.activeItems).to.equal(2);
    });
  });
});