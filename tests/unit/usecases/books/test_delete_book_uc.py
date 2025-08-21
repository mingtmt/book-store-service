# tests/usecases/books/test_delete_book_uc.py
import uuid
import pytest
from decimal import Decimal

from app.domain.entities.book import Book
from app.domain.errors import BookNotFound
from app.usecases.books.create_book import CreateBookUseCase, CreateBookCommand
from app.usecases.books.delete_book import DeleteBookUseCase


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


@pytest.mark.unit
def test_delete_book_success(fake_book_repo):
    book = seed_book(fake_book_repo)
    uc = DeleteBookUseCase(fake_book_repo)

    result = uc.execute(book.id)

    assert result is None
    assert fake_book_repo.get_by_id(book.id) is None


@pytest.mark.unit
def test_delete_book_not_found_raises(fake_book_repo):
    uc = DeleteBookUseCase(fake_book_repo)
    missing_id = uuid.uuid4()

    with pytest.raises(BookNotFound) as excinfo:
        uc.execute(missing_id)

    assert excinfo.value.context.get("book_id") == str(missing_id)


@pytest.mark.unit
def test_delete_book_twice_second_raises_not_found(fake_book_repo):
    book = seed_book(fake_book_repo)
    uc = DeleteBookUseCase(fake_book_repo)
    uc.execute(book.id)

    with pytest.raises(BookNotFound):
        uc.execute(book.id)
