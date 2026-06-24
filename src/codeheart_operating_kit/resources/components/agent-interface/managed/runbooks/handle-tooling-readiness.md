Last updated: 2026-06-24T13:51:18Z (UTC)

# Handle Tooling Readiness

Audience: hybrid

Intent:
Help an agent recover when a repository, module, or extension is blocked by missing local tooling.
The user should receive a plain-language explanation and one concrete next decision. The agent
should classify the blocker, run safe read-only checks, request approval before local changes,
recheck readiness, and return to the calling runbook when the blocker is resolved.

Success:
The required local tool is available, the calling runbook can continue, or the agent records a
clear capability blocker and stops without improvising unsafe setup.

Agent judgment boundary:
The agent may adapt the wording to the current user and platform, choose the matching baseline
tooling lane, and use module-owned install commands when the module provides them. The agent must
not silently install or repair tools, invent module-specific setup, record durable machine state,
or treat local tool readiness as external service authorization.

Stop boundary:
Stop before installing, repairing, changing PATH or shell startup files, accepting permission
prompts, using unofficial install sources, handling secrets, bypassing device policy, or touching
external services without the calling runbook's approval gate.

Use this runbook when a module onboarding or operation reaches a local environment blocker. It is
not a default setup checklist and not a command wrapper.

## Trigger

Use this runbook when:

- a module onboarding runbook reports a missing local package manager, runtime, CLI, PowerShell
  module, install path, or required local tool;
- a module operation reports the same local blocker after onboarding;
- the user explicitly asks to install or repair a tool required by the current repository or
  module;
- an agent-facing runbook declares a required tool and the read-only check reports `missing`,
  `blocked`, or `unknown`.

Do not use this runbook for module-owned live service blockers. The calling module still owns
authentication, consent, account permissions, licenses, service availability, tenant or workspace
state, and live external preflight.

## User-Facing Flow

Start with the outcome, not the internal mechanism:

```text
I need one local tool before I can continue with this module task.

The next blocker is: <plain-language blocker>.

How should we continue?
1. <Concrete install or repair action> (Recommended)
2. I will handle it another way
3. Stop here
```

Use one user-owned decision per turn by default. Ask a second question only when the answer is
required to avoid changing the wrong machine or installing the wrong tool.

For a macOS package-manager blocker, use this shape:

```text
This Mac needs Homebrew before I can install the requested tool.

How should we continue?
1. Install Homebrew (Recommended)
2. I will install it another way
3. Stop here
```

If the user chooses the install or repair action, ask for approval with local impact:

```text
I will use the official install route for <tool>. This changes local software on this computer.
I will recheck the tool afterward and then return to <module or task>.

Do you approve this local install or repair?
```

If the user chooses their own method:

```text
Okay. Please install or repair <tool> using your preferred method. When you are ready, tell me and
I will run the read-only readiness check again.
```

If the user stops:

```text
Stopped. I cannot continue this module task until <tool> is available. I have not changed local
software.
```

## Execution Path

1. Identify the calling runbook, module, extension, or repository task.
2. Identify the required capability and the missing local tool.
3. Classify the blocker as local environment or module-owned service state.
4. If it is module-owned service state, return to the calling module runbook.
5. Map the local blocker to one baseline tooling lane.
6. Run only read-only checks that are appropriate for the current platform and task.
7. Explain the blocker and present concrete user choices.
8. Ask explicit approval before local installation, repair, PATH changes, shell configuration,
   permission prompts, or sensitive reads.
9. Use an official vendor source, system package manager, or module-owned runbook for concrete
   commands.
10. Recheck readiness with a read-only check.
11. Return to the calling runbook when readiness is available.
12. Record a capability blocker and stop when readiness remains unavailable.

Do not write durable readiness state in V1. If the active task requires a run record, record only a
short non-secret blocker summary in that task's execution log or final summary.

## Baseline Tooling Catalog

This catalog is on-demand. Do not install all baseline tools during onboarding.

| Lane | Use When | Read-Only Check Examples | Ownership Boundary |
| --- | --- | --- | --- |
| `package-manager-bootstrap` | A package manager or bootstrap route is missing. | Check whether the package manager command exists and whether the module accepts another install method. | Operating Kit owns the generic choice pattern. Official sources own install instructions. |
| `powershell-runtime` | PowerShell itself is missing or not runnable. | Check `pwsh` availability and version when the command exists. | Operating Kit owns the generic readiness route. Modules own why PowerShell is needed. |
| `powershell-module` | PowerShell exists but a required module is missing. | Check whether the requested module is installed or importable. | Operating Kit owns the generic module-readiness pattern. Modules own concrete module names, versions, and install commands. |
| `node-runtime` | Node.js or its package manager is missing. | Check `node` and package-manager availability when requested by the repo or module. | Operating Kit owns the generic lane. Repositories or modules own exact version requirements. |
| `python-runtime` | Python or its package manager is missing. | Check the relevant Python command and package tool availability. | Operating Kit owns the generic lane. Repositories or modules own exact version and virtual environment requirements. |
| `browser-automation` | A local browser automation prerequisite is missing. | Check the browser or automation tool named by the calling runbook. | Operating Kit owns the generic lane. Browser tooling, plugins, or modules own concrete setup. |
| `document-pdf-tooling` | Document conversion, PDF rendering, or office-file tooling is missing. | Check the exact command or package named by the calling runbook. | Operating Kit owns the generic lane. Document/PDF skills or modules own exact tool requirements. |

## Official Source Rule

Use current official sources for install instructions. Do not transcribe stale commands into this
runbook when a vendor page or module-owned runbook is the better source of truth.

Public official starting points:

- Homebrew: `https://brew.sh/`
- PowerShell: `https://learn.microsoft.com/powershell/scripting/install/install-powershell`
- Node.js: `https://nodejs.org/en/download`
- Python: `https://www.python.org/downloads/`

When a module provides a concrete install command, treat the module as the source of truth for that
module-specific tool. Still use this runbook for user-facing approval, recheck, stop behavior, and
return-to-module handoff.

## Local Versus Service Blockers

Local environment blockers handled here:

- missing package manager;
- missing shell runtime;
- missing PowerShell runtime;
- missing PowerShell module;
- missing CLI or local executable;
- missing local browser automation prerequisite;
- missing document/PDF conversion tool;
- broken PATH or shell discovery for an installed tool.

Module-owned service blockers not handled here:

- external account sign-in;
- tenant, workspace, project, or service state;
- admin role, consent, license, permission, mailbox, site, bucket, database, or app readiness;
- API authorization;
- live external preflight;
- destructive remote cleanup.

When the blocker is service-owned, say so plainly and return to the module runbook:

```text
This is not a local tooling blocker. The next step belongs to the <module> service preflight:
<plain-language service blocker>.
```

## Evidence And Validation

For a resolved blocker, record in the current task summary or execution log:

- tool or lane checked;
- read-only check used;
- result after recheck;
- calling runbook returned to;
- any caveat that affects the next step.

For an unresolved blocker, record:

- tool or lane blocked;
- user choice or reason the install did not proceed;
- current capability that remains unavailable;
- calling runbook or module operation that must wait.

Do not record:

- secrets or tokens;
- raw command logs;
- full environment dumps;
- local absolute paths unless the user explicitly needs the path to repair the issue;
- personal account identifiers;
- tenant, customer, or service resource identifiers.

## Review Checklist

- The blocker is local environment readiness, not live service preflight.
- The user-facing explanation states the outcome before technical details.
- The user is offered concrete choices such as install, handle another way, or stop.
- The agent asks approval before local changes.
- The install or repair path uses an official source or module-owned command.
- The readiness recheck proves the local blocker is resolved before returning to the module.
- No durable machine-readiness state is committed.
- The module still owns module-specific commands and external service validation.
