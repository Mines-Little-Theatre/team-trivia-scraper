import { findAll, findOne, textContent } from "domutils";
import { parseDocument } from "htmlparser2";

export interface AnswerData {
  title: string;
  blurb: string;
  date: string;
  answer: string;
}

const freeAnswerURL = "https://teamtrivia.com/free/";

export async function fetchFreeAnswer(regionId: string): Promise<AnswerData> {
  const document = parseDocument(
    await (
      await fetch(freeAnswerURL, {
        headers: { Cookie: `new_site=Y; region_ID=${regionId}` },
      })
    ).text(),
  );

  const main = findOne((e) => e.name === "main", document.children);
  if (main) {
    const [titleEl, blurbEl, dateEl, answerEl] = findAll(
      (e) => e.name === "h1" || e.name === "p" || e.name === "h3",
      main.children,
    );
    if (titleEl && blurbEl && dateEl && answerEl) {
      return {
        title: adjustedText(titleEl),
        blurb: adjustedText(blurbEl),
        date: adjustedText(dateEl),
        answer: adjustedText(answerEl),
      };
    }
  }

  throw new Error("failed to scrape free answer");
}

function adjustedText(...args: Parameters<typeof textContent>): string {
  const text = textContent(...args);
  return text.trim().replaceAll(/\s+/g, " ");
}
