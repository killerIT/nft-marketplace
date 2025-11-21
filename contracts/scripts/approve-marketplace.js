const hre = require("hardhat");

async function main() {
    const [signer] = await hre.ethers.getSigners();

    const nftAddress = "0xD4FD24e932a9c9F3d19Ba4021132d92cb036bcd2"; // 你的 NFT 合约地址
    const marketplaceAddress = "0x89d1FfeB79155e6728a57621C55AF95888378F72";
    const tokenId = 2;

    const nft = await hre.ethers.getContractAt("IERC721", nftAddress);

    // 授权市场合约
    const tx = await nft.approve(marketplaceAddress, tokenId);
    await tx.wait();

    console.log(`NFT #${tokenId} approved for marketplace`);
}

main().catch(console.error);