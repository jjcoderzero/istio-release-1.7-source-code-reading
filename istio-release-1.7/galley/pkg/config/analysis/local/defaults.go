package local

import (
	"bytes"
	"text/template"
)

const defaultIstioIngressGateway = `
apiVersion: v1
kind: Pod
metadata:
  labels:
    istio: ingressgateway
  name: {{.ingressService}}-dummypod
  namespace: {{.namespace}}
spec:
  containers:
    - args:
      name: istio-proxy
---
apiVersion: v1
kind: Service
metadata:
  name: {{.ingressService}}
  namespace: {{.namespace}}
spec:
  ports:
  - name: http2
    nodePort: 31380
    port: 80
    protocol: TCP
    targetPort: 80
  - name: https
    nodePort: 31390
    port: 443
    protocol: TCP
    targetPort: 443
  - name: tcp
    nodePort: 31400
    port: 31400
    protocol: TCP
    targetPort: 31400
  - name: tls
    nodePort: 31447
    port: 15443
    protocol: TCP
    targetPort: 15443
  selector:
    istio: ingressgateway
`

func getDefaultIstioIngressGateway(namespace, ingressService string) (string, error) {
	result, err := generate(defaultIstioIngressGateway, map[string]string{"namespace": namespace, "ingressService": ingressService})
	if err != nil {
		return "", err
	}

	return result, nil
}

func generate(tmpl string, params map[string]string) (string, error) {
	t := template.Must(template.New("code").Parse(tmpl))

	var b bytes.Buffer
	if err := t.Execute(&b, params); err != nil {
		return "", err
	}
	return b.String(), nil
}
