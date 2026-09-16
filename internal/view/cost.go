package view

// CostPage is the standalone "Project Valuation Receipt" document served at
// /cost. It bypasses the shared marketing layout (own shell, own nav) exactly
// like public/cost/index.php on the legacy site. Assets: /assets/cost.css +
// /assets/cost.js (embedded); data: GET /cost/data/pricing.json.
const CostPage = `<!DOCTYPE html>
<html lang="en">
<head>
  <meta charset="UTF-8">
  <meta name="viewport" content="width=device-width, initial-scale=1.0">
  <title>Project Valuation Receipt | Patriot Pest Control</title>
  <meta name="robots" content="noindex, follow">
  <link rel="preconnect" href="https://fonts.googleapis.com">
  <link rel="preconnect" href="https://fonts.gstatic.com" crossorigin>
  <link href="https://fonts.googleapis.com/css2?family=Barlow:wght@400;500;600;700;800&family=Black+Ops+One&family=IBM+Plex+Mono:wght@400;500;600&display=swap" rel="stylesheet">
  <link rel="stylesheet" href="/assets/cost.css">
</head>
<body class="cost-shell">
  <canvas id="bugfield" aria-hidden="true"></canvas>
  <div class="grain" aria-hidden="true"></div>

  <nav class="cost-nav">
    <a href="/" class="brand">
      <span class="star">★</span>
      PATRIOT PEST CONTROL
    </a>
    <a href="/" class="back-link">← Return to Main Site</a>
  </nav>

  <section class="cost-hero">
    <div class="wrap">
      <div class="eyebrow">// QUARTERMASTER REPORT</div>
      <h1>Project <em>Valuation</em> Receipt</h1>
      <p class="sub">Hostile invoice incoming. This is what a custom web application with CRM integration, payment processing, and tactical design actually costs on the open market. No fluff, no hidden fees — just the raw numbers.</p>
    </div>
  </section>

  <section class="cost-content">
    <div class="wrap">

      <div class="summary-cards">
        <div class="summary-card">
          <div class="label">Agency Floor</div>
          <div class="value" id="val-min">$75,000</div>
          <div class="sub">Entry-Level Quote</div>
        </div>
        <div class="summary-card">
          <div class="label">Agency Ceiling</div>
          <div class="value" id="val-max">$150,000</div>
          <div class="sub">Premium Market Rate</div>
        </div>
      </div>

      <div class="receipt">
        <div class="receipt-header">
          <div class="receipt-logo">★ PATRIOT PEST CONTROL</div>
          <div class="receipt-meta">
            <span>REPORT #: PPC-2026-001</span>
            <span>DATE: <span id="receipt-date"></span></span>
            <span>CLASSIFICATION: UNCLASSIFIED</span>
          </div>
          <div class="hazard" style="margin:0.8rem 0 0 0"></div>
        </div>

        <div class="receipt-body">
          <div class="receipt-headings">
            <span class="rh-cat">LINE ITEM</span>
            <span class="rh-desc">SCOPE OF WORK</span>
            <span class="rh-cost">COST RANGE</span>
            <span class="rh-meter">THREAT METER</span>
          </div>

          <div class="receipt-items" id="receipt-items">
          </div>

          <div class="receipt-totals">
            <div class="total-row">
              <span>SUBTOTAL (LOW ESTIMATE)</span>
              <span class="total-amt" id="subtotal-low">$75,000</span>
            </div>
            <div class="total-row">
              <span>SUBTOTAL (HIGH ESTIMATE)</span>
              <span class="total-amt" id="subtotal-high">$150,000</span>
            </div>
            <div class="total-divider"></div>
            <div class="total-row grand">
              <span>GRAND TOTAL RANGE</span>
              <span class="total-amt grand-amt" id="grand-total">$75,000 — $150,000</span>
            </div>
            <div class="stamp-row">
              <span class="stamp">APPROVED<br>FOR REVIEW</span>
            </div>
          </div>
        </div>

        <div class="receipt-footer">
          <p><strong>PAYMENT TERMS:</strong> Net 30 upon contract signing. 50% upfront, 50% on delivery.</p>
          <p><strong>TIMELINE:</strong> <span id="receipt-timeline">3-6 months</span> — Full agency team deployment.</p>
          <p><strong>WARRANTY:</strong> 90-day bug-free guarantee. Pests and pixel bugs eliminated.</p>
        </div>
      </div>

      <div class="factors-section">
        <h3>Field Notes: Agency Pricing Factors</h3>
        <ul class="factors-list" id="factors-list">
        </ul>
      </div>

    </div>
  </section>

  <script src="/assets/cost.js"></script>
</body>
</html>
`

// CostPricingJSON is the embedded pricing dataset (public/cost/data/pricing.json).
// Served at GET /cost/data/pricing.json — cost.js fetches it relative to /cost.
const CostPricingJSON = `{
  "agency_range": { "min": 75000, "max": 150000, "currency": "USD" },
  "breakdown": [
    { "category": "Custom Framework Architecture", "min": 15000, "max": 25000, "description": "Custom PHP 8.3 MVC framework with no off-the-shelf CMS", "color": "#f4772e" },
    { "category": "CRM Integrations", "min": 10000, "max": 20000, "description": "FieldRoutes CRM + Twilio SMS, voice, voicemail, and webhooks", "color": "#ff8c3b" },
    { "category": "Payment Portal", "min": 8000, "max": 15000, "description": "Customer payment portal with secure payment processing", "color": "#c8b98c" },
    { "category": "CMS and Blog System", "min": 8000, "max": 12000, "description": "Content management system with integrated blog", "color": "#8fa05e" },
    { "category": "Admin and Customer Portals", "min": 12000, "max": 20000, "description": "Admin dashboard with full management + customer self-service", "color": "#5c6f3a" },
    { "category": "Design and Branding", "min": 10000, "max": 18000, "description": "Mobile-responsive tactical design with brand identity", "color": "#334024" },
    { "category": "SEO and Analytics", "min": 7000, "max": 12000, "description": "Search optimization and analytics integration", "color": "#26301c" },
    { "category": "Testing and Deployment", "min": 5000, "max": 8000, "description": "Quality assurance testing and production deployment", "color": "#1c2415" }
  ],
  "timeline": { "agency": "3-6 months" },
  "factors": [
    "Senior developer rates: $150-$250 per hour",
    "Project management overhead included",
    "Multiple specialists: frontend, backend, integrations, design",
    "Agency profit margin: 50-100% standard markup",
    "Ongoing maintenance and support contracts typical",
    "Hosting, SSL, and infrastructure costs additional"
  ]
}`
