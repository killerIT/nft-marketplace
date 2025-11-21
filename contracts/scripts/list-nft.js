const hre = require("hardhat");

// 新建挂单
async function main() {
    const [signer] = await hre.ethers.getSigners();

    const marketplaceAddress = "0x89d1FfeB79155e6728a57621C55AF95888378F72";
    const nftAddress = "0xD4FD24e932a9c9F3d19Ba4021132d92cb036bcd2";
    const tokenId = 2;
    const price = hre.ethers.parseEther("0.0001"); // 0.1 ETH

    const marketplace = await hre.ethers.getContractAt(
        "NFTMarketplace",
        marketplaceAddress
    );

    // 上架
    const tx = await marketplace.createMarketItem(nftAddress, tokenId, price);
    const receipt = await tx.wait();

    // 获取 ItemId
    const event = receipt.logs.find(log => {
        try {
            return marketplace.interface.parseLog(log).name === "MarketItemCreated";
        } catch {
            return false;
        }
    });

    const itemId = event ? marketplace.interface.parseLog(event).args.itemId : null;

    console.log(`NFT listed successfully!`);
    console.log(`Item ID: ${itemId}`);
    console.log(`Price: ${hre.ethers.formatEther(price)} ETH`);
    console.log(`Transaction hash: ${receipt.hash}`);
}

main().catch(console.error);