package tools

import (
	"context"
	"fmt"
	"strconv"
	"strings"

	"github.com/modelcontextprotocol/go-sdk/mcp"

	"campusapp/internal/mcp/campus"
)

type Toolset struct {
	Client *campus.Client
}

func (t *Toolset) Register(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "csu.grade",
		Description: "查询 csugo 成绩接口，返回课程成绩列表。",
	}, t.grade)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "csu.rank",
		Description: "查询 csugo 专业排名接口，返回综合成绩与排名。",
	}, t.rank)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "csu.classes",
		Description: "查询 csugo 课表接口，根据学期与周次返回排课信息。",
	}, t.classes)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "csu.bus_search",
		Description: "查询 csugo 校车接口，按起点、终点和关键字返回班次。",
	}, t.bus)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "csu.jobs",
		Description: "查询 csugo 招聘接口，支持分页与抓取招聘会地点。",
	}, t.jobs)
}

func (t *Toolset) grade(ctx context.Context, _ *mcp.CallToolRequest, args credentialArgs) (*mcp.CallToolResult, campus.GradeResponse, error) {
	if err := args.validate(); err != nil {
		return nil, campus.GradeResponse{}, err
	}
	resp, err := t.Client.Grade(ctx, args.StudentID, args.Password)
	return nil, resp, err
}

func (t *Toolset) rank(ctx context.Context, _ *mcp.CallToolRequest, args credentialArgs) (*mcp.CallToolResult, campus.RankResponse, error) {
	if err := args.validate(); err != nil {
		return nil, campus.RankResponse{}, err
	}
	resp, err := t.Client.Rank(ctx, args.StudentID, args.Password)
	return nil, resp, err
}

func (t *Toolset) classes(ctx context.Context, _ *mcp.CallToolRequest, args classArgs) (*mcp.CallToolResult, campus.ClassScheduleResponse, error) {
	if err := args.credentialArgs.validate(); err != nil {
		return nil, campus.ClassScheduleResponse{}, err
	}
	if strings.TrimSpace(args.Term) == "" {
		return nil, campus.ClassScheduleResponse{}, fmt.Errorf("term is required, e.g. 2024-2025-1")
	}
	if args.Week < 0 {
		return nil, campus.ClassScheduleResponse{}, fmt.Errorf("week must be >= 0")
	}
	resp, err := t.Client.ClassSchedule(ctx, args.StudentID, args.Password, args.Term, strconv.Itoa(args.Week))
	return nil, resp, err
}

func (t *Toolset) bus(ctx context.Context, _ *mcp.CallToolRequest, args busArgs) (*mcp.CallToolResult, campus.BusResponse, error) {
	if strings.TrimSpace(args.Start) == "" || strings.TrimSpace(args.End) == "" || strings.TrimSpace(args.Day) == "" {
		return nil, campus.BusResponse{}, fmt.Errorf("start, end, and day are required")
	}
	resp, err := t.Client.BusSearch(ctx, args.Start, args.End, args.Day)
	return nil, resp, err
}

func (t *Toolset) jobs(ctx context.Context, _ *mcp.CallToolRequest, args jobArgs) (*mcp.CallToolResult, campus.JobResponse, error) {
	category := strings.TrimSpace(args.CategoryID)
	if category == "" {
		category = "1"
	}
	page := args.Page
	if page <= 0 {
		page = 1
	}
	size := args.PageSize
	switch {
	case size <= 0:
		size = 10
	case size > 50:
		size = 50
	}
	resp, err := t.Client.JobList(ctx, category, page, size, args.IncludeSchedule)
	return nil, resp, err
}

type credentialArgs struct {
	StudentID string `json:"studentId" jsonschema:"学号，用于统一认证"`
	Password  string `json:"password" jsonschema:"教务系统/统一认证密码"`
}

func (c credentialArgs) validate() error {
	if strings.TrimSpace(c.StudentID) == "" {
		return fmt.Errorf("studentId is required")
	}
	if strings.TrimSpace(c.Password) == "" {
		return fmt.Errorf("password is required")
	}
	return nil
}

type classArgs struct {
	credentialArgs
	Term string `json:"term" jsonschema:"学年学期，例如 2024-2025-1"`
	Week int    `json:"week" jsonschema:"周次；0 表示全部周次"`
}

type busArgs struct {
	Start string `json:"start" jsonschema:"起点，如 校本部图书馆前坪"`
	End   string `json:"end" jsonschema:"终点，如 新校区艺术楼"`
	Day   string `json:"day" jsonschema:"班次关键字，例如 周一至周五/星期六"`
}

type jobArgs struct {
	CategoryID      string `json:"categoryId" jsonschema:"招聘类型：1 = 本部、2 = 湘雅、3 = 铁道、4 = 在线、5 = 事业招考"`
	Page            int    `json:"page" jsonschema:"页码，默认为 1"`
	PageSize        int    `json:"pageSize" jsonschema:"每页条数（1-50）"`
	IncludeSchedule bool   `json:"includeSchedule" jsonschema:"是否额外抓取招聘会地点（耗时更久）"`
}
