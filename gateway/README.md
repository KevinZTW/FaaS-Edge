# Function as a Service (FaaS) Monitoring

## Intro
In order to monitor pods for fine-grained management including scheduling, state management... etc.,
we use a simple mechanism to monitor function pods status which spin up by OpenFaaS or other mechanism.

## Prerequisite

## Install:

Below installation leverage multipass to create virtual nodes for k3s to create a kubernetes cluster.
Multipass works on Mac, Linux and Windows, but we only test on Mac.
If you already have a kubernetes cluster, you can skip this section.

### Install multipass and create 2 virtual nodes
```sh
 brew install --cask multipass
 multipass launch --name k3s-1 --mem 3G --disk 10G --cpus 1
 multipass launch --name k3s-2 --mem 3G --disk 10G --cpus 1

 multipass mount $(pwd) k3s-1:/root/ # mount project to virtual node
 multipass shell k3s-1
```

### Install k3s on virtual nodes
```sh
multipass shell k3s-1
# install the k3s
curl -sfL https://get.k3s.io | sh - 

# duplicate k3s.yaml to kube/config
cd ~/
mkdir .kube 
sudo cp /etc/rancher/k3s/k3s.yaml .kube/config
sudo chown $(id -u):$(id -g) .kube/config  


sudo vim /etc/systemd/system/k3s.service
```

Add below 2 lines to the end of the `k3s.service`
```
	--write-kubeconfig=/home/ubuntu/.kube/config \
    --write-kubeconfig-mode=644 \
```

Now the file should look like:
```
ExecStart=/usr/local/bin/k3s \
	server \
	--write-kubeconfig=/home/ubuntu/.kube/config \
	--write-kubeconfig-mode=644 \
```

restart the k3s and get the token, the token is for connecting other nodes to the k3s
```sh
sudo systemctl daemon-reload
sudo systemctl start k3s
sudo cat /var/lib/rancher/k3s/server/node-token
```

Enter virtual node 2 (k3s-2 in the example) and connect it to the cluster we just created
```sh
# get the node 1's IP
multipass info k3s-1 | grep IPv4 | awk '{ print $2 }' 
multipass shell k3s-2
# in node 2
k3s-2 $ curl -sfL https://get.k3s.io | K3S_URL= "<k3s-1 IP Address>"  K3S_TOKEN="<K3S token we previously got>" sh -
```
Congrats, now in node 1 (k3s-1) we should be able to see two nodes through `kubectl`
```sh
k3s-1 $ kubeclt get nodes
```

## Extend the K8S scheduler to get the pod status
In order to let us aware of the creation of function pods, we need to extend the K8S scheduler to get the pod status.
We let K8S notify us when a new pod need to be scheduled to avoid constantly polling the pod status.

### Create a KubeSchedulerConfiguration file
Create a file as below, in example the path is `/var/lib/scheduler/scheduler-config.yaml`
```sh
apiVersion: kubescheduler.config.k8s.io/v1
kind: KubeSchedulerConfiguration
clientConnection:
        kubeconfig: "/etc/rancher/k3s/k3s.yaml" # the kubeconfig file created by k3s

extenders:
        - urlPrefix: "http://localhost:8888/" # the port this application would listen to
          filterVerb: "filter"
          enableHTTPS: false
          tlsConfig:
                  insecure: true
          ignorable: true
```

The `KubeSchedulerConfiguration` file should be passed to the K8S `kube-scheduler`
with --config as mentioned [here](https://kubernetes.io/docs/reference/command-line-tools-reference/kube-scheduler/)

In K3S, the server support `kube-scheduler-arg` to pass the configuration file to `kube-scheduler`

```sh
ExecStart=/usr/local/bin/k3s \
    server \
    --kube-scheduler-arg=config=/var/lib/scheduler/scheduler-config.yaml \
    --write-kubeconfig=/home/ubuntu/.kube/config \
    --write-kubeconfig-mode=644 
```
No run this application and restart the K3S, the application should output the message when a new pod is created.
```sh
k3s-1 $ git clone https://github.com/KevinZTW/FaaS-Edge.git
k3s-1 $ cd FaaS-Edge/gateway
# install golang
k3s-1 $ wget https://go.dev/dl/go1.21.5.linux-arm64.tar.gz
k3s-1 $ sudo rm -rf /usr/local/go
k3s-1 $ sudo tar -C /usr/local -xzf go1.21.5.linux-arm64.tar.gz
 
k3s-1 $ export PATH=$PATH:/usr/local/go/bin
k3s-1 $ echo "export PATH=$PATH:/usr/local/go/bin" >> ~/.bashrc
k3s-1 $ go build  && ./gateway # build and run the application

```

### Install OpenFaaS & Verify the Kubernetes scheduler extension 
- Install OpenFaaS following the [official guide](https://docs.openfaas.com/deployment/kubernetes/#install-openfaas-with-arkade)
- Forward the gateway port to host `kubectl port-forward -n openfaas svc/gateway 8080:8080 --address='0.0.0.0' &`  


```sh
k3s-1 $ PASSWORD=$(kubectl get secret -n openfaas basic-auth -o jsonpath="{.data.basic-auth-password}" | base64 --decode; echo)
echo -n $PASSWORD | faas-cli login --username admin --password-stdin # login the OpenFaaS

k3s-1 $ faas-cli store list
k3s-1 $ faas-cli store deploy 
```









