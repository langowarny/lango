## Context

ADK v0.4.0의 메시지 처리 흐름에서 두 가지 버그가 발견되었다:

1. `SessionServiceAdapter.AppendEvent`가 `store.AppendMessage()`로 DB에만 메시지를 저장하고, `SessionAdapter.sess.History` (인메모리)를 업데이트하지 않음. ADK runner는 같은 세션 객체의 인메모리 히스토리를 읽으므로, 방금 추가한 사용자 메시지가 보이지 않아 빈 contents가 Gemini API로 전송됨.

2. `ModelAdapter.GenerateContent`가 `req.Config.SystemInstruction`을 무시하고 `req.Contents`만 provider에 전달. ADK가 설정한 시스템 프롬프트가 LLM에 도달하지 않음.

## Goals / Non-Goals

**Goals:**
- AppendEvent 호출 후 인메모리 히스토리가 즉시 업데이트되어 ADK runner가 현재 턴의 메시지를 읽을 수 있도록 함
- ADK SystemInstruction이 provider에 system 메시지로 전달되도록 함
- 두 수정 모두 테스트로 검증

**Non-Goals:**
- ADK 라이브러리 자체 수정
- 세션 저장소 인터페이스 변경
- Provider 인터페이스 변경

## Decisions

### Decision 1: AppendEvent에서 인메모리 히스토리 직접 업데이트

DB 저장 성공 후 `SessionAdapter.sess.History`에 메시지를 append한다. `sess` 파라미터를 `*SessionAdapter`로 type assertion하여 접근한다.

**대안 고려:** DB에서 다시 읽어오기 → 불필요한 I/O, 성능 저하. 직접 append가 단순하고 효율적.

### Decision 2: SystemInstruction을 system role 메시지로 변환

`genai.Content`의 text parts를 결합하여 단일 `provider.Message{Role: "system"}`로 변환하고, messages 배열 앞에 prepend한다.

**대안 고려:** GenerateParams에 SystemInstruction 필드 추가 → provider 인터페이스 변경 필요, 모든 provider 구현체 수정 필요. 기존 Messages 배열에 system role로 넣는 것이 변경 범위가 최소.

## Risks / Trade-offs

- **[인메모리/DB 불일치]** → AppendEvent 실패 시 DB에는 저장되었지만 인메모리에는 없을 수 있음. 하지만 에러를 반환하므로 호출자가 처리할 수 있고, DB 저장 성공 후에만 인메모리 업데이트하므로 일관성 보장.
- **[System message 중복]** → 일부 provider가 이미 system message를 별도로 처리할 수 있음. 현재 provider 구현을 확인한 결과 Messages 배열의 system role을 정상 처리함.
