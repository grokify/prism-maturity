package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	prism "github.com/grokify/prism-maturity"
)

// Goal Tool Input/Output types

// ListGoalsInput is the input for list_goals.
type ListGoalsInput struct {
	File   string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	Status string `json:"status,omitempty" jsonschema:"description=Filter by status (active, on_hold, completed, cancelled)"`
	Limit  int    `json:"limit,omitempty" jsonschema:"description=Maximum results (default 50)"`
}

// GoalSummary provides a brief summary of a goal.
type GoalSummary struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Status       string `json:"status"`
	Priority     int    `json:"priority"`
	CurrentLevel int    `json:"current_level"`
	TargetLevel  int    `json:"target_level"`
	Owner        string `json:"owner,omitempty"`
}

// ListGoalsOutput is the output for list_goals.
type ListGoalsOutput struct {
	File  string        `json:"file"`
	Goals []GoalSummary `json:"goals"`
	Total int           `json:"total"`
	Error string        `json:"error,omitempty"`
}

// GetGoalInput is the input for get_goal.
type GetGoalInput struct {
	File string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	ID   string `json:"id" jsonschema:"description=Goal ID"`
}

// GetGoalOutput is the output for get_goal.
type GetGoalOutput struct {
	File  string      `json:"file"`
	Goal  *prism.Goal `json:"goal,omitempty"`
	Error string      `json:"error,omitempty"`
}

// GoalProgressInput is the input for goal_progress.
type GoalProgressInput struct {
	File string `json:"file" jsonschema:"description=Path to PRISM document JSON file"`
	ID   string `json:"id" jsonschema:"description=Goal ID"`
}

// LevelProgress describes progress toward a maturity level.
type LevelProgress struct {
	Level         int     `json:"level"`
	Name          string  `json:"name"`
	SLOsMet       int     `json:"slos_met"`
	SLOsTotal     int     `json:"slos_total"`
	CriteriaMet   int     `json:"criteria_met"`
	CriteriaTotal int     `json:"criteria_total"`
	Complete      bool    `json:"complete"`
	Progress      float64 `json:"progress_percent"`
}

// GoalProgressOutput is the output for goal_progress.
type GoalProgressOutput struct {
	File           string          `json:"file"`
	GoalID         string          `json:"goal_id"`
	GoalName       string          `json:"goal_name"`
	CurrentLevel   int             `json:"current_level"`
	TargetLevel    int             `json:"target_level"`
	LevelProgress  []LevelProgress `json:"level_progress,omitempty"`
	OverallPercent float64         `json:"overall_percent"`
	Error          string          `json:"error,omitempty"`
}

// GoalStatusesInput is the input for goal_statuses.
type GoalStatusesInput struct{}

// GoalStatusesOutput is the output for goal_statuses.
type GoalStatusesOutput struct {
	Statuses []string `json:"statuses"`
}

// RegisterGoalTools registers goal-related MCP tools.
func (s *Server) RegisterGoalTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "list_goals",
		Description: "List goals from a PRISM document with optional status filtering.",
	}, s.listGoals)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_goal",
		Description: "Get details for a specific goal including maturity model definition.",
	}, s.getGoal)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "goal_progress",
		Description: "Get maturity progress for a goal including SLO and criteria status per level.",
	}, s.goalProgress)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "goal_statuses",
		Description: "List all valid goal status values.",
	}, s.goalStatuses)
}

// loadDocument loads a PRISM document from file.
func loadDocument(filePath string) (*prism.PRISMDocument, error) {
	data, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("reading file: %w", err)
	}

	var doc prism.PRISMDocument
	if err := json.Unmarshal(data, &doc); err != nil {
		return nil, fmt.Errorf("parsing document: %w", err)
	}

	return &doc, nil
}

// Tool implementations

func (s *Server) listGoals(ctx context.Context, req *mcp.CallToolRequest, input ListGoalsInput) (*mcp.CallToolResult, ListGoalsOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, ListGoalsOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	if input.Limit == 0 {
		input.Limit = 50
	}

	var summaries []GoalSummary
	for _, goal := range doc.Goals {
		// Apply status filter
		if input.Status != "" && goal.Status != input.Status {
			continue
		}

		summaries = append(summaries, GoalSummary{
			ID:           goal.ID,
			Name:         goal.Name,
			Status:       goal.Status,
			Priority:     goal.Priority,
			CurrentLevel: goal.CurrentLevel,
			TargetLevel:  goal.TargetLevel,
			Owner:        goal.Owner,
		})

		if len(summaries) >= input.Limit {
			break
		}
	}

	return nil, ListGoalsOutput{
		File:  filepath.Base(input.File),
		Goals: summaries,
		Total: len(summaries),
	}, nil
}

func (s *Server) getGoal(ctx context.Context, req *mcp.CallToolRequest, input GetGoalInput) (*mcp.CallToolResult, GetGoalOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, GetGoalOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	for i := range doc.Goals {
		if doc.Goals[i].ID == input.ID {
			return nil, GetGoalOutput{
				File: filepath.Base(input.File),
				Goal: &doc.Goals[i],
			}, nil
		}
	}

	return nil, GetGoalOutput{
		File:  filepath.Base(input.File),
		Error: fmt.Sprintf("goal not found: %s", input.ID),
	}, nil
}

func (s *Server) goalProgress(ctx context.Context, req *mcp.CallToolRequest, input GoalProgressInput) (*mcp.CallToolResult, GoalProgressOutput, error) {
	doc, err := loadDocument(input.File)
	if err != nil {
		return nil, GoalProgressOutput{
			File:  filepath.Base(input.File),
			Error: err.Error(),
		}, nil
	}

	var goal *prism.Goal
	for i := range doc.Goals {
		if doc.Goals[i].ID == input.ID {
			goal = &doc.Goals[i]
			break
		}
	}

	if goal == nil {
		return nil, GoalProgressOutput{
			File:   filepath.Base(input.File),
			GoalID: input.ID,
			Error:  fmt.Sprintf("goal not found: %s", input.ID),
		}, nil
	}

	output := GoalProgressOutput{
		File:         filepath.Base(input.File),
		GoalID:       goal.ID,
		GoalName:     goal.Name,
		CurrentLevel: goal.CurrentMaturityLevel(doc),
		TargetLevel:  goal.TargetLevel,
	}

	// Calculate progress for each level
	if goal.MaturityModel != nil {
		for _, level := range goal.MaturityModel.Levels {
			slosMet, slosTotal := goal.SLOsMetForLevel(level.Level, doc)
			critMet, critTotal := goal.CriteriaMetForLevel(level.Level, doc)

			total := slosTotal + critTotal
			met := slosMet + critMet
			var progress float64
			if total > 0 {
				progress = float64(met) / float64(total) * 100
			}

			output.LevelProgress = append(output.LevelProgress, LevelProgress{
				Level:         level.Level,
				Name:          level.Name,
				SLOsMet:       slosMet,
				SLOsTotal:     slosTotal,
				CriteriaMet:   critMet,
				CriteriaTotal: critTotal,
				Complete:      met == total && total > 0,
				Progress:      progress,
			})
		}
	}

	// Calculate overall progress toward target level
	if goal.TargetLevel > 0 {
		output.OverallPercent = float64(output.CurrentLevel) / float64(goal.TargetLevel) * 100
		if output.OverallPercent > 100 {
			output.OverallPercent = 100
		}
	}

	return nil, output, nil
}

func (s *Server) goalStatuses(ctx context.Context, req *mcp.CallToolRequest, input GoalStatusesInput) (*mcp.CallToolResult, GoalStatusesOutput, error) {
	return nil, GoalStatusesOutput{
		Statuses: prism.AllGoalStatuses(),
	}, nil
}
