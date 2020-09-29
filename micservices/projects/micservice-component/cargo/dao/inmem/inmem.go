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
package inmem

import (
	"cargo/model"
	"sync"
)

// cargo 仓储定义
type cargoRepository struct {
	mtx    sync.RWMutex
	cargos map[model.TrackingID]*model.Cargo
}

// 存储结构体内容
func (repository *cargoRepository) Store(cargo *model.Cargo) (bool, error) {
	repository.mtx.Lock()
	defer repository.mtx.Unlock()

	repository.cargos[cargo.TrackingID] = cargo
	return true, nil
}

// 查询cargoRepository
func (repository *cargoRepository) find(id model.TrackingID) (*model.Cargo, error) {
	repository.mtx.RLock()
	defer repository.mtx.RUnlock()
	if val, ok := repository.cargos[id]; ok {
		return val, nil
	}

	return nil, model.ErrUnknownCargo
}

// 查询所有的cargo信息
func (repository *cargoRepository) findAll() []*model.Cargo {
	repository.mtx.RLock()
	defer repository.mtx.RUnlock()

	categories := make([]*model.Cargo, len(repository.cargos))
	for _, val := range repository.cargos {
		categories = append(categories, val)
	}
	return categories
}

// 创建CargoRepository NewCargoRepository returns a new instance of a in-memory cargo repository.
func NewCargoRepository() *cargoRepository {
	return &cargoRepository{
		cargos: make(map[model.TrackingID]*model.Cargo),
	}
}

// 创建location仓储
type locationRepository struct {
	locations map[model.UNLocode]*model.Location
}

// 根据区域编码查询LocationRepository
func (repository *locationRepository) Find(locode model.UNLocode) (*model.Location, error) {
	if val, ok := repository.locations[locode]; ok {
		return val, nil
	}
	return nil, model.ErrUnknownLocation
}

// 根据查询所有的区域相关信息
func (repository *locationRepository) FindAll() []*model.Location {
	locations := make([]*model.Location, len(repository.locations))
	for _, val := range repository.locations {
		locations = append(locations, val)
	}
	return locations
}

// NewLocationRepository returns a new instance of a in-memory location repository.
func NewLocationRepository() model.LocationRepository {
	r := &locationRepository{
		locations: make(map[model.UNLocode]*model.Location),
	}

	r.locations[model.SESTO] = model.Stockholm
	r.locations[model.AUMEL] = model.Melbourne
	r.locations[model.CNHKG] = model.Hongkong
	r.locations[model.JNTKO] = model.Tokyo
	r.locations[model.NLRTM] = model.Rotterdam
	r.locations[model.DEHAM] = model.Hamburg

	return r
}

type voyageRepository struct {
	voyages map[model.VoyageNumber]*model.Voyage
}

func (r *voyageRepository) Find(voyageNumber model.VoyageNumber) (*model.Voyage, error) {
	if v, ok := r.voyages[voyageNumber]; ok {
		return v, nil
	}

	return nil, model.ErrUnknownVoyage
}

// NewVoyageRepository returns a new instance of a in-memory voyage repository.
func NewVoyageRepository() model.VoyageRepository {
	r := &voyageRepository{
		voyages: make(map[model.VoyageNumber]*model.Voyage),
	}

	r.voyages[model.V100.VoyageNumber] = model.V100
	r.voyages[model.V300.VoyageNumber] = model.V300
	r.voyages[model.V400.VoyageNumber] = model.V400

	r.voyages[model.V0100S.VoyageNumber] = model.V0100S
	r.voyages[model.V0200T.VoyageNumber] = model.V0200T
	r.voyages[model.V0300A.VoyageNumber] = model.V0300A
	r.voyages[model.V0301S.VoyageNumber] = model.V0301S
	r.voyages[model.V0400S.VoyageNumber] = model.V0400S

	return r
}

type handlingEventRepository struct {
	mtx    sync.RWMutex
	events map[model.TrackingID][]model.HandlingEvent
}

func (r *handlingEventRepository) Store(e model.HandlingEvent) {
	r.mtx.Lock()
	defer r.mtx.Unlock()
	// Make array if it's the first event with this tracking ID.
	if _, ok := r.events[e.TrackingID]; !ok {
		r.events[e.TrackingID] = make([]model.HandlingEvent, 0)
	}
	r.events[e.TrackingID] = append(r.events[e.TrackingID], e)
}

func (r *handlingEventRepository) QueryHandlingHistory(id model.TrackingID) model.HandlingHistory {
	r.mtx.RLock()
	defer r.mtx.RUnlock()
	return model.HandlingHistory{HandlingEvents: r.events[id]}
}

// NewHandlingEventRepository returns a new instance of a in-memory handling event repository.
func NewHandlingEventRepository() model.HandlingEventRepository {
	return &handlingEventRepository{
		events: make(map[model.TrackingID][]model.HandlingEvent),
	}
}
