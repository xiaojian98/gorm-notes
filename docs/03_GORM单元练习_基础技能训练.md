# GORM单元练习：基础技能训练

## 📚 练习说明

本练习集分为5个难度等级，每个等级包含多个练习题。建议按顺序完成，每完成一个等级后再进入下一个等级。每个练习都提供了详细的要求、提示和参考答案。

### 🎯 练习目标
- 掌握GORM的基本操作
- 理解模型定义和关联关系
- 熟练使用查询和更新功能
- 学会事务处理和性能优化
- 培养良好的编程习惯

### 📋 练习环境准备
```bash
# 创建练习项目
mkdir gorm-exercises
cd gorm-exercises
go mod init gorm-exercises

# 安装依赖
go get -u gorm.io/gorm
go get -u gorm.io/driver/sqlite  # 使用SQLite简化环境
```

---

## 🌟 Level 1: 基础入门（必须掌握）

### 练习1.1：数据库连接和基本配置

**目标**：学会初始化GORM连接和基本配置

**要求**：
1. 创建一个函数`InitDB()`，连接到SQLite数据库
2. 配置日志级别为Info
3. 启用预编译语句
4. 添加错误处理

**代码框架**：
```go
package main

import (
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// TODO: 完成InitDB函数
func InitDB() (*gorm.DB, error) {
	// 你的代码在这里
}

func main() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	
	// 测试连接
	sqlDB, _ := db.DB()
	defer sqlDB.Close()
	
	println("数据库连接成功！")
}
```

**提示**：
- 使用`sqlite.Open("test.db")`连接SQLite
- 在`gorm.Config`中设置Logger和PrepareStmt
- 记得检查连接是否成功

<details>
<summary>🔍 参考答案</summary>

```go
func InitDB() (*gorm.DB, error) {
	db, err := gorm.Open(sqlite.Open("test.db"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
		PrepareStmt: true,
	})
	
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %w", err)
	}
	
	// 测试连接
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}
	
	if err := sqlDB.Ping(); err != nil {
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}
	
	return db, nil
}
```
</details>

---

### 练习1.2：定义基本模型

**目标**：学会定义GORM模型和使用标签

**要求**：
1. 定义一个`Student`模型，包含以下字段：
   - ID（主键，自增）
   - Name（姓名，最大50字符，不能为空）
   - Email（邮箱，唯一索引）
   - Age（年龄）
   - CreatedAt（创建时间）
   - UpdatedAt（更新时间）
2. 使用适当的GORM标签
3. 实现自定义表名`students`

**代码框架**：
```go
package main

import (
	"time"
	"gorm.io/gorm"
)

// TODO: 定义Student模型
type Student struct {
	// 你的代码在这里
}

// TODO: 实现TableName方法
func (Student) TableName() string {
	// 你的代码在这里
}
```

**提示**：
- 使用`gorm:"primarykey"`定义主键
- 使用`gorm:"size:50;not null"`限制字段长度和非空
- 使用`gorm:"uniqueIndex"`创建唯一索引

<details>
<summary>🔍 参考答案</summary>

```go
type Student struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (Student) TableName() string {
	return "students"
}
```
</details>

---

### 练习1.3：自动迁移和基本CRUD

**目标**：学会数据库迁移和基本的增删改查操作

**要求**：
1. 执行自动迁移创建表
2. 创建3个学生记录
3. 查询所有学生
4. 根据ID查询单个学生
5. 更新学生信息
6. 删除学生记录

**代码框架**：
```go
func main() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	
	// TODO: 1. 自动迁移
	
	// TODO: 2. 创建学生记录
	students := []Student{
		{Name: "张三", Email: "zhang@example.com", Age: 20},
		{Name: "李四", Email: "li@example.com", Age: 21},
		{Name: "王五", Email: "wang@example.com", Age: 19},
	}
	
	// TODO: 3. 查询所有学生
	
	// TODO: 4. 根据ID查询学生
	
	// TODO: 5. 更新学生信息
	
	// TODO: 6. 删除学生记录
}
```

**提示**：
- 使用`db.AutoMigrate(&Student{})`进行迁移
- 使用`db.Create()`创建记录
- 使用`db.Find()`查询所有记录
- 使用`db.First()`查询单个记录
- 使用`db.Model().Updates()`更新记录
- 使用`db.Delete()`删除记录

<details>
<summary>🔍 参考答案</summary>

```go
func main() {
	db, err := InitDB()
	if err != nil {
		panic(err)
	}
	
	// 1. 自动迁移
	if err := db.AutoMigrate(&Student{}); err != nil {
		panic(err)
	}
	
	// 2. 创建学生记录
	students := []Student{
		{Name: "张三", Email: "zhang@example.com", Age: 20},
		{Name: "李四", Email: "li@example.com", Age: 21},
		{Name: "王五", Email: "wang@example.com", Age: 19},
	}
	
	if err := db.Create(&students).Error; err != nil {
		panic(err)
	}
	fmt.Printf("创建了 %d 个学生记录\n", len(students))
	
	// 3. 查询所有学生
	var allStudents []Student
	db.Find(&allStudents)
	fmt.Printf("总共有 %d 个学生\n", len(allStudents))
	
	// 4. 根据ID查询学生
	var student Student
	db.First(&student, 1)
	fmt.Printf("ID为1的学生：%+v\n", student)
	
	// 5. 更新学生信息
	db.Model(&student).Updates(Student{Name: "张三丰", Age: 25})
	fmt.Printf("更新后的学生：%+v\n", student)
	
	// 6. 删除学生记录
	db.Delete(&Student{}, 3)
	fmt.Println("删除了ID为3的学生")
	
	// 验证删除
	var count int64
	db.Model(&Student{}).Count(&count)
	fmt.Printf("剩余学生数量：%d\n", count)
}
```
</details>

---

## 🚀 Level 2: 查询进阶（重要技能）

### 练习2.1：条件查询和链式调用

**目标**：掌握WHERE条件查询和链式调用

**要求**：
1. 查询年龄大于20的学生
2. 查询姓名包含"张"的学生
3. 查询年龄在18-22之间的学生
4. 查询邮箱以"@example.com"结尾的学生
5. 组合条件：年龄大于19且姓名不为空的学生
6. 使用链式调用实现排序和限制

**代码框架**：
```go
func queryExercises(db *gorm.DB) {
	// TODO: 1. 年龄大于20的学生
	
	// TODO: 2. 姓名包含"张"的学生
	
	// TODO: 3. 年龄在18-22之间的学生
	
	// TODO: 4. 邮箱以"@example.com"结尾的学生
	
	// TODO: 5. 组合条件查询
	
	// TODO: 6. 链式调用：按年龄降序，限制3条记录
}
```

**提示**：
- 使用`db.Where("age > ?", 20)`进行条件查询
- 使用`LIKE`进行模糊查询
- 使用`BETWEEN`或多个条件进行范围查询
- 使用`Order()`和`Limit()`进行排序和限制

<details>
<summary>🔍 参考答案</summary>

```go
func queryExercises(db *gorm.DB) {
	// 1. 年龄大于20的学生
	var students1 []Student
	db.Where("age > ?", 20).Find(&students1)
	fmt.Printf("年龄大于20的学生：%d个\n", len(students1))
	
	// 2. 姓名包含"张"的学生
	var students2 []Student
	db.Where("name LIKE ?", "%张%").Find(&students2)
	fmt.Printf("姓名包含'张'的学生：%d个\n", len(students2))
	
	// 3. 年龄在18-22之间的学生
	var students3 []Student
	db.Where("age BETWEEN ? AND ?", 18, 22).Find(&students3)
	// 或者使用：db.Where("age >= ? AND age <= ?", 18, 22).Find(&students3)
	fmt.Printf("年龄在18-22之间的学生：%d个\n", len(students3))
	
	// 4. 邮箱以"@example.com"结尾的学生
	var students4 []Student
	db.Where("email LIKE ?", "%@example.com").Find(&students4)
	fmt.Printf("邮箱以'@example.com'结尾的学生：%d个\n", len(students4))
	
	// 5. 组合条件查询
	var students5 []Student
	db.Where("age > ? AND name != ?", 19, "").Find(&students5)
	fmt.Printf("年龄大于19且姓名不为空的学生：%d个\n", len(students5))
	
	// 6. 链式调用：按年龄降序，限制3条记录
	var students6 []Student
	db.Order("age DESC").Limit(3).Find(&students6)
	fmt.Printf("按年龄降序的前3个学生：\n")
	for _, s := range students6 {
		fmt.Printf("  %s - %d岁\n", s.Name, s.Age)
	}
}
```
</details>

---

### 练习2.2：聚合查询和分组

**目标**：学会使用聚合函数和分组查询

**要求**：
1. 统计学生总数
2. 计算平均年龄
3. 查找最大和最小年龄
4. 按年龄分组统计人数
5. 查询年龄大于平均年龄的学生

**代码框架**：
```go
type AgeGroup struct {
	Age   int `json:"age"`
	Count int `json:"count"`
}

type AgeStats struct {
	AvgAge float64 `json:"avg_age"`
	MaxAge int     `json:"max_age"`
	MinAge int     `json:"min_age"`
}

func aggregateExercises(db *gorm.DB) {
	// TODO: 实现各种聚合查询
}
```

**提示**：
- 使用`db.Model().Count()`统计数量
- 使用`db.Model().Select("AVG(age)").Scan()`计算平均值
- 使用`Group()`和`Select()`进行分组查询
- 可以使用子查询实现复杂条件

<details>
<summary>🔍 参考答案</summary>

```go
func aggregateExercises(db *gorm.DB) {
	// 1. 统计学生总数
	var count int64
	db.Model(&Student{}).Count(&count)
	fmt.Printf("学生总数：%d\n", count)
	
	// 2. 计算平均年龄
	var avgAge float64
	db.Model(&Student{}).Select("AVG(age)").Scan(&avgAge)
	fmt.Printf("平均年龄：%.2f\n", avgAge)
	
	// 3. 查找最大和最小年龄
	var stats AgeStats
	db.Model(&Student{}).Select("AVG(age) as avg_age, MAX(age) as max_age, MIN(age) as min_age").Scan(&stats)
	fmt.Printf("年龄统计：平均%.2f，最大%d，最小%d\n", stats.AvgAge, stats.MaxAge, stats.MinAge)
	
	// 4. 按年龄分组统计人数
	var ageGroups []AgeGroup
	db.Model(&Student{}).Select("age, COUNT(*) as count").Group("age").Scan(&ageGroups)
	fmt.Println("按年龄分组统计：")
	for _, group := range ageGroups {
		fmt.Printf("  %d岁：%d人\n", group.Age, group.Count)
	}
	
	// 5. 查询年龄大于平均年龄的学生
	var studentsAboveAvg []Student
	db.Where("age > (?)", db.Model(&Student{}).Select("AVG(age)")).Find(&studentsAboveAvg)
	fmt.Printf("年龄大于平均年龄的学生：%d个\n", len(studentsAboveAvg))
	for _, s := range studentsAboveAvg {
		fmt.Printf("  %s - %d岁\n", s.Name, s.Age)
	}
}
```
</details>

---

### 练习2.3：分页查询

**目标**：实现通用的分页查询功能

**要求**：
1. 实现一个通用的分页结构体
2. 创建分页查询函数
3. 支持排序和搜索的分页
4. 计算总页数和分页信息

**代码框架**：
```go
type Pagination struct {
	Page       int         `json:"page"`
	PageSize   int         `json:"page_size"`
	Total      int64       `json:"total"`
	TotalPages int         `json:"total_pages"`
	Data       interface{} `json:"data"`
}

type StudentQuery struct {
	Page     int    `json:"page"`
	PageSize int    `json:"page_size"`
	Search   string `json:"search"`
	OrderBy  string `json:"order_by"`
}

// TODO: 实现分页查询函数
func GetStudentsWithPagination(db *gorm.DB, query StudentQuery) (*Pagination, error) {
	// 你的代码在这里
}
```

**提示**：
- 使用`Offset()`和`Limit()`实现分页
- 先用`Count()`获取总数，再查询数据
- 计算总页数：`(total + pageSize - 1) / pageSize`
- 支持动态排序和搜索条件

<details>
<summary>🔍 参考答案</summary>

```go
func GetStudentsWithPagination(db *gorm.DB, query StudentQuery) (*Pagination, error) {
	// 设置默认值
	if query.Page <= 0 {
		query.Page = 1
	}
	if query.PageSize <= 0 {
		query.PageSize = 10
	}
	if query.OrderBy == "" {
		query.OrderBy = "created_at DESC"
	}
	
	// 构建查询
	queryDB := db.Model(&Student{})
	
	// 添加搜索条件
	if query.Search != "" {
		searchTerm := "%" + query.Search + "%"
		queryDB = queryDB.Where("name LIKE ? OR email LIKE ?", searchTerm, searchTerm)
	}
	
	// 获取总数
	var total int64
	if err := queryDB.Count(&total).Error; err != nil {
		return nil, err
	}
	
	// 计算偏移量
	offset := (query.Page - 1) * query.PageSize
	
	// 查询数据
	var students []Student
	err := queryDB.Order(query.OrderBy).Limit(query.PageSize).Offset(offset).Find(&students).Error
	if err != nil {
		return nil, err
	}
	
	// 计算总页数
	totalPages := int((total + int64(query.PageSize) - 1) / int64(query.PageSize))
	
	return &Pagination{
		Page:       query.Page,
		PageSize:   query.PageSize,
		Total:      total,
		TotalPages: totalPages,
		Data:       students,
	}, nil
}

// 测试分页功能
func testPagination(db *gorm.DB) {
	// 测试基本分页
	result, err := GetStudentsWithPagination(db, StudentQuery{
		Page:     1,
		PageSize: 2,
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("分页结果：第%d页，每页%d条，共%d条，共%d页\n", 
		result.Page, result.PageSize, result.Total, result.TotalPages)
	
	// 测试搜索分页
	searchResult, err := GetStudentsWithPagination(db, StudentQuery{
		Page:     1,
		PageSize: 10,
		Search:   "张",
		OrderBy:  "age ASC",
	})
	if err != nil {
		panic(err)
	}
	
	fmt.Printf("搜索'张'的结果：共%d条\n", searchResult.Total)
}
```
</details>

---

## 🔗 Level 3: 关联关系（核心技能）

### 练习3.1：一对多关系

**目标**：理解和实现一对多关系

**要求**：
1. 扩展Student模型，添加班级信息
2. 创建Class模型（一个班级有多个学生）
3. 实现关联查询
4. 使用预加载避免N+1问题

**代码框架**：
```go
// TODO: 定义Class模型
type Class struct {
	// 你的代码在这里
}

// TODO: 修改Student模型，添加班级关联
type Student struct {
	// 原有字段...
	// 添加班级关联
}

func relationshipExercises(db *gorm.DB) {
	// TODO: 1. 创建班级和学生数据
	// TODO: 2. 查询班级及其所有学生
	// TODO: 3. 查询学生及其班级信息
	// TODO: 4. 统计每个班级的学生数量
}
```

**提示**：
- 在Student中添加ClassID外键和Class关联
- 在Class中添加Students切片
- 使用`Preload()`预加载关联数据
- 使用`Association()`进行关联操作

<details>
<summary>🔍 参考答案</summary>

```go
// Class 班级模型
type Class struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Teacher   string    `json:"teacher" gorm:"size:50"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 一对多关系：一个班级有多个学生
	Students []Student `json:"students,omitempty" gorm:"foreignKey:ClassID"`
}

// 修改Student模型
type Student struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 外键
	ClassID uint  `json:"class_id" gorm:"index"`
	// 属于关系：学生属于一个班级
	Class   Class `json:"class,omitempty" gorm:"foreignKey:ClassID"`
}

func relationshipExercises(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Class{}, &Student{})
	
	// 1. 创建班级和学生数据
	classes := []Class{
		{Name: "计算机1班", Teacher: "张老师"},
		{Name: "计算机2班", Teacher: "李老师"},
	}
	db.Create(&classes)
	
	students := []Student{
		{Name: "张三", Email: "zhang@example.com", Age: 20, ClassID: classes[0].ID},
		{Name: "李四", Email: "li@example.com", Age: 21, ClassID: classes[0].ID},
		{Name: "王五", Email: "wang@example.com", Age: 19, ClassID: classes[1].ID},
		{Name: "赵六", Email: "zhao@example.com", Age: 22, ClassID: classes[1].ID},
	}
	db.Create(&students)
	
	// 2. 查询班级及其所有学生（预加载）
	var classesWithStudents []Class
	db.Preload("Students").Find(&classesWithStudents)
	
	fmt.Println("班级及其学生：")
	for _, class := range classesWithStudents {
		fmt.Printf("班级：%s（老师：%s）\n", class.Name, class.Teacher)
		for _, student := range class.Students {
			fmt.Printf("  学生：%s - %d岁\n", student.Name, student.Age)
		}
	}
	
	// 3. 查询学生及其班级信息
	var studentsWithClass []Student
	db.Preload("Class").Find(&studentsWithClass)
	
	fmt.Println("\n学生及其班级：")
	for _, student := range studentsWithClass {
		fmt.Printf("学生：%s - 班级：%s\n", student.Name, student.Class.Name)
	}
	
	// 4. 统计每个班级的学生数量
	type ClassStats struct {
		ClassName    string `json:"class_name"`
		StudentCount int    `json:"student_count"`
	}
	
	var stats []ClassStats
	db.Model(&Class{}).
		Select("classes.name as class_name, COUNT(students.id) as student_count").
		Joins("LEFT JOIN students ON classes.id = students.class_id").
		Group("classes.id, classes.name").
		Scan(&stats)
	
	fmt.Println("\n班级学生统计：")
	for _, stat := range stats {
		fmt.Printf("班级：%s - 学生数：%d\n", stat.ClassName, stat.StudentCount)
	}
}
```
</details>

---

### 练习3.2：多对多关系

**目标**：实现多对多关系（学生选课系统）

**要求**：
1. 创建Course模型
2. 实现Student和Course的多对多关系
3. 自定义中间表添加额外字段（成绩、选课时间）
4. 实现选课、退课、查询功能

**代码框架**：
```go
// TODO: 定义Course模型
type Course struct {
	// 你的代码在这里
}

// TODO: 定义中间表模型
type StudentCourse struct {
	// 你的代码在这里
}

// TODO: 修改Student模型，添加课程关联

func manyToManyExercises(db *gorm.DB) {
	// TODO: 实现选课系统功能
}
```

**提示**：
- 使用`gorm:"many2many:student_courses"`定义多对多关系
- 自定义中间表需要定义结构体并使用`SetupJoinTable`
- 使用`Association()`进行关联操作
- 可以通过中间表查询额外信息

<details>
<summary>🔍 参考答案</summary>

```go
// Course 课程模型
type Course struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Code        string    `json:"code" gorm:"size:20;uniqueIndex"`
	Credits     int       `json:"credits"`
	Teacher     string    `json:"teacher" gorm:"size:50"`
	Description string    `json:"description" gorm:"type:text"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	
	// 多对多关系：课程有多个学生
	Students []Student `json:"students,omitempty" gorm:"many2many:student_courses;"`
}

// StudentCourse 学生选课中间表（自定义字段）
type StudentCourse struct {
	StudentID    uint      `json:"student_id" gorm:"primarykey"`
	CourseID     uint      `json:"course_id" gorm:"primarykey"`
	Score        *float64  `json:"score"`        // 成绩，可为空
	EnrolledAt   time.Time `json:"enrolled_at"`  // 选课时间
	Status       string    `json:"status" gorm:"size:20;default:enrolled"` // enrolled, dropped, completed
	
	// 关联关系
	Student Student `json:"student" gorm:"foreignKey:StudentID"`
	Course  Course  `json:"course" gorm:"foreignKey:CourseID"`
}

// 修改Student模型，添加课程关联
type Student struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	Name      string    `json:"name" gorm:"size:50;not null"`
	Email     string    `json:"email" gorm:"size:100;uniqueIndex"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 班级关联
	ClassID uint  `json:"class_id" gorm:"index"`
	Class   Class `json:"class,omitempty" gorm:"foreignKey:ClassID"`
	
	// 多对多关系：学生选择多门课程
	Courses []Course `json:"courses,omitempty" gorm:"many2many:student_courses;"`
}

func manyToManyExercises(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Course{}, &StudentCourse{})
	
	// 设置自定义中间表
	db.SetupJoinTable(&Student{}, "Courses", &StudentCourse{})
	db.SetupJoinTable(&Course{}, "Students", &StudentCourse{})
	
	// 1. 创建课程数据
	courses := []Course{
		{Name: "数据结构", Code: "CS101", Credits: 3, Teacher: "王教授"},
		{Name: "算法设计", Code: "CS102", Credits: 4, Teacher: "李教授"},
		{Name: "数据库原理", Code: "CS201", Credits: 3, Teacher: "张教授"},
	}
	db.Create(&courses)
	
	// 2. 学生选课
	var student Student
	db.First(&student, 1)
	
	// 选择课程1和课程2
	enrollments := []StudentCourse{
		{
			StudentID:  student.ID,
			CourseID:   courses[0].ID,
			EnrolledAt: time.Now(),
			Status:     "enrolled",
		},
		{
			StudentID:  student.ID,
			CourseID:   courses[1].ID,
			EnrolledAt: time.Now(),
			Status:     "enrolled",
		},
	}
	db.Create(&enrollments)
	
	// 3. 查询学生的所有课程
	var studentWithCourses Student
	db.Preload("Courses").First(&studentWithCourses, student.ID)
	
	fmt.Printf("学生 %s 选择的课程：\n", studentWithCourses.Name)
	for _, course := range studentWithCourses.Courses {
		fmt.Printf("  %s (%s) - %d学分\n", course.Name, course.Code, course.Credits)
	}
	
	// 4. 查询课程的所有学生
	var courseWithStudents Course
	db.Preload("Students").First(&courseWithStudents, courses[0].ID)
	
	fmt.Printf("\n课程 %s 的学生：\n", courseWithStudents.Name)
	for _, s := range courseWithStudents.Students {
		fmt.Printf("  %s\n", s.Name)
	}
	
	// 5. 通过中间表查询详细信息
	var enrollmentDetails []StudentCourse
	db.Preload("Student").Preload("Course").Where("student_id = ?", student.ID).Find(&enrollmentDetails)
	
	fmt.Printf("\n学生 %s 的选课详情：\n", student.Name)
	for _, detail := range enrollmentDetails {
		fmt.Printf("  课程：%s，选课时间：%s，状态：%s\n", 
			detail.Course.Name, 
			detail.EnrolledAt.Format("2006-01-02 15:04:05"), 
			detail.Status)
	}
	
	// 6. 更新成绩
	db.Model(&StudentCourse{}).Where("student_id = ? AND course_id = ?", student.ID, courses[0].ID).Update("score", 85.5)
	
	// 7. 退课
	db.Model(&StudentCourse{}).Where("student_id = ? AND course_id = ?", student.ID, courses[1].ID).Update("status", "dropped")
	
	// 8. 统计每门课程的选课人数
	type CourseStats struct {
		CourseName   string `json:"course_name"`
		EnrolledCount int   `json:"enrolled_count"`
	}
	
	var courseStats []CourseStats
	db.Model(&Course{}).
		Select("courses.name as course_name, COUNT(student_courses.student_id) as enrolled_count").
		Joins("LEFT JOIN student_courses ON courses.id = student_courses.course_id AND student_courses.status = 'enrolled'").
		Group("courses.id, courses.name").
		Scan(&courseStats)
	
	fmt.Println("\n课程选课统计：")
	for _, stat := range courseStats {
		fmt.Printf("课程：%s - 选课人数：%d\n", stat.CourseName, stat.EnrolledCount)
	}
}
```
</details>

---

## ⚡ Level 4: 高级特性（进阶技能）

### 练习4.1：事务处理

**目标**：掌握GORM的事务处理机制

**要求**：
1. 实现自动事务（使用Transaction方法）
2. 实现手动事务控制
3. 处理事务回滚场景
4. 实现嵌套事务
5. 创建一个转账系统演示事务的重要性

**代码框架**：
```go
// TODO: 定义账户模型
type Account struct {
	// 你的代码在这里
}

// TODO: 实现转账功能（需要事务保证一致性）
func Transfer(db *gorm.DB, fromID, toID uint, amount float64) error {
	// 你的代码在这里
}

func transactionExercises(db *gorm.DB) {
	// TODO: 测试各种事务场景
}
```

**提示**：
- 使用`db.Transaction()`进行自动事务
- 使用`db.Begin()`、`tx.Commit()`、`tx.Rollback()`进行手动事务
- 在转账中需要检查余额、更新两个账户
- 任何步骤失败都应该回滚整个事务

<details>
<summary>🔍 参考答案</summary>

```go
// Account 账户模型
type Account struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	UserID    uint      `json:"user_id" gorm:"not null;index"`
	Balance   float64   `json:"balance" gorm:"type:decimal(10,2);default:0"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	
	// 关联用户
	User Student `json:"user" gorm:"foreignKey:UserID"`
}

// TransferRecord 转账记录
type TransferRecord struct {
	ID       uint      `json:"id" gorm:"primarykey"`
	FromID   uint      `json:"from_id" gorm:"not null"`
	ToID     uint      `json:"to_id" gorm:"not null"`
	Amount   float64   `json:"amount" gorm:"type:decimal(10,2);not null"`
	Status   string    `json:"status" gorm:"size:20;default:pending"` // pending, success, failed
	Remark   string    `json:"remark" gorm:"size:200"`
	CreatedAt time.Time `json:"created_at"`
}

// Transfer 转账功能（自动事务）
func Transfer(db *gorm.DB, fromID, toID uint, amount float64, remark string) error {
	if amount <= 0 {
		return errors.New("转账金额必须大于0")
	}
	
	return db.Transaction(func(tx *gorm.DB) error {
		// 1. 创建转账记录
		record := TransferRecord{
			FromID: fromID,
			ToID:   toID,
			Amount: amount,
			Remark: remark,
			Status: "pending",
		}
		if err := tx.Create(&record).Error; err != nil {
			return err
		}
		
		// 2. 检查转出账户余额
		var fromAccount Account
		if err := tx.First(&fromAccount, fromID).Error; err != nil {
			return errors.New("转出账户不存在")
		}
		
		if fromAccount.Balance < amount {
			return errors.New("余额不足")
		}
		
		// 3. 检查转入账户是否存在
		var toAccount Account
		if err := tx.First(&toAccount, toID).Error; err != nil {
			return errors.New("转入账户不存在")
		}
		
		// 4. 更新转出账户余额
		if err := tx.Model(&fromAccount).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
			return err
		}
		
		// 5. 更新转入账户余额
		if err := tx.Model(&toAccount).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
			return err
		}
		
		// 6. 更新转账记录状态
		if err := tx.Model(&record).Update("status", "success").Error; err != nil {
			return err
		}
		
		return nil
	})
}

// TransferManual 手动事务控制的转账
func TransferManual(db *gorm.DB, fromID, toID uint, amount float64) error {
	// 开始事务
	tx := db.Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()
	
	if tx.Error != nil {
		return tx.Error
	}
	
	// 检查余额
	var fromAccount Account
	if err := tx.First(&fromAccount, fromID).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if fromAccount.Balance < amount {
		tx.Rollback()
		return errors.New("余额不足")
	}
	
	// 更新账户余额
	if err := tx.Model(&Account{}).Where("id = ?", fromID).Update("balance", gorm.Expr("balance - ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	if err := tx.Model(&Account{}).Where("id = ?", toID).Update("balance", gorm.Expr("balance + ?", amount)).Error; err != nil {
		tx.Rollback()
		return err
	}
	
	// 提交事务
	return tx.Commit().Error
}

func transactionExercises(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Account{}, &TransferRecord{})
	
	// 创建测试账户
	accounts := []Account{
		{UserID: 1, Balance: 1000.00},
		{UserID: 2, Balance: 500.00},
		{UserID: 3, Balance: 200.00},
	}
	db.Create(&accounts)
	
	fmt.Println("初始账户余额：")
	var allAccounts []Account
	db.Find(&allAccounts)
	for _, acc := range allAccounts {
		fmt.Printf("账户%d：%.2f元\n", acc.ID, acc.Balance)
	}
	
	// 测试成功转账
	fmt.Println("\n测试转账：账户1向账户2转账100元")
	err := Transfer(db, accounts[0].ID, accounts[1].ID, 100.00, "测试转账")
	if err != nil {
		fmt.Printf("转账失败：%v\n", err)
	} else {
		fmt.Println("转账成功")
	}
	
	// 查看转账后余额
	db.Find(&allAccounts)
	fmt.Println("转账后余额：")
	for _, acc := range allAccounts {
		fmt.Printf("账户%d：%.2f元\n", acc.ID, acc.Balance)
	}
	
	// 测试余额不足的转账
	fmt.Println("\n测试余额不足：账户3向账户1转账500元")
	err = Transfer(db, accounts[2].ID, accounts[0].ID, 500.00, "余额不足测试")
	if err != nil {
		fmt.Printf("转账失败：%v\n", err)
	} else {
		fmt.Println("转账成功")
	}
	
	// 查看转账记录
	var records []TransferRecord
	db.Find(&records)
	fmt.Println("\n转账记录：")
	for _, record := range records {
		fmt.Printf("从账户%d到账户%d：%.2f元，状态：%s\n", 
			record.FromID, record.ToID, record.Amount, record.Status)
	}
	
	// 测试嵌套事务
	fmt.Println("\n测试嵌套事务")
	db.Transaction(func(tx *gorm.DB) error {
		fmt.Println("外层事务开始")
		
		// 内层事务
		tx.Transaction(func(tx2 *gorm.DB) error {
			fmt.Println("内层事务：创建新账户")
			newAccount := Account{UserID: 4, Balance: 300.00}
			return tx2.Create(&newAccount).Error
		})
		
		fmt.Println("外层事务：转账操作")
		return Transfer(tx, accounts[0].ID, 4, 50.00, "嵌套事务测试")
	})
	
	fmt.Println("嵌套事务完成")
}
```
</details>

---

### 练习4.2：钩子函数和回调

**目标**：掌握GORM的钩子函数机制

**要求**：
1. 实现所有类型的钩子函数
2. 在钩子中实现业务逻辑（如审计日志、数据验证）
3. 理解钩子的执行顺序
4. 实现条件钩子（某些情况下跳过钩子）

**代码框架**：
```go
// TODO: 定义带有完整钩子的模型
type Product struct {
	// 你的代码在这里
}

// TODO: 定义审计日志模型
type AuditLog struct {
	// 你的代码在这里
}

// TODO: 实现各种钩子函数

func hookExercises(db *gorm.DB) {
	// TODO: 测试钩子函数
}
```

**提示**：
- 实现BeforeCreate、AfterCreate、BeforeUpdate、AfterUpdate等钩子
- 在钩子中可以修改数据、记录日志、发送通知等
- 钩子返回error会中断操作
- 可以通过tx.Statement.Changed()检查字段是否被修改

<details>
<summary>🔍 参考答案</summary>

```go
// Product 产品模型（包含完整钩子）
type Product struct {
	ID          uint      `json:"id" gorm:"primarykey"`
	Name        string    `json:"name" gorm:"size:100;not null"`
	Price       float64   `json:"price" gorm:"type:decimal(10,2);not null"`
	Stock       int       `json:"stock" gorm:"default:0"`
	Status      string    `json:"status" gorm:"size:20;default:active"` // active, inactive, deleted
	Description string    `json:"description" gorm:"type:text"`
	CreatedBy   uint      `json:"created_by"`
	UpdatedBy   uint      `json:"updated_by"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	DeletedAt   *time.Time `json:"deleted_at" gorm:"index"`
}

// AuditLog 审计日志模型
type AuditLog struct {
	ID        uint      `json:"id" gorm:"primarykey"`
	TableName string    `json:"table_name" gorm:"size:50;not null"`
	RecordID  uint      `json:"record_id" gorm:"not null"`
	Action    string    `json:"action" gorm:"size:20;not null"` // create, update, delete
	OldData   string    `json:"old_data" gorm:"type:json"`
	NewData   string    `json:"new_data" gorm:"type:json"`
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

// BeforeCreate 创建前钩子
func (p *Product) BeforeCreate(tx *gorm.DB) error {
	fmt.Printf("[钩子] BeforeCreate: 准备创建产品 %s\n", p.Name)
	
	// 数据验证
	if p.Price < 0 {
		return errors.New("产品价格不能为负数")
	}
	
	// 设置默认值
	if p.Status == "" {
		p.Status = "active"
	}
	
	// 从上下文获取用户ID（模拟）
	if userID, exists := tx.Get("current_user_id"); exists {
		p.CreatedBy = userID.(uint)
	}
	
	return nil
}

// AfterCreate 创建后钩子
func (p *Product) AfterCreate(tx *gorm.DB) error {
	fmt.Printf("[钩子] AfterCreate: 产品 %s 创建成功，ID: %d\n", p.Name, p.ID)
	
	// 记录审计日志
	newData, _ := json.Marshal(p)
	auditLog := AuditLog{
		TableName: "products",
		RecordID:  p.ID,
		Action:    "create",
		NewData:   string(newData),
		UserID:    p.CreatedBy,
	}
	
	return tx.Create(&auditLog).Error
}

// BeforeUpdate 更新前钩子
func (p *Product) BeforeUpdate(tx *gorm.DB) error {
	fmt.Printf("[钩子] BeforeUpdate: 准备更新产品 ID: %d\n", p.ID)
	
	// 检查哪些字段被修改了
	if tx.Statement.Changed("Price") {
		fmt.Println("  价格被修改")
		// 价格修改需要特殊权限验证
		if p.Price < 0 {
			return errors.New("价格不能为负数")
		}
	}
	
	if tx.Statement.Changed("Stock") {
		fmt.Println("  库存被修改")
	}
	
	// 设置更新者
	if userID, exists := tx.Get("current_user_id"); exists {
		p.UpdatedBy = userID.(uint)
	}
	
	return nil
}

// AfterUpdate 更新后钩子
func (p *Product) AfterUpdate(tx *gorm.DB) error {
	fmt.Printf("[钩子] AfterUpdate: 产品 ID: %d 更新完成\n", p.ID)
	
	// 获取更新前的数据
	var oldProduct Product
	tx.Unscoped().Where("id = ?", p.ID).First(&oldProduct)
	
	// 记录审计日志
	oldData, _ := json.Marshal(oldProduct)
	newData, _ := json.Marshal(p)
	auditLog := AuditLog{
		TableName: "products",
		RecordID:  p.ID,
		Action:    "update",
		OldData:   string(oldData),
		NewData:   string(newData),
		UserID:    p.UpdatedBy,
	}
	
	return tx.Create(&auditLog).Error
}

// BeforeDelete 删除前钩子
func (p *Product) BeforeDelete(tx *gorm.DB) error {
	fmt.Printf("[钩子] BeforeDelete: 准备删除产品 ID: %d\n", p.ID)
	
	// 检查是否可以删除
	if p.Stock > 0 {
		return errors.New("有库存的产品不能删除")
	}
	
	return nil
}

// AfterDelete 删除后钩子
func (p *Product) AfterDelete(tx *gorm.DB) error {
	fmt.Printf("[钩子] AfterDelete: 产品 ID: %d 删除完成\n", p.ID)
	
	// 记录审计日志
	oldData, _ := json.Marshal(p)
	auditLog := AuditLog{
		TableName: "products",
		RecordID:  p.ID,
		Action:    "delete",
		OldData:   string(oldData),
		UserID:    p.UpdatedBy,
	}
	
	return tx.Create(&auditLog).Error
}

// AfterFind 查询后钩子
func (p *Product) AfterFind(tx *gorm.DB) error {
	// 可以在这里进行数据后处理
	if p.Status == "" {
		p.Status = "unknown"
	}
	return nil
}

func hookExercises(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Product{}, &AuditLog{})
	
	// 设置当前用户ID（模拟登录用户）
	db = db.Set("current_user_id", uint(1))
	
	fmt.Println("=== 测试创建产品（触发BeforeCreate和AfterCreate钩子）===")
	product := Product{
		Name:        "iPhone 15",
		Price:       5999.00,
		Stock:       10,
		Description: "最新款iPhone",
	}
	
	if err := db.Create(&product).Error; err != nil {
		fmt.Printf("创建失败：%v\n", err)
	} else {
		fmt.Printf("产品创建成功，ID: %d\n", product.ID)
	}
	
	fmt.Println("\n=== 测试更新产品（触发BeforeUpdate和AfterUpdate钩子）===")
	db.Model(&product).Updates(Product{Price: 5799.00, Stock: 8})
	
	fmt.Println("\n=== 测试查询产品（触发AfterFind钩子）===")
	var foundProduct Product
	db.First(&foundProduct, product.ID)
	fmt.Printf("查询到产品：%+v\n", foundProduct)
	
	fmt.Println("\n=== 测试删除限制（BeforeDelete钩子阻止删除）===")
	if err := db.Delete(&product).Error; err != nil {
		fmt.Printf("删除失败：%v\n", err)
	}
	
	fmt.Println("\n=== 清空库存后删除===")
	db.Model(&product).Update("stock", 0)
	if err := db.Delete(&product).Error; err != nil {
		fmt.Printf("删除失败：%v\n", err)
	} else {
		fmt.Println("产品删除成功")
	}
	
	fmt.Println("\n=== 查看审计日志 ===")
	var auditLogs []AuditLog
	db.Where("table_name = ? AND record_id = ?", "products", product.ID).Find(&auditLogs)
	for _, log := range auditLogs {
		fmt.Printf("操作：%s，时间：%s\n", log.Action, log.CreatedAt.Format("2006-01-02 15:04:05"))
	}
}
```
</details>

---

## 🎯 Level 5: 性能优化（专家技能）

### 练习5.1：查询优化和索引

**目标**：学会分析和优化GORM查询性能

**要求**：
1. 使用EXPLAIN分析查询计划
2. 创建和使用索引
3. 优化N+1查询问题
4. 使用原生SQL优化复杂查询
5. 实现查询缓存

**代码框架**：
```go
// TODO: 创建性能测试模型
type Order struct {
	// 你的代码在这里
}

type OrderItem struct {
	// 你的代码在这里
}

// TODO: 实现性能优化函数
func optimizationExercises(db *gorm.DB) {
	// 你的代码在这里
}
```

**提示**：
- 使用`db.Debug()`查看生成的SQL
- 使用复合索引优化多字段查询
- 使用`Preload`和`Joins`避免N+1问题
- 对比不同查询方式的性能

<details>
<summary>🔍 参考答案</summary>

```go
// Order 订单模型
type Order struct {
	ID         uint      `json:"id" gorm:"primarykey"`
	OrderNo    string    `json:"order_no" gorm:"size:32;uniqueIndex"`
	CustomerID uint      `json:"customer_id" gorm:"index"`
	Status     string    `json:"status" gorm:"size:20;index"`
	TotalAmount float64  `json:"total_amount" gorm:"type:decimal(10,2)"`
	CreatedAt  time.Time `json:"created_at" gorm:"index"`
	UpdatedAt  time.Time `json:"updated_at"`
	
	// 关联关系
	Customer   Student     `json:"customer" gorm:"foreignKey:CustomerID"`
	OrderItems []OrderItem `json:"order_items" gorm:"foreignKey:OrderID"`
}

// OrderItem 订单项模型
type OrderItem struct {
	ID        uint    `json:"id" gorm:"primarykey"`
	OrderID   uint    `json:"order_id" gorm:"index"`
	ProductID uint    `json:"product_id" gorm:"index"`
	Quantity  int     `json:"quantity"`
	Price     float64 `json:"price" gorm:"type:decimal(10,2)"`
	
	// 关联关系
	Order   Order   `json:"order" gorm:"foreignKey:OrderID"`
	Product Product `json:"product" gorm:"foreignKey:ProductID"`
}

// 添加复合索引
func (Order) TableName() string {
	return "orders"
}

func (OrderItem) TableName() string {
	return "order_items"
}

func optimizationExercises(db *gorm.DB) {
	// 自动迁移
	db.AutoMigrate(&Order{}, &OrderItem{})
	
	// 创建复合索引
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_customer_status ON orders(customer_id, status)")
	db.Exec("CREATE INDEX IF NOT EXISTS idx_orders_created_status ON orders(created_at, status)")
	
	// 创建测试数据
	createTestData(db)
	
	fmt.Println("=== 1. 查询性能对比 ===")
	
	// 不使用索引的查询
	start := time.Now()
	var orders1 []Order
	db.Where("total_amount > ?", 1000).Find(&orders1)
	fmt.Printf("无索引查询耗时：%v，结果数：%d\n", time.Since(start), len(orders1))
	
	// 使用索引的查询
	start = time.Now()
	var orders2 []Order
	db.Where("status = ?", "completed").Find(&orders2)
	fmt.Printf("有索引查询耗时：%v，结果数：%d\n", time.Since(start), len(orders2))
	
	fmt.Println("\n=== 2. N+1查询问题演示 ===")
	
	// 错误方式：N+1查询
	start = time.Now()
	var ordersN1 []Order
	db.Limit(5).Find(&ordersN1)
	for i := range ordersN1 {
		db.Where("order_id = ?", ordersN1[i].ID).Find(&ordersN1[i].OrderItems)
	}
	fmt.Printf("N+1查询耗时：%v\n", time.Since(start))
	
	// 正确方式：预加载
	start = time.Now()
	var ordersPreload []Order
	db.Preload("OrderItems").Limit(5).Find(&ordersPreload)
	fmt.Printf("预加载查询耗时：%v\n", time.Since(start))
	
	fmt.Println("\n=== 3. 使用Joins优化 ===")
	
	// 使用Joins进行关联查询
	start = time.Now()
	var ordersJoin []Order
	db.Joins("Customer").Where("customers.age > ?", 20).Find(&ordersJoin)
	fmt.Printf("Joins查询耗时：%v，结果数：%d\n", time.Since(start), len(ordersJoin))
	
	fmt.Println("\n=== 4. 原生SQL优化复杂查询 ===")
	
	// 复杂统计查询使用原生SQL
	type OrderStats struct {
		CustomerID   uint    `json:"customer_id"`
		CustomerName string  `json:"customer_name"`
		OrderCount   int     `json:"order_count"`
		TotalAmount  float64 `json:"total_amount"`
	}
	
	start = time.Now()
	var stats []OrderStats
	db.Raw(`
		SELECT 
			o.customer_id,
			s.name as customer_name,
			COUNT(o.id) as order_count,
			SUM(o.total_amount) as total_amount
		FROM orders o
		JOIN students s ON o.customer_id = s.id
		WHERE o.status = 'completed'
		GROUP BY o.customer_id, s.name
		HAVING COUNT(o.id) > 1
		ORDER BY total_amount DESC
		LIMIT 10
	`).Scan(&stats)
	fmt.Printf("原生SQL统计查询耗时：%v，结果数：%d\n", time.Since(start), len(stats))
	
	fmt.Println("\n=== 5. 查询计划分析 ===")
	
	// 分析查询计划
	var explain []map[string]interface{}
	db.Raw("EXPLAIN QUERY PLAN SELECT * FROM orders WHERE status = 'completed' AND customer_id = 1").Scan(&explain)
	fmt.Println("查询计划：")
	for _, row := range explain {
		fmt.Printf("%+v\n", row)
	}
}

// 创建测试数据
func createTestData(db *gorm.DB) {
	// 创建订单数据
	orders := make([]Order, 100)
	for i := 0; i < 100; i++ {
		orders[i] = Order{
			OrderNo:     fmt.Sprintf("ORD%06d", i+1),
			CustomerID:  uint(i%10 + 1),
			Status:      []string{"pending", "completed", "cancelled"}[i%3],
			TotalAmount: float64(100 + i*10),
			CreatedAt:   time.Now().AddDate(0, 0, -i),
		}
	}
	db.CreateInBatches(&orders, 50)
	
	// 创建订单项数据
	var orderItems []OrderItem
	for _, order := range orders {
		itemCount := 1 + order.ID%3 // 每个订单1-3个商品
		for j := uint(0); j < itemCount; j++ {
			orderItems = append(orderItems, OrderItem{
				OrderID:   order.ID,
				ProductID: j%5 + 1,
				Quantity:  int(j + 1),
				Price:     float64(50 + j*20),
			})
		}
	}
	db.CreateInBatches(&orderItems, 100)
}
```
</details>

---

## 📝 练习总结

### 🎯 学习检查清单

完成所有练习后，你应该能够：

**Level 1 - 基础入门**
- [ ] 正确配置和连接GORM数据库
- [ ] 定义模型和使用标签
- [ ] 执行基本的CRUD操作
- [ ] 理解自动迁移机制

**Level 2 - 查询进阶**
- [ ] 使用各种WHERE条件查询
- [ ] 实现聚合查询和分组
- [ ] 构建通用分页功能
- [ ] 掌握链式调用技巧

**Level 3 - 关联关系**
- [ ] 实现一对多关系
- [ ] 实现多对多关系
- [ ] 使用预加载避免N+1问题
- [ ] 通过中间表查询数据

**Level 4 - 高级特性**
- [ ] 正确使用事务保证数据一致性
- [ ] 实现完整的钩子函数
- [ ] 理解钩子执行顺序
- [ ] 在钩子中实现业务逻辑

**Level 5 - 性能优化**
- [ ] 分析和优化查询性能
- [ ] 正确使用索引
- [ ] 避免常见性能陷阱
- [ ] 使用原生SQL处理复杂查询

### 🚀 下一步学习建议

1. **实战项目**：将所学知识应用到实际项目中
2. **源码阅读**：深入研究GORM源码，理解实现原理
3. **性能调优**：在生产环境中进行性能监控和优化
4. **插件开发**：尝试开发GORM插件扩展功能

### 📚 参考资源

- [GORM官方文档](https://gorm.io/docs/)
- [Go数据库编程最佳实践](https://go.dev/doc/database/)
- [SQL性能优化指南](https://use-the-index-luke.com/)

---

**恭喜你完成了GORM单元练习！继续保持学习的热情，在实践中不断提升技能！** 🎉