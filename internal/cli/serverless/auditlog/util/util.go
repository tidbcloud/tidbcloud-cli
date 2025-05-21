package util

import "github.com/tidbcloud/tidbcloud-cli/internal/flag"

var InputDescription = map[string]string{
	flag.OutputPath:                "Input the download path, press Enter to skip and download to the current directory",
	flag.StartDate:                 "Input the start date of the download in the format of 'YYYY-MM-DD'",
	flag.EndDate:                   "Input the end date of the download in the format of 'YYYY-MM-DD'",
	flag.AuditLogFilterRuleName:    "Input the filter rule name",
	flag.AuditLogFilterRuleUsers:   "Input the filter rule users, use JSON format, e.g. [\"user1@host1\",\"user2@host2\"]",
	flag.AuditLogFilterRuleFilters: "Input the filter rule filters, use JSON format, e.g. [{\"classes\":[\"QUERY\",\"EXECUTE\"],\"tables\":[\"test.t1\"]},{\"classes\":[\"QUERY\"]}] or [{}] to filter all audit logs",
}
