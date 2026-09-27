# 🎬 Cinema Ticket Booking System

ระบบจองตั๋วภาพยนตร์แบบ Full-Stack สำหรับจัดการภาพยนตร์ โรงภาพยนตร์ โรง/Screen
ที่นั่ง รอบฉาย การจอง การชำระเงิน และระบบผู้ดูแล (Admin)

------------------------------------------------------------------------

## 📌 ภาพรวมระบบ

โปรเจกต์แบ่งออกเป็น 2 ส่วนหลัก

``` text
Cinema-Ticket-Booking-System/
├── Backend/
└── Frontend/
```

### Backend

พัฒนาโดยใช้:

-   Go 1.26.5
-   Gin
-   GORM
-   PostgreSQL
-   Redis
-   JWT
-   bcrypt
-   Google OAuth 2.0
-   WebSocket
-   Redis Pub/Sub
-   Swagger / OpenAPI

Backend ทำหน้าที่เป็น REST API สำหรับ:

-   Authentication
-   User management
-   Movie management
-   Theater management
-   Screen management
-   Seat management
-   Showtime management
-   Seat locking
-   Seat map แบบ Real-time
-   Booking
-   Payment
-   Admin
-   Audit Log
-   Password reset

### Frontend

พัฒนาโดยใช้:

-   Vue 3
-   TypeScript
-   Vite
-   Vue Router
-   Pinia
-   Axios
-   Tailwind CSS
-   Lucide Icons

Frontend แบ่งหน้าเป็น 2 กลุ่ม:

``` text
User
Admin
```

------------------------------------------------------------------------

# 🧰 Technology Stack

  ส่วน                        Technology
  -------------------------- --------------------------
  Frontend                   Vue 3
  Language                   TypeScript
  Build Tool                 Vite
  Styling                    Tailwind CSS
  State Management           Pinia
  HTTP Client                Axios
  Icons                      Lucide Vue
  Backend                    Go
  Web Framework              Gin
  ORM                        GORM
  Database                   PostgreSQL
  Cache / Distributed Lock   Redis
  Authentication             JWT
  Password Hashing           bcrypt
  Social Login               Google OAuth 2.0
  Real-time                  WebSocket
  Message Queue / Event      Redis Pub/Sub
  API Documentation          Swagger
  Container support          Dockerfile มีอยู่ใน Backend

> ใน source code ปัจจุบันไม่มี `docker-compose.yml` และ `cmd/main.go` มี
> comment ระบุว่าแนวทางการรันปัจจุบันใช้ PostgreSQL/Redis แบบ local โดยไม่ได้พึ่ง
> Docker Compose

------------------------------------------------------------------------

# 🏗️ System Architecture

``` text
                 ┌─────────────────────┐
                 │      Vue 3 App      │
                 │   TypeScript/Vite   │
                 └──────────┬──────────┘
                            │
                         Axios
                            │
                            ▼
                 ┌─────────────────────┐
                 │      Go / Gin       │
                 │      REST API       │
                 └──────┬─────┬────────┘
                        │     │
              ┌─────────┘     └─────────────┐
              ▼                             ▼
     ┌─────────────────┐           ┌─────────────────┐
     │   PostgreSQL    │           │      Redis      │
     │      GORM       │           │ Seat Lock /     │
     │                 │           │ Pub/Sub         │
     └─────────────────┘           └────────┬────────┘
                                            │
                                            ▼
                                  ┌─────────────────┐
                                  │    WebSocket    │
                                  │  Real-time Seat │
                                  │      Map        │
                                  └─────────────────┘
```

------------------------------------------------------------------------

# 📂 Backend Structure

``` text
Backend/
├── cmd/
│   └── main.go
│
├── config/
│   ├── database.go
│   ├── oauth.go
│   └── redis.go
│
├── controllers/
│   ├── admin_controller.go
│   ├── audit_log_controller.go
│   ├── auth_controller.go
│   ├── booking_controller.go
│   ├── movie_controller.go
│   ├── payment_controller.go
│   ├── screen_controller.go
│   ├── seat_controller.go
│   ├── seat_lock_controller.go
│   ├── seat_map_controller.go
│   ├── showtime_controller.go
│   ├── theater_controller.go
│   └── user_controller.go
│
├── middleware/
│   └── auth_middleware.go
│
├── models/
│   ├── audit_log.go
│   ├── booking.go
│   ├── booking_seat.go
│   ├── movie.go
│   ├── password_reset_token.go
│   ├── payment.go
│   ├── screen.go
│   ├── seat.go
│   ├── showtime.go
│   ├── theater.go
│   └── user.go
│
├── repository/
│   ├── admin_repository.go
│   ├── audit_log_repository.go
│   ├── booking_repository.go
│   ├── movie_repository.go
│   ├── password_reset_repository.go
│   ├── payment_repository.go
│   ├── screen_repository.go
│   ├── seat_repository.go
│   ├── showtime_repository.go
│   ├── theater_repository.go
│   └── user_repository.go
│
├── routes/
│   ├── admin_routes.go
│   ├── audit_log_routes.go
│   ├── auth_routes.go
│   ├── booking_routes.go
│   ├── movie_routes.go
│   ├── payment_routes.go
│   ├── screen_routes.go
│   ├── seat_lock_routes.go
│   ├── seat_map_routes.go
│   ├── seat_routes.go
│   ├── showtime_routes.go
│   ├── theater_routes.go
│   └── user_routes.go
│
├── services/
│   ├── admin_service.go
│   ├── audit_log_service.go
│   ├── auth_service.go
│   ├── booking_service.go
│   ├── message_queue_service.go
│   ├── movie_service.go
│   ├── payment_service.go
│   ├── screen_service.go
│   ├── seat_lock_service.go
│   ├── seat_map_service.go
│   ├── seat_service.go
│   ├── seat_status_service.go
│   ├── showtime_service.go
│   └── user_service.go
│
├── docs/
│   ├── docs.go
│   ├── swagger.json
│   └── swagger.yaml
│
├── uploads/
├── Dockerfile
├── go.mod
└── MIGRATION_NOTES.md
```

------------------------------------------------------------------------

# 📂 Frontend Structure

``` text
Frontend/
├── src/
│   ├── assets/
│   │   └── styles/
│   │       └── main.css
│   │
│   ├── components/
│   │   ├── admin/
│   │   │   ├── AdminSidebar.vue
│   │   │   └── AdminTopbar.vue
│   │   │
│   │   ├── common/
│   │   │   ├── AppBadge.vue
│   │   │   ├── AppButton.vue
│   │   │   ├── AppInput.vue
│   │   │   ├── AppLoading.vue
│   │   │   └── AppModal.vue
│   │   │
│   │   ├── movie/
│   │   │   └── MovieCard.vue
│   │   │
│   │   └── user/
│   │       ├── Footer.vue
│   │       └── Navbar.vue
│   │
│   ├── layouts/
│   │   ├── AdminLayout.vue
│   │   └── UserLayout.vue
│   │
│   ├── router/
│   │   ├── index.ts
│   │   └── route-meta.d.ts
│   │
│   ├── services/
│   │   ├── api.ts
│   │   └── movieService.ts
│   │
│   ├── stores/
│   │   └── auth.ts
│   │
│   ├── types/
│   │   └── movie.ts
│   │
│   ├── views/
│   │   ├── admin/
│   │   └── user/
│   │
│   ├── App.vue
│   └── main.ts
│
├── index.html
├── package.json
├── tailwind.config.ts
├── vite.config.ts
└── tsconfig*.json
```

------------------------------------------------------------------------

# 👤 User Features

Frontend มีหน้าหลักสำหรับผู้ใช้งาน ได้แก่:

-   Home
-   Movies
-   Movie Detail
-   Showtime Selection
-   Seat Selection
-   Booking Summary
-   Payment
-   Booking Success
-   My Bookings
-   Booking Detail
-   Profile
-   Login
-   Register
-   Forgot Password
-   Reset Password
-   Google OAuth Callback
-   Help Center
-   Refund Policy
-   Contact Us
-   Coming Soon

------------------------------------------------------------------------

# 🎬 Movie Flow

Flow หลักของการจองภาพยนตร์:

``` text
Home
  │
  ▼
Movies
  │
  ▼
Movie Detail
  │
  ▼
Showtime Selection
  │
  ▼
Seat Selection
  │
  ▼
Booking Summary
  │
  ▼
Payment
  │
  ▼
Booking Success
```

------------------------------------------------------------------------

# 🕐 Showtime Selection

ระบบใช้ข้อมูลจาก:

``` text
GET /movies
GET /theaters
GET /screens
GET /showtimes
```

Frontend จะนำข้อมูลมาเชื่อมกันตาม:

``` text
movie_id
theater_id
screen_id
showtime_id
```

และแสดง:

-   Movie
-   Theater
-   Screen
-   วันที่ฉาย
-   เวลาเริ่ม
-   เวลาจบ
-   ราคาที่นั่ง Standard
-   ราคาที่นั่ง Premium
-   ราคาที่นั่ง VIP
-   สถานะรอบฉาย

------------------------------------------------------------------------

# 💺 Seat Selection

ระบบรองรับ:

-   Seat Map
-   Seat Number
-   Seat Type
-   Standard
-   Premium
-   VIP
-   Available
-   Locked
-   Booked
-   Disabled
-   การเลือกหลายที่นั่ง
-   Real-time seat status

โครงสร้างที่นั่งถูกผูกกับ Screen:

``` text
Theater
   │
   └── Screen
          │
          └── Seats
```

ตัวอย่าง:

``` text
A1 A2 A3 A4 A5 A6 A7 A8 ...
B1 B2 B3 B4 B5 B6 B7 B8 ...
C1 C2 C3 C4 C5 C6 C7 C8 ...
```

------------------------------------------------------------------------

# 🔐 Distributed Seat Lock

ระบบใช้ Redis เพื่อป้องกันผู้ใช้หลายคนเลือกที่นั่งเดียวกันพร้อมกัน

Key ที่ใช้:

``` text
seat_lock:<showtime_id>:<seat_id>
```

การ Lock ใช้:

``` text
Redis SETNX
```

ค่า TTL ปัจจุบันใน `SeatLockService`:

``` text
15 minutes
```

เมื่อ User เลือก Seat:

``` text
User A
   │
   ▼
Lock Seat
   │
   ▼
Redis SETNX
   │
   ├── สำเร็จ → Seat Locked
   │
   └── ไม่สำเร็จ → Seat ถูก Lock อยู่แล้ว
```

ถ้า Seat ถูกล็อกอยู่ ระบบจะส่ง error:

``` text
seat <seat_number> is already locked
```

------------------------------------------------------------------------

# 🔄 Real-time Seat Map

Backend มี WebSocket สำหรับ Seat Map

Endpoint:

``` text
GET /showtimes/:showtime_id/seat-map
```

ระบบใช้:

``` text
Gorilla WebSocket
```

ผู้ใช้ที่กำลังดู Showtime เดียวกันสามารถรับการเปลี่ยนแปลงของ Seat ได้

ตัวอย่าง event:

``` json
{
  "type": "seat_status",
  "showtime_id": "...",
  "seat_id": "...",
  "status": "LOCKED",
  "user_id": "...",
  "message": "Seat locked"
}
```

เมื่อปลด Lock:

``` json
{
  "type": "seat_status",
  "status": "AVAILABLE"
}
```

------------------------------------------------------------------------

# 🧾 Booking

Booking รองรับ:

-   User ID
-   Showtime ID
-   Seat Numbers
-   Total Price
-   Booking Status
-   Created At
-   Updated At

สถานะหลัก:

``` text
pending
confirmed
cancelled
```

## Create Booking Flow

``` text
Select Seats
     │
     ▼
POST /bookings
     │
     ▼
ตรวจสอบ User
     │
     ▼
ตรวจสอบ Showtime
     │
     ▼
Lock Seats ด้วย Redis
     │
     ▼
สร้าง Booking
     │
     ▼
Booking = pending
```

ถ้า Lock ที่นั่งบางตัวสำเร็จแล้ว แต่ตัวถัดไป Lock ไม่ได้ ระบบจะปลด Lock
ที่ทำสำเร็จไปก่อนหน้าออก

------------------------------------------------------------------------

# 💳 Payment

Payment รองรับ:

``` text
cash
credit_card
promptpay
transfer
```

สถานะ:

``` text
pending
paid
failed
refunded
```

เมื่อ Payment มีสถานะ:

``` text
paid
```

ระบบจะเรียก:

``` text
BookingService.ConfirmBooking()
```

จากนั้น:

``` text
Payment = paid
       │
       ▼
ตรวจสอบ Booking
       │
       ▼
ตรวจสอบ Redis Seat Lock
       │
       ▼
Seat = BOOKED
       │
       ▼
Booking = confirmed
       │
       ▼
Redis Lock ถูกลบ
```

------------------------------------------------------------------------

# 📢 Redis Pub/Sub

Backend มี Message Queue Service โดยใช้ Redis Pub/Sub

Channel:

``` text
booking.success
```

Event:

``` json
{
  "event": "booking.success",
  "booking_id": "...",
  "user_id": "...",
  "showtime_id": "...",
  "seat_numbers": []
}
```

Consumer จะ subscribe channel นี้ตั้งแต่ตอน Backend startup

ปัจจุบัน consumer มี mock notification:

``` text
Booking success processed
```

------------------------------------------------------------------------

# 🔑 Authentication

ระบบรองรับ Authentication แบบ JWT

## Register

``` http
POST /auth/register
```

รับ:

``` json
{
  "name": "John",
  "email": "john@example.com",
  "password": "password123"
}
```

User ที่สมัครผ่าน endpoint นี้จะถูกสร้างด้วย:

``` text
role = user
```

## Login

``` http
POST /auth/login
```

สำเร็จแล้ว Backend จะส่ง JWT กลับมา

JWT claims ประกอบด้วย:

``` text
user_id
role
exp
iat
```

JWT ใช้:

``` text
HS256
```

อายุ Token:

``` text
24 hours
```

------------------------------------------------------------------------

# 🔒 Password Security

Password ถูก hash ด้วย:

``` text
bcrypt
```

ระบบมี:

``` text
HashPassword()
CheckPassword()
```

Password จะไม่ถูกส่งกลับผ่าน JSON response

------------------------------------------------------------------------

# 🔁 Forgot / Reset Password

มี Endpoint:

``` http
POST /auth/forgot-password
POST /auth/reset-password
```

Reset token มีอายุ:

``` text
30 minutes
```

Token ถูกเก็บใน:

``` text
password_reset_tokens
```

ใน source ปัจจุบัน `forgot-password` ส่ง `reset_token` กลับใน response
เพื่อรองรับ development/testing และมี comment ระบุว่าในอนาคตควรส่งผ่าน email

------------------------------------------------------------------------

# 🌐 Google OAuth 2.0

Source code ปัจจุบันยังมี Google OAuth implementation

Endpoints:

``` http
GET /auth/google
GET /auth/google/callback
```

Frontend มี:

``` text
OAuthCallbackView.vue
```

Backend มีการสร้าง OAuth state cookie เพื่อใช้ตรวจสอบ state ตอน callback

Environment variables ที่เกี่ยวข้อง:

``` env
GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=
```

------------------------------------------------------------------------

# 👑 Admin System

Admin routes ใช้:

``` text
JWT Authentication
+
Role = admin
```

Middleware:

``` text
AuthMiddleware()
RequireRole("admin")
```

Admin frontend มี:

-   Dashboard
-   Profile
-   Movies
-   Theaters
-   Screens
-   Seats
-   Showtimes
-   Bookings
-   Payments
-   Audit Logs

------------------------------------------------------------------------

# 🏢 Admin Management

## Movies

Admin สามารถจัดการ:

``` text
Create
Read
Update
Delete
```

ข้อมูล Movie มี:

-   Title
-   Description
-   Genre
-   Duration
-   Language
-   Subtitle
-   Rating
-   Poster URL
-   Trailer URL
-   Release Date
-   Status

Status:

``` text
NOW_SHOWING
COMING_SOON
ENDED
```

------------------------------------------------------------------------

# 🏬 Theater

Theater มีข้อมูล:

-   Name
-   Address
-   Status
-   Location
-   Capacity

Status:

``` text
active
inactive
```

CRUD:

``` text
GET
POST
PUT
DELETE
```

------------------------------------------------------------------------

# 🎥 Screen

Screen อยู่ภายใต้ Theater

ข้อมูล:

-   Screen ID
-   Theater ID
-   Name
-   Total Seats

Relationship:

``` text
Theater
  └── Screen
```

------------------------------------------------------------------------

# 💺 Admin Seat Management

Admin สามารถ:

-   ดู Seat Map
-   เพิ่ม Seat
-   แก้ไข Seat
-   ลบ Seat
-   Generate Seat หลายตัว
-   กำหนด Seat Type

Seat Type:

``` text
standard
premium
vip
```

Seat Status:

``` text
AVAILABLE
```

และสถานะอื่น ๆ ที่ใช้ในระบบ Seat Map/Booking เช่น:

``` text
LOCKED
BOOKED
DISABLED
```

------------------------------------------------------------------------

# 🎞️ Showtime Management

Showtime มี:

-   Movie
-   Theater
-   Screen
-   Show Date
-   Start Time
-   End Time
-   Standard Price
-   Premium Price
-   VIP Price
-   Status

Status:

``` text
scheduled
showing
finished
cancelled
```

ราคาที่นั่งถูกกำหนดตาม Seat Type

``` text
Standard → PriceStandard
Premium  → PricePremium
VIP      → PriceVIP
```

------------------------------------------------------------------------

# 🧑‍💼 Admin Booking Management

Admin มีหน้า:

``` text
/admin/bookings
```

สำหรับดูข้อมูล Booking

และ Backend มี:

``` http
GET /admin/bookings
```

Route นี้อยู่ภายใต้:

``` text
AuthMiddleware
RequireRole("admin")
```

------------------------------------------------------------------------

# 💰 Admin Payment Management

Frontend มี:

``` text
/admin/payments
```

Backend:

``` http
GET    /payments
POST   /payments
PUT    /payments/:payment_id
DELETE /payments/:payment_id
```

------------------------------------------------------------------------

# 📝 Audit Logs

ระบบมี Audit Log model และ API

``` http
GET    /audit_logs
POST   /audit_logs
PUT    /audit_logs/:audit_log_id
DELETE /audit_logs/:audit_log_id
```

ข้อมูล Audit Log:

-   User ID
-   Action
-   Resource
-   Resource ID
-   Details
-   Created At

------------------------------------------------------------------------

# 🗄️ Database

Database ปัจจุบันคือ:

``` text
PostgreSQL
```

ใช้:

``` text
GORM
```

และ `AutoMigrate()` ตอน Backend startup

Models ที่ migrate:

``` text
User
Movie
Theater
Screen
Seat
Showtime
Booking
BookingSeat
Payment
AuditLog
PasswordResetToken
```

------------------------------------------------------------------------

# 🔗 Database Relationships

``` text
User
 │
 └── Booking
       │
       ├── Showtime
       │      ├── Movie
       │      ├── Theater
       │      └── Screen
       │             └── Seat
       │
       ├── BookingSeat
       │       └── Seat
       │
       └── Payment
```

รายละเอียด:

``` text
Theater 1 ─── N Screen

Screen 1 ─── N Seat

Movie 1 ─── N Showtime

Theater 1 ─── N Showtime

Screen 1 ─── N Showtime

User 1 ─── N Booking

Showtime 1 ─── N Booking

Booking 1 ─── N BookingSeat

Seat 1 ─── N BookingSeat

Booking 1 ─── N Payment
```

------------------------------------------------------------------------

# 🆔 Database ID Strategy

Models ใช้ Primary Key ภายในเป็น:

``` text
ID uint
```

แต่ API ใช้ Business UUID เช่น:

``` text
movie_id
theater_id
screen_id
seat_id
showtime_id
booking_id
payment_id
user_id
```

ตัวอย่าง:

``` text
ID        → 1
MovieID   → 8f6d...-uuid
```

โดย `*RefID` ใน model เช่น:

``` text
MovieRefID
TheaterRefID
ScreenRefID
UserRefID
ShowtimeRefID
SeatRefID
BookingRefID
```

ใช้เป็น relational foreign key ภายในฐานข้อมูล และถูกซ่อนจาก JSON ด้วย
`json:"-"`

------------------------------------------------------------------------

# 🔄 PostgreSQL Migration

โปรเจกต์เคยมีโครงสร้าง MongoDB และมี `MIGRATION_NOTES.md`
สำหรับอธิบายการเปลี่ยนมาใช้ PostgreSQL

Migration หลักประกอบด้วย:

``` text
MongoDB
   │
   ▼
PostgreSQL + GORM
```

พร้อมการเพิ่ม:

-   Screen
-   Seat
-   BookingSeat
-   Payment
-   PasswordResetToken
-   Numeric relational foreign keys
-   UUID business identifiers

และมีการแก้ปัญหา GORM relationship โดยให้ child model เป็นฝ่ายประกาศ
belongs-to relationship

------------------------------------------------------------------------

# 🌍 API Endpoints

## Authentication

``` http
POST /auth/register
POST /auth/login
POST /auth/forgot-password
POST /auth/reset-password

GET  /auth/google
GET  /auth/google/callback
```

------------------------------------------------------------------------

## Users

``` http
GET    /users
POST   /users
PUT    /users/:user_id
DELETE /users/:user_id

GET    /users/me
PUT    /users/me
PUT    /users/me/password
POST   /users/me/profile-image
```

`/users/me` และ profile/password/image endpoints ใช้ JWT Authentication

------------------------------------------------------------------------

## Movies

``` http
GET    /movies
GET    /movies/:movie_id
POST   /movies
PUT    /movies/:movie_id
DELETE /movies/:movie_id
```

มี alias API:

``` http
GET    /api/movies
GET    /api/movies/:movie_id
POST   /api/movies
PUT    /api/movies/:movie_id
DELETE /api/movies/:movie_id
```

------------------------------------------------------------------------

## Theaters

``` http
GET    /theaters
POST   /theaters
PUT    /theaters/:theater_id
DELETE /theaters/:theater_id
```

------------------------------------------------------------------------

## Screens

``` http
GET    /screens
POST   /screens
PUT    /screens/:screen_id
DELETE /screens/:screen_id
```

------------------------------------------------------------------------

## Seats

``` http
GET    /seats
POST   /seats
POST   /seats/generate
PUT    /seats/:seat_id
DELETE /seats/:seat_id
```

------------------------------------------------------------------------

## Showtimes

``` http
GET    /showtimes
POST   /showtimes
PUT    /showtimes/:showtime_id
DELETE /showtimes/:showtime_id
```

------------------------------------------------------------------------

## Seat Lock

``` http
POST   /showtimes/:showtime_id/seat-lock/:seat_id
DELETE /showtimes/:showtime_id/seat-lock/:seat_id
```

------------------------------------------------------------------------

## Real-time Seat Map

``` http
GET /showtimes/:showtime_id/seat-map
```

Endpoint นี้ถูกใช้เป็น WebSocket connection

------------------------------------------------------------------------

## Bookings

``` http
GET    /bookings
POST   /bookings
PUT    /bookings/:booking_id
DELETE /bookings/:booking_id
```

------------------------------------------------------------------------

## Payments

``` http
GET    /payments
POST   /payments
PUT    /payments/:payment_id
DELETE /payments/:payment_id
```

------------------------------------------------------------------------

## Admin

``` http
GET /admin/bookings
GET /admin/me
PUT /admin/me
```

Admin routes ใช้:

``` text
JWT + role=admin
```

------------------------------------------------------------------------

## Audit Logs

``` http
GET    /audit_logs
POST   /audit_logs
PUT    /audit_logs/:audit_log_id
DELETE /audit_logs/:audit_log_id
```

------------------------------------------------------------------------

# 📚 Swagger

Backend มี Swagger/OpenAPI files:

``` text
docs/swagger.json
docs/swagger.yaml
docs/docs.go
```

และเปิด Swagger UI ผ่าน:

``` text
http://localhost:8080/swagger/index.html
```

------------------------------------------------------------------------

# 🌐 Backend Configuration

Backend อ่านค่าจาก `.env`

ตัวแปรที่ source code ใช้:

``` env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=

REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
REDIS_DB=

JWT_SECRET=

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=
```

> อย่า commit ค่า password, JWT secret หรือ OAuth secret จริงลง Git

------------------------------------------------------------------------

# 🌐 Frontend Configuration

Frontend ใช้:

``` env
VITE_API_BASE_URL=http://localhost:8080
VITE_WS_BASE_URL=ws://localhost:8080
```

Axios instance ใช้ base URL:

``` text
/api
```

และ Vite จะ proxy:

``` text
/api/*
      │
      ▼
http://localhost:8080/*
```

ตัวอย่าง:

``` text
Frontend
GET /api/movies

        ↓ Vite Proxy

Backend
GET /movies
```

------------------------------------------------------------------------

# 🔐 Frontend Session

Frontend ใช้ `sessionStorage`

สำหรับ User:

``` text
user_auth
```

สำหรับ Admin:

``` text
admin_auth
```

Axios interceptor จะตรวจสอบ URL ปัจจุบัน:

``` text
/admin/*
```

ถ้าเป็น Admin จะอ่าน:

``` text
admin_auth
```

ถ้าเป็น User จะอ่าน:

``` text
user_auth
```

จากนั้นแนบ:

``` http
Authorization: Bearer <JWT>
```

------------------------------------------------------------------------

# 🛡️ Route Guard

Vue Router มี route meta:

``` ts
meta: {
  requiresAuth: true,
  role: 'user'
}
```

หรือ:

``` ts
meta: {
  requiresAuth: true,
  role: 'admin'
}
```

ดังนั้น:

``` text
User protected route
        │
        ├── ไม่มี login → /login
        │
        └── login แล้ว → เข้าได้
```

และ:

``` text
Admin protected route
        │
        ├── ไม่มี admin session → /admin/login
        │
        └── admin session → เข้าได้
```

------------------------------------------------------------------------

# 🎨 Frontend UI

ระบบมี Layout แยก:

``` text
UserLayout
AdminLayout
```

Admin มี:

-   Sidebar
-   Topbar
-   Dashboard navigation
-   Movie management
-   Theater management
-   Screen management
-   Seat management
-   Showtime management
-   Booking management
-   Payment management
-   Audit Logs

User มี:

-   Navbar
-   Footer
-   Movie browsing
-   Showtime selection
-   Seat selection
-   Booking flow
-   Profile
-   Booking history

------------------------------------------------------------------------

# 🎟️ Seat Booking UI Concept

หน้า Seat Selection ถูกออกแบบให้มี:

``` text
Movie / Theater / Showtime
        │
        ▼
    SCREEN
───────────────

A1 A2 A3 A4 A5 ...
B1 B2 B3 B4 B5 ...
C1 C2 C3 C4 C5 ...
D1 D2 D3 D4 D5 ...
...
```

ผู้ใช้สามารถเลือกหลายที่นั่งก่อนดำเนินการต่อ

ราคาจะคำนวณตาม Seat Type:

``` text
Standard → Price Standard
Premium  → Price Premium
VIP      → Price VIP
```

------------------------------------------------------------------------

# 📱 Responsive UI

Frontend ใช้ Tailwind CSS และหน้า Seat Selection
มีการออกแบบให้รองรับหน้าจอหลายขนาด

แนวทาง UI ที่เหมาะกับระบบ:

``` text
Desktop
┌─────────────────────────────┐
│ Seat Map          Summary   │
└─────────────────────────────┘

Tablet
┌─────────────────────────────┐
│ Seat Map                    │
├─────────────────────────────┤
│ Booking Summary             │
└─────────────────────────────┘

Mobile
┌─────────────────────────────┐
│ Seat Map                    │
├─────────────────────────────┤
│ Selected / Total            │
│ [ Continue ]                │
└─────────────────────────────┘
```

------------------------------------------------------------------------

# ▶️ วิธี Run Backend

เข้าโฟลเดอร์:

``` bash
cd Backend
```

ติดตั้ง dependencies:

``` bash
go mod tidy
```

จากนั้น:

``` bash
go run ./cmd
```

Backend จะรันที่:

``` text
http://localhost:8080
```

ตรวจสอบ:

``` text
http://localhost:8080/
```

ควรได้:

``` json
{
  "message": "Cinema Ticket Booking API"
}
```

Health check:

``` text
http://localhost:8080/health
```

ควรได้:

``` json
{
  "status": "OK"
}
```

------------------------------------------------------------------------

# ▶️ วิธี Run Frontend

เข้าโฟลเดอร์:

``` bash
cd Frontend
```

ติดตั้ง dependencies:

``` bash
npm install
```

Run development server:

``` bash
npm run dev
```

จากนั้นเปิด:

``` text
http://localhost:5173
```

------------------------------------------------------------------------

# 🏗️ Frontend Build

ตรวจ TypeScript และ build:

``` bash
npm run build
```

Type check:

``` bash
npm run type-check
```

Preview production build:

``` bash
npm run preview
```

------------------------------------------------------------------------

# 🗃️ PostgreSQL

Backend ต้องสามารถเชื่อมต่อ PostgreSQL ได้

ค่าใน `.env` ต้องสอดคล้องกับ database ที่มีอยู่ เช่น:

``` env
DB_HOST=localhost
DB_PORT=5432
DB_USER=...
DB_PASSWORD=...
DB_NAME=cinema_booking
DB_SSLMODE=disable
```

ชื่อ database จริงให้ยึดตาม `.env` ของเครื่องที่นำโปรเจกต์ไปใช้งาน

เมื่อ Backend start จะเรียก:

``` go
config.DB.AutoMigrate(...)
```

เพื่อสร้าง/ปรับ schema ตาม models

------------------------------------------------------------------------

# 🧠 Redis

Redis ใช้สำหรับ:

1.  Distributed Seat Lock
2.  Pub/Sub Event

ค่า default ใน source:

``` text
Host: localhost
Port: 6379
DB: 0
```

สามารถเปลี่ยนผ่าน:

``` env
REDIS_HOST=
REDIS_PORT=
REDIS_PASSWORD=
REDIS_DB=
```

------------------------------------------------------------------------

# 🔄 Complete Booking Flow

ระบบโดยรวมทำงานประมาณนี้:

``` text
1. User Login
       │
       ▼
2. Browse Movies
       │
       ▼
3. Select Movie
       │
       ▼
4. Select Theater
       │
       ▼
5. Select Showtime
       │
       ▼
6. Open Seat Map
       │
       ▼
7. Select Seats
       │
       ▼
8. Redis Seat Lock
       │
       ▼
9. Create Booking
       │
       ▼
10. Booking = pending
       │
       ▼
11. Payment
       │
       ▼
12. Payment = paid
       │
       ▼
13. Verify Seat Lock
       │
       ▼
14. Seat = BOOKED
       │
       ▼
15. Booking = confirmed
       │
       ▼
16. Publish booking.success
       │
       ▼
17. Booking Success
```

------------------------------------------------------------------------

# ⚡ Concurrency / Double Booking Protection

จุดสำคัญของระบบคือการป้องกันการจองที่นั่งเดียวกันพร้อมกัน

ตัวอย่าง:

``` text
User A                  User B
   │                       │
   │ Lock A5               │
   ▼                       ▼
 Redis SETNX(A5)       Redis SETNX(A5)
   │                       │
   │ SUCCESS               │ FAIL
   ▼                       ▼
 A5 LOCKED             A5 already locked
```

ดังนั้น Redis ทำหน้าที่เป็น distributed lock ไม่ใช่เพียง cache

------------------------------------------------------------------------

# 🧩 Service Layer

Backend แบ่ง Business Logic เป็น Services เช่น:

``` text
AuthService
BookingService
PaymentService
MovieService
TheaterService
ScreenService
SeatService
SeatLockService
SeatMapService
ShowtimeService
UserService
AdminService
AuditLogService
MessageQueueService
```

แนวทางโดยรวม:

``` text
Controller
    │
    ▼
 Service
    │
    ▼
Repository
    │
    ▼
Database / Redis
```

------------------------------------------------------------------------

# 🗄️ Repository Layer

Repository รับผิดชอบการติดต่อ database เช่น:

``` text
Create
Find
Update
Delete
Transaction
```

และจัดการ conversion ระหว่าง:

``` text
Business UUID
        ↓
Numeric database ID
```

ตัวอย่าง:

``` text
movie_id
    ↓
MovieRepository
    ↓
Movie.ID
    ↓
MovieRefID
```

------------------------------------------------------------------------

# 🔄 Database Transactions

บาง operation ใช้ database transaction เพื่อป้องกันข้อมูลบางส่วนถูกบันทึกโดยไม่ครบ
เช่น:

### Generate Seats

``` text
Begin Transaction
      │
      ├── ตรวจสอบ Seat ซ้ำ
      │
      ├── Create Seats
      │
      └── Commit
```

ถ้าเจอ duplicate หรือ database error:

``` text
Rollback
```

### Booking

Repository มีการเชื่อม Booking กับ BookingSeat ภายใน transaction ตาม
implementation ใน repository layer

------------------------------------------------------------------------

# 📊 Main Data Entities

  Entity               หน้าที่
  -------------------- ----------------------------
  User                 ผู้ใช้งานระบบ
  Movie                ภาพยนตร์
  Theater              สาขาโรงภาพยนตร์
  Screen               โรง/ห้องฉาย
  Seat                 ที่นั่ง
  Showtime             รอบฉาย
  Booking              รายการจอง
  BookingSeat          ความสัมพันธ์ Booking กับ Seat
  Payment              การชำระเงิน
  AuditLog             ประวัติการกระทำ
  PasswordResetToken   Token สำหรับ Reset Password

------------------------------------------------------------------------

# 🧪 Current Development Notes

จาก source code ปัจจุบัน มีประเด็นที่ควรรับรู้ก่อนนำไป production:

1.  `forgot-password` ส่ง reset token กลับใน response เพื่อ
    development/testing แทนการส่งผ่าน email
2.  Google OAuth ยังมี implementation อยู่ใน Backend และ Frontend
3.  `Dockerfile` มีอยู่ แต่ไม่มี `docker-compose.yml` ในไฟล์ที่แนบมา
4.  Redis มี fallback logging หาก Redis ไม่พร้อม แต่ฟีเจอร์ Seat Lock ต้องใช้
    Redis
5.  Redis Seat Lock มี TTL 15 นาที
6.  Redis Pub/Sub consumer ปัจจุบันเป็น mock notification
7.  Swagger files มีอยู่ แต่ไม่ได้ครอบคลุมทุก route ที่มีใน source ปัจจุบัน เช่นบาง
    authentication/admin-related routes
8.  Route protection ถูกใช้กับ `/users/me`, `/admin/*` และบาง frontend
    routes แต่ endpoint CRUD บางกลุ่มยังถูกประกาศเป็น public route ใน
    `routes/*.go` ดังนั้นการเพิ่ม authorization ควรตรวจสอบตาม requirement
    ก่อนนำขึ้น production
9.  Frontend ใช้ Vite proxy `/api` ไปยัง Backend และ rewrite `/api`
    ออกก่อนส่งต่อ
10. Backend ใช้ `AutoMigrate` ทุกครั้งที่ startup

------------------------------------------------------------------------

# 🔒 Security Notes

ห้ามนำค่า secret จาก `.env` ขึ้น Git เช่น:

``` env
DB_PASSWORD
JWT_SECRET
GOOGLE_CLIENT_SECRET
```

ควรเก็บ `.env` ไว้ใน `.gitignore`

สำหรับ production ควรพิจารณาเพิ่มเติม:

-   HTTPS
-   Secure cookie / token strategy
-   Rate limiting
-   Email-based password reset
-   Authorization สำหรับทุก admin mutation
-   Authorization สำหรับ user-owned booking/payment data
-   Validation ของ ownership ก่อนอ่าน/แก้ไขข้อมูล
-   CORS ที่จำกัด origin ตาม environment
-   Secret management
-   Logging / monitoring
-   Redis availability strategy
-   WebSocket authentication

------------------------------------------------------------------------

# 📝 Project History

โปรเจกต์มี `MIGRATION_NOTES.md` ที่บันทึกการเปลี่ยนจาก:

``` text
MongoDB
   ↓
PostgreSQL + GORM
```

และมีการปรับ database design หลายรอบเพื่อให้ GORM และ PostgreSQL จัดการ
relationship ได้ถูกต้อง

แนวทางปัจจุบันใช้:

``` text
Numeric Primary Key
+
Business UUID
+
Numeric Foreign Key (*RefID)
```

เพื่อให้ database relationship มีความชัดเจน และยังรักษา API ที่ใช้ UUID ไว้ได้

------------------------------------------------------------------------

# 🚀 Quick Start

## 1. Start PostgreSQL

ตรวจสอบว่า PostgreSQL ทำงานอยู่

## 2. Start Redis

ตรวจสอบว่า Redis ทำงานอยู่

## 3. Configure Backend

สร้าง/แก้ไข:

``` text
Backend/.env
```

ใส่:

``` env
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

JWT_SECRET=

GOOGLE_CLIENT_ID=
GOOGLE_CLIENT_SECRET=
GOOGLE_REDIRECT_URL=
```

## 4. Start Backend

``` bash
cd Backend
go mod tidy
go run ./cmd
```

## 5. Start Frontend

เปิด terminal ใหม่:

``` bash
cd Frontend
npm install
npm run dev
```

## 6. Open Website

``` text
http://localhost:5173
```

------------------------------------------------------------------------

# 🔗 Important URLs

### Frontend

``` text
http://localhost:5173
```

### Backend

``` text
http://localhost:8080
```

### Backend Health

``` text
http://localhost:8080/health
```

### Swagger

``` text
http://localhost:8080/swagger/index.html
```

------------------------------------------------------------------------

# 🎯 Project Goal

ระบบนี้มีเป้าหมายเป็นระบบจองตั๋วภาพยนตร์ที่มีทั้ง:

-   User Website
-   Admin Dashboard
-   Movie Management
-   Theater Management
-   Screen Management
-   Seat Management
-   Showtime Management
-   Seat Selection
-   Redis Distributed Seat Lock
-   Real-time Seat Map
-   Booking
-   Payment
-   JWT Authentication
-   Google OAuth
-   Password Reset
-   Audit Logs
-   Redis Pub/Sub
-   PostgreSQL relational database

โดย architecture หลักคือ:

``` text
Vue 3
   ↓
Axios / REST API
   ↓
Go + Gin
   ↓
Service Layer
   ↓
Repository Layer
   ↓
PostgreSQL

          ↘ Redis
             ├── Seat Lock
             ├── Pub/Sub
             └── Real-time coordination

          ↘ WebSocket
             └── Seat Map updates
```

------------------------------------------------------------------------

# 👩‍💻 Development

Frontend:

``` bash
npm run dev
npm run type-check
npm run build
npm run preview
```

Backend:

``` bash
go mod tidy
go run ./cmd
```

------------------------------------------------------------------------

## 📌 Note

README นี้จัดทำจาก source code ที่แนบมาในรอบนี้:

-   `Backend(10).zip`
-   `Frontend(2).zip`

หาก source code เปลี่ยนแปลงภายหลัง เช่น เพิ่ม API, เปลี่ยน Authentication, เปลี่ยน
database หรือเปลี่ยน flow การจอง ควรอัปเดต README ให้ตรงกับ implementation
ล่าสุด


## Responsive Design

This web application is responsive and supports:
- Mobile
- Tablet
- Desktop
