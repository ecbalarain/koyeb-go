# TiDB Cloud Connection Fix

## Issue
TiDB Cloud username contains a dot (`.`) which must be URL-encoded in the connection string.

## Current (Broken)
```
mysql://3QQuRDHAFitYTwM.root:HWvBROQOOxQed5c1@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true
```

## Fixed (Working)
```
mysql://3QQuRDHAFitYTwM%2Eroot:HWvBROQOOxQed5c1@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true
```

**Change:** The dot in username `3QQuRDHAFitYTwM.root` becomes `3QQuRDHAFitYTwM%2Eroot`

---

## Update in Koyeb

1. Go to Koyeb dashboard
2. Click on your service `oxlook-api`
3. Go to **Settings** > **Environment**
4. Find the `DATABASE_URL` variable
5. Click **Edit**
6. Replace the value with:
   ```
   mysql://3QQuRDHAFitYTwM%2Eroot:HWvBROQOOxQed5c1@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true
   ```
7. Click **Save**
8. Koyeb will automatically redeploy with the new connection string

---

## Expected Logs After Fix

```
🚀 Starting OXLOOK API deployment...
📊 Running database migrations...
2026/02/08 XX:XX:XX Connecting to database...
2026/02/08 XX:XX:XX Running migrations...
2026/02/08 XX:XX:XX Skipping already applied migration: 000_create_migrations_table.sql
2026/02/08 XX:XX:XX Skipping already applied migration: 001_create_products.sql
...
✅ All migrations completed successfully!

    _______ __
   / ____(_) /_  ___  _____
  / /_  / / __ \/ _ \/ ___/
 / __/ / / /_/ /  __/ /
/_/   /_/_.___/\___/_/          v3.0.0
--------------------------------------------------
INFO Server started on: 	http://127.0.0.1:8000
Instance is healthy. All health checks are passing.
```

---

## Test Locally First (Optional)

Update `.env` and test:

```bash
# Update DATABASE_URL in .env
DATABASE_URL=mysql://3QQuRDHAFitYTwM%2Eroot:HWvBROQOOxQed5c1@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true

# Test migration
go run cmd/migrate/main.go

# Test API
go run main.go
```

---

## Why This Happens

TiDB Cloud uses cluster prefixes in usernames (format: `<prefix>.<username>`). When used in URLs, special characters like `.` must be URL-encoded:

- `.` becomes `%2E`
- `@` becomes `%40`
- `:` in password (if present) becomes `%3A`

The Go MySQL driver requires proper URL encoding in DSN strings.
