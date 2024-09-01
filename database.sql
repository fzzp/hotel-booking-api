CREATE SCHEMA `db_hotel_bookings` DEFAULT CHARACTER SET utf8mb4 ;

CREATE TABLE IF NOT EXISTS `users` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
    phone_number varchar(11) NOT NULL COMMENT "用户手机号码",
    password_hash varchar(60) NOT NULL COMMENT "哈希密码",
    username varchar(60) NOT NULL COMMENT "用户名称",
    avatar varchar(255) NOT NULL COMMENT "用户头像",
    role tinyint(1) NOT NULL DEFAULT 0 COMMENT "角色;0:客户,1:管理员",
    created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "更新时间",
    is_deleted tinyint(1) NOT NULL DEFAULT 1 COMMENT "是否删除, 1:正常, 0:删除",
    INDEX idx_phoneNumber(phone_number),
    INDEX idx_isDeleted(is_deleted),
    UNIQUE unq_phoneNumber(phone_number)
) COMMENT="用户表";


CREATE TABLE IF NOT EXISTS `sessions` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
	user_id INT UNSIGNED NOT NULL COMMENT "用户id标识",
	token_id char(36) NOT NULL COMMENT "token唯一标识uuid",
	refresh_token varchar(2048) NOT NULL COMMENT "刷新token",
	client_ip varchar(16) NOT NULL COMMENT "用户登陆时ip",
	user_agent varchar(255) NOT NULL COMMENT "客户端信息",
	expires_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "过期时间",
	created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
	INDEX idx_userId(user_id),
	INDEX idx_tokenId(token_id)
) COMMENT="会话信息表";


CREATE TABLE IF NOT EXISTS `hotels` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
    name varchar(200) NOT NULL COMMENT "酒店名称",
    address varchar(255) NOT NULL COMMENT "酒店地址",
	logo varchar(255) NOT NULL COMMENT "logo图片",
    created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "更新时间",
    is_deleted tinyint(1) NOT NULL DEFAULT 1 COMMENT "是否删除, 1:正常, 0:删除",
    INDEX idx_isDeleted(is_deleted)
) COMMENT="酒店表";


CREATE TABLE IF NOT EXISTS `rooms` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
    hotel_id INT UNSIGNED NOT NULL COMMENT "外键,酒店表id",
    room_no varchar(10) NOT NULL COMMENT "客房编号/门牌号",
    images varchar(2048) NOT NULL COMMENT "客房图片,可有多张,';'号分割",
    price INT UNSIGNED NOT NULL DEFAULT 0 COMMENT "客房价格(单位分)",
    capacity INT UNSIGNED NOT NULL DEFAULT 1 COMMENT "容纳人数",
    status varchar(32) NOT NULL DEFAULT "available" COMMENT "客房当前状态:available,occupied,maintain",
    room_type_id INT UNSIGNED NOT NULL COMMENT "外键,客房类型表id",
    description varchar(255) NOT NULL DEFAULT "" COMMENT "客房描述",
    created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "更新时间",
    is_deleted tinyint(1) NOT NULL DEFAULT 1 COMMENT "是否删除, 1:正常, 0:删除",
    INDEX idx_hotelId(hotel_id),
	INDEX idx_roomNumber(room_no),
    INDEX idx_roomTypeId(room_type_id),
    INDEX idx_isDeleted(is_deleted)
) COMMENT="客房表";

CREATE TABLE IF NOT EXISTS `room_types` (
	id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
	name varchar(60) NOT NULL COMMENT "客房类型",
	description varchar(255) NOT NULL DEFAULT "" COMMENT "类型描述",
    UNIQUE unq_name(name)
) COMMENT="客房类型表";

CREATE TABLE IF NOT EXISTS `bookings` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键", 
    user_id INT UNSIGNED NOT NULL COMMENT "外键,用户表id",
    check_in DATE NOT NULL COMMENT "入住日期",
    check_out DATE NOT NULL COMMENT "退房日期",
    members INT UNSIGNED NOT NULL COMMENT "入住人数",
    total_amount INT UNSIGNED NOT NULL COMMENT "总金额(单位分)",
    created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
    updated_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "更新时间",
    is_deleted tinyint(1) NOT NULL DEFAULT 1 COMMENT "是否删除, 1:正常, 0:删除",
    INDEX idx_userId(user_id),
    INDEX idx_isDeleted(is_deleted)
) COMMENT="预订表";


CREATE TABLE IF NOT EXISTS `booking_rooms` (
	id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键",
	booking_id INT UNSIGNED NOT NULL COMMENT "外键,预定id",
	room_id INT UNSIGNED NOT NULL COMMENT "外键,客房id"
) COMMENT="预订客房关联表";


CREATE TABLE IF NOT EXISTS `payments` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键", 
    booking_id INT UNSIGNED NOT NULL COMMENT "外键,预定表id",
    payment_type tinyint(1) NOT NULL DEFAULT 1 COMMENT "支付类型(1:未支付,2:现金,3:微信,4:支付宝,5:其他,6:取消支付)",
	payment_time TIMESTAMP COMMENT "支付时间",
    payment_amount INT UNSIGNED NOT NULL COMMENT "付款金额(单位分)",
    is_deleted tinyint(1) NOT NULL DEFAULT 0 COMMENT "是否删除, 0:正常, 1:删除",
    INDEX idx_bookingId(booking_id),
    INDEX idx_paymentType(payment_type)
) COMMENT="支付表";


CREATE TABLE IF NOT EXISTS `notifications` (
    id INT UNSIGNED NOT NULL PRIMARY KEY AUTO_INCREMENT COMMENT "主键", 
    user_id INT UNSIGNED NOT NULL COMMENT "外键,用户表id",
    message varchar(2000) NOT NULL COMMENT "消息内容",
    is_read tinyint(1) NOT NULL DEFAULT 0 COMMENT "是否已读;0:未读,1:已读",
    created_at TIMESTAMP NOT NULL DEFAULT NOW() COMMENT "创建时间",
    is_deleted tinyint(1) NOT NULL DEFAULT 1 COMMENT "是否删除, 1:正常, 0:删除",
    INDEX idx_userId(user_id),
    INDEX idx_isRead(is_read),
    INDEX idx_isDeleted(is_deleted)
) COMMENT="消息表";


-- >>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>>> --
-- 创建一批基础数据

INSERT INTO 
    hotels (name, address, logo)
VALUES
    ("星河驿站", "xxx路77号", "http://"),
    ("雾隐山庄", "yyy路88号", "http://"),
    ("逸境雅舍", "zzz路99号", "http://");


INSERT INTO 
    room_types(name, description)
VALUES
    ('单人客房', '单人客房-description'),
    ('标准双床房', '标准双床房-description'),
    ('豪华双床房', '豪华双床房-description');

INSERT INTO 
    rooms(
        hotel_id, room_no, images, price, room_type_id, capacity
    )
VALUES
    (1, "F06", "http://", 10000, 1, 1),
    (1, "F08", "http://", 20000, 2, 2),
    (1, "F10", "http://", 30000, 3, 2),
    (2, "A22", "http://", 10000, 1, 1),
    (2, "A66", "http://", 20000, 2, 2),
    (2, "A88", "http://", 30000, 3, 2),
    (3, "801", "http://", 10000, 1, 1),
    (3, "802", "http://", 20000, 2, 2),
    (3, "803", "http://", 30000, 3, 2);


