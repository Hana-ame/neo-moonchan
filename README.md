
# 备忘录

存储在psql中。

# tutorio

## 0. 对环境的描述
```sql
SELECT usename FROM pg_user;
SELECT rolname, rolsuper, rolinherit, rolcreaterole, rolcreatedb, rolcanlogin, rolreplication, rolbypassrls
FROM pg_roles;
```

```txt
postgres	true	true	true	true	true	true	true
sukebei	false	true	false	false	true	false	false
```
<details>
  <summary>这些字段的含义如下：</summary>
rolname: 角色名（用户名）。
rolsuper: 是否为超级用户。
rolinherit: 是否继承其他角色的权限。
rolcreaterole: 是否可以创建新角色。
rolcreatedb: 是否可以创建新数据库。
rolcanlogin: 是否可以登录。
rolreplication: 是否可以进行复制。
rolbypassrls: 是否可以绕过行级安全策略。
</details>

```sql
SELECT datname FROM pg_database WHERE datistemplate = false;
```

### 0.1 创建用户和数据库


#### 创建user
```sql
CREATE ROLE new_user WITH LOGIN PASSWORD 'user_password';
```

#### 创建db
```sql
CREATE DATABASE new_db OWNER new_user;
```

#### revoke user's privilege
```sql
DO $$ 
DECLARE 
    db_name TEXT;
BEGIN
    FOR db_name IN 
        SELECT datname 
        FROM pg_database 
        WHERE datistemplate = false AND datname != 'new_db'
    LOOP
        EXECUTE format('REVOKE ALL PRIVILEGES ON DATABASE %I FROM new_user;', db_name);
    END LOOP;
END $$;
```

#### 重命名user
```sql
ALTER ROLE new_user RENAME TO renamed_user;
```

#### 验证
```sql
SELECT rolname FROM pg_roles;
```

#### 断开所有链接
```sql
SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'new_db' AND pid <> pg_backend_pid();
```

#### 重命名db
```sql
ALTER DATABASE new_db RENAME TO renamed_db;
```
#### 验证
```sql
SELECT datname FROM pg_database;
```

#### full code
```sql
CREATE ROLE new_user WITH LOGIN PASSWORD 'user_password';
CREATE DATABASE new_db OWNER new_user;

DO $$ 
DECLARE 
    db_name TEXT;
BEGIN
    FOR db_name IN 
        SELECT datname 
        FROM pg_database 
        WHERE datistemplate = false AND datname != 'new_db'
    LOOP
        EXECUTE format('REVOKE ALL PRIVILEGES ON DATABASE %I FROM new_user;', db_name);
    END LOOP;
END $$;

ALTER ROLE new_user RENAME TO renamed_user;

SELECT pg_terminate_backend(pg_stat_activity.pid)
FROM pg_stat_activity
WHERE pg_stat_activity.datname = 'new_db' AND pid <> pg_backend_pid();

ALTER DATABASE new_db RENAME TO renamed_db;

--SELECT rolname FROM pg_roles;
--SELECT datname FROM pg_database;
```

为啥要分开运行。  
因为Ctrl + Enter是运行当前行。

#### 改密码
```sql
ALTER USER target_user WITH PASSWORD 'new_password';
```

#### 检查权限
```sql
SELECT datname, pg_get_userbyid(datdba) AS owner, has_database_privilege(rolname, datname, 'CONNECT') AS can_connect,
       has_database_privilege(rolname, datname, 'CREATE') AS can_create
FROM pg_database
CROSS JOIN pg_roles
WHERE rolname = 'renamed_user';
```

#### 附：为啥一定能connect

一定能connect的原因。
[![10:47:08](https://moonchan.xyz/icon/stackoverflow.com)postgresql - Postgres revoke database access from user - Stack Overflow](https://stackoverflow.com/questions/49206699/postgres-revoke-database-access-from-user)

#### 查看编码

```sql
SELECT pg_encoding_to_char(encoding) AS encoding
FROM pg_database
WHERE datname = 'your_database_name';
```

### 0.2 建表
```sql
CREATE TABLE accounts (
    email VARCHAR(255) PRIMARY KEY,                -- 自动递增的主键
    username VARCHAR(50) UNIQUE NOT NULL,     -- 用户名，唯一且不允许为空
    password_hash VARCHAR(255) NOT NULL,      -- 密码哈希，存储加密后的密码

    country CHAR(2) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,       

    flag VARCHAR(255) NOT NULL, -- use `|` to split.

    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP,                     -- 上次登录时间
    failed_attempts INTEGER DEFAULT 0,        -- 登录失败次数，默认值为0
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- 最后更新时间    
);

CREATE TABLE sessions (
    login_id VARCHAR(255) PRIMARY KEY,       -- 唯一标识登录记录的ID
    username VARCHAR(50) NOT NULL,              -- 用户ID，外键关联到用户表
    login_time TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- 登录时间

    country CHAR(2) NOT NULL,
    ip_address VARCHAR(45) NOT NULL,            -- 用户登录的IP地址

    user_agent TEXT NOT NULL,                   -- 用户的浏览器或客户端信息
  
    FOREIGN KEY (username) REFERENCES accounts (username) ON DELETE CASCADE -- 外键，确保引用的用户存在
    -- CONSTRAINT unique_login_per_user UNIQUE (username, login_time)  -- 确保每个用户在相同时间只能有一条登录记录
);
CREATE INDEX idx_column_username
ON sessions USING HASH (username);

```

#### 没有 on update
```sql
-- 创建表
CREATE TABLE accounts (
    email VARCHAR(255) PRIMARY KEY,                -- 主键
    username VARCHAR(50) UNIQUE NOT NULL,          -- 用户名，唯一且不允许为空
    password_hash VARCHAR(255) NOT NULL,           -- 密码哈希，存储加密后的密码

    country CHAR(2) NOT NULL,                      -- 国家代码
    ip_address VARCHAR(45) NOT NULL,               -- IP 地址

    flag VARCHAR(255) NOT NULL,                    -- 标志字段
    
    last_login TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 上次登录时间
    failed_attempts INTEGER DEFAULT 0,             -- 登录失败次数，默认值为0
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- 创建时间，默认当前时间
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP -- 最后更新时间
);

-- 创建触发器函数
CREATE OR REPLACE FUNCTION update_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = CURRENT_TIMESTAMP;
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

-- 创建触发器
CREATE TRIGGER trigger_update_updated_at
BEFORE UPDATE ON accounts
FOR EACH ROW
EXECUTE FUNCTION update_updated_at();
```

psql还真是吃屎。。怎么找准database啊。
要在DBeaver里面切。

## 1. go

### 1.0 初始化

go get -u github.com/lib/pq

### 1.1 链接db
这次就不gorm了，感觉不省事。
不如写个脚本生成代码。