package evaluation

import (
	"github.com/marqeta/pr-bot/opa/input"
)

//go:generate mockery --name ReportBuilder
type ReportBuilder interface {
	AddModuleResult(module string, result Result)
	SetInput(input *input.Model)
	SetOutcome(result Result)
	GetReport() Report
}

type reportBuilder struct {
	report Report
}

// SetInput implements Builder.
func (b *reportBuilder) SetInput(input *input.Model) {
	b.report.Input = input
	b.report.Event = input.Event
	b.report.Action = input.Action
	b.report.Title = input.PullRequest.GetTitle()
	b.report.Author = input.PullRequest.GetUser().GetLogin()
}

// AddModuleResult implements Builder.
func (b *reportBuilder) AddModuleResult(module string, result Result) {
	b.report.Breakdown[module] = result
}

// SetOutcome implements Builder.
func (b *reportBuilder) SetOutcome(outcome Result) {
	b.report.Outcome = outcome
}

// GetReport implements Builder.
func (b *reportBuilder) GetReport() Report {
	return b.report
}
