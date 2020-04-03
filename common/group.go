// Copyright 2019 The berith Authors
// This file is part of berith.
//
// berith is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// berith is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with berith. If not, see <http://www.gnu.org/licenses/>.

package common

import (
	"errors"
	"math"
)

// SequenceGroup is grouping given orders depends on sequence type such as arithmetic sequence, geometric sequence
// and first term is always 1.
// for example, assume that we use a arithmetic group and common difference is 3.
// then, sequence is grouped like {1, 4, 7, ...}.
// element == 1 will have a group order 1.
// 2 <= element <= 4 will have a group order 2.
// 5 <= element <= 7 will have a group order 3.
type SequenceGroup interface {
	GetGroupOrder(termOrder int) (int, error)       // returns included group number given rank
	GetGroupRange(groupOrder int) (int, int, error) // returns [start term order, last term order] given group order
}

// ArithmeticGroup is included common difference between groups
// e.g) 1, 6, 11, 16 ..., where first element is 1 and common difference is 5
type ArithmeticGroup struct {
	CommonDiff int // common difference
}

// GetGroupOrder returns included group number given elements order.
func (a *ArithmeticGroup) GetGroupOrder(termOrder int) (int, error) {
	if termOrder < 1 {
		return 0, errors.New("term order must be larger than 0")
	}

	if termOrder == 1 {
		return 1, nil
	}

	d := math.Ceil(float64(termOrder-1) / float64(a.CommonDiff))
	return int(d) + 1, nil
}

// GetGroupRange returns [start term order, last term order] given group order.
func (a *ArithmeticGroup) GetGroupRange(groupOrder int) (int, int, error) {
	if groupOrder < 1 {
		return 0, 0, errors.New("group order must be larger than 0")
	}

	if groupOrder == 1 {
		return 1, 1, nil
	}

	return 2 + (groupOrder-2) * a.CommonDiff, 1 + (groupOrder-1) * a.CommonDiff, nil
}
