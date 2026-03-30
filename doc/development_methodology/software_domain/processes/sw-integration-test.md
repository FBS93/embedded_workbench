# sw-integration-test

## Purpose

Define and implement software integration tests against software architecture to validate that the software detailed design and software implementation correctly realizes it.

## Input work products

- Software architecture
- Software detailed design
- Software implementation

## Output work products

- Software integration test

## Steps

1. Review the software architecture.
2. Identify software architecture elements requiring validation.
3. Define and implement the software integration tests.
4. Define traceability between software integration tests and the corresponding software architecture elements.
5. Check the software integration tests for completeness, consistency, and correctness.

## Guidelines

### Software integration test work product

The software integration test work product shall follow the [Software test definition](../resources/software_test_definition.md).

In addition to the [Software test definition](../resources/software_test_definition.md), the software integration tests shall:
- Be derived from the software architecture to validate that the software detailed design and software implementation correctly realizes it. The software detailed design may be used as complementary input to understand how software interacts and to support test implementation, but it shall not be the validation target. The software implementation shall be used only as the executable element under test.
- Define at least one integration test for each dynamic view represented as a sequence diagram, ensuring that all interactions defined through interfaces between software components are validated.
- Have a unique and well-defined objective. Multiple tests may be defined when required to validate the same interaction scenario.
- Be implemented in C using [Embedded Test Framework (ETF)](xxxyyy-ETF) or [EDF Test Framework (EDFTest)](xxxyyy-EDFTest), depending on the interaction model:
  - Event-driven interactions implemented using [Event Driven Framework (EDF)](xxxyyy-EDF) shall be tested using [EDF Test Framework (EDFTest)](xxxyyy-EDFTest).
    - The [Event Driven Framework (EDF)](xxxyyy-EDF) software tests shall be used as a reference for how to use the [EDF Test Framework (EDFTest)](xxxyyy-EDFTest).
  - All other interactions shall be tested using [Embedded Test Framework (ETF)](xxxyyy-ETF).
    - The [Embedded Middleware Framework (EMF)](xxxyyy-EMF) software tests shall be used as a reference for how to use the [Embedded Test Framework (ETF)](xxxyyy-ETF).

At the end of the software integration test process, the complete software implementation shall build as a single executable without errors or warnings, using the toolchain and build configuration defined in the platform software design. This validation shall be performed for all defined build presets, including both host and target. However, individual software integration tests may be implemented using only the subset of software components required to validate the defined interactions and may not require a complete software executable.

#### Active object integration test

For software components implemented as active objects, software integration tests shall also validate their behavior as implemented by the Hierarchical State Machine (HSM), including all event-driven interactions defined in the software architecture.

The integration tests shall ensure that:
- All event-driven interactions between active objects behave as defined.
- The complete Hierarchical State Machine (HSM) behavior is validated through event sequences.
- All relevant use cases are covered.

#### Platform integration test

The platform software design shall be validated through software integration tests for all aspects that can be verified at this level.

The integration tests shall validate:
- Correct application of toolchains and build configuration, including successful compilation and linking using all defined build presets.
- Correct integration of platform-specific startup code and linker configuration, ensuring proper initialization and memory placement.
- Correct implementation of the defined memory layout:
  - Verification that memory regions are correctly mapped according to the linker script.
  - Validation of stack usage when applicable (e.g., using stack watermarking or equivalent techniques).
- Correct behavior of stdin/stdout integration defined by the Embedded C Framework (ECF), including buffering and ISR-based execution when applicable.
- Correct configuration and accessibility of debug interfaces defined in the platform software design.

Tests shall be executed in the target environment (HiL) when validation depends on hardware-specific behavior or execution environment characteristics that cannot be reliably validated in host environments.