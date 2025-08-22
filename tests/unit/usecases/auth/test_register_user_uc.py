import uuid

import pytest

from app.domain.entities.user import User
from app.domain.errors import EmailAlreadyExists
from app.usecases.auth.register_user import RegisterUserCommand, RegisterUserUseCase


class FakeTokenService:
    def __init__(self):
        self.calls: list[str] = []

    def create_access_token(self, sub: str) -> str:
        self.calls.append(sub)
        return f"jwt:{sub}"


class FakePasswordService:
    def __init__(self):
        self.seen: str | None = None

    def hash_password(self, raw: str) -> str:
        self.seen = raw
        return "HASHED"

    def verify_password(self, raw: str, hashed: str) -> bool:
        return hashed == "HASHED"


@pytest.mark.unit
def test_register_user_success(fake_user_repo):
    token_service = FakeTokenService()
    password_service = FakePasswordService()

    cmd = RegisterUserCommand(
        email="  Ada.Lovelace@Example.com  ",  # messy to test normalization
        name="Ada Lovelace",
        age=28,
        password="S3cureP@ssw0rd",
    )
    uc = RegisterUserUseCase(
        fake_user_repo, token_service, password_service, issue_token=True
    )

    created, token = uc.execute(cmd)

    assert isinstance(created, User)
    assert isinstance(created.id, uuid.UUID)
    assert created.email == "ada.lovelace@example.com"  # normalized
    assert created.name == "Ada Lovelace"
    assert created.age == 28
    assert created.hashed_password == "HASHED"

    # token issued for created user
    assert token == f"jwt:{created.id}"
    assert token_service.calls == [str(created.id)]

    # repo lookup is case-insensitive
    assert fake_user_repo.get_by_email("ADA.LOVELACE@EXAMPLE.COM").id == created.id


@pytest.mark.unit
def test_register_user_duplicate_email_raises(fake_user_repo):
    # seed an existing user
    existing = fake_user_repo.create(
        User(
            id=None,
            email="taken@example.com",
            name="Existing",
            age=30,
            hashed_password="X",
        )
    )

    token_service = FakeTokenService()
    password_service = FakePasswordService()

    cmd = RegisterUserCommand(
        email="TAKEN@example.com",  # same email, different case
        name="Dup",
        age=20,
        password="pw",
    )
    uc = RegisterUserUseCase(
        fake_user_repo, token_service, password_service, issue_token=True
    )

    with pytest.raises(EmailAlreadyExists):
        uc.execute(cmd)

    # Ensure repository still only resolves the seeded user for both cases
    assert fake_user_repo.get_by_email("taken@example.com").id == existing.id
    assert fake_user_repo.get_by_email("TAKEN@example.com").id == existing.id
    assert token_service.calls == []


@pytest.mark.unit
def test_register_user_without_issuing_token(fake_user_repo):
    token_service = FakeTokenService()
    password_service = FakePasswordService()

    cmd = RegisterUserCommand(
        email="user@example.com",
        name="User",
        age=21,
        password="secret",
    )
    uc = RegisterUserUseCase(
        fake_user_repo, token_service, password_service, issue_token=False
    )

    created, token = uc.execute(cmd)

    assert created.email == "user@example.com"
    assert token is None
    assert token_service.calls == []


@pytest.mark.unit
def test_password_hash_called_with_plain_password(fake_user_repo):
    token_service = FakeTokenService()
    password_service = FakePasswordService()

    cmd = RegisterUserCommand(
        email="hash@test.io",
        name="Hashy",
        age=30,
        password="S3cureP@ssw0rd",
    )
    uc = RegisterUserUseCase(
        fake_user_repo, token_service, password_service, issue_token=True
    )
    uc.execute(cmd)

    assert password_service.seen == "S3cureP@ssw0rd"  # raw password passed to hasher


@pytest.mark.unit
def test_email_is_normalized_before_repo_lookup(fake_user_repo, monkeypatch):
    token_service = FakeTokenService()
    password_service = FakePasswordService()

    # spy on get_by_email arg
    calls = {}
    original_get_by_email = fake_user_repo.get_by_email

    def spy_get_by_email(email: str):
        calls["arg"] = email
        return original_get_by_email(email)

    monkeypatch.setattr(fake_user_repo, "get_by_email", spy_get_by_email)

    cmd = RegisterUserCommand(
        email="  MiXeD@Case.COM ",
        name="Mix",
        age=18,
        password="pw",
    )
    uc = RegisterUserUseCase(
        fake_user_repo, token_service, password_service, issue_token=True
    )
    created, _ = uc.execute(cmd)

    # UC should pass a normalized (lowercased, trimmed) email to repo.get_by_email
    assert calls["arg"] == "mixed@case.com"
    # Stored user is also normalized
    assert created.email == "mixed@case.com"
