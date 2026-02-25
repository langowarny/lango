## 1. Leaf Dependencies (wallet, payment)

- [x] 1.1 Add ChainID type and constants (ChainEthereumMainnet, ChainBase, ChainBaseSepolia, ChainSepolia) to wallet/wallet.go
- [x] 1.2 Add CurrencyUSDC constant to wallet/wallet.go
- [x] 1.3 Update NetworkName() to use ChainID constants
- [x] 1.4 Export WalletKeyName in wallet/create.go and update all internal references
- [x] 1.5 Add gas fee constants (DefaultBaseFeeWei, DefaultMaxPriorityFeeWei, BaseFeeMultiplier, EthAddressLength) and BalanceOfSelector to payment/tx_builder.go
- [x] 1.6 Replace magic numbers in BuildTransferTx and ValidateAddress with constants
- [x] 1.7 Add DefaultHistoryLimit and purposeX402AutoPayment to payment/service.go
- [x] 1.8 Replace BalanceOfSelector inline definition with package-level var in service.go
- [x] 1.9 Update x402/signer.go to use wallet.WalletKeyName

## 2. P2P Protocol

- [x] 2.1 Add ResponseStatus enum type with Valid() to protocol/messages.go
- [x] 2.2 Add 7 sentinel errors (ErrMissingToolName, ErrAgentCardUnavailable, etc.) to protocol/messages.go
- [x] 2.3 Change Response.Status field from string to ResponseStatus
- [x] 2.4 Replace all hardcoded status strings in handler.go (~17 occurrences)
- [x] 2.5 Add payGateStatus local constants in handler.go
- [x] 2.6 Replace error message strings with sentinel .Error() in handler.go
- [x] 2.7 Replace status comparisons in remote_agent.go with ResponseStatusOK
- [x] 2.8 Add errMsgUnknown constant in remote_agent.go
- [x] 2.9 Update handler_test.go to use typed constants

## 3. Firewall

- [x] 3.1 Add ACLAction enum type with Valid() to firewall/firewall.go
- [x] 3.2 Add WildcardAll constant to firewall/firewall.go
- [x] 3.3 Add 4 sentinel errors (ErrRateLimitExceeded, etc.) to firewall/firewall.go
- [x] 3.4 Change ACLRule.Action from string to ACLAction
- [x] 3.5 Replace all "allow"/"deny"/"*" literals in firewall.go
- [x] 3.6 Update firewall_test.go to use typed constants
- [x] 3.7 Update downstream callers (cli/p2p/firewall.go, app/tools.go, app/wiring.go) with ACLAction casts

## 4. ZKP

- [x] 4.1 Add ProofScheme enum type with Valid() to zkp/zkp.go
- [x] 4.2 Add SRSMode enum type with Valid() to zkp/zkp.go
- [x] 4.3 Add ErrUnsupportedScheme sentinel error
- [x] 4.4 Change Config.Scheme, ProverService.scheme, Proof.Scheme to ProofScheme type
- [x] 4.5 Change Config.SRSMode, ProverService.srsMode to SRSMode type
- [x] 4.6 Replace all "plonk"/"groth16"/"unsafe"/"file" literals in switch cases
- [x] 4.7 Update app/wiring.go with ProofScheme/SRSMode type casts
- [x] 4.8 Update zkp_test.go to use typed constants

## 5. Reputation

- [x] 5.1 Add FailureWeight, TimeoutWeight, BasePenalty constants to reputation/store.go
- [x] 5.2 Update CalculateScore() to use named constants

## 6. KMS Factory

- [x] 6.1 Add KMSProviderName enum type with Valid() to security/kms_factory.go
- [x] 6.2 Change NewKMSProvider parameter to KMSProviderName type
- [x] 6.3 Update callsites in app/wiring.go and cli/security/kms.go with type casts

## 7. A2A Server

- [x] 7.1 Add AgentCardRoute, ContentTypeJSON, SkillTagOrchestration, SkillTagSubAgentPrefix constants to a2a/server.go
- [x] 7.2 Replace hardcoded strings in server.go and server_test.go

## 8. Paygate

- [x] 8.1 Add DefaultQuoteExpiry constant to paygate/gate.go
- [x] 8.2 Replace "USDC" with wallet.CurrencyUSDC in gate.go and gate_test.go
- [x] 8.3 Replace 5*time.Minute with DefaultQuoteExpiry in BuildQuote

## 9. Cross-cutting USDC

- [x] 9.1 Replace "USDC" with wallet.CurrencyUSDC in app/wiring.go
- [x] 9.2 Replace "USDC" with wallet.CurrencyUSDC in app/p2p_routes.go
- [x] 9.3 Replace "USDC" with wallet.CurrencyUSDC in app/tools.go
- [x] 9.4 Replace "USDC" with wallet.CurrencyUSDC in tools/payment/payment.go
- [x] 9.5 Replace "USDC" with wallet.CurrencyUSDC in cli/payment/balance.go
- [x] 9.6 Replace "USDC" with wallet.CurrencyUSDC in cli/payment/limits.go
- [x] 9.7 Replace "USDC" with wallet.CurrencyUSDC in cli/p2p/pricing.go

## 10. CLI Display Constants

- [x] 10.1 Add maxHashDisplay, truncatedHashLen, maxPurposeDisplay, truncatedPurpLen constants to cli/payment/history.go
- [x] 10.2 Replace magic numbers with named constants
