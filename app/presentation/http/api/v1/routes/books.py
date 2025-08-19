import uuid
from fastapi import APIRouter, Depends, status
from sqlalchemy.orm import Session
from app.presentation.http.schemas.base import Envelope
from app.presentation.http.schemas.books import CreateBookRequest, CreateBookResponse, GetBookResponse
from app.usecases.books.create_book import CreateBookUseCase
from app.usecases.books.get_book import GetAllBooksUseCase, GetBookByIdUseCase
from app.usecases.books.delete_book import DeleteBookUseCase
from app.infrastructure.db.sqlalchemy.repositories.book_impl import SqlAlchemyBookRepository
from app.presentation.http.dependencies.db import get_db

router = APIRouter()

@router.post("/", response_model=Envelope[CreateBookResponse], status_code=status.HTTP_201_CREATED)
def create(data: CreateBookRequest, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = CreateBookUseCase(repo)
    book = uc.execute(data.title, data.author, data.price, data.description, data.category)
    return Envelope(
        data=CreateBookResponse(
            id=book.id,
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category
        )
    )

@router.get("/{book_id}", response_model=Envelope[GetBookResponse], status_code=status.HTTP_200_OK)
def get_by_id(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetBookByIdUseCase(repo)
    book = uc.execute(book_id)
    if book is None:
        return Envelope(data=None, message="Book not found", status_code=status.HTTP_404_NOT_FOUND)
    return Envelope(
        data=GetBookResponse(
            id=book.id,
            title=book.title,
            author=book.author,
            price=book.price,
            description=book.description,
            category=book.category
        )
    )

@router.get("/", response_model=Envelope[list[GetBookResponse]], status_code=status.HTTP_200_OK)
def get_all(db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = GetAllBooksUseCase(repo)
    books = uc.execute()
    if books is None:
        return Envelope(data=None, message="No books found", status_code=status.HTTP_404_NOT_FOUND)
    return Envelope(
        data=[
            GetBookResponse(
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

@router.delete("/{book_id}", status_code=status.HTTP_204_NO_CONTENT)
def delete(book_id: uuid.UUID, db: Session = Depends(get_db)):
    repo = SqlAlchemyBookRepository(db)
    uc = DeleteBookUseCase(repo)
    uc.execute(book_id)
    return
