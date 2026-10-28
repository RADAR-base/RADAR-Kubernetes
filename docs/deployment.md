# RADAR-Kubernetes Deployment Commands

| /!\ Version compatibility /!\ |
| ---------------------------------------- |
| This project requires Helmfile v0.169.1. Using newer versions like v1.0.0 WILL cause template processing issues, especially with the environments.yaml file. |

This document provides a reference for the main commands used to deploy the RADAR-base platform using Helmfile.

## Deployment

### Full Deployment

To deploy the complete RADAR-Kubernetes stack:

```bash
helmfile sync --concurrency 1 --wait
```

**Parameters explained:**
- `--concurrency 1`: Processes one chart at a time to avoid race conditions and dependency issues
- `--wait`: Waits for each release to be fully deployed before proceeding to the next one

This command is preferred for initial installations and updates as it ensures proper dependency resolution between components.

### Selective Deployment

To deploy specific components:

```bash
helmfile -f helmfile.yaml -e production apply --selector name=<component-name>
```

For example, to deploy only the nginx ingress controller:
```bash
helmfile -f helmfile.yaml -e production apply --selector name=nginx-ingress
```

## Cleanup

### Full Cleanup

To completely remove the RADAR-Kubernetes installation:

```bash
helmfile destroy
```

This command will uninstall all deployed Helm charts and delete their associated Kubernetes resources.

**Warning:** This will remove all components including data storage. If you want to preserve data, ensure you back up persistent volumes before running this command.

### Selective Cleanup

To remove specific components:

```bash
helmfile -f helmfile.yaml -e production destroy --selector name=<component-name>
```

## Troubleshooting

If you encounter validation webhook issues with the nginx-ingress controller during deployment, you may need to:

1. Update the nginx-ingress configuration in `production.yaml` to enable admission webhooks:
   ```yaml
   nginx_ingress:
     controller:
       admissionWebhooks:
         enabled: true
         failurePolicy: Ignore
         timeoutSeconds: 30
   ```

2. Reinstall the nginx-ingress controller:
   ```bash
   helmfile -f helmfile.yaml -e production apply --selector name=nginx-ingress
   ```

## Important Notes

1. Always use `--concurrency 1` for reliable installations, especially with interdependent components
2. The TLS configuration format in `production.yaml` may vary depending on the chart:
   - Some components expect array format: `tls: [{ secretName: ... }]`
   - Others expect object format: `tls: { secretName: ... }`
3. After a full destroy, a new installation must start from scratch, requiring all namespace creation steps 