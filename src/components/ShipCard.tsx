// src/components/ShipCard.tsx

import { Link, useNavigate } from 'react-router-dom'
import { addShipToRequest } from '../api'
import { getToken } from '../auth'
import '../resources/ShipCard.css'

export default function ShipCard({ ship }: { ship: any }) {
  const buildImgSrc = (p?: string | null) => {
    if (!p) return '/default.png'
    try { new URL(p); return p }
    catch (e) { return 'http://localhost:9000/loading-time-img/img/' + p }
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
          onError={(e:any)=>{ e.target.style.display='none' }}
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
