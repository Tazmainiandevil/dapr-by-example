# Dapr Distributed Systems Series — Sample Code & Reference Implementation

This repository contains the complete code and configuration that accompanies the 7‑part Dapr series, the Kubernetes appendix, and the bonus Aspire post. It provides runnable examples in Go and .NET, along with Dapr components, Kubernetes manifests, and documentation to help you build real‑world distributed systems without infrastructure‑specific SDKs.

The repository mirrors the flow of the series:

- Start locally  
- Add state  
- Add pub/sub  
- Add bindings  
- Add observability  
- Build an end‑to‑end service  
- Deploy to Kubernetes  
- (Optional) Orchestrate with .NET Aspire  

Whether you're following the posts or exploring independently, this repo gives you everything you need to run the examples end‑to‑end.

---

## 📚 Blog Series

Each part of the series is linked below.  
The code in this repo is organised to match these posts.

| Part | Title | Link |
|------|--------|------|
| Part 1 | What is Dapr, and Why Would You Use It? | 🔗 [Read: Part 1](https://codingwithtaz.blog/2026/02/02/part-1-introduction-to-dapr-a-practical-guide-to-reducing-glue-code-in-distributed-systems/) |
| Part 2 | Running Dapr Locally | 🔗 [Read: Part 2](https://codingwithtaz.blog/2026/02/05/part-2-running-dapr-locally-setup-run-and-debug-your-first-service/) |
| Part 3 | State Management with Dapr | 🔗 [Read: Part 3](https://codingwithtaz.blog/2026/02/09/part-3-state-management-with-dapr-redis-and-postgres-without-the-sdks/) |
| Part 4 | Event‑Driven Systems with Pub/Sub | 🔗 [Read: Part 4](https://codingwithtaz.blog/2026/02/16/part-4-event-driven-systems-with-dapr-pub-sub/) |
| Part 5 | Bindings & Storage | 🔗 [Read: Part 5](https://codingwithtaz.blog/2026/02/23/part-5-integrating-external-systems-with-dapr-bindings-and-storage/) |
| Part 6 | Observability with Dapr | 🔗 Coming Soon |
| Part 7 | Putting It All Together | 🔗 Coming Soon |
| Appendix | Real‑World Dapr Configuration for Kubernetes | 🔗 Coming Soon |
| Bonus | Using Dapr with .NET Aspire | 🔗 Coming Soon |

The repo links back to the posts, and the posts link back to the repo — a two‑way learning loop.

---

## 🧭 Who This Repo Is For

- Developers learning Dapr for the first time  
- Teams evaluating Dapr for distributed systems  
- Engineers wanting Go and .NET examples side‑by‑side  
- Platform teams exploring Dapr in Kubernetes  
- Readers following the 7‑part blog series  
- Anyone looking for a clean, infrastructure‑agnostic reference architecture  

---

## 🚀 Start Here

1. Install prerequisites: Dapr CLI, Docker, Go, .NET SDK  
2. Initialise Dapr locally: `dapr init`  
3. Run the order service (Go or .NET)  
4. Trigger an order  
5. Watch the trace appear in Zipkin  
6. Add the inventory service  
7. Explore the Dapr components  
8. Move to Kubernetes when ready  

This mirrors the learning flow of the blog series.

---

## 🧩 Architecture Overview

A typical Dapr‑enabled service in this repo looks like:

```shell
Client 
  ↓ 
Order Service (Go / .NET) 
  ↓ 
Dapr Sidecar 
  ├─ State Store (Redis / Postgres) 
  ├─ Pub/Sub Broker (Redis / Kafka) 
  └─ Storage Provider (Local / S3 / Azure Blob)
```


Your application talks only to Dapr.  
Dapr talks to the infrastructure.  
This keeps your code clean, portable, and infrastructure‑agnostic.

---

## 📁 Repository Structure

```shell
dapr-by-example/ 
├── components/                   # Local development Dapr components
|   ├── config.yaml               # Dapr Tracing configuration for observability
│   ├── pubsub.yaml               # Redis pub/sub for local dev 
│   ├── statestore.yaml           # Redis state store for local dev 
│   └── storage.yaml              # Local file storage binding 
├── docs/                         # Documentation 
│   ├── architecture.md           # System architecture overview 
│   ├── kubernetes.md             # Kubernetes deployment guide 
│   └── local-dev.md              # Local development setup 
├── k8s/                          # Kubernetes production manifests 
│   ├── components/               # Production Dapr components 
│   │   ├── azure-creds.yaml      # Example secret for Azure Storage
│   │   ├── config-jaeger.yaml    # Dapr configuration for Jaeger (tracing, metrics)
│   │   ├── pg-secret.yaml        # PostgreSQL credentials
│   │   ├── pubsub.yaml           # Kafka pub/sub broker
│   │   ├── s3-creds.yaml         # AWS S3 credentials 
│   │   ├── secretstores.yaml     # Secret store configuration 
│   │   ├── statestore.yaml       # PostgreSQL state store 
│   │   └── storage-aws.yaml      # AWS S3 storage binding 
│   │   └── storage-azure.yaml    # Azure Storage account binding 
│   └── deployments/              # Service deployment manifests
│   │   └── dotnet/
│   │       ├── inventoryservice-dotnet.yaml 
│   │       └── orderservice-dotnet.yaml 
│   │   └── go/
│   │       ├── inventoryservice-gp.yaml 
│   │       └── orderservice-go.yaml 
├── src/                          # Service implementations 
│   ├── inventoryservice-dotnet/ 
│   ├── inventoryservice-go/ 
│   ├── orderservice-dotnet/ 
│   └── orderservice-go/ 
└── README.md
```


**What lives where**

- `/src` — All Go and .NET services used in the series  
- `/components` — Dapr components for local development  
- `/k8s` — Production‑oriented manifests and components  
- `/docs` — Additional documentation and architecture notes  

---

## 🧪 Running Everything Locally

Run any service using the Dapr CLI:

Go

```shell
cd src/orderservice-go
dapr run --app-id orderservice --app-port 8080 --resources-path ../../components -- go run main.go
```

```shell
cd src/inventoryservice-go
dapr run --app-id inventoryservice --app-port 8081 --resources-path ../../components -- go run main.go
```

.NET

```shell
cd src/orderservice-dotnet
dapr run --app-id orderservice --app-port 8080 --resources-path ../../components -- dotnet run
```

```shell
cd src/inventoryservice-dotnet
dapr run --app-id inventoryservice --app-port 8081 --resources-path ../../components -- dotnet run
```

with tracing:
```
dapr run --app-id orderservice --app-port 8080 --resources-path ../../components --config ../../components/config.yaml -- go run main.go
```

### Service Endpoints

**Order Service (Port 8080):**
- `POST /orders` - Create new order
- `GET /orders/{orderId}` - Retrieve order by ID
- `GET /dapr/subscribe` - Subscription discovery (empty for publishers)
- `GET /healthz` - Health check endpoint

**Inventory Service (Port 8081):**
- `POST /orders` - Process incoming order events
- `GET /dapr/subscribe` - Subscription discovery
- `GET /healthz` - Health check endpoint
