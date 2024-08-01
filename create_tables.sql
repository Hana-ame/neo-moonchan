-- 选择当前数据库 (你可以取消注释并修改数据库名)
-- SELECT current_database();

-- 创建 accounts 表
CREATE TABLE accounts (
    email VARCHAR(255) PRIMARY KEY,                -- 主键
    username VARCHAR(50) UNIQUE NOT NULL,          -- 用户名，唯一且不允许为空
    password_hash VARCHAR(255) NOT NULL,           -- 密码哈希，存储加密后的密码

    country CHAR(2) NOT NULL,                      -- 国家代码
    ip_address VARCHAR(45) NOT NULL,               -- IP 地址

    flag VARCHAR(255) NOT NULL DEFAULT 'created',                    -- 标志字段
    
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 上次登录时间
    failed_attempts INTEGER NOT NULL DEFAULT 0,             -- 登录失败次数，默认值为0
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 最后更新时间    
);

-- 创建唯一索引，确保用户名不区分大小写的唯一性
CREATE UNIQUE INDEX unique_username 
ON accounts (lower(username));

-- 创建 sessions 表
CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,          -- 唯一标识登录记录的ID
    username VARCHAR(50) NOT NULL,                -- 用户ID，外键关联到 accounts 表的 username
    login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 登录时间

    country CHAR(2) NOT NULL,                     -- 用户所在国家代码
    ip_address VARCHAR(45) NOT NULL,              -- 用户登录的IP地址

    user_agent TEXT NOT NULL,                     -- 用户的浏览器或客户端信息
  
    FOREIGN KEY (username) REFERENCES accounts (username) ON DELETE CASCADE -- 外键，确保引用的用户存在
);

-- 创建索引，使用 HASH 索引方式以加速查询
CREATE INDEX idx_column_username
ON sessions USING HASH (username);

CREATE TABLE users (
    username VARCHAR(50) PRIMARY KEY,             -- 用户名，主键且不允许为空，与accounts中的username一致
    display_name VARCHAR(50),                     -- 展示用的用户名，可以为空
    avatar_url VARCHAR(255),                      -- 用户头像的链接，可以为空
    settings JSONB NOT NULL DEFAULT '{}',         -- 用户的设置，使用JSONB存储
    flag VARCHAR(50) NOT NULL DEFAULT 'created',   -- 用户状态，默认值为 'created'
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 最后更新时间，默认当前时间

    CONSTRAINT fk_username_accounts
    FOREIGN KEY (username) REFERENCES accounts (username) ON DELETE CASCADE -- 外键约束，确保username与accounts表中的username一致
);

-- 查询当前数据库中的所有表 (你可以取消注释以执行此查询)
-- SELECT current_database();
-- SELECT schemaname, tablename
-- FROM pg_tables;
