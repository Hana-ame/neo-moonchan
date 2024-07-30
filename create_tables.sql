
-- SELECT current_database();

CREATE TABLE accounts (
    email VARCHAR(255) PRIMARY KEY,                -- 主键
    username VARCHAR(50) UNIQUE NOT NULL,          -- 用户名，唯一且不允许为空
    password_hash VARCHAR(255) NOT NULL,           -- 密码哈希，存储加密后的密码

    country CHAR(2) NOT NULL,                      -- 国家代码
    ip_address VARCHAR(45) NOT NULL,               -- IP 地址

    flag VARCHAR(255) NOT NULL,                    -- 标志字段
    
    last_login TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 上次登录时间
    failed_attempts INTEGER NOT NULL DEFAULT 0,             -- 登录失败次数，默认值为0
    created_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP -- 最后更新时间    
);
CREATE UNIQUE INDEX unique_username ON accounts (lower(username)); -- 不允许大小写重复

CREATE TABLE sessions (
    session_id VARCHAR(255) PRIMARY KEY,       -- 唯一标识登录记录的ID
    username VARCHAR(50) NOT NULL,              -- 用户ID，外键关联到用户表
    login_time TIMESTAMP NOT NULL DEFAULT CURRENT_TIMESTAMP,  -- 登录时间

    country CHAR(2) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,            -- 用户登录的IP地址

    user_agent TEXT NOT NULL,                   -- 用户的浏览器或客户端信息
  
    FOREIGN KEY (username) REFERENCES accounts (username) ON DELETE CASCADE -- 外键，确保引用的用户存在
    -- CONSTRAINT unique_login_per_user UNIQUE (username, login_time)  -- 确保每个用户在相同时间只能有一条登录记录
);
-- create index
CREATE INDEX idx_column_username
ON sessions USING HASH (username);

-- SELECT current_database();

-- SELECT schemaname, tablename
-- FROM pg_tables;
