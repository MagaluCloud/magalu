from typing import Any, Callable, Dict, Hashable
import argparse
import yaml

OAPISchema = Dict[str, Any]


def load_yaml(path: str) -> OAPISchema:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.FullLoader)


def save_external(spec: OAPISchema, path: str):
    with open(path, "w") as fd:
        yaml.dump(spec, fd, sort_keys=False, indent=4, allow_unicode=True)


def merge_item(value: Any, new: Any, override: bool, path: list[str]) -> bool:
    if value == new:
        return False

    value_type = type(value)
    new_type = type(new)
    hasheable_types = (
        isinstance(value_type, Hashable),
        isinstance(new_type, Hashable),
    )
    if not all(hasheable_types) and value_type != new_type:
        raise NotImplementedError(
            f"{path}: cannot merge type {value_type!r} with {new_type!r}"
        )

    if value_type == dict:
        merge_dict(value, new, override, path)
        return False
    elif value_type == list:
        merge_list(value, new, override, path)
        return False
    elif not override:
        raise ValueError(f"{path}: not overriding {value!r} with {new!r}")
    else:
        return True


def merge_dict(dst: dict, extra: dict, override: bool, path: list[str]) -> dict:
    for k, new in extra.items():
        value = dst.setdefault(k, new)
        if value != new:
            replace = merge_item(value, new, override, path + [k])
            if replace:
                dst[k] = new

    return dst


def merge_list(dst: list, extra: list, override: bool, path: list[str]) -> list:
    for i, new in enumerate(extra):
        if len(dst) <= i:
            dst.append(new)
        else:
            value = dst[i]
            replace = merge_item(value, new, override, path + [str(i)])
            if replace:
                dst[i] = new

    return dst


def check_is_string(v: Any, _: Any) -> None:
    if not isinstance(v, str):
        raise ValueError("not a string")


def check_is_bool(v: Any, _: Any) -> None:
    if not isinstance(v, bool):
        raise ValueError("not a bool")


def check_is_number(v: Any, _: Any) -> None:
    if not isinstance(v, (int, float)):
        raise ValueError("not a number")


def check_unknown_default(k: str, v: Any, _: Any) -> Any:
    raise ValueError(f"unexpected customization {v!r}")


def check_any(v: Any, _: Any) -> None:
    return


def check_is_object(
    o: Any,
    base: dict | None,
    customization_checkers: Dict[str, Callable[[Any, Any], None]],
    allowed_matches: set[str] = set(),
    check_unknown: Callable[[str, Any, Any], None] = check_unknown_default,
    key_name="property",
) -> None:
    if not isinstance(o, dict):
        raise ValueError("is not a dict")

    if base is None:
        base = {}

    for k, v in o.items():
        try:
            b = base.get(k)
            if k in allowed_matches:
                if b is not None and v != b:
                    raise ValueError(f"should match {b!r}, but got {v!r}")

                continue

            checker = customization_checkers.get(k)
            if checker is None:

                def checker(a: Any, b: Any) -> None:
                    check_unknown(k, a, b)

            checker(v, b)
        except Exception as e:
            raise ValueError(f"{key_name} {k!r}: {e}") from e


def check_is_existing_key(k: Any, b: dict) -> None:
    try:
        b[k]
    except KeyError as e:
        known = ", ".join((b or {}).keys())
        raise ValueError(f"existing key is required. Existing: {known}") from e


def check_is_dict(
    o: Any,
    base: dict | None,
    check_value: Callable[[Any, Any], None],
    check_key: Callable[[Any, dict], None] = lambda k, b: None,
    key_name="key",
) -> None:
    if not isinstance(o, dict):
        raise ValueError("is not a dict")

    if base is None:
        base = {}

    for k, v in o.items():
        try:
            check_key(k, base)
            check_value(v, base.get(k))
        except Exception as e:
            raise ValueError(f"{key_name} {k!r}: {e}") from e


def check_is_list(
    o: Any,
    base: list | None,
    check_value: Callable[[Any, Any], None],
    item_name="item",
) -> None:
    if not (base is None or isinstance(base, list)):
        raise ValueError(f"base must be a list or None, got {base!r}")

    if not isinstance(o, list):
        raise ValueError(f"value must be a list, got {o!r}")

    for i, v in enumerate(o):
        if base is not None and i < len(base):
            b = base[i]
        else:
            b = None

        try:
            check_value(v, b)
        except Exception as e:
            raise ValueError(f"{item_name} {i}: {e}")


def check_is_transforms(v: Any, _: Any) -> None:
    if not isinstance(v, (dict, list)):
        raise ValueError(f"not a valid transforms specification: {v!r}")


def check_is_schema_customization_properties(v: Any, base: Any) -> None:
    # TODO: check name exists
    check_is_dict(
        v,
        base,
        check_is_schema_customization,
    )


def check_is_schema_customization_list(v: Any, base: Any) -> None:
    # TODO: check item exists
    raise NotImplementedError("TODO")


def check_is_schema_customization(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_SCHEMA_CUSTOMIZATIONS,
    )


SUPPORTED_SCHEMA_CUSTOMIZATIONS = {
    "x-mgc-name": check_is_string,
    "x-mgc-description": check_is_string,
    "x-mgc-hidden": check_is_bool,
    "x-mgc-transforms": check_is_transforms,
    "type": check_is_string,
    "properties": check_is_schema_customization_properties,
    "allOf": check_is_schema_customization_list,
    "anyOf": check_is_schema_customization_list,
    "oneOf": check_is_schema_customization_list,
    "minimum": check_is_number,
    "exclusiveMinimum": check_is_bool,
    "example": check_any,
    "default": check_any,
}


SUPPORTED_WAIT_TERMINATION_CUSTOMIZATIONS = {
    "maxRetries": check_is_number,
    "interval": check_is_string,
    "jsonPathQuery": check_is_string,
    "templateQuery": check_is_string,
    "errorJsonPathQuery": check_is_string,
    "errorTemplateQuery": check_is_string,
}


def check_is_wait_termination(v: Any, b: Any) -> None:
    check_is_object(
        v,
        b,
        SUPPORTED_WAIT_TERMINATION_CUSTOMIZATIONS,
    )
    assert isinstance(v, dict)
    if not v.get("jsonPathQuery", v.get("templateQuery")):
        raise ValueError("missing jsonPathQuery and templateQuery")


SUPPORTED_TAG_MATCH = {"name", "description"}
SUPPORTED_TAG_CUSTOMIZATIONS = {
    "x-mgc-name": check_is_string,
    "x-mgc-description": check_is_string,
    "x-mgc-hidden": check_is_bool,
    "name": check_is_string,
    "description": check_is_string,
}


def check_is_doc_tag(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_TAG_CUSTOMIZATIONS,
        SUPPORTED_TAG_MATCH,
    )


def check_is_doc_tags(v: Any, base: Any) -> None:
    check_is_list(v, base, check_is_doc_tag, item_name="tag")


def check_is_server_variable(v: Any, base: Any) -> None:
    if not isinstance(v, dict):
        raise ValueError("is not a dict")


def check_is_server_variables(v: Any, base: Any) -> None:
    check_is_dict(
        v,
        base,
        # we don't check_is_schema_customization() since we allow new schemas
        check_is_server_variable,
        key_name="variable",
    )


SUPPORTED_SERVER_MATCH = {"url", "description"}
SUPPORTED_SERVER_CUSTOMIZATIONS = {
    "url": check_is_string,
    "description": check_is_string,
    "variables": check_is_server_variables,
}


def check_is_server(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_SERVER_CUSTOMIZATIONS,
    )


def check_is_doc_servers(v: Any, base: Any) -> None:
    check_is_list(
        v,
        base,
        check_is_server,
        item_name="server",
    )


def check_is_link_parameters(v: Any, b: Any) -> None:
    check_is_dict(
        v,
        b,
        check_is_string,
        key_name="parameter",
    )


def check_is_extra_parameters(v: Any, b: Any) -> None:
    if not isinstance(v, list):
        raise ValueError("is not a list")


SUPPORTED_PATH_PARAMETER_MATCH = {"name", "description", "in"}
SUPPORTED_PATH_PARAMETER_CUSTOMIZATIONS = {
    "x-mgc-name": check_is_string,
    "x-mgc-description": check_is_string,
    "x-mgc-hidden": check_is_bool,
    "schema": check_is_schema_customization,
}


def check_is_path_parameter(v: Any, b: Any) -> None:
    check_is_object(
        v,
        b,
        SUPPORTED_PATH_PARAMETER_CUSTOMIZATIONS,
        SUPPORTED_PATH_PARAMETER_MATCH,
    )


def check_is_path_parameters(v: Any, b: Any) -> None:
    check_is_list(
        v,
        b,
        check_is_path_parameter,
    )


SUPPORTED_LINK_CUSTOMIZATIONS = {
    "id": check_is_string,
    "description": check_is_string,
    "operationId": check_is_string,
    "operationRef": check_is_string,
    "parameters": check_is_link_parameters,
    "x-mgc-wait-termination": check_is_wait_termination,
    "x-mgc-extra-parameters": check_is_extra_parameters,
    "x-mgc-hidden": check_is_bool,
}


def check_is_link(v: Any, b: Any) -> None:
    check_is_object(
        v,
        b,
        SUPPORTED_LINK_CUSTOMIZATIONS,
    )


def check_is_links(v: Any, b: Any) -> None:
    check_is_dict(
        v,
        b,
        check_is_link,
        key_name="link",
    )


SUPPORTED_RESPONSE_CUSTOMIZATIONS = {
    "links": check_is_links,
}


def check_is_response(v: Any, b: Any) -> None:
    check_is_object(
        v,
        b,
        SUPPORTED_RESPONSE_CUSTOMIZATIONS,
    )


def check_is_responses(v: Any, b: Any) -> None:
    check_is_dict(
        v,
        b,
        check_is_response,
        check_key=check_is_existing_key,
        key_name="response",
    )


def check_is_security_scopes(v: Any, b: Any) -> None:
    check_is_list(
        v,
        b,
        check_is_string,
    )


def check_is_security_item(v: Any, b: Any) -> None:
    check_is_dict(
        v,
        b,
        check_is_security_scopes,
    )


def check_is_security_list(v: Any, b: Any) -> None:
    check_is_list(
        v,
        b,
        check_is_security_item,
        item_name="security",
    )


def check_is_path_method_tag_list(v: Any, base: Any) -> None:
    check_is_list(
        v,
        base,
        check_is_string,
    )


SUPPORTED_CONFIRMABLE_CUSTOMIZATIONS = {
    "message": check_is_string,
}


def check_is_confirmable(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_CONFIRMABLE_CUSTOMIZATIONS,
    )


SUPPORTED_PROMPT_INPUT_CUSTOMIZATIONS = {
    "message": check_is_string,
    "confirmValue": check_is_string,
}


def check_is_promptInput(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_PROMPT_INPUT_CUSTOMIZATIONS,
    )


SUPPORTED_PATH_METHOD_CUSTOMIZATIONS = {
    "x-mgc-name": check_is_string,
    "x-mgc-description": check_is_string,
    "x-mgc-hidden": check_is_bool,
    "x-mgc-output-flag": check_is_string,
    "x-mgc-transforms": check_is_transforms,
    "x-mgc-wait-termination": check_is_wait_termination,
    "x-mgc-confirmable": check_is_confirmable,
    "x-mgc-confirmPrompt": check_is_promptInput,
    "parameters": check_is_path_parameters,
    "responses": check_is_responses,
    "security": check_is_security_list,
    "tags": check_is_path_method_tag_list,
}


def check_is_doc_path_method(v: Any, b: Any) -> None:
    check_is_object(
        v,
        b,
        SUPPORTED_PATH_METHOD_CUSTOMIZATIONS,
    )


def check_is_doc_path(v: Any, b: Any) -> None:
    check_is_dict(
        v,
        b,
        check_is_doc_path_method,
        check_key=check_is_existing_key,
        key_name="method",
    )


def check_is_doc_paths(v: Any, base: Any) -> None:
    check_is_dict(
        v,
        base,
        check_is_doc_path,
        check_key=check_is_existing_key,
        key_name="path",
    )


def create_check_existing_dict(check_value) -> Callable[[Any, Any], None]:
    def dict_checker(v: Any, base: Any) -> None:
        check_is_dict(
            v,
            base,
            check_value=check_value,
            check_key=check_is_existing_key,
        )

    return dict_checker


SUPPORTED_COMPONENT_SCHEMAS = {
    "schemas": create_check_existing_dict(check_is_schema_customization),
    "responses": create_check_existing_dict(check_is_response),
    "parameters": create_check_existing_dict(check_is_path_parameter),
    "links": create_check_existing_dict(check_is_link),
    # TODO: do these when we have usage
    # "requestBodies": create_check_existing_dict(check_is_request_body),
    # "headers": create_check_existing_dict(check_is_header),
    # "securitySchemes": create_check_existing_dict(check_is_security_scheme),
    # "callbacks": create_check_existing_dict(check_is_callback),
}


def check_is_doc_components(v: Any, base: Any) -> None:
    check_is_object(
        v,
        base,
        SUPPORTED_COMPONENT_SCHEMAS,
    )


SUPPORTED_DOC_CUSTOMIZATIONS = {
    "paths": check_is_doc_paths,
    "tags": check_is_doc_tags,
    "components": check_is_doc_components,
    "servers": check_is_doc_servers,
}


def validate_oapi_customizations(base: OAPISchema, c: dict):
    check_is_object(
        c,
        base,
        SUPPORTED_DOC_CUSTOMIZATIONS,
    )


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        description="Merge two YAML files",
    )
    parser.add_argument(
        "base",
        type=str,
        help="the base file to open",
    )
    # External = Viveiro in MGC context, intermediate between product and Kong
    parser.add_argument(
        "extra",
        type=str,
        help="the extra file to merge on top of base",
    )
    parser.add_argument(
        "--override",
        action="store_true",
        default=False,
        help="Override existing scalars",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        help="Path to save the new external YAML. Defaults to overwrite base",
    )
    args = parser.parse_args()

    base = load_yaml(args.base)
    extra = load_yaml(args.extra)
    validate_oapi_customizations(base, extra)

    merge_dict(base, extra, args.override, [])

    save_external(base, args.output or args.base)
