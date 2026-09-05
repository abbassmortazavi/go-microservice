# داکیومنت فنی پروژه Go Microservice

---

## 📋 فهرست مطالب

1. [مقدمه](#مقدمه)
2. [ساختار پروژه](#ساختار-پروژه)
3. [معماری سیستم](#معماری-سیستم)
4. [سرویس‌های سیستم](#سرویسهای-سیستم)
5. [تکنولوژی‌های استفاده شده](#تکنولوژیهای-استفاده-شده)
6. [دیتابیس و مدل داده](#دیتابیس-و-مدل-داده)
7. [جریان داده](#جریان-داده)
8. [امنیت](#امنیت)
9. [استقرار و Deployment](#استقرار-deployment)
10. [CI/CD Pipeline](#cicd-pipeline)
11. [مزایا و چالش‌ها](#مزایا-و-چالشها)

---

## 📌 مقدمه

این پروژه یک سیستم میکروسرویس مدرن با استفاده از زبان برنامه‌نویسی Go است که از معماری Monorepo برای مدیریت سرویس‌های مختلف استفاده می‌کند. این سیستم شامل سرویس‌های جداگانه برای احراز هویت، مدیریت پست‌ها، و اطلاع‌رسانی است که از طریق API Gateway با هم در ارتباط هستند.

### اهداف پروژه
- توسعه‌پذیری و نگهداری آسان
- مقیاس‌پذیری و بالا بودن قابلیت دسترسی (Availability)
- جدا بودن سرویس‌ها (Loose Coupling)
- استفاده از معماری Event-Driven
- پیاده‌سازی الگوهای Microservices به‌درستی

---

## 🏗️ ساختار پروژه

```
go-microservice/
├── services/                      # سرویس‌های اصلی
│   ├── api-gateway/              # دروازه API
│   ├── auth-service/             # سرویس احراز هویت
│   ├── post-service/             # سرویس پست‌ها
│   └── notification-service/     # سرویس اطلاع‌رسانی
│
├── pkg/                          # کدهای مشترک
│   ├── database/                 # لایه دیتابیس
│   ├── message/                  # پیام‌رسانی (RabbitMQ)
│   ├── utils/                    # ابزارهای عمومی
│   ├── events/                   # رویدادها
│   ├── config/                   # تنظیمات
│   └── proto/                    # فایل‌های Protocol Buffers
│
├── migrations/                   # Database migrations
├── proto/                        # فایل‌های .proto
├── infra/                        # تنظیمات زیرساخت
├── argocd/                       # تنظیمات ArgoCD
├── build/                        # فایل‌های build
├── tools/                        # ابزارهای توسعه
│
├── .github/                      # تنظیمات GitHub Actions
├── Tiltfile                      # تنظیمات Tilt برای توسعه محلی
├── Makefile                      # دستورات build و run
├── go.mod                        # مدیریت dependencyها
├── go.sum                        # checksum dependencyها
└── README.md                     # توضیحات پروژه
```

### توضیح دایرکتوری‌های کلیدی

#### services/
هر سرویس یک دایرکتوری جداگانه دارد که شامل کدهای آن سرویس است:
- **api-gateway/**: نقطه ورود اصلی سیستم
- **auth-service/**: مدیریت کاربران و احراز هویت
- **post-service/**: مدیریت پست‌ها و کامنت‌ها
- **notification-service/**: ارسال اطلاع‌رسانی‌ها

#### pkg/
کدهای مشترک بین سرویس‌ها:
- **database/**: اتصال و عملیات دیتابیس
- **message/**: ناشر و مصرف‌کننده پیام‌ها
- **utils/**: توابع عمومی و helper
- **events/**: ساختار رویدادهای سیستم
- **config/**: تنظیمات و کانفیگ‌ها

#### migrations/
فایل‌های SQL برای ساخت و تغییر ساختار دیتابیس

#### proto/
فایل‌های Protocol Buffers برای تعریف APIهای gRPC

---

## 🎯 معماری سیستم

### نمودار معماری کلی

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

### الگوهای معماری استفاده شده

#### 1. API Gateway Pattern
- API Gateway به عنوان single entry point عمل می‌کند
- درخواست‌ها را به سرویس‌های مناسب route می‌کند
- Authentication و Authorization اولیه را انجام می‌دهد
- Load balancing را مدیریت می‌کند

#### 2. Microservices Pattern
- هر سرویس مستقل است
- می‌توان هر سرویس را جداگانه توسعه و deploy کرد
- هر سرویس می‌تواند جداگانه scale شود
- Failure isolation در سرویس‌ها

#### 3. Event-Driven Architecture
- استفاده از RabbitMQ برای ارتباطات غیرهمگام
- Decoupling سرویس‌ها
- Reliability و fault tolerance
- Event notification

#### 4. Domain-Driven Design (DDD)
- هر سرویس یک domain جداگانه را مدیریت می‌کند
- Auth Domain: مدیریت کاربران و احراز هویت
- Post Domain: مدیریت محتوا
- Notification Domain: ارسال اطلاع‌رسانی‌ها

---

## 🚀 سرویس‌های سیستم

### API Gateway
**مسیر:** `services/api-gateway/`
**پورت:** 8085

**مسئولیت‌ها:**
- ✅ Routing درخواست‌ها به سرویس‌های مناسب
- ✅ Authentication و Authorization اولیه
- ✅ Load balancing
- ✅ Request validation
- ✅ Request/Response transformation

### Auth Service
**مسیر:** `services/auth-service/`
**پورت:** 9092

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
**پورت:** 9094

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
**پورت:** 9093

**مسئولیت‌ها:**
- ✅ Email notifications
- ✅ Push notifications
- ✅ SMS notifications
- ✅ Notification templates
- ✅ Notification preferences

**Events consumed:**
- User registered
- User logged in
- Post created
- Comment added
- Post liked

---

## 💻 تکنولوژی‌های استفاده شده

### زبان برنامه‌نویسی
- **Go 1.25.5**: زبان اصلی سرویس‌ها

### Frameworks و Libraryها

#### Routing
- **gorilla/mux v1.8.1**: HTTP router

#### Communication
- **gRPC v1.77.0**: برای ارتباط بین سرویس‌ها
- **Protocol Buffers v1.36.11**: برای serialization
- **RabbitMQ AMQP v1.10.0**: برای message broker

#### Authentication & Security
- **golang-jwt/jwt v5.3.0**: برای JWT token management
- **golang.org/crypto v0.46.0**: برای encryption و hashing

#### Database
- **lib/pq v1.10.9**: PostgreSQL driver
- **PostgreSQL**: دیتابیس اصلی سیستم

#### Utilities
- **google/uuid v1.6.0**: برای تولید UUID
- **go-playground/validator v10.29.0**: برای input validation

---

## 📡 gRPC Communication

### gRPC چیست؟

gRPC (Google Remote Procedure Call) یک فریمورک RPC مدرن و با کارایی بالا است که از Protocol Buffers برای تعریف سرویس‌ها و serialization داده‌ها استفاده می‌کند. در این پروژه از gRPC برای ارتباط بین سرویس‌ها (به‌ویژه Auth Service و Post Service) استفاده می‌شود.

### ✅ مزایای gRPC

#### 1. کارایی بالا (Performance)
- **HTTP/2**: استفاده از HTTP/2 برای multiplexing و reduction overhead
- **Protocol Buffers**: serialization باینری که 3-10 برابر سریع‌تر و کوچک‌تر از JSON است
- **Binary Format**: پیام‌های باینری کمتر از JSON/XML فضا اشغال می‌کنند
- **Compression**: فشرده‌سازی خودکار داده‌ها

#### 2. تایپ‌های قوی (Strong Typing)
- **Compile-time Type Safety**: خطاهای تایپ در زمان کامپایل شناسایی می‌شوند
- **Code Generation**: تولید خودکار client و server stubs از فایل‌های .proto
- **Interface Definition**: تعریف واضح و یکنواخت APIها
- **Schema Evolution**: مدیریت بهتر تغییرات در schema

#### 3. استریم دوطرفه (Bidirectional Streaming)
```protobuf
// مثال استریم دوطرفه
rpc ChatStream(stream ChatMessage) returns (stream ChatResponse);
```
- **Client Streaming**: کلاینت می‌تواند چندین پیام بفرستد
- **Server Streaming**: سرور می‌تواند چندین پیام بفرستد
- **Bidirectional Streaming**: هر دو طرف می‌توانند همزمان پیام بفرستند
- **Real-time Communication**: مناسب برای ارتباطات بلادرنگ

#### 4. Multi-language Support
- **Official Support**: پشتیبانی رسمی از Go، Java، Python، C++، Ruby، JavaScript، و ...
- **Polyglot Architecture**: سرویس‌ها می‌توانند با زبان‌های مختلف نوشته شوند
- **Consistent API**: یکسان بودن API در تمام زبان‌ها
- **Generated Code**: تولید خودکار کد برای تمام زبان‌های پشتیبانی شده

#### 5. ارتباطات کارآمد (Efficient Communication)
- **Connection Reuse**: استفاده از یک اتصال TCP برای چندین درخواست
- **Multiplexing**: ارسال همزمان چندین درخواست در یک اتصال
- **Header Compression**: فشرده‌سازی هدرها برای کاهش overhead
- **No JSON Overhead**: حذف overhead مربوط به JSON

#### 6. ابزارهای مدرن (Modern Tooling)
- **Protocol Buffers Compiler**: protoc برای تولید کد
- **gRPC Web**: پشتیبانی از مرورگرهای وب
- **Interceptors**: middleware برای logging, authentication, metrics
- **Plugins**: قابلیت افزودن افزونه‌های شخصی

### ⚠️ معایب و چالش‌های gRPC

#### 1. محدودیت در مرورگرها (Browser Support)
- **Limited Native Support**: مرورگرها به‌طور مستقیم از gRPC پشتیبانی نمی‌کنند
- **Requires gRPC-Web**: نیاز به گRPC-Web polyfill برای کلاینت‌های وب
- **Not Standard**: استاندارد وب نیست
- **Workaround Required**: نیاز به proxy یا polyfill برای وب

#### 2. دیباگ کردن سخت‌تر (Debugging Complexity)
- **Binary Format**: فرمت باینری قابل خواندن برای انسان نیست
- **Requires Tools**: نیاز به ابزارهای خاص برای نمایش پیام‌ها
- **JSON Advantage**: JSON قابل خواندن است اما gRPC نیست
- **Learning Curve**: یادگیری ابزارهای دیباگ مورد نیاز است

#### 3. منحنی یادگیری (Learning Curve)
- **Protocol Buffers Syntax**: نیاز به یادگیری سینتکس .proto
- **Toolchain**: ناشناخته بودن ابزارها برای توسعه‌دهندگان جدید
- **Code Generation**: درک فرآیند تولید کد
- **Different Paradigm**: تفاوت با RESTful APIها

#### 4. مدیریت خطا متفاوت (Error Handling)
- **No HTTP Status Codes**: استفاده از status codes گRPC به جای HTTP
- **Limited Error Context**: اطلاعات محدود در خطاها
- **Custom Status**: نیاز به تعریف statusهای سفارشی
- **Learning Curve**: یادگیری نحوه مدیریت خطاها در gRPC

#### 5. اکوسیستم کوچک‌تر (Smaller Ecosystem)
- **Fewer Tools**: ابزارهای کمتر نسبت به REST
- **Limited Frameworks**: فریمورک‌های کمتر برای توسعه سریع
- **Documentation**: مستندات کمتر نسبت به REST
- **Community**: جامعه کاربری کوچک‌تر

#### 6. چالش‌های Load Balancing
- **HTTP/2 Requirements**: نیاز به load balancerهای HTTP/2-aware
- **Connection State**: حفظ state اتصال‌ها چالش‌برانگیز است
- **Sticky Sessions**: گاهی نیاز به sticky sessions
- **Traditional LB Issues**: load balancerهای سنتی ممکن است خوب کار نکنند

#### 7. محدودیت Caching
- **No HTTP Caching**: پشتیبانی از cachingهای HTTP وجود ندارد
- **Manual Caching**: نیاز به caching سفارشی
- **CDN Limitations**: محدودیت در استفاده از CDNها
- **Cache Invalidation**: مدیریت invalidation سخت‌تر است

#### 8. پشتیبانی محدود Proxy
- **HTTP Proxies**: برخی پروکسی‌های HTTP از gRPC پشتیبانی نمی‌کنند
- **Firewall Issues**: ممکن است فایروال‌ها بلاک کنند
- **Network Restriction**: محدودیت‌های شبکه‌ای
- **Workaround Required**: نیاز به راهکارهای جایگزین

### 🎯 چه زمانی از gRPC استفاده کنیم؟

#### ✅ مناسب برای استفاده:

1. **میکروسرویس‌های داخلی (Internal Microservices)**
   - ارتباط بین سرویس‌های backend
   - محیط‌های کنترل شده
   - نیازمندی کارایی بالا

2. **ویژگی‌های بلادرنگ (Real-time Features)**
   - استریم داده
   - اطلاع‌رسانی‌های زنده
   - پیام‌رسانی بلادرنگ
   - سیستم‌های chat

3. **سناریوهای Throughput بالا (High Throughput)**
   - اپلیکیشن‌هایی با درخواست‌های همزمان زیاد
   - سیستم‌های با ترافیک سنگین
   - نیازمندی low latency

4. **قراردادهای دقیق (Strict Contracts)**
   - نیازمندی type safety
   - نیاز به schema evolution
   - پروژه‌های با APIهای پیچیده

5. **محیط‌های Polyglot (Polyglot Environments)**
   - سرویس‌های نوشته شده با زبان‌های مختلف
   - نیاز به یکنواختی API
   - تیم‌های مختلف با تخصص‌های متفاوت

6. **سیستم‌های با حجم داده بزرگ (Large Data Volume)**
   - کاهش اندازه پیام‌ها مهم است
   - نیاز به compression
   - باندویت محدود

#### ⚠️ بهتر است از REST استفاده کنیم:

1. **APIهای عمومی (Public APIs)**
   - در دسترس بودن برای توسعه‌دهندگان شخص ثالث
   - نیاز به استانداردهای وب
   - سادگی استفاده

2. **کلاینت‌های وب (Browser-based Clients)**
   - دسترسی مستقیم از مرورگرها
   - عدم نیاز به polyfill
   - سادگی integration

3. **عملیات ساده CRUD (Simple CRUD Operations)**
   - همنخوانی با معنایی HTTP
   - سادگی پیاده‌سازی
   - نیازمندی‌های ساده

4. **نیازمندی‌های Caching (Caching Requirements)**
   - استفاده از cachingهای HTTP
   - نیاز به CDN
   - نیازمندی‌های cache-aware

### 📝 پیاده‌سازی gRPC در این پروژه

#### فایل‌های Protocol Buffers
```protobuf
// proto/auth/auth.proto
syntax = "proto3";

package auth;

service AuthService {
    rpc Register(RegisterRequest) returns (RegisterResponse);
    rpc Login(LoginRequest) returns (LoginResponse);
    rpc ValidateToken(ValidateTokenRequest) returns (ValidateTokenResponse);
}

message RegisterRequest {
    string email = 1;
    string username = 2;
    string password = 3;
}

message RegisterResponse {
    bool success = 1;
    string message = 2;
    string user_id = 3;
}
```

#### ارتباط سرویس‌ها
- **API Gateway → Auth Service**: ارتباط با gRPC برای احراز هویت
- **Post Service → Auth Service**: اعتبارسنجی کاربر از طریق gRPC
- **Future Expansion**: افزودن سرویس‌های جدید با gRPC

#### پیکربندی شبکه
- **Internal Network**: gRPC در شبکه داخلی استفاده می‌شود
- **Port Configuration**: هر سرویس در پورت اختصاصی خود
- **Load Balancing**: استفاده از Kubernetes services برای balancing
- **TLS**: امکان فعال‌سازی TLS برای security

### 🔧 بهترین روش‌ها (Best Practices)

1. **مدل‌سازی داده‌ها (Data Modeling)**
   - استفاده از message types برای پیچیدگی
   - تعریف versioning برای APIها
   - استفاده از enums برای مقادیر ثابت

2. **مدیریت خطا (Error Handling)**
   - استفاده از status codes گRPC استاندارد
   - تعریف error messages واضح
   - استفاده از metadata برای context اضافه

3. **مدیریت اتصال (Connection Management)**
   - استفاده از connection pooling
   - مدیریت retry logic
   - handling timeoutها

4. **استقرار (Deployment)**
   - استفاده از Kubernetes services
   - پیکربندی health checks
   - monitoring و observability

5. **تست (Testing)**
   - نوشتن integration tests برای gRPC
   - استفاده از mock servers
   - testing error scenarios

---

### Infrastructure

#### Containerization
- **Docker**: برای containerization سرویس‌ها
- **Docker Compose**: برای مدیریت containerها

#### Orchestration
- **Kubernetes**: برای orchestration و scaling
- **Docker Desktop**: برای توسعه محلی با Kubernetes

#### CI/CD
- **GitHub Actions**: برای automation و deployment
- **ArgoCD**: برای GitOps deployment management

#### Development Tools
- **Tilt**: برای hot reload در توسعه محلی
- **Make**: برای automation دستورات build

---

## 🗄️ دیتابیس و مدل داده

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

## 🔄 جریان داده

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

## 🔒 امنیت

### Authentication Flow

1. Client → API Gateway: Credentials
2. API Gateway → Auth Service: Validate credentials
3. Auth Service → DB: Query user
4. Auth Service: Generate JWT token
5. Auth Service → Client: Return JWT token
6. Client → API Gateway: JWT token
7. API Gateway → Auth Service: Validate JWT
8. Auth Service → API Gateway: Token valid/invalid
9. API Gateway → Target Service: Forward request (with user context)

### Authorization Layers

1. **API Gateway Level**: Basic authentication check
2. **Service Level**: RBAC implementation
3. **Database Level**: Row-level security (if needed)

### Security Best Practices

#### Current Implementation
- ✅ JWT for authentication
- ✅ bcrypt for password hashing
- ✅ Role-based access control
- ✅ Database encryption at rest (Kubernetes Secrets)
- ✅ Network policies (can be enhanced)

#### Recommended Enhancements
- 🔒 OAuth 2.0 / OpenID Connect
- 🔒 Mutual TLS (mTLS) for service-to-service
- 🔒 Rate limiting and throttling
- 🔒 Input validation and sanitization
- 🔒 SQL injection prevention
- 🔒 CORS configuration

---

## 🚀 استقرار (Deployment)

### محیط‌های Deployment

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

### Kubernetes Resources

- **Deployments**: برای deployment و scaling سرویس‌ها
- **Services**: برای exposure و load balancing
- **ConfigMaps**: برای تنظیمات
- **Secrets**: برای اطلاعات حساس
- **PersistentVolumes**: برای دیتابیس
- **StatefulSets**: برای RabbitMQ

---

## 🔧 CI/CD Pipeline

### Workflow

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

### CI Stages

1. **Linting**: بررسی کد با linter
2. **Testing**: اجرای unit tests و integration tests
3. **Building**: ساختن Docker images
4. **Pushing**: pushing images به Docker registry
5. **Deploying**: deployment به Kubernetes cluster

---

## 📊 مزایا و چالش‌ها

### مزایا

#### Monorepo Benefits
- ✅ Shared code reuse
- ✅ Consistent tooling
- ✅ Atomic commits across services
- ✅ Easier dependency management
- ✅ Shared CI/CD pipelines

#### Microservices Benefits
- ✅ Scalability: هر سرویس می‌تواند جداگانه scale شود
- ✅ Maintainability: کد هر سرویس جداگانه نگهداری می‌شود
- ✅ Technology Flexibility: هر سرویس می‌تواند از تکنولوژی‌های مختلف استفاده کند
- ✅ Fault Isolation: Failure در یک سرویس بقیه را تحت تأثیر قرار نمی‌دهد
- ✅ Team Autonomy: تیم‌های مختلف می‌توانند روی سرویس‌های مختلف کار کنند

#### Event-Driven Benefits
- ✅ Decoupling: سرویس‌ها مستقل هستند
- ✅ Scalability: می‌تواند load را handle کند
- ✅ Reliability: message queue failure را handle می‌کند

### چالش‌ها

#### Complexity
- ⚠️ Distributed system: debugging و monitoring سخت‌تر است
- ⚠️ Network latency: بین سرویس‌ها وجود دارد
- ⚠️ Data consistency: در distributed systems سخت‌تر است

#### Operational Overhead
- ⚠️ Multiple services: مدیریت و monitoring بیشتر
- ⚠️ Deployment complexity: بیشتر از monolith
- ⚠️ Resource usage: بیشتر به دلیل overhead

#### Development
- ⚠️ Local development: سخت‌تر از monolith
- ⚠️ Testing: integration tests پیچیده‌تر
- ⚠️ Debugging: سخت‌تر به دلیل distributed nature

---

## 🎯 نتیجه‌گیری

این پروژه یک **modern, scalable, maintainable** microservices architecture است که:

1. ✅ از بهترین practiceها استفاده می‌کند
2. ✅ scalable و highly available است
3. ✅ separation of concerns را رعایت می‌کند
4. ✅ از event-driven patterns برای decoupling استفاده می‌کند
5. ✅ از Kubernetes برای orchestration استفاده می‌کند
6. ✅ از GitOps برای deployment management استفاده می‌کند
7. ✅ از CI/CD برای automation استفاده می‌کند

### Future Improvements
- 🔧 Add service mesh (Istio/Linkerd)
- 🔧 Implement caching layer (Redis)
- 🔧 Add rate limiting و API management
- 🔧 Improve monitoring و observability
- 🔧 Add API versioning
- 🔧 Implement circuit breakers و retries
- 🔧 Add comprehensive logging
- 🔧 Implement metrics collection
- 🔧 Add distributed tracing

---

**نسخه:** 1.0
**آخرین بروزرسانی:** 2025-09-01
**نویسنده:** Abbas Mortazavi
**لایسنس:** -