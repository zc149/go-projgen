package cmd

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	projectName  string
	groupId      string
	artifactId   string
	packageName  string
	buildTool    string
	javaVersion  int
	pushToGitHub bool
)
var springCmd = &cobra.Command{
	Use:   "spring",
	Short: "Generate a new Spring Boot project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("🚀 Generating Spring Boot project: %s\n", projectName)
		return generateSpringProject()
	},
}

func init() {
	springCmd.Flags().StringVar(&projectName, "name", "demo", "Project name")
	springCmd.Flags().StringVar(&groupId, "group", "com.example", "Group ID")
	springCmd.Flags().StringVar(&artifactId, "artifact", "demo", "Artifact ID")
	springCmd.Flags().StringVar(&packageName, "package", "com.example.demo", "Base package")
	springCmd.Flags().StringVar(&buildTool, "build", "maven", "Build tool (maven|gradle|gradle-kotlin)")
	springCmd.Flags().IntVar(&javaVersion, "java", 17, "Java version")
	springCmd.Flags().BoolVar(&pushToGitHub, "push", false, "Push project to GitHub")
	rootCmd.AddCommand(springCmd)
}

func generateSpringProject() error {
	buildType := "maven-project"
	switch buildTool {
	case "gradle":
		buildType = "gradle-project"
	case "gradle-kotlin":
		buildType = "gradle-project-kotlin"
	}

	url := fmt.Sprintf(
		"https://start.spring.io/starter.zip?type=%s&language=java&groupId=%s&artifactId=%s&name=%s&packageName=%s&javaVersion=%d&dependencies=web,actuator",
		buildType, groupId, artifactId, projectName, packageName, javaVersion,
	)

	// 1. ZIP 다운로드
	resp, err := http.Get(url)
	if err != nil {
		return fmt.Errorf("failed to download: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		return fmt.Errorf("spring initializr error: %s", resp.Status)
	}

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, resp.Body); err != nil {
		return err
	}

	// 2. 압축 해제
	zr, err := zip.NewReader(bytes.NewReader(buf.Bytes()), int64(buf.Len()))
	if err != nil {
		return err
	}

	if err := os.MkdirAll(projectName, 0755); err != nil {
		return err
	}

	for _, f := range zr.File {
		path := filepath.Join(projectName, f.Name)
		if f.FileInfo().IsDir() {
			if err := os.MkdirAll(path, f.Mode()); err != nil {
				return err
			}
			continue
		}
		// 파일 추출
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		rc, err := f.Open()
		if err != nil {
			return err
		}
		defer rc.Close()
		out, err := os.Create(path)
		if err != nil {
			return err
		}
		if _, err := io.Copy(out, rc); err != nil {
			return err
		}
		out.Close()
	}

	fmt.Println("✅ Project generated at", projectName)

	// 템플릿 파일 추가
	if err := writeTemplates(projectName); err != nil {
		return err
	}
	fmt.Println("✅ Added Dockerfile & CI workflow")

	// Helm 차트 추가
	if err := writeHelmChart(projectName); err != nil {
		return err
	}
	fmt.Println("✅ Added Helm chart")

	// GitHub Repo 생성 & Push
	if pushToGitHub {
		if err := initGitAndPushAPI(projectName); err != nil {
			return err
		}
		fmt.Println("✅ GitHub repo created & pushed")
	}

	return nil
}

func writeTemplates(projectDir string) error {
	// 1. Dockerfile
	dockerfile := `# syntax=docker/dockerfile:1
FROM eclipse-temurin:17-jre
WORKDIR /app
COPY target/*.jar app.jar
EXPOSE 8080
ENTRYPOINT ["java","-jar","/app/app.jar"]
`
	if err := os.WriteFile(filepath.Join(projectDir, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
		return err
	}

	// 2. GitHub Actions 워크플로 (CI)
	workflowDir := filepath.Join(projectDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	ci := `name: CI

on:
  push:
    branches: [ main ]
  pull_request:

jobs:
  build-test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-java@v4
        with:
          distribution: temurin
          java-version: '17'
      - name: Build with Maven
        run: mvn -q -DskipTests=false clean verify

  docker:
    needs: build-test
    runs-on: ubuntu-latest
    permissions:
      contents: read
      packages: write
    env:
      IMAGE: ghcr.io/${{ github.repository }}:$(echo ${{ github.sha }} | cut -c1-10)
    steps:
      - uses: actions/checkout@v4
      - name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}
      - name: Build jar
        run: mvn -q -DskipTests package
      - name: Build & Push Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            ${{ env.IMAGE }}
            ghcr.io/${{ github.repository }}:latest
`
	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(ci), 0644); err != nil {
		return err
	}

	return nil
}
