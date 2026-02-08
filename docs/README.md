# Deployment Documentation

## Overview
This directory contains all documentation needed to deploy OXLOOK to production.

---

## Quick Start

Follow these guides in order:

1. **[DATABASE_SETUP.md](DATABASE_SETUP.md)** - Set up production database
2. **[KOYEB_DEPLOYMENT.md](KOYEB_DEPLOYMENT.md)** - Deploy API to Koyeb
3. **[PRODUCTION_TESTING.md](PRODUCTION_TESTING.md)** - Test your deployment
4. **[MONITORING.md](MONITORING.md)** - Set up monitoring and alerts
5. **[DEPLOYMENT_CHECKLIST.md](DEPLOYMENT_CHECKLIST.md)** - Final checklist

---

## Document Summaries

### DATABASE_SETUP.md
Complete guide to choosing and setting up a production database:
- **Options:** PlanetScale (recommended), Neon, Railway, Aiven
- **Configuration:** Connection strings, SSL/TLS, backups
- **Maintenance:** Monitoring, backups, security
- **Recommendation:** PlanetScale free tier for small-medium stores

### KOYEB_DEPLOYMENT.md
Step-by-step guide to deploying the API on Koyeb:
- **Setup:** Account creation, repository connection
- **Configuration:** Environment variables, health checks
- **Custom Domain:** DNS setup with Cloudflare
- **Monitoring:** Logs, metrics, scaling
- **Troubleshooting:** Common issues and solutions

### PRODUCTION_TESTING.md
Comprehensive testing guide for production:
- **API Testing:** All endpoints, rate limiting, security
- **Frontend Testing:** Cross-browser, performance, Lighthouse
- **Load Testing:** Apache Bench, k6 examples
- **Security Testing:** SSL, headers, penetration testing basics
- **Monitoring:** Uptime checks, error tracking

### MONITORING.md
Complete monitoring and alerting setup:
- **Uptime Monitoring:** UptimeRobot, Pingdom, Checkly
- **Error Tracking:** Sentry integration
- **Performance:** RUM with Cloudflare Web Analytics
- **Alerts:** Slack, email, SMS integrations
- **Incident Response:** Escalation procedures

### DEPLOYMENT_CHECKLIST.md
Master checklist for production deployment:
- **Pre-Deployment:** Code prep, documentation
- **Database:** Setup and verification
- **Environment Variables:** Generation and configuration
- **Deployment:** Koyeb and Cloudflare Pages
- **Testing:** Full end-to-end verification
- **Post-Launch:** Monitoring and troubleshooting

---

## Additional Resources

### From Project Root

- **[PERFORMANCE.md](../PERFORMANCE.md)** - Frontend optimization guide
- **[DEPLOYMENT.md](../DEPLOYMENT.md)** - Cloudflare Pages deployment
- **[README.md](../README.md)** - Project overview
- **[plan.md](../plan.md)** - Full implementation plan

---

## Architecture

```
┌─────────────────┐
│   Cloudflare    │
│   Pages         │  Frontend (HTML/CSS/JS)
│ bhomanshah.com  │  Static files
└────────┬────────┘
         │ HTTPS
         │
         ▼
┌─────────────────┐
│     Koyeb       │
│   (Docker)      │  Go API
│ api.bhomanshah  │  Fiber v3
│      .com       │
└────────┬────────┘
         │ MySQL/TLS
         │
         ▼
┌─────────────────┐
│  PlanetScale/   │
│     Neon        │  Database
│    (MySQL)      │  Managed
└─────────────────┘
```

---

## Environment Variables

### Required for API (Koyeb)

```bash
DATABASE_URL=mysql://user:pass@host:port/db?tls=true
ADMIN_SECRET=<generate-with-openssl-rand>
CORS_ORIGIN=https://bhomanshah.com
ENVIRONMENT=production
PORT=8080
```

### Optional for Monitoring

```bash
SENTRY_DSN=https://xxx@xxx.ingest.sentry.io/xxx
NEW_RELIC_LICENSE_KEY=xxx
SENDGRID_API_KEY=xxx
```

---

## Deployment Flow

### 1. Database Setup (30 minutes)
1. Create PlanetScale account
2. Create database
3. Get connection string
4. Test connection

### 2. Koyeb Deployment (45 minutes)
1. Create Koyeb account
2. Connect GitHub repository
3. Configure environment variables
4. Deploy service
5. Wait for build and migrations

### 3. Custom Domain (15 minutes)
1. Add domain in Koyeb
2. Configure DNS in Cloudflare
3. Wait for SSL provisioning
4. Verify HTTPS working

### 4. Frontend Deployment (30 minutes)
1. Run build script
2. Create Cloudflare Pages project
3. Configure custom domain
4. Enable optimizations
5. Deploy

### 5. Testing (1 hour)
1. Test all API endpoints
2. Test frontend integration
3. Verify end-to-end flow
4. Run Lighthouse audit
5. Check security headers

### 6. Monitoring (30 minutes)
1. Set up UptimeRobot
2. Configure Cloudflare Analytics
3. (Optional) Set up Sentry
4. Configure alert channels

**Total Time: ~3-4 hours**

---

## Support

### Getting Help

1. **Check documentation** - Most issues covered in guides
2. **Review logs** - Koyeb dashboard > Logs
3. **Test locally** - Reproduce issue in development
4. **Community** - Koyeb/Cloudflare community forums
5. **Support** - Contact provider support if needed

### Common Issues

**Build fails:**
- Check Dockerfile syntax
- Verify all files committed to Git
- Check build logs for errors

**Database connection fails:**
- Verify DATABASE_URL format
- Check SSL/TLS settings
- Ensure database allows Koyeb IP ranges

**CORS errors:**
- Verify CORS_ORIGIN matches frontend URL exactly
- Include https:// protocol
- Check CORS middleware loaded

**SSL certificate not provisioning:**
- Ensure Cloudflare proxy is OFF for API subdomain
- Wait 10 minutes for provisioning
- Check DNS propagation

---

## Next Steps After Deployment

1. **Monitor** - Watch logs and metrics for first 24 hours
2. **Test** - Have friends/colleagues test the site
3. **Gather Feedback** - Note any issues or improvements
4. **Iterate** - Make improvements based on feedback
5. **Scale** - Upgrade resources as traffic grows
6. **Market** - Promote your online store!

---

## Security Reminders

- ✅ Never commit `.env` file
- ✅ Use Koyeb secrets for sensitive data
- ✅ Rotate ADMIN_SECRET periodically
- ✅ Monitor logs for suspicious activity
- ✅ Keep dependencies updated
- ✅ Use strong database passwords
- ✅ Enable 2FA on all services
- ✅ Regular security audits

---

## Cost Breakdown

### Free Tier (Suitable for Small Stores)

- **Koyeb:** Free (512MB RAM, 2M requests/month)
- **PlanetScale:** Free (5GB storage, 1B reads/month)
- **Cloudflare Pages:** Free (Unlimited requests)
- **Cloudflare DNS:** Free
- **UptimeRobot:** Free (50 monitors)
- **Cloudflare Analytics:** Free

**Total: $0/month** ✅

### When to Upgrade

- Traffic > 2M requests/month
- Database > 5GB
- Need auto-scaling
- Require advanced features
- Want premium support

---

## Maintenance Schedule

### Daily
- Check error logs
- Verify order processing

### Weekly
- Review slow queries
- Check disk usage
- Security log review

### Monthly
- Update dependencies
- Review costs
- Performance optimization
- Security audit

### Quarterly
- Backup testing
- Disaster recovery drill
- Architecture review
- Capacity planning

---

## Success! 🎉

Once deployed, your OXLOOK store will be:

- ✅ Live at https://bhomanshah.com
- ✅ API at https://api.bhomanshah.com
- ✅ Secure with HTTPS and security headers
- ✅ Fast with CDN caching
- ✅ Monitored with uptime checks
- ✅ Scalable with cloud infrastructure

Happy selling! 🛍️
