import uuid
from decimal import Decimal
from dataclasses import dataclass
from typing import Optional

@dataclass
class Book:
    id: Optional[uuid.UUID]
    title: str
    author: str
    price: Decimal
    description: str | None
    # cover_image: str
    category: str