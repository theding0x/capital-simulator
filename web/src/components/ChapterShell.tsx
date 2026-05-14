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

interface ChapterShellProps {
  activeChapterId: string;
  commodities: Commodity[];
  owners: Owner[];
  onSharedChanged: () => void;
}

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
        ) : activeChapterId === "ch01" ? (
          <Ch01Commodity commodities={commodities} onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch02" ? (
          <Ch02Exchange
            commodities={commodities}
            owners={owners}
            onSharedChanged={onSharedChanged}
          />
        ) : activeChapterId === "ch03" ? (
          <Ch03Money owners={owners} onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch04" ? (
          <Ch04Capital onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch05" ? (
          <Ch05Contradictions onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch06" ? (
          <Ch06LabourPower onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch07" ? (
          <Ch07LabourProcess onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch08" ? (
          <Ch08ConstantVariableCapital onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch09" ? (
          <Ch09RateOfSurplusValue onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch10" ? (
          <Ch10WorkingDay onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch11" ? (
          <Ch11RateAndMassOfSurplusValue onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch12" ? (
          <Ch12RelativeSurplusValue onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch13" ? (
          <Ch13Cooperation onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch14" ? (
          <Ch14Manufacture onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch15" ? (
          <Ch15Machinery onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch16" ? (
          <Ch16AbsoluteAndRelative onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch17" ? (
          <Ch17MagnitudeChanges onSharedChanged={onSharedChanged} />
        ) : activeChapterId === "ch18" ? (
          <Ch18RatesOfSurplusValue onSharedChanged={onSharedChanged} />
        ) : null}
      </div>
    </main>
  );
}
