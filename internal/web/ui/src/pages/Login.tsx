import { useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { login } from '../api'

export default function Login() {
  const [password, setPassword] = useState('')
  const [error, setError] = useState<string | null>(null)
  const [busy, setBusy] = useState(false)
  const navigate = useNavigate()

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setBusy(true)
    try {
      await login(password)
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'login failed')
    } finally {
      setBusy(false)
    }
  }

  return (
    <div className="centered">
      <form className="card" onSubmit={onSubmit}>
        <h1>⛏ mineserver manager</h1>
        <p className="muted">Sign in with the API password (api.auth.password)</p>
        <label>
          Password
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
            autoFocus
            required
          />
        </label>
        {error && <p className="error">{error}</p>}
        <button type="submit" disabled={busy || !password}>
          {busy ? 'Signing in…' : 'Sign in'}
        </button>
      </form>
    </div>
  )
}
