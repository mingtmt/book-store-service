import uuid
from fastapi import APIRouter, Depends, status, HTTPException
from sqlalchemy.orm import Session
from app.presentation.http.schemas.base import Envelope
from app.presentation.http.schemas.books import CreateBook, UpdateBook, BookOut
from app.usecases.books.create_book import CreateBookUseCase
from app.usecases.books.get_book import GetAllBooksUseCase, GetBookByIdUseCase
from app.usecases.books.update_book import UpdateBookUseCase, UpdateBookCommand
from app.usecases.books.delete_book import DeleteBookUseCase
from app.infrastructure.db.sqlalchemy.repositories.book_impl import SqlAlchemyBookRepository
from app.domain.errors import BookNotFound, ConstraintViolation
from app.presentation.http.dependencies.db import get_db

router = APIRouter()

@router.post("/", response_model=Envelope[BookOut], status_code=status.HTTP_201_CREATED)
def create(payload: CreateBook, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = CreateBookUseCase(repo)
    created_book = uc.execute(payload.title, payload.author, payload.price, payload.description, payload.category)
    return Envelope(
        data=BookOut(
            id=created_book.id,
            title=created_book.title,
            author=created_book.author,
            price=created_book.price,
            description=created_book.description,
            category=created_book.category
        )
    )

@router.get("/{book_id}", response_model=Envelope[BookOut], status_code=status.HTTP_200_OK)
def get_by_id(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetBookByIdUseCase(repo)
    book = uc.execute(book_id)
    if book is None:
        return Envelope(data=None, message="Book not found", status_code=status.HTTP_404_NOT_FOUND)
    return Envelope(
        data=BookOut(
            id=book.id,
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category
        )
    )

@router.get("/", response_model=Envelope[list[BookOut]], status_code=status.HTTP_200_OK)
def get_all(db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetAllBooksUseCase(repo)
    books = uc.execute()
    if books is None:
        return Envelope(data=None, message="No books found", status_code=status.HTTP_404_NOT_FOUND)
    return Envelope(
        data=[
            BookOut(
                id=book.id,
                title=book.title,
                author=book.author,
                price=book.price,
                description=book.description,
                category=book.category
            )
            for book in books
        ]
    )

@router.put("/{book_id}", response_model=Envelope[BookOut], status_code=status.HTTP_200_OK)
def update(book_id: uuid.UUID, payload: UpdateBook, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = UpdateBookUseCase(repo)
    cmd = UpdateBookCommand(**payload.model_dump(exclude_unset=True))

    try:
        updated_book = uc.execute(book_id, cmd)
    except BookNotFound:
        raise HTTPException(status_code=404, detail="Book not found")
    except ConstraintViolation as e:
        raise HTTPException(status_code=400, detail=str(e))
    return Envelope(
        data=BookOut(
            id=updated_book.id,
            title=updated_book.title,
            author=updated_book.author,
            price=updated_book.price,
            description=updated_book.description,
            category=updated_book.category
        )
    )

@router.delete("/{book_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = DeleteBookUseCase(repo)
    uc.execute(book_id)
    return
