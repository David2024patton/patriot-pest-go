// Page body templates — Go ports of templates/pages/*.php. Each is registered
// under the name the layout's {{template .Page .}} expects, so view.render can
// run any page through the shared shell with identical markup.
package view

func init() {
	RegisterPageTemplate("home", pageHome)
	RegisterPageTemplate("about", pageAbout)
	RegisterPageTemplate("services", pageServices)
	RegisterPageTemplate("prices", pagePrices)
	RegisterPageTemplate("faqs", pageFaqs)
	RegisterPageTemplate("contact", pageContact)
	RegisterPageTemplate("areas", pageAreas)
	RegisterPageTemplate("area-detail", pageAreaDetail)
	RegisterPageTemplate("pest", pagePest)
	RegisterPageTemplate("blog-index", pageBlogIndex)
	RegisterPageTemplate("blog-post", pageBlogPost)
	RegisterPageTemplate("help", pageHelp)
	RegisterPageTemplate("links", pageLinks)
	RegisterPageTemplate("search", pageSearch)
	RegisterPageTemplate("socials", pageSocials)
	RegisterPageTemplate("referral", pageReferral)
	RegisterPageTemplate("sitemap", pageSitemap)
	RegisterPageTemplate("legal", pageLegal)
}
