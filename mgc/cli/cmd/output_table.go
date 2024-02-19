package cmd

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"text/scanner"

	"slices"

	"magalu.cloud/core/utils"

	"github.com/jedib0t/go-pretty/v6/table"
	"github.com/jedib0t/go-pretty/v6/text"
)

type tableOutputFormatter struct{}

type alignString text.Align

var _ json.Unmarshaler = (*alignString)(nil)

func (a *alignString) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	switch strings.ToLower(s) {
	case "default", "aligndefault":
		*a = alignString(text.AlignDefault)
	case "left", "alignleft":
		*a = alignString(text.AlignLeft)
	case "center", "aligncenter":
		*a = alignString(text.AlignCenter)
	case "justify", "alignjustify":
		*a = alignString(text.AlignJustify)
	case "right", "alignright":
		*a = alignString(text.AlignRight)
	default:
		*a = alignString(text.AlignDefault)
	}

	return nil
}

type valignString text.VAlign

var _ json.Unmarshaler = (*valignString)(nil)

func (a *valignString) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	switch strings.ToLower(s) {
	case "default", "aligndefault":
		*a = valignString(text.VAlignDefault)
	case "top", "aligntop":
		*a = valignString(text.VAlignTop)
	case "middle", "alignmiddle":
		*a = valignString(text.VAlignMiddle)
	case "bottom", "alignbottom":
		*a = valignString(text.VAlignBottom)
	default:
		*a = valignString(text.VAlignDefault)
	}

	return nil
}

type colorString text.Color

var _ json.Unmarshaler = (*colorString)(nil)

func (c *colorString) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	switch strings.ToLower(s) {
	case "reset":
		*c = colorString(text.Reset)
	// Styling
	case "bold":
		*c = colorString(text.Bold)
	case "faint":
		*c = colorString(text.Faint)
	case "italic":
		*c = colorString(text.Italic)
	case "underline":
		*c = colorString(text.Underline)
	case "reversevideo":
		*c = colorString(text.ReverseVideo)
	case "concealed":
		*c = colorString(text.Concealed)
	case "crossedout":
		*c = colorString(text.CrossedOut)
	// Blinking
	case "blinkslow":
		*c = colorString(text.BlinkSlow)
	case "blinkrapid":
		*c = colorString(text.BlinkRapid)
	// Foreground colors
	case "black", "fgblack":
		*c = colorString(text.FgBlack)
	case "red", "fgred":
		*c = colorString(text.FgRed)
	case "green", "fggreen":
		*c = colorString(text.FgGreen)
	case "yellow", "fgyellow":
		*c = colorString(text.FgYellow)
	case "blue", "fgblue":
		*c = colorString(text.FgBlue)
	case "magenta", "fgmagenta":
		*c = colorString(text.FgMagenta)
	case "cyan", "fgcyan":
		*c = colorString(text.FgCyan)
	case "white", "fgwhite":
		*c = colorString(text.FgWhite)
	// Foreground high intensity colors
	case "fghiblack":
		*c = colorString(text.FgHiBlack)
	case "fghired":
		*c = colorString(text.FgHiRed)
	case "fghigreen":
		*c = colorString(text.FgHiGreen)
	case "fghiyellow":
		*c = colorString(text.FgHiYellow)
	case "fghiblue":
		*c = colorString(text.FgHiBlue)
	case "fghimagenta":
		*c = colorString(text.FgHiMagenta)
	case "fghicyan":
		*c = colorString(text.FgHiCyan)
	case "fghiwhite":
		*c = colorString(text.FgHiWhite)
	// Background colors
	case "bgblack":
		*c = colorString(text.BgBlack)
	case "bgred":
		*c = colorString(text.BgRed)
	case "bggreen":
		*c = colorString(text.BgGreen)
	case "bgyellow":
		*c = colorString(text.BgYellow)
	case "bgblue":
		*c = colorString(text.BgBlue)
	case "bgmagenta":
		*c = colorString(text.BgMagenta)
	case "bgcyan":
		*c = colorString(text.BgCyan)
	case "bgwhite":
		*c = colorString(text.BgWhite)
	// Background high intensity colors
	case "bghiblack":
		*c = colorString(text.BgHiBlack)
	case "bghired":
		*c = colorString(text.BgHiRed)
	case "bghigreen":
		*c = colorString(text.BgHiGreen)
	case "bghiyellow":
		*c = colorString(text.BgHiYellow)
	case "bghiblue":
		*c = colorString(text.BgHiBlue)
	case "bghimagenta":
		*c = colorString(text.BgHiMagenta)
	case "bghicyan":
		*c = colorString(text.BgHiCyan)
	case "bghiwhite":
		*c = colorString(text.BgHiWhite)
	default:
		return fmt.Errorf("unknown column color value: %s", s)
	}

	return nil
}

var noBorderBoxStyle = table.BoxStyle{
	BottomLeft:       "",
	BottomRight:      "",
	BottomSeparator:  "",
	EmptySeparator:   "",
	Left:             "",
	LeftSeparator:    "",
	MiddleHorizontal: "",
	MiddleSeparator:  "",
	MiddleVertical:   "",
	PaddingLeft:      " ",
	PaddingRight:     " ",
	PageSeparator:    "\n",
	Right:            "",
	RightSeparator:   "",
	TopLeft:          "",
	TopRight:         "",
	TopSeparator:     "",
	UnfinishedRow:    " ≈",
}

var noBorderStyle = tableStyleString(table.Style{
	Name:    "noBorderStyle",
	Box:     noBorderBoxStyle,
	Color:   table.ColorOptionsDefault,
	Format:  table.FormatOptionsDefault,
	HTML:    table.DefaultHTMLOptions,
	Options: table.OptionsDefault,
	Title:   table.TitleOptionsDefault,
})

type colorStrings []colorString

func (c *colorStrings) toColors() text.Colors {
	result := make(text.Colors, len(*c))
	for i, color := range *c {
		result[i] = text.Color(color)
	}
	return result
}

type columnConfig struct {
	Align        alignString  `yaml:"align"`
	AlignFooter  alignString  `yaml:"alignFooter"`
	AlignHeader  alignString  `yaml:"alignHeader"`
	AutoMerge    bool         `yaml:"autoMerge"`
	Colors       colorStrings `yaml:"colors"`
	ColorsFooter colorStrings `yaml:"colorsFooter"`
	ColorsHeader colorStrings `yaml:"colorsHeader"`
	Hidden       bool         `yaml:"hidden"`
	VAlign       valignString `yaml:"valign"`
	VAlignFooter valignString `yaml:"valignFooter"`
	VAlignHeader valignString `yaml:"valignHeader"`
	WidthMax     int          `yaml:"widthMax"`
	WidthMin     int          `yaml:"widthMin"`
}

type column struct {
	Name     string
	Parents  []string // if provided, multiple header lines will be created merging these
	JSONPath string
	Config   *columnConfig
}

func (c *column) tableColumnConfig() *table.ColumnConfig {
	t := &table.ColumnConfig{
		Name:         c.Name,
		Align:        text.Align(c.Config.Align),
		AlignFooter:  text.Align(c.Config.AlignFooter),
		AlignHeader:  text.Align(c.Config.AlignHeader),
		AutoMerge:    c.Config.AutoMerge,
		Colors:       c.Config.Colors.toColors(),
		ColorsFooter: c.Config.ColorsFooter.toColors(),
		ColorsHeader: c.Config.ColorsHeader.toColors(),
		Hidden:       c.Config.Hidden,
		VAlign:       text.VAlign(c.Config.VAlign),
		VAlignFooter: text.VAlign(c.Config.VAlignFooter),
		VAlignHeader: text.VAlign(c.Config.VAlignHeader),
		WidthMax:     c.Config.WidthMax,
		WidthMin:     c.Config.WidthMin,
	}
	return t
}

type tableRenderFormat string

const (
	Ascii    tableRenderFormat = "ascii"
	HTML     tableRenderFormat = "html"
	Markdown tableRenderFormat = "markdown"
	CSV      tableRenderFormat = "csv"
)

type tableStyleString table.Style

var _ json.Unmarshaler = (*tableStyleString)(nil)

func (t *tableStyleString) UnmarshalJSON(data []byte) (err error) {
	var s string
	err = json.Unmarshal(data, &s)
	if err != nil {
		return
	}

	switch strings.ToLower(s) {
	case "bold", "stylebold":
		*t = tableStyleString(table.StyleBold)
	case "coloredblackonbluewhite", "stylecoloredblackonbluewhite":
		*t = tableStyleString(table.StyleColoredBlackOnBlueWhite)
	case "coloredblackoncyanwhite", "stylecoloredblackoncyanwhite":
		*t = tableStyleString(table.StyleColoredBlackOnCyanWhite)
	case "coloredblackongreenwhite", "stylecoloredblackongreenwhite":
		*t = tableStyleString(table.StyleColoredBlackOnGreenWhite)
	case "coloredblackonmagentawhite", "stylecoloredblackonmagentawhite":
		*t = tableStyleString(table.StyleColoredBlackOnMagentaWhite)
	case "coloredblackonredwhite", "stylecoloredblackonredwhite":
		*t = tableStyleString(table.StyleColoredBlackOnRedWhite)
	case "coloredblackonyellowwhite", "stylecoloredblackonyellowwhite":
		*t = tableStyleString(table.StyleColoredBlackOnYellowWhite)
	case "coloredbright", "stylecoloredbright":
		*t = tableStyleString(table.StyleColoredBright)
	case "coloredcyanwhiteonblack", "stylecoloredcyanwhiteonblack":
		*t = tableStyleString(table.StyleColoredCyanWhiteOnBlack)
	case "coloreddark", "stylecoloreddark":
		*t = tableStyleString(table.StyleColoredDark)
	case "coloredgreenwhiteonblack", "stylecoloredgreenwhiteonblack":
		*t = tableStyleString(table.StyleColoredGreenWhiteOnBlack)
	case "coloredmagentawhiteonblack", "stylecoloredmagentawhiteonblack":
		*t = tableStyleString(table.StyleColoredMagentaWhiteOnBlack)
	case "coloredredwhiteonblack", "stylecoloredredwhiteonblack":
		*t = tableStyleString(table.StyleColoredRedWhiteOnBlack)
	case "coloredyellowwhiteonblack", "stylecoloredyellowwhiteonblack":
		*t = tableStyleString(table.StyleColoredYellowWhiteOnBlack)
	case "default", "styledefault":
		*t = tableStyleString(table.StyleDefault)
	case "double", "styledouble":
		*t = tableStyleString(table.StyleDouble)
	case "light", "stylelight":
		*t = tableStyleString(table.StyleLight)
	case "rounded", "stylerounded":
		*t = tableStyleString(table.StyleRounded)
	default:
		*t = tableStyleString(table.StyleDefault)
	}

	return nil
}

type tableOptions struct {
	// If true, the Columns will become Rows and the Rows will become Columns
	BuildVertically     bool
	Columns             []*column
	RowLength           int
	Format              tableRenderFormat
	Style               *tableStyleString
	StyleCustomizations func(*table.Style)
}

func concreteKind(v reflect.Value) reflect.Kind {
	kind := v.Kind()
	if kind == reflect.Interface || kind == reflect.Pointer {
		kind = v.Elem().Kind()
	}
	return kind
}

func columnsFromPointerOrInterface(v reflect.Value, prefix string) ([]*column, error) {
	if v.IsNil() {
		return []*column{{Name: "", JSONPath: prefix}}, nil
	}

	return columnsFromAny(v.Elem().Interface(), prefix)
}

func columnsFromStruct(v reflect.Value, prefix string) ([]*column, error) {
	t := v.Type()
	fieldCount := t.NumField()

	result := make([]*column, fieldCount)
	for i := 0; i < fieldCount; i++ {
		field := t.Field(i)
		if field.IsExported() {
			result[i] = &column{
				Name:     strings.ToUpper(field.Name),
				JSONPath: fmt.Sprintf("%s[%q]", prefix, field.Name),
			}
		}
	}

	return result, nil
}

func columnsFromArrayOrSlice(v reflect.Value, prefix string) ([]*column, error) {
	if v.Len() == 0 {
		return []*column{{Name: "", JSONPath: prefix}}, nil
	}

	subVal := v.Index(0)
	return columnsFromAny(subVal.Interface(), prefix+"[*]")
}

type columnGroup struct {
	name    string
	columns []*column
}

func columnsFromMap(v reflect.Value, prefix string) (result []*column, err error) {
	keys := v.MapKeys()

	var columnGroups []*columnGroup
	columnCount := 0

	for _, key := range keys {
		if key.Kind() == reflect.String {
			subVal := v.MapIndex(key)
			subKind := concreteKind(subVal)
			name := strings.ToUpper(key.String())
			jsonPath := fmt.Sprintf("%s[%q]", prefix, key)

			var columns []*column
			switch subKind {
			case reflect.Array, reflect.Slice:
				columns, err = columnsFromAny(subVal.Interface(), jsonPath)

			case reflect.Map:
				columns, err = columnsFromAny(subVal.Interface(), jsonPath)

			default:
				columns = []*column{{Name: name, JSONPath: jsonPath}}
			}

			if err != nil {
				return
			}
			columnCount += len(columns)
			columnGroups = append(columnGroups, &columnGroup{name: name, columns: columns})
		}
	}

	if len(columnGroups) == 1 {
		result = columnGroups[0].columns
		return
	}
	slices.SortFunc(columnGroups, func(l, r *columnGroup) int {
		return strings.Compare(l.name, r.name)
	})

	result = make([]*column, 0, columnCount)
	for _, group := range columnGroups {
		if len(group.columns) == 1 && group.name == group.columns[0].Name {
			// NOTE: checking group.name and column.Name avoids replicated header lines
			// but will keep 2 lines in case of "MACHINE" and "ID".
			// If one want to avoid such, remove this conditional.
			result = append(result, group.columns[0])
		} else {
			for _, c := range group.columns {
				c.Parents = append(c.Parents, group.name)
				result = append(result, c)
			}
		}
	}

	return result, nil
}

func columnsFromAny(val any, prefix string) ([]*column, error) {
	if prefix == "" {
		prefix = "$"
	}

	v := reflect.ValueOf(val)
	if !v.IsValid() {
		return []*column{{Name: "RESULT", JSONPath: prefix}}, nil
	}

	switch v.Kind() {
	case reflect.Pointer, reflect.Interface:
		return columnsFromPointerOrInterface(v, prefix)
	case reflect.Struct:
		return columnsFromStruct(v, prefix)
	case reflect.Array, reflect.Slice:
		return columnsFromArrayOrSlice(v, prefix)
	case reflect.Map:
		return columnsFromMap(v, prefix)
	default:
		return []*column{{Name: "", JSONPath: prefix}}, nil
	}
}

func splitUnquoted(str string, separator string) (result []string, err error) {
	cur := ""
	result = []string{}
	var s scanner.Scanner
	s.Init(strings.NewReader(str))
	s.Mode |= scanner.ScanStrings
	s.Error = func(s *scanner.Scanner, msg string) {} // function does nothing so errors are not printed
	for tok := s.Scan(); tok != scanner.EOF; tok = s.Scan() {
		txt := s.TokenText()
		if txt == "" {
			continue
		}
		switch txt[0] {
		case '"', '\'', '`':
			txt, err = strconv.Unquote(txt)
			if err != nil {
				return
			}
		default:
		}
		if txt == separator {
			result = append(result, cur)
			cur = ""
		} else {
			cur += txt
		}
	}
	if cur != "" {
		result = append(result, cur)
	}
	return
}

func columnsFromString(str string) ([]*column, error) {
	colStrings, err := splitUnquoted(str, ",")
	if err != nil {
		return nil, err
	}
	result := make([]*column, 0)
	for _, colString := range colStrings {
		if colString == "" {
			continue
		}
		parts := strings.SplitN(colString, ":", 2)
		if len(parts) != 2 {
			return nil, fmt.Errorf("wrong table column format: %s", colString)
		}

		name := parts[0]
		jsonPath := parts[1]

		result = append(result, &column{Name: name, JSONPath: jsonPath})
	}
	return result, nil
}

func configureWriter(w table.Writer, options *tableOptions) {
	if options.RowLength != 0 {
		w.SetAllowedRowLength(options.RowLength)
	} else {
		w.SetAllowedRowLength(getTermColumns())
	}

	if options.Style != nil {
		w.SetStyle(table.Style(*options.Style))
	} else {
		w.Style().Options.SeparateRows = true
		w.Style().Box = table.StyleBoxLight
	}

	if options.StyleCustomizations != nil {
		options.StyleCustomizations(w.Style())
	}
}

func getHeadersCount(options *tableOptions) int {
	headersCount := 0
	for _, col := range options.Columns {
		requiredColumHeaderLines := len(col.Parents)
		if col.Name != "" {
			requiredColumHeaderLines += 1
		}

		if headersCount < requiredColumHeaderLines {
			headersCount = requiredColumHeaderLines
		}
	}
	return headersCount
}

func buildTableHorizontally(writer table.Writer, val any, options *tableOptions) error {
	columnCount := len(options.Columns)
	headersCount := getHeadersCount(options)

	for headerIdx := 0; headerIdx < headersCount; headerIdx++ {
		headers := make(table.Row, columnCount)
		config := table.RowConfig{}
		for i, col := range options.Columns {
			if headerIdx < len(col.Parents) {
				headers[i] = col.Parents[headerIdx]
				config.AutoMerge = true
			} else if headerIdx == len(col.Parents) {
				headers[i] = col.Name
			} else {
				headers[i] = ""
			}
		}
		writer.AppendHeader(headers, config)
	}
	configureWriter(writer, options)

	rows := []table.Row{}
	configs := []table.ColumnConfig{}
	handleVal := func(rowIdx int, colIdx int, value any) error {
		if rowIdx >= len(rows) {
			newRow := make(table.Row, columnCount)
			for i := 0; i < columnCount; i++ {
				newRow[i] = ""
			}
			rows = append(rows, newRow)
		}

		switch value := value.(type) {
		case bool, *bool, int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64, float32, *float32, string, *string:
			rows[rowIdx][colIdx] = value
		case map[string]any, []any:
			subTable, err := buildSubTable(value, options)
			if err != nil {
				return err
			}
			rows[rowIdx][colIdx] = subTable
		default:
			// Marshall the value for easier reading when printing to console, even if inside of a table
			marshalled, err := json.Marshal(value)
			if err != nil {
				return err
			}
			rows[rowIdx][colIdx] = string(marshalled)
		}

		return nil
	}

	for colIdx, col := range options.Columns {
		result, err := utils.GetJsonPath(col.JSONPath, val)

		if err != nil {
			return err
		}

		if arr, ok := result.([]any); ok {
			for rowIdx, v := range arr {
				err := handleVal(rowIdx, colIdx, v)
				if err != nil {
					return err
				}
			}
		} else {
			err := handleVal(0, colIdx, result)
			if err != nil {
				return err
			}
		}

		if col.Config != nil {
			configs = append(configs, *col.tableColumnConfig())
		}
	}

	if len(configs) != 0 {
		writer.SetColumnConfigs(configs)
	}

	writer.AppendRows(rows)

	return nil
}

func buildSubTable(val any, parentOpts *tableOptions) (result string, err error) {
	subColumns, err := columnsFromAny(val, "$")
	if err != nil {
		return "", err
	}

	mapVal, isMap := val.(map[string]any)

	subOptions := *parentOpts
	subOptions.Columns = subColumns
	subOptions.BuildVertically = isMap && len(mapVal) > 1 && len(subOptions.Columns) > 1
	subOptions.StyleCustomizations = func(s *table.Style) {
		s.Options.DrawBorder = false
	}

	subTable, err := formatTable(val, &subOptions)
	if err != nil {
		return "", err
	}

	return renderWriterWithFormat(subTable, subOptions.Format), nil
}

func buildTableVertically(tw table.Writer, val any, options *tableOptions) error {
	configureWriter(tw, options)
	headersCount := getHeadersCount(options)

	if headersCount > 1 {
		configs := make([]table.ColumnConfig, headersCount-1)
		for headersIdx := 0; headersIdx < headersCount-1; headersIdx++ {
			configs[headersIdx] = table.ColumnConfig{Number: headersIdx + 1, AutoMerge: true}
		}
		tw.SetColumnConfigs(configs)
	}

	for _, col := range options.Columns {
		value, err := utils.GetJsonPath(col.JSONPath, val)
		if err != nil {
			return err
		}

		config := table.RowConfig{AutoMergeAlign: text.AlignLeft}
		row := make(table.Row, 0, headersCount+1)
		for headersIdx := 0; headersIdx < headersCount; headersIdx++ {
			if headersIdx < len(col.Parents) {
				row = append(row, col.Parents[headersIdx])
			} else {
				row = append(row, col.Name)
				config.AutoMerge = true
			}
		}

		switch value := value.(type) {
		case bool, *bool, int, *int, int8, *int8, int16, *int16, int32, *int32, int64, *int64, uint, *uint, uint8, *uint8, uint16, *uint16, uint32, *uint32, uint64, *uint64, float32, *float32, string, *string:
			row = append(row, value)
		case map[string]any, []any:
			subTable, err := buildSubTable(value, options)
			if err != nil {
				return err
			}
			row = append(row, subTable)
		default:
			// Marshall the value for easier reading when printing to console, even if inside of a table
			marshalled, err := json.Marshal(value)
			if err != nil {
				return err
			}
			row = append(row, string(marshalled))
		}
		tw.AppendRow(row, config)
	}

	return nil
}

func formatTable(val any, options *tableOptions) (table.Writer, error) {
	tw := table.NewWriter()
	var err error

	if options.BuildVertically {
		err = buildTableVertically(tw, val, options)
	} else {
		err = buildTableHorizontally(tw, val, options)
	}

	return tw, err
}

func formatTableWithOptions(val any, options *tableOptions) error {
	writer, err := formatTable(val, options)
	if err != nil {
		return err
	}
	fmt.Println(renderWriterWithFormat(writer, options.Format))
	return nil
}

func renderWriterWithFormat(writer table.Writer, format tableRenderFormat) string {
	switch format {
	case Ascii:
		return writer.Render()
	case HTML:
		return writer.RenderHTML()
	case CSV:
		return writer.RenderCSV()
	case Markdown:
		return writer.RenderMarkdown()
	}

	return writer.Render()
}

func (f *tableOutputFormatter) Format(val any, options string) (err error) {
	var columns []*column
	var buildVertically bool

	if options != "" {
		columns, err = columnsFromString(options)
	} else {
		columns, err = columnsFromAny(val, "$")

		if mapVal, ok := val.(map[string]any); ok {
			if len(mapVal) > 1 && len(columns) > 1 {
				buildVertically = true
			} else if len(mapVal) == 1 {
				for _, firstProp := range mapVal {
					if arr, ok := firstProp.([]any); ok {
						buildVertically = len(arr) == 1
					}
				}
			}
		}
	}

	if err != nil {
		return err
	}

	tableOptions := &tableOptions{Columns: columns, BuildVertically: buildVertically}

	if !buildVertically {
		tableOptions.Style = &noBorderStyle
	}

	return formatTableWithOptions(val, tableOptions)
}

func (*tableOutputFormatter) Description() string {
	return `Format as table using https://github.com/jedib0t/go-pretty/#table.` +
		` May be used as "table=COLNAME1:jsonpath-expression1,COLNAME2:jsonpath-expression2",` +
		` otherwise columns are automatically inferred from data layout.` +
		` For more complex specifications, see "table-file".`
}

func init() {
	outputFormatters["table"] = &tableOutputFormatter{}
}
