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
package mongo

import (
	"cargo/model"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// 创建mongo相关的仓储
type cargoRepository struct {
	db      string
	session *mgo.Session
}

// 向mongo中存入cargo相关信息
func (repository *cargoRepository) Store(cargo model.Cargo) (bool, error) {
	// copy like to new
	sess := repository.session.Copy()
	defer sess.Close()

	collection := sess.DB(repository.db).C("cargo")
	_, _ = collection.Upsert(bson.M{"trackingid": cargo.TrackingID}, bson.M{"$set": cargo})

	return true, nil
}

// 从mongo中查询cargo相关信息
func (repository *cargoRepository) Find(id model.TrackingID) (*model.Cargo, error) {
	sess := repository.session.Copy()
	defer sess.Close()

	collection := sess.DB(repository.db).C("cargo")

	var result model.Cargo
	if err := collection.Find(bson.M{"trackingid": id}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, model.ErrUnknownCargo
		}
		return nil, err
	}

	return &result, nil
}

// 从mongo中查询所有的Cargo
func (repository *cargoRepository) FindAll() []*model.Cargo {
	sess := repository.session.Copy()
	defer sess.Close()

	collection := sess.DB(repository.db).C("cargo")

	var result []*model.Cargo
	if err := collection.Find(bson.M{}).All(&result); err != nil {
		return []*model.Cargo{}
	}

	return result
}

// NewCargoRepository returns a new instance of a MongoDB cargo repository.
func NewCargoRepository(db string, session *mgo.Session) (*cargoRepository, error) {
	r := &cargoRepository{
		db:      db,
		session: session,
	}

	index := mgo.Index{
		Key:        []string{"trackingid"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("cargo")

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	return r, nil
}

type locationRepository struct {
	db      string
	session *mgo.Session
}

func (r *locationRepository) Find(locode model.UNLocode) (*model.Location, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	var result model.Location
	if err := c.Find(bson.M{"unlocode": locode}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, model.ErrUnknownLocation
		}
		return nil, err
	}

	return &result, nil
}

func (r *locationRepository) FindAll() []*model.Location {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	var result []*model.Location
	if err := c.Find(bson.M{}).All(&result); err != nil {
		return []*model.Location{}
	}

	return result
}

func (r *locationRepository) store(l *model.Location) (bool, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	_, err := c.Upsert(bson.M{"unlocode": l.UNLocode}, bson.M{"$set": l})

	return true, err
}

// NewLocationRepository returns a new instance of a MongoDB location repository.
func NewLocationRepository(db string, session *mgo.Session) (model.LocationRepository, error) {
	r := &locationRepository{
		db:      db,
		session: session,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("location")

	index := mgo.Index{
		Key:        []string{"unlocode"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	initial := []*model.Location{
		model.Stockholm,
		model.Melbourne,
		model.Hongkong,
		model.Tokyo,
		model.Rotterdam,
		model.Hamburg,
	}

	for _, l := range initial {
		r.store(l)
	}

	return r, nil
}

type voyageRepository struct {
	db      string
	session *mgo.Session
}

func (r *voyageRepository) Find(voyageNumber model.VoyageNumber) (*model.Voyage, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	var result model.Voyage
	if err := c.Find(bson.M{"number": voyageNumber}).One(&result); err != nil {
		if err == mgo.ErrNotFound {
			return nil, model.ErrUnknownVoyage
		}
		return nil, err
	}

	return &result, nil
}

func (r *voyageRepository) store(v *model.Voyage) (bool, error) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	_, err := c.Upsert(bson.M{"number": v.VoyageNumber}, bson.M{"$set": v})

	return true, err
}

// NewVoyageRepository returns a new instance of a MongoDB voyage repository.
func NewVoyageRepository(db string, session *mgo.Session) (model.VoyageRepository, error) {
	r := &voyageRepository{
		db:      db,
		session: session,
	}

	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("voyage")

	index := mgo.Index{
		Key:        []string{"number"},
		Unique:     true,
		DropDups:   true,
		Background: true,
		Sparse:     true,
	}

	if err := c.EnsureIndex(index); err != nil {
		return nil, err
	}

	initial := []*model.Voyage{
		model.V100,
		model.V300,
		model.V400,
		model.V0100S,
		model.V0200T,
		model.V0300A,
		model.V0301S,
		model.V0400S,
	}

	for _, v := range initial {
		_, _ = r.store(v)
	}

	return r, nil
}

type handlingEventRepository struct {
	db      string
	session *mgo.Session
}

func (r *handlingEventRepository) Store(e model.HandlingEvent) {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("handling_event")

	_ = c.Insert(e)
}

func (r *handlingEventRepository) QueryHandlingHistory(id model.TrackingID) model.HandlingHistory {
	sess := r.session.Copy()
	defer sess.Close()

	c := sess.DB(r.db).C("handling_event")

	var result []model.HandlingEvent
	_ = c.Find(bson.M{"trackingid": id}).All(&result)

	return model.HandlingHistory{HandlingEvents: result}
}

// NewHandlingEventRepository returns a new instance of a MongoDB handling event repository.
func NewHandlingEventRepository(db string, session *mgo.Session) model.HandlingEventRepository {
	return &handlingEventRepository{
		db:      db,
		session: session,
	}
}
