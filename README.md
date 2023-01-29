# distributed-shortest-path
Project calculating shortest path in a distributed graph.

## Local development
### Dependencies:
- Go
- kind
  - other local version of kubernetes such as `minikube` may be used, but then you're on your own with setting up the cluster and uploading docker images to it.
- kubectl
- docker

### Initialization
1. Setup cluster: `kind create cluster --config cluster/local-kind-cluster.yaml`
2. Follow [guide](https://kind.sigs.k8s.io/docs/user/loadbalancer/) to setup loadbalancer.
   - If Ip range is different from in example you have to specify the yaml file with correct address range on your own and apply it instead of https://kind.sigs.k8s.io/examples/loadbalancer/metallb-config.yaml.
3. Generate docker images for manager and worker.
   ```bash
     ./src/services/worker/new_version.sh 0.0.1
     ./src/services/manager/new_version.sh 0.0.1
     ```
4. ```bash
   ./cluster/deploy.sh --local
   ```

### Updates
#### Manager
```bash
./update_manager.sh <VERSION>
```
#### Workers
```bash
./update_workers.sh <VERSION>
```

#### DB Schema
1. Update models in `src/libs/db/model`.
2. If new models are added adjust the list in `src/libs/db/models/list.go`.
3. ```bash
    (cd src && go generate)
   ```

#### GRPC proto
1. Update grpc proto in `src/services/worker/link/proto/link.proto`.
2. ```bash
    (cd src && go generate)
    ```

#### Manual connecting to Postgres:
```bash
kubectl run postgresql-dev-client --rm --tty -i --restart='Never' --namespace postgres --image docker.io/bitnami/postgresql:14.1.0-debian-10-r80 --env="PGPASSWORD=psltest" \
--command -- psql --host postgres -U admin -d postgresdb -p 5432
```





