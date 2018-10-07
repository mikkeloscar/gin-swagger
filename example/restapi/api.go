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
	"github.com/mikkeloscar/gin-swagger/tracing"
	opentracing "github.com/opentracing/opentracing-go"
	log "github.com/sirupsen/logrus"
	ginoauth2 "github.com/zalando/gin-oauth2"
)

// Routes defines all the routes of the Server service.
type Routes struct {
	*gin.Engine
	AddOrUpdateConfigItem struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	CreateCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	CreateInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	CreateOrUpdateNodePool struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	DeleteCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	DeleteConfigItem struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	DeleteNodePool struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	GetCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	GetInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	ListClusters struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	ListInfrastructureAccounts struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	ListNodePools struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	UpdateCluster struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
	}
	UpdateInfrastructureAccount struct {
		*gin.RouterGroup
		Auth gin.HandlerFunc
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
// business logic for the Server service.
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

// initializeRoutes initializes the route structure for the Server service.
func initializeRoutes(enableAuth bool, tokenURL string, tracer opentracing.Tracer) *Routes {
	engine := gin.New()
	engine.Use(gin.Recovery())
	routes := &Routes{Engine: engine}

	routes.AddOrUpdateConfigItem.RouterGroup = routes.Group("")
	routes.AddOrUpdateConfigItem.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.AddOrUpdateConfigItem.RouterGroup.Use(tracing.InitSpan(tracer, "add_or_update_config_item"))
	}
	routes.AddOrUpdateConfigItem.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.AddOrUpdateConfigItem.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.CreateCluster.RouterGroup = routes.Group("")
	routes.CreateCluster.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.CreateCluster.RouterGroup.Use(tracing.InitSpan(tracer, "create_cluster"))
	}
	routes.CreateCluster.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.CreateCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.CreateInfrastructureAccount.RouterGroup = routes.Group("")
	routes.CreateInfrastructureAccount.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.CreateInfrastructureAccount.RouterGroup.Use(tracing.InitSpan(tracer, "create_infrastructure_account"))
	}
	routes.CreateInfrastructureAccount.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.CreateInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.CreateOrUpdateNodePool.RouterGroup = routes.Group("")
	routes.CreateOrUpdateNodePool.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.CreateOrUpdateNodePool.RouterGroup.Use(tracing.InitSpan(tracer, "create_or_update_node_pool"))
	}
	routes.CreateOrUpdateNodePool.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.CreateOrUpdateNodePool.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.DeleteCluster.RouterGroup = routes.Group("")
	routes.DeleteCluster.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.DeleteCluster.RouterGroup.Use(tracing.InitSpan(tracer, "delete_cluster"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.DeleteCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.DeleteConfigItem.RouterGroup = routes.Group("")
	routes.DeleteConfigItem.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.DeleteConfigItem.RouterGroup.Use(tracing.InitSpan(tracer, "delete_config_item"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.DeleteConfigItem.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.DeleteNodePool.RouterGroup = routes.Group("")
	routes.DeleteNodePool.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.DeleteNodePool.RouterGroup.Use(tracing.InitSpan(tracer, "delete_node_pool"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.DeleteNodePool.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.GetCluster.RouterGroup = routes.Group("")
	routes.GetCluster.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.GetCluster.RouterGroup.Use(tracing.InitSpan(tracer, "get_cluster"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.GetCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.GetInfrastructureAccount.RouterGroup = routes.Group("")
	routes.GetInfrastructureAccount.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.GetInfrastructureAccount.RouterGroup.Use(tracing.InitSpan(tracer, "get_infrastructure_account"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.GetInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.ListClusters.RouterGroup = routes.Group("")
	routes.ListClusters.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.ListClusters.RouterGroup.Use(tracing.InitSpan(tracer, "list_clusters"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.ListClusters.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.ListInfrastructureAccounts.RouterGroup = routes.Group("")
	routes.ListInfrastructureAccounts.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.ListInfrastructureAccounts.RouterGroup.Use(tracing.InitSpan(tracer, "list_infrastructure_accounts"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.ListInfrastructureAccounts.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.ListNodePools.RouterGroup = routes.Group("")
	routes.ListNodePools.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.ListNodePools.RouterGroup.Use(tracing.InitSpan(tracer, "list_node_pools"))
	}
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.ListNodePools.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.UpdateCluster.RouterGroup = routes.Group("")
	routes.UpdateCluster.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.UpdateCluster.RouterGroup.Use(tracing.InitSpan(tracer, "update_cluster"))
	}
	routes.UpdateCluster.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.UpdateCluster.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	routes.UpdateInfrastructureAccount.RouterGroup = routes.Group("")
	routes.UpdateInfrastructureAccount.RouterGroup.Use(middleware.LogrusLogger())
	if tracer != nil {
		routes.UpdateInfrastructureAccount.RouterGroup.Use(tracing.InitSpan(tracer, "update_infrastructure_account"))
	}
	routes.UpdateInfrastructureAccount.RouterGroup.Use(middleware.ContentTypes("application/json"))
	if enableAuth {

		routeTokenURL := tokenURL
		if routeTokenURL == "" {
			routeTokenURL = "https://info.services.auth.zalando.com/oauth2/tokeninfo"
		}
		routes.UpdateInfrastructureAccount.Auth = ginoauth2.Auth(
			middleware.ScopesAuth("uid"),
			oauth2.Endpoint{
				TokenURL: routeTokenURL,
			},
		)

	}

	return routes
}

// Server defines the Server service.
type Server struct {
	Routes           *Routes
	config           *Config
	server           *http.Server
	service          Service
	healthy          bool
	serviceHealthyFn func() bool
	authDisabled     bool
	Title            string
	Version          string
}

// NewServer initializes a new Server service.
func NewServer(svc Service, config *Config) *Server {
	if !config.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	server := &Server{
		Routes: initializeRoutes(
			!config.AuthDisabled,
			config.TokenURL,
			config.Tracer,
		),
		service:      svc,
		config:       config,
		Title:        "Cluster Registry",
		Version:      "0.0.1",
		authDisabled: config.AuthDisabled,
	}

	// enable pprof http endpoints in debug mode
	if config.Debug {
		pprof.Register(server.Routes.Engine)
	}

	// set logrus logger to TextFormatter with no colors
	log.SetFormatter(&log.TextFormatter{DisableColors: true})

	server.server = &http.Server{
		Addr:         config.Address,
		Handler:      server.Routes.Engine,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	server.serviceHealthyFn = svc.Healthy

	if !config.WellKnownDisabled {
		server.Routes.configureWellKnown(server.isHealthy)
	}

	// configure healthz endpoint
	server.Routes.GET("/healthz", healthHandler(server.isHealthy))

	return server
}

// isHealthy returns true if both the server and the service reports healthy.
func (s *Server) isHealthy() bool {
	return s.healthy && s.serviceHealthyFn()
}

// ConfigureRoutes starts the internal configureRoutes methode.
func (s *Server) ConfigureRoutes() {
	s.configureRoutes()
}

// configureRoutes configures the routes for the Server service.
// Configuring of routes includes setting up Auth if it is enabled.
func (s *Server) configureRoutes() {
	if !s.authDisabled {
		s.Routes.AddOrUpdateConfigItem.Use(s.Routes.AddOrUpdateConfigItem.Auth)
		s.Routes.CreateCluster.Use(s.Routes.CreateCluster.Auth)
		s.Routes.CreateInfrastructureAccount.Use(s.Routes.CreateInfrastructureAccount.Auth)
		s.Routes.CreateOrUpdateNodePool.Use(s.Routes.CreateOrUpdateNodePool.Auth)
		s.Routes.DeleteCluster.Use(s.Routes.DeleteCluster.Auth)
		s.Routes.DeleteConfigItem.Use(s.Routes.DeleteConfigItem.Auth)
		s.Routes.DeleteNodePool.Use(s.Routes.DeleteNodePool.Auth)
		s.Routes.GetCluster.Use(s.Routes.GetCluster.Auth)
		s.Routes.GetInfrastructureAccount.Use(s.Routes.GetInfrastructureAccount.Auth)
		s.Routes.ListClusters.Use(s.Routes.ListClusters.Auth)
		s.Routes.ListInfrastructureAccounts.Use(s.Routes.ListInfrastructureAccounts.Auth)
		s.Routes.ListNodePools.Use(s.Routes.ListNodePools.Auth)
		s.Routes.UpdateCluster.Use(s.Routes.UpdateCluster.Auth)
		s.Routes.UpdateInfrastructureAccount.Use(s.Routes.UpdateInfrastructureAccount.Auth)
	}

	// setup all service routes after the authenticate middleware has been
	// initialized.
	s.Routes.AddOrUpdateConfigItem.PUT(ginizePath("/kubernetes-clusters/{cluster_id}/config-items/{config_key}"), config_items.AddOrUpdateConfigItemEndpoint(s.service.AddOrUpdateConfigItem))
	s.Routes.CreateCluster.POST(ginizePath("/kubernetes-clusters"), clusters.CreateClusterEndpoint(s.service.CreateCluster))
	s.Routes.CreateInfrastructureAccount.POST(ginizePath("/infrastructure-accounts"), infrastructure_accounts.CreateInfrastructureAccountEndpoint(s.service.CreateInfrastructureAccount))
	s.Routes.CreateOrUpdateNodePool.PUT(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools/{node_pool_name}"), node_pools.CreateOrUpdateNodePoolEndpoint(s.service.CreateOrUpdateNodePool))
	s.Routes.DeleteCluster.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.DeleteClusterEndpoint(s.service.DeleteCluster))
	s.Routes.DeleteConfigItem.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}/config-items/{config_key}"), config_items.DeleteConfigItemEndpoint(s.service.DeleteConfigItem))
	s.Routes.DeleteNodePool.DELETE(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools/{node_pool_name}"), node_pools.DeleteNodePoolEndpoint(s.service.DeleteNodePool))
	s.Routes.GetCluster.GET(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.GetClusterEndpoint(s.service.GetCluster))
	s.Routes.GetInfrastructureAccount.GET(ginizePath("/infrastructure-accounts/{account_id}"), infrastructure_accounts.GetInfrastructureAccountEndpoint(s.service.GetInfrastructureAccount))
	s.Routes.ListClusters.GET(ginizePath("/kubernetes-clusters"), clusters.ListClustersEndpoint(s.service.ListClusters))
	s.Routes.ListInfrastructureAccounts.GET(ginizePath("/infrastructure-accounts"), infrastructure_accounts.ListInfrastructureAccountsEndpoint(s.service.ListInfrastructureAccounts))
	s.Routes.ListNodePools.GET(ginizePath("/kubernetes-clusters/{cluster_id}/node-pools"), node_pools.ListNodePoolsEndpoint(s.service.ListNodePools))
	s.Routes.UpdateCluster.PATCH(ginizePath("/kubernetes-clusters/{cluster_id}"), clusters.UpdateClusterEndpoint(s.service.UpdateCluster))
	s.Routes.UpdateInfrastructureAccount.PATCH(ginizePath("/infrastructure-accounts/{account_id}"), infrastructure_accounts.UpdateInfrastructureAccountEndpoint(s.service.UpdateInfrastructureAccount))
}

// Run runs the Server. It will listen on either HTTP or HTTPS depending on the
// config passed to NewServer.
func (s *Server) Run() error {
	// configure service routes
	s.configureRoutes()

	log.Infof("Serving '%s - %s' on address %s", s.Title, s.Version, s.server.Addr)
	// server is set to healthy when started.
	s.healthy = true
	if s.config.InsecureHTTP {
		return s.server.ListenAndServe()
	}
	return s.server.ListenAndServeTLS(s.config.TLSCertFile, s.config.TLSKeyFile)
}

// Shutdown will gracefully shutdown the Server server.
func (s *Server) Shutdown() error {
	// server is set to unhealthy when shutting down
	s.healthy = false
	return s.server.Shutdown(context.Background())
}

// RunWithSigHandler runs the Server server with SIGTERM handling automatically
// enabled. The server will listen for a SIGTERM signal and gracefully shutdown
// the web server.
// It's possible to optionally pass any number shutdown functions which will
// execute one by one after the webserver has been shutdown successfully.
func (s *Server) RunWithSigHandler(shutdown ...func() error) error {
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigCh
		s.Shutdown()
	}()

	err := s.Run()
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
