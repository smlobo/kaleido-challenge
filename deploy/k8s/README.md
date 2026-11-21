# Deploying to microk8s

The compose stack, in a `ticketing` namespace: Postgres, Besu, Jaeger, the
ticketing service with two replicas behind MetalLB, and the demo frontend.

| File | What |
|---|---|
| `00-namespace.yaml` | The `ticketing` namespace |
| `10-postgres.yaml` | Postgres StatefulSet, 1 replica, ephemeral (`emptyDir`); Service `postgres:5432` |
| `20-besu.yaml` | Besu StatefulSet (the single QBFT validator), ephemeral (`emptyDir`); genesis ConfigMap, validator key Secret; Service `besu:8545`; JSON-RPC on LoadBalancer **10.10.1.47:8545** |
| `30-jaeger.yaml` | Jaeger all-in-one; OTLP Service `jaeger:4317/4318`; UI on LoadBalancer **http://10.10.1.46** |
| `40-ticketing.yaml` | The service: Deployment with 2 replicas spread over nodes, a PodDisruptionBudget, and LoadBalancer **http://10.10.1.45** |
| `50-frontend.yaml` | The demo frontend: nginx serving the built page and passing `/api` to the service; LoadBalancer **http://10.10.1.48** |

The LoadBalancer IPs are fixed with MetalLB's `metallb.universe.tf/loadBalancerIPs`
annotation, from the free end of the pool (10.10.1.40-49).

For reviewers outside the LAN, a reverse proxy with TLS (10.10.1.60) forwards
`kaleido-frontend.lobo.codes`, `kaleido-backend.lobo.codes` and
`kaleido-jaeger.lobo.codes` to those three IPs. The service's `PUBLIC_BASE_URL`
is the backend hostname, so ticket URIs on new events resolve from anywhere.

**The chain and the database are ephemeral, like compose.** Their data is on
`emptyDir` volumes: it survives a container restart (a crash, a failed probe)
but not the pod going away (deleted, evicted, its node rebooted, or a changed
manifest re-applied). Requests equal limits, so under memory pressure other
pods are evicted first. If either is wiped, reset both (below). For data that
must outlive the pods, swap the `emptyDir` volumes for PersistentVolumeClaims
on `microk8s-hostpath`.

## 1. Build and push the image

The nodes are amd64; the Dockerfile cross-compiles, so building on an arm64
Mac needs no emulation. The microk8s registry listens on port 32000 of every
node, over plain HTTP, so Docker must be told to allow it once: Docker
Desktop, Settings, Docker Engine, add

```json
"insecure-registries": ["10.10.1.50:32000"]
```

then Apply & restart. From the repository root:

```bash
docker buildx build --platform linux/amd64 -t 10.10.1.50:32000/kaleido-ticketing:latest --push .
```

The pods pull it as `localhost:32000/kaleido-ticketing:latest`: the same
registry, reached through each node's own port. If the registry gives trouble,
copy the image to a node instead (and to each node a pod may land on):

```bash
docker buildx build --platform linux/amd64 -t localhost:32000/kaleido-ticketing:latest --load . && docker save localhost:32000/kaleido-ticketing:latest | ssh 10.10.1.50 microk8s ctr image import -
```

(With an imported image, set `imagePullPolicy: IfNotPresent` in `40-ticketing.yaml`.)

The frontend is a second image, built from `frontend/` (static files behind
nginx, so it is the same on every platform):

```bash
docker buildx build --platform linux/amd64 -t 10.10.1.50:32000/kaleido-ticketing-frontend:latest --push frontend
```

## 2. Create the namespace and the secret

`KEY_ENCRYPTION_KEY` and `ADMIN_TOKEN` go straight into the cluster from the
shell, never into a file. Keep the same `KEY_ENCRYPTION_KEY` for as long as
the database lives: keys encrypted under one value can't be decrypted with
another. (Postgres's password is a plain value in the manifests, as in
compose: the database is only reachable inside the cluster, and the keys in it
are encrypted under `KEY_ENCRYPTION_KEY`, which it never sees.)

```bash
kubectl apply -f deploy/k8s/00-namespace.yaml
```

```bash
kubectl -n ticketing create secret generic ticketing-secret --from-literal=key-encryption-key="$(openssl rand -base64 32)" --from-literal=admin-token="$(openssl rand -hex 16)"
```

Read the admin token back when you need it (for `create_organizer.py`):

```bash
kubectl -n ticketing get secret ticketing-secret -o jsonpath='{.data.admin-token}' | base64 -d; echo
```

A Kubernetes Secret is stored in etcd base64-encoded, not encrypted, unless
the cluster enables encryption at rest. That is weaker than "environment
only": in production the key would come from a KMS.

## 3. Apply

```bash
kubectl apply -f deploy/k8s/
```

```bash
kubectl -n ticketing get pods,svc -o wide
```

The ticketing pods wait (init container) for Postgres and Besu, then migrate
(one at a time, under an advisory lock) and become ready. Check:

```bash
curl -s http://10.10.1.45/readyz
```

## 4. Use it

```bash
./scripts/create_organizer.py "Rock Fest Ltd" --url http://10.10.1.45 --token <admin-token>
```

```bash
python3 scripts/find_activity.py --rpc http://10.10.1.47:8545
```

The frontend: **http://10.10.1.48**. Its nginx passes `/api/...` to the
`ticketing` Service, so the browser talks to one origin and there is no CORS.
One difference from `npm run dev` on localhost: browsers allow the camera only
on HTTPS or localhost, so at the gate, check in with a picture of the pass (a
screenshot) or by pasting it. Copy buttons fall back to a prompt for the same
reason. For the camera, run the frontend locally against the cluster instead:

```bash
cd frontend && TICKETING_URL=http://10.10.1.45 npm run dev
```

Traces: http://10.10.1.46, service `ticketing`. Each pod is its own `service.instance.id`, and `host.name` is the
pod's name.

## Updating the service

Push a new image (step 1), then:

```bash
kubectl -n ticketing rollout restart deployment/ticketing
```

For a new frontend image, `rollout restart deployment/frontend`; reload the
page (it is served with `no-cache`, its assets by content hash).

The rollout replaces one pod at a time (`maxUnavailable: 0`), so the API stays up.

## If the chain or the database was reset

The chain and the database must be reset together: a new chain under a
database that remembers the old one leaves deploys stuck on nonces the new
chain has never seen. Signs: `kubectl -n ticketing get pods` shows `besu-0` or
`postgres-0` younger than the other, or new events stay `deploying`.

Delete both pods (they come back empty), then restart the service, which
migrates and creates its deployer key only at startup:

```bash
kubectl -n ticketing delete pod postgres-0 besu-0 && kubectl -n ticketing rollout restart deployment/ticketing
```

Then recreate the organizer, users and event for the demo, and press Clear in
the frontend: its stored rows belong to the old database.

## Removing everything

Deletes the namespace with all its pods, services and secrets:

```bash
kubectl delete namespace ticketing
```
