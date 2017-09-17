package restapi

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"golang.org/x/oauth2"

	"github.com/gin-contrib/pprof"
	"github.com/gin-gonic/gin"
	"github.com/mikkeloscar/gin-swagger/api"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/clusters"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/config_items"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/infrastructure_accounts"
	"github.com/mikkeloscar/gin-swagger/example/restapi/operations/node_pools"
	"github.com/mikkeloscar/gin-swagger/middleware"
	log "github.com/sirupsen/logrus"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

// Routes defines all the routes of the API service.
type Routes struct {
	*gin.Engine
	AddOrUpdateConfigItem struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	CreateCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	CreateInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	CreateOrUpdateNodePool struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	DeleteCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	DeleteConfigItem struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	DeleteNodePool struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	GetCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	GetInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	ListClusters struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	ListInfrastructureAccounts struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	ListNodePools struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	UpdateCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
	UpdateInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
		Post *gin.RouterGroup
	}
}

// configureWellKnown enables and configures /.well-known endpoints.
func (r *Routes) configureWellKnown(healthFunc func() bool) {
	wellKnown := r.Group("/.well-known")
	{
		wellKnown.GET("/schema-discovery", func(ctx *gin.Context) {
			discovery := struct {
				SchemaURL  string `json:"schema_url"`
				SchemaType string `json:"schema_type"`
				UIURL      string `json:"ui_url"`
			}{
				SchemaURL:  "/swagger.json",
				SchemaType: "swagger-2.0",
			}
			ctx.JSON(http.StatusOK, &discovery)
		})
		wellKnown.GET("/health", healthHandler(healthFunc))
	}

	r.GET("/swagger.json", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, string(SwaggerJSON))
	})
}

// healthHandler is the health HTTP handler used for the /.well-known/health
// route if enabled.
func healthHandler(healthFunc func() bool) gin.HandlerFunc {
	healthy := healthFunc
	return func(ctx *gin.Context) {
		health := struct {
			Health bool `json:"health"`
		}{
			Health: healthy(),
		}

		if !health.Health {
			ctx.JSON(http.StatusServiceUnavailable, &health)
		} else {
			ctx.JSON(http.StatusOK, &health)
		}
	}
}

// Service is the interface that must be implemented in order to provide
// business logic for the API service.
type Service interface {
	Healthy() bool
	AddOrUpdateConfigItem(ctx *gin.Context, params *config_items.AddOrUpdateConfigItemParams) *api.Response
	CreateCluster(ctx *gin.Context, params *clusters.CreateClusterParams) *api.Response
	CreateInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.CreateInfrastructureAccountParams) *api.Response
	CreateOrUpdateNodePool(ctx *gin.Context, params *node_pools.CreateOrUpdateNodePoolParams) *api.Response
	DeleteCluster(ctx *gin.Context, params *clusters.DeleteClusterParams) *api.Response
	DeleteConfigItem(ctx *gin.Context, params *config_items.DeleteConfigItemParams) *api.Response
	DeleteNodePool(ctx *gin.Context, params *node_pools.DeleteNodePoolParams) *api.Response
	GetCluster(ctx *gin.Context, params *clusters.GetClusterParams) *api.Response
	GetInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.GetInfrastructureAccountParams) *api.Response
	ListClusters(ctx *gin.Context, params *clusters.ListClustersParams) *api.Response
	ListInfrastructureAccounts(ctx *gin.Context) *api.Response
	ListNodePools(ctx *gin.Context, params *node_pools.ListNodePoolsParams) *api.Response
	UpdateCluster(ctx *gin.Context, params *clusters.UpdateClusterParams) *api.Response
	UpdateInfrastructureAccount(ctx *gin.Context, params *infrastructure_accounts.UpdateInfrastructureAccountParams) *api.Response
}

func ginizePath(path string) string {
	return strings.Replace(strings.Replace(path, "{", ":", -1), "}", "", -1)
}

// configureRoutes configures the routes for the API service.
func configureRoutes(service Service, enableAuth bool) *Routes {
	engine := gin.New()
	engine.Use(gin.Recovery())
	engine.Use(middleware.LogrusLogger())
	routes := &Routes{Engine: engine}

	routes.AddOrUpdateConfigItem.RouterGroup = routes.Group("")
	routes.AddOrUpdateConfigItem.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.AddOrUpdateConfigItem.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.AddOrUpdateConfigItem.RouterGroup.Use(routes.AddOrUpdateConfigItem.Auth)

	}
	routes.AddOrUpdateConfigItem.Post = routes.AddOrUpdateConfigItem.Group("")
	routes.AddOrUpdateConfigItem.Post.PUT(ginizePath("/kubernetes-clusters/{cluster_id}/config-items/{config_key}"), config_items.BusinessLogicAddOrUpdateConfigItem(service.AddOrUpdateConfigItem))

	routes.CreateCluster.RouterGroup = routes.Group("")
	routes.CreateCluster.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.CreateCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.CreateCluster.RouterGroup.Use(routes.CreateCluster.Auth)

	}
	routes.CreateCluster.Post = routes.CreateCluster.Group("")
	routes.CreateCluster.Post.POST(ginizePath("/kubernetes-clusters"), clusters.BusinessLogicCreateCluster(service.CreateCluster))

	routes.CreateInfrastructureAccount.RouterGroup = routes.Group("")
	routes.CreateInfrastructureAccount.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.CreateInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.CreateInfrastructureAccount.RouterGroup.Use(routes.CreateInfrastructureAccount.Auth)

	}
	routes.CreateInfrastructureAccount.Post = routes.CreateInfrastructureAccount.Group("")
	routes.CreateInfrastructureAccount.Post.POST(ginizePath("/infrastructure-accounts"), infrastructure_accounts.BusinessLogicCreateInfrastructureAccount(service.CreateInfrastructureAccount))

	routes.CreateOrUpdateNodePool.RouterGroup = routes.Group("")
	routes.CreateOrUpdateNodePool.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.CreateOrUpdateNodePool.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.CreateOrUpdateNodePool.RouterGroup.Use(routes.CreateOrUpdateNodePool.Auth)

	}
	routes.CreateOrUpdateNodePool.Post = routes.CreateOrUpdateNodePool.Group("")
	routes.CreateOrUpdateNodePool.Post.PUT(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools/{node_pool_name}"), node_pools.BusinessLogicCreateOrUpdateNodePool(service.CreateOrUpdateNodePool))

	routes.DeleteCluster.RouterGroup = routes.Group("")
	routes.DeleteCluster.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.DeleteCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.DeleteCluster.RouterGroup.Use(routes.DeleteCluster.Auth)

	}
	routes.DeleteCluster.Post = routes.DeleteCluster.Group("")
	routes.DeleteCluster.Post.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.BusinessLogicDeleteCluster(service.DeleteCluster))

	routes.DeleteConfigItem.RouterGroup = routes.Group("")
	routes.DeleteConfigItem.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.DeleteConfigItem.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.DeleteConfigItem.RouterGroup.Use(routes.DeleteConfigItem.Auth)

	}
	routes.DeleteConfigItem.Post = routes.DeleteConfigItem.Group("")
	routes.DeleteConfigItem.Post.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}/config-items/{config_key}"), config_items.BusinessLogicDeleteConfigItem(service.DeleteConfigItem))

	routes.DeleteNodePool.RouterGroup = routes.Group("")
	routes.DeleteNodePool.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.DeleteNodePool.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.DeleteNodePool.RouterGroup.Use(routes.DeleteNodePool.Auth)

	}
	routes.DeleteNodePool.Post = routes.DeleteNodePool.Group("")
	routes.DeleteNodePool.Post.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools/{node_pool_name}"), node_pools.BusinessLogicDeleteNodePool(service.DeleteNodePool))

	routes.GetCluster.RouterGroup = routes.Group("")
	if enableAuth {

		routes.GetCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.GetCluster.RouterGroup.Use(routes.GetCluster.Auth)

	}
	routes.GetCluster.Post = routes.GetCluster.Group("")
	routes.GetCluster.Post.GET(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.BusinessLogicGetCluster(service.GetCluster))

	routes.GetInfrastructureAccount.RouterGroup = routes.Group("")
	if enableAuth {

		routes.GetInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.GetInfrastructureAccount.RouterGroup.Use(routes.GetInfrastructureAccount.Auth)

	}
	routes.GetInfrastructureAccount.Post = routes.GetInfrastructureAccount.Group("")
	routes.GetInfrastructureAccount.Post.GET(ginizePath("/infrastructure-accounts/{account_id}"), infrastructure_accounts.BusinessLogicGetInfrastructureAccount(service.GetInfrastructureAccount))

	routes.ListClusters.RouterGroup = routes.Group("")
	if enableAuth {

		routes.ListClusters.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.ListClusters.RouterGroup.Use(routes.ListClusters.Auth)

	}
	routes.ListClusters.Post = routes.ListClusters.Group("")
	routes.ListClusters.Post.GET(ginizePath("/kubernetes-clusters"), clusters.BusinessLogicListClusters(service.ListClusters))

	routes.ListInfrastructureAccounts.RouterGroup = routes.Group("")
	if enableAuth {

		routes.ListInfrastructureAccounts.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.ListInfrastructureAccounts.RouterGroup.Use(routes.ListInfrastructureAccounts.Auth)

	}
	routes.ListInfrastructureAccounts.Post = routes.ListInfrastructureAccounts.Group("")
	routes.ListInfrastructureAccounts.Post.GET(ginizePath("/infrastructure-accounts"), infrastructure_accounts.BusinessLogicListInfrastructureAccounts(service.ListInfrastructureAccounts))

	routes.ListNodePools.RouterGroup = routes.Group("")
	if enableAuth {

		routes.ListNodePools.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.ListNodePools.RouterGroup.Use(routes.ListNodePools.Auth)

	}
	routes.ListNodePools.Post = routes.ListNodePools.Group("")
	routes.ListNodePools.Post.GET(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools"), node_pools.BusinessLogicListNodePools(service.ListNodePools))

	routes.UpdateCluster.RouterGroup = routes.Group("")
	routes.UpdateCluster.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.UpdateCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.UpdateCluster.RouterGroup.Use(routes.UpdateCluster.Auth)

	}
	routes.UpdateCluster.Post = routes.UpdateCluster.Group("")
	routes.UpdateCluster.Post.PATCH(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.BusinessLogicUpdateCluster(service.UpdateCluster))

	routes.UpdateInfrastructureAccount.RouterGroup = routes.Group("")
	routes.UpdateInfrastructureAccount.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routes.UpdateInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: "https://info.services.auth.zalando.com/oauth2/tokeninfo",
			},
		)

		routes.UpdateInfrastructureAccount.RouterGroup.Use(routes.UpdateInfrastructureAccount.Auth)

	}
	routes.UpdateInfrastructureAccount.Post = routes.UpdateInfrastructureAccount.Group("")
	routes.UpdateInfrastructureAccount.Post.PATCH(ginizePath("/infrastructure-accounts/{account_id}"), infrastructure_accounts.BusinessLogicUpdateInfrastructureAccount(service.UpdateInfrastructureAccount))

	return routes
}

// API defines the API service.
type API struct {
	Routes  *Routes
	config  *Config
	server  *http.Server
	Title   string
	Version string
}

// NewAPI initializes a new API service.
func NewAPI(svc Service, config *Config) *API {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	api := &API{
		Routes:  configureRoutes(svc, !config.AuthDisabled),
		config:  config,
		Title:   "Cluster Registry",
		Version: "0.0.1",
	}

	// enable pprof http endpoints in debug mode
	if config.Debug {
		pprof.Register(api.Routes.Engine, nil)
	}

	api.server = &http.Server{
		Addr:         config.Address,
		Handler:      api.Routes.Engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	if !config.WellKnownDisabled {
		api.Routes.configureWellKnown(svc.Healthy)
	}

	// configure healthz endpoint
	api.Routes.GET("/healthz", healthHandler(svc.Healthy))

	return api
}

// Run runs the API server it will listen on either HTTP or HTTPS depending on
// the config passed to NewAPI.
func (a *API) Run() error {
	log.Infof("Serving '%s - %s' on address %s", a.Title, a.Version, a.server.Addr)
	if a.config.InsecureHTTP {
		return a.server.ListenAndServe()
	}
	return a.server.ListenAndServeTLS(a.config.TLSCertFile, a.config.TLSKeyFile)
}

// Shutdown will gracefully shutdown the API server.
func (a *API) Shutdown() error {
	return a.server.Shutdown(context.Background())
}

// RunWithSigHandler runs the API server with SIGTERM handling automatically
// enabled. The server will listen for a SIGTERM signal and gracefully shutdown
// the web server.
// It's possible to optionally pass any number shutdown functions which will
// execute one by one after the webserver has been shutdown successfully.
func (a *API) RunWithSigHandler(shutdown ...func() error) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGTERM)

	go func() {
		<-sigCh
		a.Shutdown()
	}()

	err := a.Run()
	if err != nil {
		if err != http.ErrServerClosed {
			return err
		}
	}

	for _, fn := range shutdown {
		err := fn()
		if err != nil {
			return err
		}
	}

	return nil
}

// vim: ft=go
