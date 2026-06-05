# 🚀 GoApp

A **modular, scalable Go-based business application framework** designed to build modern web applications with clean architecture, high extensibility, and strong separation of concerns.

GoApp helps you structure large applications into **independent, plug-and-play modules** while maintaining a simple and predictable core system.

---

## ✨ Features

- 🧩 Fully modular architecture (plugin-style system)
- ⚡ High-performance HTTP server powered by Fiber
- 🗄️ Database integration using Ent ORM
- 🔌 Automatic module bootstrap system
- 📦 Clean architecture: Handler → Service → Repository → DB
- 🔄 Event-driven communication system (EventBus)
- ⚙️ Centralized configuration via environment variables
- 🧠 Scalable core designed for large applications
- 🔐 Middleware support (logging, CORS, recovery, etc.)
- 🚀 Built for long-term maintainability and extensibility

---

## 🏗️ Architecture Overview

GoApp follows a **core + modules architecture**:

```
Core
├── Application   (Bootstrap layer)
├── Container     (Dependencies)
├── Registry      (Module loader)
├── EventBus      (Communication layer)
└── HTTP Server   (Fiber)

Modules
├── Handler       (HTTP layer)
├── Service       (Business logic)
├── Repository    (Data access layer)
└── Routes        (API definitions)
```

Each module is **independent, self-contained, and reusable**.

---

## 📦 Module System

Modules in GoApp are designed to be:

- Self-registering
- Independently testable
- Isolated from other modules
- Extensible via services or events

**Example module structure:**

```
/modules/example
├── handler/
├── service/
├── repository/
└── module.go
└── manifest.json
```

---

## 📁 Project Structure

```
/core        → Core framework (registry, config, container)
/modules     → Feature modules (business logic)
/ent         → ORM generated code
/internal    → Internal utilities
main.go      → Application entry point
```

---

## ⚙️ Configuration

GoApp uses environment-based configuration:

- `.env` file supported out of the box
- Central `Config` struct
- Supports development and production environments
- Feature toggles for enabling/disabling modules

---

## 🔄 Event System

GoApp includes a lightweight **EventBus system**:

- Publish/Subscribe pattern
- Decoupled module communication
- Useful for cross-module workflows
- Keeps modules fully independent

---

## 🚀 Getting Started

### 1. Clone the repository

```bash
git clone https://github.com/subhasundardass/goapp
cd goapp
```

### 2. Install dependencies

```bash
go mod tidy
```

### 3. Set up environment

```bash
cp .env.example .env
```

### 4. Run the application

```bash
go run cmd/web/main.go
or
make dev (hot reload)
```

---

## 📖 Documentation

Full architecture details, module development guide, and system design explanation are available here:

👉 [DOCUMENTATION.md](./DOCUMENTATION.md)

---

## 🎯 Design Principles

GoApp is built with these principles:

- **Simplicity over complexity** — avoid unnecessary abstractions
- **Modular independence** — each module stands on its own
- **Clear separation of concerns** — no logic bleeds between layers
- **Scalable architecture from day one** — designed to grow
- **Event-driven communication where needed** — loose coupling via EventBus

---

## 🧠 Ideal Use Cases

GoApp is suitable for building:

- Business applications
- SaaS platforms
- Admin dashboards
- Internal tools
- API-driven backend systems

---

## 🛣️ Roadmap

- [ ] Module dependency system
- [ ] Background job processing
- [ ] Multi-tenant support
- [ ] Plugin marketplace system
- [ ] Admin UI panel for module management
- [ ] Advanced permissions system

---

## 🤝 Contribution

GoApp is designed as a foundation framework, and contributions are welcome in:

- New modules
- Core enhancements
- Middleware extensions
- Event system improvements

---

## 📜 License

[MIT License](./LICENSE)

---

## ⭐ About

GoApp is a modular backend framework designed for building scalable business applications in Go with clean architecture and long-term extensibility in mind.
