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
	"testing"
	"time"
)

func TestConstruction(t *testing.T) {
	id := NextTrackingID()
	spec := RouteSpecification{
		Origin:          SESTO,
		Destination:     AUMEL,
		ArrivalDeadline: time.Date(2009, time.March, 13, 0, 0, 0, 0, time.UTC),
	}

	c := NewCargo(id, spec)

	if c.Delivery.RoutingStatus != NotRouted {
		t.Errorf("RoutingStatus = %v; want = %v",
			c.Delivery.RoutingStatus, NotRouted)
	}
	if c.Delivery.TransportStatus != NotReceived {
		t.Errorf("TransportStatus = %v; want = %v",
			c.Delivery.TransportStatus, NotReceived)
	}
	if c.Delivery.LastKnownLocation != "" {
		t.Errorf("LastKnownLocation = %s; want = %s",
			c.Delivery.LastKnownLocation, "")
	}
}

func TestRoutingStatus(t *testing.T) {
	good := Itinerary{
		Legs: []Leg{
			{LoadLocation: SESTO, UnloadLocation: AUMEL},
		},
	}

	bad := Itinerary{
		Legs: []Leg{
			{LoadLocation: SESTO, UnloadLocation: CNHKG},
		},
	}

	acceptOnlyGood := RouteSpecification{
		Origin:      SESTO,
		Destination: AUMEL,
	}

	c := NewCargo("ABC", RouteSpecification{})

	c.SpecifyNewRoute(acceptOnlyGood)
	if c.Delivery.RoutingStatus != NotRouted {
		t.Errorf("RoutingStatus = %v; want = %v",
			c.Delivery.RoutingStatus, NotRouted)
	}

	c.AssignToRoute(bad)
	if c.Delivery.RoutingStatus != MisRouted {
		t.Errorf("RoutingStatus = %v; want = %v",
			c.Delivery.RoutingStatus, MisRouted)
	}

	c.AssignToRoute(good)
	if c.Delivery.RoutingStatus != Routed {
		t.Errorf("RoutingStatus = %v; want = %v",
			c.Delivery.RoutingStatus, Routed)
	}
}

func TestLastKnownLocation_WhenNoEvents(t *testing.T) {
	c := NewCargo("ABC", RouteSpecification{
		Origin:      SESTO,
		Destination: CNHKG,
	})

	if c.Delivery.LastKnownLocation != "" {
		t.Errorf("should be equal")
	}
}

func TestLastKnownLocation_WhenReceived(t *testing.T) {
	c := populateCargoReceivedInStockholm()

	if c.Delivery.LastKnownLocation != SESTO {
		t.Errorf("LastKnownLocation = %s; want = %s",
			c.Delivery.LastKnownLocation, SESTO)
	}
}

func TestLastKnownLocation_WhenClaimed(t *testing.T) {
	c := populateCargoClaimedInMelbourne()

	if c.Delivery.LastKnownLocation != AUMEL {
		t.Errorf("LastKnownLocation = %s; want = %s",
			c.Delivery.LastKnownLocation, AUMEL)
	}
}

var routingStatusTests = []struct {
	routingStatus RoutingStatus
	expected      string
}{
	{NotRouted, "Not routed"},
	{MisRouted, "Misrouted"},
	{Routed, "Routed"},
	{1000, ""},
}

func TestRoutingStatus_Stringer(t *testing.T) {
	for _, tt := range routingStatusTests {
		if tt.routingStatus.String() != tt.expected {
			t.Errorf("routingStatus.String() = %s; want = %s",
				tt.routingStatus.String(), tt.expected)
		}
	}
}

var transportStatusTests = []struct {
	transportStatus TransportStatus
	expected        string
}{
	{NotReceived, "Not received"},
	{InPort, "In port"},
	{OnboardCarrier, "Onboard carrier"},
	{Claimed, "Claimed"},
	{Unknown, "Unknown"},
	{1000, ""},
}

func TestTransportStatus_Stringer(t *testing.T) {
	for _, tt := range transportStatusTests {
		if tt.transportStatus.String() != tt.expected {
			t.Errorf("transportStatus.String() = %s; want = %s",
				tt.transportStatus.String(), tt.expected)
		}
	}
}

func populateCargoReceivedInStockholm() *Cargo {
	c := NewCargo("XYZ", RouteSpecification{
		Origin:      SESTO,
		Destination: AUMEL,
	})

	e := HandlingEvent{
		TrackingID: c.TrackingID,
		Activity: HandlingActivity{
			Type:     Receive,
			Location: SESTO,
		},
	}

	hh := HandlingHistory{
		HandlingEvents: []HandlingEvent{e},
	}

	c.DeriveDeliveryProgress(hh)

	return c
}

func populateCargoClaimedInMelbourne() *Cargo {
	c := NewCargo("XYZ", RouteSpecification{
		Origin:      SESTO,
		Destination: AUMEL,
	})

	e := HandlingEvent{
		TrackingID: c.TrackingID,
		Activity: HandlingActivity{
			Type:     Claim,
			Location: AUMEL,
		},
	}

	hh := HandlingHistory{
		HandlingEvents: []HandlingEvent{e},
	}

	c.DeriveDeliveryProgress(hh)

	return c
}
