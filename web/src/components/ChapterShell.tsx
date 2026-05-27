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

const CHAPTER_PANELS: Partial<Record<string, AnyPanel>> = {
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
