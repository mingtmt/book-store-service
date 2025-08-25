import uuid

from pydantic import BaseModel


class RoleOut(BaseModel):
    id: uuid.UUID
    name: str


class CreateRoleIn(BaseModel):
    name: str
