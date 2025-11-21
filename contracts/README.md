# Sample Hardhat Project

This project demonstrates a basic Hardhat use case. It comes with a sample contract, a test for that contract, and a Hardhat Ignition module that deploys that contract.

Try running some of the following tasks:

```shell
npx hardhat help
npx hardhat test
REPORT_GAS=true npx hardhat test
npx hardhat node
npx hardhat ignition deploy ./ignition/modules/Lock.js
```
操作顺序
✅ 部署合约
npx hardhat run scripts/deploy.js --network sepolia
✅ 启动后端
cd backend
go run cmd/api/main.go
✅ 铸造 NFT
npx hardhat run scripts/mint-nft.js --network sepolia
✅ 同步 NFT 到后端
curl -X POST http://localhost:8080/api/v1/nfts \
-H "Content-Type: application/json" \
-d '{
"contract_address": "0xABCDEF1234567890...",
"token_id": "1",
"owner": "0xe68c4Aa728925085CE681A139e0D1B2e7bDc5B99",
"creator": "0xe68c4Aa728925085CE681A139e0D1B2e7bDc5B99",
"name": "Awesome NFT #1",
"description": "This is my first NFT on the marketplace",
"image_url": "https://ipfs.io/ipfs/QmXxxx.../1.png",
"metadata_uri": "https://example.com/metadata/1",
"metadata": {
"attributes": [
{"trait_type": "Background", "value": "Blue"},
{"trait_type": "Rarity", "value": "Common"}
]
}
}'
✅ 授权并上架 NFT
npx hardhat run scripts/approve-marketplace.js --network sepolia
npx hardhat run scripts/list-nft.js --network sepolia
✅ 创建挂单记录
curl -X POST http://localhost:8080/api/v1/listings \
-H "Content-Type: application/json" \
-d '{
"item_id": 1,
"nft_contract": "0xABCDEF1234567890...",
"token_id": "1",
"seller": "0xe68c4Aa728925085CE681A139e0D1B2e7bDc5B99",
"price": "100000000000000000",
"tx_hash": "0x1234567890abcdef..."
}'
✅ 购买 NFT
npx hardhat run scripts/buy-nft.js --network sepolia
✅ 查看交易和统计
