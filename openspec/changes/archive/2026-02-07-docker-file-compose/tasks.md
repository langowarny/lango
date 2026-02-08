## 1. Dockerfile 개선

- [x] 1.1 Non-root user 추가 (useradd, USER 지시문)
- [x] 1.2 Health check 추가 (curl 또는 wget으로 /health 엔드포인트 호출)
- [x] 1.3 데이터 디렉토리 권한 설정 (chown)

## 2. docker-compose.yml 생성

- [x] 2.1 lango 서비스 정의 (build, ports, volumes)
- [x] 2.2 환경변수 설정 (ANTHROPIC_API_KEY, DISCORD_BOT_TOKEN, TELEGRAM_BOT_TOKEN, SLACK_BOT_TOKEN, SLACK_APP_TOKEN)
- [x] 2.3 volume 정의 (lango-data)
- [x] 2.4 lango.json 마운트 설정 (read-only)

## 3. Docker 환경 감지 로직

- [x] 3.1 Docker 환경 감지 함수 추가 (/.dockerenv 또는 /proc/1/cgroup 확인)
- [x] 3.2 LocalCryptoProvider 초기화 시 Docker 환경 체크
- [x] 3.3 Docker 환경에서 LocalCryptoProvider 사용 시 에러 메시지 출력

## 4. 검증

- [x] 4.1 Docker 이미지 빌드 테스트
- [x] 4.2 docker-compose up 실행 확인
- [x] 4.3 Health check 작동 확인
- [x] 4.4 Docker 환경에서 RPC Provider 필수 에러 메시지 확인
