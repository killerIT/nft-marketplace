const hre = require("hardhat");

// 取消挂单
async function main() {
    const [signer] = await hre.ethers.getSigners();
    const marketplaceAddress = "0x89d1FfeB79155e6728a57621C55AF95888378F72";
    const itemId = 1; // 假设要取消 Item ID 1

    const marketplace = await hre.ethers.getContractAt(
        "NFTMarketplace",
        marketplaceAddress
    );

    const tx = await marketplace.cancelMarketItem(itemId);
    await tx.wait();

    console.log(`Listing #${itemId} cancelled successfully`);
}

main().catch(console.error);