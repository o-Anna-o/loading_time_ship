

import { Link, useLocation } from 'react-router-dom'
import '../resources/Breadcrumbs.css'  

export default function Breadcrumbs() {
  const location = useLocation()
  const parts = location.pathname.split('/').filter(Boolean)

  const shipNames: { [key: string]: string } = {
    '1': 'Ever Ace',
    '2': 'FESCO Diomid', 
    '3': 'HMM Algeciras',
    '4': 'MSC Gulsun'
  }

  const prettyPart = (part: string, index: number) => {
    if (index > 0 && parts[index - 1] === 'ship' && shipNames[part]) {
      return shipNames[part]
    }
    
    switch (part) {
      case 'ships': return 'Корабли'
      case 'ship': return 'Корабли'
      case 'request_ship': return 'Заявка'
      case 'login': return 'Вход'
      case 'register': return 'Регистрация'
      default: return part
    }
  }

  const crumbs = parts.map((part, idx) => {
    const path = part === 'ship' ? '/ships' : '/' + parts.slice(0, idx + 1).join('/')
    const displayName = prettyPart(part, idx)

    return (
      <span key={path}>
        <Link to={path}>{displayName}</Link>
        {idx < parts.length - 1 ? ' / ' : ''}
      </span>
    )
  })

  return (
    <div className="breadcrumbs">
      <Link to="/">Главная</Link>
      {parts.length > 0 ? ' / ' : ''}
      {crumbs}
    </div>
  )
}