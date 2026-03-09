#!/usr/bin/env python3
"""
Importa provas de courses_databases/ufma/ciencia_computacao para a API.

Regras:
- tenta evitar reenvio usando metadados existentes;
- se a API responder conflito de duplicidade, marca como ignorado.
"""

from __future__ import annotations

import argparse
import json
import logging
import mimetypes
import re
import socket
import time
import unicodedata
import uuid
from pathlib import Path
from typing import Any
from urllib import error, request

from common import get_api_base_url


SUPPORTED_EXTENSIONS = {".pdf", ".png", ".jpg", ".jpeg", ".webp", ".gif"}
TYPE_MAP: dict[str, str] = {
    "prova1": "PROVA1",
    "p1": "PROVA1",
    "prova2": "PROVA2",
    "p2": "PROVA2",
    "prova3": "PROVA3",
    "p3": "PROVA3",
    "prova4": "FINAL",
    "p4": "FINAL",
    "provax": "FINAL",
    "px": "FINAL",
    "reposicao": "RECUPERACAO",
    "recuperacao": "RECUPERACAO",
    "final": "FINAL",
}
PERIOD_REGEX = re.compile(r"(?<!\d)(20\d{2})[._\-\s](1|2)(?!\d)")
YEAR_ONLY_REGEX = re.compile(r"(?<!\d)((?:19|20)\d{2})(?!\d)")
UNKNOWN_PERIOD_YEAR = 9999
UNKNOWN_PERIOD_SEMESTER = 1
MAX_REAUTH_403_ATTEMPTS = 3
FILE_TOKEN_TYPE_MAP: dict[str, str] = {
    "GB": "PROVA1",
    "PR": "PROVA2",
    "PV": "PROVA3",
}


class ApiHttpError(RuntimeError):
    def __init__(self, code: int, method: str, path: str, details: str):
        super().__init__(f"Falha HTTP {code} em {method} {path}: {details}")
        self.code = code
        self.method = method
        self.path = path
        self.details = details


def parse_args() -> argparse.Namespace:
    parser = argparse.ArgumentParser(description="Importa provas da UFMA para a API.")
    parser.add_argument("--base-url", default=get_api_base_url())
    parser.add_argument("--token", help="Token Bearer (sem 'Bearer ').")
    parser.add_argument(
        "--root",
        default="courses_databases/ufma/ciencia_computacao",
        help="Pasta raiz com as disciplinas.",
    )
    parser.add_argument("--dry-run", action="store_true", help="Nao envia, apenas simula.")
    parser.add_argument("--limit", type=int, help="Limite maximo de arquivos para processar.")
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


def ask_new_token_due_forbidden(previous_token: str | None = None) -> str:
    token = input("API retornou 403. Informe outro token Bearer (sem 'Bearer '): ").strip()
    if not token:
        raise RuntimeError("Token vazio. Informe um token valido.")
    if previous_token and token == previous_token:
        raise RuntimeError(
            "Mesmo token informado novamente apos 403. Gere um novo access token e tente de novo."
        )
    return token


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
        with request.urlopen(req) as response:
            body = response.read()
            if not body:
                return None
            return json.loads(body.decode("utf-8"))
    except error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="replace")
        raise ApiHttpError(exc.code, method, path, details) from exc
    except error.URLError as exc:
        raise RuntimeError(f"Nao foi possivel conectar em {base_url}: {exc.reason}") from exc


def api_request_with_reauth(
    base_url: str,
    path: str,
    token: str,
    method: str = "GET",
    payload: dict[str, Any] | None = None,
) -> tuple[Any, str]:
    current_token = token
    forbidden_count = 0
    while True:
        try:
            response = api_request(base_url, path, current_token, method=method, payload=payload)
            return response, current_token
        except ApiHttpError as exc:
            if exc.code != 403:
                raise
            forbidden_count += 1
            if forbidden_count >= MAX_REAUTH_403_ATTEMPTS:
                raise RuntimeError(
                    f"Recebido 403 {forbidden_count} vezes em {method} {path}. "
                    "Interrompendo para evitar loop de reautenticacao."
                ) from exc
            logging.warning(
                "Recebido 403 em %s %s. Resposta: %s. Solicitando novo token...",
                method,
                path,
                (exc.details or "").strip()[:300],
            )
            current_token = ask_new_token_due_forbidden(previous_token=current_token)


def upload_exam_multipart(
    base_url: str,
    token: str,
    exam_year: int,
    semester: int,
    exam_type: str,
    subject_id: int,
    file_path: Path,
    timeout_seconds: int = 60,
    max_retries: int = 3,
) -> tuple[int, str]:
    boundary = f"----CodexBoundary{uuid.uuid4().hex}"
    content_type = mimetypes.guess_type(file_path.name)[0] or "application/octet-stream"
    file_bytes = file_path.read_bytes()

    parts: list[bytes] = []
    fields = {
        "examYear": str(exam_year),
        "semester": str(semester),
        "type": exam_type,
        "subjectId": str(subject_id),
    }

    for key, value in fields.items():
        parts.append(
            (
                f"--{boundary}\r\n"
                f'Content-Disposition: form-data; name="{key}"\r\n\r\n'
                f"{value}\r\n"
            ).encode("utf-8")
        )

    parts.append(
        (
            f"--{boundary}\r\n"
            f'Content-Disposition: form-data; name="file"; filename="{file_path.name}"\r\n'
            f"Content-Type: {content_type}\r\n\r\n"
        ).encode("utf-8")
    )
    parts.append(file_bytes)
    parts.append(f"\r\n--{boundary}--\r\n".encode("utf-8"))
    body = b"".join(parts)

    headers = {
        "Authorization": f"Bearer {token}",
        "Accept": "application/json",
        "Content-Type": f"multipart/form-data; boundary={boundary}",
        "Connection": "close",
    }
    req = request.Request(
        build_url(base_url, "/exams"),
        data=body,
        method="POST",
        headers=headers,
    )

    attempt = 0
    while True:
        attempt += 1
        try:
            with request.urlopen(req, timeout=timeout_seconds) as response:
                return int(getattr(response, "status", 201)), response.read().decode("utf-8", errors="replace")
        except error.HTTPError as exc:
            details = exc.read().decode("utf-8", errors="replace")
            return exc.code, details
        except (error.URLError, TimeoutError, socket.timeout, ConnectionResetError) as exc:
            if attempt >= max_retries:
                return 0, f"erro de conexao apos {attempt} tentativas: {exc}"
            backoff = min(8, 2 ** (attempt - 1))
            logging.warning(
                "Falha de conexao no upload (%s), tentativa %s/%s. Aguardando %ss...",
                file_path.name,
                attempt,
                max_retries,
                backoff,
            )
            time.sleep(backoff)


def upload_exam_multipart_with_reauth(
    base_url: str,
    token: str,
    exam_year: int,
    semester: int,
    exam_type: str,
    subject_id: int,
    file_path: Path,
    timeout_seconds: int = 60,
    max_retries: int = 3,
) -> tuple[int, str, str]:
    current_token = token
    forbidden_count = 0
    while True:
        status, response_body = upload_exam_multipart(
            base_url,
            current_token,
            exam_year,
            semester,
            exam_type,
            subject_id,
            file_path,
            timeout_seconds=timeout_seconds,
            max_retries=max_retries,
        )
        if status != 403:
            return status, response_body, current_token
        forbidden_count += 1
        if forbidden_count >= MAX_REAUTH_403_ATTEMPTS:
            return (
                403,
                f"Recebido 403 {forbidden_count} vezes no upload de {file_path.name}. "
                f"Ultima resposta: {(response_body or '').strip()[:300]}. "
                "Interrompendo para evitar loop de reautenticacao.",
                current_token,
            )
        logging.warning(
            "Recebido 403 no upload de %s. Resposta: %s. Solicitando novo token...",
            file_path.name,
            (response_body or "").strip()[:300],
        )
        current_token = ask_new_token_due_forbidden(previous_token=current_token)


def normalize_text(value: str) -> str:
    normalized = unicodedata.normalize("NFKD", value)
    ascii_only = "".join(ch for ch in normalized if not unicodedata.combining(ch))
    cleaned = re.sub(r"[^a-zA-Z0-9]+", "", ascii_only).lower()
    return cleaned


def resolve_exam_type(folder_name: str) -> str | None:
    normalized = normalize_text(folder_name)
    if normalized in TYPE_MAP:
        return TYPE_MAP[normalized]
    if normalized.startswith("prova1"):
        return "PROVA1"
    if normalized.startswith("prova2"):
        return "PROVA2"
    if normalized.startswith("prova3"):
        return "PROVA3"
    return None


def resolve_exam_type_from_file_name(file_name: str) -> str | None:
    stem = Path(file_name).stem.upper()
    tokens = [token for token in re.split(r"[^A-Z0-9]+", stem) if token]
    for token in tokens:
        mapped = FILE_TOKEN_TYPE_MAP.get(token)
        if mapped:
            return mapped
    return None


def parse_period(file_name: str) -> tuple[int, int] | None:
    stem = Path(file_name).stem
    match = PERIOD_REGEX.search(stem)
    if match:
        return int(match.group(1)), int(match.group(2))
    year_only = YEAR_ONLY_REGEX.search(stem)
    if year_only:
        return int(year_only.group(1)), 1
    return None


def list_exam_files(root: Path) -> list[Path]:
    return sorted(
        p for p in root.rglob("*")
        if p.is_file() and p.suffix.lower() in SUPPORTED_EXTENSIONS
    )


def list_subject_folders(root: Path) -> list[str]:
    return sorted(p.name for p in root.iterdir() if p.is_dir())


def main() -> int:
    configure_logging()
    args = parse_args()
    root = Path(args.root)
    token = ask_token(args.token)

    if not root.exists() or not root.is_dir():
        raise RuntimeError(f"Pasta invalida: {root}")

    current_user, token = api_request_with_reauth(args.base_url, "/users/me", token)

    maranhao, token = api_request_with_reauth(args.base_url, "/states/MA", token)
    universities, token = api_request_with_reauth(args.base_url, f"/states/{maranhao['id']}/universities", token)
    ufma = next((u for u in universities if (u.get("name") or "").strip().lower() == "ufma"), None)
    if ufma is None:
        raise RuntimeError("UFMA nao encontrada. Rode o script de alimentacao inicial primeiro.")

    courses, token = api_request_with_reauth(args.base_url, f"/universities/{ufma['id']}/courses", token)
    course = next((c for c in courses if normalize_text(c.get("name") or "") == normalize_text("Ciencia da Computacao")), None)
    if course is None:
        raise RuntimeError("Curso Ciencia da Computacao nao encontrado na UFMA.")

    subjects, token = api_request_with_reauth(args.base_url, f"/courses/{course['id']}/subjects", token)
    subject_by_name = {normalize_text(s["name"]): int(s["id"]) for s in subjects if s.get("name") and s.get("id")}

    folder_subjects = list_subject_folders(root)
    created_subjects = 0
    for folder_subject in folder_subjects:
        normalized = normalize_text(folder_subject)
        if normalized in subject_by_name:
            continue
        try:
            response, token = api_request_with_reauth(
                args.base_url,
                f"/courses/{course['id']}/subjects",
                token,
                method="POST",
                payload={"name": folder_subject},
            )
            if isinstance(response, dict) and response.get("id") and response.get("name"):
                subject_by_name[normalize_text(str(response["name"]))] = int(response["id"])
                created_subjects += 1
                logging.info("Disciplina criada: %s", response["name"])
        except RuntimeError as exc:
            if "HTTP 409" in str(exc):
                logging.info("Disciplina ja existente (409): %s", folder_subject)
                continue
            raise

    subjects, token = api_request_with_reauth(args.base_url, f"/courses/{course['id']}/subjects", token)
    subject_by_name = {normalize_text(s["name"]): int(s["id"]) for s in subjects if s.get("name") and s.get("id")}
    if not subject_by_name:
        raise RuntimeError("Nenhuma disciplina encontrada para o curso, mesmo apos tentativa de criacao automatica.")

    pending, token = api_request_with_reauth(args.base_url, "/exams/pending", token)
    pending_keys = {
        (int(item["subjectId"]), int(item["examYear"]), int(item["semester"]), str(item["type"]))
        for item in pending
        if item.get("subjectId") is not None and item.get("examYear") is not None and item.get("semester") is not None
    }

    exam_files = list_exam_files(root)
    if args.limit is not None:
        exam_files = exam_files[: max(args.limit, 0)]

    created = 0
    skipped = 0
    failed = 0
    unresolved_subject = 0
    unresolved_type = 0
    unresolved_period = 0
    uploaded_unknown_period = 0
    seen_payload_keys: set[tuple] = set()

    approved_cache: dict[int, set[tuple[int, int, int, str]]] = {}

    for file_path in exam_files:
        relative = file_path.relative_to(root)
        parts = relative.parts
        if len(parts) < 3:
            skipped += 1
            continue

        subject_folder = parts[0]
        type_folder = parts[1]
        subject_key = normalize_text(subject_folder)
        subject_id = subject_by_name.get(subject_key)
        if subject_id is None:
            logging.warning("Disciplina nao encontrada na API para pasta: %s", subject_folder)
            unresolved_subject += 1
            continue

        exam_type = resolve_exam_type(type_folder)
        if exam_type is None:
            exam_type = resolve_exam_type_from_file_name(file_path.name)
        if exam_type is None:
            logging.warning("Tipo de prova nao reconhecido para pasta: %s", type_folder)
            unresolved_type += 1
            continue

        parsed_period = parse_period(file_path.name)
        unknown_period = parsed_period is None
        if unknown_period:
            exam_year = UNKNOWN_PERIOD_YEAR
            semester = UNKNOWN_PERIOD_SEMESTER
            unresolved_period += 1
            logging.warning(
                "Periodo nao identificado: %s -> usando fallback %s.%s (periodo nao identificado)",
                file_path.name,
                exam_year,
                semester,
            )
        else:
            exam_year, semester = parsed_period
        payload_key = (subject_id, exam_year, semester, exam_type)
        if unknown_period:
            seen_key = (
                subject_id,
                "UNKNOWN_PERIOD",
                exam_type,
                file_path.name.lower(),
                file_path.stat().st_size,
            )
        else:
            seen_key = payload_key

        if seen_key in seen_payload_keys:
            skipped += 1
            continue

        if subject_id not in approved_cache:
            approved, token = api_request_with_reauth(args.base_url, f"/exams?subjectId={subject_id}", token)
            approved_cache[subject_id] = {
                (int(item["subjectId"]), int(item["examYear"]), int(item["semester"]), str(item["type"]))
                for item in approved
                if item.get("subjectId") is not None
            }

        if (not unknown_period) and (payload_key in approved_cache[subject_id] or payload_key in pending_keys):
            skipped += 1
            seen_payload_keys.add(seen_key)
            continue

        if args.dry_run:
            logging.info(
                "[DRY-RUN] enviaria %s | subjectId=%s type=%s period=%s.%s",
                relative,
                subject_id,
                exam_type,
                exam_year,
                semester,
            )
            created += 1
            seen_payload_keys.add(seen_key)
            continue

        status, response_body, token = upload_exam_multipart_with_reauth(
            args.base_url,
            token,
            exam_year,
            semester,
            exam_type,
            subject_id,
            file_path,
            timeout_seconds=120,
            max_retries=4,
        )

        if status in {200, 201}:
            created += 1
            seen_payload_keys.add(seen_key)
            if not unknown_period:
                approved_cache[subject_id].add(payload_key)
            else:
                uploaded_unknown_period += 1
            logging.info("Enviado: %s", relative)
        elif status == 409:
            skipped += 1
            seen_payload_keys.add(seen_key)
            logging.info("Ignorado duplicado (409): %s", relative)
        else:
            failed += 1
            logging.error("Falha no upload (%s): %s | resposta=%s", status, relative, response_body)

    print(
        json.dumps(
            {
                "root": str(root),
                "dryRun": args.dry_run,
                "processedFiles": len(exam_files),
                "createdSubjects": created_subjects,
                "created": created,
                "skipped": skipped,
                "failed": failed,
                "unresolvedSubject": unresolved_subject,
                "unresolvedType": unresolved_type,
                "unresolvedPeriod": unresolved_period,
                "uploadedUnknownPeriod": uploaded_unknown_period,
            },
            indent=2,
            ensure_ascii=True,
        )
    )
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
