from typing import Any, Callable, Dict

from oapi_types import OAPI
from spec_types import SpecTranformer


class ConvertVersionTransformer(SpecTranformer):
    """Convert the schemas to OpenAPI's 3.0

    In golang we use https://github.com/getkin/kin-openapi
    that still does not support 3.1, so we must convert to 3.0

    We do the reverse of:
    https://www.openapis.org/blog/2021/02/16/migrating-from-openapi-3-0-to-3-1-0
    """

    _dict_item_converters: Dict[str, Callable[[dict, str, Any], bool]]

    def transform(self, oapi: OAPI):
        self._fix_openapi_version(oapi)

    def _fix_openapi_version(self, oapi: OAPI):
        spec = oapi.obj
        maj_ver, min_ver = tuple(int(x) for x in spec["openapi"].split("."))[:2]
        if maj_ver != 3:
            raise ValueError(f"unsupported openapi major version {maj_ver}")

        if min_ver > 1:
            raise ValueError(f"unsupported openapi minor version {maj_ver}")

        spec["openapi"] = "3.0.3"
        self._dict_item_converters = {
            "examples": self._convert_list_of_examples,
            "anyOf": self._convert_list_of_nullable,
            "oneOf": self._convert_list_of_nullable,
            "exclusiveMinimum": self._convert_exclusive_constraint,
            "exclusiveMaximum": self._convert_exclusive_constraint,
        }
        self._convert(spec)

    def _convert(self, d: any) -> None:
        if not isinstance(d, dict):
            if isinstance(d, list):
                for v in d:
                    self._convert(v)
            return

        for k, v in list(d.items()):  # list forces a copy, we may mutate 'd'
            converter = self._dict_item_converters.get(k)
            if converter is not None and converter(d, k, v):
                continue

            self._convert(v)

    def _convert_list_of_nullable(
        self,
        d: dict,
        k: str,
        v: Any,
    ) -> bool:
        if not isinstance(v, list):
            return False

        list_of: list[dict] = v
        remaining = [item for item in list_of if item.get("type") != "null"]
        if len(remaining) == len(list_of):
            return False

        d["nullable"] = True
        if len(remaining) > 1:
            d[k] = remaining
        else:
            del d[k]
            if len(remaining) == 1:
                if examples := remaining[0].pop("examples", None):
                    remaining[0]["example"] = examples[0]

                if len(remaining) == 1:
                    if "additionalProperties" in remaining[0]:
                        list_of: list[dict] = remaining[0]["additionalProperties"][
                            "anyOf"
                        ]
                        self._convert_list_of_nullable(
                            d=remaining[0]["additionalProperties"], k=k, v=list_of
                        )

                d.update(remaining[0])
        return True

    def _convert_list_of_examples(
        self,
        d: dict,
        k: str,
        v: Any,
    ) -> bool:
        if not isinstance(v, list):
            return False

        list_of: list[Any] = v
        del d[k]
        if len(list_of) > 0:
            d["example"] = list_of[0]
        return True

    def _convert_exclusive_constraint(
        self,
        d: dict,
        k: str,
        v: Any,
    ) -> bool:
        if isinstance(v, bool):
            return False

        d[k] = True
        d[k.removeprefix("exclusive").lower()] = v
        return True
