package blueprint

import (
	"errors"
	"fmt"

	"magalu.cloud/core"
	schemaPkg "magalu.cloud/core/schema"
	"magalu.cloud/core/utils"
)

type executorSpec struct {
	ParametersSchema *schemaPkg.SchemaRef `json:"parametersSchema"`
	ConfigsSchema    *schemaPkg.SchemaRef `json:"configsSchema"`
	ResultSchema     *schemaPkg.SchemaRef `json:"resultSchema"`

	Links   map[string]*linkSpec `json:"links,omitempty"`
	linkers map[string]core.Linker

	Related          map[string]core.RefPath `json:"related,omitempty"`
	relatedExecutors map[string]core.Executor

	PositionalArgs []string `json:"positionalArgs,omitempty"`
	HiddenFlags    []string `json:"hiddenFlags,omitempty"`

	// Resolved schemas, populated by resolve():

	parametersSchema *schemaPkg.Schema
	configsSchema    *schemaPkg.Schema
	resultSchema     *schemaPkg.Schema

	// Since parameters and config are objects, this allows a convenience to declare them in the YAML.
	// all keys turn to be required properties if their schema doesn't provide a default value
	// TODO: do we want to keep this? It's helpful or just add confusion?

	Parameters map[string]*schemaPkg.SchemaRef `json:"parameters"`
	Configs    map[string]*schemaPkg.SchemaRef `json:"configs"`

	Result string `json:"result"`

	// Execution steps, processed in order
	Steps []*executeStep `json:"steps"`

	// core.Executor extensions:

	Confirm         string                      `json:"confirm,omitempty"`
	WaitTermination *core.WaitTerminationConfig `json:"waitTermination,omitempty"`
	OutputFlag      string                      `json:"outputFlag,omitempty"`
	PromptInput     *PromptInputSpec            `json:"promptInput,omitempty"`
}

type PromptInputSpec struct {
	MessageTemplate      string `json:"message,omitempty"`
	ConfirmValueTemplate string `json:"confirmValue,omitempty"`
}

func (e *executorSpec) isEmpty() bool {
	return len(e.Steps) == 0
}

func isSchemaRefEmpty(schemaRef *schemaPkg.SchemaRef) bool {
	if schemaRef == nil {
		return true
	}

	return schemaRef.Ref == "" && schemaRef.Value == nil
}

func checkObjectSchema(schema *schemaPkg.Schema) bool {
	return schema.Type == "object"
}

func checkObjectSchemaRef(schemaRef *schemaPkg.SchemaRef) bool {
	if schemaRef == nil || schemaRef.Value == nil {
		return true // for unknown, assume it is an object
	}
	return checkObjectSchema((*schemaPkg.Schema)(schemaRef.Value))
}

func (e *executorSpec) validate() (err error) {
	if len(e.Parameters) == 0 && isSchemaRefEmpty(e.ParametersSchema) {
		return errors.New("missing parameters or parametersSchema")
	}
	if !checkObjectSchemaRef(e.ParametersSchema) {
		return errors.New("parametersSchema is not an object")
	}

	if len(e.Configs) == 0 && isSchemaRefEmpty(e.ConfigsSchema) {
		return errors.New("missing configs or configsSchema")
	}
	if !checkObjectSchemaRef(e.ConfigsSchema) {
		return errors.New("configsSchema is not an object")
	}

	if isSchemaRefEmpty(e.ResultSchema) {
		return errors.New("missing resultSchema")
	}

	if e.Result != "" {
		_, err = utils.NewJsonPath(e.Result)
		if err != nil {
			return &core.ChainedError{Name: "result", Err: err}
		}
	}

	if e.Confirm != "" {
		_, err = utils.NewTemplate(e.Confirm)
		if err != nil {
			return &core.ChainedError{Name: "confirm", Err: err}
		}
	}

	if e.WaitTermination != nil {
		_, err = e.WaitTermination.Build(nil, nil)
		if err != nil {
			return &core.ChainedError{Name: "waitTermination", Err: err}
		}
	}

	if e.isEmpty() {
		return errors.New("missing steps")
	}
	for i, exec := range e.Steps {
		err = exec.validate()
		if err != nil {
			return &core.ChainedError{
				Name: "steps",
				Err:  &core.ChainedError{Name: fmt.Sprintf("%d(id=%q)", i, exec.Id), Err: err},
			}
		}
	}

	for k, v := range e.Links {
		v.name = k
		err = v.validate()
		if err != nil {
			return &core.ChainedError{
				Name: "links",
				Err:  &core.ChainedError{Name: k, Err: err},
			}
		}
	}

	return nil
}

func (e *executorSpec) resolve(refResolver *core.BoundRefPathResolver, exec core.Executor) (err error) {
	if e.parametersSchema == nil {
		e.parametersSchema, err = getSchema(refResolver, e.ParametersSchema, e.Parameters)
		if err != nil {
			return fmt.Errorf("parametersSchema: %w", err)
		}
		if !checkObjectSchema(e.parametersSchema) {
			return errors.New("parametersSchema is not an object")
		}
	}

	if e.configsSchema == nil {
		e.configsSchema, err = getSchema(refResolver, e.ConfigsSchema, e.Configs)
		if err != nil {
			return fmt.Errorf("configsSchema: %w", err)
		}
		if !checkObjectSchema(e.configsSchema) {
			return errors.New("configsSchema is not an object")
		}
	}

	if e.resultSchema == nil {
		e.resultSchema, err = getSchema(refResolver, e.ResultSchema, nil)
		if err != nil {
			return fmt.Errorf("resultSchema: %w", err)
		}
	}

	for i, step := range e.Steps {
		if step.Id == "" {
			step.Id = fmt.Sprint(i)
		}
		err := step.resolve(refResolver)
		if err != nil {
			return fmt.Errorf("invalid step %d(id=%q): %w", i, step.Id, err)
		}
	}

	if e.relatedExecutors == nil {
		e.relatedExecutors = map[string]core.Executor{}
		for k, p := range e.Related {
			e.relatedExecutors[k], err = core.ResolveExecutorPath(refResolver, p)
			if err != nil {
				return fmt.Errorf("related %q: %w", k, err)
			}
		}
	}

	if e.linkers == nil {
		e.linkers = map[string]core.Linker{}
		for k, v := range e.Links {
			err = v.resolve(refResolver)
			if err != nil {
				return fmt.Errorf("invalid link %q: %w", k, err)
			}
			e.linkers[k] = &linker{spec: v, owner: exec}
		}
	}

	return nil
}

func createSchemaFromMap(refResolver *core.BoundRefPathResolver, m map[string]*schemaPkg.SchemaRef) (result *core.Schema, err error) {
	properties := make(map[string]*schemaPkg.Schema, len(m))
	required := make([]string, 0, len(m))
	for k, v := range m {
		prop, err := core.ResolveSchemaRef(refResolver, v)
		if err != nil {
			return nil, err
		}
		properties[k] = prop
		if prop.Default == nil {
			required = append(required, k)
		}
	}

	result = schemaPkg.NewObjectSchema(properties, required)
	return
}

func getSchema(refResolver *core.BoundRefPathResolver, schemaRef *schemaPkg.SchemaRef, m map[string]*schemaPkg.SchemaRef) (result *core.Schema, err error) {
	if len(m) != 0 {
		return createSchemaFromMap(refResolver, m)
	}

	return core.ResolveSchemaRef(refResolver, schemaRef)
}
