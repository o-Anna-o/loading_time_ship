// src/pages/ShipsList.tsx
import React, { useEffect } from 'react'
import Navbar from '../components/Navbar'
import Breadcrumbs from '../components/Breadcrumbs'
import ShipCard from '../components/ShipCard'
import { useShips } from '../hooks/useShips'
import { useAppDispatch, useAppSelector } from '../store/hooks'
import { setSearch, applySearch } from '../store/slices/filterSlice'

export default function ShipsList() {
  const dispatch = useAppDispatch()
  const search = useAppSelector(state => state.filter.search)
  const appliedSearch = useAppSelector(state => state.filter.appliedSearch)

  // при загрузке/возврате синхронизируем input со значением appliedSearch
  useEffect(() => {
    dispatch(setSearch(appliedSearch))
  }, [appliedSearch, dispatch])

  const { ships, loading, error } = useShips(appliedSearch)

  const handleSubmit = (e: React.FormEvent) => {
    e.preventDefault()
    dispatch(applySearch())
  }

  return (
    <>
      <Navbar />
      <Breadcrumbs />

      <div style={{ display: 'flex', flexDirection: 'column', alignItems: 'center' }}>
        <form
          className="page__search"
          onSubmit={handleSubmit}
          style={{ marginTop: 10 }}
        >
          <input
            className="search-input page__search-input page__search-item"
            type="text"
            placeholder="Поиск контейнеровоза"
            value={search}
            onChange={e => dispatch(setSearch(e.target.value))}
          />
          <button className="btn search-btn page__search-item" type="submit">
            Найти
          </button>
        </form>

        {loading && <div style={{ marginTop: 20 }}>Загрузка...</div>}
        {error && <div style={{ marginTop: 20, color: 'red' }}>{error}</div>}

        <ul className="ship-cards">
          {ships.map((s: any) => {
            const id =
              s.ship_id ??
              s.ShipID ??
              s.id ??
              s.ID

            return (
              <li key={id}>
                <ShipCard ship={s} />
              </li>
            )
          })}
        </ul>
      </div>
    </>
  )
}