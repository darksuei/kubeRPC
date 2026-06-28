import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

// In Kubernetes, KUBERPC_CORE_URL is injected automatically via the admission webhook.
// Outside Kubernetes, set it in .env or pass it to the constructor.
const rpc = new KubeRPC();

// TARGET_SERVICE must match the kuberpc.suei.io/service annotation on the server pod
// and the name of the Kubernetes Service fronting it.
const targetService = process.env.TARGET_SERVICE ?? "sample-server";
const server = rpc.service(targetService);

const result = await server.call("fib", { n: 10 });

console.log(`fib(10) = ${result}`);

rpc.close();
