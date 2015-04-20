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
	"encoding/json"
	"github.com/cloud66/cadvisor/storage"
	"github.com/google/cadvisor/container/docker"
	info "github.com/google/cadvisor/info/v1"
	"net/http"
	"strings"
)

type httprelayStorage struct {
	machineName string
	httpPath    string
	containerId string
}

func (self *httprelayStorage) AddStats(ref info.ContainerReference, stats *info.ContainerStats) error {
	if stats == nil {
		return nil
	}
	path := self.httpPath

	// Container name
	name := ref.Name
	container_id := docker.ContainerNameToDockerId(name)

	if len(container_id) != 64 {
		return nil
	}

	if (len(self.containerId)) > 0 && (!strings.HasPrefix(container_id, self.containerId)) {
		return nil
	}

	stats_json, _ := json.Marshal(stats)

	payload := make(map[string]string)
	payload["container_id"] = container_id
	payload["data"] = string(stats_json)

	payload_json, _ := json.Marshal(payload)

	req, err := http.NewRequest("POST", path, bytes.NewBufferString(string(payload_json)))

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
	httpPath,
	containerId string,
) (storage.StorageDriver, error) {
	ret := &httprelayStorage{
		machineName: machineName,
		httpPath:    httpPath,
		containerId: containerId,
	}
	return ret, nil
}
