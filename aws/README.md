# Particles AWS Deployment

AWS CDK deployment for the particles SSH simulation using ECS.

## Deploy

1. Deploy ECR repository:
   ```bash
   cdk deploy ParticlesEcrStack
   ```

2. Build and push Docker image:
   ```bash
   ../deploy.sh
   ```

3. Deploy application:
   ```bash
   cdk deploy ParticlesStack --parameters ImageTag=latest
   ```

4. Connect:
   ```bash
   ssh <nlb-dns-name>
   ```

## Clean Up

```bash
cdk destroy ParticlesStack
cdk destroy ParticlesEcrStack
```
