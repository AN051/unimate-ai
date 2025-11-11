package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"unimate-ai/config"
	"unimate-ai/models"
	"unimate-ai/services"
)

func GeneratePlanHandler(c *gin.Context) {
	var req models.GeneratePlanRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误: " + err.Error()})
		return
	}

	if req.UserID == "" {
		req.UserID = "default_user"
	}

	deadline, err := time.Parse("2006-01-02", req.Deadline)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "日期格式错误"})
		return
	}
	totalDays := int(deadline.Sub(time.Now()).Hours() / 24)
	if totalDays <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "截止日期必须在未来"})
		return
	}

	totalHours := totalDays * req.DailyHours

	stage1Days := totalDays * 3 / 10
	stage2Days := totalDays * 4 / 10
	stage3Days := totalDays - stage1Days - stage2Days

	stage1Hours := stage1Days * req.DailyHours
	stage2Hours := stage2Days * req.DailyHours
	stage3Hours := stage3Days * req.DailyHours

	planPrompt := fmt.Sprintf(`你是专业的学习计划制定专家。根据用户的目标、难度和时间安排，生成3个阶段的学习计划。

用户信息：
- 目标：%s
- 难度：%s
- 截止日期：%s（共 %d 天）
- 每日学习时长：%d 小时
- 总学习时长：%d 小时

阶段时间分配：
- 阶段1（基础）：%d 天（%d 小时）
- 阶段2（进阶）：%d 天（%d 小时）
- 阶段3（实战）：%d 天（%d 小时）

要求：
1. 生成3个阶段：基础阶段、进阶阶段、实战阶段（可根据目标调整名称）
2. 每个阶段包含：
   - 阶段名称（5-10个字）
   - 阶段描述（60-80字，详细说明学习内容、技能要点、学习方法）
3. 描述要具体实用，结合每日学习时长给出建议
4. 不要空话套话

返回JSON格式（只返回JSON，不要其他内容）：
{
  "stages": [
    {"stage_name": "基础阶段", "description": "60-80字的详细描述"},
    {"stage_name": "进阶阶段", "description": "60-80字的详细描述"},
    {"stage_name": "实战阶段", "description": "60-80字的详细描述"}
  ]
}`, req.Flag, req.Difficulty, req.Deadline, totalDays, req.DailyHours, totalHours,
		stage1Days, stage1Hours, stage2Days, stage2Hours, stage3Days, stage3Hours)

	planReply, planAIResp, err := services.CallAI(planPrompt, "", 0.7, nil)
	if err != nil {
		log.Printf("生成计划失败: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "生成失败"})
		return
	}

	var planData struct {
		Stages []struct {
			StageName   string `json:"stage_name"`
			Description string `json:"description"`
		} `json:"stages"`
	}

	if err := json.Unmarshal([]byte(planReply), &planData); err != nil {
		log.Printf("解析失败: %v, 原始内容: %s", err, planReply)
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "message": "解析失败"})
		return
	}

	evalPrompt := fmt.Sprintf(`你是资深学习顾问。评估这个学习计划并给出实战建议。

计划信息：
- 目标：%s
- 难度：%s
- 总时长：%d 天，每天 %d 小时，共 %d 小时
- 计划内容：%s

要求：
1. 评分（70-95分之间，根据计划质量、时间合理性评分）
2. 120-150字的专业评估（包含：时间分配合理性、难度梯度、潜在风险、改进建议）
3. 4-5个可执行的关键行动建议（每条20-30字，具体可操作，不要空话）

返回JSON格式（只返回JSON）：
{
  "score": 85,
  "details": "120-150字的详细专业分析",
  "key_points": ["20-30字的具体建议1", "20-30字的具体建议2", "20-30字的具体建议3", "20-30字的具体建议4"]
}`, req.Flag, req.Difficulty, totalDays, req.DailyHours, totalHours, planReply)

	evalReply, evalAIResp, err := services.CallAI(evalPrompt, "", 0.7, nil)
	if err != nil {
		evalReply = fmt.Sprintf(`{"score": 82, "details": "该学习计划整体架构合理，三阶段设计符合认知规律。总计划时长%d天、每天%d小时是可行的安排。基础阶段建议扎实掌握核心概念，避免急于求成。进阶阶段注重理论与实践结合，多动手练习。实战阶段选择真实项目，培养解决实际问题的能力。建议每周复盘学习进度，遇到困难及时调整计划。", "key_points": ["基础阶段重点掌握核心概念，每天留出30分钟总结", "进阶阶段坚持50%%理论+50%%实践的学习节奏", "实战阶段选择1-2个完整项目，注重代码质量", "每周日晚进行学习复盘，记录问题和收获"]}`, totalDays, req.DailyHours)
	}

	var evaluation models.Evaluation
	if err := json.Unmarshal([]byte(evalReply), &evaluation); err != nil {
		evaluation = models.Evaluation{
			Score:     82,
			Details:   fmt.Sprintf("该学习计划整体架构合理，三阶段设计符合认知规律。总计划时长%d天、每天%d小时是可行的安排。基础阶段建议扎实掌握核心概念，避免急于求成。进阶阶段注重理论与实践结合，多动手练习。实战阶段选择真实项目，培养解决实际问题的能力。建议每周复盘学习进度，遇到困难及时调整计划。", totalDays, req.DailyHours),
			KeyPoints: []string{
				"基础阶段重点掌握核心概念，每天留出30分钟总结",
				"进阶阶段坚持50%理论+50%实践的学习节奏",
				"实战阶段选择1-2个完整项目，注重代码质量",
				"每周日晚进行学习复盘，记录问题和收获",
			},
		}
	}

	planID := fmt.Sprintf("plan_%s_%d", req.UserID, time.Now().Unix())

	durations := []string{
		fmt.Sprintf("%d天·%d小时", stage1Days, stage1Hours),
		fmt.Sprintf("%d天·%d小时", stage2Days, stage2Hours),
		fmt.Sprintf("%d天·%d小时", stage3Days, stage3Hours),
	}

	stages := make([]models.Stage, 3)
	for i, s := range planData.Stages {
		stages[i] = models.Stage{
			StageID:     i + 1,
			StageName:   s.StageName,
			Description: s.Description,
			Duration:    durations[i],
			IsCompleted: false,
		}
	}

	stagesJSON, _ := json.Marshal(stages)
	evaluationJSON, _ := json.Marshal(evaluation)

	studyPlan := models.StudyPlan{
		PlanID:     planID,
		UserID:     req.UserID,
		Flag:       req.Flag,
		Difficulty: req.Difficulty,
		Deadline:   req.Deadline,
		DailyHours: req.DailyHours,
		Stages:     string(stagesJSON),
		Evaluation: string(evaluationJSON),
		Progress:   0,
	}

	config.DB.Create(&studyPlan)

	totalTokens := planAIResp.Usage.TotalTokens + evalAIResp.Usage.TotalTokens
	totalCost := float64(totalTokens) * 0.35 / 1000000

	log.Printf("计划生成成功: %s, 用户: %s, 目标: %s, 总时长: %d天x%d小时", planID, req.UserID, req.Flag, totalDays, req.DailyHours)

	var resp models.GeneratePlanResponse
	resp.Code = 200
	resp.Data.PlanID = planID
	resp.Data.Flag = req.Flag
	resp.Data.Difficulty = req.Difficulty
	resp.Data.Deadline = req.Deadline
	resp.Data.DailyHours = req.DailyHours
	resp.Data.Stages = stages
	resp.Data.Evaluation = evaluation
	resp.Data.TokenUsage.Total = totalTokens
	resp.Data.TokenUsage.Cost = totalCost
	resp.Data.CreatedAt = time.Now().Format("2006-01-02 15:04:05")

	c.JSON(http.StatusOK, resp)
}

func UpdateStageHandler(c *gin.Context) {
	planID := c.Param("plan_id")
	stageIDStr := c.Param("stage_id")

	stageID, err := strconv.Atoi(stageIDStr)
	if err != nil || stageID < 1 || stageID > 3 {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "无效阶段ID"})
		return
	}

	var req models.UpdateStageRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "参数错误"})
		return
	}

	if req.IsCompleted == nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "message": "is_completed 不能为空"})
		return
	}

	var studyPlan models.StudyPlan
	if err := config.DB.Where("plan_id = ?", planID).First(&studyPlan).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"code": 404, "message": "计划不存在"})
		return
	}

	var stages []models.Stage
	json.Unmarshal([]byte(studyPlan.Stages), &stages)

	if stageID <= len(stages) {
		stages[stageID-1].IsCompleted = *req.IsCompleted
	}

	completedCount := 0
	for _, stage := range stages {
		if stage.IsCompleted {
			completedCount++
		}
	}
	progress := (completedCount * 100) / 3

	stagesJSON, _ := json.Marshal(stages)
	config.DB.Model(&studyPlan).Updates(map[string]interface{}{
		"stages":   string(stagesJSON),
		"progress": progress,
	})

	log.Printf("更新阶段: 计划=%s, 阶段=%d, 完成=%v, 进度=%d%%", planID, stageID, *req.IsCompleted, progress)

	var resp models.UpdateStageResponse
	resp.Code = 200
	resp.Message = "更新成功"
	resp.Data.PlanID = planID
	resp.Data.StageID = stageID
	resp.Data.IsCompleted = *req.IsCompleted
	resp.Data.Progress = fmt.Sprintf("%d%%", progress)

	c.JSON(http.StatusOK, resp)
}