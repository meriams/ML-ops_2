package rules

import (
	"context"
	"fmt"
	"sort"
	"sync"
	"time"

	"github.com/go-kit/log"
	opentracing "github.com/opentracing/opentracing-go"
	plabels "github.com/prometheus/prometheus/model/labels"
	"go.uber.org/zap"
)

// PromRuleTask is a promql rule executor
type PromRuleTask struct {
	name                 string
	file                 string
	frequency            time.Duration
	rules                []Rule
	seriesInPreviousEval []map[string]plabels.Labels // One per Rule.
	staleSeries          []plabels.Labels
	opts                 *ManagerOptions
	mtx                  sync.Mutex
	evaluationDuration   time.Duration
	evaluationTime       time.Duration
	lastEvaluation       time.Time

	markStale   bool
	done        chan struct{}
	terminated  chan struct{}
	managerDone chan struct{}

	pause  bool
	logger log.Logger
	notify NotifyFunc
}

// newPromRuleTask holds rules that have promql condition
// and evalutes the rule at a given frequency
func newPromRuleTask(name, file string, frequency time.Duration, rules []Rule, opts *ManagerOptions, notify NotifyFunc) *PromRuleTask {
	zap.S().Info("Initiating a new rule group:", name, "\t frequency:", frequency)

	if time.Now() == time.Now().Add(frequency) {
		frequency = DefaultFrequency
	}

	return &PromRuleTask{
		name:                 name,
		file:                 file,
		pause:                false,
		frequency:            frequency,
		rules:                rules,
		opts:                 opts,
		seriesInPreviousEval: make([]map[string]plabels.Labels, len(rules)),
		done:                 make(chan struct{}),
		terminated:           make(chan struct{}),
		notify:               notify,
		logger:               log.With(opts.Logger, "group", name),
	}
}

// Name returns the group name.
func (g *PromRuleTask) Name() string { return g.name }

// Key returns the group key
func (g *PromRuleTask) Key() string {
	return g.name + ";" + g.file
}

func (g *PromRuleTask) Type() TaskType { return TaskTypeProm }

// Rules returns the group's rules.
func (g *PromRuleTask) Rules() []Rule { return g.rules }

// Interval returns the group's interval.
func (g *PromRuleTask) Interval() time.Duration { return g.frequency }

func (g *PromRuleTask) Pause(b bool) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.pause = b
}

func (g *PromRuleTask) Run(ctx context.Context) {
	defer close(g.terminated)

	// Wait an initial amount to have consistently slotted intervals.
	evalTimestamp := g.EvalTimestamp(time.Now().UnixNano()).Add(g.frequency)
	select {
	case <-time.After(time.Until(evalTimestamp)):
	case <-g.done:
		return
	}

	ctx = NewQueryOriginContext(ctx, map[string]interface{}{
		"ruleGroup": map[string]string{
			"name": g.Name(),
		},
	})

	iter := func() {

		start := time.Now()
		g.Eval(ctx, evalTimestamp)
		timeSinceStart := time.Since(start)

		g.setEvaluationTime(timeSinceStart)
		g.setLastEvaluation(start)
	}

	// The assumption here is that since the ticker was started after having
	// waited for `evalTimestamp` to pass, the ticks will trigger soon
	// after each `evalTimestamp + N * g.frequency` occurrence.
	tick := time.NewTicker(g.frequency)
	defer tick.Stop()

	// defer cleanup
	defer func() {
		if !g.markStale {
			return
		}
		go func(now time.Time) {
			for _, rule := range g.seriesInPreviousEval {
				for _, r := range rule {
					g.staleSeries = append(g.staleSeries, r)
				}
			}
			// That can be garbage collected at this point.
			g.seriesInPreviousEval = nil

		}(time.Now())

	}()

	iter()

	// let the group iterate and run
	for {
		select {
		case <-g.done:
			return
		default:
			select {
			case <-g.done:
				return
			case <-tick.C:
				missed := (time.Since(evalTimestamp) / g.frequency) - 1
				evalTimestamp = evalTimestamp.Add((missed + 1) * g.frequency)
				iter()
			}
		}
	}
}

func (g *PromRuleTask) Stop() {
	close(g.done)
	<-g.terminated
}

func (g *PromRuleTask) hash() uint64 {
	l := plabels.New(
		plabels.Label{Name: "name", Value: g.name},
	)
	return l.Hash()
}

// PromRules returns the list of the group's promql rules.
func (g *PromRuleTask) PromRules() []*PromRule {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	var alerts []*PromRule
	for _, rule := range g.rules {
		if tr, ok := rule.(*PromRule); ok {
			alerts = append(alerts, tr)
		}
	}
	sort.Slice(alerts, func(i, j int) bool {
		return alerts[i].State() > alerts[j].State() ||
			(alerts[i].State() == alerts[j].State() &&
				alerts[i].Name() < alerts[j].Name())
	})
	return alerts
}

// HasAlertingRules returns true if the group contains at least one AlertingRule.
func (g *PromRuleTask) HasAlertingRules() bool {
	g.mtx.Lock()
	defer g.mtx.Unlock()

	for _, rule := range g.rules {
		if _, ok := rule.(*ThresholdRule); ok {
			return true
		}
	}
	return false
}

// GetEvaluationDuration returns the time in seconds it took to evaluate the rule group.
func (g *PromRuleTask) GetEvaluationDuration() time.Duration {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.evaluationDuration
}

// SetEvaluationDuration sets the time in seconds the last evaluation took.
func (g *PromRuleTask) SetEvaluationDuration(dur time.Duration) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.evaluationDuration = dur
}

// GetEvaluationTime returns the time in seconds it took to evaluate the rule group.
func (g *PromRuleTask) GetEvaluationTime() time.Duration {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.evaluationTime
}

// setEvaluationTime sets the time in seconds the last evaluation took.
func (g *PromRuleTask) setEvaluationTime(dur time.Duration) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.evaluationTime = dur
}

// GetLastEvaluation returns the time the last evaluation of the rule group took place.
func (g *PromRuleTask) GetLastEvaluation() time.Time {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	return g.lastEvaluation
}

// setLastEvaluation updates evaluationTimestamp to the timestamp of when the rule group was last evaluated.
func (g *PromRuleTask) setLastEvaluation(ts time.Time) {
	g.mtx.Lock()
	defer g.mtx.Unlock()
	g.lastEvaluation = ts
}

// EvalTimestamp returns the immediately preceding consistently slotted evaluation time.
func (g *PromRuleTask) EvalTimestamp(startTime int64) time.Time {
	var (
		offset = int64(g.hash() % uint64(g.frequency))
		adjNow = startTime - offset
		base   = adjNow - (adjNow % int64(g.frequency))
	)

	return time.Unix(0, base+offset).UTC()
}

// CopyState copies the alerting rule and staleness related state from the given group.
//
// Rules are matched based on their name and labels. If there are duplicates, the
// first is matched with the first, second with the second etc.
func (g *PromRuleTask) CopyState(fromTask Task) error {

	from, ok := fromTask.(*PromRuleTask)
	if !ok {
		return fmt.Errorf("you can only copy rule groups with same type")
	}

	g.evaluationTime = from.evaluationTime
	g.lastEvaluation = from.lastEvaluation

	ruleMap := make(map[string][]int, len(from.rules))

	for fi, fromRule := range from.rules {
		nameAndLabels := nameAndLabels(fromRule)
		l := ruleMap[nameAndLabels]
		ruleMap[nameAndLabels] = append(l, fi)
	}

	for i, rule := range g.rules {
		nameAndLabels := nameAndLabels(rule)
		indexes := ruleMap[nameAndLabels]
		if len(indexes) == 0 {
			continue
		}
		fi := indexes[0]
		g.seriesInPreviousEval[i] = from.seriesInPreviousEval[fi]
		ruleMap[nameAndLabels] = indexes[1:]

		ar, ok := rule.(*ThresholdRule)
		if !ok {
			continue
		}
		far, ok := from.rules[fi].(*ThresholdRule)
		if !ok {
			continue
		}

		for fp, a := range far.active {
			ar.active[fp] = a
		}
	}

	// Handle deleted and unmatched duplicate rules.
	g.staleSeries = from.staleSeries
	for fi, fromRule := range from.rules {
		nameAndLabels := nameAndLabels(fromRule)
		l := ruleMap[nameAndLabels]
		if len(l) != 0 {
			for _, series := range from.seriesInPreviousEval[fi] {
				g.staleSeries = append(g.staleSeries, series)
			}
		}
	}
	return nil
}

// Eval runs a single evaluation cycle in which all rules are evaluated sequentially.
func (g *PromRuleTask) Eval(ctx context.Context, ts time.Time) {
	zap.S().Info("promql rule task:", g.name, "\t eval started at:", ts)
	for i, rule := range g.rules {
		if rule == nil {
			continue
		}
		select {
		case <-g.done:
			return
		default:
		}

		func(i int, rule Rule) {
			sp, ctx := opentracing.StartSpanFromContext(ctx, "rule")

			sp.SetTag("name", rule.Name())
			defer func(t time.Time) {
				sp.Finish()

				since := time.Since(t)
				rule.SetEvaluationDuration(since)
				rule.SetEvaluationTimestamp(t)
			}(time.Now())

			_, err := rule.Eval(ctx, ts, g.opts.Queriers)
			if err != nil {
				rule.SetHealth(HealthBad)
				rule.SetLastError(err)

				zap.S().Warn("msg", "Evaluating rule failed", "rule", rule, "err", err)

				// Canceled queries are intentional termination of queries. This normally
				// happens on shutdown and thus we skip logging of any errors here.
				//! if _, ok := err.(promql.ErrQueryCanceled); !ok {
				//	level.Warn(g.logger).Log("msg", "Evaluating rule failed", "rule", rule, "err", err)
				//}
				return
			}
			rule.SendAlerts(ctx, ts, g.opts.ResendDelay, g.frequency, g.notify)

		}(i, rule)
	}
}
