# Monitoring & Alerting Setup

## Overview
Comprehensive monitoring and alerting configuration for OXLOOK production deployment.

---

## 1. Application Monitoring

### Koyeb Built-in Metrics

**Available Metrics (Free Tier):**
- HTTP requests/second
- Response time (p50, p95, p99)
- Error rate
- Memory usage
- CPU usage
- Network I/O

**Access:**
1. Koyeb Dashboard > Service > Metrics
2. Time ranges: 1h, 6h, 24h, 7d, 30d

**Grafana Integration (Paid Plans):**
- Custom dashboards
- Advanced queries
- Historical data

---

## 2. Uptime Monitoring

### Option 1: UptimeRobot (Recommended - Free)

**Setup:**
1. Sign up: https://uptimerobot.com
2. Create monitor:
   - **Monitor Type:** HTTPS
   - **Friendly Name:** OXLOOK API
   - **URL:** https://api.bhomanshah.com/health
   - **Monitoring Interval:** 5 minutes
   
3. Alert Contacts:
   - Email: your-email@domain.com
   - SMS (paid): +1234567890
   - Webhook: https://hooks.slack.com/... (optional)

4. Alert Threshold:
   - Alert when down for: 2 minutes (1 check)
   - Alert when up again: Yes

**Status Page (Optional):**
- Create public status page
- URL: https://stats.uptimerobot.com/your-id
- Share with customers

### Option 2: Pingdom

1. Sign up: https://www.pingdom.com (14-day trial)
2. Add uptime check:
   - URL: https://api.bhomanshah.com/health
   - Interval: 1 minute
   - Locations: Multiple regions
3. Configure alerts: Email, SMS, Slack

### Option 3: Checkly

1. Sign up: https://www.checklyhq.com
2. Create API check:
   - URL: https://api.bhomanshah.com/health
   - Frequency: 5 minutes
   - Locations: Frankfurt, N. Virginia, Singapore
3. Setup alerts: Email, Slack, PagerDuty

---

## 3. Error Tracking

### Option 1: Sentry (Recommended)

**Setup:**

1. **Create Account**
   - Go to https://sentry.io
   - Free tier: 5,000 events/month

2. **Create Project**
   - Platform: Go
   - Project name: oxlook-api

3. **Install Dependencies**
   ```bash
   go get github.com/getsentry/sentry-go
   ```

4. **Initialize Sentry** in main.go:
   ```go
   import (
       "github.com/getsentry/sentry-go"
       "log"
       "time"
   )

   func main() {
       // Initialize Sentry
       err := sentry.Init(sentry.ClientOptions{
           Dsn: os.Getenv("SENTRY_DSN"),
           Environment: os.Getenv("ENVIRONMENT"),
           Release: "oxlook-api@1.0.0",
           TracesSampleRate: 0.1, // 10% of requests
       })
       if err != nil {
           log.Printf("Sentry initialization failed: %v", err)
       }
       defer sentry.Flush(2 * time.Second)

       // ... rest of app
   }
   ```

5. **Add Error Reporting**
   ```go
   // In error handlers
   func (h *Handler) CreateOrder(c fiber.Ctx) error {
       // ... code ...
       if err != nil {
           sentry.CaptureException(err)
           return c.Status(500).JSON(fiber.Map{
               "error": "Internal error",
           })
       }
   }
   ```

6. **Environment Variable**
   Add to Koyeb:
   ```
   SENTRY_DSN=https://xxxxx@o123456.ingest.sentry.io/123456
   ```

**Configure Alerts:**
1. Sentry > Alerts > Create Alert
2. Rules:
   - When event count >= 10 in 5 minutes
   - When unique users affected >= 5
   - When error rate >= 1%
3. Actions: Email, Slack, PagerDuty

### Option 2: Rollbar

Similar setup to Sentry:
```bash
go get github.com/rollbar/rollbar-go
```

### Option 3: Log Aggregation

**Simple solution:** Export Koyeb logs to external service

```bash
# Stream logs to file
koyeb service logs oxlook-api --follow > app.log

# Or use log management service:
# - Logtail (free tier)
# - Papertrail (free tier)
# - Logz.io
```

---

## 4. Database Monitoring

### PlanetScale Monitoring

**Built-in Metrics:**
1. Dashboard > Insights
   - Query performance
   - Slow queries (>1s)
   - Connection count
   - Storage usage

**Alerts:**
1. Settings > Alerts
2. Configure:
   - Storage > 80%
   - Slow queries > 10/min
   - Connection errors > 5/min

### Query Analytics

**Enable in database:**
```sql
-- Check slow query log status
SHOW VARIABLES LIKE 'slow_query_log';

-- Queries taking >1 second
SELECT * FROM mysql.slow_log 
WHERE query_time > 1 
ORDER BY start_time DESC 
LIMIT 20;
```

**Monitoring Script:**
```bash
#!/bin/bash
# check-db-health.sh

DB_HOST="your-db-host"
DB_USER="your-user"
DB_PASS="your-pass"
DB_NAME="oxlook"

# Check connection count
CONN_COUNT=$(mysql -h$DB_HOST -u$DB_USER -p$DB_PASS $DB_NAME \
  -e "SHOW STATUS LIKE 'Threads_connected';" | grep Threads | awk '{print $2}')

if [ $CONN_COUNT -gt 20 ]; then
  echo "WARNING: Too many connections: $CONN_COUNT"
  # Send alert
fi

# Check slow queries
SLOW_QUERIES=$(mysql -h$DB_HOST -u$DB_USER -p$DB_PASS $DB_NAME \
  -e "SHOW STATUS LIKE 'Slow_queries';" | grep Slow | awk '{print $2}')

echo "Connections: $CONN_COUNT, Slow queries: $SLOW_QUERIES"
```

Run every 5 minutes:
```bash
chmod +x check-db-health.sh
# Add to cron: */5 * * * * /path/to/check-db-health.sh
```

---

## 5. Performance Monitoring

### Real User Monitoring (RUM)

**Cloudflare Web Analytics:**

1. **Setup**
   - Cloudflare Dashboard > Analytics > Web Analytics
   - Add site: bhomanshah.com
   - Copy JavaScript snippet

2. **Add to Frontend** (index.html):
   ```html
   <!-- Before closing </body> -->
   <script defer 
     src='https://static.cloudflareinsights.com/beacon.min.js' 
     data-cf-beacon='{"token": "YOUR_TOKEN"}'></script>
   ```

3. **Metrics Tracked:**
   - Page load time
   - Core Web Vitals (LCP, CLS, FID)
   - Browser/device breakdown
   - Geographic distribution
   - Traffic sources

**Alternative: Google Analytics 4**
```html
<!-- Global site tag (gtag.js) -->
<script async src="https://www.googletagmanager.com/gtag/js?id=G-XXXXXXXXXX"></script>
<script>
  window.dataLayer = window.dataLayer || [];
  function gtag(){dataLayer.push(arguments);}
  gtag('js', new Date());
  gtag('config', 'G-XXXXXXXXXX');
</script>
```

### Application Performance Monitoring (APM)

**New Relic (Optional - Paid):**

```bash
go get github.com/newrelic/go-agent/v3/newrelic
```

```go
import "github.com/newrelic/go-agent/v3/newrelic"

app, err := newrelic.NewApplication(
    newrelic.ConfigAppName("OXLOOK API"),
    newrelic.ConfigLicense(os.Getenv("NEW_RELIC_LICENSE_KEY")),
)

// Instrument handlers
txn := app.StartTransaction("create-order")
defer txn.End()
```

---

## 6. Alerting Channels

### Slack Integration

1. **Create Slack Webhook**
   - Slack > Apps > Incoming Webhooks
   - Add to channel: #alerts
   - Copy webhook URL

2. **Alert Script**
   ```bash
   #!/bin/bash
   # send-slack-alert.sh
   
   WEBHOOK_URL="https://hooks.slack.com/services/YOUR/WEBHOOK/URL"
   MESSAGE="$1"
   
   curl -X POST $WEBHOOK_URL \
     -H 'Content-Type: application/json' \
     -d "{\"text\":\"🚨 OXLOOK Alert: $MESSAGE\"}"
   ```

3. **Use in Monitoring**
   ```bash
   # When error detected
   ./send-slack-alert.sh "API is down!"
   ```

### Email Alerts

**Using SendGrid:**

```go
// Send alert email
import (
    "github.com/sendgrid/sendgrid-go"
    "github.com/sendgrid/sendgrid-go/helpers/mail"
)

func sendAlert(subject, message string) {
    from := mail.NewEmail("OXLOOK Alerts", "alerts@yourdomain.com")
    to := mail.NewEmail("Admin", "admin@yourdomain.com")
    content := mail.NewContent("text/plain", message)
    m := mail.NewV3MailInit(from, subject, to, content)
    
    client := sendgrid.NewSendClient(os.Getenv("SENDGRID_API_KEY"))
    client.Send(m)
}
```

### SMS Alerts (Optional)

**Using Twilio:**

```go
import (
    "github.com/twilio/twilio-go"
    twilioApi "github.com/twilio/twilio-go/rest/api/v2010"
)

func sendSMS(message string) {
    client := twilio.NewRestClient()
    params := &twilioApi.CreateMessageParams{}
    params.SetTo(os.Getenv("ALERT_PHONE"))
    params.SetFrom(os.Getenv("TWILIO_PHONE"))
    params.SetBody(message)
    
    client.Api.CreateMessage(params)
}
```

---

## 7. Health Check Endpoint Enhancement

Enhance `/health` to include system checks:

```go
// In handlers.go
func (h *Handler) HealthCheck(c fiber.Ctx) error {
    health := fiber.Map{
        "status": "ok",
        "timestamp": time.Now().Unix(),
    }
    
    // Check database
    if err := h.db.Ping(); err != nil {
        health["status"] = "degraded"
        health["database"] = "down"
        return c.Status(503).JSON(health)
    }
    health["database"] = "ok"
    
    // Check memory
    var m runtime.MemStats
    runtime.ReadMemStats(&m)
    health["memory_mb"] = m.Alloc / 1024 / 1024
    
    return c.JSON(health)
}
```

---

## 8. Monitoring Dashboard

### Grafana Cloud (Free Tier)

1. **Sign up:** https://grafana.com/products/cloud/
2. **Create Dashboard**
3. **Add Panels:**
   - Uptime percentage (last 7 days)
   - Request rate
   - Error rate
   - Response time
   - Database connections

### Create Custom Dashboard

Simple HTML dashboard:

```html
<!DOCTYPE html>
<html>
<head>
  <title>OXLOOK Status</title>
  <meta http-equiv="refresh" content="30">
</head>
<body>
  <h1>OXLOOK Status</h1>
  <div id="status"></div>
  
  <script>
    async function checkStatus() {
      try {
        const res = await fetch('https://api.bhomanshah.com/health');
        const data = await res.json();
        document.getElementById('status').innerHTML = 
          `✅ API Status: ${data.status}<br>` +
          `Database: ${data.database}<br>` +
          `Memory: ${data.memory_mb} MB`;
      } catch (err) {
        document.getElementById('status').innerHTML = 
          `❌ API is DOWN!<br>${err.message}`;
      }
    }
    
    checkStatus();
    setInterval(checkStatus, 30000);
  </script>
</body>
</html>
```

---

## 9. Alert Escalation

**Severity Levels:**

1. **P1 - Critical** (Immediate response)
   - API completely down
   - Database unavailable
   - Data loss risk
   - Action: SMS + Slack + Email

2. **P2 - High** (15 min response)
   - Error rate > 5%
   - Response time > 5s
   - Single endpoint down
   - Action: Slack + Email

3. **P3 - Medium** (1 hour response)
   - Error rate > 1%
   - Slow queries
   - Storage > 80%
   - Action: Email

4. **P4 - Low** (Next business day)
   - Performance degradation
   - Cache miss rate high
   - Action: Email digest

---

## 10. Monitoring Checklist

### Daily
- [ ] Check error rate
- [ ] Review new orders
- [ ] Check response times
- [ ] Verify uptime

### Weekly
- [ ] Review slow queries
- [ ] Check database size
- [ ] Review security logs
- [ ] Test backup restoration

### Monthly
- [ ] Review traffic patterns
- [ ] Analyze error trends
- [ ] Update dependencies
- [ ] Security audit
- [ ] Cost optimization

---

## 11. Incident Response

**When Alert Fires:**

1. **Acknowledge**
   - Acknowledge alert in monitoring tool
   - Notify team

2. **Investigate**
   - Check Koyeb dashboard
   - Review recent deployments
   - Check database status
   - Review logs

3. **Mitigate**
   - Rollback if needed
   - Scale up if needed
   - Fix critical issues

4. **Communicate**
   - Update status page
   - Notify affected users
   - Document incident

5. **Post-Mortem**
   - Root cause analysis
   - Preventive measures
   - Update runbook

---

## 12. Cost Summary

### Free Options
- **Uptime Monitoring:** UptimeRobot (50 monitors)
- **Error Tracking:** Sentry (5k events/month)
- **Analytics:** Cloudflare Web Analytics
- **Logs:** Koyeb built-in
- **Metrics:** Koyeb built-in

### Paid Options
- **Sentry Team:** $26/month
- **Pingdom:** $15/month
- **New Relic:** $99/month
- **Datadog:** $15/host/month
- **PagerDuty:** $25/user/month

**Recommendation:** Start with free tier of all services, upgrade as needed.

---

## Next Steps

1. Set up UptimeRobot
2. Configure Sentry (optional)
3. Add Cloudflare Web Analytics
4. Create Slack webhook for alerts
5. Set up email alerts
6. Document incident response plan
7. Test alert system
