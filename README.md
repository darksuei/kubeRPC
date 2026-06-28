## KubeRPC

**KubeRPC** is a **kubernetes-native remote procedure call (RPC) framework** designed to enable seamless and low-latency communication between microservices deployed within the same kubernetes cluster.

<p align="center">
  <img src="./assets/rpc.png" alt="RPC Overview" width="700" />
</p>

### **Why does it matter?**

Microservice communication is typically implemented over HTTP-based APIs (REST, GraphQL, gRPC). While these are well-established, they introduce non-negligible overhead compared to in-process calls, especially in low latency environments.

In monolithic systems, function calls are in-process and incur no network serialization, routing, or gateway overhead. In distributed systems, even internal calls must traverse these layers.

KubeRPC is designed for **internal, cluster-local service communication** where:

- Services are already co-deployed in kubernetes
- Trust boundaries are internal (not public)
- Latency is a critical constraint

This does not replace external APIs or public-facing HTTP interfaces. It is intended as a complementary mechanism for **high-frequency internal RPC-style communication with low latency requirements**.

#### **Why not just open raw TCP between services?**

Raw TCP between services is unaware of kubernetes. You would need to manually wire hostnames and ports into every service, manage that configuration across environments, and keep it in sync as services are added, moved, or scaled. kubeRPC delegates all of that to the cluster itself. The admission webhook watches your pod annotations and injects the correct configuration at scheduling time - services discover each other through kubernetes DNS and the registry without any developer-managed config. Adding a new service to the mesh is an annotation, not a configuration change across every other service that needs to reach it.

---

#### **Benchmark**

A simple benchmark was run using 10 sequential calls to `fib(40)` across services.

#### **The Result?**

KubeRPC showed approximately **~60% lower average latency** compared to equivalent HTTP-based service calls in the same cluster environment.

[Benchmark source code](https://github.com/darksuei/kubeRPC/tree/main/sdks/node/benchmark)

---

### **How it works**

1. KubeRPC deploys a **core service** within your kubernetes cluster that acts as the central service registry.
2. A **mutating admission webhook** bundled with the core watches for pods annotated with `kuberpc.suei.io/enabled: "true"` and automatically injects the runtime configuration (`KUBERPC_CORE_URL`, `KUBERPC_SERVICE_NAME`, `KUBERPC_HOST`, `KUBERPC_PORT`) at pod creation time.
3. Services **register callable methods** with the KubeRPC core on startup. No manual configuration is required inside the cluster.
4. Other services resolve and invoke those methods using the KubeRPC SDK. All RPC traffic flows directly between services over persistent TCP connections. The core is only consulted for endpoint resolution.

---

### **Setup and Deployment**

#### **Requirements**

- A kubernetes cluster (any version compatible with Helm).

#### **Deploying kubeRPC**

KubeRPC can be deployed using a Helm chart. The admission webhook and TLS certificates are included and configured automatically.

```bash
helm upgrade --install kuberpc-core \
  oci://ghcr.io/darksuei/charts/kuberpc-core \
  --version 2.0.0 \
  -n kuberpc \
  --create-namespace \
  --wait
```

---

### **Usage**

#### **Opt your pods in**

Annotate any pod you want kubeRPC to configure:

```yaml
annotations:
  kuberpc.suei.io/enabled: "true"        # required, triggers env injection
  kuberpc.suei.io/service: "my-service"  # required for servers, sets service name and host
  kuberpc.suei.io/port: "7749"           # optional, defaults to 7749
```

> The kubernetes `Service` fronting your pod must be named to match `kuberpc.suei.io/service` so that peer services can reach it via cluster DNS.

#### **Registering Methods**

Once annotated, the SDK needs no configuration inside the cluster:

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC();

await rpc.register("getUser", async ({ id }) => {
  return db.users.findById(id);
});
```

#### **Calling Methods**

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC();

const userService = rpc.service("user-service");
const user = await userService.call("getUser", { id: "123" });
```

---

## **Available SDKs**

Currently, a **TypeScript SDK** is available:

- Node.js SDK: [Source](https://github.com/darksuei/kubeRPC/tree/main/sdks/node) | [NPM](https://www.npmjs.com/package/kuberpc-node)

We welcome contributions for SDKs in other programming languages!

---

## **Future Work**

| # | Item | Status |
|---|---|---|
| 1 | **Request multiplexing** - concurrent calls to the same service are currently serialised through a promise queue. | Not started |
| 2 | **Retries, timeouts and error handling** - failed calls propagate immediately to the caller with no retry, and calls to an unresponsive service hang indefinitely. | Partially done - typed errors exist; timeouts and retries are not implemented |
| 3 | **Additional SDKs** - kubeRPC is language-agnostic at the protocol level but we currently only ship a Node.js SDK. | Not started |
| 4 | **Service TTL and heartbeat** - if a pod crashes without cleanly deregistering, its entry stays in the registry indefinitely and callers resolve a dead endpoint. | Not started |
| 5 | **RPC-level observability** - we have Prometheus metrics for core HTTP traffic and method registration counts, but no per-method call latency histograms, success/failure rates, or active connection gauges. | Not started |
| 6 | **Namespace-scoped registry** - registry keys are currently `service:<name>` with no namespace component, meaning two services with the same name in different kubernetes namespaces will collide. | Not started |
| 7 | **Security** - currently any pod in the cluster can register or deregister any service, and RPC traffic between services travels over plain TCP. | Not started |

---

## **Contributing**

Contribution is very much welcome by:

1. Building SDKs for other languages.
2. Reporting bugs or requesting features via GitHub issues.
3. Submitting pull requests to improve the core or SDKs.

## **License**

![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)
