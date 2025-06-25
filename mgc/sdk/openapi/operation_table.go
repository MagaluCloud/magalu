package openapi

import (
	"math"
	"regexp"
	"strings"

	"slices"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/stoewer/go-strcase"
)

// If slice is empty, returns a zero value
func getSecondToLastOrLastElem[T any](arr []T) (result T) {
	length := len(arr)
	if length == 1 {
		return arr[0]
	}
	if length > 1 {
		return arr[length-2]
	}

	return
}

// If slice is empty, returns a zero value
func getLastElem[T any](arr []T) (result T) {
	length := len(arr)
	if length == 0 {
		return
	}

	return arr[length-1]
}

type operationDesc struct {
	path    *openapi3.PathItem
	op      *openapi3.Operation
	method  string
	pathKey string
}

// Returns the simplest way to identify this operation. When it only has the HTTP Method in its name
// parts it returns that. Otherwise, it returns the immediately preceding name part, which is the last
// path entry
func (e *operationTableEntry) simpleNameKey() string {
	return getSecondToLastOrLastElem(e.name)
}

// Returns the HTTP Method + first name part (if any) + last name part (if any)
// If the entry only has the HTTP Method in the name, this won't be enought to disambiguate
// between its siblings. Otherwise, it will be unique to its siblings.
func (e *operationTableEntry) fullNameKey() string {
	// Last entry in name is always HTTP method
	switch length := len(e.name); length {
	case 0:
		// Should never happen
		return ""
	case 1:
		return e.name[0]
	case 2:
		return e.name[1] + "-" + e.name[0]
	default:
		// Considering two operations in the same table with more than 2 name parts, it is always safe to assume that the
		// first part in the name is different from all others, so first part + last part + http method
		// is enough to serve as an "ID".
		return e.name[length-1] + "-" + e.name[0] + "-" + e.name[length-2]
	}
}

type operationTableEntry struct {
	name      []string
	variables []string
	desc      *operationDesc
	// Key is the identifier of the entry in relation to its siblings. It may be the result of 'simpleNameKey'
	// it may be the result of 'fullNameKey' or it may be the result of 'fullNameKey' +
	// unique variables at the end. This is only set after 'table.finalizeEntryKeys()' is called.
	key string
}

// Operation Table is used to compile a structure of operations with a hierarchy based on the path entries
// of OAPI operations. This allows for sub-tables with their own set of operations, even though all of
// them are in the same overall resource (OAPI Tag)
type operationTable struct {
	name            string
	childTables     []*operationTable
	childOperations []*operationTableEntry
}

func (t *operationTable) findSibling(name []string) (siblingIdx int, sibling *operationTableEntry) {
	for i, childEntry := range t.childOperations {
		// If another entry in the table starts with the same name part,
		// a subtable for both of them should be created. Don't create subtable if the sibling
		// is already at the maximum depth possible (only has HTTP method left)
		if childEntry.name[0] == name[0] && len(childEntry.name) > 1 {
			return i, childEntry
		}
	}
	return
}

// Consider that all of the entries passed as parameters conflict with each other and set their keys
// to the fullest possible identifier, including 'entry.fullNameKey()' and all unique variables.
func (t *operationTable) setUniqueFullKeys(entries ...*operationTableEntry) {
	maxVarLength := math.MinInt
	for _, entry := range entries {
		entry.key = entry.fullNameKey()
		if varLength := len(entry.variables); varLength > maxVarLength {
			maxVarLength = varLength
		}
	}

	for i := 0; i < maxVarLength; i++ {
		commonVariable := ""
		isCommonVariable := true
		for _, entry := range entries {
			// Immediately break if variable isn't present in one if the entries, thus,
			// it's not a common variable
			if i > len(entry.variables)-1 {
				isCommonVariable = false
				break
			}

			if commonVariable == "" {
				commonVariable = entry.variables[i]
				continue
			}

			if commonVariable != entry.variables[i] {
				isCommonVariable = false
				break
			}
		}

		if isCommonVariable {
			continue
		}

		for _, entry := range entries {
			if i < len(entry.variables) {
				entry.key += "-" + entry.variables[i]
			}
		}
	}
}

func (t *operationTable) add(name, variables []string, desc *operationDesc) {
	// Should never happen
	if len(name) == 0 {
		return
	}

	// If a child table already exists for the current name, add it to that one
	for _, childTable := range t.childTables {
		if childTable.name == name[0] {
			childTable.add(name[1:], variables, desc)
			return
		}
	}

	// If a sibling is present (another entry that starts with the same first name entry), create a
	// subtable for them and add it to that one
	if siblingIdx, sibling := t.findSibling(name); sibling != nil {
		childTable := &operationTable{name: name[0]}
		childTable.add(sibling.name[1:], sibling.variables, sibling.desc)
		childTable.add(name[1:], variables, desc)
		t.childTables = append(t.childTables, childTable)

		// Remove sibiling from child operations, as it's now in subtable
		t.childOperations = append(t.childOperations[:siblingIdx], t.childOperations[siblingIdx+1:]...)
		return
	}

	// Otherwise, just add the entry normally to the current table
	entry := &operationTableEntry{name: name, variables: variables, desc: desc}
	t.childOperations = append(t.childOperations, entry)
}

// Calling this will override the current table's child tables and operations. All of them will
// be set to the child tables' values
func (t *operationTable) promoteChildTable(childTable *operationTable) {
	// Merge child table's childTables and childOperations into the current table
	t.childTables = childTable.childTables
	t.childOperations = childTable.childOperations

	// Rename the current table to include the child name
	t.name = t.name + "-" + childTable.name
}

// Simplify table implements the following rules:
//   - If a table doesn't have any operations of its own and only has one child table, the child table
//     will be merged into the parent table
//   - If a table has only one child operation, the name of the operation will be shortened to include
//     only the HTTP Method
func (t *operationTable) simplify() {
	// Recursively simplify child tables
	for _, childTable := range t.childTables {
		childTable.simplify()
	}

	if len(t.childOperations) == 0 && len(t.childTables) == 1 {
		childTable := t.childTables[0]
		t.promoteChildTable(childTable)
	}

	if len(t.childOperations) == 1 {
		entry := t.childOperations[0]
		entry.name = []string{getLastElem(entry.name)}
	}
}

// Common path endings that need the full name ('list-all', 'delete-all'...). Increment array
// as needed
var namesToBePrefixedWithHTTPMethod = []string{"all", "default"}
var httpMethodsThatEnforceFullName = []string{"delete"}

func needsFullNameKey(name []string) bool {
	return slices.Contains(namesToBePrefixedWithHTTPMethod, getSecondToLastOrLastElem(name)) || slices.Contains(httpMethodsThatEnforceFullName, getLastElem(name))
}

// Set all operation keys, disambiguating conflicting entries if needed. Conflicting entries will
// have a key with any variables that aren't common to all conflicting siblings
func (t *operationTable) finalizeEntryKeys() {
	bySimpleKey := map[string][]*operationTableEntry{}
	for _, childOperation := range t.childOperations {
		simpleKey := childOperation.simpleNameKey()
		bySimpleKey[simpleKey] = append(bySimpleKey[simpleKey], childOperation)
	}

	for simpleKey, conflictingEntries := range bySimpleKey {
		if len(conflictingEntries) > 1 {
			t.setUniqueFullKeys(conflictingEntries...)
		} else {
			entry := conflictingEntries[0]
			if needsFullNameKey(entry.name) {
				entry.key = entry.fullNameKey()
			} else {
				entry.key = simpleKey
			}
		}
	}

	for _, childTable := range t.childTables {
		childTable.finalizeEntryKeys()
	}
}

var openAPIPathArgRegex = regexp.MustCompile("[{](?P<name>[^}]+)[}]")

func getPathEntry(pathEntry string) (string, bool) {
	match := openAPIPathArgRegex.FindStringSubmatch(pathEntry)
	if len(match) > 0 {
		for i, substr := range match {
			if openAPIPathArgRegex.SubexpNames()[i] == "name" {
				return substr, true
			}
		}
	}

	return pathEntry, false
}

func renameHttpMethod(httpMethod string, endsWithVariable bool) string {
	switch httpMethod {
	case "post":
		return "create"
	case "put":
		return "replace"
	case "patch":
		return "update"
	case "get":
		// only consider "get" if ends with varable, mid-path are still list, ex:
		// GET:  /resource/{id}
		// LIST: /{containerId}/resource
		// GET:  /{containerId}/resource/{id}
		if endsWithVariable {
			return "get"
		}
		return "list"
	}

	return httpMethod
}

func isVersion(value string) bool {
	versionRegex := `^v\d+(?:[a-z]+\d+)?$`
	regex := regexp.MustCompile(versionRegex)
	return regex.MatchString(value)
}

func getOperationNameAndVariables(httpMethod, pathName string) (pathEntries []string, variables []string) {
	endsWithVariable := false
	for _, pathEntry := range strings.Split(pathName, "/") {

		pathEntry = strings.ReplaceAll(pathEntry, "_", "-")

		if pathEntry == "" {
			continue
		}

		if isVersion(pathEntry) {
			continue
		}

		if variable, isVariable := getPathEntry(pathEntry); isVariable {
			variables = append(variables, strcase.KebabCase(variable))
			endsWithVariable = true
		} else {
			pathEntries = append(pathEntries, strings.Split(strcase.KebabCase(pathEntry), "-")...)
			endsWithVariable = false
		}
	}

	pathEntries = append(pathEntries, renameHttpMethod(httpMethod, endsWithVariable))

	return
}

func newOperationTable(name string, descs []*operationDesc) *operationTable {
	table := &operationTable{name: name}
	for _, desc := range descs {
		descName, descVariables := getOperationNameAndVariables(desc.method, desc.pathKey)
		table.add(descName, descVariables, desc)
	}
	table.simplify()
	table.finalizeEntryKeys()
	return table
}
