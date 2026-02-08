## Why

현재 `LocalCryptoProvider`는 config 파일에서 passphrase를 읽는데, AI Agent가 filesystem tool로 config를 읽으면 passphrase가 노출될 수 있다. 이는 "AI는 암호화 키에 접근할 수 없어야 한다"는 보안 목표와 불일치한다.

## What Changes

- **Interactive Passphrase Prompt**: Config에서 passphrase 읽는 대신 터미널에서 직접 입력받음
- **Passphrase Checksum**: 저장된 salt와 함께 checksum을 저장하여 잘못된 passphrase 입력 시 조기 감지
- **Migration Process**: Passphrase 변경 시 기존 암호화된 데이터를 새 키로 재암호화하는 프로세스
- **Security Mode Documentation**: LocalCryptoProvider (dev/test)와 RPCProvider+Companion (production) 모드 명확히 문서화
- **Doctor Warnings**: Local provider 사용 시 개발/테스트 전용 경고 표시

## Capabilities

### New Capabilities
- `passphrase-management`: Interactive passphrase prompt, checksum validation, migration workflow

### Modified Capabilities
- `secure-signer`: LocalCryptoProvider 초기화 로직 변경 (config → interactive prompt)
- `cli-doctor`: Security provider 모드에 따른 경고/권장사항 추가

## Impact

- `internal/app/app.go`: LocalCryptoProvider 초기화 로직 수정
- `internal/config/types.go`: Passphrase 필드 deprecated 처리
- `internal/session/ent_store.go`: Checksum 저장/검증 메서드 추가
- `internal/security/local_provider.go`: Migration 로직 추가
- `internal/cli/doctor/checks/security.go`: Provider 모드 체크 추가
- `README.md`: Security 섹션에 두 모드 설명 추가
