\# UniMate AI 学习计划 API 文档



\*\*Base URL:\*\* `http://localhost:8080`



---



\## 1. 生成学习计划



\*\*接口：\*\* `POST /api/ai/generate-plan`



\*\*请求示例：\*\*



```json

{

&nbsp; "user\_id": "AN051",

&nbsp; "flag": "学会Python爬虫",

&nbsp; "difficulty": "初学者",

&nbsp; "deadline": "2025-12-31"

}

```



\*\*响应示例：\*\*



```json

{

&nbsp; "code": 200,

&nbsp; "data": {

&nbsp;   "plan\_id": "plan\_AN051\_1731340800",

&nbsp;   "flag": "学会Python爬虫",

&nbsp;   "difficulty": "初学者",

&nbsp;   "deadline": "2025-12-31",

&nbsp;   "stages": \[

&nbsp;     {

&nbsp;       "stage\_id": 1,

&nbsp;       "stage\_name": "基础阶段",

&nbsp;       "description": "掌握Python编程基础，学习数据结构和算法，了解网络编程概念",

&nbsp;       "is\_completed": false

&nbsp;     },

&nbsp;     {

&nbsp;       "stage\_id": 2,

&nbsp;       "stage\_name": "进阶阶段",

&nbsp;       "description": "深入学习Python网络库和爬虫框架，了解网页解析库，构建简单爬虫",

&nbsp;       "is\_completed": false

&nbsp;     },

&nbsp;     {

&nbsp;       "stage\_id": 3,

&nbsp;       "stage\_name": "实战阶段",

&nbsp;       "description": "结合实际项目需求，优化爬虫性能，学习反爬技巧，确保合法合规地使用爬虫技术",

&nbsp;       "is\_completed": false

&nbsp;     }

&nbsp;   ],

&nbsp;   "evaluation": {

&nbsp;     "score": 85,

&nbsp;     "details": "计划结构合理，阶段划分清晰，任务安排循序渐进。根据截止日期，时间分配基本合理。建议根据实际学习情况灵活调整进度。",

&nbsp;     "key\_points": \[

&nbsp;       "合理安排每日学习时间",

&nbsp;       "保持学习节奏避免拖延",

&nbsp;       "及时总结复习"

&nbsp;     ]

&nbsp;   },

&nbsp;   "token\_usage": {

&nbsp;     "total": 674,

&nbsp;     "cost": 0.000236

&nbsp;   },

&nbsp;   "created\_at": "2025-11-11 21:30:00"

&nbsp; }

}

```



---



\## 2. 更新阶段状态



\*\*接口：\*\* `PUT /api/ai/plans/:plan\_id/stages/:stage\_id`



\*\*请求示例：\*\*



```

PUT /api/ai/plans/plan\_AN051\_1731340800/stages/1

```



\*\*请求体：\*\*



```json

{

&nbsp; "is\_completed": true

}

```



\*\*响应示例：\*\*



```json

{

&nbsp; "code": 200,

&nbsp; "message": "更新成功",

&nbsp; "data": {

&nbsp;   "plan\_id": "plan\_AN051\_1731340800",

&nbsp;   "stage\_id": 1,

&nbsp;   "is\_completed": true,

&nbsp;   "progress": "33%"

&nbsp; }

}

```



---



\## 3. 健康检查



\*\*接口：\*\* `GET /health`



\*\*响应示例：\*\*



```json

{

&nbsp; "status": "ok"

}

```



---



\## Apifox 导入步骤



1\. 打开 Apifox

2\. 新建项目：UniMate AI

3\. 导入 → OpenAPI

4\. 复制下方 OpenAPI 规范粘贴导入



```yaml

openapi: 3.0.0

info:

&nbsp; title: UniMate AI API

&nbsp; version: 1.0.0

servers:

&nbsp; - url: http://localhost:8080

paths:

&nbsp; /api/ai/generate-plan:

&nbsp;   post:

&nbsp;     summary: 生成学习计划

&nbsp;     requestBody:

&nbsp;       content:

&nbsp;         application/json:

&nbsp;           schema:

&nbsp;             type: object

&nbsp;             required:

&nbsp;               - flag

&nbsp;             properties:

&nbsp;               user\_id:

&nbsp;                 type: string

&nbsp;                 example: AN051

&nbsp;               flag:

&nbsp;                 type: string

&nbsp;                 example: 学会Python爬虫

&nbsp;               difficulty:

&nbsp;                 type: string

&nbsp;                 example: 初学者

&nbsp;               deadline:

&nbsp;                 type: string

&nbsp;                 example: "2025-12-31"

&nbsp;     responses:

&nbsp;       '200':

&nbsp;         description: 成功

&nbsp; /api/ai/plans/{plan\_id}/stages/{stage\_id}:

&nbsp;   put:

&nbsp;     summary: 更新阶段状态

&nbsp;     parameters:

&nbsp;       - name: plan\_id

&nbsp;         in: path

&nbsp;         required: true

&nbsp;         schema:

&nbsp;           type: string

&nbsp;       - name: stage\_id

&nbsp;         in: path

&nbsp;         required: true

&nbsp;         schema:

&nbsp;           type: integer

&nbsp;     requestBody:

&nbsp;       content:

&nbsp;         application/json:

&nbsp;           schema:

&nbsp;             type: object

&nbsp;             required:

&nbsp;               - is\_completed

&nbsp;             properties:

&nbsp;               is\_completed:

&nbsp;                 type: boolean

&nbsp;     responses:

&nbsp;       '200':

&nbsp;         description: 成功

&nbsp; /health:

&nbsp;   get:

&nbsp;     summary: 健康检查

&nbsp;     responses:

&nbsp;       '200':

&nbsp;         description: 成功

```

