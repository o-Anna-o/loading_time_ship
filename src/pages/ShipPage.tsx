import React, { useEffect, useState } from 'react'
import { useParams } from 'react-router-dom'
import { getShips } from '../api'
import ShipListIcon from '../components/ShipListIcon'
import Breadcrumbs from '../components/Breadcrumbs'
import mock from '../mock'
import '../resources/ShipDetail.css'

export default function ShipPage() {
  const { id } = useParams()
  const [ship, setShip] = useState<any>(null)
  const [loading, setLoading] = useState(true)

 
  const buildImgSrc = (p?: string | null) => {
    if (!p)
      return `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`

    if (/^https?:\/\//i.test(p)) return p

    const baseImg = (import.meta.env?.VITE_IMG_BASE as string) ?? ''
    if (baseImg) return `${baseImg}/${p}`

    return `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`
  }

  useEffect(() => {
    if (!id) return

    let cancelled = false

    const load = async () => {
      setLoading(true)
      try {
        const list = await getShips()
        if (cancelled) return

        const arr: any[] = Array.isArray(list) ? list : list?.data ?? list?.ships ?? []
        const s = arr.find((x: any) =>
          String(x.ship_id ?? x.ShipID ?? x.id ?? x.ID) === String(id)
        )

        if (!cancelled) setShip(s ?? null)
      } catch (err) {
        console.warn('Backend unavailable → using mock', err)
        const s = mock.find((m: any) =>
          String(m.ship_id ?? m.ShipID ?? m.id ?? m.ID) === String(id)
        )
        if (!cancelled) setShip(s ?? null)
      } finally {
        if (!cancelled) setLoading(false)
      }
    }

    load()
    return () => { cancelled = true }
  }, [id])

  if (loading) return <div className="loading">Загрузка...</div>
  if (!ship) return <div className="loading">Корабль не найден</div>

  const photo = ship.photo_url ?? ship.PhotoURL ?? ship.photo ?? ship.Photo
  const src = buildImgSrc(photo)

  const name = ship.name ?? ship.Name ?? 'Без названия'
  const capacity = ship.capacity ?? ship.Capacity ?? ''
  const length = ship.length ?? ship.Length ?? ''
  const width = ship.width ?? ship.Width ?? ''
  const containers = ship.containers ?? ship.Containers ?? ''
  const cranes = ship.cranes ?? ship.Cranes ?? ''
  const description = ship.description ?? ship.Description ?? ''

  return (
    <>
      <ShipListIcon />
      <Breadcrumbs />

      <div className="ship-detail-wrapper">
        <div className="ship-detail-card">
          <h1 className="ship-detail-title">{name}</h1>

          <div className="ship-detail-image">
            <img
              src={src}
              alt={name}
              onError={(e: any) => {
                const fallback =
                  `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`
                if (e.target.src !== fallback) {
                  e.target.src = fallback
                }
              }}
            />
          </div>

          <div className="ship-detail-text">
            <p><b>Вместимость:</b> {capacity} TEU</p>
            <p><b>Габариты:</b> длина {length} м, ширина {width} м</p>
            <p><b>Максимальное число контейнеров на 40 футов:</b> {containers} шт.</p>
            <p><b>Максимальное число кранов:</b> {cranes} одновременно</p>
            <p><b>Особенности:</b> {description}</p>
          </div>
        </div>
      </div>
    </>
  )
}
