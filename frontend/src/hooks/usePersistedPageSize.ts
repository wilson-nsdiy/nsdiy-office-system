const STORAGE_KEY = 'table_page_size'

export function getPersistedPageSize(): number {
  if (typeof window === 'undefined') return 10
  try {
    const saved = localStorage.getItem(STORAGE_KEY)
    if (saved) {
      const size = parseInt(saved, 10)
      if (!isNaN(size) && size > 0 && size <= 100) {
        return size
      }
    }
  } catch {
    // Ignore localStorage errors
  }
  return 10
}

export function setPersistedPageSize(size: number): void {
  if (typeof window === 'undefined') return
  try {
    localStorage.setItem(STORAGE_KEY, String(size))
  } catch {
    // Ignore localStorage errors
  }
}
