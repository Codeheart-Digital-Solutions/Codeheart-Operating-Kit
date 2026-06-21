Last updated: 2026-06-21T14:53:02Z (UTC)

# Coordination Sync Pending

This kit-initialized consumer state file records portfolio coordination-home updates that could not
be applied because the configured coordination home was unavailable or unsafe to edit.

Operating Kit owns the file contract and format. This repository owns the pending-sync entries
after creation. Sync may recreate this baseline when the file is absent, but it must not overwrite
existing entries.

Local planning work can continue while coordination sync is pending. Follow
`.codeheart/kit/docs/planning-workflows/runbooks/maintain-plan-register.md` before adding,
applying, or completing pending-sync items.

## Pending Items

No pending coordination sync items recorded yet.

Use this shape when needed:

```md
## Pending Sync - YYYY-MM-DD - <affected plan ID or title>

Source repository: <member repository ID or repository name>
Target coordination register: <coordination_home_path>/<coordination_home_register_path>
Affected plan entry: <ID, title, and canonical path>
Intended change: <add | update | complete | supersede | archive | relation-update>
Reason: <why coordination-home sync is needed>
Date: YYYY-MM-DD
Session ref: <session ID, not recorded, or unavailable>
Status: pending

Notes:
- <brief note about why the coordination home was unavailable or unsafe to edit>
```

