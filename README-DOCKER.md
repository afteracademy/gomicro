# Docker Compose Setup Guide

## [Structure] Project Structure

```
gomicro/
├── docker-compose.yml               # Full stack development setup
├── docker-compose-load-balanced.yml # Production load-balanced setup
├── auth_service/
│   ├── Dockerfile
│   ├── .env
│   └── migrations/
└── blog_service/
    ├── Dockerfile
    └── .env
```

## [Quick Start]

### Development Mode (Single Instances)

Start the complete microservices stack:

```bash
# Start all services
docker compose up

# Build and start
docker compose up --build

# Run in background
docker compose up -d

# View logs
docker compose logs -f

# Stop all services
docker compose down
```

### Production Mode (Load Balanced)

Run with 2 instances of each service for load balancing:

```bash
# Start load-balanced stack
docker compose -f docker-compose-load-balanced.yml up

# Build and start
docker compose -f docker-compose-load-balanced.yml up --build

# Run in background
docker compose -f docker-compose-load-balanced.yml up -d

# Stop
docker compose -f docker-compose-load-balanced.yml down
```

## [Access] Service Access Points

### Development Mode
- **API Gateway (Kong)**: http://localhost:8000
- **Kong Admin API**: http://localhost:8001
- **NATS Client**: localhost:4222
- **NATS Monitoring**: http://localhost:8222
- **PostgreSQL**: localhost:5432
- **MongoDB**: localhost:27017
- **Redis (Auth)**: localhost:6379
- **Redis (Blog)**: localhost:6380

### Load Balanced Mode
Same access points, but traffic is distributed across:
- 2x Auth service instances (auth1, auth2)
- 2x Blog service instances (blog1, blog2)

## [Architecture]

### Network Topology

```
Internet
   ↓
Kong API Gateway (:8000)
   ├─→ Auth Service Instances → PostgreSQL + Redis (Auth)
   ├─→ Blog Service Instances → MongoDB + Redis (Blog)
   └─→ NATS (Message Queue)
```

**Key Features:**
- All services communicate via `gomicro-network` bridge network
- Services use container names for discovery (postgres, mongo, redis-auth, redis-blog, nats)
- Kong handles load balancing across service instances
- Services are not exposed externally (only Kong gateway)
- Each service has isolated Redis cache

### Service Dependencies

**Auth Service:**
- PostgreSQL (database)
- Redis (auth) - caching and session storage
- NATS - message broker

**Blog Service:**
- MongoDB (database)
- Redis (blog) - caching
- NATS - message broker

## [Operations] Common Operations

### Service Management

```bash
# View running services
docker compose ps

# View logs for specific service
docker compose logs -f auth
docker compose logs -f blog
docker compose logs -f kong

# Restart a service
docker compose restart auth
docker compose restart postgres

# Stop specific service
docker compose stop blog

# Remove all containers (keeps volumes)
docker compose down

# Remove all containers and volumes [WARNING: deletes data]
docker compose down -v
```

### Scaling Services

```bash
# Development mode - scale manually
docker compose up --scale auth=3 --scale blog=2 -d

# Production mode already has 2 instances per service
# To scale further, edit docker-compose-load-balanced.yml and add auth3, blog3, etc.
```

## [Database] Database Management

### PostgreSQL (Auth Service)

```bash
# Connect to database
docker compose exec postgres psql -U ${DB_USER} -d ${DB_NAME}

# Run migrations (automatically runs on startup)
docker compose run --rm migrate-auth

# Force re-run migrations
docker compose up migrate-auth

# Backup database
docker compose exec postgres pg_dump -U ${DB_USER} ${DB_NAME} > backup_$(date +%Y%m%d).sql

# Restore database
cat backup_20260122.sql | docker compose exec -T postgres psql -U ${DB_USER} ${DB_NAME}

# Check database health
docker compose exec postgres pg_isready -U ${DB_USER} -d ${DB_NAME}
```

### MongoDB (Blog Service)

```bash
# Connect to database
docker compose exec mongo mongosh -u ${DB_ADMIN} -p ${DB_ADMIN_PWD}

# Connect to specific database
docker compose exec mongo mongosh -u ${DB_ADMIN} -p ${DB_ADMIN_PWD} --authenticationDatabase admin blog_db

# Backup database
docker compose exec mongo mongodump --username=${DB_ADMIN} --password=${DB_ADMIN_PWD} --authenticationDatabase=admin --db=blog_db --out=/tmp/backup
docker compose exec mongo tar -czf /tmp/blog_backup_$(date +%Y%m%d).tar.gz /tmp/backup

# Restore database
docker compose exec mongo mongorestore --username=${DB_ADMIN} --password=${DB_ADMIN_PWD} --authenticationDatabase=admin /tmp/backup
```

### Redis Cache Management

```bash
# Connect to Auth Redis
docker compose exec redis-auth redis-cli -a ${REDIS_PASSWORD}

# Connect to Blog Redis
docker compose exec redis-blog redis-cli -a ${REDIS_PASSWORD}

# Monitor Redis activity
docker compose exec redis-auth redis-cli -a ${REDIS_PASSWORD} MONITOR

# Check cache stats
docker compose exec redis-auth redis-cli -a ${REDIS_PASSWORD} INFO stats

# Clear cache [WARNING: use with caution]
docker compose exec redis-auth redis-cli -a ${REDIS_PASSWORD} FLUSHALL
docker compose exec redis-blog redis-cli -a ${REDIS_PASSWORD} FLUSHALL
```

## [Configuration] Environment Configuration

### Required Environment Files

**Root `.env`** (optional - for NATS ports)
```env
NATS_CLIENT_PORT=4222
NATS_MANAGEMENT_PORT=8222
```

**`auth_service/.env`**
```env
# Database
DB_NAME=auth_db
DB_USER=auth_user
DB_USER_PWD=your_secure_password

# Redis
REDIS_PASSWORD=your_redis_password

# Application
PORT=8080
JWT_SECRET=your_jwt_secret
API_KEY=your_api_key
```

**`blog_service/.env`**
```env
# MongoDB
DB_ADMIN=admin
DB_ADMIN_PWD=your_mongo_admin_password
DB_NAME=blog_db
DB_HOST=mongo
DB_PORT=27017

# Redis
REDIS_PASSWORD=your_redis_password

# Application
PORT=8080
AUTH_SERVICE_URL=http://auth:8080
```

## [Debugging] Debugging & Troubleshooting

### Viewing Logs

```bash
# All services
docker compose logs -f

# Specific service
docker compose logs -f auth
docker compose logs -f kong

# Last 100 lines
docker compose logs --tail=100 blog

# Load balanced mode
docker compose -f docker-compose-load-balanced.yml logs -f auth1 auth2 blog1 blog2
```

### Service Health Checks

```bash
# Check all services status
docker compose ps

# Check specific service health
docker compose exec postgres pg_isready
docker compose exec redis-auth redis-cli ping
docker compose exec mongo mongosh --eval "db.adminCommand('ping')"

# Test service connectivity
docker compose exec auth nc -zv postgres 5432
docker compose exec blog nc -zv mongo 27017
docker compose exec blog nc -zv redis-blog 6379
```

### Common Issues

#### Port Already in Use
```bash
# Find process using port
lsof -i :8000
lsof -i :5432

# Kill process
kill -9 <PID>

# Or change port in .env file
```

#### Service Won't Start
```bash
# Check logs
docker compose logs [service_name]

# Rebuild container
docker compose up --build [service_name]

# Remove and recreate
docker compose down
docker compose up --build
```

#### Database Connection Failed
```bash
# Verify database is running
docker compose ps postgres mongo

# Check database logs
docker compose logs postgres
docker compose logs mongo

# Verify environment variables
docker compose exec auth env | grep DB_
docker compose exec blog env | grep DB_

# Test connection from service
docker compose exec auth ping postgres
docker compose exec blog ping mongo
```

#### Migration Failures
```bash
# Check migration logs
docker compose logs migrate-auth

# Force re-run migrations
docker compose down migrate-auth
docker compose up migrate-auth

# Manual migration
docker compose exec postgres psql -U ${DB_USER} -d ${DB_NAME} -f /path/to/migration.sql
```

#### Kong Configuration Issues
```bash
# Check Kong status
curl http://localhost:8001/status

# List configured services
curl http://localhost:8001/services

# List configured routes
curl http://localhost:8001/routes

# Check upstream targets (load balanced mode)
curl http://localhost:8001/upstreams
curl http://localhost:8001/upstreams/{upstream_name}/targets
```

### Clean Slate

```bash
# Stop all services
docker compose down

# Remove all containers, networks, volumes [WARNING: deletes all data]
docker compose down -v --remove-orphans

# Remove all images
docker compose down --rmi all

# Complete cleanup
docker compose down -v --remove-orphans --rmi all
docker system prune -a --volumes -f
```

## [Development] Development Workflow

### Building Services

```bash
# Build all services
docker compose build

# Build specific service
docker compose build auth
docker compose build blog

# Build without cache
docker compose build --no-cache

# Pull latest base images before building
docker compose build --pull
```

### Debugging Inside Containers

```bash
# Execute shell in running container
docker compose exec auth sh
docker compose exec blog sh

# Execute as root
docker compose exec -u root auth sh

# Run one-off command
docker compose exec auth ls -la /app
docker compose exec blog env
```

### Monitoring Resources

```bash
# View resource usage
docker stats

# View specific services
docker stats gomicro-auth-1 gomicro-blog-1

# Export stats to file
docker stats --no-stream > stats.txt
```

## [Deployment] Production Deployment

### Load Balanced Setup

The [docker-compose-load-balanced.yml](docker-compose-load-balanced.yml) file provides:
- 2x Auth service instances
- 2x Blog service instances
- Kong automatically load balances requests
- Shared databases and caches
- Health checks for all services

### Scaling Strategy

**Vertical Scaling (increase instance resources):**
```yaml
# Add to docker-compose-load-balanced.yml
services:
  auth1:
    deploy:
      resources:
        limits:
          cpus: '2'
          memory: 2G
        reservations:
          cpus: '1'
          memory: 1G
```

**Horizontal Scaling (add more instances):**
1. Edit [docker-compose-load-balanced.yml](docker-compose-load-balanced.yml)
2. Add `auth3`, `blog3`, etc. services
3. Update Kong's upstream configuration in [kong-load-balanced.yml](kong/kong-load-balanced.yml)
4. Restart stack

### Health Monitoring

```bash
# Check health status
docker compose ps

# Services with health checks:
# - postgres: pg_isready
# - redis-auth: redis-cli ping
# - redis-blog: redis-cli ping

# View health check logs
docker inspect --format='{{json .State.Health}}' gomicro-postgres-1 | jq
```

### Backup Strategy

```bash
# Automated backup script (save as backup.sh)
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="./backups/$DATE"
mkdir -p $BACKUP_DIR

# Backup PostgreSQL
docker compose exec postgres pg_dump -U ${DB_USER} ${DB_NAME} > $BACKUP_DIR/postgres.sql

# Backup MongoDB
docker compose exec mongo mongodump --username=${DB_ADMIN} --password=${DB_ADMIN_PWD} --out=$BACKUP_DIR/mongo

# Backup volumes
docker run --rm -v gomicro_postgres-data:/data -v $(pwd)/$BACKUP_DIR:/backup alpine tar czf /backup/postgres-data.tar.gz -C /data .
docker run --rm -v gomicro_mongo-data:/data -v $(pwd)/$BACKUP_DIR:/backup alpine tar czf /backup/mongo-data.tar.gz -C /data .

echo "Backup completed: $BACKUP_DIR"
```

## [Best Practices]

### Development
- Use `docker-compose.yml` for local development
- Keep environment files outside version control (.env in .gitignore)
- Use volume mounts for hot reloading during development
- Expose database ports for local database clients
- Use meaningful container names
- Tag images with version numbers

### Production
- Use `docker-compose-load-balanced.yml` for production
- Set restart policies to `unless-stopped` or `always`
- Implement health checks for all services
- Use secrets management (Docker Secrets, Vault)
- Monitor resource usage and set limits
- Regular backups of databases and volumes
- Use specific image versions (not `latest`)
- Implement proper logging and monitoring

### Security
- Never commit `.env` files to version control
- Use strong passwords for all services
- Change default credentials
- Run containers as non-root users
- Keep images updated
- Use private registries for custom images
- Implement network segmentation
- Enable TLS/SSL in production

## [Resources] Additional Resources

### Kong Configuration
- Kong Admin API: http://localhost:8001
- Configure services and routes through Kong Admin API
- Load balancing configured in [kong-load-balanced.yml](kong/kong-load-balanced.yml)

### Service Communication
- NATS for async messaging between services
- Services communicate via container names
- Kong handles external API routing

### Next Steps
1. [ ] Configure Kong routes for your API endpoints
2. [ ] Set up service-to-service authentication
3. [ ] Implement API rate limiting in Kong
4. [ ] Add monitoring (Prometheus/Grafana)
5. [ ] Set up CI/CD pipeline
6. [ ] Configure SSL certificates for production
7. [ ] Implement centralized logging (ELK/Loki)

## [Help] Getting Help

**Check logs first:**
```bash
docker compose logs -f [service_name]
```

**Common commands:**
```bash
# Full restart
docker compose down && docker compose up --build

# Reset everything
docker compose down -v && docker compose up --build

# Check what's running
docker compose ps
```

For more details, refer to:
- [Kong Documentation](https://docs.konghq.com/)
- [Docker Compose Documentation](https://docs.docker.com/compose/)
- [NATS Documentation](https://docs.nats.io/)
