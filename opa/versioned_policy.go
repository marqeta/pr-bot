package opa

import (
	"context"
	"errors"
	"strings"

	"github.com/marqeta/pr-bot/opa/client"
	"github.com/marqeta/pr-bot/opa/input"
	"github.com/marqeta/pr-bot/opa/rules"
	"github.com/marqeta/pr-bot/opa/types"
)

const SchemaRuleName = "schema"

var ErrInvalidSchemaVersion = errors.New("invalid schema version")

type versionedPolicy struct {
	versions map[string]Policy
	schema   rules.Rules[string]
}

// Evaluate implements Policy.
func (p *versionedPolicy) Evaluate(
	ctx context.Context, module string, ip *input.Model) (types.Result, error) {

	schema, err := p.schema.Evaluate(ctx, module, ip)
	if err != nil {
		return types.Result{}, err
	}
	schema = strings.ToLower(strings.TrimSpace(schema))
	policy, ok := p.versions[schema]
	if !ok {
		return types.Result{}, ErrInvalidSchemaVersion
	}
	return policy.Evaluate(ctx, module, ip)
}

func NewVersionedPolicyFromRules(versions map[string]Policy, schema rules.Rules[string]) Policy {
	return &versionedPolicy{
		versions: versions,
		schema:   schema,
	}
}

func NewVersionedPolicy(versions map[string]Policy, client client.Client) Policy {
	return NewVersionedPolicyFromRules(
		versions,
		rules.NewSchema(SchemaRuleName, client),
	)
}
