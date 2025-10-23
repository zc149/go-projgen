# projgen

CLI 기반 **GitOps 프로젝트 생성기**  
(Spring Boot / React 프로젝트를 자동 생성하고, Dockerfile · CI/CD 워크플로우 · Helm Chart · ArgoCD 매니페스트까지 구성)

👉 개발자는 **애플리케이션 코드 작성에만 집중**하고, 배포 파이프라인은 자동으로 세팅됩니다.

---

##  GitOps 스타일

이 템플릿은 **ArgoCD 기반 GitOps**를 따릅니다.

1. `projgen` CLI 로 새 프로젝트 생성  
2. 자동으로 다음 파일이 추가됨:
   - 프로젝트 소스 (Spring Boot / React Vite)
   - Dockerfile
   - GitHub Actions CI/CD 워크플로우
   - Helm Chart
   - ArgoCD Application 매니페스트  
3. `git push` → GitHub Actions 가 빌드 & 컨테이너 레지스트리 푸시  
4. ArgoCD 가 레포를 watch 하며 변경사항 자동 배포  

⚠️ **주의**  
`projgen` 은 **템플릿만 생성**합니다.  
실제 파이프라인이 동작하려면 **레포지토리 Secrets** 또는 `.env` 값들을 사용자가 직접 설정해야 합니다.

---

## 🚀 설치 방법

`projgen` 은 소스 빌드 없이 **릴리즈된 실행 파일**을 바로 사용하면 됩니다.  
GitHub Releases 페이지에서 운영체제에 맞는 바이너리를 다운로드하세요.

- [Releases 페이지 바로가기](https://github.com/zc149/go-projgen/releases/tag/v1.0.0)

예:
- Windows → `projgen_windows_amd64.exe`
- macOS (Intel) → `projgen_darwin_amd64`
- macOS (Apple Silicon) → `projgen_darwin_arm64`
- Linux → `projgen_linux_amd64`

다운로드 후 PATH 에 추가하거나, 원하는 디렉토리에 두고 실행하면 됩니다.

### 1) 직접 실행
압축 해제 후 실행 권한을 주고 바로 실행할 수 있습니다.
```bash
chmod +x projgen-darwin-arm64
./projgen-darwin-arm64 spring --help
```

자주 사용할 경우, 실행 파일을 PATH 경로에 옮겨두면 어디서든 projgen 명령으로 실행 가능합니다.

---

## ⚙️ 필수 환경 변수 (.env)

GitHub 에 푸시하려면 **반드시 두 개 변수가 필요**합니다.

프로젝트 생성을 위해 `.env` 파일을 **실행 파일과 같은 경로**에 준비하세요.

```env
GITHUB_TOKEN=ghp_xxx   # GitHub Personal Access Token (repo, workflow 권한 필수)
GITHUB_OWNER=kimjikwan # GitHub username 또는 org name
```
👉 GITHUB_TOKEN 은 [GitHub Developer Settings > Personal Access Token (classic)] 에서 repo, workflow 권한으로 발급하세요.
👉 .env 파일에 저장 후 CLI 실행 시 자동으로 참조됩니다.

---

## 🔧 사용법

Spring Boot 프로젝트 생성 예시:

```
./projgen-darwin-arm64 spring \
  --name my-spring-app \
  --group com.mycompany \
  --artifact my-spring-app \
  --package com.mycompany.myapp \
  --java 17 \
  --build maven \
  --push \
  --private \
  --registry ghcr
```
React 프로젝트 생성 예시:

```
./projgen-darwin-arm64 react \
  --name my-react-app \
  --node 20 \
  --push \
  --private \
  --registry ghcr
``` 

ℹ️ 위 예시는 참고용입니다.  
모든 옵션과 기본값은 `--help` 플래그로 확인하세요:

```bash
./projgen-darwin-arm64 spring --help
./projgen-darwin-arm64 react --help

