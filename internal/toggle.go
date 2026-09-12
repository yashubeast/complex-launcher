package internal

import (
	"slices"
	"fmt"
	"os/exec"
	"strings"
)

func Toggle(config Config) error {
	windowClass := config.WindowClassName

	windowId, err := getWindowId(windowClass)
	if err != nil { return err }

	// window doesn't exist -> launch it
	if windowId == "" {
		return launchWindow(config.TerminalCommand)
	}

	visible, err := isWindowVisible(windowId, windowClass)
	if err != nil { return err }

	// window exists but is hidden -> show and focus
	if !visible {
		return showWindow(windowId)
	}

	// window exists and is visible -> hide it
	return hideWindow(windowId)
}

func getWindowId(windowClass string) (string, error) {
	output, err := exec.Command(
		"xdotool",
		"search",
		"--class",
		"^" + windowClass + "$",
	).Output()
	// xdotool returns an error when no matching windows exist
	if err != nil { return "", nil }

	return strings.Split(strings.TrimSpace(string(output)), "\n")[0], nil
}

func isWindowVisible(windowId string, windowClass string) (bool, error) {
	output, err := exec.Command(
		"xdotool",
		"search",
		"--onlyvisible",
		"--class",
		"^" + windowClass + "$",
	).Output()
	if err != nil { return false, nil }

	if slices.Contains(strings.Split(strings.TrimSpace(string(output)), "\n"), windowId) { return true, nil }
	return false, nil
}

func showWindow(windowId string) error {
	err := exec.Command(
		"xdotool",
		"windowmap",
		windowId,
	).Run()
	if err != nil { return err }

	return exec.Command(
		"xdotool",
		"windowactivate",
		windowId,
	).Run()
}

func hideWindow(windowId string) error {
	return exec.Command(
		"xdotool",
		"windowunmap",
		windowId,
	).Run()
}
func hideWindowWithClass(windowClassName string) error {
	windowId, err := getWindowId(windowClassName)
	if err != nil { return err }
	return hideWindow(windowId)
}

func launchWindow(terminalCommand string) error {
	args := strings.Fields(terminalCommand)
	if len(args) == 0 {
		return fmt.Errorf("terminal command is empty")
	}

	cmd := exec.Command( args[0], args[1:]...)
	return cmd.Start()
}
