from fastapi import APIRouter, Depends, status
from app.schemas.base import Envelope
from app.schemas.book import CreateBookRequest, CreateBookResponse
from app.infrastructure.db.sqlalchemy.repos.book_impl import SqlAlchemyBookRepository
from app.use_cases.book.create_book import CreateBookUseCase
from sqlalchemy.orm import Session
from app.infrastructure.web.dependencies.db import get_db

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
