package markdown_report

import (
	"time"

	"github.com/pb33f/libopenapi"
	wcModel "github.com/pb33f/libopenapi/what-changed/model"
)

type Commit struct {
	CreatedAt         time.Time                `json:"-"`
	UpdatedAt         time.Time                `json:"-"`
	ID                uint                     `gorm:"primaryKey" json:"-"`
	Hash              string                   `json:"commitHash"`
	Message           string                   `json:"message"`
	Author            string                   `json:"author"`
	AuthorEmail       string                   `gorm:"index" json:"authorEmail"`
	CommitDate        time.Time                `json:"committed"`
	Changes           *wcModel.DocumentChanges `gorm:"-" json:"changeReport,omitempty"`
	SerializedChanges []byte                   `gorm:"-" json:"-"`
	Data              []byte                   `gorm:"-" json:"-"`
	OldData           []byte                   `gorm:"-" json:"-"`
	Document          libopenapi.Document      `gorm:"-" json:"-"`
	OldDocument       libopenapi.Document      `gorm:"-" json:"-"`
	RepoDirectory     string                   `gorm:"-" json:"-"`
	FilePath          string                   `gorm:"-" json:"-"`
}

type TreeNode struct {
	TitleString     string          `json:"titleString"`
	Title           string          `json:"title,omitempty"`
	Key             string          `json:"key"`
	IsLeaf          bool            `json:"isLeaf,omitempty"`
	Selectable      bool            `json:"selectable,omitempty"`
	TotalChanges    int             `json:"totalChanges,omitempty"`
	BreakingChanges int             `json:"breakingChanges,omitempty"`
	Change          *wcModel.Change `json:"change,omitempty"`
	Disabled        bool            `json:"disabled,omitempty"`
	Children        []*TreeNode     `json:"children,omitempty"`
}

type ChangeStatistics struct {
	Total            int               `json:"total"`
	TotalBreaking    int               `json:"totalBreaking"`
	Added            int               `json:"added"`
	Modified         int               `json:"modified"`
	Removed          int               `json:"removed"`
	BreakingAdded    int               `json:"breakingAdded"`
	BreakingModified int               `json:"breakingModified"`
	BreakingRemoved  int               `json:"breakingRemoved"`
	Commit           *CommitStatistics `json:"commit,omitempty"`
}

type CommitStatistics struct {
	Date        string `json:"date,omitempty"`
	Message     string `json:"message,omitempty"`
	Author      string `json:"author,omitempty"`
	AuthorEmail string `json:"authorEmail,omitempty"`
	Hash        string `json:"hash,omitempty"`
}

// Estruturas para relatórios em Markdown
type MarkdownChange struct {
	Type        string `json:"type"`
	Property    string `json:"property"`
	Original    string `json:"original,omitempty"`
	New         string `json:"new,omitempty"`
	Breaking    bool   `json:"breaking"`
	Description string `json:"description"`
}

type MarkdownReportItem struct {
	CommitInfo    *CommitStatistics `json:"commitInfo"`
	Statistics    *ChangeStatistics `json:"statistics"`
	Changes       []*MarkdownChange `json:"changes,omitempty"`
	Summary       string            `json:"summary"`
	BreakingCount int               `json:"breakingCount"`
	TotalCount    int               `json:"totalCount"`
	Diff          string            `json:"diff,omitempty"`
}

type MarkdownReport struct {
	Title         string                `json:"title"`
	DateGenerated string                `json:"dateGenerated"`
	ReportItems   []*MarkdownReportItem `json:"reportItems"`
	Summary       *ChangeStatistics     `json:"summary"`
}
