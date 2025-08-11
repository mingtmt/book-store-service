class AppException(Exception):
    """Base class cho custom app exceptions"""
    status_code = 500
    code = "INTERNAL_ERROR"
    message = "Internal server error"

    def __init__(self, message: str | None = None):
        if message:
            self.message = message
        super().__init__(self.message)

class NotFoundException(AppException):
    status_code = 404
    code = "NOT_FOUND"
    message = "Resource not found"

class UnauthorizedException(AppException):
    status_code = 401
    code = "UNAUTHORIZED"
    message = "Unauthorized"

class BadRequestException(AppException):
    status_code = 400
    code = "BAD_REQUEST"
    message = "Bad request"

class ConflictException(AppException):
    status_code = 409
    code = "CONFLICT"
    message = "Conflict"

class EmailAlreadyExistsException(ConflictException):
    code = "EMAIL_ALREADY_EXISTS"
    message = "Email already exists"