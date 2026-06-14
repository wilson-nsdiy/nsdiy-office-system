import { useState, useCallback, useRef, useEffect } from 'react'

interface PaginationState {
  page: number
  pageSize: number
  total: number
  totalPages: number
}

interface UseTableLoaderOptions<T, P extends Record<string, any>> {
  fetchFn: (params: P & { page: number; pageSize: number }) => Promise<{
    items: T[]
    total: number
    totalPages: number
  }>
  initialParams?: P
  pageSize?: number
  debounceMs?: number
}

interface UseTableLoaderReturn<T, P extends Record<string, any>> {
  items: T[]
  loading: boolean
  params: P
  pagination: PaginationState
  load: () => Promise<void>
  reload: () => Promise<void>
  setParams: (newParams: Partial<P>) => void
  handlePageChange: (page: number) => void
  handlePageSizeChange: (size: number) => void
}

export function useTableLoader<T, P extends Record<string, any>>(
  options: UseTableLoaderOptions<T, P>
): UseTableLoaderReturn<T, P> {
  const { fetchFn, initialParams, pageSize: initialPageSize = 10 } = options

  const [items, setItems] = useState<T[]>([])
  const [loading, setLoading] = useState(false)
  const [params, setParamsState] = useState<P>(initialParams as P)
  const [pagination, setPagination] = useState<PaginationState>({
    page: 1,
    pageSize: initialPageSize,
    total: 0,
    totalPages: 0,
  })

  const abortControllerRef = useRef<AbortController | null>(null)
  const debounceTimerRef = useRef<ReturnType<typeof setTimeout> | null>(null)

  const load = useCallback(async () => {
    if (abortControllerRef.current) {
      abortControllerRef.current.abort()
    }

    const controller = new AbortController()
    abortControllerRef.current = controller

    setLoading(true)
    try {
      const result = await fetchFn({
        ...params,
        page: pagination.page,
        pageSize: pagination.pageSize,
      })

      if (!controller.signal.aborted) {
        setItems(result.items || [])
        setPagination((prev) => ({
          ...prev,
          total: result.total || 0,
          totalPages: result.totalPages || 0,
        }))
      }
    } catch (error: any) {
      if (error?.name !== 'AbortError' && error?.code !== 'ERR_CANCELED') {
        console.error('Table load error:', error)
        throw error
      }
    } finally {
      if (abortControllerRef.current === controller) {
        setLoading(false)
      }
    }
  }, [fetchFn, params, pagination.page, pagination.pageSize])

  const reload = useCallback(() => {
    setPagination((prev) => ({ ...prev, page: 1 }))
    return load()
  }, [load])

  const setParams = useCallback(
    (newParams: Partial<P>) => {
      setParamsState((prev) => ({ ...prev, ...newParams }))
      setPagination((prev) => ({ ...prev, page: 1 }))

      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
      }
      debounceTimerRef.current = setTimeout(() => {
        load()
      }, 300)
    },
    [load]
  )

  const handlePageChange = useCallback(
    (page: number) => {
      const validPage = Math.max(1, Math.min(page, pagination.totalPages || 1))
      setPagination((prev) => ({ ...prev, page: validPage }))
    },
    [pagination.totalPages]
  )

  const handlePageSizeChange = useCallback((size: number) => {
    setPagination((prev) => ({ ...prev, pageSize: size, page: 1 }))
  }, [])

  useEffect(() => {
    return () => {
      abortControllerRef.current?.abort()
      if (debounceTimerRef.current) {
        clearTimeout(debounceTimerRef.current)
      }
    }
  }, [])

  return {
    items,
    loading,
    params,
    pagination,
    load,
    reload,
    setParams,
    handlePageChange,
    handlePageSizeChange,
  }
}
