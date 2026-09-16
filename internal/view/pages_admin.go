package view

// Page body templates for the admin console (/admin and sub-surfaces).
// Rendered through the shared layout with AppUI:true. Super-admin-only pages
// (keys) are gated by the handler, not the template.

func init() {
	RegisterPageTemplate("dash-admin", pageDashAdmin)
	RegisterPageTemplate("dash-people", pageDashPeople)
	RegisterPageTemplate("dash-keys", pageDashKeys)
	RegisterPageTemplate("dash-apikeys", pageDashApiKeys)
}

// Data keys for dash-admin: Stats []map[string]string{v,k}; Recent []map[string]string{name,district,status,date}.
const pageDashAdmin = `
<div class="app">
  <div class="wrap">
    <div class="app-head">
      <div>
        <h1>Admin Console</h1>
        <p class="sub">Command center for Patriot Pest Control operations.</p>
      </div>
      <div class="actions">
        <a class="btn btn-ghost" href="/admin/people">People</a>
        <a class="btn btn-ghost" href="/admin/api-keys">API Keys</a>
        <a class="btn btn-primary" href="/admin/keys">Keys ▸</a>
      </div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <div class="stat-cards">
      {{ range .Stats }}
      <div class="stat-card"><span class="v">{{ .v }}</span><span class="k">{{ .k }}</span></div>
      {{ end }}
    </div>

    <div class="tile-grid" style="margin-bottom:2rem">
      <a class="tile" href="/admin/people"><span class="ico">👥</span><h3>People</h3><p>Add staff, set titles &amp; permissions.</p></a>
      <a class="tile" href="/admin/keys"><span class="ico">🔑</span><h3>District Keys</h3><p>Twilio + FieldRoutes keys per district.</p></a>
      <a class="tile" href="/admin/api-keys"><span class="ico">🧩</span><h3>API Keys</h3><p>Issue &amp; revoke live API access.</p></a>
      <a class="tile" href="/admin/posts"><span class="ico">📝</span><h3>Posts</h3><p>Blog &amp; knowledge base content.</p></a>
    </div>

    <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Recent Customers</h3>
    {{ if .Recent }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Name</th><th>District</th><th>Status</th><th>Last Service</th></tr></thead>
        <tbody>
          {{ range .Recent }}
          <tr><td>{{ .name }}</td><td>{{ .district }}</td><td>{{ .status }}</td><td>{{ .date }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No customers synced yet. Run a sync to pull live districts.</p>
    {{ end }}
  </div>
</div>
`

// Data keys for dash-people: Csrf template.HTML; Staff []map[string]string{email,name,role,title,perm}.
const pageDashPeople = `
<div class="app">
  <div class="wrap">
    <div class="crumb"><a href="/admin">Admin</a><span class="sep">/</span> People</div>
    <div class="app-head">
      <div>
        <h1>People Management</h1>
        <p class="sub">Staff roster, titles and role permissions.</p>
      </div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <form method="post" action="/admin/people" novalidate>
      {{ .Csrf }}
      <div class="form-stack wide">
        <div class="form-row">
          <div class="field"><label for="name">Full Name</label><input id="name" name="name" required></div>
          <div class="field"><label for="email">Email</label><input id="email" type="email" name="email" required></div>
        </div>
        <div class="form-row">
          <div class="field"><label for="role">Role</label>
            <select id="role" name="role">
              <option value="staff">Staff</option>
              <option value="admin">Admin</option>
              <option value="super-user">Super-User</option>
            </select>
          </div>
          <div class="field"><label for="title">Title</label><input id="title" name="title" placeholder="e.g. Senior Technician"></div>
        </div>
      </div>
      <div class="form-actions"><button type="submit" class="btn btn-primary">Add Person ▸</button></div>
    </form>

    <h3 style="font-family:var(--display);color:var(--cream);margin:2rem 0 1rem">Roster</h3>
    {{ if .Staff }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Name</th><th>Email</th><th>Role</th><th>Title</th></tr></thead>
        <tbody>
          {{ range .Staff }}
          <tr><td>{{ .name }}</td><td>{{ .email }}</td><td><span class="badge role">{{ .role }}</span></td><td>{{ .title }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No people on file.</p>
    {{ end }}
  </div>
</div>
`

// Data keys for dash-keys: Districts []map[string]string{code,base,key,color}; Twilio map[string]string{sid,token,phone}.
const pageDashKeys = `
<div class="app">
  <div class="wrap">
    <div class="crumb"><a href="/admin">Admin</a><span class="sep">/</span> Keys</div>
    <div class="app-head">
      <div>
        <h1>District &amp; Twilio Keys</h1>
        <p class="sub">Super-admin only. FieldRoutes API keys per district, plus the Twilio voice/SMS line.</p>
      </div>
      <div class="actions"><span class="badge role">SUPER-ADMIN</span></div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <form method="post" action="/admin/keys" novalidate>
      {{ .Csrf }}
      <div class="form-stack wide">
        <div class="form-row">
          <div class="field"><label for="code">District Code</label><input id="code" name="code" placeholder="e.g. tx"></div>
          <div class="field"><label for="base">Base URL</label><input id="base" name="base" placeholder="https://acme.fieldroutes.com"></div>
        </div>
        <div class="form-row">
          <div class="field"><label for="key">Auth Key</label><input id="key" name="key"></div>
          <div class="field"><label for="token">Auth Token</label><input id="token" name="token"></div>
        </div>
        <div class="form-row">
          <div class="field"><label for="color">Label Color (optional)</label><input id="color" name="color" placeholder="#2563eb"></div>
        </div>
      </div>
      <div class="form-actions"><button type="submit" class="btn btn-primary">Add District ▸</button></div>
    </form>

    <h3 style="font-family:var(--display);color:var(--cream);margin:2rem 0 1rem">FieldRoutes Districts</h3>
    {{ if .Districts }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>District</th><th>Base URL</th><th>Key</th><th>Label</th></tr></thead>
        <tbody>
          {{ range .Districts }}
          <tr>
            <td><span style="display:inline-block;width:10px;height:10px;border-radius:50%;background:{{ .color }};margin-right:.4rem"></span>{{ .code }}</td>
            <td class="mono">{{ .base }}</td>
            <td class="mono">{{ .key }}</td>
            <td>{{ .code }}</td>
          </tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No district keys configured.</p>
    {{ end }}

    <h3 style="font-family:var(--display);color:var(--cream);margin:2rem 0 1rem">Twilio Line</h3>
    <dl class="kv">
      <dt>SID</dt><dd class="mono">{{ .Twilio.sid }}</dd>
      <dt>Auth Token</dt><dd class="mono">{{ .Twilio.token }}</dd>
      <dt>Phone</dt><dd class="mono">{{ .Twilio.phone }}</dd>
    </dl>
  </div>
</div>
`

// Data keys for dash-api-keys: Csrf template.HTML; Keys []map[string]string{label,token,scopes,created}.
const pageDashApiKeys = `
<div class="app">
  <div class="wrap">
    <div class="crumb"><a href="/admin">Admin</a><span class="sep">/</span> API Keys</div>
    <div class="app-head">
      <div>
        <h1>API Keys</h1>
        <p class="sub">Live keys for the /api/v1 surface. Scope-limited and revocable.</p>
      </div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <form method="post" action="/admin/api-keys" novalidate>
      {{ .Csrf }}
      <div class="form-stack">
        <div class="field"><label for="label">Key Label</label><input id="label" name="label" placeholder="e.g. Reporting"></div>
        <div class="field"><label for="scopes">Scopes</label><input id="scopes" name="scopes" placeholder="customer:read,ticket:read"></div>
      </div>
      <div class="form-actions"><button type="submit" class="btn btn-primary">Issue Key ▸</button></div>
    </form>

    <h3 style="font-family:var(--display);color:var(--cream);margin:2rem 0 1rem">Issued Keys</h3>
    {{ if .Keys }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Label</th><th>Token</th><th>Scopes</th><th>Created</th><th></th></tr></thead>
        <tbody>
          {{ range .Keys }}
          <tr><td>{{ .label }}</td><td class="mono">{{ .token }}</td><td>{{ .scopes }}</td><td>{{ .created }}</td>
            <td style="text-align:right"><form method="post" action="/admin/api-keys/revoke">{{ .Csrf }}<input type="hidden" name="token" value="{{ .token }}"><button type="submit" class="btn btn-ghost">Revoke</button></form></td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No API keys issued.</p>
    {{ end }}
  </div>
</div>
`
