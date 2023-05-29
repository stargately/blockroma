const sm = require('sitemap');
var XmlSitemap = require('xml-sitemap');
var xmlString = require('fs').readFileSync(`${__dirname}/../../build/sitemap.xml`);
var oldSitemap = new XmlSitemap(xmlString);

[
  // external
  ...[
  ].map(it => ({
    url: it,
    changefreq: 'weekly',
    priority: 0.5,
  })),

  // high priority
  ...[
  ].map(it => ({
    url: it,
    changefreq: 'weekly',
    priority: 1.0,
  })),
].forEach(it => {
  if (oldSitemap.hasUrl(it.url)) {
    const {url, ...opts} = it;
    oldSitemap.setOptionValues(it.url, opts);
  } else {
    oldSitemap.add(it)
  }
});

oldSitemap.removeOption('lastmod');

require('fs').writeFileSync(`${__dirname}/../../build/sitemap-manual.xml`, oldSitemap.xml)
require('fs').writeFileSync(`${__dirname}/../../build/sitemap.xml`, oldSitemap.xml)
