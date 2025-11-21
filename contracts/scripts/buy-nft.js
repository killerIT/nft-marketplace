const hre = require("hardhat");

async function main() {
    const signers = await hre.ethers.getSigners(); // 获取所有账户
    // 确保使用不同的账户（不是卖家账户）
    if (signers.length < 2) {
        throw new Error("Need at least 2 accounts for buyer and seller");
    }
    const buyer = signers[1]; // 使用第二个账户
    console.log(`signers: ${signers}`);
    const marketplaceAddress = "0x89d1FfeB79155e6728a57621C55AF95888378F72";
    const itemId = 2;
    const price = hre.ethers.parseEther("0.0001");

    const marketplace = await hre.ethers.getContractAt(
        "NFTMarketplace",
        marketplaceAddress
    );

    console.log(`Buyer: ${buyer.address}`);
    console.log(`Balance before: ${hre.ethers.formatEther(await hre.ethers.provider.getBalance(buyer.address))} ETH`);

    // 购买
    const tx = await marketplace.connect(buyer).createMarketSale(itemId, {
        value: price
    });
    const receipt = await tx.wait();

    console.log(`NFT purchased successfully!`);
    console.log(`Transaction hash: ${receipt.hash}`);
    console.log(`Gas used: ${receipt.gasUsed}`);

    const balanceAfter = await hre.ethers.provider.getBalance(buyer.address);
    console.log(`Balance after: ${hre.ethers.formatEther(balanceAfter)} ETH`);
}

main().catch(console.error);