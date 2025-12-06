// src/api.ts
import { getToken } from './auth'
const API_BASE = '/api'

// НЕ используем mock здесь — только бэкенд
export type ShipsFilterParams = {
  search?: string
  capacity_min?: number
  capacity_max?: number
  cranes_min?: number
  cranes_max?: number
}

export async function getShips(params?: ShipsFilterParams) {
  // логирование для диагностики
  console.log('[api.getShips] called with', params)

  const url = new URL(API_BASE + '/ships', window.location.origin)

  if (params?.search) url.searchParams.set('search', String(params.search))
  if (typeof params?.capacity_min === 'number') url.searchParams.set('capacity_min', String(params.capacity_min))
  if (typeof params?.capacity_max === 'number') url.searchParams.set('capacity_max', String(params.capacity_max))
  if (typeof params?.cranes_min === 'number') url.searchParams.set('cranes_min', String(params.cranes_min))
  if (typeof params?.cranes_max === 'number') url.searchParams.set('cranes_max', String(params.cranes_max))

  // ещё лог — покажет итоговый URL в консоли браузера
  console.log('[api.getShips] fetch ->', url.toString())

  const res = await fetch(url.toString(), {
    method: 'GET',
    credentials: 'include'
  })

  if (!res.ok) {
    const text = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (text ? ': ' + text : ''))
  }

  const j = await res.json().catch(()=>null)
  if (j && j.data && Array.isArray(j.data)) return j.data
  if (j && Array.isArray(j)) return j
  return j ?? []
}

export async function addShipToRequest(shipId:number){
  console.log('[DEBUG api] addShipToRequest token=', getToken());
  const token = getToken();
  const headers: Record<string,string> = {'Content-Type':'application/json'};
  if (token) headers['Authorization'] = 'Bearer ' + token;

  const response = await fetch(`${API_BASE}/ships/${shipId}/add-to-ship-bucket`, {
    method: 'POST',
    headers,
    credentials: 'include'
  });

  if (!response.ok) {
    let text = ''
    try { text = await response.text() } catch(e) {}
    throw new Error('HTTP ' + response.status + (text ? ': ' + text : ''))
  }
  try { return await response.json() } catch { return {} }
}

// получить одну заявку по id
export async function getRequestShip(id: number | string) {
  const token = getToken()
  const headers: Record<string,string> = {'Content-Type': 'application/json'}
  if (token) headers['Authorization'] = 'Bearer ' + token

  const url = `${API_BASE}/request_ship/${id}`
  const res = await fetch(url, {
    method: 'GET',
    headers,
    credentials: 'include'
  })

  if (!res.ok) {
    const text = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (text ? ': ' + text : ''))
  }
  const j = await res.json().catch(()=>null)
  if (!j) return {}
  // normalize common shapes
  const data = j.data ?? j
  return data
}

// получить корзину/черновик
export async function getRequestShipBasket() {
  const token = getToken()
  const headers: Record<string,string> = {'Content-Type': 'application/json'}
  if (token) headers['Authorization'] = 'Bearer ' + token

  const res = await fetch(`${API_BASE}/request_ship/basket`, {
    method: 'GET',
    headers,
    credentials: 'include'
  })
  if (!res.ok) {
    const text = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (text ? ': ' + text : ''))
  }
  const j = await res.json().catch(()=>null)
  if (!j) return null
  const payload = j.data ?? j

  // return normalized small object the redirect/consumer expects
  return {
    request_ship_id: payload.request_ship_id ?? payload.RequestShipID ?? payload.id ?? payload.requestShipId ?? null,
    ships_count: payload.ships_count ?? payload.ShipsCount ?? payload.count ?? null,
    raw: payload
  }
}

// удалить конкретный корабль из заявки
export async function deleteShipFromRequest(requestId: number | string, shipId: number | string) {
  const token = getToken()
  const headers: Record<string, string> = {}
  if (token) headers['Authorization'] = 'Bearer ' + token

  const res = await fetch(`${API_BASE}/request_ship/${requestId}/ships/${shipId}`, {
    method: 'DELETE',
    headers,
    credentials: 'include',
  })
  if (!res.ok) {
    const txt = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (txt ? ': ' + txt : ''))
  }
  return res.json().catch(()=>null)
}

// удалить всю заявку
export async function deleteRequestShip(requestId: number | string) {
  const token = getToken()
  const headers: Record<string,string> = {'Content-Type': 'application/json'}
  if (token) headers['Authorization'] = 'Bearer ' + token

  const res = await fetch(`${API_BASE}/request_ship/${requestId}`, {
    method: 'POST',
    headers,
    credentials: 'include',
    body: JSON.stringify({ _method: 'DELETE' })
  })
  if (!res.ok) {
    const txt = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (txt ? ': ' + txt : ''))
  }
  return res.json().catch(()=>null)
}

// рассчитать время погрузки (POST)
export async function calculateLoadingTime(requestId: number | string, payload: { containers_20ft?: number, containers_40ft?: number, comment?: string }) {
  const token = getToken()
  const headers: Record<string,string> = {'Content-Type': 'application/json'}
  if (token) headers['Authorization'] = 'Bearer ' + token

  const res = await fetch(`${API_BASE}/request_ship/calculate_loading_time/${requestId}`, {
    method: 'POST',
    headers,
    credentials: 'include',
    body: JSON.stringify(payload)
  })
  if (!res.ok) {
    const txt = await res.text().catch(()=> '')
    throw new Error('HTTP ' + res.status + (txt ? ': ' + txt : ''))
  }
  return res.json().catch(()=>null)
}