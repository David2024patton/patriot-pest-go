// Package view, area_content.go: researched local copy for city pages.
// Every factual claim below was verified against a named source during
// research on 2026-10-01. When only region-level sources existed, the pest
// note says so instead of inventing a city-specific claim.
package view

func init() {
	RegisterAreaCopy("airway-heights", AreaCopy{
		Intro: "Airway Heights sits on the windswept plain west of Spokane, and its identity is tied to the runway. " +
			"Fairchild Air Force Base, home of the 92nd and 141st Air Refueling Wings, anchors the town, with housing " +
			"that ranges from midcentury ranch homes near the base to newer subdivisions along US-2. " +
			"US-2 runs through town as a five-lane thoroughfare, carrying commuters between Spokane and Fairchild. " +
			"The Airway Heights Corrections Center sits on the west side of town, near Sunset Park.",
		LocalNote: "Sunset Park is the town's central green space, and the surrounding neighborhoods are where " +
			"Airway Heights shows its suburban side, with single family homes on tree lined streets. " +
			"Growth has pushed new retail and housing out along US-2, keeping the town busy between the base " +
			"and Spokane.",
		PestNote: "Local pest pros list ants and spiders around Sunset Park, and rodents near the corrections " +
			"center, matching the standard Inland Northwest pattern. " +
			"Hobo spiders are well established across the Inland Northwest, and fall brings rodents looking for " +
			"warmth in garages and crawl spaces.",
	})

	RegisterAreaCopy("cheney", AreaCopy{
		Intro: "Cheney is a college town through and through, home to Eastern Washington University, which was " +
			"founded in 1882. The campus is compact and roughly triangular, tucked near downtown with the red " +
			"turf of Roos Field visible from the surrounding streets. Turnbull National Wildlife Refuge sits " +
			"nearby, pulling waterfowl, birders, and trail traffic through town.",
		LocalNote: "Downtown Cheney clusters around First Street near the old rail depot, with small shops and " +
			"restaurants that serve the student population. Off-campus rental neighborhoods ring the university, " +
			"and family housing spreads south toward the rural edges of town.",
		PestNote: "Local pest sources name mice and rats, cockroaches, ants, spiders, wasps, and earwigs as " +
			"Cheney's regular offenders. Paper wasps and other stinging insects keep wasp control in steady " +
			"demand through the warm months.",
	})

	RegisterAreaCopy("coeur-d-alene", AreaCopy{
		Intro: "Coeur d'Alene wraps around its namesake lake, with downtown and the resort core fronting the " +
			"waterfront at the base of Tubbs Hill, a natural area that juts out into the lake. The Garden " +
			"District is the oldest neighborhood in town, lined with century old Craftsman homes and bungalows. " +
			"The Canfield area stretches to the north with newer subdivisions, and Fernan Lake sits to the " +
			"east. I-90 carries the retail and service corridor across the north side.",
		LocalNote: "Lake life defines the east side, while Government Way and Sherman Avenue split the " +
			"commercial core from the residential streets. The floating boardwalk, one of the longest in " +
			"the world, runs along the resort waterfront and draws summer crowds.",
		PestNote: "Local pest pros describe a clear seasonal cycle: ants in spring, wasps in summer, spiders " +
			"including black widows in fall, and mice seeking shelter in winter. Carpenter ants are a common " +
			"structure pest across Idaho, and they show up in Coeur d'Alene homes built into the area's " +
			"timber.",
	})

	RegisterAreaCopy("deer-park", AreaCopy{
		Intro: "Deer Park is Patriot's home base, about 15 miles north of Spokane in north Spokane County. " +
			"The town got its name when railroad surveyors spotted deer grazing here, and sawmills built the " +
			"early economy, with eight mills once operating within 10 miles of town. Arcadia Apple Orchards, " +
			"planted around 1906, put the area on the agricultural map. The annual Settlers Days fair still " +
			"celebrates that pioneer heritage.",
		LocalNote: "Mix Park anchors the town center, and the Deer Park Golf Club subdivision spreads out " +
			"around the course on the east side. Newer subdivisions are pushing outward into the timbered " +
			"land that separates Deer Park from the Spokane metro.",
		PestNote: "The standard Inland Northwest mix applies: hobo spiders, ants, fall rodents moving " +
			"indoors, and yellowjackets in late summer. Timbered lots and wooded property lines keep " +
			"spider and rodent pressure above what you see in fully built out suburbs.",
	})

	RegisterAreaCopy("hayden", AreaCopy{
		Intro: "Hayden sits directly north of Coeur d'Alene, a fast growing suburb built around Hayden Lake. " +
			"Honeysuckle Beach is the town's lakefront draw, and the Canfield Mountain trail system, with " +
			"more than 32 miles of trails, climbs the east side of town. McIntire Family Park serves the " +
			"south end, and Avondale Golf Club hosts championship play. Silverwood Theme Park is a short " +
			"drive north in Athol.",
		LocalNote: "Newer subdivisions with planned streets spread across the prairie west of US-95, while " +
			"older lakefront properties hug the Hayden Lake shoreline. The town has kept a quieter, " +
			"residential feel even as growth fills in the corridors.",
		PestNote: "Local pest sources list wasps, bees, flies, mosquitoes, ants, spiders, rodents, and bed " +
			"bugs as Hayden's regulars. Boxelder bugs and stink bugs also show up around Hayden businesses " +
			"and homes in the fall.",
	})

	RegisterAreaCopy("hermiston", AreaCopy{
		Intro: "Hermiston calls itself the Watermelon Capital of the World, and it is the largest city in " +
			"Eastern Oregon. The town sits at the junction of I-84 and I-82, with the Umatilla River running " +
			"through Riverfront Park and Hermiston Butte rising on the skyline. The area has become a " +
			"data center hub, with Amazon AWS facilities among the region's largest employers. The Hermiston " +
			"Farmers Market runs through the growing season downtown.",
		LocalNote: "Older neighborhoods cluster near downtown and the river, while newer subdivisions spread " +
			"south and west toward the highway corridors. Irrigated farmland wraps around the city limits, " +
			"keeping the town tied to the surrounding Columbia Basin agriculture.",
		PestNote: "Oregon region sources list ants, spiders, cockroaches, rodents, earwigs, and yellowjackets " +
			"as the common Hermiston-area pests. Agricultural land and irrigation nearby raise fly and ant " +
			"pressure above what denser cities see.",
	})

	RegisterAreaCopy("liberty-lake", AreaCopy{
		Intro: "Liberty Lake, incorporated in 2001, sits on the Washington Idaho border between Spokane and " +
			"Coeur d'Alene. The MeadowWood golf community anchors the east side with three courses: " +
			"MeadowWood, Liberty Lake Golf Course, and Trailhead. Liberty Lake Regional Park covers about " +
			"3,000 acres of wetlands and trails on the lake's east end. A century ago the town was Spokane's " +
			"Inland Seashore, with lake resorts and a dance pavilion drawing weekend crowds.",
		LocalNote: "The lake draws waterfront homes on the south shore, and the MeadowWood area holds the " +
			"town's higher end planned neighborhoods. Commercial growth clusters along Appleway Avenue, " +
			"serving both Liberty Lake and the Spokane Valley spillover.",
		PestNote: "Local pest sources note that Liberty Lake's dry summers and cold winters shape pest " +
			"timing: roaches, rats, mice, ants, and spiders peak on a clear seasonal cycle. Fall rodent " +
			"intrusions are the top cold weather call.",
	})

	RegisterAreaCopy("mead", AreaCopy{
		Intro: "Mead is an unincorporated community in north Spokane County, spread along US-2 and Highway " +
			"395 near Farewell Road. Mount Spokane State Park sits about 15 miles to the northeast, and " +
			"Wandermere Golf Course anchors the older part of town. The Haynes Estate and Feryn conservation " +
			"areas protect open land around the community. The top rated Mead School District draws " +
			"families to the area's newer subdivisions.",
		LocalNote: "The community mixes working ranches and rural acreage with suburban subdivisions along " +
			"the highway corridors. Larger lots and open ground mean Mead properties border wildlife " +
			"corridors more than most Spokane suburbs.",
		PestNote: "Region-level Spokane sources list ants, spiders, roaches, mosquitoes, rodents, and wasps " +
			"as the area mix, and the same pattern applies in Mead. Larger rural lots here add rodent and " +
			"wildlife pressure, especially as temperatures drop in fall.",
	})

	RegisterAreaCopy("medical-lake", AreaCopy{
		Intro: "Medical Lake takes its name from the mineral content of the lake itself, called Lac de " +
			"Medicine after Andrew Lefevre noted its properties in 1872. Eastern State Hospital sits across " +
			"the lake on the east shore, and Fairchild Air Force Base presses up against the town's north " +
			"side. Lefevre Street runs the main street, lined with midcentury ranch homes. The Channeled " +
			"Scablands shape the rolling terrain around town. The Bluegrass Festival and Founders Day are " +
			"the town's signature events.",
		LocalNote: "The town keeps a small town lake community feel, with older homes near the water and " +
			"newer development spreading toward the base. The lake and its surrounding parks draw summer " +
			"swimming and fishing traffic.",
		PestNote: "Region-level Washington sources describe the standard Spokane area mix for Medical Lake: " +
			"ants, spiders, rodents, and wasps. Lakefront properties and the nearby base landscape add " +
			"mosquito and fly pressure in summer months.",
	})

	RegisterAreaCopy("milton-freewater", AreaCopy{
		Intro: "Milton-Freewater sits on the Oregon side of the Walla Walla Valley, along OR Highway 11 " +
			"between Walla Walla and Pendleton. The town once carried the title Pea Capital of the World, " +
			"and orchards and wine country still surround it. The Muddy Frogwater Festival is the town's " +
			"signature event, complete with frog imagery everywhere. Blue Mountain Cider Company pours in " +
			"town. The city formed in 1951 when Milton and Freewater merged.",
		LocalNote: "Older homes cluster in the historic cores of the two original towns, while orchard " +
			"country wraps around the city limits. The Blue Mountains rise to the east, marking the edge " +
			"of town and the start of recreation country.",
		PestNote: "Local Oregon pest sources say termite control is the top inquiry in Milton-Freewater, " +
			"with spiders, rodents, ants, and cockroaches close behind. Boxelder bugs and spiders are noted " +
			"as regional regulars across the Walla Walla Valley.",
	})

	RegisterAreaCopy("phoenix", AreaCopy{
		Intro: "Phoenix is divided into 15 urban villages, each with its own character. Arcadia, at the " +
			"south foot of Camelback Mountain, grew out of citrus groves and still shows the ranch homes " +
			"and irrigated landscaping of that era. Ahwatukee borders South Mountain on the city's south " +
			"edge, and South Mountain Park itself is one of the largest city parks in the United States. " +
			"Camelback East Village sits between Camelback Mountain and Piestewa Peak.",
		LocalNote: "Stucco and block wall construction defines most Valley neighborhoods, and citrus trees " +
			"still dot yards across Arcadia as a reminder of the old groves. New development keeps pushing " +
			"the city outward into the desert on every side.",
		PestNote: "Arizona bark scorpions are the most venomous scorpion in the United States and are " +
			"common around low desert buildings, hiding in hollow block walls and climbing textured stucco. " +
			"Desert subterranean termites build mud tubes up stem walls to reach framing, and American and " +
			"Turkestan cockroaches come out of sewer lines and irrigation boxes. Monsoon season adds " +
			"mosquitoes and flies, while pack rats pile middens against walls along the desert preserves.",
	})

	RegisterAreaCopy("post-falls", AreaCopy{
		Intro: "Post Falls is named for Frederick Post, the German immigrant who built a lumber mill on the " +
			"Spokane River in 1871, on land purchased from Chief Andrew Seltice of the Coeur d'Alene Tribe. " +
			"The city sits on the Rathdrum Prairie, bounded by the Spokane River to the south, about 20 " +
			"miles east of Spokane. Falls Park overlooks the namesake dam and falls, and Q'emiln Park, " +
			"pronounced ka-mee-lin, offers a 78.5 acre swimming beach, boat launch, and rock climbing. " +
			"The town is often called the Industrial Heart of North Idaho for its manufacturing and " +
			"technology employers.",
		LocalNote: "Milltown and Post Falls City Center hold the older craftsman and cottage homes near the " +
			"river, with waterfront properties running up to the million dollar range. Newer subdivisions " +
			"spread north and west, and the 23 mile North Idaho Centennial Trail threads through town from " +
			"the state line to Coeur d'Alene.",
		PestNote: "The owner of the Post Falls based Pointe Pest Control says most local calls are about " +
			"ants, followed closely by spiders, stinging insects, and rodents. Rodents like mice and rats, " +
			"ants, spiders, and occasional invaders such as stink bugs and centipedes are the common " +
			"commercial and residential offenders here.",
	})

	RegisterAreaCopy("rathdrum", AreaCopy{
		Intro: "Rathdrum sits at the base of Rathdrum Mountain, part of the Selkirk Range, on the northern " +
			"edge of the Rathdrum Prairie. The town started as Westwood and was renamed in 1881 after a " +
			"place in County Wicklow, Ireland, at the suggestion of a local businessman. The railroad built " +
			"the town, and the BNSF line still runs through it, with Amtrak's Empire Builder passing " +
			"without stopping. The historic main street survived the growth waves, and the town now serves " +
			"as a bedroom community for Coeur d'Alene and Post Falls.",
		LocalNote: "Rathdrum Mountain overlooks the city with communication towers and hiking access, and " +
			"Rathdrum City Park anchors the center of town. The paved Prairie Trail runs through town and " +
			"connects into the broader North Idaho trail network. New subdivisions keep replacing old " +
			"farmland as Post Falls pushes north toward Rathdrum.",
		PestNote: "Rodents like mice and rats, ants, spiders, and occasional invaders such as beetles and " +
			"earwigs are the common Rathdrum pests. Idaho sources add carpenter ants and odorous house ants, " +
			"black widows, hobo spiders, and wasps and bees nesting around homes as regional regulars.",
	})

	RegisterAreaCopy("spokane", AreaCopy{
		Intro: "Spokane revolves around its river, which tumbles through downtown in the Spokane Falls. " +
			"Riverfront Park, the 100 acre site of Expo '74, sits just north of the downtown core, and " +
			"River Park Square anchors the shopping and entertainment district. The South Hill rises south " +
			"of downtown with neighborhoods like South Perry, home to a weekly farmers market. The " +
			"University District on the east side holds Gonzaga University and branch campuses of " +
			"Washington State and Eastern Washington Universities.",
		LocalNote: "Much of downtown was rebuilt after the Great Fire of 1889 in Romanesque Revival style " +
			"by architect Kirtland Kelsey Cutter, and that historic character still defines Riverside " +
			"Avenue. The Logan neighborhood near Gonzaga mixes student housing with some of the city's " +
			"oldest homes. North Spokane's Hillyard and Five Mile areas hold the suburban north side.",
		PestNote: "Spokane sits in one of the few US regions where hobo spiders are well established, " +
			"documented here since the late 1960s. Odorous house ants, pavement ants, and carpenter ants " +
			"are the top ant threats, with rodents moving indoors each fall. Yellow jackets and bald " +
			"faced hornets nest around eaves and sheds all summer, and moles, bed bugs, and the velvety " +
			"tree ant round out the area's pest list.",
	})

	RegisterAreaCopy("spokane-valley", AreaCopy{
		Intro: "Spokane Valley stretches east of Spokane along the Spokane River to the Idaho state line. " +
			"The Spokane River Centennial Trail runs through the valley, past the Spokane Valley Mall and " +
			"along the river corridor. Dishman Hills Natural Area preserves over 500 acres of hiking land " +
			"right inside the urban area, and the Saltese Uplands Conservation Area covers about 700 acres " +
			"to the east. Greenacres anchors the east end, with the 37 mile Centennial Trail following the " +
			"river from Nine Mile Falls to the state line.",
		LocalNote: "The valley's economy centers on retail, healthcare, and manufacturing, with the " +
			"Spokane Business and Industrial Park as a major employer. Schools fall under the Central " +
			"Valley and East Valley districts. Mirabeau Park and the riverfront trailheads serve as the " +
			"valley's outdoor hubs.",
		PestNote: "The valley shares Spokane's pest profile: ants, spiders including hobo spiders, " +
			"fall rodents, and yellowjackets. The velvety tree ant, which nests in oak, pine, alder, and " +
			"elm trees, is specifically noted as common in the Spokane Valley and is often mistaken for " +
			"a carpenter ant.",
	})
}
