/*
 * The MIT License (MIT)
 * ------------------------------------------------------------------
 * Copyright © 2020 fufeng.All Rights Reserved.
 *
 * ProjectName: thinking-in-go
 * @Author : <a href="https://github.com/lcy2013">MagicLuo(扶风)</a>
 * @date : 2020-09-28
 * @version : 1.0.0-RELEASE
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy of this software and associated documentation files (the “Software”), to deal in the Software without restriction, including without limitation the rights to use, copy, modify, merge, publish, distribute, sublicense, and/or sell copies of the Software, and to permit persons to whom the Software is furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED “AS IS”, WITHOUT WARRANTY OF ANY KIND, EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.
 *
 */
package model

import (
	"reflect"
	"testing"
)

func TestItinerary_CreateEmpty(t *testing.T) {
	i := Itinerary{}

	var legs []Leg

	if !reflect.DeepEqual(i.Legs, legs) {
		t.Errorf("should be equal")
	}
	if i.InitialDepartureLocation() != "" {
		t.Errorf("InitialDepartureLocation() = %s; want = %s",
			i.InitialDepartureLocation(), "")
	}
	if i.FinalArrivalLocation() != "" {
		t.Errorf("FinalArrivalLocation() = %s; want = %s",
			i.FinalArrivalLocation(), "")
	}
}

func TestItinerary_IsExpected_EmptyItinerary(t *testing.T) {
	i := Itinerary{}
	e := HandlingEvent{}

	if got, want := i.IsExpected(e), true; got != want {
		t.Errorf("IsExpected() = %v; want = %v", got, want)
	}
}

type eventExpectedTest struct {
	act HandlingActivity
	exp bool
}

var eventExpectedTests = []eventExpectedTest{
	{HandlingActivity{}, true},
	{HandlingActivity{Type: Receive, Location: SESTO}, true},
	{HandlingActivity{Type: Receive, Location: AUMEL}, false},
	{HandlingActivity{Type: Load, Location: AUMEL, VoyageNumber: "001A"}, true},
	{HandlingActivity{Type: Load, Location: CNHKG, VoyageNumber: "001A"}, false},
	{HandlingActivity{Type: Unload, Location: CNHKG, VoyageNumber: "001A"}, true},
	{HandlingActivity{Type: Unload, Location: SESTO, VoyageNumber: "001A"}, false},
	{HandlingActivity{Type: Claim, Location: CNHKG}, true},
	{HandlingActivity{Type: Claim, Location: SESTO}, false},
}

func TestItinerary_IsExpected(t *testing.T) {
	i := Itinerary{Legs: []Leg{
		{
			VoyageNumber:   "001A",
			LoadLocation:   SESTO,
			UnloadLocation: AUMEL,
		},
		{
			VoyageNumber:   "001A",
			LoadLocation:   AUMEL,
			UnloadLocation: CNHKG,
		},
	}}

	for _, tt := range eventExpectedTests {
		e := HandlingEvent{
			Activity: tt.act,
		}

		if got := i.IsExpected(e); got != tt.exp {
			t.Errorf("IsExpected() = %v; want = %v", got, tt.exp)
		}
	}
}
