# LinkBio Project Pack (for AI Coding)

This folder contains the missing spec files to make AI implement the project correctly:

- `docs/api_contract.md`: endpoints + payload schemas + error format + caching headers
- `docs/acceptance_tests.md`: manual acceptance tests (no-code friendly)
- `docs/payload_examples/`: example JSONs for draft/compiled/theme/patch

How to use:
1) Put these files into your repo under `docs/`.
2) When prompting AI coding, always attach or reference:
   - requirements
   - theme spec
   - db schema + migration
   - repo structure
   - api contract + payload examples + acceptance tests
3) Implement by vertical slices and validate against acceptance tests.
