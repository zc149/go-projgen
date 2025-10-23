package react

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"projgen/cmd"
	"projgen/cmd/common"

	"github.com/spf13/cobra"
)

var (
	projectName  string
	nodeVersion  string
	pushToGitHub bool
	isPrivate    bool
	registry     string
)

var reactCmd = &cobra.Command{
	Use:   "react",
	Short: "Generate a new React (Vite) project",
	RunE: func(cmd *cobra.Command, args []string) error {
		fmt.Printf("⚡ Generating React project: %s\n", projectName)
		return generateReactProject()
	},
}

func init() {
	reactCmd.Flags().StringVar(&projectName, "name", "my-react-app", "Project name")
	reactCmd.Flags().StringVar(&nodeVersion, "node", "20", "Node.js version")
	reactCmd.Flags().BoolVar(&pushToGitHub, "push", false, "Push project to GitHub")
	reactCmd.Flags().BoolVar(&isPrivate, "private", false, "Create GitHub repository as private")
	reactCmd.Flags().StringVar(&registry, "registry", "ghcr", "Container registry (ghcr|ecr)")
	cmd.RootCmd.AddCommand(reactCmd)
}

func generateReactProject() error {
	// Vite 프로젝트 생성 (기본: react-ts)
	c := exec.Command("npm", "create", "vite@latest", projectName, "--", "--template", "react-ts")
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("failed to create vite project: %w", err)
	}

	// npm install
	c = exec.Command("npm", "install")
	c.Dir = projectName
	c.Stdout, c.Stderr = os.Stdout, os.Stderr
	if err := c.Run(); err != nil {
		return fmt.Errorf("failed npm install: %w", err)
	}

	fmt.Println("✅ React project generated:", projectName)

	// Dockerfile, CI/CD 작성
	if err := writeReactTemplates(projectName); err != nil {
		return err
	}
	fmt.Println("✅ Added Dockerfile & CI workflow for React")

	// Helm 차트 추가
	if err := writeHelmChart(projectName); err != nil {
		return err
	}
	fmt.Println("✅ Added Helm chart")

	// ArgoCD 매니페스트 추가
	if err := writeArgoCDManifest(projectName); err != nil {
		return err
	}
	fmt.Println("✅ Added Argo CD manifest")

	// GitHub push
	if pushToGitHub {
		if err := common.InitGitAndPushAPI(projectName, isPrivate); err != nil {
			return err
		}
		fmt.Println("✅ GitHub repo created & pushed")
	}

	return nil
}

func writeReactTemplates(projectDir string) error {
	// 1. Dockerfile
	dockerfile := `# syntax=docker/dockerfile:1
FROM node:20 AS build
WORKDIR /app
COPY package*.json ./
RUN npm install
COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=build /app/dist /usr/share/nginx/html
EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
`
	if err := os.WriteFile(filepath.Join(projectDir, "Dockerfile"), []byte(dockerfile), 0644); err != nil {
		return err
	}

	// 2. GitHub Actions 워크플로 (CI/CD)
	workflowDir := filepath.Join(projectDir, ".github", "workflows")
	if err := os.MkdirAll(workflowDir, 0755); err != nil {
		return err
	}

	// 레지스트리 분기
	var loginStep, imageRepo string
	switch registry {
	case "ecr":
		loginStep = `- name: Configure AWS credentials
        uses: aws-actions/configure-aws-credentials@v4
        with:
          role-to-assume: ${{ secrets.AWS_ROLE }}
          aws-region: ap-northeast-2

      - name: Login to Amazon ECR
        id: login-ecr
        uses: aws-actions/amazon-ecr-login@v2`
		imageRepo = "${{ secrets.AWS_ACCOUNT_ID }}.dkr.ecr.ap-northeast-2.amazonaws.com/" + projectName
	default: // ghcr
		loginStep = `- name: Login to GHCR
        uses: docker/login-action@v3
        with:
          registry: ghcr.io
          username: ${{ github.actor }}
          password: ${{ secrets.GITHUB_TOKEN }}`
		imageRepo = "ghcr.io/${{ github.repository }}"
	}

	ci := fmt.Sprintf(`name: CI

on:
  push:
    branches: [ main ]
  pull_request:

permissions:
  contents: read
  packages: write

jobs:
  build-and-docker:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: actions/setup-node@v4
        with:
          node-version: %s

      - name: Install & Build
        run: |
          npm ci
          npm run build

      %s

      - name: Set short SHA
        run: echo "SHORT_SHA=${GITHUB_SHA::10}" >> $GITHUB_ENV

      - name: Build & Push Docker image
        uses: docker/build-push-action@v6
        with:
          context: .
          push: true
          tags: |
            %s:${{ env.SHORT_SHA }}
            %s:latest
`, nodeVersion, loginStep, imageRepo, imageRepo)

	if err := os.WriteFile(filepath.Join(workflowDir, "ci.yml"), []byte(ci), 0644); err != nil {
		return err
	}

	return nil
}
