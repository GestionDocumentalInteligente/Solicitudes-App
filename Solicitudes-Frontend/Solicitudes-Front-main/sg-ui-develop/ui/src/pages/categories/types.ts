export interface Category {
  id: number;
  name: string;
  description: string;
  is_active: boolean;
}

export interface CategoryResponse {
  data: Category[];
  error: string;
}

class RequestError extends Error {}
class UnhandledError extends Error {}

export { RequestError, UnhandledError };
