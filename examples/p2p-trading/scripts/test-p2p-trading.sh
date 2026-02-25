#!/bin/sh
set -e

# Colors for test output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

ALICE="http://localhost:18789"
BOB="http://localhost:18790"
CHARLIE="http://localhost:18791"
RPC="http://localhost:8545"

PASSED=0
FAILED=0

pass() {
  PASSED=$((PASSED + 1))
  printf "${GREEN}  PASS${NC}: %s\n" "$1"
}

fail() {
  FAILED=$((FAILED + 1))
  printf "${RED}  FAIL${NC}: %s\n" "$1"
}

section() {
  printf "\n${YELLOW}── %s ──${NC}\n" "$1"
}

# ─────────────────────────────────────────────
section "1. Health Checks"
# ─────────────────────────────────────────────
for NAME_URL in "Alice:$ALICE" "Bob:$BOB" "Charlie:$CHARLIE"; do
  NAME="${NAME_URL%%:*}"
  URL="${NAME_URL#*:}"
  if curl -sf "$URL/health" | grep -q '"status":"ok"'; then
    pass "$NAME health"
  else
    fail "$NAME health"
  fi
done

# ─────────────────────────────────────────────
section "2. P2P Status"
# ─────────────────────────────────────────────
for NAME_URL in "Alice:$ALICE" "Bob:$BOB" "Charlie:$CHARLIE"; do
  NAME="${NAME_URL%%:*}"
  URL="${NAME_URL#*:}"
  STATUS=$(curl -sf "$URL/api/p2p/status")
  if echo "$STATUS" | grep -q '"peerId"'; then
    pass "$NAME P2P status (has peerId)"
  else
    fail "$NAME P2P status"
  fi
done

# ─────────────────────────────────────────────
section "3. P2P Discovery (waiting 15s for mDNS)"
# ─────────────────────────────────────────────
sleep 15
for NAME_URL in "Alice:$ALICE" "Bob:$BOB" "Charlie:$CHARLIE"; do
  NAME="${NAME_URL%%:*}"
  URL="${NAME_URL#*:}"
  PEERS=$(curl -sf "$URL/api/p2p/peers")
  COUNT=$(echo "$PEERS" | grep -o '"count":[0-9]*' | grep -o '[0-9]*')
  if [ -n "$COUNT" ] && [ "$COUNT" -ge 2 ]; then
    pass "$NAME discovered $COUNT peers"
  else
    fail "$NAME peer discovery (count: ${COUNT:-0}, expected >= 2)"
  fi
done

# ─────────────────────────────────────────────
section "4. P2P Identity (DID)"
# ─────────────────────────────────────────────
for NAME_URL in "Alice:$ALICE" "Bob:$BOB" "Charlie:$CHARLIE"; do
  NAME="${NAME_URL%%:*}"
  URL="${NAME_URL#*:}"
  IDENTITY=$(curl -sf "$URL/api/p2p/identity")
  if echo "$IDENTITY" | grep -q '"did":"did:lango:'; then
    pass "$NAME DID starts with did:lango:"
  else
    fail "$NAME DID check ($IDENTITY)"
  fi
done

# ─────────────────────────────────────────────
section "5. USDC Balances (on-chain)"
# ─────────────────────────────────────────────
# Read USDC address from inside a running container (Docker volume not accessible from host)
USDC_ADDRESS=$(docker compose exec -T alice cat /shared/usdc-address.txt 2>/dev/null | tr -d '[:space:]')
if [ -z "$USDC_ADDRESS" ]; then
  USDC_ADDRESS=$(docker compose exec -T bob cat /shared/usdc-address.txt 2>/dev/null | tr -d '[:space:]')
fi

ALICE_ADDR="0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266"
BOB_ADDR="0x70997970C51812dc3A010C7d01b50e0d17dc79C8"
CHARLIE_ADDR="0x3C44CdDdB6a900fa2b585dd299e03d12FA4293BC"

if [ -n "$USDC_ADDRESS" ]; then
  echo "  USDC contract: $USDC_ADDRESS"
  for NAME_ADDR in "Alice:$ALICE_ADDR" "Bob:$BOB_ADDR" "Charlie:$CHARLIE_ADDR"; do
    NAME="${NAME_ADDR%%:*}"
    ADDR="${NAME_ADDR#*:}"
    # Run cast inside the anvil container (host may not have Foundry installed)
    BAL=$(docker compose exec -T anvil cast call "$USDC_ADDRESS" "balanceOf(address)(uint256)" "$ADDR" --rpc-url "http://localhost:8545" 2>/dev/null | tr -d '[:space:]')
    # 1000 USDC = 1000000000 (6 decimals)
    if echo "$BAL" | grep -q "1000000000"; then
      pass "$NAME USDC balance = 1000.00"
    else
      fail "$NAME USDC balance (got: $BAL, expected: 1000000000)"
    fi
  done
else
  fail "Could not read USDC contract address"
fi

# ─────────────────────────────────────────────
section "6. USDC Transfer (Alice → Bob, 1.00 USDC via on-chain)"
# ─────────────────────────────────────────────
if [ -n "$USDC_ADDRESS" ]; then
  # Transfer 1.00 USDC (1000000) directly on-chain using Alice's private key.
  # Note: lango CLI `payment send` requires keyfile for bootstrap, which gets
  # shredded after serve starts. Using cast for deterministic E2E testing.
  ALICE_KEY="0xac0974bec39a17e36ba4a6b4d238ff944bacb478cbed5efcae784d7bf4f2ff80"
  TRANSFER_AMOUNT="1000000"  # 1.00 USDC (6 decimals)

  docker compose exec -T anvil cast send "$USDC_ADDRESS" \
    "transfer(address,uint256)(bool)" "$BOB_ADDR" "$TRANSFER_AMOUNT" \
    --rpc-url "http://localhost:8545" \
    --private-key "$ALICE_KEY" >/dev/null 2>&1 && \
    pass "Alice transferred 1.00 USDC to Bob (on-chain)" || \
    fail "Alice USDC transfer to Bob"

  # Verify Bob's balance increased
  sleep 2  # wait for tx confirmation
  BOB_BAL=$(docker compose exec -T anvil cast call "$USDC_ADDRESS" "balanceOf(address)(uint256)" "$BOB_ADDR" --rpc-url "http://localhost:8545" 2>/dev/null | tr -d '[:space:]')
  if echo "$BOB_BAL" | grep -q "1001000000"; then
    pass "Bob balance = 1001.00 USDC (received 1.00)"
  else
    fail "Bob balance after transfer (got: $BOB_BAL, expected: 1001000000)"
  fi

  # Verify Alice's balance decreased
  ALICE_BAL=$(docker compose exec -T anvil cast call "$USDC_ADDRESS" "balanceOf(address)(uint256)" "$ALICE_ADDR" --rpc-url "http://localhost:8545" 2>/dev/null | tr -d '[:space:]')
  if echo "$ALICE_BAL" | grep -q "999000000"; then
    pass "Alice balance = 999.00 USDC (sent 1.00)"
  else
    fail "Alice balance after transfer (got: $ALICE_BAL, expected: 999000000)"
  fi
else
  fail "Skipping transfer test — USDC address unknown"
fi

# ─────────────────────────────────────────────
section "Results"
# ─────────────────────────────────────────────
TOTAL=$((PASSED + FAILED))
printf "\n${GREEN}Passed${NC}: %d / %d\n" "$PASSED" "$TOTAL"
if [ "$FAILED" -gt 0 ]; then
  printf "${RED}Failed${NC}: %d / %d\n" "$FAILED" "$TOTAL"
  exit 1
fi

printf "\n${GREEN}All tests passed!${NC}\n"
