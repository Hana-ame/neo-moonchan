```sql
CREATE TABLE profiles (
    email VARCHAR(255) PRIMARY KEY NOT NULL, -- 用户邮箱，作为主键，确保唯一性
    username VARCHAR(50) NOT NULL UNIQUE,    -- 用户名，限制最大长度为50，不允许重复
    host VARCHAR(255),                   -- 用户所在的服务器或平台
    bio TEXT,                            -- 用户签名信息，使用 TEXT 类型以适应较长文本
    avatar_url VARCHAR(1024),             -- 头像 URL 链接，使用 VARCHAR 存储，可以根据需要调整长度
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认值为当前时间
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP -- 最后更新时间，默认值为当前时间
);

-- 创建触发器自动更新 updated_at 列
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = NOW();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE OR REPLACE TRIGGER profiles_updated_at
BEFORE UPDATE ON profiles
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();

-- 添加索引可以加速查找
CREATE INDEX idx_profiles_username ON profiles(username);

````