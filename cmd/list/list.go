package list

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/go-github/v60/github"
	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

type ListOptions struct {
	Since  string
	SortBy string
	Order  string
}

func NewListOptions() *ListOptions {
	return &ListOptions{}
}

func NewCmdList() *cobra.Command {
	o := NewListOptions()

	cmd := &cobra.Command{
		Use:   "list",
		Short: "list some issue details",
		RunE: func(cmd *cobra.Command, args []string) error {
			return o.Run()
		},
	}

	cmd.Flags().StringVar(&o.Since, "since", "", "time since issue creation [today, yesterday, week, month]")

	return cmd
}

func (o *ListOptions) Run() error {

	client := github.NewClient(nil).WithAuthToken(viper.GetString("ghAuthToken"))

	var displayRows []table.Row
	displayRows = append(displayRows, table.Row{"repo", "id", "desc", "comments", "assignee", "labels"})
	repos := viper.GetStringSlice("ghRepos")
	for i := range repos {

		listOptions := &github.IssueListByRepoOptions{
			Since: parseDateTimeOption(o.Since),
		}

		issues, _, err := client.Issues.ListByRepo(context.TODO(), viper.GetString("ghOrg"), repos[i], listOptions)
		if err != nil {
			return fmt.Errorf("error getting issues by repo : %v", err)
		}

		if issues != nil {

			for j := range issues {

				displayRows = append(displayRows, table.Row{
					repos[i],
					fmt.Sprintf("#%d", *issues[j].Number),
					*issues[j].Title,
					*issues[j].Comments,
					issues[j].Assignee.GetLogin(),
					labelArrayToString(issues[j].Labels),
				})
			}
		}
	}

	printTable(displayRows)

	return nil
}

func printTable(rows []table.Row) {
	t := table.NewWriter()
	t.SetOutputMirror(os.Stdout)
	t.AppendHeader(rows[0])
	t.AppendRows(rows[1:])
	t.SetColumnConfigs([]table.ColumnConfig{
		{Number: 1, AutoMerge: true},
	})
	t.Render()
}

func labelArrayToString(labels []*github.Label) string {
	var sb strings.Builder
	for i := range labels {
		sb.WriteString(
			fmt.Sprintf("%s, ", *labels[i].Name),
		)
	}

	return strings.TrimSuffix(sb.String(), ", ")
}

func parseDateTimeOption(opt string) time.Time {
	if opt == "" {
		return time.Time{}
	}

	var t time.Time
	now := time.Now()

	switch opt {
	case "yesterday":
		t = now.Add(time.Duration(-24) * time.Hour)
	case "last-week":
		t = now.Add((time.Duration(-24) * time.Hour) * 7)
	case "last-month":
		t = now.Add((time.Duration(-24) * time.Hour) * 30)
	}

	return t
}
