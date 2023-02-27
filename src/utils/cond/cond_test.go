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

import (
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"golang.org/x/sync/errgroup"
)

func TestCond(t *testing.T) {
	assert := assert.New(t)

	cond := NewCond()
	var group errgroup.Group

	// broadcast work fine
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			c := make(chan struct{})
			go func() {
				cond.Wait()
				c <- struct{}{}
			}()

			select {
			case <-c:
				return nil
			case <-time.After(5 * time.Second):
				return fmt.Errorf("timeout")
			}
		})
	}

	cond.Broadcast()
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}

	assert.Nil(group.Wait())

	// Signal work fine
	for i := 0; i < 5; i++ {
		group.Go(func() error {
			c := make(chan struct{})
			go func() {
				cond.Wait()
				c <- struct{}{}
			}()

			select {
			case <-c:
				return fmt.Errorf("wait success")
			case <-time.After(5 * time.Second):
				return nil
			}
		})
	}

	cond.Signal()
	if err := group.Wait(); err != nil {
		fmt.Println(err)
	}

	assert.NotNil(group.Wait(), "group must return error")
}
