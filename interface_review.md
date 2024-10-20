# “意语日记”项目对接复盘
## 项目规范
### 空值不返回：
1. 方法：在json注释中加入omitempty
2. 问题： gorm中封装了字段，需要重新声明字段，然后加入omitempty吗？
3. 解决：用DTO与VO思路，进行数据转换
### DTO与VO 入参出参：
1. 概述
   
(1)dto
+ dto保证服务端数据传输和客户端响应只传递必要数据Content和TagIds
+ 定义DiaryDto 代表日记的请求模型，定义 DiaryToDto 将Diary模型转换为DiaryDto
+ CreateDiary函数使用了dto，把用户请求的数据绑定到diaryDto，只需要检查dto中的Content和TagIds是否存在。
+ 创建日记记录时，只需要提供content、tag和Id外键约束，调用DiaryToDto数据转换为dto

(2)vo
+ vo保证前端只输出必要信息，空值不返回
+ 定义DiaryVo 用于API响应的日记信息，定义Copy 直接实现赋值
+ GetDiaries中调用copy方法将绑定的数据转换为vo格式
+ 响应时，返回diaryVos，保证前端输出的是vo格式
2. 问题
+ 响应返回格式依然是models
+ 改为vo格式后，没有分页信息
3. 解决
+ 响应之前进行数据转换，用vo的copy方法将请求信息转换为vo格式，响应参数应该为vo实例
+ 目前方法是暴力copy，在vo文件中，创建结构体PaginatedDiaryVo包含DiaryVo的数据和分页信息两部分信息。实现copy方法，在controllers层实现调用。
### 用户只能看到自己id的日记列表
1. 方法：按照用户id进行条件查询
2. 问题：不需要前端用id进行请求
3. 解决：
+ 用jwt进行解析，从上下文中获取用户id，用条件查询，query := db.Where("user_id = ?", userID)
### 条件查询（时间排序）未实现
1. 方法：controllers层实现modifier函数，用if条件进行查询。根据实际，设置按照时间降序为默认排序条件
2. 问题：
+ 并未应用排序逻辑
+ 分页组件无法复用条件查询，因为不同请求的查询逻辑不同
+ 感觉controllers层和service层的逻辑有点奇怪，是不是应该交换一下？
3. 解决
+ query = query.Order("created_at desc")设置为默认按照时间降序
+ controllers层中创建modifier函数，实现修改查询的逻辑。
+ controllers层中调用分页组件，将modifier函数作为参数传递给service层中的paginate函数
+ service层中paginate函数应用modifier来修改model中的查询，查询值传给query
### 封装分页组件
1. 方法：service层中实现分页逻辑，可以复用
### 封装统一返回类组件
1. 方法：定义统一返回类，在controllers层调用
2. 问题：在解决响应缺少分页信息时，采用过一个方案，通过更改响应函数实现
3. 解决：
+ response层中，创建分页相应函数，专门响应分页信息
+ controllers层实现响应调用时，传参分别是vo信息和分页信息，实现两者的输出
+ 我理解的是一个用copy组合一个用响应组合

