# Handworks Cleaning Services System

A centralized, scalable service platform designed to optimize booking, scheduling, operations, and administration for Handworks Cleaning Services. Built on distributed microservices, the system supports cross-platform access and automates key workflows across clients, cleaners, and admin roles.

---
# How to Run
- First run the gateway
```
cd ./gateway
go run main.go
```
- Then run each service (if applicable when you test)
```
cd .services/{service name eg.. account}
go run main.go
```

## System Diagram

![Handworks Cleaning System Diagram](https://github.com/Gnashal/handworks-microservice/blob/D14N6C0/.github/diagrams/handwords_system_backend_diagram_V2.png)

---

## Overview

The system, **“Optimizing Service Delivery through a Centralized Booking and Operations Management System for Handworks Cleaning Services”**, is built using a **distributed microservices architecture**, breaking down business domains (e.g., bookings, workforce, payments, inventory) into independently developed modules. This ensures scalability, maintainability, and fault tolerance.

It supports three main user roles:

* **Clients** (web portal)
* **Cleaners/Employees** (dedicated Android app)
* **Administrators/Managers** (desktop app)

---

## System Architecture

### 🔹 Client Web Application

* **React + Vite (deployed on Netlify)**
* Fully responsive design (desktop + mobile)
* Features:

  * Service browsing & booking
  * Online payment (30% down payment required)
  * Upload job-related images
  * Post-job feedback
* **Guest bookings are not permitted**; login with verified account is required.

### 🔹 Cleaner (Employee) Mobile Application

* **Android SDK (Java + XML)**
* Proprietary, internally distributed app (not Play Store)
* Features:

  * Real-time job notifications
  * Task viewing and progress updates
  * Upload before/after photos
  * Sync with backend for coordination

### 🔹 Admin Desktop Application

* **Electron + React**
* Features:

  * Management dashboards
  * Scheduling & resource allocation
  * Reports & analytics
  * Works offline with auto-sync when reconnected

---

## Backend Microservices

Each domain (bookings, employees, payments, inventory) is an independent **Go (Golang)** microservice with its own **PostgreSQL (Neon)** database. Key features:

* Modular and scalable design
* REST, GraphQL, and gRPC APIs
* Role-based access control
* Cloud-hosted per service database

---

## Communication & Data Flow

### API Gateway (REST)

* **Go (Gin framework)**
* Functions:

  * Request validation
  * OAuth 2.0 authentication (via **Clerk**)
  * JWT-based access control
  * Rate limiting
  * Routing requests to services

### GraphQL Domain Services

* **gqlgen (Go)** per domain
* Precise data queries
* Strong typing and flexible dashboards

### gRPC Inter-Service Communication

* High-performance service-to-service calls
* Used for:

  * Cleaner assignments
  * Inventory checks
  * Real-time booking workflows

### Event-Driven Messaging

* **NATS + JetStream**
* Publish-subscribe messaging backbone
* Handles:

  * Booking creation events
  * Assignment & notification delivery
  * Payment confirmations
  * Inventory updates
* JetStream adds durability & message persistence

### Integrated TLS Security

* TLS for **REST, GraphQL, gRPC, and NATS**
* Clerk-managed JWT tokens
* Role-based enforcement at API gateway & service level

---

## Technology Stack Summary

| Tool/Framework               | Purpose                                          |
| ---------------------------- | ------------------------------------------------ |
| **Android SDK (Java + XML)** | Mobile app for cleaners (internal distribution)  |
| **Vite + React**             | Client web portal (bookings, payments, feedback) |
| **Electron + React**         | Desktop app for admins (offline sync, analytics) |
| **Go (Golang)**              | Backend microservices and business logic         |
| **PostgreSQL (Neon)**        | Cloud-hosted DB per microservice                 |
| **Gin**                      | REST API gateway                                 |
| **GraphQL + gqlgen**         | Flexible queries & reporting                     |
| **gRPC**                     | Inter-service communication                      |
| **NATS + JetStream**         | Event-driven messaging                           |
| **Clerk**                    | Authentication & role-based security             |
| **Cloudinary**               | Image storage & delivery                         |
| **Render / Netlify**         | Backend + frontend hosting                       |
| **GitHub + Actions**         | CI/CD automation & testing                       |

---

## Security Model

* **TLS encryption** across all communication channels
* **Clerk OAuth2 + JWT** for authentication & authorization
* Role-based permissions enforced at API Gateway & microservice endpoints
* Offline support for admin app with secure re-sync
* Durable messaging via **NATS JetStream**

---

## Quality Assurance

* **Unit tests** for each microservice
* **Integration tests** for workflows (e.g., booking + assignment)
* **End-to-end tests** via API Gateway
* Automated testing pipelines via **GitHub Actions**

---

## Technical Highlights

* **Scalability**: modular microservices with independent databases
* **Cross-platform**: Web (React), Mobile (Android SDK), Desktop (Electron)
* **Flexible APIs**: REST, GraphQL, gRPC
* **High performance**: Go backend + NATS JetStream for async events
* **Offline resilience**: Admin app syncs when reconnected
* **Developer-friendly**: GitHub CI/CD, ESLint, Prettier, code-first GraphQL

---

## Development Tools & Practices

* **Version control**: Git + GitHub
* **CI/CD**: GitHub Actions
* **Formatting/linting**: ESLint, Prettier
* **Media handling**: Cloudinary SDK
* **Frontend builds**: Vite
* **Deployment**: Render (backend), Netlify (frontend)
* **Monitoring & tracing**: planned

---

## Final Considerations

The Handworks Cleaning Services platform delivers a robust, modular architecture optimized for:

* Seamless client booking
* Real-time cleaner coordination
* Secure admin workflows
* Scalable, fault-tolerant expansion

Its distributed design ensures that future updates, integrations, and scaling can be achieved without disrupting core system operations.
