SCOPE_EXTENSION = "x-mgc-project-scope"
HEADER_NAME = "x-project-id"
CONFIG_NAME = "project"
DESCRIPTION = (
    "Project to scope this command to. Prefer 'mgc alpha project set', which "
    "resolves the project by id or name"
)


def project_scope_parameter() -> dict:
    # Sem default nem required: sem projeto configurado o header nao e enviado.
    return {
        "name": HEADER_NAME,
        "in": "header",
        "required": False,
        "description": DESCRIPTION,
        "x-mgc-name": CONFIG_NAME,
        "schema": {
            "type": "string",
            "title": "Project",
            "description": DESCRIPTION,
        },
    }


def is_project_scoped(spec: dict) -> bool:
    servers = spec.get("servers") or []
    return bool(servers) and bool(servers[0].get(SCOPE_EXTENSION))


def add_project_scope(spec: dict) -> int:
    if not is_project_scoped(spec):
        return 0

    added = 0
    for path in spec.get("paths", {}).values():
        parameters = path.setdefault("parameters", [])
        if any(p.get("name") == HEADER_NAME for p in parameters):
            continue
        # Sempre no fim: o yaml_merge funde listas por indice.
        parameters.append(project_scope_parameter())
        added += 1

    return added
