SCOPE_EXTENSION = "x-mgc-project-scope"
HEADER_NAME = "x-project-id"
CONFIG_NAME = "project"
DESCRIPTION = (
    "Project to scope this command to. Prefer 'mgc project set', which "
    "resolves the project by id or name and checks that it exists"
)


def project_scope_parameter() -> dict:
    """O parametro que o `in: header` transforma em config da CLI.

    Sem `default` de proposito: sem projeto configurado a chave nao chega em
    `configs` e o header nao e enviado, que e como se pede o projeto default do
    tenant. `required` tampouco, pelo mesmo motivo.
    """
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
    """Acrescenta o parametro de escopo em todo path do spec.

    No nivel do path, nao no da operacao: o loader le as duas listas e a do path
    vale para todos os metodos dele, entao uma entrada por path basta.

    Sempre no fim da lista. As customizacoes identificam parametros por indice
    (o yaml_merge funde listas por posicao), entao inserir no comeco deslocaria
    todas elas.
    """
    if not is_project_scoped(spec):
        return 0

    added = 0
    for path in spec.get("paths", {}).values():
        parameters = path.setdefault("parameters", [])
        if any(p.get("name") == HEADER_NAME for p in parameters):
            continue
        parameters.append(project_scope_parameter())
        added += 1

    return added
