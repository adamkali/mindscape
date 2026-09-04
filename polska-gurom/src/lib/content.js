// Central content for the site. Keeping copy here (rather than inline in the
// component) makes it easy to edit the facts without touching layout/markup.

export const phrase = {
  colloquial: 'Polska gurom',
  correct: 'Polska górą',
  gloss: 'Poland on top!',
  ipa: '[ˈpɔlska ˈɡurɔm]',
  respell: 'POL-ska GOO-rom'
};

export const meanings = [
  {
    heading: 'What it means',
    body:
      '"Polska górą!" is an exclamation of national pride and encouragement — roughly "Poland on top!", "Go Poland!" or "Poland rules!". It is a cheer, not a literal statement, shouted when Poland (or anything Polish) is winning or ought to win.'
  },
  {
    heading: 'The "górą" construction',
    body:
      'The word górą is the instrumental case of góra ("mountain, top"). In the fixed pattern "X górą!" it stops meaning a mountain and means "X is on top / X wins!". The pattern is fully productive: "Nasi górą!" ("Our side wins!"), "Wisła górą!" (a football chant), and so on.'
  },
  {
    heading: 'Why "gurom" and not "górą"?',
    body:
      '"gurom" is simply how "górą" sounds written out phonetically. Polish ó is pronounced like English "oo" (so gó- sounds like "goo"), and a word-final -ą is commonly pronounced [ɔm], like "-om". Spell the sound and you get gu-rom. It is a colloquial / eggcorn spelling, not the standard form.'
  }
];

export const pronunciation = [
  { label: 'Correct spelling', value: 'Polska górą' },
  { label: 'Phonetic (eggcorn) spelling', value: 'Polska gurom' },
  { label: 'IPA', value: '[ˈpɔlska ˈɡurɔm]' },
  { label: 'English respelling', value: 'POL-ska GOO-rom' }
];

export const usage = [
  {
    where: 'Sport',
    detail:
      'Football, volleyball, ski jumping, handball — chanted from the stands and in front of the TV whenever the national team plays.'
  },
  {
    where: 'Celebrations',
    detail:
      'National holidays, parades and any moment of collective pride. A short, punchy way to say "we are Polish and proud".'
  },
  {
    where: 'Online & memes',
    detail:
      'A staple caption and comment, often deliberately spelled "gurom" for a playful, informal, very-online tone.'
  }
];

export const examples = [
  { pl: 'Polska górą! Wygraliśmy!', en: 'Poland on top! We won!' },
  { pl: 'Nasi górą, brawo chłopaki!', en: 'Our side wins, well done lads!' },
  { pl: 'Do przodu, Polska górą!', en: 'Come on, go Poland!' }
];
