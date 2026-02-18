## Why

Docker에서 텔레그램 메시지를 수신하면 Gemini API가 `"GenerateContentRequest.contents: contents is not specified"` 에러를 반환한다. ADK runner가 세션에서 이벤트를 읽을 때 방금 추가한 사용자 메시지가 인메모리 히스토리에 없어서 빈 contents가 Gemini API로 전송된다. 또한 시스템 프롬프트가 provider에 전달되지 않아 에이전트 성격/지시사항이 무시된다.

## What Changes

- **AppendEvent 인메모리 히스토리 동기화**: `SessionServiceAdapter.AppendEvent`가 DB 저장 후 `SessionAdapter.sess.History`도 업데이트하여 ADK runner가 현재 턴의 사용자 메시지를 읽을 수 있도록 함
- **SystemInstruction 전달**: `ModelAdapter.GenerateContent`가 `req.Config.SystemInstruction`을 system 메시지로 변환하여 provider에 전달
- 관련 테스트 추가

## Capabilities

### New Capabilities

### Modified Capabilities
- `adk-architecture`: AppendEvent에서 인메모리 히스토리 동기화 추가, ModelAdapter에서 SystemInstruction 전달 추가

## Impact

- `internal/adk/session_service.go`: AppendEvent 메서드 수정
- `internal/adk/model.go`: GenerateContent 메서드 수정, extractSystemText 헬퍼 추가
- `internal/adk/session_service_test.go`: 신규 테스트 파일
- `internal/adk/model_test.go`: SystemInstruction 테스트 추가
- `internal/adk/state_test.go`: mockStore 수정 (DB-only 동작 시뮬레이션)
