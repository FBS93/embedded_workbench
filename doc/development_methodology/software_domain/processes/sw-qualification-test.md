# sw-qualification-test

## Purpose

Define and implement software qualification tests against software requirements to validate that the software architecture, software detailed design, and software implementation correctly realizes them.

## Input work products

- Software requirements
- Software architecture
- Software detailed design
- Software implementation

## Output work products

- Software qualification test

## Steps

1. Review the software requirements.
2. Identify software requirements requiring validation.
3. Define and implement the software qualification tests.
4. Define traceability between software qualification tests and the corresponding software requirements.
5. Check the software qualification tests for completeness, consistency, and correctness.

## Guidelines

### Software qualification test work product

The software qualification test work product shall follow the [Software test definition](../resources/software_test_definition.md).

In addition to the [Software test definition](../resources/software_test_definition.md), the software qualification tests shall:
- Be derived from the software requirements to validate that the software architecture, software detailed design, and software implementation correctly realizes them. The software architecture and software detailed design may be used as complementary input to understand how software interacts and to support test implementation, but it shall not be the validation target. The software implementation shall be used only as the executable element under test.
- Define at least one qualification test for each software requirement, ensuring full coverage of all software requirements.
- Have a unique and well-defined objective. Multiple tests may be defined when required to validate different aspects of the same software requirement.
- Be implemented in C using [Embedded Test Framework (ETF)](xxxyyy-ETF) or [EDF Test Framework (EDFTest)](xxxyyy-EDFTest), depending on the implemented architecture model.
  - Event-driven architecture model implemented using [Event Driven Framework (EDF)](xxxyyy-EDF) shall be tested using [EDF Test Framework (EDFTest)](xxxyyy-EDFTest).
    - The [Event Driven Framework (EDF)](xxxyyy-EDF) software tests shall be used as a reference for how to use the [EDF Test Framework (EDFTest)](xxxyyy-EDFTest).
  - All non event-driven architecture models shall be tested using [Embedded Test Framework (ETF)](xxxyyy-ETF).
    - The [Embedded Middleware Framework (EMF)](xxxyyy-EMF) software tests shall be used as a reference for how to use the [Embedded Test Framework (ETF)](xxxyyy-ETF).
- Validate the software behavior using the complete software executable from an external perspective, without relying on internal implementation details.

At the end of the software qualification test process, it shall be ensured that all software requirements are validated by at least one software qualification test and that the complete software behaves as specified.