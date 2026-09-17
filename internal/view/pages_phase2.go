package view

// Phase 2 page bodies — Salesforce-killer dashboards: Customer 360 timeline
// and the USA customer density heatmap (Leaflet).

func init() {
	RegisterPageTemplate("customer-360", pageCustomer360)
	RegisterPageTemplate("dash-heatmap", pageDashHeatmap)
	RegisterPageTemplate("dash-inbox", pageDashInbox)
}

// Data keys: Found bool, Customer GeoRow, Events []map{icon,kind,at,text}.
const pageCustomer360 = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Customer 360</h1>
      <p class="sub">Unified timeline — calls, texts, visits and messages in one feed.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/staff/customers">◂ All Customers</a></div>
  </div>

  {{ if not .Found }}
  <div class="panel"><p class="empty">Customer not found in the local cache. Run a FieldRoutes sync from the staff dashboard.</p></div>
  {{ else }}
  <div class="stat-cards" style="margin-bottom:1.6rem">
    <div class="stat-card"><span class="v">{{ .Customer.Name }}</span><span class="k">Account Holder</span></div>
    <div class="stat-card"><span class="v">{{ .Customer.AccountNumber }}</span><span class="k">Account No.</span></div>
    <div class="stat-card"><span class="v" style="text-transform:uppercase">{{ .Customer.District }}</span><span class="k">District</span></div>
    <div class="stat-card"><span class="v">{{ .Customer.Status }}</span><span class="k">Status</span></div>
  </div>

  <div class="panel" style="margin-bottom:1.6rem">
    <dl class="kv">
      <dt>City</dt><dd>{{ if .Customer.City }}{{ .Customer.City }}, {{ .Customer.State }}{{ else }}—{{ end }}</dd>
      <dt>Last Service</dt><dd>{{ if .Customer.LastService }}{{ .Customer.LastService }}{{ else }}—{{ end }}</dd>
      <dt>FieldRoutes ID</dt><dd class="mono">{{ .Customer.FRID }}</dd>
    </dl>
    <div style="margin-top:1rem;display:flex;gap:.6rem;flex-wrap:wrap">
      <a class="btn btn-ghost" href="/customer/messages">💬 Message Thread</a>
      <a class="btn btn-ghost" href="/admin/fieldroutes">🛰 Sync Health</a>
    </div>
  </div>

  <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Timeline</h3>
  {{ if .Events }}
  <div style="position:relative;padding-left:1.6rem;border-left:2px solid var(--olive-700);display:flex;flex-direction:column;gap:1.1rem">
    {{ range .Events }}
    <div style="position:relative">
      <span style="position:absolute;left:-2.35rem;top:.1rem;font-size:1.1rem">{{ .icon }}</span>
      <div style="font-family:var(--mono);font-size:.68rem;letter-spacing:.14em;color:var(--olive-300)">{{ .kind }} · {{ .at }}</div>
      <div style="color:var(--cream);margin-top:.15rem">{{ .text }}</div>
    </div>
    {{ end }}
  </div>
  {{ else }}
  <p class="empty">No timeline events yet — appointments and messages will appear here.</p>
  {{ end }}
  {{ end }}
</div></div>
`

// Leaflet density map. Data loads from /api/admin/heatmap.json client-side.
const pageDashHeatmap = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Customer Density Heatmap</h1>
      <p class="sub">USA drill-down — zoom from national view to Spokane / Phoenix neighborhoods.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/admin">◂ Admin Console</a></div>
  </div>

  <div class="panel" style="padding:0;overflow:hidden">
    <div id="hm-map" style="height:640px;width:100%;background:var(--olive-900)"></div>
  </div>
  <p class="muted" style="margin-top:.8rem">Circle size = customers per city. Click a circle for the account list. Cities without geocodes fall back to district centroids.</p>
</div></div>

<link rel="stylesheet" href="https://unpkg.com/leaflet@1.9.4/dist/leaflet.css">
<script src="https://unpkg.com/leaflet@1.9.4/dist/leaflet.js"></script>
<script>
(function () {
  var map = L.map('hm-map', { zoomControl: true }).setView([39.5, -98.35], 4);
  L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
    attribution: '&copy; OpenStreetMap contributors', maxZoom: 18
  }).addTo(map);

  fetch('/api/admin/heatmap.json').then(function (r) { return r.json(); }).then(function (data) {
    var spots = data.spots || [];
    var max = 1;
    spots.forEach(function (s) { if (s.count > max) max = s.count; });
    spots.forEach(function (s) {
      var radius = 8 + 26 * (s.count / max);
      var c = L.circle([s.lat, s.lng], {
        radius: radius * 400, color: '#f4772e', weight: 1.5,
        fillColor: '#f4772e', fillOpacity: 0.35
      }).addTo(map);
      // Names come from customer records: escape before building popup HTML.
      function esc(v) {
        return String(v == null ? '' : v).replace(/[&<>"']/g, function (ch) {
          return { '&': '&amp;', '<': '&lt;', '>': '&gt;', '"': '&quot;', "'": '&#39;' }[ch];
        });
      }
      var names = (s.names || []).map(esc).join('<br>');
      c.bindPopup('<b>' + esc(s.city) + ', ' + esc(s.state) + '</b><br>' + s.count + ' customer(s)<br><span style="font-size:.8em">' + names + '</span>');
    });
    if (spots.length) {
      var b = L.latLngBounds(spots.map(function (s) { return [s.lat, s.lng]; }));
      map.fitBounds(b.pad(0.4));
    }
  }).catch(function () {});
})();
</script>
`


// Data keys: Csrf, Channels []map{key,label,glyph}, Threads []map{id,channel,
// glyph,customer,preview,updated,human,bubbles[]map{from,text,at,staff}}.
const pageDashInbox = `
<div class="app"><div class="wrap">
  <div class="app-head">
    <div>
      <h1>Unified Inbox</h1>
      <p class="sub">Website, SMS, voicemail, email and social DMs in one threaded stream.</p>
    </div>
    <div class="actions"><a class="btn btn-ghost" href="/staff-dashboard">◂ Dashboard</a></div>
  </div>

  <div class="inbox-rail" style="display:flex;gap:.5rem;flex-wrap:wrap;margin-bottom:1.2rem">
    {{ range .Channels }}
    <span class="badge" style="display:inline-flex;align-items:center;gap:.35rem" title="{{ .label }}">{{ .glyph }} {{ .label }}</span>
    {{ end }}
  </div>

  {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

  {{ if not .Threads }}
  <div class="panel"><p class="empty">All caught up — no open threads. Incoming website messages, SMS, voicemails and social DMs will land here.</p></div>
  {{ else }}
  <div style="display:flex;flex-direction:column;gap:1rem">
    {{ range .Threads }}
    <div class="panel">
      <div style="display:flex;justify-content:space-between;align-items:center;gap:.8rem;flex-wrap:wrap;margin-bottom:.6rem">
        <div style="display:flex;align-items:center;gap:.6rem">
          <span style="font-size:1.2rem">{{ .glyph }}</span>
          <strong style="color:var(--cream)">{{ .customer }}</strong>
          <span class="badge">{{ .channel }}</span>
          {{ if .human }}<span class="badge active">STAFF HANDLING</span>{{ end }}
        </div>
        <span class="muted" style="font-family:var(--mono);font-size:.7rem">{{ .updated }}</span>
      </div>

      <div style="display:flex;flex-direction:column;gap:.5rem;margin-bottom:.8rem">
        {{ range .bubbles }}
        <div style="max-width:78%;padding:.55rem .8rem;border-radius:9px;border:1px solid var(--olive-700);background:{{ if .staff }}var(--olive-700){{ else }}var(--olive-900){{ end }};align-self:{{ if .staff }}flex-end{{ else }}flex-start{{ end }}">
          <div style="font-family:var(--mono);font-size:.62rem;color:var(--olive-300)">{{ .from }} · {{ .at }}</div>
          <div style="color:var(--cream);margin-top:.15rem">{{ .text }}</div>
        </div>
        {{ end }}
      </div>

      <form method="post" action="/staff/messages/reply" style="display:flex;gap:.6rem">
        {{ $.Csrf }}
        <input type="hidden" name="thread" value="{{ .id }}">
        <input type="text" name="text" required maxlength="2000" placeholder="Reply as staff — bot auto-mutes on this thread…" style="flex:1">
        <button type="submit" class="btn btn-primary">Send ▸</button>
      </form>
    </div>
    {{ end }}
  </div>
  {{ end }}
</div></div>
`
