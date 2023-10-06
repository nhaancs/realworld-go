# bhms
Hệ thống quản lý nhà trọ được phát triển với Go, Flutter, Angular, Postgres, và Kubernetes.

## Bản thử nghiệm (MVP)

### Kế hoạch phát triển
- [ ] **Xác định các tính năng cho bản thử nghiệm và thiết kế cơ sở dữ liệu &larr;**
- [ ] Thiết kế và phát triển API cho ứng dụng di động
- [ ] Thiết kế và phát triển ứng dụng Android, iOS
- [ ] Thiết kế và phát triển API cho backoffice
- [ ] Thiết kế và phát triển backoffice cho đội ngũ vận hành

### Tính năng dành cho chủ trọ
- Tạo khu trọ và quản lý danh sách phòng
- Tạo hợp đồng cho người đến thuê
- Quản lý điện, nước, và các dịch vụ khác
- Quản lý thông tin thanh toán (tài khoản ngân hàng)
- Tạo và chia sẻ hóa đơn cho từng phòng


### Thiết kế cơ sở dữ liệu

#### Bảng `users`
Bảng `users` lưu các thông tin cần thiết của người dùng.

| Tên cột         | Kiểu dữ liệu | Bắt buộc | Mô tả                                   |
|-----------------|--------------|:--------:|-----------------------------------------|
| `id`            | uuid         | &check;  | Khóa chính                              |
| `first_name`    | text         | &check;  | Tên                                     |
| `last_name`     | text         | &cross;  | Họ và chữ lót                           |
| `phone`         | text         | &check;  | Dùng để đăng nhập                       |
| `password_hash` | text         | &check;  | Mật khẩu đã được mã hóa                 |
| `roles`         | text[]       | &check;  | Phân quyền: USER, ADMIN                 |
| `status`        | text         | &check;  | Trạng thái người dùng: ACTIVE, DISABLED |
| `created_at`    | timestamp    | &check;  | Thời gian khởi tạo                      |
| `updated_at`    | timestamp    | &cross;  | Thời gian cập nhật                      |

#### Bảng `divisions`
Bảng `divisions` lưu thông tin các đơn vị hành chính của Việt Nam.

| Tên cột      | Kiểu dữ liệu | Bắt buộc | Mô tả                                                                     |
|--------------|--------------|:--------:|---------------------------------------------------------------------------|
| `id`         | serial       | &check;  | Khóa chính, tự tăng                                                       |
| `name`       | text         | &check;  | Tên đơn vị hành chính                                                     |
| `code`       | integer      | &check;  | Mã đơn vị hành chính                                                      |
| `level`      | smallint     | &check;  | Cấp đơn vị hành chính - 1: Tỉnh, thành phố, 2: quận, huyện, 3: xã, phường |
| `parent_id`  | uuid         | &cross;  | ID của cấp cao hơn                                                        |
| `created_at` | timestamp    | &check;  | Thời gian khởi tạo                                                        |
| `updated_at` | timestamp    | &cross;  | Thời gian cập nhật                                                        |

#### Bảng `properties`
Bảng `properties` lưu các thông tin cần thiết của khu trọ.

| Tên cột             | Kiểu dữ liệu | Bắt buộc | Mô tả                                |
|---------------------|--------------|:--------:|--------------------------------------|
| `id`                | uuid         | &check;  | Khóa chính                           |
| `name`              | text         | &check;  | Tên khu trọ                          |
| `address_level1_id` | serial       | &check;  | ID tỉnh, thành phố                   |
| `address_level2_id` | serial       | &check;  | ID quận, huyện                       |
| `address_level3_id` | serial       | &check;  | ID xã, phường                        |
| `street`            | text         | &check;  | Số nhà và tên đường                  |
| `manager_id`        | uuid         | &check;  | ID của chủ trọ                       |
| `status`            | text         | &check;  | Trạng thái khu trọ: ACTIVE, DISABLED |
| `created_at`        | timestamp    | &check;  | Thời gian khởi tạo                   |
| `updated_at`        | timestamp    | &cross;  | Thời gian cập nhật                   |

#### Bảng `blocks`
Bảng `blocks` lưu thông tin các dãy của khu trọ.

| Tên cột       | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|---------------|--------------|:--------:|--------------------|
| `id`          | uuid         | &check;  | Khóa chính         |
| `name`        | text         | &check;  | Tên dãy            |
| `property_id` | uuid         | &check;  | ID của khu trọ     |
| `created_at`  | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at`  | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `floors`
Bảng `floors` lưu thông tin các tầng của khu trọ.

| Tên cột      | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|--------------|--------------|:--------:|--------------------|
| `id`         | uuid         | &check;  | Khóa chính         |
| `name`       | text         | &check;  | Tên tầng           |
| `block_id`   | uuid         | &check;  | ID của dãy         |
| `created_at` | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at` | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `units`
Bảng `units` lưu thông tin các phòng của khu trọ.

| Tên cột       | Kiểu dữ liệu | Bắt buộc | Mô tả              |
|---------------|--------------|:--------:|--------------------|
| `id`          | uuid         | &check;  | Khóa chính         |
| `name`        | text         | &check;  | Tên phòng          |
| `property_id` | uuid         | &check;  | ID của khu trọ     |
| `block_id`    | uuid         | &check;  | ID của dãy         |
| `floor_id`    | uuid         | &check;  | ID của tầng        |
| `created_at`  | timestamp    | &check;  | Thời gian khởi tạo |
| `updated_at`  | timestamp    | &cross;  | Thời gian cập nhật |

#### Bảng `payment_methods`
Bảng `payment_methods` lưu thông tin thanh toán (tài khoản ngân hàng) của các chủ trọ.

| Tên cột          | Kiểu dữ liệu | Bắt buộc | Mô tả                                         |
|------------------|--------------|:--------:|-----------------------------------------------|
| `id`             | bigserial    | &check;  | Khóa chính                                    |
| `name`           | text         | &check;  | Tên ngân hàng                                 |
| `account_name`   | text         | &check;  | Tên chủ tài khoản                             |
| `account_number` | text         | &check;  | Số tài khoản                                  |
| `description`    | text         | &cross;  | Mô tả                                         |
| `enabled`        | boolean      | &check;  | Trạng thái bật/tắt của phương thức thanh toán |
| `property_id`    | uuid         | &check;  | ID của khu trọ                                |
| `created_at`     | timestamp    | &check;  | Thời gian khởi tạo                            |
| `updated_at`     | timestamp    | &cross;  | Thời gian cập nhật                            |

#### Bảng `sessions`
Bảng `sessions` lưu thông tin thuê phòng theo thời gian trên hợp đồng.

| Tên cột            | Kiểu dữ liệu | Bắt buộc | Mô tả                            |
|--------------------|--------------|:--------:|----------------------------------|
| `id`               | uuid         | &check;  | Khóa chính                       |
| `unit_id`          | uuid         | &check;  | ID của phòng                     |
| `start_at`         | timestamp    | &check;  | Thời gian bắt đầu vào ở          |
| `duration_in_days` | smallint     | &check;  | Thời hạn hợp đồng tính theo ngày |
| `price`            | bigint       | &check;  | Giá thuê (đ)                     |
| `num_of_members`   | smallint     | &check;  | Số lượng thành viên              |
| `renew_times`      | smallint     | &check;  | Số lần gia hạn hợp đồng          |
| `description`      | text         | &cross;  | Mô tả                            |
| `enabled`          | boolean      | &check;  | Trạng thái bật/tắt của session   |
| `created_at`       | timestamp    | &check;  | Thời gian khởi tạo               |
| `updated_at`       | timestamp    | &cross;  | Thời gian cập nhật               |

#### Bảng `services`
Bảng `services` lưu thông tin các dịch vụ của khu trọ.

| Tên cột        | Kiểu dữ liệu | Bắt buộc | Mô tả                                                             |
|----------------|--------------|:--------:|-------------------------------------------------------------------|
| `id`           | uuid         | &check;  | Khóa chính                                                        |
| `property_id`  | uuid         | &check;  | ID của khu trọ                                                    |
| `name`         | text         | &check;  | Tên dịch vụ                                                       |
| `price`        | bigint       | &check;  | Giá dịch vụ                                                       |
| `unit`         | text         | &check;  | Đơn vị: kwh, m3, room, member, piece, time                        |
| `invoice_type` | text         | &check;  | Đơn vị tính: PER_USAGE, PER_ROOM, PER_MEMBER, PER_PIECE, PER_TIME | 
| `enabled`      | boolean      | &check;  | Trạng thái bật/tắt dịch vụ                                        |
| `created_at`   | timestamp    | &check;  | Thời gian khởi tạo                                                |
| `updated_at`   | timestamp    | &cross;  | Thời gian cập nhật                                                |


#### Bảng `session_services`
Bảng `session_services` lưu thông tin các dịch vụ của từng session.

| Tên cột        | Kiểu dữ liệu | Bắt buộc | Mô tả          |
|----------------|--------------|:--------:|----------------|
| `session_id`   | uuid         | &check;  | ID của session |
| `service_id`   | uuid         | &check;  | ID của dịch vụ |