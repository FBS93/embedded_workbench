# sys-requirements

## Purpose

Define system requirements from specifications.

## Input work products

- Specifications

## Output work products

- System requirements

## Steps

1. Review the specifications.
2. Identify the needs to be addressed by the system requirements.
3. Define the system requirements derived from the specifications.
4. Define traceability between system requirements and specifications.
5. Check the system requirements for completeness, consistency, and correctness.

## Guidelines

### System requirements work product

The system requirements work product shall follow the [Requirements definition](../../resources/requirements_definition.md).

#### Domain attribute

System requirements shall define an additional domain attribute to specify the engineering domains to which the requirement applies.

The domain attribute shall:
- Be defined as a list of engineering domain abbreviations in lower case.
- Specify all domains responsible for implementing or fulfilling the requirement.
- Be used to support downstream allocation of system requirements to engineering domain requirements.

The domain attribute shall be included in all system requirements using the following requirement template:

```md
# SYS_<FUNC_DOMAIN>_REQ_<X>

<subject> shall <required behavior>.

Domain:
- <engineering domain abbreviation>

Upstream traceability:
- [<UPSTREAM_ELEMENT_UNIQUE_ID>](#...)

References:
- [SYS_...](#...)