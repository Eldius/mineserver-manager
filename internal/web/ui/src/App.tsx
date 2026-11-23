import type { ReactNode } from 'react'
import { Link, Navigate, Route, Routes, useNavigate } from 'react-router-dom'
import { clearToken, getToken } from './api'
import Install from './pages/Install'
import Instances from './pages/Instances'
import Login from './pages/Login'

function RequireAuth({ children }: { children: ReactNode }) {
  if (!getToken()) {
    return <Navigate to="/login" replace />
  }
  return children
}

function Layout({ children }: { children: ReactNode }) {
  const navigate = useNavigate()
  const logout = () => {
    clearToken()
    navigate('/login')
  }
  return (
    <div className="layout">
      <header className="topbar">
        <span className="brand">⛏ mineserver manager</span>
        <nav>
          <Link to="/">Instances</Link>
          <Link to="/install">Install</Link>
          <button className="link" onClick={logout}>
            Logout
          </button>
        </nav>
      </header>
      <main>{children}</main>
    </div>
  )
}

export default function App() {
  return (
    <Routes>
      <Route path="/login" element={<Login />} />
      <Route
        path="/"
        element={
          <RequireAuth>
            <Layout>
              <Instances />
            </Layout>
          </RequireAuth>
        }
      />
      <Route
        path="/install"
        element={
          <RequireAuth>
            <Layout>
              <Install />
            </Layout>
          </RequireAuth>
        }
      />
      <Route path="*" element={<Navigate to="/" replace />} />
    </Routes>
  )
}
