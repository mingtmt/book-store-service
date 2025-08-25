import uuid

from fastapi import APIRouter, Depends, Response, status
from sqlalchemy.orm import Session

from app.infrastructure.db.sqlalchemy.repositories.book_impl import (
    SqlAlchemyBookRepository,
)
from app.presentation.http.dependencies.auth import require_auth
from app.presentation.http.dependencies.db import get_db
from app.presentation.http.schemas.base import Envelope
from app.presentation.http.schemas.books import BookOut, CreateBook, UpdateBook
from app.usecases.books.create_book import CreateBookCommand, CreateBookUseCase
from app.usecases.books.delete_book import DeleteBookUseCase
from app.usecases.books.get_book import GetAllBooksUseCase, GetBookByIdUseCase
from app.usecases.books.update_book import UpdateBookCommand, UpdateBookUseCase

router = APIRouter()


@router.post(
    "/",
    response_model=Envelope[BookOut],
    status_code=status.HTTP_201_CREATED,
    dependencies=[Depends(require_auth)],
)
def create(payload: CreateBook, response: Response, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = CreateBookUseCase(repo)
    cmd = CreateBookCommand(
        title=payload.title,
        author=payload.author,
        price=payload.price,
        description=payload.description,
        category=payload.category,
    )
    created = uc.execute(cmd)
    response.headers["Location"] = f"/api/v1/books/{created.id}"
    return Envelope(data=BookOut.model_validate(created))


@router.get(
    "/{book_id}",
    response_model=Envelope[BookOut],
    status_code=status.HTTP_200_OK,
)
def get_by_id(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetBookByIdUseCase(repo)
    book = uc.execute(book_id)
    return Envelope(data=BookOut.model_validate(book))


@router.get(
    "/",
    response_model=Envelope[list[BookOut]],
    status_code=status.HTTP_200_OK,
)
def get_all(db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetAllBooksUseCase(repo)
    books = uc.execute()
    return Envelope(data=[BookOut.model_validate(book) for book in books])


@router.put(
    "/{book_id}",
    response_model=Envelope[BookOut],
    status_code=status.HTTP_200_OK,
    dependencies=[Depends(require_auth)],
)
def update(book_id: uuid.UUID, payload: UpdateBook, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = UpdateBookUseCase(repo)
    cmd = UpdateBookCommand(**payload.model_dump(exclude_unset=True))
    updated_book = uc.execute(book_id, cmd)
    return Envelope(data=BookOut.model_validate(updated_book))


@router.delete(
    "/{book_id}",
    status_code=status.HTTP_204_NO_CONTENT,
    dependencies=[Depends(require_auth)],
)
def delete(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = DeleteBookUseCase(repo)
    uc.execute(book_id)
    return
