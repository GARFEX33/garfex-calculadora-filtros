// Re-exporta las utilidades principales accesibles como $lib/...
// Importar siempre desde los sub-módulos específicos para tree-shaking óptimo:
//   import { cn } from '$lib/utils'
//   import type { ApiResult } from '$lib/types'
//   import { apiClient } from '$lib/api/client'

export { cn } from './utils/index.js';
export type { ApiResult, ApiError, ApiResponse, ApiPaginatedResponse } from './types/index.js';
