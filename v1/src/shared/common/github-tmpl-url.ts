function newGithubIssueUrl(options: Record<string, string> = {}) {
  let repoUrl;
  if (options.repoUrl) {
    repoUrl = options.repoUrl;
  } else if (options.user && options.repo) {
    repoUrl = `https://github.com/${options.user}/${options.repo}`;
  } else {
    throw new Error(
      "You need to specify either the `repoUrl` option or both the `user` and `repo` options"
    );
  }

  const url = new URL(`${repoUrl}/issues/new`);

  const types = [
    "body",
    "title",
    "labels",
    "template",
    "milestone",
    "assignee",
    "projects",
  ];

  for (const type of types) {
    let value = options[type];
    if (value === undefined) {
      // eslint-disable-next-line no-continue
      continue;
    }

    if (type === "labels" || type === "projects") {
      if (!Array.isArray(value)) {
        throw new TypeError(`The \`${type}\` option should be an array`);
      }

      value = value.join(",");
    }

    url.searchParams.set(type, value);
  }

  return url.toString();
}

export const getGithubTmplUrl = ({ symbol }: { symbol: string }) => {
  const ua = window?.navigator?.userAgent;
  return newGithubIssueUrl({
    user: "stargately",
    repo: "blockroma",
    title: `${symbol}: <Issue Title>`,
    body: `*Describe your issue here.*

### Environment
* Blockroma Version: ${new Date().toISOString()}

* User Agent: \`${ua}\`

### Steps to reproduce

*Tell us how to reproduce this issue. If possible, push up a branch to your fork with a regression test we can run to reproduce locally.*

### Expected Behaviour

*Tell us what should happen.*

### Actual Behaviour

*Tell us what happens instead.*
  `,
  });
};
