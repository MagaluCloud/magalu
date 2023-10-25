import os
from collections import defaultdict
from typing import (
    cast,
    Callable,
    Dict,
    Any,
    List,
    Literal,
    Mapping,
    NotRequired,
    Sequence,
    Set,
    TypeAlias,
    TypedDict,
    Union,
    NamedTuple,
    Optional,
)
import argparse
import yaml
import jsonschema

OAPIStats = Dict[str, Any]

ArgumentLocation: TypeAlias = Literal["query", "header", "path", "cookie"]
ArgumentStyle: TypeAlias = Literal[
    "matrix",
    "label",
    "form",
    "simple",
    "spaceDelimited",
    "pipeDelimited",
    "deepObject",
]

HttpMethod: TypeAlias = Literal["get", "put", "post", "delete", "patch"]

JSONSchema: TypeAlias = Mapping[str, Any]  # TODO


class OAPIExample(NamedTuple):
    summary: str | None
    description: str | None
    value: Any


class OAPIArgumentSchema(NamedTuple):
    location: ArgumentLocation
    required: bool
    deprecated: bool
    description: str | None
    schema: JSONSchema
    examples: Sequence[OAPIExample]
    style: ArgumentStyle
    explode: bool
    allow_empty_value: bool
    allow_reserved: bool


class OAPIHeaderSchema(NamedTuple):
    required: bool
    deprecated: bool
    description: str | None
    schema: JSONSchema
    style: ArgumentStyle
    examples: Sequence[OAPIExample]
    explode: bool


class OAPILinkSchema(NamedTuple):
    # TODO: once actions are created, materialize with action: Action
    operation_id: str | None
    operation_ref: str | None
    parameters: Mapping[str, Any]
    request_body: Any
    description: str | None
    server: str | None


class OAPIContentSchema(NamedTuple):
    schema: JSONSchema
    examples: Sequence[OAPIExample]


class OAPIResponseSchema(NamedTuple):
    description: str
    headers: Mapping[str, OAPIHeaderSchema]
    content: Mapping[str, OAPIContentSchema]
    links: Mapping[str, OAPILinkSchema]


class OAPIRequestSchema(NamedTuple):
    description: str | None
    required: bool
    content: Mapping[str, OAPIContentSchema]


class OAPISecurityRequirement(NamedTuple):
    name: str
    scopes: Sequence[str]


# --- OAPI Specification (input)


OAPIReferenceObject = TypedDict(
    "OAPIReferenceObject",
    {
        "$ref": str,
        "summary": NotRequired[str],
        "description": NotRequired[str],
    },
)


class OAPIServerVariableObject(TypedDict):
    default: str
    description: NotRequired[str]
    enum: NotRequired[Sequence[str]]


class OAPIServerObject(TypedDict):
    url: str
    description: NotRequired[str]
    variables: NotRequired[Mapping[str, OAPIServerVariableObject]]


OAPITagObject = Dict[str, Any]


class OAPIExampleObject(TypedDict):
    summary: NotRequired[str]
    description: NotRequired[str]
    value: NotRequired[Any]
    externalValue: NotRequired[str]


class OAPIHeaderObject(TypedDict):
    description: NotRequired[str]
    required: NotRequired[bool]
    deprecated: NotRequired[bool]
    style: NotRequired[ArgumentStyle]
    explode: NotRequired[bool]
    schema: NotRequired[JSONSchema]
    example: NotRequired[Any]
    examples: NotRequired[Mapping[str, OAPIExampleObject | OAPIReferenceObject]]
    content: NotRequired[Mapping[str, "OAPIMediaTypeObject"]]


class OAPIEncodingObject(TypedDict):
    contentType: NotRequired[str]
    headers: NotRequired[Mapping[str, OAPIHeaderObject | OAPIReferenceObject]]
    style: NotRequired[str]
    explode: NotRequired[bool]
    allowReserved: NotRequired[bool]


class OAPIMediaTypeObject(TypedDict):
    schema: NotRequired[JSONSchema]
    example: NotRequired[Any]
    examples: NotRequired[Mapping[str, OAPIExampleObject | OAPIReferenceObject]]
    encoding: NotRequired[Mapping[str, OAPIEncodingObject]]


OAPIParameterObject = TypedDict(
    "OAPIParameterObject",
    {
        "name": str,
        "in": ArgumentLocation,  # NOTE: reserved keyword :-(
        "description": NotRequired[str],
        "required": NotRequired[bool],
        "deprecated": NotRequired[bool],
        "allowEmptyValue": NotRequired[bool],
        "style": NotRequired[str],
        "explode": NotRequired[bool],
        "allowReserved": NotRequired[bool],
        "schema": NotRequired[JSONSchema],
        "example": NotRequired[Any],
        "examples": NotRequired[Mapping[str, OAPIExampleObject]],
        "content": NotRequired[Mapping[str, OAPIMediaTypeObject]],
    },
)


class OAPIExternalDocumentationObject(TypedDict):
    url: str
    description: NotRequired[str]


class OAPIRequestBodyObject(TypedDict):
    description: NotRequired[str]
    content: Mapping[str, OAPIMediaTypeObject]
    required: NotRequired[bool]


class OAPILinkObject(TypedDict):
    operationRef: NotRequired[str]
    operationId: NotRequired[str]
    parameters: NotRequired[Mapping[str, Any]]
    requestBody: NotRequired[Any]
    description: NotRequired[str]
    server: NotRequired[OAPIServerObject]


class OAPIResponseObject(TypedDict):
    description: str
    headers: NotRequired[Mapping[str, OAPIHeaderObject | OAPIReferenceObject]]
    content: NotRequired[Mapping[str, OAPIMediaTypeObject]]
    links: NotRequired[Mapping[str, OAPILinkObject | OAPIReferenceObject]]


OAPIResponsesObject: TypeAlias = Mapping[str, OAPIResponseObject | OAPIReferenceObject]

OAPISecurityRequirementObject: TypeAlias = Mapping[str, Sequence[str]]
OAPICallbackObject: TypeAlias = Mapping[
    str, Union["OAPIPathItemObject", OAPIReferenceObject]
]


class OAPIOperationObject(TypedDict):
    tags: NotRequired[Sequence[str]]
    summary: NotRequired[str]
    description: NotRequired[str]
    externalDocs: NotRequired[OAPIExternalDocumentationObject]
    operationId: NotRequired[str]
    parameters: NotRequired[Sequence[OAPIParameterObject | OAPIReferenceObject]]
    requestBody: NotRequired[OAPIRequestBodyObject | OAPIReferenceObject]
    responses: NotRequired[OAPIResponsesObject]
    callbacks: NotRequired[Mapping[str, OAPICallbackObject]]
    deprecated: NotRequired[bool]
    security: NotRequired[Sequence[OAPISecurityRequirementObject]]
    servers: NotRequired[Sequence[OAPIServerObject]]


OAPIPathItemObject = TypedDict(
    "OAPIPathItemObject",
    {
        "$ref": NotRequired[str],
        "summary": NotRequired[str],
        "description": NotRequired[str],
        "get": NotRequired[OAPIOperationObject],
        "put": NotRequired[OAPIOperationObject],
        "post": NotRequired[OAPIOperationObject],
        "delete": NotRequired[OAPIOperationObject],
        "options": NotRequired[OAPIOperationObject],
        "head": NotRequired[OAPIOperationObject],
        "patch": NotRequired[OAPIOperationObject],
        "trace": NotRequired[OAPIOperationObject],
        "servers": NotRequired[Sequence[OAPIServerObject]],
        "parameters": NotRequired[Sequence[OAPIParameterObject | OAPIReferenceObject]],
    },
)


class OAPIInfoObject(TypedDict):
    title: str
    version: str
    summary: NotRequired[str]
    description: NotRequired[str]


OAPISecuritySchemeApiKeyObject = TypedDict(
    "OAPISecuritySchemeApiKeyObject",
    {
        "type": Literal["apiKey"],
        "description": NotRequired[str],
        "name": str,
        "in": str,  # NOTE: reserved keyword :-(
    },
)


class OAPISecuritySchemeHttpObject(TypedDict):
    type: Literal["http"]  # noqa A003
    description: NotRequired[str]
    scheme: str
    bearerFormat: str


class OAPIOAuthFlowObject(TypedDict):
    authorizationUrl: str
    tokenUrl: str
    refreshUrl: NotRequired[str]
    scopes: Mapping[str, str]


class OAPIOAuthFlowsObject(TypedDict):
    implicit: NotRequired[OAPIOAuthFlowObject]
    password: NotRequired[OAPIOAuthFlowObject]
    clientCredentials: NotRequired[OAPIOAuthFlowObject]
    authorizationCode: NotRequired[OAPIOAuthFlowObject]


class OAPISecuritySchemeOAuth2Object(TypedDict):
    type: Literal["oauth2"]  # noqa A003
    description: NotRequired[str]
    flows: OAPIOAuthFlowsObject


class OAPISecuritySchemeOpenIdConnectObject(TypedDict):
    type: Literal["openIdConnect"]  # noqa A003
    description: NotRequired[str]
    openIdConnectUrl: str


OAPISecuritySchemeObject: TypeAlias = (
    OAPISecuritySchemeApiKeyObject
    | OAPISecuritySchemeHttpObject
    | OAPISecuritySchemeOAuth2Object
    | OAPISecuritySchemeOpenIdConnectObject
)


class OAPIComponentsObject(TypedDict):
    schemas: NotRequired[Mapping[str, JSONSchema]]
    responses: NotRequired[Mapping[str, OAPIResponseObject | OAPIReferenceObject]]
    parameters: NotRequired[Mapping[str, OAPIParameterObject | OAPIReferenceObject]]
    examples: NotRequired[Mapping[str, OAPIExampleObject | OAPIReferenceObject]]
    requestBodies: NotRequired[
        Mapping[str, OAPIRequestBodyObject | OAPIReferenceObject]
    ]
    headers: NotRequired[Mapping[str, OAPIHeaderObject | OAPIReferenceObject]]
    securitySchemes: NotRequired[
        Mapping[str, OAPISecuritySchemeObject | OAPIReferenceObject]
    ]
    links: NotRequired[Mapping[str, OAPILinkObject | OAPIReferenceObject]]
    callbacks: NotRequired[Mapping[str, OAPICallbackObject | OAPIReferenceObject]]
    pathItems: NotRequired[Mapping[str, OAPIPathItemObject | OAPIReferenceObject]]


class OAPIObject(TypedDict):
    openapi: str
    info: OAPIInfoObject
    servers: NotRequired[Sequence[OAPIServerObject]]
    paths: NotRequired[Mapping[str, OAPIPathItemObject]]
    components: NotRequired[OAPIComponentsObject]
    security: NotRequired[Sequence[OAPISecurityRequirementObject]]
    tags: NotRequired[Sequence[OAPITagObject]]
    externalDocs: NotRequired[OAPIExternalDocumentationObject]


class OAPI(NamedTuple):
    path: str
    name: str
    obj: OAPIObject
    ref_resolver: jsonschema.RefResolver

    def resolve(self, ref: str) -> Any:
        return self.ref_resolver.resolve(ref)[1]


class OAPIOperationInfo(NamedTuple):
    path: str
    method: str
    op: OAPIOperationObject

    def key(self) -> str:
        return self.method.upper() + " " + self.path


class OAPITagInfo(NamedTuple):
    name: str
    description: str
    extensions: JSONSchema


class OAPIResource(NamedTuple):
    name: str
    operations: List[OAPIOperationInfo]
    tag: Optional[OAPITagInfo]


class ResponseContext(NamedTuple):
    path: str
    method: str
    code: str


# This is used to fix list indentations, as Pyyaml doesn't indent them :/
class YamlDumper(yaml.Dumper):
    def increase_indent(self, flow=False, indentless=False):
        return super(YamlDumper, self).increase_indent(flow, False)


OPERATION_KEYS = ["get", "put", "post", "delete", "options", "head", "patch", "trace"]
CRUD_OP_KEYS = ["get", "put", "patch", "post", "delete"]
# Patch and Post are optional, as they can be mimicked with a Delete->Create op
REQUIRED_CRUD_OP_KEYS = ["get", "post", "delete"]

COMPUTED_VARIABLES = "computed-variables"
MISSING_CRUD = "missing-crud"
MISSING_OPERATION_ID = "missing-operation-id"
XOF_SCHEMA_SINGLE_ITEM = "xof-single-item"
MIXED_ENUMS = "mixed-enums"
NON_JSON_REQUESTS = "non-json-requests"
NON_JSON_RESPONSES = "non-json-responses"
NON_SNAKECASE_VALUES = "non-snakecase-values"
TAGLESS_OPERATIONS = "tagless-operations"
UNUSED_COMPONENTS = "unused-components"


class Filterer:
    filters: List[str]
    filter_out: List[str]

    def should_include(self, key: str) -> bool:
        if self.filters and key not in self.filters:
            return False

        if self.filter_out and key in self.filter_out:
            return False

        return True


filterer = Filterer()


def load_yaml(path: str) -> OAPIObject:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def load_oapi(path: str) -> Optional[OAPI]:
    f = os.path.basename(path)
    name, ext = os.path.splitext(f)
    if name == "index" or ext != ".yaml":
        print("ignored file:", f)
        return None

    obj = load_yaml(path)
    as_dict = cast(Dict[str, Any], obj)
    ref_resolver = jsonschema.RefResolver(path, as_dict)
    return OAPI(name=name, path=path, obj=obj, ref_resolver=ref_resolver)


def load_oapis(dir_or_path: str, ignore_disabled: bool) -> List[OAPI]:
    if os.path.isdir(dir_or_path):
        d = dir_or_path
        result = []
        for f in os.listdir(d):
            path = os.path.join(d, f)
            # Prevent file loading, save resources
            if ignore_disabled and ".disabled" in path:
                continue

            oapi = load_oapi(path)
            if not oapi:
                continue

            result.append(oapi)

        return result
    else:
        p = dir_or_path
        oapi = load_oapi(p)
        if oapi:
            return [oapi]
        else:
            return []


def is_tag_crud(tag: Optional[OAPITagInfo]) -> bool:
    if not tag:
        return False

    if tag.extensions.get("x-mgc-hidden"):
        return False

    return True


def single_item_xofs(schema: JSONSchema) -> List[str]:
    def filterFn(xof: str) -> bool:
        return len(schema.get(xof, [])) == 1

    return list(filter(filterFn, ["allOf", "anyOf", "oneOf"]))


def is_snake_case(s: str) -> bool:
    if any(char.isupper() for char in s):
        return False

    return "-" not in s


components_usage: dict = defaultdict(dict)


def is_ref(obj_or_ref):
    return "$ref" in obj_or_ref


def get(obj_or_ref: Any | OAPIReferenceObject, resolve: Callable[[str], Any]) -> Any:
    if is_ref(obj_or_ref):
        return resolve(obj_or_ref["$ref"])
    else:
        return obj_or_ref


def get_schema_field_names(
    schema: JSONSchema, resolve: Callable[[str], Any]
) -> Set[str]:
    result = set()
    t = schema.get("type")
    if t == "object":
        for pn, p in schema.get("properties", {}).items():
            ps: JSONSchema = get(p, resolve)
            pt = ps.get("type")
            if pt == "object":
                # Flatten out all sub fields as if top-level
                result.update(get_schema_field_names(ps, resolve))
            else:
                result.add(pn)

    elif schema.get("title") is not None:
        result.add(schema["title"])

    return result


def traverse_all_subschemas(
    path: str,
    schema_or_ref: JSONSchema | OAPIReferenceObject,
    visit: Callable[[str, JSONSchema], None],
    resolve: Callable[[str], Any] | None = None,
):
    if not is_ref(schema_or_ref):
        schema: JSONSchema = schema_or_ref
    elif resolve:  # if not resolve do not follow ref
        schema = get(schema_or_ref, resolve)
    else:
        return

    visit(path, schema)

    for pn, p_or_ref in schema.get("properties", {}).items():
        traverse_all_subschemas(path + "." + pn, p_or_ref, visit, resolve)

    if schema.get("additionalProperties"):
        s = schema["additionalProperties"]
        traverse_all_subschemas(path + "|additional_properties|", s, visit, resolve)

    xOfs = schema.get("oneOf", []) + schema.get("anyOf", []) + schema.get("allOf", [])
    for s_or_ref in xOfs:
        traverse_all_subschemas(path, s_or_ref, visit, resolve)


def get_schema_mixed_enums(
    name: str,
    schema: JSONSchema,
    resolve: Callable[[str], Any],
) -> Dict[str, List[Any]]:
    result = {}

    def visit(path: str, s: JSONSchema):
        enum = s.get("enum")
        if not enum:
            return

        types = set()
        for v in enum:
            types.add(type(v))

            if len(types) != 1:
                result[path] = enum
                break

    traverse_all_subschemas(name, schema, visit, resolve)
    return result


def fill_req_body_response_diff_stats(
    key: str,
    rb_or_ref: Optional[OAPIRequestBodyObject | OAPIReferenceObject],
    parameters: Sequence[OAPIParameterObject | OAPIReferenceObject],
    resp_or_ref: OAPIResponseObject | OAPIReferenceObject,
    dst: Dict[str, Any],
    resolve: Callable[[str], Any],
):
    def collect_content_fields(contents: Mapping[str, OAPIMediaTypeObject]) -> Set[str]:
        for c in contents.values():
            schema: JSONSchema = get(c["schema"], resolve)
            if schema:
                return get_schema_field_names(schema, resolve)
        return set()

    all_params = set()
    if rb_or_ref and rb_or_ref.get("content"):
        rb: OAPIRequestBodyObject = get(rb_or_ref, resolve)
        all_params.update(collect_content_fields(rb["content"]))

    for p_or_ref in parameters:
        p: OAPIParameterObject = get(p_or_ref, resolve)
        if p.get("name"):
            all_params.update({p["name"]})
        else:
            ps: JSONSchema = get(p.get("schema", {}), resolve)
            all_params.update(get_schema_field_names(ps, resolve))

    all_response_fields = set()
    if resp_or_ref and resp_or_ref.get("content"):
        response: OAPIResponseObject = get(resp_or_ref, resolve)
        all_response_fields.update(collect_content_fields(response["content"]))

    computed = all_response_fields.difference(all_params)
    if not computed:
        return

    values = {"computed": sorted(computed)}
    if all_params:
        values.setdefault("non-computed", sorted(all_params))

    dst.setdefault(key, values)


def fill_req_body_responses_diff_stats(
    key: str,
    rb: Optional[OAPIRequestBodyObject | OAPIReferenceObject],
    parameters: Sequence[OAPIParameterObject | OAPIReferenceObject],
    responses: OAPIResponsesObject,
    dst: OAPIStats,
    resolve: Callable[[str], Any],
):
    if not filterer.should_include(COMPUTED_VARIABLES):
        return
    if not responses:
        return

    computed_vars: Dict[str, Any] = {key: []}
    for codename, resp_or_ref in responses.items():
        code = int(codename)
        if not 200 <= code < 300:
            continue
        if not key.startswith("POST"):
            continue

        response_computed: Dict[str, Any] = {}
        fill_req_body_response_diff_stats(
            codename, rb, parameters, resp_or_ref, response_computed, resolve
        )

        if response_computed.get(codename):
            computed_vars[key].append(response_computed)

    if computed_vars[key]:
        dst.setdefault(COMPUTED_VARIABLES, []).append(computed_vars)


def fill_responses_stats(
    op: OAPIOperationInfo,
    responses: OAPIResponsesObject,
    dst: OAPIStats,
    resolve: Callable[[str], Any],
):
    for code, r_or_ref in responses.items():
        r: OAPIResponseObject = get(r_or_ref, resolve)
        for t, c in r.get("content", {}).items():
            if t != "application/json" and filterer.should_include(NON_JSON_RESPONSES):
                dst.setdefault(NON_JSON_RESPONSES, {}).setdefault(op.key(), []).append(
                    {code: t}
                )

            s: JSONSchema = get(c.get("schema", {}), resolve)

            def visit_schema(path: str, schema: JSONSchema):
                if filterer.should_include(NON_SNAKECASE_VALUES):
                    if not is_snake_case(path):
                        dst.setdefault(NON_SNAKECASE_VALUES, {}).setdefault(
                            op.key(), {}
                        ).setdefault("responses", {}).setdefault(t, []).append(path)

                if filterer.should_include(XOF_SCHEMA_SINGLE_ITEM):
                    if len(which := single_item_xofs(schema)) > 0:
                        keys = (path + " | " + xof for xof in which)
                        dst.setdefault(XOF_SCHEMA_SINGLE_ITEM, {}).setdefault(
                            op.key(), {}
                        ).setdefault("response", {}).setdefault(code, []).extend(keys)

            traverse_all_subschemas("", s, visit_schema, resolve)

            if filterer.should_include(MIXED_ENUMS):
                if s:
                    mixed = get_schema_mixed_enums("", s, resolve)
                    if mixed:
                        dst.setdefault(MIXED_ENUMS, {}).setdefault(
                            op.key(), {}
                        ).setdefault("responses", {}).setdefault("content", {})[
                            t
                        ] = mixed


def fill_req_body_stats(
    op: OAPIOperationInfo,
    rb_or_ref: OAPIRequestBodyObject | OAPIReferenceObject,
    dst: OAPIStats,
    resolve: Callable[[str], Any],
):
    r: OAPIRequestBodyObject = get(rb_or_ref, resolve)
    content = r.get("content", {})

    if content:
        for t, c in content.items():
            if t != "application/json" and filterer.should_include(NON_JSON_REQUESTS):
                dst.setdefault(NON_JSON_REQUESTS, []).append({op.key(): t})

            s: JSONSchema = get(c.get("schema", {}), resolve)

            def visit_schema(path: str, schema: JSONSchema):
                if filterer.should_include(NON_SNAKECASE_VALUES):
                    if not is_snake_case(path):
                        dst.setdefault(NON_SNAKECASE_VALUES, {}).setdefault(
                            op.key(), {}
                        ).setdefault("request-bodies", {}).setdefault(t, []).append(
                            path
                        )

                if filterer.should_include(XOF_SCHEMA_SINGLE_ITEM):
                    if len(which := single_item_xofs(schema)) > 0:
                        dst.setdefault(XOF_SCHEMA_SINGLE_ITEM, {}).setdefault(
                            op.key(), {}
                        ).setdefault("request", []).extend(which)

            traverse_all_subschemas("", s, visit_schema, resolve)

            if filterer.should_include(MIXED_ENUMS):
                if s:
                    mixed = get_schema_mixed_enums("", s, resolve)
                    if mixed:
                        dst.setdefault(MIXED_ENUMS, {}).setdefault(
                            op.key(), {}
                        ).setdefault("request-bodies", {})[t] = mixed


def fill_parameters_stats(
    op: OAPIOperationInfo,
    parameters: Sequence[OAPIParameterObject | OAPIReferenceObject],
    dst: OAPIStats,
    resolve: Callable[[str], Any],
):
    for p_or_ref in parameters:
        p: OAPIParameterObject = get(p_or_ref, resolve)
        name = p.get("name", "")
        if filterer.should_include(NON_SNAKECASE_VALUES):
            if not is_snake_case(name):
                dst.setdefault(NON_SNAKECASE_VALUES, {}).setdefault(
                    op.key(), {}
                ).setdefault("parameters", []).append(name)

        s: JSONSchema = get(p.get("schema", {}), resolve)
        if s:
            if filterer.should_include(MIXED_ENUMS):
                mixed = get_schema_mixed_enums(name, s, resolve)
                if mixed:
                    dst.setdefault(MIXED_ENUMS, {}).setdefault(op.key(), {}).setdefault(
                        "parameters", []
                    ).append(mixed)

            if filterer.should_include(NON_SNAKECASE_VALUES):

                def visit_schema(path: str, schema: JSONSchema):
                    if not is_snake_case(path):
                        dst.setdefault(NON_SNAKECASE_VALUES, {}).setdefault(
                            op.key(), {}
                        ).setdefault("parameters", []).append(path)

                traverse_all_subschemas("", s, visit_schema, resolve)


def fill_operation_stats(
    op: OAPIOperationInfo, dst: OAPIStats, resolve: Callable[[str], Any]
):
    responses = op.op.get("responses", {})
    if responses:
        fill_responses_stats(op, responses, dst, resolve)

    req_body_or_ref = op.op.get("requestBody")
    if req_body_or_ref:
        fill_req_body_stats(op, req_body_or_ref, dst, resolve)

    parameters = op.op.get("parameters", [])
    if parameters:
        fill_parameters_stats(op, parameters, dst, resolve)

    fill_req_body_responses_diff_stats(
        op.key(), req_body_or_ref, parameters, responses, dst, resolve
    )

    if "operationId" not in op.op and filterer.should_include(MISSING_OPERATION_ID):
        dst.setdefault(MISSING_OPERATION_ID, []).append(op.key())

    return


def fill_missing_crud_stats(r: OAPIResource, crud_entries: List[str], dst: OAPIStats):
    if not is_tag_crud(r.tag) or not filterer.should_include(MISSING_CRUD):
        return

    missing_crud: Dict[str, List[str]] = {}
    for crud in REQUIRED_CRUD_OP_KEYS:
        if crud not in crud_entries:
            missing_crud.setdefault(r.name, []).append(crud)

    if missing_crud:
        dst.setdefault(MISSING_CRUD, []).append(missing_crud)


def fill_resource_stats(r: OAPIResource, dst: OAPIStats, resolve: Callable[[str], Any]):
    crud_entries = []

    for op in r.operations:
        fill_operation_stats(op, dst, resolve)

        if op.method in CRUD_OP_KEYS:
            crud_entries.append(op.method)

    fill_missing_crud_stats(r, crud_entries, dst)


def get_oapi_tags(o: OAPI) -> Dict[str, OAPITagInfo]:
    result = {}
    for tag in o.obj.get("tags", {}):
        name = ""
        description = ""
        extensions = {}

        for field_name, field in tag.items():
            if field_name == "name":
                name = str(field)
            elif field_name == "description":
                description = str(field)
            else:
                extensions[field_name] = field

        if not name:
            continue

        result[name] = OAPITagInfo(
            name=name, description=description, extensions=extensions
        )

    return result


def fill_resources(o: OAPI, dst: Dict[str, OAPIResource]) -> List[OAPIOperationInfo]:
    all_tags = get_oapi_tags(o)
    tagless_ops = []
    for pn, p in o.obj.get("paths", {}).items():
        for path_field, sub_fields in p.items():
            if not isinstance(sub_fields, dict) or path_field not in OPERATION_KEYS:
                continue

            op_obj = cast(OAPIOperationObject, sub_fields)
            op = OAPIOperationInfo(path=pn, method=path_field, op=op_obj)
            tags = sub_fields.get("tags")

            if tags:
                res_name = tags[0]
                tag = all_tags.get(res_name, None)
                dst.setdefault(
                    res_name, OAPIResource(name=res_name, operations=[], tag=tag)
                ).operations.append(op)
            else:
                tagless_ops.append(op)

    return tagless_ops


def fill_server_stats(
    sv: OAPIServerObject, dst: OAPIStats, resolve: Callable[[str], Any]
):
    for name, v in sv.get("variables", {}).items():
        if filterer.should_include(NON_SNAKECASE_VALUES):
            if not is_snake_case(name):
                dst.setdefault(NON_SNAKECASE_VALUES, {}).setdefault(
                    sv.get("url", ""), {}
                ).setdefault("variables", []).append(name)


def to_ref_string(*args) -> str:
    return f"#/components/{'/'.join(args)}"


def fill_ref_usages(o: OAPI):
    # Loop through all keys to find ref usages (DFS)
    visited = set()
    # Skip "openapi" key since it will not have a "$ref"
    visited |= {"openapi"}
    nodes_to_visit: List[str] = list(o.obj)
    while nodes_to_visit:
        node_name = nodes_to_visit.pop()
        if node_name in visited:
            continue
        visited.add(node_name)
        node: OAPIObject | List[OAPIObject] | str = o.obj
        for k in node_name.split(":"):
            if k.isdigit() and isinstance(node, list):
                node = node[int(k)]  # type: ignore
            else:
                node = node[k]  # type: ignore
            if k == "$ref":
                components_usage[o.name][node] = True
                continue
        if isinstance(node, dict):
            for k in node:
                nodes_to_visit.append(f"{node_name}:{k}")
        elif isinstance(node, list):
            for idx, el in enumerate(node):
                if isinstance(el, dict):
                    for k in el:
                        nodes_to_visit.append(f"{node_name}:{idx}:{k}")


def fill_parameters_usage(
    o: OAPI, parameters: Mapping[str, OAPIParameterObject | OAPIReferenceObject]
):
    for name in parameters.keys():
        ref = to_ref_string("parameters", name)
        components_usage[o.name].setdefault(ref, False)


def fill_schemas_usage(o: OAPI, schemas: Mapping[str, JSONSchema]):
    for name, spec in schemas.items():  # type: ignore
        # components.schema.CreateResponse for example
        ref = to_ref_string("schemas", name)
        # May have been found by a ref inside a component, we don't want
        # to simply overwrite, thus check.
        if ref not in components_usage[o.name]:
            components_usage[o.name][ref] = False
        # Loop through component itself to find other components
        for prop in spec.get("properties", {}).values():
            # Already found through reference, mark as used
            if prop_ref := prop.get("$ref") or prop.get("items", {}).get("$ref"):
                components_usage[o.name][prop_ref] = True


def fill_components_usage(o: OAPI) -> None:
    if components := o.obj.get("components"):
        if parameters := components.get("parameters"):
            fill_parameters_usage(o, parameters)
        if schemas := components.get("schemas"):
            fill_schemas_usage(o, schemas)
        # TODO: check for all types, but for this we need to loop over
        # the whole OAPI tree to find the occurrences
    fill_ref_usages(o)


def get_oapi_stats(o: OAPI) -> OAPIStats:
    result: OAPIStats = {}
    resources: Dict[str, OAPIResource] = {}
    tagless_ops = fill_resources(o, resources)

    for res in resources.values():
        fill_resource_stats(res, result, o.resolve)

    for op in tagless_ops:
        fill_operation_stats(op, result, o.resolve)
        if filterer.should_include(TAGLESS_OPERATIONS):
            result.setdefault(TAGLESS_OPERATIONS, []).append(op.key())

    for sv in o.obj.get("servers", []):
        fill_server_stats(sv, result, o.resolve)

    if filterer.should_include(UNUSED_COMPONENTS):
        fill_components_usage(o)
        for name, being_used in components_usage.get(o.name, {}).items():
            if not being_used:
                if UNUSED_COMPONENTS not in result:
                    result[UNUSED_COMPONENTS] = []
                result[UNUSED_COMPONENTS].append(name)

    # TODO: Add stats for other fields

    return result


def dump_stats(stats: Dict[str, OAPIStats], output: str):
    dump = yaml.dump(
        stats, Dumper=YamlDumper, sort_keys=True, indent=2, allow_unicode=True
    )
    if output:
        with open(output, "w") as fd:
            fd.write(dump)
    else:
        print(dump)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Validate response and request bodies for all OAPI YAML"
        "files in directory"
    )
    parser.add_argument(
        "dir_or_file", type=str, help="Directory of OpenAPI files or OpenAPI file path"
    )
    parser.add_argument(
        "--filter", type=str, action="append", default=[], help="Only show these stats"
    )
    parser.add_argument(
        "--filter-out",
        type=str,
        action="append",
        default=[],
        help="Don't show these stats",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        default="",
        help="Output target file to dump results",
    )
    parser.add_argument(
        "--ignore-disabled",
        type=bool,
        default=True,
        help="Don't load OpenAPI files that have '.disabled' in their name",
    )

    args = parser.parse_args()

    filterer.filters = args.filter
    filterer.filter_out = args.filter_out

    oapis = load_oapis(args.dir_or_file, args.ignore_disabled)
    all_stats: Dict[str, OAPIStats] = {}
    for o in oapis:
        stats = get_oapi_stats(o)
        if stats:
            all_stats[o.name] = stats

    dump_stats(all_stats, args.output)
