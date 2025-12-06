// src/App.tsx

import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ShipsList from './pages/ShipsList'
import ShipPage from './pages/ShipPage'
import RequestShipPage from './pages/RequestShipPage'
import LoginPage from './pages/LoginPage'
import RegisterPage from './pages/RegisterPage'
import HomePage from './pages/HomePage'


export default function App(){
  return (
    <BrowserRouter>
    
      <div className="page-content" style={{ marginTop: '20px' }}> 
        <Routes>
            <Route path='/' element={<HomePage/>} />
            <Route path='/ships' element={<ShipsList/>} />
            <Route path='/ship/:id' element={<ShipPage/>} />
            <Route path='/request_ship/:id' element={<RequestShipPage/>} />
            <Route path='/login' element={<LoginPage/>} />
            <Route path='/register' element={<RegisterPage/>} />
        </Routes>
      </div>
    </BrowserRouter>
  )
}




