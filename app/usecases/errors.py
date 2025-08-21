from __future__ import annotations

from typing import Any, Mapping


class UseCaseError(Exception):
    code: str = "USECASE_ERROR"
    default_message: str = "Use case error"

    def __init__(
        self,
        message: str | None = None,
        *,
        context: Mapping[str, Any] | None = None,
        cause: Exception | None = None,
    ) -> None:
        self.message = message or self.default_message
        self.context = dict(context or {})
        if cause is not None:
            self.__cause__ = cause
        super().__init__(self.message)

    def to_dict(self) -> dict[str, Any]:
        data = {"code": self.code, "message": self.message}
        if self.context:
            data["context"] = self.context
        return data

    def __str__(self) -> str:
        return self.message


class NotFound(UseCaseError):
    code = "NOT_FOUND"
    default_message = "Resource not found"


class Unauthorized(UseCaseError):
    code = "UNAUTHORIZED"
    default_message = "Authentication required"


class Forbidden(UseCaseError):
    code = "FORBIDDEN"
    default_message = "Not allowed to perform this action"


class BadRequest(UseCaseError):
    code = "BAD_REQUEST"
    default_message = "Invalid request"


class Conflict(UseCaseError):
    code = "CONFLICT"
    default_message = "Request conflicts with current state"
