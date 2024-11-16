# kubeRPC

**kubeRPC** is a **Kubernetes-native remote procedure call (RPC) framework** designed to enable seamless and low-latency communication between microservices deployed within the same Kubernetes cluster or namespace.

## **Why kubeRPC?**

One of the challenges when transitioning from a monolithic architecture to microservices is **latency**. In a monolith, methods (e.g., `generateInvoice`) can be called directly with negligible overhead. In contrast, microservices require exposing APIs (HTTP, GraphQL, SOAP, etc.) that add significant latency and complexity to method calls.

**kubeRPC** solves this problem by enabling microservices to directly invoke each other's methods without relying on traditional API endpoints. This dramatically reduces latency, simplifies communication, and preserves the speed of direct method calls.

<div style="text-align: center;">
  <img src="./assets/RPC-diagram.png" alt="RPC Overview" width="700">
</div>

---

## **How kubeRPC Works**

1. kubeRPC deploys a **core service** within your Kubernetes cluster (written in Go) that acts as the central orchestrator.
2. kubeRPC **watches for all services** in the namespace and automatically registers their DNS names and other relevant metadata.
3. Microservices can **register callable methods** with the kubeRPC core service.
4. Other microservices can invoke these methods using the kubeRPC SDK, eliminating the need for HTTP or similar overhead.

---

## **Setup and Deployment**

### **Requirements**
- A Kubernetes cluster (any version compatible with Helm).
- Helm installed on your local machine.

### **Deploying kubeRPC**

kubeRPC can be deployed using a Helm chart.

```bash
helm install kuberpc-core https://github.com/darksuei/kubeRPC/raw/main/helm_chart/kuberpc-core-0.1.0.tgz --namespace <your-namespace> -f /path/to/custom-values.yaml
```

Once installed, kubeRPC will:
1. Monitor all services in the namespace.
2. Register their DNS names and metadata for service-to-service communication.

---

## **Usage**

### **Registering Methods**
To register methods, services must interact with the kubeRPC core using the kubeRPC SDK.

### **Calling Methods**
Once methods are registered, other services can directly invoke these methods using the SDK.

---

## **Available SDKs**

Currently, a **TypeScript SDK** is available:
- TypeScript SDK: [Source code](https://github.com/darksuei/kubeRPC-sdk) | [NPM](https://www.npmjs.com/package/kuberpc-sdk)

We welcome contributions for SDKs in other programming languages!

---

## **Contributing**

We encourage developers to contribute by:
1. Building SDKs for other languages.
2. Reporting bugs or requesting features via GitHub issues.
3. Submitting pull requests to improve the core or SDKs.

Check out the [Contribution Guide](#N/A) for more details.
