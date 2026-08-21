# Omnexa Locale, Language & Regionalization Standard

Status: **Canonical v1**  
Work package: **P00.03**

This standard separates language, locale, timezone, country and currency so localization never corrupts business meaning.

## 1. Language/locale identifiers

Use BCP 47 language tags for UI/content locale identifiers.

Examples:

```text
en
en-US
en-GB
tr-TR
ur-PK
ar-AE
ru-RU
```

Rules:

- preserve canonical language-tag semantics rather than inventing custom locale codes;
- technical identifiers, permission names, event names and schema names remain English and are not localized;
- a translated UI label never renames the underlying canonical entity/capability;
- locale matching/fallback must be deterministic.

## 2. Locale precedence

Unless a capability deliberately defines a stricter rule, display locale precedence is:

1. explicit request/session override;
2. user preference;
3. organization preference;
4. tenant default;
5. platform default.

The platform default is `en` unless changed through an approved platform configuration decision.

A lower-precedence locale must not overwrite a higher-precedence explicit choice.

## 3. Locale is not timezone/currency/country

Never infer persisted business facts solely from locale.

These are separate concepts:

- locale/language: BCP 47 tag;
- timezone: IANA timezone ID;
- country/region: ISO 3166-1 alpha-2 where applicable;
- fiat currency: ISO 4217 alpha-3;
- measurement/tax/legal rules: explicit domain/country-pack configuration.

For example, `en-US` does not prove USD, a US legal entity or an America timezone.

## 4. Formatting

Use Unicode CLDR-compatible locale data/runtime facilities for formatting where available rather than handwritten per-country formatting logic.

Locale-aware presentation includes:

- decimal/group separators;
- currency presentation;
- dates/times;
- plural rules;
- relative time;
- names/addresses where supported;
- list formatting;
- numbering systems.

Formatting never changes the canonical persisted value.

## 5. RTL support

Omnexa UI foundations must be direction-aware from the start.

Rules:

- support both LTR and RTL document/component flow;
- direction should be derived from active locale/script or explicit content metadata, not hard-coded globally;
- use logical layout concepts (`start`/`end`) rather than assuming `left`/`right` for semantic placement;
- icons with directional meaning may require mirroring; logos/non-directional assets do not automatically mirror;
- mixed-direction user content must remain safe and readable;
- accessibility/keyboard order must follow semantic order rather than visual hacks.

## 6. Translatable business content

Localized content must use explicit locale-keyed variants or a translation/content model owned by the relevant domain.

Do not create ad-hoc schema patterns such as:

```text
name_en
name_tr
name_ar
```

for scalable multilingual platform content.

A translatable value must preserve:

- canonical field/content identity;
- locale tag;
- value;
- translation status/source where the owning domain needs it;
- fallback behavior.

## 7. User-generated content

User-generated content is not automatically translated.

Store its original language/locale metadata when known and keep machine/human translations distinguishable from the source.

AI-generated translations must remain attributable according to later AI/content governance rules.

## 8. Country/region codes

Where a country code is required, use uppercase ISO 3166-1 alpha-2 codes such as:

```text
PK
AE
TR
GB
US
```

Do not use country calling codes, currency codes or locale tags as substitutes for country identity.

Subdivisions/regions require an explicit standard/domain decision; do not invent free-text codes as global canonical identifiers.

## 9. Names, addresses and phone numbers

International person/business data must not assume one Western format.

Rules:

- do not require first-name/last-name semantics where a display/full name is the meaningful requirement;
- address structure must be country-aware and extensible;
- phone numbers should be normalized at the communications/contact boundary using an international representation while preserving user-facing formatting separately when needed;
- canonical data must support scripts beyond Latin.

Detailed contact/address schemas belong to their owning later domain.

## 10. Collation/search

Sorting/search behavior is locale-sensitive and must not be confused with identifier equality.

Rules:

- canonical IDs/codes compare using their technical contract, not locale collation;
- human text sorting/search may use locale-aware collation/search indexes;
- case folding/normalization decisions must be explicit for fields that require uniqueness.

## 11. Fallback safety

When a translation is missing:

- use the defined fallback chain;
- never silently substitute a different business value;
- preserve the source value when the owning domain allows it;
- missing translation must remain diagnosable.

## 12. Prohibited patterns

Do not:

- infer timezone or currency from locale;
- hard-code left/right layouts as semantic direction;
- create one DB column per language for platform-scale content;
- localize stable technical IDs;
- use flag icons as language identity;
- assume Gregorian display formatting is the only future requirement while conflating display calendar with stored business date;
- assume all names/addresses fit one country-specific structure.
