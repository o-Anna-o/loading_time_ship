// src/pages/HomePage.tsx
// import React from 'react'
// import { Link } from 'react-router-dom'
import ShipListIcon from '../components/ShipListIcon'
import '../resources/HomePage.css' 

export default function HomePage() {
  return (
    <div className="home-wrapper">
      <ShipListIcon />

      <div className="home-card">
        <h1 className="home-title">Loading Time Ship</h1>
        <p className="home-text">
          Добро пожаловать в Loading Time Ship! Здесь вы можете рассчитать время погрузки контейнеровоза в порту.
        </p>
      </div>
    </div>
  )
}