# PRD: k8s-podrestart-operator

**Project Type:** Personal / Interview Portfolio  
**Language:** Go  
**Framework:** kubebuilder  
**Status:** Not Started  

---

## Overview

A Kubernetes operator that watches a custom `ScheduledRestart` CRD and automatically restarts pods matching a label selector on a cron schedule or on-demand trigger. Built to demonstrate production-grade operator fundamentals: CRDs, reconcile loops, watches, requeueing, and finalizers.

This project is intentionally scoped small. The goal is a working, well-structured operator you can demo, explain, and extend — not a feature-complete product.

---

## Goals

- Build a working Kubernetes operator using kubebuilder from scratch
- Cover every operator fundamental that will come up in a technical interview
- Produce a clean GitHub repo that can be referenced during interviews or linked on the resume
- Understand the reconcile loop deeply enough to explain it without notes

## Non-Goals

- This is not a production-hardened tool
- No UI, no API server, no Helm chart (yet)
- No multi-tenancy or RBAC scoping in v1
- Not trying to compete with existing restart tools (Reloader, Stakater, etc.)

---

## The Custom Resource: `ScheduledRestart`

```yaml
apiVersion: restart.jgiornazi.dev/v1alpha1
kind: ScheduledRestart
metadata:
  name: payments-weekly-restart
  namespace: production
spec:
  selector:
    matchLabels:
      app: payments-svc
  schedule: "0 3 * * 0"   # cron: every Sunday at 3am
  suspend: false            # set true to pause without deleting
  restartPolicy: RollingDelete  # delete pods one at a time, let Deployment recreate
status:
  lastRestartTime: "2025-03-23T03:00:00Z"
  nextRestartTime: "2025-03-30T03:00:00Z"
  observedGeneration: 1
```

### Field Reference

| Field | Type | Description |
|---|---|---|
| `spec.selector` | LabelSelector | Pods to target. Matches standard Kubernetes label selectors. |
| `spec.schedule` | string | Cron expression defining restart cadence. |
| `spec.suspend` | bool | If true, operator skips reconciliation without deleting the resource. |
| `spec.restartPolicy` | enum | `RollingDelete` (one pod at a time) or `AllAtOnce` (delete all simultaneously). |
| `status.lastRestartTime` | string | ISO timestamp of last successful restart. |
| `status.nextRestartTime` | string | ISO timestamp of next scheduled restart. |
| `status.observedGeneration` | int | Tracks spec changes for idempotency. |

---

## Architecture

```
┌─────────────────────────────────────────────────────┐
│                   Kubernetes API                     │
│                                                     │
│   ScheduledRestart CRD     Pod resources            │
└────────────┬────────────────────┬───────────────────┘
             │ Watch              │ List / Delete
             ▼                   ▼
┌─────────────────────────────────────────────────────┐
│              ScheduledRestartReconciler              │
│                                                     │
│  1. Fetch ScheduledRestart resource                 │
│  2. Check suspend flag → skip if true               │
│  3. Parse cron schedule → is it time to restart?    │
│  4. If yes → list matching pods → delete them       │
│  5. Update status.lastRestartTime                   │
│  6. Requeue for next scheduled run                  │
└─────────────────────────────────────────────────────┘
```

The reconciler is the core of the operator. Every interesting decision lives there.

---

## Reconcile Loop — Detailed Flow

This is what you need to be able to explain in an interview cold.

```
Reconcile(ctx, req) called
│
├── r.Get(ctx, req.NamespacedName, &sr)
│   └── if NotFound → return nil (resource deleted, nothing to do)
│
├── Check sr.Spec.Suspend
│   └── if true → return Result{} (stop, don't requeue)
│
├── Parse sr.Spec.Schedule (cron)
│   └── Determine: has the scheduled time passed since lastRestartTime?
│
├── if NOT time to restart yet
│   └── return Result{RequeueAfter: timeUntilNextRun}
│
├── if time to restart
│   ├── List pods matching sr.Spec.Selector
│   ├── Delete pods (RollingDelete: one at a time with sleep, or AllAtOnce)
│   ├── Update sr.Status.LastRestartTime = now
│   ├── Update sr.Status.NextRestartTime = next cron tick
│   └── return Result{RequeueAfter: timeUntilNextRun}
│
└── On any error → return Result{}, err (controller-runtime will requeue with backoff)
```

---

## Implementation Plan

### Phase 1 — Scaffold & CRD (Day 1)

- [ ] Install kubebuilder and kind (local cluster)
- [ ] `kubebuilder init --domain jgiornazi.dev --repo github.com/jgiornazi/k8s-podrestart-operator`
- [ ] `kubebuilder create api --group restart --version v1alpha1 --kind ScheduledRestart`
- [ ] Define `ScheduledRestartSpec` and `ScheduledRestartStatus` structs in `api/v1alpha1/`
- [ ] Run `make generate && make manifests` to generate CRD YAML
- [ ] Apply CRD to local kind cluster and verify with `kubectl get crds`

### Phase 2 — Reconcile Loop (Day 1-2)

- [ ] Implement `Reconcile()` in `internal/controller/scheduledrestart_controller.go`
- [ ] Fetch the `ScheduledRestart` resource, handle NotFound
- [ ] Implement suspend check
- [ ] Integrate a cron parser (`github.com/robfig/cron/v3`) to compute next/last run times
- [ ] Implement pod listing via label selector
- [ ] Implement `RollingDelete` restart strategy (delete pods one at a time with sleep between each)
- [ ] Implement `AllAtOnce` restart strategy (delete all matching pods simultaneously)
- [ ] Add strategy branching logic in reconciler to route between `RollingDelete` and `AllAtOnce`
- [ ] Update status subresource after each restart
- [ ] Return correct `RequeueAfter` durations

### Phase 3 — Finalizer (Day 2)

- [ ] Add a finalizer string: `restart.jgiornazi.dev/finalizer`
- [ ] On reconcile: if resource is being deleted (DeletionTimestamp set), run cleanup logic, remove finalizer
- [ ] Register finalizer on new resources before doing any work

### Phase 4 — Tests & Polish (Day 3)

- [ ] Write table-driven unit tests for the cron schedule logic
- [ ] Write an integration test using envtest (kubebuilder scaffolds this)
- [ ] Add a `config/samples/` directory with example CRD manifests
- [ ] Write a clear README: what it does, how to install, how to create a ScheduledRestart, how to verify it worked

---

## Key Concepts to Understand as You Build

These are the exact things an interviewer will probe. Build them, then make sure you can explain them.

### The Operator Pattern
An operator extends Kubernetes with custom domain logic. It pairs a CRD (the "what you want") with a controller (the "how to get there"). Your ScheduledRestart CRD is the desired state. The reconciler is what makes reality match it.

### The Reconcile Loop
The reconcile loop is the heart of every operator. It is:
- **Triggered** by any change to a watched resource (create, update, delete)
- **Idempotent** — it should produce the same result no matter how many times it runs
- **Level-triggered, not edge-triggered** — it doesn't care *what* changed, only what the *current state* is vs desired state

This is conceptually identical to your canary monitoring service: watch → evaluate state → act → requeue.

### RequeueAfter
When you return `ctrl.Result{RequeueAfter: duration}`, you're telling controller-runtime "call me again in N time." This is how you implement scheduled behavior without a separate goroutine or ticker.

### Finalizers
A finalizer is a string you add to a resource's `metadata.finalizers` list. Kubernetes won't delete the resource until all finalizers are removed. This gives your operator a chance to run cleanup logic before the resource disappears.

### Status Subresource
The `status` section of your CRD is updated separately from `spec`. You use `r.Status().Update()` not `r.Update()` to write status. This keeps spec (desired state) and status (observed state) cleanly separated.

---

## Learning Resources — In Order

Work through these in sequence. Don't skip ahead.

### Step 1 — Understand the operator pattern conceptually (2-3 hours)
- **The Kubernetes Operator Pattern** — kubernetes.io/docs/concepts/extend-kubernetes/operator/
  Read this first. Short, official, foundational.
- **What is a Kubernetes Operator?** — YouTube, TechWorld with Nana
  Best visual intro. She explains CRDs and controllers with diagrams. Watch before writing any code.

### Step 2 — Learn the kubebuilder framework (4-6 hours, hands-on)
- **The Kubebuilder Book** — book.kubebuilder.io
  Read chapters 1–3 and follow along. This is the primary reference. Don't just read — actually run the commands.
- **Kubebuilder Quick Start** — book.kubebuilder.io/quick-start
  Do this before chapter 1. Gets you a working scaffold in 20 minutes.

### Step 3 — Understand controller-runtime (1-2 hours)
- **controller-runtime examples/** — github.com/kubernetes-sigs/controller-runtime/tree/main/examples
  Read the `burstablecronreconciler` and `crd` examples. This is the library kubebuilder uses under the hood.
- **Deep dive: How controller-runtime works** — github.com/kubernetes-sigs/controller-runtime/blob/main/pkg/reconcile/reconcile.go
  Read the Reconciler interface definition. It's 30 lines. Understand it cold.

### Step 4 — Go concurrency fundamentals (2-3 hours, if rusty)
- **Go by Example: Goroutines, Channels, Select** — gobyexample.com
  The reconcile loop runs in a goroutine. You need to be comfortable with context cancellation and channel patterns.
- **Go Concurrency Patterns** — YouTube, Rob Pike (Google I/O 2012)
  Classic talk. Still the best explanation of Go's concurrency model.

### Step 5 — Deeper operator patterns (read while building)
- **Programming Kubernetes** — O'Reilly book, Hausenblas & Schimanski
  Don't read cover to cover. Use chapters 4–6 as a reference while building. Covers informers, work queues, and error handling in depth.
- **Kubebuilder Slack** — kubernetes.slack.com #kubebuilder
  When you get stuck (you will), this is where to ask.

### Reference Docs (bookmark these)
- controller-runtime API docs — pkg.go.dev/sigs.k8s.io/controller-runtime
- client-go label selector docs — pkg.go.dev/k8s.io/apimachinery/pkg/labels
- robfig/cron — pkg.go.dev/github.com/robfig/cron/v3 (for parsing cron expressions)

---

## Success Criteria

You're done when:

- [ ] `kubectl apply -f config/samples/scheduled_restart.yaml` creates the resource
- [ ] Operator logs show it detecting the schedule and deleting pods
- [ ] `kubectl get scheduledrestart` shows correct `lastRestartTime` and `nextRestartTime` in status
- [ ] Setting `suspend: true` stops restarts without errors
- [ ] Deleting the ScheduledRestart resource completes cleanly (finalizer works)
- [ ] README explains how to run it against a local kind cluster in under 5 steps
- [ ] You can explain the reconcile loop out loud without looking at code
- [ ] Operator handles a pod list error gracefully and requeues with backoff without panicking
- [ ] Operator handles a status update failure without re-deleting pods (idempotency under partial failure)

---

*Built as an interview portfolio project. Conceptually based on the canary monitoring service built at Splunk (watch → probe → react → requeue).*
