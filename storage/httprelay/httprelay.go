// Copyright 2014 Google Inc. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httprelay

import (
	"bytes"
	"github.com/cloud66/cadvisor/storage"
	info "github.com/google/cadvisor/info/v1"
	"net/http"
	"net/url"
	"strconv"
	"time"
)

type httprelayStorage struct {
	machineName string
	httpPath    string
}

func (self *httprelayStorage) AddStats(ref info.ContainerReference, stats *info.ContainerStats) error {
	if stats == nil {
		return nil
	}

	path := self.httpPath

	data := url.Values{}

	// Container name
	name := ref.Name
	if len(ref.Aliases) > 0 {
		name = ref.Aliases[0]
	}
	data.Set("container_name", name)

	// Timestamp
	data.Add("time", stats.Timestamp.Format(time.RFC3339))

	// Cumulative Cpu Usage
	data.Add("cpu_usage_total", strconv.FormatUint(stats.Cpu.Usage.Total, 10))

	// Cumulative Cpu Usage in system mode
	data.Add("cpu_usage_system", strconv.FormatUint(stats.Cpu.Usage.System, 10))

	// Cumulative Cpu Usage in user mode
	data.Add("cpu_usage_user", strconv.FormatUint(stats.Cpu.Usage.User, 10))

	// Memory Usage
	data.Add("memory_usage", strconv.FormatUint(stats.Memory.Usage, 10))

	// Working set size
	data.Add("memory_working_set", strconv.FormatUint(stats.Memory.WorkingSet, 10))

	// container page fault
	data.Add("memory_container_pagefault", strconv.FormatUint(stats.Memory.ContainerData.Pgfault, 10))

	// container major page fault
	data.Add("memory_container_major_pagefault", strconv.FormatUint(stats.Memory.ContainerData.Pgmajfault, 10))

	// hierarchical page fault
	data.Add("memory_hierarchical_pagefault", strconv.FormatUint(stats.Memory.HierarchicalData.Pgfault, 10))

	// hierarchical major page fault
	data.Add("memory_hierarchical_major_pagefault", strconv.FormatUint(stats.Memory.HierarchicalData.Pgmajfault, 10))

	// Network stats.
	data.Add("network_rx_bytes", strconv.FormatUint(stats.Network.RxBytes, 10))
	data.Add("network_rx__errors", strconv.FormatUint(stats.Network.RxErrors, 10))
	data.Add("network_tx_bytes", strconv.FormatUint(stats.Network.TxBytes, 10))
	data.Add("network_tx_erros", strconv.FormatUint(stats.Network.TxErrors, 10))

	req, err := http.NewRequest("POST", path, bytes.NewBufferString(data.Encode()))

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

	return nil
}

func (self *httprelayStorage) RecentStats(containerName string, numStats int) ([]*info.ContainerStats, error) {
	return nil, nil
}

func (self *httprelayStorage) Close() error {
	return nil
}

// Create a new httprelay storage driver.
// machineName: A unique identifier to identify the host that current cAdvisor
// instance is running on.
// httpPath: httpPath that we must relay data to.
func New(machineName,
	httpPath string,
) (storage.StorageDriver, error) {
	ret := &httprelayStorage{
		machineName: machineName,
		httpPath:    httpPath,
	}
	return ret, nil
}
