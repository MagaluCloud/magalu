#!/bin/bash

# Captura o input da versão
version="$1"

# Verifica se o input foi fornecido
if [ -z "$version" ]; then
    echo "Erro: Nenhuma versão fornecida."
    exit 1
fi

# Expressão regular para validar o padrão sem sufixos
if [[ "$version" =~ ^v[0-9]+\.[0-9]+\.[0-9]+$ ]]; then
    echo "Versão válida: $version"
    exit 0  # Sucesso
else
    echo "Erro: Versão inválida para publicação de sub modulos: $version"
    exit 1  # Erro
fi