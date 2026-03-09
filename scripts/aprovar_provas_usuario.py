#!/usr/bin/env python3
"""
Aprova em lote todas as provas pendentes de um usuario.

Permissoes:
- token precisa ser ADMIN, DEV ou APPROVER.
"""

from __future__ import annotations

import argparse
import json
import logging
from typing import Any
from urllib import error, request

from common import get_api_base_url


ALLOWED_ROLES = {"ADMIN", "DEV", "APPROVER"}


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Aprova em lote provas pendentes de um usuario.")
    parser.add_argument("--base-url", default=get_api_base_url())
    parser.add_argument("--token", help="Token Bearer (sem 'Bearer ').")
    parser.add_argument("--user-id", type=int, help="ID do usuario alvo (uploader).")
    parser.add_argument("--email", help="Email do usuario alvo (uploader).")
    parser.add_argument("--university-id", type=int, help="ID da universidade para filtrar pendencias.")
    parser.add_argument("--course-id", type=int, help="ID do curso para filtrar pendencias.")
    parser.add_argument("--subject-id", type=int, help="ID da disciplina para filtrar pendencias.")
    parser.add_argument("--state-id", type=int, help="ID do estado para filtrar pendencias.")
    parser.add_argument("--dry-run", action="store_true", help="Nao aplica PATCH, apenas simula.")
    parser.add_argument(
        "--note",
        default="Aprovacao em lote via script.",
        help="Observacao de revisao para registrar nas provas aprovadas.",
    )
    return parser.parse_args()


def configure_logging() -> None:
    logging.basicConfig(level=logging.INFO, format="[%(levelname)s] %(message)s")


def ask_token(cli_token: str | None = None) -> str:
    if cli_token and cli_token.strip():
        return cli_token.strip()
    token = input("Informe o token Bearer (sem 'Bearer '): ").strip()
    if not token:
        raise RuntimeError("Token vazio. Informe um token valido.")
    return token


def ask_user_id(cli_user_id: int | None = None) -> int:
    if cli_user_id is not None:
        return cli_user_id
    raw = input("Informe o user-id alvo: ").strip()
    if not raw:
        raise RuntimeError("User-id vazio. Informe um id valido.")
    try:
        return int(raw)
    except ValueError as exc:
        raise RuntimeError("User-id invalido. Informe um numero inteiro.") from exc


def ask_required_int(prompt: str, value: int | None) -> int:
    if value is not None:
        return value
    raw = input(prompt).strip()
    if not raw:
        raise RuntimeError("Valor obrigatorio nao informado.")
    try:
        return int(raw)
    except ValueError as exc:
        raise RuntimeError("Valor invalido. Informe um numero inteiro.") from exc


def build_url(base_url: str, path: str) -> str:
    return base_url.rstrip("/") + path


def api_request(
    base_url: str,
    path: str,
    token: str,
    method: str = "GET",
    payload: dict[str, Any] | None = None,
) -> Any:
    data = None
    headers = {"Authorization": f"Bearer {token}", "Accept": "application/json"}

    if payload is not None:
        data = json.dumps(payload).encode("utf-8")
        headers["Content-Type"] = "application/json"

    req = request.Request(build_url(base_url, path), data=data, method=method, headers=headers)
    try:
        with request.urlopen(req, timeout=60) as response:
            body = response.read()
            if not body:
                return None
            return json.loads(body.decode("utf-8"))
    except error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="replace")
        raise RuntimeError(f"Falha HTTP {exc.code} em {method} {path}: {details}") from exc
    except error.URLError as exc:
        raise RuntimeError(f"Nao foi possivel conectar em {base_url}: {exc.reason}") from exc


def resolve_target_user(
    base_url: str,
    token: str,
    user_id: int | None,
    email: str | None,
) -> tuple[int, str]:
    users = api_request(base_url, "/users", token)
    if not isinstance(users, list):
        raise RuntimeError("Resposta invalida ao listar usuarios.")

    if user_id is None and (email is None or not email.strip()):
        user_id = ask_user_id(None)

    if user_id is not None:
        user = next((u for u in users if int(u.get("id", -1)) == user_id), None)
        if user is None:
            raise RuntimeError(f"Usuario alvo nao encontrado para id={user_id}.")
        return int(user["id"]), str(user.get("email") or "")

    target_email = email.strip().lower()
    user = next((u for u in users if str(u.get("email") or "").strip().lower() == target_email), None)
    if user is None:
        raise RuntimeError(f"Usuario alvo nao encontrado para email={email}.")
    return int(user["id"]), str(user.get("email") or "")


def main() -> int:
    configure_logging()
    args = parse_args()
    token = ask_token(args.token)

    current_user = api_request(args.base_url, "/users/me", token)
    current_role = str(current_user.get("role") or "")
    if current_role not in ALLOWED_ROLES:
        raise RuntimeError("Token precisa ser ADMIN, DEV ou APPROVER para aprovar provas.")

    target_user_id, target_user_email = resolve_target_user(args.base_url, token, args.user_id, args.email)
    state_id = ask_required_int("Informe o state-id: ", args.state_id)
    university_id = ask_required_int("Informe o university-id: ", args.university_id)
    course_id = ask_required_int("Informe o course-id: ", args.course_id)
    subject_id = ask_required_int("Informe o subject-id: ", args.subject_id)
    pending = api_request(
        args.base_url,
        f"/exams/pending?stateId={state_id}&universityId={university_id}&courseId={course_id}&subjectId={subject_id}",
        token,
    )
    if not isinstance(pending, list):
        raise RuntimeError("Resposta invalida ao listar provas pendentes.")

    target_pending = [
        exam for exam in pending
        if isinstance(exam, dict)
        and isinstance(exam.get("uploadedBy"), dict)
        and int(exam["uploadedBy"].get("id", -1)) == target_user_id
    ]

    approved = 0
    failed = 0
    skipped = 0

    for exam in target_pending:
        exam_id = exam.get("id")
        if exam_id is None:
            skipped += 1
            continue

        if args.dry_run:
            approved += 1
            logging.info("[DRY-RUN] aprovaria examId=%s name=%s", exam_id, exam.get("name"))
            continue

        try:
            api_request(
                args.base_url,
                f"/exams/{int(exam_id)}/status",
                token,
                method="PATCH",
                payload={"status": "APPROVED", "reviewNote": args.note},
            )
            approved += 1
            logging.info("Aprovada examId=%s name=%s", exam_id, exam.get("name"))
        except RuntimeError as exc:
            failed += 1
            logging.error("Falha ao aprovar examId=%s: %s", exam_id, exc)

    print(
        json.dumps(
            {
                "targetUserId": target_user_id,
                "targetUserEmail": target_user_email,
                "stateId": state_id,
                "universityId": university_id,
                "courseId": course_id,
                "subjectId": subject_id,
                "dryRun": args.dry_run,
                "pendingFound": len(target_pending),
                "approved": approved,
                "failed": failed,
                "skipped": skipped,
            },
            indent=2,
            ensure_ascii=True,
        )
    )
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
