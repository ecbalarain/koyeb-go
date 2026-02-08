# TiDB Cloud Connection Troubleshooting

## The Problem
Error: `Missing user name prefix`

This error indicates that TiDB Cloud is not recognizing the username format. This typically happens when there's a mismatch between what TiDB Cloud expects and what we're sending.

## Solution: Get Fresh Connection String from TiDB Cloud

### Step 1: Get Official Connection String

1. Go to [TiDB Cloud Console](https://tidbcloud.com/)
2. Navigate to **Clusters** page
3. Click on your cluster name
4. Click **Connect** button (top-right)
5. In the connection dialog:
   - Select **General** connection type
   - Select **MySQL CLI** or **Go** as the client
   - Copy the full connection details

### Step 2: What to Look For

The connection dialog should show something like:

**For Go applications:**
```
Host: gateway01.ap-northeast-1.prod.aws.tidbcloud.com
Port: 4000
User: <your-prefix>.root
Password: <your-password>
Database: test
```

**The complete DSN format might look like:**
```
<prefix>.root:<password>@tcp(gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000)/test?tls=tidb
```

### Step 3: Common Issues

1. **Wrong Format**: TiDB Cloud Serverless uses format: `<cluster-prefix>.root`
   - Example: `3pTAoNNegb47Uc8.root`
   
2. **Dedicated vs Serverless**: Format differs between tiers
   - **Serverless**: Requires prefix (`xxxxx.root`)
   - **Dedicated**: May not require prefix (`root`)

3. **Organization Prefix**: Some accounts need: `<org>.<cluster>.<user>`

### Step 4: Verify Your Cluster Type

In TiDB Cloud Console:
- Look for **Tier Type**: Serverless or Dedicated
- Serverless tier ALWAYS requires the prefix

## Alternative: Use Standard Connection Format

If you have the correct details from TiDB Cloud, the connection string should be:

```
mysql://<prefix>.root:<password>@<host>:4000/test?tls=true
```

**Example with proper format:**
```
mysql://3pTAoNNegb47Uc8%2Eroot:password@gateway01.ap-northeast-1.prod.aws.tidbcloud.com:4000/test?tls=true
```

Note: The `.` in the username might need to be URL-encoded as `%2E` if you're using a URL format.

## Debugging: Test Connection Locally

Once you have the correct connection string, test it locally:

```bash
# Update .env with the new connection string
# Then test:
go run cmd/migrate/main.go
```

---

## Action Items

Please provide:
1. ✅ The exact **Host** from TiDB Cloud console
2. ✅ The exact **User** (with prefix) from TiDB Cloud console  
3. ✅ Confirm your cluster **Tier** (Serverless or Dedicated)
4. ✅ Screenshot or copy of the connection dialog (hide password!)

This will help me construct the exact connection string format your TiDB cluster expects.
