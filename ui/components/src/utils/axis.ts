// Copyright 2023 The Perses Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

import merge from 'lodash/merge';
import type { XAXisComponentOption, YAXisComponentOption } from 'echarts';
import { formatValue, FormatOptions } from '@perses-dev/core';

/*
 * Populate yAxis or xAxis properties, returns an Array since multiple axes will be supported in the future
 */
export function getFormattedAxis(axis?: YAXisComponentOption | XAXisComponentOption, unit?: FormatOptions): unknown[] {
  // TODO: support alternate yAxis that shows on right side
  const AXIS_DEFAULT = {
    type: 'value',
    boundaryGap: [0, '10%'],
    axisLabel: {
      formatter: (value: number): string => {
        return formatValue(value, unit);
      },
    },
  };
  return [merge(AXIS_DEFAULT, axis)];
}
