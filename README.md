# Student Management Application

A clean architecture-based REST API for student management system developed in Go.

## Cấu trúc source code

```
GROUP03-EX-STUDENTMANAGEMENTAPPBE/
├── common/                 # Common utilities and helpers
│   ├── common.go           # Shared functions
│   ├── db.go               # Database connection utilities
│   ├── error_messages.go   # Error message definitions
│   ├── error.go            # Error handling utils
│   ├── helper.go           # Helper functions
│   ├── http_status.go      # HTTP status code definitions
│   ├── jwt.go              # JWT authentication utilities
│   ├── key.go              # Encryption/security keys
│   ├── request.go          # Request handling utilities
│   ├── response.go         # Response formatting utilities
│   └── time.go             # Time utilities
├── config/                 # Configuration management
│   ├── config.go           # Config loading functions
│   └── config.yaml         # Application configuration
├── internal/               # Core application code
│   ├── app/                # Application setup and initialization
│   ├── handlers/           # HTTP request handlers (API Layer)
│   │   ├── auth/           # Authentication handlers
│   │   │   ├── auth_handler.go
│   │   │   └── login.go
│   │   ├── faculty/        # Faculty-related handlers
│   │   ├── student/        # Student-related handlers
│   │   └── base.go         # Base handler functionality
│   ├── models/             # Data models (Domain Layer)
│   │   ├── auth/           # Authentication models
│   │   ├── faculty/        # Faculty models
│   │   ├── student/        # Student models
│   │   └── base.go         # Base model functionality
│   ├── repositories/       # Data access layer (Repository Layer)
│   │   ├── faculty/        # Faculty repository
│   │   ├── student/        # Student repository
│   │   ├── student_status/ # Student status repository 
│   │   ├── user/           # User repository
│   │   └── base.go         # Base repository functionality
│   ├── services/           # Business logic (Service Layer)
│   │   ├── auth/           # Authentication services
│   │   ├── faculty/        # Faculty services
│   │   ├── student/        # Student services
│   │   └── base.go         # Base service functionality
│   └── middleware/         # HTTP middleware
│       ├── admin_authentication.go  # Admin auth middleware
│       └── user_authentication.go   # User auth middleware
└── script/                 # Database scripts
    └── script.sql          # SQL initialization scripts
```

## Hướng dẫn cài đặt & chạy chương trình

### Yêu cầu hệ thống
- Go 1.16 hoặc cao hơn
- PostgreSQL (hoặc cơ sở dữ liệu được cấu hình trong config.yaml)
- Git

### Cài đặt

1. Clone repository:
```bash
git clone https://github.com/your-username/GROUP03-EX-STUDENTMANAGEMENTAPPBE.git
cd GROUP03-EX-STUDENTMANAGEMENTAPPBE
```

2. Cài đặt các dependencies:
```bash
go mod download
```

3. Cấu hình cơ sở dữ liệu:
   - Cập nhật thông tin kết nối trong file `config/config.yaml`
   - Chạy script khởi tạo cơ sở dữ liệu:
   ```bash
   psql -U your_username -d your_database < script/script.sql
   ```

### Biên dịch

```bash
go build -o studentapp ./internal/app
```

### Chạy chương trình

```bash
./studentapp
```

Hoặc có thể chạy trực tiếp mà không cần biên dịch:

```bash
go run main.go
```

Ứng dụng sẽ chạy mặc định tại: `http://localhost:8080`

## Mô tả Clean Architecture trong ứng dụng

Ứng dụng này được xây dựng theo mô hình Clean Architecture với các lớp tách biệt rõ ràng:

### 1. Luồng hoạt động

1. **Client** gửi request đến API.
2. **Handlers** tiếp nhận request, xử lý các nhiệm vụ cơ bản như parsing request và validation.
3. **Middleware** xác thực người dùng và kiểm tra quyền truy cập.
4. **Handlers** gọi **Services** để xử lý logic nghiệp vụ.
5. **Services** gọi **Repositories** để thao tác với cơ sở dữ liệu.
6. Kết quả được trả về cho **Client** qua **Handlers**.

### 2. Các thành phần chính

- **Handlers (API Layer)**: Xử lý HTTP requests và responses.
- **Services (Business Logic Layer)**: Chứa logic nghiệp vụ của ứng dụng.
- **Repositories (Data Access Layer)**: Tương tác với cơ sở dữ liệu.
- **Models (Domain Layer)**: Định nghĩa cấu trúc dữ liệu của ứng dụng.
- **Middleware**: Xử lý xác thực và phân quyền.
- **Common**: Các tiện ích dùng chung.

### 3. Mô hình tương tác

```
Client
   |
   |--- [HTTP Request] ---> Middleware (Authentication)
   |                          |
   |                          |--- [Request] ---> Handlers
   |                                               |
   |                                               |--- Calls ---> Services
   |                                                                |
   |                                                                |--- Calls ---> Repositories
   |                                                                                 |
   |                                                                                 |---> Database
   |                          
   |<--- [HTTP Response] ----------------------------------------------------|
```

Kiến trúc này giúp tăng tính modular, dễ bảo trì và mở rộng của ứng dụng, đồng thời hỗ trợ testing hiệu quả cho từng thành phần.
