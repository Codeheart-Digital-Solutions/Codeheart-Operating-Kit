package commands

import (
	"fmt"
	"io"
	"strings"
	"time"

	"github.com/Codeheart-Digital-Solutions/Codeheart-Operating-Kit/internal/reconcile"
)

func RunOnboard(args []string, stdout io.Writer, stderr io.Writer) int {
	values, bools, positionals, err := parseValueArgs(args, map[string]bool{
		"--target":       true,
		"--project-name": true,
		"--purpose":      true,
		"--yes":          false,
		"--json":         false,
	})
	if err != nil {
		return writeArgError(stderr, "onboard", err)
	}
	if len(positionals) > 0 {
		return writeArgError(stderr, "onboard", fmt.Errorf("unexpected argument %q", positionals[0]))
	}
	if err := validatePurpose(values["--purpose"]); err != nil {
		return writeArgError(stderr, "onboard", err)
	}
	result, script, code, err := Onboard(values["--target"], values["--project-name"], values["--purpose"], bools["--yes"], time.Now())
	if err != nil {
		fmt.Fprintf(stderr, "codeheart-operating-kit onboard: error: %v\n", err)
		return 1
	}
	if bools["--json"] {
		payload := map[string]any{}
		for key, value := range result {
			payload[key] = value
		}
		payload["script"] = stringsAsAny(script)
		if err := writeJSON(stdout, payload); err != nil {
			fmt.Fprintf(stderr, "codeheart-operating-kit onboard: error: %v\n", err)
			return 1
		}
		return code
	}
	fmt.Fprintln(stdout, strings.Join(script, "\n"))
	if errText, ok := result["error"].(string); ok && errText != "" {
		fmt.Fprintf(stdout, "Error: %s\n", errText)
	} else if written, ok := result["written"].(bool); ok && written {
		fmt.Fprintln(stdout, "Base Operating Kit setup is complete.")
	}
	return code
}

func Onboard(target string, projectName string, purpose string, yes bool, now time.Time) (map[string]any, []string, int, error) {
	var inspection map[string]any
	mode := ""
	expandedTarget := ""
	if target != "" {
		expandedTarget = expandPath(target)
		inspection = InspectFolder(expandedTarget)
		mode = fmt.Sprint(inspection["mode"])
	} else {
		inspection = map[string]any{
			"path":   nil,
			"mode":   "target-folder-pending",
			"reason": "target folder must be chosen by the user before inspection",
		}
	}
	script := OnboardingScript(projectName, purpose, expandedTarget, mode)
	missing := requiredUserDecisionsMissing(target, projectName)
	result := map[string]any{
		"inspection":                      inspection,
		"written":                         false,
		"write_approved":                  yes,
		"required_user_decisions_missing": stringsAsAny(missing),
	}
	if yes && len(missing) > 0 {
		result["error"] = "Cannot write setup files because required user decisions are missing: " + strings.Join(missing, ", ")
		return result, script, 2, nil
	}
	if yes && inspection["mode"] == "ambiguous-folder-stop" {
		result["error"] = "Cannot write setup files because folder inspection is ambiguous."
		return result, script, 1, nil
	}
	if yes {
		var changed map[string]any
		var operationOK bool
		var err error
		if inspection["mode"] == "existing-operating-kit-repair" {
			var operationResult reconcile.Result
			changed, operationResult, err = repairOperation(expandedTarget, now, false)
			operationOK = operationResult.OK()
		} else {
			var operationResult reconcile.Result
			changed, operationResult, err = initializeOperation(expandedTarget, projectName, purpose, expandedTarget, now, false)
			operationOK = operationResult.OK()
		}
		if err != nil {
			return nil, nil, 1, err
		}
		for key, value := range changed {
			result[key] = value
		}
		if !operationOK {
			result["error"] = "The selected Operating Kit lifecycle operation was blocked. Run check for remediation."
			return result, script, 1, nil
		}
		result["written"] = true
		result["write_approved"] = true
		result["required_user_decisions_missing"] = []any{}
	}
	return result, script, 0, nil
}

func requiredUserDecisionsMissing(target string, projectName string) []string {
	missing := []string{}
	if target == "" {
		missing = append(missing, "target_folder")
	}
	if projectName == "" {
		missing = append(missing, "project_name")
	}
	return missing
}

func OnboardingScript(projectName string, purpose string, target string, mode string) []string {
	displayName := projectName
	if displayName == "" {
		displayName = "<Project-Name>"
	}
	selectedFolder := target
	if selectedFolder == "" {
		selectedFolder = "<Selected-Folder>"
	}
	recommended := "Documents > " + displayName
	lines := []string{
		"Choose setup language / Sprache waehlen / 选择设置语言:",
		"1. English",
		"2. Deutsch",
		"3. 中文",
		"",
		"Before we set up your project folder, please adjust Codex in this chat.",
		"Look at the message box on the right. In the lower-right area, open the menu for model, thinking, and speed.",
		"Set:",
		"- Model: GPT-5.5",
		"- Thinking: Extra High",
		"- Speed: Fast",
		"Tell me when this is done.",
		"",
		"Now open Codex Settings.",
		"Look at the left sidebar. At the bottom-left, click Settings.",
		"The General tab should open automatically.",
		"At the very top of the main settings screen, find Work Mode and select Coding.",
		"Directly beneath that, find Permissions and turn on all three setup options:",
		"- Default permissions",
		"- Auto review",
		"- Full access",
		"Then return to this chat. In the chat box area, check the lower-left control named Approve for me.",
		"Turn it on when it is not already selected.",
		"Tell me when this is done.",
		"",
		"Do you already know what this Codex project should be called, or should I suggest a name?",
		"",
		"If the user wants help, ask:",
		"What is this mainly for?",
		"1. Personal automation",
		"2. Company operations",
		"3. Software or product development",
		"4. I am not sure yet",
	}
	if purpose != "" {
		lines = append(lines,
			"",
			fmt.Sprintf("Explicit setup context metadata supplied: %s (%s).", purpose, purposeLabels[purpose]),
			"Use this only for naming help and next-step guidance; it does not change the standard profile.",
		)
	}
	lines = append(lines,
		"",
		"Codex needs one project folder.",
		"The folder name is important because Codex will show it in the left sidebar. Chat threads for this work will be grouped under that project name, so choose a name you will recognize later.",
		"Selected project name: "+displayName,
		"Neutral examples: Yourname-Automation; Companyname-Automation; Productname-Development; Teamname-Operations",
		"",
		"Do you already know where this project folder should be, or should I suggest a simple location?",
		"Recommended folder: "+recommended,
		"What should I use?",
		"1. Yes, use "+recommended,
		"2. Use a different folder",
		"Please tell me the folder name or location if you choose a different folder.",
		"Examples: Documents > Companyname-Automation; Documents > Yourname-Automation; Desktop > Productname-Development; Documents > Existing-Project-Name",
		"",
		"Thank you. I will check this folder now so I can prepare the setup plan: "+selectedFolder,
		"This check only looks at the folder. I will show you the setup plan before changing files.",
	)
	if mode != "" {
		lines = append(lines, modeMessage(mode), planPreview(mode))
	} else {
		lines = append(lines,
			"Folder inspection is waiting for the selected target folder.",
			"Do not write setup files until the user supplies the target folder, sees the setup plan, and approves setup.",
		)
	}
	lines = append(lines,
		"Should I continue with this setup?",
		"1. Yes, set it up",
		"2. No, stop here",
		"After successful setup, finish with: Base Operating Kit setup is complete.",
		"Foundry module setup will become available later, after the first Foundry module is released.",
	)
	return lines
}

func modeMessage(mode string) string {
	messages := map[string]string{
		"new-folder-setup":                    "This folder is ready for a new setup. I can add Codex working instructions and a small memory area.",
		"existing-folder-setup":               "This folder already contains files. I can set up Operating Kit here without replacing your existing files. I will show the exact additions before changing anything.",
		"existing-technical-project-adoption": "This is an existing technical project. I will not overwrite existing docs or instructions. I will add the managed Operating Kit area, preserve local instructions, scaffold only missing memory files, and create an adoption cleanup report for overlapping docs.",
		"existing-operating-kit-repair":       "This folder already has Operating Kit. I will check whether the managed kit files, routing, and lifecycle state need repair. I will not apply a version update unless you ask for it.",
		"ambiguous-folder-stop":               "I cannot tell whether this folder is safe to set up. Please use a different folder, or tell me more about what this folder is for.",
	}
	return messages[mode]
}

func planPreview(mode string) string {
	switch mode {
	case "existing-technical-project-adoption":
		return "Adoption plan: add .codeheart/kit/, config, lock, managed AGENTS block, missing memory files, and .codeheart/reports/adoption-cleanup-report.md without deleting overlapping docs."
	case "existing-operating-kit-repair":
		return "Repair plan: check managed Operating Kit files, agent routing, setup information, and lifecycle metadata."
	case "ambiguous-folder-stop":
		return "Stop plan: choose a different folder or provide more context before writing files."
	default:
		return "Setup plan: add .codeheart/kit/, kit config, kit lock, local user notes, AGENTS.md, docs/repo/, and missing docs/agent-memory/ files. I will not delete existing files."
	}
}
