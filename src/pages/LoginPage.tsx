import React, { useState } from 'react'
import { saveToken } from '../auth'
import { useNavigate } from 'react-router-dom'

export default function LoginPage(){
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string|null>(null)
  const navigate = useNavigate()

  const onSubmit = async (e:any) => {
    e.preventDefault()
    setError(null)
    try {
      let res = await fetch('/api/users/login', {
        method:'POST',
        headers: {'Content-Type':'application/json'},
        body: JSON.stringify({ login, password }),
      })

      if (res.status === 400 || res.status === 415) {
        res = await fetch('/api/users/login', {
          method:'POST',
          headers: {'Content-Type':'application/x-www-form-urlencoded'},
          body: new URLSearchParams({ login, password }),
        })
      }

      if (!res.ok) {
        let text = ''
        try { text = await res.text() } catch(e) {}
        setError(text || ('Login failed: ' + res.status))
        return
      }

      const json = await res.json().catch(()=> ({}))
      const token = json.token ?? json.access_token ?? json.Token ?? json.jwt ?? json.data?.token ?? null

      if (!token) {
        const err = json.error ?? JSON.stringify(json) ?? 'No token received'
        setError(String(err))
        return
      }

      saveToken(token)
      navigate('/ships')
    } catch(err:any) {
      setError(String(err))
    }
  }

  return (
    <div style={{display:'flex', flexDirection:'column', alignItems:'center'}}>
      <div style={{width:600, marginTop:40, background:'#3A3A3A', padding:30, borderRadius:6}}>
        <h2>Вход</h2>
        {error ? <div style={{color:'red'}}>{error}</div> : null}
        <form onSubmit={onSubmit}>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" placeholder="Логин" value={login} onChange={e=>setLogin(e.target.value)} />
          </div>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" type="password" placeholder="Пароль" value={password} onChange={e=>setPassword(e.target.value)} />
          </div>
          <button className="btn" type="submit">Войти</button>
        </form>
      </div>
    </div>
  )
}
