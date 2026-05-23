export interface ApiError {
  code: string
  message: string
  details?: Record<string, unknown>
}

export interface ApiResponse<T> {
  data: T
  error: ApiError | null
}

export class ApiClientError extends Error {
  code: string
  details?: Record<string, unknown>
  status: number

  constructor(message: string, code: string, status: number, details?: Record<string, unknown>) {
    super(message)
    this.name = "ApiClientError"
    this.code = code
    this.status = status
    this.details = details
  }
}
