# Handworks Cleaning Services System

A centralized, scalable service platform designed to optimize booking, scheduling, operations, and administration for Handworks Cleaning Services. Built on distributed microservices, the system supports cross-platform access and automates key workflows across clients, cleaners, and admin roles.

---

## System Diagram

![Handworks Cleaning System Diagram](https://github.com/Gnashal/handworks-microservice/blob/D14N6C0/.github/diagrams/handwords_system_backend_diagram.png)

---

## Overview

This system is built using a **distributed microservices architecture**, breaking down complex business functions (e.g. bookings, workforce, payments, inventory) into independently developed and scalable services. It supports three user roles:
- **Clients** (web/mobile bookings)
- **Cleaners/Employees** (mobile operations)
- **Administrators/Managers** (desktop app + dashboards)

---

## System Architecture

### 🔹 Client-Facing Applications
- **Mobile App** (React Native)
  - Book and pay for services
  - Upload job-related photos
  - Submit feedback
- **Web App** (Vite + React on Netlify)
  - Mirrors mobile functionality
  - Fast and accessible on any browser

### 🔹 Cleaner (Employee) Application
- **React Native mobile app**
  - Receive real-time job notifications
  - View tasks and update progress
  - Upload before/after photos
  - Syncs updates with backend

### 🔹 Admin Desktop Application
- **Electron + React**
  - Management dashboards
  - Scheduling and resource allocation
  - Works offline, syncs on reconnect
  - Advanced reporting and analytics

---

## Backend Microservices

Each domain (e.g. bookings, payments, inventory) runs as an independent Go (Golang) microservice connected to a dedicated PostgreSQL database (on Neon). Core backend features:

- Modular, scalable, and secure design
- REST/GraphQL/gRPC APIs
- Fine-grained access control per service

---

## Communication & Data Flow

### API Gateway
- Built in **Go (Gin framework)**
- Performs:
  - Request validation
  - OAuth 2.0 authentication (via Clerk)
  - Rate limiting
  - Routing to internal services

### GraphQL Domain Services
- **GraphQL endpoints** per business domain
- Precision queries and dashboards
- Strong typing with `gqlgen`

### gRPC Inter-Service Communication
- Used for high-performance operations like:
  - Cleaner assignments
  - Inventory status updates
- Powered by **Protocol Buffers**

### Event-Driven Messaging
- **NATS + JetStream**
  - Publish-subscribe architecture
  - Used for:
    - Booking notifications
    - Payment confirmations
    - Inventory updates
  - Adds fault tolerance and durability

### Integrated TLS Security
- TLS on all communication layers:
  - REST
  - GraphQL
  - gRPC
  - NATS
- JWT-based access control via Clerk
- Role-based authorization at every layer

---

## Technology Stack Summary

| Tool                     | Purpose |
|--------------------------|---------|
| **React Native (Expo)** | Mobile app for clients & cleaners |
| **Vite + React**         | Web app (fast & responsive) |
| **Electron + React**     | Desktop app for admins |
| **Go (Golang)**          | Backend business logic |
| **PostgreSQL (Neon)**    | Cloud-hosted relational DB per microservice |
| **Gin**                  | REST API Gateway |
| **GraphQL + gqlgen**     | Flexible frontend queries |
| **gRPC**                 | Inter-service communication |
| **NATS + JetStream**     | Event-driven messaging |
| **Clerk**                | OAuth 2.0, JWT auth, RBAC |
| **Cloudinary**           | Image storage & delivery |
| **Render / Netlify**     | Backend + frontend hosting |
| **GitHub + Actions**     | CI/CD pipelines, testing, deployment |

---

## Security Model

- All client and inter-service communication is encrypted (TLS).
- API Gateway manages:
  - Authentication via Clerk (OAuth2, JWT)
  - Authorization and rate limiting
- Role-based permissions are enforced at service endpoints.
- Desktop admin app supports **offline operations** with automatic sync.

---

## Quality Assurance

- **Unit tests** per microservice
- **Integration tests** for workflows (e.g. booking + assignment)
- **End-to-end tests** simulated through the API Gateway
- Automated CI/CD checks using **GitHub Actions**

---

## Technical Highlights

- **Scalability** via microservices and independent deployments
- **Cross-platform**: mobile (iOS/Android), desktop, and browser-based access
- **Flexible data access**: REST + GraphQL + gRPC
- **High performance**: Go backend + async messaging
- **Developer-friendly**: clean modular code, GitHub workflows, and type-safe API contracts

---

## Development Tools & Practices

- **Version control**: Git + GitHub
- **Formatting/linting**: ESLint, Prettier
- **Media handling**: Cloudinary SDK
- **Frontend builds**: Vite
- **Deployment**: Render (backend), Netlify (frontend)
- **Monitoring/Tracing**: (Planned)

---

## Final Considerations

The Handworks Cleaning Services platform delivers a robust, scalable architecture optimized for fast booking, real-time cleaner coordination, admin analytics, and cross-platform access. Its modular design ensures future updates, expansion, and third-party integrations can be done without disrupting core operations.

---

_This README was auto-generated from the project’s technical documentation. For more details, refer to the full project wiki or reach out to the system architect._
