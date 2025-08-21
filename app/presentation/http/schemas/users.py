import uuid

from pydantic import BaseModel


class UpdateRequest(BaseModel):
    id: uuid.UUID
    email: str | None
    password: str | None = None


class UpdateResponse(BaseModel):
    id: uuid.UUID
    email: str
