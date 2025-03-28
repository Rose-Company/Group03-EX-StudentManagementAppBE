# Student Management Application V3.0

A clean architecture-based REST API for student management system developed in Go.

## 📌 Tính năng chính

✅ Xác thực & phân quyền người dùng (Admin & Student)  
✅ Quản lý thông tin sinh viên  
✅ Quản lý giảng viên & khoa  
✅ API RESTful theo tiêu chuẩn  
✅ JWT Authentication  
✅ Kết nối PostgreSQL  
✅ Hỗ trợ logging mechanism 

- Logging: ![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1zCnBiLaXG0_FXsMJADCTP6QotH2f5O7v)
- Database: ![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1BWt2RhYNFv75lJ-AtvPTgho0oXOA_Z55)
- Các API quản lý thông tin: https://drive.google.com/file/d/1fItGjQCD1uWGDPYSrl6-TAjpGQg7c_pw/view?usp=sharing

## 📌 Tính năng chính V3.0

✅ MSSV phải là duy nhất  (DONE)
✅ Email phải thuộc một tên miền nhất định và có thể cấu hình động (configurable)
✅ Số điện thoại phải có định dạng hợp lệ theo quốc gia (configurable) 
✅ Tình trạng sinh viên chỉ có thể thay đổi theo một số quy tắc (configurable)

- Check MSSV: ![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1K31pH2YomSiaKopNwA9LRlJlm1Pj_nwu)
- Check Email:![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1o8hWLmgUpji-eu2a1d7aCNjDvT9uQ4g4)
- Check SĐT: ![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1oRy3XGR8BKNLNy16Y3-YrcZqItybEywF)
- Check tình trạng SV: ![Mô tả ảnh](https://drive.google.com/uc?export=view&id=1T7IAWbySnfEjD8XOPi6v93oOUA-luzqT)

## Cấu trúc source code

```
GROUP03-EX-STUDENTMANAGEMENTAPPBE/
├── cmd/                    # Command-line application entry points
│   ├── root.go             # Root command entry point
│   └── server.go           # Server command implementation
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
│   │   └── app.go          # Main application bootstrap
│   ├── handlers/           # HTTP request handlers (API Layer)
│   │   ├── admin/          # Admin-related handlers
│   │   │   └── handler.go  # Admin handler implementation
│   │   ├── auth/           # Authentication handlers
│   │   │   ├── handler.go  # Auth handler implementation
│   │   │   └── login.go    # Login functionality
│   │   ├── faculty/        # Faculty-related handlers
│   │   │   ├── faculty_crud.go # Faculty CRUD operations
│   │   │   └── handler.go  # Faculty handler implementation
│   │   ├── program/        # Program-related handlers
│   │   ├── student/        # Student-related handlers
│   │   │   ├── handler.go         # Main student handler
│   │   │   ├── student_edit.go    # Student edit operations
│   │   │   ├── student_info.go    # Student info operations
│   │   │   ├── student_list.go    # Student listing operations
│   │   │   └── student_statuses.go # Student status operations
│   │   └── base.go         # Base handler functionality
│   ├── models/             # Data models (Domain Layer)
│   │   ├── admin/          # Admin models
│   │   │   └── file.go     # File model for admin operations
│   │   ├── auth/           # Authentication models
│   │   ├── faculty/        # Faculty models
│   │   ├── gdrive/         # Google Drive integration models
│   │   ├── program/        # Program models
│   │   ├── student/        # Student models
│   │   ├── student_status/ # Student status models
│   │   └── base.go         # Base model functionality
│   ├── repositories/       # Data access layer (Repository Layer)
│   │   ├── admin/          # Admin repository
│   │   ├── faculty/        # Faculty repository
│   │   │   └── repository.go # Faculty repository implementation
│   │   ├── program/        # Program repository
│   │   │   └── repository.go # Program repository implementation
│   │   ├── student/        # Student repository
│   │   │   ├── repository.go    # Main student repository
│   │   │   ├── student_addresses/  # Student addresses repository
│   │   │   ├── student_documents/  # Student documents repository
│   │   │   ├── student_status/ # Student status repository 
│   │   │   └── user/       # User repository for students
│   │   └── base.go         # Base repository functionality
│   ├── services/           # Business logic (Service Layer)
│   │   ├── auth/           # Authentication services
│   │   ├── faculty/        # Faculty services
│   │   ├── program/        # Program services
│   │   ├── student/        # Student services
│   │   └── base.go         # Base service functionality
│   ├── middleware/         # HTTP middleware
│   │   ├── admin_authentication.go  # Admin auth middleware
│   │   └── user_authentication.go   # User auth middleware
│   └── keys/               # API Keys and credentials
│       └── google_service_cre.json # Google API credentials
├── pkg/                    # Public library code
├── script/                 # Database scripts
│   ├── hw2_init_data.sql   # Initial data setup
│   ├── hw2.sql             # Homework 2 scripts
│   └── script.sql          # Main SQL initialization scripts
├── .env                    # Environment variables
├── .gitignore              # Git ignore file
├── docker-compose.yaml     # Docker composition for services
├── go.mod                  # Go module definition
├── go.sum                  # Go module checksums
├── main.go                 # Main application entry point
├── README.md               # Project documentation
└── sample.env              # Sample environment configuration
```

## Hướng dẫn cài đặt & chạy chương trình

### Yêu cầu hệ thống
- Go 1.16 hoặc cao hơn
- PostgreSQL (hoặc cơ sở dữ liệu được cấu hình trong config.yaml)
- Git
- Google Key Credentials


## Cấu hình Google Key Credentials

### Truy cập Google Cloud Console
👉 [Google Cloud Console](https://console.cloud.google.com/)

### Tạo Dự án mới hoặc chọn dự án hiện có
- Mở **Google Cloud Console** và đăng nhập.
- Chọn **Create Project** để tạo dự án mới hoặc chọn một dự án có sẵn.

### Kích hoạt API cần thiết
- Điều hướng đến **APIs & Services** → **Library**.
- Tìm kiếm và kích hoạt các API cần thiết (ví dụ: **Google OAuth**, **Drive API**...).

### Tạo Key Credentials
1. Điều hướng đến **APIs & Services** → **Credentials**.
2. Chọn **Create Credentials** → **Service Account**.
3. Nhập thông tin cần thiết và tạo **Service Account**.
4. Chọn **Manage Keys** → **Add Key** → **Create New Key**.
5. Chọn định dạng **JSON** và tải file xác thực về máy.

📌 **Lưu ý:**  
- File JSON chứa thông tin xác thực cần được bảo mật, không chia sẻ công khai.  
- Cấu hình biến môi trường để ứng dụng sử dụng Google Key Credentials:
  ```bash
  export GOOGLE_DRIVE_CREDENTIALS_FILE="/path/to/your-google-key.json"
  ```


### Cài đặt Source cho V2.0

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
