from typing import Any, Dict
import yaml
import warnings
import argparse
import urllib.request
import json

SERVER_URL_MAP = {
    # VM
    "https://virtual-machine.br-ne-1.jaxyendy.com/docs": (
        "https://api-virtual-machine.br-ne-1.jaxyendy.com/"
    ),
    "https://virtual-machine.br-ne1-prod.jaxyendy.com": (
        "https://api-virtual-machine.br-ne-1.jaxyendy.com/"
    ),
    # Block Storage
    "https://block-storage.br-ne-1.jaxyendy.com/docs": (
        "https://api-block-storage.br-ne-1.jaxyendy.com/"
    ),
    # VPC
    "https://vpc.br-ne-1.jaxyendy.com/docs": ("https://api-vpc.br-ne-1.jaxyendy.com/"),
    # Object Storage
    "https://object-storage.br-ne-1.jaxyendy.com/docs": (
        "https://api-object-storage.br-ne-1.jaxyendy.com/"
    ),
    # DBaaS
    "https://dbaas.br-ne-1.jaxyendy.com/docs": (
        "https://api-dbaas.br-ne-1.jaxyendy.com/"
    ),
    # K8S
    "https://mke.br-ne-1.jaxyendy.com/docs": ("https://api-mke.br-ne-1.jaxyendy.com/"),
}


def sync_request_body(internal_spec: Dict[str, Any], external_spec: Dict[str, Any]):
    for ext_path in external_spec.get("paths", {}):
        internal_path = internal_spec.get("paths", {}).get(ext_path)
        if not internal_path:
            # No problem, it was added to Kong but not in internal
            continue

        for ext_action in ext_path:
            internal_action = internal_path.get(ext_action)
            if not internal_action:
                # Action mapped on external but not on internal,
                # should never happen
                continue

            if internal_action["requestBody"]:
                ext_action["requestBody"] = internal_action["requestBody"]


def fetch_and_parse(json_oapi_url: str) -> Dict[str, Any]:
    with urllib.request.urlopen(json_oapi_url, timeout=5) as response:
        return json.loads(response.read())


def load_yaml(path: str) -> Dict[str, Any]:
    with open(path, "r") as fd:
        return yaml.load(fd, Loader=yaml.CLoader)


def update_server_urls(spec: Dict[str, Any]):
    assert "servers" in spec, "Servers key not present in external YAML"
    for server in spec["servers"]:
        if server["url"] not in SERVER_URL_MAP:
            warnings.warn(
                f"Unrecognized url in external: {server['url']}", category=UserWarning
            )
        else:
            server["url"] = SERVER_URL_MAP[server["url"]]


def save_external(spec: Dict[str, Any], path: str):
    with open(path, "w") as fd:
        yaml.dump(spec, fd, sort_keys=False, indent=4)


if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="SyncOAPI",
        description="Sync external OAPI schema with the internal schema by "
        "fixing any mismatch of requestBody between external and "
        "internal impl. After, we change the server URL to Kong and "
        "adjust schema of error returns. The YAML generated can "
        "be used in Kong directly to serve as a ref. to external.",
    )
    # Internal = APIs generated directly from the code, always udpated
    parser.add_argument(
        "internal_spec_url",
        type=str,
        help="URL to fetch current internal OpenAPI spec, which will "
        "come in JSON format",
    )
    # External = Viveiro in MGC context, intermediate between product and Kong
    parser.add_argument(
        "external_spec_path",
        type=str,
        help="File path to current external OpenAPI spec",
    )
    parser.add_argument(
        "-o",
        "--output",
        type=str,
        default="new-ext-oapi.yaml",
        help="Path to save the new external YAML. Defaults to 'new-ext-oapi.yaml'",
    )
    args = parser.parse_args()

    # Load json into dict
    internal_spec = fetch_and_parse(args.internal_spec_url)
    # Load yaml into dict
    external_spec = load_yaml(args.external_spec_path)

    # Replace requestBody from external to the internal value if they mismatch
    sync_request_body(internal_spec, external_spec)

    # Replace server url
    update_server_urls(external_spec)

    # TODO: Replace Error Object

    # Write external to file
    save_external(external_spec, args.output)
