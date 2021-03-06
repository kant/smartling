// Copyright 2016, Fitbit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.
package main

import (
	"github.com/Fitbit/smartling/di"
	"github.com/Fitbit/smartling/logger"
	"github.com/Fitbit/smartling/model"
	"github.com/fatih/color"
	"gopkg.in/go-playground/pool.v3"
	"gopkg.in/urfave/cli.v1"
	"runtime"
	"time"
)

var pushCommand = cli.Command{
	Name:  "push",
	Usage: "Uploads translations",
	Before: func(c *cli.Context) error {
		return invokeActions([]action{
			injectDiContainerAction,
			injectProjectConfigAction,
			validateProjectConfigAction,
			injectAuthTokenAction,
		}, c)
	},
	Action: func(c *cli.Context) error {
		defer elapsedTime(time.Now())

		p := pool.NewLimited(uint(runtime.NumCPU()))

		defer p.Close()

		batch := p.Batch()

		container := c.App.Metadata[containerMetadataKey].(*di.Container)
		authToken := c.App.Metadata[authTokenMetadataKey].(*model.AuthToken)
		projectConfig := c.App.Metadata[projectConfigMetadataKey].(*model.ProjectConfig)

		go func() {
			for _, resource := range projectConfig.Resources {
				for _, path := range resource.Files() {
					batch.Queue(pushWorker(&pushWorkerParams{
						Path:        path,
						Resource:    &resource,
						Config:      projectConfig,
						AuthToken:   authToken.AccessToken,
						FileService: container.FileService,
					}))
				}
			}

			batch.QueueComplete()
		}()

		i := 0

		for results := range batch.Results() {
			result := results.Value().(*pushWorkerResult)

			if err := results.Error(); err != nil {
				logger.Errorf("%s has error %s", color.MagentaString(result.Params.Path), color.RedString(err.Error()))
			} else {
				logger.Infof("%s {Override=%t Strings=%d Words=%d}", color.MagentaString(result.Params.Path), result.Stats.OverWritten, result.Stats.StringCount, result.Stats.WordCount)
				i++
			}
		}

		logger.Infof("%d files", i)

		return nil
	},
	After: func(c *cli.Context) error {
		return invokeActions([]action{
			persistAuthTokenAction,
		}, c)
	},
}
