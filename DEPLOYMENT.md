# Deployment Guide

This guide covers deploying the Template-GoLang microservice to production environments.

## Prerequisites

- Container orchestration platform (Kubernetes, Docker Swarm, etc.)
- Access to container registry
- Kafka cluster (if using Kafka features)
- Jaeger instance (if using distributed tracing)
- Prometheus instance (for metrics collection)

## Building for Production

### Build the Docker Image

```bash
# Build the image
docker build -t template-golang:latest .

# Tag for your registry
docker tag template-golang:latest your-registry.com/template-golang:v1.0.0

# Push to registry
docker push your-registry.com/template-golang:v1.0.0
```

### Multi-Architecture Builds

```bash
docker buildx build --platform linux/amd64,linux/arm64 \
  -t your-registry.com/template-golang:v1.0.0 \
  --push .
```

## Environment Configuration

### Required Environment Variables

```bash
# Server
PORT=8080
HOST=0.0.0.0
READ_TIMEOUT=10
WRITE_TIMEOUT=10
SHUTDOWN_TIMEOUT=30
ENVIRONMENT=production

# Logging
LOG_LEVEL=info

# Tracing (Optional)
TRACING_ENABLED=true
JAEGER_ENDPOINT=http://jaeger-collector:14268/api/traces
SERVICE_NAME=template-service

# Kafka (Optional)
KAFKA_ENABLED=true
KAFKA_BROKERS=kafka-1:9092,kafka-2:9092,kafka-3:9092
KAFKA_TOPIC=items
KAFKA_GROUP_ID=template-service

# Security
JWT_SECRET=your-secret-key-here-use-strong-random-value
```

### TLS/HTTPS Configuration

```bash
TLS_ENABLED=true
TLS_CERT_FILE=/etc/tls/cert.pem
TLS_KEY_FILE=/etc/tls/key.pem
```

## Kubernetes Deployment

### Example Deployment Manifest

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: template-golang
  labels:
    app: template-golang
spec:
  replicas: 3
  selector:
    matchLabels:
      app: template-golang
  template:
    metadata:
      labels:
        app: template-golang
    spec:
      containers:
      - name: template-golang
        image: your-registry.com/template-golang:v1.0.0
        ports:
        - containerPort: 8080
          name: http
        env:
        - name: ENVIRONMENT
          value: "production"
        - name: LOG_LEVEL
          value: "info"
        - name: KAFKA_ENABLED
          value: "true"
        - name: KAFKA_BROKERS
          value: "kafka:9092"
        - name: TRACING_ENABLED
          value: "true"
        - name: JAEGER_ENDPOINT
          value: "http://jaeger-collector:14268/api/traces"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: template-golang-secrets
              key: jwt-secret
        resources:
          requests:
            memory: "128Mi"
            cpu: "100m"
          limits:
            memory: "512Mi"
            cpu: "500m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 10
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
---
apiVersion: v1
kind: Service
metadata:
  name: template-golang
  labels:
    app: template-golang
spec:
  type: ClusterIP
  ports:
  - port: 8080
    targetPort: 8080
    name: http
  selector:
    app: template-golang
```

### Secrets Configuration

```yaml
apiVersion: v1
kind: Secret
metadata:
  name: template-golang-secrets
type: Opaque
stringData:
  jwt-secret: "your-secret-key-here"
```

### ConfigMap for Additional Configuration

```yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: template-golang-config
data:
  config.yaml: |
    server:
      port: 8080
      read_timeout: 10
      write_timeout: 10
```

## Ingress Configuration

### Nginx Ingress

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: template-golang
  annotations:
    nginx.ingress.kubernetes.io/rewrite-target: /
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  ingressClassName: nginx
  tls:
  - hosts:
    - api.yourdomain.com
    secretName: template-golang-tls
  rules:
  - host: api.yourdomain.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: template-golang
            port:
              number: 8080
```

## Monitoring

### Prometheus ServiceMonitor

```yaml
apiVersion: monitoring.coreos.com/v1
kind: ServiceMonitor
metadata:
  name: template-golang
  labels:
    app: template-golang
spec:
  selector:
    matchLabels:
      app: template-golang
  endpoints:
  - port: http
    path: /metrics
    interval: 30s
```

### Grafana Dashboard

Import the Prometheus metrics into Grafana:

Key metrics to monitor:
- `http_requests_total`
- `http_request_duration_seconds`
- `items_created_total`
- `items_updated_total`
- `items_deleted_total`
- `kafka_messages_published_total`
- `kafka_messages_consumed_total`

## Horizontal Pod Autoscaling

```yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: template-golang
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: template-golang
  minReplicas: 2
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## Health Checks

The service exposes a `/health` endpoint that returns:

```json
{
  "status": "healthy",
  "version": "1.0.0"
}
```

Use this for:
- Kubernetes liveness probes
- Kubernetes readiness probes
- Load balancer health checks

## Logging

In production, logs are output in JSON format for easy parsing:

```json
{
  "level": "info",
  "time": "2024-01-01T12:00:00Z",
  "message": "HTTP request",
  "method": "GET",
  "path": "/v1/items",
  "status": 200,
  "duration": 15
}
```

Ship logs to your centralized logging system (ELK, Splunk, etc.).

## Security Best Practices

1. **Never commit secrets** - Use secret management (Kubernetes Secrets, Vault, etc.)
2. **Enable TLS** - Always use HTTPS in production
3. **Use strong JWT secrets** - Generate with: `openssl rand -base64 32`
4. **Limit resources** - Set CPU and memory limits
5. **Network policies** - Restrict pod-to-pod communication
6. **Image scanning** - Scan Docker images for vulnerabilities
7. **Keep dependencies updated** - Regularly update Go modules

## Backup and Disaster Recovery

- Ensure Kafka topics are replicated
- Back up configuration and secrets
- Document recovery procedures
- Test disaster recovery regularly

## Performance Tuning

### Go Runtime

```bash
# Set GOMAXPROCS if needed (defaults to number of CPUs)
GOMAXPROCS=4

# Tune garbage collector
GOGC=100
```

### Connection Pooling

Adjust in code for your workload:
- HTTP client timeouts
- Database connection pools
- Kafka producer/consumer settings

## Troubleshooting

### Common Issues

1. **Service not starting**
   - Check logs: `kubectl logs -f deployment/template-golang`
   - Verify environment variables
   - Check health endpoint

2. **High memory usage**
   - Review GOGC settings
   - Check for memory leaks
   - Analyze heap profiles

3. **Slow response times**
   - Check Prometheus metrics
   - Review distributed traces in Jaeger
   - Analyze database query performance

### Debug Mode

Never enable debug logging in production. If needed temporarily:

```bash
kubectl set env deployment/template-golang LOG_LEVEL=debug
```

Remember to revert:

```bash
kubectl set env deployment/template-golang LOG_LEVEL=info
```

## Rollback

If deployment fails:

```bash
kubectl rollout undo deployment/template-golang
```

Check rollout status:

```bash
kubectl rollout status deployment/template-golang
```

## Zero-Downtime Deployments

The service supports graceful shutdown. Kubernetes rolling updates will:

1. Start new pods
2. Wait for readiness probes
3. Route traffic to new pods
4. Send SIGTERM to old pods
5. Wait for graceful shutdown (default 30s)
6. Force kill if needed

## Support

For production issues:
- Check service logs
- Review Prometheus metrics
- Check Jaeger traces
- Consult runbooks
- Contact on-call team
