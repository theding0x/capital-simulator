import type { ComponentType } from "react";
import { CHAPTERS } from "../chapters/registry";
import type { Commodity, Owner } from "../types";
import { Ch01Commodity } from "../chapters/vol1/Ch01Commodity";
import { Ch02Exchange } from "../chapters/vol1/Ch02Exchange";
import { Ch03Money } from "../chapters/vol1/Ch03Money";
import { Ch04Capital } from "../chapters/vol1/Ch04Capital";
import { Ch05Contradictions } from "../chapters/vol1/Ch05Contradictions";
import { Ch06LabourPower } from "../chapters/vol1/Ch06LabourPower";
import { Ch07LabourProcess } from "../chapters/vol1/Ch07LabourProcess";
import { Ch08ConstantVariableCapital } from "../chapters/vol1/Ch08ConstantVariableCapital";
import { Ch09RateOfSurplusValue } from "../chapters/vol1/Ch09RateOfSurplusValue";
import { Ch10WorkingDay } from "../chapters/vol1/Ch10WorkingDay";
import { Ch11RateAndMassOfSurplusValue } from "../chapters/vol1/Ch11RateAndMassOfSurplusValue";
import { Ch12RelativeSurplusValue } from "../chapters/vol1/Ch12RelativeSurplusValue";
import { Ch13Cooperation } from "../chapters/vol1/Ch13Cooperation";
import { Ch14Manufacture } from "../chapters/vol1/Ch14Manufacture";
import { Ch15Machinery } from "../chapters/vol1/Ch15Machinery";
import { Ch16AbsoluteAndRelative } from "../chapters/vol1/Ch16AbsoluteAndRelative";
import { Ch17MagnitudeChanges } from "../chapters/vol1/Ch17MagnitudeChanges";
import { Ch18RatesOfSurplusValue } from "../chapters/vol1/Ch18RatesOfSurplusValue";
import { Ch19WageForm } from "../chapters/vol1/Ch19WageForm";
import { Ch20TimeWages } from "../chapters/vol1/Ch20TimeWages";
import { Ch21PieceWages } from "../chapters/vol1/Ch21PieceWages";
import { Ch22NationalWages } from "../chapters/vol1/Ch22NationalWages";
import { Ch23SimpleReproduction } from "../chapters/vol1/Ch23SimpleReproduction";
import { Ch24AccumulationOfCapital } from "../chapters/vol1/Ch24AccumulationOfCapital";
import { Ch25GeneralLaw } from "../chapters/vol1/Ch25GeneralLaw";
import { Ch26PrimitiveAccumulation } from "../chapters/vol1/Ch26PrimitiveAccumulation";
import { Ch27EnclosureEvents } from "../chapters/vol1/Ch27EnclosureEvents";
import { Ch28BloodyLegislation } from "../chapters/vol1/Ch28BloodyLegislation";
import { Ch29GenesisFarmer } from "../chapters/vol1/Ch29GenesisFarmer";
import { Ch30HomeMarket } from "../chapters/vol1/Ch30HomeMarket";
import { Ch31IndustrialCapitalist } from "../chapters/vol1/Ch31IndustrialCapitalist";
import { Ch32HistoricalTendency } from "../chapters/vol1/Ch32HistoricalTendency";
import { Ch33ModernColonisation } from "../chapters/vol1/Ch33ModernColonisation";
import { Ch01CircuitMoneyCapital } from "../chapters/vol2/Ch01CircuitMoneyCapital";
import { Ch02CircuitProductiveCapital } from "../chapters/vol2/Ch02CircuitProductiveCapital";
import { Ch03CircuitCommodityCapital } from "../chapters/vol2/Ch03CircuitCommodityCapital";
import { Ch04ThreeFormulas } from "../chapters/vol2/Ch04ThreeFormulas";
import { Ch05TimeOfCirculation } from "../chapters/vol2/Ch05TimeOfCirculation";
import { Ch06CostsOfCirculation } from "../chapters/vol2/Ch06CostsOfCirculation";
import { Ch07TurnoverTime } from "../chapters/vol2/Ch07TurnoverTime";
import { Ch08FixedAndCirculatingCapital } from "../chapters/vol2/Ch08FixedAndCirculatingCapital";
import { Ch09AggregateTurnover } from "../chapters/vol2/Ch09AggregateTurnover";
import { Ch10TheoriesPhysiocratsSmith } from "../chapters/vol2/Ch10TheoriesPhysiocratsSmith";
import { Ch11TheoriesRicardo } from "../chapters/vol2/Ch11TheoriesRicardo";
import { Ch12WorkingPeriod } from "../chapters/vol2/Ch12WorkingPeriod";
import { Ch13TimeOfProduction } from "../chapters/vol2/Ch13TimeOfProduction";
import Ch14TimeOfCirculation from "../chapters/vol2/Ch14TimeOfCirculation";
import { Ch15PriceChanges } from "../chapters/vol2/Ch15PriceChanges";
import { Ch16TurnoverOfVariableCapital } from "../chapters/vol2/Ch16TurnoverOfVariableCapital";
import { Ch17CirculationOfSurplus } from "../chapters/vol2/Ch17CirculationOfSurplus";
import { Ch18RoleOfMoneyCapital } from "../chapters/vol2/Ch18RoleOfMoneyCapital";
import { Ch19FormerPresentations } from "../chapters/vol2/Ch19FormerPresentations";
import { Ch20SimpleReproduction } from "../chapters/vol2/Ch20SimpleReproduction";
import { Ch21ExtendedReproduction } from "../chapters/vol2/Ch21ExtendedReproduction";
import { Ch01CostPriceAndProfit } from "../chapters/vol3/Ch01CostPriceAndProfit";
import { Ch02RateOfProfit } from "../chapters/vol3/Ch02RateOfProfit";
import { Ch03ProfitRateRelation } from "../chapters/vol3/Ch03ProfitRateRelation";
import { Ch04TurnoverEffect } from "../chapters/vol3/Ch04TurnoverEffect";
import { Ch05ConstantCapitalEconomy } from "../chapters/vol3/Ch05ConstantCapitalEconomy";
import { Ch06PriceFluctuation } from "../chapters/vol3/Ch06PriceFluctuation";
import { Ch07SupplementaryRemarks } from "../chapters/vol3/Ch07SupplementaryRemarks";
import { Ch08CompositionDifferences } from "../chapters/vol3/Ch08CompositionDifferences";
import { Ch09GeneralRateOfProfit } from "../chapters/vol3/Ch09GeneralRateOfProfit";
import { Ch10Equalisation } from "../chapters/vol3/Ch10Equalisation";
import { Ch11WageFluctuations } from "../chapters/vol3/Ch11WageFluctuations";
import { Ch12SupplementaryRemarks } from "../chapters/vol3/Ch12SupplementaryRemarks";
import { Ch13FallingProfitLaw } from "../chapters/vol3/Ch13FallingProfitLaw";
import { Ch14CounteractingForces } from "../chapters/vol3/Ch14CounteractingForces";
import { Ch15InternalContradictions } from "../chapters/vol3/Ch15InternalContradictions";
import { Ch16CommercialCapital } from "../chapters/vol3/Ch16CommercialCapital";
import { Ch17CommercialProfit } from "../chapters/vol3/Ch17CommercialProfit";
import { Ch18MerchantTurnover } from "../chapters/vol3/Ch18MerchantTurnover";
import { Ch19MoneyDealingCapital } from "../chapters/vol3/Ch19MoneyDealingCapital";
import { Ch20HistoricalMerchantCapital } from "../chapters/vol3/Ch20HistoricalMerchantCapital";
import { Ch21InterestBearingCapital } from "../chapters/vol3/Ch21InterestBearingCapital";
import { Ch22RateOfInterest } from "../chapters/vol3/Ch22RateOfInterest";
import { Ch23ProfitOfEnterprise } from "../chapters/vol3/Ch23ProfitOfEnterprise";
import { Ch24CompoundInterest } from "../chapters/vol3/Ch24CompoundInterest";
import { Ch25CreditFictitiousCapital } from "../chapters/vol3/Ch25CreditFictitiousCapital";
import { Ch26MoneyCapitalAccumulation } from "../chapters/vol3/Ch26MoneyCapitalAccumulation";
import { Ch27CreditRole } from "../chapters/vol3/Ch27CreditRole";
import { Ch28MediumOfCirculation } from "../chapters/vol3/Ch28MediumOfCirculation";
import { Ch29BankCapital } from "../chapters/vol3/Ch29BankCapital";
import { Ch30MoneyRealCapital } from "../chapters/vol3/Ch30MoneyRealCapital";
import { Ch31MoneyRealCapitalII } from "../chapters/vol3/Ch31MoneyRealCapitalII";
import { Ch32MoneyRealCapitalIII } from "../chapters/vol3/Ch32MoneyRealCapitalIII";
import { Ch33CirculationCredit } from "../chapters/vol3/Ch33CirculationCredit";
import { Ch34CurrencyPrinciple } from "../chapters/vol3/Ch34CurrencyPrinciple";
import { Ch35PreciousMetal } from "../chapters/vol3/Ch35PreciousMetal";
import { Ch36PreCapitalistRelationships } from "../chapters/vol3/Ch36PreCapitalistRelationships";
import { Ch37Introduction } from "../chapters/vol3/Ch37Introduction";
import { Ch38GeneralRemarks } from "../chapters/vol3/Ch38GeneralRemarks";
import { Ch39DifferentialRentI } from "../chapters/vol3/Ch39DifferentialRentI";
import { Ch40DifferentialRentII } from "../chapters/vol3/Ch40DifferentialRentII";
import { Ch41DR2ConstantPrice } from "../chapters/vol3/Ch41DR2ConstantPrice";
import { Ch42DR2FallingPrice } from "../chapters/vol3/Ch42DR2FallingPrice";
import { Ch43DR2RisingPrice } from "../chapters/vol3/Ch43DR2RisingPrice";
import { Ch44WorstSoilRent } from "../chapters/vol3/Ch44WorstSoilRent";
import { Ch45AbsoluteRent } from "../chapters/vol3/Ch45AbsoluteRent";
import { Ch46LandPrice } from "../chapters/vol3/Ch46LandPrice";
import { Ch47Genesis } from "../chapters/vol3/Ch47Genesis";
import { Ch48Trinity } from "../chapters/vol3/Ch48Trinity";
import { Ch49Analysis } from "../chapters/vol3/Ch49Analysis";
import { Ch50Illusions } from "../chapters/vol3/Ch50Illusions";
import { Ch51Distribution } from "../chapters/vol3/Ch51Distribution";
import { Ch52Classes } from "../chapters/vol3/Ch52Classes";

interface ChapterShellProps {
  activeChapterId: string;
  commodities: Commodity[];
  owners: Owner[];
  onSharedChanged: () => void;
}

/** All props that any chapter panel could receive. Panels ignore what they don't need. */
interface PanelProps {
  commodities: Commodity[];
  owners: Owner[];
  onSharedChanged: () => void;
}

// eslint-disable-next-line @typescript-eslint/no-explicit-any
type AnyPanel = ComponentType<any>;

export const CHAPTER_PANELS: Partial<Record<string, AnyPanel>> = {
  "v1-ch01": Ch01Commodity as AnyPanel,
  "v1-ch02": Ch02Exchange as AnyPanel,
  "v1-ch03": Ch03Money as AnyPanel,
  "v1-ch04": Ch04Capital as AnyPanel,
  "v1-ch05": Ch05Contradictions as AnyPanel,
  "v1-ch06": Ch06LabourPower as AnyPanel,
  "v1-ch07": Ch07LabourProcess as AnyPanel,
  "v1-ch08": Ch08ConstantVariableCapital as AnyPanel,
  "v1-ch09": Ch09RateOfSurplusValue as AnyPanel,
  "v1-ch10": Ch10WorkingDay as AnyPanel,
  "v1-ch11": Ch11RateAndMassOfSurplusValue as AnyPanel,
  "v1-ch12": Ch12RelativeSurplusValue as AnyPanel,
  "v1-ch13": Ch13Cooperation as AnyPanel,
  "v1-ch14": Ch14Manufacture as AnyPanel,
  "v1-ch15": Ch15Machinery as AnyPanel,
  "v1-ch16": Ch16AbsoluteAndRelative as AnyPanel,
  "v1-ch17": Ch17MagnitudeChanges as AnyPanel,
  "v1-ch18": Ch18RatesOfSurplusValue as AnyPanel,
  "v1-ch19": Ch19WageForm as AnyPanel,
  "v1-ch20": Ch20TimeWages as AnyPanel,
  "v1-ch21": Ch21PieceWages as AnyPanel,
  "v1-ch22": Ch22NationalWages as AnyPanel,
  "v1-ch23": Ch23SimpleReproduction as AnyPanel,
  "v1-ch24": Ch24AccumulationOfCapital as AnyPanel,
  "v1-ch25": Ch25GeneralLaw as AnyPanel,
  "v1-ch26": Ch26PrimitiveAccumulation as AnyPanel,
  "v1-ch27": Ch27EnclosureEvents as AnyPanel,
  "v1-ch28": Ch28BloodyLegislation as AnyPanel,
  "v1-ch29": Ch29GenesisFarmer as AnyPanel,
  "v1-ch30": Ch30HomeMarket as AnyPanel,
  "v1-ch31": Ch31IndustrialCapitalist as AnyPanel,
  "v1-ch32": Ch32HistoricalTendency as AnyPanel,
  "v1-ch33": Ch33ModernColonisation as AnyPanel,
  "v2-ch01": Ch01CircuitMoneyCapital as AnyPanel,
  "v2-ch02": Ch02CircuitProductiveCapital as AnyPanel,
  "v2-ch03": Ch03CircuitCommodityCapital as AnyPanel,
  "v2-ch04": Ch04ThreeFormulas as AnyPanel,
  "v2-ch05": Ch05TimeOfCirculation as AnyPanel,
  "v2-ch06": Ch06CostsOfCirculation as AnyPanel,
  "v2-ch07": Ch07TurnoverTime as AnyPanel,
  "v2-ch08": Ch08FixedAndCirculatingCapital as AnyPanel,
  "v2-ch09": Ch09AggregateTurnover as AnyPanel,
  "v2-ch10": Ch10TheoriesPhysiocratsSmith as AnyPanel,
  "v2-ch11": Ch11TheoriesRicardo as AnyPanel,
  "v2-ch12": Ch12WorkingPeriod as AnyPanel,
  "v2-ch13": Ch13TimeOfProduction as AnyPanel,
  "v2-ch14": Ch14TimeOfCirculation as AnyPanel,
  "v2-ch15": Ch15PriceChanges as AnyPanel,
  "v2-ch16": Ch16TurnoverOfVariableCapital as AnyPanel,
  "v2-ch17": Ch17CirculationOfSurplus as AnyPanel,
  "v2-ch18": Ch18RoleOfMoneyCapital as AnyPanel,
  "v2-ch19": Ch19FormerPresentations as AnyPanel,
  "v2-ch20": Ch20SimpleReproduction as AnyPanel,
  "v2-ch21": Ch21ExtendedReproduction as AnyPanel,
  "v3-ch01": Ch01CostPriceAndProfit as AnyPanel,
  "v3-ch02": Ch02RateOfProfit as AnyPanel,
  "v3-ch03": Ch03ProfitRateRelation as AnyPanel,
  "v3-ch04": Ch04TurnoverEffect as AnyPanel,
  "v3-ch05": Ch05ConstantCapitalEconomy as AnyPanel,
  "v3-ch06": Ch06PriceFluctuation as AnyPanel,
  "v3-ch07": Ch07SupplementaryRemarks as AnyPanel,
  "v3-ch08": Ch08CompositionDifferences as AnyPanel,
  "v3-ch09": Ch09GeneralRateOfProfit as AnyPanel,
  "v3-ch10": Ch10Equalisation as AnyPanel,
  "v3-ch11": Ch11WageFluctuations as AnyPanel,
  "v3-ch12": Ch12SupplementaryRemarks as AnyPanel,
  "v3-ch13": Ch13FallingProfitLaw as AnyPanel,
  "v3-ch14": Ch14CounteractingForces as AnyPanel,
  "v3-ch15": Ch15InternalContradictions as AnyPanel,
  "v3-ch16": Ch16CommercialCapital as AnyPanel,
  "v3-ch17": Ch17CommercialProfit as AnyPanel,
  "v3-ch18": Ch18MerchantTurnover as AnyPanel,
  "v3-ch19": Ch19MoneyDealingCapital as AnyPanel,
  "v3-ch20": Ch20HistoricalMerchantCapital as AnyPanel,
  "v3-ch21": Ch21InterestBearingCapital as AnyPanel,
  "v3-ch22": Ch22RateOfInterest as AnyPanel,
  "v3-ch23": Ch23ProfitOfEnterprise as AnyPanel,
  "v3-ch24": Ch24CompoundInterest as AnyPanel,
  "v3-ch25": Ch25CreditFictitiousCapital as AnyPanel,
  "v3-ch26": Ch26MoneyCapitalAccumulation as AnyPanel,
  "v3-ch27": Ch27CreditRole as AnyPanel,
  "v3-ch28": Ch28MediumOfCirculation as AnyPanel,
  "v3-ch29": Ch29BankCapital as AnyPanel,
  "v3-ch30": Ch30MoneyRealCapital as AnyPanel,
  "v3-ch31": Ch31MoneyRealCapitalII as AnyPanel,
  "v3-ch32": Ch32MoneyRealCapitalIII as AnyPanel,
  "v3-ch33": Ch33CirculationCredit as AnyPanel,
  "v3-ch34": Ch34CurrencyPrinciple as AnyPanel,
  "v3-ch35": Ch35PreciousMetal as AnyPanel,
  "v3-ch36": Ch36PreCapitalistRelationships as AnyPanel,
  "v3-ch37": Ch37Introduction as AnyPanel,
  "v3-ch38": Ch38GeneralRemarks as AnyPanel,
  "v3-ch39": Ch39DifferentialRentI as AnyPanel,
  "v3-ch40": Ch40DifferentialRentII as AnyPanel,
  "v3-ch41": Ch41DR2ConstantPrice as AnyPanel,
  "v3-ch42": Ch42DR2FallingPrice as AnyPanel,
  "v3-ch43": Ch43DR2RisingPrice as AnyPanel,
  "v3-ch44": Ch44WorstSoilRent as AnyPanel,
  "v3-ch45": Ch45AbsoluteRent as AnyPanel,
  "v3-ch46": Ch46LandPrice as AnyPanel,
  "v3-ch47": Ch47Genesis as AnyPanel,
  "v3-ch48": Ch48Trinity as AnyPanel,
  "v3-ch49": Ch49Analysis as AnyPanel,
  "v3-ch50": Ch50Illusions as AnyPanel,
  "v3-ch51": Ch51Distribution as AnyPanel,
  "v3-ch52": Ch52Classes as AnyPanel,
};

const QUOTES: Partial<Record<string, string>> = {
  "v1-ch01": "The wealth of societies in which the capitalist mode of production prevails appears as an immense collection of commodities.",
  "v1-ch02": "Commodities cannot go to market and make exchanges of their own account.",
  "v1-ch03": "Money is a crystal formed of necessity in the course of exchanges.",
  "v1-ch04": "The circulation of commodities is the starting-point of capital.",
  "v1-ch05": "Circulation, or the exchange of commodities, begets no value.",
  "v1-ch06": "The owner of labour-power is mortal. If his appearance in the market is to be continuous, the seller of labour-power must perpetuate himself.",
  "v1-ch07": "The secret of the self-expansion of capital resolves itself into having the disposal of a definite quantity of other people's unpaid labour.",
  "v1-ch08": "That part of capital which is represented by the means of production ... does not, in the process of production, undergo any quantitative alteration of value.",
  "v1-ch09": "The rate of surplus-value is therefore an exact expression for the degree of exploitation of labour-power by capital.",
  "v1-ch10": "The capitalist maintains his rights as a purchaser when he tries to make the working-day as long as possible.",
  "v1-ch11": "The mass of the surplus-value produced is equal to the amount of the variable capital advanced, multiplied by the rate of surplus-value.",
  "v1-ch12": "The shortening of the working-day is, therefore, by no means what is aimed at, in capitalist production, when labour is economised by increasing its productiveness.",
  "v1-ch13": "When the labourer co-operates systematically with others, he strips off the fetters of his individuality, and develops the capabilities of his species.",
  "v1-ch14": "While simple co-operation leaves the mode of working by the individual ... unchanged, manufacture revolutionises it ... it converts the worker into a crippled monstrosity by furthering his particular skill ... at the expense of a world of productive drives and inclinations.",
  "v1-ch15": "The machine ... enters into the value-begetting process only by bits. It never adds more value than it loses, on an average, by wear and tear.",
  "v1-ch16": "From one standpoint, any distinction between absolute and relative surplus-value appears illusory. Relative surplus-value is absolute, since it compels the absolute prolongation of the working-day beyond the labour-time necessary to the existence of the labourer himself.",
  "v1-ch17": "The value of labour-power, and the surplus-value, vary in opposite directions. The value of labour-power cannot fall, and therefore surplus-value cannot rise, without a rise in the productiveness of labour.",
  "v1-ch18": "These three formulae amount to the same thing, yet they correspond to three essentially different conceptions of the rate of surplus-value.",
  "v1-ch19": "The wage-form thus extinguishes every trace of the division of the working-day into necessary labour and surplus-labour, into paid and unpaid labour.",
  "v1-ch20": "The time-wage is only a converted form of the value, or the price, of labour-power.",
  "v1-ch21": "The piece-wage is the converted form of the time-wage, just as time-wage is the converted form of the value, or price, of labour-power.",
  "v1-ch22": "In England wages are virtually lower to the capitalist, though higher to the operative than on the Continent of Europe.",
  "v1-ch23": "Whatever the form of the process of production in a society, it must be a continuous process, must continue to go periodically through the same phases.",
  "v1-ch24": "Accumulation of capital is, therefore, increase of the proletariat.",
  "v1-ch25": "Accumulation of wealth at one pole is, therefore, at the same time accumulation of misery, agony of toil, slavery, ignorance, brutality, mental degradation, at the opposite pole.",
  "v1-ch26": "The discovery of gold and silver in America, the extirpation, enslavement and entombment in mines of the aboriginal population, the beginning of the conquest and looting of the East Indies, the turning of Africa into a warren for the commercial hunting of black-skins, signalised the rosy dawn of the era of capitalist production.",
  "v1-ch27": "The expropriation of the agricultural producer, of the peasant, from the soil is the basis of the whole process. The history of this expropriation, in different countries, assumes different aspects, and runs through its various phases in different orders of succession, and at different periods.",
  "v1-ch28": "Thus were the agricultural people, first forcibly expropriated from the soil, driven from their homes, turned into vagabonds, and then whipped, branded, tortured by laws grotesquely terrible, into the discipline necessary for the wage system.",
  "v1-ch29": "The agricultural revolution which commenced in the last third of the 15th century, and continued during almost the whole of the 16th, enriched the capitalist farmer as rapidly as it impoverished the mass of the agricultural folk.",
  "v1-ch30": "The expropriation and expulsion of the agricultural population, intermittent but renewed again and again, supplied, as we saw, the town industries with a mass of proletarians entirely unconnected with the corporate guilds and unfettered by them.",
  "v1-ch31": "If money, according to Augier, 'comes into the world with a congenital blood-stain on one cheek,' capital comes dripping from head to foot, from every pore, with blood and dirt.",
  "v1-ch32": "The knell of capitalist private property sounds. The expropriators are expropriated.",
  "v1-ch33": "Capital is not a thing, but a social relation between persons, established by the instrumentality of things.",
  "v2-ch01": "The circuit of money-capital ... M—C(Lp+Mp)…P…C′—M′. The starting point and termination of its motion is the money-form of capital.",
  "v2-ch02": "The circuit of productive capital ... begins with the act of production and ends with the same act.",
  "v2-ch03": "C′—M′—C(Lp+Mp)…P…C′. The starting-point and the finishing-point of the circuit is commodity-capital, or C′, i.e., capital-value plus surplus-value.",
  "v2-ch04": "The circuit of industrial capital … encompasses within its cycle all three circuits … by the side of one another, i.e., simultaneously.",
  "v2-ch05": "The time of circulation of capital limits its time of production and hence the creation of surplus-value.",
  "v2-ch06": "The costs of circulation which we are considering here are not costs out of which surplus-value grows, but costs which limit it, pure deductions from surplus-value.",
  "v2-ch07": "The turnover of capital has, therefore, a definite periodicity … the number of turnovers of a given capital in a given period … is equal to the length of the working period divided into the unit of time.",
  "v2-ch08": "Fixed capital and circulating capital are … not two independent categories of capital, but merely two special forms assumed by the same capital.",
  "v2-ch09": "The aggregate turnover of the advanced productive capital is the mean average of the turnover of its different component parts.",
  "v2-ch10": "Quesnay had already distinguished avances primitives … from avances annuelles … The former … correspond roughly to fixed capital; the latter … to circulating capital.",
  "v2-ch11": "Ricardo confuses fixed and circulating capital with constant and variable capital … Ricardo's mistake arises from the fact that he never distinguished between constant and variable capital.",
  "v2-ch12": "The length of the working period differs in different branches of production. In agriculture … the working period amounts to several months. In manufacture … it varies according to the articles produced.",
  "v2-ch13": "In some branches of industry the product requires a relatively long period of production-time, but a relatively short working-period. The greater the excess of the production-time over the working-period, the greater is the portion of the capital that is continuously tied up in the production process, and not available for other employment.",
  "v2-ch14": "With the development of the means of communication and transport, the time of circulation of commodities, and therefore the time during which capital is tied up in the process of circulation, is shortened.",
  "v2-ch15": "A change in the value (price) of the elements of productive capital, whether of means of production or of labour-power, changes the magnitude of the capital-value to be advanced for the reproduction of productive capital.",
  "v2-ch16": "The annual rate of surplus-value … is determined not only by the rate of surplus-value that the capital extracts from the workers in a single turnover, but by how often this exploitation process repeats itself during the year.",
  "v2-ch17": "Where does the money come from to realise the aggregate surplus-value? … formally resolved in Chs. 20–21.",
  "v2-ch18": "The circulation of money-capital, and of the whole commodity-capital as well, is now for the first time included in the circuit, for here both are within the circulation sphere.",
  "v2-ch19": "Quesnay's Tableau Économique … shows how the result of annual production circulates through consumption and thus reproduces its starting-point. Adam Smith, however, entirely omits the constant portion of value which must replace the constant capital consumed.",
  "v2-ch20": "For simple reproduction to proceed without a hitch, the variable capital plus the surplus-value of Department I must equal the constant capital of Department II: I(v+s) = II(c).",
  "v3-ch01": "The capitalist may sell a commodity at a profit even when he sells it below its value. So long as its selling price is higher than its cost-price, a portion of the surplus-value incorporated in it is always realised.",
  "v3-ch02": "Profit is nevertheless a converted form of surplus-value, a form in which its origin and the secret of its existence are obscured and extinguished.",
  "v3-ch03": "The rate of profit is related to the rate of surplus-value as the variable capital is to the total capital.",
  "v3-ch04": "The shorter the period of turnover ... the larger, therefore, the appropriated surplus-value, provided other conditions remain the same.",
  "v3-ch05": "Every economy in the conditions of production ... increases the rate of profit, since it reduces the value of the constant capital, while leaving the surplus-value untouched.",
  "v3-ch06": "The rate of profit ... may rise or fall in consequence of fluctuations in the prices of raw materials, even when these fluctuations do not affect the rate of surplus-value in the least.",
  "v3-ch07": "The rate of profit ... depends on so many factors that one might say the rate of profit is not directly determinable at all — yet a mere change in the money-value of capital leaves it wholly untouched.",
  "v3-ch08": "Equal capitals in different branches of production ... produce very different rates of profit, corresponding to the different organic compositions of these capitals.",
  "v3-ch09": "The different rates of profit in the various branches of production are, through competition, levelled into a general rate of profit. The price at which the commodities are sold ... is their price of production.",
  "v3-ch10": "Competition can only equalise the deviating profit rates into a general rate by transforming the values of commodities into prices of production. As soon as profit is distributed uniformly, it is distributed as such regardless of its origin.",
  "v3-ch11": "A general rise in wages ... would therefore lower the general rate of profit but would not alter the values of commodities.",
  "v3-ch12": "The price of production of a commodity may vary for only two reasons: first, a change in the general rate of profit; second, a change in the value of the commodity itself. The grounds advanced for a higher price create no value — they are but a claim to a larger share of the aggregate surplus-value.",
  "v3-ch13": "The progressive tendency of the general rate of profit to fall is therefore just an expression peculiar to the capitalist mode of production of the progressive development of the social productivity of labour.",
  "v3-ch14": "The same influences which produce a tendency in the general rate of profit to fall, also call forth counter-effects, which hamper, retard, and partly paralyse this fall.",
  "v3-ch15": "The real barrier of capitalist production is capital itself.",
  "v3-ch16": "Merchant's capital is nothing but capital functioning within the sphere of circulation … whose exclusive activity consists in performing this exchange process.",
  "v3-ch17": "The merchant buys and sells for many industrial capitalists, and thereby also takes over the purchases and sales for the whole society.",
  "v3-ch18": "The faster the merchant turns his capital, the smaller is the addition to the price of each individual article — but the annual profit on the total capital invested remains the same.",
  "v3-ch19": "The money-dealer's capital does not create any new value. It only contributes to equalising the general rate of profit by participating in its formation as a factor.",
  "v3-ch20": "Merchant's capital is older than the capitalist mode of production, is in fact historically the oldest free state of existence of capital.",
  "v3-ch21": "In the formula M—M′ we have the irrational form of capital, the perversion and objectification of production relations in their highest degree.",
  "v3-ch22": "The average rate of interest prevailing in a given country — as distinct from the continually fluctuating market rates — cannot be determined by any law. In this sphere there is no law except that of supply and demand.",
  "v3-ch23": "Interest is the fruit of capital as such, as distinct from profit of enterprise, which appears as the wages of superintendence for the management of capital.",
  "v3-ch26": "The accumulation of money-capital does not mean an actual accumulation of money ... It means merely an accumulation of these claims in the hands of banks.",
  "v3-ch27": "The credit system ... accelerates the material development of the productive forces and the establishment of the world market ... At the same time credit accelerates the violent eruptions of this contradiction — crises.",
  "v3-ch28": "The confusion arises from Tooke's failure to distinguish between money as medium of circulation and money as capital — between the coin function, which serves revenue expenditure, and the capital-transfer function, which moves productive capital between capitalists.",
  "v3-ch29": "The greater part of banker's capital is purely fictitious and consists of claims (bills of exchange), government securities (which represent spent capital), and stocks (drafts on future revenue).",
  "v3-ch31": "The same piece of money can effect quite different transactions. The same money serves, on the one hand, as means of circulation, on the other, as loan capital — and the mass of loanable capital is altogether different from the quantity of the circulation.",
  "v3-ch32": "The accumulation of loan capital simply expresses the fact that a portion of the money realised is transformed into loanable capital. This is very far from being identical with actual accumulation, although it may express it within certain limits — and in times of crisis the demand for loan capital is a demand for means of payment, not for productive capital.",
  "v3-ch33": "Credit reduces the need for money as a medium of circulation. The same sum of money can be used simultaneously for many transactions, owing to the speed with which it circulates — and the clearing house cancels the greater part of mutual claims without any money changing hands at all.",
  "v3-ch34": "The Bank Act of 1844 intensified rather than prevented crises: it had to be suspended in both 1847 and 1857, proving that a mechanical note-issue rule cannot substitute for the actual conditions of capitalist reproduction.",
  "v3-ch35": "Gold as world money is the pivot of the international monetary system: the central-bank reserve defends the rate of exchange, and its movements are determined by the balance of payments, not the balance of trade.",
  "v3-ch36": "Usury is the antediluvian form of capital, older than the capitalist mode of production: subordinated under modern industry it survives only as bank credit, its interest disciplined by the general rate of profit.",
  "v3-ch37": "We shall deal with landed property only in so far as a portion of the surplus-value produced by capital falls to the share of the landowner. — Marx, Capital III, Ch. 37.",
  "v3-ch38": "The surplus profit ... is here a permanent and regular affair, and forms the basis of the differential rent. — Marx, Capital III, Ch. 38.",
  "v3-ch39": "The difference between the price of production yielded by the worst soil and that yielded by any better soil constitutes the differential rent. — Marx, Capital III, Ch. 39.",
  "v3-ch40": "The second form of differential rent... arises from the different distribution of successive capital investments on the same land. — Marx, Capital III, Ch. 40.",
  "v3-ch41": "With a constant price of production... every additional investment of capital in the soil bearing rent increases the rent per acre, although the rate of this increase need not be proportional to the increase of capital. — Marx, Capital III, Ch. 41.",
  "v3-ch42": "Even with a falling price of production, additional capital can be invested with the result that the rental rises, or at least does not fall... the grain-rent rises while the money-rent remains constant. — Marx, Capital III, Ch. 42.",
  "v3-ch43": "The more capital is invested in the land, and the higher the development of agriculture and civilisation in general in a country, the more the rents rise... and the more immense becomes the tribute paid to the great landlords. — Marx/Engels, Capital III, Ch. 43.",
  "v3-ch44": "Differential rent... can also develop on the worst cultivated soil, although this soil by definition yields no differential rent in the first instance. — Marx, Capital III, Ch. 44.",
  "v3-ch45": "The lower composition of agricultural capital... is the reason why agricultural products, when sold at their value, yield a surplus over the average profit — and this surplus forms the substance of absolute rent. — Marx, Capital III, Ch. 45.",
  "v3-ch46": "The price of land is nothing but the capitalised and therefore anticipated rent. ... rent, and not the price of land, is the really decisive factor. — Marx, Capital III, Ch. 46.",
  "v3-ch47": "It is the form of rent which we have here analysed, namely ground-rent... that becomes the general, and so to say the standard form of surplus-value. — Marx, Capital III, Ch. 47.",
};

export function ChapterShell({
  activeChapterId,
  commodities,
  owners,
  onSharedChanged,
}: ChapterShellProps) {
  const chapter = CHAPTERS.find((c) => c.id === activeChapterId);
  if (!chapter) return null;

  const quote = QUOTES[activeChapterId];

  const sharedProps: PanelProps = { commodities, owners, onSharedChanged };
  const Panel = CHAPTER_PANELS[activeChapterId];

  const volumeRoman = chapter.volume === 1 ? "I" : chapter.volume === 2 ? "II" : "III";

  return (
    <main className="chapter-main">
      <header className="chapter-header">
        <div className="chapter-header-num">
          Vol. {volumeRoman} &middot; Chapter {String(chapter.number).padStart(2, "0")}
        </div>
        <h1 className="chapter-header-title">{chapter.title}</h1>
        {quote && <p className="chapter-header-quote">{quote}</p>}
      </header>

      <div className="chapter-body">
        {chapter.status === "pending" ? (
          <div className="chapter-placeholder">
            <p>Not yet implemented</p>
            <p className="small muted">Coming in a future chapter branch</p>
          </div>
        ) : Panel ? (
          <Panel {...sharedProps} />
        ) : null}
      </div>
    </main>
  );
}
