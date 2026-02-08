## Context

현재 `LocalCryptoProvider`는 config 파일의 `security.passphrase`에서 passphrase를 읽는다. AI Agent는 filesystem tool로 config 파일에 접근 가능하므로, passphrase가 노출될 위험이 있다. 이는 Zero Trust 보안 원칙에 위배된다.

또한 passphrase가 변경되거나 잘못 입력되면 기존 암호화된 데이터에 접근할 수 없게 되는 문제가 있다.

## Goals / Non-Goals

**Goals:**
- Config 파일에서 passphrase 읽는 기능 제거
- Interactive terminal prompt로 passphrase 입력받기
- Passphrase checksum으로 잘못된 입력 조기 감지
- Passphrase 변경 시 기존 데이터 마이그레이션 지원
- 보안 모드(Local vs RPC) 문서화

**Non-Goals:**
- RPCProvider/Companion 로직 변경
- 새로운 암호화 알고리즘 도입
- GUI 기반 passphrase 입력

## Decisions

### 1. Interactive Prompt 구현

**선택**: `golang.org/x/term` 패키지 사용
```go
func promptPassphrase() (string, error) {
    fmt.Print("Enter passphrase: ")
    bytes, err := term.ReadPassword(int(syscall.Stdin))
    return string(bytes), err
}
```

**대안 검토**:
- `bufio.Scanner`: 입력이 터미널에 그대로 출력됨 → 보안 취약
- Third-party library: 불필요한 의존성 추가

### 2. Checksum 저장 방식

**선택**: Salt와 함께 passphrase hash 저장
```
security_config 테이블:
- key: "default" 
- salt: <random bytes>
- checksum: SHA256(passphrase + salt)  ← 새로 추가
```

**동작**:
1. 첫 설정 시: passphrase 입력 → salt 생성 → checksum 저장
2. 이후 시작 시: passphrase 입력 → checksum 검증 → 불일치 시 에러

### 3. Migration 프로세스

**선택**: CLI 명령어로 제공
```bash
lango security migrate-passphrase
```

**동작**:
```
┌─────────────────────────────────────────┐
│ 1. 현재 passphrase 입력 (검증)          │
│ 2. 새 passphrase 입력 (2회 확인)        │
│ 3. 모든 Secret 조회                     │
│ 4. 각 Secret: old key로 복호화 →        │
│              new key로 재암호화         │
│ 5. Salt/Checksum 업데이트               │
│ 6. 완료 메시지                          │
└─────────────────────────────────────────┘
```

### 4. Config Passphrase 제거

**선택**: Deprecated 처리 후 경고
- 기존 config에 passphrase 있으면 시작 시 경고 출력
- 실제 사용하지 않음 (무시)
- 다음 major 버전에서 필드 제거

## Risks / Trade-offs

| Risk | Mitigation |
|------|------------|
| Headless 환경에서 시작 불가 | RPCProvider+Companion 사용 권장 |
| Passphrase 분실 시 데이터 손실 | 복구 방법 없음 명시, 백업 권장 문서화 |
| Migration 중 실패 시 데이터 손상 | 트랜잭션 처리 + 롤백 지원 |
| Checksum 유출 시 brute-force 가능 | Salt + 강력한 hash 알고리즘 사용 |

## Open Questions

1. Passphrase 최소 길이/복잡도 요구사항 필요한가? (Answer: required, 12 characters minimum, at least one uppercase letter, one lowercase letter, one number, and one special character)
2. Migration 실패 시 롤백을 자동으로 할 것인가, 수동 복구 가이드를 제공할 것인가? (Answer: automatically roll back all changes)
