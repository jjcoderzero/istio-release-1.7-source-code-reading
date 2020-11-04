package bootstrap

import (
	"istio.io/istio/pilot/pkg/serviceregistry/kube/controller"

	"istio.io/pkg/log"
)

// initClusterRegistries starts the secret controller to watch for remote
// clusters and initialize the multicluster structures.
func (s *Server) initClusterRegistries(args *PilotArgs) (err error) {
	if hasKubeRegistry(args.RegistryOptions.Registries) {
		log.Info("initializing Kubernetes cluster registry")
		mc, err := controller.NewMulticluster(s.kubeClient,
			args.RegistryOptions.ClusterRegistriesNamespace,
			args.RegistryOptions.KubeOptions,
			s.ServiceController(),
			s.EnvoyXdsServer,
			s.environment)

		if err != nil {
			log.Info("Unable to create new Multicluster object")
			return err
		}

		s.multicluster = mc
	}
	return nil
}
