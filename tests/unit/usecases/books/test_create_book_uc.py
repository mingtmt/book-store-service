import pytest
import uuid
from decimal import Decimal

from app.usecases.books.create_book import CreateBookUseCase, CreateBookCommand
from app.domain.entities.book import Book
from app.domain.errors import ConstraintViolation
from app.usecases.errors import BadRequest


@pytest.mark.unit
def test_create_book_success(fake_book_repo, book_payload):
    uc = CreateBookUseCase(fake_book_repo)

    cmd = CreateBookCommand(
        title=book_payload["title"],
        author=book_payload["author"],
        price=book_payload["price"],           # Decimal("19.99")
        description=book_payload["description"],
        category=book_payload["category"],
    )
    created = uc.execute(cmd)

    assert isinstance(created, Book)
    assert isinstance(created.id, uuid.UUID)
    assert created.title == book_payload["title"]
    assert created.author == book_payload["author"]
    assert created.price == Decimal("19.99")
    assert created.description == book_payload["description"]
    assert created.category == book_payload["category"]


@pytest.mark.unit
def test_create_book_trims_fields(fake_book_repo):
    uc = CreateBookUseCase(fake_book_repo)

    cmd = CreateBookCommand(
        title="   DDD Made Simple   ",
        author="  Jane   Doe ",
        price=Decimal("25.50"),
        description="  Intro to DDD  ",
        category="  software ",
    )
    created = uc.execute(cmd)

    assert created.title == "DDD Made Simple"
    assert created.author == "Jane   Doe"      # UC chỉ strip, không normalize khoảng trắng giữa
    assert created.description == "Intro to DDD"
    assert created.category == "software"


@pytest.mark.unit
def test_create_book_duplicate_title_author_raises(fake_book_repo):
    uc = CreateBookUseCase(fake_book_repo)

    first = CreateBookCommand(
        title="Domain-Driven Design",
        author="Eric Evans",
        price=Decimal("35.00"),
        description="Blue book",
        category="software",
    )
    uc.execute(first)

    # Trùng (title, author) → repo raise ConstraintViolation
    dup = CreateBookCommand(
        title="Domain-Driven Design",
        author="Eric Evans",
        price=Decimal("36.00"),
        description="Another desc",
        category="software",
    )
    with pytest.raises(ConstraintViolation):
        uc.execute(dup)


@pytest.mark.unit
def test_create_book_negative_price_raises_bad_request(fake_book_repo):
    uc = CreateBookUseCase(fake_book_repo)

    cmd = CreateBookCommand(
        title="Clean Architecture",
        author="Robert C. Martin",
        price=Decimal("-0.01"),  # UC tự check -> BadRequest
        description="Uncle Bob",
        category="software",
    )
    with pytest.raises(BadRequest):
        uc.execute(cmd)


@pytest.mark.unit
def test_create_book_empty_title_raises_bad_request(fake_book_repo):
    uc = CreateBookUseCase(fake_book_repo)
    cmd = CreateBookCommand(
        title="   ",
        author="Someone",
        price=Decimal("10.00"),
        description=None,
        category="software",
    )
    with pytest.raises(BadRequest):
        uc.execute(cmd)


@pytest.mark.unit
def test_create_book_empty_author_raises_bad_request(fake_book_repo):
    uc = CreateBookUseCase(fake_book_repo)
    cmd = CreateBookCommand(
        title="Some Title",
        author="   ",
        price=Decimal("10.00"),
        description=None,
        category="software",
    )
    with pytest.raises(BadRequest):
        uc.execute(cmd)
