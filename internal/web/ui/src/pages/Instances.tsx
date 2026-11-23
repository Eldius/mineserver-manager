import { useCallback, useEffect, useState } from 'react'
import { Link } from 'react-router-dom'
import type { Instance } from '../api'
import { instanceAction, listInstances } from '../api'

export default function Instances() {
  const [instances, setInstances] = useState<Instance[]>([])
  const [error, setError] = useState<string | null>(null)
  const [loading, setLoading] = useState(true)
  const [busyId, setBusyId] = useState<string | null>(null)
  const [notice, setNotice] = useState<string | null>(null)

  const reload = useCallback(() => {
    setLoading(true)
    listInstances()
      .then((res) => {
        setInstances(res.instances ?? [])
        setError(null)
      })
      .catch((err: unknown) => setError(err instanceof Error ? err.message : 'failed to load'))
      .finally(() => setLoading(false))
  }, [])

  useEffect(() => {
    reload()
  }, [reload])

  const run = async (id: string, action: 'start' | 'stop' | 'backup') => {
    setBusyId(id)
    setNotice(null)
    setError(null)
    try {
      await instanceAction(id, action)
      if (action === 'backup') {
        setNotice('Backup created in the backups/ folder of the API home.')
      }
    } catch (err) {
      setError(err instanceof Error ? err.message : `${action} failed`)
    } finally {
      setBusyId(null)
      reload()
    }
  }

  return (
    <div className="page">
      <div className="page-header">
        <h1>Instances</h1>
        <button onClick={reload} disabled={loading}>
          {loading ? 'Loading…' : 'Refresh'}
        </button>
      </div>

      {error && <p className="error">{error}</p>}
      {notice && <p className="notice">{notice}</p>}

      {!loading && instances.length === 0 && (
        <p className="muted">
          No instances yet. <Link to="/install">Install your first server</Link>.
        </p>
      )}

      {instances.length > 0 && (
        <table>
          <thead>
            <tr>
              <th>Name</th>
              <th>Path</th>
              <th>Installed</th>
              <th>Status</th>
              <th>Actions</th>
            </tr>
          </thead>
          <tbody>
            {instances.map((i) => (
              <tr key={i.id}>
                <td>{i.name}</td>
                <td>
                  <code>{i.path}</code>
                </td>
                <td>{i.install_date ? new Date(i.install_date).toLocaleDateString() : '—'}</td>
                <td>
                  {!i.manageable ? (
                    <span className="badge warn" title="Lives outside the API home folder">
                      unmanaged
                    </span>
                  ) : i.running ? (
                    <span className="badge ok">running</span>
                  ) : (
                    <span className="badge">stopped</span>
                  )}
                </td>
                <td className="actions">
                  {i.manageable && (
                    <>
                      {i.running ? (
                        <button onClick={() => run(i.id, 'stop')} disabled={busyId === i.id}>
                          Stop
                        </button>
                      ) : (
                        <button onClick={() => run(i.id, 'start')} disabled={busyId === i.id}>
                          Start
                        </button>
                      )}
                      <button onClick={() => run(i.id, 'backup')} disabled={busyId === i.id}>
                        Backup
                      </button>
                    </>
                  )}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      )}
    </div>
  )
}
