## Context

lango는 현재 Dockerfile은 있지만, 운영 환경 배포를 위한 docker-compose.yml이 없다. 또한 현재 Dockerfile은 headless 환경에서의 passphrase 처리를 고려하지 않아, LocalCryptoProvider 사용 시 TTY가 없어 실행이 불가능하다.

**현재 상태:**
- Dockerfile: 기본적인 multi-stage 빌드 존재
- docker-compose.yml: 없음
- Headless 문제: LocalCryptoProvider가 interactive terminal 필요

## Goals / Non-Goals

**Goals:**
- Docker 환경에서 lango를 안정적으로 운영할 수 있게 함
- RPC Provider (Companion) 사용을 통해 headless 환경에서도 암호화 가능하게 함
- 모든 채널(Discord, Telegram, Slack) 지원
- Browser tool 지원 (Chromium 포함)
- Health check를 통한 컨테이너 상태 모니터링
- 보안 강화 (non-root 사용자 실행)

**Non-Goals:**
- 멀티 환경(dev/staging/prod) 별도 구성
- Kubernetes/Helm 차트 작성
- Companion 앱 컨테이너화 (별도 프로젝트)
- LocalCryptoProvider의 Docker 환경 지원

## Decisions

### Decision 1: RPC Provider 강제 사용
Docker 환경에서는 LocalCryptoProvider를 사용하지 않고 RPC Provider (Companion 연동)를 강제한다.

**대안들:**
- A) 환경변수로 passphrase 전달 → 보안 취약 (env에 평문 노출)
- B) Docker Secret으로 passphrase 파일 마운트 → 복잡하고 LocalCryptoProvider 수정 필요
- C) **RPC Provider 강제** → 기존 보안 모델 유지, Companion과 분리된 아키텍처

**선택: Option C** - headless 환경에서 RPC Provider를 사용하는 것이 기존 설계 의도와 일치함.

### Decision 2: docker-compose.yml 구조

```yaml
services:
  lango:
    build: .
    ports: ["18789:18789"]
    volumes:
      - lango-data:/data
      - ./lango.json:/app/lango.json:ro
    environment:
      - ANTHROPIC_API_KEY
      - DISCORD_BOT_TOKEN
      - TELEGRAM_BOT_TOKEN
      - SLACK_BOT_TOKEN
      - SLACK_APP_TOKEN
```

**대안들:**
- A) 모든 설정을 환경변수로 → 복잡, lango의 config 시스템과 불일치
- B) **Config 파일 마운트 + 환경변수 치환** → 기존 설계와 일치

**선택: Option B**

### Decision 3: Dockerfile 개선 사항

1. **Health check 추가**: HTTP 엔드포인트로 상태 확인
2. **Non-root 사용자**: 보안 강화
3. **Chromium 유지**: Browser tool 필수

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Companion 미연결 시 lango 시작 불가 | 명확한 에러 메시지와 문서화 |
| Chromium으로 이미지 크기 증가 (~200MB) | Browser tool 필수 요구사항이므로 수용 |
| Config 파일 마운트 필요 | docker-compose.yml에 예시 포함 |
