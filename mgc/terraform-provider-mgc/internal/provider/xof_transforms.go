package provider

import (
	"fmt"
	"strconv"

	"github.com/getkin/kin-openapi/openapi3"
	mgcSchemaPkg "magalu.cloud/core/schema"
	mgcSdk "magalu.cloud/sdk"
)

const (
	xOfAlternativeAsProp = "x-mgc-xOfAlternativeAsProp"
	xOfPromotionKey      = "x-mgc-xOfPromotionKey"
	xOfObject            = "x-mgc-xOfObject"
)

type xOfChild struct {
	s   *mgcSchemaPkg.Schema
	key string
}

// TODO: Find better clearer names for xOfObject vs xOfAlternative
func isSchemaXOfObject(s *mgcSchemaPkg.Schema) bool {
	if s == nil {
		return false
	}
	return s.Extensions[xOfObject] == true
}

func isSchemaXOfAlternative(s *mgcSchemaPkg.Schema) bool {
	return s.Extensions[xOfAlternativeAsProp] == true
}

func hasSchemaBeenPromoted(s *mgcSchemaPkg.Schema) bool {
	return s.Extensions[xOfPromotionKey] != nil
}

func findXOfChildrenCommonType(xOfChildren []xOfChild) string {
	var commonType string
	for _, xOf := range xOfChildren {
		if commonType == "" {
			commonType = xOf.s.Type
		} else if commonType != xOf.s.Type {
			return ""
		}
	}
	return commonType
}

func canPromoteXOfChildrenProps(parentCOW *mgcSchemaPkg.COWSchema, xOfChildren []xOfChild) bool {
	switch commonType := findXOfChildrenCommonType(xOfChildren); commonType {
	case "object":
		promotedProps := map[string]*mgcSchemaPkg.Schema{}
		for _, xOf := range xOfChildren {
			for xOfPropName, xOfPropRef := range xOf.s.Properties {
				var xOfProp = (*mgcSchemaPkg.Schema)(xOfPropRef.Value)
				var parentProp *mgcSchemaPkg.Schema

				if parentPropRef, ok := parentCOW.Properties()[xOfPropName]; ok {
					parentProp = (*mgcSchemaPkg.Schema)(parentPropRef.Value)
				} else if promotedProp, ok := promotedProps[xOfPropName]; ok {
					parentProp = promotedProp
				}

				if parentProp != nil {
					if !mgcSchemaPkg.CheckSimilarJsonSchemas(parentProp, xOfProp) {
						return false
					}
				} else {
					promotedProps[xOfPropName] = xOfProp
				}
			}
		}
		return true
	default:
		return false
	}
}

func promoteXOfChildrenPropsToParent(parentCOW *mgcSchemaPkg.COWSchema, xOfChildren []xOfChild) {
	for _, xOf := range xOfChildren {
		for xOfPropName, xOfPropSchemaRef := range xOf.s.Properties {
			if _, ok := parentCOW.PropertiesCOW().Get(xOfPropName); ok {
				continue
			}

			xOfPropSchema := (*mgcSchemaPkg.Schema)(xOfPropSchemaRef.Value)

			xOfCOW := mgcSchemaPkg.NewCOWSchema(xOfPropSchema)
			xOfCOW.ExtensionsCOW().Set(xOfPromotionKey, xOf.key)

			parentCOW.PropertiesCOW().Set(xOfPropName, openapi3.NewSchemaRef("", (*openapi3.Schema)(xOfCOW.Peek())))
		}
	}
}

func promoteXOfChildrenAsNewParentProps(parentCOW *mgcSchemaPkg.COWSchema, xOfChildren []xOfChild) {
	typeCounter := map[string]int{}
	for _, xOf := range xOfChildren {
		xOfCOW := mgcSchemaPkg.NewCOWSchema(xOf.s)
		xOfCOW.ExtensionsCOW().Set(xOfAlternativeAsProp, true)
		xOfCOW.ExtensionsCOW().Set(xOfPromotionKey, xOf.key)

		typeCount := typeCounter[xOfCOW.Type()]
		typeCounter[xOfCOW.Type()] = typeCount + 1

		promotedPropName := xOfCOW.Type() + strconv.Itoa(typeCount+1)
		parentCOW.PropertiesCOW().Set(promotedPropName, openapi3.NewSchemaRef("", (*openapi3.Schema)(xOfCOW.Peek())))
	}
	parentCOW.ExtensionsCOW().Set(xOfObject, true)
}

func promoteXOfChildren(parentCOW *mgcSchemaPkg.COWSchema, xOfChildren []xOfChild) error {
	parentCOW.SetType("object")
	if canPromoteXOfChildrenProps(parentCOW, xOfChildren) {
		promoteXOfChildrenPropsToParent(parentCOW, xOfChildren)
	} else {
		if len(parentCOW.Properties()) > 0 {
			return fmt.Errorf(
				"trying to promote xOf children as new parent props but parent already has props of its own. Parent %+v Children %+v",
				parentCOW.Peek(),
				xOfChildren,
			)
		}
		promoteXOfChildrenAsNewParentProps(parentCOW, xOfChildren)
	}
	return nil
}

func promoteXOfsToProps(s *mgcSdk.Schema, key string) (*mgcSdk.Schema, error) {
	var xOfs []xOfChild

	_, err := mgcSchemaPkg.ForEachXOf(s, func(xOfS mgcSchemaPkg.XOfChildSchema) (run bool, err error) {
		childKey := key + xOfS.Field + strconv.Itoa(xOfS.Index)

		transformedChild, err := promoteXOfsToProps(xOfS.Schema, childKey)
		xOfs = append(xOfs, xOfChild{s: transformedChild, key: childKey})
		return true, err
	})
	if err != nil {
		return nil, err
	}

	if len(xOfs) == 0 {
		return s, nil
	}

	transformedCOW := mgcSchemaPkg.NewCOWSchema(s)
	err = promoteXOfChildren(transformedCOW, xOfs)
	if err != nil {
		return nil, err
	}

	transformedCOW.SetAnyOf(nil)
	transformedCOW.SetOneOf(nil)

	return transformedCOW.Peek(), nil
}

func getXOfObjectSchemaTransformed(s *mgcSdk.Schema) (*mgcSdk.Schema, error) {
	resultCOW := mgcSchemaPkg.NewCOWSchema(s)

	for propName, propSchemaRef := range resultCOW.Properties() {
		propSchema := (*mgcSchemaPkg.Schema)(propSchemaRef.Value)
		propTransformed, err := getXOfObjectSchemaTransformed(propSchema)
		if err != nil {
			return nil, err
		}
		resultCOW.PropertiesCOW().Set(propName, mgcSchemaPkg.NewSchemaRef("", propTransformed))
	}

	if notRef := resultCOW.Not(); notRef != nil {
		not := (*mgcSchemaPkg.Schema)(notRef.Value)
		notTransformed, err := getXOfObjectSchemaTransformed(not)
		if err != nil {
			return nil, err
		}
		resultCOW.NotCOW().SetValue(notTransformed)
	}

	if itemsRef := resultCOW.Items(); itemsRef != nil {
		items := (*mgcSchemaPkg.Schema)(itemsRef.Value)
		itemsTransformed, err := getXOfObjectSchemaTransformed(items)
		if err != nil {
			return nil, err
		}
		resultCOW.ItemsCOW().SetValue(itemsTransformed)
	}

	if additionalPropsData := resultCOW.AdditionalProperties(); additionalPropsData.Schema != nil {
		additionalPropsRef := additionalPropsData.Schema
		additionalProps := (*mgcSchemaPkg.Schema)(additionalPropsRef.Value)
		additionalPropsTransformed, err := getXOfObjectSchemaTransformed(additionalProps)
		if err != nil {
			return nil, err
		}
		resultCOW.SetAdditionalProperties(
			openapi3.AdditionalProperties{
				Has:    additionalPropsData.Has,
				Schema: mgcSchemaPkg.NewSchemaRef(additionalPropsRef.Ref, additionalPropsTransformed),
			},
		)
	}

	result, _ := resultCOW.Release()
	return promoteXOfsToProps(result, "")
}
