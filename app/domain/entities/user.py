import uuid
from dataclasses import dataclass, field
from typing import Optional, Set

from app.domain.entities.role import Role


@dataclass
class User:
    id: Optional[uuid.UUID]
    email: str
    name: str
    age: int
    hashed_password: str
    role: Set[Role] = field(default_factory=set)
