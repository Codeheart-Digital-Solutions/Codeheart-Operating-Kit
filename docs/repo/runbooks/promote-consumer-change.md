Last updated: 2026-06-13T22:47:44Z (UTC)

# Promote Consumer Change

Use this runbook when a useful rule, runbook, template, or workflow is discovered in a consumer
repository and may belong in Codeheart Operating Kit.

## Procedure

1. Identify the consumer-local source and its current owner.
2. Confirm the material is reusable beyond one consumer.
3. Remove customer, tenant, product-specific, environment-specific, secret, and business-private
   details.
4. Decide whether the promoted material is managed content, scaffold, template, reference, runbook,
   schema, validator, or CLI behavior.
5. Check `docs/repo/reference/placement-contract.md` for the target path.
6. Classify consumer impact with `docs/repo/reference/consumer-impact-classification.md`.
7. Add tests, fixtures, schemas, or validation when the promoted material changes behavior.
8. Leave consumer-specific exceptions in the consumer repository.
9. Add release notes or migration notes when installed consumers must act.
10. Run validation for the promoted surface.

## Stop Conditions

Stop before promotion when the material is only useful for one consumer, contains private content
that cannot be sanitized, conflicts with the placement contract, or would change consumer-owned
state without a migration strategy.
