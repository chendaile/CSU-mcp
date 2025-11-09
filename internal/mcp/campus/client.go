package campus

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
	"time"
)

type Client struct {
	baseURL string
	token   string
	http    *http.Client
}

func NewClient(baseURL, token string, timeout time.Duration) (*Client, error) {
	if _, err := url.ParseRequestURI(baseURL); err != nil {
		return nil, fmt.Errorf("invalid base url %q: %w", baseURL, err)
	}

	return &Client{
		baseURL: strings.TrimRight(baseURL, "/"),
		token:   token,
		http:    &http.Client{Timeout: timeout},
	}, nil
}

func (c *Client) Grade(ctx context.Context, id, password string) (GradeResponse, error) {
	var resp GradeResponse
	path := fmt.Sprintf("/api/v1/jwc/%s/%s/grade", escape(id), escape(password))
	if err := c.get(ctx, path, &resp); err != nil {
		return resp, err
	}
	return resp, resp.validate("grade")
}

func (c *Client) Rank(ctx context.Context, id, password string) (RankResponse, error) {
	var resp RankResponse
	path := fmt.Sprintf("/api/v1/jwc/%s/%s/rank", escape(id), escape(password))
	if err := c.get(ctx, path, &resp); err != nil {
		return resp, err
	}
	return resp, resp.validate("rank")
}

func (c *Client) ClassSchedule(ctx context.Context, id, password, term, week string) (ClassScheduleResponse, error) {
	var resp ClassScheduleResponse
	path := fmt.Sprintf("/api/v1/jwc/%s/%s/class/%s/%s",
		escape(id),
		escape(password),
		escape(term),
		escape(week),
	)
	if err := c.get(ctx, path, &resp); err != nil {
		return resp, err
	}
	return resp, resp.validate("class schedule")
}

func (c *Client) BusSearch(ctx context.Context, start, end, day string) (BusResponse, error) {
	var resp BusResponse
	path := fmt.Sprintf("/api/v1/bus/search/%s/%s/%s",
		escape(start),
		escape(end),
		escape(day),
	)
	if err := c.get(ctx, path, &resp); err != nil {
		return resp, err
	}
	return resp, resp.validate("bus search")
}

func (c *Client) JobList(ctx context.Context, typeID string, pageIndex, pageSize int, includeSchedule bool) (JobResponse, error) {
	var resp JobResponse
	path := fmt.Sprintf("/api/v1/job/%s/%s/%s/%s",
		escape(typeID),
		escape(strconv.Itoa(pageIndex)),
		escape(strconv.Itoa(pageSize)),
		escape(boolFlag(includeSchedule)),
	)
	if err := c.get(ctx, path, &resp); err != nil {
		return resp, err
	}
	return resp, resp.validate("job list")
}

func (c *Client) get(ctx context.Context, path string, target any) error {
	urlStr, err := c.buildURL(path)
	if err != nil {
		return err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, urlStr, nil)
	if err != nil {
		return fmt.Errorf("create request: %w", err)
	}
	resp, err := c.http.Do(req)
	if err != nil {
		return fmt.Errorf("call csugo: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, 2048))
		return fmt.Errorf("csugo %s returned %d: %s", path, resp.StatusCode, strings.TrimSpace(string(body)))
	}

	if err := json.NewDecoder(resp.Body).Decode(target); err != nil {
		return fmt.Errorf("decode response: %w", err)
	}
	return nil
}

func (c *Client) buildURL(path string) (string, error) {
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	u, err := url.Parse(c.baseURL + path)
	if err != nil {
		return "", fmt.Errorf("build url: %w", err)
	}
	values := u.Query()
	if c.token != "" {
		values.Set("token", c.token)
	}
	u.RawQuery = values.Encode()
	return u.String(), nil
}

type baseEnvelope struct {
	StateCode int    `json:"StateCode"`
	Error     string `json:"Error"`
}

func (b baseEnvelope) validate(action string) error {
	if b.StateCode == 1 {
		return nil
	}
	msg := strings.TrimSpace(b.Error)
	if msg == "" {
		msg = fmt.Sprintf("%s failed with state %d", action, b.StateCode)
	}
	return &Error{Code: b.StateCode, Message: msg}
}

type Error struct {
	Code    int
	Message string
}

func (e *Error) Error() string {
	return fmt.Sprintf("csugo error (code %d): %s", e.Code, e.Message)
}

type GradeResponse struct {
	baseEnvelope
	Grades []Grade `json:"Grades"`
}

type Grade struct {
	ClassNo     int    `json:"ClassNo"`
	FirstTerm   string `json:"FirstTerm"`
	GottenTerm  string `json:"GottenTerm"`
	ClassName   string `json:"ClassName"`
	MiddleGrade string `json:"MiddleGrade"`
	FinalGrade  string `json:"FinalGrade"`
	Grade       string `json:"Grade"`
	ClassScore  string `json:"ClassScore"`
	ClassType   string `json:"ClassType"`
	ClassProp   string `json:"ClassProp"`
}

type RankResponse struct {
	baseEnvelope
	Rank []RankEntry `json:"Rank"`
}

type RankEntry struct {
	Term       string `json:"Term"`
	TotalScore string `json:"TotalScore"`
	ClassRank  string `json:"ClassRank"`
	AverScore  string `json:"AverScore"`
}

type ClassScheduleResponse struct {
	baseEnvelope
	Class        [][]ClassEntry `json:"Class"`
	StartWeekDay string         `json:"StartWeekDay"`
}

type ClassEntry struct {
	ClassName string `json:"ClassName"`
	Teacher   string `json:"Teacher"`
	Weeks     string `json:"Weeks"`
	Place     string `json:"Place"`
}

type BusResponse struct {
	baseEnvelope
	Buses []Bus `json:"Buses"`
}

type Bus struct {
	StartTime string   `json:"StartTime"`
	Start     string   `json:"Start"`
	End       string   `json:"End"`
	RunTime   string   `json:"RunTime"`
	Num       string   `json:"Num"`
	Seat      string   `json:"Seat"`
	Stations  []string `json:"Stations"`
}

type JobResponse struct {
	baseEnvelope
	Jobs []Job `json:"Jobs"`
}

type Job struct {
	Link  string `json:"Link"`
	Title string `json:"Title"`
	Time  string `json:"Time"`
	Place string `json:"Place"`
}

func escape(value string) string {
	if value == "" {
		return ""
	}
	return url.PathEscape(value)
}

func boolFlag(v bool) string {
	if v {
		return "1"
	}
	return "0"
}
