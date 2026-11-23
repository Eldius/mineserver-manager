import { useEffect, useState } from 'react'
import type { FormEvent } from 'react'
import { useNavigate } from 'react-router-dom'
import { installInstance, listVersions } from '../api'

export default function Install() {
  const [name, setName] = useState('')
  const [version, setVersion] = useState('latest')
  const [versions, setVersions] = useState<string[]>([])
  const [versionsError, setVersionsError] = useState<string | null>(null)
  const [memory, setMemory] = useState('1g')
  const [motd, setMotd] = useState('')
  const [levelName, setLevelName] = useState('')
  const [seed, setSeed] = useState('')
  const [rconEnabled, setRconEnabled] = useState(false)
  const [rconPassword, setRconPassword] = useState('')
  const [busy, setBusy] = useState(false)
  const [error, setError] = useState<string | null>(null)
  const navigate = useNavigate()

  useEffect(() => {
    listVersions()
      .then((res) => setVersions(res.versions ?? []))
      .catch((err: unknown) =>
        setVersionsError(err instanceof Error ? err.message : 'failed to load versions'),
      )
  }, [])

  const onSubmit = async (e: FormEvent) => {
    e.preventDefault()
    setError(null)
    setBusy(true)
    try {
      await installInstance({
        name,
        version,
        memory_limit: memory,
        motd: motd || undefined,
        level_name: levelName || undefined,
        seed: seed || undefined,
        rcon_enabled: rconEnabled,
        rcon_password: rconEnabled ? rconPassword : undefined,
      })
      navigate('/')
    } catch (err) {
      setError(err instanceof Error ? err.message : 'install failed')
      setBusy(false)
    }
  }

  return (
    <div className="page">
      <h1>Install a new server</h1>
      <p className="muted">
        The instance is created under <code>instances/</code> in the API home folder.
      </p>
      <form className="card form" onSubmit={onSubmit}>
        <label>
          Instance name *
          <input
            value={name}
            onChange={(e) => setName(e.target.value)}
            placeholder="my-server"
            pattern="[A-Za-z0-9][A-Za-z0-9._-]{0,63}"
            title="Letters, digits, '.', '_' or '-'"
            required
          />
        </label>

        <label>
          Minecraft version
          <input
            list="mc-versions"
            value={version}
            onChange={(e) => setVersion(e.target.value)}
            placeholder="latest"
          />
          <datalist id="mc-versions">
            <option value="latest" />
            {versions.map((v) => (
              <option key={v} value={v} />
            ))}
          </datalist>
          {versionsError && <span className="muted">could not load version list: {versionsError}</span>}
        </label>

        <label>
          Memory limit
          <input value={memory} onChange={(e) => setMemory(e.target.value)} placeholder="1g" />
        </label>

        <label>
          MOTD
          <input value={motd} onChange={(e) => setMotd(e.target.value)} placeholder="A Minecraft Server" />
        </label>

        <label>
          Level name
          <input value={levelName} onChange={(e) => setLevelName(e.target.value)} placeholder="world" />
        </label>

        <label>
          Seed
          <input value={seed} onChange={(e) => setSeed(e.target.value)} />
        </label>

        <label className="inline">
          <input
            type="checkbox"
            checked={rconEnabled}
            onChange={(e) => setRconEnabled(e.target.checked)}
          />
          Enable RCON
        </label>

        {rconEnabled && (
          <label>
            RCON password *
            <input
              type="password"
              value={rconPassword}
              onChange={(e) => setRconPassword(e.target.value)}
              required={rconEnabled}
            />
          </label>
        )}

        {error && <p className="error">{error}</p>}

        <button type="submit" disabled={busy || !name}>
          {busy ? 'Installing… (downloads may take a few minutes)' : 'Install'}
        </button>
      </form>
    </div>
  )
}
