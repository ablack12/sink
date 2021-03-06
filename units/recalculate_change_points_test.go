package units

import (
	"context"
	"strconv"
	"testing"
	"time"

	"github.com/evergreen-ci/cedar"
	"github.com/evergreen-ci/cedar/model"
	"github.com/evergreen-ci/cedar/perf"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"gopkg.in/mgo.v2/bson"
)

type TestResultsAndRollups struct {
	info    *model.PerformanceResultInfo
	rollups []model.PerfRollupValue
}

func makePerfResults(length int, breakpoint int, unique string) []TestResultsAndRollups {
	var rollups []TestResultsAndRollups
	i := 0
	for i < length {
		newRollup := TestResultsAndRollups{
			info: &model.PerformanceResultInfo{
				Project:  "project" + unique,
				Variant:  "variant",
				Version:  "version" + strconv.Itoa(i+1),
				Order:    i + 1,
				TestName: "test",
				TaskName: "task",
				Mainline: true,
			},
			rollups: []model.PerfRollupValue{
				{
					Name:       "measurement",
					MetricType: model.MetricTypeSum,
					Version:    0,
				},
				{
					Name:       "measurement_another",
					MetricType: model.MetricTypeSum,
					Version:    0,
				},
			},
		}
		if i < breakpoint {
			newRollup.rollups[0].Value = float64(100)
			newRollup.rollups[1].Value = float64(100)
		} else {
			newRollup.rollups[0].Value = float64(1000)
			newRollup.rollups[1].Value = float64(1000)
		}
		rollups = append(rollups, newRollup)
		i++
	}
	return rollups
}

func init() {
	dbName := "test_cedar_signal_processing"
	env, err := cedar.NewEnvironment(context.Background(), dbName, &cedar.Configuration{
		MongoDBURI:    "mongodb://localhost:27017",
		DatabaseName:  dbName,
		SocketTimeout: time.Minute,
		NumWorkers:    2,
	})
	if err != nil {
		panic(err)
	}
	cedar.SetEnvironment(env)
}

func tearDown(env cedar.Environment) error {
	conf, session, err := cedar.GetSessionWithConfig(env)
	if err != nil {
		return errors.WithStack(err)
	}
	defer session.Close()
	return errors.WithStack(session.DB(conf.DatabaseName).DropDatabase())
}

type MockDetector struct {
	Calls [][]float64
}

func (m *MockDetector) Algorithm() perf.Algorithm {
	return perf.CreateDefaultAlgorithm()
}

func (m *MockDetector) DetectChanges(ctx context.Context, series []float64) ([]int, error) {
	m.Calls = append(m.Calls, series)
	last := series[0]
	var cps []int
	for idx, i := range series {
		if i != last {
			last = i
			cps = append(cps, idx)
		}
	}
	return cps, nil
}

func TestRecalculateChangePointsJob(t *testing.T) {
	env := cedar.GetEnvironment()
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	_ = env.GetDB().Drop(ctx)

	rollups := append(makePerfResults(100, 50, "a"), makePerfResults(100, 50, "b")...)
	for _, result := range rollups {
		performanceResult := model.CreatePerformanceResult(*result.info, nil, result.rollups)
		performanceResult.CreatedAt = time.Now().Add(time.Second * -1)
		performanceResult.Setup(env)
		err := performanceResult.SaveNew(ctx)
		if err != nil {
			panic(err)
		}
	}
	defer func() {
		assert.NoError(t, tearDown(env))
	}()

	t.Run("Recalculates", func(t *testing.T) {
		j := NewRecalculateChangePointsJob(model.PerformanceResultSeriesID{
			Project: "projecta",
			Variant: "variant",
			Task:    "task",
			Test:    "test",
		})
		mockDetector := &MockDetector{}
		j.(*recalculateChangePointsJob).changePointDetector = mockDetector
		j.Run(ctx)
		assert.True(t, j.Status().Completed)
		assert.Len(t, mockDetector.Calls, 2)
		var result []model.PerformanceResult
		filter := bson.M{
			"analysis.change_points": bson.M{"$ne": []struct{}{}},
		}
		res, err := env.GetDB().Collection("perf_results").Find(ctx, filter)
		require.NoError(t, err)
		assert.NoError(t, res.All(ctx, &result))
		require.Len(t, result, 1)
		require.Len(t, result[0].Analysis.ChangePoints, 2)
		require.NotEqual(t, result[0].Analysis.ProcessedAt, time.Time{})

		var options []model.AlgorithmOption

		for _, v := range mockDetector.Algorithm().Configuration() {
			options = append(options, model.AlgorithmOption{
				Name:  v.Name,
				Value: v.Value,
			})
		}

		//Change point 1
		require.Contains(t, []string{"measurement", "measurement_another"}, result[0].Analysis.ChangePoints[0].Measurement)
		require.Equal(t, result[0].Analysis.ChangePoints[0].Algorithm.Name, mockDetector.Algorithm().Name())
		require.Equal(t, result[0].Analysis.ChangePoints[0].Algorithm.Version, mockDetector.Algorithm().Version())
		for _, v := range options {
			require.Contains(t, result[0].Analysis.ChangePoints[0].Algorithm.Options, v)
		}
		require.Equal(t, result[0].Analysis.ChangePoints[0].Triage.TriagedOn, time.Time{})
		require.Equal(t, result[0].Analysis.ChangePoints[0].Triage.Status, model.TriageStatusUntriaged)

		//Change point 2
		require.Contains(t, []string{"measurement", "measurement_another"}, result[0].Analysis.ChangePoints[1].Measurement)
		require.Equal(t, result[0].Analysis.ChangePoints[1].Algorithm.Name, mockDetector.Algorithm().Name())
		require.Equal(t, result[0].Analysis.ChangePoints[1].Algorithm.Version, mockDetector.Algorithm().Version())
		for _, v := range options {
			require.Contains(t, result[0].Analysis.ChangePoints[1].Algorithm.Options, v)
		}
		require.Equal(t, result[0].Analysis.ChangePoints[1].Triage.TriagedOn, time.Time{})
		require.Equal(t, result[0].Analysis.ChangePoints[1].Triage.Status, model.TriageStatusUntriaged)
	})

	t.Run("DoesNothingWhenDisabled", func(t *testing.T) {
		j := NewRecalculateChangePointsJob(model.PerformanceResultSeriesID{
			Project: "projecta",
			Variant: "variant",
			Task:    "task",
			Test:    "test",
		})
		mockDetector := &MockDetector{}
		job := j.(*recalculateChangePointsJob)
		job.changePointDetector = mockDetector
		job.conf = model.NewCedarConfig(env)
		job.conf.Flags.DisableSignalProcessing = true
		j.Run(ctx)
		assert.True(t, j.Status().Completed)
		assert.Len(t, mockDetector.Calls, 0)
	})
}
