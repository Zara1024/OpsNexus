package handler

import (
	"github.com/gin-gonic/gin"

	"github.com/Zara1024/OpsNexus/cloudops-server/internal/k8s/service"
)

// Handler 聚合 k8s 模块处理器。
type Handler struct {
	services *service.Services
}

// New 创建 k8s 模块处理器。
func New(services *service.Services) *Handler {
	return &Handler{services: services}
}

// RegisterRoutes 注册 k8s HTTP 路由。
func (h *Handler) RegisterRoutes(v1 *gin.RouterGroup, jwtMiddleware gin.HandlerFunc, permissionMiddleware func(string) gin.HandlerFunc, auditMiddleware gin.HandlerFunc) {
	k8s := v1.Group("/k8s")
	k8s.Use(jwtMiddleware)
	{
		clusters := k8s.Group("/clusters")
		{
			clusters.GET("", permissionMiddleware("k8s:cluster:list"), h.ListClusters)
			clusters.POST("", permissionMiddleware("k8s:cluster:create"), auditMiddleware, h.CreateCluster)
			clusters.GET("/:id", permissionMiddleware("k8s:cluster:detail"), h.GetCluster)
			clusters.PUT("/:id", permissionMiddleware("k8s:cluster:update"), auditMiddleware, h.UpdateCluster)
			clusters.DELETE("/:id", permissionMiddleware("k8s:cluster:delete"), auditMiddleware, h.DeleteCluster)
			clusters.POST("/:id/test", permissionMiddleware("k8s:cluster:test"), h.TestCluster)

			clusters.GET("/:id/nodes", permissionMiddleware("k8s:node:list"), h.ListNodes)
			clusters.PUT("/:id/nodes/:name/cordon", permissionMiddleware("k8s:node:cordon"), auditMiddleware, h.CordonNode)
			clusters.PUT("/:id/nodes/:name/uncordon", permissionMiddleware("k8s:node:uncordon"), auditMiddleware, h.UncordonNode)
			clusters.POST("/:id/nodes/:name/drain", permissionMiddleware("k8s:node:drain"), auditMiddleware, h.DrainNode)

			clusters.GET("/:id/namespaces", permissionMiddleware("k8s:namespace:list"), h.ListNamespaces)

			clusters.GET("/:id/ns/:ns/deployments", permissionMiddleware("k8s:workload:deployment:list"), h.ListDeployments)
			clusters.POST("/:id/ns/:ns/deployments", permissionMiddleware("k8s:workload:deployment:create"), auditMiddleware, h.CreateDeployment)
			clusters.GET("/:id/ns/:ns/deployments/:name", permissionMiddleware("k8s:workload:deployment:detail"), h.GetDeployment)
			clusters.PUT("/:id/ns/:ns/deployments/:name", permissionMiddleware("k8s:workload:deployment:update"), auditMiddleware, h.UpdateDeployment)
			clusters.DELETE("/:id/ns/:ns/deployments/:name", permissionMiddleware("k8s:workload:deployment:delete"), auditMiddleware, h.DeleteDeployment)
			clusters.PUT("/:id/ns/:ns/deployments/:name/scale", permissionMiddleware("k8s:workload:deployment:scale"), auditMiddleware, h.ScaleDeployment)
			clusters.POST("/:id/ns/:ns/deployments/:name/restart", permissionMiddleware("k8s:workload:deployment:restart"), auditMiddleware, h.RestartDeployment)

			clusters.GET("/:id/ns/:ns/pods", permissionMiddleware("k8s:workload:pod:list"), h.ListPods)
			clusters.GET("/:id/ns/:ns/pods/:pod", permissionMiddleware("k8s:workload:pod:detail"), h.GetPod)
			clusters.GET("/:id/ns/:ns/pods/:pod/logs", permissionMiddleware("k8s:workload:pod:log"), h.GetPodLogs)
			clusters.DELETE("/:id/ns/:ns/pods/:pod", permissionMiddleware("k8s:workload:pod:delete"), auditMiddleware, h.DeletePod)

			clusters.GET("/:id/ns/:ns/services", permissionMiddleware("k8s:network:service:list"), h.ListServices)
			clusters.GET("/:id/ns/:ns/ingresses", permissionMiddleware("k8s:network:ingress:list"), h.ListIngresses)

			clusters.GET("/:id/ns/:ns/configmaps", permissionMiddleware("k8s:config:configmap:list"), h.ListConfigMaps)
			clusters.POST("/:id/ns/:ns/configmaps", permissionMiddleware("k8s:config:configmap:save"), auditMiddleware, h.SaveConfigMap)
			clusters.DELETE("/:id/ns/:ns/configmaps/:name", permissionMiddleware("k8s:config:configmap:delete"), auditMiddleware, h.DeleteConfigMap)
			clusters.GET("/:id/ns/:ns/secrets", permissionMiddleware("k8s:config:secret:list"), h.ListSecrets)
			clusters.POST("/:id/ns/:ns/secrets", permissionMiddleware("k8s:config:secret:save"), auditMiddleware, h.SaveSecret)
			clusters.DELETE("/:id/ns/:ns/secrets/:name", permissionMiddleware("k8s:config:secret:delete"), auditMiddleware, h.DeleteSecret)
		}
	}
}
