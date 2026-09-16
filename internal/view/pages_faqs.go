package view

// pageFaqs — templates/pages/faqs.php. Fully static content.
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
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. We use eco-friendly, low-toxicity products and apply them with precision. We'll advise you on any brief re-entry timing for specific treatments, but most are family- and pet-safe.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Are your products eco-friendly?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Absolutely. We prioritize low-toxicity, targeted products that protect your household and the environment while still eliminating pests effectively.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Pricing &amp; Guarantees</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>How much does pest control cost?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Cost depends on your home's size, the pest, and severity. We offer transparent pricing with no hidden fees. Get an accurate, no-obligation quote in minutes by phone or online.</p>
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
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Never. The price we quote is the price you pay. Free quotes, free re-treatments between visits, zero surprises.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Service &amp; Scheduling</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do you offer same-day service?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes, same-day service is available across all four states when pests can't wait. Call the line and we'll get you scheduled fast.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>What areas do you serve?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">We serve Washington, Idaho, Oregon, and Arizona, including Spokane, Coeur d'Alene, Hermiston, Phoenix, and surrounding communities. See our service areas for the full list.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do I need to prepare my home before treatment?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">We'll give you simple, specific prep instructions when we schedule. Most treatments require minimal preparation, and your technician will walk you through everything.</p>
      </details>
    </div>
    <div class="faq-cat" style="margin-bottom:2rem">
      <h2 style="font-family:var(--display);color:var(--orange);font-size:1.15rem;margin-bottom:1rem">Getting Started</h2>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>How do I get started?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Call us or request a free quote online. We'll assess the situation, recommend a plan, and schedule service, often same-day.</p>
      </details>
      <details class="faq-item" style="border:1px solid var(--olive-700);background:var(--olive-800);margin-bottom:.7rem;padding:1rem 1.2rem">
        <summary class="faq-q" style="cursor:pointer;color:var(--cream);font-weight:600;list-style:none;display:flex;justify-content:space-between;gap:1rem">
          <span>Do you handle commercial properties?</span><span style="color:var(--orange)">+</span>
        </summary>
        <p class="faq-a" style="color:var(--khaki);line-height:1.7;margin-top:.8rem">Yes. We serve both homes and businesses, including restaurants and offices, with discreet, reliable, compliance-friendly service.</p>
      </details>
    </div>
  </div>
</section>

<section class="block cta-band">
  <div class="wrap" style="text-align:center">
    <h2 style="font-family:var(--display);color:var(--cream)">Ready when <em>you are.</em></h2>
    <div class="hero-ctas" style="justify-content:center;margin-top:1.2rem">
      <a class="btn btn-primary" href="tel:+15094715767">☎ (509) 471-5767</a>
      <a class="btn btn-ghost" href="/contact">Get a Free Quote ▸</a>
    </div>
  </div>
</section>
`
