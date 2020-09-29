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
	"errors"
	"time"
)

// VoyageNumber uniquely identifies a particular Voyage.
type VoyageNumber string

// Voyage is a uniquely identifiable series of carrier movements.
type Voyage struct {
	VoyageNumber VoyageNumber
	Schedule     Schedule
}

// NewVoyage creates a voyage with a voyage number and a provided schedule.
func NewVoyage(n VoyageNumber, s Schedule) *Voyage {
	return &Voyage{VoyageNumber: n, Schedule: s}
}

// Schedule describes a voyage schedule.
type Schedule struct {
	CarrierMovements []CarrierMovement
}

// CarrierMovement is a vessel voyage from one location to another.
type CarrierMovement struct {
	DepartureLocation UNLocode
	ArrivalLocation   UNLocode
	DepartureTime     time.Time
	ArrivalTime       time.Time
}

// ErrUnknownVoyage is used when a voyage could not be found.
var ErrUnknownVoyage = errors.New("unknown voyage")

// VoyageRepository provides access a voyage store.
type VoyageRepository interface {
	Find(VoyageNumber) (*Voyage, error)
}
