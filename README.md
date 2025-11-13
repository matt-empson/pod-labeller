# Pod Labeller (Learning Project)

This is a small, simple Kubernetes controller built to learn more about Go, `client-go`, and how Kubernetes controllers actually work under the hood. It’s **not meant to be a real production controller** — just a hands-on way to play with informers, listers, workqueues, and reconcile loops.

All this controller does is watch Pods and automatically add a set of labels when they’re created or updated. The goal wasn’t to build something useful, but to understand the mechanics behind how controllers function at a low level (without Kubebuilder or controller-runtime hiding everything).

If you want to try it out on Minikube:

```sh
eval $(minikube docker-env)
docker build -t pod-labeller:latest .
kubectl apply -f configs/controller/
```

You can find the Kubernetes manifests under:

```configs/controller/```

This project is just for learning, experimenting, and getting comfortable with the controller patterns you see inside real Kubernetes projects.
