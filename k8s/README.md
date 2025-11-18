# Kubernetes Deployment

This directory contains Kubernetes manifests for deploying the ecommerce API.

## Prerequisites

- Kubernetes cluster (1.20+)
- kubectl configured
- Docker image built and pushed to a registry

## Files

- `namespace.yaml` - Creates the ecommerce namespace
- `configmap.yaml` - Non-sensitive configuration
- `secret.yaml` - Sensitive configuration (passwords, secrets)
- `deployment.yaml` - Application deployment configuration
- `service.yaml` - Service to expose the application
- `hpa.yaml` - Horizontal Pod Autoscaler for automatic scaling

## Deployment Steps

1. **Create the namespace:**
   ```bash
   kubectl apply -f namespace.yaml
   ```

2. **Create secrets (update values first!):**
   ```bash
   # Edit secret.yaml with your actual secrets
   kubectl apply -f secret.yaml
   ```

3. **Create ConfigMap:**
   ```bash
   # Edit configmap.yaml with your configuration
   kubectl apply -f configmap.yaml
   ```

4. **Deploy the application:**
   ```bash
   kubectl apply -f deployment.yaml
   ```

5. **Create the service:**
   ```bash
   kubectl apply -f service.yaml
   ```

6. **Enable autoscaling:**
   ```bash
   kubectl apply -f hpa.yaml
   ```

## Apply All Resources

To apply all resources at once:
```bash
kubectl apply -f k8s/
```

## Verify Deployment

```bash
# Check pods
kubectl get pods -n ecommerce

# Check service
kubectl get svc -n ecommerce

# Check HPA
kubectl get hpa -n ecommerce

# View logs
kubectl logs -f -n ecommerce -l app=ecommerce-api

# Check pod details
kubectl describe pod <pod-name> -n ecommerce
```

## Update Deployment

When you have a new image:

```bash
# Update the image
kubectl set image deployment/ecommerce-api api=ecommerce-api:new-version -n ecommerce

# Or edit the deployment directly
kubectl edit deployment ecommerce-api -n ecommerce

# Check rollout status
kubectl rollout status deployment/ecommerce-api -n ecommerce

# Rollback if needed
kubectl rollout undo deployment/ecommerce-api -n ecommerce
```

## Secrets Management

For production, consider using:
- **Sealed Secrets**: Encrypt secrets in Git
- **External Secrets Operator**: Sync from external secret stores
- **Vault**: HashiCorp Vault integration
- **Cloud Provider Secrets**: AWS Secrets Manager, GCP Secret Manager, Azure Key Vault

## Health Checks

The deployment includes:
- **Liveness Probe**: `/health` - Checks if the app is alive
- **Readiness Probe**: `/ready` - Checks if the app is ready to serve traffic

## Scaling

### Manual Scaling
```bash
kubectl scale deployment ecommerce-api --replicas=5 -n ecommerce
```

### Automatic Scaling
The HPA will automatically scale based on CPU and memory usage (70% and 80% respectively).

## Resource Limits

Current resource configuration per pod:
- **Requests**: 100m CPU, 128Mi memory
- **Limits**: 500m CPU, 256Mi memory

Adjust based on your application's needs.

## Cleanup

To remove all resources:
```bash
kubectl delete -f k8s/
```

Or delete the namespace (removes everything in it):
```bash
kubectl delete namespace ecommerce
```
