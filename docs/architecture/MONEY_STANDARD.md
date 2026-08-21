# Omnexa Money & Numeric Precision Standard

Status: **Canonical v1**  
Work package: **P00.03**

This standard prevents financial corruption caused by binary floating point, silent currency conversion, inconsistent rounding or ambiguous precision.

## 1. Canonical money model

A monetary value is a pair:

```text
Money {
  amount
  currency
}
```

Rules:

- `amount` is an exact decimal value, never IEEE-754 binary floating point.
- `currency` is required whenever a value represents money.
- For fiat currencies, use the uppercase ISO 4217 alpha-3 code where one exists.
- Currency must never be inferred from locale, country, tenant or user preference when persisting a monetary value.
- Two amounts in different currencies must never be added, compared as money or posted together without an explicit conversion operation.

## 2. API/contract representation

At JSON boundaries, monetary decimal values are serialized as **strings**, not JSON numbers.

Example:

```json
{
  "amount": "1250.375000",
  "currency": "AED"
}
```

This preserves exact cross-language semantics and avoids JavaScript/binary floating-point loss.

P00.04 will define the final reusable schema/envelope for API contracts.

## 3. Persistence precision

The default PostgreSQL storage type for general high-precision monetary calculations is:

```text
NUMERIC(38,18)
```

This is a platform baseline, not permission for every domain to expose 18 fractional digits to users.

Rules:

- ledgers, invoices, settlements and statutory documents quantize according to their owning accounting/country policy at an explicit boundary;
- unit prices, allocations, exchange rates and tax calculations may retain greater working precision before final quantization;
- a domain requiring precision beyond the platform baseline must document the requirement and use an ADR if it changes a stable contract/storage convention;
- schema must never use `float`, `double precision` or equivalent binary floating point for monetary values.

## 4. Currency exponent and display

Display/payment exponent is not the same thing as internal calculation precision.

Examples:

- a currency may display 0, 2 or 3 fractional digits;
- a tax allocation may require more internal precision than the final invoice line;
- a payment provider may require integer minor units even though Omnexa stores decimal money canonically.

Provider-specific minor-unit conversion belongs in the payment connector/boundary and must be checked for exactness before submission.

## 5. Rounding policy

Rounding is always explicit.

Canonical generic arithmetic default: **round half to even** when a domain has not supplied a legally/business-required rule.

However:

- tax, payroll, accounting, cash rounding and statutory rules may require a different policy;
- a country pack/domain may override the default only through an explicit named policy;
- the rounding boundary, scale and mode must be deterministic and testable;
- repeated intermediate rounding should be avoided unless the governing rule requires it;
- totals must be reconciled through an explicit allocation strategy rather than hidden penny adjustments.

Named policies should use stable identifiers such as:

```text
finance.rounding.half_even
finance.rounding.half_up
payments.cash_rounding.<policy>
```

## 6. Rates, percentages and ratios

Percentages, exchange rates, tax rates and ratios are exact decimal quantities but are **not Money**.

Rules:

- never attach a currency to a pure rate;
- never store `15%` as ambiguous `15` or `0.15` without a documented representation;
- canonical fractional representation is decimal fraction, so 15% is `0.15`;
- APIs serialize exact decimal rates as strings;
- rate precision must be explicit and must not silently inherit currency display precision.

## 7. Currency conversion

Currency conversion is an explicit business operation.

A conversion record/result must be able to identify:

- source currency and source amount;
- target currency and target amount;
- exact rate used;
- rate direction/base/quote semantics;
- effective timestamp/date;
- rate source/provider or governing rate table;
- rounding policy/boundary;
- tenant/legal-entity context when relevant.

Historical financial records must not be recomputed using a newer exchange rate merely because the current rate changed.

## 8. Digital assets and non-ISO units

The core fiat model uses ISO 4217 codes. Digital assets, loyalty points, commodities or custom units must not masquerade as ISO currency codes.

They require a separately namespaced asset/unit registry and explicit domain contract. A later payments/asset module may define such representations without weakening fiat Money semantics.

## 9. Negative values and sign semantics

Whether a negative amount is valid is domain-specific.

Examples:

- credits/adjustments may allow negative values;
- a payment authorization request may forbid them;
- accounting uses debit/credit semantics and must not reduce the ledger to sign-only meaning.

The Money primitive preserves sign but does not decide business validity.

## 10. Auditability

Material financial calculations must make it possible to reconstruct:

- original inputs;
- precision used;
- rounding policy;
- conversion rate/source if any;
- resulting posted/settled amount.

Do not rely only on a final rounded number when regulation or reconciliation requires provenance.

## 11. Prohibited patterns

Do not:

- use binary floating point for money;
- persist money without currency;
- infer currency from locale;
- silently convert currencies;
- round at arbitrary UI/service boundaries;
- assume every currency has two decimals;
- treat provider minor units as the universal Omnexa money representation;
- use formatted display strings (`"$1,250.00"`) as persisted monetary values.
