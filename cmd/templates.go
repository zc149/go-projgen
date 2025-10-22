package cmd

import (
	"fmt"
	"os"
	"path/filepath"
)

func writeHelmChart(projectDir string) error {
	chartDir := filepath.Join(projectDir, "charts", projectName)
	if err := os.MkdirAll(filepath.Join(chartDir, "templates"), 0755); err != nil {
		return err
	}

	// Chart.yaml
	chartYaml := fmt.Sprintf(`apiVersion: v2
name: %s
version: 0.1.0
appVersion: "1.0.0"
`, projectName)
	if err := os.WriteFile(filepath.Join(chartDir, "Chart.yaml"), []byte(chartYaml), 0644); err != nil {
		return err
	}

	// values.yaml
	values := fmt.Sprintf(`image:
  repository: ghcr.io/YOUR_GITHUB/%s
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80
`, projectName)
	if err := os.WriteFile(filepath.Join(chartDir, "values.yaml"), []byte(values), 0644); err != nil {
		return err
	}

	// deployment.yaml
	deployment := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: {{ include "` + projectName + `.fullname" . }}
  labels:
    app: {{ include "` + projectName + `.name" . }}
spec:
  replicas: 1
  selector:
    matchLabels:
      app: {{ include "` + projectName + `.name" . }}
  template:
    metadata:
      labels:
        app: {{ include "` + projectName + `.name" . }}
    spec:
      containers:
        - name: api
          image: "{{ .Values.image.repository }}:{{ .Values.image.tag }}"
          ports:
            - containerPort: 8080
`
	if err := os.WriteFile(filepath.Join(chartDir, "templates", "deployment.yaml"), []byte(deployment), 0644); err != nil {
		return err
	}

	// service.yaml
	service := `apiVersion: v1
kind: Service
metadata:
  name: {{ include "` + projectName + `.fullname" . }}
spec:
  type: {{ .Values.service.type }}
  selector:
    app: {{ include "` + projectName + `.name" . }}
  ports:
    - name: http
      port: {{ .Values.service.port }}
      targetPort: 8080
`
	if err := os.WriteFile(filepath.Join(chartDir, "templates", "service.yaml"), []byte(service), 0644); err != nil {
		return err
	}

	return nil
}
