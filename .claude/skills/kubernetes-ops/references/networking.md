# Kubernetes Networking & Security

> Load when: Services, Ingress, NetworkPolicies, RBAC.

## Service Types

```yaml
# ClusterIP (internal only — default)
apiVersion: v1
kind: Service
metadata:
  name: api
spec:
  selector:
    app: api
  ports:
  - port: 80
    targetPort: 3000

# NodePort (for debugging, not production)
# LoadBalancer (cloud provider provisions external LB)
```

## Ingress with TLS

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: api-ingress
  annotations:
    cert-manager.io/cluster-issuer: letsencrypt-prod
spec:
  tls:
  - hosts:
    - api.example.com
    secretName: api-tls
  rules:
  - host: api.example.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: api
            port:
              number: 80
```

## RBAC

```yaml
# ServiceAccount for the app
apiVersion: v1
kind: ServiceAccount
metadata:
  name: api-sa

# Role with minimal permissions
apiVersion: rbac.authorization.k8s.io/v1
kind: Role
metadata:
  name: api-role
rules:
- apiGroups: [""]
  resources: ["configmaps"]
  verbs: ["get", "list"]

# Bind role to service account
apiVersion: rbac.authorization.k8s.io/v1
kind: RoleBinding
metadata:
  name: api-rolebinding
subjects:
- kind: ServiceAccount
  name: api-sa
roleRef:
  kind: Role
  name: api-role
  apiGroup: rbac.authorization.k8s.io
```