### 代码协作

参考文档 [Git 分支与冲突解决探究](http://zxlmdonnie.cn/archives/1709691289692)

#### 首次获取项目
1. 从Github上克隆项目 `git clone https://github.com/Werun-backend/YiyuGoDemo.git`
2. 创建本地开发分支：`git checkout -b dev_zfr` (xxx一般为姓名首字母小写，如dev_zfr);

#### 修改代码前
3. 首先保证在自己的开发分支，每次写代码前，都要更新本地代码和远程main保持一致：
    1. 更新本地main分支：`git pull origin main`
    2. 将本地main分支合并到自己的分支 `git merge main dev_zfr`

#### 修改代码后
4. 代码提交需先push到自己的远程开发分支上：
   1. `git add 文件名` 或者使用 Goland 的图形化界面
   2. `git commit -m '日志信息'`
   3. `git push origin dev_zfr`;
6. 在GitHub的Pull Requests中手动将xxx_dev合并到main分支：点击合并请求合并分支到main，一路提交即可

### 前端对接

- 对入参、出参的 `json` 实体进行数据转换(dto vo)，以适配前端需求。
- 封装统一返回类组件 HTTP状态码
- 封装分页组件

把上述规范写一个文档（MarkDown）进行说明。

### 入参dto 出参vo
#### dto
+ dto保证服务端数据传输和客户端响应只传递必要数据Content和TagIds
+ 定义DiaryDto 代表日记的请求模型，定义 DiaryToDto 将Diary模型转换为DiaryDto
+ CreateDiary函数使用了dto，把用户请求的数据绑定到diaryDto，只需要检查dto中的Content和TagIds是否存在。
+ 创建日记记录时，只需要提供content、tag和Id外键约束，调用DiaryToDto数据转换为dto
#### vo
+ vo保证前端只输出必要信息，空值不返回
+ 定义DiaryVo 用于API响应的日记信息，定义Copy 直接实现赋值
+ GetDiaries中调用copy方法将绑定的数据转换为vo格式
+ 响应时，返回diaryVos，保证前端输出的是vo格式

### 封装统一返回类组件
+ SuccessResponse 创建一个表示成功的响应
+ UserErrorResponse 创建一个表示用户错误的响应
+ UserErrorNoMsgResponse 创建一个表示用户错误的响应，不包含额外的数据。
+ InternalErrorResponse 创建一个表示服务器内部错误的响应
+ 将响应写入HTTP，在CreateDiary和GetDiaries中使用统一返回类组件

### 封装分页组件
+ 在GetDiaries中默认排序条件是创建时间，默认降序排列
+ 在service层中实现分页组件逻辑并创建分页响应


