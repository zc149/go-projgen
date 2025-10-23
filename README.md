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

## ⚙️ 필수 환경 변수 (.env)

GitHub 에 푸시하려면 **반드시 두 개 변수가 필요**합니다.

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
go run main.go spring \
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
go run main.go react \
  --name my-react-app \
  --node 20 \
  --push \
  --private \
  --registry ghcr
``` 

ℹ️ 위 예시는 참고용입니다.  
모든 옵션과 기본값은 `--help` 플래그로 확인하세요:

```bash
go run main.go spring --help
go run main.go react --help

