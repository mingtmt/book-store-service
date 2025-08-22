import uuid
from typing import Dict, Optional

import pytest

from app.domain.entities.user import User
from app.domain.errors import ConstraintViolation
from app.utils.helper import normalize_email


class InMemoryUserRepo:
    """
    Fake IUserRepository for unit tests.
    Fields: id, email, name, age, hashed_password
    Constraints:
      - Unique email (case-insensitive)
      - 0 <= age <= 150  (aligns with your schema)
    """

    def __init__(self):
        self._data: Dict[uuid.UUID, User] = {}

    # --- helpers ---
    def _assert_unique_email(
        self, email: str, exclude_id: Optional[uuid.UUID] = None
    ) -> None:
        e = normalize_email(email)
        for uid, u in self._data.items():
            if exclude_id and uid == exclude_id:
                continue
            if normalize_email(u.email) == e:
                raise ConstraintViolation("Email already exists")

    def _ensure_age_ok(self, age: int) -> None:
        if age is None:
            return
        if not (0 <= int(age) <= 150):
            raise ConstraintViolation("Age out of range")

    # --- interface methods ---
    def create(self, user: User) -> User:
        self._assert_unique_email(user.email)
        self._ensure_age_ok(user.age)

        new_id = getattr(user, "id", None) or uuid.uuid4()
        created = User(
            id=new_id,
            email=normalize_email(user.email),
            name=user.name,
            age=user.age,
            hashed_password=user.hashed_password,
        )
        self._data[new_id] = created
        return created

    def get_by_email(self, email: str) -> Optional[User]:
        e = normalize_email(email)
        for u in self._data.values():
            if normalize_email(u.email) == e:
                return u
        return None


@pytest.fixture
def fake_user_repo():
    return InMemoryUserRepo()
