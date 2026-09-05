# معماری پروژه Go Microservice

## 📋 نمای کلی معماری

```
┌─────────────────────────────────────────────────────────────────┐
│                          Clients (HTTP)                          │
└───────────────────────────┬─────────────────────────────────────┘
                            │
                            ▼
┌─────────────────────────────────────────────────────────────────┐
│                        API Gateway                               │
│                     (Port: 8085)                                 │
│                 - Request Routing                               │
│                 - Load Balancing                                 │
│                 - Authentication Check                           │
│                 - Request/Response Transformation              │
└────────────┬───────────────┬───────────────┬─────────────────────┘
             │               │               │
             ▼               ▼               ▼
    ┌─────────────┐  ┌─────────────┐  ┌─────────────┐
    │Auth Service │  │Post Service │  │Notification │
    │  (9092)     │  │  (9094)     │  │  Service    │
    │             │  │             │  │  (9093)     │
    │- Login      │  │- CRUD       │  │- Email      │
    │- Register   │  │- Search     │  │- Push       │
    │- Token      │  │- Comments   │  │- SMS        │
    │  Mgmt       │  │             │  │             │
    └──────┬──────┘  └──────┬──────┘  └──────┬──────┘
           │                │                │
           └────────────────┼────────────────┘
                            ▼
                    ┌───────────────┐
                    │   RabbitMQ    │
                    │  (5672/15672) │
                    │   Message     │
                    │    Broker     │
                    └───────┬───────┘
                            │
                ┌───────────┼───────────┐
                ▼           ▼           ▼
          ┌─────────┐  ┌─────────┐  ┌─────────┐
          │  Auth   │  │  Post   │  │Notifi-  │
          │  DB     │  │  DB     │  │cation   │
          │         │  │         │  │Service  │
          └─────────┘  └─────────┘  └─────────┘

┌─────────────────────────────────────────────────────────────────┐
│                   Shared PostgreSQL                              │
│                      (5432)                                     │
│                 - Users Table                                    │
│                 - Posts Table                                    │
│                 - Roles Table                                    │
│                 - Permissions Table                              │
└─────────────────────────────────────────────────────────────────┘
```

---

## 🎯 الگوهای معماری استفاده شده

### 1. **API Gateway Pattern**
API Gateway به عنوان single entry point عمل می‌کند:
- ✅ Routing requests به سرویس‌های مناسب
- ✅ Authentication و Authorization
- ✅ Load balancing
- ✅ Rate limiting (قابلیت آینده)

### 2. **Microservices Pattern**
هر سرویس مستقل است:
- ✅ توسعه و deployment جداگانه
- ✅ تکنولوژی‌های مختلف (هرچند اینجا همه Go هستند)
- ✅ Scaling جداگانه
- ✅ Failure isolation

### 3. **Event-Driven Architecture**
استفاده از RabbitMQ برای ارتباطات غیرهمگام:
- ✅ Decoupling سرویس‌ها
- ✅ Reliability و fault tolerance
- ✅ Scalability
- ✅ Event notification

### 4. **Domain-Driven Design (DDD)**
هر سرویس یک domain جداگانه را مدیریت می‌کند:
- ✅ Auth Domain
- ✅ Post Domain
- ✅ Notification Domain

---

## 🔄 Data Flow

### سناریوی 1: کاربر لاگین می‌کند

```
Client
  │
  │ 1. POST /api/auth/login
  ▼
API Gateway
  │
  │ 2. Route to Auth Service (gRPC)
  ▼
Auth Service
  │
  │ 3. Query User from DB
  ▼
PostgreSQL
  │
  │ 4. Return User data
  ▼
Auth Service
  │
  │ 5. Generate JWT Token
  │
  │ 6. Publish Login Event to RabbitMQ
  ▼
RabbitMQ
  │
  │ 7. Consume Event
  ▼
Notification Service (optional)
  │
  │ 8. Send Login Notification
  ▼
Client
  │
  │ 9. Return JWT Token
  ▼
Response
```

### سناریوی 2: کاربر پست جدید می‌سازد

```
Client
  │
  │ 1. POST /api/posts (with JWT)
  ▼
API Gateway
  │
  │ 2. Verify JWT Token
  │
  │ 3. Route to Post Service (gRPC)
  ▼
Post Service
  │
  │ 4. Validate User via Auth Service (gRPC)
  ▼
Auth Service
  │
  │ 5. Validate Token and Return User ID
  ▼
Post Service
  │
  │ 6. Save Post to DB
  ▼
PostgreSQL
  │
  │ 7. Return success
  ▼
Post Service
  │
  │ 8. Publish PostCreated Event to RabbitMQ
  ▼
RabbitMQ
  │
  │ 9. Multiple Consumers can react:
  │    - Notification Service (notify followers)
  │    - Analytics Service (track metrics)
  │    - Search Service (index post)
  ▼
Client
  │
  │ 10. Return success response
  ▼
Response
```

---

## 📦 ساختار سرویس‌ها

### API Gateway
**مسیر:** `services/api-gateway/`

**مسئولیت‌ها:**
- ✅ Routing درخواست‌ها به سرویس‌های مناسب
- ✅ Authentication و Authorization اولیه
- ✅ Load balancing
- ✅ Request validation

### Auth Service
**مسیر:** `services/auth-service/`

**مسئولیت‌ها:**
- ✅ User registration
- ✅ User login
- ✅ JWT token generation و validation
- ✅ Password management (bcrypt)
- ✅ Role-based access control (RBAC)
- ✅ Permission management

**APIs:**
```protobuf
service AuthService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
    rpc GetUserInfo(GetUserInfoRequest) returns (UserInfoResponse);
}
```

### Post Service
**مسیر:** `services/post-service/`

**مسئولیت‌ها:**
- ✅ Create, Read, Update, Delete (CRUD) پست‌ها
- ✅ Comment management
- ✅ Search functionality
- ✅ Feed generation
- ✅ Post filtering

**APIs:**
```protobuf
service PostService {
    rpc CreatePost(CreatePostRequest) returns (PostResponse);
    rpc GetPost(GetPostRequest) returns (PostResponse);
    rpc UpdatePost(UpdatePostRequest) returns (PostResponse);
    rpc DeletePost(DeletePostRequest) returns (DeleteResponse);
    rpc ListPosts(ListPostsRequest) returns (ListPostsResponse);
    rpc SearchPosts(SearchRequest) returns (SearchResponse);
}
```

### Notification Service
**مسیر:** `services/notification-service/`

**مسئولیت‌ها:**
- ✅ Email notifications
- ✅ Push notifications
- ✅ SMS notifications
- ✅ Notification templates
- ✅ Notification preferences
- ✅ Real-time notifications (via WebSockets - potentially)

**Events consumed:**
- User registered
- User logged in
- Post created
- Comment added
- Post liked

---

## 🗄 Data Layer

### Database Design

#### Users Table
```sql
CREATE TABLE users (
    id UUID PRIMARY KEY,
    email VARCHAR(255) UNIQUE NOT NULL,
    username VARCHAR(100) UNIQUE NOT NULL,
    password_hash VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### Posts Table
```sql
CREATE TABLE posts (
    id UUID PRIMARY KEY,
    user_id UUID NOT NULL,
    title VARCHAR(255) NOT NULL,
    content TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

#### Roles Table (RBAC)
```sql
CREATE TABLE roles (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);
```

#### Permissions Table
```sql
CREATE TABLE permissions (
    id UUID PRIMARY KEY,
    name VARCHAR(100) UNIQUE NOT NULL,
    description TEXT
);
```

#### Role-Permission Mapping
```sql
CREATE TABLE role_permissions (
    role_id UUID NOT NULL,
    permission_id UUID NOT NULL,
    FOREIGN KEY (role_id) REFERENCES roles(id),
    FOREIGN KEY (permission_id) REFERENCES permissions(id),
    PRIMARY KEY (role_id, permission_id)
);
```

---

## 🔐 امنیت معماری

### Authentication Flow
```
1. Client → API Gateway: Credentials
2. API Gateway → Auth Service: Validate credentials
3. Auth Service → DB: Query user
4. Auth Service: Generate JWT token
5. Auth Service → Client: Return JWT token
6. Client → API Gateway: JWT token
7. API Gateway → Auth Service: Validate JWT
8. Auth Service → API Gateway: Token valid/invalid
9. API Gateway → Target Service: Forward request (with user context)
```

### Authorization Layers
1. **API Gateway Level** - Basic authentication check
2. **Service Level** - RBAC implementation
3. **Database Level** - Row-level security (if needed)

---

## 📡 Communication Protocols

### Synchronous Communication
- **HTTP/REST** (Client ↔ API Gateway)
- **gRPC** (API Gateway ↔ Services, Service ↔ Service)

### Asynchronous Communication
- **AMQP/RabbitMQ** (Event publishing/subscription)

### Protocol Comparison
| Type | Usage | Pros | Cons |
|------|-------|------|------|
| HTTP/REST | Client-Gateway | Simple, universal | More payload overhead |
| gRPC | Service-to-Service | Fast, type-safe, streaming | More complex setup |
| AMQP | Event-driven | Decoupled, reliable | Additional infrastructure |

---

## 🚀 Deployment Architecture

### Environment Levels

#### Local Development (Tilt + Docker Desktop)
```
Developer's Machine
├── Tilt (Hot reload)
├── Docker Desktop (Kubernetes)
│   ├── API Gateway (Local build)
│   ├── Auth Service (Local build)
│   ├── Post Service (Local build)
│   ├── Notification Service (Local build)
│   ├── PostgreSQL (Persistent Volume)
│   └── RabbitMQ (StatefulSet)
└── Port Forwarding to localhost
```

#### Development (Kubernetes Cluster)
```
Development Cluster
├── Namespaces: microservice-dev
├── Services: Deployed as Kubernetes Deployments
├── Database: Managed PostgreSQL
└── Message Broker: Managed RabbitMQ
```

#### Production (Kubernetes Cluster)
```
Production Cluster
├── Namespaces: microservice-prod
├── Services: Deployed as Kubernetes Deployments
│   ├── Horizontal Pod Autoscaling
│   ├── Resource limits
│   └── Health checks
├── Database: Managed PostgreSQL with HA
├── Message Broker: Managed RabbitMQ with HA
└── ArgoCD: GitOps deployment management
```

---

## 🔄 CI/CD Pipeline

```
GitHub Actions
│
├── 1. Push to develop
│   └→ run-tests.yml
│   └→ run-build.yml
│   └→ deploy-dev.yml → Development Cluster (ArgoCD)
│
├── 2. Merge to main
│   └→ run-tests.yml
│   └→ run-build.yml
│   └→ Manual trigger deploy-prod.yml
│   └→ Production Cluster (ArgoCD)
│
└── 3. Workflow dispatch
    └→ deploy-prod.yml (requires "confirm" input)
```

---

## 🎯 معماری Event-Driven

### Event Examples
```protobuf
// User Events
message UserRegistered {
    string user_id = 1;
    string email = 2;
    string username = 3;
    google.protobuf.Timestamp timestamp = 4;
}

message UserLoggedIn {
    string user_id = 1;
    google.protobuf.Timestamp timestamp = 4;
}

// Post Events
message PostCreated {
    string post_id = 1;
    string user_id = 2;
    string title = 3;
    google.protobuf.Timestamp timestamp = 4;
}

// Notification Events
message NotificationRequested {
    string user_id = 1;
    string type = 2;
    string message = 3;
    google.protobuf.Timestamp timestamp = 4;
}
```

### Event Flow
```
Publisher (e.g., Post Service)
  │
  │ 1. Publish Event to RabbitMQ
  ▼
RabbitMQ Exchange (topic/direct)
  │
  │ 2. Route to Queues
  ▼
Multiple Queues
  │
  │ 3. Multiple Consumers
  ▼
┌─────────────┬─────────────┬─────────────┐
│Notification │  Analytics  │   Search    │
│  Service    │   Service   │   Service   │
└─────────────┴─────────────┴─────────────┘
```

---

## 🔧 Service Discovery

### Kubernetes Service Discovery
```
Service Name (K8s DNS)
  │
  ▼
auth-service.microservice-dev.svc.cluster.local
  │
  │ 9092
  ▼
Auth Service Pods
```

### gRPC Service Registration
```go
// Each service registers with Kubernetes Service
// API Gateway discovers services via DNS
```

---

## 📊 Monitoring & Observability (Future Enhancements)

### Proposed Stack
- **Metrics:** Prometheus + Grafana
- **Logging:** ELK Stack (Elasticsearch, Logstash, Kibana)
- **Tracing:** Jaeger or OpenTelemetry
- **Alerting:** Alertmanager

### Health Checks
```go
// Each service exposes health endpoint
GET /health
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2025-08-26T10:00:00Z"
}
```

---

## 🔒 Security Best Practices

### Current Implementation
- ✅ JWT for authentication
- ✅ bcrypt for password hashing
- ✅ Role-based access control
- ✅ Database encryption at rest (Kubernetes Secrets)
- ✅ Network policies (can be enhanced)

### Recommended Enhancements
- 🔒 OAuth 2.0 / OpenID Connect
- 🔒 Mutual TLS (mTLS) for service-to-service
- 🔒 Rate limiting and throttling
- 🔒 Input validation and sanitization
- 🔒 SQL injection prevention
- 🔒 CORS configuration

---

## 🚦 Availability & Scalability

### High Availability
- ✅ Kubernetes with multiple replicas
- ✅ StatefulSet for RabbitMQ
- ✅ Persistent volumes for PostgreSQL
- ✅ Health checks and liveness probes

### Scalability
- ✅ Horizontal scaling via Kubernetes HPA
- ✅ Vertical scaling via resource requests/limits
- ✅ Database connection pooling
- ✅ Message queue buffering

---

## 🎨 Architectural Principles

### SOLID Principles Applied
- ✅ **S**ingle Responsibility: Each service has one clear purpose
- ✅ **O**pen/Closed: Extensible without modification
- ✅ **L**iskov Substitution: Services can be replaced
- ✅ **I**nterface Segregation: Small, focused gRPC services
- ✅ **D**ependency Inversion: Depend on abstractions (interfaces)

### 12-Factor App Principles
- ✅ **I. Codebase:** One codebase per service
- ✅ **II. Dependencies:** Explicit dependencies (go.mod)
- ✅ **III. Config:** Config via environment variables
- ✅ **IV. Backing Services:** Treat attached services as resources
- ✅ **V. Build, Release, Run:** Strict separation
- ✅ **VI. Processes:** Stateless processes
- ✅ **VII. Port Binding:** Export services via port binding
- ✅ **VIII. Concurrency:** Scale via process model
- ✅ **IX. Disposability:** Maximize robustness with fast startup
- ✅ **X. Dev/Prod Parity:** Keep dev, staging, and prod similar
- ✅ **XI. Logs:** Treat logs as event streams
- ✅ **XII. Admin Processes:** Run admin/management tasks as one-off processes

---

## 📝 Conclusion

این معماری یک **modern, scalable, maintainable** microservices architecture است که:

1. ✅ از بهترین实践‌ها استفاده می‌کند
2. ✅ scalable و highly available است
3. ✅ separation of concerns را رعایت می‌کند
4. ✅ از event-driven patterns برای decoupling استفاده می‌کند
5. ✅ از Kubernetes برای orchestration استفاده می‌کند
6. ✅ از GitOps برای deployment management استفاده می‌کند
7. ✅ از CI/CD برای automation استفاده می‌کند

**Future Improvements:**
- 🔧 Add service mesh (Istio/Linkerd)
- 🔧 Implement caching layer (Redis)
- 🔧 Add rate limiting and API management
- 🔧 Improve monitoring and observability
- 🔧 Add API versioning
- 🔧 Implement circuit breakers and retries

---

**نسخه:** 1.0
**آخرین بروزرسانی:** 2025-08-26
