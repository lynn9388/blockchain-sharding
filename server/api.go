/*
 * Copyright © 2018 Lynn <lynn9388@gmail.com>
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package server

import (
	"net/http"

	"github.com/emicklei/go-restful"
	"github.com/emicklei/go-restful-openapi"
	"github.com/go-openapi/spec"
	"github.com/lynn9388/blockchain-sharding/common"
)

type serverAPI int

func startAPIServer() {
	restful.DefaultContainer.Add(new(serverAPI).webService())

	config := restfulspec.Config{
		WebServices: restful.RegisteredWebServices(),
		APIPath:     "/api.json",
		PostBuildSwaggerObjectHandler: enrichSwaggerObject}
	restful.DefaultContainer.Add(restfulspec.NewOpenAPIService(config))

	// Optionally, you can install the Swagger Service which provides a nice Web UI on your REST API
	// You need to download the Swagger HTML5 assets and change the FilePath location in the config below.
	http.Handle("/api/", http.StripPrefix("/api/", http.FileServer(http.Dir("/home/lynn/Documents/Git/swagger-ui/dist"))))

	addr := common.GetServerInfo().APIAddr
	common.Logger.Infof("API server listening at: %v, see API doc on http://%v/api/?url=http://%v/api.json", addr, addr, addr)
	go http.ListenAndServe(addr, nil)
}

// WebService creates a new service that can handle REST requests for Server resources.
func (s *serverAPI) webService() *restful.WebService {
	ws := new(restful.WebService)
	ws.
		Path("/server").
		Consumes(restful.MIME_JSON, restful.MIME_XML).
		Produces(restful.MIME_JSON, restful.MIME_XML)

	tags := []string{"server"}

	ws.Route(ws.GET("/").To(s.findServerConfig).
		Doc("get server config").
		Metadata(restfulspec.KeyOpenAPITags, tags).
		Writes(common.Config{}).
		Returns(200, "OK", common.Config{}))

	return ws
}

func (s *serverAPI) findServerConfig(request *restful.Request, response *restful.Response) {
	response.WriteEntity(common.GetConfig())
}

func enrichSwaggerObject(swo *spec.Swagger) {
	swo.Info = &spec.Info{
		InfoProps: spec.InfoProps{
			Title:       "API Service",
			Description: "Resource for managing Server",
			Contact: &spec.ContactInfo{
				Name:  "Lynn",
				Email: "lynn9388@gmail.com",
				URL:   "https://lynn9388.com/",
			},
			License: &spec.License{
				Name: "Apache-2.0",
				URL:  "http://www.apache.org/licenses/",
			},
			Version: "0.1",
		},
	}
	swo.Tags = []spec.Tag{{TagProps: spec.TagProps{
		Name:        "server",
		Description: "Managing server"}}}
}
