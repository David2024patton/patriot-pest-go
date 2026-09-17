package view

// Phase 1 page bodies — real replacements for the legacy JSON stubs
// (PLANS/MASTERPLAN_GO_REWRITE.md §Phase 1). App-shell pages render with
// AppUI:true; auth surfaces reuse the authx design system.

func init() {
	RegisterPageTemplate("account", pageAccount)
	RegisterPageTemplate("customer-messages", pageCustomerMessages)
	RegisterPageTemplate("su-login", pageSuLogin)
	RegisterPageTemplate("su-verify", pageSuVerify)
	RegisterPageTemplate("tech-ask", pageTechAsk)
	RegisterPageTemplate("dash-fieldroutes", pageDashFieldRoutes)
	RegisterPageTemplate("admin-post-edit", pageAdminPostEdit)
	RegisterPageTemplate("admin-staff-edit", pageAdminStaffEdit)
	RegisterPageTemplate("admin-settings", pageAdminSettings)
	RegisterPageTemplate("apikey-audit", pageApiKeyAudit)
	RegisterPageTemplate("customer-portal", pageCustomerPortal)
}

// Data keys: Type ("staff"|"customer"), Name, RoleLabel, Record map, Flash.
const pageAccount = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>My Account</h1>
      <div class="sub">Signed in as {{ .Name }} · <span class="badge role">{{ .RoleLabel }}</span></div>
    </div>
    <div class="actions">
      {{ if eq .Type "staff" }}<a class="btn btn-ghost" href="/staff-dashboard">◂ Dashboard</a>
      {{ else }}<a class="btn btn-ghost" href="/customer-dashboard">◂ Dashboard</a>{{ end }}
      <a class="btn btn-ghost" href="/logout">Sign Out</a>
    </div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  {{ if not .Record }}
    <div class="panel"><p class="empty">We couldn't load your account record.</p></div>
  {{ else if eq .Type "staff" }}
    <div class="panel">
      <h3>Staff Profile</h3>
      <dl class="kv" style="margin-top:.8rem">
        <dt>Name</dt><dd>{{ index .Record "name" }}</dd>
        <dt>Email</dt><dd>{{ index .Record "email" }}</dd>
        <dt>Role</dt><dd><span class="badge role">{{ index .Record "role" }}</span></dd>
        <dt>Status</dt><dd><span class="badge active">Active</span></dd>
      </dl>
    </div>
    <p class="muted" style="margin-top:1rem;line-height:1.7">Sign-in is passwordless: we email a one-time code each time. To change your name or email, contact an administrator.</p>
  {{ else }}
    <div class="panel">
      <h3>Account Details</h3>
      <dl class="kv" style="margin-top:.8rem">
        <dt>Account #</dt><dd class="mono">{{ index .Record "account_number" }}</dd>
        <dt>Name</dt><dd>{{ index .Record "name" }}</dd>
        <dt>Email</dt><dd>{{ index .Record "email" }}</dd>
        <dt>Phone</dt><dd>{{ index .Record "phone" }}</dd>
        <dt>District</dt><dd>{{ index .Record "district" }}</dd>
        <dt>Status</dt><dd><span class="badge active">{{ index .Record "status" }}</span></dd>
      </dl>
    </div>
    <p class="muted" style="margin-top:1rem;line-height:1.7">Need to update your contact info or service address? Call <a href="{{ .PhoneHref }}">{{ .PhoneDisplay }}</a> or use the <a href="/contact">contact form</a>.</p>
  {{ end }}
</div></div>
`

// Data keys: Csrf, Flash, Messages []map{from,text,at,mine}.
const pageCustomerMessages = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Messages</h1>
      <p class="sub">Talk to the Patriot dispatch desk — we answer fast.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/customer-dashboard">◂ Dashboard</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  <div class="panel" style="margin-bottom:1.4rem">
    {{ if .Messages }}
    <div style="display:flex;flex-direction:column;gap:.8rem">
      {{ range .Messages }}
      <div style="max-width:70%;padding:.7rem .9rem;border-radius:10px;background:{{ if .mine }}var(--olive-700){{ else }}var(--olive-900){{ end }};border:1px solid var(--olive-700);align-self:{{ if .mine }}flex-end{{ else }}flex-start{{ end }}">
        <div style="font-size:.72rem;color:var(--olive-300);font-family:var(--mono)">{{ .from }} · {{ .at }}</div>
        <div style="color:var(--cream);margin-top:.2rem">{{ .text }}</div>
      </div>
      {{ end }}
    </div>
    {{ else }}
    <p class="empty">No messages yet. Send us a note below.</p>
    {{ end }}
  </div>

  <form method="post" action="/api/customer/messages" class="panel">
    {{ .Csrf }}
    <div class="form-group">
      <label for="msg">New message</label>
      <textarea id="msg" name="message" rows="3" required maxlength="2000" placeholder="Question about your next visit, a pest sighting, billing…"></textarea>
    </div>
    <button type="submit" class="btn btn-primary">Send Message ▸</button>
  </form>
</div></div>
`

// Data keys: Csrf, FlashError, SentTo.
const pageSuLogin = `
<div class="authx">
  <div class="authx-form" style="margin:4rem auto">
    <div class="authx-card2">
      <span class="step">ELEVATED ACCESS</span>
      <h1>Superuser Sign In</h1>
      <p class="sub">Enter your email to receive a secure 8-digit sign-in code. This surface is restricted to command-level accounts.</p>

      {{ if .FlashError }}<div class="notice error">{{ .FlashError }}</div>{{ end }}
      {{ if .SentTo }}<div class="notice info">Code sent to <b>{{ .SentTo }}</b>. Enter it below.</div>{{ end }}

      <form method="post" action="/su" novalidate>
        {{ .Csrf }}
        <div class="authx-field">
          <label for="email">Email address</label>
          <input type="email" id="email" name="email" required autofocus autocomplete="email" placeholder="you@example.com" maxlength="254">
        </div>
        <button type="submit" class="authx-btn">Send Code ▸</button>
      </form>

      <p class="authx-foot"><a href="/login">Staff &amp; customer sign in</a></p>
    </div>
  </div>
</div>
`

// Data keys: Csrf, FlashError, SentTo.
const pageSuVerify = `
<div class="authx">
  <div class="authx-form" style="margin:4rem auto">
    <div class="authx-card2">
      <span class="step">COMMAND VERIFICATION</span>
      <h1>Enter your 8-digit code</h1>
      <p class="sub">{{ if .SentTo }}We sent an 8-digit code to <span class="sent">{{ .SentTo }}</span>. It expires in 5 minutes and works once.{{ else }}Type the 8-digit code we sent you.{{ end }}</p>

      {{ if .FlashError }}<div class="notice error">{{ .FlashError }}</div>{{ end }}

      <form method="post" action="/su/verify" novalidate>
        {{ .Csrf }}
        <div class="authx-field">
          <label for="code">8-digit code</label>
          <input type="text" id="code" name="code" class="authx-code" inputmode="numeric" autocomplete="one-time-code" pattern="[0-9]{8}" maxlength="8" required autofocus placeholder="••••••••">
        </div>
        <button type="submit" class="authx-btn">Verify &amp; Sign In ▸</button>
      </form>

      <p class="authx-foot"><a href="/su">◂ Use a different email</a></p>
    </div>
  </div>
</div>
`

// Data keys: Q, Hits []map{title,url,snippet}.
const pageTechAsk = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Field Tech Copilot</h1>
      <p class="sub">Search the pest knowledge base, treatment guides and blog corpus.</p>
    </div>
  </div>

  <form method="get" action="/tech/ask" class="panel" style="margin-bottom:1.4rem">
    <div class="form-group">
      <label for="q">Ask the copilot</label>
      <input type="text" id="q" name="q" value="{{ .Q }}" required autofocus placeholder="e.g. carpenter ant treatment, rodent bait stations, wasp nest removal">
    </div>
    <button type="submit" class="btn btn-primary">Search ▸</button>
  </form>

  {{ if .Q }}
    {{ if .Hits }}
    <div class="panel">
      <h3>Results for “{{ .Q }}”</h3>
      <div style="display:flex;flex-direction:column;gap:1rem;margin-top:.8rem">
        {{ range .Hits }}
        <div>
          <a href="{{ .url }}" style="color:var(--orange);font-weight:700">{{ .title }}</a>
          <p class="muted" style="margin:.2rem 0 0">{{ .snippet }}</p>
        </div>
        {{ end }}
      </div>
    </div>
    {{ else }}
    <div class="panel"><p class="empty">No matches in the knowledge base for “{{ .Q }}”. Try a broader pest name.</p></div>
    {{ end }}
  {{ end }}
</div></div>
`

// Data keys: Districts []map{code,base,configured,healthy,latency,customers}, Flash.
const pageDashFieldRoutes = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>FieldRoutes Sync Health</h1>
      <p class="sub">Live district connectivity — WA and AZ sync status.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin">◂ Admin Console</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  <div class="table-wrap">
    <table class="data">
      <thead><tr><th>District</th><th>Base URL</th><th>Keys</th><th>Health</th><th>Latency</th><th>Customers</th></tr></thead>
      <tbody>
        {{ range .Districts }}
        <tr>
          <td style="text-transform:uppercase;font-weight:700">{{ .code }}</td>
          <td class="mono">{{ .base }}</td>
          <td>{{ if .configured }}<span class="badge active">CONFIGURED</span>{{ else }}<span class="badge cancelled">MISSING</span>{{ end }}</td>
          <td>{{ if .healthy }}<span class="badge active">HEALTHY</span>{{ else }}<span class="badge cancelled">DOWN</span>{{ end }}</td>
          <td>{{ .latency }}</td>
          <td>{{ .customers }}</td>
        </tr>
        {{ end }}
      </tbody>
    </table>
  </div>
  <p class="muted" style="margin-top:1rem">Dashboards query the local SQLite cache first — a red row means background sync is degraded, not that the app is down.</p>
</div></div>
`

// Data keys: Csrf, Post map{slug,title,excerpt,body_html,author,season,pest_category,published_at}, IsNew, Flash.
const pageAdminPostEdit = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>{{ if .IsNew }}New Blog Post{{ else }}Edit Post{{ end }}</h1>
      <p class="sub">CMS — writes to the SQLite posts catalog.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin/posts">◂ All Posts</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  {{ if .Post }}
  <form method="post" class="panel" style="display:flex;flex-direction:column;gap:1rem;max-width:860px">
    {{ .Csrf }}
    <div class="form-group"><label for="title">Title</label>
      <input type="text" id="title" name="title" required value="{{ index .Post "title" }}"></div>
    <div class="form-group"><label for="slug">Slug</label>
      <input type="text" id="slug" name="slug" required value="{{ index .Post "slug" }}" placeholder="why-quarterly-pest-control"></div>
    <div style="display:grid;grid-template-columns:1fr 1fr;gap:1rem">
      <div class="form-group"><label for="author">Author</label>
        <input type="text" id="author" name="author" value="{{ index .Post "author" }}"></div>
      <div class="form-group"><label for="pest_category">Pest Category</label>
        <input type="text" id="pest_category" name="pest_category" value="{{ index .Post "pest_category" }}" placeholder="ants"></div>
      <div class="form-group"><label for="season">Season</label>
        <input type="text" id="season" name="season" value="{{ index .Post "season" }}" placeholder="summer"></div>
      <div class="form-group"><label for="published_at">Publish Date</label>
        <input type="date" id="published_at" name="published_at" value="{{ index .Post "published_at" }}"></div>
    </div>
    <div class="form-group"><label for="excerpt">Excerpt</label>
      <textarea id="excerpt" name="excerpt" rows="2">{{ index .Post "excerpt" }}</textarea></div>
    <div class="form-group"><label for="body_html">Body (HTML)</label>
      <textarea id="body_html" name="body_html" rows="14" style="font-family:var(--mono);font-size:.85rem">{{ index .Post "body_html" }}</textarea></div>
    <div><button type="submit" class="btn btn-primary">{{ if .IsNew }}Create Post{{ else }}Save Changes{{ end }} ▸</button></div>
  </form>
  {{ end }}
</div></div>
`

// Data keys: Csrf, StaffRow map{email,name,role,title}, IsNew, Flash.
const pageAdminStaffEdit = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>{{ if .IsNew }}Add Staff Member{{ else }}Edit Staff Member{{ end }}</h1>
      <p class="sub">Passwordless accounts — staff sign in with an emailed one-time code.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin/people">◂ People</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  {{ if .StaffRow }}
  <form method="post" action="/admin/staff" class="panel" style="max-width:560px;display:flex;flex-direction:column;gap:1rem">
    {{ .Csrf }}
    <div class="form-group"><label for="email">Email</label>
      <input type="email" id="email" name="email" required value="{{ index .StaffRow "email" }}" {{ if not .IsNew }}readonly{{ end }}></div>
    <div class="form-group"><label for="name">Name</label>
      <input type="text" id="name" name="name" required value="{{ index .StaffRow "name" }}"></div>
    <div class="form-group"><label for="role">Role</label>
      <select id="role" name="role">
        <option value="staff"{{ if or .IsNew (eq (index .StaffRow "role") "staff") }} selected{{ end }}>staff</option>
        <option value="admin"{{ if eq (index .StaffRow "role") "admin" }} selected{{ end }}>admin</option>
        <option value="super-user"{{ if eq (index .StaffRow "role") "super-user" }} selected{{ end }}>super-user</option>
      </select></div>
    <div class="form-group"><label for="title">Title</label>
      <input type="text" id="title" name="title" value="{{ index .StaffRow "title" }}" placeholder="Dispatch Coordinator"></div>
    <div><button type="submit" class="btn btn-primary">{{ if .IsNew }}Create Staff{{ else }}Save Changes{{ end }} ▸</button></div>
  </form>
  {{ end }}
</div></div>
`

// Data keys: Csrf, Settings map, Flash.
const pageAdminSettings = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Site Settings</h1>
      <p class="sub">Runtime settings persisted to storage/settings.json (hot-reloaded on save).</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin">◂ Admin Console</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  {{ if .Settings }}
  <form method="post" action="/admin/settings" class="panel" style="max-width:680px;display:flex;flex-direction:column;gap:1rem">
    {{ .Csrf }}
    <div class="form-group"><label for="site_name">Site Name</label>
      <input type="text" id="site_name" name="site_name" value="{{ index .Settings "site_name" }}"></div>
    <div class="form-group"><label for="fb_pixel">Meta Pixel ID (global)</label>
      <input type="text" id="fb_pixel" name="fb_pixel" value="{{ index .Settings "fb_pixel" }}" placeholder="123456789012345"></div>
    <div class="form-group"><label for="gtag_id">GA4 Measurement ID</label>
      <input type="text" id="gtag_id" name="gtag_id" value="{{ index .Settings "gtag_id" }}" placeholder="G-XXXXXXXXXX"></div>
    <div class="form-group"><label for="gads_id">Google Ads ID</label>
      <input type="text" id="gads_id" name="gads_id" value="{{ index .Settings "gads_id" }}" placeholder="AW-XXXXXXXXXX"></div>
    <div class="form-group"><label for="clarity_id">MS Clarity ID</label>
      <input type="text" id="clarity_id" name="clarity_id" value="{{ index .Settings "clarity_id" }}"></div>
    <div><button type="submit" class="btn btn-primary">Save Settings ▸</button></div>
  </form>
  {{ end }}
</div></div>
`

// Data keys: Csrf, Keys []map{id,label,prefix,scopes,created}, Audit []map{at,action,label,detail}, Flash.
const pageApiKeyAudit = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>API Key Audit</h1>
      <p class="sub">Lifecycle controls + usage audit trail for ppc_live_ keys.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin/api-keys">◂ API Keys</a></div>
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Issued Keys</h3>
  <div class="table-wrap" style="margin-bottom:2rem">
    <table class="data">
      <thead><tr><th>Label</th><th>Token</th><th>Scopes</th><th>Created</th><th>Actions</th></tr></thead>
      <tbody>
        {{ range .Keys }}
        <tr>
          <td>{{ .label }}</td>
          <td class="mono">{{ .prefix }}…</td>
          <td>{{ .scopes }}</td>
          <td>{{ .created }}</td>
          <td style="white-space:nowrap">
            <form method="post" action="/admin/api-keys/{{ .id }}/rotate" style="display:inline">{{ $.Csrf }}<button class="btn btn-ghost" type="submit">Rotate</button></form>
            <form method="post" action="/admin/api-keys/{{ .id }}/revoke" style="display:inline" onsubmit="return confirm('Revoke this key?')">{{ $.Csrf }}<button class="btn btn-ghost" type="submit">Revoke</button></form>
          </td>
        </tr>
        {{ end }}
      </tbody>
    </table>
  </div>

  <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Audit Trail</h3>
  {{ if .Audit }}
  <div class="table-wrap">
    <table class="data">
      <thead><tr><th>When</th><th>Action</th><th>Key</th><th>Detail</th></tr></thead>
      <tbody>
        {{ range .Audit }}
        <tr><td class="mono">{{ .at }}</td><td>{{ .action }}</td><td>{{ .label }}</td><td>{{ .detail }}</td></tr>
        {{ end }}
      </tbody>
    </table>
  </div>
  {{ else }}
  <p class="empty">No audit events recorded yet.</p>
  {{ end }}
</div></div>
`

// Landing alias page for /customer-portal.
const pageCustomerPortal = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Customer Portal</h1>
      <p class="sub">Your account, appointments, invoices and messages in one place.</p>
    </div>
  </div>
  <div class="stat-cards">
    <a class="stat-card" href="/customer-dashboard" style="text-decoration:none"><span class="v">▸</span><span class="k">My Dashboard</span></a>
    <a class="stat-card" href="/customer/messages" style="text-decoration:none"><span class="v">▸</span><span class="k">Messages</span></a>
    <a class="stat-card" href="/account" style="text-decoration:none"><span class="v">▸</span><span class="k">Account Details</span></a>
    <a class="stat-card" href="/prices" style="text-decoration:none"><span class="v">▸</span><span class="k">Plans &amp; Pricing</span></a>
  </div>
  <p class="muted" style="margin-top:1.4rem">Payments are handled securely by FieldRoutes — Patriot never stores card numbers. Use “Make a Payment” inside your dashboard.</p>
</div></div>
`
