// Copyright 2019-present Open Networking Foundation.
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

package kpimon

import (
	"context"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"text/tabwriter"
	"time"

	prototypes "github.com/gogo/protobuf/types"
	"github.com/prometheus/common/log"

	kpimonapi "github.com/onosproject/onos-api/go/onos/kpimon"
	"github.com/onosproject/onos-lib-go/pkg/cli"
	"github.com/spf13/cobra"
)

func getListMetricsCommand() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "metrics",
		Short: "Get metrics",
		RunE:  runListMetricsCommand,
	}
	cmd.Flags().Bool("no-headers", false, "disables output headers")
	return cmd
}

func runListMetricsCommand(cmd *cobra.Command, args []string) error {
	var types []string
	results := make(map[string]map[uint64]map[string]string)

	noHeaders, _ := cmd.Flags().GetBool("no-headers")
	conn, err := cli.GetConnection(cmd)
	if err != nil {
		return err
	}
	defer conn.Close()
	outputWriter := cli.GetOutput()
	writer := new(tabwriter.Writer)
	writer.Init(outputWriter, 0, 0, 3, ' ', tabwriter.FilterHTML)

	request := kpimonapi.GetRequest{}
	client := kpimonapi.NewKpimonClient(conn)

	respGetMeasurement, err := client.ListMeasurements(context.Background(), &request)
	if err != nil {
		return err
	}

	attr := make(map[string]string)
	for key, measItems := range respGetMeasurement.GetMeasurements() {
		for _, measItem := range measItems.MeasurementItems {
			for _, measRecord := range measItem.MeasurementRecords {
				timeStamp := measRecord.Timestamp
				measName := measRecord.MeasurementName
				measValue := measRecord.MeasurementValue

				if _, ok := attr[measName]; !ok {
					attr[measName] = measName
				}

				if _, ok1 := results[key]; !ok1 {
					results[key] = make(map[uint64]map[string]string)
				}
				if _, ok2 := results[key][timeStamp]; !ok2 {
					results[key][timeStamp] = make(map[string]string)
				}

				var value interface{}

				switch {
				case prototypes.Is(measValue, &kpimonapi.IntegerValue{}):
					v := kpimonapi.IntegerValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.RealValue{}):
					v := kpimonapi.RealValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				case prototypes.Is(measValue, &kpimonapi.NoValue{}):
					v := kpimonapi.NoValue{}
					err := prototypes.UnmarshalAny(measValue, &v)
					if err != nil {
						log.Warn(err)
					}
					value = v.GetValue()

				}

				results[key][timeStamp][measName] = fmt.Sprintf("%v", value)
			}
		}
	}

	for key := range attr {
		types = append(types, key)
	}
	sort.Strings(types)

	header := "Node ID\tCell ID\tTime"

	for _, key := range types {
		tmpHeader := header
		header = fmt.Sprintf("%s\t%s", tmpHeader, key)
	}

	if !noHeaders {
		_, _ = fmt.Fprintln(writer, header)
	}

	keys := make([]string, 0, len(results))
	for k := range results {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, keyID := range keys {
		metrics := results[keyID]
		// sort 2nd map with timestamp
		timeKeySlice := make([]uint64, 0, len(metrics))
		for timeStampKey := range metrics {
			timeKeySlice = append(timeKeySlice, timeStampKey)
		}

		sort.Slice(timeKeySlice, func(i, j int) bool { return timeKeySlice[i] < timeKeySlice[j] })

		for _, timeStamp := range timeKeySlice {
			timeObj := time.Unix(0, int64(timeStamp))
			tsFormat := fmt.Sprintf("%02d:%02d:%02d.%d", timeObj.Hour(), timeObj.Minute(), timeObj.Second(), timeObj.Nanosecond()/1000000)

			ids := strings.Split(keyID, ":")
			nodeID, cellID := ids[0], ids[1]
			// parse string to int in order to print as hex
			cellNum, err := strconv.Atoi(cellID)
			if err != nil {
				return err
			}
			resultLine := fmt.Sprintf("%s\t%s\t%s", nodeID, fmt.Sprintf("%x", cellNum), tsFormat)
			for _, typeValue := range types {
				tmpResultLine := resultLine
				var tmpValue string
				if _, ok := metrics[timeStamp][typeValue]; !ok {
					tmpValue = "N/A"
				} else {
					tmpValue = metrics[timeStamp][typeValue]
				}
				resultLine = fmt.Sprintf("%s\t%s", tmpResultLine, tmpValue)
			}
			_, _ = fmt.Fprintln(writer, resultLine)
		}
	}

	_ = writer.Flush()

	return nil
}
