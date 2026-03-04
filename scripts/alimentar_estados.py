#!/usr/bin/env python3
"""
Usa a API para garantir a carga dos estados, UFMA e Ciencia da Computacao.

Requer token Bearer de um usuario com papel ADMIN ou DEV.
"""

from __future__ import annotations

import argparse
import json
import os
import sys
from typing import Any
from urllib import error, parse, request


EXPECTED_STATE_COUNT = 27
STATES: list[tuple[str, str]] = [
    ("AC", "Acre"),
    ("AL", "Alagoas"),
    ("AP", "Amapa"),
    ("AM", "Amazonas"),
    ("BA", "Bahia"),
    ("CE", "Ceara"),
    ("DF", "Distrito Federal"),
    ("ES", "Espirito Santo"),
    ("GO", "Goias"),
    ("MA", "Maranhao"),
    ("MT", "Mato Grosso"),
    ("MS", "Mato Grosso do Sul"),
    ("MG", "Minas Gerais"),
    ("PA", "Para"),
    ("PB", "Paraiba"),
    ("PR", "Parana"),
    ("PE", "Pernambuco"),
    ("PI", "Piaui"),
    ("RJ", "Rio de Janeiro"),
    ("RN", "Rio Grande do Norte"),
    ("RS", "Rio Grande do Sul"),
    ("RO", "Rondonia"),
    ("RR", "Roraima"),
    ("SC", "Santa Catarina"),
    ("SP", "Sao Paulo"),
    ("SE", "Sergipe"),
    ("TO", "Tocantins"),
]


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Alimenta estados e estrutura academica via API.")
    parser.add_argument("--base-url", default=os.getenv("API_BASE_URL", "http://localhost:8080"))
    parser.add_argument("--token", default=os.getenv("API_TOKEN"))
    args = parser.parse_args()

    if not args.token:
        parser.error("o token e obrigatorio. Use --token ou API_TOKEN.")

    return args


def build_url(base_url: str, path: str) -> str:
    return base_url.rstrip("/") + path


def api_request(base_url: str, path: str, token: str, method: str = "GET", payload: dict[str, Any] | None = None) -> Any:
    data = None
    headers = {
        "Authorization": f"Bearer {token}",
        "Accept": "application/json",
    }

    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = request.Request(build_url(base_url, path), data=data, method=method, headers=headers)

    try:
        with request.urlopen(req) as response:
            body = response.read()
            if not body:
                return None
            return json.loads(body.decode("utf-8"))
    except error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"Falha HTTP {exc.code} em {path}: {details}") from exc
    except error.URLError as exc:
        raise RuntimeError(f"Nao foi possivel conectar em {base_url}: {exc.reason}") from exc


def get_states(base_url: str, token: str) -> list[dict[str, Any]]:
    response = api_request(base_url, "/states", token)
    if not isinstance(response, list):
        raise RuntimeError("Resposta invalida ao listar estados.")
    return response


def get_current_user(base_url: str, token: str) -> dict[str, Any]:
    response = api_request(base_url, "/users/me", token)
    if not isinstance(response, dict):
        raise RuntimeError("Resposta invalida ao validar token.")
    return response


def get_state_by_code(base_url: str, token: str, code: str) -> dict[str, Any]:
    response = api_request(base_url, f"/states/{parse.quote(code)}", token)
    if not isinstance(response, dict):
        raise RuntimeError(f"Resposta invalida ao buscar estado {code}.")
    return response


def create_state(base_url: str, token: str, code: str, name: str) -> dict[str, Any]:
    response = api_request(
        base_url,
        "/states",
        token,
        method="POST",
        payload={"code": code, "name": name},
    )
    if not isinstance(response, dict):
        raise RuntimeError(f"Resposta invalida ao criar estado {code}.")
    return response


def get_universities(base_url: str, token: str, state_id: int) -> list[dict[str, Any]]:
    response = api_request(base_url, f"/states/{state_id}/universities", token)
    if not isinstance(response, list):
        raise RuntimeError("Resposta invalida ao listar universidades.")
    return response


def create_university(base_url: str, token: str, state_id: int, name: str) -> dict[str, Any]:
    response = api_request(
        base_url,
        f"/states/{state_id}/universities",
        token,
        method="POST",
        payload={"name": name},
    )
    if not isinstance(response, dict):
        raise RuntimeError(f"Resposta invalida ao criar universidade {name}.")
    return response


def get_courses(base_url: str, token: str, university_id: int) -> list[dict[str, Any]]:
    response = api_request(base_url, f"/universities/{university_id}/courses", token)
    if not isinstance(response, list):
        raise RuntimeError("Resposta invalida ao listar cursos.")
    return response


def create_course(base_url: str, token: str, university_id: int, name: str) -> dict[str, Any]:
    response = api_request(
        base_url,
        f"/universities/{university_id}/courses",
        token,
        method="POST",
        payload={"name": name},
    )
    if not isinstance(response, dict):
        raise RuntimeError(f"Resposta invalida ao criar curso {name}.")
    return response


def main() -> int:
    args = parse_args()
    base_url = args.base_url
    token = args.token

    current_user = get_current_user(base_url, token)
    if current_user.get("role") not in {"ADMIN", "DEV"}:
        raise RuntimeError("O token informado precisa pertencer a um usuario ADMIN ou DEV.")

    states_before = get_states(base_url, token)
    existing_state_codes = {state.get("code") for state in states_before}
    states_inserted = 0

    for code, name in STATES:
        if code in existing_state_codes:
            continue
        create_state(base_url, token, code, name)
        states_inserted += 1

    maranhao = get_state_by_code(base_url, token, "MA")
    universities = get_universities(base_url, token, maranhao["id"])
    ufma = next((u for u in universities if u.get("name") == "UFMA"), None)
    ufma_created = False
    if ufma is None:
        ufma = create_university(base_url, token, maranhao["id"], "UFMA")
        ufma_created = True

    courses = get_courses(base_url, token, ufma["id"])
    computer_science_created = False
    if not any(c.get("name") == "Ciencia da Computacao" for c in courses):
        create_course(base_url, token, ufma["id"], "Ciencia da Computacao")
        computer_science_created = True

    courses = get_courses(base_url, token, ufma["id"])
    if not any(c.get("name") == "Ciencia da Computacao" for c in courses):
        raise RuntimeError("Curso Ciencia da Computacao nao encontrado apos criacao.")

    states_after = get_states(base_url, token)
    print(json.dumps(
        {
            "statesBefore": len(states_before),
            "statesAfter": len(states_after),
            "statesExpected": EXPECTED_STATE_COUNT,
            "statesInserted": states_inserted,
            "ufmaCreated": ufma_created,
            "computerScienceCreated": computer_science_created,
            "validated": {
                "userRole": current_user.get("role"),
                "stateCode": "MA",
                "university": "UFMA",
                "course": "Ciencia da Computacao",
            },
        },
        ensure_ascii=True,
        indent=2,
    ))
    return 0


if __name__ == "__main__":
    raise SystemExit(main())
