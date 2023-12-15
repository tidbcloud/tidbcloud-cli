// Copyright 2022 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package flag

const (
	AwsRoleArn            string = "aws-role-arn"
	BackslashEscape       string = "backslash-escape"
	ClientName            string = "client"
	CloudProvider         string = "cloud-provider"
	ClusterID             string = "cluster-id"
	ClusterIDShort        string = "c"
	DisplayName           string = "display-name"
	DisplayNameShort      string = "n"
	ClusterType           string = "cluster-type"
	BranchID              string = "branch-id"
	BranchIDShort         string = "b"
	Database              string = "database"
	DataFormat            string = "data-format"
	Debug                 string = "debug"
	DebugShort            string = "D"
	Delimiter             string = "delimiter"
	Force                 string = "force"
	ImportID              string = "import-id"
	NoColor               string = "no-color"
	OperatingSystem       string = "operating-system"
	Output                string = "output"
	OutputShort           string = "o"
	Password              string = "password"
	ProjectID             string = "project-id"
	ProjectIDShort        string = "p"
	ProfileName           string = "profile-name"
	Profile               string = "profile"
	ProfileShort          string = "P"
	PublicKey             string = "public-key"
	PrivateKey            string = "private-key"
	Query                 string = "query"
	QueryShort            string = "q"
	Region                string = "region"
	RegionShort           string = "r"
	RootPassword          string = "root-password"
	Separator             string = "separator"
	SkipCreateTable       string = "skip-create-table"
	SourceDatabase        string = "source-database"
	SourceHost            string = "source-host"
	SourcePassword        string = "source-password"
	SourcePort            string = "source-port"
	SourceTable           string = "source-table"
	SourceUser            string = "source-user"
	SourceUrl             string = "source-url"
	TargetDatabase        string = "target-database"
	TargetPassword        string = "target-password"
	TargetTable           string = "target-table"
	TargetUser            string = "target-user"
	TrimLastSeparator     string = "trim-last-separator"
	User                  string = "user"
	UserShort             string = "u"
	View                  string = "view"
	ViewShort             string = "v"
	SpendingLimitMonthly  string = "spending-limit-monthly"
	ServerlessLabels      string = "labels"
	ServerlessAnnotations string = "annotations"
	Monthly               string = "monthly"
	BackupID              string = "backup-id"
	RestoreMode           string = "restore-mode"
	BackupTime            string = "backup-time"
)

const (
	BasicView string = "BASIC"
	FullView  string = "FULL"
)

const (
	RestoreModeSnapshot    string = "Snapshot"
	RestoreModePointInTime string = "PointInTime"
)
