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
    "reposicao": "RECUPERACAO",
    "recuperacao": "RECUPERACAO",
    "final": "FINAL",
}
PERIOD_REGEX = re.compile(r"(20\d{2})[._\-\s](1|2)\b")


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
        raise RuntimeError(f"Falha HTTP {exc.code} em {method} {path}: {details}") from exc
    except error.URLError as exc:
        raise RuntimeError(f"Nao foi possivel conectar em {base_url}: {exc.reason}") from exc


def upload_exam_multipart(
    base_url: str,
    token: str,
    exam_year: int,
    semester: int,
    exam_type: str,
    subject_id: int,
    file_path: Path,
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
    }
    req = request.Request(
        build_url(base_url, "/exams"),
        data=body,
        method="POST",
        headers=headers,
    )

    try:
        with request.urlopen(req) as response:
            return int(getattr(response, "status", 201)), response.read().decode("utf-8", errors="replace")
    except error.HTTPError as exc:
        details = exc.read().decode("utf-8", errors="replace")
        return exc.code, details


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


def parse_period(file_name: str) -> tuple[int, int] | None:
    stem = Path(file_name).stem
    match = PERIOD_REGEX.search(stem)
    if not match:
        return None
    return int(match.group(1)), int(match.group(2))


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

    current_user = api_request(args.base_url, "/users/me", token)
    if current_user.get("role") != "DEV":
        raise RuntimeError("Este script so pode ser executado com token de usuario DEV.")

    maranhao = api_request(args.base_url, "/states/MA", token)
    universities = api_request(args.base_url, f"/states/{maranhao['id']}/universities", token)
    ufma = next((u for u in universities if (u.get("name") or "").strip().lower() == "ufma"), None)
    if ufma is None:
        raise RuntimeError("UFMA nao encontrada. Rode o script de alimentacao inicial primeiro.")

    courses = api_request(args.base_url, f"/universities/{ufma['id']}/courses", token)
    course = next((c for c in courses if normalize_text(c.get("name") or "") == normalize_text("Ciencia da Computacao")), None)
    if course is None:
        raise RuntimeError("Curso Ciencia da Computacao nao encontrado na UFMA.")

    subjects = api_request(args.base_url, f"/courses/{course['id']}/subjects", token)
    subject_by_name = {normalize_text(s["name"]): int(s["id"]) for s in subjects if s.get("name") and s.get("id")}

    folder_subjects = list_subject_folders(root)
    created_subjects = 0
    for folder_subject in folder_subjects:
        normalized = normalize_text(folder_subject)
        if normalized in subject_by_name:
            continue
        try:
            response = api_request(
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

    subjects = api_request(args.base_url, f"/courses/{course['id']}/subjects", token)
    subject_by_name = {normalize_text(s["name"]): int(s["id"]) for s in subjects if s.get("name") and s.get("id")}
    if not subject_by_name:
        raise RuntimeError("Nenhuma disciplina encontrada para o curso, mesmo apos tentativa de criacao automatica.")

    pending = api_request(args.base_url, "/exams/pending", token)
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
    seen_payload_keys: set[tuple[int, int, int, str]] = set()

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
            logging.warning("Tipo de prova nao reconhecido para pasta: %s", type_folder)
            unresolved_type += 1
            continue

        parsed_period = parse_period(file_path.name)
        if parsed_period is None:
            logging.warning("Periodo nao identificado no nome do arquivo: %s", file_path.name)
            unresolved_period += 1
            continue

        exam_year, semester = parsed_period
        payload_key = (subject_id, exam_year, semester, exam_type)
        if payload_key in seen_payload_keys:
            skipped += 1
            continue

        if subject_id not in approved_cache:
            approved = api_request(args.base_url, f"/exams?subjectId={subject_id}", token)
            approved_cache[subject_id] = {
                (int(item["subjectId"]), int(item["examYear"]), int(item["semester"]), str(item["type"]))
                for item in approved
                if item.get("subjectId") is not None
            }

        if payload_key in approved_cache[subject_id] or payload_key in pending_keys:
            skipped += 1
            seen_payload_keys.add(payload_key)
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
            seen_payload_keys.add(payload_key)
            continue

        status, response_body = upload_exam_multipart(
            args.base_url,
            token,
            exam_year,
            semester,
            exam_type,
            subject_id,
            file_path,
        )

        if status in {200, 201}:
            created += 1
            seen_payload_keys.add(payload_key)
            approved_cache[subject_id].add(payload_key)
            logging.info("Enviado: %s", relative)
        elif status == 409:
            skipped += 1
            seen_payload_keys.add(payload_key)
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
            },
            indent=2,
            ensure_ascii=True,
        )
    )
    return 0 if failed == 0 else 1


if __name__ == "__main__":
    raise SystemExit(main())
