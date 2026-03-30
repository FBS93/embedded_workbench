# Hardware domain overview

## Processes

Currently, the hardware domain is limited to the definition of the `hardware-software interface`; therefore, no processes are defined.

@todo: Refine the HW/SW interface definition to specify in detail the expected content of the document. This shall include, for example, the definition of available test points, debug pins, and other hardware observability mechanisms. Providing this level of detail will enable the software domain to define more accurate and measurable integration and system tests, particularly for platform validation (e.g., using logic analyzers or similar tools). Once the HW/SW interface is fully specified, this information shall be leveraged to improve software processes by defining more precise and mandatory test strategies.

## Work products organization

The hardware domain shall be organized by process. Where applicable, processes shall have dedicated folders containing their generated work products. The structure shall be as follows:

```text
hw/
  hw_sw_if/
```

The `hw/hw_sw_if/` folder contains the work products generated for the `hardware-software interface`.