package opa

import (
	"context"
	"fmt"

	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/httplog"
	"github.com/marqeta/pr-bot/opa/evaluation"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/types"
)

// Evaluator evaluates the policy for each module in the bundle.
//
//go:generate mockery --name Evaluator
type Evaluator interface {
	Evaluate(ctx context.Context, input input.GHE) (types.Result, error)
}

type evaluator struct {
	modules      []string
	policy       Policy
	inputFactory input.Factory
	manager      evaluation.Manager
}

func NewEvaluator(modules []string, policy Policy, factory input.Factory, manager evaluation.Manager) Evaluator {
	return &evaluator{
		modules:      modules,
		policy:       policy,
		inputFactory: factory,
		manager:      manager,
	}
}

// Evaluate implements Evaluator.
// Evaluate evaluates the policy for each module in the bundle.
// returns coalesced result.
// RequestChanges > comment > approve > skip
func (e *evaluator) Evaluate(ctx context.Context, ghe input.GHE) (types.Result, error) {

	oplog := httplog.LogEntry(ctx)
	coalesced := types.Result{}
	model, err := e.inputFactory.CreateModel(ctx, ghe)
	if err != nil {
		oplog.Err(err).Msg("failed to create input model for policy evaluation")
		return types.Result{}, err
	}
	report := e.newReportBuilder(ctx, ghe)
	defer e.storeReport(ctx, report)
	report.SetInput(model)
	for _, module := range e.modules {
		result, err := e.policy.Evaluate(ctx, module, model)
		report.AddModuleResult(module, evaluation.Result{
			Result: result,
			Err:    err,
		})
		if err != nil {
			oplog.Err(err).Msgf("failed to evaluate policy for module %s", module)
			// single module eval failed stop evaluating other modules
			// module A -> request changes
			// module B -> approve
			// coalesced result should be request changes
			// if module A start failing after an update to policies
			// we should not start approving the PR
			// therefore fail the entire evaluation if a single module fails
			report.SetOutcome(evaluation.Result{
				Err: err,
			})
			return types.Result{}, err
		}
		oplog.Info().Msgf("%v policy evaluation result: %+v", module, result)
		if !result.Track {
			oplog.Info().Msgf("track is false, skipping result for module %s", module)
			// module ignores a PR, need to continue to evaluate other modules
			continue
		}
		if result.Review.Type >= coalesced.Review.Type {
			// higher priority review
			coalesced = result
		}
	}

	report.SetOutcome(evaluation.Result{
		Result: coalesced,
		Err:    err,
	})
	oplog.Info().Interface("coalesced result", coalesced).Msg("coalesced policy evaluation result")
	return coalesced, nil
}

func (e *evaluator) newReportBuilder(ctx context.Context, ghe input.GHE) evaluation.ReportBuilder {
	reqID := middleware.GetReqID(ctx)
	deliveryID := evaluation.GetDeliveryID(ctx)
	pr := fmt.Sprintf("%s/%d", ghe.Repository.GetFullName(), ghe.PullRequest.GetNumber())
	return e.manager.NewReportBuilder(ctx, pr, reqID, deliveryID)
}

func (e *evaluator) storeReport(ctx context.Context, report evaluation.ReportBuilder) {
	oplog := httplog.LogEntry(ctx)
	err := e.manager.StoreReport(ctx, report)
	if err != nil {
		oplog.Err(err).Msg("failed to store policy evaluation report")
	}
}
