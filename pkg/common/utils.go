/*
 * Copyright (c) 2024, NVIDIA CORPORATION.  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package common

import (
	"fmt"
	"os"
	"reflect"
	"slices"
	"sync"
	"time"
)

func WaitWithTimeout(wg *sync.WaitGroup, timeout time.Duration) error {
	c := make(chan struct{})
	go func() {
		defer close(c)
		wg.Wait()
	}()
	select {
	case <-c:
		return nil
	case <-time.After(timeout):
		return fmt.Errorf("timeout waiting for WaitGroup")
	}
}

func GetHostname(config *Config) (string, error) {
	hostname := ""
	var err error
	if !config.NoHostname {
		if nodeName := os.Getenv("NODE_NAME"); nodeName != "" {
			hostname = nodeName
		} else {
			hostname, err = os.Hostname()
			if err != nil {
				return "", err
			}
		}
	}
	return hostname, nil
}

func IsMetricsTypeEnabled(counters []Counter, metricsType string) bool {
	return slices.ContainsFunc(counters, func(c Counter) bool {
		return c.FieldName == metricsType
	})
}

func IsNil(i interface{}) bool {
	iv := reflect.ValueOf(i)
	if !iv.IsValid() {
		return true
	}
	switch iv.Kind() {
	case reflect.Ptr, reflect.Slice, reflect.Map, reflect.Func, reflect.Interface:
		return iv.IsNil()
	default:
		return false
	}
}
