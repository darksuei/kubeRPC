## KubeRPC

**KubeRPC** is a **Kubernetes-native remote procedure call (RPC) framework** designed to enable seamless and low-latency communication between microservices deployed within the same Kubernetes cluster.

<p align="center">
  <img src="./assets/rpc.png" alt="RPC Overview" width="700" />
</p>

### **Why does it matter?**

Microservice communication is typically implemented over HTTP-based APIs (REST, GraphQL, gRPC). While these are well-established, they introduce non-negligible overhead compared to in-process calls, especially in low latency environments.

In monolithic systems, function calls are in-process and incur no network serialization, routing, or gateway overhead. In distributed systems, even internal calls must traverse these layers.

KubeRPC is designed for **internal, cluster-local service communication** where:

- Services are already co-deployed in Kubernetes
- Trust boundaries are internal (not public APIs)
- Latency and call overhead are critical constraints

It does not replace external APIs or public-facing HTTP interfaces. It is intended as a complementary mechanism for **high-frequency internal RPC-style interactions with low latency requirements**.

---

#### **Benchmark**

A simple benchmark was run using 10 sequential calls to `fib(40)` across services.

#### **The Result?**

KubeRPC showed approximately **~60% lower average latency** compared to equivalent HTTP-based service calls in the same cluster environment.

[Benchmark source code](https://github.com/darksuei/kubeRPC/tree/main/sdks/node/benchmark)

---

### **How it works**

1. KubeRPC deploys a **core service** within your Kubernetes cluster that acts as the central orchestrator.
2. Services can **register public functions** with the KubeRPC core service.
3. Other services can then invoke these functions using the KubeRPC SDK, eliminating the need for API implementations.

---

### **Setup and Deployment**

#### **Requirements**

- A Kubernetes cluster (any version compatible with Helm).

#### **Deploying kubeRPC**

KubeRPC can be deployed using a helm chart.

```bash
helm upgrade --install kuberpc-core \
  oci://ghcr.io/darksuei/charts/kuberpc-core \
  --version 1.0.0 \
  -n KubeRPC \
  --create-namespace \
  --wait
```

---

### **Usage**

#### **Registering Methods**

To register methods, services must interact with the KubeRPC core using the KubeRPC SDK.

#### **Calling Methods**

Once methods are registered, other services can directly invoke these methods using the SDK.

---

## **Available SDKs**

Currently, a **TypeScript SDK** is available:

- Node.js SDK: [Source](https://github.com/darksuei/kubeRPC/tree/main/sdks/node) | [NPM](https://www.npmjs.com/package/kuberpc-sdk)

We welcome contributions for SDKs in other programming languages!

---

## **Contributing**

Contribution is very much welcome by:

1. Building SDKs for other languages.
2. Reporting bugs or requesting features via GitHub issues.
3. Submitting pull requests to improve the core or SDKs.

## **License**

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
