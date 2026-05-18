import type { ComponentType } from "react";
import { CHAPTERS } from "../chapters/registry";
import type { Commodity, Owner } from "../types";
import { Ch01Commodity } from "../chapters/Ch01Commodity";
import { Ch02Exchange } from "../chapters/Ch02Exchange";
import { Ch03Money } from "../chapters/Ch03Money";
import { Ch04Capital } from "../chapters/Ch04Capital";
import { Ch05Contradictions } from "../chapters/Ch05Contradictions";
import { Ch06LabourPower } from "../chapters/Ch06LabourPower";
import { Ch07LabourProcess } from "../chapters/Ch07LabourProcess";
import { Ch08ConstantVariableCapital } from "../chapters/Ch08ConstantVariableCapital";
import { Ch09RateOfSurplusValue } from "../chapters/Ch09RateOfSurplusValue";
import { Ch10WorkingDay } from "../chapters/Ch10WorkingDay";
import { Ch11RateAndMassOfSurplusValue } from "../chapters/Ch11RateAndMassOfSurplusValue";
import { Ch12RelativeSurplusValue } from "../chapters/Ch12RelativeSurplusValue";
import { Ch13Cooperation } from "../chapters/Ch13Cooperation";
import { Ch14Manufacture } from "../chapters/Ch14Manufacture";
import { Ch15Machinery } from "../chapters/Ch15Machinery";
import { Ch16AbsoluteAndRelative } from "../chapters/Ch16AbsoluteAndRelative";
import { Ch17MagnitudeChanges } from "../chapters/Ch17MagnitudeChanges";
import { Ch18RatesOfSurplusValue } from "../chapters/Ch18RatesOfSurplusValue";
import { Ch19WageForm } from "../chapters/Ch19WageForm";
import { Ch20TimeWages } from "../chapters/Ch20TimeWages";
import { Ch21PieceWages } from "../chapters/Ch21PieceWages";
import { Ch22NationalWages } from "../chapters/Ch22NationalWages";
import { Ch23SimpleReproduction } from "../chapters/Ch23SimpleReproduction";
import { Ch24AccumulationOfCapital } from "../chapters/Ch24AccumulationOfCapital";
import { Ch25GeneralLaw } from "../chapters/Ch25GeneralLaw";
import { Ch26PrimitiveAccumulation } from "../chapters/Ch26PrimitiveAccumulation";
import { Ch27EnclosureEvents } from "../chapters/Ch27EnclosureEvents";
import { Ch28BloodyLegislation } from "../chapters/Ch28BloodyLegislation";
import { Ch29GenesisFarmer } from "../chapters/Ch29GenesisFarmer";
import { Ch30HomeMarket } from "../chapters/Ch30HomeMarket";
import { Ch31IndustrialCapitalist } from "../chapters/Ch31IndustrialCapitalist";
import { Ch32HistoricalTendency } from "../chapters/Ch32HistoricalTendency";
import { Ch33ModernColonisation } from "../chapters/Ch33ModernColonisation";

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
  ch01: Ch01Commodity as AnyPanel,
  ch02: Ch02Exchange as AnyPanel,
  ch03: Ch03Money as AnyPanel,
  ch04: Ch04Capital as AnyPanel,
  ch05: Ch05Contradictions as AnyPanel,
  ch06: Ch06LabourPower as AnyPanel,
  ch07: Ch07LabourProcess as AnyPanel,
  ch08: Ch08ConstantVariableCapital as AnyPanel,
  ch09: Ch09RateOfSurplusValue as AnyPanel,
  ch10: Ch10WorkingDay as AnyPanel,
  ch11: Ch11RateAndMassOfSurplusValue as AnyPanel,
  ch12: Ch12RelativeSurplusValue as AnyPanel,
  ch13: Ch13Cooperation as AnyPanel,
  ch14: Ch14Manufacture as AnyPanel,
  ch15: Ch15Machinery as AnyPanel,
  ch16: Ch16AbsoluteAndRelative as AnyPanel,
  ch17: Ch17MagnitudeChanges as AnyPanel,
  ch18: Ch18RatesOfSurplusValue as AnyPanel,
  ch19: Ch19WageForm as AnyPanel,
  ch20: Ch20TimeWages as AnyPanel,
  ch21: Ch21PieceWages as AnyPanel,
  ch22: Ch22NationalWages as AnyPanel,
  ch23: Ch23SimpleReproduction as AnyPanel,
  ch24: Ch24AccumulationOfCapital as AnyPanel,
  ch25: Ch25GeneralLaw as AnyPanel,
  ch26: Ch26PrimitiveAccumulation as AnyPanel,
  ch27: Ch27EnclosureEvents as AnyPanel,
  ch28: Ch28BloodyLegislation as AnyPanel,
  ch29: Ch29GenesisFarmer as AnyPanel,
  ch30: Ch30HomeMarket as AnyPanel,
  ch31: Ch31IndustrialCapitalist as AnyPanel,
  ch32: Ch32HistoricalTendency as AnyPanel,
  ch33: Ch33ModernColonisation as AnyPanel,
};

const QUOTES: Partial<Record<string, string>> = {
  ch01: "The wealth of societies in which the capitalist mode of production prevails appears as an immense collection of commodities.",
  ch02: "Commodities cannot go to market and make exchanges of their own account.",
  ch03: "Money is a crystal formed of necessity in the course of exchanges.",
  ch04: "The circulation of commodities is the starting-point of capital.",
  ch05: "Circulation, or the exchange of commodities, begets no value.",
  ch06: "The owner of labour-power is mortal. If his appearance in the market is to be continuous, the seller of labour-power must perpetuate himself.",
  ch07: "The secret of the self-expansion of capital resolves itself into having the disposal of a definite quantity of other people's unpaid labour.",
  ch08: "That part of capital which is represented by the means of production ... does not, in the process of production, undergo any quantitative alteration of value.",
  ch09: "The rate of surplus-value is therefore an exact expression for the degree of exploitation of labour-power by capital.",
  ch10: "The capitalist maintains his rights as a purchaser when he tries to make the working-day as long as possible.",
  ch11: "The mass of the surplus-value produced is equal to the amount of the variable capital advanced, multiplied by the rate of surplus-value.",
  ch12: "The shortening of the working-day is, therefore, by no means what is aimed at, in capitalist production, when labour is economised by increasing its productiveness.",
  ch13: "When the labourer co-operates systematically with others, he strips off the fetters of his individuality, and develops the capabilities of his species.",
  ch14: "While simple co-operation leaves the mode of working by the individual ... unchanged, manufacture revolutionises it ... it converts the worker into a crippled monstrosity by furthering his particular skill ... at the expense of a world of productive drives and inclinations.",
  ch15: "The machine ... enters into the value-begetting process only by bits. It never adds more value than it loses, on an average, by wear and tear.",
  ch16: "From one standpoint, any distinction between absolute and relative surplus-value appears illusory. Relative surplus-value is absolute, since it compels the absolute prolongation of the working-day beyond the labour-time necessary to the existence of the labourer himself.",
  ch17: "The value of labour-power, and the surplus-value, vary in opposite directions. The value of labour-power cannot fall, and therefore surplus-value cannot rise, without a rise in the productiveness of labour.",
  ch18: "These three formulae amount to the same thing, yet they correspond to three essentially different conceptions of the rate of surplus-value.",
  ch19: "The wage-form thus extinguishes every trace of the division of the working-day into necessary labour and surplus-labour, into paid and unpaid labour.",
  ch20: "The time-wage is only a converted form of the value, or the price, of labour-power.",
  ch21: "The piece-wage is the converted form of the time-wage, just as time-wage is the converted form of the value, or price, of labour-power.",
  ch22: "In England wages are virtually lower to the capitalist, though higher to the operative than on the Continent of Europe.",
  ch23: "Whatever the form of the process of production in a society, it must be a continuous process, must continue to go periodically through the same phases.",
  ch24: "Accumulation of capital is, therefore, increase of the proletariat.",
  ch25: "Accumulation of wealth at one pole is, therefore, at the same time accumulation of misery, agony of toil, slavery, ignorance, brutality, mental degradation, at the opposite pole.",
  ch26: "Capital is not a thing, but a social relation between persons, established by the instrumentality of things.",
  ch27: "The spoliation of the Church's property, the fraudulent alienation of the State domains, the robbery of the common lands, the usurpation of feudal and clan property, and its transformation into modern private property under circumstances of reckless terrorism, were just so many idyllic methods of primitive accumulation.",
  ch28: "The advance of capitalist production develops a working class which by education, tradition, and habit looks upon the requirements of that mode of production as self-evident natural laws.",
  ch29: "The progressive fall in the value of the precious metals ... lowered wages. A portion of the latter was now added to profits.",
  ch30: "The expropriation of the agricultural producer, of the peasant, from the soil, is the basis of the whole process. The history of this expropriation, in different countries, assumes different aspects, and runs through its various phases in different orders of succession.",
  ch31: "The discovery of gold and silver in America, the extirpation, enslavement and entombment in mines of the aboriginal population, the beginning of the conquest and looting of the East Indies, the turning of Africa into a warren for the commercial hunting of black-skins, signalised the rosy dawn of the era of capitalist production.",
  ch32: "The knell of capitalist private property sounds. The expropriators are expropriated.",
  ch33: "Capital is not a thing, but a social relation between persons, established by the instrumentality of things.",
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

  return (
    <main className="chapter-main">
      <header className="chapter-header">
        <div className="chapter-header-num">
          Vol. I &middot; Chapter {String(chapter.number).padStart(2, "0")}
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
