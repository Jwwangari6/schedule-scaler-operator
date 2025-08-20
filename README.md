# Schedule Scaler Operator

A Kubernetes operator built with [Kubebuilder](https://book.kubebuilder.io/) that automates **cost savings** by scaling workloads (Deployments, StatefulSets) up or down based on custom schedules.  

Use cases:
- Scale **non-production environments** down after working hours.
- Automatically restore replicas in the morning.
- Reduce cloud costs by shutting down idle workloads.

---

## 🚀 Features
- Define schedules per namespace via a custom resource (`ScheduleScaler`).
- Scale Deployments (StatefulSets support coming soon).
- Simple cron-like configuration (`scaleUp`/`scaleDown` times + days).
- Status reporting in the CRD.
- Built using Kubebuilder and the controller-runtime library.

---

## 🏗️ Install

### Prerequisites
- Kubernetes cluster (minikube, kind, or real cluster).
- `kubectl` configured for your cluster.
- Go 1.22+ (for local dev).
- [Kubebuilder](https://book.kubebuilder.io/quick-start.html) installed.

### Steps

1. Install CRDs:
   ```bash
   make install
