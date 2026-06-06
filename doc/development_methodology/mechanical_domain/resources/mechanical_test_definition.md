# Mechanical test definition

## Purpose

Define the common test definition across the mechanical domain.

## Test structure

Each mechanical test shall:
- Be documented with a Markdown document named `<mech_test_name>.md`, following the template defined in [python_coding_guidelines.md](../../resources/python_coding_guidelines.md#test-documentation), with the following additional information:
  - A unique identifier defined as a top Markdown heading (all other headings defined in the [python_coding_guidelines.md](../../resources/python_coding_guidelines.md#test-documentation) template shall be shifted one level below this top heading), following the pattern `MECH_<TEST_TYPE>_TEST_<X>`.
    - `<TEST_TYPE>` identifies the test type in upper case. Possible values are `INTEGRATION` for integration test and `QUALIFICATION` for qualification test.
    - `<X>` is a monotonically increasing number that shall never be reused.
  - The `<Test name> test overview` chapter placeholder shall be replaced with the mechanical test name used in the document name, using the following pattern `<Mechanical test name> overview`. The chapter shall contain the following additional information:
    - Upstream traceability to the mechanical traceable elements being tested. This shall be defined as a list of Markdown links referencing the corresponding mechanical traceable elements.
  - A status indicating the approval state of the test. Shall use exactly one of these values: `not approved`, `approved`.
- Be implemented in Python using [pytest](https://docs.pytest.org/).
- Follow the rules defined in [cadquery_guidelines.md](cadquery_guidelines.md).
- Be designed to validate in a fully automated way whenever possible. When full automation is not feasible, tests shall define a structured sequence of manual steps, where each step requires human confirmation (`OK`/`NOK`) for the observed result.
- Define automated and manual-assisted tests separately to allow independent execution.

The following additions to the Markdown document template shall be applied for mechanical tests, in addition to the template defined in [python_coding_guidelines.md](../../resources/python_coding_guidelines.md#test-documentation):

```md
# MECH_<TEST_TYPE>_TEST_<X>

## <Mechanical test name> overview

<Description of the purpose of the test>.

Upstream traceability:
- [<UPSTREAM_ELEMENT_UNIQUE_ID>](#...)

Status: <status>
```
