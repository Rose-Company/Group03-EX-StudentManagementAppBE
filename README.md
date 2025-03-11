# MÔ TẢ 1 CÁCH TỔNG QUÁT NHẤT CỦA MÔ HÌNH CLEAN ARCHITECTURE TRONG SOURCE NÀY.
Dưới đây là mô tả tổng quát về cách một API trong kiến trúc **Clean Architecture** hoạt động, bao gồm các lớp chính và tương tác với bên thứ ba.

---

## 1. **Luồng Tổng Quát**

1. **Client** gửi một request đến API (ví dụ: `POST /api/resource`).
2. **Handlers** tiếp nhận request, thực hiện các nhiệm vụ cơ bản như:
   - Parse request (dữ liệu JSON, form).
   - Xác thực và tiền xử lý (middleware có thể can thiệp trước).
3. **Handlers** gọi **Services** để xử lý logic nghiệp vụ.
4. **Services** có thể:
   - Gọi **Repositories** để thao tác với cơ sở dữ liệu.
   - Gọi API bên thứ ba qua HTTP hoặc sử dụng các package bên ngoài.
   - Xử lý kết quả và áp dụng các quy tắc nghiệp vụ.
5. Kết quả được trả về **Handlers**, sau đó gửi lại cho client dưới dạng JSON hoặc phản hồi khác.

---

## 2. **Cấu Trúc Lớp Tổng Quát**

- **Handlers (API Layer)**:
  - Tiếp nhận yêu cầu từ client và chuyển tiếp tới **Services**.
  - Xử lý phản hồi và gửi kết quả về client.
- **Middleware**:
  - Xác thực, kiểm tra quyền truy cập, logging.
- **Services (Business Logic Layer)**:
  - Xử lý logic nghiệp vụ của ứng dụng.
  - Gọi đến **Repositories** hoặc các dịch vụ bên thứ ba.
- **Repositories (Data Access Layer)**:
  - Tương tác với cơ sở dữ liệu để thực hiện các thao tác CRUD.
- **External Services (Third-Party API)**:
  - Các dịch vụ bên ngoài được gọi qua HTTP hoặc SDK.
- **Common Utilities**:
  - Các hàm và công cụ dùng chung như logging, JWT, xử lý lỗi.

---

## 3. **Mô Hình Tương Tác Tổng Quát**

```plaintext
Client
   |
   |--- [HTTP Request] ---> Middleware (Validation, Authentication)
   |                          |
   |                          |--- [Request] ---> Handlers
   |                          |                    |
   |                          |                    |--- Calls ---> Services
   |                          |                                   |
   |                          |                                   |--- Calls ---> Repositories (Database)
   |                          |                                             |
   |                          |                                             |---> Third-Party API
   |                          |                                                  (Gemini API, Payment API, etc.)
   |                          |                                   |
   |                          |<--- Returns Response -------------|
   |                          |
   |<--- [HTTP Response] -----|


