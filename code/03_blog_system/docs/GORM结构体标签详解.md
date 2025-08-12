# GORM 结构体标签详解 📚

## 🎯 什么是结构体标签（Struct Tags）？

在Go语言中，结构体标签是附加在结构体字段后面的字符串，用于为字段提供元数据信息。GORM使用这些标签来定义数据库表的结构和约束。

## 📋 User模型字段详解

让我们逐个分析User结构体中的每个字段：

### 1. Username 字段 👤
```go
Username string `json:"username" gorm:"size:50;uniqueIndex:idx_user_username;not null" validate:"required,min=3,max=50"`
```

#### 标签解析：
- **`json:"username"`** 📤
  - **作用**：JSON序列化时的字段名
  - **含义**：当结构体转换为JSON时，这个字段会显示为"username"
  - **示例**：`{"username": "张三"}`

- **`gorm:"size:50;uniqueIndex:idx_user_username;not null"`** 🗄️
  - **`size:50`**：数据库字段长度限制为50个字符
  - **`uniqueIndex:idx_user_username`**：创建唯一索引，索引名为"idx_user_username"
  - **`not null`**：字段不能为空

- **`validate:"required,min=3,max=50"`** ✅
  - **`required`**：必填字段
  - **`min=3`**：最小长度3个字符
  - **`max=50`**：最大长度50个字符

### 2. Email 字段 📧
```go
Email string `json:"email" gorm:"size:100;uniqueIndex:idx_user_email;not null" validate:"required,email"`
```

#### 标签解析：
- **`json:"email"`**：JSON字段名为"email"
- **`gorm:"size:100;uniqueIndex:idx_user_email;not null"`**
  - **`size:100`**：最大长度100字符
  - **`uniqueIndex:idx_user_email`**：邮箱唯一索引
  - **`not null`**：不能为空
- **`validate:"required,email"`**
  - **`required`**：必填
  - **`email`**：必须是有效的邮箱格式

### 3. Password 字段 🔐
```go
Password string `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
```

#### 标签解析：
- **`json:"-"`** 🚫
  - **作用**：JSON序列化时忽略此字段
  - **安全性**：防止密码在API响应中泄露
- **`gorm:"size:255;not null"`**
  - **`size:255`**：支持长密码（加密后）
  - **`not null`**：密码必须存在
- **`validate:"required,min=6"`**
  - **`required`**：必填
  - **`min=6`**：最少6位密码

### 4. Nickname 字段 😊
```go
Nickname string `json:"nickname" gorm:"size:50"`
```

#### 标签解析：
- **`json:"nickname"`**：JSON字段名
- **`gorm:"size:50"`**：最大50字符
- **注意**：没有`not null`，说明昵称可以为空

### 5. Avatar 字段 🖼️
```go
Avatar string `json:"avatar" gorm:"size:255"`
```

#### 标签解析：
- **`json:"avatar"`**：JSON字段名
- **`gorm:"size:255"`**：存储头像URL，最大255字符

### 6. Status 字段 📊
```go
Status string `json:"status" gorm:"size:20;default:active;index" validate:"oneof=active inactive banned"`
```

#### 标签解析：
- **`json:"status"`**：JSON字段名
- **`gorm:"size:20;default:active;index"`**
  - **`size:20`**：状态字符串最大20字符
  - **`default:active`**：默认值为"active"
  - **`index`**：创建普通索引（非唯一）
- **`validate:"oneof=active inactive banned"`**
  - **`oneof`**：只能是指定值中的一个
  - 有效值：active（活跃）、inactive（非活跃）、banned（封禁）

### 7. LastLoginAt 字段 ⏰
```go
LastLoginAt *time.Time `json:"last_login_at"`
```

#### 标签解析：
- **`*time.Time`**：指针类型，可以为nil（表示从未登录）
- **`json:"last_login_at"`**：JSON字段名
- **没有gorm标签**：使用GORM默认设置

### 8. LoginCount 字段 🔢
```go
LoginCount int `json:"login_count" gorm:"default:0"`
```

#### 标签解析：
- **`json:"login_count"`**：JSON字段名
- **`gorm:"default:0"`**：默认值为0
- **用途**：记录用户登录次数

## 🔗 关联关系字段详解

### Profile 关联 👤
```go
Profile *Profile `json:"profile,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
```

#### 标签解析：
- **`json:"profile,omitempty"`**
  - **`omitempty`**：如果为空则在JSON中省略
- **`gorm:"foreignKey:UserID;references:ID;constraint:fk_profiles_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`**
  - **`foreignKey:UserID`**：Profile表中的外键字段名
  - **`references:ID`**：引用User表的ID字段
  - **`constraint:fk_profiles_user_id`**：外键约束名称
  - **`OnUpdate:CASCADE`**：用户ID更新时级联更新
  - **`OnDelete:CASCADE`**：用户删除时级联删除资料

### Posts 关联 📝
```go
Posts []*Post `json:"posts,omitempty" gorm:"foreignKey:UserID;references:ID;constraint:fk_posts_user_id,OnUpdate:CASCADE,OnDelete:CASCADE;"`
```

#### 标签解析：
- **`[]*Post`**：一对多关系，一个用户可以有多篇文章
- **外键配置**：类似Profile，但是一对多关系

## 📖 GORM标签完整参考

### 🏗️ 字段定义标签

| 标签 | 作用 | 示例 | 说明 |
|------|------|------|------|
| `column` | 指定列名 | `gorm:"column:user_name"` | 数据库列名与字段名不同时使用 |
| `type` | 指定数据类型 | `gorm:"type:varchar(100)"` | 自定义数据库字段类型 |
| `size` | 字段大小 | `gorm:"size:255"` | 字符串字段的最大长度 |
| `primaryKey` | 主键 | `gorm:"primaryKey"` | 设置为主键字段 |
| `unique` | 唯一约束 | `gorm:"unique"` | 字段值必须唯一 |
| `default` | 默认值 | `gorm:"default:0"` | 字段的默认值 |
| `precision` | 精度 | `gorm:"precision:10"` | 数值字段的精度 |
| `scale` | 小数位数 | `gorm:"scale:2"` | 数值字段的小数位数 |
| `not null` | 非空约束 | `gorm:"not null"` | 字段不能为空 |
| `autoIncrement` | 自增 | `gorm:"autoIncrement"` | 自动递增（通常用于ID） |
| `autoCreateTime` | 自动创建时间 | `gorm:"autoCreateTime"` | 创建记录时自动设置时间 |
| `autoUpdateTime` | 自动更新时间 | `gorm:"autoUpdateTime"` | 更新记录时自动设置时间 |

### 🔍 索引标签

| 标签 | 作用 | 示例 | 说明 |
|------|------|------|------|
| `index` | 普通索引 | `gorm:"index"` | 提高查询性能 |
| `uniqueIndex` | 唯一索引 | `gorm:"uniqueIndex:idx_name"` | 唯一性约束+索引 |
| `index:,composite` | 复合索引 | `gorm:"index:idx_name,composite:name_age"` | 多字段组合索引 |

### 🔗 关联标签

| 标签 | 作用 | 示例 | 说明 |
|------|------|------|------|
| `foreignKey` | 外键字段 | `gorm:"foreignKey:UserID"` | 指定外键字段名 |
| `references` | 引用字段 | `gorm:"references:ID"` | 引用的主表字段 |
| `constraint` | 约束设置 | `gorm:"constraint:OnDelete:CASCADE"` | 外键约束行为 |
| `many2many` | 多对多关系 | `gorm:"many2many:user_roles"` | 指定中间表名 |
| `polymorphic` | 多态关联 | `gorm:"polymorphic:Owner"` | 多态关联设置 |

### 🚫 序列化标签

| 标签 | 作用 | 示例 | 说明 |
|------|------|------|------|
| `-` | 忽略字段 | `gorm:"-"` | GORM完全忽略此字段 |
| `-:migration` | 忽略迁移 | `gorm:"-:migration"` | 迁移时忽略 |
| `-:all` | 忽略所有 | `gorm:"-:all"` | 所有操作都忽略 |

## 🧠 记忆技巧和最佳实践

### 🎯 记忆口诀

1. **JSON标签**："JSON见名知意" - `json:"字段名"`
2. **GORM标签**："数据库三要素" - 类型(size)、约束(not null)、索引(index)
3. **验证标签**："前端后端双保险" - `validate:"规则"`
4. **关联关系**："外键引用约束" - `foreignKey` + `references` + `constraint`

### 📝 常用组合模板

#### 🔑 主键字段
```go
ID uint `json:"id" gorm:"primaryKey;autoIncrement"`
```

#### 📧 邮箱字段
```go
Email string `json:"email" gorm:"size:100;uniqueIndex;not null" validate:"required,email"`
```

#### 🔐 密码字段
```go
Password string `json:"-" gorm:"size:255;not null" validate:"required,min=6"`
```

#### ⏰ 时间字段
```go
CreatedAt time.Time `json:"created_at" gorm:"autoCreateTime"`
UpdatedAt time.Time `json:"updated_at" gorm:"autoUpdateTime"`
```

#### 🔗 外键关联
```go
UserID uint `json:"user_id" gorm:"not null;index"`
User   User `json:"user" gorm:"foreignKey:UserID;references:ID"`
```

### 🎨 命名规范

1. **索引命名**：`idx_表名_字段名`
   ```go
   gorm:"uniqueIndex:idx_user_email"
   ```

2. **外键约束命名**：`fk_表名_字段名`
   ```go
   gorm:"constraint:fk_posts_user_id"
   ```

3. **JSON字段命名**：使用snake_case
   ```go
   json:"created_at"
   ```

### 🔍 调试技巧

#### 1. 查看生成的SQL
```go
// 开启SQL日志
db.Logger = logger.Default.LogMode(logger.Info)

// 查看建表SQL
db.Migrator().CreateTable(&User{})
```

#### 2. 验证约束
```go
// 测试唯一约束
user1 := User{Email: "test@example.com"}
user2 := User{Email: "test@example.com"} // 应该失败
```

#### 3. 检查关联关系
```go
// 预加载关联数据
var user User
db.Preload("Profile").Preload("Posts").First(&user, 1)
```

### 🚀 性能优化建议

1. **合理使用索引**
   - 经常查询的字段加索引
   - 避免过多索引影响写入性能

2. **字段长度优化**
   - 根据实际需求设置合适的size
   - 避免过大的字段长度

3. **关联查询优化**
   - 使用Preload预加载
   - 避免N+1查询问题

## 🎉 总结

GORM的结构体标签是连接Go结构体和数据库的桥梁，通过合理使用这些标签，我们可以：

✅ **定义清晰的数据库结构**  
✅ **确保数据完整性和一致性**  
✅ **优化查询性能**  
✅ **简化API开发**  
✅ **提高代码可维护性**  

记住：**标签是声明式的配置，一次定义，处处生效！** 🎊

### 🔗 相关资源

- [GORM官方文档](https://gorm.io/docs/)
- [Go结构体标签详解](https://golang.org/ref/spec#Struct_types)
- [数据库设计最佳实践](https://www.postgresql.org/docs/current/ddl.html)