# Accessing Headlamp in RADAR-Kubernetes

Headlamp is a general-purpose, extensible Kubernetes web UI. It is deployed as the replacement for the legacy Kubernetes
Dashboard and gives a read/browse view of cluster resources (Deployments, Pods, Nodes, PersistentVolumes, RBAC objects,
etc.) without needing to run individual `kubectl` commands.

Headlamp is deployed via the `radar/headlamp` Helm chart (release `headlamp`, namespace `headlamp`) from
`helmfile.d/30-custom.yaml`. It is only exposed inside the cluster (ClusterIP Service) — there is no public Ingress —
so access is via `kubectl port-forward` plus a ServiceAccount token.

## Prerequisites

- `kubectl` is configured with a context pointing at the target cluster.
- Your kubeconfig user has permission to port-forward to the `headlamp` namespace and to request tokens for the
  `headlamp` ServiceAccount (`create` on `serviceaccounts/token`).
- The `headlamp` release is deployed on the cluster you're targeting (namespace `headlamp`).

## 1) Port-forward the Headlamp Service

```bash
kubectl port-forward -n headlamp service/headlamp 8080:80
```

Leave this running in a terminal. Headlamp is now reachable at:

```
http://localhost:8080/c/main
```

## 2) Get an access token

Headlamp authenticates using a Kubernetes ServiceAccount token for the `headlamp` ServiceAccount, which is bound to the
`view` ClusterRole plus a small `headlamp-cluster-read` ClusterRole for read-only access to cluster-scoped resources
(Nodes, PersistentVolumes, StorageClasses, metrics, RBAC objects). It deliberately has no access to Secrets.

**Ephemeral token (recommended, no extra objects created)**

```bash
kubectl create token headlamp -n headlamp --duration=24h
```

This uses the TokenRequest API to mint a short-lived token for the `headlamp` ServiceAccount. Nothing is persisted in
the cluster — re-run the command whenever the token expires.

## 3) Log in

1. Open `http://localhost:8080/c/main/token` in your browser.
2. Paste the token from step 2 into the token field.
3. Submit — you should land on the Headlamp cluster overview.

Tokens obtained via `kubectl create token` expire after the requested `--duration` and will need to be
regenerated; the Headlamp deployment itself sets a `-session-ttl=86400` (24h) session timeout regardless of token
lifetime.
