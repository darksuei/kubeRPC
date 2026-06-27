import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

const rpc = new KubeRPC({
  coreURL: process.env.KUBERPC_CORE_URL,
  serviceName: "client",
});

const result = await rpc.call("server", "fib", { n: 10 });

console.log(`fib(10) = ${result}`);

rpc.close();
