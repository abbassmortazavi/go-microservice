# 🔍 مشکل کندی لاگین - تحلیل و راهکار

## ⚠️ مشکلات شناسایی شده

### مشکل 1: تعداد زیاد درخواست‌های دیتابیس در هر لاگین

در تابع `Login` از `auth_service.go`:
```go
func (a *AuthService) Login(ctx context.Context, email, password string) (*response.TokenResponseResult, error) {
    // 1. دیتابیس query برای پیدا کردن کاربر
    user, err := a.userRepo.FindByEmail(ctx, email)
    
    // 2. مقایسه رمز عبور (بسته به bcrypt cost می‌تواند کند باشد)
    if err := a.hasher.Compare(user.Password, password); err == false {
        return nil, errors.New("password is wrong")
    }
    
    // 3. تابع GenerateToken صدا زده می‌شود که 4 درخواست دیتابیس دیگر انجام می‌دهد
    tokens, err := a.TokenService.GenerateToken(user.ID, user.Name)
}
```

### مشکل 2: تابع GenerateToken خیلی کند است

در تابع `GenerateToken` از `token_service.go`:
```go
func (j *JWT) GenerateToken(userID int64, name string) (response.TokenResponse, error) {
    // 1. دیتابیس query برای پیدا کردن کاربر (غیرضروری!)
    user, err := j.FindByUserId(ctx, userID)
    
    // 2. دیتابیس query برای پیدا کردن توکن‌های قبلی کاربر
    res, err := j.TokenRepository.FindByUserId(ctx, userID)
    
    // 3. دیتابیس DELETE برای حذف تمام توکن‌های کاربر
    err := j.TokenRepository.RevokeAllUserTokens(ctx, userID)
    
    // 4. دیتابیس INSERT برای ایجاد access token
    err = j.TokenRepository.Create(ctx, &reqAccessToken)
    
    // 5. دیتابیس INSERT برای ایجاد refresh token
    err = j.TokenRepository.Create(ctx, &reqRefreshToken)
}
```

**مجموعاً:** هر لاگین **6 درخواست دیتابیس** انجام می‌دهد! 🐢

### مشکل 3: Bug در تابع Create token_repository.go

در خط 116 از `token_repository.go`:
```go
row := r.db.QueryRowContext(ctx, query, token.UserID, token.TokenType, token.HashToken, token.CreatedAt, token.IsRevoked)
//                                                                                     ^^^^^^^^^^^^^^^^
//                                                                                     باید token.ExpiredAt باشد
```

---

## 🐌 مسیر کامل یک درخواست لاگین

```
Client Request
    ↓
[1] API Gateway (نامشخص زمان)
    ↓
[2] Auth Service: Login() ← userRepo.FindByEmail() [DB Query 1]
    ↓
[3] Bcrypt Password Compare [CPU intensive]
    ↓
[4] GenerateToken() ← FindByUserId() [DB Query 2] ⚠️ غیرضروری!
    ↓
[5] FindByUserId() ← TokenRepository.FindByUserId() [DB Query 3]
    ↓
[6] RevokeAllUserTokens() [DB Query 4] ⚠️ می‌تواند کند باشد
    ↓
[7] TokenRepository.Create() [DB Query 5] ⚠️ Bug!
    ↓
[8] TokenRepository.Create() [DB Query 6]
    ↓
Response (1-10 ثانیه)
```

---

## ⚡ راه‌حل‌ها

### راه‌حل 1: حذف دیتابیس query غیرضروری در GenerateToken

**مسئله:** اطلاعات کاربر قبلاً در Login به دست آمده است، ولی دوباره در GenerateToken دیتابیس query می‌شود.

**راه‌حل:** اطلاعات کاربر را به عنوان پارامتر به GenerateToken پاس بدهیم:

```go
// قبل
func (j *JWT) GenerateToken(userID int64, name string) (response.TokenResponse, error)

// بعد
func (j *JWT) GenerateToken(user *entity.User) (response.TokenResponse, error)
```

**مزیت:** حذف 1 درخواست دیتابیس ⚡

### راه‌حل 2: بهینه‌سازی حذف توکن‌های قبلی

**مسئله:** در هر لاگین، تمام توکن‌های قبلی کاربر حذف می‌شود که می‌تواند کند باشد.

**راه‌حل 1:** از soft delete استفاده کنید:
```sql
UPDATE tokens SET is_revoked = true WHERE user_id = $1 AND is_revoked = false
```

**راه‌حل 2:** به جای حذف، فقط tokenهای منقضی شده را حذف کنید:
```sql
DELETE FROM tokens WHERE user_id = $1 AND expired_at < NOW()
```

**راه‌حل 3 (بهترین):** حذف را به صورت asynchronous انجام دهید:
```go
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    j.TokenRepository.RevokeAllUserTokens(ctx, userID)
}()
```

**مزیت:** حذف شدن از مسیر اصلی درخواست ⚡⚡⚡

### راه‌حل 3: رفع Bug در Create

```go
// قبل (خط 116)
row := r.db.QueryRowContext(ctx, query, token.UserID, token.TokenType, token.HashToken, token.CreatedAt, token.IsRevoked)

// بعد
row := r.db.QueryRowContext(ctx, query, token.UserID, token.TokenType, token.HashToken, token.ExpiredAt, token.IsRevoked)
```

### راه‌حل 4: Connection Pool Optimization

به نظر می‌رسد که connection pool به درستی تنظیم نشده است. به دیتابیس connection نگاه کنید:

```go
// تنظیم connection pool
db.SetMaxOpenConns(25)
db.SetMaxIdleConns(10)
db.SetConnMaxLifetime(5 * time.Minute)
```

### راه‌حل 5: Database Indexes

اطمینان حاصل کنید که این indexها روی دیتابیس وجود دارند:

```sql
-- برای سرعت‌بخش به queryهای لاگین
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id ON tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_tokens_user_id_type ON tokens(user_id, token_type);
CREATE INDEX IF NOT EXISTS idx_tokens_hash_token ON tokens(hash_token);

-- برای RevokeAllUserTokens
CREATE INDEX IF NOT EXISTS idx_tokens_user_id_revoked ON tokens(user_id, is_revoked);
```

### راه‌حل 6: بهینه‌سازی Bcrypt Cost

bcrypt cost در تنظیمات خودتون رو چک کنید:

```go
// اگر cost خیلی بالاست، می‌تواند کند باشد
// معمولاً 10-12 کافی است
cost := 10
hasher := security.NewBcryptPasswordHasher(cost)
```

---

## 🎯 پیشنهاد پیاده‌سازی سریع

### گام 1: رفع Bug فوری
```go
// در token_repository.go خط 116
row := r.db.QueryRowContext(ctx, query, token.UserID, token.TokenType, token.HashToken, token.ExpiredAt, token.IsRevoked)
```

### گام 2: بهینه‌سازی GenerateToken
```go
// در token_service.go خط 47
func (j *JWT) GenerateToken(user *entity.User) (response.TokenResponse, error) {
    // حذف دیتابیس query غیرضروری
    // user := j.FindByUserId(ctx, userID) ← حذف این خط
    
    // استفاده از کاربر از پارامتر
    userInfo := entity.User{
        ID:    user.ID,
        Name:  user.Name,
        Email: user.Email,
        Role:  user.Role,
    }
    // ... ادامه کد
}
```

### گام 3: Async token cleanup
```go
// در token_service.go خط 105
// قبل
err := j.TokenRepository.RevokeAllUserTokens(ctx, userID)

// بعد
go func() {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    j.TokenRepository.RevokeAllUserTokens(ctx, userID)
}()
```

### گام 4: به‌روزرسانی فراخوانی در auth_service.go
```go
// در auth_service.go خط 74
// قبل
tokens, err := a.TokenService.GenerateToken(user.ID, user.Name)

// بعد
tokens, err := a.TokenService.GenerateToken(&userEntity)
```

---

## 📊 نتایج مورد انتظار

**قبل از بهینه‌سازی:**
- تعداد queryهای دیتابیس: 6
- زمان لاگین: 5-10 ثانیه

**بعد از بهینه‌سازی:**
- تعداد queryهای دیتابیس: 3-4
- زمان لاگین: 0.5-2 ثانیه

**بعد از async cleanup:**
- تعداد queryهای دیتابیس در مسیر اصلی: 2-3
- زمان لاگین: 0.2-1 ثانیه ⚡⚡⚡

---

## 🔍 ابزارها برای دیباگ

### 1. Enable Query Logging
```go
// در init دیتابیس
db.SetLogger(log.New(os.Stdout, "DB: ", log.LstdFlags))
db.SetMaxOpenConns(1) // برای دیدن تمام queryها
```

### 2. Add Timing to Functions
```go
func (a *AuthService) Login(ctx context.Context, email, password string) (*response.TokenResponseResult, error) {
    start := time.Now()
    defer log.Printf("Login took: %v", time.Since(start))
    
    user, err := a.userRepo.FindByEmail(ctx, email)
    if err != nil {
        return nil, err
    }
    log.Printf("FindByEmail took: %v", time.Since(start))
    
    // ... ادامه کد
}
```

### 3. Database Profiling
```sql
-- اجرا در PostgreSQL
EXPLAIN ANALYZE SELECT * FROM users WHERE email = 'test@example.com';
EXPLAIN ANALYZE DELETE FROM tokens WHERE user_id = 1;
```

---

## 🚀 اولویت‌بندی فوری

1. **فوری:** رفع Bug در Create (خط 116) ⚠️
2. **خیلی فوری:** Async token cleanup (خط 105) ⚠️⚠️
3. **فوری:** حذف دیتابیس query غیرضروری ⚠️
4. **مهم:** Database indexes
5. **مهم:** Connection pool optimization

---

**وضعیت:** ⚠️ **کریتیکال** - مشکلات چندگانه که به طور جدی عملکرد را تحت تأثیر قرار می‌دهند
