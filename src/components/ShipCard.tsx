// src/components/ShipCard.tsx
import { Link, useNavigate } from 'react-router-dom'
import { addShipToRequest } from '../api'
import { getToken } from '../auth'
import '../resources/ShipCard.css'

export default function ShipCard({ ship }: { ship: any }) {
  const buildImgSrc = (p?: string | null) => {
    // если нет фото — берем default из public, с учётом BASE_URL
    if (!p) return `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`

    // абсолютный URL 
    if (/^https?:\/\//i.test(p)) return p

    const baseImg = (import.meta.env?.VITE_IMG_BASE as string) ?? ''
    if (baseImg) return `${baseImg}/${p}`

    // fallback — default
    return `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`
  }

  const src = buildImgSrc(ship.photo_url ?? ship.PhotoURL)
  const id = ship.ship_id ?? ship.ShipID ?? ship.id ?? ship.ID
  const name = ship.name ?? ship.Name
  const capacity = ship.capacity ?? ship.Capacity
  const cranes = ship.cranes ?? ship.Cranes

  const navigate = useNavigate()

  const handleAdd = async (e: any) => {
    e.preventDefault()
    const token = getToken()
    if (!token) {
      navigate('/login')
      return
    }
    try {
      await addShipToRequest(Number(id))
      window.dispatchEvent(new Event('lt:basket:refresh'))
    } catch (err: any) {
      alert('Ошибка: ' + (err.message || 'unknown'))
    }
  }

  const isAuthenticated = Boolean(getToken())

  return (
    <div className="ship-item">
      <div className="ship-image">
        <img
          src={src}
          alt={name}
          onError={(e: any) => {
            const fallback = `${import.meta.env.BASE_URL ?? '/loading-time-frontend/'}default.png`
            if (e.target.src !== fallback) {
              e.target.src = fallback     // ← подставляем default.png
            }
          }}
        />
      </div>

      <h2 className="ship-title">
        <Link to={'/ship/' + id}>{name}</Link>
      </h2>

      <div className="ship-item-text">
        <p><b>Вместимость:</b> {capacity} TEU</p>
        <p><b>Краны:</b> {cranes} одновременно</p>
      </div>

      {isAuthenticated && (
        <form onSubmit={handleAdd}>
          <button type="submit" className="btn card-btn ship-btn">
            Добавить
          </button>
        </form>
      )}
    </div>
  )
}
