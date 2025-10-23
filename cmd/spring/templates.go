package spring

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

	var imageRepo string
	switch registry {
	case "ecr":
		imageRepo = fmt.Sprintf("${AWS_ACCOUNT_ID}.dkr.ecr.ap-northeast-2.amazonaws.com/%s", projectName)
	default: // ghcr
		imageRepo = fmt.Sprintf("ghcr.io/${GITHUB_REPOSITORY}/%s", projectName)
	}

	// values.yaml
	values := fmt.Sprintf(`image:
  repository: %s
  tag: latest
  pullPolicy: IfNotPresent

service:
  type: ClusterIP
  port: 80
`, imageRepo)
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

func writeArgoCDManifest(projectDir string) error {
	manifestDir := filepath.Join(projectDir, "manifests")
	if err := os.MkdirAll(manifestDir, 0755); err != nil {
		return err
	}

	appYaml := fmt.Sprintf(`apiVersion: argoproj.io/v1alpha1
kind: Application
metadata:
  name: %s
  namespace: argocd
spec:
  project: default
  destination:
    server: https://kubernetes.default.svc
    namespace: default
  source:
    repoURL: https://github.com/${GITHUB_REPOSITORY}
    targetRevision: main
    path: charts/%s
  syncPolicy:
    automated:
      prune: true
      selfHeal: true
`, projectName, projectName)

	if err := os.WriteFile(filepath.Join(manifestDir, "argocd-app.yaml"), []byte(appYaml), 0644); err != nil {
		return err
	}

	return nil
}
