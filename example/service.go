package main

import (
	"net/http"

	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/clusters"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/config_items"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/infrastructure_accounts"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/node_pools"

	"gopkg.in/gin-gonic/gin.v1"
	"github.com/mikkeloscar/gin-swagger/api"
)

type ExampleService struct {
	Health bool
}

func (s *ExampleService) Healthy() bool {
	return s.Health
}
func (s *ExampleService) AddOrUpdateConfigItem(ctx *gin.Context, params *config_items.AddOrUpdateConfigItemParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) CreateCluster(ctx *gin.Context, params *clusters.CreateClusterParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) CreateInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.CreateInfrastructureAccountParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) CreateOrUpdateNodePool(ctx *gin.Context, params *node_pools.CreateOrUpdateNodePoolParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) DeleteCluster(ctx *gin.Context, params *clusters.DeleteClusterParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) DeleteConfigItem(ctx *gin.Context, params *config_items.DeleteConfigItemParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) DeleteNodePool(ctx *gin.Context, params *node_pools.DeleteNodePoolParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) GetCluster(ctx *gin.Context, params *clusters.GetClusterParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) GetHealth(ctx *gin.Context) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) GetInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.GetInfrastructureAccountParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) ListClusters(ctx *gin.Context, params *clusters.ListClustersParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) ListInfrastructureAccounts(ctx *gin.Context) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) ListNodePools(ctx *gin.Context, params *node_pools.ListNodePoolsParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) UpdateCluster(ctx *gin.Context, params *clusters.UpdateClusterParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
func (s *ExampleService) UpdateInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.UpdateInfrastructureAccountParams) *api.Response {
	return &api.Response{Code: http.StatusNotImplemented, Body: "Not Implemented"}
}
