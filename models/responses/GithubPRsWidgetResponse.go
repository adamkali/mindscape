package responses

import (
	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type GithubPRData struct {
	Number    int    `json:"number"`
	Title     string `json:"title"`
	HTMLURL   string `json:"html_url"`
	Author    string `json:"author"`
	Repo      string `json:"repo"`
	CreatedAt string `json:"created_at"`
} // @name GithubPRData

type GithubPRsWidgetData struct {
	ReviewRequested []GithubPRData `json:"review_requested"`
	Authored        []GithubPRData `json:"authored"`
	Mentioned       []GithubPRData `json:"mentioned"`
	Assigned        []GithubPRData `json:"assigned"`
} // @name GithubPRsWidgetData

type GithubPRsWidgetResponse struct {
	Data    GithubPRsWidgetData `json:"data"`
	Success bool                `json:"success"`
	Message string              `json:"message"`
} // @name GithubPRsWidgetResponse

func toPRData(items []clients.PullRequest) []GithubPRData {
	out := make([]GithubPRData, 0, len(items))
	for _, pr := range items {
		out = append(out, GithubPRData{
			Number:    pr.Number,
			Title:     pr.Title,
			HTMLURL:   pr.HTMLURL,
			Author:    pr.User.Login,
			Repo:      pr.RepositoryURL,
			CreatedAt: pr.CreatedAt,
		})
	}
	return out
}

func NewGithubPRsWidgetResponse() *GithubPRsWidgetResponse {
	return &GithubPRsWidgetResponse{
		Data:    GithubPRsWidgetData{},
		Success: true,
		Message: "Ok",
	}
}

func (r *GithubPRsWidgetResponse) Fail(ctx echo.Context, code int, err error) error {
	r.Success = false
	r.Message = err.Error()
	return ctx.JSON(code, r)
}

func (r *GithubPRsWidgetResponse) Successful(
	ctx echo.Context,
	reviewRequested, authored, mentioned, assigned []clients.PullRequest,
) error {
	r.Success = true
	r.Data = GithubPRsWidgetData{
		ReviewRequested: toPRData(reviewRequested),
		Authored:        toPRData(authored),
		Mentioned:       toPRData(mentioned),
		Assigned:        toPRData(assigned),
	}
	return ctx.JSON(200, r)
}
