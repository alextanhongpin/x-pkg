// Package stopwords allows checking if a given word is a stopwords, and
// customization to the stopwords list.
package stopwords

import (
	"strings"
)

var stopwords = `i me my myself we our ours ourselves you your yours yourself yourselves he him his himself she her hers herself it its itself they them their theirs themselves what which who whom this that these those am is are was were be been being have has had having do does did doing a an the and but if or because as until while of at by for with about against between into through during before after above below to from up down in out on off over under again further then once here there when where why how all any both each few more most other some such no nor not only own same so than too very s t can will just don should now you're you've you'll you'd she's it's that'll don't should've d ll m o re ve y ain aren aren't couldn couldn't didn didn't doesn doesn't hadn hadn't hasn hasn't haven haven't isn isn't ma mightn mightn't mustn mustn't needn needn't shan shan't shouldn shouldn't wasn wasn't weren weren't won won't wouldn wouldn't could he'd he'll he's here's how's i'd i'll i'm i've let's ought she'd she'll that's there's they'd they'll they're they've we'd we'll we're we've what's when's where's who's why's would able abst accordance according accordingly across act actually added adj affected affecting affects afterwards ah almost alone along already also although always among amongst announce another anybody anyhow anymore anyone anything anyway anyways anywhere apparently approximately arent arise around aside ask asking auth available away awfully b back became become becomes becoming beforehand begin beginning beginnings begins behind believe beside besides beyond biol brief briefly c ca came cannot can't cause causes certain certainly co com come comes contain containing contains couldnt date different done downwards due e ed edu effect eg eight eighty either else elsewhere end ending enough especially et etc even ever every everybody everyone everything everywhere ex except f far ff fifth first five fix followed following follows former formerly forth found four furthermore g gave get gets getting give given gives giving go goes gone got gotten h happens hardly hed hence hereafter hereby herein heres hereupon hes hi hid hither home howbeit however hundred id ie im immediate immediately importance important inc indeed index information instead invention inward itd it'll j k keep keeps kept kg km know known knows l largely last lately later latter latterly least less lest let lets like liked likely line little 'll look looking looks ltd made mainly make makes many may maybe mean means meantime meanwhile merely mg might million miss ml moreover mostly mr mrs much mug must n na name namely nay nd near nearly necessarily necessary need needs neither never nevertheless new next nine ninety nobody non none nonetheless noone normally nos noted nothing nowhere obtain obtained obviously often oh ok okay old omitted one ones onto ord others otherwise outside overall owing p page pages part particular particularly past per perhaps placed please plus poorly possible possibly potentially pp predominantly present previously primarily probably promptly proud provides put q que quickly quite qv r ran rather rd readily really recent recently ref refs regarding regardless regards related relatively research respectively resulted resulting results right run said saw say saying says sec section see seeing seem seemed seeming seems seen self selves sent seven several shall shed shes show showed shown showns shows significant significantly similar similarly since six slightly somebody somehow someone somethan something sometime sometimes somewhat somewhere soon sorry specifically specified specify specifying still stop strongly sub substantially successfully sufficiently suggest sup sure take taken taking tell tends th thank thanks thanx thats that've thence thereafter thereby thered therefore therein there'll thereof therere theres thereto thereupon there've theyd theyre think thou though thoughh thousand throug throughout thru thus til tip together took toward towards tried tries truly try trying ts twice two u un unfortunately unless unlike unlikely unto upon ups us use used useful usefully usefulness uses using usually v value various 've via viz vol vols vs w want wants wasnt way wed welcome went werent whatever what'll whats whence whenever whereafter whereas whereby wherein wheres whereupon wherever whether whim whither whod whoever whole who'll whomever whos whose widely willing wish within without wont words world wouldnt www x yes yet youd youre z zero a's ain't allow allows apart appear appreciate appropriate associated best better c'mon c's cant changes clearly concerning consequently consider considering corresponding course currently definitely described despite entirely exactly example going greetings hello help hopefully ignored inasmuch indicate indicated indicates inner insofar it'd novel presumably reasonably second secondly sensible serious seriously t's third thorough thoroughly three well wonder amoungst amount bill bottom call con cry de describe detail eleven empty fifteen fify fill find fire forty front full hasnt interest mill mine move side sincere sixty system ten thickv thin top twelve twenty A B C D E F G H I J K L M N O P Q R S T U V W X Y Z op research-articl pagecount cit ibid les le au est pas el los u201d well-b http volumtype par 0o 0s 3a 3b 3d 6b 6o a1 a2 a3 a4 ab ac ad ae af ag aj al ao ap ar av aw ax ay az b1 b2 b3 ba bc bd bi bj bk bl bn bp br bs bt bu bx c1 c2 c3 cc cd ce cf cg ch ci cj cl cm cn cp cq cr cs ct cu cv cx cy cz d2 da dc dd df di dj dk dl dp dr ds dt du dx dy e2 e3 ea ec ee ef ei ej em en eo ep eq er es eu ev ey f2 fa fc fi fj fl fn fo fr fs ft fu fy ga ge gi gj gl gr gs gy h2 h3 hh hj ho hr hs hu hy i2 i3 i4 i6 i7 i8 ia ib ic ig ih ii ij il io ip iq ir iv ix iy iz jj jr js jt ju ke kj ko l2 la lb lc lf lj ln lo lr ls lt m2 mn mo ms mt mu n2 nc ne ng ni nj nl nn nr ns nt ny oa ob oc od og oi oj ol om oo oq os ot ou ow ox oz p1 p2 p3 pc pd pe pf ph pi pj pk pl pm pn po pq pr ps pt pu py qj qu r2 ra rc rf rh ri rj rl rm rn ro rq rr rs rt ru rv ry s2 sa sc sd se sf si sj sl sm sn sp sq sr ss st sy sz t1 t2 t3 tb tc td te tf ti tj tl tm tn tp tq tr tt tv tx ue ui uj uk um uo ur ut va wa vd wi vj vo wo vq vt vu x1 x2 x3 xf xi xj xk xl xn xo xs xt xv xx y2 yj yl yr ys yt zi zz`

var stops Stopwords

func init() {
	stops = New()
}
func Has(word string) bool {
	return stops.Has(word)
}
func Add(words ...string) {
	stops.Add(words...)
}

type Stopwords map[string]bool

func New() Stopwords {
	m := make(map[string]bool)
	words := strings.Fields(stopwords)
	for _, word := range words {
		m[word] = true
	}
	return m
}

// Has checks if the word exists in the stopwords list.
func (s Stopwords) Has(word string) bool {
	word = strings.TrimSpace(word)
	word = strings.ToLower(word)
	return s[word]
}

// Add adds words to the stopwords list.
func (s Stopwords) Add(words ...string) {
	for _, word := range words {
		s[word] = true
	}
}
