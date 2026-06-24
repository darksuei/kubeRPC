# kubeRPC - Node SDK

Node.js SDK for [kubeRPC](https://github.com/darksuei/kubeRPC). Enables services to register callable methods and invoke methods on other services over raw TCP.

---

## Installation

```bash
npm install kuberpc-node
```

---

## Setup

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC({
  coreURL: "http://kuberpc-core:8080", // kubeRPC core endpoint
  serviceName: "my-service", // unique name for this service
  port: 8082, // port this service listens on for incoming RPC calls
  host: "my-service", // hostname other services use to reach this one (K8s service name or IP)
});
```

> `host` should be the address other services in the cluster can connect to - typically the Kubernetes service name. Defaults to `0.0.0.0` if omitted.

---

## Register a method

Expose a method that other services can call. The TCP listener starts on first `register()` call.

```js
await rpc.register("getUser", async ({ id }) => {
  return db.users.findById(id);
});
```

---

## Call a method

Invoke a method on another service. The first call resolves the endpoint from kubeRPC core and opens a persistent TCP connection. Subsequent calls reuse it.

```js
const user = await rpc.call("user-service", "getUser", { id: "123" });
```

---

## Full sample

**service-a** (registers a method):

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC({
  coreURL: "http://kuberpc-core:8080",
  serviceName: "service-a",
  port: 8082,
  host: "service-a",
});

await rpc.register("greet", async ({ name }) => {
  return `Hello, ${name}!`;
});
```

**service-b** (calls the method):

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC({
  coreURL: "http://kuberpc-core:8080",
  serviceName: "service-b",
  port: 8083,
  host: "service-b",
});

const message = await rpc.call("service-a", "greet", { name: "world" });
console.log(message); // Hello, world!
```
