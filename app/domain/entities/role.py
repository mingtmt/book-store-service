import uuid
from dataclasses import dataclass
from typing import Optional


@dataclass(frozen=True)
class Role:
    id: Optional[uuid.UUID]
    name: str
