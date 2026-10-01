package view

// pageFaqs, port of templates/pages/faqs.php. Fully static content.
const pageFaqs = `<section class="block">
  <div class="wrap">
    <div class="eyebrow">INTEL // FAQ</div>
    <h1 style="font-family:var(--display);color:var(--cream);font-size:clamp(2rem,6vw,3rem);margin:.4rem 0 .8rem">Questions, <em>answered.</em></h1>
    <p class="lead">Everything you need to know about safety, pricing, guarantees, and what to expect. Still curious? <a href="/contact" style="color:var(--orange)">Reach out</a>. We're happy to help.</p>
  </div>
</section>

<section class="block alt">
  <div class="wrap" style="max-width:840px">
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Safety &amp; Methods</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Are your treatments safe for kids and pets?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. We use low-toxicity products chosen to be tough on pests and easy on households, and we apply them with precision in cracks, crevices, and other spots kids and pets cannot reach, rather than broadcasting them across open surfaces. Baits go into tamper-resistant stations placed out of reach. Your technician will tell you about any brief re-entry wait for the specific treatment used and what to do if your pet shows interest in a treated area. If anyone in the home has allergies, asthma, or chemical sensitivities, mention it when you book so we can plan the treatment around it.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Are your products eco-friendly?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Absolutely. We work from an integrated approach: inspection, exclusion, and sanitation come first, and products are used only where they are needed, in the smallest effective amount. Applications are targeted at pest harborages and entry points instead of blanket spraying the whole property, which keeps pollinators, soil, and waterways out of the picture. It still gets the job done, because the product goes where the pests actually live. Ask your technician what is being applied and why, and you will get a straight answer.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Pricing &amp; Guarantees</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>How much does pest control cost?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">It depends on the pest, how far the infestation has spread, the size of your home, and whether you need a one-time knockdown or ongoing protection. A single wasp nest and a whole-house rodent exclusion are very different jobs, so we do not do one-size-fits-all pricing. Every quote is free, itemized, and given before any work starts, so you see exactly what is included. Call or request a quote online and you will get a firm number, not a range that moves later.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do you offer a guarantee?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. Every treatment is backed by our 90-day warranty and 100% satisfaction guarantee. If pests return between scheduled visits, we re-treat at no additional cost.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Are there any hidden fees?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Never. The price we quote is the price you pay, and we confirm it before work starts so there is nothing to argue about later. Quotes are free, and re-treatments between scheduled visits are covered, not billed as extras. No fuel surcharges, no surprise line items on the invoice. If a job ever needs something outside the original quote, we tell you first and get your OK before we do it.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Service &amp; Scheduling</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do you offer same-day service?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. When pests cannot wait, like a wasp nest by the front door or mice in the pantry, we offer same-day visits across Washington, Idaho, Oregon, and Arizona. Scheduling depends on technician availability and routing, so call as early in the day as you can to get a slot. If the day is fully booked, we will set the next available visit and tell you what to do in the meantime.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>What areas do you serve?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">We cover four states. In Washington: Spokane, Spokane Valley, Cheney, Liberty Lake, Airway Heights, Medical Lake, Deer Park, and Mead. In Idaho: Coeur d'Alene, Post Falls, Hayden, and Rathdrum. In Oregon: Hermiston and Milton-Freewater. In Arizona: Phoenix. See our service areas page for the full list, and if you are just outside these areas, call and ask.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do I need to prepare my home before treatment?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Usually very little. When you schedule, we give you prep instructions matched to your pest and treatment, for example clearing under-sink cabinets before an ant treatment or trimming grass back from the foundation before an exterior visit. Most interior work just needs clear access along baseboards and to the problem areas. Your technician walks you through anything left on arrival, so you are never guessing.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Getting Started</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>How do I get started?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Call us at (509) 818-0993 or request a free quote through the contact page. We will ask what you are seeing, where, and how long it has been going on, then recommend a plan and a price before any work happens. Many jobs can be scheduled the same day you call. There is no obligation and no pressure; if all you need is advice, we will say so.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do you handle commercial properties?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. We treat restaurants, offices, retail spaces, warehouses, and multi-unit housing as well as homes. Visits are scheduled around your operating hours so customers and staff are not disrupted. We also provide the service records that health inspections and audits ask for, so compliance stays simple.</p>
      </details>
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Ready when <em>you are.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15098180993">☎ (509) 818-0993</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>
`
