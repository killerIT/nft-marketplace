const { ethers } = require("hardhat");

async function main() {
  console.log("Deploying NFTMarketplace contract...");

  // 获取合约工厂
  const NFTMarketplace = await ethers.getContractFactory("NFTMarketplace");

  // 获取部署者地址
  const [deployer] = await ethers.getSigners();
  console.log("Deploying contracts with the account:", deployer.address);

  // 修复：使用正确的方式获取余额
  const balance = await ethers.provider.getBalance(deployer.address);
  console.log("Account balance:", balance.toString());

  // 部署合约 - 注意：你的合约构造函数可能不需要参数
  const marketplace = await NFTMarketplace.deploy();

  await marketplace.waitForDeployment();

  const marketplaceAddress = await marketplace.getAddress();
  console.log("NFTMarketplace deployed to:", marketplaceAddress);

  // 验证部署
  const owner = await marketplace.owner();
  const platformFee = await marketplace.platformFee();

  console.log("Contract owner:", owner);
  console.log("Platform fee:", platformFee.toString());

  // 保存部署地址到文件
  const fs = require("fs");
  const deploymentInfo = {
    marketplaceAddress: marketplaceAddress,
    deployer: deployer.address,
    timestamp: new Date().toISOString()
  };

  fs.writeFileSync(
      "deployment-info.json",
      JSON.stringify(deploymentInfo, null, 2)
  );

  console.log("Deployment info saved to deployment-info.json");
}

main()
    .then(() => process.exit(0))
    .catch((error) => {
      console.error(error);
      process.exit(1);
    });