import uuid
from decimal import Decimal

import pytest


@pytest.fixture
def money():
    def _m(x) -> Decimal:
        return Decimal(str(x))

    return _m


@pytest.fixture
def random_email():
    return f"u_{uuid.uuid4().hex[:8]}@example.com"
