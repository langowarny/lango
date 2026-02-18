## Why

Docker 컨테이너 환경에서 lango를 운영 가능하게 만들기 위해 Dockerfile 개선과 docker-compose.yml 추가가 필요하다. 현재 Dockerfile은 headless 환경에서의 passphrase 처리와 Companion RPC Provider 연동을 고려하지 않으며, docker-compose가 없어 배포 및 운영이 어렵다.

## What Changes

- **Dockerfile 개선**: Health check 추가, non-root 사용자 실행, 보안 강화
- **docker-compose.yml 추가**: lango 서비스 정의, volume 마운트, 환경변수 설정
- **RPC Provider 모드 강제**: Docker 환경에서는 LocalCryptoProvider 대신 RPC Provider (Companion) 사용
- **채널 설정**: Discord, Telegram, Slack 모든 채널 지원
- **Browser Tool 유지**: Chromium 포함 (tool-browser 사용 필수)

## Capabilities

### New Capabilities
- `docker-deployment`: Docker 및 docker-compose를 통한 lango 배포 관련 요구사항

### Modified Capabilities
- `secure-signer`: Docker 환경에서 LocalCryptoProvider 비활성화 및 RPC Provider 강제 사용

## Impact

- `Dockerfile`: 개선 및 보안 강화
- `docker-compose.yml`: 신규 생성
- `internal/security`: headless 환경 감지 로직 확인 필요
- 운영 문서: Docker 배포 가이드 추가 필요
