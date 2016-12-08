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
	"github.com/mdreizin/smartling/service"
	"gopkg.in/go-playground/pool.v3"
)

func pushJob(req *pushRequest) pool.WorkFunc {
	return func(wu pool.WorkUnit) (interface{}, error) {
		if wu.IsCancelled() {
			return nil, nil
		}

		directives := req.Resource.Directives.WithPrefix()

		params := &service.FilePushParams{
			ProjectID:  req.Config.Project.ID,
			FileURI:    req.Config.FileURI(req.Path),
			FilePath:   req.Path,
			FileType:   req.Resource.Type,
			Authorize:  req.Resource.AuthorizeContent,
			Directives: directives,
			AuthToken:  req.AuthToken,
		}

		stats, err := req.FileService.Push(params)

		return &pushResponse{
			Stats:  stats,
			Params: params,
		}, err
	}
}
