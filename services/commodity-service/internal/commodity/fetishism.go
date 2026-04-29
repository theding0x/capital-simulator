package commodity

// SocialRelations is the simulation's answer to §4 - "The Fetishism of
// Commodities and the Secret Thereof."
//
// In bourgeois daily experience, exchange-value appears as a property of
// things ("this coat is worth 20 yards of linen"). The relation between the
// labours of the tailor and the weaver is hidden behind the relation between
// their products. Marx calls this displacement the fetish-character of the
// commodity-form: "a definite social relation between men ... assumes, in
// their eyes, the fantastic form of a relation between things."
//
// SocialRelations is the un-mystified view. Given a Commodity it surfaces:
//
//   - the abstract labour-time congealed in one unit (the value);
//   - the kind of concrete labour that produced it;
//   - the equivalent labour-time in any other commodity it exchanges for.
//
// The two views (exchange-form vs. social-relations form) compute the same
// numbers - only the framing differs. That is exactly Marx's point: the
// fetish is not a calculation error, it is a social form. We treat exposing
// it as a first-class API affordance rather than a footnote.
type SocialRelations struct {
	// Subject is the commodity being analyzed.
	Subject Commodity `json:"subject"`

	// LabourPerUnit is how much abstract human labour is congealed in one
	// unit of Subject (== Subject.SNLTPerUnit, but named to make the social
	// reading explicit).
	LabourPerUnit LabourMinutes `json:"labour_per_unit"`

	// ConcreteLabour names the kind of useful labour that produced Subject.
	// The use-value side of the dual character of the labour.
	ConcreteLabour ConcreteLabour `json:"concrete_labour"`

	// LabourRelations expresses the labour-time equivalence with each other
	// commodity in the population. "x of A is the product of T minutes of
	// abstract human labour - so is y of B."
	LabourRelations []LabourRelation `json:"labour_relations"`

	// Note is a short editorial line surfacing the fetishism critique to the
	// API consumer. It is intentionally part of the response so the political
	// economy of the data shape is not erased.
	Note string `json:"note"`
}

// LabourRelation expresses the equivalence between one unit of Subject and
// some quantity of Counterpart in terms of abstract human labour-time.
type LabourRelation struct {
	Counterpart   Commodity     `json:"counterpart"`
	CounterpartQty Quantity     `json:"counterpart_qty"`
	LabourTime    LabourMinutes `json:"labour_time"`
}

const fetishismNote = "" +
	"Exchange-value appears as a property of things; in fact it expresses a " +
	"relation between human labours. The quantities below are the " +
	"abstract-human-labour times congealed in equal-value bundles of these " +
	"commodities. (Capital I, Ch. 1, §4 - The Fetishism of Commodities and " +
	"the Secret Thereof.)"

// SocialRelationsOf computes the un-mystified view of subject against the
// rest of the population.
func SocialRelationsOf(subject Commodity, population []Commodity) SocialRelations {
	out := SocialRelations{
		Subject:        subject,
		LabourPerUnit:  subject.SNLTPerUnit,
		ConcreteLabour: subject.ConcreteLabour,
		Note:           fetishismNote,
	}
	for _, c := range population {
		if c.ID == subject.ID || c.SNLTPerUnit <= 0 {
			continue
		}
		ratio, err := ExchangeRatioBetween(subject, c, 1)
		if err != nil {
			continue
		}
		out.LabourRelations = append(out.LabourRelations, LabourRelation{
			Counterpart:    c,
			CounterpartQty: ratio.QuoteQty,
			LabourTime:     ratio.CommonValue,
		})
	}
	return out
}
