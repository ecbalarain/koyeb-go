# Database Setup Guide

## Overview
This guide covers setting up a production MySQL database for the OXLOOK API.

## Option 1: PlanetScale (Recommended)

PlanetScale is a MySQL-compatible serverless database with excellent scalability and built-in branching.

### Setup Steps

1. **Create Account**
   - Go to https://planetscale.com
   - Sign up for free tier (5GB storage, 1 billion row reads/month)

2. **Create Database**
   ```bash
   # Via CLI (optional)
   pscale db create oxlook --region us-east
   
   # Or use web dashboard
   # Database name: oxlook
   # Region: Choose closest to your Koyeb deployment
   ```

3. **Create Branch** (optional, for development)
   ```bash
   pscale branch create oxlook development
   ```

4. **Get Connection String**
   - Dashboard > Database > Connect
   - Select "Connect with: Go"
   - Copy the connection string
   - Format: `mysql://username:password@host/database?tls=true`

5. **Create Database Password**
   - Dashboard > Settings > Passwords
   - Create new password
   - Select "main" branch for production
   - Copy the generated credentials

### Connection String Format
```
mysql://username:pscale_pw_xxxxx@aws.connect.psdb.cloud/oxlook?tls=true
```

### PlanetScale Benefits
- ✅ Auto-scaling
- ✅ Automatic backups
- ✅ 99.99% uptime SLA (paid plans)
- ✅ Branching for development
- ✅ No need to run migrations manually (supports schema changes)

---

## Option 2: Neon (PostgreSQL Alternative)

If you prefer PostgreSQL over MySQL, Neon is an excellent choice.

### Important: Switch to PostgreSQL Driver

1. **Update go.mod**
   ```bash
   go get github.com/jackc/pgx/v5
   ```

2. **Update database connection code** in `internal/database/database.go`
   Replace MySQL driver with PostgreSQL driver

3. **Convert SQL migrations** from MySQL to PostgreSQL syntax

### Neon Setup

1. **Create Account**
   - Go to https://neon.tech
   - Sign up for free tier (0.5 GB storage, 1 compute hour)

2. **Create Project**
   - Project name: oxlook
   - Region: Choose closest to Koyeb

3. **Get Connection String**
   - Dashboard > Connection Details
   - Select "Pooled connection"
   - Copy the connection string

### Connection String Format
```
postgres://username:password@ep-xxx.region.neon.tech/neondb?sslmode=require
```

---

## Option 3: Railway (Easy Setup)

Railway provides simple MySQL/PostgreSQL hosting with automatic deployments.

### Setup Steps

1. **Create Account**
   - Go to https://railway.app
   - Sign up with GitHub

2. **Create MySQL Database**
   - New Project > Add MySQL
   - Wait for provisioning

3. **Get Connection Details**
   - Click on MySQL service
   - Variables tab shows all connection details
   - Copy DATABASE_URL or construct manually

### Connection String Format
```
mysql://root:password@containers-us-west-xx.railway.app:port/railway?tls=true
```

---

## Option 4: Aiven (Production-Grade)

Aiven offers managed MySQL with excellent reliability and compliance.

### Setup Steps

1. **Create Account**
   - Go to https://aiven.io
   - Free trial available

2. **Create MySQL Service**
   - Service: MySQL 8
   - Cloud: AWS/GCP/Azure
   - Region: us-east-1 (or closest to Koyeb)
   - Plan: Hobbyist (free) or Startup-4

3. **Get Connection String**
   - Service > Overview
   - Copy MySQL URI

### Connection String Format
```
mysql://avnadmin:password@mysql-xxx.aivencloud.com:port/defaultdb?ssl-mode=REQUIRED
```

---

## Database Configuration for Koyeb

### Environment Variable

Add to Koyeb service:

```
DATABASE_URL=mysql://username:password@host:port/database?tls=true
```

### Connection Pool Settings

Already configured in `internal/database/database.go`:
- Max Open Connections: 25
- Max Idle Connections: 5
- Connection Max Lifetime: 5 minutes
- Connection Max Idle Time: 10 minutes

### SSL/TLS Requirements

Most providers require TLS. Ensure your connection string includes:
- `?tls=true` or `?sslmode=require`

---

## Initial Database Setup

### Migrations Run Automatically

The Docker container runs migrations on startup via `docker-entrypoint.sh`.

Migrations applied in order:
1. `000_create_migrations_table.sql` - Migration tracking
2. `001_create_products.sql` - Products table
3. `002_create_variants.sql` - Variants table
4. `003_create_orders.sql` - Orders table
5. `004_create_order_items.sql` - Order items table
6. `005_seed_products.sql` - Demo products
7. `006_seed_variants.sql` - Demo variants

### Manual Migration (if needed)

```bash
# Local testing
go run cmd/migrate/main.go

# Docker
docker exec -it container_name ./migrate
```

---

## Database Maintenance

### Backups

**PlanetScale:** Automatic backups (7-30 days retention depending on plan)

**Neon:** Automatic backups (7 days retention on free tier)

**Railway:** Manual backups via dashboard

**Aiven:** Automatic backups (configurable retention)

### Manual Backup

```bash
# MySQL
mysqldump -h host -u username -p database > backup.sql

# Restore
mysql -h host -u username -p database < backup.sql
```

### Monitoring

1. **Query Performance**
   - Enable slow query log in database settings
   - Monitor queries > 1 second

2. **Connection Pool**
   - Watch for "too many connections" errors
   - Adjust pool settings if needed

3. **Disk Usage**
   - Monitor via provider dashboard
   - Set up alerts at 80% usage

---

## Security Checklist

- [ ] Use strong password (min 32 characters)
- [ ] Enable SSL/TLS connections
- [ ] Restrict IP access (if provider supports)
- [ ] Use connection pooling
- [ ] Set up automated backups
- [ ] Enable monitoring/alerting
- [ ] Rotate database credentials periodically
- [ ] Use read replicas for high traffic (paid plans)

---

## Cost Estimation

### Free Tiers
- **PlanetScale:** Free forever (5GB, 1B reads/month)
- **Neon:** Free forever (0.5GB, 1 compute hour)
- **Railway:** $5 credit/month
- **Aiven:** 30-day trial, then paid

### Paid Plans (Monthly)
- **PlanetScale Scaler:** $29/month (10GB, 10B reads)
- **Neon Scale:** $69/month (50GB)
- **Railway:** Pay as you go (~$5-20/month for small DB)
- **Aiven Startup-4:** $49/month (2GB RAM, 80GB storage)

---

## Recommended Choice

**For OXLOOK:** PlanetScale Free Tier

**Reasons:**
- ✅ Generous free tier (sufficient for small-medium stores)
- ✅ MySQL compatible (no code changes needed)
- ✅ Excellent performance
- ✅ Automatic backups
- ✅ Easy scaling when needed
- ✅ Built-in connection pooling
- ✅ Great developer experience

---

## Next Steps

1. Create database account (PlanetScale recommended)
2. Create database and get connection string
3. Add DATABASE_URL to Koyeb environment variables
4. Deploy application (migrations run automatically)
5. Verify database connection in Koyeb logs
6. Test API endpoints

## Troubleshooting

### Connection Failed

```
Error: Failed to connect to database: Error 1045: Access denied
```

**Solution:** Verify username/password in DATABASE_URL

### SSL/TLS Errors

```
Error: x509: certificate signed by unknown authority
```

**Solution:** Add `?tls=skip-verify` (development only) or ensure `?tls=true`

### Too Many Connections

```
Error: Error 1040: Too many connections
```

**Solution:** 
- Reduce `MaxOpenConns` in database.go
- Upgrade database plan
- Check for connection leaks

### Migration Failed

```
Error: Migration failed: Table 'products' already exists
```

**Solution:**
- Check migrations table: `SELECT * FROM migrations;`
- Migrations are idempotent (safe to re-run)
- If needed, manually mark as run in migrations table
