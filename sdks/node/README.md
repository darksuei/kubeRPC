# kubeRPC - Node.js SDK

Node.js SDK for [kubeRPC](https://github.com/darksuei/kubeRPC). Register callable methods on your service and invoke methods on other services over a persistent TCP connection with MessagePack framing.

---

## Installation

```bash
npm install kuberpc-node
```

---

## Kubernetes usage (zero-config)

When kubeRPC core is deployed with the admission webhook enabled, annotate your pod and all required environment variables are injected automatically at pod creation time:

```yaml
annotations:
  kuberpc.suei.io/enabled: "true"       # inject KUBERPC_CORE_URL
  kuberpc.suei.io/service: "my-service" # inject KUBERPC_SERVICE_NAME, KUBERPC_HOST, KUBERPC_PORT
  kuberpc.suei.io/port: "7749"          # optional - defaults to 7749
```

With those env vars present, the constructor takes no arguments:

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC();
```

> The Kubernetes Service fronting your pod **must be named to match** `kuberpc.suei.io/service`. That name is what kubeRPC core stores as the reachable host for your service.

---

## Non-Kubernetes usage

Pass configuration explicitly. These values override any environment variables.

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC({
  coreURL: "http://localhost:8080", // kubeRPC core endpoint
  serviceName: "my-service",        // unique name for this service
  host: "localhost",                // address other services use to reach this one
  port: 7749,                       // TCP port this service listens on for inbound RPC calls
});
```

---

## Register a method

Expose a callable method. The TCP listener starts on the first `register()` call, and the service is registered with kubeRPC core.

```js
await rpc.register("getUser", async ({ id }) => {
  return db.users.findById(id);
});
```

---

## Call a method

Get a proxy for a target service, then call a method on it. The first call resolves the endpoint from kubeRPC core and opens a persistent TCP connection. Subsequent calls reuse the connection.

```js
const userService = rpc.service("user-service");

const user = await userService.call("getUser", { id: "123" });
```

---

## Full example

**service-a** - registers a method:

```js
import { KubeRPC } from "kuberpc-node";

// In Kubernetes: no args needed, env vars are injected by the webhook.
// Outside Kubernetes: pass config explicitly.
const rpc = new KubeRPC({
  coreURL: "http://localhost:8080",
  serviceName: "service-a",
  host: "localhost",
  port: 7749,
});

await rpc.register("greet", async ({ name }) => {
  return `Hello, ${name}!`;
});
```

**service-b** - calls the method:

```js
import { KubeRPC } from "kuberpc-node";

const rpc = new KubeRPC({
  coreURL: "http://localhost:8080",
  serviceName: "service-b",
  host: "localhost",
  port: 7750,
});

const serviceA = rpc.service("service-a");
const message = await serviceA.call("greet", { name: "world" });

console.log(message); // Hello, world!

rpc.close();
```

---

## API reference

### `new KubeRPC(config?)`

All fields are optional. Environment variables are the fallback when a field is omitted.

| Field | Env var | Default | Description |
|---|---|---|---|
| `coreURL` | `KUBERPC_CORE_URL` | - | kubeRPC core base URL. **Required** (via config or env). |
| `serviceName` | `KUBERPC_SERVICE_NAME` | `""` | Name this service registers under. Required when calling `register()`. |
| `host` | `KUBERPC_HOST` | `"localhost"` | Address other services use to reach this one. |
| `port` | `KUBERPC_PORT` | `7749` | TCP port for inbound RPC calls. |

---

### `rpc.register(name, handler)`

Registers a method and starts the TCP listener if not already running. Returns a `Promise<void>` that resolves once the method is registered with kubeRPC core.

```ts
await rpc.register("methodName", async (args) => {
  return result;
});
```

---

### `rpc.service(name): ServiceProxy`

Returns a lightweight proxy for a named service. Does not make any network calls.

```ts
const svc = rpc.service("other-service");
```

---

### `ServiceProxy.call(method, args?): Promise<any>`

Invokes a method on the target service. Resolves the endpoint from kubeRPC core on the first call (cached for subsequent calls), then sends the request over a persistent TCP connection.

```ts
const result = await svc.call("methodName", { key: "value" });
```

---

### `rpc.close()`

Closes all outbound TCP connections, destroys the inbound TCP listener, and clears internal state. Call this on graceful shutdown.

```ts
rpc.close();
```

---

## Transport behaviour

- **Wire format**: length-prefixed MessagePack frames. Each frame is a 4-byte big-endian length followed by a MessagePack-encoded positional array.
- **Connection pooling**: one persistent TCP socket per target service, opened on first call and reused for all subsequent calls.
- **Endpoint cache**: service host/port is resolved from kubeRPC core once and cached. The cache is invalidated if the connection drops.
- **Concurrency**: concurrent calls to the same service are serialised through an internal queue. If you need parallel throughput across services, each target service gets its own independent queue and connection.
- **Reconnect**: if the socket drops, the next call attempts a single reconnect. There are no automatic retries beyond that - error handling is left to the caller.
