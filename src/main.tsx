// src/main.tsx

import ReactDOM from 'react-dom/client'

import 'bootstrap/dist/css/bootstrap.min.css'

import './resources/index_style.css'
import './resources/request_ship_style.css'
import './resources/ship_style.css'

import { Provider } from 'react-redux'
import { store } from './store'

import { RouterProvider, createBrowserRouter } from 'react-router-dom'
import routes from './App'
import React from 'react'

const router = createBrowserRouter(routes, {
  basename: '/loading-time-frontend/'
})

ReactDOM.createRoot(document.getElementById('root')!).render(
  <React.StrictMode>
    <Provider store={store}>
      <RouterProvider router={router} />
    </Provider>
  </React.StrictMode>
)

if ("serviceWorker" in navigator) {
  window.addEventListener("load", () => {
    navigator.serviceWorker
      .register("/loading-time-frontend/service-worker.js")
      .then(reg => console.log("Service worker registered:", reg))
      .catch(err => console.error("SW registration failed:", err));
  });
}