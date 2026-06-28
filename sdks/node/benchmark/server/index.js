import express from "express";
import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

const HTTP_PORT = Number(process.env.HTTP_PORT || 8081);
const serviceName = process.env.KUBERPC_SERVICE_NAME ?? "benchmark-server";

console.log(`[${serviceName}] starting - core: ${process.env.KUBERPC_CORE_URL}`);

const rpc = new KubeRPC();

function fib(n) {
  if (n <= 1) return n;
  let a = 0, b = 1;
  for (let i = 2; i <= n; i++) [a, b] = [b, a + b];
  return b;
}

function generate(count) {
  return Array.from({ length: count }, () => Math.random());
}

const app = express();

app.get("/", (_, res) => res.send("ok"));
app.get("/fib", (req, res) => {
  const n = Number(req.query.n);
  if (!Number.isInteger(n) || n < 0) return res.status(400).json({ error: "invalid n" });
  res.json({ result: fib(n) });
});
app.get("/generate", (req, res) => {
  const count = Number(req.query.count);
  if (!Number.isInteger(count) || count < 1) return res.status(400).json({ error: "invalid count" });
  res.json({ result: generate(count) });
});
app.get("/ping", (_, res) => res.json({ result: "pong" }));

app.listen(HTTP_PORT, async () => {
  console.log(`[${serviceName}] HTTP listening on :${HTTP_PORT}`);

  try {
    await rpc.register("fib", async ({ n }) => {
      const result = fib(n);
      console.log(`[${serviceName}] rpc:fib       n=${n}  →  ${result}`);
      return result;
    });
    await rpc.register("generate", async ({ count }) => {
      const result = generate(count);
      console.log(`[${serviceName}] rpc:generate  count=${count}  →  ${count} floats`);
      return result;
    });
    await rpc.register("ping", async () => "pong");

    console.log(`[${serviceName}] RPC ready on :${process.env.KUBERPC_PORT ?? 7749}  (methods: fib, generate, ping)`);
  } catch (err) {
    console.error(`[${serviceName}] register failed:`, err.message);
    process.exit(1);
  }
});
