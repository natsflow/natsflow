A flow with similar functionality to [bitnami-labs/kubewatch](https://github.com/bitnami-labs/kubewatch).

Listens to events in all kube clusters and outputs them to slack.

<img src="https://github.com/natsflow/natsflow/blob/master/examples/slack-kubewatch/output.png" width="300">

1. Run [slack-nats](https://github.com/natsflow/slack-nats) & [kube-nats](https://github.com/natsflow/kube-nats)
1. Run slackkubewatch.js flow.

Inside kubernetes (expects nats at "nats-cluster" - see deployment.yaml):

```
skaffold dev
```

...or external to cluster:
```
kubectl port-forward nats-cluster-1 4222:4222
npm install
node slackkubewatch.js
```
