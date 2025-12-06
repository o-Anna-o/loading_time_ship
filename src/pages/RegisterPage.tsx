import React, { useState } from 'react'
import { useNavigate } from 'react-router-dom'

export default function RegisterPage(){
  const [fio, setFio] = useState('')
  const [login, setLogin] = useState('')
  const [password, setPassword] = useState('')
  const [contacts, setContacts] = useState('')
  const [error, setError] = useState<string|null>(null)
  const navigate = useNavigate()

  const onSubmit = async (e:any) => {
    e.preventDefault()
    setError(null)
    try {
      const res = await fetch('/api/users/register', {
        method:'POST',
        headers: {'Content-Type':'application/x-www-form-urlencoded'},
        body: new URLSearchParams({fio, login, password, contacts})
      })
      if(!res.ok) {
        const text = await res.text()
        setError(text || 'Register failed')
        return
      }
      // success - go to login
      navigate('/login')
    } catch(err:any) {
      setError(String(err))
    }
  }

  return (
    <div style={{display:'flex', flexDirection:'column', alignItems:'center'}}>
      <div style={{width:600, marginTop:40, background:'#3A3A3A', padding:30, borderRadius:6}}>
        <h2>Регистрация</h2>
        {error ? <div style={{color:'red'}}>{error}</div> : null}
        <form onSubmit={onSubmit}>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" placeholder="ФИО" value={fio} onChange={e=>setFio(e.target.value)} />
          </div>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" placeholder="Логин" value={login} onChange={e=>setLogin(e.target.value)} />
          </div>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" type="password" placeholder="Пароль" value={password} onChange={e=>setPassword(e.target.value)} />
          </div>
          <div style={{marginBottom:10}}>
            <input className="request__cnt-input" placeholder="Контакты" value={contacts} onChange={e=>setContacts(e.target.value)} />
          </div>
          <button className="btn" type="submit">Зарегистрироваться</button>
        </form>
      </div>
    </div>
  )
}
