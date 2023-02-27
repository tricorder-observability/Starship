// Copyright (C) 2023  Tricorder Observability
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU Affero General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU Affero General Public License for more details.
//
// You should have received a copy of the GNU Affero General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package cond

import "sync"

type Cond struct {
	cond      *sync.Cond
	sharedRsc bool
}

func NewCond() *Cond {
	return &Cond{
		cond: sync.NewCond(new(sync.Mutex)),
	}
}

func (c *Cond) Wait() {
	c.cond.L.Lock()
	for !c.sharedRsc {
		c.cond.Wait()
	}
	c.cond.L.Unlock()
}

func (c *Cond) Signal() {
	c.sharedRsc = true
	c.cond.Signal()
}

func (c *Cond) Broadcast() {
	c.sharedRsc = true
	c.cond.Broadcast()
}
