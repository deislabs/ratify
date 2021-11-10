/*
Copyright The Ratify Authors.
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

package httpserver

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/docker/distribution/registry/api/errcode"
	"github.com/sirupsen/logrus"
)

type Error struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// ContextHandler defines a http handler with a context input
type ContextHandler func(ctx context.Context, w http.ResponseWriter, r *http.Request) error

// contextHandler is http handler wrappered by a context
type contextHandler struct {
	context context.Context
	handler ContextHandler
}

// ServeHTTP serves an HTTP request and implements the http.Handler interface.
func (ch *contextHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	logrus.Infof("received request %v %v ", r.Method, r.URL)
	if err := ch.handler(ch.context, w, r); err != nil {
		logrus.Errorf("request %v %v failed with error %v", r.Method, r.URL, err)
		if serveErr := errcode.ServeJSON(w, err); serveErr != nil {
			// TODO log the error
			logrus.Errorf("request %v %v failed to send with error  %v", r.Method, r.URL, serveErr)
		}
	}
}

func serveErrorJSON(w http.ResponseWriter, err error) error {
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	errObj := Error{
		Code:    "UNKNOWN",
		Message: err.Error(),
	}
	if err := json.NewEncoder(w).Encode(errObj); err != nil {
		return err
	}
	return nil
}
