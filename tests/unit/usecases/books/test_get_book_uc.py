import uuid
from decimal import Decimal

import pytest

from app.domain.entities.book import Book
from app.domain.errors import BookNotFound
from app.usecases.books.create_book import CreateBookCommand, CreateBookUseCase
from app.usecases.books.get_book import GetAllBooksUseCase, GetBookByIdUseCase


# ---------- GetBookByIdUseCase ----------
@pytest.mark.unit
def test_get_book_by_id_success(fake_book_repo, book_payload):
    create_uc = CreateBookUseCase(fake_book_repo)
    created = create_uc.execute(
        CreateBookCommand(
            title=book_payload["title"],
            author=book_payload["author"],
            price=book_payload["price"],  # e.g. Decimal("19.99")
            description=book_payload["description"],
            category=book_payload["category"],
        )
    )

    uc = GetBookByIdUseCase(fake_book_repo)
    got = uc.execute(created.id)

    assert isinstance(got, Book)
    assert got.id == created.id
    assert got.title == created.title
    assert got.author == created.author
    assert got.price == created.price
    assert got.description == created.description
    assert got.category == created.category


@pytest.mark.unit
def test_get_book_by_id_not_found_raises(fake_book_repo):
    uc = GetBookByIdUseCase(fake_book_repo)
    missing_id = uuid.uuid4()

    with pytest.raises(BookNotFound) as excinfo:
        uc.execute(missing_id)

    assert isinstance(excinfo.value, BookNotFound)
    assert excinfo.value.context.get("book_id") == str(missing_id)


# ---------- GetAllBooksUseCase ----------
@pytest.mark.unit
def test_get_all_books_returns_list(fake_book_repo):
    create_uc = CreateBookUseCase(fake_book_repo)
    b1 = create_uc.execute(
        CreateBookCommand(
            title="Domain-Driven Design",
            author="Eric Evans",
            price=Decimal("35.00"),
            description="Blue book",
            category="software",
        )
    )
    b2 = create_uc.execute(
        CreateBookCommand(
            title="Clean Architecture",
            author="Robert C. Martin",
            price=Decimal("19.99"),
            description="Uncle Bob",
            category="software",
        )
    )

    uc = GetAllBooksUseCase(fake_book_repo)
    items = uc.execute()

    assert isinstance(items, list)
    assert all(isinstance(x, Book) for x in items)
    ids = {x.id for x in items}
    assert b1.id in ids and b2.id in ids
    assert len(items) == 2


@pytest.mark.unit
def test_get_all_books_empty_list(fake_book_repo):
    uc = GetAllBooksUseCase(fake_book_repo)
    items = uc.execute()
    assert items == []
