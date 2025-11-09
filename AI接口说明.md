# AI助手接口说明

## 接口列表

### 1. AI对话
- **地址**: `POST /api/ai/chat`
- **用途**: 与AI助手进行对话
- **参数**: user_id, message, conversation_id(可选)
- **返回**: AI回复、对话ID、费用统计

### 2. 生成学习计划
- **地址**: `POST /api/ai/generate-plan`
- **用途**: 根据目标生成详细学习计划
- **参数**: goal, deadline, current_level, daily_hours
- **返回**: 分阶段计划、建议

### 3. 目标拆解
- **地址**: `POST /api/ai/breakdown-goal`
- **用途**: 将大目标拆解为小任务
- **参数**: goal, deadline
- **返回**: 子任务列表（带优先级）

### 4. 评估计划
- **地址**: `POST /api/ai/evaluate-plan`
- **用途**: 评估学习计划可行性
- **参数**: goal, plan_content, daily_hours, deadline
- **返回**: 评分、优缺点、建议

### 5. 获取对话列表
- **地址**: `GET /api/ai/conversations?user_id=AN051`
- **返回**: 该用户的所有对话

### 6. 删除对话
- **地址**: `DELETE /api/ai/conversations/:conversation_id`
- **返回**: 删除结果

## 使用示例

```javascript
// 前端调用示例
const response = await fetch('http://localhost:8080/api/ai/chat', {
  method: 'POST',
  headers: { 'Content-Type': 'application/json' },
  body: JSON.stringify({
    user_id: 'AN051',
    message: '我想准备考研数学'
  })
});

const data = await response.json();
console.log(data.data.reply);  // AI的回复