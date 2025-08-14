from fastapi import APIRouter
from . import auth
from . import users
from . import book

api_router = APIRouter()
api_router.include_router(auth.router, prefix="/auth", tags=["Auth"])
api_router.include_router(users.router, prefix="/users", tags=["Users"])
api_router.include_router(book.router, prefix="/book", tags=["Book"])
