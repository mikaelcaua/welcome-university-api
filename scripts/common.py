#!/usr/bin/env python3
"""
Constantes e utilitarios compartilhados para scripts.
"""

from __future__ import annotations

import os


API_BASE_URL = "https://welcome-university.duckdns.org"


def get_api_base_url() -> str:
    return os.getenv("API_BASE_URL", API_BASE_URL)
