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
1. `kind create cluster --config cluster/kind-setup.yaml`
2. Follow [guide](https://kind.sigs.k8s.io/docs/user/loadbalancer/) to setup loadbalancer.
   - If Ip range is different from in example you have to specify the yaml file with correct address range on your own and apply it instead of https://kind.sigs.k8s.io/examples/loadbalancer/metallb-config.yaml.
3. Generate docker images for manager and worker.
   ```bash
     cd src/services/worker/
     ./new_version.sh 0.0.1
     cd ../manager
     ./new_version.sh 0.0.1
     cd ../../..
     ```
4. ```bash
   cd cluster
   ./deploy.sh --local
   cd ..
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








