// Copyright 2024 PingCAP, Inc.
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
	ClusterID            string = "cluster-id"
	ClusterIDShort       string = "c"
	LocalConcurrency     string = "local.concurrency"
	CSVBackslashEscape   string = "csv.backslash-escape"
	CSVDelimiter         string = "csv.delimiter"
	CSVSeparator         string = "csv.separator"
	CSVTrimLastSeparator string = "csv.trim-last-separator"
	CSVNullValue         string = "csv.null-value"
	CSVSkipHeader        string = "csv.skip-header"
	DisplayName          string = "display-name"
	DisplayNameShort     string = "n"
	ClusterType          string = "cluster-type"
	BranchID             string = "branch-id"
	BranchIDShort        string = "b"
	Debug                string = "debug"
	DebugShort           string = "D"
	LocalFilePath        string = "local.file-path"
	Force                string = "force"
	ImportID             string = "import-id"
	NoColor              string = "no-color"
	Output               string = "output"
	OutputShort          string = "o"
	Password             string = "password"
	ProjectID            string = "project-id"
	ProjectIDShort       string = "p"
	ProfileName          string = "profile-name"
	Profile              string = "profile"
	ProfileShort         string = "P"
	PublicKey            string = "public-key"
	PrivateKey           string = "private-key"
	Query                string = "query"
	QueryShort           string = "q"
	Region               string = "region"
	RegionShort          string = "r"
	LocalTargetDatabase  string = "local.target-database"
	LocalTargetTable     string = "local.target-table"
	User                 string = "user"
	UserShort            string = "u"
	SpendingLimitMonthly string = "spending-limit-monthly"
	ServerlessLabels     string = "labels"
	Monthly              string = "monthly"
	BackupID             string = "backup-id"
	BackupTime           string = "backup-time"
	// External storage
	S3URI                  string = "s3.uri"
	S3AccessKeyID          string = "s3.access-key-id"
	S3SecretAccessKey      string = "s3.secret-access-key"
	S3RoleArn              string = "s3.role-arn"
	GCSURI                 string = "gcs.uri"
	GCSServiceAccountKey   string = "gcs.service-account-key"
	AzureBlobURI           string = "azblob.uri"
	AzureBlobSASToken      string = "azblob.sas-token"
	TargetType             string = "target-type"
	FileType               string = "file-type"
	ExportID               string = "export-id"
	ExportIDShort          string = "e"
	OutputPath             string = "output-path"
	Encryption             string = "encryption"
	Compression            string = "compression"
	SourceType             string = "source-type"
	UserRole               string = "role"
	AddRole                string = "add-role"
	DeleteRole             string = "delete-role"
	Concurrency            string = "concurrency"
	SQL                    string = "sql"
	TableWhere             string = "where"
	TableFilter            string = "filter"
	ParentID               string = "parent-id"
	ParentTimestamp        string = "parent-timestamp"
	PublicEndpointDisabled string = "disable-public-endpoint"
	ParquetCompression     string = "parquet.compression"
)

const OutputHelp = "Output format, one of [\"human\" \"json\"]. For the complete result, please use json format."
