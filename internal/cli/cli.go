package cli

import (
	"fmt"
	"io"
	"strings"

	commandimpl "github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/commands"
	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/version"
)

const prog = "codeheart-operating-kit"

var commands = []command{
	{
		name:        "onboard",
		help:        "Guide first-run Operating Kit onboarding",
		usage:       "[--target TARGET] [--project-name PROJECT_NAME] [--purpose {private-automation,company-automation,software-product}] [--yes] [--json]",
		description: "Render the agent-guided first-run onboarding script and setup plan.",
		options: []option{
			{flag: "--target TARGET", help: "Explicit folder to inspect or set up. Required with --yes."},
			{flag: "--project-name PROJECT_NAME", help: "Explicit Codex project folder name. Required with --yes."},
			{flag: "--purpose {private-automation,company-automation,software-product}", help: "Optional backward-compatible setup metadata; it does not select a different profile."},
			{flag: "--yes", help: "Write setup files after explicit decisions are supplied"},
			{flag: "--json"},
		},
		epilog: "First-run user setup should run without --yes so Codex can ask for language, project name, target folder, setup-plan review, and write approval in chat. Non-interactive --yes writes require explicit --target and --project-name values. Base onboarding does not install, offer, or implicitly check optional native capabilities.",
	},
	{
		name:  "inspect",
		help:  "Inspect a folder before setup",
		usage: "[--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--json"},
		},
	},
	{
		name:  "init",
		help:  "Initialize Operating Kit in a folder",
		usage: "[--project-name PROJECT_NAME] [--purpose {private-automation,company-automation,software-product}] [--selected-folder SELECTED_FOLDER] [--dry-run] [--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--project-name PROJECT_NAME"},
			{flag: "--purpose {private-automation,company-automation,software-product}", help: "Optional backward-compatible setup metadata; omitted for new generic setups."},
			{flag: "--selected-folder SELECTED_FOLDER"},
			{flag: "--dry-run", help: "Show the deterministic change plan without writing"},
			{flag: "--json"},
		},
	},
	{
		name:  "repair",
		help:  "Repair a compatible current-version installation",
		usage: "[--dry-run] [--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--dry-run", help: "Show the deterministic change plan without writing"},
			{flag: "--json"},
		},
	},
	{
		name:  "sync",
		help:  "Refresh managed Operating Kit files",
		usage: "[--release-manifest RELEASE_MANIFEST] [--dry-run] [--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--release-manifest RELEASE_MANIFEST", help: "Optional release manifest fixture to validate before sync"},
			{flag: "--dry-run", help: "Show the deterministic change plan without writing"},
			{flag: "--json"},
		},
	},
	{
		name:  "check",
		help:  "Check installed Operating Kit state",
		usage: "[--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--json"},
		},
	},
	{
		name:  "update-check",
		help:  "Check latest version metadata without applying updates",
		usage: "[--latest-version LATEST_VERSION] [--metadata-url METADATA_URL] [--now NOW] [--agent-notification] [--dry-run] [--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--latest-version LATEST_VERSION"},
			{flag: "--metadata-url METADATA_URL", help: "Latest release metadata URL; defaults to the public GitHub latest-release endpoint"},
			{flag: "--now NOW"},
			{flag: "--agent-notification"},
			{flag: "--dry-run", help: "Show the metadata change without writing"},
			{flag: "--json"},
		},
	},
	{
		name:  "upgrade",
		help:  "Verify and apply an explicitly approved version upgrade",
		usage: "--version VERSION [--catalog CATALOG] (--dry-run | --yes) [--json] [path]",
		options: []option{
			{flag: "path"},
			{flag: "--version VERSION"},
			{flag: "--catalog CATALOG", help: "External release catalog URL or local path"},
			{flag: "--dry-run", help: "Verify the candidate without changing installed state"},
			{flag: "--yes", help: "Approve the verified version change"},
			{flag: "--json"},
		},
	},
	{
		name:        "plans",
		help:        "Validate, list, inventory, and migrate repository plans",
		usage:       "{validate,list,inventory,migrate} ...",
		description: "Read canonical plan metadata and guarded legacy migration evidence.",
		options: []option{
			{flag: "--json", help: "Emit deterministic structured output where supported"},
		},
	},
	{
		name:        "portfolio",
		help:        "Configure and refresh branch-aware portfolio coordination",
		usage:       "{configure,scan} ...",
		description: "Discover exactly enrolled repositories and build a complete source-derived plan catalog.",
	},
}

var planSubcommands = []command{
	{name: "validate", help: "Validate repository plan metadata and catalog-mode rules", usage: "[--remote-overlays] [--json] [path]", options: []option{{flag: "path"}, {flag: "--remote-overlays", help: "Refresh scanner-owned mirrors and validate pushed remote plan overlays"}, {flag: "--json"}}},
	{name: "list", help: "Generate the current repository-scoped plan view", usage: "[--format {text,json}] [--json] [path]", options: []option{{flag: "path"}, {flag: "--format {text,json}", help: "Choose text or deterministic JSON output"}, {flag: "--json", help: "Backward-compatible alias for --format json"}}},
	{name: "inventory", help: "Record Git-backed migration evidence without changing plans", usage: "--output OUTPUT [--remote-overlays] [--json] [path]", options: []option{{flag: "path"}, {flag: "--output OUTPUT", help: "Required inventory artifact destination"}, {flag: "--remote-overlays", help: "Refresh scanner-owned mirrors and inventory pushed remote plan overlays"}, {flag: "--json"}}},
	{name: "migrate", help: "Apply an explicitly reviewed migration ledger", usage: "--ledger LEDGER (--dry-run | --yes) [--json] [path]", options: []option{{flag: "path"}, {flag: "--ledger LEDGER"}, {flag: "--dry-run"}, {flag: "--yes"}, {flag: "--json"}}},
}

var portfolioSubcommands = []command{
	{name: "configure", help: "Configure member or coordination-home identity and discovery scope", usage: "--role {member,coordination-home} --member-repository-id ID --coordination-home-id ID [--github-owner OWNER]... [--local-root ROOT]... (--dry-run | --yes) [--json] [path]", options: []option{{flag: "path"}, {flag: "--role {member,coordination-home}"}, {flag: "--member-repository-id ID"}, {flag: "--coordination-home-id ID"}, {flag: "--github-owner OWNER"}, {flag: "--local-root ROOT"}, {flag: "--dry-run"}, {flag: "--yes"}, {flag: "--json"}}},
	{name: "scan", help: "Refresh the complete source-derived portfolio catalog", usage: "[--format {text,json}] [--json] [path]", options: []option{{flag: "path"}, {flag: "--format {text,json}"}, {flag: "--json"}}},
}

type command struct {
	name        string
	help        string
	usage       string
	description string
	options     []option
	epilog      string
}

type option struct {
	flag string
	help string
}

// Run executes the CLI and returns a process exit code.
func Run(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printRootHelp(stderr)
		fmt.Fprintf(stderr, "%s: error: the following arguments are required: command\n", prog)
		return 2
	}
	if args[0] == "__upgrade-reconcile" {
		return commandimpl.RunUpgradeReconcile(args[1:], stdout, stderr)
	}
	if args[0] == "__upgrade-handoff" {
		return commandimpl.RunUpgradeHandoff(args[1:], stderr)
	}
	if args[0] == "__verify-content-identity" {
		return commandimpl.RunVerifyContentIdentity(args[1:], stderr)
	}
	if args[0] == "__verify-release-evidence" {
		return commandimpl.RunVerifyReleaseEvidence(args[1:], stderr)
	}
	if args[0] == "__cleanup-upgrade-handoff" {
		return commandimpl.RunCleanupUpgradeHandoff(args[1:], stderr)
	}

	switch args[0] {
	case "-h", "--help":
		printRootHelp(stdout)
		return 0
	case "--version":
		fmt.Fprintf(stdout, "%s %s\n", prog, version.Version)
		return 0
	}

	cmd, ok := findCommand(args[0])
	if !ok {
		printRootHelp(stderr)
		fmt.Fprintf(stderr, "%s: error: argument command: invalid choice: '%s' (choose from onboard, inspect, init, repair, sync, check, update-check, upgrade, plans, portfolio)\n", prog, args[0])
		return 2
	}
	if cmd.name == "plans" {
		return runPlansGroup(args[1:], stdout, stderr)
	}
	if cmd.name == "portfolio" {
		return runPortfolioGroup(args[1:], stdout, stderr)
	}

	if containsCallableHelp(cmd, args[1:]) {
		printCommandHelp(stdout, cmd)
		return 0
	}

	switch cmd.name {
	case "onboard":
		return commandimpl.RunOnboard(args[1:], stdout, stderr)
	case "inspect":
		return commandimpl.RunInspect(args[1:], stdout, stderr)
	case "init":
		return commandimpl.RunInit(args[1:], stdout, stderr)
	case "repair":
		return commandimpl.RunRepair(args[1:], stdout, stderr)
	case "sync":
		return commandimpl.RunSync(args[1:], stdout, stderr)
	case "check":
		return commandimpl.RunCheck(args[1:], stdout, stderr)
	case "update-check":
		return commandimpl.RunUpdateCheck(args[1:], stdout, stderr)
	case "upgrade":
		return commandimpl.RunUpgrade(args[1:], stdout, stderr)
	case "plans":
		return commandimpl.RunPlans(args[1:], stdout, stderr)
	case "portfolio":
		return commandimpl.RunPortfolio(args[1:], stdout, stderr)
	default:
		fmt.Fprintf(stderr, "%s: unknown command %q\n", prog, cmd.name)
		return 2
	}
}

func containsCallableHelp(cmd command, args []string) bool {
	valueOptions := cmd.valueOptions()
	for index := 0; index < len(args); index++ {
		arg := args[index]
		if arg == "-h" || arg == "--help" {
			return true
		}
		optionName, _, hasInlineValue := strings.Cut(arg, "=")
		if valueOptions[optionName] && !hasInlineValue && index+1 < len(args) {
			if strings.HasPrefix(args[index+1], "-") {
				return false
			}
			index++
		}
	}
	return false
}

func (cmd command) valueOptions() map[string]bool {
	values := map[string]bool{}
	for _, opt := range cmd.options {
		if !strings.HasPrefix(opt.flag, "--") {
			continue
		}
		name, _, hasValueName := strings.Cut(opt.flag, " ")
		if hasValueName {
			values[name] = true
		}
	}
	return values
}

func findCommand(name string) (command, bool) {
	for _, cmd := range commands {
		if cmd.name == name {
			return cmd, true
		}
	}
	return command{}, false
}

func printRootHelp(w io.Writer) {
	fmt.Fprintf(w, "usage: %s [-h] [--version] {onboard,inspect,init,repair,sync,check,update-check,upgrade,plans,portfolio} ...\n\n", prog)
	fmt.Fprintln(w, "positional arguments:")
	fmt.Fprintln(w, "  {onboard,inspect,init,repair,sync,check,update-check,upgrade,plans,portfolio}")
	for _, cmd := range commands {
		fmt.Fprintf(w, "    %-12s %s\n", cmd.name, cmd.help)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "options:")
	fmt.Fprintln(w, "  -h, --help            show this help message and exit")
	fmt.Fprintln(w, "  --version             show program's version number and exit")
}

func runPlansGroup(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printPlansHelp(stderr)
		fmt.Fprintf(stderr, "%s plans: error: the following arguments are required: subcommand\n", prog)
		return 2
	}
	if args[0] == "-h" || args[0] == "--help" {
		printPlansHelp(stdout)
		return 0
	}
	subcommand, ok := findPlanSubcommand(args[0])
	if !ok {
		printPlansHelp(stderr)
		fmt.Fprintf(stderr, "%s plans: error: invalid subcommand %q (choose from validate, list, inventory, migrate)\n", prog, args[0])
		return 2
	}
	if containsCallableHelp(subcommand, args[1:]) {
		printCommandHelp(stdout, command{name: "plans " + subcommand.name, help: subcommand.help, usage: subcommand.usage, description: subcommand.description, options: subcommand.options, epilog: subcommand.epilog})
		return 0
	}
	return commandimpl.RunPlans(args, stdout, stderr)
}

func findPlanSubcommand(name string) (command, bool) {
	for _, cmd := range planSubcommands {
		if cmd.name == name {
			return cmd, true
		}
	}
	return command{}, false
}

func printPlansHelp(w io.Writer) {
	fmt.Fprintf(w, "usage: %s plans [-h] {validate,list,inventory,migrate} ...\n\n", prog)
	fmt.Fprintln(w, "positional arguments:")
	fmt.Fprintln(w, "  {validate,list,inventory,migrate}")
	for _, cmd := range planSubcommands {
		fmt.Fprintf(w, "    %-12s %s\n", cmd.name, cmd.help)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "options:")
	fmt.Fprintln(w, "  -h, --help            show this help message and exit")
}

func runPortfolioGroup(args []string, stdout io.Writer, stderr io.Writer) int {
	if len(args) == 0 {
		printPortfolioHelp(stderr)
		fmt.Fprintf(stderr, "%s portfolio: error: the following arguments are required: subcommand\n", prog)
		return 2
	}
	if args[0] == "-h" || args[0] == "--help" {
		printPortfolioHelp(stdout)
		return 0
	}
	subcommand, ok := findPortfolioSubcommand(args[0])
	if !ok {
		printPortfolioHelp(stderr)
		fmt.Fprintf(stderr, "%s portfolio: error: invalid subcommand %q (choose from configure, scan)\n", prog, args[0])
		return 2
	}
	if containsCallableHelp(subcommand, args[1:]) {
		printCommandHelp(stdout, command{name: "portfolio " + subcommand.name, help: subcommand.help, usage: subcommand.usage, description: subcommand.description, options: subcommand.options, epilog: subcommand.epilog})
		return 0
	}
	return commandimpl.RunPortfolio(args, stdout, stderr)
}

func findPortfolioSubcommand(name string) (command, bool) {
	for _, cmd := range portfolioSubcommands {
		if cmd.name == name {
			return cmd, true
		}
	}
	return command{}, false
}

func printPortfolioHelp(w io.Writer) {
	fmt.Fprintf(w, "usage: %s portfolio [-h] {configure,scan} ...\n\n", prog)
	fmt.Fprintln(w, "positional arguments:")
	fmt.Fprintln(w, "  {configure,scan}")
	for _, cmd := range portfolioSubcommands {
		fmt.Fprintf(w, "    %-12s %s\n", cmd.name, cmd.help)
	}
	fmt.Fprintln(w)
	fmt.Fprintln(w, "options:")
	fmt.Fprintln(w, "  -h, --help            show this help message and exit")
}

func printCommandHelp(w io.Writer, cmd command) {
	fmt.Fprintf(w, "usage: %s %s [-h]", prog, cmd.name)
	if cmd.usage != "" {
		fmt.Fprintf(w, " %s", cmd.usage)
	}
	fmt.Fprintln(w)
	if cmd.description != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, cmd.description)
	}
	if len(cmd.options) > 0 {
		fmt.Fprintln(w)
		fmt.Fprintln(w, "options:")
		fmt.Fprintln(w, "  -h, --help            show this help message and exit")
		for _, opt := range cmd.options {
			if opt.help == "" {
				fmt.Fprintf(w, "  %-22s\n", opt.flag)
				continue
			}
			fmt.Fprintf(w, "  %-22s %s\n", opt.flag, opt.help)
		}
	}
	if cmd.epilog != "" {
		fmt.Fprintln(w)
		fmt.Fprintln(w, cmd.epilog)
	}
}

func CommandNames() []string {
	names := make([]string, 0, len(commands))
	for _, cmd := range commands {
		names = append(names, cmd.name)
	}
	return names
}

func IsKnownCommand(name string) bool {
	_, ok := findCommand(name)
	return ok
}

func NormalizeOutput(value string) string {
	return strings.ReplaceAll(value, "\r\n", "\n")
}
