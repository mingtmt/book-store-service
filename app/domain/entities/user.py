import uuid
from dataclasses import dataclass
from typing import Optional

@dataclass
class User:
    id: Optional[uuid.UUID]
    email: str
    hashed_password: str
