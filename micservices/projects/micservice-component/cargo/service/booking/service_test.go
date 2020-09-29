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
package booking

import (
	"cargo/mock"
	"cargo/model"
	"testing"
	"time"
)

func TestBookNewCargo(t *testing.T) {
	var (
		origin      = model.SESTO
		destination = model.AUMEL
		deadline    = time.Date(2020, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	var cargos mockCargoRepository

	s := NewService(&cargos, nil, nil, nil)

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	c, err := cargos.Find(id)
	if err != nil {
		t.Fatal(err)
	}

	if c.TrackingID != id {
		t.Errorf("c.TrackingID = %s; want = %s", c.TrackingID, id)
	}
	if c.Origin != origin {
		t.Errorf("c.Origin = %s; want = %s", c.Origin, origin)
	}
	if c.RouteSpecification.Destination != destination {
		t.Errorf("c.RouteSpecification.Destination = %s; want = %s",
			c.RouteSpecification.Destination, destination)
	}
	if c.RouteSpecification.ArrivalDeadline != deadline {
		t.Errorf("c.RouteSpecification.ArrivalDeadline = %s; want = %s",
			c.RouteSpecification.ArrivalDeadline, deadline)
	}
}

type stubRoutingService struct{}

func (s *stubRoutingService) FetchRoutesForSpecification(rs model.RouteSpecification) []model.Itinerary {
	legs := []model.Leg{
		{LoadLocation: rs.Origin, UnloadLocation: rs.Destination},
	}

	return []model.Itinerary{
		{Legs: legs},
	}
}

func TestRequestPossibleRoutesForCargo(t *testing.T) {
	var (
		origin      = model.SESTO
		destination = model.AUMEL
		deadline    = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	var cargos mockCargoRepository

	var rs stubRoutingService

	s := NewService(&cargos, nil, nil, &rs)

	r := s.RequestPossibleRoutesForCargo("no_such_id")

	if len(r) != 0 {
		t.Errorf("len(r) = %d; want = %d", len(r), 0)
	}

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	i := s.RequestPossibleRoutesForCargo(id)

	if len(i) != 1 {
		t.Errorf("len(i) = %d; want = %d", len(i), 1)
	}
}

func TestAssignCargoToRoute(t *testing.T) {
	var cargos mockCargoRepository

	var rs stubRoutingService

	s := NewService(&cargos, nil, nil, &rs)

	var (
		origin      = model.SESTO
		destination = model.AUMEL
		deadline    = time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)
	)

	id, err := s.BookNewCargo(origin, destination, deadline)
	if err != nil {
		t.Fatal(err)
	}

	i := s.RequestPossibleRoutesForCargo(id)

	if len(i) != 1 {
		t.Errorf("len(i) = %d; want = %d", len(i), 1)
	}

	if _, err := s.AssignCargoToRoute(id, i[0]); err != nil {
		t.Fatal(err)
	}

	if _, err := s.AssignCargoToRoute("no_such_id", model.Itinerary{}); err != ErrInvalidArgument {
		t.Errorf("err = %s; want = %s", err, ErrInvalidArgument)
	}
}

func TestChangeCargoDestination(t *testing.T) {
	var cargos mockCargoRepository
	var locations mock.LocationRepository

	locations.FindFn = func(loc model.UNLocode) (*model.Location, error) {
		if loc != model.AUMEL {
			return nil, model.ErrUnknownLocation
		}
		return model.Melbourne, nil
	}

	var rs stubRoutingService

	s := NewService(&cargos, &locations, nil, &rs)

	c := model.NewCargo("ABC", model.RouteSpecification{
		Origin:          model.SESTO,
		Destination:     model.CNHKG,
		ArrivalDeadline: time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC),
	})

	if _, err := s.ChangeDestination("no_such_id", model.SESTO); err != model.ErrUnknownCargo {
		t.Errorf("err = %s; want = %s", err, model.ErrUnknownCargo)
	}

	if _, err := cargos.Store(c); err != nil {
		t.Fatal(err)
	}

	if _, err := s.ChangeDestination(c.TrackingID, "no_such_unlocode"); err != model.ErrUnknownLocation {
		t.Errorf("err = %s; want = %s", err, model.ErrUnknownLocation)
	}

	if c.RouteSpecification.Destination != model.CNHKG {
		t.Errorf("c.RouteSpecification.Destination = %s; want = %s",
			c.RouteSpecification.Destination, model.CNHKG)
	}

	if _, err := s.ChangeDestination(c.TrackingID, model.AUMEL); err != nil {
		t.Fatal(err)
	}

	uc, err := cargos.Find(c.TrackingID)
	if err != nil {
		t.Fatal(err)
	}

	if uc.RouteSpecification.Destination != model.AUMEL {
		t.Errorf("uc.RouteSpecification.Destination = %s; want = %s",
			uc.RouteSpecification.Destination, model.AUMEL)
	}
}

func TestLoadCargo(t *testing.T) {
	deadline := time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC)

	var cargos mock.CargoRepository
	cargos.FindFn = func(id model.TrackingID) (*model.Cargo, error) {
		return &model.Cargo{
			TrackingID: "test_id",
			Origin:     model.SESTO,
			RouteSpecification: model.RouteSpecification{
				Origin:          model.SESTO,
				Destination:     model.AUMEL,
				ArrivalDeadline: deadline,
			},
		}, nil
	}

	s := NewService(&cargos, nil, nil, nil)

	c, err := s.LoadCargo("test_id")
	if err != nil {
		t.Fatal(err)
	}

	if c.TrackingID != "test_id" {
		t.Errorf("c.TrackingID = %s; want = %s", c.TrackingID, "test_id")
	}
	if c.Origin != "SESTO" {
		t.Errorf("c.Origin = %s; want = %s", c.Origin, "SESTO")
	}
	if c.Destination != "AUMEL" {
		t.Errorf("c.Destination = %s; want = %s", c.Origin, "AUMEL")
	}
	if c.ArrivalDeadline != deadline {
		t.Errorf("c.ArrivalDeadline = %s; want = %s", c.ArrivalDeadline, deadline)
	}
	if c.Misrouted {
		t.Errorf("cargo should not be misrouted")
	}
	if c.Routed {
		t.Errorf("cargo should not have been routed")
	}
	if len(c.Legs) != 0 {
		t.Errorf("len(c.Legs) = %d; want = %d", len(c.Legs), 0)
	}
}

type mockCargoRepository struct {
	cargo *model.Cargo
}

func (r *mockCargoRepository) Store(c *model.Cargo) (bool, error) {
	r.cargo = c
	return true, nil
}

func (r *mockCargoRepository) Find(id model.TrackingID) (*model.Cargo, error) {
	if r.cargo != nil {
		return r.cargo, nil
	}
	return nil, model.ErrUnknownCargo
}

func (r *mockCargoRepository) FindAll() []*model.Cargo {
	return []*model.Cargo{r.cargo}
}
