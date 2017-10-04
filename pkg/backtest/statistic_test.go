package backtest

import (
	"errors"
	"reflect"
	"testing"
	"time"
)

func TestTotalEquityReturn(t *testing.T) {
	// set up test cases for getting first equity point
	var testCases = []struct {
		msg       string
		stat      Statistic
		expReturn float64
		expErr    error
	}{
		{"testing for multiple entryPoints",
			Statistic{
				equity: []equityPoint{
					{equity: 100, equityReturn: 0},
					{equity: 120, equityReturn: 0.2},
				},
			},
			0.2,
			nil},
		{"testing for multiple entryPoints with same value",
			Statistic{
				equity: []equityPoint{
					{equity: 100, equityReturn: 0},
					{equity: 100, equityReturn: 0},
				},
			},
			0,
			nil},
		{"testing for last entryPoints with 0 equity",
			Statistic{
				equity: []equityPoint{
					{equity: 100, equityReturn: 0.1},
					{equity: 0, equityReturn: -1},
				},
			},
			-1,
			nil},
		{"testing for nil entryPoints",
			Statistic{},
			0,
			errors.New("could not calculate totalEquityReturn, no equity points found")},
	}

	for _, tc := range testCases {
		r, err := tc.stat.TotalEquityReturn()
		if (r != tc.expReturn) || (reflect.TypeOf(err) != reflect.TypeOf(tc.expErr)) {
			t.Errorf("%v TotalEquityReturn(): \nexpected %#v %#v, \nactual   %#v %#v",
				tc.msg, tc.expReturn, tc.expErr, r, err)
		}
	}
}
func TestGetEquityPoint(t *testing.T) {
	var statCases = map[string]Statistic{
		"multiple": {
			equity: []equityPoint{
				{equity: 100, equityReturn: 0.1},
				{equity: 110, equityReturn: 0.2},
				{equity: 120, equityReturn: 0.3},
			},
		},
		"single": {
			equity: []equityPoint{
				{equity: 150, equityReturn: 0.25},
			},
		},
		"empty": {
			equity: []equityPoint{},
		},
	}

	// define test cases struct
	type testCase struct {
		msg   string
		stat  Statistic
		expEP equityPoint
		expOk bool
	}

	// set up test cases for getting first equity point
	var testCasesFirst = []testCase{
		{"testing first for multiple entryPoints",
			statCases["multiple"],
			equityPoint{equity: 100, equityReturn: 0.1},
			true},
		{"testing first for single entryPoints",
			statCases["single"],
			equityPoint{equity: 150, equityReturn: 0.25},
			true},
		{"testing first for nil entryPoints",
			statCases["empty"],
			equityPoint{},
			false},
	}

	for _, tc := range testCasesFirst {
		ep, ok := tc.stat.firstEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}

	// set up test cases for getting last equity point
	var testCasesLast = []testCase{
		{"testing last for multiple entryPoints",
			statCases["multiple"],
			equityPoint{equity: 120, equityReturn: 0.3},
			true},
		{"testing last for single entryPoints",
			statCases["single"],
			equityPoint{equity: 150, equityReturn: 0.25},
			true},
		{"testing last for nil entryPoints",
			statCases["empty"],
			equityPoint{},
			false},
	}

	for _, tc := range testCasesLast {
		ep, ok := tc.stat.lastEquityPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (ok != tc.expOk) {
			t.Errorf("%v firstEquityPoint(): \nexpected %#v %v, \nactual   %#v %v",
				tc.msg, tc.expEP, tc.expOk, ep, ok)
		}
	}
}

func TestCalcEquityReturn(t *testing.T) {
	var testCases = []struct {
		msg   string
		stat  Statistic
		ep    equityPoint
		expEP equityPoint
	}{
		{"testing equity return with single equity points",
			Statistic{
				equity: []equityPoint{
					{equity: 100},
				},
			},
			equityPoint{equity: 90},
			equityPoint{
				equity:       90,
				equityReturn: -0.1,
			},
		},
		{"testing equity return with multiple equity points",
			Statistic{
				equity: []equityPoint{
					{equity: 100},
					{equity: 90},
					{equity: 110},
				},
			},
			equityPoint{equity: 100},
			equityPoint{
				equity:       100,
				equityReturn: -0.0909,
			},
		},
		{"testing equity return with single equity points but 0 equity",
			Statistic{
				equity: []equityPoint{
					{equity: 0},
				},
			},
			equityPoint{equity: 100},
			equityPoint{
				equity:       100,
				equityReturn: 1,
			},
		},
		{"testing equity return with nil equity points",
			Statistic{
				equity: []equityPoint{},
			},
			equityPoint{equity: 100},
			equityPoint{
				equity:   100,
				drawdown: 0,
			},
		},
	}

	for _, tc := range testCases {
		ep := tc.stat.calcEquityReturn(tc.ep)
		if !reflect.DeepEqual(ep, tc.expEP) {
			t.Errorf("%v calcEquityReturn(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.ep, tc.expEP, ep)
		}
	}

}

func TestCalcDrawdown(t *testing.T) {
	var testCases = []struct {
		msg   string
		stat  Statistic
		ep    equityPoint
		expEP equityPoint
	}{
		{"testing drawdown with simple high equityPoint",
			Statistic{
				high: equityPoint{equity: 100},
			},
			equityPoint{equity: 90},
			equityPoint{
				equity:   90,
				drawdown: -0.1,
			},
		},
		{"testing drawdown with simple high equityPoint equal equity",
			Statistic{
				high: equityPoint{equity: 100},
			},
			equityPoint{equity: 100},
			equityPoint{
				equity:   100,
				drawdown: 0,
			},
		},
		{"testing drawdown with simple high equityPoint lower equity",
			Statistic{
				high: equityPoint{equity: 90},
			},
			equityPoint{equity: 100},
			equityPoint{
				equity:   100,
				drawdown: 0,
			},
		},
		{"testing drawdown with empty high equityPoint",
			Statistic{},
			equityPoint{equity: 100},
			equityPoint{
				equity:   100,
				drawdown: 0,
			},
		},
	}

	for _, tc := range testCases {
		ep := tc.stat.calcDrawdown(tc.ep)
		if !reflect.DeepEqual(ep, tc.expEP) {
			t.Errorf("%v calcDrawdown(%v): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.ep, tc.expEP, ep)
		}
	}

}

func TestMaxDrawdown(t *testing.T) {
	var time1, _ = time.Parse("2006-01-02", "2017-09-25")
	var time2, _ = time.Parse("2006-01-02", "2017-09-26")
	var time3, _ = time.Parse("2006-01-02", "2017-09-27")
	var time4, _ = time.Parse("2006-01-02", "2017-09-28")
	var time5, _ = time.Parse("2006-01-02", "2017-09-29")

	// set up test cases for getting the max drawdown point
	var testCases = []struct {
		msg     string
		stat    Statistic
		expEP   equityPoint
		expInt  int
		expMax  float64
		expTime time.Time
		expDur  float64 // duration in hours
	}{
		{"testing maxdrawdown for multiple entryPoints",
			Statistic{
				equity: []equityPoint{
					{timestamp: time1, equity: 100, drawdown: 0},
					{timestamp: time2, equity: 110, drawdown: 0},
					{timestamp: time3, equity: 105, drawdown: -0.0455},
					{timestamp: time4, equity: 95, drawdown: -0.1364},
					{timestamp: time5, drawdown: 0},
				},
			},
			equityPoint{timestamp: time4, equity: 95, drawdown: -0.1364},
			3,
			-0.1364,
			time4,
			48,
		},
		{"testing maxdrawdown for single entryPoints",
			Statistic{
				equity: []equityPoint{
					{timestamp: time1, equity: 100, drawdown: 0},
				},
			},
			equityPoint{timestamp: time1, equity: 100, drawdown: 0},
			0,
			0,
			time1,
			0,
		},
		{"testing maxdrawdown for nil entryPoints",
			Statistic{},
			equityPoint{},
			0,
			0,
			time.Time{},
			0,
		},
	}

	// testing for max drawdown equity point
	for _, tc := range testCases {
		i, ep := tc.stat.maxDrawdownPoint()
		if !reflect.DeepEqual(ep, tc.expEP) || (i != tc.expInt) {
			t.Errorf("%v maxDrawdownPoint(): \nexpected %d %#v, \nactual   %d %#v",
				tc.msg, tc.expInt, tc.expEP, i, ep)
		}
	}

	// testing for max drawdown value
	for _, tc := range testCases {
		max := tc.stat.MaxDrawdown()
		if !reflect.DeepEqual(max, tc.expMax) {
			t.Errorf("%v MaxDrawdown(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expMax, max)
		}
	}

	// testing for max drawdown time
	for _, tc := range testCases {
		time := tc.stat.MaxDrawdownTime()
		if !reflect.DeepEqual(time, tc.expTime) {
			t.Errorf("%v MaxDrawdownTime(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expTime, time)
		}
	}

	// testing for max drawdown duration
	for _, tc := range testCases {
		duration := tc.stat.MaxDrawdownDuration()
		if !reflect.DeepEqual(duration.Hours(), tc.expDur) {
			t.Errorf("%v MaxDrawdownDuration(): \nexpected %#v, \nactual   %#v",
				tc.msg, tc.expDur, duration.Hours())
		}
	}
}
