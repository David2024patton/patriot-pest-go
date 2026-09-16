package view

// Page body templates for the authenticated app shell (customer + staff dashboards).
// Rendered through the shared layout with AppUI:true so /assets/admin.css loads.
// Each template is a page body only; the layout supplies nav/footer/SEO.

func init() {
	RegisterPageTemplate("dash-customer", pageDashCustomer)
	RegisterPageTemplate("dash-staff", pageDashStaff)
}

// Data keys for dash-customer: Name, Email, Phone, AccountNumber, Address string;
// Appointments []map[string]string{date,start,type,status}; Invoices []map[string]string{id,date,amount,status}.
const pageDashCustomer = `
<div class="app">
  <div class="wrap">
    <div class="app-head">
      <div>
        <h1>My Account</h1>
        <p class="sub">Your Patriot Pest Control profile, appointments and invoices.</p>
      </div>
      <div class="actions">
        <a class="btn btn-ghost" href="/customer/messages">Messages</a>
        <a class="btn btn-primary" href="/customer/account">Account Details ▸</a>
      </div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <div class="stat-cards">
      <div class="stat-card"><span class="v">{{ .Name }}</span><span class="k">Account Holder</span></div>
      <div class="stat-card"><span class="v">{{ .AccountNumber }}</span><span class="k">Account No.</span></div>
    </div>

    <div class="panel" style="margin-bottom:2rem">
      <h3>Profile</h3>
      <dl class="kv">
        <dt>Email</dt><dd>{{ .Email }}</dd>
        <dt>Phone</dt><dd>{{ .Phone }}</dd>
        <dt>Address</dt><dd>{{ .Address }}</dd>
      </dl>
    </div>

    <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Upcoming Appointments</h3>
    {{ if .Appointments }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Date</th><th>Start</th><th>Type</th><th>Status</th></tr></thead>
        <tbody>
          {{ range .Appointments }}
          <tr><td>{{ .date }}</td><td>{{ .start }}</td><td>{{ .type }}</td><td>{{ .status }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No upcoming appointments on file.</p>
    {{ end }}

    <h3 style="font-family:var(--display);color:var(--cream);margin:2rem 0 1rem">Invoices</h3>
    {{ if .Invoices }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Invoice</th><th>Date</th><th>Amount</th><th>Status</th></tr></thead>
        <tbody>
          {{ range .Invoices }}
          <tr><td>{{ .id }}</td><td>{{ .date }}</td><td>{{ .amount }}</td><td>{{ .status }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No invoices issued yet.</p>
    {{ end }}
  </div>
</div>
`

// Data keys for dash-staff: Stats []map[string]string{v,k}; Appointments []map[string]string{date,start,type,status,district}.
const pageDashStaff = `
<div class="app">
  <div class="wrap">
    <div class="app-head">
      <div>
        <h1>Staff Dashboard</h1>
        <p class="sub">Operations overview across every district.</p>
      </div>
      <div class="actions">
        <a class="btn btn-ghost" href="/admin">Admin Console</a>
        <a class="btn btn-primary" href="/staff/appointments">Appointments ▸</a>
      </div>
    </div>

    {{ if .Flash }}<div class="notice info">{{ .Flash }}</div>{{ end }}

    <div class="stat-cards">
      {{ range .Stats }}
      <div class="stat-card"><span class="v">{{ .v }}</span><span class="k">{{ .k }}</span></div>
      {{ end }}
    </div>

    <h3 style="font-family:var(--display);color:var(--cream);margin-bottom:1rem">Latest Appointments</h3>
    {{ if .Appointments }}
    <div class="table-wrap">
      <table class="data">
        <thead><tr><th>Date</th><th>Start</th><th>Type</th><th>District</th><th>Status</th></tr></thead>
        <tbody>
          {{ range .Appointments }}
          <tr><td>{{ .date }}</td><td>{{ .start }}</td><td>{{ .type }}</td><td>{{ .district }}</td><td>{{ .status }}</td></tr>
          {{ end }}
        </tbody>
      </table>
    </div>
    {{ else }}
    <p class="empty">No appointments recorded this week.</p>
    {{ end }}
  </div>
</div>
`
