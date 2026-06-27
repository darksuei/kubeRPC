import { KubeRPC } from "kuberpc-node";
import dotenv from "dotenv";

dotenv.config();

const SERVER_HTTP = process.env.SERVER_HTTP_URL || "http://localhost:8081";

const rpc = new KubeRPC({
  coreURL: process.env.KUBERPC_CORE_URL,
  serviceName: "benchmark-client",
});

async function bench(fn, runs = 10) {
  await fn(); // warmup
  let total = 0;
  for (let i = 0; i < runs; i++) {
    const t = performance.now();
    await fn();
    total += performance.now() - t;
  }
  return total / runs;
}

async function burstTotal(fn, count) {
  const t = performance.now();
  for (let i = 0; i < count; i++) await fn();
  return performance.now() - t;
}

function report(rpcMs, httpMs) {
  const pct = (((httpMs - rpcMs) / httpMs) * 100).toFixed(1);
  console.log(`  kubeRPC  ${rpcMs.toFixed(2)}ms`);
  console.log(`  HTTP     ${httpMs.toFixed(2)}ms`);
  console.log(`  kubeRPC was ${pct}% faster`);
}

try {
  // ── 1. computation ──────────────────────────────────────────────────────────
  console.log("\n── 1. computation  fib(40), 10 runs ──");
  report(
    await bench(() => rpc.call("benchmark-server", "fib", { n: 40 })),
    await bench(async () => (await (await fetch(`${SERVER_HTTP}/fib?n=40`)).json()).result),
  );

  // ── 2. data transfer ────────────────────────────────────────────────────────
  console.log("\n── 2. data transfer  generate(5000 floats), 5 runs ──");
  report(
    await bench(() => rpc.call("benchmark-server", "generate", { count: 5000 }), 5),
    await bench(async () => (await (await fetch(`${SERVER_HTTP}/generate?count=5000`)).json()).result, 5),
  );

  // ── 3. burst ────────────────────────────────────────────────────────────────
  const PINGS = 100;
  console.log(`\n── 3. burst  ${PINGS}x sequential ping ──`);
  const rpcBurst = await burstTotal(() => rpc.call("benchmark-server", "ping"), PINGS);
  const httpBurst = await burstTotal(async () => (await fetch(`${SERVER_HTTP}/ping`)).json(), PINGS);
  const pct = (((httpBurst - rpcBurst) / httpBurst) * 100).toFixed(1);
  console.log(`  kubeRPC  ${rpcBurst.toFixed(2)}ms total`);
  console.log(`  HTTP     ${httpBurst.toFixed(2)}ms total`);
  console.log(`  kubeRPC was ${pct}% faster`);

  console.log("\n── done ──\n");
} catch (err) {
  console.error("benchmark failed:", err.message);
  process.exit(1);
}

rpc.close();
