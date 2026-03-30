# k8s-podrestart-operator

A Kubernetes operator that watches a custom `ScheduledRestart` CRD and automatically restarts pods matching a label selector on a cron schedule. Built with kubebuilder to demonstrate production-grade operator fundamentals: CRDs, reconcile loops, requeueing, and finalizers.

## How it works

The operator watches for `ScheduledRestart` resources. On each reconcile it checks whether the cron schedule has elapsed since the last restart. If yes, it deletes the matching pods (letting the Deployment's ReplicaSet recreate them) and updates the status with the last and next restart times.

## Prerequisites

- Go v1.21+
- kubectl
- [kind](https://kind.sigs.k8s.io/) (local Kubernetes cluster)
- [kubebuilder](https://book.kubebuilder.io/)

## Running locally against kind

**1. Create a local cluster**
```sh
kind create cluster --name podrestart
```

**2. Install the CRD**
```sh
make manifests && kubectl apply -f config/crd/bases/
```

**3. Run the operator**
```sh
make run
```

**4. Create a ScheduledRestart resource**
```sh
kubectl apply -f config/samples/restart_v1alpha1_scheduledrestart.yaml
```

**5. Verify it's working**
```sh
# Watch the operator logs in the terminal running make run
# Check status
kubectl get scheduledrestart scheduledrestart-sample -o yaml
```

## Example resource

```yaml
apiVersion: restart.jgiornazi.dev/v1alpha1
kind: ScheduledRestart
metadata:
  name: payments-weekly-restart
spec:
  selector:
    matchLabels:
      app: payments-svc
  schedule: "0 3 * * 0"   # every Sunday at 3am
  suspend: false
  restartPolicy: RollingDelete
```

## Fields

| Field | Description |
|---|---|
| `spec.selector` | Label selector targeting pods to restart |
| `spec.schedule` | Cron expression defining restart cadence |
| `spec.suspend` | Set true to pause restarts without deleting the resource |
| `spec.restartPolicy` | `RollingDelete` (one at a time) or `AllAtOnce` |
| `status.lastRestartTime` | Timestamp of last successful restart |
| `status.nextRestartTime` | Timestamp of next scheduled restart |

## Uninstall

```sh
kubectl delete -f config/samples/restart_v1alpha1_scheduledrestart.yaml
kubectl delete -f config/crd/bases/
kind delete cluster --name podrestart
```

## License

Apache 2.0
