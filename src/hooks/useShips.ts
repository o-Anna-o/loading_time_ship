// src/hooks/useShips.ts
import { useEffect, useState } from 'react'
import { getShips, ShipsFilterParams } from '../api'
import mock from '../mock'
import { Ship } from '../types'

type ShipForFilter = { Name?: string; name?: string; [key: string]: any }

export function useShips(appliedSearch?: string) {
  const [ships, setShips] = useState<Ship[]>([])
  const [loading, setLoading] = useState(false)
  const [error, setError] = useState<string | null>(null)

  useEffect(() => {
    let cancelled = false

    const fetchShips = async () => {
      setLoading(true)
      setError(null)

      const params: ShipsFilterParams | undefined = appliedSearch
        ? { search: appliedSearch }
        : undefined

      try {
        const res = await getShips(params)

        // бекенд возвращает в поле data или массив
        let arr = Array.isArray(res) ? res : res?.data ?? []

        if (params?.search) {
          const searchLower = params.search.toLowerCase()
          arr = arr.filter((ship: ShipForFilter) => {
            const name =
              ship.Name ??
              ship.name ??
              ''
            return name.toLowerCase().includes(searchLower)
          })
        }

        if (!cancelled) setShips(arr)
      } catch (err) {
        console.warn('Backend unavailable, switching to mock...', err)

        // fallback на mock
        let arr = mock

        if (appliedSearch) {
          const s = appliedSearch.toLowerCase()
          arr = arr.filter(ship =>
            ship.name.toLowerCase().includes(s)
          )
        }

        if (!cancelled) setShips(arr)
          
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    fetchShips()
    return () => { cancelled = true }
  }, [appliedSearch])

  return { ships, loading, error }
}
