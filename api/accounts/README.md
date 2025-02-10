```sql

-- sudo su postgres
-- psql
-- \c moonchan
GRANT ALL PRIVILEGES ON DATABASE moonchan TO lumin;
GRANT ALL PRIVILEGES ON SCHEMA public TO lumin;

-- \c moonchan

CREATE TABLE accounts (
    email VARCHAR(255) PRIMARY KEY NOT NULL,
    username VARCHAR(255) NOT NULL UNIQUE,
    password VARCHAR(255) NOT NULL

    
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

```

curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/register -d '{"email":"admin@moonchan.xyz","password":"admin"}' -v
curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/login -d '{"email":"admin@moonchan.xyz","password":"admin"}' -v
curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/login -d '{"email":"user@moonchan.xyz","password":"admin"}' -v
curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/login -d '{"email":"admin@moonchan.xyz","password":"invalid"}' -v
curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/update -d '{"email":"admin@moonchan.xyz","password":"admin", "newpassword":"user"}' -v
curl -x "" -X POST http://127.24.7.29:8080/api/chan/accounts/update -d '{"email":"admin@moonchan.xyz","password":"user", "newpassword":"admin"}' -v


```sql
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
  NEW.updated_at = NOW();
  RETURN NEW;
END;
$$ LANGUAGE plpgsql;
```