# tests/usecases/books/test_update_book_uc.py
import uuid
from decimal import Decimal
import pytest

from app.domain.entities.book import Book
from app.domain.errors import BookNotFound, ConstraintViolation
from app.usecases.books.create_book import CreateBookUseCase, CreateBookCommand
from app.usecases.books.update_book import UpdateBookUseCase, UpdateBookCommand


def seed_book(repo, **overrides) -> Book:
    create_uc = CreateBookUseCase(repo)
    cmd = CreateBookCommand(
        title=overrides.get("title", "Domain-Driven Design"),
        author=overrides.get("author", "Eric Evans"),
        price=overrides.get("price", Decimal("35.00")),
        description=overrides.get("description", "Blue book"),
        category=overrides.get("category", "software"),
    )
    return create_uc.execute(cmd)


# ---------- Success paths ----------
@pytest.mark.unit
def test_update_book_partial_success_keeps_unspecified_fields(fake_book_repo):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    updated = uc.execute(
        existing.id,
        UpdateBookCommand(
            title="DDD Distilled",
            price=Decimal("29.99"),
        ),
    )

    assert updated.id == existing.id
    assert updated.title == "DDD Distilled"
    assert updated.price == Decimal("29.99")
    assert updated.author == existing.author
    assert updated.description == existing.description
    assert updated.category == existing.category


@pytest.mark.unit
def test_update_book_trims_and_updates_fields(fake_book_repo):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    updated = uc.execute(
        existing.id,
        UpdateBookCommand(
            title="   Clean Architecture   ",
            author="  Robert   C.   Martin ",
            description="  Uncle Bob  ",
            category="  software ",
        ),
    )

    assert updated.title == "Clean Architecture"
    assert updated.author == "Robert   C.   Martin"
    assert updated.description == "Uncle Bob"
    assert updated.category == "software"


@pytest.mark.unit
def test_update_book_no_changes_returns_same_values(fake_book_repo):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    updated = uc.execute(existing.id, UpdateBookCommand())

    assert updated.id == existing.id
    assert updated.title == existing.title
    assert updated.author == existing.author
    assert updated.price == existing.price
    assert updated.description == existing.description
    assert updated.category == existing.category


# ---------- Error paths ----------
@pytest.mark.unit
def test_update_book_not_found_raises(fake_book_repo):
    uc = UpdateBookUseCase(fake_book_repo)
    missing_id = uuid.uuid4()

    with pytest.raises(BookNotFound) as excinfo:
        uc.execute(missing_id, UpdateBookCommand(title="Whatever"))

    assert excinfo.value.context.get("book_id") == str(missing_id)


@pytest.mark.unit
def test_update_book_empty_title_raises_constraint_violation(fake_book_repo):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    with pytest.raises(ConstraintViolation):
        uc.execute(existing.id, UpdateBookCommand(title="   "))


@pytest.mark.unit
def test_update_book_empty_author_raises_constraint_violation(fake_book_repo):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    with pytest.raises(ConstraintViolation):
        uc.execute(existing.id, UpdateBookCommand(author="   "))


@pytest.mark.unit
@pytest.mark.parametrize("bad_price", [Decimal("0.00"), Decimal("-0.01")])
def test_update_book_non_positive_price_raises(fake_book_repo, bad_price):
    existing = seed_book(fake_book_repo)
    uc = UpdateBookUseCase(fake_book_repo)

    with pytest.raises(ConstraintViolation):
        uc.execute(existing.id, UpdateBookCommand(price=bad_price))


@pytest.mark.unit
def test_update_book_to_duplicate_title_author_raises(fake_book_repo):
    b1 = seed_book(fake_book_repo, title="Refactoring", author="Martin Fowler")
    b2 = seed_book(fake_book_repo, title="Patterns of Enterprise", author="Martin Fowler")
    uc = UpdateBookUseCase(fake_book_repo)

    with pytest.raises(ConstraintViolation):
        uc.execute(
            b2.id,
            UpdateBookCommand(title="Refactoring", author="Martin Fowler"),
        )
