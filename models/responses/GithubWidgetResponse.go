package responses

import (
	"time"

	"github.com/adamkali/mindscape/clients"
	"github.com/labstack/echo/v4"
)

type GithubWidgetProfileData struct {
	AvatarURL   string `json:"avatar_url"`
	HTMLURL     string `json:"html_url"`
	Name        string `json:"name"`
	Company     string `json:"company"`
	Bio         string `json:"bio"`
	PublicRepos int    `json:"public_repos"`
	Followers   int    `json:"followers"`
}

type GithubWidgetCommitsDayData struct {
	Date    string `json:"date"`
	Count   int    `json:"count"`
	Percent int    `json:"percent"`
	Color   string `json:"color"`
}

type GithubWidgetCommitsWeekData struct {
	WeekStart string                     `json:"week_start"`
	Monday    GithubWidgetCommitsDayData `json:"monday"`
	Tuesday   GithubWidgetCommitsDayData `json:"tuesday"`
	Wednesday GithubWidgetCommitsDayData `json:"wednesday"`
	Thursday  GithubWidgetCommitsDayData `json:"thursday"`
	Friday    GithubWidgetCommitsDayData `json:"friday"`
	Saturday  GithubWidgetCommitsDayData `json:"saturday"`
	Sunday    GithubWidgetCommitsDayData `json:"sunday"`
}

type GithubWidgetCommitsData struct {
	Total int                           `json:"total"`
	Weeks []GithubWidgetCommitsWeekData `json:"weeks"`
}

type GithubWidgetData struct {
	Profile GithubWidgetProfileData `json:"profile"`
	Commits GithubWidgetCommitsData `json:"commits"`
}

type GithubWidgetResponse struct {
	Data    GithubWidgetData `json:"data"`
	Success bool             `json:"success"`
	Message string           `json:"message"`
} // @name GithubResponse

func NewGithubWidgetProfileData(profile clients.GitHubProfile) GithubWidgetProfileData {
	return GithubWidgetProfileData{
		AvatarURL:   profile.AvatarURL,
		HTMLURL:     profile.HTMLURL,
		Name:        profile.Name,
		Company:     profile.Company,
		Bio:         profile.Bio,
		PublicRepos: profile.PublicRepos,
		Followers:   profile.Followers,
	}
}

// getWeekStart returns the Monday of the week for a given date
func getWeekStart(date time.Time) time.Time {
	weekday := date.Weekday()
	if weekday == time.Sunday {
		weekday = 7
	}
	daysToMonday := int(weekday) - int(time.Monday)
	return date.AddDate(0, 0, -daysToMonday)
}

// newEmptyDayData creates an empty day entry with default values
func newEmptyDayData(date string) GithubWidgetCommitsDayData {
	return GithubWidgetCommitsDayData{
		Date:    date,
		Count:   0,
		Percent: 0,
		Color:   "#cc1e1e", // Base color for 0 commits
	}
}

// newDayDataFromEntry creates a day entry from a DateEntry
func newDayDataFromEntry(entry *clients.DateEntry) GithubWidgetCommitsDayData {
	return GithubWidgetCommitsDayData{
		Date:    entry.Date,
		Count:   entry.Count,
		Percent: entry.Percent,
		Color:   entry.Color,
	}
}

// newEmptyWeekData creates an empty week with all days initialized
func newEmptyWeekData(weekStart time.Time) GithubWidgetCommitsWeekData {
	return GithubWidgetCommitsWeekData{
		WeekStart: weekStart.Format("2006-01-02"),
		Monday:    newEmptyDayData(weekStart.Format("2006-01-02")),
		Tuesday:   newEmptyDayData(weekStart.AddDate(0, 0, 1).Format("2006-01-02")),
		Wednesday: newEmptyDayData(weekStart.AddDate(0, 0, 2).Format("2006-01-02")),
		Thursday:  newEmptyDayData(weekStart.AddDate(0, 0, 3).Format("2006-01-02")),
		Friday:    newEmptyDayData(weekStart.AddDate(0, 0, 4).Format("2006-01-02")),
		Saturday:  newEmptyDayData(weekStart.AddDate(0, 0, 5).Format("2006-01-02")),
		Sunday:    newEmptyDayData(weekStart.AddDate(0, 0, 6).Format("2006-01-02")),
	}
}

// setDayInWeek sets the appropriate day field in the week struct
func setDayInWeek(week *GithubWidgetCommitsWeekData, weekday time.Weekday, dayData GithubWidgetCommitsDayData) {
	switch weekday {
	case time.Monday:
		week.Monday = dayData
	case time.Tuesday:
		week.Tuesday = dayData
	case time.Wednesday:
		week.Wednesday = dayData
	case time.Thursday:
		week.Thursday = dayData
	case time.Friday:
		week.Friday = dayData
	case time.Saturday:
		week.Saturday = dayData
	case time.Sunday:
		week.Sunday = dayData
	}
}

func NewGithubWidgetCommitsData(commits *clients.FinalDateEntry) GithubWidgetCommitsData {
	// Calculate total commits
	total := 0
	commits.Iterate(func(date string, entry *clients.DateEntry) {
		total += entry.Count
	})

	// Get sorted entries (oldest to newest)
	entries := commits.GetSortedEntries()
	println(len(entries))

	if len(entries) == 0 {
		return GithubWidgetCommitsData{
			Total: 0,
			Weeks: []GithubWidgetCommitsWeekData{},
		}
	}

	// Map to store weeks by their start date
	weeksMap := make(map[string]*GithubWidgetCommitsWeekData)
	var weekStartDates []string // To maintain order

	for _, entry := range entries {
		// Parse the date string
		date, err := time.Parse("2006-01-02", entry.Date)
		if err != nil {
			continue
		}

		// Get the Monday of this week
		weekStart := getWeekStart(date)
		weekStartStr := weekStart.Format("2006-01-02")

		// Create week if it doesn't exist
		if _, exists := weeksMap[weekStartStr]; !exists {
			week := newEmptyWeekData(weekStart)
			weeksMap[weekStartStr] = &week
			weekStartDates = append(weekStartDates, weekStartStr)
		}

		// Set the day data in the appropriate weekday slot
		dayData := newDayDataFromEntry(entry)
		setDayInWeek(weeksMap[weekStartStr], date.Weekday(), dayData)
	}

	// Build the final weeks slice in order
	weeks := make([]GithubWidgetCommitsWeekData, 0, len(weekStartDates))
	for _, weekStartStr := range weekStartDates {
		weeks = append(weeks, *weeksMap[weekStartStr])
	}

	return GithubWidgetCommitsData{
		Total: total,
		Weeks: weeks,
	}
}


func NewGithubWidgetData(
	profile *clients.GitHubProfile,
	commits *clients.FinalDateEntry,
) *GithubWidgetData {
	return &GithubWidgetData{
		Profile: NewGithubWidgetProfileData(*profile),
		Commits: NewGithubWidgetCommitsData(commits),
	}
}

func NewGithubWidgetResponse(
) *GithubWidgetResponse {
	return &GithubWidgetResponse{
		Data:    GithubWidgetData{},
		Success: true,
		Message: "Ok",
	}
}

func (w *GithubWidgetResponse) Fail(ctx echo.Context, code int, err error) error {
	w.Success = false
	w.Message = err.Error()
	return ctx.JSON(code, w)
}

func (w *GithubWidgetResponse) Successful(
	ctx echo.Context,
	profile *clients.GitHubProfile,
	commits *clients.FinalDateEntry,
) error {
	w.Success = true
	w.Data = *NewGithubWidgetData(profile, commits)
	return ctx.JSON(200, w)
}

