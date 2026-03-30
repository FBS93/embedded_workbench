# specs-definition

## Purpose

Define project specifications from stakeholder specifications.

## Input work products

- Stakeholder specifications

## Output work products

- Specifications
- Tools

## Steps

1. Review the stakeholder specifications.
2. Identify specification elements by extracting the smallest meaningful text fragments.
3. Define the specifications.
4. Check the specifications for completeness, consistency, and correctness.
5. Identify tools available for project development, validation and debugging.
6. Define the tools.
7. Check the tools for completeness, consistency, and correctness.

## Guidelines

### Specifications work product

A dedicated `specs.md` document shall be defined to capture all specifications relevant to the system design by decomposing stakeholder specifications into traceable specifications.

Each specification of the `specs.md` document shall:
- Be derived directly from stakeholder specifications without modifying their original content.
- Be defined by extracting the smallest possible meaningful text fragments.
- Preserve the original text exactly as defined in stakeholder specifications.
- Ensure each specification element is self-contained and understandable.
- Define a unique identifier for each specification element following the pattern: `SPEC_<X>`
  - `<X>` is a monotonically increasing number that shall never be reused.

The following specification template shall be used:

```md
# SPEC_<X>

<Exact text extracted from stakeholder specifications>.
```

### Tools work product

A dedicated `tools.md` document shall be defined to capture all available tools for project development, validation and debugging by reviewing stakeholder specifications. 

Each tool of the `tools.md` document shall:
- Define all project-specific tools available for implementation, validation, and debugging.
- Describe how each tool shall be used. 
  - Reference how it is integrated into the [Embedded C Workbench (xxxyy aquest nom pot cmbiar)](xxxyyy-link) when applicable instead of redefining it.
  - Reference tool-specific manuals or stakeholder documentation that support its usage.

The following tool template shall be used:

```md
# <Tool name>

## Overview

<Description of the tool and in which cases is useful>.

## Usage

<Usage description. May include references to manufacturer manuals or stakeholder specifications explaining how to use the tool>.
```