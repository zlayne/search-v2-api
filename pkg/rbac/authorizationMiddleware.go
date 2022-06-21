package rbac

import (
	"net/http"

	"k8s.io/klog/v2"
)

func AuthorizeUser(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		_, err := cacheInst.ClusterScopedResources(r.Context())
		if err != nil {
			klog.Warning("Unexpected error while obtaining cluster-scoped resources.", err)
		}
		klog.Info("Finished getting resources. Now gettng user namespaces..")

		clientToken := r.Context().Value(ContextAuthTokenKey).(string)
		_, newerr := cacheInst.NamespacedResources(r.Context(), clientToken)
		if newerr != nil {
			klog.Warning("Unexpected error while obtaining namesapces.", newerr)
		}

		klog.Info("Finished getting resources. Now Authorizing..")
		next.ServeHTTP(w, r.WithContext(r.Context()))

	})
}