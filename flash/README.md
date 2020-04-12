#### Inject & Delete Linkerd

```
kubectl get -n kn deploy -o yaml \
  | linkerd inject - \
  | kubectl apply -f -
```

```
linkerd install --ignore-cluster | kubectl delete -f -
```

#### Docker Credentials

Call `docker login` locally and then execute the following command:

```
kubectl create secret generic regcred \
--from-file=.dockerconfigjson=/Users/kananrahimov/.docker/config.json \
--type=kubernetes.io/dockerconfigjson

kubectl create secret docker-registry dhregcred --docker-server=https://index.docker.io/v1/ --docker-username=kenanbek --docker-password=<PASSWORD> --docker-email=mail@kenanbek.me
```

#### Pod shell

```
kubectl -n kn exec -it flash-cddf44465-w4tnj -- /bin/sh
```

#### Pod internal IPs:

```
k -n kn describe pods | grep IP
```
