import uuid
from dataclasses import dataclass
from typing import Optional


@dataclass
class User:
    id: Optional[uuid.UUID]
    email: str
    name: str
    age: int
    hashed_password: str
