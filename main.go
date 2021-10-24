package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/fatih/color"
)

func main() {
	var check bool
	flag.BoolVar(&check, "check", false, "Show branches to delete")
	flag.Parse()

	whiteBold := color.New(color.FgWhite, color.Bold).SprintFunc()

	if check {
		fmt.Fprintf(color.Output, "%s\n", whiteBold("== DRY RUN ==\n"))
	}

	conn := &ConnectionImpl{}

	branches, err := GetBranches(conn)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}

	if !check {
		branches = DeleteBranches(branches, conn)
	}

	var deletedStates []BranchState
	var notDeletedStates []BranchState
	if check {
		deletedStates = []BranchState{Deletable}
		notDeletedStates = []BranchState{NotDeletable}
	} else {
		deletedStates = []BranchState{Deleted}
		notDeletedStates = []BranchState{Deletable, NotDeletable}
	}

	fmt.Fprintf(color.Output, "%s\n", whiteBold("Deleted branches"))
	printBranches(getBranches(branches, deletedStates))
	fmt.Println()

	fmt.Fprintf(color.Output, "%s\n", whiteBold("Not deleted branches"))
	printBranches(getBranches(branches, notDeletedStates))
	fmt.Println()
}

func printBranches(branches []Branch) {
	hiBlack := color.New(color.FgHiBlack).SprintFunc()

	if len(branches) == 0 {
		fmt.Fprintf(color.Output, "%s\n", hiBlack("  There are no branches in the current directory"))
	}

	white := color.New(color.FgWhite).SprintFunc()
	green := color.New(color.FgGreen).SprintFunc()

	for _, branch := range branches {
		branchName := fmt.Sprintf("%s", branch.Name)
		if branch.Head {
			fmt.Fprintf(color.Output, "* %s\n", green(branchName))
		} else {
			fmt.Fprintf(color.Output, "  %s\n", white(branchName))
		}

		for i, pr := range branch.PullRequests {
			number := fmt.Sprintf("#%v", pr.Number)
			issueNoColor := getIssueNoColor(pr.State)
			var line string
			if i == len(branch.PullRequests)-1 {
				line = "└─"
			} else {
				line = "├─"
			}

			fmt.Fprintf(color.Output, "    %s %v  %s\n",
				line,
				color.New(issueNoColor).SprintFunc()(number),
				white(pr.Url),
			)
		}
	}
}

func getIssueNoColor(state PullRequestState) color.Attribute {
	switch state {
	case Open:
		return color.FgGreen
	case Merged:
		return color.FgMagenta
	case Closed:
		return color.FgRed
	default:
		return color.FgHiBlack
	}
}

func getBranches(branches []Branch, states []BranchState) []Branch {
	results := []Branch{}
	for _, branch := range branches {
		if contains(branch.State, states) {
			results = append(results, branch)
		}
	}
	return results
}

func contains(state BranchState, states []BranchState) bool {
	for _, s := range states {
		if s == state {
			return true
		}
	}
	return false
}