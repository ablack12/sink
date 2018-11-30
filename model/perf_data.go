package model

import (
	"fmt"
	"time"

	"github.com/evergreen-ci/sink"
	"github.com/mongodb/anser/bsonutil"
	"github.com/mongodb/ftdc/events"
	"github.com/mongodb/grip"
	"github.com/montanaflynn/stats"
	"github.com/pkg/errors"
	mgo "gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

const (
	idField      = "_id"
	rollupsField = "rollups"
	nameField    = "name"
	valField     = "val"
	verField     = "version"
	userField    = "user"
	defaultVer   = 1
)

type PerfRollupValue struct {
	Name          string      `bson:"name"`
	Value         interface{} `bson:"val"`
	Version       int         `bson:"version"`
	UserSubmitted bool        `bson:"user"`
}

type PerfRollups struct {
	Stats       []PerfRollupValue `bson:"stats"`
	ProcessedAt time.Time         `bson:"processed_at"`
	Count       int               `bson:"count"`
	Valid       bool              `bson:"valid"`

	dirty     bool // nolint
	populated bool
	id        string
	env       sink.Environment
}

func (v *PerfRollupValue) getIntLong() (int64, error) {
	if val, ok := v.Value.(int64); ok {
		return val, nil
	} else if val, ok := v.Value.(int32); ok {
		return int64(val), nil
	} else if val, ok := v.Value.(int); ok {
		return int64(val), nil
	}
	return 0, errors.Errorf("mismatched type for name %s", v.Name)
}

func (v *PerfRollupValue) getFloat() (float64, error) {
	if val, ok := v.Value.(float64); ok {
		return val, nil
	}
	return 0, errors.Errorf("mismatched type for name %s", v.Name)
}

func (r *PerfRollups) Setup(env sink.Environment) {
	r.env = env
}

func (r *PerfRollups) Add(name string, version int, userSubmitted bool, value interface{}) error {
	if !r.populated {
		return errors.New("rollups have not been populated")
	}

	rollup := PerfRollupValue{
		Name:          name,
		Version:       version,
		Value:         value,
		UserSubmitted: userSubmitted,
	}

	err := r.updateExistingEntry(rollup)
	if err != nil && errors.Cause(err) == mgo.ErrNotFound {
		return r.insertNewEntry(rollup)
	}
	if err != nil {
		return errors.Wrap(err, "error updating entry") //NOTE: if the version is outdated, update should return an error I think
	}
	return nil
}

func (r *PerfRollups) insertNewEntry(rollup PerfRollupValue) error {
	conf, session, err := sink.GetSessionWithConfig(r.env)
	if err != nil {
		return errors.Wrap(err, "error connecting")
	}
	defer session.Close()

	insert := bson.M{rollupsField: rollup}
	search := bson.M{idField: r.id}
	c := session.DB(conf.DatabaseName).C(perfResultCollection)
	err = c.Update(search, bson.M{"$push": insert})
	if err != nil {
		return errors.Wrap(err, "error pushing new entry")
	}
	r.Stats = append(r.Stats, rollup)
	r.Count++
	return nil
}

func (r *PerfRollups) updateExistingEntry(rollup PerfRollupValue) error {
	conf, session, err := sink.GetSessionWithConfig(r.env)
	if err != nil {
		return errors.Wrap(err, "error connecting")
	}
	defer session.Close()
	search := bson.M{
		idField: r.id,
		bsonutil.GetDottedKeyName(rollupsField, nameField): rollup.Name,
	}
	update := bson.M{
		bsonutil.GetDottedKeyName(rollupsField, "$[elem]", verField):  rollup.Version,
		bsonutil.GetDottedKeyName(rollupsField, "$[elem]", valField):  rollup.Value,
		bsonutil.GetDottedKeyName(rollupsField, "$[elem]", userField): rollup.UserSubmitted,
	}
	arrayFilter := bson.M{
		bsonutil.GetDottedKeyName("elem", verField): bson.M{"$lt": rollup.Version},
	}
	change := mgo.Change{
		ArrayFilters: []bson.M{arrayFilter},
		Update:       bson.M{"$set": update},
	}

	c := session.DB(conf.DatabaseName).C(perfResultCollection)
	var result PerfRollups
	changeInfo, err := c.Find(search).Apply(change, &result)
	if err != nil {
		return errors.Wrap(err, "error updating an existing entry")
	}

	if changeInfo.Matched > 0 && changeInfo.Updated == 0 { // TODO: mess around with this
		return errors.New("error updating with outdated version")
	}
	for i := range r.Stats {
		if r.Stats[i].Name == rollup.Name {
			r.Stats[i].Version = rollup.Version
			r.Stats[i].Value = rollup.Value
			r.Stats[i].UserSubmitted = rollup.UserSubmitted
			return nil
		}
	}
	r.Stats = append(r.Stats, rollup)
	r.Count++
	return nil
}

func (r *PerfRollups) GetInt(name string) (int, error) {
	for _, rollup := range r.Stats {
		if rollup.Name == name {
			if val, ok := rollup.Value.(int); ok {
				return val, nil
			} else if val, ok := rollup.Value.(int32); ok {
				return int(val), nil
			} else {
				return 0, errors.Errorf("mismatched type for name %s", name)
			}
		}
	}
	return 0, errors.Errorf("name %s does not exist", name)
}

func (r *PerfRollups) GetInt32(name string) (int32, error) {
	val, err := r.GetInt(name)
	return int32(val), err
}

func (r *PerfRollups) GetInt64(name string) (int64, error) {
	for _, rollup := range r.Stats {
		if rollup.Name == name {
			return rollup.getIntLong()
		}
	}
	return 0, errors.Errorf("name %s does not exist", name)
}

func (r *PerfRollups) GetFloat(name string) (float64, error) {
	for _, rollup := range r.Stats {
		if rollup.Name == name {
			return rollup.getFloat()
		}
	}
	return 0, errors.Errorf("name %s does not exist", name)
}

func (r *PerfRollups) Validate() error {
	if len(r.Stats) != r.Count {
		return errors.New("number of stats and count of stats not equal")
	}
	return nil
}

func (r *PerfRollups) Map() map[string]int64 {
	result := make(map[string]int64)
	for _, rollup := range r.Stats {
		val, err := rollup.getIntLong()
		if err == nil {
			result[rollup.Name] = val
		}
	}
	return result
}

func (r *PerfRollups) MapFloat() map[string]float64 {
	result := make(map[string]float64)
	for _, rollup := range r.Stats {
		if val, err := rollup.getFloat(); err == nil {
			result[rollup.Name] = val
		} else if val, err := rollup.getIntLong(); err == nil {
			result[rollup.Name] = float64(val)
		}
	}
	return result
}

////////////////////////////////////////////////////////////////////////
//
// Legacy Rollup Calculations. TODO: port to new system asap

// PerformanceStatistics is an intermediate form used to calculate
// statistics from a sequence of reports. It is never persisted to the
// database, but provides access to application level aggregations.
type performanceStatistics struct {
	counters struct {
		operations stats.Float64Data
		size       stats.Float64Data
		errors     stats.Float64Data
	}

	totalCount struct {
		operations int64
		size       int64
		errors     int64
	}

	timers struct {
		duration stats.Float64Data
		total    stats.Float64Data
	}

	totalTime struct {
		duration time.Duration
		waiting  time.Duration
	}

	state struct {
		workers stats.Float64Data
		failed  bool
	}

	samples int
	span    time.Duration
}

type performanceMetricSummary struct {
	counters struct {
		operations float64
		size       float64
		errors     float64
	}

	totalCount struct {
		operations int64
		size       int64
		errors     int64
	}

	totalTime struct {
		duration time.Duration
		waiting  time.Duration
	}

	timers struct {
		duration float64
		total    float64
	}

	guages struct {
		workers float64
		failed  bool
	}

	span       time.Duration
	samples    int
	metricType string
}

// PerformanceTimeSeries provides an expanded, in-memory value
// reflecting the results of a single value. These series reflect,
// generally the data stored in the FTDC chunks that are persisted in
// off-line storage. While the FTDC data often requires additional data
// than a sequence of points, these series are the basis of the
// reporting that this application does across test values.
type PerformanceTimeSeries []*events.Performance

// statistics converts a a series into an intermediate format that we
// can use to calculate means/totals/etc to create default rollups
func (ts PerformanceTimeSeries) statistics() (*performanceStatistics, error) {
	out := &performanceStatistics{
		samples: len(ts),
	}

	out.counters.operations = make(stats.Float64Data, len(ts))
	out.counters.size = make(stats.Float64Data, len(ts))
	out.counters.errors = make(stats.Float64Data, len(ts))
	out.state.workers = make(stats.Float64Data, len(ts))
	out.timers.duration = make(stats.Float64Data, len(ts))
	out.timers.total = make(stats.Float64Data, len(ts))

	var lastPoint time.Time

	for idx, point := range ts {
		if idx == 0 {
			lastPoint = point.Timestamp
		} else if point.Timestamp.Before(lastPoint) {
			return nil, errors.New("data is not valid")
		} else {
			out.span += point.Timestamp.Sub(lastPoint)
			lastPoint = point.Timestamp
		}

		out.counters.operations[idx] = float64(point.Counters.Operations)
		out.counters.size[idx] = float64(point.Counters.Size)
		out.counters.errors[idx] = float64(point.Counters.Errors)
		out.state.workers[idx] = float64(point.Guages.Workers)

		out.totalCount.errors += point.Counters.Errors
		out.totalCount.operations += point.Counters.Operations
		out.totalCount.size += point.Counters.Size

		if point.Guages.Failed {
			out.state.failed = true
		}

		// Handle time differently: Negative duration values should be ignored
		if point.Timers.Duration > 0 {
			out.timers.duration[idx] = float64(point.Timers.Duration)
			out.totalTime.waiting += point.Timers.Duration
		}
		if point.Timers.Total > 0 {
			out.timers.total[idx] = float64(point.Timers.Total)
			out.totalTime.duration += point.Timers.Total
		}
	}

	return out, nil
}

func (perf *performanceStatistics) mean() (performanceMetricSummary, error) {
	var err error
	catcher := grip.NewBasicCatcher()
	out := performanceMetricSummary{
		samples:    perf.samples,
		span:       perf.span,
		metricType: "mean",
	}
	out.totalTime.duration = perf.totalTime.duration
	out.totalTime.waiting = perf.totalTime.waiting
	out.totalCount.errors = perf.totalCount.errors
	out.totalCount.size = perf.totalCount.size
	out.totalCount.operations = perf.totalCount.operations
	out.guages.failed = perf.state.failed

	out.counters.size, err = perf.counters.size.Mean()
	catcher.Add(err)

	out.counters.operations, err = perf.counters.operations.Mean()
	catcher.Add(err)

	out.counters.errors, err = perf.counters.errors.Mean()
	catcher.Add(err)

	out.guages.workers, err = perf.state.workers.Mean()
	catcher.Add(err)

	out.timers.duration, err = perf.timers.duration.Mean()
	catcher.Add(err)

	out.timers.total, err = perf.timers.total.Mean()
	catcher.Add(err)

	return out, catcher.Resolve()
}

func (perf *performanceStatistics) p90() (performanceMetricSummary, error) {
	return perf.percentile(90.0)
}

func (perf *performanceStatistics) p50() (performanceMetricSummary, error) {
	return perf.percentile(50.0)
}

func (perf *performanceStatistics) percentile(pval float64) (performanceMetricSummary, error) {
	var err error
	catcher := grip.NewBasicCatcher()
	out := performanceMetricSummary{
		samples:    perf.samples,
		span:       perf.span,
		metricType: fmt.Sprintf("percentile_%.2f", pval),
	}

	out.guages.failed = perf.state.failed
	out.totalTime.duration = perf.totalTime.duration
	out.totalTime.waiting = perf.totalTime.waiting
	out.totalCount.errors = perf.totalCount.errors
	out.totalCount.size = perf.totalCount.size
	out.totalCount.operations = perf.totalCount.operations

	out.counters.size, err = perf.counters.size.Percentile(pval)
	catcher.Add(err)

	out.counters.operations, err = perf.counters.operations.Percentile(pval)
	catcher.Add(err)

	out.counters.errors, err = perf.counters.errors.Percentile(pval)
	catcher.Add(err)

	out.guages.workers, err = perf.state.workers.Percentile(pval)
	catcher.Add(err)

	out.timers.duration, err = perf.timers.duration.Percentile(pval)
	catcher.Add(err)

	out.timers.total, err = perf.timers.total.Percentile(pval)
	catcher.Add(err)

	return out, catcher.Resolve()
}

func (r *PerformanceResult) UpdateThroughputOps(perf performanceMetricSummary) error {
	if float64(perf.span) == 0 {
		return errors.New("cannot divide by zero")
	}
	name := fmt.Sprintf("throughputOps_%s", perf.metricType)
	val := perf.counters.operations / perf.span.Seconds()
	// save to database
	err := r.Rollups.Add(name, defaultVer, false, val)
	if err != nil {
		return errors.Wrapf(err, "error calculating %s", name)
	}
	return nil
}

func (r *PerformanceResult) UpdateThroughputSize(perf performanceMetricSummary) error {
	if float64(perf.span) == 0 {
		return errors.New("cannot divide by zero")
	}
	name := fmt.Sprintf("throughputSize_%s", perf.metricType)
	val := perf.counters.size / perf.span.Seconds()
	err := r.Rollups.Add(name, defaultVer, false, val)
	if err != nil {
		return errors.Wrapf(err, "error calculating %s", name)
	}
	return nil
}

func (r *PerformanceResult) UpdateErrorRate(perf performanceMetricSummary) error {
	if float64(perf.span) == 0 {
		return errors.New("cannot divide by zero")
	}
	name := fmt.Sprintf("errorRate_%s", perf.metricType)
	val := perf.counters.errors / perf.span.Seconds()
	err := r.Rollups.Add(name, defaultVer, false, val)
	if err != nil {
		return errors.Wrapf(err, "error calculating %s", name)
	}
	return nil
}

// TotalTime stored in seconds
func (r *PerformanceResult) UpdateTotalTime(perf performanceMetricSummary) error {
	val := perf.totalTime.duration
	err := r.Rollups.Add("totalTime", defaultVer, false, val.Seconds())
	if err != nil {
		return errors.Wrap(err, "error calculating totalTime")
	}
	return nil
}

// Latency stored in seconds
func (r *PerformanceResult) UpdateLatency(perf performanceMetricSummary) error {
	if time.Duration(perf.totalCount.operations) == 0 {
		return errors.New("cannot divide by zero duration")
	}
	val := perf.totalTime.duration / time.Duration(perf.totalCount.operations)
	err := r.Rollups.Add("latency", defaultVer, false, val.Seconds())
	if err != nil {
		return errors.Wrap(err, "error calculating latency")
	}
	return nil
}

func (r *PerformanceResult) UpdateTotalSamples(perf performanceMetricSummary) error {
	err := r.Rollups.Add("totalSamples", defaultVer, false, perf.samples)
	if err != nil {
		return errors.Wrap(err, "error calculating totalSamples")
	}
	return nil
}

// UpdateDefaultRollups adds throughput and error rate means and 50th/90th percentiles.
// Additionally computes latency, total time, and total number of samples.
// These are either added or updated within r.Rollups.
func (r *PerformanceResult) UpdateDefaultRollups(ts PerformanceTimeSeries) error {
	perfStats, err := ts.statistics()
	if err != nil {
		return errors.WithStack(err)
	}
	catcher := grip.NewBasicCatcher()
	means, err := perfStats.mean()
	catcher.Add(err)
	percentile50, err := perfStats.p50()
	catcher.Add(err)
	percentile90, err := perfStats.p90()
	catcher.Add(err)
	if catcher.HasErrors() {
		return catcher.Resolve()
	}

	for _, perf := range []performanceMetricSummary{means, percentile50, percentile90} {
		err = r.UpdateErrorRate(perf)
		catcher.Add(err)
		err = r.UpdateThroughputOps(perf)
		catcher.Add(err)
		err = r.UpdateThroughputSize(perf)
		catcher.Add(err)
	}

	err = r.UpdateTotalSamples(means)
	catcher.Add(err)
	err = r.UpdateTotalTime(means)
	catcher.Add(err)
	err = r.UpdateLatency(means)
	catcher.Add(err)
	return catcher.Resolve()
}
