package ui

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/a-h/templ"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/httplog"
	prbot "github.com/marqeta/pr-bot"
	"github.com/marqeta/pr-bot/metrics"
	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/ui/templates"
	"github.com/marqeta/pr-bot/ui/templates/components"
	"github.com/marqeta/pr-bot/ui/templates/pages"
)

const (
	Title = "PR-Bot"
)

type controller struct {
	metrics           metrics.Emitter
	staticContents    http.FileSystem
	evaluationManager evaluation.Manager
}

func (c *controller) FileServer(w http.ResponseWriter, r *http.Request) {
	if strings.HasSuffix(r.URL.Path, "/") {
		http.NotFound(w, r)
		return
	}
	// Remove the path prefix from the request path.
	// /ui/static/css/style.css => /css/style.css
	fs := http.StripPrefix(prbot.UI+"/static/", http.FileServer(c.staticContents))
	fs.ServeHTTP(w, r)
}

func (c *controller) ListReports(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oplog := httplog.LogEntry(ctx)
	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	prNumber := chi.URLParam(r, "pr")
	pr := fmt.Sprintf("%s/%s/%s", owner, repo, prNumber)

	reports, err := c.evaluationManager.ListReports(ctx, pr)
	if err != nil {
		oplog.Err(err).Msg("error listing reports")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if len(reports) == 0 {
		oplog.Err(err).Msg("no reports found for pr")
		c.NotFound(w, r)
		return
	}

	oplog.Info().Msgf("%v reports found", len(reports))
	navbar := pages.Navbar(Title)
	lottie := components.LottieAutomation()
	PRDetails := pages.PRDetails(reports[0])
	metadata := pages.Metadata(lottie, PRDetails)
	eventsTable := pages.EventsTable(reports)
	body := pages.Container(pages.ListReportsPage(navbar, metadata, eventsTable))
	html := templates.Html(Title, body)
	templ.Handler(html).ServeHTTP(w, r)
}

func (c *controller) ReportDetails(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	oplog := httplog.LogEntry(ctx)

	owner := chi.URLParam(r, "owner")
	repo := chi.URLParam(r, "repo")
	prNumber := chi.URLParam(r, "pr")
	pr := fmt.Sprintf("%s/%s/%s", owner, repo, prNumber)
	deliveryID := chi.URLParam(r, "deliveryID")

	report, err := c.evaluationManager.GetReport(ctx, pr, deliveryID)
	if errors.Is(err, evaluation.ErrReportNotFound) {
		oplog.Err(err).Msg("error getting report")
		c.NotFound(w, r)
		return
	} else if err != nil {
		oplog.Err(err).Msg("error getting report")
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	navbar := pages.Navbar(Title)
	lottie := components.Lottie(report.Outcome.Result, report.Outcome.Err)
	EventDetails := pages.EventDetails(report)
	metadata := pages.Metadata(lottie, EventDetails)

	modules := make([]string, 0, len(report.Breakdown))
	for k := range report.Breakdown {
		modules = append(modules, k)
	}

	breakdownSection := pages.BreakdownSection(modules, report)
	inputSection := pages.InputSection(report.Input)
	body := pages.Container(pages.DetailsPage(navbar, metadata, breakdownSection, inputSection))
	html := templates.Html(Title, body)
	templ.Handler(html).ServeHTTP(w, r)
}

func (c *controller) NotFound(w http.ResponseWriter, r *http.Request) {
	lottie := components.LottieNotFound()
	body := pages.Container(lottie)
	html := templates.Html(Title, body)
	templ.Handler(html).ServeHTTP(w, r)
}
