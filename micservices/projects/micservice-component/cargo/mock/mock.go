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
package mock

import "cargo/model"

// CargoRepository is a mock cargo repository.
type CargoRepository struct {
	StoreFn      func(c *model.Cargo) error
	StoreInvoked bool

	FindFn      func(id model.TrackingID) (*model.Cargo, error)
	FindInvoked bool

	FindAllFn      func() []*model.Cargo
	FindAllInvoked bool
}

// Store calls the StoreFn.
func (r *CargoRepository) Store(c *model.Cargo) (bool, error) {
	r.StoreInvoked = true
	return true, r.StoreFn(c)
}

// Find calls the FindFn.
func (r *CargoRepository) Find(id model.TrackingID) (*model.Cargo, error) {
	r.FindInvoked = true
	return r.FindFn(id)
}

// FindAll calls the FindAllFn.
func (r *CargoRepository) FindAll() []*model.Cargo {
	r.FindAllInvoked = true
	return r.FindAllFn()
}

// LocationRepository is a mock location repository.
type LocationRepository struct {
	FindFn      func(model.UNLocode) (*model.Location, error)
	FindInvoked bool

	FindAllFn      func() []*model.Location
	FindAllInvoked bool
}

// Find calls the FindFn.
func (r *LocationRepository) Find(locode model.UNLocode) (*model.Location, error) {
	r.FindInvoked = true
	return r.FindFn(locode)
}

// FindAll calls the FindAllFn.
func (r *LocationRepository) FindAll() []*model.Location {
	r.FindAllInvoked = true
	return r.FindAllFn()
}

// VoyageRepository is a mock voyage repository.
type VoyageRepository struct {
	FindFn      func(model.VoyageNumber) (*model.Voyage, error)
	FindInvoked bool
}

// Find calls the FindFn.
func (r *VoyageRepository) Find(number model.VoyageNumber) (*model.Voyage, error) {
	r.FindInvoked = true
	return r.FindFn(number)
}

// HandlingEventRepository is a mock handling events repository.
type HandlingEventRepository struct {
	StoreFn      func(model.HandlingEvent)
	StoreInvoked bool

	QueryHandlingHistoryFn      func(model.TrackingID) model.HandlingHistory
	QueryHandlingHistoryInvoked bool
}

// Store calls the StoreFn.
func (r *HandlingEventRepository) Store(e model.HandlingEvent) {
	r.StoreInvoked = true
	r.StoreFn(e)
}

// QueryHandlingHistory calls the QueryHandlingHistoryFn.
func (r *HandlingEventRepository) QueryHandlingHistory(id model.TrackingID) model.HandlingHistory {
	r.QueryHandlingHistoryInvoked = true
	return r.QueryHandlingHistoryFn(id)
}

// RoutingService provides a mock routing service.
type RoutingService struct {
	FetchRoutesFn      func(model.RouteSpecification) []model.Itinerary
	FetchRoutesInvoked bool
}

// FetchRoutesForSpecification calls the FetchRoutesFn.
func (s *RoutingService) FetchRoutesForSpecification(rs model.RouteSpecification) []model.Itinerary {
	s.FetchRoutesInvoked = true
	return s.FetchRoutesFn(rs)
}
