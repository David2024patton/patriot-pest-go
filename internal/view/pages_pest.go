package view

// PestCopy — category-specific copy for the pest "threat file" page.
// Verbatim port of the $byCat table in templates/pages/pest.php, so each
// threat file reads as tailored (insect / rodent / wildlife) rather than generic.
type PestCopy struct {
	Signs []string
	Treat string
	Prev  []string
}

var pestCopies = map[string]PestCopy{
	"insect": {
		Signs: []string{
			"Live insects or shed skins near entry points, kitchens, or baseboards",
			"Small droppings, smear marks, or nesting material in hidden corners",
			"Damaged food packaging, wood, or fabrics",
			"Unusual odors or faint rustling sounds at night",
		},
		Treat: "Targeted crack-and-crevice and baiting programs that reach the colony at its source, not just the insects you see. We use low-toxicity, family- and pet-safe products, then monitor to confirm elimination.",
		Prev: []string{
			"Seal cracks around the foundation, windows, and utility entry points",
			"Store food in airtight containers and clean up crumbs and spills promptly",
			"Eliminate standing water and fix leaky fixtures",
			"Keep vegetation and mulch pulled back from the foundation",
		},
	},
	"rodent": {
		Signs: []string{
			"Droppings along walls, in cabinets, or near food sources",
			"Gnaw marks on wiring, wood, or food packaging",
			"Grease rub marks along baseboards and entry routes",
			"Scratching or scurrying noises in walls or attics at night",
		},
		Treat: "A complete exclusion-plus-removal program: we trap and remove active rodents, then seal entry points larger than a quarter-inch so new ones can’t get in. Attic and crawlspace decontamination available.",
		Prev: []string{
			"Seal gaps around pipes, vents, and the foundation with steel wool and caulk",
			"Store food and pet food in sealed, rodent-proof containers",
			"Keep garbage in tightly sealed bins and remove clutter",
			"Trim tree branches and vegetation away from the roofline",
		},
	},
	"wildlife": {
		Signs: []string{
			"Noises in the attic, walls, or chimney at dawn and dusk",
			"Entry holes, torn vents, or damaged soffits and fascia",
			"Droppings or nesting material in attics and crawlspaces",
			"Damaged insulation, wiring, or stored items",
		},
		Treat: "Humane removal and exclusion. We safely remove the animals, then install one-way doors and seal entry points so they can’t return, with cleanup and decontamination of affected areas.",
		Prev: []string{
			"Cap chimneys and vent openings with wildlife-proof covers",
			"Repair damaged soffits, fascia, and roof vents",
			"Keep trash secured and remove outdoor food sources",
			"Trim overhanging branches that provide roof access",
		},
	},
}

// PestCopyFor returns the copy block for a pest category, falling back to insect.
func PestCopyFor(cat string) PestCopy {
	if c, ok := pestCopies[cat]; ok {
		return c
	}
	return pestCopies["insect"]
}

// pagePest — templates/pages/pest.php. The unified threat-file page for one pest.
// Data: .Pest (data.Pest), .Related ([]data.Pest).
const pagePest = `<!-- ===== THREAT FILE HEADER ===== -->
<section class="block">
  <div class="wrap">
    <div class="eyebrow">THREAT FILE #{{p3 (add .Pest.SortOrder 1)}} // {{upper .Pest.Category}}</div>
    <div class="split" style="margin-top:1.2rem;align-items:start">
      <div>
        <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3.2rem);line-height:1.05">{{.Pest.Name}} <span style="color:var(--orange)">Control</span></h1>
        {{if .Pest.ScientificName}}<div class="sci" style="font-family:var(--mono);font-style:italic;color:var(--khaki);margin:.4rem 0 1rem">{{.Pest.ScientificName}}</div>{{end}}
        <p class="lead">{{.Pest.Description}}</p>
        <div class="meter" style="max-width:420px;margin:1.4rem 0">
          <div class="label"><span>Regional Threat Level</span><span>{{.Pest.ThreatLevel}}%</span></div>
          <div class="bar"><div class="fill" data-lvl="{{.Pest.ThreatLevel}}"></div></div>
        </div>
        <div class="hero-ctas" style="margin-top:1.6rem">
          <a class="btn btn-primary" href="tel:+15098180993">☎ Get a Free Quote</a>
          <a class="btn btn-ghost" href="/contact">Book Service ▸</a>
        </div>
      </div>
      <div>
        <span class="pphoto s-lg"><img src="{{pestimg .Pest.Filename}}" alt="{{.Pest.Name}}" loading="eager"><i class="ret" aria-hidden="true"></i></span>
      </div>
    </div>
  </div>
</section>

<!-- ===== SIGNS / TREATMENT / PREVENTION ===== -->
<section class="block alt">
  <div class="wrap">
    <div class="grid g3">
      {{with pestcopy .Pest.Category}}
      <div class="card">
        <h3 style="font-family:var(--display);color:var(--cream)">⚠ Signs of Activity</h3>
        <ul style="margin:.6rem 0 0 1.1rem;color:var(--khaki);line-height:1.7">
          {{range .Signs}}<li>{{.}}</li>{{end}}
        </ul>
      </div>
      <div class="card">
        <h3 style="font-family:var(--display);color:var(--cream)">🎯 Our Treatment</h3>
        <p style="color:var(--khaki);line-height:1.7;margin-top:.6rem">{{.Treat}}</p>
      </div>
      <div class="card">
        <h3 style="font-family:var(--display);color:var(--cream)">🛡 Prevention</h3>
        <ul style="margin:.6rem 0 0 1.1rem;color:var(--khaki);line-height:1.7">
          {{range .Prev}}<li>{{.}}</li>{{end}}
        </ul>
      </div>
      {{end}}
    </div>
    <div class="promise" style="margin-top:1.6rem"><b>🛡️ 90-Day Warranty:</b> if {{lower .Pest.Name}} return between scheduled visits, we re-treat at no additional cost. Licensed, bonded, and insured across Washington, Idaho, Oregon &amp; Arizona.</div>
  </div>
</section>

<!-- ===== RELATED THREATS ===== -->
{{if .Related}}
<section class="block">
  <div class="wrap">
    <div class="eyebrow">RELATED THREAT FILES</div>
    <h2 style="font-family:var(--display);color:var(--cream);margin:.4rem 0 1.4rem">Also operating <em>in your region.</em></h2>
    <div class="grid g3">
      {{range $r := .Related}}
      <a class="card" href="/pest/{{$r.Slug}}" style="text-decoration:none;color:inherit">
        <span class="pphoto"><img src="{{pestimg $r.Filename}}" alt="{{$r.Name}}" loading="lazy"><i class="ret" aria-hidden="true"></i></span>
        <h3 style="font-family:var(--display);color:var(--cream);margin-top:.7rem">{{$r.Name}} ▸</h3>
      </a>
      {{end}}
    </div>
  </div>
</section>
{{end}}

<!-- ===== CTA ===== -->
<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <div class="eyebrow">FINAL ORDERS</div>
    <h2 style="font-family:var(--display);color:var(--cream);margin:.4rem 0 1rem">Seeing {{lower .Pest.Name}}? <em>Let's end it.</em></h2>
    <p class="lead">Same-day service available. Free quotes, transparent pricing, 90-day warranty.</p>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.4rem">
      <a class="btn btn-primary" href="tel:+15098180993">☎ (509) 818-0993 <small>WA, ID, OR</small></a>
      <a class="btn btn-ghost" href="tel:+16027558414">☎ (602) 755-8414 <small>ARIZONA</small></a>
    </div>
  </div>
</section>`
