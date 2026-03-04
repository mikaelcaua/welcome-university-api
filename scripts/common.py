#!/usr/bin/env python3
"""
Constantes e utilitarios compartilhados para scripts.
"""

from __future__ import annotations

import os


API_BASE_HOST = "http://mikael-dev"
API_BASE_PORT = 8081
API_BASE_URL = f"{API_BASE_HOST}:{API_BASE_PORT}"


def get_api_base_url() -> str:
    return os.getenv("API_BASE_URL", API_BASE_URL)
