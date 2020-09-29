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

// A set of sample voyages.
var (
	V100 = NewVoyage("V100", Schedule{
		[]CarrierMovement{
			{DepartureLocation: CNHKG, ArrivalLocation: JNTKO},
			{DepartureLocation: JNTKO, ArrivalLocation: USNYC},
		},
	})

	V300 = NewVoyage("V300", Schedule{
		[]CarrierMovement{
			{DepartureLocation: JNTKO, ArrivalLocation: NLRTM},
			{DepartureLocation: NLRTM, ArrivalLocation: DEHAM},
			{DepartureLocation: DEHAM, ArrivalLocation: AUMEL},
			{DepartureLocation: AUMEL, ArrivalLocation: JNTKO},
		},
	})

	V400 = NewVoyage("V400", Schedule{
		[]CarrierMovement{
			{DepartureLocation: DEHAM, ArrivalLocation: SESTO},
			{DepartureLocation: SESTO, ArrivalLocation: FIHEL},
			{DepartureLocation: FIHEL, ArrivalLocation: DEHAM},
		},
	})
)

// These voyages are hard-coded into the current pathfinder. Make sure
// they exist.
var (
	V0100S = NewVoyage("0100S", Schedule{[]CarrierMovement{}})
	V0200T = NewVoyage("0200T", Schedule{[]CarrierMovement{}})
	V0300A = NewVoyage("0300A", Schedule{[]CarrierMovement{}})
	V0301S = NewVoyage("0301S", Schedule{[]CarrierMovement{}})
	V0400S = NewVoyage("0400S", Schedule{[]CarrierMovement{}})
)
