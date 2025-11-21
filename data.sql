-- ============================================
-- NFT Marketplace 数据库架构 postgresql数据库
-- PostgreSQL 15+
-- ============================================

-- 创建数据库（如果不存在）
-- CREATE DATABASE nft_marketplace;

-- 连接到数据库
\c nft_marketplace;

-- 启用扩展
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm"; -- 用于文本搜索

-- ============================================
-- 1. NFTs 表 - 存储 NFT 元数据
-- ============================================
CREATE TABLE IF NOT EXISTS nfts (
    id BIGSERIAL PRIMARY KEY,
    contract_address VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL, -- 支持大数字
    owner VARCHAR(42) NOT NULL,
    creator VARCHAR(42),
    name VARCHAR(255),
    description TEXT,
    image_url TEXT,
    metadata_uri TEXT,
    metadata JSONB, -- 存储完整的 metadata JSON
    
    -- 索引字段
    status VARCHAR(20) DEFAULT 'active', -- active, burned, transferred
    
    -- 统计字段
    view_count BIGINT DEFAULT 0,
    like_count BIGINT DEFAULT 0,
    
    -- 时间戳
    minted_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    -- 唯一约束
    CONSTRAINT uk_nfts_contract_token UNIQUE(contract_address, token_id)
);

-- NFTs 索引
CREATE INDEX idx_nfts_owner ON nfts(owner);
CREATE INDEX idx_nfts_creator ON nfts(creator);
CREATE INDEX idx_nfts_contract ON nfts(contract_address);
CREATE INDEX idx_nfts_status ON nfts(status);
CREATE INDEX idx_nfts_created_at ON nfts(created_at DESC);
CREATE INDEX idx_nfts_metadata_gin ON nfts USING gin(metadata); -- JSONB 索引

-- NFTs 表注释
COMMENT ON TABLE nfts IS 'NFT 元数据表';
COMMENT ON COLUMN nfts.contract_address IS 'NFT 合约地址';
COMMENT ON COLUMN nfts.token_id IS 'NFT Token ID（大数字字符串）';
COMMENT ON COLUMN nfts.metadata IS 'NFT 完整元数据 JSON';

-- ============================================
-- 2. Listings 表 - 市场挂单
-- ============================================
CREATE TABLE IF NOT EXISTS listings (
    id BIGSERIAL PRIMARY KEY,
    item_id BIGINT NOT NULL UNIQUE, -- 链上 item ID
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL,
    seller VARCHAR(42) NOT NULL,
    buyer VARCHAR(42), -- 买家地址（售出后填充）
    price VARCHAR(78) NOT NULL, -- Wei 单位的价格字符串
    price_numeric NUMERIC(78, 0), -- 用于排序和计算的数值类型
    
    -- 状态管理
    status VARCHAR(20) NOT NULL DEFAULT 'active', -- active, sold, cancelled
    
    -- 交易信息
    tx_hash VARCHAR(66), -- 创建交易哈希
    sale_tx_hash VARCHAR(66), -- 成交交易哈希
    
    -- 时间戳
    listed_at TIMESTAMP WITH TIME ZONE NOT NULL,
    sold_at TIMESTAMP WITH TIME ZONE,
    cancelled_at TIMESTAMP WITH TIME ZONE,
    expires_at TIMESTAMP WITH TIME ZONE, -- 可选的过期时间
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Listings 索引
CREATE INDEX idx_listings_item_id ON listings(item_id);
CREATE INDEX idx_listings_seller ON listings(seller);
CREATE INDEX idx_listings_buyer ON listings(buyer);
CREATE INDEX idx_listings_nft_contract ON listings(nft_contract);
CREATE INDEX idx_listings_token_id ON listings(token_id);
CREATE INDEX idx_listings_status ON listings(status);
CREATE INDEX idx_listings_listed_at ON listings(listed_at DESC);
CREATE INDEX idx_listings_price ON listings(price_numeric);
CREATE INDEX idx_listings_active_price ON listings(status, price_numeric) WHERE status = 'active'; -- 部分索引
CREATE INDEX idx_listings_contract_status ON listings(nft_contract, status);

-- Listings 表注释
COMMENT ON TABLE listings IS '市场挂单表';
COMMENT ON COLUMN listings.item_id IS '链上市场项 ID';
COMMENT ON COLUMN listings.price_numeric IS '价格数值类型（用于排序）';
COMMENT ON COLUMN listings.status IS '挂单状态：active-活跃, sold-已售, cancelled-已取消';

-- ============================================
-- 3. Transactions 表 - 交易记录
-- ============================================
CREATE TABLE IF NOT EXISTS transactions (
    id BIGSERIAL PRIMARY KEY,
    tx_hash VARCHAR(66) NOT NULL UNIQUE,
    block_number BIGINT NOT NULL,
    block_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- 交易类型
    tx_type VARCHAR(20) NOT NULL, -- list, sale, cancel, transfer, mint
    
    -- 关联信息
    listing_id BIGINT REFERENCES listings(id),
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL,
    
    -- 参与方
    from_address VARCHAR(42) NOT NULL,
    to_address VARCHAR(42),
    
    -- 金额信息
    value VARCHAR(78), -- Wei 单位
    value_numeric NUMERIC(78, 0),
    gas_price VARCHAR(78),
    gas_used BIGINT,
    
    -- 平台费用
    platform_fee VARCHAR(78),
    platform_fee_numeric NUMERIC(78, 0),
    
    -- 交易状态
    status VARCHAR(20) DEFAULT 'confirmed', -- pending, confirmed, failed
    
    -- 元数据
    log_index INTEGER,
    transaction_index INTEGER,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Transactions 索引
CREATE INDEX idx_transactions_tx_hash ON transactions(tx_hash);
CREATE INDEX idx_transactions_block ON transactions(block_number DESC);
CREATE INDEX idx_transactions_type ON transactions(tx_type);
CREATE INDEX idx_transactions_listing ON transactions(listing_id);
CREATE INDEX idx_transactions_nft ON transactions(nft_contract, token_id);
CREATE INDEX idx_transactions_from ON transactions(from_address);
CREATE INDEX idx_transactions_to ON transactions(to_address);
CREATE INDEX idx_transactions_timestamp ON transactions(block_timestamp DESC);
CREATE INDEX idx_transactions_value ON transactions(value_numeric DESC);

-- Transactions 表注释
COMMENT ON TABLE transactions IS '交易记录表';
COMMENT ON COLUMN transactions.tx_type IS '交易类型：list-挂单, sale-成交, cancel-取消, transfer-转账, mint-铸造';

-- ============================================
-- 4. Users 表 - 用户信息
-- ============================================
CREATE TABLE IF NOT EXISTS users (
    id BIGSERIAL PRIMARY KEY,
    address VARCHAR(42) NOT NULL UNIQUE,
    username VARCHAR(50),
    email VARCHAR(255),
    bio TEXT,
    avatar_url TEXT,
    banner_url TEXT,
    website TEXT,
    twitter_handle VARCHAR(50),
    discord_handle VARCHAR(50),
    
    -- 验证状态
    is_verified BOOLEAN DEFAULT FALSE,
    verified_at TIMESTAMP WITH TIME ZONE,
    
    -- 统计信息
    nfts_owned INTEGER DEFAULT 0,
    nfts_created INTEGER DEFAULT 0,
    nfts_sold INTEGER DEFAULT 0,
    total_sales_value VARCHAR(78) DEFAULT '0',
    total_purchases_value VARCHAR(78) DEFAULT '0',
    
    -- 设置
    email_notifications BOOLEAN DEFAULT TRUE,
    preferences JSONB DEFAULT '{}',
    
    -- 时间戳
    first_seen_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    last_active_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Users 索引
CREATE INDEX idx_users_address ON users(address);
CREATE INDEX idx_users_username ON users(username);
CREATE INDEX idx_users_verified ON users(is_verified);
CREATE INDEX idx_users_last_active ON users(last_active_at DESC);
CREATE INDEX idx_users_username_trgm ON users USING gin(username gin_trgm_ops); -- 模糊搜索

-- Users 表注释
COMMENT ON TABLE users IS '用户信息表';

-- ============================================
-- 5. Collections 表 - NFT 系列
-- ============================================
CREATE TABLE IF NOT EXISTS collections (
    id BIGSERIAL PRIMARY KEY,
    contract_address VARCHAR(42) NOT NULL UNIQUE,
    name VARCHAR(255) NOT NULL,
    symbol VARCHAR(20),
    description TEXT,
    
    -- 媒体
    logo_url TEXT,
    banner_url TEXT,
    featured_image_url TEXT,
    
    -- 创建者
    creator_address VARCHAR(42),
    owner_address VARCHAR(42),
    
    -- 链信息
    chain_id INTEGER NOT NULL DEFAULT 1,
    contract_type VARCHAR(20), -- ERC721, ERC1155
    
    -- 社交媒体
    website TEXT,
    discord_url TEXT,
    twitter_url TEXT,
    instagram_url TEXT,
    telegram_url TEXT,
    
    -- 统计信息
    total_supply BIGINT DEFAULT 0,
    total_owners INTEGER DEFAULT 0,
    floor_price VARCHAR(78) DEFAULT '0',
    total_volume VARCHAR(78) DEFAULT '0',
    
    -- 状态
    is_verified BOOLEAN DEFAULT FALSE,
    is_featured BOOLEAN DEFAULT FALSE,
    status VARCHAR(20) DEFAULT 'active', -- active, inactive, suspended
    
    -- Royalty
    royalty_percentage NUMERIC(5, 2) DEFAULT 0.00, -- 0.00 - 100.00
    royalty_recipient VARCHAR(42),
    
    -- 元数据
    metadata JSONB DEFAULT '{}',
    
    -- 时间戳
    deployed_at TIMESTAMP WITH TIME ZONE,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Collections 索引
CREATE INDEX idx_collections_address ON collections(contract_address);
CREATE INDEX idx_collections_creator ON collections(creator_address);
CREATE INDEX idx_collections_chain ON collections(chain_id);
CREATE INDEX idx_collections_verified ON collections(is_verified);
CREATE INDEX idx_collections_featured ON collections(is_featured);
CREATE INDEX idx_collections_floor_price ON collections(floor_price);
CREATE INDEX idx_collections_name_trgm ON collections USING gin(name gin_trgm_ops);

-- Collections 表注释
COMMENT ON TABLE collections IS 'NFT 系列表';

-- ============================================
-- 6. Offers 表 - 出价/报价
-- ============================================
CREATE TABLE IF NOT EXISTS offers (
    id BIGSERIAL PRIMARY KEY,
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL,
    
    -- 出价方
    offerer VARCHAR(42) NOT NULL,
    
    -- 价格
    price VARCHAR(78) NOT NULL,
    price_numeric NUMERIC(78, 0),
    
    -- 状态
    status VARCHAR(20) DEFAULT 'active', -- active, accepted, rejected, expired, cancelled
    
    -- 过期时间
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    
    -- 交易信息
    tx_hash VARCHAR(66),
    accepted_tx_hash VARCHAR(66),
    
    -- 时间戳
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    accepted_at TIMESTAMP WITH TIME ZONE,
    
    CONSTRAINT chk_offers_expires CHECK (expires_at > created_at)
);

-- Offers 索引
CREATE INDEX idx_offers_nft ON offers(nft_contract, token_id);
CREATE INDEX idx_offers_offerer ON offers(offerer);
CREATE INDEX idx_offers_status ON offers(status);
CREATE INDEX idx_offers_expires ON offers(expires_at);
CREATE INDEX idx_offers_price ON offers(price_numeric DESC);
CREATE INDEX idx_offers_created ON offers(created_at DESC);

-- Offers 表注释
COMMENT ON TABLE offers IS '出价表';

-- ============================================
-- 7. Activities 表 - 活动记录
-- ============================================
CREATE TABLE IF NOT EXISTS activities (
    id BIGSERIAL PRIMARY KEY,
    activity_type VARCHAR(20) NOT NULL, -- mint, list, sale, transfer, offer, cancel_listing, cancel_offer
    
    -- NFT 信息
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL,
    
    -- 参与方
    from_address VARCHAR(42),
    to_address VARCHAR(42),
    
    -- 价格（如果适用）
    price VARCHAR(78),
    price_numeric NUMERIC(78, 0),
    
    -- 交易信息
    tx_hash VARCHAR(66),
    
    -- 时间
    occurred_at TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Activities 索引
CREATE INDEX idx_activities_type ON activities(activity_type);
CREATE INDEX idx_activities_nft ON activities(nft_contract, token_id);
CREATE INDEX idx_activities_from ON activities(from_address);
CREATE INDEX idx_activities_to ON activities(to_address);
CREATE INDEX idx_activities_occurred ON activities(occurred_at DESC);
CREATE INDEX idx_activities_tx_hash ON activities(tx_hash);

-- Activities 表注释
COMMENT ON TABLE activities IS '活动记录表（用于展示交易历史）';

-- ============================================
-- 8. Favorites 表 - 收藏/点赞
-- ============================================
CREATE TABLE IF NOT EXISTS favorites (
    id BIGSERIAL PRIMARY KEY,
    user_address VARCHAR(42) NOT NULL,
    nft_contract VARCHAR(42) NOT NULL,
    token_id VARCHAR(78) NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_favorites UNIQUE(user_address, nft_contract, token_id)
);

-- Favorites 索引
CREATE INDEX idx_favorites_user ON favorites(user_address);
CREATE INDEX idx_favorites_nft ON favorites(nft_contract, token_id);
CREATE INDEX idx_favorites_created ON favorites(created_at DESC);

-- Favorites 表注释
COMMENT ON TABLE favorites IS '用户收藏表';

-- ============================================
-- 9. Events 表 - 区块链事件日志
-- ============================================
CREATE TABLE IF NOT EXISTS events (
    id BIGSERIAL PRIMARY KEY,
    event_name VARCHAR(50) NOT NULL,
    contract_address VARCHAR(42) NOT NULL,
    
    -- 区块信息
    block_number BIGINT NOT NULL,
    block_timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    transaction_hash VARCHAR(66) NOT NULL,
    log_index INTEGER NOT NULL,
    
    -- 事件数据
    event_data JSONB NOT NULL,
    
    -- 处理状态
    processed BOOLEAN DEFAULT FALSE,
    processed_at TIMESTAMP WITH TIME ZONE,
    
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    
    CONSTRAINT uk_events UNIQUE(transaction_hash, log_index)
);

-- Events 索引
CREATE INDEX idx_events_name ON events(event_name);
CREATE INDEX idx_events_contract ON events(contract_address);
CREATE INDEX idx_events_block ON events(block_number DESC);
CREATE INDEX idx_events_tx_hash ON events(transaction_hash);
CREATE INDEX idx_events_processed ON events(processed) WHERE NOT processed;
CREATE INDEX idx_events_timestamp ON events(block_timestamp DESC);
CREATE INDEX idx_events_data_gin ON events USING gin(event_data);

-- Events 表注释
COMMENT ON TABLE events IS '区块链事件日志表';

-- ============================================
-- 10. Sync_State 表 - 同步状态
-- ============================================
CREATE TABLE IF NOT EXISTS sync_state (
    id SERIAL PRIMARY KEY,
    contract_address VARCHAR(42) NOT NULL UNIQUE,
    last_synced_block BIGINT NOT NULL DEFAULT 0,
    last_synced_at TIMESTAMP WITH TIME ZONE,
    sync_status VARCHAR(20) DEFAULT 'active', -- active, paused, error
    error_message TEXT,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Sync_State 索引
CREATE INDEX idx_sync_state_contract ON sync_state(contract_address);
CREATE INDEX idx_sync_state_status ON sync_state(sync_status);

-- Sync_State 表注释
COMMENT ON TABLE sync_state IS '区块链同步状态表';

-- ============================================
-- 视图：活跃挂单统计
-- ============================================
CREATE OR REPLACE VIEW v_active_listings_stats AS
SELECT 
    nft_contract,
    COUNT(*) as total_listings,
    MIN(price_numeric) as floor_price,
    MAX(price_numeric) as ceiling_price,
    AVG(price_numeric) as avg_price,
    COUNT(DISTINCT seller) as unique_sellers
FROM listings
WHERE status = 'active'
GROUP BY nft_contract;

COMMENT ON VIEW v_active_listings_stats IS '活跃挂单统计视图';

-- ============================================
-- 视图：用户交易统计
-- ============================================
CREATE OR REPLACE VIEW v_user_stats AS
SELECT 
    address,
    nfts_owned,
    nfts_created,
    nfts_sold,
    (SELECT COUNT(*) FROM listings WHERE seller = users.address AND status = 'active') as active_listings,
    (SELECT COUNT(*) FROM listings WHERE seller = users.address AND status = 'sold') as completed_sales,
    (SELECT COUNT(*) FROM transactions WHERE from_address = users.address) as total_transactions,
    total_sales_value,
    total_purchases_value
FROM users;

COMMENT ON VIEW v_user_stats IS '用户交易统计视图';

-- ============================================
-- 视图：系列统计
-- ============================================
CREATE OR REPLACE VIEW v_collection_stats AS
SELECT 
    c.id,
    c.contract_address,
    c.name,
    c.floor_price,
    c.total_volume,
    COUNT(DISTINCT n.owner) as unique_owners,
    COUNT(n.id) as total_nfts,
    COUNT(l.id) as active_listings,
    (SELECT COUNT(*) FROM transactions t WHERE t.nft_contract = c.contract_address) as total_transactions
FROM collections c
LEFT JOIN nfts n ON n.contract_address = c.contract_address
LEFT JOIN listings l ON l.nft_contract = c.contract_address AND l.status = 'active'
GROUP BY c.id, c.contract_address, c.name, c.floor_price, c.total_volume;

COMMENT ON VIEW v_collection_stats IS '系列统计视图';

-- ============================================
-- 触发器：自动更新 updated_at
-- ============================================
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 为所有表添加 updated_at 触发器
CREATE TRIGGER update_nfts_updated_at BEFORE UPDATE ON nfts
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_listings_updated_at BEFORE UPDATE ON listings
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_transactions_updated_at BEFORE UPDATE ON transactions
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_collections_updated_at BEFORE UPDATE ON collections
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_offers_updated_at BEFORE UPDATE ON offers
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

CREATE TRIGGER update_sync_state_updated_at BEFORE UPDATE ON sync_state
    FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();

-- ============================================
-- 触发器：自动将 price 转换为 price_numeric
-- ============================================
CREATE OR REPLACE FUNCTION sync_price_numeric()
RETURNS TRIGGER AS $$
BEGIN
    IF NEW.price IS NOT NULL THEN
        NEW.price_numeric = NEW.price::NUMERIC;
    END IF;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER sync_listings_price_numeric BEFORE INSERT OR UPDATE ON listings
    FOR EACH ROW EXECUTE FUNCTION sync_price_numeric();

CREATE TRIGGER sync_transactions_value_numeric BEFORE INSERT OR UPDATE ON transactions
    FOR EACH ROW WHEN (NEW.value IS NOT NULL)
    EXECUTE FUNCTION sync_price_numeric();

-- ============================================
-- 函数：计算系列地板价
-- ============================================
CREATE OR REPLACE FUNCTION calculate_floor_price(p_contract_address VARCHAR)
RETURNS VARCHAR AS $$
DECLARE
    v_floor_price VARCHAR;
BEGIN
    SELECT MIN(price)
    INTO v_floor_price
    FROM listings
    WHERE nft_contract = p_contract_address
    AND status = 'active';
    
    RETURN COALESCE(v_floor_price, '0');
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 函数：更新系列统计
-- ============================================
CREATE OR REPLACE FUNCTION update_collection_stats(p_contract_address VARCHAR)
RETURNS VOID AS $$
BEGIN
    UPDATE collections
    SET 
        floor_price = calculate_floor_price(p_contract_address),
        total_supply = (SELECT COUNT(*) FROM nfts WHERE contract_address = p_contract_address),
        total_owners = (SELECT COUNT(DISTINCT owner) FROM nfts WHERE contract_address = p_contract_address),
        total_volume = (
            SELECT COALESCE(SUM(value_numeric), 0)::VARCHAR
            FROM transactions
            WHERE nft_contract = p_contract_address
            AND tx_type = 'sale'
        ),
        updated_at = CURRENT_TIMESTAMP
    WHERE contract_address = p_contract_address;
END;
$$ LANGUAGE plpgsql;

-- ============================================
-- 插入测试数据（可选）
-- ============================================

-- 插入示例系列
INSERT INTO collections (contract_address, name, symbol, description, chain_id, is_verified)
VALUES 
    ('0x1234567890123456789012345678901234567890', 'Bored Ape NFT Club', 'BANC', 'A collection of unique ape NFTs', 1, true),
    ('0x2345678901234567890123456789012345678901', 'CryptoPunks Tribute', 'CPT', 'Tribute to the OG NFT collection', 1, true)
ON CONFLICT (contract_address) DO NOTHING;

-- 插入示例用户
INSERT INTO users (address, username, is_verified)
VALUES 
    ('0xabcdef1234567890abcdef1234567890abcdef12', 'NFTCollector', true),
    ('0xbcdef1234567890abcdef1234567890abcdef123', 'CryptoArtist', true)
ON CONFLICT (address) DO NOTHING;

-- ============================================
-- 授权（根据需要调整）
-- ============================================
-- GRANT ALL PRIVILEGES ON ALL TABLES IN SCHEMA public TO nft_marketplace_user;
-- GRANT ALL PRIVILEGES ON ALL SEQUENCES IN SCHEMA public TO nft_marketplace_user;

-- ============================================
-- 完成
-- ============================================
SELECT 'Database schema created successfully!' as status;