import jsonschema

import argparse
from typing import List, cast
from urllib import parse
from oapi_types import OAPI
from transform_helpers import fetch_and_parse, load_yaml, save_external, add_spec_uid
from spec_add_tags_block import AddTagsBlockTransformer
from spec_types import SpecTranformer
from spec_version_convert import ConvertVersionTransformer
from spec_remove_param import RemoveParamTransformer
from spec_remove_path import RemovePathTransformer
from spec_remove_component import RemoveComponentTransformer
from spec_update_error import UpdateErrorTransformer
from spec_add_security import AddSecurityTransformer
from spec_fix_links import FixLinksTransformer
from spec_create_links import CreateLinks

if __name__ == "__main__":
    parser = argparse.ArgumentParser(
        prog="Transform OAPI Spec",
        description="Transforms a product OAPI schema removing internal "
        "elements and making basic ajustments necessary to create a "
        '"public" schema that can be used by the MGC SDK.',
    )
    # Internal = APIs generated directly from the code, always udpated
    parser.add_argument(
        "product_name",
        type=str,
        help="Product name as it will appear in CLI",
    )
    parser.add_argument(
        "spec_file",
        type=str,
        help="Raw product open API schema, can be yaml or json",
    )
    parser.add_argument(
        "spec_uid",
        type=str,
        help="Universal identifier for the specification, this will be used "
        "to identify operations between different products, should be an URL",
    )
    parser.add_argument(
        "-o",
        "--output",
        required=True,
        type=str,
        help="Path to save the new external YAML",
    )
    args = parser.parse_args()

    if parse.urlparse(args.spec_file).scheme != "":
        # If is a valid url fetch and load
        product_spec = fetch_and_parse(args.spec_file)
    else:
        # Load spec into dict
        product_spec = load_yaml(args.spec_file)

    ref_resolver = jsonschema.RefResolver(args.spec_file, cast(dict, product_spec))
    oapi = OAPI(
        path=args.spec_file,
        name=args.product_name,
        obj=product_spec,
        ref_resolver=ref_resolver,
    )

    # Perform changes in the spec
    transformers: List[SpecTranformer] = [
        ConvertVersionTransformer(),
        UpdateErrorTransformer(),
        RemovePathTransformer(".*xaas.*"),
        RemoveComponentTransformer(".*(xaas|XAAS|Xaas).*"),
        RemoveParamTransformer("x-tenant-id"),
        AddSecurityTransformer(args.product_name),
        # AddParameterTypes(),
        AddTagsBlockTransformer(),
        FixLinksTransformer(),
        CreateLinks(),
    ]
    for t in transformers:
        t.transform(oapi)

    # Write external to output file
    add_spec_uid(product_spec, args.spec_uid)
    save_external(oapi.obj, args.output)
