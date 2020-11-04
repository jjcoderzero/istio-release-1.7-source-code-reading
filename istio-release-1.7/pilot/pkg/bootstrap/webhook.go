package bootstrap

import (
	"crypto/tls"
	"net/http"
	"net/url"
	"time"

	"istio.io/pkg/log"
)

const (
	HTTPSHandlerReadyPath = "/httpsReady"
)

func (s *Server) initSecureWebhookServer(args *PilotArgs) {
	if s.kubeClient == nil {
		return
	}

	log.Info("initializing secure webhook server for istiod webhooks")
	// create the https server for hosting the k8s injectionWebhook handlers.
	s.httpsMux = http.NewServeMux()
	s.httpsServer = &http.Server{
		Addr:    args.ServerOptions.HTTPSAddr,
		Handler: s.httpsMux,
		TLSConfig: &tls.Config{
			GetCertificate: s.getIstiodCertificate,
		},
	}

	// setup our readiness handler and the corresponding client we'll use later to check it with.
	s.httpsMux.HandleFunc(HTTPSHandlerReadyPath, func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusOK)
	})
	s.httpsReadyClient = &http.Client{
		Timeout: time.Second,
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}
	s.addReadinessProbe("Secure Webhook Server", s.webhookReadyHandler)
}

func (s *Server) webhookReadyHandler() (bool, error) {
	req := &http.Request{
		Method: http.MethodGet,
		URL: &url.URL{
			Scheme: "https",
			Host:   s.httpsServer.Addr,
			Path:   HTTPSHandlerReadyPath,
		},
	}

	response, err := s.httpsReadyClient.Do(req)
	if err != nil {
		return false, err
	}
	defer response.Body.Close()

	return response.StatusCode == http.StatusOK, nil
}
