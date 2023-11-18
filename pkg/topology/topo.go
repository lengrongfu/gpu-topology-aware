/*
Copyright 2022 The Kangaroo Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package topology

import "context"

type Matrix [][]uint64

type Type string

const (
	LinkType      = "Link"
	BandwidthType = "Bandwidth"
)

// Collector define collect single node all gpu device topology matrix info.
type Collector interface {
	// Collect function return matrix info, if collect error, return error.
	Collect(ctx context.Context) (Matrix, error)
}
