# System test definition

## Purpose

Define the common test definition across the system domain.

## Test structure

The system tests shall:
- Be documented with a Markdown document named `<sys_test_name>.md`, following a consistent and structured format, with the following information:
  - A unique identifier defined as a top Markdown heading, following the pattern `SYS_<TEST_TYPE>_TEST_<X>`.
    - `<TEST_TYPE>` identifies the test type in upper case. Possible values are `INTEGRATION` and `QUALIFICATION`.
    - `<X>` is a monotonically increasing number that shall never be reused.
  - An overview chapter using the following pattern `<System test name> overview`. The chapter shall contain:
    - A description of the purpose of the test.
    - Upstream traceability to the system traceable elements being tested. This shall be defined as a list of Markdown links referencing the corresponding system traceable elements.
- Be implemented in Python using [pytest](https://docs.pytest.org/).
- Follow the rules defined in [python_guidelines.md](xxxyyy-todo).
- Be designed to stimulate the system externally and observe its behavior in a fully automated way whenever possible. When full automation is not feasible, tests shall define a structured sequence of manual steps executed through the available interfaces (e.g., terminal), where each step requires user confirmation (OK/NOK) for the observed behavior. Automated and manual-assisted tests shall be defined separately to allow independent execution.

The following template shall be used:

```md
# SYS_<TEST_TYPE>_TEST_<X>

## <System test name> overview

<Description of the purpose of the test>.

Upstream traceability:
- [<UPSTREAM_ELEMENT_UNIQUE_ID>](#...)
```