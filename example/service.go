package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExampleService struct {
	Health bool
}

func (s *ExampleService) Healthy() bool {
	return s.Health
}
func (s *ExampleService) AddOrUpdateConfigItem(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) DeleteConfigItem(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) CreateCluster(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) CreateInfrastructureAccount(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) CreateOrUpdateNodePool(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) DeleteCluster(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) DeleteConfigValue(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) DeleteNodePool(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) GetCluster(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) GetHealth(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) GetInfrastructureAccount(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) ListClusters(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) ListInfrastructureAccounts(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) ListNodePools(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) UpdateCluster(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
func (s *ExampleService) UpdateInfrastructureAccount(ctx *gin.Context) {
	ctx.String(http.StatusNotImplemented, "Not implemented")
}
