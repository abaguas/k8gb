## Start the failover playground
```bash
K8GB_LOCAL_VERSION=test FULL_LOCAL_SETUP_WITH_APPS=true make deploy-full-local-setup
```

## Start the backend server:
```bash
node server.js
```

## Start the frontend:
```bash
live-server
```
