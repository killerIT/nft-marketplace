const hre = require("hardhat");

async function main() {
    const [signer] = await hre.ethers.getSigners();

    // 部署一个简单的 NFT 合约
    const NFT = await hre.ethers.getContractFactory("NFT");
    const nft = await NFT.deploy("MyNFT", "MNFT");
    await nft.waitForDeployment();

    const nftAddress = await nft.getAddress();
    console.log("NFT Contract deployed to:", nftAddress);

    // 铸造 3 个 NFT
    for (let i = 1; i <= 3; i++) {
        const tx = await nft.mint(`https://example.com/metadata/${i}`);
        await tx.wait();
        console.log(`NFT #${i} minted to ${signer.address}`);
    }

    console.log("\nMinting Summary:");
    console.log("NFT Contract:", nftAddress);
    console.log("Owner:", signer.address);
    console.log("Total minted: 3");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
        console.error(error);
        process.exit(1);
    });