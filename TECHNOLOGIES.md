# پروژه Go Microservice - مستندسازی تکنولوژی‌ها

## 📋 معرفی پروژه

این پروژه یک معماری میکروسرویس با استفاده از زبان Go است که شامل سرویس‌های مختلف برای مدیریت احراز هویت، پست‌ها، نوتیفیکیشن‌ها و API Gateway می‌باشد.

---

## 🛠 تکنولوژی‌های استفاده شده

### 🔧 زبان برنامه‌نویسی

| تکنولوژی | نسخه | توضیحات |
|---------|------|---------|
| **Go** | 1.25.5 | زبان اصلی برنامه‌نویسی سرویس‌ها |

---

### 🏗 معماری سرویس‌ها

| سرویس | پورت | توضیحات |
|-------|------|---------|
| **API Gateway** | 8085 | دروازه ورودی و مدیریت روتینگ درخواست‌ها |
| **Auth Service** | 9092 | سرویس احراز هویت و مدیریت توکن‌ها |
| **Notification Service** | 9093 | سرویس ارسال نوتیفیکیشن |
| **Post Service** | 9094 | سرویس مدیریت پست‌ها |
| **Postgres Service** | 5432 | سرویس دیتابیس PostgreSQL |

---

### 📡 ارتباطات سرویس‌ها

#### gRPC & Protobuf
- **google.golang.org/grpc** v1.77.0 - فریم‌ورک gRPC برای ارتباطات بین سرویس‌ها
- **google.golang.org/protobuf** v1.36.11 - کامپایلر و کتابخانه Protobuf
- **github.com/golang/protobuf** v1.5.4 - پشتیبانی از Protobuf در Go

#### فایل‌های Proto موجود:
- `auth.proto` - تعریف سرویس احراز هویت
- `post.proto` - تعریف سرویس پست‌ها
- `events.proto` - تعریف رویدادها
- `permission.proto` - تعریف مجوزها
- `role.proto` - تعریف نقش‌ها
- `rbac.proto` - کنترل دسترسی مبتنی بر نقش

#### HTTP
- **github.com/gorilla/mux** v1.8.1 - Router قدرتمند برای REST API

---

### 📨 پیام‌رسانی

| تکنولوژی | نسخه | توضیحات |
|---------|------|---------|
| **RabbitMQ** | 3-management | Message Broker برای ارتباط غیرهمگام |
| **amqp091-go** | v1.10.0 | کلاینت RabbitMQ برای Go |
| پورت‌ها | 5672 (AMQP), 15672 (Management UI) | |

---

### 🗄 پایگاه داده

| تکنولوژی | نسخه | توضیحات |
|---------|------|---------|
| **PostgreSQL** | 15 | دیتابیس اصلی رابطه‌ای |
| **lib/pq** | v1.10.9 | Driver PostgreSQL برای Go |

#### مدیریت Migration
- **migrate** - ابزار مدیریت migration دیتابیس
  - دستورات: `migrate-up`, `migrate-down`, `migrate-create`, `migrate-status`, `migrate-force`

---

### 🔐 امنیت و احراز هویت

| کتابخانه | نسخه | کاربرد |
|----------|------|-------|
| **JWT** | v5.3.0 | تولید و اعتبارسنجی توکن‌های JSON Web Token |
| **Crypto (bcrypt)** | v0.46.0 | هش کردن رمز عبور با الگوریتم bcrypt |
| **UUID** | v1.6.0 | تولید شناسه‌های منحصر به فرد |

---

### ✅ اعتبارسنجی

| کتابخانه | نسخه | کاربرد |
|----------|------|-------|
| **Validator** | v10.29.0 | اعتبارسنجی ساختارها و داده‌ها |

---

### 🐳 کانتینر و اورکستریشن

#### Docker
- استفاده برای کانتینر کردن سرویس‌ها
- Dockerfile اختصاصی برای هر سرویس

#### Kubernetes
- اورکستریشن سرویس‌ها در محیط‌های مختلف:
  - **Local** - محیط توسعه محلی
  - **Development** - محیط توسعه
  - **Production** - محیط تولید

#### منابع Kubernetes:
- Deployments
- Services
- StatefulSets (برای RabbitMQ)
- PersistentVolumeClaims (برای ذخیره‌سازی دیتابیس)
- ConfigMaps
- Secrets

---

### 🎯 GitOps و CI/CD

#### ArgoCD
- مدیریت دیپلویoment با استفاده از GitOps
- Application manifests برای محیط‌های مختلف:
  - `infra/production/argocd/application.yaml`
  - `infra/development/argocd/application.yaml`

#### GitHub Actions
- Workflows برای اتوماسیون:
  - `run-tests.yml` - اجرای تست‌ها
  - `run-build.yml` - بیلد پروژه
  - `deploy.yml` - دیپلوی عمومی
  - `deploy-dev.yml` - دیپلوی به محیط توسعه
  - `deploy-prod.yml` - دیپلوی به محیط تولید (نیاز به تایید)

---

### 🛠 ابزارهای توسعه

| ابزار | کاربرد |
|------|-------|
| **Tilt** | توسعه سریع و اتوماتیک در محیط Kubernetes |
| **Make** | اتوماسیون وظایف (generate-proto, migrations, etc.) |
| **protoc** | کامپایل فایل‌های .proto به کد Go |
| **Go test** | فریم‌ورک تست داخلی Go |
| **testify** | v1.10.0 - ابزارهای کمکی برای نوشتن تست |

---

### 📁 ساختار پروژه

```
go-microservice/
├── services/              # سرویس‌های میکروسرویس
│   ├── api-gateway/      # دروازه API
│   ├── auth-service/     # سرویس احراز هویت
│   ├── notification-service/ # سرویس نوتیفیکیشن
│   └── post-service/     # سرویس پست‌ها
├── proto/                # فایل‌های Protobuf
├── pkg/                  # پکیج‌های مشترک
├── migrations/           # فایل‌های migration دیتابیس
├── infra/                # کانفیگ‌های زیرساخت
│   ├── local/            # محیط محلی
│   ├── development/      # محیط توسعه
│   └── production/       # محیط تولید
├── build/                # فایل‌های بیلد شده
├── tools/                # ابزارهای اضافی
├── argocd/               # کانفیگ‌های ArgoCD
└── .github/              # Workflows گیت‌هاب
    └── workflows/        # فایل‌های CI/CD
```

---

### 🔧 دستورات Make مهم

| دستور | توضیحات |
|-------|---------|
| `make generate-proto` | کامپایل فایل‌های proto |
| `make migrate-up` | اجرای migrationها |
| `make migrate-down` | بازگشت به نسخه قبلی |
| `make migrate-create` | ساخت migration جدید |
| `make migrate-status` | نمایش وضعیت migrationها |
| `make migrate-force` | فورس کردن نسخه migration |

---

### 🌐 پورت‌های سرویس‌ها

| سرویس | پورت | توضیحات |
|-------|------|---------|
| API Gateway | 8085 | دروازه ورودی |
| Auth Service | 9092 | سرویس احراز هویت |
| Notification Service | 9093 | سرویس نوتیفیکیشن |
| Post Service | 9094 | سرویس پست‌ها |
| PostgreSQL | 5432 | دیتابیس |
| RabbitMQ (AMQP) | 5672 | پیام‌رسانی |
| RabbitMQ (Management) | 15672 | پنل مدیریتی |

---

### 📦 وابستگی‌های اصلی

```go
require (
    github.com/go-playground/validator/v10 v10.29.0
    github.com/golang-jwt/jwt/v5 v5.3.0
    github.com/golang/protobuf v1.5.4
    github.com/google/uuid v1.6.0
    github.com/gorilla/mux v1.8.1
    github.com/lib/pq v1.10.9
    github.com/rabbitmq/amqp091-go v1.10.0
    golang.org/x/crypto v0.46.0
    google.golang.org/grpc v1.77.0
    google.golang.org/protobuf v1.36.11
)
```

---

### 🎨 ویژگی‌های کلیدی

- ✅ معماری میکروسرویس با decoupled components
- ✅ ارتباط gRPC برای performance بالا
- ✅ RabbitMQ برای async messaging
- ✅ PostgreSQL به عنوان دیتابیس اصلی
- ✅ JWT برای احراز هویت
- ✅ Kubernetes برای اورکستریشن
- ✅ ArgoCD برای GitOps
- ✅ CI/CD با GitHub Actions
- ✅ Tilt برای توسعه سریع
- ✅ Database migrations با ابزار migrate
- ✅ Health checks و probes در Kubernetes
- ✅ Resource limiting در Kubernetes

---

### 📝 نکات مهم

1. **محیط‌های مختلف**: پروژه از سه محیط پشتیبانی می‌کند (local, development, production)
2. **GitOps**: استفاده از ArgoCD برای مدیریت دیپلویoment‌ها
3. **CI/CD**: اتوماسیون کامل تست، بیلد و دیپلوی
4. **High Availability**: استفاده از StatefulSets برای RabbitMQ و PVC برای PostgreSQL
5. **Resource Management**: تعیین منابع CPU و Memory برای هر سرویس

---

**نسخه مستندات:** 1.0
**آخرین بروزرسانی:** 2025-08-26
