// src/App.tsx

import ShipsList from './pages/ShipsList'
import ShipPage from './pages/ShipPage'
import RequestShipPage from './pages/RequestShipPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import HomePage from './pages/HomePage'

export const routes = [
  { path: '/', element: <HomePage /> },
  { path: '/ships', element: <ShipsList /> },
  { path: '/ship/:id', element: <ShipPage /> },
  { path: '/request_ship/:id', element: <RequestShipPage /> },
  { path: '/login', element: <LoginPage /> },
  { path: '/register', element: <RegisterPage /> },
]

export default routes
