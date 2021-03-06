// Copyright 2016, Fitbit, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License"); you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software distributed under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and limitations under the License.
package service

import (
	"github.com/Fitbit/smartling/model"
	"github.com/Fitbit/smartling/rest"
	"github.com/Fitbit/smartling/test"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDefaultAuthService_Authenticate(t *testing.T) {
	ts := test.MockServer()

	defer ts.Close()

	authToken := model.AuthToken{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	authService := DefaultAuthService{
		Client: rest.Client(false),
	}

	resp, err := authService.Authenticate(&model.UserToken{
		ID:     "userId",
		Secret: "userSecret",
	})

	a := assert.New(t)

	a.NoError(err)
	a.EqualValues(&authToken, resp)
}

func TestDefaultAuthService_Refresh(t *testing.T) {
	ts := test.MockServer()

	defer ts.Close()

	authToken := model.AuthToken{
		AccessToken:  "accessToken",
		RefreshToken: "refreshToken",
	}

	authService := DefaultAuthService{
		Client: rest.Client(false),
	}

	resp, err := authService.Refresh("refreshToken")

	a := assert.New(t)

	a.NoError(err)
	a.EqualValues(&authToken, resp)
}
